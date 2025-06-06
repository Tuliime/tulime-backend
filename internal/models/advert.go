package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ad *Advert) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (ad *Advert) Create(advert Advert) (Advert, error) {
	result := db.Create(&advert)

	if result.Error != nil {
		return advert, result.Error
	}
	return advert, nil
}

// Find retrieves one record matching provided id
// and doesn't include joins
func (ad *Advert) Find(id string) (Advert, error) {
	var advert Advert
	db.Where("id = ?", id).First(&advert)
	return advert, nil
}

// Find retrieves one record matching provided id
// and includes joins via preload
func (ad *Advert) FindOne(id string) (Advert, error) {
	var advert Advert
	query := db.Preload("Store").Preload("AdvertImage").Preload("AdvertPrice").Preload("AdvertInventory")
	query.Where("id = ?", id).First(&advert)
	return advert, nil
}

func (ad *Advert) FindByName(name string) (Advert, error) {
	var advert Advert
	db.Preload("AdvertImage").First(&advert, "name = ?", name)

	return advert, nil
}

func (ad *Advert) FindByUser(userID string, limit float64, cursor string) ([]Advert, error) {
	var adverts []Advert
	query := db.Order("\"createdAt\" DESC").Limit(int(limit))
	query = query.Preload("AdvertImage").Preload("AdvertPrice").Preload("AdvertInventory")

	if cursor != "" {
		var lastAdvert Advert
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastAdvert).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createdAt\" < ?", lastAdvert.CreatedAt)
	}
	if err := query.Where("\"userID\" = ?", userID).Find(&adverts).Error; err != nil {
		return nil, err
	}

	return adverts, nil
}

func (ad *Advert) FindByStore(storeID string, limit float64, cursor string) ([]Advert, error) {
	var adverts []Advert
	query := db.Order("\"createdAt\" DESC").Limit(int(limit))
	query = query.Preload("AdvertImage").Preload("AdvertPrice").Preload("AdvertInventory")

	if cursor != "" {
		var lastAdvert Advert
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastAdvert).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createdAt\" < ?", lastAdvert.CreatedAt)
	}
	if err := query.Where("\"storeID\" = ?", storeID).Find(&adverts).Error; err != nil {
		return nil, err
	}

	return adverts, nil
}

func (ad *Advert) FindAll(limit float64, cursor string, includeCursor bool, direction string) ([]Advert, error) {
	var adverts []Advert

	if direction == "FORWARD" {
		advertsInAscOrder, err := ad.FindAllInAscOrder(limit, cursor, includeCursor)
		if err != nil {
			return adverts, err
		}
		adverts = advertsInAscOrder

	} else if direction == "BACKWARD" {
		advertsInDescOrder, err := ad.FindAllInDescOrder(limit, cursor, includeCursor)
		if err != nil {
			return adverts, err
		}
		adverts = advertsInDescOrder
	} else {
		return adverts, errors.New("invalid direction value")
	}

	return adverts, nil
}

func (ad *Advert) FindAllInDescOrder(limit float64, cursor string, includeCursor bool) ([]Advert, error) {
	var adverts []Advert
	query := db.Order("\"createdAt\" DESC").Limit(int(limit))
	query = query.Preload("AdvertImage").Preload("AdvertPrice").Preload("AdvertInventory")

	if cursor != "" {
		var lastAdvert Advert
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastAdvert).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"createdAt\" <= ?", lastAdvert.CreatedAt)
		} else {
			query = query.Where("\"createdAt\" < ?", lastAdvert.CreatedAt)
		}
	}

	if err := query.Find(&adverts).Error; err != nil {
		return nil, err
	}

	return adverts, nil
}

func (ad *Advert) FindAllInAscOrder(limit float64, cursor string, includeCursor bool) ([]Advert, error) {
	var adverts []Advert

	query := db.Order("\"createdAt\" ASC").Limit(int(limit))
	query = query.Preload("AdvertImage").Preload("AdvertPrice").Preload("AdvertInventory")

	if cursor != "" {
		var lastAdvert Advert
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastAdvert).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"createdAt\" >= ?", lastAdvert.CreatedAt)
		} else {
			query = query.Where("\"createdAt\" > ?", lastAdvert.CreatedAt)
		}
	}

	if err := query.Find(&adverts).Error; err != nil {
		return nil, err
	}

	return adverts, nil
}

func (ad *Advert) Search(searchQuery string) ([]Advert, error) {
	var adverts []Advert

	query := db.Order("\"createdAt\" DESC").Limit(20)
	query = query.Preload("AdvertImage").Preload("AdvertPrice").Preload("AdvertInventory")

	query = query.Where("\"isPublished\" = ? AND (\"productName\" ILIKE ? OR \"productDescription\" ILIKE ?)",
		true, "%"+searchQuery+"%", "%"+searchQuery+"%")

	if err := query.Find(&adverts).Error; err != nil {
		return adverts, err
	}
	return adverts, nil
}

func (ad *Advert) Update() (Advert, error) {
	db.Save(&ad)

	return *ad, nil
}

// TODO: consider soft deleting the Advert
func (ad *Advert) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Advert{}).Error; err != nil {
		return err
	}

	return nil
}
