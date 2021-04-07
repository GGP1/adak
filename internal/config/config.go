package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
)

var (
	configFileName    = "config"
	configFileExt     = ".yml"
	configType        = "yaml"
	appName           = "adak"
	configDir         = filepath.Join(getConfigDir(runtime.GOOS), appName)
	configFileAbsPath = filepath.Join(configDir, configFileName)
)

// Config constains all the server configurations.
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Email    EmailConfig
	Stripe   StripeConfig
}

// DatabaseConfig hols the database attributes.
type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

// CacheConfig is the LRU-cache configuration.
type CacheConfig struct {
	Size int
}

// ServerConfig holds the server attributes.
type ServerConfig struct {
	Host    string
	Port    string
	Timeout struct {
		Read     time.Duration
		Write    time.Duration
		Shutdown time.Duration
	}
}

// EmailConfig holds email attributes.
type EmailConfig struct {
	Host        string
	Port        string
	Sender      string
	Password    string
	AdminEmails string
}

// StripeConfig hold stripe attributes
type StripeConfig struct {
	SecretKey string
	Logger    struct {
		Level stripe.Level
	}
}

// New sets up the configuration with the values the user gave.
func New() (*Config, error) {
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configType)

	path := os.Getenv("ADAK_CONFIG")
	if path != "" {
		if filepath.Ext(path) == "" {
			path = filepath.Join(path, ".env")
		}

		if err := godotenv.Load(path); err != nil {
			return nil, errors.Wrap(err, "env loading failed")
		}
	}

	// Bind envs
	for k, v := range envVars {
		viper.BindEnv(k, v)
	}

	// Set defaults
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	// Read or create configuration file
	if err := loadConfig(); err != nil {
		return nil, errors.Wrap(err, "read configuration failed")
	}

	// Auto read env variables
	viper.AutomaticEnv()
	config := &Config{}

	// Unmarshal config file to struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, errors.Wrap(err, "unmarshal configuration failed")
	}

	return config, nil
}

// read configuration from file.
func loadConfig() error {
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		// if file does not exist, simply create one
		if _, err := os.Stat(configFileAbsPath + configFileExt); os.IsNotExist(err) {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return errors.New("failed creating folder")
			}
			f, err := os.Create(configFileAbsPath + configFileExt)
			if err != nil {
				return errors.New("failed creating file")
			}
			f.Close()
		} else {
			return err
		}

		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}

	return nil
}

func getConfigDir(osName string) string {
	if os.Getenv("SV_DIR") != "" {
		return os.Getenv("SV_DIR")
	}

	switch osName {
	case "windows":
		return os.Getenv("APPDATA")
	case "darwin":
		return os.Getenv("HOME") + "/Library/Application Support"
	case "linux":
		return os.Getenv("HOME") + "/.config"
	default:
		dir, _ := os.Getwd()
		return dir
	}
}

var (
	defaults = map[string]interface{}{
		// Admins
		"admin.emails": []string{},
		// Database
		"database.username": "postgres",
		"database.password": "password",
		"database.host":     "localhost",
		"database.port":     "5432",
		"database.name":     "postgres",
		"database.sslmode":  "disable",
		// Cache
		"cache.size": 100,
		// Server
		"server.host":             "localhost",
		"server.port":             "7070",
		"server.dir":              "../",
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
		// Database
		"database.username": "POSTGRES_USERNAME",
		"database.password": "POSTGRES_PASSWORD",
		"database.host":     "POSTGRES_HOST",
		"database.port":     "POSTGRES_PORT",
		"database.name":     "POSTGRES_NAME",
		"database.sslmode":  "POSTGRES_SSL",
		// Cache
		"cache.size": "CACHE_SIZE",
		// Server
		"server.host":             "SV_HOST",
		"server.port":             "SV_PORT",
		"server.dir":              "SV_DIR",
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
