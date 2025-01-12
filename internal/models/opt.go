package models

import (
	"errors"
	"time"

	"github.com/Tuliime/tulime-backend/internal/pkg"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (o *OTP) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()

	tx.Statement.SetColumn("ID", uuid)
	return nil
}

// creates, encodes, saves and return a new OTP
func (o *OTP) Create(otp OTP) (string, error) {
	if otp.UserID == "" {
		return "", errors.New("no userID provided in the OTP struct")
	}

	randNum := pkg.NewRandomNumber()
	generatedOTP := randNum.D6()
	encodedOTP := pkg.EncodeToHexString(generatedOTP)

	otp.ExpiresAt = time.Now().Add(20 * time.Minute)
	otp.OTP = encodedOTP

	// Expire all other existing user opts
	if err := o.ExpireUserOTP(otp.UserID); err != nil {
		return "", err
	}

	result := db.Create(&otp)

	if result.Error != nil {
		return "", result.Error
	}
	return generatedOTP, nil
}

func (o *OTP) FindOne(id string) (OTP, error) {
	var opt OTP
	db.First(&opt, "id = ?", id)

	return opt, nil
}

func (o *OTP) FindByOTP(otpString string) (OTP, error) {
	var otp OTP
	encodedOTP := pkg.EncodeToHexString(otpString)
	db.Where("\"OTP\" = ?", encodedOTP).First(&otp)

	return otp, nil
}

func (o *OTP) FindByUser(userID string) ([]OTP, error) {
	var otp []OTP
	db.Find(&otp, "\"userID\" = ?", userID)

	return otp, nil
}

func (o *OTP) ExpireUserOTP(userID string) error {
	userOTPs, err := o.FindByUser(userID)
	if err != nil {
		return err
	}
	if len(userOTPs) == 0 {
		return nil
	}
	if err := db.Where("\"userID\" = ?", userID).Where("\"expiresAt\" > ?", time.Now()).Model(&OTP{}).Update("\"expiresAt\"", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (o *OTP) Update() (OTP, error) {
	db.Save(&o)

	otp, err := o.FindOne(o.ID)
	if err != nil {
		return otp, err
	}

	return otp, nil
}
