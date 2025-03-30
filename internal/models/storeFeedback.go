package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (f *StoreFeedback) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (f *StoreFeedback) Create(feedback StoreFeedback) (StoreFeedback, error) {
	result := db.Create(&feedback)

	if result.Error != nil {
		return feedback, result.Error
	}
	return feedback, nil
}

func (f *StoreFeedback) FindOne(id string) (StoreFeedback, error) {
	var feedback StoreFeedback
	db.Preload("File").Preload("User").First(&feedback, "id = ?", id)

	return feedback, nil
}

func (f *StoreFeedback) FindReply(reply string) ([]StoreFeedback, error) {
	var feedback []StoreFeedback
	err := db.Preload("File").Preload("User").Where("id = ?", reply).Find(&feedback).Error
	if err != nil {
		return feedback, err
	}
	return feedback, nil
}

func (f *StoreFeedback) FindByStore(storeID string, limit float64,
	cursor string) ([]StoreFeedback, error) {
	var feedback []StoreFeedback

	query := db.Preload("File").Preload("User").Order("\"createdAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastFeedback StoreFeedback
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastFeedback).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createdAt\" < ?", lastFeedback.CreatedAt)
	}

	if err := query.Where("\"storeID\" = ?", storeID).Find(&feedback).Error; err != nil {
		return nil, err
	}

	return feedback, nil
}

// Update updates one Messenger in the database, using the information
// stored in the receiver u
func (f *StoreFeedback) Update() (StoreFeedback, error) {
	if err := db.Save(&f).Error; err != nil {
		return *f, err
	}
	return *f, nil
}

func (f *StoreFeedback) Delete(id string) error {
	//   TODO: to softly delete StoreFeedback
	return nil
}
