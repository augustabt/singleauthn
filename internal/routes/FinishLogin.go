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

func FinishLogin(webAuthn *webauthn.WebAuthn, rpid string, store *sessions.CookieStore, db *storm.DB) http.HandlerFunc {
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

		cred, err := webAuthn.FinishLogin(validUser, sessionData, r)
		if err != nil {
			log.Println("Error logging in:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if cred.Authenticator.CloneWarning {
			log.Println("Cloned authenticator tried to use for authentication")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		validUser.UpdateCredentialSignCount(cred.ID, cred.Authenticator.SignCount)
		db.Set("user", "valid", validUser)

		// Create a new session that is valid for 12 hours
		newSession := models.GenerateNewSession(43200)
		jsonSession, err := json.Marshal(newSession)
		if err != nil {
			log.Println("Error marshaling session:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionStore, err := store.Get(r, "auth-session")
		if err != nil {
			log.Println("Error getting auth session cookie:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		sessionStore.Options.MaxAge = 43200
		sessionStore.Options.Domain = rpid
		sessionStore.Values["session"] = jsonSession
		sessionStore.Save(r, w)
		if err != nil {
			log.Println("Error saving jsonSession to the session:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = newSession.SaveSession(db)
		if err != nil {
			log.Println("Error saving session to db:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
