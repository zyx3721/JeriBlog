/*
项目名称：JeriBlog
文件名称：setting.go
创建时间：2026-04-16 15:00:50

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：系统设置数据传输对象
*/

package dto

// ============ 配置管理请求 ============

// UpdateSettingGroupRequest 更新配置分组请求
type UpdateSettingGroupRequest struct {
	Settings map[string]string `json:"settings" binding:"required"`
}

// ============ 配置管理响应 ============

// SettingItemResponse 单个配置项响应
type SettingItemResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SettingGroupResponse 配置分组响应
type SettingGroupResponse struct {
	Group    string                `json:"group"`
	Settings []SettingItemResponse `json:"settings"`
}
