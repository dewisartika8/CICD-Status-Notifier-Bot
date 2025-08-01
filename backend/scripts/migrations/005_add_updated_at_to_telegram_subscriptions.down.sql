-- Migration 005 down: Remove updated_at column from telegram_subscriptions table
-- This migration removes the updated_at column and trigger

-- Drop trigger for updated_at column
DROP TRIGGER IF EXISTS update_telegram_subscriptions_updated_at ON telegram_subscriptions;

-- Drop updated_at column
ALTER TABLE telegram_subscriptions DROP COLUMN IF EXISTS updated_at;
