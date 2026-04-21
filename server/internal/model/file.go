/*
项目名称：JeriBlog
文件名称：file.go
创建时间：2026-04-16 15:00:36

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：文件数据模型
*/

package model

import (
	"time"
)

// File 文件模型
type File struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	FileName       string    `json:"file_name" gorm:"type:varchar(255);not null"`
	OriginalName   string    `json:"original_name" gorm:"type:varchar(255);not null"`
	FilePath       string    `json:"file_path" gorm:"type:varchar(500);not null"`
	FileSize       int64     `json:"file_size"`
	FileType       string    `json:"file_type" gorm:"type:varchar(100)"`
	UploadType     string    `json:"upload_type" gorm:"type:varchar(20);index"`
	StorageType    string    `json:"storage_type" gorm:"type:varchar(20);index"`
	UserID         *uint     `json:"user_id" gorm:"index"`
	FileURL        string    `json:"file_url" gorm:"type:varchar(500)"`
	Status         int       `json:"status" gorm:"default:0;index"` // 0:未使用 1:使用中
	ReferenceCount int       `json:"reference_count" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
