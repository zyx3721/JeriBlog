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
	"sort"
	"strings"
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

// GetByURLPattern 通过URL模式查找文件（支持精确匹配和模糊匹配）
// 优先返回精确匹配的结果，如果没有则返回模糊匹配的结果
func (r *FileRepository) GetByURLPattern(urlPattern string) (*model.File, error) {
	var files []model.File
	// 使用 OR 条件同时查询精确匹配和模糊匹配
	err := r.db.Where("file_url = ? OR file_url LIKE ?", urlPattern, "%"+urlPattern).
		Find(&files).Error
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// 优先返回精确匹配的结果
	for _, file := range files {
		if file.FileURL == urlPattern {
			return &file, nil
		}
	}

	// 没有精确匹配，返回第一个模糊匹配的结果
	return &files[0], nil
}

// Delete 删除文件记录
func (r *FileRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&model.File{}, id).Error
}

// ============ 查询方法 ============

// List 获取文件列表
func (r *FileRepository) List(offset, limit int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	if err := r.db.Model(&model.File{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&files).Error

	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// FileListFilter 文件列表筛选条件
type FileListFilter struct {
	Keyword    string
	FileType   string
	Status     *int
	UploadType string
	MinSize    int64
	MaxSize    int64
	StartTime  string
	EndTime    string
}

// GetByFilter 根据筛选条件获取文件列表
func (r *FileRepository) GetByFilter(filter *FileListFilter, offset, limit int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := r.db.Model(&model.File{})

	if filter.Keyword != "" {
		query = query.Where("file_name LIKE ? OR original_name LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	if filter.FileType != "" {
		query = query.Where("file_type LIKE ?", filter.FileType+"%")
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.UploadType != "" {
		query = query.Where("upload_type LIKE ?", "%"+filter.UploadType+"%")
	}
	if filter.MinSize > 0 {
		query = query.Where("file_size >= ?", filter.MinSize)
	}
	if filter.MaxSize > 0 {
		query = query.Where("file_size <= ?", filter.MaxSize)
	}
	if filter.StartTime != "" {
		query = query.Where("created_at >= ?", filter.StartTime)
	}
	if filter.EndTime != "" {
		query = query.Where("created_at <= ?", filter.EndTime+" 23:59:59")
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
	return r.db.Model(&model.File{}).
		Where("file_url = ?", url).
		Update("status", status).Error
}

// CountByURL 统计指定URL的文件记录数量
// 支持相对路径和完整URL的匹配
func (r *FileRepository) CountByURL(url string) (int64, error) {
	var count int64
	// 同时匹配完整URL和相对路径（使用 LIKE 匹配以 url 结尾的记录）
	err := r.db.Model(&model.File{}).
		Where("file_url = ? OR file_url LIKE ?", url, "%"+url).
		Count(&count).Error
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
// 支持相对路径和完整URL的匹配
func (r *FileRepository) IncrementReferenceCount(url string) error {
	if url == "" {
		return nil
	}
	return r.db.Model(&model.File{}).
		Where("file_url = ? OR file_url LIKE ?", url, "%"+url).
		UpdateColumn("reference_count", gorm.Expr("reference_count + 1")).Error
}

// DecrementReferenceCount 减少文件引用计数
// 支持相对路径和完整URL的匹配
func (r *FileRepository) DecrementReferenceCount(url string) error {
	if url == "" {
		return nil
	}
	return r.db.Model(&model.File{}).
		Where("(file_url = ? OR file_url LIKE ?) AND reference_count > 0", url, "%"+url).
		UpdateColumn("reference_count", gorm.Expr("reference_count - 1")).Error
}

// UpdateStatusByReferenceCount 根据引用计数更新文件状态
// 支持相对路径和完整URL的匹配
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
		WHERE file_url = ? OR file_url LIKE ?
	`, url, "%"+url).Error
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

// AddUploadType 添加用途到 upload_type 字段（逗号分隔，不去重，支持同用途多次引用）
// 支持相对路径和完整URL的匹配
func (r *FileRepository) AddUploadType(url string, usageType string) error {
	if url == "" || usageType == "" {
		return nil
	}

	var file model.File
	if err := r.db.Where("file_url = ? OR file_url LIKE ?", url, "%"+url).First(&file).Error; err != nil {
		return err
	}

	// 解析现有用途（保留所有记录，不去重）
	var types []string
	if file.UploadType != "" {
		for _, t := range strings.Split(file.UploadType, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				types = append(types, trimmed)
			}
		}
	}

	// 直接添加新用途（允许重复）
	types = append(types, usageType)
	sort.Strings(types) // 排序保证一致性

	return r.db.Model(&model.File{}).
		Where("file_url = ?", url).
		Update("upload_type", strings.Join(types, ",")).
		Error
}

// RemoveUploadType 从 upload_type 字段移除指定用途（仅移除一次，支持同用途多次引用）
// 支持相对路径和完整URL的匹配
func (r *FileRepository) RemoveUploadType(url string, usageType string) error {
	if url == "" || usageType == "" {
		return nil
	}

	var file model.File
	if err := r.db.Where("file_url = ? OR file_url LIKE ?", url, "%"+url).First(&file).Error; err != nil {
		return err
	}

	// 解析现有用途（保留所有记录）
	var types []string
	if file.UploadType != "" {
		for _, t := range strings.Split(file.UploadType, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				types = append(types, trimmed)
			}
		}
	}

	// 仅移除第一个匹配的用途（支持同用途多次引用）
	removed := false
	var newTypes []string
	for _, t := range types {
		if !removed && t == usageType {
			removed = true
			continue
		}
		newTypes = append(newTypes, t)
	}

	sort.Strings(newTypes) // 排序保证一致性

	return r.db.Model(&model.File{}).
		Where("file_url = ? OR file_url LIKE ?", url, "%"+url).
		Update("upload_type", strings.Join(newTypes, ",")).
		Error
}

