// Package middleware RBAC角色权限中间件
package middleware

import (
	"life-rpg/utils"

	"github.com/gin-gonic/gin"
)

// RoleRequired 角色权限校验中间件
// roles: 允许访问的角色Key列表，如 "admin", "user"
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleKey := GetCurrentRoleKey(c)
		if roleKey == "" {
			utils.Forbidden(c, "未获取到角色信息")
			c.Abort()
			return
		}

		// admin 角色拥有所有权限
		if roleKey == "admin" {
			c.Next()
			return
		}

		// 检查当前角色是否在允许的角色列表中
		allowed := false
		for _, role := range roles {
			if roleKey == role {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminRequired 仅管理员可访问
func AdminRequired() gin.HandlerFunc {
	return RoleRequired("admin")
}

// UserRequired 用户及以上角色可访问
func UserRequired() gin.HandlerFunc {
	return RoleRequired("admin", "user")
}
