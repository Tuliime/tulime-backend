package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (msgrf *MessengerFile) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (msgrf *MessengerFile) Create(messengerFile MessengerFile) (MessengerFile, error) {
	result := db.Create(&messengerFile)

	if result.Error != nil {
		return messengerFile, result.Error
	}
	return messengerFile, nil
}

func (msgrf *MessengerFile) FindOne(id string) (MessengerFile, error) {
	var messengerFile MessengerFile
	db.First(&messengerFile, "id = ?", id)

	return messengerFile, nil
}

// TODO: add pagination for all select queries that return many results
func (msgrf *MessengerFile) FindAll(limit float64) ([]MessengerFile, error) {
	var messengerFile []MessengerFile
	db.Limit(int(limit)).Find(&messengerFile)

	return messengerFile, nil
}

// Update updates one MessengerTag in the database, using the information
// stored in the receiver u
func (msgrf *MessengerFile) Update() (MessengerFile, error) {
	db.Save(&msgrf)

	messengerFile, err := msgrf.FindOne(msgrf.ID)
	if err != nil {
		return messengerFile, err
	}

	return messengerFile, nil
}

func (msgrf *MessengerFile) Delete(id string) error {
	//   TODO: to softly delete MessengerFile
	return nil
}
