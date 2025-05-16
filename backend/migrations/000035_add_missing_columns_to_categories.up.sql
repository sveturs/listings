-- Добавляем недостающие колонки в marketplace_categories
DO $$ 
BEGIN
    -- Добавляем sort_order если не существует
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'marketplace_categories' AND column_name = 'sort_order'
    ) THEN
        ALTER TABLE marketplace_categories ADD COLUMN sort_order INT DEFAULT 0;
    END IF;
    
    -- Добавляем level если не существует
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'marketplace_categories' AND column_name = 'level'
    ) THEN
        ALTER TABLE marketplace_categories ADD COLUMN level INT DEFAULT 0;
    END IF;
    
    -- Добавляем count если не существует
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'marketplace_categories' AND column_name = 'count'
    ) THEN
        ALTER TABLE marketplace_categories ADD COLUMN count INT DEFAULT 0;
    END IF;
    
    -- Добавляем external_id если не существует
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'marketplace_categories' AND column_name = 'external_id'
    ) THEN
        ALTER TABLE marketplace_categories ADD COLUMN external_id VARCHAR(255);
    END IF;
END $$;

-- Создаем индексы для новых колонок
CREATE INDEX IF NOT EXISTS idx_marketplace_categories_sort_order ON marketplace_categories(sort_order);
CREATE INDEX IF NOT EXISTS idx_marketplace_categories_external_id ON marketplace_categories(external_id);