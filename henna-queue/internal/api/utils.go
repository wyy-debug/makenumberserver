package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
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
