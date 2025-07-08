-- Тестовые данные для полного E2E тестирования поиска
-- Данные охватывают различные категории, города и языки

-- Недвижимость на кириллице
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, original_language, created_at, updated_at)
VALUES 
    (1, 1100, 'Студио у центру Београда', 'Комфоран студио апартман у самом центру града. Потпуно намештен, са свим потребним апаратима. Идеалан за студенте или младе парове.', 45000, 'good', 'active', 'Стариград, Београд', 44.8178, 20.4569, 'Београд', 'Србија', 'sr', NOW(), NOW()),
    
    (1, 1100, 'Двособан стан на Новом Београду', 'Светао двособан стан на 5. спрату са лифтом. Стан има велику терасу са погледом на реку. Паркинг место у гаражи.', 85000, 'excellent', 'active', 'Блок 45, Нови Београд', 44.8201, 20.4011, 'Београд', 'Србија', 'sr', NOW(), NOW()),
    
    (2, 1300, 'Кућа са баштом у Земуну', 'Породична кућа са великом баштом. Три спаваће собе, две купатила. Тиха улица, близу школе и вртића.', 150000, 'good', 'active', 'Земун, Београд', 44.8434, 20.4018, 'Београд', 'Србија', 'sr', NOW(), NOW()),
    
    (1, 1200, 'Соба за издавање', 'Велика соба у кући, засебан улаз. Намештена, интернет, све рачуне плаћа власник. Близу факултета.', 250, 'good', 'active', 'Вождовац, Београд', 44.7651, 20.4917, 'Београд', 'Србија', 'sr', NOW(), NOW());

-- Недвижимость на латинице  
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, original_language, created_at, updated_at)
VALUES 
    (1, 1100, 'Luksuzni apartman u Novom Sadu', 'Prekrasan trosoban apartman u centru Novog Sada. Kompletno renoviran, sa modernim nameštajem i aparatima.', 95000, 'excellent', 'active', 'Centar, Novi Sad', 45.2551, 19.8452, 'Novi Sad', 'Srbija', 'sr', NOW(), NOW()),
    
    (2, 1300, 'Vila na Zlatiboru', 'Luksuzna vila sa pogledom na planinu. Četiri spavaće sobe, sauna, jacuzzi. Idealno za odmor i relaksaciju.', 220000, 'new', 'active', 'Zlatibor', 43.7289, 19.7081, 'Zlatibor', 'Srbija', 'sr', NOW(), NOW()),
    
    (1, 1200, 'Soba za studente', 'Udobna soba u studentskom domu. Deljeno kupatilo i kuhinja. Sve usluge uključene u cenu.', 180, 'good', 'active', 'Studentski grad, Novi Sad', 45.2396, 19.8227, 'Novi Sad', 'Srbija', 'sr', NOW(), NOW());

-- Объявления на русском языке
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, original_language, created_at, updated_at)
VALUES 
    (2, 1100, 'Квартира в центре Белграда', 'Прекрасная двухкомнатная квартира в историческом центре города. Высокие потолки, большие окна, паркет.', 75000, 'excellent', 'active', 'Центр, Белград', 44.8178, 20.4569, 'Белград', 'Сербия', 'ru', NOW(), NOW()),
    
    (1, 1300, 'Дом с садом', 'Уютный дом с большим садом и гаражом. Три спальни, две ванные комнаты. Тихий район, хорошие соседи.', 120000, 'good', 'active', 'Пригород Белграда', 44.7868, 20.4489, 'Белград', 'Сербия', 'ru', NOW(), NOW()),
    
    (2, 1200, 'Комната для аренды', 'Светлая комната в квартире с хорошим ремонтом. Все удобства, интернет, рядом остановки общественного транспорта.', 300, 'good', 'active', 'Новый Белград', 44.8201, 20.4011, 'Белград', 'Сербия', 'ru', NOW(), NOW());

-- Разные категории товаров
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, original_language, created_at, updated_at)
VALUES 
    (1, 2000, 'Хонда Civic 2018', 'Одличан аутомобил у беспрекорном стању. Први власник, редовно сервисиран. Мала потрошња горива.', 12000, 'excellent', 'active', 'Београд', 44.8178, 20.4569, 'Београд', 'Србија', 'sr', NOW(), NOW()),
    
    (2, 2000, 'BMW X3 2020', 'Luksuzni SUV sa punim opremon. Kožne sedišta, navigacija, automatik. Uvoz iz Nemačke.', 35000, 'new', 'active', 'Novi Sad', 45.2551, 19.8452, 'Novi Sad', 'Srbija', 'sr', NOW(), NOW()),
    
    (1, 3000, 'iPhone 13 Pro', 'Телефон в отличном состоянии. Полный комплект, без царапин. Батарея держит весь день.', 800, 'excellent', 'active', 'Белград', 44.8178, 20.4569, 'Белград', 'Сербия', 'ru', NOW(), NOW()),
    
    (2, 3000, 'Samsung Galaxy S22', 'Нов телефон, непакован. Гаранција 2 године. Све боје доступне.', 650, 'new', 'active', 'Ниш', 43.3209, 21.8958, 'Ниш', 'Србија', 'sr', NOW(), NOW()),
    
    (1, 4000, 'Кухонный гарнитур', 'Модерна кухиња од масивног дрвета. Укључује све елементе и апарате. Могућа испорука и монтажа.', 2500, 'new', 'active', 'Београд', 44.8178, 20.4569, 'Београд', 'Србија', 'sr', NOW(), NOW()),
    
    (2, 4000, 'Диван-кровать', 'Удобный диван-кровать в отличном состоянии. Механизм работает идеально. Можно использовать как спальное место.', 350, 'good', 'active', 'Нови Сад', 45.2551, 19.8452, 'Нови Сад', 'Сербия', 'ru', NOW(), NOW());

-- Объявления с опечатками для тестирования fuzzy search
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, original_language, created_at, updated_at)
VALUES 
    (1, 1100, 'Апартман у Бегораду', 'Стан у центру Београда са грешком у наслову за тестирање нечетког претраге.', 50000, 'good', 'active', 'Београд', 44.8178, 20.4569, 'Београд', 'Србија', 'sr', NOW(), NOW()),
    
    (2, 2000, 'Аутомобил Фолксваген', 'Продајем Volkswagen Golf са намерном грешком у куцању за тестирање.', 8000, 'good', 'active', 'Нови Сад', 45.2551, 19.8452, 'Нови Сад', 'Србија', 'sr', NOW(), NOW()),
    
    (1, 3000, 'Мобилни телефон', 'Продајем мобилни телфон (намерна грешка) за тестирање алгоритма.', 200, 'used', 'active', 'Ниш', 43.3209, 21.8958, 'Ниш', 'Србија', 'sr', NOW(), NOW());

-- Объявления для тестирования транслитерации
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, original_language, created_at, updated_at)
VALUES 
    (1, 1100, 'Beograd centar stan', 'Stan u centru Beograda sa latiničnim naslovom ali ćirilicom u opisu: центар града, лепо место.', 60000, 'good', 'active', 'Стариград, Београд', 44.8178, 20.4569, 'Београд', 'Србија', 'sr', NOW(), NOW()),
    
    (2, 2000, 'Автомобил Тойота', 'Toyota Corolla са мешањем ћирилице и латинице у наслову за тестирање.', 9000, 'excellent', 'active', 'Београд', 44.8178, 20.4569, 'Београд', 'Србија', 'sr', NOW(), NOW()),
    
    (1, 3000, 'Laptop kompjuter', 'Лаптоп компјутер Dell са мешањем писама за тестирање транслитерације.', 500, 'used', 'active', 'Нови Сад', 45.2551, 19.8452, 'Нови Сад', 'Србија', 'sr', NOW(), NOW());

-- Объявления в разных городах
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, original_language, created_at, updated_at)
VALUES 
    (1, 1100, 'Стан у Крагујевцу', 'Леп стан у Крагујевцу, близу центра. Две собе, кухиња, купатило.', 40000, 'good', 'active', 'Крагујевац', 44.0165, 20.9114, 'Крагујевац', 'Србија', 'sr', NOW(), NOW()),
    
    (2, 1100, 'Квартира в Суботице', 'Прекрасная квартира в Суботице рядом с центром. Недавно сделан ремонт.', 55000, 'excellent', 'active', 'Суботица', 46.1008, 19.6677, 'Суботица', 'Сербия', 'ru', NOW(), NOW()),
    
    (1, 2000, 'Ауто у Панчеву', 'Продајем аутомобил у Панчеву. Опел Астра, добро стање.', 4500, 'good', 'active', 'Панчево', 44.8704, 20.6423, 'Панчево', 'Србија', 'sr', NOW(), NOW());

-- Объявления с особыми символами и длинными запросами
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, original_language, created_at, updated_at)
VALUES 
    (1, 4000, 'Намештај & декорација', 'Продајем намештај и декорацију за дом. Столице, столови, лампе, слике, огледала и остале ставке за уређење дома.', 150, 'used', 'active', 'Београд', 44.8178, 20.4569, 'Београд', 'Србија', 'sr', NOW(), NOW()),
    
    (2, 3000, 'Електроника - телевизор, радио, CD плејер', 'Веома дугачак наслов за тестирање претраге са више кључних речи и специјалних карактера као што су - цртица и запети.', 300, 'used', 'active', 'Нови Сад', 45.2551, 19.8452, 'Нови Сад', 'Србија', 'sr', NOW(), NOW()),
    
    (1, 5000, 'Модна одећа: јакне, панталоне, мајице', 'Широк избор модне одеће за све узрасте и прилике. Јакне, панталоне, мајице, хаљине, обућа и додаци.', 50, 'good', 'active', 'Ниш', 43.3209, 21.8958, 'Ниш', 'Србија', 'sr', NOW(), NOW());