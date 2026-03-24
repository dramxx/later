package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	Gist GistConfig
}

type GistConfig struct {
	Token  string
	GistID string
}

const template = `[gist]
token = ""
gist_id = ""
`

func Path() string {
	var configDir string

	if runtime.GOOS == "windows" {
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	} else {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}

	return filepath.Join(configDir, "later", "config.toml")
}

func EnsureFile() (string, error) {
	configPath := Path()
	dir := filepath.Dir(configPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create config directory: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.WriteFile(configPath, []byte(template), 0644); err != nil {
			return "", fmt.Errorf("create config file: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("stat config file: %w", err)
	}

	return configPath, nil
}

func Save(cfg *Config) error {
	configPath, err := EnsureFile()
	if err != nil {
		return err
	}

	content := fmt.Sprintf("[gist]\ntoken = %q\ngist_id = %q\n", cfg.Gist.Token, cfg.Gist.GistID)
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

func Load() (*Config, error) {
	configPath := Path()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("config not found at %s — run 'later config --init'", configPath)
	}

	cfg := &Config{}
	err = parseConfig(string(data), cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Gist.Token == "" {
		return nil, fmt.Errorf("missing 'token' in config — run 'later config --init'")
	}
	if cfg.Gist.GistID == "" {
		return nil, fmt.Errorf("missing 'gist_id' in config — run 'later config --init'")
	}

	return cfg, nil
}

func parseConfig(content string, cfg *Config) error {
	currentSection := ""

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			continue
		}

		if currentSection == "gist" {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\"")

				switch key {
				case "token":
					cfg.Gist.Token = value
				case "gist_id":
					cfg.Gist.GistID = value
				}
			}
		}
	}

	return nil
}
