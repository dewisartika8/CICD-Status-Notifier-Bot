-- Migration 005: Add updated_at column to telegram_subscriptions table
-- This migration adds the missing updated_at column and trigger

-- Add updated_at column to telegram_subscriptions
ALTER TABLE telegram_subscriptions 
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Update existing records to have updated_at = created_at
UPDATE telegram_subscriptions SET updated_at = created_at WHERE updated_at IS NULL;

-- Add trigger for updated_at column on telegram_subscriptions
CREATE TRIGGER update_telegram_subscriptions_updated_at 
    BEFORE UPDATE ON telegram_subscriptions 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
