package repository

import (
	"context"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

// SubscriberRepository 订阅者仓储
type SubscriberRepository struct {
	db *gorm.DB
}

// NewSubscriberRepository 创建订阅者仓储
func NewSubscriberRepository(db *gorm.DB) *SubscriberRepository {
	return &SubscriberRepository{db: db}
}

// Create 创建订阅者
func (r *SubscriberRepository) Create(ctx context.Context, subscriber *model.Subscriber) error {
	return r.db.WithContext(ctx).Create(subscriber).Error
}

// GetByEmail 根据邮箱获取订阅者
func (r *SubscriberRepository) GetByEmail(ctx context.Context, email string) (*model.Subscriber, error) {
	var subscriber model.Subscriber
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&subscriber).Error
	return &subscriber, err
}

// GetByToken 根据令牌获取订阅者
func (r *SubscriberRepository) GetByToken(ctx context.Context, token string) (*model.Subscriber, error) {
	var subscriber model.Subscriber
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&subscriber).Error
	return &subscriber, err
}

// Update 更新订阅者
func (r *SubscriberRepository) Update(ctx context.Context, subscriber *model.Subscriber) error {
	return r.db.WithContext(ctx).Save(subscriber).Error
}

// GetActiveSubscribers 获取所有活跃订阅者
func (r *SubscriberRepository) GetActiveSubscribers(ctx context.Context) ([]*model.Subscriber, error) {
	var subscribers []*model.Subscriber
	err := r.db.WithContext(ctx).Where("active = ?", true).Find(&subscribers).Error
	return subscribers, err
}

// List 获取订阅者列表（后台管理）
func (r *SubscriberRepository) List(ctx context.Context, offset, limit int) ([]model.Subscriber, int64, error) {
	var subscribers []model.Subscriber
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Subscriber{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&subscribers).Error
	return subscribers, total, err
}

// Delete 删除订阅者
func (r *SubscriberRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Subscriber{}, id).Error
}
