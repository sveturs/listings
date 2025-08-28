-- Подключаем ID атрибуты к категории автомобилей (Lični automobili - id: 1301)

-- Подключаем car_make_id к категории
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 
    1301 as category_id,
    ca.id as attribute_id,
    true as is_required,
    1 as sort_order
FROM category_attributes ca
WHERE ca.name = 'car_make_id'
ON CONFLICT (category_id, attribute_id) DO UPDATE
SET is_required = EXCLUDED.is_required,
    sort_order = EXCLUDED.sort_order;

-- Подключаем car_model_id к категории
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 
    1301 as category_id,
    ca.id as attribute_id,
    true as is_required,
    2 as sort_order
FROM category_attributes ca
WHERE ca.name = 'car_model_id'
ON CONFLICT (category_id, attribute_id) DO UPDATE
SET is_required = EXCLUDED.is_required,
    sort_order = EXCLUDED.sort_order;

-- Подключаем car_generation_id к категории (необязательный)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT 
    1301 as category_id,
    ca.id as attribute_id,
    false as is_required,
    3 as sort_order
FROM category_attributes ca
WHERE ca.name = 'car_generation_id'
ON CONFLICT (category_id, attribute_id) DO UPDATE
SET is_required = EXCLUDED.is_required,
    sort_order = EXCLUDED.sort_order;

-- Обновляем существующие атрибуты car_make и car_model, чтобы они не были обязательными
-- (так как теперь будут использоваться ID атрибуты)
UPDATE category_attribute_mapping cam
SET is_required = false
FROM category_attributes ca
WHERE cam.attribute_id = ca.id 
  AND cam.category_id = 1301
  AND ca.name IN ('car_make', 'car_model');

-- Обновляем конфигурацию car_make_id атрибута для использования таблицы car_makes
UPDATE category_attributes
SET 
    data_source = 'database',
    data_source_config = jsonb_build_object(
        'table', 'car_makes',
        'valueField', 'id',
        'labelField', 'name',
        'sortField', 'popularity_rs',
        'sortOrder', 'DESC'
    ),
    input_type = 'select',
    display_name = 'Marka automobila'
WHERE name = 'car_make_id';

-- Обновляем конфигурацию car_model_id атрибута для использования таблицы car_models
UPDATE category_attributes
SET 
    data_source = 'database',
    data_source_config = jsonb_build_object(
        'table', 'car_models',
        'valueField', 'id', 
        'labelField', 'name',
        'dependsOn', 'car_make_id',
        'dependsOnField', 'make_id',
        'sortField', 'name',
        'sortOrder', 'ASC'
    ),
    input_type = 'select',
    display_name = 'Model automobila'
WHERE name = 'car_model_id';

-- Обновляем конфигурацию car_generation_id атрибута
UPDATE category_attributes
SET 
    data_source = 'database',
    data_source_config = jsonb_build_object(
        'table', 'car_generations',
        'valueField', 'id',
        'labelField', 'name', 
        'dependsOn', 'car_model_id',
        'dependsOnField', 'model_id',
        'sortField', 'production_start',
        'sortOrder', 'DESC'
    ),
    input_type = 'select',
    display_name = 'Generacija'
WHERE name = 'car_generation_id';

-- Добавляем переводы для новых атрибутов
INSERT INTO translations (key, language, value)
VALUES 
    ('attribute.car_make_id', 'sr', 'Marka automobila'),
    ('attribute.car_make_id', 'en', 'Car Make'),
    ('attribute.car_make_id', 'ru', 'Марка автомобиля'),
    ('attribute.car_model_id', 'sr', 'Model automobila'),
    ('attribute.car_model_id', 'en', 'Car Model'),
    ('attribute.car_model_id', 'ru', 'Модель автомобиля'),
    ('attribute.car_generation_id', 'sr', 'Generacija'),
    ('attribute.car_generation_id', 'en', 'Generation'),
    ('attribute.car_generation_id', 'ru', 'Поколение')
ON CONFLICT (key, language) DO UPDATE 
SET value = EXCLUDED.value;