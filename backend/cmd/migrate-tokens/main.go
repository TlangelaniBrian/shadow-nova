package main

import (
	"context"
	"encoding/base64"
	"log/slog"
	"os"
	"time"

	"shadow-nova/backend/internal/crypto"
	"shadow-nova/backend/internal/logging"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type tokenRecord struct {
	id           int
	userID       int
	accessToken  string
	refreshToken *string
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using environment variables")
	}

	// Initialize structured logging
	logging.Init()

	// Initialize encryption
	if err := crypto.Init(); err != nil {
		logging.Error("failed to initialize encryption", err)
		os.Exit(1)
	}

	// Connect to database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logging.Error("DATABASE_URL environment variable not set", nil)
		os.Exit(1)
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		logging.Error("unable to parse database URL", err)
		os.Exit(1)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logging.Error("unable to create connection pool", err)
		os.Exit(1)
	}
	defer pool.Close()

	ctx := context.Background()

	// Check if we can connect
	if err := pool.Ping(ctx); err != nil {
		logging.Error("unable to ping database", err)
		os.Exit(1)
	}

	logging.Info("connected to database successfully")

	// Get all github integrations
	query := `SELECT id, user_id, access_token, refresh_token FROM github_integrations`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		logging.Error("failed to query github integrations", err)
		os.Exit(1)
	}
	defer rows.Close()

	var records []tokenRecord
	for rows.Next() {
		var rec tokenRecord
		if err := rows.Scan(&rec.id, &rec.userID, &rec.accessToken, &rec.refreshToken); err != nil {
			logging.Error("failed to scan record", err)
			os.Exit(1)
		}
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		logging.Error("error iterating records", err)
		os.Exit(1)
	}

	logging.Info("found github integrations to process", "count", len(records))

	if len(records) == 0 {
		logging.Info("no records to migrate")
		return
	}

	// Process each record
	encrypted := 0
	alreadyEncrypted := 0
	failed := 0

	for _, rec := range records {
		// Check if token is already encrypted (base64 encoded and longer than typical OAuth token)
		if isLikelyEncrypted(rec.accessToken) {
			alreadyEncrypted++
			logging.Info("record appears already encrypted, skipping", "record_id", rec.id, "user_id", rec.userID)
			continue
		}

		// Encrypt access token
		encryptedAccess, err := crypto.Encrypt(rec.accessToken)
		if err != nil {
			logging.Error("failed to encrypt access token", err, "record_id", rec.id)
			failed++
			continue
		}

		// Encrypt refresh token if present
		var encryptedRefresh *string
		if rec.refreshToken != nil && *rec.refreshToken != "" {
			if isLikelyEncrypted(*rec.refreshToken) {
				encryptedRefresh = rec.refreshToken
			} else {
				encrypted_token, err := crypto.Encrypt(*rec.refreshToken)
				if err != nil {
					logging.Error("failed to encrypt refresh token", err, "record_id", rec.id)
					failed++
					continue
				}
				encryptedRefresh = &encrypted_token
			}
		}

		// Update the record
		updateQuery := `UPDATE github_integrations SET access_token = $1, refresh_token = $2, updated_at = $3 WHERE id = $4`
		_, err = pool.Exec(ctx, updateQuery, encryptedAccess, encryptedRefresh, time.Now(), rec.id)
		if err != nil {
			logging.Error("failed to update record", err, "record_id", rec.id)
			failed++
			continue
		}

		encrypted++
		logging.Info("successfully encrypted tokens", "record_id", rec.id, "user_id", rec.userID)
	}

	logging.Info("migration summary",
		"total_records", len(records),
		"encrypted", encrypted,
		"already_encrypted", alreadyEncrypted,
		"failed", failed,
	)

	if failed > 0 {
		logging.Error("migration completed with errors", nil)
		os.Exit(1)
	}

	logging.Info("migration completed successfully")
}

// isLikelyEncrypted checks if a token appears to be already encrypted
// Encrypted tokens are base64 encoded and contain the nonce + ciphertext
// They should be longer than typical OAuth tokens and have specific characteristics
func isLikelyEncrypted(token string) bool {
	// If it's too short, it's definitely not encrypted
	if len(token) < 50 {
		return false
	}

	// Try to base64 decode it
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		// Not valid base64, so probably not encrypted by our system
		return false
	}

	// Our encrypted tokens should be at least nonce size (12) + some ciphertext + tag (16)
	// So minimum around 40-50 bytes decoded
	if len(decoded) < 40 {
		return false
	}

	// If we can successfully decrypt it, it's already encrypted
	_, err = crypto.Decrypt(token)
	return err == nil
}
