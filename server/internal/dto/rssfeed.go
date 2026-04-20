/*
项目名称：JeriBlog
文件名称：rssfeed.go
创建时间：2026-04-16 15:00:50

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：RSS 订阅数据传输对象
*/

package dto

import "jeri_blog/pkg/utils"

// ListRssArticleRequest RSS文章列表请求
type ListRssArticleRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword   string `form:"keyword"`
	IsRead    *bool  `form:"is_read"`
	IsDeleted *bool  `form:"is_deleted"`
	FriendID  *uint  `form:"friend_id"`
}

// RssArticleResponse RSS文章响应
type RssArticleResponse struct {
	ID          uint            `json:"id"`
	FriendID    uint            `json:"friend_id"`
	FriendName  string          `json:"friend_name"`
	FriendURL   string          `json:"friend_url"`
	Title       string          `json:"title"`
	Link        string          `json:"link"`
	Description string          `json:"description"`
	PublishedAt *utils.JSONTime `json:"published_at,omitempty"`
	IsRead      bool            `json:"is_read"`
	IsDeleted   bool            `json:"is_deleted"`
	UpdateType  string          `json:"update_type"`
	CreatedAt   *utils.JSONTime `json:"created_at"`
}

// RssArticleListResponse RSS文章列表响应
type RssArticleListResponse struct {
	List        []RssArticleResponse `json:"list"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
	UnreadCount int64                `json:"unread_count"`
}

// MarkArticleReadRequest 标记文章已读请求
type MarkArticleReadRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// MarkAllReadRequest 全部标记已读请求
type MarkAllReadRequest struct{}
