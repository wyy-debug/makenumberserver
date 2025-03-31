package model

import (
	"time"
)

type Shop struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	Address       string    `gorm:"type:varchar(255);not null" json:"address"`
	Latitude      float64   `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude     float64   `gorm:"type:decimal(10,7)" json:"longitude"`
	Phone         string    `gorm:"type:varchar(20)" json:"phone"`
	BusinessHours string    `gorm:"type:varchar(100)" json:"business_hours"`
	Description   string    `gorm:"type:text" json:"description"`
	CoverImage    string    `gorm:"type:varchar(255)" json:"cover_image"`
	Rating        float64   `gorm:"type:decimal(2,1);default:5.0" json:"rating"`
	Status        int8      `gorm:"type:tinyint;default:1" json:"status"` // 1: 营业中, 0: 休息中
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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