package cmd

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var GithubOAuthConfig *oauth2.Config

func OAuthInit() {
	cfg := &oauth2.Config{
		ClientID:     AppConfig.GitHub.ClientID,
		ClientSecret: AppConfig.GitHub.ClientSecret,
		RedirectURL:  AppConfig.GitHub.RedirectURL,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}

	GithubOAuthConfig = cfg
}
