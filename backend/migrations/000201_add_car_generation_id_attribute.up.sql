-- Добавляем атрибут car_generation_id для категории Автомобили
INSERT INTO category_attributes (
    name,
    display_name,
    attribute_type,
    is_searchable,
    is_filterable,
    is_required,
    sort_order,
    show_in_card,
    show_in_list
) VALUES (
    'car_generation_id',
    'Car Generation ID',
    'number',
    true,
    true,
    false,
    25,
    false,
    false
);

-- Получаем ID созданного атрибута и добавляем его в категорию Автомобили (ID: 1003)
INSERT INTO category_attribute_mapping (category_id, attribute_id)
SELECT 1003, id FROM category_attributes WHERE name = 'car_generation_id';