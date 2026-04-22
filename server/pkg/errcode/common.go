/*
项目名称：JeriBlog
文件名称：common.go
创建时间：2026-04-16 14:59:17

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：通用错误码定义
*/

package errcode

// 通用错误码 (10000-10999)
var (
	ServerError   = NewError(10000, "服务内部错误")
	InvalidParams = NewError(10001, "入参错误")
	NotFound      = NewError(10002, "资源不存在")
	Unauthorized  = NewError(10003, "未授权访问")
)

// 用户/认证错误码 (20000-20999)
var (
	TokenError = NewError(20004, "Token错误或已过期")
)

// 文件上传错误码 (40000-40999)
var (
	FileUploadError  = NewError(40001, "文件上传失败")
	FileNotFound     = NewError(40003, "文件不存在")
	FileProcessError = NewError(40006, "文件处理失败")
	FileInUseError   = NewError(40007, "文件正在被使用,无法删除")
)
