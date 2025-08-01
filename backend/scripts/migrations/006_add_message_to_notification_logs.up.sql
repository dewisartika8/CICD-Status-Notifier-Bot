-- Migration 006: Add message column to notification_logs
-- This migration adds the message column to store actual notification content

-- Add message column to notification_logs table
ALTER TABLE notification_logs ADD COLUMN IF NOT EXISTS message TEXT;

-- Create index for message column for potential search functionality
CREATE INDEX IF NOT EXISTS idx_notification_logs_message ON notification_logs USING gin(to_tsvector('english', message));

-- Backfill existing records with a generic message
-- Note: This is for existing records only, new records will have proper messages
UPDATE notification_logs 
SET message = 'Legacy notification' 
WHERE message IS NULL;
