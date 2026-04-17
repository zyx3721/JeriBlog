/*
项目名称：JeriBlog
文件名称：file.go
创建时间：2026-04-16 15:00:50

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：文件数据传输对象
*/

package dto

import (
	"flec_blog/pkg/upload"
	"flec_blog/pkg/utils"
)

// ============ 通用文件请求 ============

// UploadFileRequest 文件上传请求
type UploadFileRequest struct {
	Type string `form:"type" binding:"required"`
}

// ============ 通用文件响应 ============

// FileUploadForWebResponse 文件上传响应
type FileUploadForWebResponse struct {
	OriginalName string `json:"original_name"`
	FileURL      string `json:"file_url"`
}

// ============ 后台文件管理请求 ============

// ListFilesRequest 文件列表请求
type ListFilesRequest struct {
	Type       string `form:"type"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=1000"`
	Keyword    string `form:"keyword"`     // 搜索关键词（文件名、原始文件名）
	Status     *int   `form:"status"`      // 状态筛选（0=未使用，1=使用中）
	UploadType string `form:"upload_type"` // 上传类型筛选
}

// ============ 后台文件管理响应 ============

// FileResponse 文件信息响应
type FileResponse struct {
	ID           uint           `json:"id"`
	OriginalName string         `json:"original_name"`
	FileName     string         `json:"file_name"`
	FileSize     int64          `json:"file_size"`
	FileType     string         `json:"file_type"`
	FileURL      string         `json:"file_url"`
	UploadType   upload.Type    `json:"upload_type"`
	UserID       *uint          `json:"user_id"`
	Status       int            `json:"status"` // 文件状态：0=未使用 1=使用中
	UploadTime   utils.JSONTime `json:"upload_time"`
}
