package service

import (
	"fmt"
	"oms/domain"
	"oms/model"
	"oms/types"
)

// CityService provides business logic for city operations
type cityService struct {
	cityRepository domain.CityRepository
}

func NewCityService(cityRepository domain.CityRepository) domain.CityService {
	return &cityService{cityRepository: cityRepository}
}

func (cs cityService) CreateCity(city types.CityCreateRequest) error {
	existing, err := cs.cityRepository.GetCityByName(city.Name)
	if err == nil && existing.ID != 0 {
		return fmt.Errorf("city with name '%s' already exists", city.Name)
	}

	newCity := model.City{
		Name:            city.Name,
		BaseDeliveryFee: city.BaseDeliveryFee,
	}

	err = cs.cityRepository.CreateCity(newCity)
	if err != nil {
		return err
	}

	return nil
}

func (cs cityService) GetCityByID(id int64) (types.CityResponse, error) {
	existingCity, err := cs.cityRepository.GetCityByID(id)
	if err != nil {
		return types.CityResponse{}, err
	}

	return types.CityResponse{
		Name:            existingCity.Name,
		BaseDeliveryFee: existingCity.BaseDeliveryFee,
		UpdatedAt:       existingCity.UpdatedAt,
	}, nil
}

func (cs cityService) GetAllCities(limit, offset int) ([]types.CityResponse, error) {
	existingCities, err := cs.cityRepository.GetAllCities(limit, offset)
	if err != nil {
		return nil, err
	}

	var result []types.CityResponse

	for _, existingCity := range existingCities {
		result = append(result, types.CityResponse{
			Id:              existingCity.ID,
			Name:            existingCity.Name,
			BaseDeliveryFee: existingCity.BaseDeliveryFee,
			UpdatedAt:       existingCity.UpdatedAt,
		})
	}

	return result, nil
}

func (cs cityService) GetCityByName(name string) (types.CityResponse, error) {
	existingCity, err := cs.cityRepository.GetCityByName(name)
	if err != nil {
		return types.CityResponse{}, err
	}

	return types.CityResponse{
		Name:            existingCity.Name,
		BaseDeliveryFee: existingCity.BaseDeliveryFee,
		UpdatedAt:       existingCity.UpdatedAt,
	}, nil
}

func (cs cityService) UpdateCity(city types.CityUpdateRequest) error {
	existingCity, err := cs.cityRepository.GetCityByID(city.ID)
	if err != nil {
		return err
	}

	if city.Name != existingCity.Name {
		existing, err := cs.cityRepository.GetCityByName(city.Name)
		if err == nil && existing.ID != 0 {
			return fmt.Errorf("city with name '%s' already exists", city.Name)
		}

		existingCity.Name = city.Name
	}

	if city.BaseDeliveryFee != existingCity.BaseDeliveryFee {
		existingCity.BaseDeliveryFee = city.BaseDeliveryFee
	}

	err = cs.cityRepository.UpdateCity(existingCity)
	if err != nil {
		return err
	}

	return nil
}

func (cs cityService) DeleteCity(id int64) error {
	return cs.cityRepository.DeleteCity(id)
}
