package config

import (
	"embed"
	"os"
	"path/filepath"
	"time"

	"github.com/GGP1/adak/internal/logger"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v72"
)

// Config constains all the server configurations.
type Config struct {
	Admins      []string
	Development bool

	Email       Email
	Memcached   Memcached
	Postgres    Postgres
	RateLimiter RateLimiter
	Redis       Redis
	Server      Server
	Session     Session
	Static      Static
	Stripe      Stripe
}

// Email holds email attributes.
type Email struct {
	Host     string
	Port     string
	Sender   string
	Password string
}

// Memcached is the LRU-cache configuration.
type Memcached struct {
	Servers []string
}

// Postgres hols the database attributes.
type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

// RateLimiter configuration.
type RateLimiter struct {
	Enabled bool
	Rate    int
}

// Redis configuration.
type Redis struct {
	Host     string
	Port     string
	Password string
}

// Server holds the server attributes.
type Server struct {
	Host string
	Port string
	TLS  struct {
		KeyFile  string
		CertFile string
	}
	Timeout struct {
		Read     time.Duration
		Write    time.Duration
		Shutdown time.Duration
	}
}

// Session contains the session configuration.
type Session struct {
	Attempts int64
	Delay    int64
	Length   int
}

// Static contains the static file system.
type Static struct {
	FS embed.FS
}

// Stripe hold stripe attributes
type Stripe struct {
	SecretKey string
	Logger    struct {
		Level stripe.Level
	}
}

// New sets up the configuration with the values the user gave.
// Defaults and env variables are placed at the end to make the config easier to read.
func New() (Config, error) {
	path := getConfigPath()
	viper.SetConfigFile(path)

	// Bind envs
	for k, v := range envVars {
		viper.BindEnv(k, v)
	}

	// Read or create configuration file
	if err := loadConfig(path); err != nil {
		return Config{}, errors.Wrap(err, "couldn't read the configuration file")
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return Config{}, errors.Wrap(err, "unmarshal configuration failed")
	}

	return *config, nil
}

// read configuration from file.
func loadConfig(path string) error {
	// if file does not exist, create one
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return errors.Wrap(err, "creating folder")
		}
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			return errors.Wrap(err, "creating file")
		}
		for k, v := range defaults {
			viper.SetDefault(k, v)
		}
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "reading configuration")
	}

	return nil
}

// getConfigPath returns the location of the configuration file.
func getConfigPath() string {
	if path := os.Getenv("ADAK_CONFIG"); path != "" {
		ext := filepath.Ext(path)
		if ext != "" && ext != "." {
			viper.SetConfigType(ext[1:])
		}

		logger.Infof("Using customized configuration: %s", path)
		return path
	}

	viper.SetConfigType("yaml")
	logger.Info("Using default configuration")
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "config.yml")
}

var (
	defaults = map[string]interface{}{
		// Admins
		"admins": []string{},
		// Development
		"development": true,
		// Email
		"email.host":     "smtp.default.com",
		"email.port":     "587",
		"email.sender":   "default@adak.com",
		"email.password": "default",
		"email.admins":   "../pkg/auth/",
		// Google
		"google.client.id":     "id",
		"google.client.secret": "secret",
		// Memcached
		"memcached.servers": []string{"localhost:11211"},
		// Postgres
		"postgres.username": "postgres",
		"postgres.password": "password",
		"postgres.host":     "localhost",
		"postgres.port":     "5432",
		"postgres.name":     "postgres",
		"postgres.sslmode":  "disable",
		// Rate limiter
		"ratelimiter.enabled": true,
		"ratelimiter.rate":    5, // Per minute
		// Redis
		"redis.host":     "localhost",
		"redis.port":     "6379",
		"redis.password": "redis",
		// Server
		"server.host":             "localhost",
		"server.port":             "7070",
		"server.tls.keyfile":      "",
		"server.tls.certfile":     "",
		"server.timeout.read":     5,
		"server.timeout.write":    5,
		"server.timeout.shutdown": 5,
		// Session
		"session.attempts": 5,
		"session.delay":    0,
		"session.length":   0,
		// Stripe
		"stripe.secretkey":    "sk_test_default",
		"stripe.logger.level": "4",
		// Token
		"token.secretkey": "secretkey",
	}

	envVars = map[string]string{
		// Admins
		"admins": "ADAK_ADMINS",
		// Development
		"development": "DEVELOPMENT",
		// Email
		"email.host":     "EMAIL_HOST",
		"email.port":     "EMAIL_PORT",
		"email.sender":   "EMAIL_SENDER",
		"email.password": "EMAIL_PASSWORD",
		// Google
		"google.client.id":     "GOOGLE_CLIENT_ID",
		"google.client.secret": "GOOGLE_CLIENT_SECRET",
		// Memcached
		"memcached.servers": "MEMCACHED_SERVERS",
		// Postgres
		"postgres.username": "POSTGRES_USERNAME",
		"postgres.password": "POSTGRES_PASSWORD",
		"postgres.host":     "POSTGRES_HOST",
		"postgres.port":     "POSTGRES_PORT",
		"postgres.name":     "POSTGRES_DB",
		"postgres.sslmode":  "POSTGRES_SSL",
		// Rate limiter
		"ratelimiter.enabled": "RATELIMITER_ENABLED",
		"ratelimiter.rate":    "RATELIMITER_RATE",
		// Redis
		"redis.host":     "REDIS_HOST",
		"redis.port":     "REDIS_PORT",
		"redis.password": "REDIS_PASSWORD",
		// Server
		"server.host":             "SV_HOST",
		"server.port":             "SV_PORT",
		"server.tls.keyfile":      "SV_TLS_KEYFILE",
		"server.tls.certfile":     "SV_TLS_CERTFILE",
		"server.timeout.read":     "SV_TIMEOUT_READ",
		"server.timeout.write":    "SV_TIMEOUT_WRITE",
		"server.timeout.shutdown": "SV_TIMEOUT_SHUTDOWN",
		// Session
		"session.attempts": "SESSION_ATTEMPTS",
		"session.delay":    "SESSION_DELAY",
		"session.length":   "SESSION_LENGTH",
		// Stripe
		"stripe.secretkey":    "STRIPE_SECRET_KEY",
		"stripe.logger.level": "STRIPE_LOGGER_LEVEL",
		// Token
		"token.secretkey": "TOKEN_SECRET_KEY",
	}
)
