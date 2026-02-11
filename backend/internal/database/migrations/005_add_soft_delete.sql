-- Migration: Add soft delete support and updated_at columns
-- This migration adds deleted_at and updated_at columns to support soft deletes and track updates

-- Add deleted_at and updated_at to learning_paths
ALTER TABLE learning_paths
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Add deleted_at and updated_at to modules
ALTER TABLE modules
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Add deleted_at and updated_at to lessons
ALTER TABLE lessons
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Add deleted_at and updated_at to projects
ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Create indexes for soft delete queries (improves performance)
CREATE INDEX IF NOT EXISTS idx_learning_paths_deleted_at ON learning_paths(deleted_at);
CREATE INDEX IF NOT EXISTS idx_modules_deleted_at ON modules(deleted_at);
CREATE INDEX IF NOT EXISTS idx_lessons_deleted_at ON lessons(deleted_at);
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects(deleted_at);

-- Update trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for automatic updated_at updates
DROP TRIGGER IF EXISTS update_learning_paths_updated_at ON learning_paths;
CREATE TRIGGER update_learning_paths_updated_at
    BEFORE UPDATE ON learning_paths
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_modules_updated_at ON modules;
CREATE TRIGGER update_modules_updated_at
    BEFORE UPDATE ON modules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_lessons_updated_at ON lessons;
CREATE TRIGGER update_lessons_updated_at
    BEFORE UPDATE ON lessons
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
CREATE TRIGGER update_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Note: Soft delete behavior
-- When a learning_path is soft-deleted (deleted_at IS NOT NULL):
--   - The application layer handles cascading soft deletes to modules and lessons
--   - Foreign key constraints remain in place (ON DELETE CASCADE)
--   - Hard deletes (if ever needed) will cascade automatically
--   - Soft deletes provide an audit trail and allow recovery
