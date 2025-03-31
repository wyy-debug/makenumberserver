package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"henna-queue/internal/middleware"
	"henna-queue/internal/service"
	"henna-queue/internal/util/response"
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

// GetAdminProfile 获取管理员信息
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

// UpdateAdminProfile 更新管理员信息
func UpdateAdminProfile(c *gin.Context) {
	// 从上下文获取管理员ID
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
	role := c.GetInt8(middleware.ContextRole)
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
	role := c.GetInt8(middleware.ContextRole)
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
	role := c.GetInt8(middleware.ContextRole)
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
