package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"example.com/henna-queue/internal/api"
	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/pkg/db"
)

func main() {
	// 加载配置
	loadConfig()

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
	server := fmt.Sprintf(":%s", port)

	go func() {
		log.Printf("Server is running on port %s", port)
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
			shops.GET("/:id/services", api.GetShopServices)
		}

		// 图案相关
		designs := public.Group("/designs")
		{
			designs.GET("", api.GetDesigns)
			designs.GET("/:id", api.GetDesign)
		}
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
			queue.GET("/current", api.GetCurrentQueue)
		}

		// 图案收藏
		user.POST("/designs/:id/like", api.ToggleFavorite)
		user.GET("/users/favorites", api.GetUserFavorites)
	}

	// 管理员API
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AdminRequired())
	{
		admin.GET("/queue", api.GetAdminQueue)
		admin.PUT("/queue/:id", api.UpdateQueueStatus)
		admin.POST("/queue/next", api.CallNextNumber)
		admin.POST("/queue/pause", api.ToggleQueuePause)

		admin.GET("/statistics", api.GetStatistics)

		admin.GET("/designs", api.GetAdminDesigns)
		admin.POST("/designs", api.CreateDesign)
		admin.PUT("/designs/:id", api.UpdateDesign)
		admin.DELETE("/designs/:id", api.DeleteDesign)

		admin.GET("/shop", api.GetAdminShop)
		admin.PUT("/shop", api.UpdateShop)
		admin.GET("/services", api.GetAdminServices)
		admin.POST("/services", api.CreateService)
		admin.PUT("/services/:id", api.UpdateService)
		admin.DELETE("/services/:id", api.DeleteService)
	}

	// 管理员认证
	adminAuth := r.Group("/api/v1/admin/auth")
	{
		adminAuth.POST("/login", api.AdminLogin)
		adminAuth.POST("/logout", api.AdminLogout)
	}
}
