package sso

import (
	"vitaliiPsl/synthesizer/internal/users"
	"golang.org/x/oauth2"
)

type SSOProvider interface {
	AuthCodeURL(state string) string
	Exchange(code string) (*oauth2.Token, error)
	FetchUserInfo(token *oauth2.Token) (*users.UserDto, error)
}
