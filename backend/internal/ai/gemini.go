package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type AIService struct {
	apiKey     string
	baseURL    string
	HTTPClient *http.Client
}

func NewAIService() *AIService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: GEMINI_API_KEY not set")
	}
	return &AIService{
		apiKey:     apiKey,
		baseURL:    "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent",
		HTTPClient: &http.Client{},
	}
}

func (s *AIService) SetBaseURL(url string) {
	s.baseURL = url
}

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type AIResult struct {
	Summary    string   `json:"summary"`
	Tags       []string `json:"tags"`
	Difficulty string   `json:"difficulty"`
}

// GenerateSummary uses Gemini to summarize content and generate tags
func (s *AIService) GenerateSummary(ctx context.Context, title, description string) (*AIResult, error) {
	url := fmt.Sprintf("%s?key=%s", s.baseURL, s.apiKey)

	prompt := fmt.Sprintf(`
		Analyze the following content metadata (Title and Description) for a learning platform.
		
		Title: %s
		Description: %s
		
		Please provide a JSON response with the following fields:
		1. "summary": A concise 2-sentence summary of what this content teaches.
		2. "tags": A list of 3-5 relevant technical tags (e.g., "Go", "React", "Database").
		3. "difficulty": One of "Beginner", "Intermediate", "Advanced".
		
		Return ONLY the JSON.
	`, title, description)

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini api returned status: %d", resp.StatusCode)
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content generated")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	
	// Clean up markdown code blocks if present
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	var result AIResult
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI result: %w", err)
	}

	return &result, nil
}
