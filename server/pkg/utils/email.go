package utils

import "strings"

// VirtualEmailSuffix 虚拟邮箱后缀（用于不提供邮箱的第三方登录）
const VirtualEmailSuffix = "@virtual.local"

// IsVirtualEmail 判断是否为虚拟邮箱
func IsVirtualEmail(email string) bool {
	return strings.HasSuffix(strings.ToLower(email), VirtualEmailSuffix)
}
