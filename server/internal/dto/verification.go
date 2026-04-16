/*
项目名称：JeriBlog
文件名称：verification.go
创建时间：2026-04-16 15:00:50

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：验证码数据传输对象
*/

package dto

// ============ 通用验证请求 ============

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required,len=6"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}
