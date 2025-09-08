-- Добавление ключевых слов для категории "Audio i video oprema" (id=10191)
-- для поддержки определения DJ оборудования и музыкальных инструментов

-- Основные ключевые слова для DJ оборудования
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source)
VALUES 
    -- DJ оборудование - английский
    (10191, 'dj', 'en', 10.0, 'main', 'manual'),
    (10191, 'dj equipment', 'en', 10.0, 'main', 'manual'),
    (10191, 'dj controller', 'en', 9.0, 'main', 'manual'),
    (10191, 'dj mixer', 'en', 9.0, 'main', 'manual'),
    (10191, 'turntable', 'en', 8.0, 'main', 'manual'),
    (10191, 'cdj', 'en', 8.0, 'main', 'manual'),
    (10191, 'dj case', 'en', 7.0, 'synonym', 'manual'),
    (10191, 'flight case', 'en', 7.0, 'synonym', 'manual'),
    (10191, 'equipment case', 'en', 6.0, 'synonym', 'manual'),
    
    -- DJ оборудование - русский
    (10191, 'dj', 'ru', 10.0, 'main', 'manual'),
    (10191, 'диджей', 'ru', 10.0, 'main', 'manual'),
    (10191, 'dj оборудование', 'ru', 10.0, 'main', 'manual'),
    (10191, 'диджейское оборудование', 'ru', 10.0, 'main', 'manual'),
    (10191, 'dj контроллер', 'ru', 9.0, 'main', 'manual'),
    (10191, 'dj микшер', 'ru', 9.0, 'main', 'manual'),
    (10191, 'dj пульт', 'ru', 9.0, 'main', 'manual'),
    (10191, 'кейс', 'ru', 7.0, 'synonym', 'manual'),
    (10191, 'кейс для оборудования', 'ru', 7.0, 'synonym', 'manual'),
    (10191, 'транспортировочный кейс', 'ru', 7.0, 'synonym', 'manual'),
    
    -- DJ оборудование - сербский
    (10191, 'dj', 'sr', 10.0, 'main', 'manual'),
    (10191, 'dj oprema', 'sr', 10.0, 'main', 'manual'),
    (10191, 'dj kontroler', 'sr', 9.0, 'main', 'manual'),
    (10191, 'dj mikser', 'sr', 9.0, 'main', 'manual'),
    (10191, 'gramofon', 'sr', 8.0, 'main', 'manual'),
    (10191, 'kofer za opremu', 'sr', 7.0, 'synonym', 'manual'),
    
    -- Аудио оборудование - английский  
    (10191, 'audio', 'en', 8.0, 'main', 'manual'),
    (10191, 'audio equipment', 'en', 8.0, 'main', 'manual'),
    (10191, 'sound system', 'en', 8.0, 'main', 'manual'),
    (10191, 'speakers', 'en', 7.0, 'main', 'manual'),
    (10191, 'amplifier', 'en', 7.0, 'main', 'manual'),
    (10191, 'headphones', 'en', 6.0, 'main', 'manual'),
    (10191, 'microphone', 'en', 6.0, 'main', 'manual'),
    
    -- Аудио оборудование - русский
    (10191, 'аудио', 'ru', 8.0, 'main', 'manual'),
    (10191, 'аудио оборудование', 'ru', 8.0, 'main', 'manual'),
    (10191, 'звуковое оборудование', 'ru', 8.0, 'main', 'manual'),
    (10191, 'колонки', 'ru', 7.0, 'main', 'manual'),
    (10191, 'акустика', 'ru', 7.0, 'main', 'manual'),
    (10191, 'усилитель', 'ru', 7.0, 'main', 'manual'),
    (10191, 'наушники', 'ru', 6.0, 'main', 'manual'),
    (10191, 'микрофон', 'ru', 6.0, 'main', 'manual'),
    
    -- Аудио оборудование - сербский
    (10191, 'audio', 'sr', 8.0, 'main', 'manual'),
    (10191, 'audio oprema', 'sr', 8.0, 'main', 'manual'),
    (10191, 'zvučnici', 'sr', 7.0, 'main', 'manual'),
    (10191, 'pojačalo', 'sr', 7.0, 'main', 'manual'),
    (10191, 'slušalice', 'sr', 6.0, 'main', 'manual'),
    (10191, 'mikrofon', 'sr', 6.0, 'main', 'manual'),
    
    -- Видео оборудование
    (10191, 'video', 'en', 6.0, 'main', 'manual'),
    (10191, 'видео', 'ru', 6.0, 'main', 'manual'),
    (10191, 'video', 'sr', 6.0, 'main', 'manual'),
    
    -- Музыкальное оборудование
    (10191, 'music', 'en', 5.0, 'synonym', 'manual'),
    (10191, 'музыка', 'ru', 5.0, 'synonym', 'manual'),
    (10191, 'музыкальный', 'ru', 5.0, 'synonym', 'manual'),
    (10191, 'muzika', 'sr', 5.0, 'synonym', 'manual'),
    (10191, 'muzički', 'sr', 5.0, 'synonym', 'manual')
ON CONFLICT (category_id, keyword, language) DO UPDATE 
SET weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type,
    updated_at = CURRENT_TIMESTAMP;

-- Добавляем также ключевые слова для категории "TV i audio" (id=1103) как резервную
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source)
VALUES 
    (1103, 'dj', 'en', 6.0, 'synonym', 'manual'),
    (1103, 'dj equipment', 'en', 6.0, 'synonym', 'manual'),
    (1103, 'dj', 'ru', 6.0, 'synonym', 'manual'),
    (1103, 'диджей', 'ru', 6.0, 'synonym', 'manual'),
    (1103, 'dj', 'sr', 6.0, 'synonym', 'manual'),
    (1103, 'dj oprema', 'sr', 6.0, 'synonym', 'manual')
ON CONFLICT (category_id, keyword, language) DO UPDATE 
SET weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type,
    updated_at = CURRENT_TIMESTAMP;