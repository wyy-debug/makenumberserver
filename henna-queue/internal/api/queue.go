package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/internal/service"
	"example.com/henna-queue/internal/util/response"
)

var queueService = service.NewQueueService()

// GetQueueStatus 获取当前排队状态
func GetQueueStatus(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetString(middleware.ContextUserID)

	// 从上下文或查询参数获取shopID
	shopIDStr := c.Query("shop_id")
	shopID, err := strconv.ParseUint(shopIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的店铺ID")
		return
	}

	status, err := queueService.GetQueueStatus(userID, uint(shopID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, status)
}

// GetQueueNumber 取号排队
func GetQueueNumber(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetString(middleware.ContextUserID)

	// 绑定参数
	var req struct {
		ShopID    uint `json:"shop_id" binding:"required"`
		ServiceID uint `json:"service_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	queue, err := queueService.GetQueueNumber(userID, req.ShopID, req.ServiceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, queue)
}

// CancelQueue 取消排队
func CancelQueue(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetString(middleware.ContextUserID)

	// 绑定参数
	var req struct {
		ShopID uint `json:"shop_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	err := queueService.CancelQueue(userID, req.ShopID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetCurrentQueue 获取当前叫号情况
func GetCurrentQueue(c *gin.Context) {
	// 从参数获取店铺ID
	shopIDStr := c.Query("shop_id")
	shopID, err := strconv.ParseUint(shopIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的店铺ID")
		return
	}

	status, err := queueService.GetCurrentQueue(uint(shopID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, status)
}

// GetAdminQueue 获取管理后台队列
func GetAdminQueue(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	queues, err := queueService.GetAdminQueue(shopID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, queues)
}

// UpdateQueueStatus 更新排队状态
func UpdateQueueStatus(c *gin.Context) {
	// 从路径获取排队ID
	queueIDStr := c.Param("id")
	queueID, err := strconv.ParseUint(queueIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的排队ID")
		return
	}

	// 绑定参数
	var req struct {
		Status int8 `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	queue, err := queueService.UpdateQueueStatus(uint(queueID), shopID, req.Status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, queue)
}

// CallNextNumber 叫号下一位
func CallNextNumber(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	queue, err := queueService.CallNextNumber(shopID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, queue)
}

// ToggleQueuePause 暂停/恢复排队
func ToggleQueuePause(c *gin.Context) {
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)

	shop, err := queueService.ToggleQueuePause(shopID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, shop)
}
