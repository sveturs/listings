# Задача 022: Рефакторинг файлов витрин

## Дата: 2025-06-18

## Описание
Произведен рефакторинг структуры файлов витрин (storefronts) в backend:
1. Удалены старые файлы модели и репозитория
2. Переименованы новые файлы с удалением суффикса "_new"
3. Обновлены все импорты в коде

## Выполненные действия

### 1. Удаление старых файлов
- ✅ Удален `/backend/internal/domain/models/storefront.go` (старая модель)
- ✅ Удален `/backend/internal/storage/postgres/storefront.go` (старый репозиторий)
- ✅ Удален модуль `/backend/internal/proj/storefront/` целиком

### 2. Переименование новых файлов
- ✅ `storefront_new.go` → `storefront.go` (модель)
- ✅ `storefront_new.go` → `storefront.go` (репозиторий)
- ✅ `storefront_new_methods.go` → `storefront_methods.go`

### 3. Обновление кода
- ✅ Заменены все упоминания `StorefrontNew` на `Storefront`
- ✅ Обновлены импорты с `storefront/` на `storefronts/`
- ✅ Добавлена реализация методов Storefront в Database struct
- ✅ Исправлены типы данных (float32 → float64 для rating)
- ✅ Удалены неиспользуемые модели ImportSource и ImportHistory
- ✅ Создан основной handler.go для модуля storefronts
- ✅ Добавлены методы UploadLogo и UploadBanner

### 4. Исправленные проблемы
- ✅ Удален дублирующийся MapCluster из storefront.go
- ✅ Добавлен ErrNotFound в storefront.go
- ✅ Исправлен импорт с minio на filestorage
- ✅ Обновлены вызовы Upload на UploadFile
- ✅ Убраны упоминания scheduleService
- ✅ Создана правильная структура Handler с методами GetPrefix и RegisterRoutes

## Структура после рефакторинга
```
backend/
├── internal/
│   ├── domain/models/
│   │   └── storefront.go (новая объединенная модель)
│   ├── storage/postgres/
│   │   ├── storefront.go (интерфейс репозитория)
│   │   └── storefront_methods.go (методы репозитория)
│   └── proj/storefronts/
│       ├── handler/
│       │   ├── handler.go (основной handler)
│       │   ├── storefront_handler.go
│       │   └── staff_handler.go
│       └── service/
│           └── storefront_service.go
```

## Статус
✅ Задача выполнена успешно
✅ Компиляция backend прошла без ошибок