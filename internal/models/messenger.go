package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (msgr *Messenger) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (msgr *Messenger) Create(message Messenger) (Messenger, error) {
	result := db.Create(&message)

	if result.Error != nil {
		return message, result.Error
	}
	return message, nil
}

func (cr *Messenger) FindOne(id string) (Messenger, error) {
	var message Messenger
	db.First(&message, "id = ?", id)

	return message, nil
}

func (msgr *Messenger) FindReply(reply string) (Messenger, error) {
	var message Messenger
	err := db.Preload("File").Preload("Tag").Where("id = ?", reply).First(&message).Error
	if err != nil {
		return message, err
	}
	return message, nil
}

func (msgr *Messenger) FindByRoom(roomID string, limit float64, cursor string,
	includeCursor bool, direction string) ([]Messenger, error) {
	var messages []Messenger

	if direction == "FORWARD" {
		messagesInAscOrder, err := msgr.FindInAscOrderByRoom(roomID, limit, cursor, includeCursor)
		if err != nil {
			return messages, err
		}
		messages = messagesInAscOrder

	} else if direction == "BACKWARD" {
		messagesInDescOrder, err := msgr.FindInDescOrderByRoom(roomID, limit, cursor, includeCursor)
		if err != nil {
			return messages, err
		}
		messages = messagesInDescOrder
	} else {
		return messages, errors.New("invalid direction value")
	}

	return messages, nil
}

func (msgr *Messenger) FindInDescOrderByRoom(roomID string, limit float64,
	cursor string, includeCursor bool) ([]Messenger, error) {
	var messages []Messenger

	query := db.Preload("File").Preload("Tag").Order("\"arrivedAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastMessage Messenger
		if err := db.Select("\"arrivedAt\"").Where("id = ?", cursor).First(&lastMessage).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"arrivedAt\" <= ?", lastMessage.ArrivedAt)
		} else {
			query = query.Where("\"arrivedAt\" < ?", lastMessage.ArrivedAt)
		}
	}

	if err := query.Where("\"messengerRoomID\" = ?", roomID).Find(&messages).Error; err != nil {
		return nil, err
	}

	// Reverse the slice to return in ascending order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (msgr *Messenger) FindInAscOrderByRoom(roomID string, limit float64,
	cursor string, includeCursor bool) ([]Messenger, error) {
	var messages []Messenger

	query := db.Preload("File").Preload("Tag").Order("\"arrivedAt\" ASC").Limit(int(limit))

	if cursor != "" {
		var lastMessage Messenger
		if err := db.Select("\"arrivedAt\"").Where("id = ?", cursor).First(&lastMessage).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"arrivedAt\" >= ?", lastMessage.ArrivedAt)
		} else {
			query = query.Where("\"arrivedAt\" > ?", lastMessage.ArrivedAt)
		}
	}

	if err := query.Where("\"messengerRoomID\" = ?", roomID).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// FindRoomsByUser fetches messages with sender and recipient
// details in descending order
func (msgr *Messenger) FindRoomsByUser(userID string, limit float64,
	cursor string) ([]Messenger, error) {
	var messengerRooms []Messenger

	query := db.Table("messengers").Select("DISTINCT ON (\"messengerRoomID\") *").
		Order("\"messengerRoomID\", \"arrivedAt\" DESC").
		Limit(int(limit))

	query = query.Preload("File").Preload("Tag").Preload("Sender").Preload("Recipient")

	if cursor != "" {
		var lastMessage Messenger
		if err := db.Select("\"arrivedAt\"").Where("id = ?",
			cursor).First(&lastMessage).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"arrivedAt\" < ?", lastMessage.ArrivedAt)
	}
	if err := query.Where("\"senderID\" = ? OR \"recipientID\" = ?",
		userID, userID).Find(&messengerRooms).Error; err != nil {
		return nil, err
	}

	// Reverse the slice to return in descending order
	for i, j := 0, len(messengerRooms)-1; i < j; i, j = i+1, j-1 {
		messengerRooms[i], messengerRooms[j] = messengerRooms[j], messengerRooms[i]
	}
	return messengerRooms, nil
}

// Update updates one Messenger in the database, using the information
// stored in the receiver u
func (msgr *Messenger) Update() (Messenger, error) {
	if err := db.Save(&msgr).Error; err != nil {
		return *msgr, err
	}
	return *msgr, nil
}

func (msgr *Messenger) Delete(id string) error {
	//   TODO: to softly delete Messenger together with its file and tag
	return nil
}
