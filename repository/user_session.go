package repository

import (
	"errors"
	"fmt"
	"oms/domain"
	"oms/model"
	"time"

	"gorm.io/gorm"
)

type userSessionRepository struct {
	masterDb  *gorm.DB
	replicaDb *gorm.DB
}

func NewUserSessionRepository(masterDB, replicaDB *gorm.DB) domain.UserSessionRepository {
	return &userSessionRepository{
		masterDb:  masterDB,
		replicaDb: replicaDB,
	}
}

func (r *userSessionRepository) CreateUserSession(session model.UserSession) error {
	return r.masterDb.Create(&session).Error
}

func (r *userSessionRepository) GetUserSessionByAccessToken(accessToken string) (model.UserSession, error) {
	var session model.UserSession
	err := r.replicaDb.Where("access_token = ?", accessToken).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.UserSession{}, fmt.Errorf("user session not found")
		}
		return model.UserSession{}, err
	}

	return session, nil
}

func (r *userSessionRepository) DeleteExpiredSessions() error {
	result := r.masterDb.Where("expires_at <= ?", time.Now()).Delete(&model.UserSession{})
	return result.Error
}

func (r *userSessionRepository) InvalidateSession(accessToken string) error {
	result := r.masterDb.Where("access_token = ?", accessToken).Delete(&model.UserSession{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}
