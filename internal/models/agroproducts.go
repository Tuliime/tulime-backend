package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ap *Agroproduct) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (ap *Agroproduct) Create(agroproduct Agroproduct) (Agroproduct, error) {
	result := db.Create(&agroproduct)

	if result.Error != nil {
		return agroproduct, result.Error
	}
	// return agroproduct.ID, nil
	return agroproduct, nil
}

func (ap *Agroproduct) FindOne(id string) (Agroproduct, error) {
	var agroproduct Agroproduct
	db.First(&agroproduct, "id = ?", id)

	return agroproduct, nil
}

func (ap *Agroproduct) FindByName(name string) (Agroproduct, error) {
	var agroproduct Agroproduct
	db.First(&agroproduct, "name = ?", name)

	return agroproduct, nil
}

// TODO: add pagination for all select queries that return many results
func (ap *Agroproduct) FindByCategory(name string) ([]Agroproduct, error) {
	var agroproducts []Agroproduct
	db.Find(&agroproducts, "name = ?", name)

	return agroproducts, nil
}

// TODO: add pagination for all select queries that return many results
func (ap *Agroproduct) FindAll() ([]Agroproduct, error) {
	var agroproducts []Agroproduct
	db.Find(&agroproducts)

	return agroproducts, nil
}

// Update updates one Agroproduct in the database, using the information
// stored in the receiver u
func (ap *Agroproduct) Update() error {
	db.Save(&ap)

	return nil
}

func (ap *Agroproduct) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Agroproduct{}).Error; err != nil {
		return err
	}

	return nil
}

func (ap *Agroproduct) ValidCategory(category string) bool {
	categories := []string{"crop", "livestock", "poultry", "fish"}

	for _, r := range categories {
		if r == category {
			return true
		}
	}

	return false
}
