/*
项目名称：JeriBlog
文件名称：setting.go
创建时间：2026-04-16 15:00:36

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：系统设置数据模型
*/

package model

import "time"

// Setting KV 配置表
type Setting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"` // 配置键，如 email.smtp_host
	Value     string    `gorm:"type:text" json:"value"`                   // 配置值，统一存字符串
	Group     string    `gorm:"index;size:50" json:"group"`               // 配置分组，如 email, upload
	IsPublic  bool      `gorm:"default:true" json:"is_public"`            // 是否公开可见（前台是否可见）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SettingGroup 配置分组常量
const (
	SettingGroupBasic        = "basic"
	SettingGroupBlog         = "blog"
	SettingGroupNotification = "notification"
	SettingGroupUpload       = "upload"
	SettingGroupAI           = "ai"
	SettingGroupOAuth        = "oauth"
	SettingGroupWeChat       = "wechat"
)
