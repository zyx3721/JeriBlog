/*
项目名称：JeriBlog
文件名称：article.go
创建时间：2026-04-16 15:00:03

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：文章业务逻辑
*/

package service

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"flec_blog/config"
	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/logger"
	"flec_blog/pkg/random"
	"flec_blog/pkg/utils"
	"flec_blog/pkg/wechatmp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gorm.io/gorm"
)

// ArticleService 文章服务
type ArticleService struct {
	articleRepo       *repository.ArticleRepository
	tagRepo           *repository.TagRepository
	categoryRepo      *repository.CategoryRepository
	commentRepo       *repository.CommentRepository
	fileService       *FileService
	subscriberService *SubscriberService
	db                *gorm.DB
	config            *config.Config // 配置对象（支持热重载）
	md                goldmark.Markdown
	httpClient        *http.Client
}

// NewArticleService 创建文章服务实例
func NewArticleService(articleRepo *repository.ArticleRepository, tagRepo *repository.TagRepository, categoryRepo *repository.CategoryRepository, commentRepo *repository.CommentRepository, fileService *FileService, db *gorm.DB, cfg *config.Config) *ArticleService {
	// 初始化 goldmark（用于微信导出时渲染 Markdown）
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	return &ArticleService{
		articleRepo:  articleRepo,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
		commentRepo:  commentRepo,
		fileService:  fileService,
		db:           db,
		config:       cfg,
		md:           md,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// SetSubscriberService 设置订阅者服务（避免循环依赖）
func (s *ArticleService) SetSubscriberService(subscriberService *SubscriberService) {
	s.subscriberService = subscriberService
}

// ============ 通用服务 ============

// Get 获取文章详情
func (s *ArticleService) Get(_ context.Context, id uint) (*dto.ArticleAdminDetailResponse, error) {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("获取文章失败: %w", err)
	}

	response := &dto.ArticleAdminDetailResponse{
		ID:          article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Content:     article.Content,
		Summary:     article.Summary,
		AISummary:   article.AISummary,
		Cover:       article.Cover,
		Location:    article.Location,
		IsPublish:   article.IsPublish,
		IsTop:       article.IsTop,
		IsEssence:   article.IsEssence,
		IsOutdated:  article.IsOutdated,
		PublishTime: utils.ToJSONTime(article.PublishTime),
		UpdateTime:  utils.ToJSONTime(article.UpdateTime),
	}

	// 填充分类信息
	response.Category.ID = article.Category.ID
	response.Category.Name = article.Category.Name

	// 填充标签信息
	for _, tag := range article.Tags {
		response.Tags = append(response.Tags, struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}{tag.ID, tag.Name})
	}

	return response, nil
}

// ============ 前台服务 ============

// ListForWeb 获取前台文章列表
func (s *ArticleService) ListForWeb(ctx context.Context, req *dto.ListArticlesRequest) ([]dto.ArticleWebResponse, int64, error) {
	articles, total, err := s.articleRepo.ListForWeb(req.Page, req.PageSize, req.Year, req.Month, req.Category, req.Tag)
	if err != nil {
		return nil, 0, err
	}

	// 批量获取文章评论数
	articleSlugs := make([]string, len(articles))
	for i, article := range articles {
		articleSlugs[i] = article.Slug
	}

	commentCounts := make(map[string]int64)
	if len(articleSlugs) > 0 && s.commentRepo != nil {
		commentCounts, err = s.commentRepo.CountByTargetKeys(ctx, "article", articleSlugs)
		if err != nil {
			// 如果获取评论数失败，不影响主流程，只记录错误
			commentCounts = make(map[string]int64)
		}
	}

	// 转换为前台响应格式
	response := make([]dto.ArticleWebResponse, 0)
	for _, article := range articles {
		item := dto.ArticleWebResponse{
			ID:           article.ID,
			Title:        article.Title,
			Summary:      article.Summary,
			Cover:        article.Cover,
			Location:     article.Location,
			IsTop:        article.IsTop,
			IsEssence:    article.IsEssence,
			IsOutdated:   article.IsOutdated,
			URL:          fmt.Sprintf("/posts/%s", article.Slug),
			CommentCount: commentCounts[article.Slug],
			PublishTime:  utils.ToJSONTime(article.PublishTime),
			UpdateTime:   utils.ToJSONTime(article.UpdateTime),
		}

		// 填充分类信息
		if article.Category.ID > 0 {
			item.Category.ID = article.Category.ID
			item.Category.Name = article.Category.Name
			item.Category.URL = fmt.Sprintf("/category/%s", article.Category.Slug)
		}

		// 填充标签信息
		for _, tag := range article.Tags {
			item.Tags = append(item.Tags, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
				URL  string `json:"url"`
			}{
				ID:   tag.ID,
				Name: tag.Name,
				URL:  fmt.Sprintf("/tag/%s", tag.Slug),
			})
		}

		response = append(response, item)
	}

	return response, total, nil
}

// Search 搜索文章
func (s *ArticleService) Search(ctx context.Context, req *dto.SearchArticlesRequest) ([]dto.ArticleWebResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	articles, total, err := s.articleRepo.Search(req.Keyword, offset, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 批量获取文章评论数
	articleSlugs := make([]string, len(articles))
	for i, article := range articles {
		articleSlugs[i] = article.Slug
	}

	commentCounts := make(map[string]int64)
	if len(articleSlugs) > 0 && s.commentRepo != nil {
		commentCounts, err = s.commentRepo.CountByTargetKeys(ctx, "article", articleSlugs)
		if err != nil {
			// 如果获取评论数失败，不影响主流程，只记录错误
			commentCounts = make(map[string]int64)
		}
	}

	response := make([]dto.ArticleWebResponse, 0)
	for _, article := range articles {
		item := dto.ArticleWebResponse{
			ID:           article.ID,
			Title:        article.Title,
			Summary:      article.Summary,
			Cover:        article.Cover,
			Location:     article.Location,
			IsTop:        article.IsTop,
			IsEssence:    article.IsEssence,
			URL:          fmt.Sprintf("/posts/%s", article.Slug),
			Excerpt:      utils.GenerateExcerpt(article.Content, req.Keyword, 40), // 生成包含关键词的摘录
			CommentCount: commentCounts[article.Slug],
			PublishTime:  utils.ToJSONTime(article.PublishTime),
			UpdateTime:   utils.ToJSONTime(article.UpdateTime),
		}

		if article.Category.ID > 0 {
			item.Category.ID = article.Category.ID
			item.Category.Name = article.Category.Name
			item.Category.URL = fmt.Sprintf("/category/%s", article.Category.Slug)
		}

		for _, tag := range article.Tags {
			item.Tags = append(item.Tags, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
				URL  string `json:"url"`
			}{
				ID:   tag.ID,
				Name: tag.Name,
				URL:  fmt.Sprintf("/tag/%s", tag.Slug),
			})
		}

		response = append(response, item)
	}

	return response, total, nil
}

// GetBySlug 通过slug获取文章详情
func (s *ArticleService) GetBySlug(ctx context.Context, slug string) (*dto.ArticleDetailResponse, error) {
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// 异步增加浏览数
	go func() {
		_ = s.articleRepo.IncrementViewCount(article.ID)
	}()

	// 获取文章评论数
	var commentCount int64
	if s.commentRepo != nil {
		commentCounts, err := s.commentRepo.CountByTargetKeys(ctx, "article", []string{article.Slug})
		if err == nil {
			commentCount = commentCounts[article.Slug]
		}
	}

	response := &dto.ArticleDetailResponse{
		ID:           article.ID,
		Title:        article.Title,
		Slug:         article.Slug,
		Content:      article.Content,
		Summary:      article.Summary,
		AISummary:    article.AISummary,
		Cover:        article.Cover,
		Location:     article.Location,
		IsTop:        article.IsTop,
		IsEssence:    article.IsEssence,
		IsOutdated:   article.IsOutdated,
		ViewCount:    article.ViewCount,
		CommentCount: commentCount,
		URL:          fmt.Sprintf("/posts/%s", article.Slug),
		PublishTime:  utils.ToJSONTime(article.PublishTime),
		UpdateTime:   utils.ToJSONTime(article.UpdateTime),
	}

	// 填充分类信息
	if article.Category.ID > 0 {
		response.Category.ID = article.Category.ID
		response.Category.Name = article.Category.Name
		response.Category.URL = fmt.Sprintf("/category/%s", article.Category.Slug)
	}

	// 填充标签信息
	for _, tag := range article.Tags {
		response.Tags = append(response.Tags, struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			ID:   tag.ID,
			Name: tag.Name,
			URL:  fmt.Sprintf("/tag/%s", tag.Slug),
		})
	}

	// 查询上一篇文章
	if prevArticle, err := s.articleRepo.GetPrevArticle(article.PublishTime); err == nil {
		response.Prev = &struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		}{
			Title: prevArticle.Title,
			URL:   fmt.Sprintf("/posts/%s", prevArticle.Slug),
		}
	}

	// 查询下一篇文章
	if nextArticle, err := s.articleRepo.GetNextArticle(article.PublishTime); err == nil {
		response.Next = &struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		}{
			Title: nextArticle.Title,
			URL:   fmt.Sprintf("/posts/%s", nextArticle.Slug),
		}
	}

	return response, nil
}

// ============ 后台管理服务 ============

// List 获取文章列表
func (s *ArticleService) List(ctx context.Context, req *dto.ListArticlesRequest) ([]dto.ArticleListResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	articles, total, err := s.articleRepo.List(offset, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 批量获取文章评论数
	articleSlugs := make([]string, len(articles))
	for i, article := range articles {
		articleSlugs[i] = article.Slug
	}

	commentCounts := make(map[string]int64)
	if len(articleSlugs) > 0 && s.commentRepo != nil {
		commentCounts, err = s.commentRepo.CountByTargetKeys(ctx, "article", articleSlugs)
		if err != nil {
			// 如果获取评论数失败，不影响主流程
			commentCounts = make(map[string]int64)
		}
	}

	// 转换为后台列表响应格式
	var response []dto.ArticleListResponse
	for _, article := range articles {
		item := dto.ArticleListResponse{
			ID:           article.ID,
			Title:        article.Title,
			Cover:        article.Cover,
			Location:     article.Location,
			IsPublish:    article.IsPublish,
			IsTop:        article.IsTop,
			IsEssence:    article.IsEssence,
			IsOutdated:   article.IsOutdated,
			ViewCount:    article.ViewCount,
			CommentCount: commentCounts[article.Slug],
			PublishTime:  utils.ToJSONTime(article.PublishTime),
			UpdateTime:   utils.ToJSONTime(article.UpdateTime),
		}

		item.Category.ID = article.Category.ID
		item.Category.Name = article.Category.Name

		for _, tag := range article.Tags {
			item.Tags = append(item.Tags, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{tag.ID, tag.Name})
		}

		response = append(response, item)
	}

	return response, total, nil
}

// Create 创建文章
func (s *ArticleService) Create(ctx context.Context, req *dto.CreateArticleRequest) (*dto.ArticleAdminDetailResponse, error) {
	// 验证分类是否存在
	if req.CategoryID != nil && *req.CategoryID > 0 {
		_, err := s.categoryRepo.Get(ctx, *req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("分类不存在: %w", err)
		}
	}

	article := &model.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		Location:   req.Location,
		CategoryID: req.CategoryID,
	}

	// 设置置顶状态
	if req.IsTop != nil {
		article.IsTop = *req.IsTop
	}

	// 设置精选状态
	if req.IsEssence != nil {
		article.IsEssence = *req.IsEssence
	}

	// 设置过时状态
	if req.IsOutdated != nil {
		article.IsOutdated = *req.IsOutdated
	}

	// 设置发布状态
	if req.IsPublish != nil {
		article.IsPublish = *req.IsPublish
	}

	// 如果是发布状态，自动设置发布时间
	if article.IsPublish {
		now := utils.Now().Time
		article.PublishTime = &now
	}

	// 处理 slug：如果用户提供了则使用，否则自动生成
	if req.Slug != "" {
		// 检查用户提供的 slug 是否已存在
		exists, err := s.articleRepo.CheckSlugExists(req.Slug)
		if err != nil {
			return nil, fmt.Errorf("检查 slug 失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("slug 已存在: %s", req.Slug)
		}
		article.Slug = req.Slug
	} else {
		// 自动生成唯一 slug
		generatedSlug, err := random.UniqueCode(8, s.articleRepo.CheckSlugExists)
		if err != nil {
			return nil, fmt.Errorf("生成 slug 失败: %w", err)
		}
		article.Slug = generatedSlug
	}

	// 创建文章并关联标签
	if err := s.articleRepo.Create(article, req.TagIDs); err != nil {
		return nil, err
	}

	// 标记封面为使用中
	if req.Cover != "" && s.fileService != nil {
		_ = s.fileService.MarkAsUsed(req.Cover)
	}

	// 标记内容中的图片为使用中
	s.markContentImagesAsUsed(req.Content)

	// 如果是发布状态，异步发送订阅推送
	if article.IsPublish && s.subscriberService != nil {
		go func() {
			if err := s.subscriberService.SendArticleNotification(context.Background(), article); err != nil {
				logger.Warn("发送文章推送失败 (文章ID: %d): %v", article.ID, err)
			}
		}()
	}

	return s.Get(ctx, article.ID)
}

// Update 更新文章
func (s *ArticleService) Update(ctx context.Context, id uint, req *dto.UpdateArticleRequest) (*dto.ArticleAdminDetailResponse, error) {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return nil, err
	}

	// 验证新分类是否存在
	if req.CategoryID != nil && *req.CategoryID > 0 {
		if _, err := s.categoryRepo.Get(ctx, *req.CategoryID); err != nil {
			return nil, fmt.Errorf("分类不存在: %w", err)
		}
	}

	// 保存旧值用于后续处理
	oldCover := article.Cover
	oldContent := article.Content
	oldIsPublish := article.IsPublish

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}

	// 处理 slug 更新
	if req.Slug != "" && req.Slug != article.Slug {
		// 检查新 slug 是否已被其他文章使用
		exists, err := s.articleRepo.CheckSlugExists(req.Slug)
		if err != nil {
			return nil, fmt.Errorf("检查 slug 失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("slug 已存在: %s", req.Slug)
		}
		article.Slug = req.Slug
	}

	article.Summary = req.Summary
	article.AISummary = req.AISummary
	article.Cover = req.Cover
	article.Location = req.Location
	article.CategoryID = req.CategoryID
	if req.IsTop != nil {
		article.IsTop = *req.IsTop
	}

	// 处理精选状态
	if req.IsEssence != nil {
		article.IsEssence = *req.IsEssence
	}

	// 处理过时状态
	if req.IsOutdated != nil {
		article.IsOutdated = *req.IsOutdated
	}

	// 处理发布状态
	if req.IsPublish != nil {
		article.IsPublish = *req.IsPublish
	}

	// 先处理请求中的 PublishTime（仅当传入非空时间时才更新）
	if req.PublishTime != nil && !req.PublishTime.IsZero() {
		article.PublishTime = utils.FromJSONTime(req.PublishTime)
	}

	// 如果是发布状态且没有发布时间，自动设置发布时间
	if article.IsPublish && article.PublishTime == nil {
		now := utils.Now().Time
		article.PublishTime = &now
	}
	if req.UpdateTime != nil {
		article.UpdateTime = utils.FromJSONTime(req.UpdateTime)
	}

	if err := s.articleRepo.Update(article, req.TagIDs); err != nil {
		return nil, err
	}

	// 处理封面变化
	if s.fileService != nil && oldCover != req.Cover {
		if oldCover != "" {
			_ = s.fileService.MarkAsUnused(oldCover)
		}
		if req.Cover != "" {
			_ = s.fileService.MarkAsUsed(req.Cover)
		}
	}

	// 处理内容图片变化
	if req.Content != "" {
		s.updateContentFileStatus(oldContent, req.Content)
	}

	// 如果从草稿变为发布状态，异步发送订阅推送
	if !oldIsPublish && article.IsPublish && s.subscriberService != nil {
		go func() {
			if err := s.subscriberService.SendArticleNotification(context.Background(), article); err != nil {
				logger.Warn("发送文章推送失败 (文章ID: %d): %v", article.ID, err)
			}
		}()
	}

	return s.Get(ctx, id)
}

// Delete 删除文章
func (s *ArticleService) Delete(ctx context.Context, id uint) error {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return err
	}

	// 标记封面为未使用
	if s.fileService != nil && article.Cover != "" {
		_ = s.fileService.MarkAsUnused(article.Cover)
	}

	// 标记内容中的图片为未使用
	s.markContentImagesAsUnused(article.Content)

	return s.articleRepo.Delete(id)
}

// ============ 辅助方法 ============

// extractContentImages 从 Markdown/HTML 内容中提取所有图片 URL
func extractContentImages(content string) []string {
	var urls []string
	seen := make(map[string]bool)

	// 提取 Markdown 图片: ![alt](url)
	mdImageRe := regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
	matches := mdImageRe.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.TrimSpace(match[1])
			if url != "" && !seen[url] {
				seen[url] = true
				urls = append(urls, url)
			}
		}
	}

	// 提取 HTML img 标签: <img src="url" />
	htmlImageRe := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["'][^>]*>`)
	matches = htmlImageRe.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.TrimSpace(match[1])
			if url != "" && !seen[url] {
				seen[url] = true
				urls = append(urls, url)
			}
		}
	}

	return urls
}

// markContentImagesAsUsed 标记内容中的图片为已使用
func (s *ArticleService) markContentImagesAsUsed(content string) {
	if s.fileService == nil {
		return
	}
	for _, url := range extractContentImages(content) {
		_ = s.fileService.MarkAsUsed(url)
	}
}

// markContentImagesAsUnused 标记内容中的图片为未使用
func (s *ArticleService) markContentImagesAsUnused(content string) {
	if s.fileService == nil {
		return
	}
	for _, url := range extractContentImages(content) {
		_ = s.fileService.MarkAsUnused(url)
	}
}

// updateContentFileStatus 对比新旧内容，更新图片文件状态
func (s *ArticleService) updateContentFileStatus(oldContent, newContent string) {
	if s.fileService == nil {
		return
	}

	oldImages := make(map[string]bool)
	for _, url := range extractContentImages(oldContent) {
		oldImages[url] = true
	}

	newImages := make(map[string]bool)
	for _, url := range extractContentImages(newContent) {
		newImages[url] = true
		// 新增的图片标记为使用中
		if !oldImages[url] {
			_ = s.fileService.MarkAsUsed(url)
		}
	}

	// 移除的图片标记为未使用
	for url := range oldImages {
		if !newImages[url] {
			_ = s.fileService.MarkAsUnused(url)
		}
	}
}

// ============ 数据导入导出方法 ============

// ImportArticles 导入文章（支持 Hexo 和 Markdown 格式）
func (s *ArticleService) ImportArticles(ctx context.Context, files map[string]string, sourceType string, uploadImages bool, host string) (*dto.ImportArticlesResult, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("没有找到有效的文章数据")
	}

	result := &dto.ImportArticlesResult{
		Total: len(files),
	}

	// 缓存已创建的分类和标签
	categoryCache := make(map[string]*model.Category)
	tagCache := make(map[string]*model.Tag)

	// 处理每篇文章
	for filename, content := range files {
		if err := s.importSingleArticle(ctx, filename, content, sourceType, uploadImages, host, categoryCache, tagCache); err != nil {
			result.AddError(filename, extractTitle(content), err.Error())
		} else {
			result.Success++
		}
	}

	result.CategoriesAdded = len(categoryCache)
	result.TagsAdded = len(tagCache)

	return result, nil
}

// importSingleArticle 导入单篇文章
func (s *ArticleService) importSingleArticle(
	ctx context.Context,
	filename string,
	content string,
	sourceType string,
	uploadImages bool,
	host string,
	categoryCache map[string]*model.Category,
	tagCache map[string]*model.Tag,
) error {
	// 解析文章
	var parsed *HexoParsedArticle
	var err error

	switch sourceType {
	case "hexo":
		parsed, err = parseHexoArticle(content)
	case "markdown":
		parsed, err = parseMarkdownArticle(filename, content)
	default:
		return fmt.Errorf("不支持的来源类型: %s", sourceType)
	}

	if err != nil {
		return fmt.Errorf("解析失败: %w", err)
	}

	// 处理图片上传
	if uploadImages {
		parsed.Content, err = s.downloadAndUploadImages(ctx, parsed.Content, host)
		if err != nil {
			return fmt.Errorf("图片处理失败: %w", err)
		}
	}

	// 处理分类
	var categoryID *uint
	if parsed.Category != "" {
		category, err := s.getOrCreateCategory(ctx, parsed.Category, categoryCache)
		if err != nil {
			return fmt.Errorf("分类处理失败: %w", err)
		}
		categoryID = &category.ID
	}

	// 处理标签
	var tagIDs []uint
	for _, tagName := range parsed.Tags {
		tag, err := s.getOrCreateTag(ctx, tagName, tagCache)
		if err != nil {
			return fmt.Errorf("标签处理失败: %w", err)
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	// 处理 slug：优先使用原有的，否则生成新的
	articleSlug := parsed.Slug
	if articleSlug != "" {
		if exists, _ := s.articleRepo.CheckSlugExists(articleSlug); exists {
			articleSlug = "" // slug 已存在，需要生成新的
		}
	}
	if articleSlug == "" {
		articleSlug, _ = random.UniqueCode(8, s.articleRepo.CheckSlugExists)
	}

	// 创建文章
	article := &model.Article{
		Title:       parsed.Title,
		Slug:        articleSlug,
		Content:     parsed.Content,
		Summary:     parsed.Summary,
		Cover:       parsed.Cover,
		IsPublish:   false, // 导入的文章默认为草稿
		IsTop:       false,
		CategoryID:  categoryID,
		PublishTime: parsed.PublishTime,
		UpdateTime:  parsed.UpdateTime,
	}

	if err := s.articleRepo.Create(article, tagIDs); err != nil {
		return fmt.Errorf("保存失败: %w", err)
	}

	return nil
}

// ImportFromHexo 从Hexo格式导入文章（保留向后兼容）
func (s *ArticleService) ImportFromHexo(ctx context.Context, files map[string]string) (*dto.ImportArticlesResult, error) {
	return s.ImportArticles(ctx, files, "hexo", false, "")
}

// getOrCreateCategory 获取或创建分类
func (s *ArticleService) getOrCreateCategory(ctx context.Context, name string, cache map[string]*model.Category) (*model.Category, error) {
	// 检查缓存
	if category, exists := cache[name]; exists {
		return category, nil
	}

	// 尝试从数据库获取
	category, err := s.categoryRepo.GetBySlug(ctx, name)
	if err == nil {
		cache[name] = category
		return category, nil
	}

	// 不存在则创建
	category = &model.Category{
		Name:        name,
		Slug:        name,
		Description: "",
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	cache[name] = category
	return category, nil
}

// getOrCreateTag 获取或创建标签
func (s *ArticleService) getOrCreateTag(ctx context.Context, name string, cache map[string]*model.Tag) (*model.Tag, error) {
	// 检查缓存
	if tag, exists := cache[name]; exists {
		return tag, nil
	}

	// 尝试从数据库获取
	tag, err := s.tagRepo.GetBySlug(ctx, name)
	if err == nil {
		cache[name] = tag
		return tag, nil
	}

	// 不存在则创建
	tag = &model.Tag{
		Name:        name,
		Slug:        name,
		Description: "",
	}

	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}

	cache[name] = tag
	return tag, nil
}

// extractTitle 从内容中提取标题
func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "title:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "title:"))
		}
	}
	return "未知标题"
}

// HexoParsedArticle 解析后的Hexo文章
type HexoParsedArticle struct {
	Title       string
	Slug        string
	Content     string
	Summary     string
	Cover       string
	Category    string
	Tags        []string
	PublishTime *time.Time
	UpdateTime  *time.Time
}

// parseHexoArticle 解析Hexo文章格式（Front Matter + Markdown）
func parseHexoArticle(content string) (*HexoParsedArticle, error) {
	var frontMatter string
	var markdown string

	// 检查是否包含Front Matter标记
	if strings.HasPrefix(content, "---") {
		// 分割Front Matter和内容
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			frontMatter = parts[1]
			markdown = strings.TrimSpace(parts[2])
		} else {
			// Front Matter 格式不完整，当作纯 Markdown 处理
			markdown = strings.TrimSpace(content)
		}
	} else {
		// 纯 Markdown 文件，没有 Front Matter
		markdown = strings.TrimSpace(content)
	}

	// 转换 HTML img 标签为 Markdown 格式
	markdown = convertHTMLImagesToMarkdown(markdown)

	// 解析Front Matter
	parsed := &HexoParsedArticle{
		Content: markdown,
	}

	// 如果有 Front Matter，解析它
	if frontMatter != "" {
		lines := strings.Split(frontMatter, "\n")
		var tagLines []string
		inTags := false

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 处理标签数组
			if inTags {
				if strings.HasPrefix(line, "-") {
					tagValue := strings.TrimSpace(strings.TrimPrefix(line, "-"))
					tagValue = strings.Trim(tagValue, "\"'")
					if tagValue != "" {
						tagLines = append(tagLines, tagValue)
					}
				} else {
					inTags = false
				}
			}

			// 解析键值对
			if strings.Contains(line, ":") && !strings.HasPrefix(line, "-") {
				parts := strings.SplitN(line, ":", 2)
				key := strings.TrimSpace(parts[0])
				value := ""
				if len(parts) > 1 {
					value = strings.TrimSpace(parts[1])
					value = strings.Trim(value, "\"'")
				}

				switch key {
				case "title":
					parsed.Title = value
				case "date":
					if t, err := parseHexoDate(value); err == nil {
						parsed.PublishTime = t
					}
				case "updated":
					if t, err := parseHexoDate(value); err == nil {
						parsed.UpdateTime = t
					}
				case "categories", "category":
					if value != "" {
						parsed.Category = value
					}
					// 如果value为空，可能是数组格式，下一行开始
				case "tags":
					if value != "" {
						// 内联格式: tags: [tag1, tag2]
						value = strings.Trim(value, "[]")
						for _, tag := range strings.Split(value, ",") {
							tag = strings.TrimSpace(tag)
							tag = strings.Trim(tag, "\"'")
							if tag != "" {
								parsed.Tags = append(parsed.Tags, tag)
							}
						}
					} else {
						// 数组格式
						inTags = true
					}
				case "cover", "thumbnail":
					parsed.Cover = value
				case "description", "excerpt":
					parsed.Summary = value
				case "slug", "abbrlink":
					parsed.Slug = value
				}
			}
		}

		// 添加收集的标签
		if len(tagLines) > 0 {
			parsed.Tags = append(parsed.Tags, tagLines...)
		}
	}

	// 如果没有标题，尝试从 Markdown 内容中提取第一个标题
	if parsed.Title == "" {
		parsed.Title = extractTitleFromMarkdown(markdown)
	}

	// 如果还是没有标题，返回错误
	if parsed.Title == "" {
		return nil, fmt.Errorf("文章缺少标题（需要在 Front Matter 中指定 title 或在内容中使用 # 标题）")
	}

	// 如果没有摘要，从内容中生成
	if parsed.Summary == "" {
		parsed.Summary = generateSummary(parsed.Content, 200)
	}

	return parsed, nil
}

// parseHexoDate 解析Hexo日期格式
func parseHexoDate(dateStr string) (*time.Time, error) {
	// 支持多种日期格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("无法解析日期: %s", dateStr)
}

// generateSummary 从内容生成摘要
func generateSummary(content string, maxLen int) string {
	// 移除Markdown标记
	content = strings.ReplaceAll(content, "#", "")
	content = strings.ReplaceAll(content, "*", "")
	content = strings.ReplaceAll(content, "`", "")
	content = strings.ReplaceAll(content, "\n", " ")
	content = strings.TrimSpace(content)

	// 截取指定长度
	runes := []rune(content)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return content
}

// extractTitleFromMarkdown 从 Markdown 内容中提取第一个标题
func extractTitleFromMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 匹配 # 标题
		if strings.HasPrefix(line, "#") {
			title := strings.TrimSpace(strings.TrimLeft(line, "#"))
			if title != "" {
				return title
			}
		}
	}
	return ""
}

// parseMarkdownArticle 解析纯 Markdown 格式文章
func parseMarkdownArticle(filename, content string) (*HexoParsedArticle, error) {
	parsed := &HexoParsedArticle{
		Tags:        []string{},
		PublishTime: nil,
		UpdateTime:  nil,
	}

	// 从文件名提取标题
	if filename != "" {
		lowerName := strings.ToLower(filename)
		if strings.HasSuffix(lowerName, ".md") {
			parsed.Title = strings.TrimSpace(filename[:len(filename)-3])
		} else if strings.HasSuffix(lowerName, ".markdown") {
			parsed.Title = strings.TrimSpace(filename[:len(filename)-9])
		} else {
			parsed.Title = strings.TrimSpace(filename)
		}
	}

	// 如果文件名没有标题，尝试从内容提取
	if parsed.Title == "" {
		parsed.Title = extractTitleFromMarkdown(content)
	}

	// 如果还是没有标题，使用默认值
	if parsed.Title == "" {
		parsed.Title = "未命名文章"
	}

	parsed.Summary = generateSummary(content, 200)
	parsed.Content = content

	return parsed, nil
}

// downloadAndUploadImages 下载并上传文章中的图片
func (s *ArticleService) downloadAndUploadImages(ctx context.Context, content string, host string) (string, error) {
	if s.fileService == nil {
		return content, nil
	}

	// 提取所有图片 URL
	imageURLs := extractContentImages(content)
	if len(imageURLs) == 0 {
		return content, nil
	}

	// 去重
	uniqueURLs := make(map[string]bool)
	for _, url := range imageURLs {
		uniqueURLs[url] = true
	}

	// 并发下载上传图片
	replacements := make(map[string]string)
	for url := range uniqueURLs {
		// 跳过相对路径和本地路径
		if strings.HasPrefix(url, "./") || strings.HasPrefix(url, "../") || strings.HasPrefix(url, "/") {
			continue
		}

		// 下载并上传图片
		newURL, err := s.downloadAndUploadSingleImage(ctx, url, host)
		if err == nil {
			replacements[url] = newURL
		}
	}

	// 替换内容中的图片 URL
	for oldURL, newURL := range replacements {
		content = strings.ReplaceAll(content, oldURL, newURL)
	}

	return content, nil
}

// downloadAndUploadSingleImage 下载并上传单张图片
func (s *ArticleService) downloadAndUploadSingleImage(ctx context.Context, imgURL string, host string) (string, error) {
	if s.fileService == nil || imgURL == "" {
		return imgURL, nil
	}

	// 跳过相对路径
	if strings.HasPrefix(imgURL, "./") || strings.HasPrefix(imgURL, "../") || strings.HasPrefix(imgURL, "/") {
		return imgURL, nil
	}

	// 下载图片
	data, ext, err := s.fetchImage(ctx, imgURL)
	if err != nil {
		return imgURL, fmt.Errorf("下载图片失败: %w", err)
	}

	// 生成文件名（使用 SHA256 哈希避免重复）
	hashBytes := sha256.Sum256(data)
	hashStr := fmt.Sprintf("%x", hashBytes)[:12]
	filename := fmt.Sprintf("import_%s%s", hashStr, ext)

	// 确定 MIME 类型
	mimeType := "image/jpeg"
	switch strings.ToLower(ext) {
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	case ".avif":
		mimeType = "image/avif"
	case ".svg":
		mimeType = "image/svg+xml"
	case ".bmp":
		mimeType = "image/bmp"
	case ".tiff", ".tif":
		mimeType = "image/tiff"
	}

	// 上传图片
	reader := bytes.NewReader(data)
	uploadedURL, err := s.fileService.UploadFromReader(reader, filename, mimeType, "文章图片", 0, host)
	if err != nil {
		return imgURL, fmt.Errorf("上传图片失败: %w", err)
	}

	// 标记文件为已使用
	if err := s.fileService.MarkAsUsed(uploadedURL); err != nil {
		logger.Warn("标记文件状态失败: %v", err)
	}

	return uploadedURL, nil
}

// ============ 微信公众号导出 ============

// ExportToWeChat 导出文章到微信公众号
func (s *ArticleService) ExportToWeChat(ctx context.Context, id uint) *dto.WeChatExportResult {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return &dto.WeChatExportResult{Success: false}
	}

	// 预处理并渲染 Markdown
	processed := wechatmp.ConvertCustomBlocks(article.Content)
	processed = wechatmp.ConvertLinksToFootnotes(processed)
	processed = wechatmp.PreprocessMarkdown(processed)

	var htmlBuf bytes.Buffer
	if err := s.md.Convert([]byte(processed), &htmlBuf); err != nil {
		return &dto.WeChatExportResult{Success: false}
	}

	// 转换为公众号格式
	result, err := wechatmp.ConvertMarkdownToWeChatHTML(htmlBuf.String())
	if err != nil {
		return &dto.WeChatExportResult{Success: false}
	}
	html := result.HTML

	// 检查微信配置
	if s.config.WeChat.AppID == "" || s.config.WeChat.AppSecret == "" {
		return &dto.WeChatExportResult{Success: false, HTML: html}
	}

	// 创建微信客户端
	client, err := wechatmp.NewClient(wechatmp.Config{
		AppID:     s.config.WeChat.AppID,
		AppSecret: s.config.WeChat.AppSecret,
		BaseURL:   s.config.WeChat.TokenURL,
	})
	if err != nil {
		return &dto.WeChatExportResult{Success: false, HTML: html}
	}

	// 上传图片
	htmlContent := result.HTML
	var warnings []string
	for _, img := range result.Images {
		newURL, err := s.uploadImageToWeChat(ctx, client, img.OriginalURL)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("图片 %s 上传失败", img.OriginalURL))
			continue
		}
		htmlContent = wechatmp.ReplaceImageURL(htmlContent, img.OriginalURL, newURL)
	}

	// 上传封面
	coverURL := article.Cover
	if coverURL == "" {
		coverURL = "https://api.pearktrue.cn/api/bing/"
	}
	thumbMediaID, err := s.uploadCoverToWeChat(ctx, client, coverURL)
	if err != nil {
		return &dto.WeChatExportResult{Success: false, HTML: html, Warnings: warnings}
	}

	// 创建草稿
	author := s.config.Basic.Author
	if author == "" {
		author = s.config.Blog.Title
	}
	draftResult, err := client.CreateDraft(ctx, []wechatmp.DraftArticle{{
		Title:            article.Title,
		Author:           author,
		Content:          htmlContent,
		Digest:           truncateString(article.Summary, 120),
		ContentSourceURL: s.buildArticleURL(article),
		ThumbMediaID:     thumbMediaID,
		NeedOpenComment:  1,
	}})
	if err != nil {
		return &dto.WeChatExportResult{Success: false, HTML: html, Warnings: warnings}
	}

	return &dto.WeChatExportResult{Success: true, MediaID: draftResult.MediaID, Warnings: warnings}
}

// uploadImageToWeChat 上传文章内图片到微信
func (s *ArticleService) uploadImageToWeChat(ctx context.Context, client *wechatmp.Client, imgURL string) (string, error) {
	data, ext, err := s.fetchImage(ctx, imgURL)
	if err != nil {
		return "", err
	}

	filename := "image" + ext
	result, err := client.UploadImage(ctx, filename, data)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

// uploadCoverToWeChat 上传封面图到微信素材库
func (s *ArticleService) uploadCoverToWeChat(ctx context.Context, client *wechatmp.Client, coverURL string) (string, error) {
	data, ext, err := s.fetchImage(ctx, coverURL)
	if err != nil {
		return "", fmt.Errorf("下载封面图失败: %w", err)
	}

	const maxImageSize = 10 * 1024 * 1024
	if len(data) > maxImageSize {
		return "", fmt.Errorf("封面图片过大（%d MB），微信限制最大 10MB", len(data)/1024/1024)
	}

	result, err := client.AddThumbMaterial(ctx, "cover"+ext, data)
	if err != nil {
		return "", fmt.Errorf("上传封面到微信失败: %w", err)
	}

	return result.MediaID, nil
}

// fetchImage 下载图片，返回数据和扩展名
func (s *ArticleService) fetchImage(ctx context.Context, imgURL string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imgURL, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// 从 URL 或 Content-Type 获取扩展名
	ext := ".jpg"
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		switch ct {
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		}
	} else if idx := strings.LastIndex(imgURL, "."); idx > 0 {
		if e := imgURL[idx:]; len(e) <= 5 {
			ext = e
		}
	}

	return data, ext, nil
}

// buildArticleURL 构建文章链接
func (s *ArticleService) buildArticleURL(article *model.Article) string {
	if s.config.Basic.BlogURL != "" {
		return s.config.Basic.BlogURL + "/posts/" + article.Slug
	}
	return ""
}

// truncateString 截断字符串
func truncateString(str string, maxLen int) string {
	runes := []rune(str)
	if len(runes) <= maxLen {
		return str
	}
	return string(runes[:maxLen-3]) + "..."
}

// ============ 文章下载导出 ============

// imageDownloadResult 图片下载结果
type imageDownloadResult struct {
	url      string
	data     []byte
	ext      string
	filename string
	err      error
}

// extractFilenameFromURL 从 URL 中提取文件名并清理非法字符
func extractFilenameFromURL(imgURL string) string {
	// 移除查询参数
	if idx := strings.Index(imgURL, "?"); idx > 0 {
		imgURL = imgURL[:idx]
	}
	// 提取路径最后一部分
	var filename string
	if idx := strings.LastIndex(imgURL, "/"); idx >= 0 && idx < len(imgURL)-1 {
		filename = imgURL[idx+1:]
	}
	if filename == "" {
		return ""
	}
	// 清理文件名中的非法字符
	filename = strings.Map(func(r rune) rune {
		if strings.ContainsRune("<>:\"/\\|?*", r) {
			return '_'
		}
		return r
	}, filename)
	return filename
}

// DownloadZip 下载文章为压缩包
func (s *ArticleService) DownloadZip(ctx context.Context, id uint) ([]byte, string, error) {
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return nil, "", err
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	imageMap := make(map[string]string)

	// 收集所有需要下载的图片 URL（封面 + 内容图片）
	var imageURLs []string
	if article.Cover != "" {
		imageURLs = append(imageURLs, article.Cover)
	}
	imageURLs = append(imageURLs, s.extractImageURLs(article.Content)...)

	// 去重
	seen := make(map[string]bool)
	var uniqueURLs []string
	for _, url := range imageURLs {
		if !seen[url] {
			seen[url] = true
			uniqueURLs = append(uniqueURLs, url)
		}
	}

	// 如果没有图片，直接生成 Markdown 文件
	if len(uniqueURLs) == 0 {
		frontMatter := s.buildYAMLFrontMatter(article, imageMap)
		mdContent := frontMatter + "\n" + article.Content
		filename := s.sanitizeFilename(article.Title) + ".md"
		if w, _ := zipWriter.Create(filename); w != nil {
			w.Write([]byte(mdContent))
		}
		zipWriter.Close()
		return buf.Bytes(), s.sanitizeFilename(article.Title) + ".zip", nil
	}

	// 并发下载图片（限制并发数为 10）
	const maxConcurrency = 10
	results := make(chan imageDownloadResult, len(uniqueURLs))
	sem := make(chan struct{}, maxConcurrency)

	// 预先为每个 URL 分配文件名（避免并发竞态）
	filenameMap := make(map[string]string)
	filenameCounter := make(map[string]int)
	for _, url := range uniqueURLs {
		// 从 URL 提取原始文件名
		originalName := extractFilenameFromURL(url)
		if originalName == "" {
			// 从 fetchImage 获取扩展名（这里先使用默认）
			originalName = "image.jpg"
		}

		// 处理文件名冲突
		finalName := "assets/" + originalName
		if count, exists := filenameCounter[originalName]; exists {
			// 文件名冲突，添加序号
			nameWithoutExt := originalName
			ext := ""
			if idx := strings.LastIndex(originalName, "."); idx > 0 {
				nameWithoutExt = originalName[:idx]
				ext = originalName[idx:]
			}
			finalName = fmt.Sprintf("assets/%s_%d%s", nameWithoutExt, count+1, ext)
			filenameCounter[originalName] = count + 1
		} else {
			filenameCounter[originalName] = 1
		}

		// 封面图特殊处理
		if url == article.Cover {
			finalName = "assets/cover.jpg" // 默认扩展名，后续会根据实际类型调整
		}

		filenameMap[url] = finalName
	}

	// 并发下载
	for _, url := range uniqueURLs {
		go func(imgURL string) {
			sem <- struct{}{}
			defer func() { <-sem }()

			result := imageDownloadResult{url: imgURL}
			if data, ext, err := s.fetchImage(ctx, imgURL); err == nil {
				result.data = data
				result.ext = ext

				// 获取预分配的文件名，并根据实际扩展名调整
				filename := filenameMap[imgURL]
				// 替换扩展名
				if idx := strings.LastIndex(filename, "."); idx > 0 {
					filename = filename[:idx] + ext
				}
				result.filename = filename
			} else {
				result.err = err
			}
			results <- result
		}(url)
	}

	// 收集结果并写入 zip
	for range uniqueURLs {
		result := <-results
		if result.err != nil {
			continue
		}
		if w, _ := zipWriter.Create(result.filename); w != nil {
			w.Write(result.data)
			imageMap[result.url] = result.filename
		}
	}

	// 替换图片链接
	content := article.Content
	for url, path := range imageMap {
		content = strings.ReplaceAll(content, url, path)
	}

	// 写入 Markdown 文件
	frontMatter := s.buildYAMLFrontMatter(article, imageMap)
	mdContent := frontMatter + "\n" + content
	filename := s.sanitizeFilename(article.Title) + ".md"
	if w, _ := zipWriter.Create(filename); w != nil {
		w.Write([]byte(mdContent))
	}

	zipWriter.Close()
	return buf.Bytes(), s.sanitizeFilename(article.Title) + ".zip", nil
}

// buildYAMLFrontMatter 构建 YAML Front Matter
func (s *ArticleService) buildYAMLFrontMatter(article *model.Article, imageMap map[string]string) string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %q\n", article.Title)
	fmt.Fprintf(&b, "slug: %s\n", article.Slug)

	if article.Summary != "" {
		fmt.Fprintf(&b, "summary: %q\n", article.Summary)
	}
	if article.Cover != "" {
		if path, ok := imageMap[article.Cover]; ok {
			fmt.Fprintf(&b, "cover: %s\n", path)
		} else {
			fmt.Fprintf(&b, "cover: %s\n", article.Cover)
		}
	}
	if article.Location != "" {
		fmt.Fprintf(&b, "location: %q\n", article.Location)
	}

	fmt.Fprintf(&b, "published: %t\n", article.IsPublish)
	fmt.Fprintf(&b, "top: %t\n", article.IsTop)
	fmt.Fprintf(&b, "essence: %t\n", article.IsEssence)
	fmt.Fprintf(&b, "outdated: %t\n", article.IsOutdated)

	if article.Category.ID > 0 {
		fmt.Fprintf(&b, "category: %q\n", article.Category.Name)
	}
	if len(article.Tags) > 0 {
		b.WriteString("tags:\n")
		for _, tag := range article.Tags {
			fmt.Fprintf(&b, "  - %q\n", tag.Name)
		}
	}
	if article.PublishTime != nil {
		fmt.Fprintf(&b, "date: %s\n", article.PublishTime.Format("2006-01-02 15:04:05"))
	}
	if article.UpdateTime != nil {
		fmt.Fprintf(&b, "updated: %s\n", article.UpdateTime.Format("2006-01-02 15:04:05"))
	}

	b.WriteString("---\n")
	return b.String()
}

// extractImageURLs 提取 Markdown 中的图片 URL
func (s *ArticleService) extractImageURLs(content string) []string {
	re := regexp.MustCompile(`!\[.*?\]\((https?://[^)]+)\)`)
	matches := re.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	urls := make([]string, 0, len(matches))
	for _, m := range matches {
		if !seen[m[1]] {
			seen[m[1]] = true
			urls = append(urls, m[1])
		}
	}
	return urls
}

// sanitizeFilename 清理文件名
func (s *ArticleService) sanitizeFilename(name string) string {
	result := strings.Map(func(r rune) rune {
		if strings.ContainsRune("<>:\"/\\|?*", r) {
			return '_'
		}
		return r
	}, name)

	if len([]rune(result)) > 100 {
		result = string([]rune(result)[:100])
	}
	return result
}

// convertHTMLImagesToMarkdown 将 HTML <img> 标签转换为 Markdown 格式
func convertHTMLImagesToMarkdown(content string) string {
	// 正则表达式匹配 <img> 标签，提取 src 和 alt 属性
	// 支持多种格式：
	// <img src="url" alt="text" />
	// <img alt="text" src="url" />
	// <img src="url" />
	// <img src='url' alt='text' style="..." />
	imgRegex := regexp.MustCompile(`<img\s+[^>]*?>`)

	result := imgRegex.ReplaceAllStringFunc(content, func(imgTag string) string {
		// 提取 src 属性
		srcRegex := regexp.MustCompile(`src\s*=\s*["']([^"']+)["']`)
		srcMatch := srcRegex.FindStringSubmatch(imgTag)
		if len(srcMatch) < 2 {
			// 没有找到 src 属性，保持原样
			return imgTag
		}
		src := srcMatch[1]

		// 提取 alt 属性（可选）
		altRegex := regexp.MustCompile(`alt\s*=\s*["']([^"']*?)["']`)
		altMatch := altRegex.FindStringSubmatch(imgTag)
		alt := ""
		if len(altMatch) >= 2 {
			alt = altMatch[1]
		}

		// 转换为 Markdown 格式：![alt](src)
		return fmt.Sprintf("![%s](%s)", alt, src)
	})

	return result
}
