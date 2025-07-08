-- Создание таблицы для конфигурации весов поиска
-- Эта таблица содержит веса для различных полей поиска, которые используются
-- для расчета релевантности результатов поиска в OpenSearch и PostgreSQL

CREATE TABLE IF NOT EXISTS search_weights (
    id BIGSERIAL PRIMARY KEY,
    
    -- Название поля поиска (title, description, category, price, location, etc.)
    field_name VARCHAR(100) NOT NULL,
    
    -- Вес поля (0.0 - 1.0, где 1.0 - максимальная важность)
    weight FLOAT NOT NULL CHECK (weight >= 0.0 AND weight <= 1.0),
    
    -- Тип поиска: 'fulltext', 'fuzzy', 'exact'
    search_type VARCHAR(20) NOT NULL DEFAULT 'fulltext' 
        CHECK (search_type IN ('fulltext', 'fuzzy', 'exact')),
    
    -- Тип элемента: 'marketplace', 'storefront', 'global' (для всех типов)
    item_type VARCHAR(20) NOT NULL DEFAULT 'global'
        CHECK (item_type IN ('marketplace', 'storefront', 'global')),
    
    -- ID категории (NULL для глобальных весов)
    category_id INTEGER REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    
    -- Описание назначения поля
    description TEXT,
    
    -- Флаг активности веса
    is_active BOOLEAN DEFAULT true,
    
    -- Версия конфигурации весов для отслеживания изменений
    version INTEGER DEFAULT 1,
    
    -- Временные метки
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Кто создал/обновил запись
    created_by INTEGER REFERENCES admin_users(id) ON DELETE SET NULL,
    updated_by INTEGER REFERENCES admin_users(id) ON DELETE SET NULL
);

-- Уникальные индексы для предотвращения дублирования
CREATE UNIQUE INDEX idx_search_weights_unique_global 
    ON search_weights(field_name, item_type, search_type) 
    WHERE category_id IS NULL;

CREATE UNIQUE INDEX idx_search_weights_unique_category 
    ON search_weights(field_name, item_type, search_type, category_id) 
    WHERE category_id IS NOT NULL;

-- Индексы для быстрого поиска
CREATE INDEX idx_search_weights_field_name ON search_weights(field_name);
CREATE INDEX idx_search_weights_item_type ON search_weights(item_type);
CREATE INDEX idx_search_weights_category_id ON search_weights(category_id) WHERE category_id IS NOT NULL;
CREATE INDEX idx_search_weights_is_active ON search_weights(is_active) WHERE is_active = true;
CREATE INDEX idx_search_weights_version ON search_weights(version);

-- Таблица для истории изменений весов (для аудита и rollback)
CREATE TABLE IF NOT EXISTS search_weights_history (
    id BIGSERIAL PRIMARY KEY,
    
    -- Ссылка на исходную запись весов
    weight_id BIGINT NOT NULL REFERENCES search_weights(id) ON DELETE CASCADE,
    
    -- Старое значение веса
    old_weight FLOAT NOT NULL,
    
    -- Новое значение веса
    new_weight FLOAT NOT NULL,
    
    -- Причина изменения: 'manual', 'optimization', 'rollback'
    change_reason VARCHAR(50) NOT NULL DEFAULT 'manual'
        CHECK (change_reason IN ('manual', 'optimization', 'rollback', 'initialization')),
    
    -- Дополнительная информация об изменении
    change_metadata JSONB DEFAULT '{}',
    
    -- Кто внес изменение
    changed_by INTEGER REFERENCES admin_users(id) ON DELETE SET NULL,
    
    -- Время изменения
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для истории изменений
CREATE INDEX idx_search_weights_history_weight_id ON search_weights_history(weight_id);
CREATE INDEX idx_search_weights_history_changed_at ON search_weights_history(changed_at);
CREATE INDEX idx_search_weights_history_changed_by ON search_weights_history(changed_by);
CREATE INDEX idx_search_weights_history_reason ON search_weights_history(change_reason);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_search_weights_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_search_weights_updated_at
    BEFORE UPDATE ON search_weights
    FOR EACH ROW
    EXECUTE FUNCTION update_search_weights_updated_at();

-- Триггер для автоматического логирования изменений весов
CREATE OR REPLACE FUNCTION log_search_weight_changes()
RETURNS TRIGGER AS $$
BEGIN
    -- Логируем только изменения веса
    IF OLD.weight <> NEW.weight THEN
        INSERT INTO search_weights_history (
            weight_id, 
            old_weight, 
            new_weight, 
            change_reason,
            changed_by
        ) VALUES (
            NEW.id, 
            OLD.weight, 
            NEW.weight,
            'manual',  -- По умолчанию, может быть переопределено в коде
            NEW.updated_by
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_search_weight_changes
    AFTER UPDATE ON search_weights
    FOR EACH ROW
    EXECUTE FUNCTION log_search_weight_changes();

-- Вставка начальных весов для поиска
INSERT INTO search_weights (field_name, weight, search_type, item_type, description, version) VALUES
-- Глобальные веса для всех типов контента
('title', 0.9, 'fulltext', 'global', 'Название/заголовок - самое важное поле для поиска', 1),
('title', 0.85, 'fuzzy', 'global', 'Название/заголовок для нечеткого поиска', 1),
('description', 0.7, 'fulltext', 'global', 'Описание - второе по важности поле', 1),
('description', 0.6, 'fuzzy', 'global', 'Описание для нечеткого поиска', 1),
('category', 0.8, 'exact', 'global', 'Категория товара/услуги', 1),
('location', 0.75, 'exact', 'global', 'Местоположение', 1),
('price', 0.5, 'exact', 'global', 'Цена (важна для фильтрации)', 1),
('brand', 0.6, 'exact', 'global', 'Бренд/производитель', 1),
('tags', 0.4, 'fulltext', 'global', 'Теги и ключевые слова', 1),

-- Специфичные веса для marketplace
('title', 0.95, 'fulltext', 'marketplace', 'Название объявления в маркетплейсе', 1),
('user_name', 0.3, 'fuzzy', 'marketplace', 'Имя продавца', 1),
('condition', 0.4, 'exact', 'marketplace', 'Состояние товара', 1),

-- Специфичные веса для storefront
('title', 0.9, 'fulltext', 'storefront', 'Название товара в магазине', 1),
('store_name', 0.5, 'fuzzy', 'storefront', 'Название магазина', 1),
('sku', 0.6, 'exact', 'storefront', 'Артикул товара', 1),
('availability', 0.3, 'exact', 'storefront', 'Наличие товара', 1);

-- Комментарии к таблицам и колонкам
COMMENT ON TABLE search_weights IS 'Конфигурация весов для различных полей поиска';
COMMENT ON COLUMN search_weights.field_name IS 'Название поля поиска (title, description, category, etc.)';
COMMENT ON COLUMN search_weights.weight IS 'Вес поля (0.0 - 1.0, где 1.0 - максимальная важность)';
COMMENT ON COLUMN search_weights.search_type IS 'Тип поиска: fulltext, fuzzy, exact';
COMMENT ON COLUMN search_weights.item_type IS 'Тип элемента: marketplace, storefront, global';
COMMENT ON COLUMN search_weights.category_id IS 'ID категории (NULL для глобальных весов)';
COMMENT ON COLUMN search_weights.version IS 'Версия конфигурации весов для отслеживания изменений';

COMMENT ON TABLE search_weights_history IS 'История изменений весов поиска для аудита и rollback';
COMMENT ON COLUMN search_weights_history.change_reason IS 'Причина изменения: manual, optimization, rollback, initialization';
COMMENT ON COLUMN search_weights_history.change_metadata IS 'Дополнительная информация об изменении в JSON формате';