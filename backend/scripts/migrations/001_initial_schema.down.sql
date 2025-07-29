-- Migration: 001_initial_schema.down.sql
-- Rollback initial database schema

-- Drop triggers
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_notification_logs_created_at;
DROP INDEX IF EXISTS idx_notification_logs_status;
DROP INDEX IF EXISTS idx_notification_logs_chat_id;
DROP INDEX IF EXISTS idx_notification_logs_build_event_id;

DROP INDEX IF EXISTS idx_telegram_subscriptions_is_active;
DROP INDEX IF EXISTS idx_telegram_subscriptions_chat_id;
DROP INDEX IF EXISTS idx_telegram_subscriptions_project_id;

DROP INDEX IF EXISTS idx_build_events_branch;
DROP INDEX IF EXISTS idx_build_events_event_type;
DROP INDEX IF EXISTS idx_build_events_status;
DROP INDEX IF EXISTS idx_build_events_created_at;
DROP INDEX IF EXISTS idx_build_events_project_id;

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS notification_logs;
DROP TABLE IF EXISTS telegram_subscriptions;
DROP TABLE IF EXISTS build_events;
DROP TABLE IF EXISTS projects;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp";
