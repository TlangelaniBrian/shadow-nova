package collector

import (
	"context"
	"fmt"
	"log"
	"shadow-nova/backend/internal/ai"
	"shadow-nova/backend/internal/database"
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

	log.Printf("Processing %d items with AI...", len(items))

	for _, item := range items {
		// Call Gemini
		result, err := s.ai.GenerateSummary(ctx, item.Title, item.Description)
		if err != nil {
			log.Printf("Failed to generate summary for item %d: %v", item.ID, err)
			continue
		}

		// Update item
		item.AISummary = result.Summary
		item.AITags = result.Tags
		item.AIDifficulty = result.Difficulty
		item.ProcessedByAI = true

		if err := s.db.UpdateContentItemAI(ctx, &item); err != nil {
			log.Printf("Failed to update item %d: %v", item.ID, err)
		}
		
		// Sleep briefly to be nice to the API
		time.Sleep(1 * time.Second)
	}

	return nil
}

// CollectAll fetches content from all registered sources
func (s *Service) CollectAll(ctx context.Context) error {
	sources, err := s.db.GetContentSources(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sources: %w", err)
	}

	log.Printf("Starting collection for %d sources...", len(sources))

	for _, source := range sources {
		if err := s.processSource(ctx, &source); err != nil {
			log.Printf("Error processing source %s: %v", source.Name, err)
			// Continue with next source
		}
	}

	return nil
}

func (s *Service) processSource(ctx context.Context, source *models.ContentSource) error {
	log.Printf("Fetching %s (%s)...", source.Name, source.URL)

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

	log.Printf("Found %d items for %s", len(items), source.Name)

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
	
	log.Printf("Saved %d new items for %s", count, source.Name)
	return nil
}
