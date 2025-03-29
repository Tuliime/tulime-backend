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

func (ad *Advert) FindOne(id string) (Advert, error) {
	var advert Advert
	db.Preload("AdvertImage").Preload("Store").First(&advert, "id = ?", id)

	return advert, nil
}

func (ad *Advert) FindByName(name string) (Advert, error) {
	var advert Advert
	db.First(&advert, "name = ?", name)

	return advert, nil
}

func (ad *Advert) FindByUSer(userID string, limit float64, cursor string) ([]Advert, error) {
	var adverts []Advert
	query := db.Order("\"createAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastAdvert Advert
		if err := db.Select("\"createAt\"").Where("id = ?", cursor).First(&lastAdvert).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createAt\" < ?", lastAdvert.CreatedAt)
	}
	if err := query.Where("\"userID\" = ?", userID).Find(&adverts).Error; err != nil {
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
	query := db.Order("\"createAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastAdvert Advert
		if err := db.Select("\"createAt\"").Where("id = ?", cursor).First(&lastAdvert).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"createAt\" <= ?", lastAdvert.CreatedAt)
		} else {
			query = query.Where("\"createAt\" < ?", lastAdvert.CreatedAt)
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
