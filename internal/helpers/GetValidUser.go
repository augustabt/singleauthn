package helpers

import (
	"github.com/asdine/storm/v3"
	"github.com/augustabt/SingleAuthN/models"
)

func GetValidUser(db *storm.DB, generate bool) *models.ValidUser {
	validUser := &models.ValidUser{}
	err := db.Get("user", "valid", validUser)

	// If this is the first time this function has run, create and save a user
	if err != nil {
		if generate {
			validUser = models.GenerateValidUser()
			db.Set("user", "valid", validUser)
		} else {
			return nil
		}
	}

	return validUser
}
