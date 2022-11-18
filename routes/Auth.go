package routes

import (
	"fmt"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/augustabt/SingleAuthN/helpers"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/sessions"
)

func Auth(webAuthn *webauthn.WebAuthn, store *sessions.CookieStore, db *storm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := helpers.GetAuthSession(r, store)
		if err != nil {
			redirectToLogin(w, r)
			return
		}

		if !session.ValidateSession(db) {
			redirectToLogin(w, r)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}

	}
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://valid.augustabt.com/login.html?rd=%s://%s%s", r.Header.Get("x-forwarded-proto"), r.Header.Get("x-forwarded-host"), r.Header.Get("x-forwarded-uri")), http.StatusSeeOther)
}
