package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"slices"

	"github.com/joho/godotenv"
)

var EnvVars *EnvConfig

type EnvConfig struct {
	Environment  string
	Port         int
	Domain       string
	CookieSecure bool
	DBUrl        string
}

func NewEnvConfig() (*EnvConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf(".env file not found")
	}

	cfg := &EnvConfig{}
	validEnvs := []string{"development", "testing", "production"}

	environment := os.Getenv("ENVIRONMENT")
	port := os.Getenv("PORT")
	domain := os.Getenv("DOMAIN")
	cookieSecure := os.Getenv("COOKIE_SECURE")
	dbUrl := os.Getenv("DATABASE_URL")

	environment = strings.ToLower(environment)
	isValid := slices.Contains(validEnvs, environment)
	if !isValid {
		return nil, fmt.Errorf("Invalid ENVIRONMENT value: %s", environment)
	}
	cfg.Environment = environment
	if port == "" {
		return nil, fmt.Errorf("PORT environment variable is missing.")
	}
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("Invalid PORT value: %w", err)
	}
	cfg.Port = portNum
	if domain == "" {
		return nil, fmt.Errorf("DOMAIN environment variable is missing.")
	}
	cfg.Domain = domain
	if cookieSecure == "" {
		return nil, fmt.Errorf("COOKIE_SECURE environment variable is missing.")
	}
	cookieSecurityBool, err := strconv.ParseBool(cookieSecure)
	if err != nil {
		return nil, fmt.Errorf("Invalid COOKIE_SECURE value %s", cookieSecure)
	}
	cfg.CookieSecure = cookieSecurityBool
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is missing.")
	}
	cfg.DBUrl = dbUrl

	return cfg, nil
}
