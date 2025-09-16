package model

import "time"

type MigrationRecord struct {
	Version   string    `gorm:"primaryKey;size:255"`
	AppliedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
