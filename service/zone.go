package service

import (
	"fmt"
	"oms/domain"
	"oms/model"
	"oms/types"
)

type zoneService struct {
	zoneRepository domain.ZoneRepository
	cityRepository domain.CityRepository
}

func NewZoneService(zoneRepository domain.ZoneRepository, cityRepository domain.CityRepository) domain.ZoneService {
	return &zoneService{
		zoneRepository: zoneRepository,
		cityRepository: cityRepository,
	}
}

func (zs zoneService) CreateZone(zone types.ZoneCreateRequest) error {
	// Check if city exists
	_, err := zs.cityRepository.GetCityByID(zone.CityID)
	if err != nil {
		return fmt.Errorf("city with ID %d does not exist", zone.CityID)
	}

	// Check if zone with same name already exists in the same city
	existing, err := zs.zoneRepository.GetZoneByNameAndCityID(zone.Name, zone.CityID)
	if err == nil && existing.ID != 0 {
		return fmt.Errorf("zone with name '%s' already exists in this city", zone.Name)
	}

	newZone := model.Zone{
		CityID: zone.CityID,
		Name:   zone.Name,
	}

	err = zs.zoneRepository.CreateZone(newZone)
	if err != nil {
		return err
	}

	return nil
}

func (zs zoneService) GetZoneByID(id int64) (types.ZoneResponse, error) {
	existingZone, err := zs.zoneRepository.GetZoneByID(id)
	if err != nil {
		return types.ZoneResponse{}, err
	}

	return types.ZoneResponse{
		ID:        existingZone.ID,
		CityID:    existingZone.CityID,
		Name:      existingZone.Name,
		CreatedAt: existingZone.CreatedAt,
		UpdatedAt: existingZone.UpdatedAt,
	}, nil
}

func (zs zoneService) GetAllZones(limit, offset int) ([]types.ZoneResponse, error) {
	existingZones, err := zs.zoneRepository.GetAllZones(limit, offset)
	if err != nil {
		return nil, err
	}

	var result []types.ZoneResponse

	for _, existingZone := range existingZones {
		result = append(result, types.ZoneResponse{
			ID:        existingZone.ID,
			CityID:    existingZone.CityID,
			Name:      existingZone.Name,
			CreatedAt: existingZone.CreatedAt,
			UpdatedAt: existingZone.UpdatedAt,
		})
	}

	return result, nil
}

func (zs zoneService) GetZonesByCityID(cityID int64, limit, offset int) ([]types.ZoneResponse, error) {
	// Check if city exists
	_, err := zs.cityRepository.GetCityByID(cityID)
	if err != nil {
		return nil, fmt.Errorf("city with ID %d does not exist", cityID)
	}

	existingZones, err := zs.zoneRepository.GetZonesByCityID(cityID, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []types.ZoneResponse

	for _, existingZone := range existingZones {
		result = append(result, types.ZoneResponse{
			ID:        existingZone.ID,
			CityID:    existingZone.CityID,
			Name:      existingZone.Name,
			CreatedAt: existingZone.CreatedAt,
			UpdatedAt: existingZone.UpdatedAt,
		})
	}

	return result, nil
}

func (zs zoneService) UpdateZone(zone types.ZoneUpdateRequest) error {
	existingZone, err := zs.zoneRepository.GetZoneByID(zone.ID)
	if err != nil {
		return err
	}

	if zone.Name != existingZone.Name {
		existing, err := zs.zoneRepository.GetZoneByNameAndCityID(zone.Name, existingZone.CityID)
		if err == nil && existing.ID != 0 && existing.ID != existingZone.ID {
			return fmt.Errorf("zone with name '%s' already exists in this city", zone.Name)
		}

		existingZone.Name = zone.Name
	}

	err = zs.zoneRepository.UpdateZone(existingZone)
	if err != nil {
		return err
	}

	return nil
}

func (zs zoneService) DeleteZone(id int64) error {
	existingZone, err := zs.zoneRepository.GetZoneByID(id)
	if err != nil || existingZone.ID == 0 {
		return fmt.Errorf("zone does not exist")
	}

	return zs.zoneRepository.DeleteZone(id)
}
