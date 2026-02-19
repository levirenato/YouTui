// Package config to read config file
package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Theme ThemeConfig `toml:"theme"`
	UI    UIConfig    `toml:"ui"`
}

type ThemeConfig struct {
	Active     string `toml:"active"`
	CustomPath string `toml:"custom_path,omitempty"`
}

type UIConfig struct {
	Language string `toml:"language, omitempty"`
}

func GetConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "youtui")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "youtui")
}

func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), "youtui.conf")
}

func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()

	cfg := &Config{
		Theme: ThemeConfig{
			Active: "catppuccin-mocha",
		},
		UI: UIConfig{
			Language: detectDefaultLanguage(),
		},
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	_, err := toml.DecodeFile(configPath, cfg)
	if err != nil {
		return cfg, err
	}

	if strings.TrimSpace(cfg.UI.Language) == "" {
		cfg.UI.Language = detectDefaultLanguage()
	}

	cfg.UI.Language = normalizeLang(cfg.UI.Language)

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

func detectDefaultLanguage() string {
	envs := []string{
		os.Getenv("LC_ALL"),
		os.Getenv("LC_MESSAGES"),
		os.Getenv("LANG"),
	}
	for _, v := range envs {
		if v = strings.ToLower(v); v != "" {
			if strings.HasPrefix(v, "en") {
				return "en"
			}
			if strings.HasPrefix(v, "pt") {
				return "pt"
			}
		}
	}
	return "pt"
}

func normalizeLang(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	switch {
	case s == "en", strings.HasPrefix(s, "en-"), strings.HasPrefix(s, "en_"):
		return "en"
	case s == "pt", strings.HasPrefix(s, "pt-"), strings.HasPrefix(s, "pt_"), s == "br":
		return "pt"
	default:
		return "pt"
	}
}
