package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Redis *redis.Client
var Ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := Redis.Ping(Ctx).Result()
	return err
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if Redis != nil {
		Redis.Close()
	}
} 