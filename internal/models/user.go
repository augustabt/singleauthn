package models

import (
	"encoding/binary"

	"gorm.io/gorm"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// Struct that implements the duo-labs/webauthn user interface to store valid credentials
type User struct {
	gorm.Model
	Username    string `gorm:"primaryKey"`
	IsAdmin     bool
	Credentials []Credential
}

func (u User) WebAuthnID() []byte {
	webauthnID := make([]byte, 8)
	binary.LittleEndian.PutUint64(webauthnID, uint64(u.ID))

	return webauthnID
}

func (u User) WebAuthnName() string {
	return u.Username
}

func (u User) WebAuthnDisplayName() string {
	return u.Username
}

func (u User) WebAuthnIcon() string {
	return ""
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	credentials := make([]webauthn.Credential, len(u.Credentials))

	for _, credential := range u.Credentials {
		credentials = append(credentials, credential.WebauthnCredential)
	}

	return credentials
}

func (u User) CredentialExclusionList() []protocol.CredentialDescriptor {
	credentials := u.WebAuthnCredentials()
	exclusionList := []protocol.CredentialDescriptor{}

	for _, credential := range credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: credential.ID,
		}

		exclusionList = append(exclusionList, descriptor)
	}

	return exclusionList
}
