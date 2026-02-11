package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// StoreIdempotentResponse stores a response for a given idempotency key.
// This allows future requests with the same key to return the cached response.
func (s *service) StoreIdempotentResponse(
	ctx context.Context,
	key string,
	userID int,
	path, method string,
	status int,
	body string,
	expiresAt time.Time,
) error {
	query := `
		INSERT INTO idempotency_keys (key, user_id, request_path, request_method, response_status, response_body, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (key) DO NOTHING
	`

	_, err := s.db.Exec(ctx, query, key, userID, path, method, status, body, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

// GetIdempotentResponse retrieves a cached response for a given idempotency key.
// Returns nil if no cached response is found or if the response has expired.
func (s *service) GetIdempotentResponse(ctx context.Context, key string, userID int) (*IdempotentResponse, error) {
	query := `
		SELECT response_status, response_body
		FROM idempotency_keys
		WHERE key = $1 AND user_id = $2 AND expires_at > NOW()
	`

	var response IdempotentResponse
	response.Key = key

	err := s.db.QueryRow(ctx, query, key, userID).Scan(&response.StatusCode, &response.Body)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &response, nil
}

// DeleteExpiredIdempotencyKeys removes all expired idempotency keys from the database.
// Returns the number of keys deleted.
func (s *service) DeleteExpiredIdempotencyKeys(ctx context.Context) (int64, error) {
	query := `DELETE FROM idempotency_keys WHERE expires_at <= NOW()`

	result, err := s.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}
