package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ols *OnlineStatus) BeforeSave(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (ols *OnlineStatus) FindAll() ([]OnlineStatus, error) {
	//TODO: To fetch statuses < 1min
	var statuses []OnlineStatus
	db.Find(&statuses)

	return statuses, nil
}

// Update updates one OnlineStatus in the database, using the information
// stored in the receiver ols
func (ols *OnlineStatus) Update() (OnlineStatus, error) {
	db.Save(&ols)
	return *ols, nil
}
