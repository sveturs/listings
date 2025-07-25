# Улучшение панели управления витринами (Storefront Dashboard)

## Текущее состояние

### Обзор существующей реализации
Dashboard витрины расположен по адресу `/[locale]/storefronts/[slug]/dashboard` и содержит следующие элементы:

#### 1. Статистические карточки (Quick Stats)
- Активные товары (48 из 52)
- Ожидающие заказы (12 требуют действий)
- Непрочитанные сообщения (7 от клиентов)
- Низкий запас (3 товара на исходе)

#### 2. Быстрые действия (Quick Actions)
- Добавить товар
- Управление заказами
- Товары
- Клиенты
- Акции
- Доставка

#### 3. Последние заказы
Таблица с заказами, показывающая:
- ID заказа
- Клиент
- Количество товаров
- Сумма
- Статус (Ожидает/Доставляется/Выполнен)
- Действия

#### 4. Предупреждение о низком запасе
Список товаров с критически низким запасом:
- iPhone 15 Pro Max (осталось 2)
- MacBook Pro M3 (осталось 1)
- Sony WH-1000XM5 (Нет в наличии)

#### 5. Боковая панель (правая колонка)
- Уведомления (5 новых)
- Сообщения (7 непрочитанных) с превью последних
- Статус магазина (переключатели)

### Проблемы текущей реализации

1. **Фиктивные данные**: Все данные в dashboard захардкожены (48 товаров, 12 заказов и т.д.)

2. **Отсутствие реальной интеграции с API**: 
   - Нет загрузки реальной статистики
   - Нет реальных заказов
   - Нет реальных уведомлений

3. **Недостающие переводы**:
   - `common.rating`
   - `common.phone`
   - `common.email`
   - `common.share`
   - `products.inStock`

4. **Мобильная версия**:
   - Таблица заказов плохо адаптирована для мобильных устройств
   - Горизонтальная прокрутка таблицы не очевидна

## Предложения по улучшению

### 1. Интеграция с реальными данными

#### A. Создать API endpoints для dashboard статистики
```go
// GET /api/v1/storefronts/{slug}/dashboard/stats
type DashboardStats struct {
    ActiveProducts   int `json:"active_products"`
    TotalProducts    int `json:"total_products"`
    PendingOrders    int `json:"pending_orders"`
    UnreadMessages   int `json:"unread_messages"`
    LowStockProducts int `json:"low_stock_products"`
}

// GET /api/v1/storefronts/{slug}/dashboard/recent-orders
type RecentOrdersResponse struct {
    Orders []Order `json:"orders"`
    Total  int     `json:"total"`
}

// GET /api/v1/storefronts/{slug}/dashboard/low-stock
type LowStockResponse struct {
    Products []ProductStock `json:"products"`
}

// GET /api/v1/storefronts/{slug}/dashboard/notifications
type NotificationsResponse struct {
    Notifications []Notification `json:"notifications"`
    UnreadCount   int            `json:"unread_count"`
}
```

#### B. Создать соответствующие Redux actions
```typescript
// storefrontSlice.ts additions
export const fetchDashboardStats = createAsyncThunk(
  'storefronts/fetchDashboardStats',
  async (slug: string) => {
    const response = await api.get(`/storefronts/${slug}/dashboard/stats`);
    return response.data;
  }
);

export const fetchRecentOrders = createAsyncThunk(
  'storefronts/fetchRecentOrders',
  async (slug: string) => {
    const response = await api.get(`/storefronts/${slug}/dashboard/recent-orders`);
    return response.data;
  }
);
```

### 2. Улучшение UI/UX

#### A. Добавить графики и визуализацию
- График продаж за последние 7/30 дней
- График посещений витрины
- Топ-5 популярных товаров

#### B. Улучшить мобильную версию
- Заменить таблицу заказов на карточки для мобильных устройств
- Добавить свайп для переключения между статистическими карточками
- Использовать выпадающее меню для быстрых действий

#### C. Добавить фильтры и сортировку
- Фильтр заказов по статусу
- Фильтр по датам
- Быстрый поиск по заказам

### 3. Новые функции

#### A. Быстрые действия с заказами
- Кнопка "Подтвердить" для новых заказов
- Кнопка "Отправить" для подтвержденных
- Массовые действия с заказами

#### B. Управление запасами прямо из dashboard
- Быстрое пополнение запасов
- Установка минимального уровня запасов
- Автоматические уведомления о низком запасе

#### C. Интеграция с чатом
- Быстрые ответы на частые вопросы
- Шаблоны ответов
- Статус "онлайн/офлайн" для витрины

#### D. Экспорт данных
- Экспорт заказов в CSV/Excel
- Экспорт статистики
- Генерация отчетов

### 4. Технические улучшения

#### A. Оптимизация производительности
- Lazy loading для секций dashboard
- Кеширование статистики на 5 минут
- WebSocket для real-time обновлений

#### B. Добавить тесты
- Unit тесты для компонентов dashboard
- Integration тесты для API
- E2E тесты для критических user flows

### 5. Приоритизация задач

#### Фаза 1 (Критично)
1. Интеграция реальных данных через API
2. Исправление недостающих переводов
3. Базовая адаптация для мобильных устройств

#### Фаза 2 (Важно)
1. Добавление графиков продаж
2. Улучшение мобильной версии
3. Быстрые действия с заказами

#### Фаза 3 (Желательно)
1. Экспорт данных
2. WebSocket обновления
3. Расширенная аналитика

## Следующие шаги

1. Создать API endpoints для dashboard данных
2. Обновить Redux store для работы с реальными данными
3. Рефакторинг компонента dashboard для использования реальных данных
4. Добавить недостающие переводы
5. Улучшить мобильную версию