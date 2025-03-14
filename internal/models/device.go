package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (d *Device) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (d *Device) Create(device Device) (Device, error) {
	var devices []Device
	if err := db.Create(&device).Error; err != nil {
		return device, err
	}

	// TODO: To improve algorithm for writing device data to cache
	db.Find(&devices, "\"userID\" = ?", device.UserID)
	if len(devices) > 0 {
		d.writeToCache(devices)
	}

	return device, nil
}

func (d *Device) FindOne(id string) (Device, error) {
	var device Device
	db.First(&device, "id = ?", id)

	return device, nil
}

func (d *Device) FindByTokenAndUser(token, userID string) (Device, error) {
	var device Device

	db.Find(&device, "token = ? AND \"userID\" = ?", token, userID)
	return device, nil
}

func (d *Device) FindByUser(userID string) ([]Device, error) {
	var devices []Device

	devices, err := d.readFromCache(userID)
	if err != nil {
		return devices, err
	}
	if len(devices) > 0 {
		return devices, nil
	}
	db.Find(&devices, "\"userID\" = ?", userID)

	return devices, nil
}

func (d *Device) FindAll() ([]Device, error) {
	var devices []Device
	devices, err := d.readAllFromCache()
	if err != nil {
		return devices, err
	}
	if len(devices) > 0 {
		return devices, nil
	}
	db.Find(&devices)

	return devices, nil
}

// Update updates one user in the database, using the information
// stored in the receiver d and update the cache
func (d *Device) Update() (Device, error) {
	var devices []Device
	if err := db.Save(&d).Error; err != nil {
		return *d, err
	}

	db.Find(&devices, "\"userID\" = ?", d.UserID)
	if len(devices) > 0 {
		d.writeToCache(devices)
	}

	return *d, nil
}

func (d *Device) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Device{}).Error; err != nil {
		return err
	}

	return nil
}

func (d *Device) DeleteByUser(userID string) error {
	if err := db.Unscoped().Where("\"userID\" = ?", userID).Delete(&Device{}).Error; err != nil {
		return err
	}
	if err := d.deleteFromCache(userID); err != nil {
		return err
	}

	return nil
}
