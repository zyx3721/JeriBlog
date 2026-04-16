/*
项目名称：JeriBlog
文件名称：tag.go
创建时间：2026-04-16 15:00:36

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：标签数据模型
*/

package model

import (
	"time"
)

// Tag 标签模型
type Tag struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name"`
	Slug        string    `gorm:"uniqueIndex;size:50" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Count       int       `gorm:"default:0" json:"count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
