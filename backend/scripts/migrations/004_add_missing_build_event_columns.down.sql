-- Migration: 004_add_missing_build_event_columns.down.sql
-- Rollback: Remove the added columns from build_events table

-- Drop trigger
DROP TRIGGER IF EXISTS update_build_events_updated_at ON build_events;

-- Drop index
DROP INDEX IF EXISTS idx_build_events_deleted_at;

-- Remove columns
ALTER TABLE build_events DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE build_events DROP COLUMN IF EXISTS updated_at;
