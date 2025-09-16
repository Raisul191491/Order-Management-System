package service

import (
	"fmt"
	"oms/domain"
	"oms/types"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository     domain.UserRepository
	userSessionService domain.UserSessionService
}

func NewAuthService(userRepository domain.UserRepository, userSessionService domain.UserSessionService) domain.AuthService {
	return &authService{
		userRepository:     userRepository,
		userSessionService: userSessionService,
	}
}

func (as authService) Login(loginRequest types.UserLoginRequest) (types.UserLoginResponse, error) {
	user, err := as.userRepository.GetUserByEmail(loginRequest.Email)
	if err != nil {
		return types.UserLoginResponse{}, fmt.Errorf("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password))
	if err != nil {
		return types.UserLoginResponse{}, fmt.Errorf("invalid email or password")
	}

	session, err := as.userSessionService.CreateUserSession(user.ID)
	if err != nil {
		return types.UserLoginResponse{}, fmt.Errorf("failed to create session")
	}

	return types.UserLoginResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		TokenType:    "Bearer",
	}, nil
}

func (as authService) Logout(accessToken string) error {
	_, err := as.userSessionService.ValidateSession(accessToken)
	if err != nil {
		return fmt.Errorf("invalid access token")
	}

	// Invalidate the session
	err = as.userSessionService.InvalidateSession(accessToken)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}
