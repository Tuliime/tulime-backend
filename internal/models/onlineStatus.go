package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ols *OnlineStatus) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (ols *OnlineStatus) Create(onlineStatus OnlineStatus) (OnlineStatus, error) {
	if err := db.Create(&onlineStatus).Error; err != nil {
		return onlineStatus, err
	}
	return onlineStatus, nil
}

func (ols *OnlineStatus) FindByUser(userID string) (OnlineStatus, error) {
	var onlineStatus OnlineStatus
	db.Last(&onlineStatus, "\"userID\" = ?", userID)

	return onlineStatus, nil
}

func (ols *OnlineStatus) FindByUserList(userIDList []string) ([]OnlineStatus, error) {
	var statuses []OnlineStatus

	if err := db.Table("online_statuses").
		Select("DISTINCT ON (\"userID\") *").
		Where("\"userID\" IN (?)", userIDList).
		Order("\"userID\", \"createdAt\" DESC").
		Find(&statuses).Error; err != nil {
		return nil, err
	}

	return statuses, nil
}

func (ols *OnlineStatus) FindAll() ([]OnlineStatus, error) {
	var statuses []OnlineStatus
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)

	err := db.Where("\"updatedAt\" >= ?", oneMinuteAgo).Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (ols *OnlineStatus) Update() (OnlineStatus, error) {
	ols.UpdatedAt = time.Now()

	if err := db.Model(&OnlineStatus{}).Where("id = ?", ols.ID).
		UpdateColumn("\"updatedAt\"", ols.UpdatedAt).Error; err != nil {
		return *ols, err
	}
	return *ols, nil
}
