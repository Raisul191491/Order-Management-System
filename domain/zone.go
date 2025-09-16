package domain

import (
	"oms/model"
	"oms/types"
)

type ZoneRepository interface {
	CreateZone(zone model.Zone) error
	GetZoneByID(id int64) (model.Zone, error)
	GetAllZones(limit, offset int) ([]model.Zone, error)
	UpdateZone(zone model.Zone) error
	DeleteZone(id int64) error
	GetZoneByName(name string) (model.Zone, error)
	GetZonesByCityID(cityID int64, limit, offset int) ([]model.Zone, error)
	GetZoneByNameAndCityID(name string, cityID int64) (model.Zone, error)
	CountZonesByCity(cityID int64) (int64, error)
}

type ZoneService interface {
	CreateZone(zone types.ZoneCreateRequest) error
	GetZoneByID(id int64) (types.ZoneResponse, error)
	GetAllZones(limit, offset int) ([]types.ZoneResponse, error)
	GetZonesByCityID(cityID int64, limit, offset int) ([]types.ZoneResponse, error)
	UpdateZone(zone types.ZoneUpdateRequest) error
	DeleteZone(id int64) error
}
