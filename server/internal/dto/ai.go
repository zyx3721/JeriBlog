/*
项目名称：JeriBlog
文件名称：ai.go
创建时间：2026-04-16 15:00:50

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：AI 数据传输对象
*/

package dto

// ============ AI功能请求 ============

// AISummaryRequest 生成文章摘要请求
type AISummaryRequest struct {
	Content string `json:"content" binding:"required"` // 文章内容
}

// AIAISummaryRequest 生成AI摘要请求
type AIAISummaryRequest struct {
	Content string `json:"content" binding:"required"` // 文章内容
}

// AITitleRequest 生成标题请求
type AITitleRequest struct {
	Content string `json:"content" binding:"required"` // 文章内容
}

// ============ AI功能响应 ============

// AISummaryResponse 摘要生成响应
type AISummaryResponse struct {
	Summary string `json:"summary"` // 生成的摘要
}

// AIAISummaryResponse AI摘要生成响应
type AIAISummaryResponse struct {
	Summary string `json:"summary"` // 生成的AI摘要
}

// AITitleResponse 标题生成响应
type AITitleResponse struct {
	Title string `json:"title"` // 生成的标题
}

// AITestRequest 测试AI配置请求
type AITestRequest struct {
	BaseURL string `json:"base_url" binding:"required"`
	APIKey  string `json:"api_key" binding:"required"`
	Model   string `json:"model" binding:"required"`
}
