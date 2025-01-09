package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (app *AgroproductPrice) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (app *AgroproductPrice) Create(agroproductPrice AgroproductPrice) (AgroproductPrice, error) {
	result := db.Create(&agroproductPrice)

	if result.Error != nil {
		return agroproductPrice, result.Error
	}
	return agroproductPrice, nil
}

func (app *AgroproductPrice) FindOne(id string) (AgroproductPrice, error) {
	var AgroproductPrice AgroproductPrice
	db.First(&AgroproductPrice, "id = ?", id)

	return AgroproductPrice, nil
}

func (app *AgroproductPrice) FindByAgroProduct(agroproductID string) ([]AgroproductPrice, error) {
	var AgroproductPrices []AgroproductPrice
	db.Order("\"updatedAt\" desc").Find(&AgroproductPrices, "\"agroproductID\" = ?", agroproductID)

	return AgroproductPrices, nil
}

func (app *AgroproductPrice) FindAll() ([]AgroproductPrice, error) {
	var AgroproductPrices []AgroproductPrice
	db.Find(&AgroproductPrices)

	return AgroproductPrices, nil
}

// Update updates one AgroproductPrice in the database, using the information
// stored in the receiver u
func (app *AgroproductPrice) Update() error {
	db.Save(&app)

	return nil
}

func (app *AgroproductPrice) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&AgroproductPrice{}).Error; err != nil {
		return err
	}

	return nil
}

func (app *AgroproductPrice) DeleteByAgroProduct(agroProductID string) error {
	if err := db.Unscoped().Where("\"agroproductID\" = ?", agroProductID).Delete(&AgroproductPrice{}).Error; err != nil {
		return err
	}

	return nil
}
