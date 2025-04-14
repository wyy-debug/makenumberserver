package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/internal/service"
	"example.com/henna-queue/internal/util/response"
)

var shopService = service.NewShopService()

// GetShop 获取店铺信息
func GetShop(c *gin.Context) {
	// 从路径参数获取店铺ID
	shopIDStr := c.Param("id")
	shopID, err := strconv.ParseUint(shopIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的店铺ID")
		return
	}

	shop, err := shopService.GetShop(uint(shopID))
	if err != nil {
		response.NotFound(c, "店铺不存在")
		return
	}

	response.Success(c, shop)
}

// GetShopServices 获取店铺服务
func GetShopServices(c *gin.Context) {
	// 从路径参数获取店铺ID
	shopIDStr := c.Param("id")
	shopID, err := strconv.ParseUint(shopIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的店铺ID")
		return
	}

	services, err := shopService.GetShopServices(uint(shopID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, services)
}

// GetAdminShop 获取管理员店铺信息
func GetAdminShop(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	shop, err := shopService.GetShop(shopID)
	if err != nil {
		response.NotFound(c, "店铺不存在")
		return
	}

	response.Success(c, shop)
}

// UpdateShop 更新店铺信息
func UpdateShop(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	var req struct {
		Name          string  `json:"name"`
		Address       string  `json:"address"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		Phone         string  `json:"phone"`
		BusinessHours string  `json:"business_hours"`
		Description   string  `json:"description"`
		CoverImage    string  `json:"cover_image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	shop, err := shopService.UpdateShop(shopID, &req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, shop)
}

// GetAdminServices 获取店铺全部服务
func GetAdminServices(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	services, err := shopService.GetAllServices(shopID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, services)
}

// CreateService 创建服务
func CreateService(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	var req struct {
		Name        string `json:"name" binding:"required"`
		Duration    int    `json:"duration" binding:"required"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	service, err := shopService.CreateService(shopID, &req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, service)
}

// UpdateService 更新服务
func UpdateService(c *gin.Context) {
	// 从路径参数获取服务ID
	serviceIDStr := c.Param("id")
	serviceID, err := strconv.ParseUint(serviceIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务ID")
		return
	}

	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	var req struct {
		Name        string `json:"name"`
		Duration    int    `json:"duration"`
		Description string `json:"description"`
		Status      *int8  `json:"status"`
		SortOrder   int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	service, err := shopService.UpdateService(uint(serviceID), shopID, &req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, service)
}

// DeleteService 删除服务
func DeleteService(c *gin.Context) {
	// 从路径参数获取服务ID
	serviceIDStr := c.Param("id")
	serviceID, err := strconv.ParseUint(serviceIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务ID")
		return
	}

	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	err = shopService.DeleteService(uint(serviceID), shopID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetShopStats 获取店铺统计数据
func GetShopStats(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	// 获取日期范围参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	stats, err := shopService.GetShopStats(shopID, startDate, endDate)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetPublicServices 获取所有公开的服务列表
func GetPublicServices(c *gin.Context) {
	// 可选参数：店铺ID
	shopIDStr := c.Query("shop_id")
	
	var shopID uint
	if shopIDStr != "" {
		id, err := strconv.ParseUint(shopIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的店铺ID")
			return
		}
		shopID = uint(id)
	}

	services, err := shopService.GetPublicServices(shopID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, services)
}
