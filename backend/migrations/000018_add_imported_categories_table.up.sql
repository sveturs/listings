-- /data/hostel-booking-system/backend/migrations/000018_add_imported_categories_table.up.sql
-- Создание таблицы для отслеживания категорий, импортированных из внешних источников

CREATE TABLE IF NOT EXISTS imported_categories (
    id SERIAL PRIMARY KEY,
    source_id INT NOT NULL REFERENCES import_sources(id) ON DELETE CASCADE,
    source_category VARCHAR(255) NOT NULL,
    category_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_id, source_category)
);

-- Создаем индекс для быстрого поиска по source_id
CREATE INDEX IF NOT EXISTS idx_imported_categories_source_id ON imported_categories(source_id);