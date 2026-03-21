package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Channels ChannelsConfig `json:"channels"`
	Editor   string         `json:"editor"`
}

type ChannelsConfig struct {
	AutoSync bool                `json:"auto_sync"`
	CacheTTL string              `json:"cache_ttl"`
	Pins     []string            `json:"pins"`
	Groups   map[string][]string `json:"groups"`
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ctv")
}

func Load() (*Config, error) {
	return LoadFrom(ConfigDir())
}

func LoadFrom(dir string) (*Config, error) {
	path := filepath.Join(dir, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	cfg := defaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func Save(cfg *Config) error {
	return SaveTo(cfg, ConfigDir())
}

func SaveTo(cfg *Config, dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)
}

func ParseDuration(s string) time.Duration {
	if s == "" {
		return 24 * time.Hour
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}

func defaultConfig() *Config {
	return &Config{
		Channels: ChannelsConfig{
			AutoSync: true,
			CacheTTL: "24h",
			Groups:   map[string][]string{},
		},
	}
}
