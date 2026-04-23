/*
项目名称：JeriBlog
文件名称：ip.go
创建时间：2026-04-16 14:59:17

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：IP 地址处理
*/

package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const ipAPIURL = "http://ip-api.com/json/%s?lang=zh-CN"

var httpClient = &http.Client{Timeout: 5 * time.Second}

type ipAPIResponse struct {
	Status     string `json:"status"`
	Country    string `json:"country"`
	RegionName string `json:"regionName"`
	City       string `json:"city"`
	Message    string `json:"message"`
}

// InitIPSearcher 初始化 IP 地址搜索器（在线模式无需初始化）
func InitIPSearcher(_ string) error {
	return nil
}

// CloseIPSearcher 关闭 IP 搜索器（在线模式无需关闭）
func CloseIPSearcher() {}

// GetRealIP 获取真实客户端 IP 地址
// 优先级：阿里云 ESA > X-Forwarded-For > X-Real-IP > ClientIP
func GetRealIP(c *gin.Context) string {
	// 1. 优先使用阿里云 ESA 的真实 IP 头部
	if ip := c.GetHeader("Ali-CDN-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.GetHeader("X-Client-IP"); ip != "" {
		return ip
	}

	// 2. 使用 X-Forwarded-For (取第一个非内网 IP)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if parsedIP := net.ParseIP(ip); parsedIP != nil {
				if !parsedIP.IsLoopback() && !parsedIP.IsPrivate() {
					return ip
				}
			}
		}
	}

	// 3. 使用 X-Real-IP
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}

	// 4. 最后使用 Gin 的默认方法
	return c.ClientIP()
}

// GetIPLocation 获取 IP 地址的地理位置
func GetIPLocation(ip string) string {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return "未知"
	}

	if ipAddr.IsLoopback() || ipAddr.IsPrivate() {
		return "本地"
	}

	resp, err := httpClient.Get(fmt.Sprintf(ipAPIURL, ip))
	if err != nil {
		return "未知"
	}
	defer resp.Body.Close()

	var result ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "未知"
	}

	if result.Status != "success" {
		return "未知"
	}

	var parts []string
	if result.Country != "" {
		parts = append(parts, result.Country)
	}
	if result.RegionName != "" && result.RegionName != result.Country {
		parts = append(parts, result.RegionName)
	}
	if result.City != "" && result.City != result.RegionName {
		parts = append(parts, result.City)
	}

	if len(parts) == 0 {
		return "未知"
	}
	return strings.Join(parts, " ")
}
