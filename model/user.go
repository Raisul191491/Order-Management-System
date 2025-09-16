package model

import (
	"time"
)

type User struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	Email        string     `json:"email" gorm:"type:varchar(255);unique;not null"`
	PasswordHash string     `json:"-" gorm:"type:varchar(255);not null"` // Hidden from JSON
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
