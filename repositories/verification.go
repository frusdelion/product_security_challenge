package repositories

import (
	"errors"
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	"github.com/jinzhu/gorm"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/snwfdhmp/errlog"
	"time"
)

var (
	ErrExpiredVerificationCode = errors.New("expired verification code")
)

type VerificationRepository interface {
	CreateVerificationCode(user *models.User, purposeType models.PurposeType, expiryDuration time.Duration) (*models.Verification, error)
	RetrieveVerificationCode(user *models.User, purposeType models.PurposeType) (*models.Verification, error)
	RetrieveVerificationByCode(Code string, purposeType models.PurposeType) (*models.Verification, error)
	DeleteVerificationCode(verification *models.Verification) error
}

func NewVerificationRepository(db *gorm.DB) VerificationRepository {
	return &verificationRepository{db: db}
}

type verificationRepository struct {
	db *gorm.DB
}

func (v verificationRepository) RetrieveVerificationByCode(Code string, purposeType models.PurposeType) (*models.Verification, error) {
	vr := &models.Verification{}
	if err := v.db.Preload("User").Find(vr, models.Verification{Code: Code, Purpose: purposeType}).Error; errlog.Debug(err) {
		return nil, err
	}

	if time.Now().After(vr.ExpiresOn) {
		return vr, ErrExpiredVerificationCode
	}

	return vr, nil
}

func (v verificationRepository) CreateVerificationCode(user *models.User, purposeType models.PurposeType, expiryDuration time.Duration) (*models.Verification, error) {
	vr, err := v.RetrieveVerificationCode(user, purposeType)
	if errlog.Debug(err) && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	if err == nil {
		if err := v.DeleteVerificationCode(vr); errlog.Debug(err) {
			return nil, err
		}
	}

	var code string
	switch purposeType {
	case models.PurposeMFA:
		code, err = gonanoid.Generate("0123456789", 8)
		break
	default:
		code, err = gonanoid.Nanoid(48)

	}
	if errlog.Debug(err) {
		return nil, err
	}

	vrNew := &models.Verification{UserID: user.ID, Code: code, Purpose: purposeType, ExpiresOn: time.Now().Add(expiryDuration)}
	return vrNew, v.db.Create(vrNew).Error
}

func (v verificationRepository) RetrieveVerificationCode(user *models.User, purposeType models.PurposeType) (*models.Verification, error) {
	vrNew := &models.Verification{}
	if err := v.db.Preload("User").Find(vrNew, models.Verification{UserID: user.ID, Purpose: purposeType}).Error; errlog.Debug(err) {
		return nil, err
	}

	if time.Now().After(vrNew.ExpiresOn) {
		return nil, ErrExpiredVerificationCode
	}

	return vrNew, nil
}

func (v verificationRepository) DeleteVerificationCode(verification *models.Verification) error {
	return v.db.Model(&models.Verification{}).Delete(&models.Verification{}, verification.ID).Error
}
