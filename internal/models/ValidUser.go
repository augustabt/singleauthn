package models

import (
	"bytes"
	"crypto/rand"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

// Struct that implements the duo-labs/webauthn user interface to store valid credentials
type ValidUser struct {
	ID          []byte `storm:"id"`
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
}

// Creates and returns a default struct for when there is not already a database file
func GenerateValidUser() *ValidUser {
	validUser := &ValidUser{}

	validUser.ID = make([]byte, 8)
	rand.Read(validUser.ID)

	validUser.Name = "ValidUser"
	validUser.DisplayName = "ValidUser"
	validUser.Credentials = []webauthn.Credential{}

	return validUser
}

func (valid ValidUser) WebAuthnID() []byte {
	return valid.ID
}

func (valid ValidUser) WebAuthnName() string {
	return valid.Name
}

func (valid ValidUser) WebAuthnDisplayName() string {
	return valid.DisplayName
}

func (valid ValidUser) WebAuthnIcon() string {
	return ""
}

func (valid ValidUser) WebAuthnCredentials() []webauthn.Credential {
	return valid.Credentials
}

func (valid *ValidUser) AddCredentials(newCredential webauthn.Credential) {
	valid.Credentials = append(valid.Credentials, newCredential)
}

func (valid *ValidUser) UpdateCredentialSignCount(credID []byte, signCount uint32) {
	for i, cred := range valid.Credentials {
		if bytes.Equal(cred.ID, credID) {
			cred.Authenticator.SignCount = signCount
			valid.Credentials[i] = cred
		}
	}
}

func (valid ValidUser) CredentialExclusionList() []protocol.CredentialDescriptor {
	exclusionList := []protocol.CredentialDescriptor{}

	for _, credential := range valid.Credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: credential.ID,
		}

		exclusionList = append(exclusionList, descriptor)
	}

	return exclusionList
}
