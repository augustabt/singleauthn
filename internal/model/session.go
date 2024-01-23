package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ExpirationTime time.Time `gorm:"serializer:json"`
	SessionKey     string    `gorm:"primaryKey"`
}

func GenerateSession(validityLength time.Duration) Session {
	return Session{
		ExpirationTime: time.Now().Add(time.Second * validityLength),
		SessionKey:     uuid.NewString(),
	}
}

func (session Session) ValidateSession(db *gorm.DB) bool {
	storedSession := Session{}
	err := db.First(&storedSession, "id = ?", session.SessionKey)

	if err != nil || storedSession.SessionKey != session.SessionKey || storedSession.ExpirationTime != session.ExpirationTime || session.IsExpired() {
		return false
	}

	return true
}

func (session Session) IsExpired() bool {
	return session.ExpirationTime.Before(time.Now())
}
