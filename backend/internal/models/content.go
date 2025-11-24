package models

import "time"

type ContentSource struct {
	ID            int       `json:"id"`
	Name          string    `json:"name" validate:"required"`
	Type          string    `json:"type" validate:"required,oneof=youtube_channel blog_rss twitter_handle"`
	URL           string    `json:"url" validate:"required,url"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type ContentItem struct {
	ID            int       `json:"id"`
	SourceID      int       `json:"source_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	URL           string    `json:"url"`
	ImageURL      string    `json:"image_url"`
	PublishedAt   time.Time `json:"published_at"`
	FetchedAt     time.Time `json:"fetched_at"`
	
	// AI Metadata
	AISummary     string    `json:"ai_summary"`
	AITags        []string  `json:"ai_tags"`
	AIDifficulty  string    `json:"ai_difficulty"`
	ProcessedByAI bool      `json:"processed_by_ai"`
}

type CreateSourceRequest struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required,oneof=youtube_channel blog_rss twitter_handle"`
	URL  string `json:"url" validate:"required,url"`
}
