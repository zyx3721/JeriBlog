package service

import (
	"fmt"
	"io"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/logger"
	"flec_blog/pkg/upload"
	"flec_blog/pkg/utils"
)

// FileService 文件服务
type FileService struct {
	fileRepo      *repository.FileRepository
	uploadManager *upload.Manager
}

// NewFileService 创建文件服务
func NewFileService(fileRepo *repository.FileRepository, uploadManager *upload.Manager) *FileService {
	return &FileService{
		fileRepo:      fileRepo,
		uploadManager: uploadManager,
	}
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

// MarkAsUsed 标记文件为使用中
func (s *FileService) MarkAsUsed(fileUrl string) error {
	if fileUrl == "" {
		return nil
	}
	return s.fileRepo.UpdateStatus(fileUrl, 1)
}

// MarkAsUnused 标记文件为未使用
func (s *FileService) MarkAsUnused(fileUrl string) error {
	if fileUrl == "" {
		return nil
	}
	return s.fileRepo.UpdateStatus(fileUrl, 0)
}

// ============ 前台服务 ============

// UploadForWeb 前台文件上传
func (s *FileService) UploadForWeb(req *upload.Request, host string) (*dto.FileUploadForWebResponse, error) {
	// 验证上传类型
	if string(req.UploadType) == "" {
		return nil, fmt.Errorf("上传类型不能为空")
	}

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
	// 验证上传类型
	if string(req.UploadType) == "" {
		return nil, fmt.Errorf("上传类型不能为空")
	}

	// 调用通用上传方法（传递 host）
	file, err := s.handleUpload(req, host)
	if err != nil {
		return nil, err
	}

	return &dto.FileResponse{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		FileName:     file.FileName,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		FileURL:      file.FileURL,
		UploadType:   upload.Type(file.UploadType),
		UserID:       file.UserID,
		Status:       file.Status,
		UploadTime:   utils.NewJSONTime(file.CreatedAt),
	}, nil
}

// List 获取文件列表
func (s *FileService) List(req *dto.ListFilesRequest) ([]dto.FileResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	var files []model.File
	var total int64
	var err error

	// 根据类型过滤
	if req.Type != "" {
		files, total, err = s.fileRepo.GetByUploadType(req.Type, offset, req.PageSize)
	} else {
		files, total, err = s.fileRepo.List(offset, req.PageSize)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("获取文件列表失败: %w", err)
	}

	// 转换为响应格式
	fileResponses := make([]dto.FileResponse, len(files))
	for i, file := range files {
		fileResponses[i] = dto.FileResponse{
			ID:           file.ID,
			OriginalName: file.OriginalName,
			FileName:     file.FileName,
			FileSize:     file.FileSize,
			FileType:     file.FileType,
			FileURL:      file.FileURL,
			UploadType:   upload.Type(file.UploadType),
			UserID:       file.UserID,
			Status:       file.Status,
			UploadTime:   utils.NewJSONTime(file.CreatedAt),
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
		ID:           file.ID,
		OriginalName: file.OriginalName,
		FileName:     file.FileName,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		FileURL:      file.FileURL,
		UploadType:   upload.Type(file.UploadType),
		UserID:       file.UserID,
		Status:       file.Status,
		UploadTime:   utils.NewJSONTime(file.CreatedAt),
	}, nil
}

// Delete 删除文件
func (s *FileService) Delete(id uint) error {
	file, err := s.fileRepo.Get(id)
	if err != nil {
		return fmt.Errorf("文件不存在: %w", err)
	}

	// 根据文件的存储类型删除物理文件
	if err := s.uploadManager.DeleteFileByStorageType(file.FilePath, file.StorageType); err != nil {
		return fmt.Errorf("删除存储文件失败: %w", err)
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

// CleanupUnusedFiles 清理未使用的文件（超过15天未使用）
func (s *FileService) CleanupUnusedFiles() error {
	// 获取超过15天未使用的文件
	files, err := s.fileRepo.GetUnusedFiles(15)
	if err != nil {
		return fmt.Errorf("获取未使用文件失败: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	// 删除物理文件
	var deletedIDs []uint
	for _, file := range files {
		if err := s.uploadManager.DeleteFileByStorageType(file.FilePath, file.StorageType); err != nil {
			logger.Warn("删除物理文件失败 %s: %v", file.FilePath, err)
			continue
		}
		deletedIDs = append(deletedIDs, file.ID)
	}

	// 批量删除数据库记录
	if len(deletedIDs) > 0 {
		if err := s.fileRepo.DeleteByIDs(deletedIDs); err != nil {
			return fmt.Errorf("删除文件记录失败: %w", err)
		}
	}

	return nil
}
