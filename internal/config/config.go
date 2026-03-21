package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Scan ScanConfig `mapstructure:"scan"`
}

type ScanConfig struct {
	Roots  []string `mapstructure:"roots"`
	Ignore []string `mapstructure:"ignore"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ctv")
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir())

	viper.SetDefault("scan.roots", []string{})
	viper.SetDefault("scan.ignore", []string{"node_modules", ".git", "vendor"})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("config read: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}
	return &cfg, nil
}

func Save() error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	return viper.WriteConfigAs(filepath.Join(dir, "config.yaml"))
}

func AddScanRoot(path string) error {
	roots := viper.GetStringSlice("scan.roots")
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	for _, r := range roots {
		if r == abs {
			return fmt.Errorf("path already registered: %s", abs)
		}
	}
	roots = append(roots, abs)
	viper.Set("scan.roots", roots)
	return Save()
}

func RemoveScanRoot(path string) error {
	roots := viper.GetStringSlice("scan.roots")
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	filtered := make([]string, 0, len(roots))
	found := false
	for _, r := range roots {
		if r == abs {
			found = true
			continue
		}
		filtered = append(filtered, r)
	}
	if !found {
		return fmt.Errorf("path not found: %s", abs)
	}
	viper.Set("scan.roots", filtered)
	return Save()
}
