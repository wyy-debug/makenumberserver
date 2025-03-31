package model

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	UnionID   string    `gorm:"type:varchar(50)" json:"union_id"`
	Nickname  string    `gorm:"type:varchar(50)" json:"nickname"`
	AvatarURL string    `gorm:"type:varchar(255)" json:"avatar_url"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} 