/*
项目名称：JeriBlog
文件名称：subscriber.go
创建时间：2026-04-16 15:00:36

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：订阅者数据模型
*/

package model

import "time"

// Subscriber 订阅者模型
type Subscriber struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Active    bool      `gorm:"default:true;index" json:"active"` // true=已订阅, false=已退订
	Token     string    `gorm:"uniqueIndex;size:64" json:"-"`     // 退订令牌
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
