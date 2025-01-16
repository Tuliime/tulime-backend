package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (cb *Chatbot) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (cb *Chatbot) Create(chatbot Chatbot) (Chatbot, error) {
	result := db.Create(&chatbot)

	if result.Error != nil {
		return chatbot, result.Error
	}
	return chatbot, nil
}

func (cb *Chatbot) FindOne(id string) (Chatbot, error) {
	var chatbot Chatbot
	db.First(&chatbot, "id = ?", id)

	return chatbot, nil
}

func (cb *Chatbot) FindByUser(userID string, limit float64, cursor string) ([]Chatbot, error) {
	var chatbot []Chatbot
	query := db.Order("\"updatedAt\" ASC").Where("\"userID\" =?", userID).Limit(int(limit))

	if cursor != "" {
		var lastChat Chatbot
		if err := db.Select("\"updatedAt\"").Where("id = ?", cursor).First(&lastChat).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"updatedAt\" < ?", lastChat.UpdatedAt)
	}

	query.Find(&chatbot)

	return chatbot, nil
}

// Update updates one Chatbot in the database, using the information
// stored in the receiver u
func (cb *Chatbot) Update() (Chatbot, error) {
	db.Save(&cb)

	chatbot, err := cb.FindOne(cb.ID)
	if err != nil {
		return chatbot, err
	}

	return chatbot, nil
}

func (cb *Chatbot) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Chatbot{}).Error; err != nil {
		return err
	}

	return nil
}
