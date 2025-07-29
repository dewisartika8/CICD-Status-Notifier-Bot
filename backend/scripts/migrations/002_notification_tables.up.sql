-- Migration 002: Additional notification system enhancements
-- This migration adds additional features to the existing notification system

-- Add additional columns to notification_logs if needed (example)
-- ALTER TABLE notification_logs ADD COLUMN IF NOT EXISTS retry_count INTEGER DEFAULT 0;
-- ALTER TABLE notification_logs ADD COLUMN IF NOT EXISTS channel VARCHAR(50);

-- Add any additional indexes or constraints here if needed
-- CREATE INDEX IF NOT EXISTS idx_notification_logs_retry_count ON notification_logs(retry_count);

-- Placeholder migration - all tables already exist in 001_initial_schema
-- This file can be used for future notification system enhancements
