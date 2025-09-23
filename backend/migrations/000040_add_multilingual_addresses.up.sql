-- Добавляем поле для хранения мультиязычных адресов в JSON формате
ALTER TABLE marketplace_listings
ADD COLUMN IF NOT EXISTS address_multilingual JSONB;

-- Создаем индекс для быстрого поиска по мультиязычным адресам
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_address_multilingual
ON marketplace_listings USING gin (address_multilingual);

-- Добавляем комментарий для документации
COMMENT ON COLUMN marketplace_listings.address_multilingual IS
'Мультиязычные версии адреса в формате JSON: {"sr": "...", "en": "...", "ru": "..."}';

-- Для существующего объявления 325 добавим мультиязычные адреса
UPDATE marketplace_listings
SET address_multilingual = jsonb_build_object(
    'sr', 'Васе Стајића, Нови Сад, Јужнобачки округ',
    'en', 'Vase Stajica, Novi Sad, South Bačka District',
    'ru', 'Васе Стайича, Нови Сад, Южно-Бачский округ'
)
WHERE id = 325;