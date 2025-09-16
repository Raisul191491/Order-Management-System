package model

import (
	"time"
)

type Zone struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	CityID    int64      `json:"city_id" gorm:"not null;index"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
