/*
项目名称：JeriBlog
文件名称：user.go
创建时间：2026-04-16 15:00:20

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：用户数据访问层
*/

package repository

import (
	"jeri_blog/internal/model"
	"time"

	"gorm.io/gorm"
)

// UserRepository 用户仓储
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ============ 基础CRUD ============

// Create 创建新用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Get 获取用户
func (r *UserRepository) Get(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 物理删除用户
func (r *UserRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&model.User{}, id).Error
}

// ============ 查询方法 ============

// GetByEmail 通过邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// GetGuestByEmail 通过邮箱获取游客用户
func (r *UserRepository) GetGuestByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND role = ?", email, model.RoleGuest).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表（后台管理）
func (r *UserRepository) List(offset, limit int, keyword string, role model.UserRole) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 构建查询
	query := r.db.Unscoped().Model(&model.User{})

	// 关键词搜索（昵称、邮箱、网站）
	if keyword != "" {
		query = query.Where("nickname LIKE ? OR email LIKE ? OR website LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 角色筛选
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// 统计总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	err = query.
		Select("id, email, nickname, avatar, badge, website, is_enabled, role, last_login, created_at, updated_at, deleted_at, has_password, github_id, google_id, qq_id, feishu_open_id").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ExistsByAvatar 检查是否有用户头像引用该文件
func (r *UserRepository) ExistsByAvatar(url string) (bool, error) {
	var count int64
	// 只查询未删除的用户，已删除用户的头像文件可以删除
	err := r.db.Model(&model.User{}).Where("avatar = ?", url).Count(&count).Error
	return count > 0, err
}

// ============ 辅助方法 ============

// UpdateAvatar 更新用户头像
func (r *UserRepository) UpdateAvatar(userID uint, avatarURL string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error
}

// UpdatePassword 更新用户密码
func (r *UserRepository) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}

// ============ Token黑名单 ============

// AddTokenToBlacklist 添加token到黑名单
func (r *UserRepository) AddTokenToBlacklist(tokenHash string, userID uint, expiresAt time.Time) error {
	blacklist := &model.TokenBlacklist{
		TokenHash: tokenHash,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(blacklist).Error
}

// IsTokenBlacklisted 检查token是否在黑名单中
func (r *UserRepository) IsTokenBlacklisted(tokenHash string) bool {
	var count int64
	r.db.Model(&model.TokenBlacklist{}).
		Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).
		Count(&count)
	return count > 0
}

// CleanupExpiredTokens 清理过期的黑名单记录
func (r *UserRepository) CleanupExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&model.TokenBlacklist{}).Error
}

// RevokeAllUserTokens 撤销某用户的所有token
func (r *UserRepository) RevokeAllUserTokens(userID uint) error {
	return r.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Delete(&model.TokenBlacklist{}).Error
}

// ============ OAuth 相关 ============

// GetByOAuthID 通过 OAuth ID 获取用户
func (r *UserRepository) GetByOAuthID(provider, providerID string) (*model.User, error) {
	var user model.User
	var query string

	// 根据提供商选择查询字段
	switch provider {
	case "github":
		query = "github_id = ?"
	case "google":
		query = "google_id = ?"
	case "qq":
		query = "qq_id = ?"
	case "microsoft":
		query = "microsoft_id = ?"
	case "feishu":
		query = "feishu_open_id = ?"
	default:
		return nil, gorm.ErrRecordNotFound
	}

	// 执行查询
	err := r.db.Where(query, providerID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateOAuthBinding 更新 OAuth 绑定
func (r *UserRepository) UpdateOAuthBinding(userID uint, provider, providerID string) error {
	// 根据提供商选择更新字段
	switch provider {
	case "github":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("github_id", providerID).Error
	case "google":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("google_id", providerID).Error
	case "qq":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("qq_id", providerID).Error
	case "microsoft":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("microsoft_id", providerID).Error
	case "feishu":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("feishu_open_id", providerID).Error
	default:
		return gorm.ErrInvalidData
	}
}

// ClearOAuthBinding 清除 OAuth 绑定
func (r *UserRepository) ClearOAuthBinding(userID uint, provider string) error {
	// 根据提供商选择清除字段
	switch provider {
	case "github":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("github_id", "").Error
	case "google":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("google_id", "").Error
	case "qq":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("qq_id", "").Error
	case "microsoft":
		return r.db.Model(&model.User{}).Where("id = ?", userID).Update("microsoft_id", "").Error
	default:
		return gorm.ErrInvalidData
	}
}
