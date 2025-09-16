package repository

import (
	"errors"
	"fmt"
	"oms/domain"
	"oms/model"

	"gorm.io/gorm"
)

type userRepository struct {
	masterDb  *gorm.DB
	replicaDb *gorm.DB
}

func NewUserRepository(masterDB, replicaDB *gorm.DB) domain.UserRepository {
	return &userRepository{
		masterDb:  masterDB,
		replicaDb: replicaDB,
	}
}

func (r *userRepository) CreateUser(user model.User) error {
	return r.masterDb.Create(&user).Error
}

func (r *userRepository) GetUserByID(id int64) (model.User, error) {
	var user model.User
	err := r.replicaDb.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, fmt.Errorf("user with ID %d not found", id)
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetAllUsers(limit, offset int) ([]model.User, error) {
	var users []model.User
	err := r.replicaDb.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) UpdateUserEmail(user model.User) error {
	result := r.masterDb.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	return nil
}

func (r *userRepository) DeleteUser(id int64) error {
	result := r.masterDb.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	err := r.replicaDb.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, fmt.Errorf("user with email '%s' not found", email)
		}
		return model.User{}, err
	}

	return user, nil
}
