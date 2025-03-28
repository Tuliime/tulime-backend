package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (l *Location) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (l *Location) Create(location Location) (Location, error) {
	result := db.Create(&location)

	if result.Error != nil {
		return location, result.Error
	}
	return location, nil
}

func (l *Location) FindOne(id string) (Location, error) {
	var location Location
	db.First(&location, "id = ?", id)

	return location, nil
}

// GetLocationByIP fetches a location record where Info.query matches the given IP
func (l *Location) FindByIP(ip string) (Location, error) {
	var location Location

	// Query using PostgreSQL JSONB extraction
	db.Where("info->>'query' = ?", ip).First(&location)

	return location, nil
}

func (l *Location) FindByUser(userID string, limit float64, cursor string) ([]Location, error) {
	var locations []Location
	query := db.Order("\"createAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastLocation Location
		if err := db.Select("\"createAt\"").Where("id = ?", cursor).First(&lastLocation).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createAt\" < ?", lastLocation.CreatedAt)
	}
	if err := query.Where("\"userID\" = ?", userID).Find(&locations).Error; err != nil {
		return nil, err
	}

	return locations, nil
}

func (l *Location) Update() (Location, error) {
	db.Save(&l)

	return *l, nil
}

// TODO: consider soft deleting the Location
func (l *Location) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Location{}).Error; err != nil {
		return err
	}

	return nil
}
