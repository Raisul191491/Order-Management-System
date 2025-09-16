package repository

import (
	"errors"
	"fmt"
	"oms/domain"
	"oms/model"

	"gorm.io/gorm"
)

type deliveryTypeRepository struct {
	masterDb  *gorm.DB
	replicaDb *gorm.DB
}

func NewDeliveryTypeRepository(masterDB, replicaDB *gorm.DB) domain.DeliveryTypeRepository {
	return &deliveryTypeRepository{
		masterDb:  masterDB,
		replicaDb: replicaDB,
	}
}

func (r *deliveryTypeRepository) CreateDeliveryType(deliveryType model.DeliveryType) error {
	return r.masterDb.Create(&deliveryType).Error
}

func (r *deliveryTypeRepository) GetDeliveryTypeByID(id int64) (model.DeliveryType, error) {
	var deliveryType model.DeliveryType
	err := r.replicaDb.First(&deliveryType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.DeliveryType{}, fmt.Errorf("delivery type with ID %d not found", id)
		}
		return model.DeliveryType{}, err
	}

	return deliveryType, nil
}

func (r *deliveryTypeRepository) GetAllDeliveryTypes(limit, offset int) ([]model.DeliveryType, error) {
	var deliveryTypes []model.DeliveryType
	err := r.replicaDb.Order("name ASC").Limit(limit).Offset(offset).Find(&deliveryTypes).Error
	return deliveryTypes, err
}

func (r *deliveryTypeRepository) UpdateDeliveryType(deliveryType model.DeliveryType) error {
	result := r.masterDb.Save(&deliveryType)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("delivery type with ID %d not found", deliveryType.ID)
	}

	return nil
}

func (r *deliveryTypeRepository) DeleteDeliveryType(id int64) error {
	result := r.masterDb.Delete(&model.DeliveryType{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("delivery type with ID %d not found", id)
	}

	return nil
}

func (r *deliveryTypeRepository) GetDeliveryTypeByName(name string) (model.DeliveryType, error) {
	var deliveryType model.DeliveryType
	err := r.replicaDb.Where("name = ?", name).First(&deliveryType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.DeliveryType{}, fmt.Errorf("delivery type with name '%s' not found", name)
		}
		return model.DeliveryType{}, err
	}

	return deliveryType, nil
}
