package model

import (
	"time"
)

type UserSession struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	UserID       int64      `json:"user_id" gorm:"not null"`
	AccessToken  string     `json:"access_token" gorm:"type:varchar(255);not null"`
	RefreshToken string     `json:"refresh_token" gorm:"type:varchar(255);not null"`
	ExpiresAt    time.Time  `json:"expires_at" gorm:"not null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
