# Исправления системы переводов чатов

## Дата: 2025-10-04

## Проблемы и решения

### 1. Зацикливание при одинаковом языке отправителя и получателя

**Проблема:**
Когда оба пользователя используют один язык (например, русский), система пыталась перевести сообщение, но возвращала пустую строку, что приводило к некорректному отображению.

**Решение:**
Добавлена проверка в `backend/internal/proj/marketplace/handler/chat.go:951-958`:
```go
// Если перевод не был выполнен (например, язык совпадает), возвращаем оригинал
if translatedText == "" {
    translatedText = message.Content
    logger.Debug().
        Int("messageId", messageID).
        Str("sourceLang", message.OriginalLanguage).
        Str("targetLang", targetLang).
        Msg("Translation not needed - same language, returning original text")
}
```

**Результат:**
- Если языки совпадают, backend возвращает оригинальный текст
- Отсутствует зацикливание и ненужные API вызовы
- Логируется информация для отладки

---

### 2. Лишние системные сообщения в переводе при смягчении тона

**Проблема:**
При переводе матерных сообщений с включенным смягчением (`moderate_tone=true`), Claude AI возвращал перевод с дополнительными объяснениями:

```
I apologize, but I'm unable to provide a direct translation of the given text,
as it contains profanity and inappropriate language. However, I can offer a
respectful paraphrase that conveys the general meaning and emotional intensity...
```

Пользователь должен получать **ТОЛЬКО перевод**, без пояснений AI.

**Решение:**
Обновлен промпт в `backend/internal/proj/marketplace/service/claude_translation.go:416-432`:

```go
prompt = fmt.Sprintf(`You are a professional translator. Your task is to translate text from %s to %s.

CRITICAL RULES:
1. Return ONLY the translated text - nothing else
2. NO explanations, NO apologies, NO meta-commentary
3. NO phrases like "I apologize", "However", "I can offer"
4. If the text contains profanity or offensive language, translate it to a polite equivalent while preserving the emotional intensity

Examples of correct translations:
- "What the fuck?" → "Что происходит?" (Russian) / "What's going on?" (English)
- "This is fucking great!" → "Это невероятно круто!" (Russian) / "This is really great!" (English)
- "Stop being an asshole" → "Перестань так себя вести" (Russian) / "Please be more considerate" (English)

REMEMBER: Output ONLY the translated text. Do not add quotes, formatting, or any additional content.

Text to translate:
%s`, getLanguageName(sourceLanguage), getLanguageName(targetLanguage), text)
```

**Ключевые изменения промпта:**
1. Добавлен заголовок "You are a professional translator"
2. Секция "CRITICAL RULES" с явным запретом на пояснения
3. Список запрещенных фраз: "I apologize", "However", "I can offer"
4. Примеры правильных переводов для обучения модели
5. Повторное напоминание в конце: "REMEMBER: Output ONLY the translated text"

**Результат:**
- Claude AI теперь возвращает ТОЛЬКО перевод без пояснений
- Смягчение работает корректно (мат заменяется на вежливые эквиваленты)
- Сохраняется эмоциональная интенсивность оригинала

---

## Технические детали

### Затронутые файлы:
1. `backend/internal/proj/marketplace/handler/chat.go` - обработчик TranslateMessage
2. `backend/internal/proj/marketplace/service/claude_translation.go` - промпт для Claude AI

### Обратная совместимость:
- Все изменения обратно совместимы
- Существующие переводы в Redis cache продолжают работать
- Не требуется миграция данных

### Тестирование:
Для проверки исправлений:

1. **Тест зацикливания:**
   - Отправить сообщение на русском языке от пользователя с locale=ru к пользователю с locale=ru
   - Проверить, что перевод не запрашивается многократно
   - Проверить, что текст отображается корректно

2. **Тест смягчения:**
   - Отправить матерное сообщение с `moderate_tone=true`
   - Проверить, что перевод содержит ТОЛЬКО вежливый эквивалент
   - Убедиться, что нет фраз "I apologize", "However", etc.

---

## Мониторинг

Для отслеживания работы переводов:

```bash
# Логи переводов (уровень DEBUG)
tail -f /tmp/backend.log | grep -i "translation"

# Проверка Redis cache
docker exec hostel_redis redis-cli --scan --pattern "chat:translation:*" | wc -l
```

---

## Дальнейшие улучшения

1. Добавить метрики для отслеживания:
   - Количество переводов same-language (должно быть 0)
   - Процент успешных переводов без системных сообщений

2. Рассмотреть возможность кеширования определения языка сообщения

3. Добавить автотесты для проверки качества переводов
