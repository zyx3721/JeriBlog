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
)
