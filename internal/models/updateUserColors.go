package models

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/packages"
)

func UpdateUserColors() {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}

	for _, user := range users {
		profileColor, err := packages.GetRandomColor()
		if err != nil {
			log.Printf("Error getting random color for user %s: %v", user.ID, err)
			continue
		}
		user.ProfileBgColor = profileColor

		chatroomColor, err := packages.GetRandomColor()
		if err != nil {
			log.Printf("Error getting random color for user %s: %v", user.ID, err)
			continue
		}
		user.ChatroomColor = chatroomColor

		if err := db.Save(&user).Error; err != nil {
			log.Printf("Error updating user %s: %v", user.ID, err)
		}
	}

	log.Println("Updated all user colors")
}

func init() {
	UpdateUserColors()
}
