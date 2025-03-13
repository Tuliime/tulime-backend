package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (n *Notification) Create(notification Notification) (Notification, error) {
	if err := db.Create(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (n *Notification) FindOne(id string) (Notification, error) {
	var notification Notification
	db.First(&notification, "id = ?", id)

	return notification, nil
}

func (n *Notification) FindByUser(userID string) ([]Notification, error) {
	var notifications []Notification
	db.Find(&notifications, "\"userID\" = ?", userID)

	return notifications, nil
}

func (n *Notification) FindUnreadByUser(userID string) ([]Notification, error) {
	var notifications []Notification
	err := db.Where("\"userID\" = ? AND \"isRead\" = ?", userID, false).Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (n *Notification) FindUnreadByUserAndType(userID, Type string) ([]Notification, error) {
	var notifications []Notification
	err := db.Where("\"userID\" = ? AND type = ? AND \"isRead\" = ?", userID, Type, false).Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (n *Notification) FindUnreadCountByUser(userID string) (int64, error) {
	var count int64
	if err := db.Model(&Notification{}).Where("\"userID\" = ? AND \"isRead\" = ?",
		userID, false).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (n *Notification) FindUnreadCountByUserAndType(userID, Type string) (int64, error) {
	var count int64
	if err := db.Model(&Notification{}).Where("\"userID\" = ? AND type = ? AND \"isRead\" = ?",
		userID, Type, false).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (n *Notification) FindAll() ([]Notification, error) {
	var notifications []Notification
	db.Find(&notifications)

	return notifications, nil
}

// Update updates one user in the database, using the information
// stored in the receiver d
func (n *Notification) Update() (Notification, error) {
	if err := db.Save(&n).Error; err != nil {
		return *n, err
	}

	return *n, nil
}

func (n *Notification) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Notification{}).Error; err != nil {
		return err
	}

	return nil
}
