package domain

import (
	"oms/model"
	"oms/types"
)

type ItemTypeRepository interface {
	CreateItemType(itemType model.ItemType) error
	GetItemTypeByID(id int64) (model.ItemType, error)
	GetAllItemTypes(limit, offset int) ([]model.ItemType, error)
	UpdateItemType(itemType model.ItemType) error
	DeleteItemType(id int64) error
	GetItemTypeByName(name string) (model.ItemType, error)
}

type ItemTypeService interface {
	CreateItemType(itemType types.ItemTypeCreateRequest) error
	GetItemTypeByID(id int64) (types.ItemTypeResponse, error)
	GetAllItemTypes(limit, offset int) ([]types.ItemTypeResponse, error)
	UpdateItemType(itemType types.ItemTypeUpdateRequest) error
	DeleteItemType(id int64) error
}
