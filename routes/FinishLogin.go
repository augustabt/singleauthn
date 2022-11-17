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

func FinishLogin(webAuthn *webauthn.WebAuthn, store *sessions.CookieStore, db *storm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validUser := helpers.GetValidUser(db, false)
		if validUser == nil {
			log.Println("Invalid login, no valid user")
			helpers.SendJsonResponse(w, "Invalid login", http.StatusBadRequest)
			return
		}

		session, err := store.Get(r, "webauthn-session")
		if err != nil {
			log.Println("Error getting the webauthn-session:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonSessionData, ok := session.Values["authentication"].([]byte)
		if !ok {
			log.Println("Invalid login, unable to retrieve jsonSessionData")
			helpers.SendJsonResponse(w, "Invalid login", http.StatusBadRequest)
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

		_, err = webAuthn.FinishLogin(validUser, sessionData, r)
		if err != nil {
			log.Println("Error logging in:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
