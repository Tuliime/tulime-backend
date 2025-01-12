package models

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := u.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func (u *User) Create(user User) (string, error) {
	result := db.Create(&user)

	if result.Error != nil {
		return "", result.Error
	}
	return user.ID, nil
}

func (u *User) FindOne(id string) (User, error) {
	var user User
	db.First(&user, "id = ?", id)

	return user, nil
}

func (u *User) FindByTelNumber(telNumber int) (User, error) {
	var user User
	db.First(&user, "\"telNumber\" = ?", telNumber)

	return user, nil
}

func (u *User) FindAll() ([]User, error) {
	var users []User
	db.Find(&users)

	return users, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update() (User, error) {
	db.Save(&u)

	user, err := u.FindOne(u.ID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) Delete(id string) error {
	db.Delete(&User{}, id)

	return nil
}

// ResetPassword is the method updates user's password in db
func (u *User) ResetPassword(password string) error {
	hashedPassword, err := u.HashPassword(password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	db.Save(&u)

	return nil
}

func (u *User) PasswordMatches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// Converts plain text password into hashed string
func (u *User) HashPassword(plainTextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (u *User) ValidRole(role string) bool {
	roles := []string{"user", "sys_admin"}

	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func (u *User) SetRole(role string) error {
	isValidRole := u.ValidRole(role)

	if !isValidRole {
		return errors.New("invalid user role")
	}

	u.Role = role
	return nil
}
