package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/augustabt/SingleAuthN/helpers"
	"github.com/augustabt/SingleAuthN/models"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/sessions"
)

func FinishRegistration(webAuthn *webauthn.WebAuthn, store *sessions.CookieStore, db *storm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validUser := &models.ValidUser{}
		err := db.Get("user", "valid", validUser)
		if err != nil {
			log.Println("Invalid Challenge")
			helpers.SendJsonResponse(w, "Invalid Challenge", http.StatusBadRequest)
			return
		}

		session, err := store.Get(r, "webauthn-session")
		if err != nil {
			log.Println(err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonSessionData, ok := session.Values["registration"].([]byte)
		if !ok {
			log.Println("Invalid Challenge")
			helpers.SendJsonResponse(w, "Invalid Challenge", http.StatusBadRequest)
			return
		}
		sessionData := webauthn.SessionData{}
		err = json.Unmarshal(jsonSessionData, &sessionData)
		if err != nil {
			log.Println(err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "webauthn-session", MaxAge: -1, Path: "/"})

		credential, err := webAuthn.FinishRegistration(validUser, sessionData, r)
		if err != nil {
			log.Println("Error creating new authenticator")
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		validUser.AddCredentials(*credential)
		db.Set("user", "valid", validUser)
		w.WriteHeader(http.StatusOK)
	}
}
