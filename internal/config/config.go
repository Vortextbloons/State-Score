package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	AppName     = "StateScore"
	DefaultHost = "127.0.0.1"
	DefaultPort = 8787
	Version     = "0.1.0"
)

// Config holds runtime configuration for the local application.
type Config struct {
	Host         string
	Port         int
	DataDir      string
	DatabasePath string
	OpenBrowser  bool
}

// Load returns the default configuration and ensures the data directory exists.
// STATESCORE_PORT overrides the listen port (used for dual-server frontend development).
func Load() (*Config, error) {
	dataDir, err := DataDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}

	port := DefaultPort
	if raw := os.Getenv("STATESCORE_PORT"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 65535 {
			return nil, fmt.Errorf("invalid STATESCORE_PORT %q", raw)
		}
		port = parsed
	}

	openBrowser := true
	if os.Getenv("STATESCORE_NO_BROWSER") == "1" {
		openBrowser = false
	}

	return &Config{
		Host:         DefaultHost,
		Port:         port,
		DataDir:      dataDir,
		DatabasePath: filepath.Join(dataDir, "statescore.db"),
		OpenBrowser:  openBrowser,
	}, nil
}

// DataDir returns the OS-specific application data directory.
func DataDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		base := os.Getenv("LOCALAPPDATA")
		if base == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			base = filepath.Join(home, "AppData", "Local")
		}
		return filepath.Join(base, AppName), nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", AppName), nil
	default:
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			return filepath.Join(xdg, "statescore"), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".local", "share", "statescore"), nil
	}
}
