package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		desc         string
		configEnvVar string
		expected     string
	}{
		{
			desc:         "Default",
			configEnvVar: "",
			expected:     "default@adak.com",
		},
		{
			desc:         "Custom",
			configEnvVar: "testdata/mock_config.yml",
			expected:     "test@testing.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			viper.Reset()
			os.Setenv("ADAK_CONFIG", tc.configEnvVar)

			config, err := New()
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, config.Email.Sender)
		})
	}

	dir, _ := os.UserConfigDir()
	assert.NoError(t, os.Remove(filepath.Join(dir, "config.yml")))
}

func TestLoadConfig(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		path := "config.yml"
		viper.SetConfigFile(path)
		assert.NoError(t, loadConfig(path))
		assert.NoError(t, os.Remove(path))
	})

	t.Run("Read", func(t *testing.T) {
		path := "testdata/mock_config.yml"
		viper.SetConfigFile(path)
		assert.NoError(t, loadConfig(path))
	})
}

func TestGetConfigPath(t *testing.T) {
	cfgDir, _ := os.UserConfigDir()

	testCases := []struct {
		desc     string
		expected string
		path     string
	}{
		{
			desc:     "Default",
			expected: filepath.Join(cfgDir, "config.yml"),
		},
		{
			desc:     "Custom",
			path:     "/custom/path/config.yml",
			expected: "/custom/path/config.yml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			os.Setenv("ADAK_CONFIG", "")
			if tc.path != "" {
				os.Setenv("ADAK_CONFIG", tc.path)
			}

			got := getConfigPath()
			assert.Equal(t, tc.expected, got)
		})
	}
}
