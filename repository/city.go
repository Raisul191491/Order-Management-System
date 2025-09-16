package repository

import (
	"errors"
	"fmt"
	"oms/domain"
	"oms/model"

	"gorm.io/gorm"
)

type cityRepository struct {
	masterDb  *gorm.DB
	replicaDb *gorm.DB
}

func NewCityRepository(masterDB, replicaDB *gorm.DB) domain.CityRepository {
	return &cityRepository{
		masterDb:  masterDB,
		replicaDb: replicaDB,
	}
}

func (r *cityRepository) CreateCity(city model.City) error {
	return r.masterDb.Create(&city).Error
}

func (r *cityRepository) GetCityByID(id int64) (model.City, error) {
	var city model.City
	err := r.replicaDb.First(&city, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.City{}, fmt.Errorf("city with ID %d not found", id)
		}
		return model.City{}, err
	}

	return city, nil
}

func (r *cityRepository) GetAllCities(limit, offset int) ([]model.City, error) {
	var cities []model.City
	err := r.replicaDb.Limit(limit).Offset(offset).Find(&cities).Error
	return cities, err
}

func (r *cityRepository) UpdateCity(city model.City) error {
	result := r.masterDb.Save(&city)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("city with ID %d not found", city.ID)
	}

	return nil
}

func (r *cityRepository) DeleteCity(id int64) error {
	result := r.masterDb.Delete(&model.City{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *cityRepository) GetCityByName(name string) (model.City, error) {
	var city model.City
	err := r.replicaDb.Where("name = ?", name).First(&city).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.City{}, fmt.Errorf("city with name '%s' not found", name)
		}
		return model.City{}, err
	}

	return city, nil
}
