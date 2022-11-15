package routes

import (
	"log"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/augustabt/SingleAuthN/helpers"
	"github.com/augustabt/SingleAuthN/models"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

func BeginRegistration(webAuthn *webauthn.WebAuthn, db *storm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validUser := &models.ValidUser{}
		db.Get("user", "valid", validUser)

		// If this is the first time this function has run, create and save a user
		if validUser == nil {
			validUser = models.GenerateValidUser()
			db.Set("user", "valid", validUser)
		}

		registrationOptions := func(creationOptions *protocol.PublicKeyCredentialCreationOptions) {
			creationOptions.CredentialExcludeList = validUser.CredentialExclusionList()
		}

		options, sessionData, err := webAuthn.BeginRegistration(validUser, registrationOptions)
		if err != nil {
			log.Println(err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db.Set("registrationSessions", sessionData.Challenge, sessionData)
		if err != nil {
			log.Println(err)
			helpers.SendJsonResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		helpers.SendJsonResponse(w, options, http.StatusOK)
	}
}
