package model

import (
	"time"
)

type TattooDesign struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ShopID      uint      `gorm:"not null" json:"shop_id"`
	Title       string    `gorm:"type:varchar(100);not null" json:"title"`
	Category    string    `gorm:"type:varchar(50);not null" json:"category"`
	ImageURL    string    `gorm:"type:varchar(255);not null" json:"image_url"`
	Description string    `gorm:"type:text" json:"description"`
	Likes       int       `gorm:"default:0" json:"likes"`
	Status      int8      `gorm:"type:tinyint;default:1" json:"status"` // 1: 可用, 0: 不可用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(50);not null" json:"user_id"`
	DesignID  uint      `gorm:"not null" json:"design_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Design TattooDesign `gorm:"foreignKey:DesignID" json:"design,omitempty"`
}

// DesignResponse 带有是否收藏信息的图案响应
type DesignResponse struct {
	TattooDesign
	IsLiked bool `json:"is_liked"`
} 