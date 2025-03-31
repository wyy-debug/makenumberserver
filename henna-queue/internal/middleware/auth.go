package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"example.com/henna-queue/internal/util/response"
	"example.com/henna-queue/pkg/jwt"
)

// ContextUserID 上下文用户ID键
const ContextUserID = "user_id"

// ContextIsSuperAdmin 上下文是否超级管理员键
const ContextIsSuperAdmin = "is_super_admin"

// AuthRequired 用户认证中间件
func AuthRequired() gin.HandlerFunc {
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

		// 解析用户声明
		claims, err := jwt.ParseUserToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "无效的token")
			c.Abort()
			return
		}

		// 设置上下文
		c.Set(ContextUserID, claims.UserID)
		c.Next()
	}
}
