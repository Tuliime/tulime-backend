package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (n *News) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (n *News) Create(news News) (News, error) {
	result := db.Create(&news)

	if result.Error != nil {
		return news, result.Error
	}
	return news, nil
}

func (n *News) FindOne(id string) (News, error) {
	var news News
	db.First(&news, "id = ?", id)

	return news, nil
}

// TODO: add pagination for all select queries that return many results
func (n *News) FindByCategory(category string) ([]News, error) {
	var news []News
	db.Find(&news, "category = ?", category)

	return news, nil
}

func (n *News) FindAll(limit float64, category string, cursor string) ([]News, error) {
	var news []News
	query := db.Order("\"updatedAt\" DESC").Limit(int(limit))

	if category != "" {
		query.Where("category = ?", category)
	}

	if cursor != "" {
		var lastNews News
		if err := db.Select("\"updatedAt\"").Where("id = ?", cursor).First(&lastNews).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"updatedAt\" < ?", lastNews.UpdatedAt)
	}
	query.Find(&news)

	return news, nil
}

func (n *News) Search(searchQuery string) ([]News, error) {
	var news []News

	query := db.Order("\"createdAt\" DESC").Limit(20)

	query = query.Where("title ILIKE ? OR description ILIKE ?",
		"%"+searchQuery+"%", "%"+searchQuery+"%")

	if err := query.Find(&news).Error; err != nil {
		return news, err
	}

	return news, nil
}

// Update updates one News in the database, using the information
// stored in the receiver u
func (n *News) Update() (News, error) {
	db.Save(&n)

	news, err := n.FindOne(n.ID)
	if err != nil {
		return news, err
	}

	return news, nil
}

func (n *News) Delete(id string) error {
	if err := db.Unscoped().Where("id = ?", id).Delete(&News{}).Error; err != nil {
		return err
	}

	return nil
}
