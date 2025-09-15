# 📋 План улучшения системы чатов

## 📊 Результаты аудита

### ✅ Сильные стороны
- Полнофункциональная WebSocket реализация
- Хорошая архитектура (handler → service → storage)
- Поддержка мультиязычности
- Адаптивный дизайн
- Оптимизированная структура БД с индексами

### ⚠️ Выявленные проблемы

#### 🔴 Критические
1. **Смешанная стилизация** - используются DaisyUI + кастомные CSS + чистый Tailwind
2. **Отсутствие Error Boundaries** в компонентах чата
3. **Memory leaks** в WebSocket при быстром переподключении
4. **Нет ограничений** на размер загружаемых файлов

#### 🟡 Средние
5. **Дублирование логики** между useChat hook и Redux slice
6. **Неэффективная пагинация** без виртуализации для больших списков
7. **Отсутствие rate limiting** для WebSocket сообщений
8. **Избыточные API вызовы** при загрузке контактов

## 🎯 План миграции UI на DaisyUI

### Фаза 1: Унификация базовых компонентов (2-3 дня)

#### 1. ChatWindow.tsx
```tsx
// Было (смешанные стили):
<div className="p-3 sm:p-4 border-b border-base-300 bg-base-100">
  <button className="btn btn-ghost btn-sm btn-circle">

// Станет (чистый DaisyUI):
<div className="navbar bg-base-100 border-b">
  <div className="navbar-start">
    <button className="btn btn-ghost btn-circle">
```

#### 2. MessageItem.tsx
```tsx
// Было (кастомные классы):
<div className={`message-bubble ${isOwn ? 'own' : 'other'}`}>

// Станет (DaisyUI chat):
<div className={`chat ${isOwn ? 'chat-end' : 'chat-start'}`}>
  <div className="chat-bubble chat-bubble-primary">
    {message.content}
  </div>
  <div className="chat-footer opacity-50">
    {formatTime(message.created_at)}
  </div>
</div>
```

#### 3. ChatList.tsx
```tsx
// Было:
<div className="chat-list-item hover:bg-gray-50">

// Станет:
<div className="card card-compact hover:bg-base-200 transition-colors cursor-pointer">
  <div className="card-body">
    <div className="flex items-center gap-3">
      <div className="avatar">
        <div className="w-12 rounded-full">
          <img src={avatar} />
        </div>
      </div>
      <div className="flex-1">
        <h3 className="card-title text-sm">{title}</h3>
        <p className="text-xs opacity-70">{lastMessage}</p>
      </div>
      {unreadCount > 0 && (
        <div className="badge badge-primary">{unreadCount}</div>
      )}
    </div>
  </div>
</div>
```

#### 4. MessageInput.tsx
```tsx
// Было:
<div className="input-group">
  <input className="input input-bordered flex-1">
  <button className="btn btn-primary">

// Станет:
<div className="join w-full">
  <input className="input input-bordered join-item flex-1" />
  <button className="btn btn-primary join-item">
    <svg className="w-5 h-5">...</svg>
  </button>
</div>
```

### Фаза 2: Улучшение UX компонентов (2 дня)

#### 5. EmojiPicker.tsx
```tsx
// Использовать DaisyUI dropdown:
<div className="dropdown dropdown-top dropdown-end">
  <label tabIndex={0} className="btn btn-ghost btn-circle">
    😊
  </label>
  <div tabIndex={0} className="dropdown-content card card-compact w-64 p-2 shadow bg-base-100">
    <div className="grid grid-cols-8 gap-1">
      {emojis.map(emoji => (
        <button className="btn btn-ghost btn-sm">{emoji}</button>
      ))}
    </div>
  </div>
</div>
```

#### 6. FileUploadProgress.tsx
```tsx
// DaisyUI progress:
<div className="card bg-base-200">
  <div className="card-body p-3">
    <div className="flex items-center gap-3">
      <span className="text-sm">{fileName}</span>
      <progress className="progress progress-primary flex-1" value={progress} max="100" />
      <span className="text-xs">{progress}%</span>
    </div>
  </div>
</div>
```

#### 7. ChatAttachments.tsx
```tsx
// DaisyUI carousel для изображений:
<div className="carousel carousel-center max-w-md space-x-2">
  {images.map((img, idx) => (
    <div key={idx} className="carousel-item">
      <div className="card card-compact">
        <figure className="px-2 pt-2">
          <img src={img} className="rounded-lg" />
        </figure>
      </div>
    </div>
  ))}
</div>
```

### Фаза 3: Удаление кастомных стилей (1 день)

#### Файлы для удаления:
- `/src/styles/chat-bubble.css`
- `/src/styles/chat-patterns.css`
- `/src/styles/chat-animations.css`

#### Замена паттернов фона:
```tsx
// Было (кастомный CSS):
.chat-pattern-hexagon {
  background-image: url("data:image/svg+xml...");
}

// Станет (DaisyUI utilities):
<div className="bg-base-200 bg-opacity-50">
  {/* или использовать hero-pattern из DaisyUI */}
  <div className="hero-overlay bg-opacity-10"></div>
</div>
```

## 🚀 План технических улучшений

### Приоритет 1: Критические исправления (1 неделя)

#### 1. Error Boundaries
```tsx
// components/Chat/ChatErrorBoundary.tsx
export class ChatErrorBoundary extends Component {
  state = { hasError: false };

  static getDerivedStateFromError(error) {
    return { hasError: true };
  }

  componentDidCatch(error, errorInfo) {
    console.error('Chat error:', error, errorInfo);
    // Отправить в систему мониторинга
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="alert alert-error">
          <svg className="stroke-current shrink-0 h-6 w-6">...</svg>
          <span>Произошла ошибка в чате. Попробуйте обновить страницу.</span>
        </div>
      );
    }
    return this.props.children;
  }
}
```

#### 2. WebSocket Memory Leak Fix
```tsx
// hooks/useChat.ts
useEffect(() => {
  let ws: WebSocket | null = null;
  let reconnectTimeout: NodeJS.Timeout;
  let isComponentMounted = true;

  const connect = () => {
    if (!isComponentMounted) return;

    ws = new WebSocket(wsUrl);
    // ... handlers
  };

  return () => {
    isComponentMounted = false;
    clearTimeout(reconnectTimeout);
    ws?.close();
  };
}, []);
```

#### 3. File Upload Limits
```go
// backend/internal/proj/marketplace/handler/chat.go
const MaxFileSize = 10 * 1024 * 1024 // 10MB

func (h *Handler) UploadAttachment(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return err
    }

    if file.Size > MaxFileSize {
        return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File too large")
    }
    // ...
}
```

### Приоритет 2: Оптимизация производительности (1 неделя)

#### 4. Виртуализация списков
```tsx
// Установить: yarn add react-window
import { VariableSizeList } from 'react-window';

const VirtualizedMessages = ({ messages }) => {
  const getItemSize = (index) => {
    // Вычислить высоту сообщения
    const msg = messages[index];
    return msg.content.length > 100 ? 120 : 80;
  };

  return (
    <VariableSizeList
      height={600}
      itemCount={messages.length}
      itemSize={getItemSize}
      width="100%"
    >
      {({ index, style }) => (
        <div style={style}>
          <MessageItem message={messages[index]} />
        </div>
      )}
    </VariableSizeList>
  );
};
```

#### 5. Rate Limiting
```go
// backend/internal/middleware/ratelimit.go
func ChatRateLimit() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        10,
        Expiration: 60 * time.Second,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.Locals("user_id").(string)
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(429).JSON(fiber.Map{
                "error": "Too many messages, please slow down",
            })
        },
    })
}
```

#### 6. Redis кеширование
```go
// backend/internal/proj/marketplace/service/chat.go
func (s *Service) GetOnlineUsers(ctx context.Context) ([]int, error) {
    // Проверить кеш
    cached, err := s.redis.Get(ctx, "online_users").Result()
    if err == nil {
        var users []int
        json.Unmarshal([]byte(cached), &users)
        return users, nil
    }

    // Получить из БД
    users := s.getOnlineUsersFromDB()

    // Сохранить в кеш на 30 секунд
    data, _ := json.Marshal(users)
    s.redis.Set(ctx, "online_users", data, 30*time.Second)

    return users, nil
}
```

### Приоритет 3: Дополнительные возможности (2 недели)

#### 7. Поддержка Markdown
```tsx
// Установить: yarn add react-markdown
import ReactMarkdown from 'react-markdown';

const MessageContent = ({ content }) => (
  <div className="chat-bubble">
    <ReactMarkdown
      components={{
        p: ({children}) => <p className="mb-2">{children}</p>,
        code: ({children}) => <code className="badge badge-neutral">{children}</code>,
      }}
    >
      {content}
    </ReactMarkdown>
  </div>
);
```

#### 8. Реакции на сообщения
```tsx
// Добавить в MessageItem
const MessageReactions = ({ reactions, onReact }) => (
  <div className="flex gap-1 mt-1">
    {reactions.map(reaction => (
      <button
        key={reaction.emoji}
        className="btn btn-xs btn-ghost"
        onClick={() => onReact(reaction.emoji)}
      >
        {reaction.emoji} {reaction.count}
      </button>
    ))}
    <div className="dropdown dropdown-end">
      <label tabIndex={0} className="btn btn-xs btn-ghost">+</label>
      <div className="dropdown-content">
        {/* Emoji picker */}
      </div>
    </div>
  </div>
);
```

## 📅 Timeline

### Неделя 1
- ✅ Аудит системы (выполнено)
- 🔄 Миграция UI на DaisyUI (3 дня)
- 🔄 Критические исправления (2 дня)

### Неделя 2
- Оптимизация производительности
- Тестирование изменений

### Неделя 3-4
- Дополнительные возможности
- Документация
- Развертывание

## 🎯 Метрики успеха

### Производительность
- ⬇️ Время загрузки чата < 1 сек
- ⬇️ Использование памяти -30%
- ⬆️ FPS при скролле > 30

### UX
- ⬆️ Единообразный интерфейс DaisyUI
- ⬆️ Поддержка 1000+ сообщений без лагов
- ⬆️ Мгновенная отправка сообщений

### Надежность
- ⬇️ Количество ошибок -80%
- ⬆️ Автовосстановление WebSocket
- ⬆️ Graceful degradation при сбоях

## 📝 Заметки

- Все изменения должны быть обратно совместимы
- Добавить feature flags для постепенного развертывания
- Провести A/B тестирование новых возможностей
- Мониторить метрики после каждого этапа