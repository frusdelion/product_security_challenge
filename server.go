package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/danielkov/gin-helmet"
	"github.com/frusdelion/zendesk-product_security_challenge/config"
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	server2 "github.com/frusdelion/zendesk-product_security_challenge/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	jwemail "github.com/jordan-wright/email"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/snwfdhmp/errlog"
	"github.com/utrack/gin-csrf"
	"golang.org/x/time/rate"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dvwright/xss-mw"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/utrack/gin-merry"
	limit "github.com/yangxikun/gin-limit-by-key"
)

func NewServer(log *logrus.Logger) server2.Server {
	var sc config.ServerConfiguration

	err := envconfig.Process("app", &sc)
	if errlog.Debug(err) {
		log.Fatal(err.Error())
	}

	r := gin.Default()
	store := cookie.NewStore([]byte(sc.CookieSecret))
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   sc.Domain,
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	r.Use(sessions.Sessions(sc.CookieName, store))
	r.Use(csrf.Middleware(csrf.Options{
		Secret:        sc.CSRFSecret,
		IgnoreMethods: []string{"GET"},
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
		TokenGetter: func(c *gin.Context) string {
			return c.PostForm("__csrf")
		},
	}))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("https://%s", sc.Domain)},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	}))

	merryMl := &ginMerry.Middleware{
		Debug:        false,
		GenericError: "We have encountered an error.",
		LogFunc: func(err string, code int, vals map[string]interface{}) {
			log.Errorf("[%d] %s (%v)", code, err, vals)
		},
	}
	r.Use(merryMl.Handler())

	xssMdlwr := &xss.XssMw{
		FieldsToSkip: []string{"password"},
	}
	r.Use(xssMdlwr.RemoveXss())

	r.Use(limit.NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP() // limit rate by client ip
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(rate.Every(100*time.Millisecond), 10), time.Hour // limit 10 qps/clientIp and permit bursts of at most 10 tokens, and the limiter liveness time duration is 1 hour
	}, func(c *gin.Context) {
		c.AbortWithStatus(429) // handle exceed rate limit request
	}))

	r.Use(helmet.Default())

	ctx, cancelFunc := context.WithCancel(context.Background())

	emailHost := fmt.Sprintf("%s:%d", sc.SMTPHost, sc.SMTPPort)
	emailPool, err := jwemail.NewPool(emailHost, 4, smtp.PlainAuth("", sc.SMTPUsername, sc.SMTPPassword, sc.SMTPHost))
	if errlog.Debug(err) {
		log.Panic(err)
	}

	db, err := gorm.Open("sqlite3", "./db.sqlite3")
	db.LogMode(true)
	if errlog.Debug(err) {
		log.Panic(err)
	}

	return &server{
		log:        log,
		http:       r,
		config:     sc,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		emailPool:  emailPool,
		db:         db,
	}
}

type server struct {
	db         *gorm.DB
	log        *logrus.Logger
	http       *gin.Engine
	config     config.ServerConfiguration
	ctx        context.Context
	cancelFunc context.CancelFunc
	emailPool  *jwemail.Pool
	validator  *validator.Validate
}

func (s server) Validator() *validator.Validate {
	return s.validator
}

func (s server) GreeterBanner() {
	s.Log().Infoln("                                         S.;..                   ;;;;           ")
	s.Log().Infoln("  .  . .  .  . .  .  . .  .  . .  .  . . 88S;   .  . .  .  . .  .X88S . .  . .  ")
	s.Log().Infoln("   S@@@@@@@;. .%888t     .X88@;    .;@@8t88%; . %888t    .%8@8t. @@88  :X8t    .")
	s.Log().Infoln(" ..X@%%SXX8 .%SX@;SX.% .;@X@t8@X;.%888X%888%; X%S8;SXt8 t8SS8X88:XX@X.8%8X; .   ")
	s.Log().Infoln("   ....;88t:8X8S. :X8. XS88 .SX@   8 . .X 8@;8S8;. :%@%:888S..@: X888.8%8.   .  ")
	s.Log().Infoln("  .  8.8t8  ;88Sttt@X8S888t  .X@;:8X8  . 88S: @88tttSX88  S88X %:@88X8X:. .    .")
	s.Log().Infoln("  .8 8St;tt:88X.@%XS%t.888S   @8:. 8 %:;888S;X8t%8SSX;X %;:88t8X8%88;8;8:;  .   ")
	s.Log().Infoln("  : 8X8X@88 .@;S8S888% 88.: . .tt  %888@:t8  .S S88X88X t:88SSt:.XX%%  XXt8  .  ")
	s.Log().Infoln("  .  . .  .  . .  .  . .  .  . .  .  . . . .    .  . .  .  . .  . . . . .  . .  ")
	s.Log().Infoln()
	s.Log().Infoln("Zendesk Product Security Challenge, March 2020")
	s.Log().Infoln("Done by Lionell Yip")
}

func (s server) AutoMigrate() {
	s.DB().AutoMigrate(
		&models.User{},
		&models.FailedLogin{},
		&models.Verification{},
	)
}

func (s server) Log() *logrus.Logger {
	return s.log
}

func (s server) DB() *gorm.DB {
	return s.db
}

func (s server) Email() *jwemail.Pool {
	return s.emailPool
}

func (s server) Context() context.Context {
	return s.ctx
}

func (s server) ContextCancel() context.CancelFunc {
	return s.cancelFunc
}

func (s server) Config() config.ServerConfiguration {
	return s.config
}

func (s server) Run() {
	s.GreeterBanner()
	s.AutoMigrate()
	s.Routes()

	// Setup validator
	s.validator = validator.New()
	s.validator.SetTagName("binding")

	serveMux := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Config().Port),
		Handler: s.http,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		},
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	go func() { // http proxy
		s.Log().Infof("Visit us at %s", fmt.Sprintf("https://%s:%d", s.Config().Domain, s.Config().Port))
		if err := serveMux.ListenAndServeTLS("localhost+1.pem", "localhost+1-key.pem"); err != http.ErrServerClosed {
			s.Log().Fatalf("failed to Gin serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.emailPool.Close()

	s.Log().Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	if err := serveMux.Shutdown(ctx); err != nil {
		s.Log().Fatal("Gin shutdown: ", err)
	}

}
