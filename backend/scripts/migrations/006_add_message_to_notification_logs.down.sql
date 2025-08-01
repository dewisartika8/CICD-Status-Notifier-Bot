-- Migration 006: Rollback - Remove message column from notification_logs

-- Remove the message column from notification_logs table
ALTER TABLE notification_logs DROP COLUMN IF EXISTS message;

-- Drop the index (will be automatically dropped with the column, but being explicit)
DROP INDEX IF EXISTS idx_notification_logs_message;
