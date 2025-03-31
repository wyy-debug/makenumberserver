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
	Today      *DailyStatistics   `json:"today"`     // 当天数据
	Yesterday  *DailyStatistics   `json:"yesterday"` // 昨天数据
	Trend      []*DailyStatistics `json:"trend"`     // 趋势数据
	Comparison struct {
		ServedPercentage float64 `json:"served_percentage"` // 服务人数环比
		QueuePercentage  float64 `json:"queue_percentage"`  // 排队人数环比
		WaitPercentage   float64 `json:"wait_percentage"`   // 等待时间环比
	} `json:"comparison"`
}
