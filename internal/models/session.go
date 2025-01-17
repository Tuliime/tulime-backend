package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (cr *Session) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (s *Session) Create(session Session) (Session, error) {
	result := db.Create(&session)

	if result.Error != nil {
		return session, result.Error
	}
	return session, nil
}

func (s *Session) FindOne(id string) (Session, error) {
	var session Session
	db.First(&session, "id = ?", id)

	return session, nil
}

func (s *Session) FindByAccessToken(accessToken string) (Session, error) {
	var session Session
	db.First(&session, "\"accessToken\" = ?", accessToken)

	return session, nil
}

func (s *Session) FindByRefreshToken(refreshToken string) (Session, error) {
	var session Session
	db.First(&session, "\"refreshToken\" = ?", refreshToken)

	return session, nil
}

func (s *Session) FindByUser(userID string, limit float64, cursor string) ([]Session, error) {
	var session []Session

	query := db.Order("\"createdAt\" DESC").Where("\"userID\" = ?", cursor).Limit(int(limit))

	if cursor != "" {
		var lastSession Session
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastSession).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createdAt\" < ?", lastSession.CreatedAt)
	}

	if err := query.Find(&session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Session) FindAll(limit float64, cursor string) ([]Session, error) {
	var sessions []Session

	query := db.Order("\"createdAt\" DESC").Limit(int(limit))

	if cursor != "" {
		var lastSession Session
		if err := db.Select("\"createdAt\"").Where("id = ?", cursor).First(&lastSession).Error; err != nil {
			return nil, err
		}
		query = query.Where("\"createdAt\" < ?", lastSession.CreatedAt)
	}

	if err := query.Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

// Update updates one Session in the database, using the information
// stored in the receiver u
func (s *Session) Update() (Session, error) {
	db.Save(&s)

	session, err := s.FindOne(s.ID)
	if err != nil {
		return session, err
	}

	return session, nil
}
