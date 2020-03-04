package config

type ServerConfiguration struct {
	Port int `required:"true" default:"443"`

	Domain string `required:"true"`

	CookieSecret string `required:"true" split_words:"true"`
	CookieName   string `required:"true" split_words:"true"`
	CSRFSecret   string `required:"true" split_words:"true"`
	JwtSecret    string `required:"true" split_words:"true"`

	// Email
	SMTPHost     string `required:"true" split_words:"true"`
	SMTPPort     int    `required:"true" split_words:"true"`
	SMTPUsername string `required:"true" split_words:"true"`
	SMTPPassword string `required:"true" split_words:"true"`
	SMTPFrom     string `required:"true" split_words:"true"`

	// Authentication Policy
	MaximumFailedAttempts int `required:"true" split_words:"true" default:"0"`
	ValidJWTLengthHours   int `split_words:"true" default:"6"`

	// Password Policy
	PasswordDatabaseURL string `default:"./common-passwords.db"`

	// ReCaptcha Settings
	RecaptchaSiteKey   string `required:"true" split_words:"true"`
	RecaptchaSecretKey string `required:"true" split_words:"true"`
}
