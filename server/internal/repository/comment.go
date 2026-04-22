/*
项目名称：JeriBlog
文件名称：comment.go
创建时间：2026-04-16 15:00:20

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：评论数据访问层
*/

package repository

import (
	"context"
	"strings"

	"jeri_blog/internal/model"

	"gorm.io/gorm"
)

// CommentRepository 评论仓储
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// ============ 查询方法 ============

// GetByTarget 获取目标的顶级评论列表
func (r *CommentRepository) GetByTarget(ctx context.Context, targetType, targetKey string, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	// 查询顶级评论（包含已删除的记录）
	query := r.db.WithContext(ctx).Unscoped().Model(&model.Comment{}).
		Where("target_type = ? AND target_key = ? AND root_id IS NULL", targetType, targetKey)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		})

	// 只有当 page 和 pageSize 都大于 0 时才应用分页
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err := query.Find(&comments).Error
	return comments, total, err
}

// GetRepliesByRootIDs 批量获取多个根评论的所有回复
func (r *CommentRepository) GetRepliesByRootIDs(ctx context.Context, rootIDs []uint) ([]model.Comment, error) {
	if len(rootIDs) == 0 {
		return []model.Comment{}, nil
	}

	var replies []model.Comment
	err := r.db.WithContext(ctx).Unscoped().
		Where("root_id IN ?", rootIDs).
		Order("created_at ASC").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		}).
		Preload("ReplyUser", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		}).
		Find(&replies).Error

	return replies, err
}

// GetForWeb 获取评论详情（前台）
func (r *CommentRepository) GetForWeb(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.WithContext(ctx).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		}).
		Preload("ReplyUser", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, nickname, avatar, email, website, badge, role, deleted_at")
		}).
		First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// List 获取评论列表（后台管理）
func (r *CommentRepository) List(ctx context.Context, offset, limit int, keyword string, status *int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	// 包含软删除的记录
	query := r.db.WithContext(ctx).Unscoped().Model(&model.Comment{})

	// 关键词搜索：评论内容、用户昵称、用户邮箱、文章标题、评论来源类型
	if keyword != "" {
		// 中文类型映射
		typeMap := map[string]string{
			"文章": "article",
			"页面": "page",
		}

		// 页面标题映射（target_type='page' 时的 target_key）
		pageKeyMap := map[string]string{
			"留言": "message",
			"动态": "moment",
			"友链": "friend",
		}

		// 构建搜索条件
		query = query.Joins("LEFT JOIN users ON users.id = comments.user_id").
			Joins("LEFT JOIN articles ON comments.target_type = 'article' AND articles.slug = comments.target_key")

		// 基础搜索条件
		conditions := "comments.content LIKE ? OR users.nickname LIKE ? OR users.email LIKE ? OR articles.title LIKE ?"
		args := []interface{}{
			"%" + keyword + "%",
			"%" + keyword + "%",
			"%" + keyword + "%",
			"%" + keyword + "%",
		}

		// 如果关键词匹配中文类型，添加对应的英文类型搜索
		if englishType, ok := typeMap[keyword]; ok {
			conditions += " OR comments.target_type = ?"
			args = append(args, englishType)
		}

		// 如果关键词匹配页面标题，添加页面类型+key的搜索
		if pageKey, ok := pageKeyMap[keyword]; ok {
			conditions += " OR (comments.target_type = 'page' AND comments.target_key = ?)"
			args = append(args, pageKey)
		}

		// 特殊处理："动态"可能是 target_type='moment' 或 target_type='page' AND target_key='moment'
		if keyword == "动态" {
			conditions += " OR comments.target_type = 'moment'"
		}

		// 特殊处理："文章已删除" - 搜索 target_type='article' 且文章不存在的评论
		if keyword == "文章已删除" || strings.Contains(keyword, "已删除") {
			conditions += " OR (comments.target_type = 'article' AND articles.id IS NULL)"
		}

		query = query.Where(conditions, args...)
	}

	// 状态过滤
	if status != nil {
		query = query.Where("comments.status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("comments.created_at DESC").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, email, nickname, avatar, badge, deleted_at")
		}).
		Offset(offset).
		Limit(limit).
		Find(&comments).Error

	return comments, total, err
}

// Get 获取评论详情（后台管理）
func (r *CommentRepository) Get(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.WithContext(ctx).Unscoped().
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().Select("id, email, nickname, avatar, badge, deleted_at")
		}).
		First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ExistsByContentURL 检查是否有评论内容引用该文件（包含软删除评论，便于恢复后仍可展示）
func (r *CommentRepository) ExistsByContentURL(url string) (bool, error) {
	var count int64
	err := r.db.Unscoped().Model(&model.Comment{}).
		Where("status = ? AND content LIKE ?", 1, "%"+url+"%").
		Count(&count).Error
	return count > 0, err
}

// FindByContentURL 查找内容引用该文件的评论列表
func (r *CommentRepository) FindByContentURL(url string) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Unscoped().
		Where("status = ? AND content LIKE ?", 1, "%"+url+"%").
		Find(&comments).Error
	return comments, err
}

// ============ 基础CRUD ============

// Create 创建评论
func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// Update 更新评论
func (r *CommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// UpdateStatus 更新评论状态
func (r *CommentRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
	return r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 软删除评论
func (r *CommentRepository) Delete(ctx context.Context, id uint) error {
	// 只删除评论本身，子评论保留
	return r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

// HardDelete 硬删除评论
func (r *CommentRepository) HardDelete(ctx context.Context, id uint) error {
	// 永久删除评论本身，子评论保留
	return r.db.WithContext(ctx).Unscoped().Delete(&model.Comment{}, id).Error
}

// CountByTargetKeys 批量获取多个目标的评论数
func (r *CommentRepository) CountByTargetKeys(ctx context.Context, targetType string, targetKeys []string) (map[string]int64, error) {
	if len(targetKeys) == 0 {
		return make(map[string]int64), nil
	}

	type CountResult struct {
		TargetKey string
		Count     int64
	}

	var results []CountResult
	err := r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Select("target_key, COUNT(*) as count").
		Where("target_type = ? AND target_key IN ? AND status = 1 AND deleted_at IS NULL", targetType, targetKeys).
		Group("target_key").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 构建映射表
	countMap := make(map[string]int64, len(targetKeys))
	for _, key := range targetKeys {
		countMap[key] = 0
	}
	for _, result := range results {
		countMap[result.TargetKey] = result.Count
	}

	return countMap, nil
}
