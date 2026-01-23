# CHANGELOG - Listings Service

## [Unreleased]

### Fixed - 2026-01-23 (b077baeeb)

**Расширен cleanup старых кэшей на ВСЕ jobs в CI workflow**

#### Проблема
- ❌ Validate job failed on checkout: `File was unable to be removed Error: EACCES: permission denied`
- ❌ Первоначальный fix был применён только к Build job, но Validate job (запускается первым) имел ту же проблему
- ❌ Все jobs с checkout нуждались в cleanup read-only `.gomodcache` файлов

#### Решение
- ✅ Добавлен шаг "Cleanup old caches" ко ВСЕМ jobs: Validate, Lint, Test, Build, Docker, Integration-test
- ✅ Используется `chmod -R +w` перед удалением старых кэшей в каждом job
- ✅ Отключен `clean: true` в checkout для ВСЕХ jobs (установлен `clean: false`)
- ✅ Добавлен GOMODCACHE env var в integration-test job для консистентности

#### Файлы
- `.github/workflows/ci.yml` - cleanup step добавлен в 6 jobs (Validate, Lint, Test, Build, Docker, Integration-test)

---

### Fixed - 2026-01-23 (88271b6e3)

**Добавлен cleanup старых кэшей перед Build job**

#### Проблема
- ❌ Build job failed on checkout: `File was unable to be removed Error: EACCES: permission denied`
- ❌ Checkout не может удалить `.gomodcache` из-за read-only файлов Go модулей
- ❌ `clean: true` в checkout конфликтует с Go module cache permissions

#### Решение
- ✅ Добавлен шаг "Cleanup old caches" в Build job
- ✅ Используется `chmod -R +w` перед удалением старых кэшей
- ✅ Отключен `clean: true` в checkout для Build job

#### Файлы
- `.github/workflows/ci.yml` - добавлен cleanup step в Build job

---

### Fixed - 2026-01-23 (b28a6848f)

**Использование локального GOMODCACHE для предотвращения проблем с правами доступа**

#### Проблема
- ❌ Test/Lint jobs failed: `open /home/svetu/go/pkg/mod/.../file.go: no such file or directory`
- ❌ Protobuf errors: `pattern editions_defaults.binpb: no matching files found`
- ❌ Системный GOMODCACHE (`/home/svetu/go/pkg/mod/`) недоступен для чтения

#### Решение
- ✅ Добавлен `GOMODCACHE: ${{ github.workspace }}/.gomodcache` для всех jobs
- ✅ Создаётся локальная директория `.gomodcache` в workspace
- ✅ Go модули скачиваются в контролируемую директорию с правильными правами

#### Файлы
- `.github/workflows/ci.yml` - настроен локальный GOMODCACHE

---

### Fixed - 2026-01-23 (e235346bb)

**Отключение Go module cache в CI для совместимости с TMPDIR**

#### Проблема
- ❌ Test job failed: `open /home/svetu/go/pkg/mod/.../file.go: no such file or directory`
- ❌ Protobuf errors: `pattern editions_defaults.binpb: no matching files found`
- Конфликт между `cache: true` в setup-go и TMPDIR=/opt/github-runner/tmp

#### Решение
- ✅ Отключен `cache: true` в Test и Build jobs (установлен `cache: false`)
- ✅ Сохранена безопасная конфигурация TMPDIR для предотвращения выполнения кода в /tmp

#### Файлы
- `.github/workflows/ci.yml` - отключен cache в Test и Build jobs

---

### Fixed - 2026-01-23 (1af08a95d)

**Исправление дубликата метода в chat_service_test.go**

#### Проблема
- ❌ CI checks failed: typecheck error - method MockMessageRepository.GetLatestMessage declared twice
- ❌ Locations: chat_service_test.go:153 and chat_service_test.go:181

#### Решение
- ✅ Удален дублирующий метод GetLatestMessage (строки 181-187)
- ✅ Оставлено одно корректное объявление метода

#### Файлы
- `internal/service/chat_service_test.go` - удален дубликат метода

---

### Fixed - 2026-01-23 (167263b0c)

**Исправление интеграции с Auth Service и добавление поддержки переводов**

#### Изменённые компоненты
- `internal/service/chat_service.go` - обновлен для использования нового API Auth Service
- `internal/domain/listing.go` - добавлена поддержка многоязычности (title/description translations)
- `internal/domain/jsonb_types.go` - новый файл с JSONB типами для PostgreSQL
- `internal/indexer/listing_indexer.go` - добавлена индексация переводов в OpenSearch
- `internal/repository/opensearch/client.go` - обновления для поддержки многоязычного поиска
- `internal/repository/postgres/repository.go` - изменения в работе с БД
- `internal/transport/grpc/converters.go` - обновление gRPC конвертеров
- Proto файлы - регенерация всех .pb.go файлов
- Множество репозиториев - обновления зависимостей и импортов

#### Проблемы решены
- ❌ Ошибка компиляции: `s.authService.UserService undefined`
- ✅ Решение: Auth Service API изменился, UserService() больше не существует
- ✅ Обновлено: теперь GetUser() вызывается напрямую на AuthService
- ❌ Конфликты при слиянии веток в listing_indexer.go
- ✅ Решение: сохранена логика индексации переводов при разрешении конфликта

#### Изменения в БД
- База данных: `listings_dev_db`
- Добавлена поддержка JSONB полей для переводов
- OpenSearch индекс `marketplace_listings` обновлён для многоязычного поиска

#### Дополнительно
- Исправлены все 56 изменённых файлов (1 новый, 55 обновлённых)
- Все proto файлы регенерированы
- Код скомпилирован успешно

---

### Fixed - 2026-01-18 (3a3042184)

**Storefront Data Integrity: Parse settings and create related records**

Добавлен парсинг `settings` JSONB для автоматического создания связанных записей при создании витрины.

#### Problem
- Listings microservice получал `settings` JSONB от монолита, сохранял в `storefronts.settings`, но НЕ создавал связанные записи
- Таблицы `storefront_hours`, `storefront_payment_methods`, `storefront_delivery_options` оставались пустыми
- Данные были в settings, но не использовались

#### Solution
- `internal/service/listings/storefront_service.go` - Добавлен `processStorefrontSettings()` после `CreateStorefront`:
  - ✅ Unmarshal settings JSONB → map[string]interface{}
  - ✅ Парсит `businessHours` → вызывает `processBusinessHours()` → создаёт 7 записей
  - ✅ Парсит `paymentMethods` → вызывает `processPaymentMethods()` → создаёт N записей
  - ✅ Парсит `deliveryOptions` → вызывает `processDeliveryOptions()` → создаёт M записей
  - ✅ Graceful error handling - продолжает при ошибках парсинга

#### Files Changed
- `internal/service/listings/storefront_service.go` (4 new methods):
  - `processStorefrontSettings()` - main parser
  - `processBusinessHours()` - parses businessHours array → creates storefront_hours records
  - `processPaymentMethods()` - parses paymentMethods array → creates storefront_payment_methods records
  - `processDeliveryOptions()` - parses deliveryOptions array → creates storefront_delivery_options records

#### Result
- ✅ Автоматическое создание связанных записей при создании витрины
- ✅ 100% data integrity для settings данных
- ✅ Repository методы (SetWorkingHours, SetPaymentMethods, SetDeliveryOptions) reviewed и работают корректно

---

### Fixed - 2026-01-12 (53b7a8981)

**OpenSearch: Fix image indexing and parsing field name mismatch**

#### Problem
When changing primary image via edit page, the OpenSearch index was not reflecting the change. Root cause: `buildProductDocument()` was writing images with `file_path` field, but `parseImages()` was looking for `public_url` field.

#### Solution
```go
// Before (inconsistent field names)
buildProductDocument() → images[].file_path
parseImages() → looking for public_url → ❌ not found

// After (consistent field names)
buildProductDocument() → images[].public_url ✅
parseImages() → looking for public_url ✅
```

#### Files Changed
- `internal/indexer/listing_indexer.go:266` - Changed `"file_path": img.URL` → `"public_url": img.URL`

#### Result
- ✅ Primary image changes now reflect in OpenSearch immediately
- ✅ Frontend displays correct images from search results
- ✅ No data migration needed (reindexing fixes legacy data)
