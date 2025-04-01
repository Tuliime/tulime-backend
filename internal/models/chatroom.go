package models

import (
	"errors"

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
	var chatRoom Chatroom
	err := db.Preload("File").Preload("Mention").Preload("User").Where("id = ?", id).First(&chatRoom).Error
	if err != nil {
		return chatRoom, err
	}
	return chatRoom, nil
}

func (cr *Chatroom) FindReply(reply string) (Chatroom, error) {
	var chatRoom Chatroom
	err := db.Preload("File").Preload("Mention").Preload("User").Where("id = ?", reply).First(&chatRoom).Error
	if err != nil {
		return chatRoom, err
	}
	return chatRoom, nil
}

func (cr *Chatroom) FindAll(limit float64, cursor string, includeCursor bool, direction string) ([]Chatroom, error) {
	var chatRooms []Chatroom

	if direction == "FORWARD" {
		chatRoomsInAscOrder, err := cr.FindAllInAscendingOrder(limit, cursor, includeCursor)
		if err != nil {
			return chatRooms, err
		}
		chatRooms = chatRoomsInAscOrder

	} else if direction == "BACKWARD" {
		chatRoomsInDescOrder, err := cr.FindAllInDescendingOrder(limit, cursor, includeCursor)
		if err != nil {
			return chatRooms, err
		}
		chatRooms = chatRoomsInDescOrder
	} else {
		return chatRooms, errors.New("invalid direction value")
	}

	return chatRooms, nil
}

func (cr *Chatroom) FindAllInDescendingOrder(limit float64, cursor string, includeCursor bool) ([]Chatroom, error) {
	var chatRooms []Chatroom

	query := db.Preload("File").Preload("Mention").Preload("User").Order("\"arrivedAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastChatroom Chatroom
		if err := db.Select("\"arrivedAt\"").Where("id = ?", cursor).First(&lastChatroom).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"arrivedAt\" <= ?", lastChatroom.ArrivedAt)
		} else {
			query = query.Where("\"arrivedAt\" < ?", lastChatroom.ArrivedAt)
		}
	}

	if err := query.Find(&chatRooms).Error; err != nil {
		return nil, err
	}

	// Reverse the slice to return in ascending order
	for i, j := 0, len(chatRooms)-1; i < j; i, j = i+1, j-1 {
		chatRooms[i], chatRooms[j] = chatRooms[j], chatRooms[i]
	}

	return chatRooms, nil
}

func (cr *Chatroom) FindAllInAscendingOrder(limit float64, cursor string, includeCursor bool) ([]Chatroom, error) {
	var chatRooms []Chatroom

	query := db.Preload("File").Preload("Mention").Preload("User").Order("\"arrivedAt\" ASC").Limit(int(limit))

	if cursor != "" {
		var lastChatroom Chatroom
		if err := db.Select("\"arrivedAt\"").Where("id = ?", cursor).First(&lastChatroom).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"arrivedAt\" >= ?", lastChatroom.ArrivedAt)
		} else {
			query = query.Where("\"arrivedAt\" > ?", lastChatroom.ArrivedAt)
		}
	}

	if err := query.Find(&chatRooms).Error; err != nil {
		return nil, err
	}

	return chatRooms, nil
}

// Update updates one Agroproduct in the database, using the information
// stored in the receiver u
func (cr *Chatroom) Update() (Chatroom, error) {
	if err := db.Save(&cr).Error; err != nil {
		return *cr, err
	}
	return *cr, nil
}

func (cr *Chatroom) Delete(id string) error {
	//   TODO: to softly delete chatroom together with its file and mention
	return nil
}
