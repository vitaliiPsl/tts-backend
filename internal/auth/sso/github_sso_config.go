package sso

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func GithubSSOConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  os.Getenv("SSO_GITHUB_REDIRECT_URL"),
		ClientID:     os.Getenv("SSO_GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("SSO_GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}
