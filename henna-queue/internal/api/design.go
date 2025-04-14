package api

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	
	"example.com/henna-queue/internal/util/response"
	"example.com/henna-queue/internal/service"
	"example.com/henna-queue/internal/middleware"

)

type QueueRepository struct{}

func NewQueueRepository() *QueueRepository {
	return &QueueRepository{}
}

var designService = service.NewDesignService()

// GetDesigns 获取设计列表
func GetDesigns(c *gin.Context) {
	category := c.Query("category_id")
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	
	designs, total, err := designService.GetDesigns(1, category, page, pageSize)

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	log.Println("designs",designs)
	log.Println("total",total)
	// 调试输出找到的队列数据
	
	// 构造响应数据
	result := map[string]interface{}{
		"designs": designs,
		"total": total,
		"page": page,
		"page_size": pageSize,
	}
	
	response.Success(c, result)
}

// GetDesign 获取单个设计详情
func GetDesign(c *gin.Context) {
	// 获取设计ID
	designIDStr := c.Param("id")
	designID, err := strconv.ParseUint(designIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的设计ID")
		return
	}
	
	// 获取用户ID（如果有）
	var userID string
	if value, exists := c.Get(middleware.ContextUserID); exists {
		userID = value.(string)
	}
	
	// 调用服务获取设计详情
	design, err := designService.GetDesign(uint(designID), userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	
	response.Success(c, design)
}

// ToggleFavorite 收藏/取消收藏设计
func ToggleFavorite(c *gin.Context) {
	// 获取设计ID
	designIDStr := c.Param("id")
	designID, err := strconv.ParseUint(designIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的设计ID")
		return
	}
	
	// 获取用户ID
	userID, exists := c.Get(middleware.ContextUserID)
	if !exists {
		response.Unauthorized(c, "需要登录")
		return
	}
	
	// 调用服务切换收藏状态
	isFavorite, err := designService.ToggleFavorite(userID.(string), uint(designID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	
	// 构造响应
	result := map[string]interface{}{
		"is_favorite": isFavorite,
		"design_id": designID,
	}
	response.Success(c, result)
}

// GetUserFavorites 获取用户收藏的设计列表
func GetUserFavorites(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get(middleware.ContextUserID)
	if !exists {
		response.Unauthorized(c, "需要登录")
		return
	}
	
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	
	// 调用服务获取收藏列表
	designs, total, err := designService.GetUserFavorites(userID.(string), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	
	// 构造响应数据
	result := map[string]interface{}{
		"designs": designs,
		"total": total,
		"page": page,
		"page_size": pageSize,
	}
	
	response.Success(c, result)
}

// CreateDesign 创建设计
func CreateDesign(c *gin.Context) {
	// 获取店铺ID（优先从中间件获取，没有则从请求参数获取）
	var shopID uint
	shopID = 1
	// 获取设计基本信息
	name := c.PostForm("name")
	if name == "" {
		response.BadRequest(c, "设计名称不能为空")
		return
	}
	
	categoryID := c.PostForm("category_id")
	if categoryID == "" {
		response.BadRequest(c, "分类ID不能为空")
		return
	}
	
	description := c.PostForm("description")
	//price := c.PostForm("price")
	
	
	// 处理文件上传
	file, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, "图片上传失败: "+err.Error())
		return
	}
	
	// 创建上传目录
	uploadPath := viper.GetString("upload.path")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	
	// 确保目录存在
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		response.ServerError(c, "创建上传目录失败: "+err.Error())
		return
	}
	
	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	allowedTypes := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowedTypes[ext] {
		response.BadRequest(c, "不支持的图片格式，仅支持jpg、jpeg、png、gif")
		return
	}
	
	// 生成文件名：时间戳 + 随机数 + 扩展名
	fileName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), strconv.FormatUint(uint64(shopID), 10), ext)
	filePath := filepath.Join(uploadPath, fileName)
	
	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		response.ServerError(c, "创建文件失败: "+err.Error())
		return
	}
	defer dst.Close()
	
	src, err := file.Open()
	if err != nil {
		response.ServerError(c, "打开上传文件失败: "+err.Error())
		return
	}
	defer src.Close()
	
	if _, err = io.Copy(dst, src); err != nil {
		response.ServerError(c, "保存文件失败: "+err.Error())
		return
	}
	
	// 构造请求数据
	reqData := &struct {
		Title       string `json:"title" binding:"required"`
		Category    string `json:"category" binding:"required"`
		ImageURL    string `json:"image_url" binding:"required"`
		Description string `json:"description"`
	}{
		Title:       name,
		Category:    categoryID,
		ImageURL:    "/uploads/" + fileName, // 保存相对路径
		Description: description,
	}
	// 调用服务创建设计
	design, err := designService.CreateDesign(shopID, reqData)
	if err != nil {
		// 创建失败，删除已上传的文件
		os.Remove(filePath)
		response.ServerError(c, err.Error())
		return
	}
	// 成功响应
	response.Success(c, design)
}
