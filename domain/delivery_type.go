package domain

import (
	"oms/model"
	"oms/types"
)

// DeliveryTypeRepository defines the interface for delivery type data operations
type DeliveryTypeRepository interface {
	CreateDeliveryType(deliveryType model.DeliveryType) error
	GetDeliveryTypeByID(id int64) (model.DeliveryType, error)
	GetAllDeliveryTypes(limit, offset int) ([]model.DeliveryType, error)
	UpdateDeliveryType(deliveryType model.DeliveryType) error
	DeleteDeliveryType(id int64) error
	GetDeliveryTypeByName(name string) (model.DeliveryType, error)
}

type DeliveryTypeService interface {
	CreateDeliveryType(deliveryType types.DeliveryTypeCreateRequest) error
	GetDeliveryTypeByID(id int64) (types.DeliveryTypeResponse, error)
	GetAllDeliveryTypes(limit, offset int) ([]types.DeliveryTypeResponse, error)
	UpdateDeliveryType(deliveryType types.DeliveryTypeUpdateRequest) error
	DeleteDeliveryType(id int64) error
}
