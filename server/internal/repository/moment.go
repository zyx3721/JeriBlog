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
	"encoding/json"
	"regexp"
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

// List 获取动态列表
func (r *MomentRepository) List(
	ctx context.Context,
	page, pageSize int,
	Tags []string,
	keyword, location string,
	isPublish, hasImages, hasVideo, hasMusic, hasLink *bool,
	startTime, endTime string,
) ([]model.Moment, int64, error) {
	var moments []model.Moment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Moment{})

	// 关键词搜索
	if keyword != "" {
		pattern := `"text"\s*:\s*"[^"]*` + regexp.QuoteMeta(keyword) + `[^"]*"`
		query = query.Where("content::text ~* ?", pattern)
	}

	// 根据标签筛选（多选，任意匹配）
	if len(Tags) > 0 {
		tagConditions := make([]string, len(Tags))
		tagValues := make([]interface{}, len(Tags))
		for i, tag := range Tags {
			tagConditions[i] = "content::jsonb->>'tags' = ?"
			tagValues[i] = tag
		}
		query = query.Where(strings.Join(tagConditions, " OR "), tagValues...)
	}

	// 按发布地点筛选
	if location != "" {
		pattern := `"location"\s*:\s*"[^"]*` + regexp.QuoteMeta(location) + `[^"]*"`
		query = query.Where("content::text ~* ?", pattern)
	}

	// 按发布状态筛选
	if isPublish != nil {
		query = query.Where("is_publish = ?", *isPublish)
	}

	// 按是否有图片筛选
	if hasImages != nil {
		if *hasImages {
			query = query.Where("content::text ~* ?", `"images"\s*:\s*\[.+\]`)
		} else {
			query = query.Where(
				"content::text IS NULL OR content::text !~* ? OR content::text ~* ?",
				`"images"`,
				`"images"\s*:\s*(null|\[\s*\])`,
			)
		}
	}

	// 按是否有视频筛选
	if hasVideo != nil {
		if *hasVideo {
			query = query.Where("content::text ILIKE '%\"video\":%'")
		} else {
			query = query.Where("content::text NOT ILIKE '%\"video\":%' OR content::text ILIKE '%\"video\":{}%'")
		}
	}

	// 按是否有音乐筛选（包含音乐和音频）
	if hasMusic != nil {
		if *hasMusic {
			query = query.Where("(content::text ILIKE '%\"music\":%' AND content::text NOT ILIKE '%\"music\":{}%') OR (content::text ILIKE '%\"audio\":%' AND content::text NOT ILIKE '%\"audio\":{}%')")
		} else {
			query = query.Where("(content::text NOT ILIKE '%\"music\":%' OR content::text ILIKE '%\"music\":{}%') AND (content::text NOT ILIKE '%\"audio\":%' OR content::text ILIKE '%\"audio\":{}%')")
		}
	}

	// 按是否有链接筛选
	if hasLink != nil {
		if *hasLink {
			query = query.Where("content::text ILIKE '%\"link\":%'")
		} else {
			query = query.Where("content::text NOT ILIKE '%\"link\":%' OR content::text ILIKE '%\"link\":{}%'")
		}
	}

	// 按发布时间范围筛选
	if startTime != "" {
		query = query.Where("publish_time >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("publish_time <= ?", endTime+" 23:59:59")
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

// MomentFileReference 动态文件引用信息
type MomentFileReference struct {
	Moment   model.Moment
	FileType string // "动态配图" / "动态视频" / "动态音频"
}

// FindByContentURLWithType 查找内容引用该文件的动态列表（带文件类型）
func (r *MomentRepository) FindByContentURLWithType(url string) ([]MomentFileReference, error) {
	var moments []model.Moment
	var references []MomentFileReference

	// 先查找所有包含该 URL 的动态
	err := r.db.Where("content::text LIKE ?", "%"+url+"%").Find(&moments).Error
	if err != nil {
		return nil, err
	}

	// 解析每个动态的 content JSON，判断文件类型
	for _, moment := range moments {
		var contentData map[string]interface{}
		if err := json.Unmarshal([]byte(moment.Content), &contentData); err != nil {
			continue
		}

		fileType := ""

		// 判断是视频（video 是对象，需要提取 url 字段）
		if videoObj, ok := contentData["video"].(map[string]interface{}); ok {
			if videoURL, ok := videoObj["url"].(string); ok && videoURL == url {
				fileType = "动态视频"
			}
		}

		// 判断是音频（audio 是对象，需要提取 url 字段）
		if fileType == "" {
			if audioObj, ok := contentData["audio"].(map[string]interface{}); ok {
				if audioURL, ok := audioObj["url"].(string); ok && audioURL == url {
					fileType = "动态音频"
				}
			}
		}

		// 判断是配图
		if fileType == "" {
			if images, ok := contentData["images"].([]interface{}); ok {
				for _, img := range images {
					if imgURL, ok := img.(string); ok && imgURL == url {
						fileType = "动态配图"
						break
					}
				}
			}
		}

		if fileType != "" {
			references = append(references, MomentFileReference{
				Moment:   moment,
				FileType: fileType,
			})
		}
	}

	return references, nil
}

// Delete 删除动态
func (r *MomentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Moment{}, id).Error
}
