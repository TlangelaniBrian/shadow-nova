package collector

import (
	"testing"
	"time"
)

func TestExtractMetaTag(t *testing.T) {
	html := `
		<html>
			<head>
				<meta property="og:title" content="Test Title" />
				<meta name="description" content="Test Description" />
			</head>
		</html>
	`

	title := extractMetaTag(html, "og:title")
	if title != "Test Title" {
		t.Errorf("expected 'Test Title', got '%s'", title)
	}

	desc := extractMetaTag(html, "description")
	if desc != "Test Description" {
		t.Errorf("expected 'Test Description', got '%s'", desc)
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
		hasError bool
	}{
		{"Mon, 02 Jan 2006 15:04:05 -0700", time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*3600)), false},
		{"2023-10-27T10:00:00Z", time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC), false},
		{"invalid-date", time.Now(), true},
	}

	for _, test := range tests {
		result, err := parseDate(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("expected error for input '%s', got nil", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input '%s': %v", test.input, err)
			}
			if !result.Equal(test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		}
	}
}

func TestRSSParser(t *testing.T) {
	// Mock RSS Feed
	// Mock RSS Feed
	// rssXML := `...` (omitted for brevity)
	
	// We can't easily test FetchFeed without mocking HTTP, but we can test the parsing logic if we extract it.
	// For now, let's just test the helper functions which are critical.
	// Ideally, we should refactor FetchFeed to take an io.Reader.
	
	// We can't easily test FetchFeed without mocking HTTP, but we can test the parsing logic if we extract it.
	// For now, let's just test the helper functions which are critical.
	// Ideally, we should refactor FetchFeed to take an io.Reader.
}
