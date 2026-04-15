package middleware

import (
	"time"

	"flec_blog/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger HTTP 请求日志中间件
// 功能：
//   - 生成唯一请求ID并注入上下文
//   - 记录请求信息（方法、路径、IP、UA、状态码、耗时等）
//
// 使用: router.Use(middleware.Logger())
// 获取请求ID: c.GetString("request_id")
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := uuid.New().String()[:8]
		c.Set("request_id", requestID)

		// 记录开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()

		// 处理请求
		c.Next()

		// 拼接完整路径
		fullPath := path
		if query != "" {
			fullPath = path + "?" + query
		}

		// 记录 HTTP 请求日志
		logger.HTTPRequest(
			requestID,
			method,
			fullPath,
			c.Writer.Status(),
			clientIP,
			time.Since(start),
			c.Writer.Size(),
			userAgent,
			referer,
		)

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			logger.HTTPError(requestID, c.Errors.String())
		}
	}
}
