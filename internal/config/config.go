package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// FirstRun indicates if this is the first time running the app
	FirstRun bool `mapstructure:"first_run"`

	// Alias is the shell alias configured for graft (e.g., "gg")
	Alias string `mapstructure:"alias"`

	// DefaultBranch is the default base branch for new branches
	DefaultBranch string `mapstructure:"default_branch"`

	// Editor is the preferred editor for commit messages
	Editor string `mapstructure:"editor"`

	// PushAfterCommit prompts to push after committing
	PushAfterCommit bool `mapstructure:"push_after_commit"`

	// ShowCommitPreview shows commit diff preview
	ShowCommitPreview bool `mapstructure:"show_commit_preview"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		FirstRun:          true,
		Alias:             "",
		DefaultBranch:     "main",
		Editor:            "",
		PushAfterCommit:   true,
		ShowCommitPreview: true,
	}
}

// configDir returns the config directory path
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "graft"), nil
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return DefaultConfig(), err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	// Set defaults
	viper.SetDefault("first_run", true)
	viper.SetDefault("alias", "")
	viper.SetDefault("default_branch", "main")
	viper.SetDefault("editor", "")
	viper.SetDefault("push_after_commit", true)
	viper.SetDefault("show_commit_preview", true)

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, use defaults
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return DefaultConfig(), err
	}

	return &cfg, nil
}

// Save saves the configuration to disk
func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	viper.Set("first_run", cfg.FirstRun)
	viper.Set("alias", cfg.Alias)
	viper.Set("default_branch", cfg.DefaultBranch)
	viper.Set("editor", cfg.Editor)
	viper.Set("push_after_commit", cfg.PushAfterCommit)
	viper.Set("show_commit_preview", cfg.ShowCommitPreview)

	configPath := filepath.Join(dir, "config.yaml")
	return viper.WriteConfigAs(configPath)
}

// MarkOnboardingComplete marks the onboarding as complete
func (c *Config) MarkOnboardingComplete() {
	c.FirstRun = false
}

// SetAlias sets the shell alias
func (c *Config) SetAlias(alias string) {
	c.Alias = alias
}
