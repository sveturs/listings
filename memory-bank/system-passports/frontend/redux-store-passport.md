# Redux Store System Passport
## Sve Tu Platform Frontend State Management

### 📊 Обзор архитектуры

Redux Store построен на основе **Redux Toolkit (RTK)** с типизацией TypeScript и содержит 6 основных слайсов и 1 middleware для WebSocket соединений.

```
Store Architecture:
├── store/
│   ├── index.ts          # Конфигурация store
│   ├── hooks.ts          # Типизированные хуки
│   ├── slices/           # Состояние по доменам
│   │   ├── chatSlice.ts       # Чаты и сообщения + WebSocket
│   │   ├── importSlice.ts     # Импорт товаров
│   │   ├── paymentSlice.ts    # Платежная система
│   │   ├── productSlice.ts    # Управление товарами
│   │   ├── reviewsSlice.ts    # Отзывы и рейтинги
│   │   └── storefrontSlice.ts # Витрины магазинов
│   └── middleware/
│       └── websocketMiddleware.ts # WebSocket обработка
```

---

## 🏪 Store Configuration (`index.ts`)

### Конфигурация
```typescript
export const store = configureStore({
  reducer: {
    chat: chatReducer,
    reviews: reviewsReducer, 
    storefronts: storefrontsReducer,
    import: importReducer,
    products: productReducer,
    payment: paymentReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // Игнорируем WebSocket, Set и File объекты
        ignoredActions: [...],
        ignoredPaths: [...]
      }
    }).concat(websocketMiddleware)
});
```

### Типизация
- `RootState` - корневое состояние store
- `AppDispatch` - типизированный dispatch
- Типизированные хуки: `useAppDispatch`, `useAppSelector`

---

## 💬 Chat Slice - Система чатов и сообщений

### Назначение
Управление чатами, сообщениями, WebSocket соединениями и файловыми вложениями в real-time.

### Структура State
```typescript
interface ChatState {
  // Основные данные
  chats: MarketplaceChat[];
  currentChat: MarketplaceChat | null;
  messages: Record<number, MarketplaceMessage[]>; // chatId -> messages
  unreadCount: number;

  // WebSocket соединение
  ws: WebSocket | null;
  typingUsers: Record<number, number[]>; // chatId -> userIds
  onlineUsers: number[];
  userLastSeen: Record<number, string>;
  currentUserId: number | null;

  // Пагинация
  chatsPage: number;
  messagesPage: Record<number, number>;
  hasMoreChats: boolean;
  hasMoreMessages: Record<number, boolean>;
  messagesLoaded: Record<number, boolean>;

  // Загрузка файлов
  uploadingFiles: Record<string, UploadingFile>;
  attachments: Record<number, ChatAttachment[]>;

  // UI состояние
  pendingChatId: number | null;
  isLoading: boolean;
  error: string | null;
}
```

### Async Thunks
- `loadChats(page)` - Загрузка списка чатов с пагинацией
- `loadMessages(params)` - Загрузка сообщений чата
- `sendMessage(payload)` - Отправка сообщения с созданием нового чата при необходимости
- `markMessagesAsRead({chatId, messageIds})` - Пометка сообщений как прочитанных
- `archiveChat(chatId)` - Архивирование чата
- `uploadFiles({messageId, files})` - Загрузка файловых вложений с прогрессом
- `deleteAttachment(attachmentId)` - Удаление вложения
- `refreshMessageWithAttachments({chatId, messageId})` - Обновление сообщения с вложениями

### Синхронные Actions
- `setCurrentChat` - Установка активного чата
- `setWebSocket` - Сохранение WebSocket соединения
- `setUserTyping` - Управление индикатором набора текста
- `handleNewMessage` - Обработка нового сообщения через WebSocket
- `handleMessageRead` - Обработка прочтения сообщений
- `handleUserOnline/Offline` - Управление статусом пользователей
- `clearAllData` - Очистка всех данных при выходе

### Особенности реализации
1. **File Upload Manager**: Файлы хранятся вне Redux для производительности
2. **Real-time Updates**: Интеграция с WebSocket middleware
3. **Optimistic Updates**: Сообщения добавляются сразу, затем обновляются
4. **Race Condition Prevention**: Дублирование сообщений от текущего пользователя предотвращается
5. **Auto Chat Selection**: Автоматический выбор нового чата после создания

### Selectors
```typescript
export const selectChats = (state: RootState) => state.chat.chats;
export const selectCurrentChat = (state: RootState) => state.chat.currentChat;
export const selectMessages = (state: RootState, chatId?: number) => ...;
export const selectUnreadCount = (state: RootState) => state.chat.unreadCount;
export const selectOnlineUsers = (state: RootState) => state.chat.onlineUsers;
export const selectTypingUsers = (state: RootState, chatId: number) => ...;
```

---

## 📥 Import Slice - Система импорта товаров

### Назначение
Управление импортом товаров из различных источников (CSV, XML, ZIP файлы, URL).

### Структура State
```typescript
interface ImportState {
  // Задачи импорта
  jobs: ImportJob[];
  currentJob: ImportJob | null;
  
  // Загрузка файлов
  isUploading: boolean;
  uploadProgress: UploadProgress | null;
  validationErrors: string[];
  formats: ImportFormats | null;

  // UI состояние
  isLoading: boolean;
  error: string | null;
  selectedFiles: File[];
  importUrl: string;
  selectedFileType: 'xml' | 'csv' | 'zip' | '';
  updateMode: 'create_only' | 'update_only' | 'upsert';
  categoryMappingMode: 'auto' | 'manual' | 'skip';

  // Модальные окна
  isImportModalOpen: boolean;
  isJobDetailsModalOpen: boolean;
  isErrorsModalOpen: boolean;
}
```

### Async Thunks
- `fetchImportFormats()` - Получение поддерживаемых форматов
- `fetchImportJobs(params)` - Загрузка списка задач импорта
- `importFromUrl(params)` - Импорт из URL (поддержка slug-based API)
- `importFromFile(params)` - Импорт из файла с отслеживанием прогресса
- `validateImportFile(params)` - Валидация файла перед импортом
- `fetchJobStatus(jobId)` - Получение статуса задачи
- `fetchJobDetails(jobId)` - Детали задачи импорта
- `downloadCsvTemplate()` - Скачивание шаблона CSV
- `cancelImportJob(jobId)` - Отмена задачи
- `retryImportJob(jobId)` - Повтор неудачной задачи

### Синхронные Actions
- `setIsUploading` - Управление состоянием загрузки
- `setUploadProgress` - Обновление прогресса загрузки
- `setSelectedFiles` - Выбор файлов для импорта
- UI управление: модальные окна, настройки импорта
- `updateJobStatus` - Обновление статуса задачи в списке
- `resetForm` - Сброс формы после импорта

### Особенности реализации
1. **Progress Tracking**: Реальное отслеживание прогресса загрузки файлов
2. **File Validation**: Предварительная валидация файлов
3. **Slug-based API**: Поддержка API через slug витрины
4. **Job Management**: Полный lifecycle управления задачами импорта
5. **Template Download**: Автоматическое скачивание файлов-шаблонов

---

## 💳 Payment Slice - Платежная система

### Назначение
Управление платежными операциями, кошельком пользователя и методами оплаты.

### Структура State
```typescript
interface PaymentState {
  // Текущий checkout
  checkoutData: {
    listingId: string | null;
    amount: number;
    currency: string;
    paymentMethod: string | null;
    commission: number;
    total: number;
  } | null;

  // Транзакции
  transactions: Transaction[];
  transactionsLoading: boolean;
  transactionsError: string | null;

  // Кошелек
  wallet: {
    balance: number;
    pendingBalance: number;
    currency: string;
  } | null;
  walletLoading: boolean;

  // Методы оплаты
  paymentMethods: PaymentMethod[];
  paymentMethodsLoading: boolean;

  // Процесс оплаты
  paymentProcessing: boolean;
  paymentError: string | null;
  lastPaymentId: string | null;
}
```

### Async Thunks
- `createPayment(data)` - Создание платежа
- `fetchTransactions(params?)` - Получение истории транзакций
- `fetchWallet()` - Получение баланса кошелька
- `fetchPaymentMethods()` - Загрузка доступных методов оплаты
- `requestWithdrawal(data)` - Запрос на вывод средств
- `confirmPayment(paymentId)` - Подтверждение платежа
- `refundPayment({paymentId, reason})` - Возврат платежа

### Синхронные Actions
- `setCheckoutData` - Установка данных для оплаты (с автоматическим расчетом комиссии)
- `setPaymentMethod` - Выбор метода оплаты (с перерасчетом доплат)
- `clearCheckout` - Очистка данных checkout
- `clearPaymentError` - Сброс ошибок

### Особенности реализации
1. **Commission Calculation**: Автоматический расчет комиссии по категориям
2. **Dynamic Pricing**: Перерасчет суммы при смене метода оплаты
3. **Transaction History**: Фильтрация и пагинация транзакций
4. **Wallet Integration**: Управление балансом и выводом средств

---

## 📦 Product Slice - Управление товарами

### Назначение
Массовые операции с товарами, фильтрация, сортировка и экспорт.

### Структура State
```typescript
interface ProductState {
  products: StorefrontProduct[];
  selectedIds: number[];
  loading: boolean;
  error: string | null;

  // Фильтры
  filters: {
    search: string;
    categoryId: number | null;
    minPrice: number | null;
    maxPrice: number | null;
    stockStatus: 'all' | 'in_stock' | 'low_stock' | 'out_of_stock';
    isActive: boolean | null;
  };

  // Пагинация
  pagination: {
    page: number;
    limit: number;
    total: number;
    hasMore: boolean;
  };

  // Массовые операции
  bulkOperation: {
    isProcessing: boolean;
    progress: number;
    total: number;
    errors: BulkOperationError[];
    successCount: number;
    currentOperation: 'idle' | 'delete' | 'update' | 'status' | 'export';
  };

  // UI состояния
  ui: {
    isSelectMode: boolean;
    viewMode: 'grid' | 'list' | 'table';
    sortBy: 'name' | 'price' | 'created_at' | 'stock_quantity';
    sortOrder: 'asc' | 'desc';
  };
}
```

### Async Thunks
- `bulkDeleteProducts({storefrontSlug, productIds})` - Массовое удаление
- `bulkUpdateStatus({storefrontSlug, productIds, isActive})` - Массовое изменение статуса
- `exportProducts({storefrontSlug, productIds?, format})` - Экспорт в CSV/XML

### Синхронные Actions
- **Selection Management**: `toggleProductSelection`, `selectAll`, `clearSelection`, `selectByFilter`
- **UI Control**: `toggleSelectMode`, `setViewMode`, `setSortBy`
- **Filters**: `setFilters`, `resetFilters`
- **Products**: `setProducts`, `appendProducts`, `updateProduct`, `removeProducts`
- **Bulk Operations**: `setBulkOperationProgress`, `addBulkOperationError`, `resetBulkOperation`
- **Pagination**: `setPagination`

### Особенности реализации
1. **Smart Selection**: Фильтрация выбора товаров по предикатам
2. **Progress Tracking**: Отслеживание прогресса массовых операций
3. **Error Handling**: Детальная отчетность об ошибках операций
4. **Toast Notifications**: Интеграция с react-hot-toast
5. **View Modes**: Поддержка различных режимов отображения
6. **Export Features**: Экспорт выбранных товаров или всех

---

## ⭐ Reviews Slice - Система отзывов

### Назначение
Управление отзывами, рейтингами и голосованием за полезность отзывов.

### Структура State
```typescript
interface ReviewsState {
  reviews: Review[];
  currentReview: Review | null;
  stats: ReviewStats | null;
  aggregatedRating: AggregatedRating | null;
  canReview: CanReviewResponse | null;
  filters: ReviewsFilter;
  loading: boolean;
  error: string | null;
  totalPages: number;
  currentPage: number;
}
```

### Async Thunks
- `fetchReviews(filters)` - Загрузка отзывов с фильтрацией
- `fetchReviewStats({entityType, entityId})` - Статистика отзывов
- `fetchAggregatedRating({entityType, entityId})` - Агрегированный рейтинг
- `checkCanReview({entityType, entityId})` - Проверка возможности оставить отзыв
- `createReview(reviewData)` - Создание нового отзыва
- `voteReview({reviewId, voteType})` - Голосование за полезность
- `confirmReview({reviewId, notes?})` - Подтверждение отзыва
- `disputeReview({reviewId, reason})` - Оспаривание отзыва

### Синхронные Actions
- `setFilters` - Установка фильтров отзывов
- `clearError` - Сброс ошибок
- `setCurrentReview` - Установка текущего отзыва
- `updateReviewInList` - Обновление отзыва в списке

### Особенности реализации
1. **Entity-based Reviews**: Универсальная система для любых сущностей
2. **Vote Management**: Подсчет голосов за полезность с предотвращением дублирования
3. **Permission System**: Проверка прав на создание отзывов
4. **Dispute System**: Механизм оспаривания отзывов
5. **Real-time Updates**: Мгновенное обновление счетчиков голосов

---

## 🏪 Storefront Slice - Управление витринами

### Назначение
CRUD операции с витринами, фильтрация, геолокация и аналитика.

### Структура State
```typescript
interface StorefrontState {
  // Основные данные
  storefronts: Storefront[];
  currentStorefront: Storefront | null;
  myStorefronts: Storefront[];

  // UI состояние
  isLoading: boolean;
  error: string | null;
  isCreating: boolean;
  isUpdating: boolean;
  isDeleting: boolean;

  // Фильтры и пагинация
  filters: StorefrontFilters;
  pagination: PaginationParams;
  totalCount: number;
  hasMore: boolean;

  // Дополнительные данные
  analytics: StorefrontAnalytics[] | null;
  ratingSummary: StorefrontRatingSummary | null;
  deliveryOptions: StorefrontDeliveryOption[];
  paymentMethods: StorefrontPaymentMethod[];
  workingHours: StorefrontHours[];

  // Состояние загрузки
  isLoadingAnalytics: boolean;
  isLoadingSettings: boolean;
}
```

### Async Thunks
- `fetchStorefronts({filters?, pagination?})` - Поиск витрин с фильтрами
- `fetchStorefrontById(id)` - Получение витрины по ID
- `fetchStorefrontBySlug(slug)` - Получение витрины по slug (с авторизацией)
- `fetchMyStorefronts()` - Мои витрины (требует авторизации)
- `createStorefront(data)` - Создание новой витрины
- `updateStorefront({id, data})` - Обновление витрины
- `deleteStorefront(id)` - Удаление витрины
- `fetchStorefrontAnalytics({id, from?, to?})` - Аналитика витрины

### Синхронные Actions
- **Filters**: `setFilters`, `clearFilters`
- **Pagination**: `setPagination`
- **Current Storefront**: `setCurrentStorefront`, `clearCurrentStorefront`
- **Error Handling**: `clearError`
- **Data Sync**: `updateStorefrontInLists` - синхронизация витрины во всех списках

### Особенности реализации
1. **Advanced Filtering**: Геолокация, рейтинги, методы оплаты, расписание работы
2. **Token Management**: Интеграция с tokenManager для авторизованных запросов
3. **Data Synchronization**: Автоматическое обновление витрины во всех списках
4. **Cookie Support**: Поддержка httpOnly cookies для авторизации
5. **Memoized Selectors**: Оптимизированные селекторы с createSelector
6. **Graceful Error Handling**: Детальная обработка ошибок сети

### Селекторы
```typescript
// Базовые селекторы
export const selectStorefronts, selectCurrentStorefront, selectMyStorefronts, ...

// Мемоизированные селекторы
export const selectVerifiedStorefronts = createSelector(...);
export const selectStorefrontsByCity = createSelector(...);
export const selectActiveStorefronts = createSelector(...);
```

---

## 🔌 WebSocket Middleware

### Назначение
Обработка WebSocket соединений для real-time обновлений чатов.

### Архитектура
```typescript
export const websocketMiddleware: Middleware = (store) => (next) => (action) => {
  // Инициализация WebSocket
  if (action.type === 'chat/initWebSocket') {
    // Создание соединения, установка обработчиков
  }
  
  // Закрытие WebSocket
  if (action.type === 'chat/closeWebSocket') {
    // Очистка соединения и таймеров
  }
  
  return next(action);
};
```

### Обрабатываемые события
- `new_message` - Новое сообщение с обработкой вложений
- `message_read` - Прочтение сообщений
- `user_typing` - Индикатор набора текста
- `user_online/offline` - Статус пользователей
- `attachment_upload/delete` - Управление вложениями
- `online_users_list` - Список онлайн пользователей

### Особенности
1. **Heartbeat System**: Поддержание соединения через ping/pong
2. **Auto Reconnect**: Автоматическое переподключение
3. **User Status Tracking**: Запрос статуса всех пользователей при подключении
4. **Attachment Handling**: Отложенная загрузка вложений для новых сообщений
5. **Error Handling**: Обработка ошибок парсинга сообщений

---

## 🔧 Технические особенности

### Serializable Check Configuration
```typescript
serializableCheck: {
  ignoredActions: [
    'chat/initWebSocket',
    'chat/setWebSocket', 
    'import/setSelectedFiles',
    // ... другие non-serializable actions
  ],
  ignoredPaths: [
    'chat.ws',
    'chat.onlineUsers',
    'import.selectedFiles',
    'products.selectedIds',
    // ... пути с non-serializable данными
  ]
}
```

### TypeScript Integration
- Полная типизация всех slice'ов
- Типизированные хуки (`useAppDispatch`, `useAppSelector`)
- Использование сгенерированных API типов
- Типизация async thunk параметров и возвращаемых значений

### Performance Optimizations
1. **File Upload Manager**: Файлы хранятся вне Redux
2. **Memoized Selectors**: Оптимизированные селекторы для вычислений
3. **Efficient Updates**: Использование immer для immutable обновлений
4. **Lazy Loading**: Постепенная загрузка данных с пагинацией

### Integration Points
- **API Client**: Интеграция с unified HTTP client
- **Token Manager**: Автоматическое управление авторизацией
- **Toast Notifications**: Интеграция с react-hot-toast
- **File Manager**: Отдельный менеджер для загрузки файлов
- **WebSocket Service**: Централизованный сервис для WebSocket соединений

---

## 📈 Метрики и мониторинг

### Состояние загрузки
Каждый slice имеет индикаторы состояния:
- `isLoading` - общая загрузка
- Специфичные флаги (`isCreating`, `isUpdating`, `isDeleting`)
- `error` - обработка ошибок

### Прогресс операций
- Upload progress для файлов
- Bulk operation progress для массовых операций
- WebSocket connection status

### Error Handling Strategy
1. **Graceful Degradation**: Приложение продолжает работать при ошибках
2. **User Feedback**: Ошибки отображаются пользователю через toast
3. **Retry Mechanisms**: Возможность повтора неудачных операций
4. **Offline Support**: Базовая поддержка работы оффлайн

---

*Последнее обновление: 29 декабря 2024*  
*Версия Redux Toolkit: ^2.0.0*  
*TypeScript: ^5.0.0*