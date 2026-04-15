package feishu

import (
	"encoding/json"
	"fmt"
)

// Card 飞书卡片 2.0 结构
type Card struct {
	Schema string      `json:"schema"`
	Config *CardConfig `json:"config,omitempty"`
	Header *CardHeader `json:"header,omitempty"`
	Body   *CardBody   `json:"body,omitempty"`
}

// CardConfig 卡片配置
type CardConfig struct {
	UpdateMulti bool       `json:"update_multi,omitempty"`
	Style       *CardStyle `json:"style,omitempty"`
}

// CardStyle 卡片样式
type CardStyle struct {
	TextSize *TextSizeStyle `json:"text_size,omitempty"`
}

// TextSizeStyle 文本大小样式
type TextSizeStyle struct {
	NormalV2 *SizeConfig `json:"normal_v2,omitempty"`
}

// SizeConfig 大小配置
type SizeConfig struct {
	Default string `json:"default,omitempty"`
	PC      string `json:"pc,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
}

// CardHeader 卡片头部
type CardHeader struct {
	Template string      `json:"template"`
	Title    *TextObject `json:"title"`
	Subtitle *TextObject `json:"subtitle,omitempty"`
	Padding  string      `json:"padding,omitempty"`
}

// CardBody 卡片主体
type CardBody struct {
	Direction string        `json:"direction,omitempty"`
	Padding   string        `json:"padding,omitempty"`
	Elements  []interface{} `json:"elements"`
}

// TextObject 文本对象
type TextObject struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

// MarkdownElement markdown 元素
type MarkdownElement struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// InputElement 输入框元素
type InputElement struct {
	Tag         string                 `json:"tag"`
	Name        string                 `json:"name,omitempty"`
	Placeholder *TextObject            `json:"placeholder"`
	Width       string                 `json:"width,omitempty"`
	Value       map[string]interface{} `json:"value,omitempty"`
}

// ButtonElement 按钮元素
type ButtonElement struct {
	Tag      string                 `json:"tag"`
	Text     *TextObject            `json:"text"`
	Type     string                 `json:"type,omitempty"`
	Width    string                 `json:"width,omitempty"`
	Size     string                 `json:"size,omitempty"`
	MultiURL *URLElement            `json:"multi_url,omitempty"`
	Value    map[string]interface{} `json:"value,omitempty"`
}

// URLElement URL跳转元素
type URLElement struct {
	URL   string `json:"url,omitempty"`
	PCURL string `json:"pc_url,omitempty"`
}

// 默认卡片配置
func defaultCardConfig() *CardConfig {
	return &CardConfig{
		UpdateMulti: true,
		Style: &CardStyle{
			TextSize: &TextSizeStyle{
				NormalV2: &SizeConfig{
					Default: "normal",
					PC:      "normal",
					Mobile:  "heading",
				},
			},
		},
	}
}

// 默认卡片头部
func defaultCardHeader(title, template string) *CardHeader {
	return &CardHeader{
		Template: template,
		Title: &TextObject{
			Content: title,
			Tag:     "plain_text",
		},
		Subtitle: &TextObject{
			Content: "",
			Tag:     "plain_text",
		},
		Padding: "12px 12px 12px 12px",
	}
}

// 默认卡片主体
func defaultCardBody(elements []interface{}) *CardBody {
	return &CardBody{
		Direction: "vertical",
		Padding:   "12px 12px 12px 12px",
		Elements:  elements,
	}
}

// 创建 markdown 元素
func newMarkdownElement(content string) *MarkdownElement {
	return &MarkdownElement{Tag: "markdown", Content: content}
}

// 创建带交互动作的 input 元素
func newInputElementWithAction(placeholder string, value map[string]interface{}) *InputElement {
	return &InputElement{
		Tag:         "input",
		Name:        "reply_content",
		Placeholder: &TextObject{Content: placeholder, Tag: "plain_text"},
		Width:       "fill",
		Value:       value,
	}
}

// 创建带跳转链接的 button 元素
func newButtonElementWithURL(content, btnType, url string) *ButtonElement {
	return &ButtonElement{
		Tag:      "button",
		Text:     &TextObject{Content: content, Tag: "plain_text"},
		Type:     btnType,
		Width:    "fill",
		Size:     "medium",
		MultiURL: &URLElement{URL: url, PCURL: url},
	}
}

// 创建带交互回调的 button 元素
func newButtonElementWithAction(content, btnType string, value map[string]interface{}) *ButtonElement {
	return &ButtonElement{
		Tag:   "button",
		Text:  &TextObject{Content: content, Tag: "plain_text"},
		Type:  btnType,
		Width: "fill",
		Size:  "medium",
		Value: value,
	}
}

// buildCard 构建卡片（公共逻辑）
func buildCard(title, template string, elements []interface{}) (string, error) {
	card := &Card{
		Schema: "2.0",
		Config: defaultCardConfig(),
		Header: defaultCardHeader(title, template),
		Body:   defaultCardBody(elements),
	}

	cardJSON, err := json.Marshal(card)
	if err != nil {
		return "", fmt.Errorf("序列化卡片失败: %w", err)
	}
	return string(cardJSON), nil
}

// BuildCommentCard 构建评论通知卡片
func BuildCommentCard(commentID uint, pageTitle, pageLink, senderName, commentContent, blogURL, adminURL string) (string, error) {
	content := fmt.Sprintf("来源：%s\n评论者：%s\n评论内容：%s", pageTitle, senderName, commentContent)
	elements := []interface{}{
		newMarkdownElement(content),
		newInputElementWithAction("回复评论", map[string]interface{}{"action": "reply_comment", "comment_id": commentID}),
		newButtonElementWithURL("查看", "primary", fmt.Sprintf("%s%s", blogURL, pageLink)),
		newButtonElementWithURL("详情", "default", fmt.Sprintf("%s/comments", adminURL)),
	}
	return buildCard("📬 收到了新的评论通知", "blue", elements)
}

// BuildFriendApplyCard 构建友链申请卡片
func BuildFriendApplyCard(friendID uint, siteName, siteURL, siteDescription, siteLogo, siteScreenshot, adminURL string) (string, error) {
	content := fmt.Sprintf("名称：%s\n地址：[%s](%s)\n描述：%s\n头像：[%s](%s)\n截图：[%s](%s)",
		siteName, siteURL, siteURL, siteDescription, siteLogo, siteLogo, siteScreenshot, siteScreenshot)
	elements := []interface{}{
		newMarkdownElement(content),
		newButtonElementWithAction("通过", "primary_filled", map[string]interface{}{"action": "approve_friend", "friend_id": friendID}),
		newButtonElementWithURL("详情", "primary", fmt.Sprintf("%s/friends", adminURL)),
	}
	return buildCard("🔗 收到了新的友链申请", "green", elements)
}

// BuildFeedbackCard 构建反馈通知卡片
func BuildFeedbackCard(feedbackID uint, ticketNo, reportType, adminURL string) (string, error) {
	content := fmt.Sprintf("工单号：%s\n举报类型：%s", ticketNo, reportType)
	elements := []interface{}{
		newMarkdownElement(content),
		newButtonElementWithURL("详情", "primary", fmt.Sprintf("%s/feedback/%d", adminURL, feedbackID)),
	}
	return buildCard("📝 收到了新的反馈投诉", "orange", elements)
}

// BuildRssFeedCard 构建RSS订阅推送卡片
func BuildRssFeedCard(articles []RssArticleItem, adminURL string) (string, error) {
	var content string
	for i, article := range articles {
		content += fmt.Sprintf("%d. 《[%s](%s)》 %s\n\n", i+1, article.Title, article.Link, article.FriendName)
	}

	elements := []interface{}{
		newMarkdownElement(content),
		newButtonElementWithAction("全部已读", "primary_filled", map[string]interface{}{"action": "rss_mark_all_read"}),
		newButtonElementWithURL("查看详情", "primary", fmt.Sprintf("%s/rssfeed", adminURL)),
	}

	return buildCard("📰 RSS订阅日报", "purple", elements)
}

// RssArticleItem RSS文章项
type RssArticleItem struct {
	Title      string
	Link       string
	FriendName string
}
