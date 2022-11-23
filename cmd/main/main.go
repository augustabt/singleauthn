package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/asdine/storm/v3"
	"github.com/augustabt/SingleAuthN/helpers"
	"github.com/augustabt/SingleAuthN/routes"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	origin, rpid := helpers.ParseDomain()

	registrationMode := os.Getenv("REGISTRATION")

	// Opening the database
	db, err := storm.Open("../../data/storage.db")
	if err != nil {
		log.Fatal("Error opening or creating the database file:", err)
	}
	defer db.Close()

	// Signal handler for saving the database when the program is terminated manually
	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigint
		fmt.Println()
		log.Println("Saving database and exiting")
		db.Close()
		os.Exit(0)
	}()

	// Creating the webauthn object from duo-labs
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Change this later",
		RPID:          rpid,
		RPOrigin:      origin,
	})
	if err != nil {
		log.Fatal("Failed to create WebAuthn based on the provided config:", err)
	}

	// Setup a session store for WebAuthN registration and auth sessions
	authKey, encryptionKey := helpers.GetSessionKeys(db)
	store := sessions.NewCookieStore(authKey, encryptionKey)

	router := mux.NewRouter()

	if registrationMode == "TRUE" {
		router.HandleFunc("/registration", routes.ServeRegistration(origin)).Methods(http.MethodGet)
		router.HandleFunc("/registration/start", routes.StartRegistration(webAuthn, store, db)).Methods("GET")
		router.HandleFunc("/registration/finish", routes.FinishRegistration(webAuthn, store, db)).Methods("POST")
	}

	router.HandleFunc("/login", routes.ServeLogin()).Methods(http.MethodGet)
	router.HandleFunc("/login/start", routes.StartLogin(webAuthn, store, db)).Methods(http.MethodGet)
	router.HandleFunc("/login/finish", routes.FinishLogin(webAuthn, rpid, store, db)).Methods(http.MethodPost)
	router.HandleFunc("/auth", routes.Auth(webAuthn, origin, store, db)).Methods(http.MethodGet)
	router.HandleFunc("/forbidden", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusForbidden) }).Methods(http.MethodGet)

	serverAddress := ":7633"
	log.Println("Starting server listening on port", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
