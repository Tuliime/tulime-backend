package models

import (
	// "context"
	"context"
	"time"

	"gorm.io/gorm"
	// "gorm.io/gorm"
)

var db = Db()
var DB = db
var redisClient = RedisClient()
var ctx = context.Background()

type User struct {
	ID             string       `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name           string       `gorm:"column:name;not null;index" json:"name"`
	TelNumber      int          `gorm:"column:telNumber;unique;not null;index" json:"telNumber"`
	Password       string       `gorm:"column:password;not null" json:"password"`
	Role           string       `gorm:"column:role;default:'user';not null" json:"role"`
	ImageUrl       string       `gorm:"column:imageUrl;default:null" json:"imageUrl"`
	ImagePath      string       `gorm:"column:imagePath;default:null" json:"imagePath"`
	ProfileBgColor string       `gorm:"column:profileBgColor;default:null" json:"profileBgColor"`
	ChatroomColor  string       `gorm:"column:chatroomColor;default:null" json:"chatroomColor"`
	OPT            []OTP        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"OPT"`
	OnlineStatus   OnlineStatus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"onlineStatus"`
	FarmManager    FarmManager  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"farmManager"`
	// VetDoctor       VetDoctor         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vetDoctor"`
	VetDoctor       *VetDoctor        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vetDoctor"`
	Chatroom        []Chatroom        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"chatroom"`
	ChatroomMention []ChatroomMention `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"chatroomMention"`
	Chatbot         []Chatbot         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"chatbot"`
	Session         []Session         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"session"`
	Device          []Device          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"device"`
	Notification    []Notification    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"notification"`
	CreatedAt       time.Time         `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time         `gorm:"column:updatedAt" json:"updatedAt"`
}

type OnlineStatus struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"column:userID;unique;not null;index" json:"userID"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type Agroproduct struct {
	ID        string             `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name      string             `gorm:"column:name;unique;not null;index" json:"name"`
	Category  string             `gorm:"column:category;not null;index" json:"category"`
	ImageUrl  string             `gorm:"column:imageUrl;not null" json:"imageUrl"`
	ImagePath string             `gorm:"column:imagePath;default:null" json:"imagePath"`
	Price     []AgroproductPrice `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"Price"`
	CreatedAt time.Time          `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time          `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type AgroproductPrice struct {
	ID            string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AgroproductID string    `gorm:"column:agroproductID;not null;index" json:"agroproductID"`
	Amount        float64   `gorm:"column:amount;not null;index" json:"amount"`
	Currency      string    `gorm:"column:currency;not null;index" json:"currency"`
	CreatedAt     time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type News struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Title       string    `gorm:"column:title;unique;not null;index" json:"title"`
	Description string    `gorm:"column:description;not null" json:"description"`
	Category    string    `gorm:"column:category;not null;index" json:"category"`
	Source      string    `gorm:"column:source;not null;index" json:"source"`
	ImageUrl    string    `gorm:"column:imageUrl;not null" json:"imageUrl"`
	ImagePath   string    `gorm:"column:imagePath;not null" json:"imagePath"`
	PostedAt    time.Time `gorm:"column:postedAt;default:CURRENT_TIMESTAMP;index" json:"postedAt"`
	CreatedAt   time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type FarmInputs struct {
	ID            string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name          string    `gorm:"column:name;unique;not null;index" json:"name"`
	Purpose       string    `gorm:"column:purpose;not null" json:"purpose"`
	Category      string    `gorm:"column:category;not null;index" json:"category"`
	ImageUrl      string    `gorm:"column:imageUrl;not null" json:"imageUrl"`
	ImagePath     string    `gorm:"column:imagePath;not null" json:"imagePath"`
	Price         float64   `gorm:"column:price;not null" json:"price"`
	PriceCurrency string    `gorm:"column:priceCurrency;not null" json:"priceCurrency"`
	Source        string    `gorm:"column:source;not null" json:"Source"`
	SourceUrl     string    `gorm:"column:sourceUrl;default:null" json:"sourceUrl"`
	CreatedAt     time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type OTP struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID     string    `gorm:"column:userID;not null;index" json:"userID"`
	OTP        string    `gorm:"column:OTP;not null;index" json:"OTP"`
	IsUsed     bool      `gorm:"column:isUsed;default:false" json:"isUsed"`
	IsVerified bool      `gorm:"column:isVerified;default:false" json:"isVerified"`
	ExpiresAt  time.Time `gorm:"column:expiresAt;not null;index" json:"expiresAt"`
	CreatedAt  time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type FarmManager struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID     string    `gorm:"column:userID;unique;not null;index" json:"userID"`
	Name       string    `gorm:"column:name;not null;index" json:"name"`
	Gender     string    `gorm:"column:gender;not null" json:"gender"`
	RegNo      string    `gorm:"column:regNo;not null" json:"regNo"`
	Email      string    `gorm:"column:email;unique;not null" json:"email"`
	TelNumber  int       `gorm:"column:telNumber;not null" json:"telNumber"`
	IsVerified bool      `gorm:"column:isVerified;default:false" json:"isVerified"`
	CreatedAt  time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type VetDoctor struct {
	ID            string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID        string    `gorm:"column:userID;unique;not null;index" json:"userID"`
	Name          string    `gorm:"column:name;not null;index" json:"name"`
	Gender        string    `gorm:"column:gender;not null" json:"gender"`
	LicenseNumber string    `gorm:"column:licenseNumber;not null" json:"licenseNumber"`
	Email         string    `gorm:"column:email;unique;not null" json:"email"`
	TelNumber     int       `gorm:"column:telNumber;not null" json:"telNumber"`
	IsVerified    bool      `gorm:"column:isVerified;default:false" json:"isVerified"`
	CreatedAt     time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	User          *User     `gorm:"foreignKey:UserID;references:ID" json:"user"`
}

type Chatroom struct {
	ID        string            `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID    string            `gorm:"column:userID;not null;index" json:"userID"`
	Text      string            `gorm:"column:text;default:null" json:"text"`
	Reply     string            `gorm:"column:reply;default:null;index" json:"reply"`
	File      ChatroomFile      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"file"`
	Mention   []ChatroomMention `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"mention"`
	SentAt    time.Time         `gorm:"column:sentAt;not null;index" json:"sentAt"`
	ArrivedAt time.Time         `gorm:"column:arrivedAt;not null;index" json:"arrivedAt"`
	CreatedAt time.Time         `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time         `gorm:"column:updatedAt;index" json:"updatedAt"`
	DeletedAt gorm.DeletedAt    `gorm:"column:deletedAt;index" json:"deletedAt"`
}

type ChatroomFile struct {
	ID         string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ChatroomID string         `gorm:"column:chatroomID;not null;index" json:"chatroomID"`
	URL        string         `gorm:"column:url;not null" json:"url"`
	Path       string         `gorm:"column:path;not null" json:"path"`
	CreatedAt  time.Time      `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updatedAt;index" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt"`
}

type ChatroomMention struct {
	ID         string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ChatroomID string         `gorm:"column:chatroomID;not null;index" json:"chatroomID"`
	UserID     string         `gorm:"column:userID;not null;index" json:"userID"`
	CreatedAt  time.Time      `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updatedAt;index" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deletedAt;index" json:"deletedAt"`
}

type Chatbot struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"column:userID;type:uuid;not null;index" json:"userID"`
	Message   string    `gorm:"column:message;not null" json:"message"`
	WrittenBy string    `gorm:"column:writtenBy;not null" json:"writtenBy"` // "user" or "bot"
	PostedAt  time.Time `gorm:"column:postedAt;index" json:"postedAt"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type Session struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID       string    `gorm:"column:userID;not null;index" json:"userID"`
	AccessToken  string    `gorm:"column:accessToken;not null;index" json:"accessToken"`
	RefreshToken string    `gorm:"column:refreshToken;not null;index" json:"refreshToken"`
	GeneratedVia string    `gorm:"column:generatedVia;not null;index" json:"generatedVia"`
	IsRevoked    bool      `gorm:"column:isRevoked;default:false" json:"isRevoked"`
	CreatedAt    time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type Device struct {
	ID                   string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID               string    `gorm:"column:userID;not null;index" json:"userID"`
	Token                string    `gorm:"column:token;not null;index" json:"token"`
	TokenType            string    `gorm:"column:tokenType;not null" json:"tokenType"`
	Name                 string    `gorm:"column:name;not null" json:"name"`
	NotificationDisabled bool      `gorm:"column:notificationDisabled;default:false" json:"notificationDisabled"`
	CreatedAt            time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt            time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type Notification struct {
	ID             string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID         string    `gorm:"column:userID;not null;index" json:"userID"`
	Title          string    `gorm:"column:title;not null" json:"title"`
	Body           string    `gorm:"column:body;not null" json:"body"`
	Data           string    `gorm:"column:data" json:"data"` // stringified json
	Icon           string    `gorm:"column:icon" json:"icon"`
	Attachments    string    `gorm:"column:attachments" json:"attachments"` // stringified json
	IsRead         bool      `gorm:"column:isRead;default:false;index" json:"isRead"`
	SendStatusCode int       `gorm:"column:sendStatusCode;not null" json:"sendStatusCode"`
	Type           string    `gorm:"column:type;not null;index" json:"type"`
	CreatedAt      time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

// Other Types
type SendNotification = struct {
	Notification Notification
	DeviceToken  string
}
