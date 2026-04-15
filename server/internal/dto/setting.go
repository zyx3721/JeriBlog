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
