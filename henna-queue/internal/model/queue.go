package model

import (
	"time"
)

type Queue struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ShopID            uint      `gorm:"not null" json:"shop_id"`
	QueueNumber       string    `gorm:"type:varchar(10);not null" json:"queue_number"`
	UserID            string    `gorm:"type:varchar(50);not null" json:"user_id"`
	ServiceID         uint      `gorm:"not null" json:"service_id"`
	Status            int8      `gorm:"type:tinyint;default:0" json:"status"` // 0: 等待中, 1: 就位中, 2: 服务中, 3: 已完成, 4: 已取消
	EstimatedWaitTime int       `json:"estimated_wait_time"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// 关联
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Service Service `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
}

// QueueStatusResponse 排队状态响应
type QueueStatusResponse struct {
	HasNumber      bool   `json:"has_number"`
	QueueNumber    string `json:"queue_number,omitempty"`
	PeopleAhead    int    `json:"people_ahead"`
	WaitTime       int    `json:"wait_time"`
	CurrentServing string `json:"current_serving"`
	CurrentWaiting string `json:"current_waiting"`
	TotalServed    int    `json:"total_served"`
}

// CurrentQueueResponse 当前叫号情况
type CurrentQueueResponse struct {
	CurrentServing string `json:"current_serving"`
	CurrentWaiting string `json:"current_waiting"`
	TotalServed    int    `json:"total_served"`
} 