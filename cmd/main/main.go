package main

import (
	"log"

	"github.com/asdine/storm/v3"
)

func main() {
	db, err := storm.Open("./data/storage.db")
	if err != nil {
		log.Fatal("Error opening or creating the database file")
	}
	defer db.Close()

}
