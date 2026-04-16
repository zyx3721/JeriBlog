/*
项目名称：JeriBlog
文件名称：stats.go
创建时间：2026-04-16 15:00:36

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：统计数据模型
*/

package model

import "time"

// Visit 访问记录模型
type Visit struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VisitorID string    `gorm:"size:64;not null" json:"visitor_id"`         // 访客唯一标识
	IP        string    `gorm:"size:45;not null" json:"ip"`                 // 访客IP
	PageURL   string    `gorm:"size:500" json:"page_url"`                   // 访问页面URL
	ArticleID *uint     `json:"article_id"`                                 // 文章ID（可选）
	UserAgent string    `gorm:"type:text" json:"user_agent"`                // 浏览器UA
	Location  string    `gorm:"size:100" json:"location"`                   // 地理位置
	Browser   string    `gorm:"size:50" json:"browser"`                     // 浏览器
	OS        string    `gorm:"size:50" json:"os"`                          // 操作系统
	Referer   string    `gorm:"size:500" json:"referer"`                    // 来源页面
	VisitDate time.Time `gorm:"type:date;not null;index" json:"visit_date"` // 访问日期
	CreatedAt time.Time `json:"created_at"`
}
