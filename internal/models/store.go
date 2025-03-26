package models

import (
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

func (s *Store) FindByUSer(userID string) (Store, error) {
	var store Store
	db.Find(&store, "\"userID\" = ?", userID)

	return store, nil
}

func (s *Store) FindAll(limit float64, cursor string) ([]Store, error) {
	var store []Store
	query := db.Order("\"updatedAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastStore Store
		if err := db.Select("\"updatedAt\"").Where("id = ?", cursor).First(&lastStore).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"updatedAt\" < ?", lastStore.UpdatedAt)
	}

	query.Find(&store)

	return store, nil
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
