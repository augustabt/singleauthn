package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/augustabt/singleauthn/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func getDataPath() string {
	// Running in docker
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "/data"
	}

	return "./data"
}

func main() {
	// origin, rpid := helpers.ParseDomain()

	// Opening the database
	db, err := gorm.Open(sqlite.Open(getDataPath()+"/storage.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening or creating the database file:", err)
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Credential{})

	// Signal handler for saving the database when the program is terminated manually
	sigint := make(chan os.Signal, 2)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigint
		fmt.Println()
		log.Println("Received interrupt signal, exiting")
		os.Exit(0)
	}()
}
