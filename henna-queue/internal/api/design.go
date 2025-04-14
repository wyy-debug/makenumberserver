package api

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/internal/service"
	"example.com/henna-queue/internal/util/response"
)

var designService = service.NewDesignService()

// GetDesigns 获取图案列表
func GetDesigns(c *gin.Context) {
	// 获取查询参数
	shopIDStr := c.Query("shop_id")
	var shopID uint
	
	// 如果提供了shop_id，则验证其有效性
	if shopIDStr != "" {
		id, err := strconv.ParseUint(shopIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的店铺ID")
			return
		}
		shopID = uint(id)
	} else {
		// 如果没有提供shop_id，设置一个默认值或错误处理
		// 这里我们用默认值1
		shopID = 1
	}

	// 处理其他查询参数
	category := c.Query("category")
	search := c.Query("search")
	priceRange := c.Query("price_range")
	
	// 添加日志以便调试
	log.Printf("GetDesigns 参数: shop_id=%d, category=%s, search=%s, price_range=%s", 
		shopID, category, search, priceRange)
	
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 获取用户ID (如果已登录)
	userID := ""
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		userID = c.GetString(middleware.ContextUserID)
	}

	designs, total, err := designService.GetDesigns(shopID, category, userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"designs":   designs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetDesign 获取单个图案
func GetDesign(c *gin.Context) {
	// 从路径参数获取图案ID
	designIDStr := c.Param("id")
	designID, err := strconv.ParseUint(designIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的图案ID")
		return
	}

	// 获取用户ID (如果已登录)
	userID := ""
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		userID = c.GetString(middleware.ContextUserID)
	}

	design, err := designService.GetDesign(uint(designID), userID)
	if err != nil {
		response.NotFound(c, "图案不存在")
		return
	}

	response.Success(c, design)
}

// ToggleFavorite 收藏/取消收藏图案
func ToggleFavorite(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetString(middleware.ContextUserID)

	// 从路径参数获取图案ID
	designIDStr := c.Param("id")
	designID, err := strconv.ParseUint(designIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的图案ID")
		return
	}

	isLiked, err := designService.ToggleFavorite(userID, uint(designID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"is_liked": isLiked,
	})
}

// GetUserFavorites 获取用户收藏的图案
func GetUserFavorites(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetString(middleware.ContextUserID)

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	designs, total, err := designService.GetUserFavorites(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"designs":   designs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetAdminDesigns 获取管理后台图案列表
func GetAdminDesigns(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	// 获取查询参数
	category := c.Query("category")
	status := c.Query("status")

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	designs, total, err := designService.GetAdminDesigns(shopID, category, status, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"designs":   designs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateDesign 创建图案
func CreateDesign(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	var req struct {
		Title       string `json:"title" binding:"required"`
		Category    string `json:"category" binding:"required"`
		ImageURL    string `json:"image_url" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	design, err := designService.CreateDesign(shopID, &req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, design)
}

// UpdateDesign 更新图案
func UpdateDesign(c *gin.Context) {
	// 从路径参数获取图案ID
	designIDStr := c.Param("id")
	designID, err := strconv.ParseUint(designIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的图案ID")
		return
	}

	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	var req struct {
		Title       string `json:"title"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
		Description string `json:"description"`
		Status      *int8  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	design, err := designService.UpdateDesign(uint(designID), shopID, &req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, design)
}

// DeleteDesign 删除图案
func DeleteDesign(c *gin.Context) {
	// 从路径参数获取图案ID
	designIDStr := c.Param("id")
	designID, err := strconv.ParseUint(designIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的图案ID")
		return
	}

	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	err = designService.DeleteDesign(uint(designID), shopID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
