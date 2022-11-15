package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/asdine/storm/v3"
	"github.com/augustabt/SingleAuthN/routes"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
)

func main() {
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
		log.Println("Saving database and exiting")
		db.Close()
		os.Exit(0)
	}()

	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Change this later",
		RPID:          "This should be the domain",
		RPOrigin:      "Origin URL for requests",
	})
	if err != nil {
		log.Fatal("Failed to create WebAuthn based on the provided config:", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/register/begin", routes.BeginRegistration(webAuthn, db)).Methods("GET")

	serverAddress := ":7633"
	log.Println("Starting server listening on port", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
