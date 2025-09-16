package model

import "time"

type City struct {
	ID              int64      `json:"id" gorm:"primaryKey:autoIncrement"`
	Name            string     `json:"name" gorm:"type:varchar(100);not null"`
	BaseDeliveryFee float64    `json:"base_delivery_fee" gorm:"type:decimal(10,2);default:100.00"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
