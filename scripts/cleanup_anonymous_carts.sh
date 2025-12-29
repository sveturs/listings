#!/bin/bash
#
# Скрипт для очистки старых анонимных корзин
# Использование: ./cleanup_anonymous_carts.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SQL_FILE="$SCRIPT_DIR/cleanup_anonymous_carts.sql"

# DB credentials (из docker-compose.yml)
DB_HOST="${VONDILISTINGS_DB_HOST:-localhost}"
DB_PORT="${VONDILISTINGS_DB_PORT:-35434}"
DB_USER="${VONDILISTINGS_DB_USER:-listings_user}"
DB_PASSWORD="${VONDILISTINGS_DB_PASSWORD:-listings_secret}"
DB_NAME="${VONDILISTINGS_DB_NAME:-listings_dev_db}"

CONNECTION_STRING="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME"

echo "🧹 Cleaning up anonymous carts (older than 7 days)..."
echo "📍 Database: $DB_NAME@$DB_HOST:$DB_PORT"
echo ""

psql "$CONNECTION_STRING" -f "$SQL_FILE"

echo ""
echo "✅ Cleanup completed!"
