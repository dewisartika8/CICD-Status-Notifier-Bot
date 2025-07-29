-- Rollback migration 002: Additional notification system enhancements
-- This rollback removes any enhancements added in migration 002

-- Remove any additional columns that were added (example)
-- ALTER TABLE notification_logs DROP COLUMN IF EXISTS retry_count;
-- ALTER TABLE notification_logs DROP COLUMN IF EXISTS channel;

-- Remove any additional indexes that were added
-- DROP INDEX IF EXISTS idx_notification_logs_retry_count;

-- Placeholder rollback - no changes to undo in this migration
