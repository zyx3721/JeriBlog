package dto

import "flec_blog/pkg/utils"

// SubscriberQueryRequest 订阅者查询请求
type SubscriberQueryRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

// SubscriberResponse 订阅者响应
type SubscriberResponse struct {
	ID        uint           `json:"id"`
	Email     string         `json:"email"`
	Active    bool           `json:"active"`
	CreatedAt utils.JSONTime `json:"created_at"`
	UpdatedAt utils.JSONTime `json:"updated_at"`
}
