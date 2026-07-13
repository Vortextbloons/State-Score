package config_test

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/isaac/statescore/internal/config"
)

func TestDataDir(t *testing.T) {
	t.Parallel()

	dir, err := config.DataDir()
	if err != nil {
		t.Fatalf("DataDir: %v", err)
	}
	if dir == "" {
		t.Fatal("DataDir returned empty path")
	}

	base := filepath.Base(dir)
	switch runtime.GOOS {
	case "windows", "darwin":
		if base != config.AppName {
			t.Fatalf("base = %q, want %q", base, config.AppName)
		}
	default:
		if base != "statescore" {
			t.Fatalf("base = %q, want statescore", base)
		}
	}

	if runtime.GOOS == "windows" && !strings.Contains(strings.ToLower(dir), "statescore") {
		t.Fatalf("unexpected windows data dir: %s", dir)
	}
}

func TestLoadCreatesDataDir(t *testing.T) {
	// Load uses the real OS data directory; only verify it succeeds and paths are set.
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DatabasePath == "" {
		t.Fatal("DatabasePath empty")
	}
	if filepath.Base(cfg.DatabasePath) != "statescore.db" {
		t.Fatalf("DatabasePath = %q", cfg.DatabasePath)
	}
}
