package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"shadow-nova/backend/internal/crypto"

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
		log.Println("No .env file found, using environment variables")
	}

	// Initialize encryption
	if err := crypto.Init(); err != nil {
		log.Fatalf("Failed to initialize encryption: %v", err)
	}

	// Connect to database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	// Check if we can connect
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Get all github integrations
	query := `SELECT id, user_id, access_token, refresh_token FROM github_integrations`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("Failed to query github integrations: %v", err)
	}
	defer rows.Close()

	var records []tokenRecord
	for rows.Next() {
		var rec tokenRecord
		if err := rows.Scan(&rec.id, &rec.userID, &rec.accessToken, &rec.refreshToken); err != nil {
			log.Fatalf("Failed to scan record: %v", err)
		}
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating records: %v", err)
	}

	log.Printf("Found %d GitHub integrations to process\n", len(records))

	if len(records) == 0 {
		log.Println("No records to migrate")
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
			log.Printf("Record %d (user %d) appears already encrypted, skipping\n", rec.id, rec.userID)
			continue
		}

		// Encrypt access token
		encryptedAccess, err := crypto.Encrypt(rec.accessToken)
		if err != nil {
			log.Printf("Failed to encrypt access token for record %d: %v\n", rec.id, err)
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
					log.Printf("Failed to encrypt refresh token for record %d: %v\n", rec.id, err)
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
			log.Printf("Failed to update record %d: %v\n", rec.id, err)
			failed++
			continue
		}

		encrypted++
		log.Printf("Successfully encrypted tokens for record %d (user %d)\n", rec.id, rec.userID)
	}

	log.Println("\nMigration Summary:")
	log.Printf("Total records: %d\n", len(records))
	log.Printf("Encrypted: %d\n", encrypted)
	log.Printf("Already encrypted: %d\n", alreadyEncrypted)
	log.Printf("Failed: %d\n", failed)

	if failed > 0 {
		log.Fatal("Migration completed with errors")
	}

	log.Println("Migration completed successfully!")
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
