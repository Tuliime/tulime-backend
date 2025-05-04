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
		if user.ProfileBgColor != "" && user.ChatroomColor != "" {
			log.Printf("User :%s has both ProfileBgColor :%s and ChatroomColor :%s",
				user.Name, user.ProfileBgColor, user.ChatroomColor)
			continue
		}

		if user.ProfileBgColor == "" {
			profileColor, err := packages.GetRandomColor()
			if err != nil {
				log.Printf("Error getting random profileColor for user %s: %v", user.ID, err)
				continue
			}
			user.ProfileBgColor = profileColor
		}

		if user.ChatroomColor == "" {
			chatroomColor, err := packages.GetRandomColor()
			if err != nil {
				log.Printf("Error getting random chatroomColor for user %s: %v", user.ID, err)
				continue
			}
			user.ChatroomColor = chatroomColor
		}

		if err := db.Save(&user).Error; err != nil {
			log.Printf("Error updating user %s: %v", user.ID, err)
		}
	}

	log.Println("Updated all user colors")
}

// func init() {
// 	UpdateUserColors()
// }
