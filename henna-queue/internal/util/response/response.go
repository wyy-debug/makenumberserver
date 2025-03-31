package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Code:    status,
		Message: message,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未授权"
	}
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "禁止访问"
	}
	Error(c, http.StatusForbidden, message)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	Error(c, http.StatusNotFound, message)
}

// ServerError 服务器错误
func ServerError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	Error(c, http.StatusInternalServerError, message)
} 