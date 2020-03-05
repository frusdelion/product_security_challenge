package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	gorm.Model

	Email         string `gorm:"unique;not null"`
	EmailVerified bool   `gorm:"default:'false'"`
	Password      string `gorm:"not null"`
	FirstName     string
	LastName      string
	Username      string `gorm:"unique;not null"`
}

func (user *User) BeforeSave(scope *gorm.Scope) error {

	if user.Password != "" {
		if pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err == nil {
			if err := scope.SetColumn("password", pw); err != nil {
				return err
			}
		}
	}

	return nil
}

func (user *User) ComparePassword(givenPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(givenPassword)); err == nil {
		return true
	}

	return false
}

type FailedLogin struct {
	BaseModel
	UserID             uint
	User               User
	BrowserFingerprint string
	BrowserUserAgent   string
	IPAddress          string
	Attempts           int `gorm:"default:'0'"`
	BannedUntil        *time.Time
}

type AuthenticateUser struct {
	Username           string `form:"username" binding:"required,alphanum"`
	Password           string `form:"password" binding:"required"`
	BrowserFingerprint string `form:"browser_fingerprint" binding:"required"`
	BrowserUserAgent   string
	IPAddress          string
}

type AuthenticateClaims struct {
	UserID uint   `binding:"required"`
	Email  string `binding:"required"`
	jwt.StandardClaims
}

type AuthenticatedResult struct {
	AccessToken string
	User        *User
}

type AuthenticationRegister struct {
	Email              string `form:"email" binding:"required,email"`
	Username           string `form:"username" binding:"required,alphanum"`
	Password           string `form:"password" binding:"required,min=7"`
	ConfirmPassword    string `form:"confirm_password" binding:"required,min=7,eqfield=Password"`
	FirstName          string `form:"first_name" binding:"required,alpha"`
	LastName           string `form:"last_name" binding:"required,alpha"`
	BrowserFingerprint string `form:"browser_fingerprint" binding:"required"`
	BrowserUserAgent   string
	IPAddress          string
}

type AuthenticationForgot struct {
	Email string `form:"email" binding:"required,email"`
}

type AuthenticationNewPassword struct {
	Password        string `form:"password" binding:"required,min=7"`
	ConfirmPassword string `form:"confirm_password" binding:"required,min=7,eqfield=Password"`
}
