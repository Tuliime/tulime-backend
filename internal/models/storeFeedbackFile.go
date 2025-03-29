package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (fl *StoreFeedbackFile) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (fl *StoreFeedbackFile) Create(file StoreFeedbackFile) (StoreFeedbackFile, error) {
	result := db.Create(&file)

	if result.Error != nil {
		return file, result.Error
	}
	return file, nil
}

func (fl *StoreFeedbackFile) CreateMany(files []StoreFeedbackFile) ([]StoreFeedbackFile, error) {
	result := db.Create(&files)

	if result.Error != nil {
		return files, result.Error
	}
	return files, nil
}

func (fl *StoreFeedbackFile) FindOne(id string) (StoreFeedbackFile, error) {
	var file StoreFeedbackFile
	db.First(&file, "id = ?", id)

	return file, nil
}

// Update updates one StoreFeedbackFile in the database, using the information
// stored in the receiver u
func (fl *StoreFeedbackFile) Update() (StoreFeedbackFile, error) {
	db.Save(&fl)

	return *fl, nil
}

func (fl *StoreFeedbackFile) Delete(id string) error {
	//   TODO: to softly delete StoreFeedbackFile
	return nil
}
