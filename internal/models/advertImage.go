package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (adi *AdvertImage) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (adi *AdvertImage) Create(image AdvertImage) (AdvertImage, error) {
	result := db.Create(&image)

	if result.Error != nil {
		return image, result.Error
	}
	return image, nil
}

func (adi *AdvertImage) CreateMany(images []AdvertImage) ([]AdvertImage, error) {
	result := db.Create(&images)

	if result.Error != nil {
		return images, result.Error
	}
	return images, nil
}

func (adi *AdvertImage) FindOne(id string) (AdvertImage, error) {
	var image AdvertImage
	db.First(&image, "id = ?", id)

	return image, nil
}

func (adi *AdvertImage) FindByAdvert(advertID string, limit float64, cursor string) ([]AdvertImage, error) {
	var images []AdvertImage
	query := db.Order("\"createAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastImage AdvertImage
		if err := db.Select("\"createAt\"").Where("id = ?", cursor).First(&lastImage).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createAt\" < ?", lastImage.CreatedAt)
	}
	if err := query.Where("\"advertID\" = ?", advertID).Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (adi *AdvertImage) FindAll(limit float64, cursor string, includeCursor bool, direction string) ([]AdvertImage, error) {
	var images []AdvertImage

	if direction == "FORWARD" {
		advertsInAscOrder, err := adi.FindAllInAscOrder(limit, cursor, includeCursor)
		if err != nil {
			return images, err
		}
		images = advertsInAscOrder

	} else if direction == "BACKWARD" {
		advertsInDescOrder, err := adi.FindAllInDescOrder(limit, cursor, includeCursor)
		if err != nil {
			return images, err
		}
		images = advertsInDescOrder
	} else {
		return images, errors.New("invalid direction value")
	}

	return images, nil
}

func (adi *AdvertImage) FindAllInDescOrder(limit float64, cursor string, includeCursor bool) ([]AdvertImage, error) {
	var images []AdvertImage
	query := db.Order("\"createAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastImage Advert
		if err := db.Select("\"createAt\"").Where("id = ?", cursor).First(&lastImage).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"createAt\" <= ?", lastImage.CreatedAt)
		} else {
			query = query.Where("\"createAt\" < ?", lastImage.CreatedAt)
		}
	}

	if err := query.Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (adi *AdvertImage) FindAllInAscOrder(limit float64, cursor string, includeCursor bool) ([]AdvertImage, error) {
	var images []AdvertImage

	query := db.Order("\"createdAt\" ASC").Limit(int(limit))

	if cursor != "" {
		var lastImage AdvertImage
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastImage).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"createdAt\" >= ?", lastImage.CreatedAt)
		} else {
			query = query.Where("\"createdAt\" > ?", lastImage.CreatedAt)
		}
	}

	if err := query.Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (adi *AdvertImage) Update() (AdvertImage, error) {
	db.Save(&adi)

	return *adi, nil
}

// TODO: consider soft deleting the AdvertImage
func (ad *AdvertImage) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Advert{}).Error; err != nil {
		return err
	}

	return nil
}
