package server

import (
	"context"
	"github.com/frusdelion/zendesk-product_security_challenge/config"
	"github.com/jinzhu/gorm"
	jwemail "github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type Server interface {
	Log() *logrus.Logger
	DB() *gorm.DB
	Email() *jwemail.Pool
	Validator() *validator.Validate
	Routes()
	Context() context.Context
	ContextCancel() context.CancelFunc
	Config() config.ServerConfiguration

	AutoMigrate()
	GreeterBanner()

	Run()
}
