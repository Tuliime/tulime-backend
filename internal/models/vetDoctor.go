package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (vd *VetDoctor) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()

	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (vd *VetDoctor) Create(vetDoctor VetDoctor) (string, error) {
	result := db.Create(&vetDoctor)

	if result.Error != nil {
		return "", result.Error
	}
	return vetDoctor.ID, nil
}

func (vd *VetDoctor) FindOne(id string) (VetDoctor, error) {
	var vetDoctor VetDoctor
	db.First(&vetDoctor, "id = ?", id)

	return vetDoctor, nil
}

func (vd *VetDoctor) FindByName(name string) (VetDoctor, error) {
	var vetDoctor VetDoctor
	db.First(&vetDoctor, "name = ?", name)

	return vetDoctor, nil
}

func (vd *VetDoctor) FindByUser(userID string) (VetDoctor, error) {
	var vetDoctor VetDoctor
	db.First(&vetDoctor, "\"userID\" = ?", userID)

	return vetDoctor, nil
}

func (vd *VetDoctor) FindAll(limit float64) ([]VetDoctor, error) {
	var vetDoctors []VetDoctor

	err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, \"imageUrl\"")
	}).
		Limit(int(limit)).
		Find(&vetDoctors).Error

	if err != nil {
		return nil, err
	}

	return vetDoctors, nil
}

func (vd *VetDoctor) Update() (VetDoctor, error) {
	db.Save(&vd)

	vetDoctor, err := vd.FindOne(vd.ID)
	if err != nil {
		return vetDoctor, err
	}
	return vetDoctor, nil
}

func (vd *VetDoctor) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&VetDoctor{}).Error; err != nil {
		return err
	}
	return nil
}
