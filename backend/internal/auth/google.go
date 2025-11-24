package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthService struct {
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
	jwtSecret []byte
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}


func NewGoogleAuthService() (*GoogleAuthService, error) {
	ctx := context.Background()
	
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	
	if clientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID not set")
	}
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/auth/callback"
	}
	
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}
	
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	
	return &GoogleAuthService{
		config:   config,
		verifier: verifier,
		jwtSecret: []byte(jwtSecret),
	}, nil
}

// GetAuthURL returns the Google OAuth URL for login
func (s *GoogleAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// VerifyGoogleToken verifies a Google ID token and returns user info
func (s *GoogleAuthService) VerifyGoogleToken(ctx context.Context, idToken string) (*GoogleUserInfo, error) {
	token, err := s.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	
	var userInfo GoogleUserInfo
	if err := token.Claims(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}
	
	return &userInfo, nil
}

// ExchangeCodeForToken exchanges authorization code for tokens
func (s *GoogleAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

	// Generate our JWT
	// Note: Google User ID is a string, but our system uses int IDs.
	// This is a mismatch. In a real system, we'd look up the user in DB to get their INT ID.
	// For now, let's assume we have a way to get the INT ID or we change GenerateJWT to accept string.
	// But GenerateJWT in auth.go accepts int.
	// We need to fix this architectural mismatch.
	// The GoogleCallback in handlers/auth.go (which calls this) should handle the DB lookup/creation and then generate JWT.
	// THIS method `HandleGoogleCallback` inside the service seems to be doing too much (HTTP handling + Logic).
	// It should probably just return the user info.
	
	// However, to fix the immediate compilation error without major refactor:
	// We'll comment out the JWT generation here and let the handler do it, 
	// OR we temporarily convert string to int (bad) or change GenerateJWT to string (better but affects other parts).
	
	// Actually, looking at routes.go, `authHandler.GoogleCallback` is used.
	// Let's check `backend/internal/handlers/auth.go` to see what it does.
	// If `GoogleAuthService` has `HandleGoogleCallback`, maybe it's being used directly?
	// `routes.go` says: `r.Get("/auth/google/callback", authHandler.GoogleCallback)`
	
	// So `GoogleAuthService.HandleGoogleCallback` might be unused or legacy code if `authHandler` implements it.
	// Let's check `handlers/auth.go`.
	
	// For now, I will remove the JWT methods from here to avoid conflict/duplication.
