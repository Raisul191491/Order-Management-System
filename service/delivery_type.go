package service

import (
	"fmt"
	"oms/domain"
	"oms/model"
	"oms/types"
	"strings"
)

type deliveryTypeService struct {
	deliveryTypeRepository domain.DeliveryTypeRepository
}

func NewDeliveryTypeService(deliveryTypeRepository domain.DeliveryTypeRepository) domain.DeliveryTypeService {
	return &deliveryTypeService{deliveryTypeRepository: deliveryTypeRepository}
}

func (dts deliveryTypeService) CreateDeliveryType(deliveryType types.DeliveryTypeCreateRequest) error {
	// Normalize name (trim spaces)
	normalizedName := strings.TrimSpace(deliveryType.Name)
	if normalizedName == "" {
		return fmt.Errorf("delivery type name cannot be empty")
	}

	// Check if delivery type with same name already exists
	existing, err := dts.deliveryTypeRepository.GetDeliveryTypeByName(normalizedName)
	if err == nil && existing.ID != 0 {
		return fmt.Errorf("delivery type with name '%s' already exists", normalizedName)
	}

	newDeliveryType := model.DeliveryType{
		Name: normalizedName,
	}

	err = dts.deliveryTypeRepository.CreateDeliveryType(newDeliveryType)
	if err != nil {
		return err
	}

	return nil
}

func (dts deliveryTypeService) GetDeliveryTypeByID(id int64) (types.DeliveryTypeResponse, error) {
	existingDeliveryType, err := dts.deliveryTypeRepository.GetDeliveryTypeByID(id)
	if err != nil {
		return types.DeliveryTypeResponse{}, err
	}

	return types.DeliveryTypeResponse{
		ID:        existingDeliveryType.ID,
		Name:      existingDeliveryType.Name,
		CreatedAt: existingDeliveryType.CreatedAt,
		UpdatedAt: existingDeliveryType.UpdatedAt,
	}, nil
}

func (dts deliveryTypeService) GetAllDeliveryTypes(limit, offset int) ([]types.DeliveryTypeResponse, error) {
	existingDeliveryTypes, err := dts.deliveryTypeRepository.GetAllDeliveryTypes(limit, offset)
	if err != nil {
		return nil, err
	}

	var result []types.DeliveryTypeResponse

	for _, existingDeliveryType := range existingDeliveryTypes {
		result = append(result, types.DeliveryTypeResponse{
			ID:        existingDeliveryType.ID,
			Name:      existingDeliveryType.Name,
			CreatedAt: existingDeliveryType.CreatedAt,
			UpdatedAt: existingDeliveryType.UpdatedAt,
		})
	}

	return result, nil
}

func (dts deliveryTypeService) UpdateDeliveryType(deliveryType types.DeliveryTypeUpdateRequest) error {
	existingDeliveryType, err := dts.deliveryTypeRepository.GetDeliveryTypeByID(deliveryType.ID)
	if err != nil {
		return err
	}

	// Normalize name
	normalizedName := strings.TrimSpace(deliveryType.Name)
	if normalizedName == "" {
		return fmt.Errorf("delivery type name cannot be empty")
	}

	// If name is being changed, check for duplicates
	if normalizedName != existingDeliveryType.Name {
		existing, err := dts.deliveryTypeRepository.GetDeliveryTypeByName(normalizedName)
		if err == nil && existing.ID != 0 && existing.ID != existingDeliveryType.ID {
			return fmt.Errorf("delivery type with name '%s' already exists", normalizedName)
		}

		existingDeliveryType.Name = normalizedName
	}

	err = dts.deliveryTypeRepository.UpdateDeliveryType(existingDeliveryType)
	if err != nil {
		return err
	}

	return nil
}

func (dts deliveryTypeService) DeleteDeliveryType(id int64) error {
	existingDeliveryType, err := dts.deliveryTypeRepository.GetDeliveryTypeByID(id)
	if err != nil || existingDeliveryType.ID == 0 {
		return fmt.Errorf("delivery type does not exist")
	}

	return dts.deliveryTypeRepository.DeleteDeliveryType(id)
}
