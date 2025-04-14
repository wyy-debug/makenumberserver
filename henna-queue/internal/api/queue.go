package api

import (
	"fmt"
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

// GetQueues 获取队列列表
func GetQueues(c *gin.Context) {
	// 获取查询参数
	status := c.Query("status")
	serviceIDStr := c.Query("service_id")
	date := c.Query("date")
	
	// 记录请求参数（调试用）
	requestLog := map[string]string{
		"status":     status,
		"service_id": serviceIDStr,
		"date":       date,
		"shop_id":    c.Query("shop_id"),
	}
	
	c.Set("request_log", requestLog)
	
	// 从上下文获取商店ID
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
			// 在查询参数中不存在shop_id
			userID, exists := c.Get(middleware.ContextUserID)
			if exists && userID != nil {
				// 如果用户已登录，提示必须提供shop_id
				response.BadRequest(c, "缺少店铺ID参数")
				return
			} else {
				// 对于未登录用户，返回相同错误
				response.BadRequest(c, "缺少店铺ID参数")
				return
			}
		}
	}
	
	// 解析服务ID
	var serviceID uint
	if serviceIDStr != "" {
		id, err := strconv.ParseUint(serviceIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "无效的服务ID")
			return
		}
		serviceID = uint(id)
	}
	
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	
	queues, total, err := queueService.GetQueues(shopID, status, serviceID, date, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	
	// 调试输出找到的队列数据
	fmt.Printf("查询结果: 店铺ID=%d, 状态=%s, 服务ID=%d, 日期=%s, 找到%d条数据\n", 
		shopID, status, serviceID, date, len(queues))
	
	for i, q := range queues {
		if i < 5 { // 最多输出5条，避免日志过长
			// 使用类型断言访问map中的字段
			id, _ := q["id"]
			customerName, _ := q["customer_name"]
			status, _ := q["status"]
			createdAt, _ := q["created_at"]
			
			fmt.Printf("队列[%d]: ID=%v, 客户=%v, 状态=%v, 创建时间=%v\n", 
				i, id, customerName, status, createdAt)
		}
	}
	
	response.Success(c, gin.H{
		"queues":    queues,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateQueue 创建队列（批量导入或管理员手动创建）
func CreateQueue(c *gin.Context) {
	// 获取商店ID - 既可以从管理员上下文获取，也可以从请求体获取
	var shopID uint
	
	// 尝试从上下文获取(管理员路由)
	if value, exists := c.Get(middleware.ContextShopID); exists {
		shopID = value.(uint)
	}
	
	var req struct {
		ShopID       uint   `json:"shop_id"`
		CustomerName string `json:"customer_name" binding:"required"`
		Phone        string `json:"phone"`
		ServiceID    uint   `json:"service_id" binding:"required"`
		Note         string `json:"note"`
		Status       int8   `json:"status"`
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
	
	queue, err := queueService.CreateQueueByAdmin(shopID, req.CustomerName, req.Phone, req.ServiceID, req.Note, req.Status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	
	// 调试输出创建的队列信息
	// 使用类型断言访问map中的字段
	id, _ := queue["id"]
	qShopID, _ := queue["shop_id"]
	customerName, _ := queue["customer_name"]
	status, _ := queue["status"]
	createdAt, _ := queue["created_at"]
	
	fmt.Printf("成功创建队列: ID=%v, 店铺ID=%v, 客户=%v, 状态=%v, 创建时间=%v\n", 
		id, qShopID, customerName, status, createdAt)
	
	response.Success(c, queue)
}
