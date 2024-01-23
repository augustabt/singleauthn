package models

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"strconv"

	"gorm.io/gorm"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// Struct that implements the webauthn user interface to store valid credentials
type User struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Username    string
	IsAdmin     bool
	Credentials []Credential
}

func (u User) WebAuthnID() []byte {
	userID, err := strconv.ParseUint(u.ID, 10, 64)
	if err != nil {
		log.Fatal("Unable to parse userID", err)
	}

	webauthnID := make([]byte, 8)
	binary.LittleEndian.PutUint64(webauthnID, userID)

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
	exclusionList := make([]protocol.CredentialDescriptor, len(credentials))

	for _, credential := range credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: credential.ID,
		}

		exclusionList = append(exclusionList, descriptor)
	}

	return exclusionList
}

func GenerateUserID() string {
	id := make([]byte, 8)
	n, err := rand.Read(id)

	if (err != nil) || (n != 8) {
		log.Fatal("Unable to generate a random user id", err)
	}

	userID := binary.LittleEndian.Uint64(id)

	return strconv.FormatUint(userID, 10)
}
