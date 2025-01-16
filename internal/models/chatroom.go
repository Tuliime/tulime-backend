package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (cr *Chatroom) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (cr *Chatroom) Create(chatroom Chatroom) (Chatroom, error) {
	result := db.Create(&chatroom)

	if result.Error != nil {
		return chatroom, result.Error
	}
	return chatroom, nil
}

func (cr *Chatroom) FindOne(id string) (Chatroom, error) {
	var chatroom Chatroom
	db.First(&chatroom, "id = ?", id)

	return chatroom, nil
}

// TODO: to include file and mention
func (cr *Chatroom) FindAll(limit float64, cursor string) ([]Chatroom, error) {
	var chatRooms []Chatroom

	query := db.Order("\"arrivedAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastChatroom Chatroom
		if err := db.Select("\"arrivedAt\"").Where("id = ?", cursor).First(&lastChatroom).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"arrivedAt\" < ?", lastChatroom.ArrivedAt)
	}

	if err := query.Find(&chatRooms).Error; err != nil {
		return nil, err
	}

	return chatRooms, nil
}

// Update updates one Agroproduct in the database, using the information
// stored in the receiver u
func (cr *Chatroom) Update() (Chatroom, error) {
	db.Save(&cr)

	chatRoom, err := cr.FindOne(cr.ID)
	if err != nil {
		return chatRoom, err
	}

	return chatRoom, nil
}

func (cr *Chatroom) Delete(id string) error {
	//   TODO: to softly delete chatroom together with its file and mention
	return nil
}
