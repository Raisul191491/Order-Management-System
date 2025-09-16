package repository

import (
	"errors"
	"fmt"
	"oms/domain"
	"oms/model"

	"gorm.io/gorm"
)

type itemTypeRepository struct {
	masterDb  *gorm.DB
	replicaDb *gorm.DB
}

func NewItemTypeRepository(masterDB, replicaDB *gorm.DB) domain.ItemTypeRepository {
	return &itemTypeRepository{
		masterDb:  masterDB,
		replicaDb: replicaDB,
	}
}

func (r *itemTypeRepository) CreateItemType(itemType model.ItemType) error {
	return r.masterDb.Create(&itemType).Error
}

func (r *itemTypeRepository) GetItemTypeByID(id int64) (model.ItemType, error) {
	var itemType model.ItemType
	err := r.replicaDb.First(&itemType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ItemType{}, fmt.Errorf("item type with ID %d not found", id)
		}
		return model.ItemType{}, err
	}

	return itemType, nil
}

func (r *itemTypeRepository) GetAllItemTypes(limit, offset int) ([]model.ItemType, error) {
	var itemTypes []model.ItemType
	err := r.replicaDb.Order("name ASC").Limit(limit).Offset(offset).Find(&itemTypes).Error
	return itemTypes, err
}

func (r *itemTypeRepository) UpdateItemType(itemType model.ItemType) error {
	result := r.masterDb.Save(&itemType)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("item type with ID %d not found", itemType.ID)
	}

	return nil
}

func (r *itemTypeRepository) DeleteItemType(id int64) error {
	result := r.masterDb.Delete(&model.ItemType{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("item type with ID %d not found", id)
	}

	return nil
}

func (r *itemTypeRepository) GetItemTypeByName(name string) (model.ItemType, error) {
	var itemType model.ItemType
	err := r.replicaDb.Where("name = ?", name).First(&itemType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ItemType{}, fmt.Errorf("item type with name '%s' not found", name)
		}
		return model.ItemType{}, err
	}

	return itemType, nil
}
