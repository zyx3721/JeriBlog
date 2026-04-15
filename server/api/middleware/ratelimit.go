package middleware

import (
	"flec_blog/pkg/errcode"
	"flec_blog/pkg/response"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var (
	store = memory.NewStore()
)

// RateLimit 通用限流中间件
// 使用: middleware.RateLimit(limit, minutes, keyType, message)
// 参数:
//
//	limit: 限制次数
//	minutes: 时间窗口（分钟）
//	keyType: 限流key类型 ("ip"=按IP, "user"=按用户ID, "global"=全局)
//	message: 触发限流时的提示信息（可选）
//
// 示例: RateLimit(10, 1, "ip", "登录过于频繁")
func RateLimit(limit int64, minutes int, keyType string, message ...string) gin.HandlerFunc {
	rate := limiter.Rate{Period: time.Duration(minutes) * time.Minute, Limit: limit}
	lim := limiter.New(store, rate)

	defaultMsg := "请求过于频繁，请稍后重试"
	if len(message) > 0 && message[0] != "" {
		defaultMsg = message[0]
	}

	return func(c *gin.Context) {
		var key string

		// 根据keyType生成限流key
		switch keyType {
		case "ip":
			key = c.ClientIP()
		case "user":
			userID, exists := c.Get("user_id")
			if !exists {
				key = c.ClientIP()
			} else {
				key = fmt.Sprintf("user:%v", userID)
			}
		case "global":
			key = "global"
		default:
			key = c.ClientIP()
		}

		context, err := lim.Get(c, key)
		if err != nil {
			c.Next()
			return
		}

		// 设置限流响应头
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

		if context.Reached {
			retryMsg := fmt.Sprintf("%s，请在 %d 秒后重试", defaultMsg, context.Reset)
			response.Error(c, errcode.NewError(429, retryMsg))
			c.Abort()
			return
		}

		c.Next()
	}
}
