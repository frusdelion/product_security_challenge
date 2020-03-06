package controllers

import (
	"fmt"
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	server2 "github.com/frusdelion/zendesk-product_security_challenge/server"
	"github.com/frusdelion/zendesk-product_security_challenge/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/snwfdhmp/errlog"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

type AuthenticationController interface {
	PostRegister(ctx *gin.Context)
	PostLogin(ctx *gin.Context)
	PostMFA(ctx *gin.Context)
	PostForgetPassword(ctx *gin.Context)

	PostVerifyRegistration(c *gin.Context)
	PostVerifyForgetPassword(c *gin.Context)
	GetNewPasswordForm(c *gin.Context)
	PostNewPassword(c *gin.Context)
}

func NewAuthenticationController(s server2.Server, as services.AuthenticationService, us services.UserService, vs services.VerificationService) AuthenticationController {
	return &authenticationController{s: s, as: as, vs: vs, us: us}
}

type authenticationController struct {
	s  server2.Server
	us services.UserService
	as services.AuthenticationService
	vs services.VerificationService
}

func (a authenticationController) GetNewPasswordForm(c *gin.Context) {
	code := c.Param("code")

	_, err := a.vs.VerifyForgotPassword(code)
	if errlog.Debug(err) {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	session := sessions.Default(c)
	errorFlash := session.Flashes("error")
	messageFlash := session.Flashes("message")
	session.Save()

	c.HTML(http.StatusOK, "newpassword", gin.H{
		"csrf":    csrf.GetToken(c),
		"error":   errorFlash,
		"message": messageFlash,
	})
}

func (a authenticationController) PostNewPassword(c *gin.Context) {
	code := c.Param("code")

	vr, err := a.vs.VerifyForgotPassword(code)
	if errlog.Debug(err) {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	anp := &models.AuthenticationNewPassword{}
	if err := c.ShouldBind(anp); errlog.Debug(err) {
		c.Error(err)
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/newpassword/%s", code))
		return
	}

	usr, err := a.us.UpdateUser(&vr.User, models.User{Password: anp.Password})
	if errlog.Debug(err) {
		c.Error(err)
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/newpassword/%s", code))
		return
	}

	a.s.Log().Infof("[SECURITY EVENT] updated password for %d, ip:%s, browserUA:%s, browserFingerprint:%s", usr.ID, c.ClientIP(), c.GetHeader("User-Agent"), anp.BrowserFingerprint)

	if err := a.vs.DeleteVerification(vr); errlog.Debug(err) {
		c.Error(err)
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/newpassword/%s", code))
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Your password has been changed successfully.", "message")
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}

func (a authenticationController) PostVerifyRegistration(c *gin.Context) {
	code := c.Param("code")

	vr, err := a.vs.VerifyRegistration(code)
	if errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if err := a.vs.DeleteVerification(vr); errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if err := a.us.ActivateEmail(&vr.User); errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	session := sessions.Default(c)
	session.AddFlash("Thanks for verifying!", "message")
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}

func (a authenticationController) PostVerifyForgetPassword(c *gin.Context) {
	code := c.Param("code")

	_, err := a.vs.VerifyForgotPassword(code)
	if errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/newpassword/%s", code))
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

	ac, user, err := a.as.AuthenticateUser(au)
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

	if err := a.vs.SendMFAVerification(user, c.GetHeader("User-Agent"), c.ClientIP()); errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Error(err)
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	session := sessions.Default(c)
	if au.Remember == "" {
		session.Options(sessions.Options{MaxAge: 3 * 60})
	}

	session.Set("user", jwtKey)
	session.Set("mfa", false)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/2fa")
}

func (a authenticationController) PostMFA(c *gin.Context) {
	code := c.PostForm("code")
	u1, ok := c.Get("user")
	if !ok {
		c.Redirect(http.StatusSeeOther, "/logout")
		return
	}

	user := u1.(*models.User)

	vr, err := a.vs.VerifyMFA(code, user)
	if errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Error(err)
		c.Redirect(http.StatusSeeOther, "/2fa")
		return
	}

	if err := a.vs.DeleteVerification(vr); errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Error(err)
		c.Redirect(http.StatusSeeOther, "/2fa")
		return
	}

	session := sessions.Default(c)
	session.Set("mfa", true)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/")
}

func (a authenticationController) PostForgetPassword(c *gin.Context) {
	af := &models.AuthenticationForgot{}

	if err := c.ShouldBind(af); errlog.Debug(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/forget")
		return
	}

	user, err := a.us.FindUserByEmail(af.Email)
	if errlog.Debug(err) && !gorm.IsRecordNotFoundError(err) {
		session := sessions.Default(c)
		session.AddFlash(err.Error(), "error")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/forget")
		return
	}

	foundUser := true
	if gorm.IsRecordNotFoundError(err) {
		foundUser = false
	}

	if foundUser {
		if err := a.vs.SendForgotPasswordVerification(user, c.GetHeader("User-Agent"), c.ClientIP()); errlog.Debug(err) {
			session := sessions.Default(c)
			session.AddFlash(err.Error(), "error")
			session.Save()
			c.Redirect(http.StatusSeeOther, "/forget")
			return
		}
	}

	session := sessions.Default(c)
	session.AddFlash("If this email has been registered, you will receive an email.", "message")
	session.Save()
	c.Redirect(http.StatusSeeOther, "/forget")
}
