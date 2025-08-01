-- Migration: 004_add_missing_build_event_columns.up.sql
-- Add missing columns to build_events table to match domain model

-- Add updated_at column
ALTER TABLE build_events 
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Add deleted_at column for soft delete functionality
ALTER TABLE build_events 
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

-- Create index for deleted_at column (required by GORM)
CREATE INDEX IF NOT EXISTS idx_build_events_deleted_at ON build_events(deleted_at);

-- Add trigger for updated_at column
CREATE TRIGGER update_build_events_updated_at 
    BEFORE UPDATE ON build_events 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Update existing records to have updated_at = created_at
UPDATE build_events SET updated_at = created_at WHERE updated_at IS NULL;
