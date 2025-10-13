# 🔄 ДЕТАЛЬНЫЙ ПЛАН РЕФАКТОРИНГА `marketplace.go`

**Дата создания:** 2025-10-13
**Последнее обновление:** 2025-10-13 (начало работы)
**Статус проекта:** Pre-production (НЕ в продакшне)
**Файл:** `backend/internal/proj/c2c/storage/postgres/marketplace.go`
**Текущий размер:** 3,761 строк, 46 функций

---

## ⚠️ КРИТИЧЕСКИ ВАЖНЫЕ ПРАВИЛА

### 🚫 ЗАПРЕЩЕНО:
1. ❌ **Оставлять рудименты** - проект НЕ в продакшне, обратная совместимость НЕ нужна!
2. ❌ **Работать без обновления этого плана** - каждый коммит = обновление статуса, обновлять этот файл плана
3. ❌ **Пропускать функциональные тесты** - после каждого этапа тестирование обязательно!
4. ❌ **Забывать применять миграции** - все изменения БД через `./migrator up`

### ✅ ОБЯЗАТЕЛЬНО:
1. ✅ **Актуализация плана** после каждой выполненной задачи
2. ✅ **Функциональные тесты** с токеном из `/tmp/token` после каждого этапа
3. ✅ **Применение миграций** при изменении схемы БД
4. ✅ **Голосовые уведомления** о прогрессе через `say`
5. ✅ **Создание новых файлов** вместо "TODO: создать" - делаем сразу!
6. ✅ **Удаление старого кода** после переноса - не оставляем дубликатов!

---

## 📊 EXECUTIVE SUMMARY

### Текущее состояние:
- **Размер файла:** 3,761 строк (God Object anti-pattern!)
- **Количество функций:** 46 методов
- **Проблемы:**
  - 🔴 Нарушение Single Responsibility Principle
  - 🔴 Сложность тестирования и поддержки
  - 🔴 Невозможность параллельной работы команды
  - 🔴 Высокий риск конфликтов при merge

### Целевое состояние:
- **8 отдельных файлов** по доменным областям
- **Средний размер:** ~470 строк на файл
- **4-8 методов** в каждом файле
- **100% покрытие тестами** критичных методов

### Трудозатраты:
- **Общее время:** ~8 часов
- **Разбиение и тестирование:** ~6 часов
- **Миграции и оптимизация:** ~2 часа

---

## 🗂️ ДОМЕННЫЙ АНАЛИЗ ФАЙЛА

### Выявленные доменные области:

| Домен | Методы | Строки | Приоритет |
|-------|--------|--------|-----------|
| **1. Listings CRUD** | 10 методов | ~1,100 строк | 🔴 HIGH |
| **2. Images** | 4 метода | ~280 строк | 🟡 MEDIUM |
| **3. Attributes** | 6 методов | ~780 строк | 🔴 HIGH |
| **4. Categories** | 8 методов | ~750 строк | 🟡 MEDIUM |
| **5. Favorites** | 6 методов | ~420 строк | 🟢 LOW |
| **6. Variants** | 4 метода | ~180 строк | 🟢 LOW |
| **7. Search & Queries** | 4 метода | ~240 строк | 🟡 MEDIUM |
| **8. Utilities** | 4 метода | ~100 строк | 🟢 LOW |

**ИТОГО:** 46 методов, ~3,850 строк (с учетом пустых строк и комментариев)

---

## 📋 ДЕТАЛЬНАЯ КАРТА МЕТОДОВ

### 1️⃣ LISTINGS CRUD (Основные операции)
**Целевой файл:** `listings_crud.go` (~1,100 строк) (также внедряем разделение c2c и b2c разделение)

```go
// Строки 122-332 (210 строк)
func (s *Storage) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error)

// Строки 445-814 (370 строк) - ОГРОМНЫЙ МЕТОД!
func (s *Storage) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error)

// Строки 1323-1458 (136 строк)
func (s *Storage) DeleteListing(ctx context.Context, id int, userID int) error

// Строки 1459-1590 (132 строк)
func (s *Storage) DeleteListingAdmin(ctx context.Context, id int) error

// Строки 1591-1695 (105 строк)
func (s *Storage) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error

// Строки 2946-3269 (324 строк) - ОГРОМНЫЙ МЕТОД!
func (s *Storage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)

// Строки 3270-3290 (21 строка)
func (s *Storage) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error)

// Строки 3291-3305 (15 строк)
func (s *Storage) IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error)

// Строки 3306-3335 (30 строк)
func (s *Storage) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error)

// Строки 3336-3395 (60 строк)
func (s *Storage) getStorefrontProductAsListing(ctx context.Context, id int) (*models.MarketplaceListing, error)
```

**Проблемы:**
- 🔴 `GetListings` - **370 строк** (нужно разбить!)
- 🔴 `GetListingByID` - **324 строки** (нужно разбить!)
- 🟡 Нет транзакций в `CreateListing`
- 🟡 Дублирование логики в Delete методах

**Приоритет:** 🔴 CRITICAL

---

### 2️⃣ IMAGES (Работа с изображениями)
**Целевой файл:** `listings_images.go` (~280 строк)

```go
// Строки 78-121 (44 строки) - utility функция
func buildFullImageURL(relativeURL string) string

// Строки 333-348 (16 строк)
func (s *Storage) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error)

// Строки 349-398 (50 строк)
func (s *Storage) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error)

// Строки 399-419 (21 строка)
func (s *Storage) DeleteListingImage(ctx context.Context, imageID string) (string, error)

// Строки 3589-3644 (56 строк)
func (s *Storage) GetB2CProductImages(ctx context.Context, productID int) ([]models.MarketplaceImage, error)
```

**Проблемы:**
- 🟡 `buildFullImageURL` - utility функция, возможно перенести в отдельный пакет `utils`
- 🟢 Методы простые, проблем нет

**Приоритет:** 🟡 MEDIUM

---

### 3️⃣ ATTRIBUTES (Атрибуты товаров)
**Целевой файл:** `listings_attributes.go` (~780 строк)

```go
// Строки 1696-1738 (43 строки) - utility функция
func sanitizeAttributeValue(attr *models.ListingAttributeValue)

// Строки 1739-1749 (11 строк) - utility функция
func isZeroValidValue(attrName string) bool

// Строки 1750-1939 (190 строк) - БОЛЬШОЙ МЕТОД!
func (s *Storage) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error

// Строки 1940-1963 (24 строки)
func (s *Storage) GetFormattedAttributeValue(ctx context.Context, attr models.ListingAttributeValue, language string) string

// Строки 1964-2205 (242 строки) - ОГРОМНЫЙ МЕТОД!
func (s *Storage) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error)

// Строки 2206-2341 (136 строк)
func (s *Storage) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error)

// Строки 2342-2356 (15 строк)
func (s *Storage) InvalidateAttributesCache(categoryID int)

// Строки 2357-2566 (210 строк) - ОГРОМНЫЙ МЕТОД!
func (s *Storage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error)
```

**Проблемы:**
- 🔴 `GetCategoryAttributes` - **210 строк** (нужно разбить!)
- 🔴 `GetListingAttributes` - **242 строки** (нужно разбить!)
- 🔴 `SaveListingAttributes` - **190 строк** (транзакции, валидация)
- 🟡 Кэширование в глобальных переменных (строки 40-48)

**Приоритет:** 🔴 HIGH

---

### 4️⃣ CATEGORIES (Категории)
**Целевой файл:** `categories.go` (~750 строк)

```go
// Строки 815-1003 (189 строк) - БОЛЬШОЙ МЕТОД!
func (s *Storage) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error)

// Строки 2567-2696 (130 строк)
func (s *Storage) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error)

// Строки 2697-2802 (106 строк)
func (s *Storage) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error)

// Строки 2803-2897 (95 строк)
func (s *Storage) GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error)

// Строки 2898-2945 (48 строк)
func (s *Storage) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error)

// Строки 3473-3588 (116 строк)
func (s *Storage) SearchCategories(ctx context.Context, query string, limit int) ([]models.MarketplaceCategory, error)
```

**Проблемы:**
- 🔴 `GetCategoryTree` - **189 строк** (рекурсия, сложная логика)
- 🟡 Дублирование SQL запросов между методами
- 🟡 Обработка переводов дублируется

**Приоритет:** 🟡 MEDIUM

---

### 5️⃣ FAVORITES (Избранное)
**Целевой файл:** `listings_favorites.go` (~420 строк)

```go
// Строки 1004-1012 (9 строк)
func (s *Storage) AddToFavorites(ctx context.Context, userID int, listingID int) error

// Строки 1013-1021 (9 строк)
func (s *Storage) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error

// Строки 1022-1031 (10 строк)
func (s *Storage) AddStorefrontToFavorites(ctx context.Context, userID int, productID int) error

// Строки 1032-1040 (9 строк)
func (s *Storage) RemoveStorefrontFromFavorites(ctx context.Context, userID int, productID int) error

// Строки 1041-1179 (139 строк)
func (s *Storage) GetUserStorefrontFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)

// Строки 1180-1298 (119 строк)
func (s *Storage) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)

// Строки 1299-1322 (24 строки)
func (s *Storage) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error)
```

**Проблемы:**
- 🟡 Дублирование между C2C и Storefront методами
- 🟢 Методы простые, проблем нет

**Приоритет:** 🟢 LOW

---

### 6️⃣ VARIANTS (Варианты товаров)
**Целевой файл:** `listings_variants.go` (~180 строк)

```go
// Строки 3645-3677 (33 строки)
func (s *Storage) CreateListingVariants(ctx context.Context, listingID int, variants []models.MarketplaceListingVariant) error

// Строки 3678-3719 (42 строки)
func (s *Storage) GetListingVariants(ctx context.Context, listingID int) ([]models.MarketplaceListingVariant, error)

// Строки 3720-3746 (27 строк)
func (s *Storage) UpdateListingVariant(ctx context.Context, variant *models.MarketplaceListingVariant) error

// Строки 3747-3761 (15 строк)
func (s *Storage) DeleteListingVariant(ctx context.Context, variantID int) error
```

**Проблемы:**
- 🟢 Методы простые и короткие
- 🟢 Хорошая структура

**Приоритет:** 🟢 LOW

---

### 7️⃣ SEARCH & QUERIES (Поиск и популярные запросы)
**Целевой файл:** `search_queries.go` (~240 строк)

```go
// Строки 3396-3440 (45 строк)
func (s *Storage) GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]service.SearchQuery, error)

// Строки 3441-3472 (32 строки)
func (s *Storage) SaveSearchQuery(ctx context.Context, query, normalizedQuery string, resultsCount int, language string) error
```

**Проблемы:**
- 🟢 Методы простые
- 🟡 `SearchCategories` логически относится сюда (перенесен в Categories выше для полноты)

**Приоритет:** 🟡 MEDIUM

---

### 8️⃣ UTILITIES (Вспомогательные функции)
**Целевой файл:** `storage_utils.go` (~100 строк)

```go
// Строки 51-56 (6 строк)
func init()

// Строки 64-70 (7 строк)
func NewStorage(...)

// Строки 73-75 (3 строки)
func (s *Storage) SetUserService(userService *authservice.UserService)

// Строки 420-444 (25 строк)
func (s *Storage) processTranslations(rawTranslations interface{}) models.TranslationMap
```

**Проблемы:**
- 🟡 Глобальный кэш (строки 40-48) - нужно переделать в Storage поля
- 🟢 Остальные методы простые

**Приоритет:** 🟢 LOW

---

## 🎯 ПЛАН ВЫПОЛНЕНИЯ (ПО ЭТАПАМ)

### 📦 ЭТАП 0: ПОДГОТОВКА (30 минут)

**Цель:** Создать инфраструктуру для рефакторинга

#### Задачи:
- [x] ✅ Создать план рефакторинга
- [ ] 📋 Создать backup оригинального файла
- [ ] 🧪 Подготовить тестовые сценарии с токеном `/tmp/token`
- [ ] 📁 Создать структуру новых файлов
- [ ] 🔧 Проверить, что все зависимости установлены

#### Команды:
```bash
# Backup
cp backend/internal/proj/c2c/storage/postgres/marketplace.go \
   backend/internal/proj/c2c/storage/postgres/marketplace.go.backup-2025-10-13

# Проверка токена
test -f /tmp/token && echo "✅ Token exists" || echo "❌ Token missing"

# Тестовый запрос
TOKEN="$(cat /tmp/token)" && \
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/unified/listings?limit=5' | jq '.data | length'
```

**Критерий выполнения:**
- ✅ Backup создан
- ✅ Токен работает
- ✅ Backend отвечает на запросы

**Время:** 30 минут
**Статус:** 🟡 В ПРОЦЕССЕ

---

### 📦 ЭТАП 1: UTILITIES & STORAGE INFRASTRUCTURE (1 час)

**Цель:** Создать базовую инфраструктуру Storage и вспомогательные функции

#### Задачи:
- [ ] Создать `storage_utils.go` с конструктором и utilities
- [ ] Перенести глобальный кэш в поля структуры Storage
- [ ] Создать базовый `storage.go` с общими интерфейсами
- [ ] Обновить импорты в зависимых файлах

#### Новые файлы:
```
backend/internal/proj/c2c/storage/postgres/
├── storage.go                    # NEW: Storage struct + NewStorage()
├── storage_utils.go              # NEW: processTranslations, cache management
└── storage_images_utils.go       # NEW: buildFullImageURL
```

#### Содержимое `storage.go`:
```go
package postgres

import (
	"sync"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	authservice "github.com/sveturs/auth/pkg/http/service"
	"backend/internal/proj/c2c/service"
)

// Storage представляет слой доступа к данным marketplace
type Storage struct {
	pool               *pgxpool.Pool
	translationService service.TranslationServiceInterface
	userService        *authservice.UserService

	// Cache для атрибутов
	attributeCacheMutex sync.RWMutex
	attributeCache      map[int][]models.CategoryAttribute
	attributeCacheTime  map[int]time.Time

	// Cache для ranges
	rangesCacheMutex sync.RWMutex
	rangesCache      map[int]map[string]map[string]interface{}
	rangesCacheTime  map[int]time.Time

	cacheTTL time.Duration
}

// NewStorage создает новый экземпляр Storage
func NewStorage(pool *pgxpool.Pool, translationService service.TranslationServiceInterface, userService *authservice.UserService) *Storage {
	return &Storage{
		pool:               pool,
		translationService: translationService,
		userService:        userService,
		attributeCache:     make(map[int][]models.CategoryAttribute),
		attributeCacheTime: make(map[int]time.Time),
		rangesCache:        make(map[int]map[string]map[string]interface{}),
		rangesCacheTime:    make(map[int]time.Time),
		cacheTTL:           30 * time.Minute,
	}
}

// SetUserService устанавливает UserService для Storage
func (s *Storage) SetUserService(userService *authservice.UserService) {
	s.userService = userService
}
```

#### Функциональные тесты:
```bash
# 1. Проверка создания Storage
cd backend && go test -v -run TestNewStorage ./internal/proj/c2c/storage/postgres/

# 2. Проверка что backend компилируется
cd backend && go build ./cmd/api/

# 3. Проверка API запроса
TOKEN="$(cat /tmp/token)" && \
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/marketplace/categories' | jq '.success'
# Ожидаем: true
```

**Критерий выполнения:**
- ✅ Все файлы созданы и компилируются
- ✅ Backend успешно запускается
- ✅ API тесты проходят
- ✅ План обновлен (этап 1 отмечен как выполнен)

**Трудозатраты:** 1 час
**Статус:** 🔴 НЕ НАЧАТО

---

### 📦 ЭТАП 2: IMAGES MODULE (45 минут)

**Цель:** Выделить работу с изображениями в отдельный модуль

#### Задачи:
- [ ] Создать `listings_images.go`
- [ ] Перенести методы работы с изображениями
- [ ] Обновить тесты
- [ ] Удалить перенесенный код из `marketplace.go`

#### Новый файл:
```
backend/internal/proj/c2c/storage/postgres/
└── listings_images.go            # NEW: 4 метода, ~280 строк
```

#### Методы для переноса:
```go
func buildFullImageURL(relativeURL string) string
func (s *Storage) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error)
func (s *Storage) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error)
func (s *Storage) DeleteListingImage(ctx context.Context, imageID string) (string, error)
func (s *Storage) GetB2CProductImages(ctx context.Context, productID int) ([]models.MarketplaceImage, error)
```

#### Функциональные тесты:
```bash
TOKEN="$(cat /tmp/token)"

# 1. Получить изображения листинга
LISTING_ID=$(curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/unified/listings?limit=1' | jq -r '.data[0].id')

curl -s -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/listings/${LISTING_ID}/images" | \
jq '.success'
# Ожидаем: true

# 2. Проверка URL изображений
curl -s -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/listings/${LISTING_ID}/images" | \
jq '.data[0].url' | grep -E "^https?://"
# Ожидаем: URL начинается с http:// или https://
```

**Критерий выполнения:**
- ✅ `listings_images.go` создан и работает
- ✅ Старый код удален из `marketplace.go`
- ✅ API тесты проходят
- ✅ Plan обновлен

**Трудозатраты:** 45 минут
**Статус:** 🔴 НЕ НАЧАТО

---

### 📦 ЭТАП 3: VARIANTS MODULE (30 минут)

**Цель:** Выделить работу с вариантами товаров

#### Задачи:
- [ ] Создать `listings_variants.go`
- [ ] Перенести методы вариантов
- [ ] Обновить тесты
- [ ] Удалить код из `marketplace.go`

#### Новый файл:
```
backend/internal/proj/c2c/storage/postgres/
└── listings_variants.go          # NEW: 4 метода, ~180 строк
```

#### Методы для переноса:
```go
func (s *Storage) CreateListingVariants(ctx context.Context, listingID int, variants []models.MarketplaceListingVariant) error
func (s *Storage) GetListingVariants(ctx context.Context, listingID int) ([]models.MarketplaceListingVariant, error)
func (s *Storage) UpdateListingVariant(ctx context.Context, variant *models.MarketplaceListingVariant) error
func (s *Storage) DeleteListingVariant(ctx context.Context, variantID int) error
```

#### Функциональные тесты:
```bash
TOKEN="$(cat /tmp/token)"

# 1. Получить варианты листинга
LISTING_ID=$(curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/unified/listings?limit=1' | jq -r '.data[0].id')

curl -s -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/listings/${LISTING_ID}/variants" | \
jq '.success'
# Ожидаем: true

# 2. Создать вариант (если есть права)
curl -s -X POST -H "Authorization: Bearer ${TOKEN}" \
-H "Content-Type: application/json" \
"http://localhost:3000/api/v1/admin/marketplace/listings/${LISTING_ID}/variants" \
-d '{"name": "Test Variant", "price": 100}' | jq '.success'
```

**Критерий выполнения:**
- ✅ `listings_variants.go` создан
- ✅ Все методы работают
- ✅ Старый код удален
- ✅ План обновлен

**Трудозатраты:** 30 минут
**Статус:** 🔴 НЕ НАЧАТО

---

### 📦 ЭТАП 4: FAVORITES MODULE (45 минут)

**Цель:** Выделить функционал избранного

#### Задачи:
- [ ] Создать `listings_favorites.go`
- [ ] Перенести методы favorites
- [ ] Рефакторинг дублирующегося кода (C2C vs Storefront)
- [ ] Удалить код из `marketplace.go`

#### Новый файл:
```
backend/internal/proj/c2c/storage/postgres/
└── listings_favorites.go         # NEW: 6 методов, ~420 строк
```

#### Методы для переноса:
```go
func (s *Storage) AddToFavorites(ctx context.Context, userID int, listingID int) error
func (s *Storage) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error
func (s *Storage) AddStorefrontToFavorites(ctx context.Context, userID int, productID int) error
func (s *Storage) RemoveStorefrontFromFavorites(ctx context.Context, userID int, productID int) error
func (s *Storage) GetUserStorefrontFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)
func (s *Storage) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)
func (s *Storage) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error)
```

#### Функциональные тесты:
```bash
TOKEN="$(cat /tmp/token)"

# 1. Получить избранное пользователя
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/marketplace/favorites' | jq '.success'
# Ожидаем: true

# 2. Добавить в избранное
LISTING_ID=$(curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/unified/listings?limit=1' | jq -r '.data[0].id')

curl -s -X POST -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/favorites/${LISTING_ID}" | jq '.success'
# Ожидаем: true

# 3. Удалить из избранного
curl -s -X DELETE -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/favorites/${LISTING_ID}" | jq '.success'
# Ожидаем: true
```

**Критерий выполнения:**
- ✅ `listings_favorites.go` создан
- ✅ Дублирующийся код устранен
- ✅ API тесты проходят
- ✅ План обновлен

**Трудозатраты:** 45 минут
**Статус:** 🔴 НЕ НАЧАТО

---

### 📦 ЭТАП 5: SEARCH & QUERIES MODULE (30 минут)

**Цель:** Выделить поисковые запросы

#### Задачи:
- [ ] Создать `search_queries.go`
- [ ] Перенести методы поиска
- [ ] Удалить код из `marketplace.go`

#### Новый файл:
```
backend/internal/proj/c2c/storage/postgres/
└── search_queries.go             # NEW: 2 метода, ~77 строк
```

#### Методы для переноса:
```go
func (s *Storage) GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]service.SearchQuery, error)
func (s *Storage) SaveSearchQuery(ctx context.Context, query, normalizedQuery string, resultsCount int, language string) error
```

#### Функциональные тесты:
```bash
TOKEN="$(cat /tmp/token)"

# 1. Получить популярные запросы
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/marketplace/search/popular?q=auto&limit=10' | \
jq '.success'
# Ожидаем: true

# 2. Сохранить поисковый запрос
curl -s -X POST -H "Authorization: Bearer ${TOKEN}" \
-H "Content-Type: application/json" \
'http://localhost:3000/api/v1/marketplace/search/save' \
-d '{"query": "test search", "results_count": 5}' | jq '.success'
# Ожидаем: true
```

**Критерий выполнения:**
- ✅ `search_queries.go` создан
- ✅ API тесты проходят
- ✅ План обновлен

**Трудозатраты:** 30 минут
**Статус:** 🔴 НЕ НАЧАТО

---

### 📦 ЭТАП 6: CATEGORIES MODULE (1.5 часа)

**Цель:** Выделить работу с категориями и оптимизировать

#### Задачи:
- [ ] Создать `categories.go`
- [ ] Перенести методы категорий
- [ ] Оптимизировать `GetCategoryTree` (189 строк!)
- [ ] Устранить дублирование SQL запросов
- [ ] Удалить код из `marketplace.go`

#### Новый файл:
```
backend/internal/proj/c2c/storage/postgres/
└── categories.go                 # NEW: 6 методов, ~750 строк (после оптимизации ~600)
```

#### Методы для переноса:
```go
func (s *Storage) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error)
func (s *Storage) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
func (s *Storage) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
func (s *Storage) GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error)
func (s *Storage) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error)
func (s *Storage) SearchCategories(ctx context.Context, query string, limit int) ([]models.MarketplaceCategory, error)
```

#### Оптимизации:
```go
// БЫЛО: GetCategoryTree - 189 строк с повторяющейся логикой
// БУДЕТ: Разбить на helper функции
func (s *Storage) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
	// Вызов helper'а
	categories, err := s.fetchCategoriesForTree(ctx)
	if err != nil {
		return nil, err
	}
	return s.buildCategoryTree(categories), nil
}

func (s *Storage) fetchCategoriesForTree(ctx context.Context) ([]models.MarketplaceCategory, error) {
	// SQL запрос (извлечено из GetCategoryTree)
}

func (s *Storage) buildCategoryTree(categories []models.MarketplaceCategory) []models.CategoryTreeNode {
	// Построение дерева (извлечено из GetCategoryTree)
}
```

#### Функциональные тесты:
```bash
TOKEN="$(cat /tmp/token)"

# 1. Получить дерево категорий
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/marketplace/categories/tree' | jq '.success'
# Ожидаем: true

# 2. Получить все категории
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/marketplace/categories' | jq '.data | length'
# Ожидаем: число > 0

# 3. Поиск категорий
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/marketplace/categories/search?q=auto' | \
jq '.data[0].name'
# Ожидаем: название категории

# 4. Популярные категории
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/marketplace/categories/popular?limit=5' | \
jq '.data | length'
# Ожидаем: 5
```

**Критерий выполнения:**
- ✅ `categories.go` создан
- ✅ `GetCategoryTree` оптимизирован (< 100 строк основной функции)
- ✅ Дублирование SQL устранено
- ✅ Все API тесты проходят
- ✅ План обновлен

**Трудозатраты:** 1.5 часа
**Статус:** 🔴 НЕ НАЧАТО

---

### 📦 ЭТАП 7: ATTRIBUTES MODULE (2 часа)

**Цель:** Выделить работу с атрибутами и оптимизировать

**⚠️ СЛОЖНОСТЬ: HIGH** - атрибуты используются везде!

#### Задачи:
- [ ] Создать `listings_attributes.go`
- [ ] Перенести методы атрибутов
- [ ] Оптимизировать `GetListingAttributes` (242 строки!)
- [ ] Оптимизировать `GetCategoryAttributes` (210 строк!)
- [ ] Рефакторинг `SaveListingAttributes` (190 строк)
- [ ] Удалить код из `marketplace.go`

#### Новый файл:
```
backend/internal/proj/c2c/storage/postgres/
└── listings_attributes.go        # NEW: 8 методов, ~780 строк (после оптимизации ~600)
```

#### Методы для переноса:
```go
func sanitizeAttributeValue(attr *models.ListingAttributeValue)
func isZeroValidValue(attrName string) bool
func (s *Storage) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error
func (s *Storage) GetFormattedAttributeValue(ctx context.Context, attr models.ListingAttributeValue, language string) string
func (s *Storage) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error)
func (s *Storage) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error)
func (s *Storage) InvalidateAttributesCache(categoryID int)
func (s *Storage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error)
```

#### Оптимизации:
```go
// БЫЛО: GetListingAttributes - 242 строки монолитный метод
// БУДЕТ: Разбить на helper'ы
func (s *Storage) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error) {
	rows, err := s.queryAttributeValues(ctx, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.processAttributeRows(ctx, rows)
}

func (s *Storage) queryAttributeValues(ctx context.Context, listingID int) (pgx.Rows, error) {
	// SQL запрос (извлечено)
}

func (s *Storage) processAttributeRows(ctx context.Context, rows pgx.Rows) ([]models.ListingAttributeValue, error) {
	// Обработка rows (извлечено)
}
```

#### Функциональные тесты:
```bash
TOKEN="$(cat /tmp/token)"

# 1. Получить атрибуты листинга
LISTING_ID=$(curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/unified/listings?limit=1' | jq -r '.data[0].id')

curl -s -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/listings/${LISTING_ID}/attributes" | \
jq '.success'
# Ожидаем: true

# 2. Получить атрибуты категории
CATEGORY_ID=$(curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/marketplace/categories?limit=1' | jq -r '.data[0].id')

curl -s -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/categories/${CATEGORY_ID}/attributes" | \
jq '.data | length'
# Ожидаем: число > 0

# 3. Получить ranges атрибутов
curl -s -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/categories/${CATEGORY_ID}/attributes/ranges" | \
jq '.success'
# Ожидаем: true

# 4. Создать листинг с атрибутами (admin)
curl -s -X POST -H "Authorization: Bearer ${TOKEN}" \
-H "Content-Type: application/json" \
'http://localhost:3000/api/v1/admin/marketplace/listings' \
-d '{
  "title": "Test Listing",
  "category_id": '${CATEGORY_ID}',
  "attributes": [
    {"attribute_id": 1, "value": "test"}
  ]
}' | jq '.success'
```

**Критерий выполнения:**
- ✅ `listings_attributes.go` создан
- ✅ Все методы оптимизированы (< 100 строк каждый)
- ✅ Кэш перенесен в Storage поля
- ✅ Все API тесты проходят
- ✅ План обновлен

**Трудозатраты:** 2 часа
**Статус:** 🔴 НЕ НАЧАТО

---

### 📦 ЭТАП 8: LISTINGS CRUD MODULE (2.5 часа)

**Цель:** Выделить основные CRUD операции с листингами

**⚠️ СЛОЖНОСТЬ: CRITICAL** - самые большие методы!

#### Задачи:
- [ ] Создать `listings_crud.go`
- [ ] Перенести CRUD методы
- [ ] **КРИТИЧНО:** Оптимизировать `GetListings` (370 строк!)
- [ ] **КРИТИЧНО:** Оптимизировать `GetListingByID` (324 строки!)
- [ ] Рефакторинг Delete методов (дублирование)
- [ ] Удалить код из `marketplace.go`

#### Новый файл:
```
backend/internal/proj/c2c/storage/postgres/
└── listings_crud.go              # NEW: 10 методов, ~1100 строк (после оптимизации ~800)
```

#### Методы для переноса:
```go
func (s *Storage) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error)
func (s *Storage) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error)
func (s *Storage) DeleteListing(ctx context.Context, id int, userID int) error
func (s *Storage) DeleteListingAdmin(ctx context.Context, id int) error
func (s *Storage) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error
func (s *Storage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
func (s *Storage) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error)
func (s *Storage) IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error)
func (s *Storage) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error)
func (s *Storage) getStorefrontProductAsListing(ctx context.Context, id int) (*models.MarketplaceListing, error)
```

#### Оптимизации:

**1. GetListings (370 строк → ~150 строк):**
```go
// БЫЛО: Огромный метод с вложенными условиями
func (s *Storage) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	// 370 строк монолитного кода
}

// БУДЕТ: Разбить на helper'ы
func (s *Storage) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	query, args := s.buildListingsQuery(filters, limit, offset)
	countQuery, countArgs := s.buildListingsCountQuery(filters)

	listings, err := s.executeListingsQuery(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.executeCountQuery(ctx, countQuery, countArgs)
	if err != nil {
		return nil, 0, err
	}

	return listings, total, nil
}

func (s *Storage) buildListingsQuery(filters map[string]string, limit, offset int) (string, []interface{}) {
	// Построение SQL запроса
}

func (s *Storage) buildListingsCountQuery(filters map[string]string) (string, []interface{}) {
	// Построение COUNT запроса
}

func (s *Storage) executeListingsQuery(ctx context.Context, query string, args []interface{}) ([]models.MarketplaceListing, error) {
	// Выполнение и обработка результатов
}
```

**2. GetListingByID (324 строки → ~120 строк):**
```go
// БЫЛО: Монолитный метод
func (s *Storage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	// 324 строки
}

// БУДЕТ: Разбить на helper'ы
func (s *Storage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	listing, err := s.fetchListingBase(ctx, id)
	if err != nil {
		return nil, err
	}

	// Параллельная загрузка связанных данных
	s.enrichListingData(ctx, listing)

	return listing, nil
}

func (s *Storage) fetchListingBase(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	// Базовый запрос листинга
}

func (s *Storage) enrichListingData(ctx context.Context, listing *models.MarketplaceListing) error {
	// Загрузка images, attributes, variants параллельно
}
```

**3. Объединить Delete методы:**
```go
// БЫЛО: DeleteListing (136 строк) + DeleteListingAdmin (132 строки) = 268 строк дублирования
// БУДЕТ: Один метод с параметром isAdmin
func (s *Storage) deleteListing(ctx context.Context, id int, userID int, isAdmin bool) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if !isAdmin {
		// Проверка владельца
		if err := s.checkListingOwnership(ctx, tx, id, userID); err != nil {
			return err
		}
	}

	// Удаление (общая логика)
	if err := s.executeListingDeletion(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Storage) DeleteListing(ctx context.Context, id int, userID int) error {
	return s.deleteListing(ctx, id, userID, false)
}

func (s *Storage) DeleteListingAdmin(ctx context.Context, id int) error {
	return s.deleteListing(ctx, id, 0, true)
}
```

#### Функциональные тесты:
```bash
TOKEN="$(cat /tmp/token)"

# 1. GetListings - базовый запрос
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/unified/listings?limit=10' | \
jq '.data | length'
# Ожидаем: 10

# 2. GetListings - с фильтрами
curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/unified/listings?category_id=1&min_price=100&max_price=1000' | \
jq '.success'
# Ожидаем: true

# 3. GetListingByID
LISTING_ID=$(curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/unified/listings?limit=1' | jq -r '.data[0].id')

curl -s -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/listings/${LISTING_ID}" | \
jq '.data.id'
# Ожидаем: LISTING_ID

# 4. CreateListing (admin)
NEW_LISTING=$(curl -s -X POST -H "Authorization: Bearer ${TOKEN}" \
-H "Content-Type: application/json" \
'http://localhost:3000/api/v1/admin/marketplace/listings' \
-d '{
  "title": "Test Refactoring Listing",
  "description": "Created during refactoring test",
  "price": 500,
  "category_id": 1
}')

echo "$NEW_LISTING" | jq '.success'
# Ожидаем: true

NEW_LISTING_ID=$(echo "$NEW_LISTING" | jq -r '.data.id')

# 5. UpdateListing
curl -s -X PUT -H "Authorization: Bearer ${TOKEN}" \
-H "Content-Type: application/json" \
"http://localhost:3000/api/v1/marketplace/listings/${NEW_LISTING_ID}" \
-d '{
  "title": "Updated Test Listing",
  "price": 600
}' | jq '.success'
# Ожидаем: true

# 6. DeleteListing (admin)
curl -s -X DELETE -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/admin/marketplace/listings/${NEW_LISTING_ID}" | \
jq '.success'
# Ожидаем: true

# 7. GetListingBySlug
SLUG=$(curl -s -H "Authorization: Bearer ${TOKEN}" \
'http://localhost:3000/api/v1/unified/listings?limit=1' | jq -r '.data[0].slug')

curl -s -H "Authorization: Bearer ${TOKEN}" \
"http://localhost:3000/api/v1/marketplace/listings/slug/${SLUG}" | \
jq '.success'
# Ожидаем: true
```

**Критерий выполнения:**
- ✅ `listings_crud.go` создан
- ✅ `GetListings` оптимизирован (< 150 строк)
- ✅ `GetListingByID` оптимизирован (< 120 строк)
- ✅ Delete методы объединены
- ✅ Все API тесты проходят
- ✅ План обновлен

**Трудозатраты:** 2.5 часа
**Статус:** 🔴 НЕ НАЧАТО

---

### 📦 ЭТАП 9: ФИНАЛЬНАЯ ОЧИСТКА (30 минут)

**Цель:** Удалить старый `marketplace.go` и финальная проверка

#### Задачи:
- [ ] Убедиться, что `marketplace.go` пустой (только package и import)
- [ ] Удалить `marketplace.go` (или оставить как stub с redirect комментариями)
- [ ] Проверка компиляции всего проекта
- [ ] Запуск всех функциональных тестов
- [ ] Проверка размеров новых файлов

#### Команды:
```bash
# 1. Проверить размер оригинального файла
wc -l backend/internal/proj/c2c/storage/postgres/marketplace.go
# Ожидаем: ~50 строк (только package + imports + комментарий)

# 2. Проверить новые файлы
fd "\.go$" backend/internal/proj/c2c/storage/postgres -x wc -l
# Ожидаем: 8 новых файлов, каждый < 800 строк

# 3. Компиляция
cd backend && go build ./cmd/api/
# Ожидаем: успешная компиляция без ошибок

# 4. Запуск backend
/home/dim/.local/bin/kill-port-3000.sh
screen -dmS backend-3000 bash -c 'cd /data/hostel-booking-system/backend && go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'

# Ждем 5 секунд
sleep 5

# 5. Проверка что backend запустился
curl -s http://localhost:3000/ | jq '.'
# Ожидаем: {"message": "Svetu API ...", "status": "ok"}
```

#### Финальные тесты:
```bash
# Запустить все тесты из предыдущих этапов
bash /data/hostel-booking-system/docs/test-marketplace-refactoring.sh
```

#### Удаление старого файла:
```bash
# Создать stub файл
cat > backend/internal/proj/c2c/storage/postgres/marketplace.go <<'EOF'
// backend/internal/proj/c2c/storage/postgres/marketplace.go
package postgres

/*
ЭТОТ ФАЙЛ УДАЛЕН В РАМКАХ РЕФАКТОРИНГА 2025-10-13

God Object (3,761 строк) разбит на модули:

├── storage.go                    # Storage struct + NewStorage()
├── storage_utils.go              # Utilities и cache management
├── listings_crud.go              # CRUD операции с листингами
├── listings_images.go            # Работа с изображениями
├── listings_attributes.go        # Атрибуты товаров
├── listings_favorites.go         # Избранное
├── listings_variants.go          # Варианты товаров
├── categories.go                 # Категории
└── search_queries.go             # Поисковые запросы

См. docs/MARKETPLACE_GO_REFACTORING_PLAN.md для деталей.
*/
EOF

# Или полностью удалить
rm backend/internal/proj/c2c/storage/postgres/marketplace.go
```

**Критерий выполнения:**
- ✅ `marketplace.go` удален или является stub'ом
- ✅ Все новые файлы существуют и работают
- ✅ Проект компилируется без ошибок
- ✅ Backend запускается успешно
- ✅ Все API тесты проходят
- ✅ План обновлен как завершенный

**Трудозатраты:** 30 минут
**Статус:** 🔴 НЕ НАЧАТО

---

### 📦 ЭТАП 10: ДОКУМЕНТАЦИЯ И КОММИТ (30 минут)

**Цель:** Оформить результаты рефакторинга

#### Задачи:
- [ ] Обновить основной план технического долга
- [ ] Создать migration guide для команды
- [ ] Создать коммит с изменениями
- [ ] Обновить метрики в плане

#### Команды:
```bash
# 1. Статистика изменений
echo "=== MARKETPLACE.GO REFACTORING STATS ===" > /tmp/refactoring-stats.txt
echo "" >> /tmp/refactoring-stats.txt

echo "БЫЛО:" >> /tmp/refactoring-stats.txt
echo "  marketplace.go: 3,761 строк, 46 функций" >> /tmp/refactoring-stats.txt
echo "" >> /tmp/refactoring-stats.txt

echo "СТАЛО:" >> /tmp/refactoring-stats.txt
fd "\.go$" backend/internal/proj/c2c/storage/postgres -x wc -l | \
  awk '{print "  " $2 ": " $1 " строк"}' >> /tmp/refactoring-stats.txt

cat /tmp/refactoring-stats.txt

# 2. Создать коммит
cd /data/hostel-booking-system
git add backend/internal/proj/c2c/storage/postgres/*.go
git add docs/MARKETPLACE_GO_REFACTORING_PLAN.md
git commit -m "refactor(backend): разбить marketplace.go God Object (3761 строк → 8 модулей)

- Разделение на доменные модули:
  - storage.go: базовая инфраструктура
  - listings_crud.go: CRUD операции
  - listings_images.go: работа с изображениями
  - listings_attributes.go: атрибуты
  - listings_favorites.go: избранное
  - listings_variants.go: варианты
  - categories.go: категории
  - search_queries.go: поисковые запросы

- Оптимизации:
  - GetListings: 370 → 150 строк
  - GetListingByID: 324 → 120 строк
  - GetCategoryTree: 189 → 100 строк
  - Объединены дублирующиеся Delete методы

- Улучшения:
  - Кэш перенесен в Storage поля (thread-safe)
  - Все методы < 150 строк
  - Улучшенная читаемость и maintainability
  - 100% функциональных тестов пройдено

См. docs/MARKETPLACE_GO_REFACTORING_PLAN.md"

# 3. Голосовое уведомление
say "рефакторинг маркетплейс го завершен успешно - восемь модулей созданы"
```

#### Обновление основного плана:
```bash
# Отметить задачу 19 как выполненную в COMPLETE_TECHNICAL_DEBT_AUDIT
# Обновить метрики качества Backend с 5.8/10 → 6.5/10
```

**Критерий выполнения:**
- ✅ Коммит создан
- ✅ Основной план обновлен
- ✅ Статистика сохранена
- ✅ Голосовое уведомление отправлено

**Трудозатраты:** 30 минут
**Статус:** 🔴 НЕ НАЧАТО

---

## 📊 ИТОГОВАЯ СТАТИСТИКА

### Целевая структура:
```
backend/internal/proj/c2c/storage/postgres/
├── storage.go                    # ~100 строк   (Storage struct + NewStorage)
├── storage_utils.go              # ~80 строк    (processTranslations, cache)
├── listings_crud.go              # ~800 строк   (10 методов CRUD)
├── listings_images.go            # ~280 строк   (5 методов images)
├── listings_attributes.go        # ~600 строк   (8 методов attributes)
├── listings_favorites.go         # ~420 строк   (6 методов favorites)
├── listings_variants.go          # ~180 строк   (4 метода variants)
├── categories.go                 # ~600 строк   (6 методов categories)
└── search_queries.go             # ~77 строк    (2 метода search)
```

**ИТОГО:** 8 файлов, ~3,137 строк (оптимизация ~624 строки = -16.6%)

### Метрики до/после:

| Метрика | ДО | ПОСЛЕ | Улучшение |
|---------|-------|-------|-----------|
| **Количество файлов** | 1 | 8 | +700% |
| **Средний размер файла** | 3,761 строк | ~390 строк | -89.6% |
| **Максимальный метод** | 370 строк | ~150 строк | -59.5% |
| **Maintainability Index** | 3/10 | 8/10 | +166% |
| **Testability** | 4/10 | 9/10 | +125% |
| **Code Quality** | 5/10 | 8/10 | +60% |

### Преимущества:
✅ **Модульность** - каждый домен в отдельном файле
✅ **Читаемость** - методы < 150 строк
✅ **Тестируемость** - легко писать unit тесты
✅ **Параллелизм** - команда может работать без конфликтов
✅ **Performance** - оптимизированы самые медленные методы
✅ **Maintainability** - легко найти и исправить баги

---

## 🎯 ОБЩИЙ ПРОГРЕСС

### Этапы выполнения:

- [ ] **ЭТАП 0:** Подготовка (30 мин) - 🟡 В ПРОЦЕССЕ
- [ ] **ЭТАП 1:** Utilities & Infrastructure (1 час) - 🔴 НЕ НАЧАТО
- [ ] **ЭТАП 2:** Images Module (45 мин) - 🔴 НЕ НАЧАТО
- [ ] **ЭТАП 3:** Variants Module (30 мин) - 🔴 НЕ НАЧАТО
- [ ] **ЭТАП 4:** Favorites Module (45 мин) - 🔴 НЕ НАЧАТО
- [ ] **ЭТАП 5:** Search & Queries (30 мин) - 🔴 НЕ НАЧАТО
- [ ] **ЭТАП 6:** Categories Module (1.5 часа) - 🔴 НЕ НАЧАТО
- [ ] **ЭТАП 7:** Attributes Module (2 часа) - 🔴 НЕ НАЧАТО
- [ ] **ЭТАП 8:** Listings CRUD (2.5 часа) - 🔴 НЕ НАЧАТО
- [ ] **ЭТАП 9:** Финальная очистка (30 мин) - 🔴 НЕ НАЧАТО
- [ ] **ЭТАП 10:** Документация (30 мин) - 🔴 НЕ НАЧАТО

**ОБЩИЙ ПРОГРЕСС:** 0/10 этапов (0%)
**ВРЕМЯ:** 0/8 часов
**ETA:** ~8 часов работы

---

## 📝 ИСТОРИЯ ИЗМЕНЕНИЙ

### 2025-10-13 (начало работы)
- ✅ Создан детальный план рефакторинга
- ✅ Проанализирована структура `marketplace.go`
- ✅ Выявлены 8 доменных областей
- ✅ Определены трудозатраты (~8 часов)
- 📋 Начат ЭТАП 0 (подготовка)

---

## ⚠️ РИСКИ И МИТИГАЦИЯ

### Риск 1: Поломка существующего функционала
**Вероятность:** MEDIUM
**Воздействие:** HIGH
**Митигация:**
- ✅ Backup оригинального файла
- ✅ Функциональные тесты после каждого этапа
- ✅ Постепенное разбиение (этап за этапом)
- ✅ Rollback план (восстановление из backup)

### Риск 2: Большой размер PR
**Вероятность:** HIGH
**Воздействие:** MEDIUM
**Митигация:**
- ✅ Разделение на 10 этапов
- ✅ Коммиты после каждого этапа (если нужно)
- ✅ Детальная документация изменений

### Риск 3: Забыть обновить план
**Вероятность:** MEDIUM
**Воздействие:** MEDIUM
**Митигация:**
- ✅ Голосовые напоминания через `say`
- ✅ Обязательное обновление после каждого этапа
- ✅ Чеклисты в каждом этапе

---

## 🔗 СВЯЗАННЫЕ ДОКУМЕНТЫ

- **Основной план технического долга:** `/data/hostel-booking-system/docs/COMPLETE_TECHNICAL_DEBT_AUDIT_2025-10-12.md`
- **CLAUDE.md:** `/data/hostel-booking-system/CLAUDE.md`
- **Тестовый скрипт:** `/data/hostel-booking-system/docs/test-marketplace-refactoring.sh` (будет создан)
- **Backup оригинала:** `marketplace.go.backup-2025-10-13` (будет создан)

---

**Последнее обновление:** 2025-10-13 (план создан)
**Автор:** Claude Code
**Статус:** 🟡 ПОДГОТОВКА К ВЫПОЛНЕНИЮ
