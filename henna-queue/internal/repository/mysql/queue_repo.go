package mysql

import (
	"time"
	"log"
	"example.com/henna-queue/internal/model"
	"example.com/henna-queue/pkg/db"
	"errors"
)

type QueueRepository struct {
}

func NewQueueRepository() *QueueRepository {
	// 不在这里 panic，允许延迟初始化
	return &QueueRepository{}
}

// GetByID 根据ID获取队列
func (r *QueueRepository) GetByID(id uint) (*model.Queue, error) {

	var queue model.Queue
	result := db.DB.Preload("User").Preload("Service").Where("id = ?", id).First(&queue)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queue, nil
}

// GetActiveByUserID 获取用户当前活跃的排队
func (r *QueueRepository) GetActiveByUserID(userID string, shopID uint) (*model.Queue, error) {

	var queue model.Queue
	result := db.DB.Preload("Service").
		Where("user_id = ? AND shop_id = ? AND status < ?", userID, shopID, 3).
		First(&queue)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queue, nil
}

// GetByStatus 获取指定状态的队列
func (r *QueueRepository) GetByStatus(shopID uint, status int8) (*model.Queue, error) {

	var queue model.Queue
	result := db.DB.Preload("Service").Preload("User").
		Where("shop_id = ? AND status = ?", shopID, status).
		First(&queue)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queue, nil
}

// GetNextWaiting 获取下一个等待中的队列
func (r *QueueRepository) GetNextWaiting(shopID uint) (*model.Queue, error) {
	
	var queue model.Queue
	result := db.DB.Preload("Service").Preload("User").
		Where("shop_id = ? AND status = ?", shopID, 0).
		Order("created_at").
		First(&queue)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queue, nil
}

// GetWaitingCount 获取等待中的人数
func (r *QueueRepository) GetWaitingCount(shopID uint) (int, error) {
	var count int64
	result := db.DB.Model(&model.Queue{}).
		Where("shop_id = ? AND status = ?", shopID, 0).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

// GetPeopleAhead 获取前方等待人数
func (r *QueueRepository) GetPeopleAhead(shopID uint, queueNumber string) (int, error) {
	
	var count int64
	result := db.DB.Model(&model.Queue{}).
		Where("shop_id = ? AND status = ? AND created_at < (SELECT created_at FROM queues WHERE shop_id = ? AND queue_number = ? LIMIT 1)",
			shopID, 0, shopID, queueNumber).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

// GetDailyCount 获取当日队列数量
func (r *QueueRepository) GetDailyCount(shopID uint, date string) (int, error) {
	
	startTime, _ := time.Parse("20060102", date)
	endTime := startTime.AddDate(0, 0, 1)

	var count int64
	result := db.DB.Model(&model.Queue{}).
		Where("shop_id = ? AND created_at >= ? AND created_at < ?", shopID, startTime, endTime).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

// GetTodayServedCount 获取今日已服务人数
func (r *QueueRepository) GetTodayServedCount(shopID uint) (int, error) {
	
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endTime := startTime.AddDate(0, 0, 1)

	var count int64
	result := db.DB.Model(&model.Queue{}).
		Where("shop_id = ? AND status = ? AND updated_at >= ? AND updated_at < ?", shopID, 3, startTime, endTime).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

// GetActiveQueues 获取所有活跃队列
func (r *QueueRepository) GetActiveQueues(shopID uint) ([]*model.Queue, error) {
	
	var queues []*model.Queue
	result := db.DB.Preload("User").Preload("Service").
		//Where("shop_id = ? AND status < ?", shopID, 3).
		Where("shop_id = ?",shopID).
		//Order("status ASC, created_at ASC").
		Find(&queues)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Printf("queue: %v", shopID)
	return queues, nil
}

// Create 创建队列
func (r *QueueRepository) Create(queue *model.Queue) error {
	
	return db.DB.Create(queue).Error
}

// Update 更新队列
func (r *QueueRepository) Update(queue *model.Queue) error {

	return db.DB.Save(queue).Error
}
