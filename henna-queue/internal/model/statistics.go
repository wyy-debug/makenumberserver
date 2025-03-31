package model

import (
	"time"
)

type Statistic struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ShopID      uint      `gorm:"not null" json:"shop_id"`
	Date        time.Time `gorm:"type:date;not null" json:"date"`
	ServedCount int       `gorm:"default:0" json:"served_count"`
	QueueCount  int       `gorm:"default:0" json:"queue_count"`
	CancelCount int       `gorm:"default:0" json:"cancel_count"`
	AvgWaitTime int       `gorm:"default:0" json:"avg_wait_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DailyStatistics 每日统计数据
type DailyStatistics struct {
	Date        string `json:"date"`
	ServedCount int    `json:"served_count"`
	QueueCount  int    `json:"queue_count"`
	CancelCount int    `json:"cancel_count"`
	AvgWaitTime int    `json:"avg_wait_time"`
}

// StatisticsResponse 统计数据响应
type StatisticsResponse struct {
	Today      DailyStatistics   `json:"today"`
	Yesterday  DailyStatistics   `json:"yesterday"`
	Trend      []DailyStatistics `json:"trend"`
	Comparison struct {
		ServedPercentage float64 `json:"served_percentage"`
		QueuePercentage  float64 `json:"queue_percentage"`
		WaitPercentage   float64 `json:"wait_percentage"`
	} `json:"comparison"`
} 