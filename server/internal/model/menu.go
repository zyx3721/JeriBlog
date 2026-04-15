package model

import "time"

// Menu 菜单项
type Menu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Type      string    `gorm:"size:50;not null;index" json:"type"` // 菜单类型: aggregate(网站聚合), navigation(顶部导航), footer(页脚菜单)
	ParentID  *uint     `gorm:"index" json:"parent_id"`             // 父菜单ID
	Title     string    `gorm:"size:100;not null" json:"title"`     // 菜单标题
	URL       string    `gorm:"size:500" json:"url"`                // 链接地址，主菜单可为空，子菜单需填写
	Icon      string    `gorm:"size:500" json:"icon"`               // 图标：remixicon类名(如ri-home-line) 或 图片URL
	Sort      int       `gorm:"default:5;index" json:"sort"`        // 排序，数字越小越靠前，范围1-10
	IsEnabled bool      `gorm:"default:true" json:"is_enabled"`     // 是否启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联关系
	Parent   *Menu  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`   // 父菜单
	Children []Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子菜单列表
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}
