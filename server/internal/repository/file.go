/*
项目名称：JeriBlog
文件名称：file.go
创建时间：2026-04-16 15:00:20

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：文件数据访问层
*/

package repository

import (
	"jeri_blog/internal/model"
	"time"

	"gorm.io/gorm"
)

// FileRepository 文件仓储
type FileRepository struct {
	db *gorm.DB
}

// NewFileRepository 创建文件仓储
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

// ============ 基础CRUD ============

// Create 创建文件记录
func (r *FileRepository) Create(file *model.File) error {
	return r.db.Create(file).Error
}

// Get 获取文件信息
func (r *FileRepository) Get(id uint) (*model.File, error) {
	var file model.File
	err := r.db.First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// Delete 删除文件记录
func (r *FileRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&model.File{}, id).Error
}

// ============ 查询方法 ============

// List 获取文件列表
func (r *FileRepository) List(offset, limit int, keyword string, status *int, uploadType string) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	// 构建查询
	query := r.db.Model(&model.File{})

	// 关键词搜索（文件名、原始文件名）
	if keyword != "" {
		query = query.Where("file_name LIKE ? OR original_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 上传类型筛选
	if uploadType != "" {
		query = query.Where("upload_type LIKE ?", "%"+uploadType+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&files).Error

	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// GetByUploadType 根据上传类型获取文件列表
func (r *FileRepository) GetByUploadType(uploadType string, offset, limit int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := r.db.Model(&model.File{})

	// 如果是 "image"，则查询所有图片类型的文件（通过 file_type 字段）
	if uploadType == "image" {
		query = query.Where("file_type LIKE ?", "image/%")
	} else {
		// 否则按 upload_type 精确匹配
		query = query.Where("upload_type = ?", uploadType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&files).Error

	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// GetByStatus 根据状态获取文件列表
func (r *FileRepository) GetByStatus(status int) ([]model.File, error) {
	var files []model.File
	err := r.db.Where("status = ?", status).Order("created_at ASC").Find(&files).Error
	return files, err
}

// ExistsByURLExcludingID 检查是否存在其他文件记录使用相同的URL（排除指定ID）
func (r *FileRepository) ExistsByURLExcludingID(url string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.File{}).
		Where("file_url = ? AND id != ?", url, excludeID).
		Count(&count).Error
	return count > 0, err
}

// ============ 辅助方法 ============

// UpdateStatus 更新文件使用状态
func (r *FileRepository) UpdateStatus(url string, status int) error {
	result := r.db.Model(&model.File{}).
		Where("file_url = ?", url).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CountByURL 统计指定URL的文件记录数量
func (r *FileRepository) CountByURL(url string) (int64, error) {
	var count int64
	err := r.db.Model(&model.File{}).Where("file_url = ?", url).Count(&count).Error
	return count, err
}

// UpdateFileStatusByUrls 批量更新文件状态
func (r *FileRepository) UpdateFileStatusByUrls(urls []string, status int) error {
	if len(urls) == 0 {
		return nil
	}

	return r.db.Model(&model.File{}).
		Where("file_url IN ?", urls).
		Update("status", status).Error
}

// ============ 维护方法 ============

// GetUnusedFiles 获取超过指定天数未使用的文件
func (r *FileRepository) GetUnusedFiles(days int) ([]model.File, error) {
	var files []model.File
	cutoffTime := time.Now().AddDate(0, 0, -days)

	err := r.db.Where("status = ? AND created_at < ?", 0, cutoffTime).
		Find(&files).Error

	return files, err
}

// DeleteByIDs 批量删除文件记录
func (r *FileRepository) DeleteByIDs(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Unscoped().Delete(&model.File{}, ids).Error
}

// ============ 引用计数方法 ============

// IncrementReferenceCount 增加文件引用计数
func (r *FileRepository) IncrementReferenceCount(url string) error {
	if url == "" {
		return nil
	}
	return r.db.Model(&model.File{}).
		Where("file_url = ?", url).
		UpdateColumn("reference_count", gorm.Expr("reference_count + 1")).Error
}

// DecrementReferenceCount 减少文件引用计数
func (r *FileRepository) DecrementReferenceCount(url string) error {
	if url == "" {
		return nil
	}
	return r.db.Model(&model.File{}).
		Where("file_url = ? AND reference_count > 0", url).
		UpdateColumn("reference_count", gorm.Expr("reference_count - 1")).Error
}

// UpdateStatusByReferenceCount 根据引用计数更新文件状态
func (r *FileRepository) UpdateStatusByReferenceCount(url string) error {
	if url == "" {
		return nil
	}
	// 使用子查询根据引用计数更新状态
	return r.db.Exec(`
		UPDATE files
		SET status = CASE
			WHEN reference_count > 0 THEN 1
			ELSE 0
		END,
		updated_at = NOW()
		WHERE file_url = ?
	`, url).Error
}

// GetReferenceCount 获取文件引用计数
func (r *FileRepository) GetReferenceCount(url string) (int, error) {
	if url == "" {
		return 0, nil
	}
	var file model.File
	err := r.db.Select("reference_count").Where("file_url = ?", url).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return file.ReferenceCount, nil
}

