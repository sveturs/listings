# Отчёт: Обновление OpenSearch индексатора (ФАЗА 3)

## Выполнено

### 1. Добавлены зависимости
- Импортирован `github.com/google/uuid` для работы с UUID
- Импортирован `strings` и `time` для обработки данных
- Добавлен `postgres.VariantRepository` для загрузки вариантов

### 2. Обновлена структура ListingIndexer
- Добавлено поле `variantRepo *postgres.VariantRepository`
- Репозиторий инициализируется в `NewListingIndexer()`

### 3. Реализованы helper функции

#### `extractBrand(attrs []AttributeForIndex) string`
Извлекает бренд из атрибутов товара по коду "brand"

#### `calculatePopularityScore(listing *domain.Listing) float64`
Рассчитывает популярность по формуле: `views*0.3 + favorites*0.5`

#### `isNewArrival(createdAt time.Time) bool`
Проверяет, что товар добавлен менее 7 дней назад

#### `unique(slice []string) []string`
Удаляет дубликаты из массива строк

#### `loadStorefrontData(ctx, storefrontID) (*StorefrontData, error)`
Загружает данные магазина: name, slug, rating, is_verified

#### `loadVariantsForListing(ctx, listing) ([]*domain.ProductVariantV2, error)`
Загружает все активные варианты товара с атрибутами

### 4. Расширен buildListingDocument()

#### Новые поля в OpenSearch документе:

**Популярность и флаги:**
- `popularity_score` - рассчитанный score
- `is_new_arrival` - новинка (< 7 дней)
- `is_promoted` - продвигается (TODO)
- `is_featured` - рекомендованный (TODO)

**Бренд:**
- `brand` - извлекается из атрибутов

**Теги:**
- `tags` - теги товара (уникальные)

**Переводы локации:**
- `country_sr`, `country_en`, `country_ru`
- `city_sr`, `city_en`, `city_ru`
- `address_sr`, `address_en`, `address_ru`

**Данные магазина (B2C):**
- `storefront_name`
- `storefront_slug`
- `storefront_rating`
- `seller_verified`

**Варианты товара (B2C):**
- `has_variants` (boolean)
- `variants_count` (int)
- `variant_skus` (array of strings)
- `variant_colors` (array, extracted from attributes)
- `variant_sizes` (array, extracted from attributes)
- `total_stock` (int, сумма доступного стока всех вариантов)
- `min_price`, `max_price` (float64)
- `variants` (nested array) - полные данные вариантов:
  - `sku`
  - `stock` (available = quantity - reserved)
  - `price`
  - `compare_at_price`
  - `attributes` (map атрибутов варианта)

### 5. Увеличен batch size
Изменён с 100 на 500 в `ReindexAllWithAttributes()`

### 6. Компиляция
Код успешно компилируется: `go build ./...` - без ошибок

## TODO (будущие улучшения)

- Реализовать поля `is_promoted`, `is_featured` (требуют логики промо-акций)
- Добавить `shipping_available`, `shipping_free`, `shipping_price` (требуют интеграцию с delivery service)
- Добавить `rating`, `review_count` (требуют интеграцию с reviews)
- Добавить `condition` (new/used/refurbished) - добавить в domain.Listing
- Добавить `old_price`, `discount_percent`, `has_discount` - рассчитывать на основе CompareAtPrice

## Итоги

Индексатор теперь поддерживает:
1. Полную индексацию вариантов B2C товаров
2. Извлечение данных магазина
3. Многоязычные переводы локации
4. Расчёт popularity score
5. Определение новинок
6. Извлечение бренда из атрибутов
7. Увеличенный batch size для производительности

Все изменения совместимы с существующим кодом, используют существующие репозитории и не ломают компиляцию.
