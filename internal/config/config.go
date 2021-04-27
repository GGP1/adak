package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
)

// Config constains all the server configurations.
type Config struct {
	Postgres Postgres
	Server   Server
	Email    Email
	Stripe   Stripe
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

// Cache is the LRU-cache configuration.
type Cache struct {
	Size int
}

// Server holds the server attributes.
type Server struct {
	Host    string
	Port    string
	Timeout struct {
		Read     time.Duration
		Write    time.Duration
		Shutdown time.Duration
	}
}

// Email holds email attributes.
type Email struct {
	Host        string
	Port        string
	Sender      string
	Password    string
	AdminEmails string
}

// Stripe hold stripe attributes
type Stripe struct {
	SecretKey string
	Logger    struct {
		Level stripe.Level
	}
}

// New sets up the configuration with the values the user gave.
func New() (*Config, error) {
	path := getConfigDir()
	viper.SetConfigFile(path)

	// Bind envs
	for k, v := range envVars {
		viper.BindEnv(k, v)
	}

	// Read or create configuration file
	if err := loadConfig(path); err != nil {
		return nil, errors.Wrap(err, "read configuration failed")
	}

	config := &Config{}

	// Unmarshal config file to struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, errors.Wrap(err, "unmarshal configuration failed")
	}

	return config, nil
}

// read configuration from file.
func loadConfig(path string) error {
	// if file does not exist, simply create one
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Info().Msgf("Configuration file not found, creating it on: %s", path)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return errors.New("failed creating folder")
		}
		f, err := os.Create(path)
		if err != nil {
			return errors.New("failed creating file")
		}
		f.Close()
		// Set defaults
		for k, v := range defaults {
			viper.SetDefault(k, v)
		}
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}

	return viper.ReadInConfig()
}

func getConfigDir() string {
	if path := os.Getenv("ADAK_CONFIG"); path != "" {
		log.Info().Msgf("Using customized configuration: %s", path)
		if ext := filepath.Ext(path); ext != "" && ext != "." {
			viper.SetConfigType(ext[1:])
		}
		return path
	}

	log.Info().Msg("Using default configuration")
	viper.SetConfigType("yaml")
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "config.yml")
}

var (
	defaults = map[string]interface{}{
		// Admins
		"admin.emails": []string{},
		// Postgres
		"postgres.host":     "postgres",
		"postgres.port":     "5432",
		"postgres.name":     "postgres",
		"postgres.username": "postgres",
		"postgres.password": "postgres",
		"postgres.sslmode":  "disable",
		// Cache
		"cache.size": 100,
		// Server
		"server.host":             "0.0.0.0",
		"server.port":             "4000",
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
		// JWT
		"jwt.secretkey": "secretkey",
	}

	envVars = map[string]string{
		// Admins
		"admin.emails": "ADMIN_EMAILS",
		// Postgres
		"postgres.host":     "POSTGRES_HOST",
		"postgres.port":     "POSTGRES_PORT",
		"postgres.name":     "POSTGRES_DB",
		"postgres.username": "POSTGRES_USER",
		"postgres.password": "POSTGRES_PASSWORD",
		"postgres.sslmode":  "POSTGRES_SSL",
		// Cache
		"cache.size": "CACHE_SIZE",
		// Server
		"server.host":             "SV_HOST",
		"server.port":             "SV_PORT",
		"server.timeout.read":     "SV_TIMEOUT_READ",
		"server.timeout.write":    "SV_TIMEOUT_WRITE",
		"server.timeout.shutdown": "SV_TIMEOUT_SHUTDOWN",
		// Email
		"email.host":     "EMAIL_HOST",
		"email.port":     "EMAIL_PORT",
		"email.sender":   "EMAIL_SENDER",
		"email.password": "EMAIL_PASSWORD",
		"email.admins":   "EMAIL_ADMINS_PATH",
		// Stripe
		"stripe.secretkey":    "STRIPE_SECRET_KEY",
		"stripe.logger.level": "STRIPE_LOGGER_LEVEL",
		// JWT
		"jwt.secretkey": "TOKEN_SECRET_KEY",
	}
)
