package dto

import (
	"encoding/json"
	"flec_blog/pkg/utils"
)

// ============ 通知请求 ============

// NotificationListRequest 通知列表请求
type NotificationListRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

// ============ 通知响应 ============

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID       uint   `json:"id"`
	Type     string `json:"type"`      // 类型原始值
	TypeText string `json:"type_text"` // 类型中文文本（前端直接显示）

	// 前端显示字段
	Title   string `json:"title"`   // 通知标题
	Content string `json:"content"` // 通知内容
	Link    string `json:"link"`    // 跳转链接

	Data      json.RawMessage `json:"data" swaggertype:"object"` // 结构化数据（JSON 对象）
	TargetID  *uint           `json:"target_id,omitempty"`       // 目标对象 ID
	IsRead    bool            `json:"is_read"`
	ReadAt    *utils.JSONTime `json:"read_at,omitempty"`
	CreatedAt utils.JSONTime  `json:"created_at"`
	Sender    *string         `json:"sender,omitempty"` // 发送者昵称
}

// NotificationListResponse 通知列表响应
type NotificationListResponse struct {
	List        []NotificationResponse `json:"list"`
	Total       int64                  `json:"total"`
	Page        int                    `json:"page"`
	PageSize    int                    `json:"page_size"`
	UnreadCount int64                  `json:"unread_count"` // 未读数量
}

// MarkAsReadRequest 标记已读请求
type MarkAsReadRequest struct {
	ID uint `uri:"id" binding:"required"`
}
