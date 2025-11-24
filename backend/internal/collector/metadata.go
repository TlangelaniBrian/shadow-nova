package collector

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type ContentMetadata struct {
	Title       string
	Description string
	URL         string
	ImageURL    string
	PublishedAt time.Time
}

// FetchMetadata scrapes a URL for OpenGraph tags to get metadata
func FetchMetadata(url string) (*ContentMetadata, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Limit reader to avoid reading huge files (e.g. 5MB)
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	html := string(bodyBytes)

	meta := &ContentMetadata{}

	// Simple regex for OpenGraph tags (robust enough for most sites without full HTML parser)
	// <meta property="og:title" content="..." />
	meta.Title = extractMetaTag(html, "og:title")
	if meta.Title == "" {
		meta.Title = extractTitleTag(html)
	}

	meta.Description = extractMetaTag(html, "og:description")
	if meta.Description == "" {
		meta.Description = extractMetaTag(html, "description")
	}

	meta.ImageURL = extractMetaTag(html, "og:image")

	// Try to parse date (very basic, usually complex)
	// For now, default to current time if not found
	meta.PublishedAt = time.Now()

	return meta, nil
}

func extractMetaTag(html, property string) string {
	// Matches <meta property="og:title" content="Title" /> or <meta name="description" content="..." />
	// This is a simplified regex and might miss edge cases, but works for standard OG tags.
	// We look for property="property" or name="property"
	re := regexp.MustCompile(fmt.Sprintf(`(?i)<meta\s+(?:property|name)=["']%s["']\s+content=["'](.*?)["']`, regexp.QuoteMeta(property)))
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return htmlUnescape(matches[1])
	}
	return ""
}

func extractTitleTag(html string) string {
	re := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return htmlUnescape(matches[1])
	}
	return ""
}

func htmlUnescape(s string) string {
	// Basic unescape for common entities
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&#39;", "'")
	return s
}
