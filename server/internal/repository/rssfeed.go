/*
项目名称：JeriBlog
文件名称：rssfeed.go
创建时间：2026-04-16 15:00:20

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：RSS 订阅数据访问层
*/

package repository

import (
	"context"
	"time"

	"jeri_blog/internal/model"

	"gorm.io/gorm"
)

// RssFeedRepository RSS订阅仓储
type RssFeedRepository struct {
	db *gorm.DB
}

// NewRssFeedRepository 创建RSS订阅仓储
func NewRssFeedRepository(db *gorm.DB) *RssFeedRepository {
	return &RssFeedRepository{db: db}
}

// List 获取RSS文章列表
func (r *RssFeedRepository) List(ctx context.Context, page, pageSize int, keyword string, isRead *bool, isDeleted *bool, friendID *uint) ([]model.RssArticle, int64, error) {
	var articles []model.RssArticle
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RssArticle{}).Preload("Friend")

	// 关键词搜索（标题）
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	// 已读状态筛选
	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}

	// 已删除状态筛选
	if isDeleted != nil {
		query = query.Where("is_deleted = ?", *isDeleted)
	}

	// 来源筛选
	if friendID != nil {
		query = query.Where("friend_id = ?", *friendID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("published_at DESC NULLS LAST, created_at DESC")

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetByID 根据ID获取RSS文章
func (r *RssFeedRepository) GetByID(ctx context.Context, id uint) (*model.RssArticle, error) {
	var article model.RssArticle
	if err := r.db.WithContext(ctx).First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// CreateBatch 批量创建RSS文章（忽略重复）
func (r *RssFeedRepository) CreateBatch(ctx context.Context, articles []model.RssArticle) error {
	if len(articles) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(articles, 100).Error
}

// MarkRead 标记文章已读
func (r *RssFeedRepository) MarkRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllRead 全部标记已读
func (r *RssFeedRepository) MarkAllRead(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("is_read = ?", false).Update("is_read", true)
	return result.RowsAffected, result.Error
}

// CountUnread 统计未读数量
func (r *RssFeedRepository) CountUnread(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("is_read = ?", false).Count(&count).Error
	return count, err
}

// ListUnread 获取未读文章列表
func (r *RssFeedRepository) ListUnread(ctx context.Context, limit int) ([]model.RssArticle, error) {
	var articles []model.RssArticle
	query := r.db.WithContext(ctx).Model(&model.RssArticle{}).
		Preload("Friend").
		Where("is_read = ?", false).
		Order("published_at DESC NULLS LAST, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

// DeleteOrphaned 删除孤立文章（友链不存在、失效或RSS地址为空）
func (r *RssFeedRepository) DeleteOrphaned(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("friend_id NOT IN (SELECT id FROM friends) OR friend_id IN (SELECT id FROM friends WHERE is_invalid = true OR rss_url = '' OR rss_url IS NULL)").
		Delete(&model.RssArticle{})
	return result.RowsAffected, result.Error
}

// DeleteByFriendID 删除指定友链的所有RSS文章
func (r *RssFeedRepository) DeleteByFriendID(ctx context.Context, friendID uint) (int64, error) {
	result := r.db.WithContext(ctx).Where("friend_id = ?", friendID).Delete(&model.RssArticle{})
	return result.RowsAffected, result.Error
}

// ExistsByLink 检查链接是否已存在
func (r *RssFeedRepository) ExistsByLink(ctx context.Context, link string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("link = ?", link).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetFriendsWithRSS 获取所有配置了RSS且未失效的友链
func (r *RssFeedRepository) GetFriendsWithRSS(ctx context.Context) ([]model.Friend, error) {
	var friends []model.Friend
	err := r.db.WithContext(ctx).Where("rss_url != '' AND rss_url IS NOT NULL AND is_invalid = ?", false).Find(&friends).Error
	return friends, err
}

// UpdateFriendRSSLatime 更新友链RSS最后更新时间
func (r *RssFeedRepository) UpdateFriendRSSLatime(ctx context.Context, friendID uint, latime time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Friend{}).Where("id = ?", friendID).Update("rss_latime", latime).Error
}

// GetLatestPublishedTime 获取指定友链最新文章的发布时间
func (r *RssFeedRepository) GetLatestPublishedTime(ctx context.Context, friendID uint) (*time.Time, error) {
	var article model.RssArticle
	err := r.db.WithContext(ctx).
		Where("friend_id = ?", friendID).
		Where("published_at IS NOT NULL").
		Order("published_at DESC").
		Limit(1).
		Find(&article).Error
	if err != nil {
		return nil, err
	}
	if article.ID == 0 {
		return nil, nil
	}
	return article.PublishedAt, nil
}

// GetByLink 根据链接获取RSS文章
func (r *RssFeedRepository) GetByLink(ctx context.Context, link string) (*model.RssArticle, error) {
	var article model.RssArticle
	err := r.db.WithContext(ctx).Where("link = ?", link).First(&article).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// MarkDeletedByFriendAndLinks 标记指定友链中不在RSS源的文章为已删除
func (r *RssFeedRepository) MarkDeletedByFriendAndLinks(ctx context.Context, friendID uint, rssLinks map[string]bool) error {
	// 获取该友链的所有文章
	var articles []model.RssArticle
	if err := r.db.WithContext(ctx).Where("friend_id = ?", friendID).Find(&articles).Error; err != nil {
		return err
	}

	// 标记不在RSS源中的文章为已删除
	for _, article := range articles {
		if !rssLinks[article.Link] && !article.IsDeleted {
			if err := r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("id = ?", article.ID).Update("is_deleted", true).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// RestoreArticle 恢复已删除的文章为未读状态
func (r *RssFeedRepository) RestoreArticle(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.RssArticle{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": false,
		"is_read":    false,
	}).Error
}
