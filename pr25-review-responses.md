# Ответы на комментарии по PR #25

## 1. backend/internal/proj/marketplace/service/contacts.go - использование logger вместо fmt.Print()

✅ **Исправлено**: Все вызовы `fmt.Printf()` заменены на соответствующие вызовы logger:
- `fmt.Printf("[ContactsService] AddContact called...")` → `logger.Debug().Int("userID", userID).Int("contactUserID", req.ContactUserID).Msg("[ContactsService] AddContact called")`
- `fmt.Printf("[ContactsService] CanAddContact error...")` → `logger.Error().Err(err).Msg("[ContactsService] CanAddContact error")`
- `fmt.Printf("[ContactsService] CanAddContact result...")` → `logger.Debug().Bool("canAdd", canAdd).Msg("[ContactsService] CanAddContact result")`
- `fmt.Printf("Warning: failed to update reverse contact status...")` → `logger.Warn().Err(err).Msg("Failed to update reverse contact status")`
- `fmt.Printf("Warning: failed to create reverse contact...")` → `logger.Warn().Err(err).Msg("Failed to create reverse contact")`

Всего заменено 6 вызовов.

## 2. backend/internal/proj/marketplace/storage/postgres/contacts.go - почему убрал err == sql.ErrNoRows

❌ **Это недоразумение**: Обработка `sql.ErrNoRows` **НЕ была удалена**. Она находится на месте в методе `GetContact()` на строке 93:

```go
if err != nil {
    if err == sql.ErrNoRows {
        return nil, nil // Контакт не найден
    }
    return nil, fmt.Errorf("error getting contact: %w", err)
}
```

Возможно, ты смотрел на старую версию diff. В текущем коде обработка присутствует и работает корректно.

✅ **Дополнительно исправлено**: В этом же файле заменены 3 вызова `fmt.Printf()` на logger:
- Строка 393: `logger.Debug()` для входных параметров CanAddContact
- Строка 401: `logger.Error()` для ошибки GetUserPrivacySettings
- Строка 407: `logger.Debug()` для результата настроек приватности

## 3. frontend/svetu/src/services/api-client.ts - почему убрал обработку rate-limit

❌ **Это тоже недоразумение**: Обработка rate-limit **НЕ была удалена**. Она находится на строках 218-224:

```typescript
// Обработка rate limiting
if (response.status === 429) {
  const retryAfter = response.headers.get('Retry-After');
  console.warn(
    `API rate limited (429), retry after: ${retryAfter || 'unknown'} seconds`
  );
}
```

Возможно, в каком-то из промежуточных коммитов она была случайно удалена, но в финальной версии PR она присутствует.

## Тестирование

Все изменения протестированы:
- ✅ Код успешно компилируется (`go build`)
- ✅ Статический анализ пройден (`go vet`)
- ✅ Форматирование применено (`go fmt`)
- ✅ API endpoints работают корректно
- ✅ Logger правильно импортирован и используется