# Паспорт системы: Redux Store

## 📋 Метаданные
- **Название**: Redux Store  
- **Путь**: `frontend/svetu/src/store/`
- **Роль**: Централизованное управление состоянием приложения
- **Уровень**: Архитектурный слой состояния

## 🎯 Назначение
Система управления состоянием, построенная на Redux Toolkit, обеспечивающая централизованное хранение данных приложения, асинхронные операции с API и real-time обновления через WebSocket.

## 🏗️ Архитектура Store

### Основная конфигурация
**Путь**: `store/index.ts`

```typescript
export const store = configureStore({
  reducer: {
    chat: chatSlice.reducer,
    import: importSlice.reducer,
    payment: paymentSlice.reducer,
    product: productSlice.reducer,
    reviews: reviewsSlice.reducer,
    storefront: storefrontSlice.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ['websocket/messageReceived'],
        ignoredPaths: ['chat.fileUploads.*.file'],
      },
    }).concat(websocketMiddleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
```

### Типизированные хуки
**Путь**: `store/hooks.ts`

```typescript
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;

// Селекторы с мемоизацией
export const useOptimizedSelector = <T>(
  selector: (state: RootState) => T,
  equalityFn?: (left: T, right: T) => boolean
) => useAppSelector(selector, equalityFn);
```

## 📊 Слайсы состояния (6 модулей)

### 1. 💬 Chat Slice
**Путь**: `store/slices/chatSlice.ts`

**Назначение**: Управление real-time чатами и сообщениями

**Структура состояния**:
```typescript
interface ChatState {
  chats: Chat[];
  currentChatId: string | null;
  messages: { [chatId: string]: Message[] };
  isConnected: boolean;
  fileUploads: { [fileId: string]: FileUpload };
  unreadCounts: { [chatId: string]: number };
  typing: { [chatId: string]: string[] };
  error: string | null;
  loading: boolean;
}

interface Message {
  id: string;
  chatId: string;
  senderId: number;
  content: string;
  type: 'text' | 'image' | 'file' | 'emoji';
  timestamp: string;
  status: 'sending' | 'sent' | 'delivered' | 'read' | 'failed';
  attachments?: MessageAttachment[];
}
```

**Async Thunks**:
```typescript
// Загрузка чатов пользователя
loadUserChats: createAsyncThunk('chat/loadUserChats', async () => {
  const response = await chatApi.getUserChats();
  return response.data;
});

// Отправка сообщения
sendMessage: createAsyncThunk('chat/sendMessage', async (messageData) => {
  const response = await chatApi.sendMessage(messageData);
  return response.data;
});

// Загрузка истории сообщений
loadMessages: createAsyncThunk('chat/loadMessages', async ({ chatId, page }) => {
  const response = await chatApi.getMessages(chatId, page);
  return { chatId, messages: response.data, hasMore: response.hasMore };
});
```

**Особенности**:
- WebSocket интеграция для real-time сообщений
- File Upload Manager с прогрессом загрузки
- Optimistic updates для отправленных сообщений
- Typing indicators
- Unread message counter

### 2. 📥 Import Slice  
**Путь**: `store/slices/importSlice.ts`

**Назначение**: Управление импортом товаров из внешних источников

**Структура состояния**:
```typescript
interface ImportState {
  jobs: ImportJob[];
  currentJob: ImportJob | null;
  progress: { [jobId: string]: ImportProgress };
  errors: { [jobId: string]: ImportError[] };
  filters: ImportFilters;
  loading: boolean;
  error: string | null;
}

interface ImportJob {
  id: string;
  type: 'file' | 'url' | 'api';
  source: string;
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';
  totalItems: number;
  processedItems: number;
  successItems: number;
  errorItems: number;
  createdAt: string;
  completedAt?: string;
}
```

**Async Thunks**:
```typescript
// Создание задачи импорта
createImportJob: createAsyncThunk('import/createJob', async (jobData) => {
  const response = await importApi.createJob(jobData);
  return response.data;
});

// Загрузка задач импорта
loadImportJobs: createAsyncThunk('import/loadJobs', async (filters) => {
  const response = await importApi.getJobs(filters);
  return response.data;
});

// Отмена задачи
cancelImportJob: createAsyncThunk('import/cancelJob', async (jobId) => {
  await importApi.cancelJob(jobId);
  return jobId;
});
```

**Особенности**:
- Real-time прогресс импорта
- Валидация данных перед импортом
- Обработка ошибок с детализацией
- Возможность отмены задач
- Экспорт отчетов об ошибках

### 3. 💳 Payment Slice
**Путь**: `store/slices/paymentSlice.ts`

**Назначение**: Управление платежами и балансом пользователя

**Структура состояния**:
```typescript
interface PaymentState {
  balance: UserBalance | null;
  transactions: Transaction[];
  currentPayment: PaymentSession | null;
  methods: PaymentMethod[];
  escrowAccounts: EscrowAccount[];
  loading: boolean;
  error: string | null;
}

interface PaymentSession {
  id: string;
  amount: number;
  currency: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  method: 'card' | 'bank' | 'wallet';
  listingId?: string;
  escrowId?: string;
}
```

**Async Thunks**:
```typescript
// Загрузка баланса пользователя
loadUserBalance: createAsyncThunk('payment/loadBalance', async () => {
  const response = await paymentApi.getUserBalance();
  return response.data;
});

// Инициация платежа
initiatePayment: createAsyncThunk('payment/initiate', async (paymentData) => {
  const response = await paymentApi.initiatePayment(paymentData);
  return response.data;
});

// Загрузка истории транзакций
loadTransactions: createAsyncThunk('payment/loadTransactions', async (filters) => {
  const response = await paymentApi.getTransactions(filters);
  return response.data;
});
```

**Особенности**:
- Интеграция с AllSecure payment gateway
- Escrow система для безопасных покупок
- Поддержка множественных валют
- Real-time обновления баланса
- История транзакций с фильтрацией

### 4. 📦 Product Slice
**Путь**: `store/slices/productSlice.ts`

**Назначение**: Управление товарами в витринах

**Структура состояния**:
```typescript
interface ProductState {
  products: Product[];
  selectedProducts: string[];
  bulkOperation: BulkOperation | null;
  filters: ProductFilters;
  pagination: Pagination;
  loading: boolean;
  error: string | null;
}

interface BulkOperation {
  type: 'delete' | 'update' | 'export';
  status: 'pending' | 'processing' | 'completed' | 'failed';
  progress: number;
  totalItems: number;
  processedItems: number;
  errors: string[];
}
```

**Async Thunks**:
```typescript
// Загрузка товаров витрины
loadStorefrontProducts: createAsyncThunk('product/loadProducts', async ({ 
  storefrontId, 
  filters, 
  page 
}) => {
  const response = await productApi.getStorefrontProducts(storefrontId, filters, page);
  return response.data;
});

// Массовое удаление товаров
bulkDeleteProducts: createAsyncThunk('product/bulkDelete', async (productIds) => {
  const response = await productApi.bulkDelete(productIds);
  return response.data;
});

// Массовое обновление товаров
bulkUpdateProducts: createAsyncThunk('product/bulkUpdate', async ({ 
  productIds, 
  updateData 
}) => {
  const response = await productApi.bulkUpdate(productIds, updateData);
  return response.data;
});
```

**Особенности**:
- Массовые операции с товарами
- Прогресс выполнения операций
- Фильтрация и поиск товаров
- Кэширование результатов
- Оптимистичные обновления

### 5. ⭐ Reviews Slice
**Путь**: `store/slices/reviewsSlice.ts`

**Назначение**: Управление отзывами и рейтингами

**Структура состояния**:
```typescript
interface ReviewsState {
  reviews: Review[];
  stats: ReviewStats | null;
  currentReview: Review | null;
  filters: ReviewFilters;
  pagination: Pagination;
  loading: boolean;
  error: string | null;
}

interface Review {
  id: string;
  listingId: string;
  authorId: number;
  rating: number;
  title: string;
  content: string;
  images: ReviewImage[];
  status: 'pending' | 'approved' | 'rejected';
  createdAt: string;
  helpful: number;
  reported: boolean;
}

interface ReviewStats {
  averageRating: number;
  totalReviews: number;
  ratingDistribution: { [rating: number]: number };
  recentTrend: 'up' | 'down' | 'stable';
}
```

**Async Thunks**:
```typescript
// Загрузка отзывов для объявления
loadListingReviews: createAsyncThunk('reviews/loadReviews', async ({ 
  listingId, 
  filters, 
  page 
}) => {
  const response = await reviewsApi.getListingReviews(listingId, filters, page);
  return response.data;
});

// Создание нового отзыва
createReview: createAsyncThunk('reviews/create', async (reviewData) => {
  const response = await reviewsApi.createReview(reviewData);
  return response.data;
});

// Загрузка статистики отзывов
loadReviewStats: createAsyncThunk('reviews/loadStats', async (listingId) => {
  const response = await reviewsApi.getReviewStats(listingId);
  return response.data;
});
```

**Особенности**:
- Система модерации отзывов
- Загрузка изображений к отзывам
- Статистика и тренды рейтингов
- Система жалоб на отзывы
- Helpful votes для отзывов

### 6. 🏪 Storefront Slice
**Путь**: `store/slices/storefrontSlice.ts`

**Назначение**: Управление витринами магазинов

**Структура состояния**:
```typescript
interface StorefrontState {
  storefronts: Storefront[];
  currentStorefront: Storefront | null;
  nearbyStorefronts: Storefront[];
  mapData: MapStorefront[];
  filters: StorefrontFilters;
  userLocation: UserLocation | null;
  loading: boolean;
  error: string | null;
}

interface Storefront {
  id: string;
  name: string;
  description: string;
  category: string;
  location: Location;
  contact: ContactInfo;
  businessHours: BusinessHours[];
  products: Product[];
  rating: number;
  reviewsCount: number;
  verified: boolean;
  createdAt: string;
}
```

**Async Thunks**:
```typescript
// Загрузка витрин пользователя
loadUserStorefronts: createAsyncThunk('storefront/loadUserStorefronts', async () => {
  const response = await storefrontApi.getUserStorefronts();
  return response.data;
});

// Поиск витрин поблизости
loadNearbyStorefronts: createAsyncThunk('storefront/loadNearby', async ({
  lat,
  lng,
  radius,
  filters
}) => {
  const response = await storefrontApi.getNearbyStorefronts(lat, lng, radius, filters);
  return response.data;
});

// Создание новой витрины
createStorefront: createAsyncThunk('storefront/create', async (storefrontData) => {
  const response = await storefrontApi.createStorefront(storefrontData);
  return response.data;
});
```

**Особенности**:
- Геолокация и поиск по карте
- Фильтрация по категориям и расстоянию
- Управление рабочим временем
- Верификация витрин
- Интеграция с картографическими сервисами

## 🔧 Middleware

### WebSocket Middleware
**Путь**: `store/middleware/websocketMiddleware.ts`

**Назначение**: Обработка real-time событий через WebSocket

```typescript
interface WebSocketMessage {
  type: 'chat_message' | 'notification' | 'status_update';
  payload: any;
  userId?: number;
  chatId?: string;
}

const websocketMiddleware: Middleware = (store) => (next) => (action) => {
  if (action.type === 'websocket/connect') {
    // Установка WebSocket соединения
    websocket = new WebSocket(action.payload.url);
    
    websocket.onmessage = (event) => {
      const message: WebSocketMessage = JSON.parse(event.data);
      
      switch (message.type) {
        case 'chat_message':
          store.dispatch(chatSlice.actions.messageReceived(message.payload));
          break;
        case 'notification':
          store.dispatch(notificationSlice.actions.notificationReceived(message.payload));
          break;
        case 'status_update':
          handleStatusUpdate(store, message.payload);
          break;
      }
    };
  }
  
  return next(action);
};
```

**Обрабатываемые события**:
- Новые сообщения в чате
- Typing indicators
- Статусы доставки сообщений
- Системные уведомления
- Обновления статусов операций

## 🔍 Селекторы

### Мемоизированные селекторы
```typescript
// Общее количество непрочитанных сообщений
export const selectTotalUnreadCount = createSelector(
  [(state: RootState) => state.chat.unreadCounts],
  (unreadCounts) => Object.values(unreadCounts).reduce((sum, count) => sum + count, 0)
);

// Фильтрованные товары
export const selectFilteredProducts = createSelector(
  [
    (state: RootState) => state.product.products,
    (state: RootState) => state.product.filters
  ],
  (products, filters) => {
    return products.filter(product => {
      if (filters.category && product.category !== filters.category) return false;
      if (filters.priceRange && 
          (product.price < filters.priceRange.min || product.price > filters.priceRange.max)) {
        return false;
      }
      return true;
    });
  }
);

// Активные задачи импорта
export const selectActiveImportJobs = createSelector(
  [(state: RootState) => state.import.jobs],
  (jobs) => jobs.filter(job => job.status === 'processing' || job.status === 'pending')
);
```

## ⚡ Оптимизации

### 1. Серializability
```typescript
// Исключение несериализуемых данных
const serializableCheck = {
  ignoredActions: ['websocket/messageReceived', 'chat/fileUpload'],
  ignoredPaths: ['chat.fileUploads.*.file', 'chat.websocketConnection']
};
```

### 2. Мемоизация
- Использование `createSelector` для дорогих вычислений
- Shallow равенство для примитивных значений
- Custom equality functions для сложных объектов

### 3. Code Splitting
```typescript
// Lazy loading для крупных slice'ов
const LazyAdminSlice = lazy(() => import('./slices/adminSlice'));
```

## 🎯 Примеры использования

### Отправка сообщения в чате
```typescript
const ChatComponent = () => {
  const dispatch = useAppDispatch();
  const { messages, currentChatId } = useAppSelector(state => state.chat);
  
  const handleSendMessage = async (content: string) => {
    if (currentChatId) {
      await dispatch(sendMessage({
        chatId: currentChatId,
        content,
        type: 'text'
      }));
    }
  };
  
  return (
    <div>
      {messages[currentChatId || '']?.map(message => (
        <MessageItem key={message.id} message={message} />
      ))}
      <MessageInput onSend={handleSendMessage} />
    </div>
  );
};
```

### Массовые операции с товарами
```typescript
const ProductsManager = () => {
  const dispatch = useAppDispatch();
  const { selectedProducts, bulkOperation } = useAppSelector(state => state.product);
  
  const handleBulkDelete = async () => {
    if (selectedProducts.length > 0) {
      await dispatch(bulkDeleteProducts(selectedProducts));
    }
  };
  
  return (
    <div>
      <button 
        onClick={handleBulkDelete}
        disabled={bulkOperation?.status === 'processing'}
      >
        Delete Selected ({selectedProducts.length})
      </button>
      {bulkOperation && (
        <ProgressBar 
          progress={bulkOperation.progress} 
          status={bulkOperation.status}
        />
      )}
    </div>
  );
};
```

### Создание отзыва
```typescript
const ReviewForm = ({ listingId }: { listingId: string }) => {
  const dispatch = useAppDispatch();
  
  const handleSubmitReview = async (reviewData: CreateReviewData) => {
    const result = await dispatch(createReview({
      ...reviewData,
      listingId
    }));
    
    if (createReview.fulfilled.match(result)) {
      // Успешно создан отзыв
      toast.success('Review submitted successfully');
    }
  };
  
  return <ReviewFormComponent onSubmit={handleSubmitReview} />;
};
```

## 🐛 Обработка ошибок

### Error Handling Pattern
```typescript
// В каждом async thunk
const loadDataThunk = createAsyncThunk(
  'slice/loadData',
  async (params, { rejectWithValue }) => {
    try {
      const response = await api.loadData(params);
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message || 'Unknown error');
    }
  }
);

// В reducers
extraReducers: (builder) => {
  builder
    .addCase(loadDataThunk.pending, (state) => {
      state.loading = true;
      state.error = null;
    })
    .addCase(loadDataThunk.fulfilled, (state, action) => {
      state.loading = false;
      state.data = action.payload;
    })
    .addCase(loadDataThunk.rejected, (state, action) => {
      state.loading = false;
      state.error = action.payload as string;
    });
}
```

## 🔒 Безопасность

### Sanitization
- Санитизация всех пользовательских данных
- Валидация действий перед отправкой
- Проверка прав доступа для чувствительных операций

### Data Protection  
- Не хранение чувствительных данных в store
- Шифрование файлов перед загрузкой
- Очистка состояния при выходе пользователя

## 📈 Мониторинг

### Redux DevTools
```typescript
const store = configureStore({
  // ... reducers
  devTools: process.env.NODE_ENV !== 'production' && {
    trace: true,
    traceLimit: 25,
    actionCreators: true,
  }
});
```

### Performance Monitoring
- Логирование времени выполнения thunk'ов
- Отслеживание размера state
- Мониторинг частоты re-renders

## 🌟 Лучшие практики

1. **Нормализация данных** - плоская структура для сложных связанных данных
2. **Immutability** - использование Immer через RTK
3. **Type Safety** - полная типизация TypeScript
4. **Error Boundaries** - graceful handling ошибок
5. **Optimistic Updates** - мгновенная обратная связь для пользователя
6. **Cache Invalidation** - своевременное обновление устаревших данных