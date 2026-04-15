package notification

import (
	"flec_blog/config"
	"flec_blog/pkg/email"
	"flec_blog/pkg/feishu"
)

// Service 通知服务（多渠道）
type Service struct {
	emailNotifier  *EmailNotifier
	feishuNotifier *FeishuNotifier
}

// NewService 创建通知服务
func NewService(emailClient *email.Client, feishuClient *feishu.Client, cfg *config.Config) *Service {
	var emailNotifier *EmailNotifier
	if emailClient != nil && cfg != nil {
		emailNotifier = NewEmailNotifier(emailClient, cfg)
	}

	var feishuNotifier *FeishuNotifier
	if feishuClient != nil && cfg != nil {
		feishuNotifier = NewFeishuNotifier(feishuClient, cfg)
	}

	return &Service{
		emailNotifier:  emailNotifier,
		feishuNotifier: feishuNotifier,
	}
}

// Data 通知数据
type Data struct {
	Title      string
	Content    string
	Link       string
	Type       string                 // 通知类型
	SenderName string                 // 发件人名称（评论者昵称）
	Data       map[string]interface{} // 详细数据
}

// SendEmail 发送邮件通知
func (s *Service) SendEmail(email string, data Data) error {
	if email != "" && s.emailNotifier != nil {
		return s.emailNotifier.Send(email, data)
	}
	return nil
}

// SendFeishu 发送飞书通知
func (s *Service) SendFeishu(data Data) error {
	if s.feishuNotifier != nil {
		return s.feishuNotifier.Send(data)
	}
	return nil
}
