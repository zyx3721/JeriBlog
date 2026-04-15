package middleware

import (
	"flec_blog/internal/model"
	"flec_blog/pkg/errcode"
	"flec_blog/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRoles 角色权限中间件
// 检查用户是否具有指定角色之一，超级管理员拥有所有权限
// 使用: router.Use(middleware.RequireRoles(model.RoleAdmin, model.RoleUser))
// 依赖: 需在 Auth 中间件之后使用
func RequireRoles(roles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户
		user, exists := c.Get("user")
		if !exists {
			response.Error(c, errcode.NewError(http.StatusUnauthorized, "未找到用户信息"))
			c.Abort()
			return
		}

		// 将接口类型转换为User类型
		currentUser, ok := user.(*model.User)
		if !ok {
			response.Error(c, errcode.NewError(http.StatusInternalServerError, "用户信息类型错误"))
			c.Abort()
			return
		}

		// 如果用户是超级管理员，直接通过
		if currentUser.Role == model.RoleSuperAdmin {
			c.Next()
			return
		}

		// 检查用户是否具有所需的任一角色
		hasRole := false
		for _, role := range roles {
			if currentUser.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Error(c, errcode.NewError(http.StatusForbidden, "权限不足，需要管理员权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// IsSuperAdmin 超级管理员权限
// 使用: router.Use(middleware.IsSuperAdmin())
func IsSuperAdmin() gin.HandlerFunc {
	return RequireRoles(model.RoleSuperAdmin)
}

// IsAdminOrAbove 管理员及以上权限
// 使用: router.Use(middleware.IsAdminOrAbove())
func IsAdminOrAbove() gin.HandlerFunc {
	return RequireRoles(model.RoleSuperAdmin, model.RoleAdmin)
}

// IsUserOrAbove 普通用户及以上权限
// 使用: router.Use(middleware.IsUserOrAbove())
func IsUserOrAbove() gin.HandlerFunc {
	return RequireRoles(model.RoleSuperAdmin, model.RoleAdmin, model.RoleUser)
}
