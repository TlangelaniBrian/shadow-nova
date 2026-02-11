




package database

import (
	"context"
	"errors"
	"fmt"
	"shadow-nova/backend/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *service) CreateContentSource(ctx context.Context, source *models.ContentSource) error {
	query := `
		INSERT INTO content_sources (name, type, url)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := s.db.QueryRow(ctx, query, source.Name, source.Type, source.URL).Scan(&source.ID, &source.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create content source: %w", err)
	}

	return nil
}

func (s *service) GetContentSources(ctx context.Context, limit, offset int) ([]models.ContentSource, error) {
	query := `
		SELECT id, name, type, url, last_fetched_at, created_at
		FROM content_sources
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query content sources: %w", err)
	}
	defer rows.Close()

	var sources []models.ContentSource
	for rows.Next() {
		var src models.ContentSource
		var lastFetched *time.Time

		if err := rows.Scan(&src.ID, &src.Name, &src.Type, &src.URL, &lastFetched, &src.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan content source: %w", err)
		}

		if lastFetched != nil {
			src.LastFetchedAt = *lastFetched
		}
		sources = append(sources, src)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating content sources: %w", err)
	}
	return sources, nil
}

func (s *service) GetContentSourcesCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM content_sources`
	err := s.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (s *service) CreateContentItem(ctx context.Context, item *models.ContentItem) error {
	query := `
		INSERT INTO content_items (source_id, title, description, url, image_url, published_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (url) DO NOTHING
		RETURNING id
	`
	err := s.db.QueryRow(ctx, query, item.SourceID, item.Title, item.Description, item.URL, item.ImageURL, item.PublishedAt).Scan(&item.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // ON CONFLICT DO NOTHING: item already exists
		}
		return fmt.Errorf("failed to create content item: %w", err)
	}

	return nil
}

func (s *service) GetUnprocessedItems(ctx context.Context, limit int) ([]models.ContentItem, error) {
	query := `
		SELECT id, source_id, title, description, url, image_url, published_at
		FROM content_items
		WHERE processed_by_ai = FALSE
		LIMIT $1
	`

	rows, err := s.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query unprocessed items: %w", err)
	}
	defer rows.Close()

	var items []models.ContentItem
	for rows.Next() {
		var item models.ContentItem
		if err := rows.Scan(&item.ID, &item.SourceID, &item.Title, &item.Description, &item.URL, &item.ImageURL, &item.PublishedAt); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating unprocessed items: %w", err)
	}

	return items, nil
}

func (s *service) UpdateContentItemAI(ctx context.Context, item *models.ContentItem) error {
	query := `
		UPDATE content_items
		SET ai_summary = $1, ai_tags = $2, ai_difficulty = $3, processed_by_ai = TRUE
		WHERE id = $4
	`
	
	_, err := s.db.Exec(ctx, query, item.AISummary, item.AITags, item.AIDifficulty, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update content item AI: %w", err)
	}
	return nil
}
