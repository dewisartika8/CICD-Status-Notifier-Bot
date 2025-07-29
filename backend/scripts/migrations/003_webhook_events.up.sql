-- Create webhook_events table for storing GitHub webhook events
-- This table tracks all webhook events received from GitHub for auditing and reprocessing

CREATE TABLE IF NOT EXISTS webhook_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    signature VARCHAR(255) NOT NULL,
    delivery_id VARCHAR(255),
    processed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_webhook_events_project_id 
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Create indexes for performance (PostgreSQL syntax)
CREATE INDEX IF NOT EXISTS idx_webhook_events_project_id ON webhook_events(project_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_event_type ON webhook_events(event_type);
CREATE INDEX IF NOT EXISTS idx_webhook_events_delivery_id ON webhook_events(delivery_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_processed_at ON webhook_events(processed_at);
CREATE INDEX IF NOT EXISTS idx_webhook_events_created_at ON webhook_events(created_at);

-- Add unique constraint for delivery_id to prevent duplicate processing
-- Note: delivery_id can be NULL for testing or when GitHub doesn't provide it
CREATE UNIQUE INDEX IF NOT EXISTS uk_webhook_events_delivery_id ON webhook_events(delivery_id) WHERE delivery_id IS NOT NULL;

-- Comments for documentation
COMMENT ON TABLE webhook_events IS 'Stores GitHub webhook events for auditing and reprocessing';
COMMENT ON COLUMN webhook_events.id IS 'Unique identifier for the webhook event (UUID)';
COMMENT ON COLUMN webhook_events.project_id IS 'Reference to the project that received the webhook (UUID foreign key)';
COMMENT ON COLUMN webhook_events.event_type IS 'Type of GitHub event (workflow_run, push, pull_request)';
COMMENT ON COLUMN webhook_events.payload IS 'Complete GitHub webhook payload as JSON';
COMMENT ON COLUMN webhook_events.signature IS 'GitHub webhook signature for verification';
COMMENT ON COLUMN webhook_events.delivery_id IS 'GitHub delivery ID for idempotency';
COMMENT ON COLUMN webhook_events.processed_at IS 'Timestamp when the webhook was successfully processed';
COMMENT ON COLUMN webhook_events.created_at IS 'Timestamp when the webhook was received';
