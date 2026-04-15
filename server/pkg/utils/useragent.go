package utils

import "github.com/mssola/user_agent"

// ParseUserAgent 解析 User-Agent 获取浏览器和操作系统信息
func ParseUserAgent(uaString string) (browser, os string) {
	if uaString == "" {
		return "Unknown", "Unknown"
	}

	ua := user_agent.New(uaString)

	// 浏览器
	name, version := ua.Browser()
	if name != "" {
		browser = name
		if version != "" {
			browser += " " + version
		}
	} else {
		browser = "Unknown"
	}

	// 操作系统
	if osInfo := ua.OS(); osInfo != "" {
		os = osInfo
	} else {
		os = "Unknown"
	}

	return
}
