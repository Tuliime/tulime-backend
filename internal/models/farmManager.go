package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (fm *FarmManager) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()

	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (fm *FarmManager) Create(farmManager FarmManager) (string, error) {

	result := db.Create(&farmManager)

	if result.Error != nil {
		return "", result.Error
	}
	return farmManager.ID, nil
}

func (fm *FarmManager) FindOne(id string) (FarmManager, error) {
	var farmManager FarmManager
	db.First(&farmManager, "id = ?", id)

	return farmManager, nil
}

func (fm *FarmManager) FindByName(name string) (FarmManager, error) {
	var farmManager FarmManager
	db.First(&farmManager, "name = ?", name)

	return farmManager, nil
}

func (fm *FarmManager) FindByUser(userID string) (FarmManager, error) {
	var farmManager FarmManager
	db.Find(&farmManager, "\"userID\" = ?", userID)

	return farmManager, nil
}

func (fm *FarmManager) FindAll(limit float64) ([]FarmManager, error) {
	var farmManagers []FarmManager
	db.Limit(int(limit)).Find(&farmManagers)

	return farmManagers, nil
}

func (fm *FarmManager) Update() (FarmManager, error) {
	db.Save(&fm)

	farmManager, err := fm.FindOne(fm.ID)
	if err != nil {
		return farmManager, err
	}
	return farmManager, nil
}

func (ap *FarmManager) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&FarmManager{}).Error; err != nil {
		return err
	}
	return nil
}
