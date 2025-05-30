package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/internal/service"
	"example.com/henna-queue/internal/util/response"
)

var adminService = service.NewAdminService()

// GetStatistics 获取统计数据
func GetStatistics(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	// 获取日期范围
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := adminService.GetStatistics(shopID, days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetAdminProfile 获取管理员个人信息
func GetAdminProfile(c *gin.Context) {
	// 从上下文获取管理员ID
	adminID := c.GetUint(middleware.ContextAdminID)

	admin, err := adminService.GetAdmin(adminID)
	if err != nil {
		response.NotFound(c, "管理员不存在")
		return
	}

	response.Success(c, admin)
}

// UpdateAdminProfile 更新管理员个人信息
func UpdateAdminProfile(c *gin.Context) {
	adminID := c.GetUint(middleware.ContextAdminID)

	var req struct {
		Username    string `json:"username"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	admin, err := adminService.UpdateAdmin(adminID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, admin)
}

// CreateAdmin 创建管理员(超级管理员使用)
func CreateAdmin(c *gin.Context) {
	// 检查超级管理员权限
	roleValue, exists := c.Get(middleware.ContextRole)
	if !exists {
		response.Forbidden(c, "未找到角色信息")
		return
	}

	role, ok := roleValue.(int8)
	if !ok {
		// 尝试将 interface{} 转换为 int8
		if roleInt, ok := roleValue.(int); ok {
			role = int8(roleInt)
		} else {
			response.ServerError(c, "角色类型错误")
			return
		}
	}

	if role != 2 {
		response.Forbidden(c, "需要超级管理员权限")
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		ShopID   *uint  `json:"shop_id"`
		Role     int8   `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	admin, err := adminService.CreateAdmin(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, admin)
}

// GetAdmins 获取管理员列表(超级管理员使用)
func GetAdmins(c *gin.Context) {
	// 检查超级管理员权限
	roleValue, exists := c.Get(middleware.ContextRole)
	if !exists {
		response.Forbidden(c, "未找到角色信息")
		return
	}

	role, ok := roleValue.(int8)
	if !ok {
		// 尝试将 interface{} 转换为 int8
		if roleInt, ok := roleValue.(int); ok {
			role = int8(roleInt)
		} else {
			response.ServerError(c, "角色类型错误")
			return
		}
	}

	if role != 2 {
		response.Forbidden(c, "需要超级管理员权限")
		return
	}

	admins, err := adminService.GetAllAdmins()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, admins)
}

// DeleteAdmin 删除管理员(超级管理员使用)
func DeleteAdmin(c *gin.Context) {
	// 检查超级管理员权限
	roleValue, exists := c.Get(middleware.ContextRole)
	if !exists {
		response.Forbidden(c, "未找到角色信息")
		return
	}

	role, ok := roleValue.(int8)
	if !ok {
		// 尝试将 interface{} 转换为 int8
		if roleInt, ok := roleValue.(int); ok {
			role = int8(roleInt)
		} else {
			response.ServerError(c, "角色类型错误")
			return
		}
	}

	if role != 2 {
		response.Forbidden(c, "需要超级管理员权限")
		return
	}

	// 从路径参数获取管理员ID
	adminIDStr := c.Param("id")
	adminID, err := strconv.ParseUint(adminIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的管理员ID")
		return
	}

	err = adminService.DeleteAdmin(uint(adminID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// CheckAdminExists 检查管理员是否存在
func CheckAdminExists(c *gin.Context) {
	exists, err := adminService.CheckAdminExists()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"exists": exists})
}

// RegisterAdmin 注册管理员
func RegisterAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		ShopName string `json:"shop_name" binding:"required"`
		ShopDesc string `json:"shop_desc"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 创建管理员账号，无论是否已存在管理员
	admin, err := adminService.RegisterAdmin(req.Username, req.Password, req.ShopName, req.ShopDesc, req.Phone)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, admin)
}
