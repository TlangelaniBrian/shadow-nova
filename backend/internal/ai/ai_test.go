package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateSummary(t *testing.T) {
	// Mock Gemini API Response
	mockResponse := GeminiResponse{
		Candidates: []Candidate{
			{
				Content: Content{
					Parts: []Part{
						{
							Text: "```json\n{\"summary\": \"This is a test summary.\", \"tags\": [\"Go\", \"Test\"], \"difficulty\": \"Beginner\"}\n```",
						},
					},
				},
			},
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Verify query param
		key := r.URL.Query().Get("key")
		if key != "test-api-key" {
			t.Errorf("expected api key 'test-api-key', got '%s'", key)
		}

		// Send mock response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Initialize AIService with mock server URL
	service := NewAIService()
	service.apiKey = "test-api-key"
	service.baseURL = server.URL
	service.HTTPClient = server.Client()

	// Execute
	result, err := service.GenerateSummary(context.Background(), "Test Title", "Test Description")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Assert
	if result.Summary != "This is a test summary." {
		t.Errorf("expected summary 'This is a test summary.', got '%s'", result.Summary)
	}
	if len(result.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result.Tags))
	}
	if result.Difficulty != "Beginner" {
		t.Errorf("expected difficulty 'Beginner', got '%s'", result.Difficulty)
	}
}

func TestGenerateSummary_Error(t *testing.T) {
	// Create a test server that returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := NewAIService()
	service.baseURL = server.URL
	service.HTTPClient = server.Client()

	_, err := service.GenerateSummary(context.Background(), "Title", "Desc")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
