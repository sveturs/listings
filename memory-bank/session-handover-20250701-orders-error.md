# Session Handover: Orders System Error - 2025-07-01

## 🎯 Текущее состояние

### Выполненная работа
1. ✅ Исправлена ошибка 500 в UpdateCartItem - метод теперь реализован
2. ✅ Исправлена ссылка на checkout - использует правильный роутинг Next.js
3. ✅ Добавлены все недостающие переводы для страницы checkout
4. ✅ Создана миграция 000066_create_storefront_orders_only.up.sql
5. ✅ Таблица storefront_orders создана в БД

### 🐛 Текущая проблема
При попытке создать заказ возникает ошибка:
```
ERROR: column "user_id" of relation "storefront_orders" does not exist (SQLSTATE 42703)
```

## 📊 Анализ проблемы

### Логи ошибки
```
INFO: 2025/07/01 16:17:53.349163 order_service.go:67: Creating order%!(EXTRA string=user_id, int=7, string=storefront_id, int=4)
{"level":"error","error":"failed to create order: failed to create order: ERROR: column \"user_id\" of relation \"storefront_orders\" does not exist (SQLSTATE 42703)"}
```

### Структура таблицы storefront_orders
В миграции 000066 колонка называется `customer_id`, а не `user_id`:
```sql
CREATE TABLE IF NOT EXISTS storefront_orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(32) UNIQUE NOT NULL,
    storefront_id INTEGER REFERENCES storefronts(id) ON DELETE RESTRICT,
    customer_id INTEGER REFERENCES users(id) ON DELETE RESTRICT,  -- ⚠️ НЕ user_id!
    ...
);
```

### Вероятная причина
Код в `order_service.go` или `order_repository.go` использует имя колонки `user_id` вместо `customer_id`.

## 🔧 Что нужно исправить

1. **Проверить order_repository.go**
   - Найти SQL запрос INSERT для создания заказа
   - Заменить `user_id` на `customer_id`

2. **Проверить модель StorefrontOrder**
   - Убедиться что поле правильно замаплено
   - Возможно нужен тег `db:"customer_id"`

3. **Альтернативное решение**
   - Можно переименовать колонку в миграции на `user_id`
   - Но лучше придерживаться семантики `customer_id`

## 📁 Файлы для проверки

1. `/backend/internal/proj/orders/storage/postgres/order_repository.go`
2. `/backend/internal/domain/models/storefront_order.go`
3. `/backend/internal/proj/orders/service/order_service.go`

## 🛠️ Рабочие компоненты

### ✅ Корзина
- Добавление товаров
- Обновление количества
- Удаление товаров
- Очистка корзины

### ✅ UI компоненты
- ShoppingCartModal с полным функционалом
- Страница checkout с формой заказа
- Все переводы на месте

### ✅ Backend endpoints
- GET /api/v1/storefronts/{id}/cart
- POST /api/v1/storefronts/{id}/cart/items
- PUT /api/v1/storefronts/{id}/cart/items/{itemId}
- DELETE /api/v1/storefronts/{id}/cart/items/{itemId}
- DELETE /api/v1/storefronts/{id}/cart

### ❌ Не работает
- POST /api/v1/orders - ошибка с колонкой user_id

## 📋 Следующие шаги

1. Найти и исправить использование `user_id` → `customer_id` в коде
2. Протестировать создание заказа
3. Проверить переход на страницу оплаты
4. Убедиться что заказ сохраняется в БД

## 🔗 Полезные ссылки для тестирования

- Витрина: http://localhost:3001/storefronts/tech-store-dmitry
- Товар: http://localhost:3001/storefronts/tech-store-dmitry/products/1
- Checkout: http://localhost:3001/checkout?storefront=4

## 💡 Примечания

- Миграции 063-065 имеют некоторое дублирование (корзины)
- Миграция 063 не применялась полностью из-за конфликтов
- Создана новая миграция 066 только для таблиц заказов
- Mock payment система уже готова (миграция 064)