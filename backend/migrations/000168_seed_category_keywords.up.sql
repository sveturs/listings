-- Начальные ключевые слова для категорий
-- Используем ON CONFLICT для безопасного добавления

-- Шины и колеса (1304)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type) VALUES
-- Русский
(1304, 'шина', 'ru', 10.0, 'main'),
(1304, 'резина', 'ru', 9.0, 'synonym'),
(1304, 'покрышка', 'ru', 8.0, 'synonym'),
(1304, 'колесо', 'ru', 8.0, 'synonym'),
(1304, 'диск', 'ru', 7.0, 'context'),
(1304, 'летняя', 'ru', 5.0, 'attribute'),
(1304, 'зимняя', 'ru', 5.0, 'attribute'),
(1304, 'всесезонная', 'ru', 5.0, 'attribute'),
(1304, 'протектор', 'ru', 4.0, 'context'),
(1304, 'балансировка', 'ru', 3.0, 'context'),
(1304, 'шиномонтаж', 'ru', 3.0, 'context'),
-- English
(1304, 'tire', 'en', 10.0, 'main'),
(1304, 'tyre', 'en', 10.0, 'main'),
(1304, 'wheel', 'en', 8.0, 'synonym'),
(1304, 'rim', 'en', 7.0, 'context'),
(1304, 'summer', 'en', 5.0, 'attribute'),
(1304, 'winter', 'en', 5.0, 'attribute'),
(1304, 'all-season', 'en', 5.0, 'attribute'),
(1304, 'all season', 'en', 5.0, 'attribute'),
(1304, 'tread', 'en', 4.0, 'context'),
-- Serbian
(1304, 'guma', 'sr', 10.0, 'main'),
(1304, 'točak', 'sr', 8.0, 'synonym'),
(1304, 'felna', 'sr', 7.0, 'context'),
(1304, 'letnja', 'sr', 5.0, 'attribute'),
(1304, 'zimska', 'sr', 5.0, 'attribute'),
(1304, 'celogodišnja', 'sr', 5.0, 'attribute'),
-- Бренды (универсальные)
(1304, 'michelin', '*', 3.0, 'brand'),
(1304, 'continental', '*', 3.0, 'brand'),
(1304, 'bridgestone', '*', 3.0, 'brand'),
(1304, 'pirelli', '*', 3.0, 'brand'),
(1304, 'goodyear', '*', 3.0, 'brand'),
(1304, 'dunlop', '*', 3.0, 'brand'),
(1304, 'yokohama', '*', 3.0, 'brand'),
-- Размеры (паттерны)
(1304, '205/55', '*', 2.0, 'pattern'),
(1304, '195/65', '*', 2.0, 'pattern'),
(1304, '225/45', '*', 2.0, 'pattern'),
(1304, 'R16', '*', 2.0, 'pattern'),
(1304, 'R17', '*', 2.0, 'pattern'),
(1304, 'R15', '*', 2.0, 'pattern')
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type;

-- Исключающие слова для шин
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative) VALUES
(1304, 'велосипед', 'ru', 8.0, 'context', true),
(1304, 'bicycle', 'en', 8.0, 'context', true),
(1304, 'bicikl', 'sr', 8.0, 'context', true),
(1304, 'ластик', 'ru', 10.0, 'context', true),
(1304, 'eraser', 'en', 10.0, 'context', true),
(1304, 'сапоги', 'ru', 8.0, 'context', true),
(1304, 'boots', 'en', 8.0, 'context', true)
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    is_negative = EXCLUDED.is_negative;

-- Двигатель и детали двигателя (1305)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type) VALUES
-- Русский
(1305, 'двигатель', 'ru', 10.0, 'main'),
(1305, 'мотор', 'ru', 10.0, 'main'),
(1305, 'движок', 'ru', 8.0, 'synonym'),
(1305, 'поршень', 'ru', 7.0, 'context'),
(1305, 'клапан', 'ru', 7.0, 'context'),
(1305, 'коленвал', 'ru', 7.0, 'context'),
(1305, 'распредвал', 'ru', 7.0, 'context'),
(1305, 'турбина', 'ru', 6.0, 'context'),
(1305, 'карбюратор', 'ru', 6.0, 'context'),
(1305, 'инжектор', 'ru', 6.0, 'context'),
-- English
(1305, 'engine', 'en', 10.0, 'main'),
(1305, 'motor', 'en', 10.0, 'main'),
(1305, 'piston', 'en', 7.0, 'context'),
(1305, 'valve', 'en', 7.0, 'context'),
(1305, 'crankshaft', 'en', 7.0, 'context'),
(1305, 'camshaft', 'en', 7.0, 'context'),
(1305, 'turbo', 'en', 6.0, 'context'),
(1305, 'carburetor', 'en', 6.0, 'context'),
-- Serbian
(1305, 'motor', 'sr', 10.0, 'main'),
(1305, 'klip', 'sr', 7.0, 'context'),
(1305, 'ventil', 'sr', 7.0, 'context'),
(1305, 'radilica', 'sr', 7.0, 'context'),
(1305, 'bregasta', 'sr', 7.0, 'context')
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type;

-- Автомобили (1301)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type) VALUES
-- Русский
(1301, 'автомобиль', 'ru', 10.0, 'main'),
(1301, 'машина', 'ru', 10.0, 'main'),
(1301, 'авто', 'ru', 9.0, 'synonym'),
(1301, 'тачка', 'ru', 7.0, 'synonym'),
(1301, 'седан', 'ru', 6.0, 'attribute'),
(1301, 'хэтчбек', 'ru', 6.0, 'attribute'),
(1301, 'универсал', 'ru', 6.0, 'attribute'),
(1301, 'кроссовер', 'ru', 6.0, 'attribute'),
(1301, 'внедорожник', 'ru', 6.0, 'attribute'),
-- English
(1301, 'car', 'en', 10.0, 'main'),
(1301, 'automobile', 'en', 10.0, 'main'),
(1301, 'vehicle', 'en', 9.0, 'synonym'),
(1301, 'sedan', 'en', 6.0, 'attribute'),
(1301, 'hatchback', 'en', 6.0, 'attribute'),
(1301, 'wagon', 'en', 6.0, 'attribute'),
(1301, 'crossover', 'en', 6.0, 'attribute'),
(1301, 'SUV', 'en', 6.0, 'attribute'),
-- Serbian
(1301, 'automobil', 'sr', 10.0, 'main'),
(1301, 'auto', 'sr', 10.0, 'main'),
(1301, 'kola', 'sr', 9.0, 'synonym'),
(1301, 'vozilo', 'sr', 9.0, 'synonym'),
(1301, 'limuzina', 'sr', 6.0, 'attribute'),
(1301, 'karavan', 'sr', 6.0, 'attribute'),
-- Бренды автомобилей
(1301, 'volkswagen', '*', 4.0, 'brand'),
(1301, 'toyota', '*', 4.0, 'brand'),
(1301, 'mercedes', '*', 4.0, 'brand'),
(1301, 'bmw', '*', 4.0, 'brand'),
(1301, 'audi', '*', 4.0, 'brand'),
(1301, 'ford', '*', 4.0, 'brand'),
(1301, 'opel', '*', 4.0, 'brand'),
(1301, 'peugeot', '*', 4.0, 'brand'),
(1301, 'renault', '*', 4.0, 'brand'),
(1301, 'skoda', '*', 4.0, 'brand'),
(1301, 'škoda', '*', 4.0, 'brand')
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type;

-- Недвижимость - квартиры (1401)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type) VALUES
-- Русский
(1401, 'квартира', 'ru', 10.0, 'main'),
(1401, 'квартиру', 'ru', 10.0, 'main'),
(1401, 'однокомнатная', 'ru', 8.0, 'attribute'),
(1401, 'двухкомнатная', 'ru', 8.0, 'attribute'),
(1401, 'трехкомнатная', 'ru', 8.0, 'attribute'),
(1401, 'студия', 'ru', 8.0, 'attribute'),
(1401, 'новостройка', 'ru', 7.0, 'context'),
(1401, 'вторичка', 'ru', 7.0, 'context'),
(1401, 'ремонт', 'ru', 5.0, 'context'),
(1401, 'балкон', 'ru', 4.0, 'context'),
(1401, 'этаж', 'ru', 4.0, 'context'),
-- English
(1401, 'apartment', 'en', 10.0, 'main'),
(1401, 'flat', 'en', 10.0, 'main'),
(1401, 'condo', 'en', 9.0, 'synonym'),
(1401, 'studio', 'en', 8.0, 'attribute'),
(1401, 'bedroom', 'en', 6.0, 'context'),
(1401, 'balcony', 'en', 4.0, 'context'),
(1401, 'floor', 'en', 4.0, 'context'),
-- Serbian
(1401, 'stan', 'sr', 10.0, 'main'),
(1401, 'jednosoban', 'sr', 8.0, 'attribute'),
(1401, 'dvosoban', 'sr', 8.0, 'attribute'),
(1401, 'trosoban', 'sr', 8.0, 'attribute'),
(1401, 'garsonjera', 'sr', 8.0, 'attribute'),
(1401, 'novogradnja', 'sr', 7.0, 'context'),
(1401, 'balkon', 'sr', 4.0, 'context'),
(1401, 'sprat', 'sr', 4.0, 'context')
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type;

-- Электроника - смартфоны (1101)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type) VALUES
-- Русский
(1101, 'телефон', 'ru', 10.0, 'main'),
(1101, 'смартфон', 'ru', 10.0, 'main'),
(1101, 'мобильный', 'ru', 8.0, 'context'),
(1101, 'мобильник', 'ru', 8.0, 'synonym'),
(1101, 'айфон', 'ru', 7.0, 'brand'),
(1101, 'самсунг', 'ru', 7.0, 'brand'),
(1101, 'сяоми', 'ru', 7.0, 'brand'),
-- English
(1101, 'phone', 'en', 10.0, 'main'),
(1101, 'smartphone', 'en', 10.0, 'main'),
(1101, 'mobile', 'en', 8.0, 'context'),
(1101, 'cellphone', 'en', 8.0, 'synonym'),
(1101, 'iphone', 'en', 7.0, 'brand'),
(1101, 'samsung', 'en', 7.0, 'brand'),
(1101, 'xiaomi', 'en', 7.0, 'brand'),
-- Serbian
(1101, 'telefon', 'sr', 10.0, 'main'),
(1101, 'mobilni', 'sr', 8.0, 'context'),
(1101, 'ajfon', 'sr', 7.0, 'brand')
ON CONFLICT (category_id, keyword, language) DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type;

-- Функция для автоматического обновления usage_count
CREATE OR REPLACE FUNCTION increment_keyword_usage(
    p_category_id INTEGER,
    p_keywords TEXT[],
    p_language VARCHAR(2)
) RETURNS VOID AS $$
BEGIN
    UPDATE category_keywords
    SET usage_count = usage_count + 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE category_id = p_category_id
        AND keyword = ANY(p_keywords)
        AND (language = p_language OR language = '*');
END;
$$ LANGUAGE plpgsql;

-- Комментарий к функции
COMMENT ON FUNCTION increment_keyword_usage IS 'Увеличивает счетчик использования ключевых слов при поиске категории';