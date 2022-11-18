package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/augustabt/SingleAuthN/models"
	"github.com/gorilla/sessions"
)

func GetAuthSession(r *http.Request, store *sessions.CookieStore) (models.Session, error) {
	session := models.Session{}
	sessionStore, err := store.Get(r, "auth-session")
	if err != nil {
		return session, err
	}

	jsonSession, ok := sessionStore.Values["session"].([]byte)
	if !ok {
		return session, err
	}
	err = json.Unmarshal(jsonSession, &session)
	if err != nil {
		return session, err
	}

	return session, err
}
