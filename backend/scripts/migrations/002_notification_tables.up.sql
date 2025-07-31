-- Migration 002: Additional notification system enhancements
-- This migration adds notification templates and retry configurations

-- Create notification templates table
CREATE TABLE IF NOT EXISTS notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_type VARCHAR(50) NOT NULL,
    channel VARCHAR(50) NOT NULL,
    subject TEXT NOT NULL,
    body_template TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(template_type, channel)
);

-- Create retry configurations table
CREATE TABLE IF NOT EXISTS retry_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,
    channel VARCHAR(50) NOT NULL,
    max_retries INTEGER NOT NULL DEFAULT 3,
    base_delay_seconds INTEGER NOT NULL DEFAULT 5,
    max_delay_seconds INTEGER NOT NULL DEFAULT 300,
    backoff_factor DECIMAL(3,2) NOT NULL DEFAULT 2.0,
    retryable_errors TEXT[], -- PostgreSQL array of text
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(project_id, channel),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Add additional columns to notification_logs if needed
ALTER TABLE notification_logs ADD COLUMN IF NOT EXISTS retry_count INTEGER DEFAULT 0;
ALTER TABLE notification_logs ADD COLUMN IF NOT EXISTS channel VARCHAR(50);
ALTER TABLE notification_logs ADD COLUMN IF NOT EXISTS template_id UUID;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_notification_templates_type ON notification_templates(template_type);
CREATE INDEX IF NOT EXISTS idx_notification_templates_channel ON notification_templates(channel);
CREATE INDEX IF NOT EXISTS idx_notification_templates_active ON notification_templates(is_active);

CREATE INDEX IF NOT EXISTS idx_retry_configurations_project ON retry_configurations(project_id);
CREATE INDEX IF NOT EXISTS idx_retry_configurations_channel ON retry_configurations(channel);
CREATE INDEX IF NOT EXISTS idx_retry_configurations_active ON retry_configurations(is_active);

CREATE INDEX IF NOT EXISTS idx_notification_logs_retry_count ON notification_logs(retry_count);
CREATE INDEX IF NOT EXISTS idx_notification_logs_channel ON notification_logs(channel);
CREATE INDEX IF NOT EXISTS idx_notification_logs_template ON notification_logs(template_id);

-- Add foreign key for template reference
ALTER TABLE notification_logs ADD CONSTRAINT fk_notification_logs_template 
    FOREIGN KEY (template_id) REFERENCES notification_templates(id) ON DELETE SET NULL;
