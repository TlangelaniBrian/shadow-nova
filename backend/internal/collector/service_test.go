package collector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"shadow-nova/backend/internal/ai"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/models"
	"testing"
	"time"
)

func TestCollectAll(t *testing.T) {
	// Mock DB
	mockDB := &database.MockService{
		GetContentSourcesFunc: func(ctx context.Context) ([]models.ContentSource, error) {
			return []models.ContentSource{
				{ID: 1, Name: "Test Source", Type: "blog_rss", URL: "http://test.com/feed"},
			}, nil
		},
		CreateContentItemFunc: func(ctx context.Context, item *models.ContentItem) error {
			if item.Title != "Test Item" {
				t.Errorf("expected item title 'Test Item', got '%s'", item.Title)
			}
			return nil
		},
	}

	// Mock FetchFeed
	originalFetchFeed := FetchFeed
	defer func() { FetchFeed = originalFetchFeed }()
	
	FetchFeed = func(url string) ([]ContentMetadata, error) {
		return []ContentMetadata{
			{Title: "Test Item", Description: "Desc", URL: "http://test.com/item", PublishedAt: time.Now()},
		}, nil
	}

	// Service
	service := New(mockDB, nil)

	// Execute
	err := service.CollectAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProcessUnprocessedItems(t *testing.T) {
	// Mock DB
	mockDB := &database.MockService{
		GetUnprocessedItemsFunc: func(ctx context.Context, limit int) ([]models.ContentItem, error) {
			return []models.ContentItem{
				{ID: 1, Title: "Test Item", Description: "Desc"},
			}, nil
		},
		UpdateContentItemAIFunc: func(ctx context.Context, item *models.ContentItem) error {
			if item.AISummary != "Summary" {
				t.Errorf("expected summary 'Summary', got '%s'", item.AISummary)
			}
			return nil
		},
	}

	// Mock AI Service with Test Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"candidates": [{
				"content": {
					"parts": [{"text": "{\"summary\": \"Summary\", \"tags\": [\"Tag\"], \"difficulty\": \"Beginner\"}"}]
				}
			}]
		}`))
	}))
	defer server.Close()

	aiService := ai.NewAIService()
	aiService.SetBaseURL(server.URL)
	aiService.HTTPClient = server.Client()

	// Service
	service := New(mockDB, aiService)

	// Execute
	err := service.ProcessUnprocessedItems(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
