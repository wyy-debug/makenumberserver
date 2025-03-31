package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/internal/service"
	"example.com/henna-queue/internal/util/response"
)

var categoryService = service.NewCategoryService()

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	// 获取查询参数
	shopIDStr := c.Query("shop_id")
	if shopIDStr != "" {
		parsedID, err := strconv.ParseUint(shopIDStr, 10, 32)
		if err == nil {
			shopID = uint(parsedID)
		}
	}

	categories, err := categoryService.GetCategories(shopID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, categories)
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	var req struct {
		Name      string `json:"name" binding:"required"`
		Code      string `json:"code" binding:"required"`
		SortOrder int    `json:"sort_order"`
		ShopID    uint   `json:"shop_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 如果请求中指定了店铺ID，且当前用户为超级管理员，则使用请求中的店铺ID
	if req.ShopID > 0 && c.GetBool(middleware.ContextIsSuperAdmin) {
		shopID = req.ShopID
	}

	category, err := categoryService.CreateCategory(shopID, req.Name, req.Code, req.SortOrder)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, category)
}

// UpdateCategory 更新分类
func UpdateCategory(c *gin.Context) {
	// 从路径参数获取分类ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	var req struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	category, err := categoryService.UpdateCategory(uint(categoryID), shopID, req.Name, req.SortOrder)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, category)
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	// 从路径参数获取分类ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	err = categoryService.DeleteCategory(uint(categoryID), shopID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
