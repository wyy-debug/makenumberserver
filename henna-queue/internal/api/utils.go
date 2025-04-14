package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	
	"example.com/henna-queue/internal/middleware"
	"example.com/henna-queue/internal/util/response"
)

// GetInt8 获取 int8 类型的参数
func GetInt8(c *gin.Context, key string) (int8, bool) {
	val := c.Param(key)
	if val == "" {
		val = c.Query(key)
	}
	if val == "" {
		return 0, false
	}

	i, err := strconv.ParseInt(val, 10, 8)
	if err != nil {
		return 0, false
	}

	return int8(i), true
}

// GetInt8FromForm 从表单中获取 int8 类型的参数
func GetInt8FromForm(c *gin.Context, key string) (int8, bool) {
	val := c.PostForm(key)
	if val == "" {
		return 0, false
	}

	i, err := strconv.ParseInt(val, 10, 8)
	if err != nil {
		return 0, false
	}

	return int8(i), true
}

// GetPublicSettings 获取公共设置
func GetPublicSettings(c *gin.Context) {
	// 公共设置内容，可以从配置或数据库获取
	settings := map[string]interface{}{
		"app_name":        "Henna排队系统",
		"app_version":     "1.0.0",
		"max_queue_size":  100,
		"maintenance_mode": false,
		"notice":          "",
	}

	response.Success(c, settings)
}

// GetUserBackups 获取用户备份数据
func GetUserBackups(c *gin.Context) {
	// 从上下文获取用户ID
	// 如果函数内不需要使用userID，可以通过_接收或直接删除
	_ = c.GetString(middleware.ContextUserID)
	
	// 这里应该根据实际情况从数据库获取用户备份数据
	// 这里只是示例实现
	backups := []map[string]interface{}{
		{
			"id": 1,
			"name": "自动备份",
			"created_at": "2023-10-01 10:00:00",
			"size": "2.5MB",
		},
	}

	response.Success(c, backups)
}
