-- Скрипт для очистки старых анонимных корзин
--
-- Назначение: Удаление анонимных корзин (user_id IS NULL), старше 7 дней
-- Запуск: psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -f cleanup_anonymous_carts.sql
-- Или через cron: 0 3 * * * psql ... -f cleanup_anonymous_carts.sql

BEGIN;

-- Удаляем cart_items из старых анонимных корзин (старше 7 дней)
DELETE FROM cart_items
WHERE cart_id IN (
  SELECT id FROM shopping_carts
  WHERE user_id IS NULL
    AND updated_at < NOW() - INTERVAL '7 days'
);

-- Удаляем сами анонимные корзины (старше 7 дней)
DELETE FROM shopping_carts
WHERE user_id IS NULL
  AND updated_at < NOW() - INTERVAL '7 days';

COMMIT;

-- Показать результат
SELECT
  COUNT(CASE WHEN user_id IS NULL THEN 1 END) as anonymous_carts,
  COUNT(CASE WHEN user_id IS NOT NULL THEN 1 END) as user_carts,
  COUNT(*) as total_carts
FROM shopping_carts;
