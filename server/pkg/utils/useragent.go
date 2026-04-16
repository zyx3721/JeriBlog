/*
项目名称：JeriBlog
文件名称：useragent.go
创建时间：2026-04-16 14:59:17

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：User-Agent 解析
*/

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
