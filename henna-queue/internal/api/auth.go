package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"henna-queue/internal/model"
	"henna-queue/internal/service"
	"henna-queue/internal/util/response"
	"henna-queue/internal/util/wechat"
	"henna-queue/pkg/jwt"
)

var authService = service.NewAuthService()

// Login 微信登录
func Login(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}
	
	// 调用微信API获取openid和session_key
	session, err := wechat.Code2Session(req.Code)
	if err != nil {
		response.ServerError(c, "获取微信用户信息失败")
		return
	}
	
	// 查找或创建用户
	user, err := authService.GetOrCreateUser(session.OpenID, session.UnionID)
	if err != nil {
		response.ServerError(c, "用户创建失败")
		return
	}
	
	// 生成token
	token, err := jwt.GenerateUserToken(user.ID)
	if err != nil {
		response.ServerError(c, "生成token失败")
		return
	}
	
	// 缓存token
	authService.CacheUserToken(user.ID, token)
	
	response.Success(c, gin.H{
		"token": token,
		"user": user,
	})
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}
	
	// 验证管理员
	admin, err := authService.VerifyAdmin(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}
	
	// 生成token
	var shopID uint
	if admin.ShopID != nil {
		shopID = *admin.ShopID
	}
	
	token, err := jwt.GenerateAdminToken(admin.ID, shopID, admin.Role)
	if err != nil {
		response.ServerError(c, "生成token失败")
		return
	}
	
	// 更新最后登录时间
	authService.UpdateAdminLastLogin(admin.ID)
	
	response.Success(c, gin.H{
		"token": token,
		"admin": admin,
	})
}

// AdminLogout 管理员登出
func AdminLogout(c *gin.Context) {
	// 这里可以做一些Token失效处理
	response.Success(c, nil)
}

// Register 管理员注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=4,max=20"`
		Password string `json:"password" binding:"required,min=6"`
		Role     int8   `json:"role"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}
	
	// 检查是否已存在超级管理员
	adminExists, err := adminService.CheckSuperAdminExists()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	
	// 如果已存在超级管理员，且请求中的角色是超级管理员，则拒绝
	if adminExists && req.Role == 2 {
		response.BadRequest(c, "系统已存在超级管理员，不能再注册超级管理员账号")
		return
	}
	
	// 如果没有管理员，强制设置为超级管理员
	if !adminExists {
		req.Role = 2
	}
	
	// 创建管理员账号
	admin, err := adminService.CreateAdmin(req.Username, req.Password, req.Role)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	
	response.Success(c, admin)
}

// CheckAdminExists 检查是否存在管理员账号
func CheckAdminExists(c *gin.Context) {
	exists, err := adminService.CheckSuperAdminExists()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	
	response.Success(c, gin.H{
		"exists": exists,
	})
} 