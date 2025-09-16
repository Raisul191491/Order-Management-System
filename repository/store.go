package repository

import (
	"errors"
	"fmt"
	"oms/domain"
	"oms/model"

	"gorm.io/gorm"
)

type storeRepository struct {
	masterDb  *gorm.DB
	replicaDb *gorm.DB
}

func NewStoreRepository(masterDB, replicaDB *gorm.DB) domain.StoreRepository {
	return &storeRepository{
		masterDb:  masterDB,
		replicaDb: replicaDB,
	}
}

func (r *storeRepository) CreateStore(store model.Store) error {
	return r.masterDb.Create(&store).Error
}

func (r *storeRepository) GetStoreByID(id int64) (model.Store, error) {
	var store model.Store
	err := r.replicaDb.First(&store, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Store{}, fmt.Errorf("store with ID %d not found", id)
		}
		return model.Store{}, err
	}

	return store, nil
}

func (r *storeRepository) GetAllStores(limit, offset int) ([]model.Store, error) {
	var stores []model.Store
	err := r.replicaDb.Limit(limit).Offset(offset).Find(&stores).Error
	return stores, err
}

func (r *storeRepository) UpdateStore(store model.Store) error {
	result := r.masterDb.Save(&store)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("store with ID %d not found", store.ID)
	}

	return nil
}

func (r *storeRepository) DeleteStore(id int64) error {
	result := r.masterDb.Delete(&model.Store{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("store with ID %d not found", id)
	}

	return nil
}

func (r *storeRepository) GetStoreByName(name string) (model.Store, error) {
	var store model.Store
	err := r.replicaDb.Where("name = ?", name).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Store{}, fmt.Errorf("store with name '%s' not found", name)
		}
		return model.Store{}, err
	}

	return store, nil
}
