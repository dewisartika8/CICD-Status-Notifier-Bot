-- Rollback migration 002: Additional notification system enhancements

-- Remove foreign key constraints
ALTER TABLE notification_logs DROP CONSTRAINT IF EXISTS fk_notification_logs_template;

-- Remove indexes
DROP INDEX IF EXISTS idx_notification_logs_template;
DROP INDEX IF EXISTS idx_notification_logs_channel;
DROP INDEX IF EXISTS idx_notification_logs_retry_count;

DROP INDEX IF EXISTS idx_retry_configurations_active;
DROP INDEX IF EXISTS idx_retry_configurations_channel;
DROP INDEX IF EXISTS idx_retry_configurations_project;

DROP INDEX IF EXISTS idx_notification_templates_active;
DROP INDEX IF EXISTS idx_notification_templates_channel;
DROP INDEX IF EXISTS idx_notification_templates_type;

-- Remove added columns from notification_logs
ALTER TABLE notification_logs DROP COLUMN IF EXISTS template_id;
ALTER TABLE notification_logs DROP COLUMN IF EXISTS channel;
ALTER TABLE notification_logs DROP COLUMN IF EXISTS retry_count;

-- Drop new tables
DROP TABLE IF EXISTS retry_configurations;
DROP TABLE IF EXISTS notification_templates;
