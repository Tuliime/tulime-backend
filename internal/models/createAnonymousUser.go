package models

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/constants"
	"github.com/Tuliime/tulime-backend/internal/packages"
)

func CreateAnonymousUser() {
	user := User{Name: "anonymous", Role: "anonymous",
		TelNumber: constants.AnonymousTelNumber,
		Password:  packages.NewRandomNumber().D10()}

	userSaved, err := user.FindByTelNumber(user.TelNumber)
	if err != nil {
		log.Printf("Error finding anonymous user: %v", err)
		return
	}

	if userSaved.ID != "" {
		log.Printf("Anonymous user already exists: %v", userSaved)
		return
	}

	newUser, err := user.Create(user)
	if err != nil {
		log.Printf("Error creating anonymous user: %v", err)
		return
	}
	log.Printf("New Anonymous User Created Successfully: %v", newUser)
}

func init() {
	CreateAnonymousUser()
}
