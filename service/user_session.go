package service

import (
	"fmt"
	"oms/config"
	"oms/domain"
	"oms/model"
	"oms/types"
	"oms/utility"
	"time"
)

type userSessionService struct {
	userSessionRepository domain.UserSessionRepository
	config                config.Config
}

func NewUserSessionService(
	userSessionRepository domain.UserSessionRepository,
	config config.Config,
) domain.UserSessionService {
	return &userSessionService{
		userSessionRepository: userSessionRepository,
		config:                config,
	}
}

func (uss userSessionService) CreateUserSession(userID int64) (model.UserSession, error) {
	accessToken, err := utility.GenerateJWT(userID, uss.config.AccessTokenExpirationTime)
	if err != nil {
		return model.UserSession{}, err
	}

	refreshToken, err := utility.GenerateJWT(userID, uss.config.RefreshTokenExpirationTime)
	if err != nil {
		return model.UserSession{}, err
	}

	session := model.UserSession{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().UTC().Add(uss.config.AccessTokenExpirationTime),
	}

	err = uss.userSessionRepository.CreateUserSession(session)
	if err != nil {
		return model.UserSession{}, err
	}

	return session, nil
}

func (uss userSessionService) ValidateSession(accessToken string) (types.UserSessionResponse, error) {
	session, err := uss.userSessionRepository.GetUserSessionByAccessToken(accessToken)
	if err != nil {
		return types.UserSessionResponse{}, err
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		_ = uss.userSessionRepository.InvalidateSession(accessToken)
		return types.UserSessionResponse{}, fmt.Errorf("session expired")
	}

	return types.UserSessionResponse{
		ID:        session.ID,
		UserID:    session.UserID,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}, nil
}

func (uss userSessionService) CleanupExpiredSessions() error {
	return uss.userSessionRepository.DeleteExpiredSessions()
}

func (uss userSessionService) InvalidateSession(accessToken string) error {
	return uss.userSessionRepository.InvalidateSession(accessToken)
}
