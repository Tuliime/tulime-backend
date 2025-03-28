package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (av *AdvertView) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (av *AdvertView) Create(advertViews []AdvertView) ([]AdvertView, error) {
	result := db.Create(&advertViews)

	if result.Error != nil {
		return advertViews, result.Error
	}
	return advertViews, nil
}

func (av *AdvertView) FindOne(id string) (AdvertView, error) {
	var advertView AdvertView
	db.First(&advertView, "id = ?", id)

	return advertView, nil
}

func (av *AdvertView) FindByAdvert(advertID string, limit float64, cursor string) ([]AdvertView, error) {
	var advertViews []AdvertView
	query := db.Order("\"createAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastAdvertView AdvertView
		if err := db.Select("\"createAt\"").Where("id = ?", cursor).First(&lastAdvertView).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createAt\" < ?", lastAdvertView.CreatedAt)
	}
	if err := query.Where("\"advertID\" = ?", advertID).Find(&advertViews).Error; err != nil {
		return nil, err
	}

	return advertViews, nil
}

// TODO: To add counts based on time
func (av *AdvertView) FindCountByAdvert(advertID string) (int64, error) {
	var viewCount int64

	err := db.Model(&AdvertView{}).Where("advertID = ?", advertID).Count(&viewCount).Error
	if err != nil {
		return 0, err
	}

	return viewCount, nil
}

func (av *AdvertView) Update() (AdvertView, error) {
	db.Save(&av)

	return *av, nil
}

// TODO: consider soft deleting the Location
func (av *AdvertView) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&AdvertView{}).Error; err != nil {
		return err
	}

	return nil
}
