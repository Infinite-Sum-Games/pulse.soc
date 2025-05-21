package cmd

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var GithubOAuthConfig *oauth2.Config

func OAuthInit() {
	cfg := &oauth2.Config{
		ClientID:     EnvVars.GhClientId,
		ClientSecret: EnvVars.GhClientSecret,
		RedirectURL:  EnvVars.GhRedirectUrl,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}

	GithubOAuthConfig = cfg
}
