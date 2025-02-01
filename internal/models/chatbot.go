package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (cb *Chatbot) BeforeCreate(tx *gorm.DB) error {
	var ID string
	if cb.ID == "" {
		ID = uuid.New().String()
	} else {
		ID = cb.ID
	}
	tx.Statement.SetColumn("ID", ID)
	return nil
}

func (cb *Chatbot) Create(chatbot Chatbot) (Chatbot, error) {
	if !cb.ValidateWrittenBy(chatbot.WrittenBy) {
		return chatbot, errors.New("invalid writtenBy value")
	}
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
// stored in the receiver cb
func (cb *Chatbot) Update() (Chatbot, error) {
	if !cb.ValidateWrittenBy(cb.WrittenBy) {
		return *cb, errors.New("invalid writtenBy value")
	}
	db.Save(&cb)
	return *cb, nil
}

func (cb *Chatbot) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Chatbot{}).Error; err != nil {
		return err
	}

	return nil
}

func (cb *Chatbot) ValidateWrittenBy(writtenBy string) bool {
	var allowedWrittenBy = []string{"user", "bot"}
	for _, awb := range allowedWrittenBy {
		if awb == writtenBy {
			return true
		}
	}
	return false
}
