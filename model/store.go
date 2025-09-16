package model

import (
	"time"
)

type Store struct {
	ID           int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string     `json:"name" gorm:"type:varchar(255);not null"`
	ContactPhone string     `json:"contact_phone" gorm:"type:varchar(20)"`
	Address      string     `json:"address" gorm:"type:text"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
}
