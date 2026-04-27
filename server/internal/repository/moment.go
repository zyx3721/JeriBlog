/*
项目名称：JeriBlog
文件名称：moment.go
创建时间：2026-04-16 15:00:20

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：动态数据访问层
*/

package repository

import (
	"context"
	"strings"

	"jeri_blog/internal/model"

	"gorm.io/gorm"
)

// MomentRepository 动态仓储
type MomentRepository struct {
	db *gorm.DB
}

// NewMomentRepository 创建动态仓储
func NewMomentRepository(db *gorm.DB) *MomentRepository {
	return &MomentRepository{db: db}
}

// ============ 基础CRUD ============

// ListParams 动态列表查询参数
type ListParams struct {
	IsPublish *bool
	Keyword   string
	Tags      []string
	Location  string
	HasImages *bool
	HasVideo  *bool
	HasMusic  *bool
	HasLink   *bool
	StartTime string
	EndTime   string
}

// List 获取动态列表
func (r *MomentRepository) List(ctx context.Context, page, pageSize int, params ListParams) ([]model.Moment, int64, error) {
	var moments []model.Moment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Moment{})

	// 根据发布状态过滤
	if params.IsPublish != nil {
		query = query.Where("is_publish = ?", *params.IsPublish)
	}

	// 根据关键词搜索（按内容模糊搜索）
	if params.Keyword != "" {
		query = query.Where("content::text LIKE ?", "%"+params.Keyword+"%")
	}

	// 根据标签筛选（多选，任意匹配）
	if len(params.Tags) > 0 {
		tagConditions := make([]string, len(params.Tags))
		tagValues := make([]interface{}, len(params.Tags))
		for i, tag := range params.Tags {
			tagConditions[i] = "content::jsonb->>'tags' = ?"
			tagValues[i] = tag
		}
		query = query.Where(strings.Join(tagConditions, " OR "), tagValues...)
	}

	// 根据发布地点筛选
	if params.Location != "" {
		query = query.Where("content::jsonb->>'location' LIKE ?", "%"+params.Location+"%")
	}

	// 根据是否包含图片筛选
	if params.HasImages != nil {
		if *params.HasImages {
			query = query.Where("jsonb_array_length(content::jsonb->'images') > 0")
		} else {
			query = query.Where("(content::jsonb->'images' IS NULL OR jsonb_array_length(content::jsonb->'images') = 0)")
		}
	}

	// 根据是否包含视频筛选
	if params.HasVideo != nil {
		if *params.HasVideo {
			query = query.Where("content::jsonb->'video' IS NOT NULL")
		} else {
			query = query.Where("content::jsonb->'video' IS NULL")
		}
	}

	// 根据是否包含音乐筛选
	if params.HasMusic != nil {
		if *params.HasMusic {
			query = query.Where("content::jsonb->'music' IS NOT NULL")
		} else {
			query = query.Where("content::jsonb->'music' IS NULL")
		}
	}

	// 根据是否包含链接筛选
	if params.HasLink != nil {
		if *params.HasLink {
			query = query.Where("content::jsonb->'link' IS NOT NULL")
		} else {
			query = query.Where("content::jsonb->'link' IS NULL")
		}
	}

	// 根据发布时间范围筛选
	if params.StartTime != "" {
		query = query.Where("DATE(publish_time) >= ?", params.StartTime)
	}
	if params.EndTime != "" {
		query = query.Where("DATE(publish_time) <= ?", params.EndTime)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 排序：优先按发布时间，没有发布时间则按创建时间倒序
	query = query.Order("COALESCE(publish_time, created_at) DESC")

	// 分页处理
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err = query.Find(&moments).Error
	if err != nil {
		return nil, 0, err
	}

	return moments, total, nil
}

// Get 获取动态详情
func (r *MomentRepository) Get(ctx context.Context, id uint) (*model.Moment, error) {
	var moment model.Moment
	err := r.db.WithContext(ctx).First(&moment, id).Error
	if err != nil {
		return nil, err
	}
	return &moment, nil
}

// Create 创建动态
func (r *MomentRepository) Create(ctx context.Context, moment *model.Moment) error {
	return r.db.WithContext(ctx).Create(moment).Error
}

// Update 更新动态
func (r *MomentRepository) Update(ctx context.Context, moment *model.Moment) error {
	return r.db.WithContext(ctx).Save(moment).Error
}

// ExistsByContentURL 检查是否有动态内容引用该文件
func (r *MomentRepository) ExistsByContentURL(url string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Moment{}).Where("content::text LIKE ?", "%"+url+"%").Count(&count).Error
	return count > 0, err
}

// FindByContentURL 查找内容引用该文件的动态列表
func (r *MomentRepository) FindByContentURL(url string) ([]model.Moment, error) {
	var moments []model.Moment
	err := r.db.Where("content::text LIKE ?", "%"+url+"%").Find(&moments).Error
	return moments, err
}

// Delete 删除动态
func (r *MomentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Moment{}, id).Error
}
