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
	// TODO: to validate the category here
	result := db.Create(&agroproduct)

	if result.Error != nil {
		return agroproduct, result.Error
	}
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
func (ap *Agroproduct) FindByCategory(category string) ([]Agroproduct, error) {
	var agroproducts []Agroproduct
	db.Find(&agroproducts, "category = ?", category)

	for i, agroproduct := range agroproducts {
		var agroProductPrice []AgroproductPrice
		db.Order("\"createdAt\" desc").Limit(1).Find(&agroProductPrice, "\"agroproductID\" = ?", agroproduct.ID)
		agroproduct.Price = agroProductPrice
		agroproducts[i] = agroproduct
	}

	return agroproducts, nil
}

func (ap *Agroproduct) FindAll(limit float64, category string, cursor string) ([]Agroproduct, error) {
	var agroproducts []Agroproduct
	query := db.Order("\"updatedAt\" DESC").Limit(int(limit))

	if category != "" {
		query.Where("category = ?", category)
	}

	if cursor != "" {
		var lastAgroproduct Agroproduct
		if err := db.Select("\"updatedAt\"").Where("id = ?", cursor).First(&lastAgroproduct).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"updatedAt\" < ?", lastAgroproduct.UpdatedAt)
	}

	query.Find(&agroproducts)

	for i, agroproduct := range agroproducts {
		var agroProductPrice []AgroproductPrice
		db.Order("\"createdAt\" desc").Limit(1).Find(&agroProductPrice, "\"agroproductID\" = ?", agroproduct.ID)
		agroproduct.Price = agroProductPrice
		agroproducts[i] = agroproduct
	}

	return agroproducts, nil
}

// Update updates one Agroproduct in the database, using the information
// stored in the receiver u
func (ap *Agroproduct) Update() (Agroproduct, error) {
	db.Save(&ap)

	agroProduct, err := ap.FindOne(ap.ID)
	if err != nil {
		return agroProduct, err
	}

	return agroProduct, nil
}

func (ap *Agroproduct) Delete(id string) error {
	agroProductPrice := AgroproductPrice{AgroproductID: id}

	if err := agroProductPrice.DeleteByAgroProduct(id); err != nil {
		return err
	}

	if err := db.Unscoped().Where("id = ?", id).Delete(&Agroproduct{}).Error; err != nil {
		return err
	}

	return nil
}

func (ap *Agroproduct) ValidCategory(category string) bool {
	categories := []string{"crop", "livestock", "poultry", "fish"}

	for _, c := range categories {
		if c == category {
			return true
		}
	}

	return false
}

// HasImagePath is temporary function to check if an agroproduct
// has imagePath value
func (ap *Agroproduct) HasImagePath(agroProduct Agroproduct) bool {
	return agroProduct.ImagePath != ""
}
