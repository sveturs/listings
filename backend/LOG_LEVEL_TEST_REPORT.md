# Отчет о тестировании LOG_LEVEL=warn

## Статус работы LOG_LEVEL=warn
**Частично работает** - INFO сообщения корректно фильтруются, но есть проблемы с DEBUG сообщениями.

## Обнаруженные проблемы

### 1. Неправильное использование уровней логирования
В коде найдены места, где DEBUG сообщения выводятся через `logger.Info()`:
```go
// backend/internal/proj/marketplace/storage/opensearch/repository.go:1232
logger.Info().Msgf("DEBUG: Listing %d has no storefront_id", listing.ID)
```

Это приводит к тому, что сообщения с префиксом "DEBUG:" отображаются даже при LOG_LEVEL=warn.

### 2. Отсутствие инициализации логгера в main.go
В файле `cmd/api/main.go` не вызывается `logger.Init()` для установки уровня логирования из конфигурации. Логгер использует дефолтный уровень из `zerolog.SetGlobalLevel()`.

## Сравнение логов

### При LOG_LEVEL=warn
- ✅ INFO сообщения типа `{"level":"info","message":"REQUEST"}` НЕ отображаются
- ❌ DEBUG сообщения типа `DEBUG: Listing X has no storefront_id` отображаются (из-за logger.Info())
- ✅ Сообщения запуска сервера (Fiber banner) отображаются
- Всего строк логов при запуске: ~15

### При LOG_LEVEL=info  
- ✅ INFO сообщения отображаются
- ✅ DEBUG сообщения отображаются
- ✅ Подробные логи инициализации (~40+ строк)
- Всего строк логов при запуске: ~50+

## Примеры отфильтрованных логов при LOG_LEVEL=warn

Следующие логи НЕ отображаются при LOG_LEVEL=warn:
```json
{"level":"info","gitCommit":"unknown","buildTime":"unknown","config":{...},"time":"2025-07-07T12:28:53+02:00","message":"Config loaded successfully"}
{"level":"info","provider":"minio","time":"2025-07-07T12:28:53+02:00","message":"Инициализация хранилища файлов."}
{"level":"info","time":"2025-07-07T12:28:53+02:00","message":"Успешное подключение к OpenSearch"}
{"level":"info","method":"GET","path":"/api/v1/marketplace/listings","status":200,"duration":13.490388,"time":"2025-07-07T12:24:26+02:00","message":"RESPONSE"}
```

## Оценка улучшения производительности

1. **Уменьшение объема логов**: ~70% меньше логов при LOG_LEVEL=warn
2. **Производительность**: 
   - Меньше операций I/O на запись логов
   - Меньше сериализации JSON объектов
   - Ожидаемое улучшение: 5-10% для высоконагруженных endpoints
3. **Размер лог-файлов**: Существенное уменьшение (в 3-4 раза)

## Рекомендации

1. **Исправить использование уровней логирования**:
   ```go
   // Неправильно:
   logger.Info().Msgf("DEBUG: Listing %d has no storefront_id", listing.ID)
   
   // Правильно:
   logger.Debug().Msgf("Listing %d has no storefront_id", listing.ID)
   ```

2. **Добавить инициализацию логгера в main.go**:
   ```go
   // После загрузки конфигурации
   if err := logger.Init(cfg.Environment, cfg.LogLevel); err != nil {
       log.Fatal("Failed to init logger:", err)
   }
   ```

3. **Провести аудит всех логов** и убедиться, что они используют правильные уровни:
   - Debug: для отладочной информации
   - Info: для важных событий
   - Warn: для предупреждений
   - Error: для ошибок

## Заключение

LOG_LEVEL=warn частично работает, но требует доработки кода для полноценного функционирования. После исправлений можно ожидать существенного улучшения производительности и уменьшения нагрузки на систему логирования.