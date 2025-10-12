-- Откат миграции: удаление материализованного представления category_listing_counts

DROP MATERIALIZED VIEW IF EXISTS category_listing_counts CASCADE;
