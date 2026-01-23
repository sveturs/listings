# CHANGELOG - Listings Service

## [Unreleased]

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
