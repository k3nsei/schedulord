package config

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	GitHubToken string     `mapstructure:"gh_token"`
	Tasks       []TaskSpec `mapstructure:"tasks"`
	TimeZone    string     `mapstructure:"tz"`
}

func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("gh_token", "")
	v.SetDefault("tasks", []TaskSpec{})
	v.SetDefault("tz", "Europe/Warsaw")

	// Load config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/schedulord")
	v.AddConfigPath("$HOME/.config/schedulord")

	// Load dotenv
	_ = godotenv.Load()

	// Override with environment variables
	v.SetEnvPrefix("SCHEDULORD")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	slog.Info("Configuration has been loaded")

	return &cfg, nil
}

func ReadGithubToken(cfg *Config) (string, error) {
	token := cfg.GitHubToken

	if token != "" {
		return cfg.GitHubToken, nil
	}

	tokenBytes, err := os.ReadFile("/run/secrets/gh_token")
	if err != nil {
		return "", err
	}

	token = string(tokenBytes)

	if token == "" {
		return "", errors.New("GitHub toke is missing")
	}

	return token, nil
}
