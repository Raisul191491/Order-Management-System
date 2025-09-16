package service

import (
	"fmt"
	"oms/domain"
	"oms/model"
	"oms/types"
	"strings"
)

type itemTypeService struct {
	itemTypeRepository domain.ItemTypeRepository
}

func NewItemTypeService(itemTypeRepository domain.ItemTypeRepository) domain.ItemTypeService {
	return &itemTypeService{itemTypeRepository: itemTypeRepository}
}

func (its itemTypeService) CreateItemType(itemType types.ItemTypeCreateRequest) error {
	// Normalize name (trim spaces and convert to proper case)
	normalizedName := strings.TrimSpace(itemType.Name)
	if normalizedName == "" {
		return fmt.Errorf("item type name cannot be empty")
	}

	// Check if item type with same name already exists (case-insensitive)
	existing, err := its.itemTypeRepository.GetItemTypeByName(normalizedName)
	if err == nil && existing.ID != 0 {
		return fmt.Errorf("item type with name '%s' already exists", normalizedName)
	}

	newItemType := model.ItemType{
		Name: normalizedName,
	}

	err = its.itemTypeRepository.CreateItemType(newItemType)
	if err != nil {
		return err
	}

	return nil
}

func (its itemTypeService) GetItemTypeByID(id int64) (types.ItemTypeResponse, error) {
	existingItemType, err := its.itemTypeRepository.GetItemTypeByID(id)
	if err != nil {
		return types.ItemTypeResponse{}, err
	}

	return types.ItemTypeResponse{
		ID:        existingItemType.ID,
		Name:      existingItemType.Name,
		CreatedAt: existingItemType.CreatedAt,
		UpdatedAt: existingItemType.UpdatedAt,
	}, nil
}

func (its itemTypeService) GetAllItemTypes(limit, offset int) ([]types.ItemTypeResponse, error) {
	existingItemTypes, err := its.itemTypeRepository.GetAllItemTypes(limit, offset)
	if err != nil {
		return nil, err
	}

	var result []types.ItemTypeResponse

	for _, existingItemType := range existingItemTypes {
		result = append(result, types.ItemTypeResponse{
			ID:        existingItemType.ID,
			Name:      existingItemType.Name,
			CreatedAt: existingItemType.CreatedAt,
			UpdatedAt: existingItemType.UpdatedAt,
		})
	}

	return result, nil
}

func (its itemTypeService) UpdateItemType(itemType types.ItemTypeUpdateRequest) error {
	existingItemType, err := its.itemTypeRepository.GetItemTypeByID(itemType.ID)
	if err != nil {
		return err
	}

	// Normalize name
	normalizedName := strings.TrimSpace(itemType.Name)
	if normalizedName == "" {
		return fmt.Errorf("item type name cannot be empty")
	}

	// If name is being changed, check for duplicates
	if normalizedName != existingItemType.Name {
		existing, err := its.itemTypeRepository.GetItemTypeByName(normalizedName)
		if err == nil && existing.ID != 0 && existing.ID != existingItemType.ID {
			return fmt.Errorf("item type with name '%s' already exists", normalizedName)
		}

		existingItemType.Name = normalizedName
	}

	err = its.itemTypeRepository.UpdateItemType(existingItemType)
	if err != nil {
		return err
	}

	return nil
}

func (its itemTypeService) DeleteItemType(id int64) error {
	existingItemType, err := its.itemTypeRepository.GetItemTypeByID(id)
	if err != nil || existingItemType.ID == 0 {
		return fmt.Errorf("item type does not exist")
	}

	return its.itemTypeRepository.DeleteItemType(id)
}
