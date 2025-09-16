package domain

import (
	"oms/model"
	"oms/types"
)

type CityRepository interface {
	CreateCity(city model.City) error
	GetCityByID(id int64) (model.City, error)
	GetAllCities(limit, offset int) ([]model.City, error)
	GetCityByName(name string) (model.City, error)
	UpdateCity(city model.City) error
	DeleteCity(id int64) error
}

type CityService interface {
	CreateCity(city types.CityCreateRequest) error
	GetCityByID(id int64) (types.CityResponse, error)
	GetAllCities(limit, offset int) ([]types.CityResponse, error)
	GetCityByName(name string) (types.CityResponse, error)
	UpdateCity(city types.CityUpdateRequest) error
	DeleteCity(id int64) error
}
