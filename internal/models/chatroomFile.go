package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (crf *ChatroomFile) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (crf *ChatroomFile) Create(chatroomFile ChatroomFile) (ChatroomFile, error) {
	result := db.Create(&chatroomFile)

	if result.Error != nil {
		return chatroomFile, result.Error
	}
	return chatroomFile, nil
}

func (crf *ChatroomFile) FindOne(id string) (ChatroomFile, error) {
	var chatroomFile ChatroomFile
	db.First(&chatroomFile, "id = ?", id)

	return chatroomFile, nil
}

func (crf *ChatroomFile) FindAll(limit float64) ([]ChatroomFile, error) {
	var chatroomFiles []ChatroomFile
	db.Limit(int(limit)).Find(&chatroomFiles)

	return chatroomFiles, nil
}

// Update updates one Agroproduct in the database, using the information
// stored in the receiver u
func (crf *ChatroomFile) Update() (ChatroomFile, error) {
	db.Save(&crf)

	return *crf, nil
}

func (crf *ChatroomFile) Delete(id string) error {
	//   TODO: to softly delete chatroom together with its file and mention
	return nil
}
