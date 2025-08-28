-- Добавление ключевых слов для категории "Pića" (Beverages) ID=1802
-- Эти ключевые слова помогут AI правильно определять категорию для напитков

-- Ключевые слова на сербском (латиница)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- Пиво
(1802, 'pivo', 'sr', 1.0, 'main', 'manual'),
(1802, 'točeno pivo', 'sr', 0.9, 'main', 'manual'),
(1802, 'kraft pivo', 'sr', 0.9, 'main', 'manual'),
(1802, 'lager', 'sr', 0.8, 'main', 'manual'),
(1802, 'ale', 'sr', 0.8, 'main', 'manual'),
(1802, 'stout', 'sr', 0.8, 'main', 'manual'),
(1802, 'guinness', 'sr', 0.8, 'brand', 'manual'),
(1802, 'heineken', 'sr', 0.8, 'brand', 'manual'),
(1802, 'tuborg', 'sr', 0.8, 'brand', 'manual'),
(1802, 'jelen', 'sr', 0.8, 'brand', 'manual'),
(1802, 'lav', 'sr', 0.8, 'brand', 'manual'),
(1802, 'zaječarsko', 'sr', 0.8, 'brand', 'manual'),

-- Вино
(1802, 'vino', 'sr', 1.0, 'main', 'manual'),
(1802, 'crveno vino', 'sr', 0.9, 'main', 'manual'),
(1802, 'belo vino', 'sr', 0.9, 'main', 'manual'),
(1802, 'roze vino', 'sr', 0.9, 'main', 'manual'),
(1802, 'šampanjac', 'sr', 0.8, 'main', 'manual'),
(1802, 'penušavo vino', 'sr', 0.8, 'main', 'manual'),

-- Крепкий алкоголь
(1802, 'viski', 'sr', 0.9, 'main', 'manual'),
(1802, 'whiskey', 'sr', 0.9, 'main', 'manual'),
(1802, 'votka', 'sr', 0.9, 'main', 'manual'),
(1802, 'rakija', 'sr', 1.0, 'main', 'manual'),
(1802, 'šljivovica', 'sr', 0.9, 'main', 'manual'),
(1802, 'lozovača', 'sr', 0.9, 'main', 'manual'),
(1802, 'konjak', 'sr', 0.8, 'main', 'manual'),
(1802, 'rom', 'sr', 0.8, 'main', 'manual'),
(1802, 'džin', 'sr', 0.8, 'main', 'manual'),
(1802, 'tekila', 'sr', 0.8, 'main', 'manual'),
(1802, 'liker', 'sr', 0.8, 'main', 'manual'),

-- Безалкогольные напитки
(1802, 'sok', 'sr', 0.9, 'main', 'manual'),
(1802, 'prirodni sok', 'sr', 0.9, 'main', 'manual'),
(1802, 'gazirana pića', 'sr', 0.8, 'main', 'manual'),
(1802, 'koka kola', 'sr', 0.8, 'brand', 'manual'),
(1802, 'pepsi', 'sr', 0.8, 'brand', 'manual'),
(1802, 'fanta', 'sr', 0.8, 'brand', 'manual'),
(1802, 'sprite', 'sr', 0.8, 'brand', 'manual'),
(1802, 'mineralna voda', 'sr', 0.9, 'main', 'manual'),
(1802, 'voda', 'sr', 0.8, 'main', 'manual'),
(1802, 'flaširana voda', 'sr', 0.8, 'main', 'manual'),

-- Горячие напитки
(1802, 'kafa', 'sr', 0.9, 'main', 'manual'),
(1802, 'čaj', 'sr', 0.9, 'main', 'manual'),
(1802, 'kakao', 'sr', 0.8, 'main', 'manual'),
(1802, 'topla čokolada', 'sr', 0.8, 'main', 'manual'),

-- Энергетические напитки
(1802, 'energetski napitak', 'sr', 0.9, 'main', 'manual'),
(1802, 'red bull', 'sr', 0.8, 'brand', 'manual'),
(1802, 'monster', 'sr', 0.8, 'brand', 'manual'),

-- Ключевые слова на русском
(1802, 'пиво', 'ru', 1.0, 'main', 'manual'),
(1802, 'разливное пиво', 'ru', 0.9, 'main', 'manual'),
(1802, 'крафтовое пиво', 'ru', 0.9, 'main', 'manual'),
(1802, 'лагер', 'ru', 0.8, 'main', 'manual'),
(1802, 'эль', 'ru', 0.8, 'main', 'manual'),
(1802, 'стаут', 'ru', 0.8, 'main', 'manual'),
(1802, 'темное пиво', 'ru', 0.8, 'main', 'manual'),
(1802, 'светлое пиво', 'ru', 0.8, 'main', 'manual'),
(1802, 'гиннесс', 'ru', 0.8, 'brand', 'manual'),
(1802, 'хайнекен', 'ru', 0.8, 'brand', 'manual'),

-- Вино на русском
(1802, 'вино', 'ru', 1.0, 'main', 'manual'),
(1802, 'красное вино', 'ru', 0.9, 'main', 'manual'),
(1802, 'белое вино', 'ru', 0.9, 'main', 'manual'),
(1802, 'розовое вино', 'ru', 0.9, 'main', 'manual'),
(1802, 'шампанское', 'ru', 0.9, 'main', 'manual'),
(1802, 'игристое вино', 'ru', 0.8, 'main', 'manual'),

-- Крепкий алкоголь на русском
(1802, 'виски', 'ru', 0.9, 'main', 'manual'),
(1802, 'водка', 'ru', 0.9, 'main', 'manual'),
(1802, 'коньяк', 'ru', 0.9, 'main', 'manual'),
(1802, 'ром', 'ru', 0.8, 'main', 'manual'),
(1802, 'джин', 'ru', 0.8, 'main', 'manual'),
(1802, 'текила', 'ru', 0.8, 'main', 'manual'),
(1802, 'ликер', 'ru', 0.8, 'main', 'manual'),
(1802, 'ракия', 'ru', 0.9, 'main', 'manual'),

-- Безалкогольные на русском
(1802, 'сок', 'ru', 0.9, 'main', 'manual'),
(1802, 'натуральный сок', 'ru', 0.9, 'main', 'manual'),
(1802, 'газированные напитки', 'ru', 0.8, 'main', 'manual'),
(1802, 'кока-кола', 'ru', 0.8, 'brand', 'manual'),
(1802, 'пепси', 'ru', 0.8, 'brand', 'manual'),
(1802, 'минеральная вода', 'ru', 0.9, 'main', 'manual'),
(1802, 'вода', 'ru', 0.8, 'main', 'manual'),

-- Горячие напитки на русском
(1802, 'кофе', 'ru', 0.9, 'main', 'manual'),
(1802, 'чай', 'ru', 0.9, 'main', 'manual'),
(1802, 'какао', 'ru', 0.8, 'main', 'manual'),
(1802, 'горячий шоколад', 'ru', 0.8, 'main', 'manual'),

-- Энергетики на русском
(1802, 'энергетический напиток', 'ru', 0.9, 'main', 'manual'),
(1802, 'энергетик', 'ru', 0.9, 'main', 'manual'),

-- Ключевые слова на английском
(1802, 'beer', 'en', 1.0, 'main', 'manual'),
(1802, 'draft beer', 'en', 0.9, 'main', 'manual'),
(1802, 'craft beer', 'en', 0.9, 'main', 'manual'),
(1802, 'lager', 'en', 0.8, 'main', 'manual'),
(1802, 'ale', 'en', 0.8, 'main', 'manual'),
(1802, 'stout', 'en', 0.8, 'main', 'manual'),
(1802, 'pilsner', 'en', 0.8, 'main', 'manual'),
(1802, 'ipa', 'en', 0.8, 'main', 'manual'),

-- Вино на английском
(1802, 'wine', 'en', 1.0, 'main', 'manual'),
(1802, 'red wine', 'en', 0.9, 'main', 'manual'),
(1802, 'white wine', 'en', 0.9, 'main', 'manual'),
(1802, 'rose wine', 'en', 0.9, 'main', 'manual'),
(1802, 'champagne', 'en', 0.9, 'main', 'manual'),
(1802, 'sparkling wine', 'en', 0.8, 'main', 'manual'),

-- Крепкий алкоголь на английском
(1802, 'whiskey', 'en', 0.9, 'main', 'manual'),
(1802, 'vodka', 'en', 0.9, 'main', 'manual'),
(1802, 'rum', 'en', 0.8, 'main', 'manual'),
(1802, 'gin', 'en', 0.8, 'main', 'manual'),
(1802, 'tequila', 'en', 0.8, 'main', 'manual'),
(1802, 'cognac', 'en', 0.8, 'main', 'manual'),
(1802, 'brandy', 'en', 0.8, 'main', 'manual'),
(1802, 'liqueur', 'en', 0.8, 'main', 'manual'),

-- Безалкогольные на английском
(1802, 'juice', 'en', 0.9, 'main', 'manual'),
(1802, 'soft drink', 'en', 0.8, 'main', 'manual'),
(1802, 'soda', 'en', 0.8, 'main', 'manual'),
(1802, 'cola', 'en', 0.8, 'main', 'manual'),
(1802, 'mineral water', 'en', 0.9, 'main', 'manual'),
(1802, 'water', 'en', 0.8, 'main', 'manual'),
(1802, 'bottled water', 'en', 0.8, 'main', 'manual'),

-- Горячие напитки на английском
(1802, 'coffee', 'en', 0.9, 'main', 'manual'),
(1802, 'tea', 'en', 0.9, 'main', 'manual'),
(1802, 'cocoa', 'en', 0.8, 'main', 'manual'),
(1802, 'hot chocolate', 'en', 0.8, 'main', 'manual'),

-- Энергетики на английском
(1802, 'energy drink', 'en', 0.9, 'main', 'manual'),

-- Общие термины
(1802, 'напиток', 'ru', 0.7, 'main', 'manual'),
(1802, 'напитки', 'ru', 0.7, 'main', 'manual'),
(1802, 'алкоголь', 'ru', 0.8, 'main', 'manual'),
(1802, 'алкогольные напитки', 'ru', 0.8, 'main', 'manual'),
(1802, 'безалкогольные напитки', 'ru', 0.8, 'main', 'manual'),
(1802, 'napitak', 'sr', 0.7, 'main', 'manual'),
(1802, 'napitci', 'sr', 0.7, 'main', 'manual'),
(1802, 'alkohol', 'sr', 0.8, 'main', 'manual'),
(1802, 'alkoholna pića', 'sr', 0.8, 'main', 'manual'),
(1802, 'bezalkoholna pića', 'sr', 0.8, 'main', 'manual'),
(1802, 'beverage', 'en', 0.7, 'main', 'manual'),
(1802, 'beverages', 'en', 0.7, 'main', 'manual'),
(1802, 'drink', 'en', 0.7, 'main', 'manual'),
(1802, 'drinks', 'en', 0.7, 'main', 'manual'),
(1802, 'alcohol', 'en', 0.8, 'main', 'manual'),
(1802, 'alcoholic beverages', 'en', 0.8, 'main', 'manual'),
(1802, 'non-alcoholic beverages', 'en', 0.8, 'main', 'manual');

-- Обновляем счетчик использования для существующих ключевых слов
UPDATE category_keywords 
SET usage_count = usage_count + 1, 
    success_rate = 0.95,
    updated_at = CURRENT_TIMESTAMP
WHERE category_id = 1802;