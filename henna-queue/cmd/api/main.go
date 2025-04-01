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

	// 检查静态文件夹是否存在
	if _, err := os.Stat(absStaticPath); os.IsNotExist(err) {
		log.Printf("警告: 静态文件夹不存在，将创建它: %s", absStaticPath)
		if err := os.MkdirAll(filepath.Join(absStaticPath, "html"), 0755); err != nil {
			log.Fatalf("无法创建静态文件夹: %v", err)
		}
		// 创建一个测试HTML文件
		testHTML := `<!DOCTYPE html>
<html>
<head>
    <title>测试页面</title>
</head>
<body>
    <h1>测试成功！</h1>
    <p>静态文件服务器工作正常。</p>
</body>
</html>`
		if err := os.WriteFile(filepath.Join(absStaticPath, "html", "test.html"), []byte(testHTML), 0644); err != nil {
			log.Fatalf("无法创建测试HTML文件: %v", err)
		}
		log.Printf("已创建测试HTML文件: %s", filepath.Join(absStaticPath, "html", "test.html"))
	} else {
		log.Printf("静态文件夹存在: %s", absStaticPath)
		// 列出html目录中的文件
		files, err := os.ReadDir(filepath.Join(absStaticPath, "html"))
		if err != nil {
			log.Printf("警告: 无法读取html目录: %v", err)
		} else {
			log.Printf("html目录中的文件:")
			for _, file := range files {
				log.Printf("  - %s", file.Name())
			}
		}
	}

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
}
