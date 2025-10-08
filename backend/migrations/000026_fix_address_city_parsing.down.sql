-- Откат миграции исправления парсинга адресов
-- Удаляем созданные функции

DROP FUNCTION IF EXISTS extract_city_from_location(TEXT);
DROP FUNCTION IF EXISTS extract_country_from_location(TEXT);

-- Примечание: Мы не откатываем обновленные данные address_city и address_country,
-- так как невозможно восстановить предыдущие (неправильные) значения.
-- Если необходим полный откат, восстановите данные из бэкапа.
