package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (msgrt *MessengerTag) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (msgrf *MessengerTag) Create(messengerTag MessengerTag) (MessengerTag, error) {
	result := db.Create(&messengerTag)

	if result.Error != nil {
		return messengerTag, result.Error
	}
	return messengerTag, nil
}

func (msgrf *MessengerTag) FindOne(id string) (MessengerTag, error) {
	var messengerTag MessengerTag
	db.First(&messengerTag, "id = ?", id)

	return messengerTag, nil
}

// TODO: add pagination for all select queries that return many results
func (msgrf *MessengerTag) FindAll(limit float64) ([]MessengerTag, error) {
	var messengerTag []MessengerTag
	db.Limit(int(limit)).Find(&messengerTag)

	// TODO: to include full details here

	return messengerTag, nil
}

// Update updates one MessengerTag in the database, using the information
// stored in the receiver u
func (msgrf *MessengerTag) Update() (MessengerTag, error) {
	db.Save(&msgrf)

	messengerTag, err := msgrf.FindOne(msgrf.ID)
	if err != nil {
		return messengerTag, err
	}

	return messengerTag, nil
}

func (msgrf *MessengerTag) Delete(id string) error {
	//   TODO: to softly delete MessengerTag together
	return nil
}
