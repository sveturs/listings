-- Миграция 000043: Добавление атрибутов для категории Jobs
-- Дата: 03.09.2025
-- Цель: Расширить функциональность для вакансий и работы

-- Добавляем специфичные атрибуты для вакансий
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, is_active, created_at) VALUES
('job_type', 'Employment Type', 'Тип занятости', 'select', 'regular', true, CURRENT_TIMESTAMP),
('job_category', 'Job Category', 'Категория работы', 'select', 'regular', true, CURRENT_TIMESTAMP),
('salary_min', 'Min Salary', 'Минимальная зарплата', 'number', 'regular', true, CURRENT_TIMESTAMP),
('salary_max', 'Max Salary', 'Максимальная зарплата', 'number', 'regular', true, CURRENT_TIMESTAMP),
('salary_currency', 'Currency', 'Валюта', 'select', 'regular', true, CURRENT_TIMESTAMP),
('experience_years', 'Experience Years', 'Лет опыта', 'number', 'regular', true, CURRENT_TIMESTAMP),
('education_level', 'Education', 'Образование', 'select', 'regular', true, CURRENT_TIMESTAMP),
('remote_work', 'Remote Work', 'Удалённая работа', 'boolean', 'regular', true, CURRENT_TIMESTAMP),
('work_schedule', 'Work Schedule', 'График работы', 'select', 'regular', true, CURRENT_TIMESTAMP),
('company_name', 'Company', 'Компания', 'text', 'regular', true, CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;

-- Обновляем опции для атрибутов типа select (хранятся в JSONB поле options)
UPDATE unified_attributes
SET options = CASE code 
    WHEN 'job_type' THEN '[
        {"value": "full_time", "label": "Полная занятость"},
        {"value": "part_time", "label": "Частичная занятость"},
        {"value": "contract", "label": "Контракт"},
        {"value": "temporary", "label": "Временная"},
        {"value": "internship", "label": "Стажировка"},
        {"value": "freelance", "label": "Фриланс"},
        {"value": "volunteer", "label": "Волонтёрство"}
    ]'::jsonb
    WHEN 'job_category' THEN '[
        {"value": "it", "label": "IT и программирование"},
        {"value": "sales", "label": "Продажи"},
        {"value": "marketing", "label": "Маркетинг"},
        {"value": "finance", "label": "Финансы"},
        {"value": "hr", "label": "HR"},
        {"value": "admin", "label": "Администрирование"},
        {"value": "education", "label": "Образование"},
        {"value": "healthcare", "label": "Здравоохранение"},
        {"value": "engineering", "label": "Инженерия"},
        {"value": "construction", "label": "Строительство"},
        {"value": "transport", "label": "Транспорт"},
        {"value": "hospitality", "label": "Гостиничный бизнес"},
        {"value": "retail", "label": "Розничная торговля"},
        {"value": "manufacturing", "label": "Производство"},
        {"value": "other", "label": "Другое"}
    ]'::jsonb
    WHEN 'salary_currency' THEN '[
        {"value": "RSD", "label": "Динары (RSD)"},
        {"value": "EUR", "label": "Евро (EUR)"},
        {"value": "USD", "label": "Доллары (USD)"}
    ]'::jsonb
    WHEN 'education_level' THEN '[
        {"value": "none", "label": "Без образования"},
        {"value": "primary", "label": "Начальное"},
        {"value": "secondary", "label": "Среднее"},
        {"value": "vocational", "label": "Среднее специальное"},
        {"value": "bachelor", "label": "Высшее (бакалавр)"},
        {"value": "master", "label": "Высшее (магистр)"},
        {"value": "phd", "label": "Учёная степень"}
    ]'::jsonb
    WHEN 'work_schedule' THEN '[
        {"value": "standard", "label": "Стандартный (5/2)"},
        {"value": "shift", "label": "Сменный"},
        {"value": "flexible", "label": "Гибкий"},
        {"value": "rotating", "label": "Вахтовый"},
        {"value": "weekend", "label": "Только выходные"},
        {"value": "night", "label": "Ночной"}
    ]'::jsonb
    ELSE options
END
WHERE code IN ('job_type', 'job_category', 'salary_currency', 'education_level', 'work_schedule');

-- Привязываем атрибуты к категории Jobs и её подкатегориям
INSERT INTO unified_category_attributes (category_id, attribute_id, is_required, is_enabled, sort_order)
SELECT 
    c.id as category_id,
    ua.id as attribute_id,
    CASE 
        WHEN ua.code IN ('job_type', 'job_category', 'company_name') THEN true
        ELSE false
    END as is_required,
    true as is_enabled,
    ROW_NUMBER() OVER (PARTITION BY c.id ORDER BY 
        CASE ua.code
            WHEN 'job_type' THEN 1
            WHEN 'job_category' THEN 2
            WHEN 'company_name' THEN 3
            WHEN 'salary_min' THEN 4
            WHEN 'salary_max' THEN 5
            WHEN 'salary_currency' THEN 6
            WHEN 'experience_years' THEN 7
            WHEN 'education_level' THEN 8
            WHEN 'work_schedule' THEN 9
            WHEN 'remote_work' THEN 10
            ELSE 11
        END
    ) as sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE (c.id = 1018 OR c.parent_id = 1018) -- Jobs и её подкатегории
  AND ua.code IN ('job_type', 'job_category', 'salary_min', 'salary_max', 'salary_currency', 
                  'experience_years', 'education_level', 'remote_work', 'work_schedule', 'company_name')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавляем переводы для новых атрибутов
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) 
SELECT 'unified_attribute', ua.id, t.lang, 'display_name', t.translation
FROM unified_attributes ua
CROSS JOIN (VALUES
    ('job_type', 'en', 'Employment Type'),
    ('job_type', 'sr', 'Tip zaposlenja'),
    ('job_category', 'en', 'Job Category'),
    ('job_category', 'sr', 'Kategorija posla'),
    ('salary_min', 'en', 'Min Salary'),
    ('salary_min', 'sr', 'Minimalna plata'),
    ('salary_max', 'en', 'Max Salary'),
    ('salary_max', 'sr', 'Maksimalna plata'),
    ('salary_currency', 'en', 'Currency'),
    ('salary_currency', 'sr', 'Valuta'),
    ('experience_years', 'en', 'Years of Experience'),
    ('experience_years', 'sr', 'Godine iskustva'),
    ('education_level', 'en', 'Education Level'),
    ('education_level', 'sr', 'Nivo obrazovanja'),
    ('remote_work', 'en', 'Remote Work'),
    ('remote_work', 'sr', 'Rad od kuće'),
    ('work_schedule', 'en', 'Work Schedule'),
    ('work_schedule', 'sr', 'Raspored rada'),
    ('company_name', 'en', 'Company'),
    ('company_name', 'sr', 'Kompanija')
) AS t(attr_code, lang, translation)
WHERE ua.code = t.attr_code
ON CONFLICT DO NOTHING;

-- Логирование результата
DO $$
DECLARE
    attr_count INTEGER;
    cat_count INTEGER;
    option_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO attr_count 
    FROM unified_attributes 
    WHERE code IN ('job_type', 'job_category', 'salary_min', 'salary_max', 'salary_currency', 
                   'experience_years', 'education_level', 'remote_work', 'work_schedule', 'company_name');
    
    SELECT COUNT(DISTINCT category_id) INTO cat_count
    FROM unified_category_attributes uca
    JOIN unified_attributes ua ON uca.attribute_id = ua.id
    WHERE ua.code IN ('job_type', 'job_category', 'salary_min', 'salary_max', 'salary_currency', 
                      'experience_years', 'education_level', 'remote_work', 'work_schedule', 'company_name');
    
    -- Опции теперь хранятся в JSONB, поэтому просто считаем атрибуты с типом select
    SELECT COUNT(*) INTO option_count
    FROM unified_attributes
    WHERE code IN ('job_type', 'job_category', 'salary_currency', 'education_level', 'work_schedule');
    
    RAISE NOTICE 'Добавлено % атрибутов для вакансий с % опциями, привязано к % категориям', 
                 attr_count, option_count, cat_count;
END $$;