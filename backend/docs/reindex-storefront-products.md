# Переиндексация товаров витрин в OpenSearch

## Описание

Система поддерживает полнотекстовый поиск товаров витрин через OpenSearch. При создании или обновлении товаров они автоматически индексируются, но иногда может потребоваться полная переиндексация.

## Когда нужна переиндексация

1. После изменения структуры индекса
2. После массового импорта товаров через SQL
3. При рассинхронизации данных между PostgreSQL и OpenSearch
4. После восстановления базы данных из бэкапа

## Способы переиндексации

### 1. Через скрипт (рекомендуется)

```bash
cd backend
./reindex-products.sh
```

Опции:
- `--create` - пересоздать индекс перед переиндексацией
- `--batch SIZE` - размер пакета для индексации (по умолчанию 100)
- `--index NAME` - имя индекса (по умолчанию storefront_products)

### 2. Через Makefile

```bash
cd backend
make reindex_products         # Использует скрипт
make reindex_products_direct  # Запускает напрямую Go утилиту
```

### 3. Через Go утилиту напрямую

```bash
cd backend
go run ./cmd/reindex-products/main.go --batch 100
```

## Проверка результатов

### Количество документов в индексе
```bash
curl -s "http://localhost:9200/storefront_products/_count" | jq '.count'
```

### Поиск конкретного товара
```bash
curl -s "http://localhost:9200/storefront_products/_search?q=НАЗВАНИЕ_ТОВАРА" | jq '.hits.total.value'
```

### Проверка через API
```bash
curl -s "http://localhost:3001/api/v1/search?query=НАЗВАНИЕ_ТОВАРА&product_types=storefront" | jq '.total'
```

## Структура индекса

Индекс `storefront_products` содержит следующие поля:
- `product_id` - ID товара
- `storefront_id` - ID витрины
- `name` - название товара
- `description` - описание
- `price` - цена
- `category_id` - ID категории
- `attributes` - атрибуты товара
- `inventory` - информация о наличии
- `search_keywords` - ключевые слова для поиска

## Решение проблем

### Товары не находятся через поиск

1. Убедитесь, что товары активны (`is_active = true`)
2. Проверьте, что товары проиндексированы:
   ```bash
   curl -s "http://localhost:9200/storefront_products/_search?q=product_id:ID_ТОВАРА"
   ```
3. Запустите переиндексацию

### Ошибка подключения к OpenSearch

1. Проверьте, что OpenSearch запущен:
   ```bash
   docker ps | grep opensearch
   ```
2. Проверьте настройки в `.env`:
   ```
   OPENSEARCH_URL=http://localhost:9200
   OPENSEARCH_USERNAME=admin
   OPENSEARCH_PASSWORD=admin
   ```

### Медленная индексация

Увеличьте размер пакета:
```bash
./reindex-products.sh --batch 500
```