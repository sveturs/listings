# Задание на устранение недостатков системы чатов (Обновленная версия)

## Резюме анализа

После детального анализа системы чатов на всех уровнях (БД, Backend, Frontend в svetu, MinIO) выявлены следующие ключевые аспекты:

### Сильные стороны:
- **База данных**: Отлично спроектированная структура с оптимизированными индексами
- **Backend**: Чистая архитектура с разделением на слои, WebSocket для real-time
- **Frontend (svetu)**: 
  - Modern stack: Next.js 15, TypeScript, Redux Toolkit
  - Полная поддержка загрузки файлов с прогрессом
  - WebSocket интеграция с индикаторами онлайн статуса
  - Поддержка прямых сообщений между пользователями
  - Интернационализация через next-intl
- **MinIO**: Организованное хранение файлов с валидацией

### Основные недостатки:
1. **🔴 Критично - Безопасность файлов**: публичный доступ к файлам через прямые URL
2. **Отсутствующий функционал**:
   - Редактирование/удаление сообщений
   - Поиск по истории чатов
   - Реакции на сообщения
   - Распознавание телефонов/email (было в старом фронте)
   - Звуковые уведомления
3. **Производительность**:
   - Нет оптимизации изображений
   - Отсутствуют превью для видео
   - Нет lazy loading для старых сообщений
4. **UX проблемы**:
   - Нет сохранения черновиков
   - Отсутствует индикатор "сообщение доставлено/прочитано"
   - Нет группировки сообщений по датам

## План устранения недостатков

### Фаза 1: Критические проблемы безопасности (1-2 недели)

#### 1.1 Приватный доступ к файлам
**Проблема**: Файлы чатов доступны публично через прямые URL
**Решение**:
```go
// Backend: новый endpoint для безопасного доступа к файлам
// POST /api/v1/marketplace/attachments/{id}/presigned-url
func (s *ChatService) GeneratePresignedURL(ctx context.Context, userID, attachmentID int) (string, error) {
    // 1. Проверить что пользователь - участник чата
    attachment, err := s.storage.GetAttachment(ctx, attachmentID)
    if err != nil {
        return "", err
    }
    
    // 2. Проверить доступ
    hasAccess, err := s.storage.UserHasAccessToChat(ctx, userID, attachment.ChatID)
    if !hasAccess {
        return "", ErrAccessDenied
    }
    
    // 3. Генерировать presigned URL с TTL 1 час
    url, err := s.minioClient.PresignedGetObject(ctx, "chat-files", attachment.FilePath, time.Hour)
    
    // 4. Логировать доступ
    s.logger.Info("File access granted", "user_id", userID, "attachment_id", attachmentID)
    
    return url, err
}
```

**Frontend изменения**:
```typescript
// svetu: обновить компонент отображения файлов
const FileAttachment: React.FC<{attachment: Attachment}> = ({attachment}) => {
  const [url, setUrl] = useState<string | null>(null);
  
  useEffect(() => {
    // Получаем безопасный URL при монтировании
    chatService.getAttachmentURL(attachment.id)
      .then(setUrl)
      .catch(console.error);
  }, [attachment.id]);
  
  if (!url) return <Skeleton />;
  
  return <img src={url} alt={attachment.name} />;
};
```

#### 1.2 Антивирусная проверка
**Реализация**:
- Интегрировать ClamAV через Docker
- Добавить очередь проверки файлов (Redis + asynq)
- Блокировать доступ до завершения сканирования

### Фаза 2: Улучшение функциональности (2-3 недели)

#### 2.1 Редактирование и удаление сообщений
**База данных**:
```sql
ALTER TABLE marketplace_messages 
ADD COLUMN edited_at TIMESTAMP,
ADD COLUMN deleted_at TIMESTAMP,
ADD COLUMN edit_history JSONB DEFAULT '[]'::jsonb,
ADD COLUMN original_content TEXT;

CREATE INDEX idx_messages_deleted_at ON marketplace_messages(deleted_at) 
WHERE deleted_at IS NULL;
```

**Backend API**:
```go
// PUT /api/v1/marketplace/messages/{id}
func (h *ChatHandler) EditMessage(c *fiber.Ctx) error {
    messageID := c.Params("id")
    userID := c.Locals("user_id").(int)
    
    var req EditMessageRequest
    if err := c.BodyParser(&req); err != nil {
        return err
    }
    
    // Проверяем что это автор сообщения
    // Сохраняем историю изменений
    // Отправляем WebSocket событие всем участникам
}
```

**Frontend (svetu)**:
```typescript
// Добавить в Redux slice
editMessage: (state, action) => {
  const {messageId, newContent} = action.payload;
  const message = state.messages[messageId];
  if (message) {
    message.content = newContent;
    message.edited_at = new Date().toISOString();
    message.edit_history.push({
      content: message.content,
      edited_at: message.edited_at
    });
  }
}
```

#### 2.2 Распознавание контактов (портировать из старого фронта)
```typescript
// Утилита для распознавания телефонов и email
export const detectContacts = (text: string) => {
  // Паттерн для телефонов
  const phoneRegex = /(\+\d{1,3}[-.\s]?)?\(?\d{1,4}\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}/g;
  
  // Паттерн для email
  const emailRegex = /[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}/g;
  
  const phones = text.match(phoneRegex) || [];
  const emails = text.match(emailRegex) || [];
  
  return {phones, emails};
};

// Компонент сообщения с подсветкой контактов
const MessageContent: React.FC<{content: string}> = ({content}) => {
  const contacts = detectContacts(content);
  
  // Подсветка и делаем кликабельными
  let processedContent = content;
  contacts.phones.forEach(phone => {
    processedContent = processedContent.replace(
      phone, 
      `<a href="tel:${phone}" class="text-primary">${phone}</a>`
    );
  });
  
  return <div dangerouslySetInnerHTML={{__html: processedContent}} />;
};
```

#### 2.3 Поиск по сообщениям
**Backend**:
```sql
-- Добавить полнотекстовый поиск
ALTER TABLE marketplace_messages 
ADD COLUMN search_vector tsvector;

CREATE INDEX idx_messages_search ON marketplace_messages 
USING GIN(search_vector);

-- Триггер для обновления search_vector
CREATE TRIGGER update_messages_search 
BEFORE INSERT OR UPDATE ON marketplace_messages
FOR EACH ROW EXECUTE FUNCTION 
tsvector_update_trigger(search_vector, 'pg_catalog.english', content);
```

**Frontend компонент поиска**:
```typescript
const ChatSearch: React.FC = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<Message[]>([]);
  
  const handleSearch = useDebouncedCallback(async (q: string) => {
    if (q.length < 3) return;
    
    const data = await chatService.searchMessages(chatId, q);
    setResults(data);
  }, 300);
  
  return (
    <div className="relative">
      <input 
        type="text"
        value={query}
        onChange={(e) => {
          setQuery(e.target.value);
          handleSearch(e.target.value);
        }}
        className="input input-bordered w-full"
        placeholder="Поиск по сообщениям..."
      />
      {results.length > 0 && (
        <div className="absolute top-full mt-2 w-full bg-base-100 rounded-lg shadow-lg max-h-96 overflow-y-auto">
          {results.map(msg => (
            <SearchResult key={msg.id} message={msg} onClick={() => scrollToMessage(msg.id)} />
          ))}
        </div>
      )}
    </div>
  );
};
```

### Фаза 3: Оптимизация производительности (2-3 недели)

#### 3.1 Оптимизация изображений
**Backend сервис**:
```go
import "github.com/disintegration/imaging"

func (s *ImageService) ProcessChatImage(ctx context.Context, file io.Reader, filename string) error {
    img, err := imaging.Decode(file)
    if err != nil {
        return err
    }
    
    // Генерируем размеры
    sizes := map[string]int{
        "thumb": 150,
        "preview": 400,
        "full": 1200,
    }
    
    for sizeName, width := range sizes {
        resized := imaging.Resize(img, width, 0, imaging.Lanczos)
        
        // Конвертируем в WebP для экономии места
        path := fmt.Sprintf("%s_%s.webp", filename, sizeName)
        err = s.uploadToMinio(ctx, resized, path)
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

#### 3.2 Виртуализация списка сообщений
```typescript
import { VariableSizeList } from 'react-window';

const VirtualizedChat: React.FC = () => {
  const messages = useAppSelector(selectChatMessages);
  const listRef = useRef<VariableSizeList>(null);
  
  // Кеш высот сообщений
  const rowHeights = useRef<{[key: string]: number}>({});
  
  const getRowHeight = (index: number) => {
    return rowHeights.current[messages[index].id] || 100;
  };
  
  const setRowHeight = (index: number, height: number) => {
    if (rowHeights.current[messages[index].id] !== height) {
      rowHeights.current[messages[index].id] = height;
      listRef.current?.resetAfterIndex(index);
    }
  };
  
  return (
    <VariableSizeList
      ref={listRef}
      height={600}
      itemCount={messages.length}
      itemSize={getRowHeight}
      width="100%"
      overscanCount={5}
    >
      {({index, style}) => (
        <div style={style}>
          <MessageRow 
            message={messages[index]} 
            onHeightChange={(height) => setRowHeight(index, height)}
          />
        </div>
      )}
    </VariableSizeList>
  );
};
```

### Фаза 4: Улучшение UX (1-2 недели)

#### 4.1 Звуковые уведомления
```typescript
// Добавить в ChatContext
const playNotificationSound = useCallback(() => {
  if (!userSettings.soundEnabled) return;
  
  const audio = new Audio('/notification.mp3');
  audio.volume = 0.5;
  audio.play().catch(console.error);
}, [userSettings.soundEnabled]);

// В обработчике новых сообщений
useEffect(() => {
  const handleNewMessage = (message: Message) => {
    if (message.sender_id !== currentUserId && !document.hasFocus()) {
      playNotificationSound();
      
      // Push уведомление если разрешено
      if (Notification.permission === 'granted') {
        new Notification(`Новое сообщение от ${message.sender_name}`, {
          body: message.content.substring(0, 100),
          icon: '/logo192.png',
        });
      }
    }
  };
  
  chatService.on('new_message', handleNewMessage);
  return () => chatService.off('new_message', handleNewMessage);
}, [currentUserId, playNotificationSound]);
```

#### 4.2 Статусы доставки и прочтения
```typescript
// Компонент статуса сообщения
const MessageStatus: React.FC<{message: Message}> = ({message}) => {
  if (message.sender_id !== currentUserId) return null;
  
  return (
    <div className="text-xs text-base-content/60 mt-1">
      {message.is_read ? (
        <span className="flex items-center gap-1">
          <CheckCheck className="w-4 h-4 text-primary" />
          {message.read_at && formatTime(message.read_at)}
        </span>
      ) : message.delivered_at ? (
        <span className="flex items-center gap-1">
          <Check className="w-4 h-4" />
          Доставлено
        </span>
      ) : (
        <span className="flex items-center gap-1">
          <Clock className="w-4 h-4" />
          Отправляется...
        </span>
      )}
    </div>
  );
};
```

#### 4.3 Сохранение черновиков
```typescript
// Hook для автосохранения черновиков
const useDraftMessage = (chatId: string) => {
  const [draft, setDraft] = useState('');
  const draftKey = `chat_draft_${chatId}`;
  
  // Загружаем черновик при монтировании
  useEffect(() => {
    const saved = localStorage.getItem(draftKey);
    if (saved) setDraft(saved);
  }, [chatId]);
  
  // Сохраняем черновик при изменении
  const updateDraft = useDebouncedCallback((text: string) => {
    if (text.trim()) {
      localStorage.setItem(draftKey, text);
    } else {
      localStorage.removeItem(draftKey);
    }
  }, 500);
  
  const clearDraft = () => {
    setDraft('');
    localStorage.removeItem(draftKey);
  };
  
  return {draft, setDraft: (text: string) => {
    setDraft(text);
    updateDraft(text);
  }, clearDraft};
};
```

### Фаза 5: Дополнительные улучшения (2-3 недели)

#### 5.1 Типы сообщений и rich content
```typescript
enum MessageType {
  TEXT = 'text',
  IMAGE = 'image',
  VIDEO = 'video',
  DOCUMENT = 'document',
  VOICE = 'voice',
  LOCATION = 'location',
  PRODUCT = 'product', // Ссылка на товар
  SYSTEM = 'system'
}

// Компонент для разных типов контента
const MessageContent: React.FC<{message: Message}> = ({message}) => {
  switch (message.type) {
    case MessageType.PRODUCT:
      return <ProductPreview productId={message.metadata.product_id} />;
      
    case MessageType.LOCATION:
      return <LocationMap coords={message.metadata.coordinates} />;
      
    case MessageType.VOICE:
      return <VoicePlayer url={message.content} duration={message.metadata.duration} />;
      
    default:
      return <TextMessage content={message.content} />;
  }
};
```

#### 5.2 Оффлайн поддержка с Service Worker
```typescript
// service-worker.ts
import { precacheAndRoute } from 'workbox-precaching';
import { registerRoute } from 'workbox-routing';
import { StaleWhileRevalidate, NetworkFirst } from 'workbox-strategies';

// Кеширование статики
precacheAndRoute(self.__WB_MANIFEST);

// Кеширование API запросов чатов
registerRoute(
  ({url}) => url.pathname.startsWith('/api/v1/marketplace/chat'),
  new NetworkFirst({
    cacheName: 'chat-cache',
    plugins: [
      {
        cacheWillUpdate: async ({response}) => {
          // Кешируем только успешные ответы
          if (response && response.status === 200) {
            return response;
          }
          return null;
        }
      }
    ]
  })
);

// Очередь для отправки сообщений оффлайн
const messageQueue = [];

self.addEventListener('fetch', (event) => {
  if (event.request.url.includes('/messages') && event.request.method === 'POST') {
    event.respondWith(
      fetch(event.request.clone()).catch(() => {
        // Сохраняем в очередь если оффлайн
        return event.request.json().then(body => {
          messageQueue.push({
            url: event.request.url,
            body,
            headers: Object.fromEntries(event.request.headers)
          });
          
          // Возвращаем фейковый успешный ответ
          return new Response(JSON.stringify({
            success: true,
            data: {...body, id: `temp_${Date.now()}`, pending: true}
          }));
        });
      })
    );
  }
});

// Синхронизация при восстановлении соединения
self.addEventListener('sync', (event) => {
  if (event.tag === 'sync-messages') {
    event.waitUntil(
      Promise.all(
        messageQueue.map(req => 
          fetch(req.url, {
            method: 'POST',
            headers: req.headers,
            body: JSON.stringify(req.body)
          })
        )
      ).then(() => {
        messageQueue.length = 0;
      })
    );
  }
});
```

## Метрики успеха

### Безопасность:
- **0 инцидентов** с несанкционированным доступом к файлам
- **100% файлов** проверены антивирусом перед доступом

### Производительность:
- **60-80% снижение** размера изображений через оптимизацию
- **<100ms** время загрузки новых сообщений
- **<2s** время загрузки истории чата

### Пользовательский опыт:
- **95%+ сообщений** доставляются в течение 1 секунды
- **30% увеличение** времени в чатах
- **50% снижение** жалоб на функциональность чатов

### Технические метрики:
- **90%+ покрытие** тестами критических функций
- **0 критических** уязвимостей в зависимостях
- **<5% CPU** использование при 1000 активных чатах

## Приоритеты реализации

### Критические (Sprint 1):
1. ✅ Безопасность файлов через presigned URLs
2. ✅ Базовая антивирусная проверка

### Высокий приоритет (Sprint 2-3):
1. Редактирование и удаление сообщений
2. Распознавание контактов 
3. Оптимизация изображений
4. Звуковые уведомления

### Средний приоритет (Sprint 4-5):
1. Поиск по сообщениям
2. Виртуализация для производительности
3. Статусы доставки/прочтения
4. Сохранение черновиков

### Низкий приоритет (Backlog):
1. Rich типы сообщений
2. Оффлайн поддержка
3. Голосовые сообщения
4. End-to-end шифрование

## Технологический стек

### Backend:
- **Обработка изображений**: `github.com/disintegration/imaging`
- **Антивирус**: ClamAV + `github.com/dutchcoders/go-clamd`
- **Очереди задач**: Redis + `github.com/hibiken/asynq`
- **Полнотекстовый поиск**: PostgreSQL FTS

### Frontend (svetu):
- **Виртуализация**: `react-window`
- **Service Workers**: `workbox`
- **Состояние**: Redux Toolkit (уже используется)
- **UI компоненты**: DaisyUI + Tailwind CSS
- **WebSocket**: нативная реализация (уже есть)

## Заключение

Система чатов в новом фронтенде (svetu) имеет современную архитектуру и хорошую основу. Основные приоритеты:
1. Устранить критическую проблему безопасности с публичным доступом к файлам
2. Добавить базовый функционал из старого фронта (распознавание контактов)
3. Реализовать must-have функции (редактирование, поиск)
4. Оптимизировать производительность для масштабирования

Предложенный план позволит создать современную, безопасную и удобную систему чатов, соответствующую лучшим практикам.

**Общее время реализации**: 8-12 недель
**Рекомендуемая команда**: 1-2 backend + 1-2 frontend разработчика