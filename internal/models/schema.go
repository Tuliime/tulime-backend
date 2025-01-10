package models

import (
	// "context"
	"time"
	// "gorm.io/gorm"
)

var db = Db()
var DB = db

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
