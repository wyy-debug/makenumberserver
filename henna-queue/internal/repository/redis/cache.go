package redis

import (
	"encoding/json"
	"fmt"
	"henna-queue/internal/model"
	"henna-queue/pkg/db"
	"time"
)

type CacheRepository struct{}

func NewCacheRepository() *CacheRepository {
	return &CacheRepository{}
}

// UpdateQueueStatus 更新队列状态缓存
func (r *CacheRepository) UpdateQueueStatus(shopID uint) error {
	key := fmt.Sprintf("shop:%d:queue_status", shopID)
	
	// 从MySQL查询当前队列状态
	var serving, waiting model.Queue
	db.DB.Where("shop_id = ? AND status = ?", shopID, 2).First(&serving)
	db.DB.Where("shop_id = ? AND status = ?", shopID, 1).First(&waiting)
	
	// 构建缓存数据
	status := struct {
		CurrentServing string `json:"current_serving"`
		CurrentWaiting string `json:"current_waiting"`
		UpdatedAt      int64  `json:"updated_at"`
	}{
		UpdatedAt: time.Now().Unix(),
	}
	
	if serving.ID > 0 {
		status.CurrentServing = serving.QueueNumber
	}
	
	if waiting.ID > 0 {
		status.CurrentWaiting = waiting.QueueNumber
	}
	
	// 序列化并存储
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	
	return db.Redis.Set(db.Ctx, key, string(data), 24*time.Hour).Err()
}

// GetQueueStatus 获取队列状态缓存
func (r *CacheRepository) GetQueueStatus(shopID uint) (string, string, error) {
	key := fmt.Sprintf("shop:%d:queue_status", shopID)
	
	// 从Redis获取缓存
	data, err := db.Redis.Get(db.Ctx, key).Result()
	if err != nil {
		return "", "", err
	}
	
	// 解析数据
	var status struct {
		CurrentServing string `json:"current_serving"`
		CurrentWaiting string `json:"current_waiting"`
		UpdatedAt      int64  `json:"updated_at"`
	}
	
	if err := json.Unmarshal([]byte(data), &status); err != nil {
		return "", "", err
	}
	
	return status.CurrentServing, status.CurrentWaiting, nil
}

// CacheToken 缓存用户token
func (r *CacheRepository) CacheToken(userID string, token string, expiration time.Duration) error {
	key := fmt.Sprintf("token:user:%s", userID)
	return db.Redis.Set(db.Ctx, key, token, expiration).Err()
}

// GetCachedToken 获取缓存的token
func (r *CacheRepository) GetCachedToken(userID string) (string, error) {
	key := fmt.Sprintf("token:user:%s", userID)
	return db.Redis.Get(db.Ctx, key).Result()
}

// InvalidateToken 使token失效
func (r *CacheRepository) InvalidateToken(userID string) error {
	key := fmt.Sprintf("token:user:%s", userID)
	return db.Redis.Del(db.Ctx, key).Err()
} 