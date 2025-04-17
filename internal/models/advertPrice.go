package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ap *AdvertPrice) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (ap *AdvertPrice) Create(advertPrice AdvertPrice) (AdvertPrice, error) {
	result := db.Create(&advertPrice)

	if result.Error != nil {
		return advertPrice, result.Error
	}
	return advertPrice, nil
}

func (ap *AdvertPrice) FindOne(id string) (AdvertPrice, error) {
	var advertPrice AdvertPrice
	db.First(&advertPrice, "id = ?", id)

	return advertPrice, nil
}

func (ai *AdvertPrice) FindByAdvert(advertID string) (AdvertPrice, error) {
	var advertPrice AdvertPrice

	db.First(&advertPrice, "\"advertID\" = ?", advertID)

	return advertPrice, nil
}

func (ap *AdvertPrice) Update() (AdvertPrice, error) {
	db.Save(&ap)

	return *ap, nil
}

func (ap *AdvertPrice) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&AdvertPrice{}).Error; err != nil {
		return err
	}

	return nil
}
