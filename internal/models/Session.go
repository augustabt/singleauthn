package models

import (
	"time"

	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
)

type Session struct {
	ExpirationTime time.Time `json:"expirationTime"`
	SessionKey     string    `json:"sessionKey"`
}

func GenerateNewSession(validTime time.Duration) Session {
	return Session{
		ExpirationTime: time.Now().Add(time.Second * validTime),
		SessionKey:     uuid.NewString(),
	}
}

func (session Session) SaveSession(db *storm.DB) error {
	return db.Set("auth-sessions", session.SessionKey, &session)
}

func (session Session) ValidateSession(db *storm.DB) bool {
	storedSession := Session{}
	err := db.Get("auth-sessions", session.SessionKey, &storedSession)

	if err != nil || storedSession.SessionKey != session.SessionKey || storedSession.ExpirationTime != session.ExpirationTime || session.IsExpired() {
		return false
	}

	return true
}

func (session Session) IsExpired() bool {
	return session.ExpirationTime.Before(time.Now())
}
