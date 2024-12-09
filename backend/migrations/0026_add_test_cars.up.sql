-- backend/migrations/0026_add_test_cars.up.sql
-- Добавляем тестовые автомобили
INSERT INTO cars (
    make, model, year, price_per_day, location, latitude, longitude,
    description, seats, transmission, fuel_type, availability
) VALUES 
    -- Axia Neo
    ('Axia', 'Neo', 2023, 3500, 'Novi Sad, Serbia', 45.2671, 19.8335,
    'Axia Neo — это идеальный выбор для городских поездок. Современный дизайн и экологичный двигатель делают его популярным среди молодых профессионалов. Лёгкий в управлении, тихий и экономичный, этот компактный автомобиль станет вашим верным спутником в мегаполисе.',
    4, 'automatic', 'electric', true),

    -- Velaris Elegance
    ('Velaris', 'Elegance', 2022, 8000, 'Novi Sad, Serbia', 45.2551, 19.8452,
    'Velaris Elegance сочетает в себе роскошь и инновации. Мощный двигатель и просторный салон обеспечивают непревзойдённый комфорт даже на дальних поездках. Эта модель создаёт ощущение премиум-класса для тех, кто ценит утончённость и технологии.',
    5, 'automatic', 'petrol', true),

    -- Terra Cruiser X8
    ('Terra', 'Cruiser X8', 2023, 6500, 'Novi Sad, Serbia', 45.2460, 19.8235,
    'Terra Cruiser X8 — это надежный внедорожник, идеально подходящий для семейных поездок и активного отдыха. Благодаря просторному салону, гибридной системе и современным технологиям, он станет незаменимым помощником в любой ситуации.',
    7, 'automatic', 'hybrid', true),

    -- Navis Venture 9
    ('Navis', 'Venture 9', 2024, 5500, 'Novi Sad, Serbia', 45.2541, 19.8401,
    'Navis Venture 9 создан для больших компаний и комфортных путешествий. Благодаря вместительности, продуманной эргономике и современным технологиям, он обеспечивает максимальное удобство на дороге.',
    9, 'automatic', 'diesel', true);

-- Добавляем изображения для автомобилей
INSERT INTO car_images (
    car_id, file_path, file_name, file_size, content_type, is_main
) VALUES 
    (1, 'Axia_Neo.jpg', 'Axia_Neo.jpg', 1024, 'image/jpeg', true),
    (2, 'Velaris_Elegance.jpg', 'Velaris_Elegance.jpg', 1024, 'image/jpeg', true),
    (3, 'Terra_Cruiser.jpg', 'Terra_Cruiser.jpg', 1024, 'image/jpeg', true),
    (4, 'Navis_Venture.jpg', 'Navis_Venture.jpg', 1024, 'image/jpeg', true);

-- Добавляем особенности для каждого автомобиля
WITH new_cars AS (
  SELECT id, make, model FROM cars ORDER BY id DESC LIMIT 4
)
INSERT INTO car_feature_links (car_id, feature_id)
SELECT
  c.id,
  f.id
FROM new_cars c
CROSS JOIN car_features f
WHERE 
  (c.make = 'Axia' AND f.name IN (
    'Кондиционер', 'Климат-контроль', 'Навигация', 'Bluetooth', 'USB', 'Электропривод окон'
  ))
  OR 
  (c.make = 'Velaris' AND f.name IN (
    'Кондиционер', 'Климат-контроль', 'Круиз-контроль', 'Парктроники', 'Камера заднего вида',
    'Кожаный салон', 'Люк', 'Панорамная крыша', 'Подогрев сидений', 'Электропривод сидений',
    'Навигация', 'Bluetooth', 'USB', 'AUX', 'MP3'
  ))
  OR 
  (c.make = 'Terra' AND f.name IN (
    'Кондиционер', 'Климат-контроль', 'Круиз-контроль', 'Парктроники', 'Камера заднего вида',
    'Кожаный салон', 'Подогрев сидений', 'Навигация', 'Bluetooth', 'USB'
  ))
  OR 
  (c.make = 'Navis' AND f.name IN (
    'Кондиционер', 'Климат-контроль', 'Круиз-контроль', 'Парктроники', 'Камера заднего вида',
    'Навигация', 'Bluetooth', 'USB', 'AUX', 'MP3', 'Электропривод окон'
  ));