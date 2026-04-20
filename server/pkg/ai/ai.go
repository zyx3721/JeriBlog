/*
项目名称：JeriBlog
文件名称：ai.go
创建时间：2026-04-16 14:58:57

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：AI 服务接口定义
*/

package ai

import (
	"jeri_blog/config"
	"fmt"
)

// Provider AI服务提供商接口
type Provider interface {
	GenerateSummary(content string) (string, error)
	GenerateAISummary(content string) (string, error)
	GenerateTitle(content string) ([]string, error)
	Test() error
}

// GetProvider 根据配置获取AI服务提供商
func GetProvider(cfg *config.AIConfig) (Provider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("AI配置未设置")
	}

	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("AI BaseURL 未配置")
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("AI API Key 未配置")
	}

	if cfg.Model == "" {
		return nil, fmt.Errorf("AI Model 未配置")
	}

	return NewOpenAIClientWithConfig(cfg), nil
}
