#!/bin/bash

# Подключение к PostgreSQL внутри Docker-контейнера
CONTAINER_NAME="hostel_db"
DB_USER="postgres"
DB_NAME="hostel_db"

# Выполнить команды в psql внутри контейнера
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Список всех таблиц
SELECT table_schema, table_name
FROM information_schema.tables
WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
ORDER BY table_schema, table_name;

-- Информация о столбцах всех таблиц
SELECT table_name, column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
ORDER BY table_name, ordinal_position;

-- Список ограничений
SELECT table_name, constraint_type, constraint_name
FROM information_schema.table_constraints
WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
ORDER BY table_name;

-- Список индексов
SELECT
    t.relname AS table_name,
    i.relname AS index_name,
    a.attname AS column_name
FROM
    pg_class t,
    pg_class i,
    pg_index ix,
    pg_attribute a
WHERE
    t.oid = ix.indrelid
    AND i.oid = ix.indexrelid
    AND a.attrelid = t.oid
    AND a.attnum = ANY(ix.indkey)
ORDER BY t.relname, i.relname;

-- Завершение
EOF

