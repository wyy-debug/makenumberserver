// Add these routes to the existing router setup

// 管理员API组
adminGroup := v1.Group("/admin")
adminGroup.Use(middleware.AdminAuth())
{
	// 分类管理相关路由
	adminGroup.GET("/categories", api.GetCategories)
	adminGroup.POST("/categories", api.CreateCategory)
	adminGroup.PUT("/categories/:id", api.UpdateCategory)
	adminGroup.DELETE("/categories/:id", api.DeleteCategory)

	// 店铺相关路由
	adminGroup.GET("/shop", api.GetAdminShop)
	adminGroup.PUT("/shop", api.UpdateShop)
	adminGroup.GET("/shop/stats", api.GetShopStats)
}

// 公开的分类API
apiGroup.GET("/categories", api.GetCategories)

// 无需认证的API组
publicGroup := r.Group("/api")
{
	publicGroup.POST("/admin/login", api.Login)
	publicGroup.POST("/admin/register", api.Register)
	publicGroup.GET("/admin/check-exists", api.CheckSuperAdminExists)
	// 其他公开API...
} 