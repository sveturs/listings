-- Drop webhook audit log table
DROP TABLE IF EXISTS webhook_audit_log;

-- Drop failed webhooks table
DROP TABLE IF EXISTS failed_webhooks;

-- Drop the update trigger function
DROP FUNCTION IF EXISTS update_failed_webhooks_updated_at();