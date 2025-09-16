package model

import (
	"time"
)

type Order struct {
	ID                 int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	ConsignmentID      string     `json:"consignment_id" gorm:"type:varchar(50);uniqueIndex;not null"`
	UserID             int64      `json:"user_id" gorm:"index"`
	StoreID            int64      `json:"store_id" gorm:"index"`
	MerchantOrderID    string     `json:"merchant_order_id" gorm:"type:varchar(100)"`
	RecipientName      string     `json:"recipient_name" gorm:"type:varchar(255);not null"`
	RecipientPhone     string     `json:"recipient_phone" gorm:"type:varchar(20);not null"`
	RecipientAddress   string     `json:"recipient_address" gorm:"type:text;not null"`
	RecipientCity      int64      `json:"recipient_city" gorm:"index"`
	RecipientZone      int64      `json:"recipient_zone" gorm:"index"`
	RecipientArea      string     `json:"recipient_area" gorm:"type:text"`
	OrderType          string     `json:"order_type" gorm:"type:order_type_enum;not null;default:'delivery'"`
	DeliveryTypeID     int64      `json:"delivery_type_id" gorm:"index"`
	ItemType           int64      `json:"item_type" gorm:"index"`
	ItemQuantity       int        `json:"item_quantity" gorm:"not null;default:1"`
	ItemWeight         float64    `json:"item_weight" gorm:"type:decimal(8,2);not null"`
	ItemDescription    string     `json:"item_description" gorm:"type:text"`
	SpecialInstruction string     `json:"special_instruction" gorm:"type:text"`
	OrderAmount        float64    `json:"order_amount" gorm:"type:decimal(10,2);not null"`
	AmountToCollect    float64    `json:"amount_to_collect" gorm:"type:decimal(10,2);not null"`
	DeliveryFee        float64    `json:"delivery_fee" gorm:"type:decimal(10,2);not null"`
	CodFee             float64    `json:"cod_fee" gorm:"type:decimal(10,2);not null;default:0"`
	PromoDiscount      float64    `json:"promo_discount" gorm:"type:decimal(10,2);default:0"`
	Discount           float64    `json:"discount" gorm:"type:decimal(10,2);default:0"`
	TotalFee           float64    `json:"total_fee" gorm:"type:decimal(10,2);not null"`
	OrderStatus        string     `json:"order_status" gorm:"type:order_status_enum;not null;default:'pending'"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          *time.Time `json:"deleted_at" gorm:"index"`
}
