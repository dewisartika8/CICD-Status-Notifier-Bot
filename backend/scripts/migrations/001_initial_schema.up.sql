-- Migration: 001_initial_schema.up.sql
-- Create initial database schema for CI/CD Status Notifier Bot

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    repository_url VARCHAR(500) NOT NULL,
    webhook_secret VARCHAR(255),
    telegram_chat_id BIGINT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Build status events
CREATE TABLE build_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- build_started, build_success, build_failed, test_started, test_passed, test_failed, deployment_started, deployment_success, deployment_failed
    status VARCHAR(20) NOT NULL, -- success, failed, pending, running
    branch VARCHAR(255) NOT NULL,
    commit_sha VARCHAR(40),
    commit_message TEXT,
    author_name VARCHAR(255),
    author_email VARCHAR(255),
    build_url VARCHAR(500),
    duration_seconds INTEGER,
    webhook_payload JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Telegram subscriptions
CREATE TABLE telegram_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    chat_id BIGINT NOT NULL,
    user_id BIGINT,
    username VARCHAR(255),
    event_types TEXT[], -- Array of event types to subscribe to
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Notification logs
CREATE TABLE notification_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    build_event_id UUID NOT NULL REFERENCES build_events(id) ON DELETE CASCADE,
    chat_id BIGINT NOT NULL,
    message_id INTEGER,
    status VARCHAR(20) NOT NULL, -- sent, failed, pending
    error_message TEXT,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance optimization
CREATE INDEX idx_build_events_project_id ON build_events(project_id);
CREATE INDEX idx_build_events_created_at ON build_events(created_at DESC);
CREATE INDEX idx_build_events_status ON build_events(status);
CREATE INDEX idx_build_events_event_type ON build_events(event_type);
CREATE INDEX idx_build_events_branch ON build_events(branch);

CREATE INDEX idx_telegram_subscriptions_project_id ON telegram_subscriptions(project_id);
CREATE INDEX idx_telegram_subscriptions_chat_id ON telegram_subscriptions(chat_id);
CREATE INDEX idx_telegram_subscriptions_is_active ON telegram_subscriptions(is_active);

CREATE INDEX idx_notification_logs_build_event_id ON notification_logs(build_event_id);
CREATE INDEX idx_notification_logs_chat_id ON notification_logs(chat_id);
CREATE INDEX idx_notification_logs_status ON notification_logs(status);
CREATE INDEX idx_notification_logs_created_at ON notification_logs(created_at DESC);

-- Unique constraint for telegram subscriptions
ALTER TABLE telegram_subscriptions 
ADD CONSTRAINT unique_project_chat UNIQUE (project_id, chat_id);

-- Add triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_projects_updated_at 
    BEFORE UPDATE ON projects 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
