/*
项目名称：JeriBlog
文件名称：rssfeed.go
创建时间：2026-04-16 15:00:36

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：RSS 订阅数据模型
*/

package model

import (
	"time"
)

// RssArticle RSS文章模型
type RssArticle struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	FriendID    uint       `gorm:"not null;index" json:"friend_id"`
	Friend      *Friend    `gorm:"foreignKey:FriendID" json:"friend"`
	Title       string     `gorm:"size:500;not null;uniqueIndex:idx_friend_title" json:"title"`
	Link        string     `gorm:"size:1000;not null" json:"link"`
	PublishedAt *time.Time `json:"published_at"`
	IsRead      bool       `gorm:"default:false" json:"is_read"`
	IsDeleted   bool       `gorm:"default:false" json:"is_deleted"`
	UpdateType  string     `gorm:"size:20;default:''" json:"update_type"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
