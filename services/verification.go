package services

import (
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	"github.com/frusdelion/zendesk-product_security_challenge/repositories"
	server2 "github.com/frusdelion/zendesk-product_security_challenge/server"
	"github.com/jinzhu/gorm"
	"github.com/snwfdhmp/errlog"
	"time"
)

type VerificationService interface {
	SendRegistrationVerification(user *models.User) error
	SendForgotPasswordVerification(user *models.User, userAgent string, ipAddress string) error
	SendMFAVerification(user *models.User, userAgent string, ipAddress string) error

	VerifyRegistration(Code string) (*models.Verification, error)
	VerifyForgotPassword(Code string) (*models.Verification, error)
	VerifyMFA(Code string, user *models.User) (*models.Verification, error)

	DeleteVerification(vr *models.Verification) error
}

func NewVerificationService(r repositories.VerificationRepository, c CommunicationsService, s server2.Server) VerificationService {
	return &verificationService{r: r, s: s, c: c}
}

type verificationService struct {
	r repositories.VerificationRepository
	s server2.Server
	c CommunicationsService
}

func (v verificationService) VerifyRegistration(Code string) (*models.Verification, error) {
	return v.r.RetrieveVerificationByCode(Code, models.PurposeRegistration)
}

func (v verificationService) VerifyForgotPassword(Code string) (*models.Verification, error) {
	return v.r.RetrieveVerificationByCode(Code, models.PurposeForgotPassword)
}

func (v verificationService) VerifyMFA(Code string, user *models.User) (*models.Verification, error) {
	vr, err := v.r.RetrieveVerificationByCode(Code, models.PurposeMFA)
	if errlog.Debug(err) {
		return nil, err
	}

	if vr.UserID != user.ID {
		return nil, gorm.ErrRecordNotFound
	}

	return vr, err
}

func (v verificationService) DeleteVerification(vr *models.Verification) error {
	return v.r.DeleteVerificationCode(vr)
}

func (v verificationService) SendRegistrationVerification(user *models.User) error {
	vc, err := v.r.CreateVerificationCode(user, models.PurposeRegistration, 1*time.Hour)
	if errlog.Debug(err) {
		return err
	}

	if err := v.c.SendRegistrationValidationEmail(user, vc.Code, vc.ExpiresOn); errlog.Debug(err) {
		return err
	}

	v.s.Log().Infof("[SECURITY EVENT] registration verification for %d", vc.UserID)

	return nil
}

func (v verificationService) SendForgotPasswordVerification(user *models.User, userAgent string, ipAddress string) error {
	vc, err := v.r.CreateVerificationCode(user, models.PurposeForgotPassword, 1*time.Hour)
	if errlog.Debug(err) {
		return err
	}

	if err := v.c.SendPasswordForgotEmail(user, vc.Code, vc.ExpiresOn, userAgent, ipAddress); errlog.Debug(err) {
		return err
	}

	v.s.Log().Infof("[SECURITY EVENT] forgot password for %d, ip:%s, browserUA:%s", vc.UserID, ipAddress, userAgent)

	return nil
}

func (v verificationService) SendMFAVerification(user *models.User, userAgent string, ipAddress string) error {
	vc, err := v.r.CreateVerificationCode(user, models.PurposeMFA, 1*time.Hour)
	if errlog.Debug(err) {
		return err
	}

	if err := v.c.SendMFAEmail(user, vc.Code, vc.ExpiresOn, userAgent, ipAddress); errlog.Debug(err) {
		return err
	}

	v.s.Log().Infof("[SECURITY EVENT] mfa verification for %d, ip:%s, browserUA:%s", vc.UserID, ipAddress, userAgent)
	return nil
}
