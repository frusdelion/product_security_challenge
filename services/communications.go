package services

import (
	"fmt"
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	server2 "github.com/frusdelion/zendesk-product_security_challenge/server"
	jwemail "github.com/jordan-wright/email"
	"github.com/matcornic/hermes"
	"github.com/snwfdhmp/errlog"
	"time"
)

type CommunicationsService interface {
	SendMFAEmail(user *models.User, Code string, expiresOn time.Time, userAgent string, ipAddress string) error
	SendPasswordForgotEmail(user *models.User, Code string, expiresOn time.Time, userAgent string, ipAddress string) error
	SendRegistrationValidationEmail(user *models.User, Code string, expiresOn time.Time) error
}

func NewCommunicationsService(s server2.Server) CommunicationsService {
	hm := hermes.Hermes{
		Theme:         new(hermes.Default),
		TextDirection: hermes.TDLeftToRight,
		Product: hermes.Product{
			Name: "Zendesk Product Security Challenge, March 2020",
			Link: fmt.Sprintf("https://%s/", s.Config().Domain),
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}
	return &communicationsService{s: s, hm: hm}
}

type communicationsService struct {
	s  server2.Server
	hm hermes.Hermes
}

func (c communicationsService) SendMFAEmail(user *models.User, Code string, expiresOn time.Time, userAgent string, ipAddress string) error {
	email := hermes.Email{
		Body: hermes.Body{
			Name: "Multi-factor Authentication",
			Intros: []string{
				"We have received a sign in to your account. For your security, here is your sign-in code:",
			},

			Dictionary: []hermes.Entry{
				{
					Key:   "Code",
					Value: Code,
				},
				{
					Key:   "Browser",
					Value: userAgent,
				},
				{
					Key:   "IP Address",
					Value: ipAddress,
				},
			},

			Outros: []string{
				fmt.Sprintf("This code will expire at %s.", expiresOn.Format("Mon Jan 2 15:04:05 -0700 MST 2006")),
				"If this is not you, you may disregard this message.",
			},
		},
	}

	ht, err := c.hm.GenerateHTML(email)
	if errlog.Debug(err) {
		return err
	}

	tt, err := c.hm.GeneratePlainText(email)
	if errlog.Debug(err) {
		return err
	}

	return c.s.Email().Send(&jwemail.Email{
		To:      []string{user.Email},
		From:    c.s.Config().SMTPFrom,
		Subject: "Request to sign in",
		Text:    []byte(tt),
		HTML:    []byte(ht),
	}, 3*time.Second)
}

func (c communicationsService) SendPasswordForgotEmail(user *models.User, Code string, expiresOn time.Time, userAgent string, ipAddress string) error {
	email := hermes.Email{
		Body: hermes.Body{
			Name: "Password reset request",
			Intros: []string{
				"We have received a request to reset your password. The following is the details of the request:",
			},

			Dictionary: []hermes.Entry{
				{
					Key:   "Browser",
					Value: userAgent,
				},
				{
					Key:   "IP Address",
					Value: ipAddress,
				},
			},

			Actions: []hermes.Action{
				{
					Instructions: "To continue to reset your password, click on the following link:",
					Button: hermes.Button{
						Text: "Continue resetting your password",
						Link: fmt.Sprintf("https://%s/verify/resetpassword/%s", c.s.Config().Domain, Code),
					},
				},
			},

			Outros: []string{
				fmt.Sprintf("This link will expire on %s.", expiresOn.Format(expiresOn.Format("Mon Jan 2 15:04:05 -0700 MST 2006"))),
				"If this is not you, you may disregard this message.",
			},
		},
	}

	ht, err := c.hm.GenerateHTML(email)
	if errlog.Debug(err) {
		return err
	}

	tt, err := c.hm.GeneratePlainText(email)
	if errlog.Debug(err) {
		return err
	}

	return c.s.Email().Send(&jwemail.Email{
		To:      []string{user.Email},
		From:    c.s.Config().SMTPFrom,
		Subject: "Request to reset your password",
		Text:    []byte(tt),
		HTML:    []byte(ht),
	}, 3*time.Second)
}

func (c communicationsService) SendRegistrationValidationEmail(user *models.User, Code string, expiresOn time.Time) error {
	email := hermes.Email{
		Body: hermes.Body{
			Name: "Complete your Registration",
			Intros: []string{
				fmt.Sprintf("Welcome, %s!", user.FirstName),
				"We are so excited to have you!",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To continue your registration, click on the following link:",
					Button: hermes.Button{
						Text: "Continue resetting your password",
						Link: fmt.Sprintf("https://%s/verify/registration/%s", c.s.Config().Domain, Code),
					},
				},
			},

			Outros: []string{
				fmt.Sprintf("This link expires on %s.", expiresOn.Format(expiresOn.Format("Mon Jan 2 15:04:05 -0700 MST 2006"))),
				"If this is not you, you may disregard this message.",
			},
		},
	}

	ht, err := c.hm.GenerateHTML(email)
	if errlog.Debug(err) {
		return err
	}

	tt, err := c.hm.GeneratePlainText(email)
	if errlog.Debug(err) {
		return err
	}

	return c.s.Email().Send(&jwemail.Email{
		To:      []string{user.Email},
		From:    c.s.Config().SMTPFrom,
		Subject: "Complete your registration",
		Text:    []byte(tt),
		HTML:    []byte(ht),
	}, 3*time.Second)
}
