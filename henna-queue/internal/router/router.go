// Add these routes to the existing router setup

// 分类管理相关路由
adminGroup.GET("/categories", api.GetCategories)
adminGroup.POST("/categories", api.CreateCategory)
adminGroup.PUT("/categories/:id", api.UpdateCategory)
adminGroup.DELETE("/categories/:id", api.DeleteCategory)

// 公开的分类API
apiGroup.GET("/categories", api.GetCategories)

// 无需认证的API组
publicGroup := r.Group("/api")
{
    publicGroup.POST("/admin/login", api.Login)
    publicGroup.POST("/admin/register", api.Register)
    publicGroup.GET("/admin/check-exists", api.CheckAdminExists)
    // 其他公开API...
} 