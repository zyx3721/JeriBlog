/*
项目名称：JeriBlog
文件名称：openai.go
创建时间：2026-04-16 14:59:17

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：OpenAI API 客户端
*/

package ai

import (
	"flec_blog/config"
	"strings"
)

// OpenAIClient OpenAI 兼容 API 客户端
type OpenAIClient struct {
	BaseURL         string
	APIKey          string
	Model           string
	SummaryPrompt   string
	AISummaryPrompt string
	TitlePrompt     string
}

// OpenAIRequest OpenAI API 请求结构
type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

// OpenAIMessage 消息结构
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse OpenAI API 响应结构
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewOpenAIClient 创建 OpenAI 兼容客户端
func NewOpenAIClient(baseURL, apiKey, model string) *OpenAIClient {
	return &OpenAIClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		Model:   model,
	}
}

// NewOpenAIClientWithConfig creates an OpenAI-compatible client with custom prompts.
func NewOpenAIClientWithConfig(cfg *config.AIConfig) *OpenAIClient {
	if cfg == nil {
		return &OpenAIClient{}
	}

	return &OpenAIClient{
		BaseURL:         strings.TrimRight(cfg.BaseURL, "/"),
		APIKey:          cfg.APIKey,
		Model:           cfg.Model,
		SummaryPrompt:   cfg.SummaryPrompt,
		AISummaryPrompt: cfg.AISummaryPrompt,
		TitlePrompt:     cfg.TitlePrompt,
	}
}
