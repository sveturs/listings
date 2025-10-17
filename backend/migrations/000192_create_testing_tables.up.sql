-- Create testing module tables
-- Migration: 000192_create_testing_tables

-- Test runs table
CREATE TABLE test_runs (
    id BIGSERIAL PRIMARY KEY,
    run_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    test_suite VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled')),
    started_by_user_id INT NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    duration_ms INT,
    total_tests INT DEFAULT 0,
    passed_tests INT DEFAULT 0,
    failed_tests INT DEFAULT 0,
    skipped_tests INT DEFAULT 0,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for test_runs
CREATE UNIQUE INDEX idx_test_runs_uuid ON test_runs(run_uuid);
CREATE INDEX idx_test_runs_suite ON test_runs(test_suite);
CREATE INDEX idx_test_runs_status ON test_runs(status);
CREATE INDEX idx_test_runs_started_at ON test_runs(started_at DESC);
CREATE INDEX idx_test_runs_user_id ON test_runs(started_by_user_id);

-- Test results table
CREATE TABLE test_results (
    id BIGSERIAL PRIMARY KEY,
    test_run_id BIGINT NOT NULL REFERENCES test_runs(id) ON DELETE CASCADE,
    test_name VARCHAR(255) NOT NULL,
    test_suite VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('passed', 'failed', 'skipped')),
    duration_ms INT NOT NULL,
    error_msg TEXT,
    stack_trace TEXT,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for test_results
CREATE INDEX idx_test_results_run_id ON test_results(test_run_id);
CREATE INDEX idx_test_results_status ON test_results(status);
CREATE INDEX idx_test_results_suite ON test_results(test_suite);

-- Test logs table
CREATE TABLE test_logs (
    id BIGSERIAL PRIMARY KEY,
    test_run_id BIGINT NOT NULL REFERENCES test_runs(id) ON DELETE CASCADE,
    level VARCHAR(20) NOT NULL CHECK (level IN ('debug', 'info', 'warn', 'error')),
    message TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for test_logs
CREATE INDEX idx_test_logs_run_id ON test_logs(test_run_id);
CREATE INDEX idx_test_logs_level ON test_logs(level);
CREATE INDEX idx_test_logs_timestamp ON test_logs(timestamp DESC);

-- Add comments
COMMENT ON TABLE test_runs IS 'Stores test suite execution runs';
COMMENT ON TABLE test_results IS 'Stores individual test execution results';
COMMENT ON TABLE test_logs IS 'Stores test execution logs';

COMMENT ON COLUMN test_runs.run_uuid IS 'Unique identifier for test run';
COMMENT ON COLUMN test_runs.test_suite IS 'Name of test suite (e.g., api-endpoints)';
COMMENT ON COLUMN test_runs.status IS 'Current execution status';
COMMENT ON COLUMN test_runs.metadata IS 'Additional metadata for test run (JSON)';

COMMENT ON COLUMN test_results.test_name IS 'Name of individual test';
COMMENT ON COLUMN test_results.error_msg IS 'Error message if test failed';
COMMENT ON COLUMN test_results.stack_trace IS 'Stack trace if test failed';

COMMENT ON COLUMN test_logs.level IS 'Log level (debug, info, warn, error)';
COMMENT ON COLUMN test_logs.message IS 'Log message';
