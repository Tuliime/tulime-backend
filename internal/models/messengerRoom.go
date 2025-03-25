package models

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (msgrr *MessengerRoom) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

// Creates new MessengerRoom, determine order of UserOneID and
// UserTwoID based on the user's createdAt.
// Older user takes UserOneID and young user takes UserTwoID
func (msgrr *MessengerRoom) Create(messengerRoom MessengerRoom) (MessengerRoom, error) {
	user := User{}
	userOne, err := user.FindOne(messengerRoom.UserOneID)
	if err != nil {
		return messengerRoom, err
	}
	userTwo, err := user.FindOne(messengerRoom.UserTwoID)
	if err != nil {
		return messengerRoom, err
	}

	if userOne.CreatedAt.Before(userTwo.CreatedAt) {
		messengerRoom.UserOneID = userOne.ID
		messengerRoom.UserTwoID = userTwo.ID
	} else {
		messengerRoom.UserOneID = userTwo.ID
		messengerRoom.UserTwoID = userOne.ID
	}

	if result := db.Create(&messengerRoom); result.Error != nil {
		return messengerRoom, result.Error
	}
	return messengerRoom, nil
}

func (msgrr *MessengerRoom) FindOne(id string) (MessengerRoom, error) {
	var messengerRoom MessengerRoom
	db.First(&messengerRoom, "id = ?", id)

	return messengerRoom, nil
}

func (msgrr *MessengerRoom) FindByUsers(userOneID, userTwoID string) (MessengerRoom, error) {
	var messengerRoom MessengerRoom

	db.Where("\"userOneID\" = ? AND \"userTwoID\" = ?", userOneID, userTwoID).First(&messengerRoom)

	if messengerRoom.ID == "" {
		log.Println("using second query to get messengerRoom")
		db.Where("\"userOneID\" = ? AND \"userTwoID\" = ?", userTwoID, userOneID).First(&messengerRoom)
	}
	return messengerRoom, nil
}

// TODO: add pagination for all select queries that return many results
func (msgrr *MessengerRoom) FindAll(limit float64) ([]MessengerRoom, error) {
	var messengerRoom []MessengerRoom
	db.Limit(int(limit)).Find(&messengerRoom)

	// TODO: to include full details here

	return messengerRoom, nil
}

// Update updates one MessengerRoom in the database, using the information
// stored in the receiver u
func (msgrr *MessengerRoom) Update() (MessengerRoom, error) {
	db.Save(&msgrr)

	messengerRoom, err := msgrr.FindOne(msgrr.ID)
	if err != nil {
		return messengerRoom, err
	}

	return messengerRoom, nil
}

func (msgrr *MessengerRoom) Delete(id string) error {
	//   TODO: to softly delete MessengerRoom together
	return nil
}
