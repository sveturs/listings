-- Создание таблицы import_jobs для отслеживания импорта товаров
CREATE TABLE import_jobs (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL,
    file_name VARCHAR(255),
    file_type VARCHAR(10) NOT NULL CHECK (file_type IN ('xml', 'csv', 'zip')),
    file_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    total_records INTEGER DEFAULT 0,
    processed_records INTEGER DEFAULT 0,
    successful_records INTEGER DEFAULT 0,
    failed_records INTEGER DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_import_jobs_storefront ON import_jobs(storefront_id);
CREATE INDEX idx_import_jobs_status ON import_jobs(status);
CREATE INDEX idx_import_jobs_created ON import_jobs(created_at DESC);
CREATE INDEX idx_import_jobs_user ON import_jobs(user_id);

-- Таблица для детальных ошибок импорта
CREATE TABLE import_errors (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES import_jobs(id) ON DELETE CASCADE,
    line_number INTEGER,
    field_name VARCHAR(100),
    error_message TEXT NOT NULL,
    raw_data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_import_errors_job ON import_errors(job_id);

-- Trigger для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_import_jobs_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_import_jobs_updated_at
    BEFORE UPDATE ON import_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_import_jobs_updated_at();
