package main

import (
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/frusdelion/zendesk-product_security_challenge/controllers"
	"github.com/frusdelion/zendesk-product_security_challenge/repositories"
	"github.com/frusdelion/zendesk-product_security_challenge/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

func (s *server) Routes() {
	//new template engine
	s.http.HTMLRender = ginview.New(goview.Config{
		Root:      "./project",
		Extension: ".html",
		Master:    "./layouts/master",
	})

	// TODO: go.rice
	s.http.Static("/assets", "./project/assets")

	ac := controllers.NewAuthenticationController(s, services.NewAuthenticationService(repositories.NewAuthenticationRepository(s.DB()), repositories.NewUserRepository(s.DB()), s))

	r := s.http.Group("/")
	r.Use(func(c *gin.Context) {
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

		//All is OK, can pass
		c.Next()
	})
	{
		r.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "hello")
		})
	}

	s.http.POST("/login", ac.PostLogin)
	s.http.POST("/register", ac.PostRegister)
	s.http.POST("/forget", ac.PostForgetPassword)

	// Verification by Email Group
	s.http.GET("/verify/registration/:code")
	s.http.GET("/verify/resetpassword/:code")

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
