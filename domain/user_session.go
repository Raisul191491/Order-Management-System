package domain

import (
	"oms/model"
	"oms/types"
)

type UserSessionRepository interface {
	CreateUserSession(session model.UserSession) error
	GetUserSessionByAccessToken(accessToken string) (model.UserSession, error)
	DeleteExpiredSessions() error
	InvalidateSession(tokenHash string) error
}

type UserSessionService interface {
	CreateUserSession(userID int64) (model.UserSession, error)
	ValidateSession(tokenHash string) (types.UserSessionResponse, error)
	CleanupExpiredSessions() error
	InvalidateSession(tokenHash string) error
}
