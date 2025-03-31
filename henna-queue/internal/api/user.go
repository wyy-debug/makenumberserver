package api

import (
	"github.com/gin-gonic/gin"

	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/internal/util/response"
)

// GetUserProfile 获取用户个人信息
func GetUserProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetString(middleware.ContextUserID)

	user, err := authService.GetUser(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// UpdateUserProfile 更新用户个人信息
func UpdateUserProfile(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetString(middleware.ContextUserID)

	// 获取当前用户
	user, err := authService.GetUser(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	var req struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	err = authService.UpdateUser(user)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateUserPhone 更新用户手机号
func UpdateUserPhone(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetString(middleware.ContextUserID)

	// 获取当前用户
	user, err := authService.GetUser(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 验证短信验证码
	// TODO: 实现短信验证码验证
	// 此处为简化处理，直接更新手机号

	// 更新手机号
	user.Phone = req.Phone

	err = authService.UpdateUser(user)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, user)
}
