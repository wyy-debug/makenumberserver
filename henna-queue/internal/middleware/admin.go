package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"example.com/henna-queue/internal/util/response"
	"example.com/henna-queue/pkg/jwt"
)

// 上下文键
const (
	ContextAdminID = "admin_id"
	ContextShopID  = "shop_id"
	ContextRole    = "role"
)

// AdminRequired 管理员认证中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 解析token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "无效的认证头")
			c.Abort()
			return
		}

		// 解析管理员声明
		claims, err := jwt.ParseAdminToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "无效的token")
			c.Abort()
			return
		}

		// 设置上下文
		c.Set(ContextAdminID, claims.AdminID)
		c.Set(ContextShopID, claims.ShopID)
		c.Set(ContextRole, claims.Role)

		c.Next()
	}
}

// SuperAdminRequired 超级管理员认证中间件
func SuperAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先经过普通管理员认证
		AdminRequired()(c)

		// 如果已经终止了，直接返回
		if c.IsAborted() {
			return
		}

		// 检查角色
		roleValue, exists := c.Get(ContextRole)
		if !exists {
			response.Forbidden(c, "未找到角色信息")
			c.Abort()
			return
		}

		var role int8
		switch v := roleValue.(type) {
		case int8:
			role = v
		case int:
			role = int8(v)
		default:
			response.ServerError(c, "角色类型错误")
			c.Abort()
			return
		}

		if role != 2 { // 2 是超级管理员
			response.Forbidden(c, "需要超级管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
