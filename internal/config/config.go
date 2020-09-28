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
	configuration          *Configuration
	configFileName         = "config"
	configFileExt          = ".yml"
	configType             = "yaml"
	appName                = "palo"
	configurationDirectory = filepath.Join(osConfigDirectory(runtime.GOOS), appName)
	configFileAbsPath      = filepath.Join(configurationDirectory, configFileName)
)

// Configuration constains all the server configurations.
type Configuration struct {
	Database DatabaseConfiguration
	Server   ServerConfiguration
	Email    EmailConfiguration
	Stripe   StripeConfiguration
}

// DatabaseConfiguration hols the database attributes.
type DatabaseConfiguration struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

// ServerConfiguration holds the server attributes.
type ServerConfiguration struct {
	Host    string
	Port    string
	Timeout struct {
		Read     time.Duration
		Write    time.Duration
		Shutdown time.Duration
	}
}

// EmailConfiguration holds email attributes.
type EmailConfiguration struct {
	Host        string
	Port        string
	Sender      string
	Password    string
	AdminEmails string
}

// StripeConfiguration hold stripe attributes
type StripeConfiguration struct {
	SecretKey string
	Logger    struct {
		Level stripe.Level
	}
}

// New sets up the configuration with the values the user gave.
// Defaults and env variables are placed at the end to make the config easier to read.
func New() (*Configuration, error) {
	viper.AddConfigPath(configurationDirectory)
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configType)

	path := os.Getenv("PALO_CONFIG") + "/.env"

	if err := godotenv.Load(path); err != nil {
		return nil, errors.Wrap(err, "env loading failed")
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
	if err := readConfiguration(); err != nil {
		return nil, errors.Wrap(err, "read configuration failed")
	}

	// Auto read env variables
	viper.AutomaticEnv()

	// Unmarshal config file to struct
	if err := viper.Unmarshal(&configuration); err != nil {
		return nil, errors.Wrap(err, "unmarshal configuration failed")
	}

	return configuration, nil
}

// read configuration from file
func readConfiguration() error {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// if file does not exist, simply create one
		if _, err := os.Stat(configFileAbsPath + configFileExt); os.IsNotExist(err) {
			os.MkdirAll(configurationDirectory, 0755)
			os.Create(configFileAbsPath + configFileExt)
		} else {
			return err
		}

		// Write defaults
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}

	return nil
}

func osConfigDirectory(osName string) string {
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
		// Database
		"database.username": "postgres",
		"database.password": "password",
		"database.host":     "localhost",
		"database.port":     "5432",
		"database.name":     "postgres",
		"database.sslmode":  "disable",
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
		"email.sender":   "default@palo.com",
		"email.password": "default",
		"email.admins":   "../pkg/auth/",
		// Stripe
		"stripe.secretkey":    "sk_test_default",
		"stripe.logger.level": "4",
		// JWT
		"jwt.secretkey": "secretkey",
	}

	envVars = map[string]string{
		// Database
		"database.username": "DB_USERNAME",
		"database.password": "DB_PASSWORD",
		"database.host":     "DB_HOST",
		"database.port":     "DB_PORT",
		"database.name":     "DB_NAME",
		"database.sslmode":  "DB_SSL",
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
