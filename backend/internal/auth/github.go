package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubAuthService struct {
	config *oauth2.Config
}

func NewGitHubAuthService() *GitHubAuthService {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURL := os.Getenv("GITHUB_REDIRECT_URL")

	if clientID == "" || clientSecret == "" {
		// Log warning or return error depending on strictness
		// For now, we'll allow it to be empty but methods will fail if called
		fmt.Println("Warning: GITHUB_CLIENT_ID or GITHUB_CLIENT_SECRET not set")
	}

	if redirectURL == "" {
		redirectURL = "http://localhost:3000/api/auth/github/callback"
	}

	return &GitHubAuthService{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"read:user", "user:email", "repo"}, // Request repo access
			Endpoint:     github.Endpoint,
		},
	}
}

func (s *GitHubAuthService) GetLoginURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *GitHubAuthService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (s *GitHubAuthService) GetUser(ctx context.Context, token *oauth2.Token) (*GitHubUser, error) {
	client := s.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status: %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &user, nil
}
