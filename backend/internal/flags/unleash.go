package flags

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Unleash/unleash-client-go/v4"
)

type Service interface {
	IsEnabled(feature string) bool
	Close()
}

type service struct{}

func New() (Service, error) {
	appName := os.Getenv("UNLEASH_APP_NAME")
	url := os.Getenv("UNLEASH_URL")
	token := os.Getenv("UNLEASH_TOKEN")

	if appName == "" {
		appName = "shadow-nova-backend"
	}
	if url == "" {
		// Default to a placeholder or local instance if not set
		url = "http://localhost:4242/api"
	}

	err := unleash.Initialize(
		unleash.WithAppName(appName),
		unleash.WithUrl(url),
		unleash.WithCustomHeaders(http.Header{"Authorization": []string{token}}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize unleash: %w", err)
	}

	return &service{}, nil
}

func (s *service) IsEnabled(feature string) bool {
	return unleash.IsEnabled(feature)
}

func (s *service) Close() {
	unleash.Close()
}
