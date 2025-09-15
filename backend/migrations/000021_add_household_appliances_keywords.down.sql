-- Удаляем добавленные ключевые слова для бытовой техники

DELETE FROM category_keywords
WHERE category_id = 1005
AND keyword IN (
    -- English
    'vacuum', 'cleaner', 'vacuum cleaner', 'hoover', 'dust', 'suction',
    'miele', 'dyson', 'bosch', 'samsung', 'electrolux', 'philips',
    'appliance', 'domestic', 'iron', 'washing', 'dishwasher',
    'microwave', 'refrigerator', 'fridge', 'oven', 'stove', 'cooker',
    'blender', 'mixer', 'toaster', 'kettle', 'coffee', 'kitchen',
    -- Russian
    'пылесос', 'уборка', 'пыль', 'чистка', 'мешок', 'фильтр', 'щетка',
    'техника', 'бытовая', 'домашняя', 'утюг', 'стиральная', 'посудомоечная',
    'микроволновка', 'холодильник', 'плита', 'духовка', 'варочная',
    'блендер', 'миксер', 'тостер', 'чайник', 'кофеварка', 'кофемашина',
    'кухня', 'кухонная',
    -- Serbian
    'usisivač', 'čišćenje', 'prašina', 'vrećica', 'četka',
    'aparat', 'kućni', 'pegla', 'veš', 'mašina', 'mikrotalasna',
    'frižider', 'šporet', 'rerna', 'blender', 'mikser', 'toster',
    'kuvalo', 'kuhinja'
);

-- Откатываем обновления счетчиков
UPDATE category_keywords
SET usage_count = usage_count - 1, success_rate = 0.80
WHERE category_id = 1005
AND keyword IN ('home', 'garden', 'дом', 'сад', 'household', 'быт')
AND usage_count > 0;