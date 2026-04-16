/*
项目名称：JeriBlog
文件名称：email.go
创建时间：2026-04-16 14:59:17

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：邮箱验证工具
*/

package utils

import "strings"

// VirtualEmailSuffix 虚拟邮箱后缀（用于不提供邮箱的第三方登录）
const VirtualEmailSuffix = "@virtual.local"

// IsVirtualEmail 判断是否为虚拟邮箱
func IsVirtualEmail(email string) bool {
	return strings.HasSuffix(strings.ToLower(email), VirtualEmailSuffix)
}
