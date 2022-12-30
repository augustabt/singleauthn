package routes

import (
	"fmt"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/augustabt/SingleAuthN/helpers"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/sessions"
)

func Auth(webAuthn *webauthn.WebAuthn, origin string, store *sessions.CookieStore, db *storm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := helpers.GetAuthSession(r, store)
		if err != nil {
			redirectToLogin(w, r, origin)
			return
		}

		if !session.ValidateSession(db) {
			redirectToLogin(w, r, origin)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}

	}
}

func redirectToLogin(w http.ResponseWriter, r *http.Request, origin string) {
	r.Header.Set("Cache-Control", "no-store")
	http.Redirect(w, r, fmt.Sprintf("%s/login?rd=%s://%s%s", origin, r.Header.Get("x-forwarded-proto"), r.Header.Get("x-forwarded-host"), r.Header.Get("x-forwarded-uri")), http.StatusSeeOther)
}
