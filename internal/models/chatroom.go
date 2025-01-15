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

// TODO: add pagination for all select queries that return many results
func (cr *Chatroom) FindAll(limit float64) ([]Chatroom, error) {
	var chatRoom []Chatroom
	db.Limit(int(limit)).Find(&chatRoom)

	// TODO: to fetch file and mention for each chat
	// for i, agroproduct := range chatRoom {
	// 	var Chatroom []Chatroom
	// 	db.Order("\"createdAt\" desc").Limit(1).Find(&agroProductPrice, "\"agroproductID\" = ?", agroproduct.ID)
	// 	agroproduct.Price = agroProductPrice
	// 	chatRoom[i] = agroproduct
	// }

	return chatRoom, nil
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
