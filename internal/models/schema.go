package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var db = Db()
var DB = db
var redisClient = RedisClient()
var ctx = context.Background()

type User struct {
	ID                   string             `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name                 string             `gorm:"column:name;not null;index" json:"name"`
	TelNumber            int                `gorm:"column:telNumber;unique;not null;index" json:"telNumber"`
	Password             string             `gorm:"column:password;not null" json:"password"`
	Role                 string             `gorm:"column:role;default:'user';not null" json:"role"`
	Email                string             `gorm:"column:email;index" json:"email"` //Enforce email uniqueness at the application level
	Gender               string             `gorm:"column:gender;default:null" json:"gender"`
	DateOfBirth          string             `gorm:"column:dateOfBirth;default:null" json:"dateOfBirth"`
	Country              string             `gorm:"column:country;default:null" json:"country"`
	ImageUrl             string             `gorm:"column:imageUrl;default:null" json:"imageUrl"`
	ImagePath            string             `gorm:"column:imagePath;default:null" json:"imagePath"`
	ProfileBgColor       string             `gorm:"column:profileBgColor;default:null" json:"profileBgColor"`
	ChatroomColor        string             `gorm:"column:chatroomColor;default:null" json:"chatroomColor"`
	OPT                  []OTP              `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"OPT,omitempty"`
	OnlineStatus         OnlineStatus       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"onlineStatus,omitempty"`
	FarmManager          FarmManager        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"farmManager,omitempty"`
	VetDoctor            *VetDoctor         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vetDoctor,omitempty"`
	Chatroom             []*Chatroom        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"chatroom,omitempty"`
	ChatroomMention      []*ChatroomMention `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"chatroomMention,omitempty"`
	Chatbot              []Chatbot          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"chatbot,omitempty"`
	Session              []*Session         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"session,omitempty"`
	Device               []Device           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"device,omitempty"`
	Notification         []Notification     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"notification,omitempty"`
	MessengerRoomUserOne []*MessengerRoom   `gorm:"foreignKey:UserOneID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"messengerRoomUserOne,omitempty"`
	MessengerRoomUserTwo []*MessengerRoom   `gorm:"foreignKey:UserTwoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"messengerRoomUserTwo,omitempty"`
	Sender               []*Messenger       `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"sender,omitempty"`
	Recipient            []*Messenger       `gorm:"foreignKey:RecipientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"recipient,omitempty"`
	Store                []*Store           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"store,omitempty"`
	Advert               []*Advert          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advert,omitempty"`
	AdvertView           []*AdvertView      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advertView,omitempty"`
	Location             []*Location        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"location,omitempty"`
	StoreFeedback        []*StoreFeedback   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"storeFeedback,omitempty"`
	SearchQuery          []*SearchQuery     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"searchQuery,omitempty"`
	CreatedAt            time.Time          `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt            time.Time          `gorm:"column:updatedAt" json:"updatedAt"`
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
	Description string    `gorm:"column:description;not null;index" json:"description"`
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
	File      ChatroomFile      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"file,omitempty"`
	Mention   []ChatroomMention `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"mention,omitempty"`
	SentAt    time.Time         `gorm:"column:sentAt;not null;index" json:"sentAt"`
	ArrivedAt time.Time         `gorm:"column:arrivedAt;not null;index" json:"arrivedAt"`
	CreatedAt time.Time         `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time         `gorm:"column:updatedAt;index" json:"updatedAt"`
	DeletedAt gorm.DeletedAt    `gorm:"column:deletedAt;index" json:"deletedAt"`
	User      *User             `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user,omitempty"`
}

type ChatroomFile struct {
	ID         string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ChatroomID string         `gorm:"column:chatroomID;not null;index" json:"chatroomID"`
	URL        string         `gorm:"column:url;not null" json:"url"`
	Path       string         `gorm:"column:path;not null" json:"path"`
	Dimensions JSONB          `gorm:"column:dimensions;type:jsonb" json:"dimensions"`
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
	User       *User          `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user,omitempty"`
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
	Device       string    `gorm:"column:Device;default:'Unknown Device';index" json:"device"`
	LocationID   string    `gorm:"column:locationID;default:null" json:"locationID"`
	IsRevoked    bool      `gorm:"column:isRevoked;default:false" json:"isRevoked"`
	CreatedAt    time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	User         *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Location     *Location `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"location,omitempty"`
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

type Store struct {
	ID                  string           `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID              string           `gorm:"column:userID;not null;index" json:"userID"`
	Name                string           `gorm:"column:name;not null;index" json:"name"`
	Description         string           `gorm:"column:description;default:null;index" json:"description"`
	Website             string           `gorm:"column:website;default:null" json:"website"`
	Email               string           `gorm:"column:email;unique;default:null" json:"email"`
	Location            string           `gorm:"column:location;default:null" json:"location"` //TODO: To use json data type
	LogoUrl             string           `gorm:"column:logoUrl;default:null" json:"logoUrl"`
	LogoPath            string           `gorm:"column:logoPath;default:null" json:"logoPath"`
	BackgroundImageUrl  string           `gorm:"column:backgroundImageUrl;default:null" json:"backgroundImageUrl"`
	BackgroundImagePath string           `gorm:"column:backgroundImagePath;default:null" json:"backgroundImagePath"`
	Type                string           `gorm:"column:type;default:'INDIVIDUAL'" json:"type"`
	CreatedAt           time.Time        `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt           time.Time        `gorm:"column:updatedAt;index" json:"updatedAt"`
	User                *User            `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Advert              []*Advert        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"adverts,omitempty"`
	StoreFeedback       []*StoreFeedback `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"feedback,omitempty"`
}

type Advert struct {
	ID                 string              `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	StoreID            string              `gorm:"column:storeID;not null;index" json:"storeID"`
	UserID             string              `gorm:"column:userID;not null;index" json:"userID"`
	ProductName        string              `gorm:"column:productName;not null;index" json:"productName"`
	ProductDescription string              `gorm:"column:productDescription;not null;index" json:"productDescription"`
	IsPublished        bool                `gorm:"column:isPublished;default:false;index" json:"isPublished"`
	AdvertImage        []*AdvertImage      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"images"`
	AdvertPrice        *AdvertPrice        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"price"`
	AdvertInventory    *AdvertInventory    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"inventory"`
	MessengerTag       []*MessengerTag     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"tags"`
	AdvertView         []*AdvertView       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"views"`
	AdvertImpression   []*AdvertImpression `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"impressions"`
	CreatedAt          time.Time           `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt          time.Time           `gorm:"column:updatedAt;index" json:"updatedAt"`
	User               *User               `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Store              *Store              `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"store,omitempty"`
}

type AdvertImage struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AdvertID  string    `gorm:"column:advertID;not null;index" json:"advertID"`
	URL       string    `gorm:"column:url;not null" json:"url"`
	Path      string    `gorm:"column:path;not null" json:"path"`
	IsPrimary bool      `gorm:"column:isPrimary;default:false" json:"isPrimary"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	Advert    *Advert   `gorm:"foreignKey:AdvertID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advert,omitempty"`
}

type AdvertPrice struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AdvertID  string    `gorm:"column:advertID;not null;index" json:"advertID"`
	Amount    float64   `gorm:"column:amount;not null" json:"amount"`
	Currency  string    `gorm:"column:currency;not null" json:"currency"` // Stringified json
	Unit      string    `gorm:"column:unit;not null" json:"unit"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	Advert    *Advert   `gorm:"foreignKey:AdvertID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advert,omitempty"`
}

type AdvertInventory struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AdvertID  string    `gorm:"column:advertID;not null;index" json:"advertID"`
	Quantity  float64   `gorm:"column:quantity;not null" json:"quantity"`
	Unit      string    `gorm:"column:unit;not null" json:"unit"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	Advert    *Advert   `gorm:"foreignKey:AdvertID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advert,omitempty"`
}

type AdvertView struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AdvertID   string    `gorm:"column:advertID;not null;index" json:"advertID"`
	UserID     string    `gorm:"column:userID;not null;index" json:"userID"`
	LocationID string    `gorm:"column:location;not null" json:"locationID"`
	Device     string    `gorm:"column:device;default:'UNKNOWN'" json:"device"`
	CreatedAt  time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	Advert     *Advert   `gorm:"foreignKey:AdvertID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advert,omitempty"`
	User       *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Location   *Location `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"location,omitempty"`
}

type AdvertImpression struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	AdvertID   string    `gorm:"column:advertID;not null;index" json:"advertID"`
	UserID     string    `gorm:"column:userID;not null;index" json:"userID"`
	LocationID string    `gorm:"column:location;not null" json:"locationID"`
	Device     string    `gorm:"column:device;default:'UNKNOWN'" json:"device"`
	CreatedAt  time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	Advert     *Advert   `gorm:"foreignKey:AdvertID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advert,omitempty"`
	User       *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Location   *Location `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"location,omitempty"`
}

type StoreFeedback struct {
	ID          string              `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	StoreID     string              `gorm:"column:storeID;not null;index" json:"storeID"`
	UserID      string              `gorm:"column:userID;not null;index" json:"userID"`
	Experience  string              `gorm:"column:experience;not null" json:"experience"`
	Title       string              `gorm:"column:title;not null" json:"title"`
	Description string              `gorm:"column:description;not null" json:"description"`
	Reply       string              `gorm:"column:reply;default:null;index" json:"reply"`
	File        []StoreFeedbackFile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"files"`
	CreatedAt   time.Time           `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt   time.Time           `gorm:"column:updatedAt;index" json:"updatedAt"`
	Store       *Store              `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	User        *User               `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"store,omitempty"`
}

type StoreFeedbackFile struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	StoreFeedbackID string    `gorm:"column:storeFeedbackID;not null;index" json:"storeFeedbackID"`
	URL             string    `gorm:"column:url;not null" json:"url"`
	Path            string    `gorm:"column:path;not null" json:"path"`
	CreatedAt       time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type MessengerRoom struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserOneID string    `gorm:"column:userOneID;not null;index" json:"userOneID"`
	UserTwoID string    `gorm:"column:userTwoID;not null;index" json:"userTwoID"`
	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	UserOne   *User     `gorm:"foreignKey:UserOneID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"userOne,omitempty"`
	UserTwo   *User     `gorm:"foreignKey:UserTwoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"userTwo,omitempty"`
}

type Messenger struct {
	ID              string         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	MessengerRoomID string         `gorm:"column:messengerRoomID;not null;index" json:"messengerRoomID"`
	SenderID        string         `gorm:"column:senderID;not null;index" json:"senderID"`
	RecipientID     string         `gorm:"column:recipientID;not null;index" json:"recipientID"`
	Text            string         `gorm:"column:text;default:null" json:"text"`
	Reply           string         `gorm:"column:reply;default:null;index" json:"reply"`
	File            MessengerFile  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"file"`
	Tag             []MessengerTag `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tag"`
	IsRead          bool           `gorm:"column:isRead;default:false" json:"isRead"`
	SentAt          time.Time      `gorm:"column:sentAt;not null;index" json:"sentAt"`
	ArrivedAt       time.Time      `gorm:"column:arrivedAt;not null;index" json:"arrivedAt"`
	CreatedAt       time.Time      `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updatedAt;index" json:"updatedAt"`
	Sender          *User          `gorm:"foreignKey:SenderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"sender,omitempty"`
	Recipient       *User          `gorm:"foreignKey:RecipientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"recipient,omitempty"`
}

type MessengerFile struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	MessengerID string    `gorm:"column:messengerID;not null;index" json:"messengerID"`
	URL         string    `gorm:"column:url;not null" json:"url"`
	Path        string    `gorm:"column:path;not null" json:"path"`
	Dimensions  JSONB     `gorm:"column:dimensions;type:jsonb" json:"dimensions"`
	CreatedAt   time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
}

type MessengerTag struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	MessengerID string    `gorm:"column:messengerID;not null;index" json:"messengerID"`
	AdvertID    string    `gorm:"column:advertID;not null;index" json:"advertID"`
	CreatedAt   time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	Advert      *Advert   `gorm:"foreignKey:AdvertID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advert,omitempty"`
}

type Location struct {
	ID               string              `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID           string              `gorm:"column:userID;not null;index" json:"userID"`
	Info             JSONB               `gorm:"column:info;type:jsonb;not null;" json:"info"`
	CreatedAt        time.Time           `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt        time.Time           `gorm:"column:updatedAt;index" json:"updatedAt"`
	User             *User               `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advert,omitempty"`
	AdvertView       []*AdvertView       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advertView,omitempty"`
	AdvertImpression []*AdvertImpression `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"advertImpression,omitempty"`
	Session          []*Session          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"session,omitempty"`
	SearchQuery      []*SearchQuery      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"searchQuery,omitempty"`
}

type SearchQuery struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID     string    `gorm:"column:userID;not null;index" json:"userID"`
	LocationID string    `gorm:"column:location;not null" json:"locationID"`
	Device     string    `gorm:"column:device;default:'UNKNOWN'" json:"device"`
	Query      string    `gorm:"column:query;not null" json:"query"`
	CreatedAt  time.Time `gorm:"column:createdAt;index" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
	User       *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Location   *Location `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"location,omitempty"`
}

// Other Types
type SendNotification = struct {
	Notification Notification
	DeviceToken  string
}

type ImageDimensions = struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
