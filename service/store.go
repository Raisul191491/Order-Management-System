package service

import (
	"fmt"
	"oms/domain"
	"oms/model"
	"oms/types"
)

type storeService struct {
	storeRepository domain.StoreRepository
}

func NewStoreService(storeRepository domain.StoreRepository) domain.StoreService {
	return &storeService{storeRepository: storeRepository}
}

func (ss storeService) CreateStore(store types.StoreCreateRequest) error {
	existing, err := ss.storeRepository.GetStoreByName(store.Name)
	if err == nil && existing.ID != 0 {
		return fmt.Errorf("store with name '%s' already exists", store.Name)
	}

	newStore := model.Store{
		Name:         store.Name,
		ContactPhone: store.ContactPhone,
		Address:      store.Address,
	}

	err = ss.storeRepository.CreateStore(newStore)
	if err != nil {
		return err
	}

	return nil
}

func (ss storeService) GetStoreByID(id int64) (types.StoreResponse, error) {
	existingStore, err := ss.storeRepository.GetStoreByID(id)
	if err != nil {
		return types.StoreResponse{}, err
	}

	return types.StoreResponse{
		ID:           existingStore.ID,
		Name:         existingStore.Name,
		ContactPhone: existingStore.ContactPhone,
		Address:      existingStore.Address,
		UpdatedAt:    existingStore.UpdatedAt,
	}, nil
}

func (ss storeService) GetAllStores(limit, offset int) ([]types.StoreResponse, error) {
	existingStores, err := ss.storeRepository.GetAllStores(limit, offset)
	if err != nil {
		return nil, err
	}

	var result []types.StoreResponse

	for _, existingStore := range existingStores {
		result = append(result, types.StoreResponse{
			ID:           existingStore.ID,
			Name:         existingStore.Name,
			ContactPhone: existingStore.ContactPhone,
			Address:      existingStore.Address,
			UpdatedAt:    existingStore.UpdatedAt,
		})
	}

	return result, nil
}

func (ss storeService) UpdateStore(store types.StoreUpdateRequest) error {
	existingStore, err := ss.storeRepository.GetStoreByID(store.ID)
	if err != nil {
		return err
	}

	if store.Name != existingStore.Name {
		existing, err := ss.storeRepository.GetStoreByName(store.Name)
		if err == nil && existing.ID != 0 {
			return fmt.Errorf("store with name '%s' already exists", store.Name)
		}

		existingStore.Name = store.Name
	}

	if store.ContactPhone != nil && *store.ContactPhone != existingStore.ContactPhone {
		existingStore.ContactPhone = *store.ContactPhone
	}

	if store.Address != existingStore.Address {
		existingStore.Address = store.Address
	}

	err = ss.storeRepository.UpdateStore(existingStore)
	if err != nil {
		return err
	}

	return nil
}

func (ss storeService) DeleteStore(id int64) error {
	existingStore, err := ss.storeRepository.GetStoreByID(id)
	if err != nil || existingStore.ID == 0 {
		return fmt.Errorf("store does not exist")
	}

	return ss.storeRepository.DeleteStore(id)
}
