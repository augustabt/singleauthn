package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/augustabt/SingleAuthN/helpers"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/sessions"
)

func StartRegistration(webAuthn *webauthn.WebAuthn, store *sessions.CookieStore, db *storm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validUser := helpers.GetValidUser(db, true)

		registrationOptions := func(creationOptions *protocol.PublicKeyCredentialCreationOptions) {
			creationOptions.CredentialExcludeList = validUser.CredentialExclusionList()
		}

		options, sessionData, err := webAuthn.BeginRegistration(validUser, registrationOptions)
		if err != nil {
			log.Println("Error calling webAuthn.BeginRegistration:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonSessionData, err := json.Marshal(sessionData)
		if err != nil {
			log.Println("Error marshaling sessionData:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := store.Get(r, "webauthn-session")
		if err != nil {
			log.Println("Error getting the webauthn-session:", err)
			http.SetCookie(w, &http.Cookie{Name: "webauthn-session", MaxAge: -1})
			helpers.SendJsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		session.Values["registration"] = jsonSessionData
		session.Options.MaxAge = 120 // 2 Minutes to register
		session.Options.Path = "/registration"
		err = session.Save(r, w)
		if err != nil {
			log.Println("Error saving jsonSessionData to the session:", err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		helpers.SendJsonResponse(w, options, http.StatusOK)
	}
}
