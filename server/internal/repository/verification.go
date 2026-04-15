package repository

import (
	"time"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// VerificationRepository 验证码仓储
type VerificationRepository struct {
	db *gorm.DB
}

// NewVerificationRepository 创建验证码仓储
func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

// ============ 基础CRUD ============

// Create 创建验证码记录
func (r *VerificationRepository) Create(verification *model.Verification) error {
	return r.db.Create(verification).Error
}

// ============ 查询方法 ============

// GetByEmailAndCode 根据邮箱和验证码查询验证码
func (r *VerificationRepository) GetByEmailAndCode(email, code string) (*model.Verification, error) {
	var v model.Verification
	err := r.db.Where("email = ? AND code = ?", email, code).
		Order("created_at DESC").
		First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// GetLatestByEmail 获取指定邮箱的最新验证码
func (r *VerificationRepository) GetLatestByEmail(email string) (*model.Verification, error) {
	var v model.Verification
	err := r.db.Where("email = ?", email).
		Order("created_at DESC").
		First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// ============ 更新方法 ============

// MarkAsUsed 标记验证码为已使用
func (r *VerificationRepository) MarkAsUsed(id uint) error {
	return r.db.Model(&model.Verification{}).
		Where("id = ?", id).
		Update("used", true).Error
}

// IncrementFailedCount 增加验证失败次数
func (r *VerificationRepository) IncrementFailedCount(id uint) error {
	return r.db.Model(&model.Verification{}).
		Where("id = ?", id).
		UpdateColumn("failed_count", gorm.Expr("failed_count + 1")).Error
}

// ============ 维护方法 ============

// CleanExpired 清理过期的验证码记录（7天前）
func (r *VerificationRepository) CleanExpired() error {
	return r.db.Where("expires_at < ?", time.Now().AddDate(0, 0, -7)).
		Delete(&model.Verification{}).Error
}
