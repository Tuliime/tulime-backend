package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ai *AdvertImpression) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (ai *AdvertImpression) Create(advertImpressions []AdvertImpression) ([]AdvertImpression, error) {
	result := db.Create(&advertImpressions)

	if result.Error != nil {
		return advertImpressions, result.Error
	}
	return advertImpressions, nil
}

func (ai *AdvertImpression) FindOne(id string) (AdvertImpression, error) {
	var advertImpression AdvertImpression
	db.First(&advertImpression, "id = ?", id)

	return advertImpression, nil
}

func (ai *AdvertImpression) FindByAdvert(advertID string, limit float64, cursor string) ([]AdvertImpression, error) {
	var advertImpressions []AdvertImpression
	query := db.Order("\"createAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastAdvertImpression AdvertImpression
		if err := db.Select("\"createAt\"").Where("id = ?", cursor).First(&lastAdvertImpression).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createAt\" < ?", lastAdvertImpression.CreatedAt)
	}
	if err := query.Where("\"advertID\" = ?", advertID).Find(&advertImpressions).Error; err != nil {
		return nil, err
	}

	return advertImpressions, nil
}

// TODO: To add counts based on time
func (ai *AdvertImpression) FindCountByAdvert(advertID string) (int64, error) {
	var viewCount int64

	err := db.Model(&AdvertImpression{}).Where("\"advertID\" = ?", advertID).Count(&viewCount).Error
	if err != nil {
		return 0, err
	}

	return viewCount, nil
}

func (ai *AdvertImpression) Update() (AdvertImpression, error) {
	db.Save(&ai)

	return *ai, nil
}

// TODO: consider soft deleting the AdvertImpression
func (ai *AdvertImpression) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&AdvertImpression{}).Error; err != nil {
		return err
	}

	return nil
}
