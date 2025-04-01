package model

import (
	"time"
)

// Shop 店铺模型
type Shop struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"size:100;not null"`
	Description   string    `json:"description" gorm:"type:text"`
	Phone         string    `json:"phone" gorm:"size:20"`
	BusinessHours string    `json:"business_hours" gorm:"size:100"`
	Address       string    `json:"address" gorm:"size:200"`
	Latitude      float64   `json:"latitude" gorm:"type:decimal(10,7)"`
	Longitude     float64   `json:"longitude" gorm:"type:decimal(10,7)"`
	CoverImage    string    `json:"cover_image" gorm:"size:255"`
	Rating        float64   `json:"rating" gorm:"type:decimal(2,1);default:5.0"`
	Status        int       `json:"status" gorm:"default:1"` // 1: 营业中, 0: 已关闭
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ShopStats 店铺统计信息
type ShopStats struct {
	TodayQueueCount     int     `json:"today_queue_count"`
	TodayCompletedCount int     `json:"today_completed_count"`
	AvgWaitTime         int     `json:"avg_wait_time"` // 平均等待时间（分钟）
	CancelRate          float64 `json:"cancel_rate"`   // 取消率（百分比）
}

type Service struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ShopID      uint      `gorm:"not null" json:"shop_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Duration    int       `gorm:"not null" json:"duration"` // 服务时长(分钟)
	Description string    `gorm:"type:text" json:"description"`
	Status      int8      `gorm:"type:tinyint;default:1" json:"status"` // 1: 可用, 0: 不可用
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
