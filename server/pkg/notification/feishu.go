package notification

import (
	"context"
	"encoding/json"
	"flec_blog/config"
	"flec_blog/pkg/feishu"
)

// FeishuNotifier 飞书通知器
type FeishuNotifier struct {
	client *feishu.Client
	config *config.Config
}

// NewFeishuNotifier 创建飞书通知器
func NewFeishuNotifier(client *feishu.Client, cfg *config.Config) *FeishuNotifier {
	return &FeishuNotifier{client: client, config: cfg}
}

// Send 发送飞书通知
func (f *FeishuNotifier) Send(data Data) error {
	if f.client == nil || !f.client.IsEnabled() {
		return nil
	}

	cardContent, err := f.buildCard(data)
	if err != nil {
		return err
	}

	return f.client.SendMessage(context.Background(), cardContent)
}

// buildCard 构建卡片内容
func (f *FeishuNotifier) buildCard(data Data) (string, error) {
	switch data.Type {
	case "comment_new":
		return f.buildCommentCard(data)
	case "friend_apply":
		return f.buildFriendApplyCard(data)
	case "feedback_new":
		return f.buildFeedbackCard(data)
	case "rss_feed_daily":
		return f.buildRssFeedCard(data)
	default:
		return "", nil
	}
}

// buildCommentCard 构建评论通知卡片
func (f *FeishuNotifier) buildCommentCard(data Data) (string, error) {
	return feishu.BuildCommentCard(
		uint(getFloat64(data.Data, "comment_id")),
		getString(data.Data, "page_title"),
		getString(data.Data, "page_link"),
		getString(data.Data, "sender_name"),
		getString(data.Data, "comment_content"),
		f.config.Basic.BlogURL,
		f.config.Basic.AdminURL,
	)
}

// buildFriendApplyCard 构建友链申请卡片
func (f *FeishuNotifier) buildFriendApplyCard(data Data) (string, error) {
	return feishu.BuildFriendApplyCard(
		uint(getFloat64(data.Data, "friend_id")),
		getString(data.Data, "site_name"),
		getString(data.Data, "site_url"),
		getString(data.Data, "site_description"),
		getString(data.Data, "site_logo"),
		getString(data.Data, "site_screenshot"),
		f.config.Basic.AdminURL,
	)
}

// buildFeedbackCard 构建反馈通知卡片
func (f *FeishuNotifier) buildFeedbackCard(data Data) (string, error) {
	return feishu.BuildFeedbackCard(
		uint(getFloat64(data.Data, "feedback_id")),
		getString(data.Data, "ticket_no"),
		f.mapReportType(getString(data.Data, "report_type")),
		f.config.Basic.AdminURL,
	)
}

// buildRssFeedCard 构建RSS订阅推送卡片
func (f *FeishuNotifier) buildRssFeedCard(data Data) (string, error) {
	var articles []feishu.RssArticleItem
	if articlesData, ok := data.Data["articles"].([]feishu.RssArticleItem); ok {
		articles = articlesData
	}

	return feishu.BuildRssFeedCard(articles, f.config.Basic.AdminURL)
}

// mapReportType 映射举报类型
func (f *FeishuNotifier) mapReportType(reportType string) string {
	typeMap := map[string]string{
		"copyright":     "版权侵权内容投诉",
		"inappropriate": "不当内容举报投诉",
		"summary":       "文章摘要问题反馈",
		"suggestion":    "功能建议优化反馈",
	}
	if mapped, ok := typeMap[reportType]; ok {
		return mapped
	}
	return reportType
}

// getString 从 map 中获取字符串值
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// getFloat64 从 map 中获取 float64 值
func getFloat64(m map[string]interface{}, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	// 处理 json.Number 类型
	if v, ok := m[key].(json.Number); ok {
		if f, err := v.Float64(); err == nil {
			return f
		}
	}
	// 处理 int 类型
	if v, ok := m[key].(int); ok {
		return float64(v)
	}
	// 处理 int64 类型
	if v, ok := m[key].(int64); ok {
		return float64(v)
	}
	// 处理 uint 类型
	if v, ok := m[key].(uint); ok {
		return float64(v)
	}
	// 处理 uint64 类型
	if v, ok := m[key].(uint64); ok {
		return float64(v)
	}
	return 0
}
