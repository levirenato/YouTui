// Package config to read config file
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Theme ThemeConfig `toml:"theme"`
}

type ThemeConfig struct {
	Active     string `toml:"active"`
	CustomPath string `toml:"custom_path,omitempty"`
}

func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "youtui", "youtui.conf")
}

func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()

	cfg := &Config{
		Theme: ThemeConfig{
			Active: "catppuccin-mocha",
		},
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	_, err := toml.DecodeFile(configPath, cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func SaveConfig(cfg *Config) (err error) {
	configPath := GetConfigPath()

	dir := filepath.Dir(configPath)
	if err = os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	encoder := toml.NewEncoder(f)
	if err = encoder.Encode(cfg); err != nil {
		return err
	}
	return nil
}
