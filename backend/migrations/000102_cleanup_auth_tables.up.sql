-- Migration: Cleanup legacy auth tables after migration to Auth Service microservice
-- 
-- These tables have been migrated to the Auth Service and are no longer needed in the main database
-- Auth Service now handles all authentication, sessions, OAuth, and refresh tokens

-- Drop sessions table (moved to Auth Service)
DROP TABLE IF EXISTS sessions CASCADE;

-- Drop oauth_states table (moved to Auth Service)  
DROP TABLE IF EXISTS oauth_states CASCADE;

-- Drop refresh_tokens table (moved to Auth Service)
DROP TABLE IF EXISTS refresh_tokens CASCADE;

-- Drop auth_audit_logs table (moved to Auth Service)
DROP TABLE IF EXISTS auth_audit_logs CASCADE;

-- Clean up any remaining auth-related indexes that might exist
DROP INDEX IF EXISTS idx_sessions_token;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_oauth_states_state;
DROP INDEX IF EXISTS idx_oauth_states_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_auth_audit_logs_user_id;
DROP INDEX IF EXISTS idx_auth_audit_logs_created_at;