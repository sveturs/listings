-- Rollback testing module tables
-- Migration: 000192_create_testing_tables

-- Drop tables in reverse order (children first due to foreign keys)
DROP TABLE IF EXISTS test_logs;
DROP TABLE IF EXISTS test_results;
DROP TABLE IF EXISTS test_runs;
