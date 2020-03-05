package main

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/frusdelion/zendesk-product_security_challenge/controllers"
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	"github.com/frusdelion/zendesk-product_security_challenge/repositories"
	"github.com/frusdelion/zendesk-product_security_challenge/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/snwfdhmp/errlog"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
	"path/filepath"
)

func (s *server) Routes() {
	//new template engine
	gvw := ginview.New(goview.Config{
		Root:      "./project",
		Extension: ".html",
		Master:    "./layouts/master",
	})

	tplFiles := rice.MustFindBox("./project")
	gvw.SetFileHandler(func(config goview.Config, tplFile string) (content string, err error) {
		path := filepath.Clean(tplFile + config.Extension)

		//s.Log().Infof("%s", path)
		content, err = tplFiles.String(path)
		if errlog.Debug(err) {
			return "", fmt.Errorf("ViewEngine render read name:%v, path:%v, error: %v", tplFile, path, err)
		}
		return
	})

	s.http.HTMLRender = gvw

	// TODO: go.rice
	s.http.StaticFS("/assets", rice.MustFindBox("./project/assets").HTTPBox())
	//s.http.Static("/assets", "./project/assets")

	ur := repositories.NewUserRepository(s.DB())
	s.Log().Info("Loaded ur service")

	us := services.NewUserService(ur)

	s.Log().Info("Loaded us service")
	vs := services.NewVerificationService(repositories.NewVerificationRepository(s.DB()), services.NewCommunicationsService(s), s)

	s.Log().Info("Loaded vs service")
	as := services.NewAuthenticationService(
		repositories.NewAuthenticationRepository(s.DB()),
		ur,
		vs,
		s,
	)
	s.Log().Info("Loaded as service")

	ac := controllers.NewAuthenticationController(s, as, us, vs)
	s.Log().Info("Loaded ac")

	r := s.http.Group("/")
	r.Use(s.middlewareFullAuthorization(as))
	{
		r.GET("/", func(c *gin.Context) {
			user, ok := c.Get("user")
			if !ok {
				c.Redirect(http.StatusSeeOther, "/logout")
				return
			}

			user2 := user.(*models.User)

			c.HTML(http.StatusOK, "welcome", gin.H{
				"User": user2,
			})
		})
	}

	s.http.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("mfa")
		session.Delete("user")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
	})

	s.http.POST("/login", ac.PostLogin)
	s.http.POST("/register", ac.PostRegister)
	s.http.POST("/forget", ac.PostForgetPassword)

	// 2FA
	fa2 := s.http.Group("/")
	fa2.Use(s.middlewareMFAAuthorization(as))
	{
		fa2.GET("/2fa", func(c *gin.Context) {
			user, ok := c.Get("user")
			if !ok {
				c.Redirect(http.StatusSeeOther, "/logout")
				return
			}

			user2 := user.(*models.User)

			c.HTML(http.StatusOK, "mfa", gin.H{
				"csrf":       csrf.GetToken(c),
				"first_name": user2.FirstName,
			})
		})
		fa2.POST("/2fa", ac.PostMFA)
	}

	// Verification by Email Group
	s.http.GET("/verify/registration/:code", ac.PostVerifyRegistration)
	s.http.GET("/verify/resetpassword/:code", ac.PostVerifyForgetPassword)

	s.http.GET("/newpassword/:code", ac.GetNewPasswordForm)
	s.http.POST("/newpassword/:code", ac.PostNewPassword)

	s.http.GET("/login", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		errorFlash := session.Flashes("error")
		messageFlash := session.Flashes("message")
		session.Save()

		ctx.HTML(http.StatusOK, "index", gin.H{
			"csrf":    csrf.GetToken(ctx),
			"error":   errorFlash,
			"message": messageFlash,
		})

	})

	s.http.GET("/register", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		errorFlash := session.Flashes("error")
		messageFlash := session.Flashes("message")
		session.Save()

		ctx.HTML(http.StatusOK, "register", gin.H{
			"csrf":          csrf.GetToken(ctx),
			"error":         errorFlash,
			"message":       messageFlash,
			"recaptchaSite": s.Config().RecaptchaSiteKey,
		})
	})

	s.http.GET("/forget", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		errorFlash := session.Flashes("error")
		messageFlash := session.Flashes("message")
		session.Save()
		ctx.HTML(http.StatusOK, "forget", gin.H{
			"csrf":    csrf.GetToken(ctx),
			"error":   errorFlash,
			"message": messageFlash,
		})
	})
}

func (s *server) middlewareFullAuthorization(as services.AuthenticationService) func(c *gin.Context) {
	return func(c *gin.Context) {

		// Check if user's user and mfa flag is set
		session := sessions.Default(c)
		userJwt := session.Get("user")
		if userJwt == nil {
			session.Delete("mfa")
			session.Delete("user")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		mfa, ok := session.Get("mfa").(bool)
		if !ok {
			session.Delete("mfa")
			session.Delete("user")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		if !mfa {
			c.Redirect(http.StatusSeeOther, "/2fa")
			return
		}

		ac, err := as.ValidateJWTKey(userJwt.(string))
		if errlog.Debug(err) {
			session.Delete("mfa")
			session.Delete("user")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		user, err := as.RetrieveUserFromClaims(ac)
		c.Set("user", user)

		//All is OK, can pass
		c.Next()
	}
}

func (s *server) middlewareMFAAuthorization(as services.AuthenticationService) func(c *gin.Context) {
	return func(c *gin.Context) {

		// Check if user's user and mfa flag is set
		session := sessions.Default(c)
		userJwt := session.Get("user")
		if userJwt == nil {
			session.Delete("mfa")
			session.Delete("user")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		mfa, ok := session.Get("mfa").(bool)
		if !ok {
			session.Delete("mfa")
			session.Delete("user")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		if mfa {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}

		ac, err := as.ValidateJWTKey(userJwt.(string))
		if errlog.Debug(err) {
			session.Delete("mfa")
			session.Delete("user")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		user, err := as.RetrieveUserFromClaims(ac)
		c.Set("user", user)

		//All is OK, can pass
		c.Next()
	}
}
