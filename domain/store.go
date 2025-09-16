package domain

import (
	"oms/model"
	"oms/types"
)

type StoreRepository interface {
	CreateStore(store model.Store) error
	GetStoreByID(id int64) (model.Store, error)
	GetAllStores(limit, offset int) ([]model.Store, error)
	GetStoreByName(name string) (model.Store, error)
	UpdateStore(store model.Store) error
	DeleteStore(id int64) error
}

type StoreService interface {
	CreateStore(store types.StoreCreateRequest) error
	GetStoreByID(id int64) (types.StoreResponse, error)
	GetAllStores(limit, offset int) ([]types.StoreResponse, error)
	UpdateStore(store types.StoreUpdateRequest) error
	DeleteStore(id int64) error
}
