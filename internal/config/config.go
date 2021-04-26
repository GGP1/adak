package config

import (
	"embed"
	"os"
	"path/filepath"
	"time"

	"github.com/GGP1/adak/internal/logger"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
)

// Config constains all the server configurations.
type Config struct {
	Development bool

	Admin     Admin
	Postgres  Postgres
	Memcached Memcached
	Server    Server
	Email     Email
	Stripe    Stripe
	Static    Static
}

// Admin contains the admins emails.
type Admin struct {
	Emails []string
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

// Memcached is the LRU-cache configuration.
type Memcached struct {
	Servers []string
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

// Email holds email attributes.
type Email struct {
	Host     string
	Port     string
	Sender   string
	Password string
}

// Stripe hold stripe attributes
type Stripe struct {
	SecretKey string
	Logger    struct {
		Level stripe.Level
	}
}

// Static contains the static file system.
type Static struct {
	FS embed.FS
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

		logger.Log.Infof("Using customized configuration: %s", path)
		return path
	}

	viper.SetConfigType("yaml")
	logger.Log.Info("Using default configuration")
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "config.yml")
}

var (
	defaults = map[string]interface{}{
		// Admins
		"admin.emails": []string{},
		// Postgres
		"postgres.username": "postgres",
		"postgres.password": "password",
		"postgres.host":     "localhost",
		"postgres.port":     "5432",
		"postgres.name":     "postgres",
		"postgres.sslmode":  "disable",
		// Development
		"development": true,
		// Memcached
		"memcached.servers": []string{"localhost:11211"},
		// Server
		"server.host":             "localhost",
		"server.port":             "7070",
		"server.tls.keyfile":      "",
		"server.tls.certfile":     "",
		"server.timeout.read":     5,
		"server.timeout.write":    5,
		"server.timeout.shutdown": 5,
		// Email
		"email.host":     "smtp.default.com",
		"email.port":     "587",
		"email.sender":   "default@adak.com",
		"email.password": "default",
		"email.admins":   "../pkg/auth/",
		// Stripe
		"stripe.secretkey":    "sk_test_default",
		"stripe.logger.level": "4",
		// Token
		"token.secretkey": "secretkey",
	}

	envVars = map[string]string{
		// Admins
		"admin.emails": "ADMIN_EMAILS",
		// Postgres
		"postgres.username": "POSTGRES_USERNAME",
		"postgres.password": "POSTGRES_PASSWORD",
		"postgres.host":     "POSTGRES_HOST",
		"postgres.port":     "POSTGRES_PORT",
		"postgres.name":     "POSTGRES_DB",
		"postgres.sslmode":  "POSTGRES_SSL",
		// Development
		"development": "DEVELOPMENT",
		// Memcached
		"memcached.servers": "MEMCACHED_SERVERS",
		// Server
		"server.host":             "SV_HOST",
		"server.port":             "SV_PORT",
		"server.tls.keyfile":      "SV_TLS_KEYFILE",
		"server.tls.certfile":     "SV_TLS_CERTFILE",
		"server.timeout.read":     "SV_TIMEOUT_READ",
		"server.timeout.write":    "SV_TIMEOUT_WRITE",
		"server.timeout.shutdown": "SV_TIMEOUT_SHUTDOWN",
		// Email
		"email.host":     "EMAIL_HOST",
		"email.port":     "EMAIL_PORT",
		"email.sender":   "EMAIL_SENDER",
		"email.password": "EMAIL_PASSWORD",
		// Stripe
		"stripe.secretkey":    "STRIPE_SECRET_KEY",
		"stripe.logger.level": "STRIPE_LOGGER_LEVEL",
		// Token
		"token.secretkey": "TOKEN_SECRET_KEY",
		// Google
		"google.client.id":     "GOOGLE_CLIENT_ID",
		"google.client.secret": "GOOGLE_CLIENT_SECRET",
		// Github
		"github.client.id":     "GITHUB_CLIENT_ID",
		"github.client.secret": "GITHUB_CLIENT_SECRET",
	}
)
