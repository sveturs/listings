-- Drop admin_users table as we now use Auth Service for role management

-- First, drop foreign key constraints
ALTER TABLE search_optimization_sessions DROP CONSTRAINT IF EXISTS search_optimization_sessions_created_by_fkey;
ALTER TABLE search_weights DROP CONSTRAINT IF EXISTS search_weights_created_by_fkey;
ALTER TABLE search_weights_history DROP CONSTRAINT IF EXISTS search_weights_history_changed_by_fkey;
ALTER TABLE search_weights DROP CONSTRAINT IF EXISTS search_weights_updated_by_fkey;

-- Now drop the table
DROP TABLE IF EXISTS admin_users;
