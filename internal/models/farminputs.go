package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (f *FarmInputs) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (f *FarmInputs) Create(farmInputs FarmInputs) (FarmInputs, error) {
	result := db.Create(&farmInputs)

	if result.Error != nil {
		return farmInputs, result.Error
	}
	return farmInputs, nil
}

func (f *FarmInputs) FindOne(id string) (FarmInputs, error) {
	var FarmInputs FarmInputs
	db.First(&FarmInputs, "id = ?", id)

	return FarmInputs, nil
}

// TODO: add pagination for all select queries that return many results
func (f *FarmInputs) FindByCategory(category string) ([]FarmInputs, error) {
	var FarmInputs []FarmInputs
	db.Find(&FarmInputs, "category = ?", category)

	return FarmInputs, nil
}

// TODO: add pagination for all select queries that return many results
func (f *FarmInputs) FindAll(limit float64) ([]FarmInputs, error) {
	var FarmInputs []FarmInputs
	db.Limit(int(limit)).Find(&FarmInputs)

	return FarmInputs, nil
}

// Update updates one FarmInputs in the database, using the information
// stored in the receiver u
func (f *FarmInputs) Update() (FarmInputs, error) {
	db.Save(&f)

	FarmInputs, err := f.FindOne(f.ID)
	if err != nil {
		return FarmInputs, err
	}

	return FarmInputs, nil
}

func (f *FarmInputs) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&FarmInputs{}).Error; err != nil {
		return err
	}

	return nil
}
