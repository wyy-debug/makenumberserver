package model

import (
	"time"
)

// Category 图案分类
type Category struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	ShopID    uint      `json:"shop_id" gorm:"index:idx_shop_code"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	Code      string    `json:"code" gorm:"size:50;not null;index:idx_shop_code"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} 