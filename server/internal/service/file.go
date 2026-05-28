/*
项目名称：JeriBlog
文件名称：file.go
创建时间：2026-04-16 15:00:03

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：文件管理业务逻辑
*/

package service

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"jeri_blog/config"
	"jeri_blog/internal/dto"
	"jeri_blog/internal/model"
	"jeri_blog/internal/repository"
	"jeri_blog/pkg/logger"
	"jeri_blog/pkg/upload"
	"jeri_blog/pkg/utils"

	"gorm.io/gorm"
)

var reconciledSettingImageKeys = []string{
	KeyBasicAuthorAvatar,
	KeyBasicAuthorPhoto,
	KeyBlogFavicon,
	KeyBlogBackgroundImage,
	KeyBlogAboutExhibition,
	KeyBlogScreenshot,
	KeyBlogWechatQrCode,
	KeyBlogAlipayQrCode,
	KeyBlogWechatOffAccounts,
}

// FileUsageChecker 文件引用检查器
type FileUsageChecker struct {
	articleRepo  *repository.ArticleRepository
	friendRepo   *repository.FriendRepository
	momentRepo   *repository.MomentRepository
	settingRepo  *repository.SettingRepository
	userRepo     *repository.UserRepository
	menuRepo     *repository.MenuRepository
	feedbackRepo *repository.FeedbackRepository
	commentRepo  *repository.CommentRepository
}

// NewFileUsageChecker 创建文件引用检查器
func NewFileUsageChecker(
	articleRepo *repository.ArticleRepository,
	friendRepo *repository.FriendRepository,
	momentRepo *repository.MomentRepository,
	settingRepo *repository.SettingRepository,
	userRepo *repository.UserRepository,
	menuRepo *repository.MenuRepository,
	feedbackRepo *repository.FeedbackRepository,
	commentRepo *repository.CommentRepository,
) *FileUsageChecker {
	return &FileUsageChecker{
		articleRepo:  articleRepo,
		friendRepo:   friendRepo,
		momentRepo:   momentRepo,
		settingRepo:  settingRepo,
		userRepo:     userRepo,
		menuRepo:     menuRepo,
		feedbackRepo: feedbackRepo,
		commentRepo:  commentRepo,
	}
}

// IsActuallyUsed 检查文件是否仍被业务引用
func (c *FileUsageChecker) IsActuallyUsed(fileURL string) (bool, string, error) {
	checks := []struct {
		name string
		fn   func(string) (bool, error)
	}{
		{name: "文章封面", fn: c.articleRepo.ExistsByCover},
		{name: "文章正文", fn: c.articleRepo.ExistsByContentURL},
		{name: "友链图片", fn: c.friendRepo.ExistsByAvatarOrScreenshot},
		{name: "动态内容", fn: c.momentRepo.ExistsByContentURL},
		{name: "设置图片", fn: func(url string) (bool, error) {
			return c.settingRepo.ExistsByValueAndKeys(url, reconciledSettingImageKeys)
		}},
		{name: "用户头像", fn: c.userRepo.ExistsByAvatar},
		{name: "菜单图标", fn: c.menuRepo.ExistsByIcon},
		{name: "反馈附件", fn: c.feedbackRepo.ExistsByAttachmentURL},
		{name: "评论内容", fn: c.commentRepo.ExistsByContentURL},
	}

	for _, check := range checks {
		used, err := check.fn(fileURL)
		if err != nil {
			return false, "", err
		}
		if used {
			return true, check.name, nil
		}
	}

	return false, "", nil
}

// FileService 文件服务
type FileService struct {
	fileRepo      *repository.FileRepository
	uploadManager *upload.Manager
	usageChecker  *FileUsageChecker
	config        *config.Config
}

// NewFileService 创建文件服务
func NewFileService(fileRepo *repository.FileRepository, uploadManager *upload.Manager, cfg *config.Config) *FileService {
	return &FileService{
		fileRepo:      fileRepo,
		uploadManager: uploadManager,
		config:        cfg,
	}
}

// SetUsageChecker 设置文件引用检查器
func (s *FileService) SetUsageChecker(checker *FileUsageChecker) {
	s.usageChecker = checker
}

// ============ 通用服务 ============

// UploadFromReader 从Reader上传文件
func (s *FileService) UploadFromReader(reader io.Reader, originalName, fileType string, uploadType upload.Type, userID uint, host string) (string, error) {
	// 读取文件数据并计算hash
	data, fileHash, err := s.uploadManager.HandleUploadFromReader(reader)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 保存文件数据（使用相对路径）
	fileInfo, err := s.uploadManager.SaveFileData(data, fileHash, originalName, fileType, uploadType, userID, host)
	if err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 创建数据库记录
	file := s.createFileFromUploadInfo(fileInfo)
	file.Status = 0 // 默认未使用

	if err := s.fileRepo.Create(file); err != nil {
		_ = s.uploadManager.DeleteFile(fileInfo.FilePath)
		return "", fmt.Errorf("保存记录失败: %w", err)
	}

	return file.FileURL, nil
}

// MarkAsUsed 增加文件引用计数并添加用途
// normalizeFileURL 规范化文件URL
// 如果是相对路径，尝试查找对应的完整URL
func (s *FileService) normalizeFileURL(fileUrl string) (string, error) {
	if fileUrl == "" {
		return "", nil
	}

	// 如果已经是完整URL（http/https开头），直接返回
	if strings.HasPrefix(fileUrl, "http://") || strings.HasPrefix(fileUrl, "https://") {
		return fileUrl, nil
	}

	// 如果是相对路径，通过模糊匹配查找完整URL
	file, err := s.fileRepo.GetByURLPattern(fileUrl)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 找不到记录，返回原始URL
			return "", nil
		}
		return "", err
	}

	// 找到了完整URL，返回
	logger.Info("URL规范化: %s -> %s", fileUrl, file.FileURL)
	return file.FileURL, nil
}

func (s *FileService) MarkAsUsed(fileUrl string, usageType string) error {
	if fileUrl == "" {
		return nil
	}

	// 规范化URL
	normalizedURL, err := s.normalizeFileURL(fileUrl)
	if err != nil {
		return err
	}

	if normalizedURL == "" {
		return nil
	}

	count, err := s.fileRepo.CountByURL(normalizedURL)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	// 增加引用计数
	if err := s.fileRepo.IncrementReferenceCount(normalizedURL); err != nil {
		return err
	}

	// 添加用途到 upload_type 字段
	if usageType != "" {
		if err := s.fileRepo.AddUploadType(normalizedURL, usageType); err != nil {
			return err
		}
	}

	// 根据引用计数更新状态
	return s.fileRepo.UpdateStatusByReferenceCount(normalizedURL)
}

// MarkAsUnused 减少文件引用计数并移除用途
func (s *FileService) MarkAsUnused(fileUrl string, usageType string) error {
	if fileUrl == "" {
		return nil
	}

	// 规范化URL
	normalizedURL, err := s.normalizeFileURL(fileUrl)
	if err != nil {
		return err
	}

	if normalizedURL == "" {
		return nil
	}

	count, err := s.fileRepo.CountByURL(normalizedURL)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	// 减少引用计数
	if err := s.fileRepo.DecrementReferenceCount(normalizedURL); err != nil {
		return err
	}

	// 从 upload_type 字段移除用途
	if usageType != "" {
		if err := s.fileRepo.RemoveUploadType(normalizedURL, usageType); err != nil {
			return err
		}
	}

	// 根据引用计数更新状态
	return s.fileRepo.UpdateStatusByReferenceCount(normalizedURL)
}

// ============ 前台服务 ============

// UploadForWeb 前台文件上传
func (s *FileService) UploadForWeb(req *upload.Request, host string) (*dto.FileUploadForWebResponse, error) {
	// 文件大小限制（从配置获取，单位MB）
	maxFileSizeMB := s.uploadManager.GetMaxFileSize()
	if maxFileSizeMB <= 0 {
		maxFileSizeMB = 5 // 默认 5MB
	}
	maxWebFileSize := maxFileSizeMB * 1024 * 1024
	if req.File.Size > maxWebFileSize {
		return nil, fmt.Errorf("文件大小超出限制，前台上传最大允许 %dMB", maxFileSizeMB)
	}

	// 文件类型白名单验证（具体场景限制由前端控制）
	contentType := req.File.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg":         true,
		"image/jpg":          true,
		"image/png":          true,
		"image/gif":          true,
		"image/webp":         true,
		"image/avif":         true,
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("不支持的文件类型: %s", contentType)
	}

	// 调用通用上传方法（传递 host）
	file, err := s.handleUpload(req, host)
	if err != nil {
		return nil, err
	}

	// 返回简化响应
	return &dto.FileUploadForWebResponse{
		OriginalName: file.OriginalName,
		FileURL:      file.FileURL,
	}, nil
}

// ============ 后台管理服务 ============

// Upload 文件上传
func (s *FileService) Upload(req *upload.Request, host string) (*dto.FileResponse, error) {
	// 调用通用上传方法（传递 host）
	file, err := s.handleUpload(req, host)
	if err != nil {
		return nil, err
	}

	return &dto.FileResponse{
		ID:             file.ID,
		OriginalName:   file.OriginalName,
		FileName:       file.FileName,
		FileSize:       file.FileSize,
		FileType:       file.FileType,
		FileURL:        file.FileURL,
		UploadType:     upload.Type(file.UploadType),
		UserID:         file.UserID,
		Status:         file.Status,
		ReferenceCount: file.ReferenceCount,
		UploadTime:     utils.NewJSONTime(file.CreatedAt),
	}, nil
}

// List 获取文件列表
func (s *FileService) List(req *dto.ListFilesRequest) ([]dto.FileResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := &repository.FileListFilter{
		Keyword:    req.Keyword,
		FileType:   req.FileType,
		Status:     req.Status,
		UploadType: req.UploadType,
		MinSize:    req.MinSize,
		MaxSize:    req.MaxSize,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}

	// 调用仓储层查询（支持关键词、状态、文件类型、上传类型筛选）
	files, total, err := s.fileRepo.GetByFilter(filter, offset, req.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取文件列表失败: %w", err)
	}

	// 转换为响应格式
	fileResponses := make([]dto.FileResponse, len(files))
	for i, file := range files {
		fileResponses[i] = dto.FileResponse{
			ID:             file.ID,
			OriginalName:   file.OriginalName,
			FileName:       file.FileName,
			FileSize:       file.FileSize,
			FileType:       file.FileType,
			FileURL:        file.FileURL,
			UploadType:     upload.Type(file.UploadType),
			UserID:         file.UserID,
			Status:         file.Status,
			ReferenceCount: file.ReferenceCount,
			UploadTime:     utils.NewJSONTime(file.CreatedAt),
		}
	}

	return fileResponses, total, nil
}

// Get 获取文件详情
func (s *FileService) Get(id uint) (*dto.FileResponse, error) {
	file, err := s.fileRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %w", err)
	}

	return &dto.FileResponse{
		ID:             file.ID,
		OriginalName:   file.OriginalName,
		FileName:       file.FileName,
		FileSize:       file.FileSize,
		FileType:       file.FileType,
		FileURL:        file.FileURL,
		UploadType:     upload.Type(file.UploadType),
		UserID:         file.UserID,
		Status:         file.Status,
		ReferenceCount: file.ReferenceCount,
		UploadTime:     utils.NewJSONTime(file.CreatedAt),
	}, nil
}

// GetReferences 获取文件引用详情
func (s *FileService) GetReferences(id uint) ([]dto.FileReferenceResponse, error) {
	file, err := s.fileRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %w", err)
	}

	// 获取博客前台地址
	BlogURL := strings.TrimSuffix(s.config.Basic.BlogURL, "/")

	if s.usageChecker == nil {
		return []dto.FileReferenceResponse{}, nil
	}

	references := make([]dto.FileReferenceResponse, 0)

	// 检查文章封面
	articles, err := s.usageChecker.articleRepo.FindByCover(file.FileURL)
	if err == nil {
		for _, article := range articles {
			references = append(references, dto.FileReferenceResponse{
				Type:  "article",
				ID:    article.ID,
				Title: article.Title,
				Field: "文章封面",
				URL:   "",
			})
		}
	}

	// 检查文章正文（带文件类型）
	articleRefs, err := s.usageChecker.articleRepo.FindByContentURLWithType(file.FileURL)
	if err == nil {
		for _, ref := range articleRefs {
			references = append(references, dto.FileReferenceResponse{
				Type:  "article",
				ID:    ref.Article.ID,
				Title: ref.Article.Title,
				Field: string(ref.FileType), // 使用文件类型：文章配图/文章视频/文章音频/文章附件
				URL:   "", // 前端通过 Title 搜索，不需要 URL
			})
		}
	}

	// 检查友链头像
	friends, err := s.usageChecker.friendRepo.FindByAvatar(file.FileURL)
	if err == nil {
		for _, friend := range friends {
			references = append(references, dto.FileReferenceResponse{
				Type:  "friend",
				ID:    friend.ID,
				Title: friend.Name,
				Field: "友情链接A",
				URL:   "", // 前端通过 Title 搜索，不需要 URL
			})
		}
	}

	// 检查友链截图
	friends, err = s.usageChecker.friendRepo.FindByScreenshot(file.FileURL)
	if err == nil {
		for _, friend := range friends {
			references = append(references, dto.FileReferenceResponse{
				Type:  "friend",
				ID:    friend.ID,
				Title: friend.Name,
				Field: "友情链接S",
				URL:   "", // 前端通过 Title 搜索，不需要 URL
			})
		}
	}

	// 检查用户头像
	users, err := s.usageChecker.userRepo.FindByAvatar(file.FileURL)
	if err == nil {
		for _, user := range users {
			references = append(references, dto.FileReferenceResponse{
				Type:  "user",
				ID:    user.ID,
				Title: user.Nickname,
				Field: "用户头像",
				URL:   "", // 前端通过 Title 搜索，不需要 URL
			})
		}
	}

	// 检查系统设置
	settings, err := s.usageChecker.settingRepo.FindByValueAndKeys(file.FileURL, reconciledSettingImageKeys)
	if err == nil {
		for _, setting := range settings {
			fieldName := getSettingFieldName(setting.Key)
			references = append(references, dto.FileReferenceResponse{
				Type:  "setting",
				ID:    setting.ID,
				Title: "系统设置",
				Field: fieldName,
				URL:   "", // 前端直接指定 URL，不需要跳转链接
			})
		}
	}

	// 检查动态配图、动态视频和动态音频
	momentRefs, err := s.usageChecker.momentRepo.FindByContentURLWithType(file.FileURL)
	if err == nil {
		for _, ref := range momentRefs {
			// 解析动态内容 JSON，提取 text 字段作为标题
			var contentData map[string]interface{}
			title := "动态" // 默认标题
			if err := json.Unmarshal([]byte(ref.Moment.Content), &contentData); err == nil {
				if text, ok := contentData["text"].(string); ok && text != "" {
					title = text
				}
			}

			references = append(references, dto.FileReferenceResponse{
				Type:  "moment",
				ID:    ref.Moment.ID,
				Title: title,
				Field: ref.FileType, // 直接使用返回的文件类型
				URL:   "", // 前端通过 Title 搜索，不需要 URL
			})
		}
	}

	// 检查评论配图
	comments, err := s.usageChecker.commentRepo.FindByContentURL(file.FileURL)
	if err == nil {
		for _, comment := range comments {
			// 显示评论内容作为标题
			title := comment.Content
			var url string

			// 根据评论目标类型确定跳转链接
			switch comment.TargetType {
			case "article":
				// 优先通过 TargetID 查询文章
				if comment.TargetID != nil && *comment.TargetID > 0 {
					article, err := s.usageChecker.articleRepo.Get(*comment.TargetID)
					if err == nil {
						url = fmt.Sprintf("%s/posts/%s/", BlogURL, article.Slug)
					}
				}
				// 兼容旧数据：通过 TargetKey（slug）查询
				if url == "" {
					article, err := s.usageChecker.articleRepo.GetBySlug(comment.TargetKey)
					if err == nil {
						url = fmt.Sprintf("%s/posts/%s/", BlogURL, article.Slug)
					} else {
						url = "" // 文章已删除
					}
				}
			case "page":
				// 根据 TargetKey 确定具体页面
				switch comment.TargetKey {
				case "moment":
					url = fmt.Sprintf("%s/moment", BlogURL)
				case "message":
					url = fmt.Sprintf("%s/message", BlogURL)
				case "friend":
					url = fmt.Sprintf("%s/friend", BlogURL)
				default:
					url = ""
				}
			}

			references = append(references, dto.FileReferenceResponse{
				Type:       "comment",
				ID:         comment.ID,
				Title:      title,
				Field:      "评论贴图",
				URL:        url,
				TargetType: comment.TargetType,
				TargetKey:  comment.TargetKey,
			})
		}
	}

	// 检查菜单图标
	menus, err := s.usageChecker.menuRepo.FindByIcon(file.FileURL)
	if err == nil {
		for _, menu := range menus {
			references = append(references, dto.FileReferenceResponse{
				Type:  "menu",
				ID:    menu.ID,
				Title: menu.Title,
				Field: "菜单图标",
				URL:   "", // 前端直接指定 URL，不需要跳转链接
			})
		}
	}

	// 检查反馈附件
	feedbacks, err := s.usageChecker.feedbackRepo.FindByAttachmentURL(file.FileURL)
	if err == nil {
		for _, feedback := range feedbacks {
			references = append(references, dto.FileReferenceResponse{
				Type:  "feedback",
				ID:    feedback.ID,
				Title: feedback.TicketNo,
				Field: "反馈投诉",
				URL:   "", // 前端通过 Title（工单号）搜索，不需要 URL
			})
		}
	}

	return references, nil
}

// getSettingFieldName 获取设置字段的友好名称
func getSettingFieldName(key string) string {
	fieldNames := map[string]string{
		KeyBasicAuthorAvatar:   "站长头像",
		KeyBasicAuthorPhoto:    "站长形象",
		KeyBlogFavicon:         "博客图标",
		KeyBlogBackgroundImage: "博客背景",
		KeyBlogAboutExhibition: "展览图片",
		KeyBlogScreenshot:      "博客截图",
		KeyBlogWechatQrCode:    "微信收款码",
		KeyBlogAlipayQrCode:    "支付宝收款码",
		KeyBlogWechatOffAccounts: "公众号二维码",
	}
	if name, ok := fieldNames[key]; ok {
		return name
	}
	return key
}

// Delete 删除文件
func (s *FileService) Delete(id uint) error {
	file, err := s.fileRepo.Get(id)
	if err != nil {
		return fmt.Errorf("文件不存在: %w", err)
	}

	// 检查文件是否被引用
	if s.usageChecker != nil {
		used, source, err := s.usageChecker.IsActuallyUsed(file.FileURL)
		if err != nil {
			return fmt.Errorf("检查文件引用失败: %w", err)
		}
		if used {
			return fmt.Errorf("文件正在被使用，无法删除 (引用来源: %s)", source)
		}
	}

	// 检查是否有其他文件记录使用相同的 URL
	otherFilesExist, err := s.fileRepo.ExistsByURLExcludingID(file.FileURL, id)
	if err != nil {
		return fmt.Errorf("检查文件记录失败: %w", err)
	}

	// 只有当没有其他文件记录使用相同 URL 时，才删除物理文件
	if !otherFilesExist {
		if err := s.uploadManager.DeleteFileByStorageType(file.FilePath, file.StorageType); err != nil {
			// 对于默认头像等系统生成的文件，如果物理文件不存在，只记录警告，不阻止删除数据库记录
			return fmt.Errorf("删除存储文件失败: %w", err)
		}
	}

	// 删除数据库记录
	if err := s.fileRepo.Delete(id); err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	return nil
}

// ============ 辅助方法 ============

// handleUpload 处理文件上传
func (s *FileService) handleUpload(req *upload.Request, host string) (*model.File, error) {
	result, err := s.uploadManager.HandleUpload(req, host)
	if err != nil {
		return nil, fmt.Errorf("文件上传失败: %w", err)
	}

	if !result.Success || result.FileInfo == nil {
		return nil, fmt.Errorf("文件上传失败: %s", result.Message)
	}

	// 创建文件记录
	file := s.createFileFromUploadInfo(result.FileInfo)
	file.Status = 0 // 默认未使用

	if err := s.fileRepo.Create(file); err != nil {
		_ = s.uploadManager.DeleteFile(result.FileInfo.FilePath)
		return nil, fmt.Errorf("保存记录失败: %w", err)
	}

	return file, nil
}

// createFileFromUploadInfo 从上传信息创建文件模型
func (s *FileService) createFileFromUploadInfo(info *upload.FileInfo) *model.File {
	// 处理 UserID：0 表示匿名上传，转为 nil
	var userID *uint
	if info.UserID > 0 {
		userID = &info.UserID
	}

	return &model.File{
		FileName:     info.FileName,
		OriginalName: info.OriginalName,
		FilePath:     info.FilePath,
		FileSize:     info.FileSize,
		FileType:     info.FileType,
		UploadType:   string(info.UploadType),
		StorageType:  info.StorageType,
		UserID:       userID,
		FileURL:      info.FileURL,
	}
}

// ============ 定时任务方法 ============

// DeleteUnusedFiles 删除未使用文件，先纠正误标，再清理超过15天仍未使用的文件
func (s *FileService) DeleteUnusedFiles() error {
	if s.usageChecker == nil {
		return nil
	}

	files, err := s.fileRepo.GetByStatus(0)
	if err != nil {
		return fmt.Errorf("获取未使用文件失败: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	deleteBefore := time.Now().AddDate(0, 0, -15)
	usedURLs := make([]string, 0)
	deletableFiles := make([]model.File, 0)

	for _, file := range files {
		used, source, err := s.usageChecker.IsActuallyUsed(file.FileURL)
		if err != nil {
			logger.Warn("检查文件引用失败 %s: %v", file.FileURL, err)
			continue
		}
		if used {
			usedURLs = append(usedURLs, file.FileURL)
			logger.Info("文件引用自检纠正成功 %s -> %s", file.FileURL, source)
			continue
		}

		if file.CreatedAt.Before(deleteBefore) {
			deletableFiles = append(deletableFiles, file)
		}
	}

	if err := s.fileRepo.UpdateFileStatusByUrls(usedURLs, 1); err != nil {
		return fmt.Errorf("批量纠正文件状态失败: %w", err)
	}

	deletedIDs := make([]uint, 0, len(deletableFiles))
	for _, file := range deletableFiles {
		if err := s.uploadManager.DeleteFileByStorageType(file.FilePath, file.StorageType); err != nil {
			logger.Warn("删除物理文件失败 %s: %v", file.FilePath, err)
			continue
		}
		deletedIDs = append(deletedIDs, file.ID)
	}

	if err := s.fileRepo.DeleteByIDs(deletedIDs); err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	return nil
}
