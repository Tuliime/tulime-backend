package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (crm *ChatroomMention) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (crm *ChatroomMention) Create(chatroomMention ChatroomMention) (ChatroomMention, error) {
	result := db.Create(&chatroomMention)

	if result.Error != nil {
		return chatroomMention, result.Error
	}
	return chatroomMention, nil
}

func (crm *ChatroomMention) FindOne(id string) (ChatroomMention, error) {
	var chatroomMention ChatroomMention
	db.First(&chatroomMention, "id = ?", id)

	return chatroomMention, nil
}

// TODO: add pagination for all select queries that return many results
func (crm *ChatroomMention) FindAll(limit float64) ([]ChatroomMention, error) {
	var chatroomMention []ChatroomMention
	db.Limit(int(limit)).Find(&chatroomMention)

	// TODO: to include full details here

	return chatroomMention, nil
}

// Update updates one Agroproduct in the database, using the information
// stored in the receiver u
func (crm *ChatroomMention) Update() (ChatroomMention, error) {
	db.Save(&crm)

	chatroomMention, err := crm.FindOne(crm.ID)
	if err != nil {
		return chatroomMention, err
	}

	return chatroomMention, nil
}

func (crm *ChatroomMention) Delete(id string) error {
	//   TODO: to softly delete chatroom together with its file and mention
	return nil
}
