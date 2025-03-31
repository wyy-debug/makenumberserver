package model

import (
	"time"
)

type Admin struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"type:varchar(50);not null;unique" json:"username"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	ShopID       *uint      `json:"shop_id"`
	Role         int8       `gorm:"type:tinyint;not null;default:1" json:"role"` // 1: 普通管理员, 2: 超级管理员
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// 关联
	Shop *Shop `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
} 