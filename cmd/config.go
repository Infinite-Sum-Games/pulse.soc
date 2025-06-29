package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

type EnvConfig struct {
	Environment string `mapstructure:"environment"`
	Port        int    `mapstructure:"port"`
	DatabaseURL string `mapstructure:"database_url"`
	JWTSecret   string `mapstructure:"jwt_secret"`
	FrontendURL string `mapstructure:"frontend_url"`

	Valkey struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapsstructure:"port"`
		Username string `mapsstructure:"string"`
		Password string `mapsstructure:"password"`
	} `mapstructure:"valkey"`

	SMTP struct {
		Host        string `mapstructure:"host"`
		Port        int    `mapstructure:"port"`
		Username    string `mapstructure:"username"`
		AppPassword string `mapstructure:"app_password"`
	} `mapstructure:"smtp"`

	GitHub struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
		PAT          string `mapstructure:"personal_access_token"`
	} `mapstructure:"github"`
}

var AppConfig *EnvConfig

// LoadConfig reads configuration from config.toml file
func LoadConfig() (*EnvConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	// Set environment variables as override
	viper.SetEnvPrefix("PULSE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &EnvConfig{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *EnvConfig) error {
	validEnvs := []string{"development", "testing", "production"}
	result := slices.Contains(validEnvs, config.Environment)

	// Validate server and database configurations
	if !result {
		return fmt.Errorf("invalid environment: %s", config.Environment)
	}
	if config.Port <= 0 {
		return fmt.Errorf("invalid port number")
	}
	if config.DatabaseURL == "" {
		return fmt.Errorf("database URL is required")
	}
	if config.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	if config.FrontendURL == "" {
		return fmt.Errorf("frontend URL is required")
	}

	// Validate Valkey configuration
	if config.Valkey.Host == "" {
		return fmt.Errorf("valkey host is required")
	}
	if config.Valkey.Port <= 0 {
		return fmt.Errorf("invalid Valkey port")
	}

	// Validate SMTP configuration
	if config.SMTP.Host == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if config.SMTP.Port <= 0 {
		return fmt.Errorf("invalid SMTP port")
	}
	if config.SMTP.Username == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if config.SMTP.AppPassword == "" {
		return fmt.Errorf("SMTP app password is required")
	}

	// Validate GitHub configuration
	if config.GitHub.ClientID == "" {
		return fmt.Errorf("GitHub client ID is required")
	}
	if config.GitHub.ClientSecret == "" {
		return fmt.Errorf("GitHub client secret is required")
	}
	if config.GitHub.RedirectURL == "" {
		return fmt.Errorf("GitHub redirect URL is required")
	}
	if config.GitHub.PAT == "" {
		return fmt.Errorf("GitHub PAT is required")
	}

	return nil
}
