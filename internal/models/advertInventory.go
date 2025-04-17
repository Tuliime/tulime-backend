package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ai *AdvertInventory) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (ai *AdvertInventory) Create(advertInventory AdvertInventory) (AdvertInventory, error) {
	result := db.Create(&advertInventory)

	if result.Error != nil {
		return advertInventory, result.Error
	}
	return advertInventory, nil
}

func (ai *AdvertInventory) FindOne(id string) (AdvertInventory, error) {
	var advertInventory AdvertInventory
	db.First(&advertInventory, "id = ?", id)

	return advertInventory, nil
}

func (ai *AdvertInventory) FindByAdvert(advertID string) (AdvertInventory, error) {
	var advertInventory AdvertInventory

	db.First(&advertInventory, "\"advertID\" = ?", advertID)

	return advertInventory, nil
}

func (ai *AdvertInventory) Update() (AdvertInventory, error) {
	db.Save(&ai)

	return *ai, nil
}

func (ai *AdvertInventory) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&AdvertInventory{}).Error; err != nil {
		return err
	}

	return nil
}
