# План реализации функционала витрин маркетплейса

## Фаза 1: Базовая функциональность товаров (Приоритет: Высокий)

### 1. API endpoints для управления товарами витрины
- [ ] GET `/api/v1/storefronts/{id}/products` - список товаров витрины
- [ ] GET `/api/v1/storefronts/{id}/products/{product_id}` - детали товара
- [ ] POST `/api/v1/storefronts/{id}/products` - создание товара
- [ ] PUT `/api/v1/storefronts/{id}/products/{product_id}` - обновление товара
- [ ] DELETE `/api/v1/storefronts/{id}/products/{product_id}` - удаление товара
- [ ] POST `/api/v1/storefronts/{id}/products/{product_id}/images` - загрузка изображений
- [ ] DELETE `/api/v1/storefronts/{id}/products/{product_id}/images/{image_id}` - удаление изображения

### 2. Структура данных товара
```typescript
interface StorefrontProduct {
  id: number;
  storefront_id: number;
  name: string;
  description: string;
  price: number;
  currency: string;
  category_id: number;
  sku?: string;
  barcode?: string;
  stock_quantity: number;
  stock_status: 'in_stock' | 'low_stock' | 'out_of_stock';
  images: ProductImage[];
  attributes: Record<string, any>;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
```

### 3. Страница добавления товара `/storefronts/[slug]/products/new`
- [ ] Форма с полями: название, описание, цена, категория, количество
- [ ] Загрузка до 10 изображений с drag & drop
- [ ] Выбор категории с динамическими атрибутами
- [ ] Управление складскими остатками
- [ ] Предпросмотр перед сохранением

### 4. Страница списка товаров `/storefronts/[slug]/products`
- [ ] Таблица/сетка товаров с фильтрами
- [ ] Быстрое редактирование цены и остатков
- [ ] Массовые действия (активация/деактивация, удаление)
- [ ] Поиск по названию и SKU
- [ ] Фильтры по категории, статусу, остаткам

## Фаза 2: Управление заказами (Приоритет: Высокий)

### 1. API endpoints для заказов
- [ ] GET `/api/v1/storefronts/{id}/orders` - список заказов
- [ ] GET `/api/v1/storefronts/{id}/orders/{order_id}` - детали заказа
- [ ] PUT `/api/v1/storefronts/{id}/orders/{order_id}/status` - изменение статуса
- [ ] POST `/api/v1/storefronts/{id}/orders/{order_id}/messages` - сообщение покупателю
- [ ] GET `/api/v1/storefronts/{id}/orders/stats` - статистика заказов

### 2. Структура данных заказа
```typescript
interface StorefrontOrder {
  id: number;
  order_number: string;
  storefront_id: number;
  customer_id: number;
  status: 'pending' | 'processing' | 'shipped' | 'delivered' | 'cancelled';
  items: OrderItem[];
  total_amount: number;
  shipping_address: Address;
  payment_method: string;
  payment_status: 'pending' | 'paid' | 'refunded';
  tracking_number?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}
```

### 3. Страница управления заказами `/storefronts/[slug]/orders`
- [ ] Список заказов с фильтрами по статусу
- [ ] Детальный просмотр заказа
- [ ] Изменение статуса заказа
- [ ] Печать накладной/чека
- [ ] Добавление трекинг-номера
- [ ] История изменений заказа

## Фаза 3: Управление складом (Приоритет: Высокий)

### 1. API endpoints для склада
- [ ] GET `/api/v1/storefronts/{id}/inventory` - состояние склада
- [ ] PUT `/api/v1/storefronts/{id}/inventory/{product_id}` - обновление остатков
- [ ] GET `/api/v1/storefronts/{id}/inventory/movements` - история движений
- [ ] POST `/api/v1/storefronts/{id}/inventory/bulk-update` - массовое обновление

### 2. Страница управления складом `/storefronts/[slug]/inventory`
- [ ] Таблица товаров с текущими остатками
- [ ] Предупреждения о низких остатках
- [ ] История движений товара
- [ ] Массовое обновление остатков
- [ ] Экспорт/импорт данных склада

## Фаза 4: Расширенная функциональность (Приоритет: Средний)

### 1. Система уведомлений
- [ ] WebSocket подключение для real-time уведомлений
- [ ] Уведомления о новых заказах
- [ ] Уведомления о низких остатках
- [ ] Уведомления о новых сообщениях
- [ ] Настройки уведомлений

### 2. Интеграция чата
- [ ] Чат с покупателями на странице дашборда
- [ ] История переписки по заказам
- [ ] Быстрые ответы/шаблоны
- [ ] Уведомления о новых сообщениях

### 3. Управление статусом магазина
- [ ] API для изменения статуса (открыт/закрыт)
- [ ] API для режима отпуска
- [ ] Автоматическое закрытие по расписанию
- [ ] Уведомления покупателям о статусе

### 4. Страница настроек `/storefronts/[slug]/settings`
- [ ] Общие настройки витрины
- [ ] Настройки доставки и оплаты
- [ ] Управление персоналом
- [ ] Настройки уведомлений
- [ ] Интеграции с внешними сервисами

## Фаза 5: Продвинутые функции (Приоритет: Низкий)

### 1. Промо-акции и скидки
- [ ] API для создания промо-акций
- [ ] Купоны и промокоды
- [ ] Скидки по категориям
- [ ] Временные акции
- [ ] Страница управления `/storefronts/[slug]/promotions`

### 2. Управление клиентами
- [ ] База клиентов витрины
- [ ] История покупок клиента
- [ ] Программа лояльности
- [ ] Email-рассылки
- [ ] Страница клиентов `/storefronts/[slug]/customers`

### 3. Расширенная аналитика
- [ ] Экспорт отчетов в PDF/Excel
- [ ] Сравнение периодов
- [ ] Прогнозирование продаж
- [ ] A/B тестирование
- [ ] Интеграция с Google Analytics

## Технические требования

### Backend (Go)
- Использовать существующую архитектуру с Fiber
- Добавить новые таблицы в PostgreSQL
- Интегрировать с MinIO для изображений товаров
- Добавить индексы в OpenSearch для товаров
- Реализовать middleware для проверки владельца витрины

### Frontend (Next.js)
- Использовать Redux Toolkit для состояния
- Реализовать оптимистичные обновления
- Добавить skeleton loaders для загрузки
- Использовать react-hook-form для форм
- Интегрировать с Chart.js для графиков

### Безопасность
- Проверка прав доступа на каждом endpoint
- Валидация всех входных данных
- Rate limiting для API
- Логирование всех действий
- Защита от CSRF атак

## Приоритеты реализации

1. **Неделя 1-2**: Товары (API + UI)
2. **Неделя 3-4**: Заказы (API + UI)
3. **Неделя 5**: Склад (API + UI)
4. **Неделя 6**: Уведомления и чат
5. **Неделя 7**: Настройки и статус магазина
6. **Неделя 8+**: Промо-акции и клиенты

## Метрики успеха

- Время загрузки страниц < 1 сек
- Uptime API > 99.9%
- Покрытие тестами > 80%
- Отсутствие критических уязвимостей
- Положительные отзывы пользователей