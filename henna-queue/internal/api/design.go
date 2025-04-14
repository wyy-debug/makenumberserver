package api

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
	// 记录请求内容类型
	contentType := c.GetHeader("Content-Type")
	log.Printf("创建图案请求内容类型: %s", contentType)
	
	// 获取店铺ID
	var shopID uint
	
	// 尝试从路由中间件获取商店ID（管理员）
	if value, exists := c.Get(middleware.ContextShopID); exists {
		shopID = value.(uint)
		log.Printf("从上下文获取到店铺ID: %d", shopID)
	}
	
	// 处理表单数据
	title := c.PostForm("title")
	if title == "" {
		response.BadRequest(c, "图案标题不能为空")
		return
	}
	
	category := c.PostForm("category")
	if category == "" {
		response.BadRequest(c, "图案分类不能为空")
		return
	}
	
	description := c.PostForm("description")
	
	// 处理文件上传
	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("获取上传文件失败: %v", err)
		response.BadRequest(c, "请上传图案图片")
		return
	}
	
	// 创建上传目录
	uploadDir := "./static/uploads/designs/"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("创建目录失败: %v", err)
		response.ServerError(c, "上传目录创建失败")
		return
	}
	
	// 生成唯一文件名
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], fileExt)
	filePath := filepath.Join(uploadDir, fileName)
	
	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("保存文件失败: %v", err)
		response.ServerError(c, "文件保存失败")
		return
	}
	
	// 生成可访问的URL
	imageURL := "/static/uploads/designs/" + fileName
	
	// 如果是公共API调用，需要从请求中获取shopID
	if shopID == 0 {
		shopIDStr := c.PostForm("shop_id")
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
	
	log.Printf("创建图案，店铺ID: %d, 标题: %s, 分类: %s, 图片URL: %s", 
		shopID, title, category, imageURL)
	
	// 创建设计
	reqData := struct {
		Title       string `json:"title" binding:"required"`
		Category    string `json:"category" binding:"required"`
		ImageURL    string `json:"image_url" binding:"required"`
		Description string `json:"description"`
	}{
		Title:       title,
		Category:    category,
		ImageURL:    imageURL,
		Description: description,
	}
	
	design, err := designService.CreateDesign(shopID, &reqData)
	if err != nil {
		log.Printf("创建图案失败: %v", err)
		response.ServerError(c, err.Error())
		return
	}
	
	log.Printf("创建图案成功: ID=%d, 标题=%s", design.ID, design.Title)
	response.Success(c, design)
}

// UpdateDesign 更新图案
func UpdateDesign(c *gin.Context) {
	// 记录请求内容类型
	contentType := c.GetHeader("Content-Type")
	log.Printf("更新图案请求内容类型: %s", contentType)
	
	// 从路径参数获取图案ID
	designIDStr := c.Param("id")
	designID, err := strconv.ParseUint(designIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的图案ID")
		return
	}
	
	// 从上下文获取店铺ID
	shopID := c.GetUint(middleware.ContextShopID)
	log.Printf("从上下文获取到店铺ID: %d", shopID)
	
	// 处理表单数据
	title := c.PostForm("title")
	category := c.PostForm("category")
	description := c.PostForm("description")
	
	// 处理状态
	var status *int8
	statusStr := c.PostForm("status")
	if statusStr != "" {
		statusVal, err := strconv.Atoi(statusStr)
		if err == nil {
			statusInt8 := int8(statusVal)
			status = &statusInt8
		}
	}
	
	// 初始化请求数据
	reqData := struct {
		Title       string `json:"title"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
		Description string `json:"description"`
		Status      *int8  `json:"status"`
	}{
		Title:       title,
		Category:    category,
		Description: description,
		Status:      status,
	}
	
	// 处理文件上传（如果有）
	file, err := c.FormFile("image")
	if err == nil {
		// 创建上传目录
		uploadDir := "./static/uploads/designs/"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("创建目录失败: %v", err)
			response.ServerError(c, "上传目录创建失败")
			return
		}
		
		// 生成唯一文件名
		fileExt := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], fileExt)
		filePath := filepath.Join(uploadDir, fileName)
		
		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			log.Printf("保存文件失败: %v", err)
			response.ServerError(c, "文件保存失败")
			return
		}
		
		// 更新图片URL
		reqData.ImageURL = "/static/uploads/designs/" + fileName
	}
	
	log.Printf("更新图案，ID: %d, 店铺ID: %d", designID, shopID)
	
	design, err := designService.UpdateDesign(uint(designID), shopID, &reqData)
	if err != nil {
		log.Printf("更新图案失败: %v", err)
		response.ServerError(c, err.Error())
		return
	}
	
	log.Printf("更新图案成功: ID=%d", design.ID)
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
