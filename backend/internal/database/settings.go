package database

import (
	"context"
	"fmt"
)

func (s *service) GetSystemSetting(ctx context.Context, key string) (string, error) {
	query := `SELECT value FROM system_settings WHERE key = $1`
	
	var value string
	err := s.db.QueryRow(ctx, query, key).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("failed to get setting %s: %w", key, err)
	}
	
	return value, nil
}

func (s *service) UpdateSystemSetting(ctx context.Context, key, value string) error {
	query := `
		INSERT INTO system_settings (key, value, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (key) 
		DO UPDATE SET value = $2, updated_at = CURRENT_TIMESTAMP
	`
	
	_, err := s.db.Exec(ctx, query, key, value)
	if err != nil {
		return fmt.Errorf("failed to update setting %s: %w", key, err)
	}
	
	return nil
}
