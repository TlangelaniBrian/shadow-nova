-- Add deleted_at column to users table for soft deletes
-- This migration can be run on existing databases

-- Add the column if it doesn't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL;

-- Add index for performance on deleted_at queries
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Update existing queries to exclude soft-deleted users
-- Note: The application code already handles this, but we add the index for performance
