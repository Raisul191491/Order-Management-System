package repository

import (
	"errors"
	"fmt"
	"oms/domain"
	"oms/model"

	"gorm.io/gorm"
)

type zoneRepository struct {
	masterDb  *gorm.DB
	replicaDb *gorm.DB
}

func NewZoneRepository(masterDB, replicaDB *gorm.DB) domain.ZoneRepository {
	return &zoneRepository{
		masterDb:  masterDB,
		replicaDb: replicaDB,
	}
}

func (r *zoneRepository) CreateZone(zone model.Zone) error {
	return r.masterDb.Create(&zone).Error
}

func (r *zoneRepository) GetZoneByID(id int64) (model.Zone, error) {
	var zone model.Zone
	err := r.replicaDb.First(&zone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Zone{}, fmt.Errorf("zone with ID %d not found", id)
		}
		return model.Zone{}, err
	}

	return zone, nil
}

func (r *zoneRepository) GetAllZones(limit, offset int) ([]model.Zone, error) {
	var zones []model.Zone
	err := r.replicaDb.Limit(limit).Offset(offset).Find(&zones).Error
	return zones, err
}

func (r *zoneRepository) UpdateZone(zone model.Zone) error {
	result := r.masterDb.Save(&zone)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("zone with ID %d not found", zone.ID)
	}

	return nil
}

func (r *zoneRepository) DeleteZone(id int64) error {
	result := r.masterDb.Delete(&model.Zone{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("zone with ID %d not found", id)
	}

	return nil
}

func (r *zoneRepository) GetZoneByName(name string) (model.Zone, error) {
	var zone model.Zone
	err := r.replicaDb.Where("name = ?", name).First(&zone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Zone{}, fmt.Errorf("zone with name '%s' not found", name)
		}
		return model.Zone{}, err
	}

	return zone, nil
}

func (r *zoneRepository) GetZonesByCityID(cityID int64, limit, offset int) ([]model.Zone, error) {
	var zones []model.Zone
	err := r.replicaDb.Where("city_id = ?", cityID).
		Limit(limit).Offset(offset).Find(&zones).Error
	return zones, err
}

func (r *zoneRepository) GetZoneByNameAndCityID(name string, cityID int64) (model.Zone, error) {
	var zone model.Zone
	err := r.replicaDb.Where("name = ? AND city_id = ?", name, cityID).First(&zone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Zone{}, fmt.Errorf("zone with name '%s' in city %d not found", name, cityID)
		}
		return model.Zone{}, err
	}

	return zone, nil
}

func (r *zoneRepository) CountZonesByCity(cityID int64) (int64, error) {
	var count int64
	err := r.replicaDb.Model(&model.Zone{}).Where("city_id = ?", cityID).Count(&count).Error
	return count, err
}
