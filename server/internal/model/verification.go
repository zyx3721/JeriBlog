package model

import "time"

// Verification 邮箱验证码模型
type Verification struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Email       string    `gorm:"type:varchar(100);not null;index:idx_email_code" json:"email"`
	Code        string    `gorm:"type:varchar(6);not null;index:idx_email_code" json:"code"`
	ExpiresAt   time.Time `gorm:"not null;index:idx_expires_at" json:"expires_at"`
	Used        bool      `gorm:"default:false" json:"used"`
	FailedCount int       `gorm:"default:0" json:"failed_count"`
	CreatedAt   time.Time `json:"created_at"`
}
