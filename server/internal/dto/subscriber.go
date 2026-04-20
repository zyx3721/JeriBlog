/*
项目名称：JeriBlog
文件名称：subscriber.go
创建时间：2026-04-16 15:00:50

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：订阅者数据传输对象
*/

package dto

import "jeri_blog/pkg/utils"

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
