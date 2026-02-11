package database

import (
	"context"
	"fmt"
	"time"

	"shadow-nova/backend/internal/errors"
	"shadow-nova/backend/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *service) CreateUser(ctx context.Context, user *models.User) error {
	return s.WithTx(ctx, func(tx pgx.Tx) error {
		// Default role to 'user' if not set
		if user.Role == "" {
			user.Role = "user"
		}

		query := `
			INSERT INTO users (email, username, password_hash, user_role)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at
		`

		err := tx.QueryRow(ctx, query, user.Email, user.Username, user.PasswordHash, user.Role).Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Future: Initialize user preferences, welcome email queue, etc.
		// These would all be done within the same transaction

		return nil
	})
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, user_role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound(fmt.Sprintf("user with email %s not found", email))
		}
		return nil, errors.DatabaseError(err, "failed to get user by email")
	}

	return user, nil
}

// BlacklistToken adds a token to the blacklist
func (s *service) BlacklistToken(ctx context.Context, jti string, userID int, expiresAt time.Time, reason string) error {
	query := `INSERT INTO token_blacklist (jti, user_id, expires_at, reason) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(ctx, query, jti, userID, expiresAt, reason)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	return nil
}

// IsTokenBlacklisted checks if a token is in the blacklist
func (s *service) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE jti = $1 AND expires_at > NOW())`
	err := s.db.QueryRow(ctx, query, jti).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	return exists, nil
}

// BlacklistAllUserTokens blacklists all tokens for a user
func (s *service) BlacklistAllUserTokens(ctx context.Context, userID int, reason string) error {
	// This requires tracking all issued tokens - implement later
	// For now, just log that user requested logout from all devices
	return nil
}

// DeleteExpiredBlacklistedTokens removes expired tokens from the blacklist
func (s *service) DeleteExpiredBlacklistedTokens(ctx context.Context) (int64, error) {
	query := `DELETE FROM token_blacklist WHERE expires_at <= NOW()`
	result, err := s.db.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired blacklisted tokens: %w", err)
	}
	return result.RowsAffected(), nil
}
