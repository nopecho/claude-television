package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, err := LoadFrom(tmpDir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !cfg.Channels.AutoSync {
		t.Error("expected auto_sync default true")
	}
	if cfg.Channels.CacheTTL != "24h" {
		t.Errorf("expected cache_ttl 24h, got %s", cfg.Channels.CacheTTL)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{
		Channels: ChannelsConfig{
			AutoSync: true,
			CacheTTL: "12h",
			Pins:     []string{"my-project"},
			Groups:   map[string][]string{"work": {"api", "web"}},
		},
		Editor: "vim",
	}
	if err := SaveTo(cfg, tmpDir); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadFrom(tmpDir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Editor != "vim" {
		t.Errorf("expected vim, got %s", loaded.Editor)
	}
	if len(loaded.Channels.Pins) != 1 || loaded.Channels.Pins[0] != "my-project" {
		t.Errorf("unexpected pins: %v", loaded.Channels.Pins)
	}
	if len(loaded.Channels.Groups["work"]) != 2 {
		t.Errorf("unexpected groups: %v", loaded.Channels.Groups)
	}
}

func TestConfigDir(t *testing.T) {
	dir := ConfigDir()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "ctv")
	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}
