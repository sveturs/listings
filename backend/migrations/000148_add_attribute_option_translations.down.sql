-- Удаление всех добавленных переводов опций атрибутов
DELETE FROM attribute_option_translations WHERE attribute_name IN (
    'condition',
    'color',
    'transmission',
    'fuel_type',
    'body_type',
    'storage',
    'operating_system',
    'ram',
    'storage_type',
    'rooms',
    'warranty',
    'return_policy',
    'service_type',
    'availability',
    'service_area'
);