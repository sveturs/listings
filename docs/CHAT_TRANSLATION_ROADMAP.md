# 🚀 ROADMAP - ПЛАН УЛУЧШЕНИЙ ПЕРЕВОДОВ В ЧАТЕ

**Дата создания:** 2025-10-04
**Автор:** Claude (по результатам аудита)
**Статус:** 🔴 КРИТИЧНЫЕ УЛУЧШЕНИЯ ТРЕБУЮТСЯ

---

## 🔍 РЕЗУЛЬТАТЫ АУДИТА (2025-10-04)

### ❌ КРИТИЧЕСКИЕ ПРОБЛЕМЫ ТЕКУЩЕЙ РЕАЛИЗАЦИИ

#### Проблема #1: Client-side вместо Server-side переводы

**Как работает сейчас:**
```
1. Backend → Frontend: сообщение в оригинале "Привет"
2. Frontend показывает оригинал (~300мс видно пользователю)
3. Frontend запрашивает перевод через API
4. Backend переводит и возвращает "Hello"
5. Frontend заменяет текст (эффект "прыгания")
```

**Проблемы:**
- ❌ Пользователь видит "прыгающие" сообщения
- ❌ При обновлении страницы ~50 API запросов к Claude (по одному на сообщение)
- ❌ Redis кеш существует, но НЕ используется в GetMessages
- ❌ Локаль только в localStorage (не синхронизирована с сервером)
- ❌ WebSocket отправляет сообщения БЕЗ переводов

#### Проблема #2: useEffect вызывает перевод при каждом рендере

**Код в MessageItem.tsx:**
```typescript
useEffect(() => {
  if (autoTranslate && !translatedText) {
    handleTranslate(); // ⚠️ API запрос при каждом монтировании!
  }
}, [message.id, translatedText, handleTranslate]);
```

**Последствия:**
- При открытии чата с 50 сообщениями → 50 API запросов
- При обновлении страницы → снова 50 API запросов
- Redis кеш игнорируется на уровне GetMessages

#### Проблема #3: Локаль не хранится на сервере

**Текущее состояние:**
- `chat_auto_translate` только в localStorage
- Backend НЕ знает preferred_language пользователя
- При WebSocket broadcast невозможно отправить перевод

#### Проблема #4: Redis кеш не используется в GetMessages

```go
// backend/internal/proj/marketplace/handler/chat.go
func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
    messages, _ := h.services.Storage().GetMessages(ctx, params)

    // ❌ Переводы НЕ добавляются!
    // Сообщения возвращаются БЕЗ translations

    return utils.SuccessResponse(c, messages)
}
```

#### Проблема #5: WebSocket без переводов

```go
// WebSocket broadcast
ws.Send(newMessage) // ❌ Только оригинал
```

**Последствия:**
- Новые сообщения всегда показываются в оригинале
- Frontend запрашивает перевод через API (300-500мс лаг)

---

## 🚀 ПЛАН УЛУЧШЕНИЙ (5 ФАЗ, 8-12 ДНЕЙ)

### 🎯 ЦЕЛЕВАЯ АРХИТЕКТУРА (Server-side Translation)

```
┌─────────────────────────────────────────────────────────────┐
│                    USER AUTHENTICATION                      │
│  1. Login → Backend saves preferred_language in DB          │
│  2. Change locale → Frontend syncs to backend               │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                    GET MESSAGES (HTTP)                      │
│  1. Frontend: GET /api/v2/marketplace/chat/messages         │
│  2. Backend: Loads messages from DB                         │
│  3. Backend: Gets user's preferred_language from DB         │
│  4. Backend: TranslateBatch() - checks Redis first!         │
│  5. Backend: Returns messages WITH translations             │
│  6. Frontend: Shows translated text IMMEDIATELY (0ms)       │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                  NEW MESSAGE (WebSocket)                    │
│  1. User A sends: "Привет"                                  │
│  2. Backend detects language: "ru"                          │
│  3. Backend saves with original_language="ru"               │
│  4. Backend broadcasts per participant:                     │
│     - User A (ru): "Привет" (original)                      │
│     - User B (en): "Hello" (translated from Redis/API)      │
│  5. Both users see correct version INSTANTLY                │
└─────────────────────────────────────────────────────────────┘
```

---

## 📋 ДЕТАЛЬНЫЙ ПЛАН РЕАЛИЗАЦИИ

### Phase 1: Backend - Хранение локали пользователя (1-2 дня)

#### Task 1.1: Обновить модели

**Файл:** `backend/internal/domain/models/marketplace_chat.go`

```go
type ChatUserSettings struct {
    AutoTranslate     bool   `json:"auto_translate_chat"`
    PreferredLanguage string `json:"preferred_language"` // "ru", "en", "sr"
    ShowLanguageBadge bool   `json:"show_original_language_badge"`
    ModerateTone      bool   `json:"chat_tone_moderation"` // NEW
}
```

**Миграция БД:** НЕ ТРЕБУЕТСЯ
(`user_privacy_settings.settings` уже JSONB - просто добавляем новые ключи)

#### Task 1.2: Endpoint для управления настройками

**Новые endpoints:**
- `PUT /api/v1/users/chat-settings` - обновить настройки
- `GET /api/v1/users/chat-settings` - получить настройки

**Файл:** `backend/internal/proj/users/handler/user.go`

```go
func (h *UserHandler) UpdateChatSettings(c *fiber.Ctx) error {
    userID, _ := authMiddleware.GetUserID(c)

    var req models.ChatUserSettings
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidRequest")
    }

    // Валидация языка
    if !isValidLanguage(req.PreferredLanguage) {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidLanguage")
    }

    // Сохраняем в user_privacy_settings.settings
    err := h.services.User().UpdateChatSettings(c.Context(), userID, &req)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.updateSettingsFailed")
    }

    return utils.SuccessResponse(c, req)
}
```

#### Task 1.3: Service методы

**Файл:** `backend/internal/proj/users/service/user.go`

```go
func (s *UserService) UpdateChatSettings(ctx context.Context, userID int, settings *models.ChatUserSettings) error {
    settingsJSON, _ := json.Marshal(settings)

    query := `
        INSERT INTO user_privacy_settings (user_id, settings)
        VALUES ($1, $2)
        ON CONFLICT (user_id) DO UPDATE
        SET settings = user_privacy_settings.settings || $2::jsonb
    `

    _, err = s.db.ExecContext(ctx, query, userID, settingsJSON)
    return err
}

func (s *UserService) GetChatSettings(ctx context.Context, userID int) (*models.ChatUserSettings, error) {
    query := `
        SELECT
            settings->'preferred_language' as preferred_language,
            settings->'auto_translate_chat' as auto_translate,
            settings->'show_original_language_badge' as show_badge,
            settings->'chat_tone_moderation' as moderate_tone
        FROM user_privacy_settings
        WHERE user_id = $1
    `

    var settings models.ChatUserSettings
    err := s.db.QueryRowContext(ctx, query, userID).Scan(
        &settings.PreferredLanguage,
        &settings.AutoTranslate,
        &settings.ShowLanguageBadge,
        &settings.ModerateTone,
    )

    if err == sql.ErrNoRows {
        // Возвращаем defaults
        return &models.ChatUserSettings{
            AutoTranslate:     true,
            PreferredLanguage: "en",
            ShowLanguageBadge: true,
            ModerateTone:      true,
        }, nil
    }

    return &settings, err
}
```

---

### Phase 2: Backend - Server-side переводы в GetMessages (2-3 дня)

#### Task 2.1: Обновить GetMessages

**Файл:** `backend/internal/proj/marketplace/handler/chat.go`

```go
func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
    userID, _ := authMiddleware.GetUserID(c)

    // ... existing code для загрузки сообщений ...
    messages, err := h.services.Storage().GetMessages(c.Context(), params)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getMessagesFailed")
    }

    // ✅ НОВЫЙ КОД: Автоматический перевод
    userSettings, err := h.services.User().GetChatSettings(c.Context(), userID)
    if err != nil {
        logger.Warn().Err(err).Msg("Failed to get user chat settings")
    } else if userSettings.AutoTranslate && userSettings.PreferredLanguage != "" {
        logger.Debug().
            Int("userId", userID).
            Str("preferredLang", userSettings.PreferredLanguage).
            Int("messagesCount", len(messages)).
            Msg("Auto-translating messages")

        // Batch перевод с использованием Redis кеша
        err = h.services.ChatTranslation().TranslateBatch(
            c.Context(),
            messages,
            userSettings.PreferredLanguage,
            userSettings.ModerateTone,
        )
        if err != nil {
            logger.Warn().Err(err).Msg("Batch translation failed")
        }
    }

    return utils.SuccessResponse(c, map[string]interface{}{
        "messages": messages,
        "total":    len(messages),
        "page":     params.Page,
        "limit":    params.Limit,
    })
}
```

#### Task 2.2: Обновить ChatTranslationService.GetUserTranslationSettings

**Файл:** `backend/internal/proj/marketplace/service/chat_translation.go`

```go
func (s *ChatTranslationService) GetUserTranslationSettings(
    ctx context.Context,
    userID int,
) (*models.ChatUserSettings, error) {
    // ✅ ИЗМЕНЕНО: Теперь загружаем из БД, а не возвращаем defaults
    return s.userService.GetChatSettings(ctx, userID)
}
```

---

### Phase 3: Frontend - Синхронизация локали (1-2 дня)

#### Task 3.1: API методы

**Файл:** `frontend/svetu/src/services/chat.ts`

```typescript
// Новые методы для работы с настройками
async getChatSettings(): Promise<ChatUserSettings> {
  const response = await this.request<{
    data: ChatUserSettings;
    success: boolean;
  }>('/settings'); // BFF proxy → /api/v1/users/chat-settings
  return response.data;
}

async updateChatSettings(settings: ChatUserSettings): Promise<void> {
  await this.request<void>('/settings', {
    method: 'PUT',
    body: JSON.stringify(settings),
  });
}
```

**Файл:** `frontend/svetu/src/types/chat.ts`

```typescript
export interface ChatUserSettings {
  auto_translate_chat: boolean;
  preferred_language: 'ru' | 'en' | 'sr';
  show_original_language_badge: boolean;
  chat_tone_moderation: boolean; // NEW
}
```

#### Task 3.2: Синхронизация локали при смене языка

**Файл:** `frontend/svetu/src/app/[locale]/layout.tsx`

```typescript
'use client';

import { useLocale } from 'next-intl';
import { useEffect } from 'react';
import { chatService } from '@/services/chat';

export default function LocaleLayout({ children }: { children: React.ReactNode }) {
  const locale = useLocale();

  useEffect(() => {
    const syncLocale = async () => {
      try {
        const currentSettings = await chatService.getChatSettings();

        // Синхронизируем локаль с сервером
        if (currentSettings.preferred_language !== locale) {
          await chatService.updateChatSettings({
            ...currentSettings,
            preferred_language: locale as 'ru' | 'en' | 'sr',
          });

          console.log(`Locale synced: ${locale}`);
        }
      } catch (error) {
        console.error('Failed to sync locale:', error);
      }
    };

    syncLocale();
  }, [locale]);

  return <>{children}</>;
}
```

#### Task 3.3: Упростить MessageItem

**Файл:** `frontend/svetu/src/components/Chat/MessageItem.tsx`

```typescript
export default function MessageItem({ message, isOwn }: MessageItemProps) {
  const locale = useLocale();
  const t = useTranslations('chat');

  const [showOriginal, setShowOriginal] = useState(false);

  // ✅ ИЗМЕНЕНО: Проверяем готовый перевод из backend
  const hasTranslation = message.translations && message.translations[locale];

  // ✅ ИЗМЕНЕНО: Просто показываем (БЕЗ API запроса!)
  const displayText = showOriginal
    ? message.content
    : (hasTranslation || message.content);

  // ✅ УДАЛЕНО: useEffect с handleTranslate()
  // ✅ УДАЛЕНО: isTranslating, translatedText states
  // ✅ УДАЛЕНО: API запросы из компонента

  const shouldShowToggleButton = !isOwn && hasTranslation;

  return (
    <div className={`chat ${isOwn ? 'chat-end' : 'chat-start'} mb-2`}>
      <div className="chat-bubble">
        <p className="whitespace-pre-wrap">{displayText}</p>
      </div>

      {/* Кнопка только для переключения (БЕЗ API!) */}
      {shouldShowToggleButton && (
        <button
          onClick={() => setShowOriginal(!showOriginal)}
          className="btn btn-xs btn-ghost"
        >
          {showOriginal ? t('translation.showTranslation') : t('translation.showOriginal')}
        </button>
      )}
    </div>
  );
}
```

---

### Phase 4: Backend - WebSocket с переводами (1-2 дня)

#### Task 4.1: Обновить broadcast сообщений

**Файл:** `backend/internal/proj/marketplace/handler/websocket.go`

```go
func (h *ChatHandler) broadcastMessageToParticipants(ctx context.Context, message *models.MarketplaceMessage) {
    // Получаем участников чата
    chat, err := h.services.Storage().GetChatByID(ctx, message.ChatID)
    if err != nil {
        logger.Error().Err(err).Msg("Failed to get chat")
        return
    }

    participants := []int{chat.BuyerID, chat.SellerID}

    for _, participantID := range participants {
        // Клонируем сообщение для каждого участника
        msgCopy := *message

        // Получаем настройки участника из БД
        settings, err := h.services.User().GetChatSettings(ctx, participantID)
        if err != nil {
            logger.Warn().Err(err).Int("userId", participantID).Msg("Failed to get settings")
            settings = &models.ChatUserSettings{
                AutoTranslate:     false,
                PreferredLanguage: "en",
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

        // Отправляем сообщение С переводом
        h.sendToUser(participantID, &msgCopy)
    }
}
```

---

### Phase 5: Оптимизация и тестирование (2-3 дня)

#### Task 5.1: Unit tests

**Файлы для тестирования:**
- `backend/internal/proj/marketplace/service/chat_translation_test.go`
- `backend/internal/proj/users/service/user_test.go`
- `frontend/svetu/src/services/chat.test.ts`

**Тест-кейсы:**
1. UpdateChatSettings сохраняет в БД
2. GetChatSettings загружает из БД
3. TranslateBatch использует Redis кеш
4. WebSocket broadcast переводит для каждого участника

#### Task 5.2: Integration tests

**Тест:** GetMessages с переводами

```go
func TestGetMessages_WithAutoTranslation(t *testing.T) {
    // 1. Создать 2 пользователей (RU и EN)
    // 2. User RU отправляет сообщение "Привет"
    // 3. User EN делает GET /messages
    // 4. Проверить: message.translations["en"] == "Hello"
    // 5. Проверить: metadata.cache_hit == true (при повторном запросе)
}
```

#### Task 5.3: E2E tests

**Сценарий:** User1 (RU) → User2 (EN) через WebSocket

```typescript
test('should translate WebSocket messages', async () => {
  // 1. User1 (ru) connects to WebSocket
  // 2. User2 (en) connects to WebSocket
  // 3. User1 sends "Привет"
  // 4. User2 receives message with translations["en"] == "Hello"
  // 5. User2 sees translated text immediately (no delay)
});
```

#### Task 5.4: Load testing

**Цель:** Проверить Redis cache hit rate

```bash
# Симуляция 1000 пользователей, 50 сообщений каждый
# Ожидаемый cache hit rate: >80%
```

---

## ✅ ЧЕКЛИСТ ВНЕДРЕНИЯ

### Backend

**Phase 1: Хранение локали**
- [ ] Обновить модель ChatUserSettings (добавить ModerateTone)
- [ ] Endpoint PUT `/api/v1/users/chat-settings`
- [ ] Endpoint GET `/api/v1/users/chat-settings`
- [ ] Service методы UpdateChatSettings / GetChatSettings
- [ ] Unit tests для UserService

**Phase 2: Server-side переводы**
- [ ] Обновить GetMessages - автоматический batch перевод
- [ ] Обновить ChatTranslationService.GetUserTranslationSettings (использовать БД)
- [ ] Integration tests для GetMessages с переводами

**Phase 4: WebSocket**
- [ ] Обновить broadcast - переводы per participant
- [ ] E2E tests для WebSocket с переводами

### Frontend

**Phase 3: Синхронизация локали**
- [ ] chatService методы getChatSettings / updateChatSettings
- [ ] Синхронизация локали при смене языка
- [ ] Обновить MessageItem - убрать useEffect, показ готовых переводов
- [ ] Обновить ChatSettings - синхронизация с сервером

### Testing

**Phase 5: Тестирование**
- [ ] Unit tests: ChatTranslationService
- [ ] Unit tests: UserService (chat settings)
- [ ] Integration tests: GetMessages с переводами
- [ ] E2E tests: User RU → User EN (WebSocket + HTTP)
- [ ] Load testing: Redis cache hit rate >80%
- [ ] Manual testing: обновление страницы, смена локали

---

## 🎯 ОЖИДАЕМЫЕ РЕЗУЛЬТАТЫ

### До улучшений (текущее состояние)

- ❌ При обновлении страницы: **~50 API запросов** к Claude
- ❌ Пользователь видит "прыгание": оригинал → перевод (**~300мс задержка**)
- ❌ Redis кеш существует, но **НЕ используется** в GetMessages
- ❌ Локаль только в localStorage (**не работает на других устройствах**)
- ❌ WebSocket сообщения приходят БЕЗ переводов (**300-500мс лаг**)

### После улучшений

- ✅ При обновлении страницы: **0 API запросов** (все из Redis/БД)
- ✅ Пользователь видит перевод **СРАЗУ (0мс задержка)**
- ✅ Redis кеш **используется эффективно** в GetMessages и WebSocket
- ✅ Локаль **синхронизирована на сервере** (работает на всех устройствах)
- ✅ WebSocket сообщения приходят **УЖЕ с переводами** (мгновенно)
- ✅ При новом сообщении: **1 API запрос** → кеш → broadcast всем

### Экономия ресурсов

**Текущая реализация:**
- Chat со 100 сообщениями, 10 пользователей
- При каждом обновлении страницы: 100 сообщений × 10 пользователей = **1000 API запросов**

**После улучшений:**
- При первом открытии чата: 100 сообщений → **100 API запросов** (только если нет в Redis)
- При повторном открытии: **0 API запросов** (все из Redis)
- **Экономия: 90-95% API запросов**

### Улучшение UX

**Текущий UX:**
1. Открыть чат → вижу "Привет" (300мс)
2. Текст меняется на "Hello" (прыгание)
3. Обновить страницу → снова "Привет" → "Hello"

**Новый UX:**
1. Открыть чат → вижу "Hello" (0мс, сразу!)
2. Нет прыгания
3. Обновить страницу → снова "Hello" (0мс, из кеша)

---

## 📊 МЕТРИКИ УСПЕХА

После внедрения нужно отслеживать:

1. **Redis cache hit rate**
   - Цель: >80%
   - Prometheus: `chat_translation_cache_hit_rate`

2. **Translation latency**
   - Цель: <100ms (с кешем)
   - Prometheus: `chat_translation_duration_seconds`

3. **API requests count**
   - Цель: снижение на 90-95%
   - Prometheus: `chat_translation_requests_total`

4. **User experience**
   - Цель: 0 "прыгающих" сообщений
   - A/B тестирование: старая vs новая архитектура

---

## 📚 СВЯЗАННЫЕ ДОКУМЕНТЫ

- [Выполненная функциональность](./CHAT_TRANSLATION_COMPLETED.md) - что уже реализовано
- [Исходный план](./CHAT_TRANSLATION_IMPLEMENTATION_PLAN.md) - полный первоначальный план

---

## 🚨 КРИТИЧНОСТЬ УЛУЧШЕНИЙ

**Приоритет:** 🔴 ВЫСОКИЙ

**Причины:**
1. Каждое обновление страницы = десятки лишних API запросов
2. Плохой UX - пользователи видят "прыгающие" сообщения
3. Redis кеш не используется эффективно
4. Локаль не синхронизирована между устройствами

**Рекомендуемый срок реализации:** 2-3 недели

---

**Дата последнего обновления:** 2025-10-04
**Автор аудита:** Claude
**Статус:** 🔴 ОЖИДАЕТ РЕАЛИЗАЦИИ
