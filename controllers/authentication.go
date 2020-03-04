package controllers

import (
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	server2 "github.com/frusdelion/zendesk-product_security_challenge/server"
	"github.com/frusdelion/zendesk-product_security_challenge/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/snwfdhmp/errlog"
	"net/http"
)

type AuthenticationController interface {
	PostRegister(ctx *gin.Context)
	PostLogin(ctx *gin.Context)
	PostMFA(ctx *gin.Context)
	PostForgetPassword(ctx *gin.Context)
}

func NewAuthenticationController(s server2.Server, as services.AuthenticationService) AuthenticationController {
	return &authenticationController{s: s, as: as}
}

type authenticationController struct {
	s  server2.Server
	as services.AuthenticationService
}

func (a authenticationController) PostRegister(ctx *gin.Context) {
	authenticateRegister := &models.AuthenticationRegister{
		IPAddress:        ctx.ClientIP(),
		BrowserUserAgent: ctx.GetHeader("User-Agent"),
	}
	err := ctx.ShouldBind(authenticateRegister)
	if errlog.Debug(err) {
		a.s.Log().Error(err)
		session := sessions.Default(ctx)
		session.AddFlash(err.Error(), "error")
		if err := session.Save(); errlog.Debug(err) {
			a.s.Log().Error(err)
		}
		ctx.Redirect(http.StatusSeeOther, "/register")
		return
	}

	_, err = a.as.Register(authenticateRegister)
	if errlog.Debug(err) {
		a.s.Log().Error(err)
		session := sessions.Default(ctx)
		session.AddFlash(err.Error(), "error")
		if err := session.Save(); errlog.Debug(err) {
			a.s.Log().Error(err)
		}
		ctx.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/register")
		return
	}

	session := sessions.Default(ctx)
	session.AddFlash("Look for your validation email in your inbox!", "message")
	if err := session.Save(); errlog.Debug(err) {
		a.s.Log().Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/register")
}

func (a authenticationController) PostLogin(c *gin.Context) {
	au := &models.AuthenticateUser{
		BrowserUserAgent: c.GetHeader("User-Agent"),
		IPAddress:        c.ClientIP(),
	}

	if err := c.ShouldBind(au); errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Error(err)
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	ac, _, err := a.as.AuthenticateUser(au)
	if errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Error(err)
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	jwtKey, err := a.as.GenerateJWTKey(ac)
	if errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Error(err)
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	session := sessions.Default(c)
	session.Set("user", jwtKey)
	session.Set("mfa", false)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/2fa")
}

func (a authenticationController) PostMFA(ctx *gin.Context) {
	panic("implement me")
}

func (a authenticationController) PostForgetPassword(ctx *gin.Context) {
	panic("implement me")
}
