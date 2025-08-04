# Реальные проблемы с переводами

## 1. Модуль storefront не загружается

**Проблема**: Страницы используют `useTranslations('storefront')` и `useTranslations('storefronts')`, но эти модули не загружаются.

**Файлы с проблемой**:
- `/profile/storefronts/page.tsx` - использует `storefronts`
- Все компоненты витрин

**Решение**: Добавить загрузку модуля в соответствующих местах.

## 2. Неправильное обращение к ключам authError

**Проблема**: Компоненты ищут `common.errors.authError.*`, но эти ключи находятся в `auth.errors.*`

**Ошибки в консоли**:
```
IntlError: MISSING_MESSAGE: Could not resolve `common.errors.authError.title`
IntlError: MISSING_MESSAGE: Could not resolve `common.errors.authError.description`
IntlError: MISSING_MESSAGE: Could not resolve `common.errors.authError.details`
IntlError: MISSING_MESSAGE: Could not resolve `common.errors.authError.reload`
```

## 3. Отсутствующие ключи в common

**Проблема**: Используются ключи, которых нет в common.json:
- `common.all`
- `common.reviews`
- `common.notSpecified`

## 4. Неправильное использование вложенных namespace

**Проблема**: 91 файл использует паттерн типа `useTranslations('module.submodule')`, что не поддерживается.

**Примеры**:
- `useTranslations('storefronts.products')`
- `useTranslations('storefronts.business_types')`
- `useTranslations('create_storefront.staff_setup')`

## 5. Snake_case ключи

**Проблема**: 2052 ключа используют snake_case вместо camelCase.

**Примеры**:
- `fetch_failed` → `fetchFailed`
- `user_not_found` → `userNotFound`
- `validation_failed` → `validationFailed`

## Итоговая статистика

- **Критичных проблем**: ~100 (отсутствующие модули и ключи)
- **Средних проблем**: ~91 (неправильные namespace)
- **Низких проблем**: ~2052 (snake_case)
- **Ложных срабатываний**: ~6500+ (дубликаты из-за ошибки в скрипте)

**Реальное количество проблем**: ~2243 (а не 8846)