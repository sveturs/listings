-- Добавляем уникальный индекс на title_hash для поддержки ON CONFLICT
CREATE UNIQUE INDEX idx_ai_decisions_unique_title_hash ON ai_category_decisions(title_hash);