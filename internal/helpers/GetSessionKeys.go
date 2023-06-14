package helpers

import (
	"crypto/rand"

	"github.com/asdine/storm/v3"
)

func GetSessionKeys(db *storm.DB) ([]byte, []byte) {
	authKey := make([]byte, 64)
	encryptionKey := make([]byte, 32)

	err := db.Get("session_keys", "auth", &authKey)
	if err != nil {
		rand.Read(authKey)
		db.Set("session_keys", "auth", authKey)
	}

	err = db.Get("session_keys", "encryption", &encryptionKey)
	if err != nil {
		rand.Read(encryptionKey)
		db.Set("session_keys", "encryption", encryptionKey)
	}

	return authKey, encryptionKey
}
