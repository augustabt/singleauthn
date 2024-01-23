package model

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model
	Name               string
	WebauthnCredential webauthn.Credential `gorm:"serializer:json"`
	UserID             string
}

func (c *Credential) UpdateSignCount(db *gorm.DB, signCount uint32) {
	c.WebauthnCredential.Authenticator.SignCount = signCount
	db.Save(c)
}
