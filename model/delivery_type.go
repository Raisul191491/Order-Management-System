package model

import (
	"time"
)

type DeliveryType struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"type:varchar(50);not null;uniqueIndex"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
