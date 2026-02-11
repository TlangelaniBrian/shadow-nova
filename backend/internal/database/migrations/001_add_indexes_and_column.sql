-- Migration: Add performance indexes and github_username column
-- Created: 2026-02-11

-- Add missing column
ALTER TABLE users ADD COLUMN IF NOT EXISTS github_username VARCHAR(100);

-- Add performance indexes
CREATE INDEX IF NOT EXISTS idx_modules_path_id ON modules(path_id);
CREATE INDEX IF NOT EXISTS idx_lessons_module_id ON lessons(module_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_user_completed ON user_progress(user_id, completed);
CREATE INDEX IF NOT EXISTS idx_content_items_processed ON content_items(processed_by_ai) WHERE processed_by_ai = FALSE;
CREATE INDEX IF NOT EXISTS idx_content_items_source_id ON content_items(source_id);
CREATE INDEX IF NOT EXISTS idx_project_submissions_project_id ON project_submissions(project_id);
CREATE INDEX IF NOT EXISTS idx_project_submissions_status ON project_submissions(status, submitted_at DESC);
CREATE INDEX IF NOT EXISTS idx_github_integrations_github_user_id ON github_integrations(github_user_id);
