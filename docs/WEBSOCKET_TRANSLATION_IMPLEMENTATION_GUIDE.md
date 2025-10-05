# WebSocket Translation Implementation Guide

## Статус: ✅ РЕАЛИЗОВАНО (2025-10-05)

## Описание

Server-side переводы для WebSocket broadcast сообщений успешно реализованы.

## Текущая проблема

- WebSocket отправляет сообщения БЕЗ переводов
- Frontend должен запрашивать перевод через отдельный API запрос (300-500ms задержка)
- Нет персонализации по языку для каждого участника чата

## Целевое решение

WebSocket broadcast должен отправлять каждому участнику персонализированную версию сообщения:
- User A (ru preference) получает "Привет"
- User B (en preference) получает "Hello"
- User C (sr preference) получает Serbian перевод

## Реализация

### 1. Найти WebSocket broadcast функцию

Необходимо найти где происходит broadcast новых сообщений (вероятно в одном из файлов):
- `backend/internal/proj/marketplace/service/chat.go`
- `backend/internal/proj/marketplace/handler/chat.go`
- Отдельный WebSocket handler file

### 2. Создать функцию broadcastMessageToParticipants

```go
// backend/internal/proj/marketplace/handler/chat.go или websocket.go

func (h *ChatHandler) broadcastMessageToParticipants(ctx context.Context, message *models.MarketplaceMessage) {
    // 1. Получить участников чата
    chat, err := h.services.Storage().GetChatByID(ctx, message.ChatID)
    if err != nil {
        logger.Error().Err(err).Msg("Failed to get chat")
        return
    }

    participants := []int{chat.BuyerID, chat.SellerID}

    // 2. Для каждого участника - персонализированное сообщение
    for _, participantID := range participants {
        // Клонируем сообщение
        msgCopy := *message

        // Получаем настройки участника
        settings, err := h.services.User().GetChatSettings(ctx, participantID)
        if err != nil {
            logger.Warn().Err(err).Int("userId", participantID).Msg("Failed to get settings")
            settings = &models.ChatUserSettings{
                AutoTranslate:     false,
                PreferredLanguage: "en",
                ModerateTone:      true,
            }
        }

        // Если нужен перевод
        if settings.AutoTranslate &&
           msgCopy.OriginalLanguage != settings.PreferredLanguage {

            // Переводим (используя Redis кеш!)
            err = h.services.ChatTranslation().TranslateMessage(
                ctx,
                &msgCopy,
                settings.PreferredLanguage,
                settings.ModerateTone,
            )
            if err != nil {
                logger.Warn().Err(err).Msg("WebSocket translation failed")
            }
        }

        // Отправляем персонализированное сообщение
        h.sendToUser(participantID, &msgCopy)
    }
}
```

### 3. Интеграция

Заменить старый broadcast:
```go
// ❌ СТАРЫЙ КОД
ws.Send(newMessage) // Отправка всем одинакового сообщения

// ✅ НОВЫЙ КОД
h.broadcastMessageToParticipants(ctx, newMessage)
```

## Преимущества

1. ✅ **Мгновенный перевод**: Каждый участник получает сообщение на своем языке сразу
2. ✅ **Использование кеша**: TranslateMessage проверяет Redis кеш перед API запросом
3. ✅ **Персонализация**: Учитываются индивидуальные настройки каждого пользователя
4. ✅ **Нет лишних запросов**: Frontend не нужно делать отдельный API вызов для перевода

## Зависимости

**Уже реализовано:**
- ✅ ChatUserSettings model (с ModerateTone)
- ✅ User.GetChatSettings() service method
- ✅ ChatTranslation.TranslateMessage() service method
- ✅ Redis кеш для переводов

**Требуется найти/реализовать:**
- 🔍 WebSocket broadcast функцию
- 🔍 Метод `sendToUser(userID, message)` для отправки персонализированных сообщений

## Тестирование

После реализации протестировать:
1. User RU отправляет "Привет" → User EN получает "Hello" через WebSocket
2. Проверить Redis cache hit rate (должен быть >80%)
3. Проверить отсутствие дополнительных API запросов от frontend

## Файлы для проверки

1. Найти где определен WebSocket handler:
   - `backend/internal/server/router.go` (поиск `/ws` route)
   - `backend/internal/proj/marketplace/handler/` (websocket.go или chat.go)

2. Найти где происходит broadcast новых сообщений

## Статус реализации

- ✅ Phase 1-3: Backend settings, Frontend sync - ЗАВЕРШЕНО
- ✅ Phase 4: WebSocket broadcast - РЕАЛИЗОВАНО (2025-10-05)
- ⏳ Phase 5: Testing - ТРЕБУЕТ ТЕСТИРОВАНИЯ

---

## ✅ Детали реализации (2025-10-05)

### Изменённые файлы:

**Backend:**
```
backend/internal/proj/marketplace/service/chat.go
  - Добавлены поля chatTranslationSvc и userService в структуру ChatService
  - Добавлены методы SetChatTranslationService() и SetUserService()
  - Создана функция BroadcastMessageWithTranslations() (строка 259)
  - Заменён вызов BroadcastMessage() на BroadcastMessageWithTranslations() в SendMessage()

backend/internal/proj/marketplace/service/chat_interface.go
  - Добавлен метод BroadcastMessageWithTranslations() в интерфейс

backend/internal/proj/marketplace/handler/chat.go
  - Заменён вызов BroadcastMessage() на BroadcastMessageWithTranslations() в UploadAttachments()

backend/internal/proj/global/service/service.go
  - Добавлена установка зависимостей chatTranslationSvc и userService в ChatService (строка 147-152)
```

### Как работает:

1. **При отправке сообщения** (`SendMessage` или `UploadAttachments`):
   - Создаётся сообщение в БД
   - Вызывается `BroadcastMessageWithTranslations(ctx, msg)`

2. **BroadcastMessageWithTranslations** делает:
   - Получает список участников (sender и receiver)
   - Для каждого участника:
     - Клонирует сообщение
     - Загружает `GetChatSettings()` из БД
     - Если `AutoTranslate == true` И язык отличается - переводит через `TranslateMessage()`
     - Отправляет персонализированную версию через существующий механизм subscribers

3. **TranslateMessage** использует:
   - Redis кеш (30 дней TTL)
   - Модерацию тона если включена
   - Fallback на оригинал при ошибке

### Преимущества реализации:

✅ **0ms задержка на frontend** - переводы приходят сразу с WebSocket
✅ **Redis кеш** - повторные переводы берутся из кеша
✅ **Персонализация** - каждый участник получает сообщение на своём языке
✅ **Fallback** - если зависимости не установлены, используется старый BroadcastMessage
✅ **Модерация тона** - учитываются настройки пользователя

### Логи подтверждения:

```
2025/10/05 15:58:58 ChatService dependencies set (translation & user service)
```

Backend успешно запущен и работает с новой реализацией.

---

**Дата создания:** 2025-10-04
**Дата реализации:** 2025-10-05
**Автор:** Claude
**Related:** CHAT_TRANSLATION_ROADMAP.md, CHAT_TRANSLATION_IMPLEMENTATION_SUMMARY.md
