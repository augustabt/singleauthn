package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/augustabt/SingleAuthN/helpers"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/sessions"
)

func FinishRegistration(webAuthn *webauthn.WebAuthn, store *sessions.CookieStore, db *storm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validUser := helpers.GetValidUser(db, false)
		if validUser == nil {
			log.Println("Invalid Challenge, no valid user")
			helpers.SendJsonResponse(w, "Invalid Challenge", http.StatusBadRequest)
			return
		}

		session, err := store.Get(r, "webauthn-session")
		if err != nil {
			log.Println("Error getting the webauthn-session:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonSessionData, ok := session.Values["registration"].([]byte)
		if !ok {
			log.Println("Invalid Challenge, unable to retrieve jsonSessionData")
			helpers.SendJsonResponse(w, "Invalid Challenge", http.StatusBadRequest)
			return
		}
		sessionData := webauthn.SessionData{}
		err = json.Unmarshal(jsonSessionData, &sessionData)
		if err != nil {
			log.Println("Error un-marshaling jsonSessionData:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "webauthn-session", MaxAge: -1})

		credential, err := webAuthn.FinishRegistration(validUser, sessionData, r)
		if err != nil {
			log.Println("Error creating new authenticator:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		validUser.AddCredentials(*credential)
		db.Set("user", "valid", validUser)
		w.WriteHeader(http.StatusOK)
	}
}
