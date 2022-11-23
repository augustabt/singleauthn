package routes

import (
	"log"
	"net/http"
	"os"
)

func ServeLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginPage, err := os.ReadFile("../../html/login.html")
		if err != nil {
			log.Fatal("Unable to read login page:", err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write(loginPage)
	}
}
