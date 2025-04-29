package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"example.com/henna-queue/internal/api"
	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/pkg/db"
)

var globalUIPath string
var currentDir string

func main() {
	// 加载配置
	loadConfig()

	// 获取当前工作目录并打印
	var err error
	currentDir, err = os.Getwd()
	if err != nil {
		log.Fatalf("无法获取当前工作目录: %v", err)
	}
	log.Printf("当前工作目录: %s", currentDir)

	// 计算静态文件夹的绝对路径
	staticPath := filepath.Join(currentDir, "./static")
	absStaticPath, err := filepath.Abs(staticPath)
	if err != nil {
		log.Fatalf("无法获取静态文件夹的绝对路径: %v", err)
	}
	log.Printf("静态文件夹路径: %s", absStaticPath)


	// 计算UI文件夹的绝对路径
	uiPath := filepath.Join(currentDir, "../../UI")
	absUIPath, err := filepath.Abs(uiPath)
	if err != nil {
		log.Fatalf("无法获取UI文件夹的绝对路径: %v", err)
	}
	log.Printf("UI文件夹路径: %s", absUIPath)

	// 检查UI文件夹是否存在
	if _, err := os.Stat(absUIPath); os.IsNotExist(err) {
		log.Printf("警告: UI文件夹不存在: %s", absUIPath)
	} else {
		log.Printf("UI文件夹存在: %s", absUIPath)
	}

	// 保存UI路径供后续使用
	globalUIPath = absUIPath

	// 初始化数据库连接
	if err := db.InitMySQL(); err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.CloseMySQL()

	if err := db.InitRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer db.CloseRedis()

	// 设置Gin模式
	if viper.GetString("app.mode") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 注册路由
	setupRoutes(r)

	// 启动服务器
	port := viper.GetString("app.port")
	if port == "" {
		port = "8080" // 使用默认端口
	}
	server := fmt.Sprintf("0.0.0.0:%s", port)

	go func() {
		log.Printf("Server is running on 0.0.0.0:%s", port)
		if err := r.Run(server); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// 公共API
	public := r.Group("/api/v1")
	{
		// 认证相关
		auth := public.Group("/auth")
		{
			auth.POST("/login", api.Login)
		}

		// 店铺相关
		shops := public.Group("/shops")
		{
			shops.GET("/:id", api.GetShop)
			shops.GET("/:id/services", api.GetSvcsByShopId)
		}

		// 图案相关
		designs := public.Group("/designs")
		{
			designs.GET("", api.GetDesigns)
			//designs.GET("/:id", api.GetDesign)
			designs.POST("", api.CreateDesign)
		}

		// 添加缺失的公共API路由
		public.GET("/services", api.GetPublicSvcs)
		public.POST("/services", api.CreateSvc)
		public.GET("/settings", api.GetPublicSettings)
		
		// 公共队列查询API
		public.GET("/queues", api.GetQueues)
		public.POST("/queues", api.CreateQueue)
		public.GET("/queue/current", api.GetCurrentQueue)
	}

	// 需要用户认证的API
	user := r.Group("/api/v1")
	user.Use(middleware.AuthRequired())
	{
		// 用户相关
		user.GET("/users/profile", api.GetUserProfile)
		user.PUT("/users/profile", api.UpdateUserProfile)
		user.POST("/users/phone", api.UpdateUserPhone)

		// 排队相关
		queue := user.Group("/queue")
		{
			queue.GET("/status", api.GetQueueStatus)
			queue.POST("/number", api.GetQueueNumber)
			queue.DELETE("/number", api.CancelQueue)
		}

		// 图案收藏
		//user.POST("/designs/:id/like", api.ToggleFavorite)
		//user.GET("/users/favorites", api.GetUserFavorites)

		// 添加用户设置相关API
		settings := user.Group("/settings")
		{
			settings.GET("/backups", api.GetUserBackups)
		}
	}

	// 管理员API
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AdminRequired())
	{
		// 管理员个人资料
		admin.GET("/profile", api.GetAdminProfile)
		admin.PUT("/profile", api.UpdateAdminProfile)

		admin.GET("/queue", api.GetAdminQueue)
		admin.PUT("/queue/:id", api.UpdateQueueStatus)
		admin.POST("/queue/next", api.CallNextNumber)
		admin.POST("/queue/pause", api.ToggleQueuePause)
		
		// 添加管理员队列管理API
		admin.GET("/queues", api.GetQueues)
		admin.POST("/queues", api.CreateQueue)

		admin.GET("/statistics", api.GetStatistics)

		//admin.GET("/designs", api.GetAdminDesigns)
		//admin.POST("/designs", api.CreateDesign)
		//admin.PUT("/designs/:id", api.UpdateDesign)
		//admin.DELETE("/designs/:id", api.DeleteDesign)

		admin.GET("/shop", api.GetAdminShop)
		admin.PUT("/shop", api.UpdateShop)
		// 添加店铺统计API
		admin.GET("/shop/stats", api.GetShopStats)
		
		admin.GET("/services", api.GetAdminSvcs)
		admin.POST("/services", api.CreateSvc)
		admin.PUT("/services/:id", api.UpdateSvc)
		admin.DELETE("/services/:id", api.DeleteSvc)
	}

	// 管理员认证
	adminAuth := r.Group("/api/v1/admin/auth")
	{
		adminAuth.POST("/login", api.AdminLogin)
		adminAuth.POST("/logout", api.AdminLogout)
	}

	// 管理员注册相关API（不需要认证）
	adminRegister := r.Group("/api/v1/admin")
	{
		adminRegister.GET("/check-exists", api.CheckAdminExists)
		adminRegister.POST("/register", api.RegisterAdmin)
	}

	// 配置静态文件服务（放在最后）
	absStaticPath := filepath.Join(currentDir, "./static")
	absStaticPath, _ = filepath.Abs(absStaticPath)
	log.Printf("配置静态文件服务: 路径=%s, URL前缀=/static", absStaticPath)
	r.Static("/static", absStaticPath)
	
	// 配置上传文件目录为静态文件服务
	uploadsPath := viper.GetString("upload.path")
	if uploadsPath == "" {
		uploadsPath = "./uploads"
	}
	absUploadsPath, _ := filepath.Abs(uploadsPath)
	log.Printf("配置上传目录静态服务: 路径=%s, URL前缀=/uploads", absUploadsPath)
	// 确保目录存在
	if err := os.MkdirAll(absUploadsPath, 0755); err != nil {
		log.Printf("警告: 创建上传目录失败: %v", err)
	}
	r.Static("/uploads", absUploadsPath)
}
