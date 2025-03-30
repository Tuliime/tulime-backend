package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Store) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (s *Store) Create(store Store) (Store, error) {
	result := db.Create(&store)

	if result.Error != nil {
		return store, result.Error
	}
	return store, nil
}

func (s *Store) FindOne(id string) (Store, error) {
	var store Store
	db.First(&store, "id = ?", id)

	return store, nil
}

func (s *Store) FindByName(name string) (Store, error) {
	var store Store
	db.First(&store, "name = ?", name)

	return store, nil
}

func (s *Store) FindByUSer(userID string, limit float64, cursor string) ([]Store, error) {
	var stores []Store
	query := db.Order("\"createdAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastStore Store
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastStore).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createdAt\" < ?", lastStore.CreatedAt)
	}
	if err := query.Where("\"userID\" = ?", userID).Find(&stores).Error; err != nil {
		return nil, err
	}

	return stores, nil
}

func (cr *Store) FindAll(limit float64, cursor string, includeCursor bool, direction string) ([]Store, error) {
	var stores []Store

	if direction == "FORWARD" {
		storesInAscOrder, err := cr.FindAllInASCOrder(limit, cursor, includeCursor)
		if err != nil {
			return stores, err
		}
		stores = storesInAscOrder

	} else if direction == "BACKWARD" {
		chatRoomsInDescOrder, err := cr.FindAllInDESCOrder(limit, cursor, includeCursor)
		if err != nil {
			return stores, err
		}
		stores = chatRoomsInDescOrder
	} else {
		return stores, errors.New("invalid direction value")
	}

	return stores, nil
}

func (s *Store) FindAllInDESCOrder(limit float64, cursor string, includeCursor bool) ([]Store, error) {
	var stores []Store
	query := db.Order("\"createdAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastStore Store
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastStore).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"createdAt\" <= ?", lastStore.CreatedAt)
		} else {
			query = query.Where("\"createdAt\" < ?", lastStore.CreatedAt)
		}
	}

	if err := query.Find(&stores).Error; err != nil {
		return nil, err
	}

	return stores, nil
}

func (s *Store) FindAllInASCOrder(limit float64, cursor string, includeCursor bool) ([]Store, error) {
	var stores []Store

	query := db.Order("\"createdAt\" ASC").Limit(int(limit))

	if cursor != "" {
		var lastStore Store
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastStore).Error; err != nil {
			return nil, err
		}
		if includeCursor {
			query = query.Where("\"createdAt\" >= ?", lastStore.CreatedAt)
		} else {
			query = query.Where("\"createdAt\" > ?", lastStore.CreatedAt)
		}
	}

	if err := query.Find(&stores).Error; err != nil {
		return nil, err
	}

	return stores, nil
}

func (s *Store) Update() (Store, error) {
	db.Save(&s)

	return *s, nil
}

// TODO: consider soft deleting the store
func (s *Store) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&Store{}).Error; err != nil {
		return err
	}

	return nil
}
