package server

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config contains the server configuration
type Config struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Timeout struct {
			Read     time.Duration `yaml:"read"`
			Write    time.Duration `yaml:"write"`
			Shutdown time.Duration `yaml:"shutdown"`
		} `yaml:"timeout"`
	} `yaml:"server"`
}

// NewConfig returns a new configuration
func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "../config.yml", "path to config file")
	flag.Parse()

	err := validateConfigPath(configPath)
	if err != nil {
		return "", err
	}

	return configPath, nil
}

// validateConfigPath checks if the flag is provided correctly
func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", path)
	}
	return nil
}
