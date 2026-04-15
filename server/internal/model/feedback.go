package model

import "time"

// Feedback 投诉举报模型
type Feedback struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	TicketNo     string     `gorm:"type:varchar(50);unique;not null" json:"ticket_no"` // 工单号
	ReportUrl    string     `gorm:"type:varchar(500);not null" json:"report_url"`      // 举报URL
	ReportType   string     `gorm:"type:varchar(50);not null" json:"report_type"`      // 举报类型
	FormContent  string     `gorm:"type:text;not null" json:"form_content"`            // 表单内容(JSON格式)
	Email        string     `gorm:"type:varchar(255)" json:"email"`                    // 联系邮箱(可选)
	Status       string     `gorm:"type:varchar(20);default:'pending'" json:"status"`  // 状态
	AdminReply   string     `gorm:"type:text" json:"admin_reply"`                      // 管理员回复
	ReplyTime    *time.Time `json:"reply_time"`                                        // 回复时间
	UserAgent    string     `gorm:"type:varchar(500)" json:"user_agent"`               // 用户代理
	IP           string     `gorm:"type:varchar(45)" json:"ip"`                        // IP地址
	FeedbackTime time.Time  `gorm:"not null" json:"feedback_time"`                     // 反馈时间
	CreatedAt    time.Time  `json:"created_at"`                                        // 创建时间
	UpdatedAt    time.Time  `json:"updated_at"`                                        // 更新时间
}
