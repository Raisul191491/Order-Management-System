package domain

import (
	"oms/model"
	"oms/types"
)

type UserRepository interface {
	CreateUser(user model.User) error
	GetUserByID(id int64) (model.User, error)
	GetAllUsers(limit, offset int) ([]model.User, error)
	UpdateUserEmail(user model.User) error
	DeleteUser(id int64) error
	GetUserByEmail(email string) (model.User, error)
}

type UserService interface {
	CreateUser(user types.UserCreateRequest) error
	GetUserByID(id int64) (types.UserResponse, error)
	GetAllUsers(limit, offset int) ([]types.UserResponse, error)
	UpdateUserEmail(user types.UserUpdateRequest) error
	DeleteUser(id int64) error
	GetUserByEmail(email string) (types.UserResponse, error)
	VerifyUserCredentials(email, password string) bool
}
