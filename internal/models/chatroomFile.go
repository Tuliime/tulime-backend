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

// TODO: add pagination for all select queries that return many results
func (crf *ChatroomFile) FindAll(limit float64) ([]ChatroomFile, error) {
	var chatroomFile []ChatroomFile
	db.Limit(int(limit)).Find(&chatroomFile)

	// TODO: to fetch file and mention for each chat
	// for i, agroproduct := range chatRoom {
	// 	var Chatroom []Chatroom
	// 	db.Order("\"createdAt\" desc").Limit(1).Find(&agroProductPrice, "\"agroproductID\" = ?", agroproduct.ID)
	// 	agroproduct.Price = agroProductPrice
	// 	chatRoom[i] = agroproduct
	// }

	return chatroomFile, nil
}

// Update updates one Agroproduct in the database, using the information
// stored in the receiver u
func (crf *ChatroomFile) Update() (ChatroomFile, error) {
	db.Save(&crf)

	chatRoomFile, err := crf.FindOne(crf.ID)
	if err != nil {
		return chatRoomFile, err
	}

	return chatRoomFile, nil
}

func (crf *ChatroomFile) Delete(id string) error {
	//   TODO: to softly delete chatroom together with its file and mention
	return nil
}
