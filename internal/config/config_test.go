package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		desc         string
		configEnvVar string
		expected     string
	}{
		{
			desc:         "Defaults",
			configEnvVar: "",
			expected:     "default@adak.com",
		},
		{
			desc:         "Mock config",
			configEnvVar: "testdata/.mock_env",
			expected:     "test@testing.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			os.Setenv("ADAK_CONFIG", tc.configEnvVar)

			config, err := New()
			assert.NoError(t, err, "New()")

			assert.Equal(t, tc.expected, config.Email.Sender)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	err := loadConfig()
	assert.NoError(t, err, "loadConfig()")
}

func TestLoadConfigCreate(t *testing.T) {
	old, _ := os.Getwd()

	os.Chdir("testdata")
	dir, _ := os.Getwd()
	configDir = dir
	configFilename = "mock_config"

	err := loadConfig()
	assert.NoError(t, err, "loadConfig()")

	// Go back to the last directory
	os.Chdir(old)
}

func TestGetConfigDir(t *testing.T) {
	linux := filepath.Join(os.Getenv("HOME"), ".config")
	darwin := filepath.Join(os.Getenv("HOME"), "Library/Application Support")
	wd, _ := os.Getwd()
	appData := os.Getenv("APPDATA")

	testCases := []struct {
		desc     string
		osName   string
		expected string
		svDirEnv string
	}{
		{
			desc:     "Linux",
			osName:   "linux",
			expected: linux,
		},
		{
			desc:     "Darwin",
			osName:   "darwin",
			expected: darwin,
		},
		{
			desc:     "Windows",
			osName:   "windows",
			expected: appData,
		},
		{
			desc:     "Default",
			osName:   "",
			expected: wd,
		},
		{
			desc:     "SV_DIR environment variable",
			svDirEnv: "/custom/path",
			expected: "/custom/path",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			os.Setenv("SV_DIR", "")
			if tc.svDirEnv != "" {
				os.Setenv("SV_DIR", tc.svDirEnv)
			}

			got := getConfigDir(tc.osName)
			assert.Equal(t, tc.expected, got)
		})
	}
}
