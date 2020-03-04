package repositories

import (
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	"github.com/jinzhu/gorm"
	"github.com/snwfdhmp/errlog"
	"time"
)

type AuthenticationRepository interface {
	GetFailedLogin(user *models.User, browserFingerprint string, browserUA string, ipAddress string) (*models.FailedLogin, error)
	RecordFailedLogin(fl *models.FailedLogin) error
	ResetFailedLogins(user *models.User, browserFingerprint string, browserUA string, ipAddress string) error
	BanUserUntil(failedLogin *models.FailedLogin, duration time.Duration) error
}

func NewAuthenticationRepository(db *gorm.DB) AuthenticationRepository {
	return &authenticationRepository{db: db}
}

type authenticationRepository struct {
	db *gorm.DB
}

func (a authenticationRepository) BanUserUntil(failedLogin *models.FailedLogin, duration time.Duration) error {
	durationTime := time.Now().Add(duration)
	return a.db.Model(failedLogin).Updates(models.FailedLogin{
		BannedUntil: &durationTime,
	}).Error
}

func (a authenticationRepository) GetFailedLogin(user *models.User, browserFingerprint string, browserUA string, ipAddress string) (*models.FailedLogin, error) {
	failedLogin := &models.FailedLogin{}
	return failedLogin, a.db.FirstOrCreate(failedLogin, models.FailedLogin{UserID: user.ID, BrowserFingerprint: browserFingerprint, BrowserUserAgent: browserUA, IPAddress: ipAddress}).Error
}

func (a authenticationRepository) RecordFailedLogin(failedLogins *models.FailedLogin) error {
	failedLogins.Attempts += 1
	return a.db.Model(failedLogins).Updates(models.FailedLogin{Attempts: failedLogins.Attempts}).Error
}

func (a authenticationRepository) ResetFailedLogins(user *models.User, browserFingerprint string, browserUA string, ipAddress string) error {
	if err := a.db.Model(&models.FailedLogin{}).Delete(&models.FailedLogin{}, models.FailedLogin{UserID: user.ID, BrowserFingerprint: browserFingerprint, BrowserUserAgent: browserUA, IPAddress: ipAddress}).Error; errlog.Debug(err) {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		return err
	}

	return nil
}
