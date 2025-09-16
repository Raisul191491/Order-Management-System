package service

import (
	"fmt"
	"oms/domain"
	"oms/model"
	"oms/types"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) domain.UserService {
	return &userService{userRepository: userRepository}
}

func (us userService) CreateUser(user types.UserCreateRequest) error {
	// Normalize email (trim spaces and convert to lowercase)
	normalizedEmail := strings.ToLower(strings.TrimSpace(user.Email))
	if normalizedEmail == "" {
		return fmt.Errorf("email cannot be empty")
	}

	// Check if user with same email already exists
	existing, err := us.userRepository.GetUserByEmail(normalizedEmail)
	if err == nil && existing.ID != 0 {
		return fmt.Errorf("user with email '%s' already exists", normalizedEmail)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := model.User{
		Email:        normalizedEmail,
		PasswordHash: string(hashedPassword),
	}

	err = us.userRepository.CreateUser(newUser)
	if err != nil {
		return err
	}

	return nil
}

func (us userService) GetUserByID(id int64) (types.UserResponse, error) {
	existingUser, err := us.userRepository.GetUserByID(id)
	if err != nil {
		return types.UserResponse{}, err
	}

	return types.UserResponse{
		ID:        existingUser.ID,
		Email:     existingUser.Email,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt,
	}, nil
}

func (us userService) GetAllUsers(limit, offset int) ([]types.UserResponse, error) {
	existingUsers, err := us.userRepository.GetAllUsers(limit, offset)
	if err != nil {
		return nil, err
	}

	var result []types.UserResponse

	for _, existingUser := range existingUsers {
		result = append(result, types.UserResponse{
			ID:        existingUser.ID,
			Email:     existingUser.Email,
			CreatedAt: existingUser.CreatedAt,
			UpdatedAt: existingUser.UpdatedAt,
		})
	}

	return result, nil
}

func (us userService) UpdateUserEmail(user types.UserUpdateRequest) error {
	existingUser, err := us.userRepository.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(user.Email))
	if normalizedEmail == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if normalizedEmail != existingUser.Email {
		existing, err := us.userRepository.GetUserByEmail(normalizedEmail)
		if err == nil && existing.ID != 0 && existing.ID != existingUser.ID {
			return fmt.Errorf("user with email '%s' already exists", normalizedEmail)
		}
		existingUser.Email = normalizedEmail
	}

	isVerified := us.VerifyUserCredentials(user.Email, user.Password)
	if !isVerified {
		return fmt.Errorf("user with email '%s' has not been verified", user.Email)
	}

	err = us.userRepository.UpdateUserEmail(existingUser)
	if err != nil {
		return err
	}

	return nil
}

func (us userService) DeleteUser(id int64) error {
	existingUser, err := us.userRepository.GetUserByID(id)
	if err != nil || existingUser.ID == 0 {
		return fmt.Errorf("user does not exist")
	}

	return us.userRepository.DeleteUser(id)
}

func (us userService) GetUserByEmail(email string) (types.UserResponse, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	existingUser, err := us.userRepository.GetUserByEmail(normalizedEmail)
	if err != nil {
		return types.UserResponse{}, err
	}

	return types.UserResponse{
		ID:        existingUser.ID,
		Email:     existingUser.Email,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt,
	}, nil
}

func (us userService) VerifyUserCredentials(email, password string) bool {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	existingUser, err := us.userRepository.GetUserByEmail(normalizedEmail)
	if err != nil {
		return false
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
