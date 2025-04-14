package api

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/internal/service"
	"example.com/henna-queue/internal/util/response"
)

var serviceService = service.NewServiceService()

// GetSvcById 获取单个服务
func GetSvcById(c *gin.Context) {
	// 获取服务ID
	serviceIDStr := c.Param("id")
	serviceID, err := strconv.ParseUint(serviceIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务ID")
		return
	}

	// 获取服务
	service, err := serviceService.GetService(uint(serviceID))
	if err != nil {
		response.NotFound(c, "服务不存在")
		return
	}

	response.Success(c, service)
}

// GetSvcsByShopId 获取店铺的所有服务
func GetSvcsByShopId(c *gin.Context) {
	// 获取店铺ID
	var shopID uint
	
	// 尝试从路由中间件获取商店ID（管理员）
	if value, exists := c.Get(middleware.ContextShopID); exists {
		shopID = value.(uint)
	} else {
		// 否则尝试从查询参数获取
		shopIDStr := c.Query("shop_id")
		if shopIDStr != "" {
			id, err := strconv.ParseUint(shopIDStr, 10, 32)
			if err != nil {
				response.BadRequest(c, "无效的店铺ID")
				return
			}
			shopID = uint(id)
		} else {
			response.BadRequest(c, "缺少店铺ID参数")
			return
		}
	}

	// 获取服务列表
	services, err := serviceService.GetShopServices(shopID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, services)
}

// CreateSvc 创建服务
func CreateSvc(c *gin.Context) {
	// 记录请求的内容类型
	contentType := c.GetHeader("Content-Type")
	log.Printf("请求内容类型: %s", contentType)
	
	// 获取店铺ID
	var shopID uint
	
	// 尝试从路由中间件获取商店ID（管理员）
	if value, exists := c.Get(middleware.ContextShopID); exists {
		shopID = value.(uint)
		log.Printf("从上下文获取到店铺ID: %d", shopID)
	}
	
	// 绑定表单数据
	name := c.PostForm("name")
	if name == "" {
		response.BadRequest(c, "服务名称不能为空")
		return
	}
	
	durationStr := c.PostForm("duration")
	if durationStr == "" {
		response.BadRequest(c, "服务时长不能为空")
		return
	}
	
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		response.BadRequest(c, "服务时长必须是数字")
		return
	}
	
	description := c.PostForm("description")
	
	// 尝试获取shop_id
	shopIDStr := c.PostForm("shop_id")
	if shopID == 0 && shopIDStr != "" {
		id, err := strconv.ParseUint(shopIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的店铺ID")
			return
		}
		shopID = uint(id)
		log.Printf("从表单中获取到店铺ID: %d", shopID)
	}
	
	shopID = 1

	// 如果仍然没有shopID，则报错
	if shopID == 0 {
		log.Printf("缺少店铺ID参数")
		response.BadRequest(c, "缺少店铺ID参数")
		return
	}
	
	// 获取排序和状态
	sortOrderStr := c.PostForm("sort_order")
	sortOrder := 0
	if sortOrderStr != "" {
		sortOrder, _ = strconv.Atoi(sortOrderStr)
	}
	
	statusStr := c.PostForm("status")
	status := int8(1) // 默认启用
	if statusStr != "" {
		statusVal, err := strconv.Atoi(statusStr)
		if err == nil {
			status = int8(statusVal)
		}
	}
	
	log.Printf("创建服务，店铺ID: %d, 名称: %s, 时长: %d, 状态: %d, 排序: %d", 
		shopID, name, duration, status, sortOrder)
	
	// 创建服务
	service, err := serviceService.CreateService(
		shopID,
		name,
		duration,
		description,
		status,
		sortOrder,
	)
	
	if err != nil {
		log.Printf("创建服务失败: %v", err)
		response.ServerError(c, err.Error())
		return
	}
	
	log.Printf("创建服务成功: ID=%d, 名称=%s", service.ID, service.Name)
	response.Success(c, service)
}

// UpdateSvc 更新服务
func UpdateSvc(c *gin.Context) {
	// 获取服务ID
	serviceIDStr := c.Param("id")
	serviceID, err := strconv.ParseUint(serviceIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务ID")
		return
	}
	
	// 获取店铺ID
	var shopID uint
	
	// 尝试从路由中间件获取商店ID（管理员）
	if value, exists := c.Get(middleware.ContextShopID); exists {
		shopID = value.(uint)
	}
	
	// 绑定请求数据
	var req struct {
		ShopID      uint   `json:"shop_id"`
		Name        string `json:"name" binding:"required"`
		Duration    int    `json:"duration" binding:"required"`
		Description string `json:"description"`
		Status      int8   `json:"status"`
		SortOrder   int    `json:"sort_order"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}
	
	// 如果上下文中没有shopID，则从请求中获取
	if shopID == 0 {
		if req.ShopID == 0 {
			response.BadRequest(c, "缺少店铺ID参数")
			return
		}
		shopID = req.ShopID
	}
	
	// 更新服务
	service, err := serviceService.UpdateService(
		uint(serviceID),
		shopID,
		req.Name,
		req.Duration,
		req.Description,
		req.Status,
		req.SortOrder,
	)
	
	if err != nil {
		if err.Error() == "无权操作该服务" {
			response.Forbidden(c, err.Error())
			return
		}
		response.ServerError(c, err.Error())
		return
	}
	
	response.Success(c, service)
}

// DeleteSvc 删除服务
func DeleteSvc(c *gin.Context) {
	// 获取服务ID
	serviceIDStr := c.Param("id")
	serviceID, err := strconv.ParseUint(serviceIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务ID")
		return
	}
	
	// 获取店铺ID
	var shopID uint
	
	// 尝试从路由中间件获取商店ID（管理员）
	if value, exists := c.Get(middleware.ContextShopID); exists {
		shopID = value.(uint)
	} else {
		// 从查询参数获取
		shopIDStr := c.Query("shop_id")
		if shopIDStr == "" {
			response.BadRequest(c, "缺少店铺ID参数")
			return
		}
		id, err := strconv.ParseUint(shopIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的店铺ID")
			return
		}
		shopID = uint(id)
	}
	
	// 删除服务
	err = serviceService.DeleteService(uint(serviceID), shopID)
	if err != nil {
		if err.Error() == "无权操作该服务" {
			response.Forbidden(c, err.Error())
			return
		}
		response.ServerError(c, err.Error())
		return
	}
	
	response.Success(c, nil)
}

// GetAdminSvcs 管理员获取服务列表
func GetAdminSvcs(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)
	if shopID == 0 {
		response.BadRequest(c, "无法获取店铺ID")
		return
	}
	
	// 获取服务列表
	services, err := serviceService.GetShopServices(shopID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	
	response.Success(c, services)
}

// GetPublicSvcs 公开服务列表接口
func GetPublicSvcs(c *gin.Context) {
	// 获取查询参数
	shopIDStr := c.Query("shop_id")
	categoryIDStr := c.Query("category_id")
	priceRange := c.Query("price_range")
	status := c.Query("status")
	
	log.Printf("查询服务列表，店铺ID: %s, 分类ID: %s, 价格范围: %s, 状态: %s", 
		shopIDStr, categoryIDStr, priceRange, status)
	
	var services []*model.Service
	var err error
	
	// 如果提供了店铺ID，则按店铺ID查询
	if shopIDStr != "" {
		shopID, err := strconv.ParseUint(shopIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的店铺ID")
			return
		}
		
		// 获取指定店铺的服务列表
		services, err = serviceService.GetShopServices(uint(shopID))
		if err != nil {
			response.ServerError(c, err.Error())
			return
		}
	} else {
		// 获取所有可用服务
		services, err = serviceService.GetAllServices()
		if err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}
	
	// 根据其他参数过滤（这里可以实现更复杂的过滤逻辑）
	// TODO: 根据分类、价格和状态进行过滤
	
	response.Success(c, services)
} 