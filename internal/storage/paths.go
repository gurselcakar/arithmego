package storage

import (
	"os"
	"path/filepath"
)

const (
	configDirName  = "arithmego"
	statisticsFile = "statistics.json"
	settingsFile   = "settings.json"
)

// configDirOverride allows tests to use a temporary directory.
var configDirOverride string

// SetConfigDirForTesting sets a custom config directory for tests.
// Pass an empty string to restore default behavior.
func SetConfigDirForTesting(path string) {
	configDirOverride = path
}

// ConfigDir returns the path to the ArithmeGo config directory.
// Creates the directory if it doesn't exist.
func ConfigDir() (string, error) {
	var dir string
	if configDirOverride != "" {
		dir = configDirOverride
	} else {
		configHome, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(configHome, configDirName)
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return dir, nil
}

// StatisticsPath returns the path to the statistics file.
func StatisticsPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, statisticsFile), nil
}

// SettingsPath returns the path to the settings file.
func SettingsPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, settingsFile), nil
}
