package repository

import (
	"errors"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// SettingRepository 配置仓库
type SettingRepository struct {
	db *gorm.DB
}

// NewSettingRepository 创建配置仓库
func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetByGroup 获取某个分组的所有配置
func (r *SettingRepository) GetByGroup(group string, isPublicOnly ...bool) (map[string]string, error) {
	var settings []model.Setting
	query := r.db.Where("\"group\" = ?", group)

	// 如果指定只返回公开配置，则添加过滤条件
	if len(isPublicOnly) > 0 && isPublicOnly[0] {
		query = query.Where("is_public = ?", true)
	}

	err := query.Order("id ASC").Find(&settings).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

// UpdateGroup 更新配置
func (r *SettingRepository) UpdateGroup(updates map[string]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range updates {
			var setting model.Setting
			err := tx.Where("\"key\" = ?", key).First(&setting).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			if err != nil {
				return err
			}
			// 存在则更新 value
			if err := tx.Model(&setting).Update("value", value).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByKeys 根据键列表获取配置
func (r *SettingRepository) GetByKeys(keys []string) (map[string]string, error) {
	var settings []model.Setting
	err := r.db.Where("\"key\" IN ?", keys).Find(&settings).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}
