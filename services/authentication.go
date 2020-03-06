package services

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	"github.com/frusdelion/zendesk-product_security_challenge/repositories"
	server2 "github.com/frusdelion/zendesk-product_security_challenge/server"
	"github.com/snwfdhmp/errlog"
	"time"
)

var (
	ErrBadCredentials   = errors.New("bad credentials")
	ErrAccountLocked    = errors.New("account locked")
	ErrInvalidToken     = errors.New("invalid token")
	ErrEmailNotVerified = errors.New("unverified email")
)

type AuthenticationService interface {
	AuthenticateUser(login *models.AuthenticateUser) (*models.AuthenticateClaims, *models.User, error)
	GenerateJWTKey(claim *models.AuthenticateClaims) (string, error)
	Register(register *models.AuthenticationRegister) (*models.User, error)
	RetrieveUserFromClaims(claim *models.AuthenticateClaims) (*models.User, error)
	GetUserFromJWT(token string) (*models.User, error)
	ValidateJWTKey(tokenString string) (*models.AuthenticateClaims, error)
}

func NewAuthenticationService(r repositories.AuthenticationRepository, u repositories.UserRepository, v VerificationService, s server2.Server) AuthenticationService {
	return &authenticationService{r: r, u: u, s: s, v: v}
}

type authenticationService struct {
	r repositories.AuthenticationRepository
	u repositories.UserRepository
	s server2.Server
	v VerificationService
}

func (a authenticationService) RetrieveUserFromClaims(claim *models.AuthenticateClaims) (*models.User, error) {
	if claim == nil {
		return nil, errors.New("claim cannot be nil")
	}

	return a.u.FindUserByUserID(claim.UserID)
}

func (a authenticationService) AuthenticateUser(login *models.AuthenticateUser) (*models.AuthenticateClaims, *models.User, error) {
	user, err := a.u.FindUserByUsername(login.Username)
	if errlog.Debug(err) {
		return nil, nil, err
	}

	// Check if password given is correct
	if user.ComparePassword(login.Password) {
		if !user.EmailVerified {
			return nil, nil, ErrEmailNotVerified
		}

		// If within the acceptable attempt count, we can reset upon successful
		if err := a.r.ResetFailedLogins(user, login.BrowserFingerprint, login.BrowserUserAgent, login.IPAddress); errlog.Debug(err) {
			return nil, nil, err
		}

		expirationTime := time.Now().Add(time.Duration(a.s.Config().ValidJWTLengthHours) * time.Hour)

		a.s.Log().Infof("[SECURITY EVENT] log in user for %d, ip:%s, browserUA:%s, browserFingerprint:%s", user.ID, login.IPAddress, login.BrowserUserAgent, login.BrowserFingerprint)

		return &models.AuthenticateClaims{
			UserID: user.ID,
			Email:  user.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
				Issuer:    a.s.Config().Domain,
				Audience:  "website",
				Subject:   fmt.Sprint(user.ID),
			},
		}, user, nil
	}

	// Get failed login attempts
	fl, err := a.r.GetFailedLogin(user, login.BrowserFingerprint, login.BrowserUserAgent, login.IPAddress)
	if errlog.Debug(err) {
		return nil, nil, err
	}

	// If a ban is placed, ensure that the ban is enforced
	if fl.BannedUntil != nil && time.Now().Before(*fl.BannedUntil) {
		return nil, nil, ErrAccountLocked
	} else if fl.BannedUntil != nil && time.Now().After(*fl.BannedUntil) {
		// Ban has been lifted, reset failed logins and try again
		a.s.Log().Info("[SECURITY EVENT] reset banned client for user id: %d, ip:%s, browserUA:%s, browserFingerprint:%s", user.ID, login.IPAddress, login.BrowserUserAgent, login.BrowserFingerprint)

		if err := a.r.ResetFailedLogins(user, login.BrowserFingerprint, login.BrowserUserAgent, login.IPAddress); errlog.Debug(err) {
			return nil, nil, err
		}

		fl.Attempts = 0
	}

	// If maximum attempts reached, lock account
	if fl.Attempts+1 > a.s.Config().MaximumFailedAttempts {
		if err := a.r.BanUserUntil(fl, 5*time.Minute); errlog.Debug(err) {
			return nil, nil, err
		}
		a.s.Log().Info("[SECURITY EVENT] banned client for user id: %d, ip:%s, browserUA:%s, browserFingerprint:%s, expiresOn:%d", user.ID, login.IPAddress, login.BrowserUserAgent, login.BrowserFingerprint, time.Now().Add(5*time.Minute).Unix())

		return nil, nil, ErrAccountLocked
	}

	a.s.Log().Info("[SECURITY EVENT] failed login: %d, ip:%s, browserUA:%s, browserFingerprint:%s", user.ID, login.IPAddress, login.BrowserUserAgent, login.BrowserFingerprint)

	if err := a.r.RecordFailedLogin(fl); errlog.Debug(err) {
		return nil, nil, err
	}

	return nil, nil, ErrBadCredentials

}

func (a authenticationService) GetUserFromJWT(token string) (*models.User, error) {
	claims, err := a.ValidateJWTKey(token)
	if errlog.Debug(err) {
		return nil, err
	}

	return a.RetrieveUserFromClaims(claims)
}

func (a authenticationService) ValidateJWTKey(tokenString string) (*models.AuthenticateClaims, error) {
	jwtKey := []byte(a.s.Config().JwtSecret)

	claims := &models.AuthenticateClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	if errlog.Debug(err) {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (a authenticationService) GenerateJWTKey(claim *models.AuthenticateClaims) (string, error) {
	jwtKey := []byte(a.s.Config().JwtSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)

	a.s.Log().Infof("[SECURITY EVENT] generated user key for %d", claim.UserID)

	return tokenString, err
}

func (a authenticationService) Register(register *models.AuthenticationRegister) (*models.User, error) {
	usr, err := a.u.CreateUser(&models.User{
		Email:     register.Email,
		Password:  register.Password,
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Username:  register.Username,
	})
	if errlog.Debug(err) {
		return nil, err
	}

	if err := a.v.SendRegistrationVerification(usr); errlog.Debug(err) {
		return nil, err
	}

	a.s.Log().Infof("[SECURITY EVENT] registered new user for %d, ip:%s, browserUA:%s, browserFingerprint:%s", usr.ID, register.IPAddress, register.BrowserUserAgent, register.BrowserFingerprint)

	return usr, nil
}
