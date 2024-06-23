package sso

import (
	"context"
	"encoding/json"
	"io"

	"golang.org/x/oauth2"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/users"
)

const (
	USER_INFO_ENDPOINT  = "https://api.github.com/user"
	USER_EMAIL_ENDPOINT = "https://api.github.com/user/emails"
)

type GithubProvider struct {
	config *oauth2.Config
}

func NewGithubProvider(cfg *oauth2.Config) *GithubProvider {
	return &GithubProvider{config: cfg}
}

func (p *GithubProvider) AuthCodeURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *GithubProvider) Exchange(code string) (*oauth2.Token, error) {
	return p.config.Exchange(context.Background(), code)
}

func (p *GithubProvider) FetchUserInfo(token *oauth2.Token) (*users.UserDto, error) {
	client := p.config.Client(context.Background(), token)
	resp, err := client.Get(USER_INFO_ENDPOINT)
	if err != nil {
		return nil, service_errors.NewErrBadGateway("Failed to fetch user info")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, service_errors.NewErrInternalServer("Failed to read user info")
	}

	var user *users.UserDto
	user, err = p.buildUserModel(data)
	if err != nil {
		return nil, err
	}

	if user.Email == "" {
		user.Email, err = p.fetchUserEmail(token)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (p *GithubProvider) fetchUserEmail(token *oauth2.Token) (string, error) {
	client := p.config.Client(context.Background(), token)
	resp, err := client.Get(USER_EMAIL_ENDPOINT)
	if err != nil {
		return "", service_errors.NewErrBadGateway("Failed to fetch user email")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", service_errors.NewErrInternalServer("Failed to read user email")
	}

	var githubEmails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(data, &githubEmails); err != nil {
		return "", service_errors.NewErrInternalServer("Failed to unmarshal user email")
	}

	for _, email := range githubEmails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	return "", service_errors.NewErrNotFound("User email not found")
}

func (p *GithubProvider) buildUserModel(data []byte) (*users.UserDto, error) {
	var githubUser struct {
		Login     string `json:"login"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.Unmarshal(data, &githubUser); err != nil {
		return nil, service_errors.NewErrInternalServer("Failed to unmarshal user info")
	}

	return &users.UserDto{
		Email:      githubUser.Email,
		Username:   githubUser.Login,
		PictureUrl: githubUser.AvatarURL,
	}, nil
}
