-- Migration: Add token blacklist table for JWT revocation
-- Date: 2026-02-11

-- Token blacklist for revoked JWTs
CREATE TABLE IF NOT EXISTS token_blacklist (
    jti VARCHAR(36) PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    blacklisted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reason VARCHAR(100)
);

-- Indexes for efficient token blacklist lookups
CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires ON token_blacklist(expires_at);
CREATE INDEX IF NOT EXISTS idx_token_blacklist_user ON token_blacklist(user_id);
