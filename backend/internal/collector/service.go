package collector

import (
	"context"
	"fmt"
	"shadow-nova/backend/internal/ai"
	"shadow-nova/backend/internal/database"
	"shadow-nova/backend/internal/logging"
	"shadow-nova/backend/internal/metrics"
	"shadow-nova/backend/internal/models"
	"time"
)

type Service struct {
	db database.Service
	ai *ai.AIService
}

func New(db database.Service, ai *ai.AIService) *Service {
	return &Service{db: db, ai: ai}
}

// ProcessUnprocessedItems fetches items without AI metadata and processes them
func (s *Service) ProcessUnprocessedItems(ctx context.Context) error {
	// Fetch 5 items at a time to respect rate limits
	items, err := s.db.GetUnprocessedItems(ctx, 5)
	if err != nil {
		return fmt.Errorf("failed to get unprocessed items: %w", err)
	}

	if len(items) == 0 {
		return nil
	}

	logging.Info("processing items with ai", "count", len(items))

	for _, item := range items {
		// Track AI processing duration
		start := time.Now()

		// Call Gemini
		result, err := s.ai.GenerateSummary(ctx, item.Title, item.Description)
		duration := time.Since(start).Seconds()

		if err != nil {
			logging.Warn("failed to generate summary", "item_id", item.ID, "error", err)
			metrics.AIProcessingErrors.Inc()
			continue
		}

		metrics.AIProcessingDuration.Observe(duration)

		// Update item
		item.AISummary = result.Summary
		item.AITags = result.Tags
		item.AIDifficulty = result.Difficulty
		item.ProcessedByAI = true

		if err := s.db.UpdateContentItemAI(ctx, &item); err != nil {
			logging.Error("failed to update item with ai metadata", err, "item_id", item.ID)
		}

		// Sleep briefly to be nice to the API
		time.Sleep(1 * time.Second)
	}

	return nil
}

// CollectAll fetches content from all registered sources
func (s *Service) CollectAll(ctx context.Context) error {
	// Get all sources (no pagination needed for collector)
	sources, err := s.db.GetContentSources(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to get sources: %w", err)
	}

	logging.Info("starting content collection", "source_count", len(sources))

	for _, source := range sources {
		if err := s.processSource(ctx, &source); err != nil {
			logging.Error("failed to process source", err, "source_name", source.Name, "source_url", source.URL)
			// Continue with next source
		}
	}

	return nil
}

func (s *Service) processSource(ctx context.Context, source *models.ContentSource) error {
	logging.Info("fetching content from source", "source_name", source.Name, "source_url", source.URL, "source_type", source.Type)

	var items []ContentMetadata
	var err error

	// Determine how to fetch based on type
	switch source.Type {
	case "youtube_channel", "blog_rss":
		// For YouTube, we assume the URL is the RSS feed URL or we convert it
		// YouTube Channel RSS: https://www.youtube.com/feeds/videos.xml?channel_id=CHANNEL_ID
		// If user provided regular URL, we might need to resolve it.
		// For now, assume user provides the RSS URL or we handle it simply.
		items, err = FetchFeed(source.URL)
	default:
		// For generic URL, just fetch metadata for that single page?
		// Or maybe it's not supported yet.
		return fmt.Errorf("unsupported source type: %s", source.Type)
	}

	if err != nil {
		return err
	}

	logging.Info("fetched items from source", "source_name", source.Name, "item_count", len(items))

	// Save items
	count := 0
	for _, meta := range items {
		// Skip if too old? (e.g. older than 30 days)
		// For now, save all.

		item := &models.ContentItem{
			SourceID:    source.ID,
			Title:       meta.Title,
			Description: meta.Description,
			URL:         meta.URL,
			ImageURL:    meta.ImageURL,
			PublishedAt: meta.PublishedAt,
		}

		// Temporary fix: I'll need to update metadata.go to include URL/Link.
		// For now, let's assume I'll fix it in the next step.

		if err := s.db.CreateContentItem(ctx, item); err == nil {
			count++
		}
	}

	metrics.ContentItemsCollected.Add(float64(count))
	logging.Info("saved new items from source", "source_name", source.Name, "saved_count", count)
	return nil
}
