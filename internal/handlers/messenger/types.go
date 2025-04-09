package messenger

import (
	"time"

	"github.com/Tuliime/tulime-backend/internal/models"
	"gorm.io/gorm"
)

type User struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	TelNumber      int       `json:"telNumber"`
	Role           string    `json:"role"`
	ImageUrl       string    `json:"imageUrl"`
	ImagePath      string    `json:"imagePath"`
	ProfileBgColor string    `json:"profileBgColor"`
	ChatroomColor  string    `json:"chatroomColor"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Message struct {
	ID              string                `json:"id"`
	MessengerRoomID string                `json:"messengerRoomID"`
	SenderID        string                `json:"senderID"`
	RecipientID     string                `json:"recipientID"`
	Text            string                `json:"text"`
	Reply           string                `json:"reply"`
	RepliedMessage  any                   `json:"repliedMessage"`
	File            any                   `json:"file"`
	Tag             []models.MessengerTag `json:"tag"`
	IsRead          bool                  `json:"isRead"`
	SentAt          time.Time             `json:"sentAt"`
	ArrivedAt       time.Time             `json:"arrivedAt"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	Sender          any                   `json:"sender"`
	Recipient       any                   `json:"recipient"`
}

type File struct {
	ID          string                 `json:"id"`
	MessengerID string                 `json:"messengerID"`
	URL         string                 `json:"url"`
	Path        string                 `json:"path"`
	Dimensions  models.ImageDimensions `json:"dimensions"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt         `json:"deletedAt"`
}
