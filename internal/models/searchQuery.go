package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (sq *SearchQuery) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (sq *SearchQuery) Create(searchQuery SearchQuery) (SearchQuery, error) {
	result := db.Create(&searchQuery)

	if result.Error != nil {
		return searchQuery, result.Error
	}
	return searchQuery, nil
}

func (sq *SearchQuery) FindOne(id string) (SearchQuery, error) {
	var query SearchQuery
	db.First(&query, "id = ?", id)

	return query, nil
}

func (n *SearchQuery) FindAll(limit float64, cursor string) ([]SearchQuery, error) {
	var searchQueries []SearchQuery
	query := db.Order("\"updatedAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastSearchQuery SearchQuery
		if err := db.Select("\"updatedAt\"").Where("id = ?", cursor).First(&lastSearchQuery).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"updatedAt\" < ?", lastSearchQuery.UpdatedAt)
	}
	query.Find(&searchQueries)

	return searchQueries, nil
}

func (sq *SearchQuery) Update() (SearchQuery, error) {
	db.Save(&sq)
	return *sq, nil
}

func (sq *SearchQuery) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&SearchQuery{}).Error; err != nil {
		return err
	}

	return nil
}
