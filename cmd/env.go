package cmd

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var EnvVars *EnvConfig

type EnvConfig struct {
	Environment    string
	Port           int
	DBUrl          string
	TokenSecret    string
	SmtpHost       string
	SmtpPort       int
	GmailUser      string
	AppPassword    string
	GhClientId     string // github
	GhClientSecret string
	GhRedirectUrl  string
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
	dbUrl := os.Getenv("DATABASE_URL")
	tokenSecret := os.Getenv("JWT_SECRET")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	gmailUser := os.Getenv("GMAIL_USERNAME")
	appPwd := os.Getenv("GMAIL_APP_PASSWORD")
	ghClientId := os.Getenv("GITHUB_CLIENT_ID")
	ghClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	ghRedirectUrl := os.Getenv("GITHUB_REDIRECT_URL")

	// Environment
	environment = strings.ToLower(environment)
	isValid := slices.Contains(validEnvs, environment)
	if !isValid {
		return nil, fmt.Errorf("Invalid ENVIRONMENT value: %s", environment)
	}
	cfg.Environment = environment
	// Port
	if port == "" {
		return nil, fmt.Errorf("PORT environment variable is missing.")
	}
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("Invalid PORT value: %w", err)
	}
	cfg.Port = portNum
	// Database URL
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is missing.")
	}
	cfg.DBUrl = dbUrl
	// Token secret
	if tokenSecret == "" {
		return nil, fmt.Errorf("TOKEN_SECRET environment variable is missing.")
	}
	cfg.TokenSecret = tokenSecret
	// SMTP host and port
	if smtpHost == "" {
		return nil, fmt.Errorf("GMAIL_USERNAME environment variable is missing.")
	}
	cfg.SmtpHost = smtpHost
	if smtpPort == "" {
		return nil, fmt.Errorf("SMTP_PORT environment variable is missing.")
	}
	cfg.SmtpPort, err = strconv.Atoi(smtpPort)
	if err != nil {
		return nil, err
	}
	// Gmail user
	if gmailUser == "" {
		return nil, fmt.Errorf("GMAIL_USERNAME environment variable is missing.")
	}
	cfg.GmailUser = gmailUser
	// App password
	if appPwd == "" {
		return nil, fmt.Errorf("GMAIL_APP_PASSWORD environment variable is missing.")
	}
	cfg.AppPassword = appPwd
	// GitHub Client Id
	if ghClientId == "" {
		return nil, fmt.Errorf("GITHUB_CLIENT_ID environment variable is missing.")
	}
	cfg.GhClientId = ghClientId
	// GitHub Client Secret
	if ghClientSecret == "" {
		return nil, fmt.Errorf("GITHUB_CLIENT_SECRET environment variable is missing.")
	}
	cfg.GhClientSecret = ghClientSecret
	// GitHub Redirect Url
	if ghRedirectUrl == "" {
		return nil, fmt.Errorf("GITHUB_REDIRECT_URL environment variable is missing.")
	}
	cfg.GhRedirectUrl = ghRedirectUrl

	return cfg, nil
}
