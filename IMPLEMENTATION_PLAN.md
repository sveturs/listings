# План внедрения системы покупки товаров в витринах

## 📋 Текущий контекст

### ✅ Что уже реализовано:

#### Database Layer (100% готово)
- ✅ **Migration 000063**: Полная схема заказов витрин
  - `shopping_carts` - корзины покупок
  - `shopping_cart_items` - позиции корзин
  - `storefront_orders` - заказы
  - `storefront_order_items` - позиции заказов
  - `inventory_reservations` - резервирование товаров
  - Расширена `payment_transactions` для поддержки заказов

#### Models Layer (100% готово)
- ✅ **storefront_order.go**: Все модели для заказов
  - `StorefrontOrder`, `StorefrontOrderItem`
  - `ShoppingCart`, `ShoppingCartItem`
  - `InventoryReservation`
  - DTO для API запросов/ответов
- ✅ **payment_gateway.go**: Расширен для унифицированной системы
  - Добавлены `PaymentSource` (marketplace_listing/storefront_order)
  - Поля `SourceType`, `SourceID`, `StorefrontID`

#### Service Layer (100% готово)
- ✅ **OrderService**: Полная бизнес-логика заказов
  - Создание заказов из корзины или прямых позиций
  - Резервирование товаров
  - Подтверждение после оплаты
  - Отмена заказов
  - Управление статусами
- ✅ **InventoryManager**: Управление остатками
  - Резервирование/освобождение товаров
  - Списание при подтверждении
  - Очистка истекших резервирований
  - Аналитика остатков

#### Repository Layer (100% готово)
- ✅ **OrderRepository**: CRUD операции заказов
- ✅ **CartRepository**: Управление корзинами
- ✅ **InventoryRepository**: Работа с резервированиями

### 🎯 Бизнес-модель (задокументирована в B.MD)
- **Единая escrow система** для marketplace и vitrin
- **Адаптивные сроки**: 7-30 дней для marketplace, 3-7 для витрин
- **Комиссионная модель витрин**: зависит от тарифного плана
- **AllSecure интеграция**: готова в mock-режиме

---

## 🚀 План дальнейших действий

### Фаза 1: Backend Integration (3-4 дня)

#### 1.1 Расширение Payment Service (1 день)
**Статус**: Pending
**Файлы**: 
- `backend/internal/proj/payments/service/allsecure_service.go`
- `backend/internal/proj/payments/handler/payment_handler.go`

**Задачи**:
- [ ] Добавить поддержку `PaymentSource.StorefrontOrder`
- [ ] Модифицировать `CreatePayment` для заказов витрин
- [ ] Интегрировать с `OrderService.ConfirmOrder()` при успешной оплате
- [ ] Адаптировать escrow логику для витрин (3-7 дней вместо 7-30)
- [ ] Обновить webhook handler для обработки заказов

**Ключевые изменения**:
```go
// В CreatePaymentRequest добавить:
OrderID      *int64 `json:"order_id,omitempty"`
StorefrontID *int   `json:"storefront_id,omitempty"`

// В AllSecureService:
func (s *AllSecureService) CreateOrderPayment(ctx context.Context, order *models.StorefrontOrder) (*models.PaymentTransaction, error)
```

#### 1.2 HTTP Handlers для заказов (1 день)
**Статус**: Pending
**Файлы**: 
- `backend/internal/proj/orders/handler/order_handler.go`
- `backend/internal/proj/orders/handler/cart_handler.go`

**Endpoints для реализации**:
```
# Корзина
POST   /api/v1/cart/{storefront_id}/items     - добавить товар
PUT    /api/v1/cart/{storefront_id}/items/{id} - изменить количество  
DELETE /api/v1/cart/{storefront_id}/items/{id} - удалить товар
GET    /api/v1/cart/{storefront_id}           - получить корзину

# Заказы (покупатель)
POST   /api/v1/orders                         - создать заказ
GET    /api/v1/orders                         - мои заказы
GET    /api/v1/orders/{id}                    - детали заказа
PUT    /api/v1/orders/{id}/cancel             - отменить заказ

# Заказы (продавец)
GET    /api/v1/storefronts/{id}/orders        - заказы витрины
PUT    /api/v1/storefronts/{id}/orders/{id}/status - обновить статус
```

#### 1.3 Интеграция с существующими сервисами (1 день)
**Статус**: Pending
**Файлы**:
- `backend/internal/proj/storefronts/service/storefront_service.go`
- `backend/internal/proj/storefronts/service/product_service.go`

**Задачи**:
- [ ] Добавить методы получения commission rate
- [ ] Интегрировать inventory updates в product service
- [ ] Добавить валидацию активности витрины при создании заказов

#### 1.4 Dependency Injection и роутинг (0.5 дня)
**Статус**: Pending
**Файлы**:
- `backend/cmd/api/main.go`
- `backend/internal/server/server.go`

**Задачи**:
- [ ] Подключить новые handlers к роутеру
- [ ] Настроить DI для OrderService и зависимостей
- [ ] Добавить middleware для авторизации заказов

#### 1.5 Swagger документация (0.5 дня)
**Статус**: Pending
**Файлы**: Обновить swagger аннотации

---

### Фаза 2: Testing Backend (2 дня)

#### 2.1 Unit тесты (1 день)
**Статус**: Pending
**Файлы для создания**:
- `backend/internal/proj/orders/service/order_service_test.go`
- `backend/internal/proj/orders/service/inventory_manager_test.go`
- `backend/internal/proj/orders/repository/*_test.go`

**Тестовые сценарии**:
- [ ] Создание заказа из корзины
- [ ] Резервирование и списание товаров
- [ ] Расчет комиссий по тарифным планам
- [ ] Escrow сроки для разных типов витрин
- [ ] Отмена заказов и освобождение резервирований

#### 2.2 Integration тесты (1 день)
**Статус**: Pending
**Задачи**:
- [ ] Тест полного flow: корзина → заказ → оплата → подтверждение
- [ ] Тест отмены на разных стадиях
- [ ] Тест истечения резервирований
- [ ] Тест недостатка товара на складе

---

### Фаза 3: Frontend Implementation (4-5 дней)

#### 3.1 Shopping Cart Components (2 дня)
**Статус**: Pending
**Файлы для создания**:
- `frontend/svetu/src/components/cart/ShoppingCart.tsx`
- `frontend/svetu/src/components/cart/CartItem.tsx`
- `frontend/svetu/src/components/cart/CartSummary.tsx`
- `frontend/svetu/src/hooks/useShoppingCart.ts`
- `frontend/svetu/src/store/slices/cartSlice.ts`

**Функции**:
- [ ] Добавление товаров в корзину
- [ ] Изменение количества
- [ ] Удаление позиций
- [ ] Отображение итоговой суммы
- [ ] Сохранение состояния (Redux + localStorage)

#### 3.2 Checkout Flow (2 дня)
**Статус**: Pending
**Файлы для создания**:
- `frontend/svetu/src/app/[locale]/checkout/page.tsx`
- `frontend/svetu/src/components/checkout/CheckoutForm.tsx`
- `frontend/svetu/src/components/checkout/ShippingAddressForm.tsx`
- `frontend/svetu/src/components/checkout/OrderSummary.tsx`
- `frontend/svetu/src/hooks/useCheckout.ts`

**Функции**:
- [ ] Форма адреса доставки
- [ ] Выбор способа доставки
- [ ] Подтверждение заказа
- [ ] Интеграция с AllSecure payment
- [ ] Обработка ошибок

#### 3.3 Order Management (1 день)
**Статус**: Pending
**Файлы для создания**:
- `frontend/svetu/src/app/[locale]/orders/page.tsx`
- `frontend/svetu/src/app/[locale]/orders/[id]/page.tsx`
- `frontend/svetu/src/components/orders/OrderList.tsx`
- `frontend/svetu/src/components/orders/OrderDetails.tsx`
- `frontend/svetu/src/hooks/useOrders.ts`

**Функции**:
- [ ] Список заказов пользователя
- [ ] Детали заказа
- [ ] Отслеживание статуса
- [ ] Отмена заказов

---

### Фаза 4: Storefront Integration (2 дня)

#### 4.1 Product Page Integration (1 день)
**Статус**: Pending
**Файлы для модификации**:
- `frontend/svetu/src/app/[locale]/storefronts/[slug]/page.tsx`
- `frontend/svetu/src/components/storefronts/ProductCard.tsx`

**Задачи**:
- [ ] Добавить кнопки "В корзину" на товары
- [ ] Показать доступное количество с учетом резервирований
- [ ] Индикатор корзины в header витрины

#### 4.2 Seller Order Management (1 день)
**Статус**: Pending
**Файлы для создания**:
- `frontend/svetu/src/app/[locale]/storefronts/[slug]/orders/page.tsx`
- `frontend/svetu/src/components/storefronts/orders/OrderManagement.tsx`

**Функции**:
- [ ] Список заказов витрины
- [ ] Обновление статусов заказов
- [ ] Добавление tracking номеров

---

### Фаза 5: End-to-End Testing (2 дня)

#### 5.1 Manual Testing (1 день)
**Тестовые сценарии**:
- [ ] **Happy Path**: Выбор товара → корзина → checkout → оплата → получение
- [ ] **Cancellation Flow**: Отмена на разных стадиях
- [ ] **Inventory Management**: Резервирование при нехватке товара
- [ ] **Multiple Storefronts**: Отдельные корзины для разных витрин

#### 5.2 Automated E2E Tests (1 день)
**Статус**: Pending
**Файлы для создания**:
- `frontend/svetu/src/e2e/storefront-purchase.spec.ts`

**Сценарии**:
- [ ] Полный flow покупки с mock payment
- [ ] Тест корзины и checkout форм
- [ ] Тест управления заказами

---

## 🧪 План тестирования

### 1. Database Testing

#### Миграция тестинг:
```bash
# 1. Применить миграцию
cd backend && go run ./cmd/migrate/main.go up

# 2. Проверить структуру таблиц
psql -d hostel_booking_system -c "\d storefront_orders"
psql -d hostel_booking_system -c "\d shopping_carts"

# 3. Тест triggers и functions
INSERT INTO storefront_orders (storefront_id, customer_id, total_amount, commission_amount, seller_amount) 
VALUES (1, 1, 100.00, 3.00, 97.00);
```

#### Тест данных:
```sql
-- Создание тестового заказа
INSERT INTO storefront_orders (storefront_id, customer_id, subtotal_amount, total_amount, commission_amount, seller_amount, currency, status, escrow_days)
VALUES (1, 1, 100.00, 105.00, 3.15, 101.85, 'RSD', 'pending', 3);

-- Проверка автогенерации order_number и escrow_release_date
SELECT order_number, escrow_release_date FROM storefront_orders WHERE id = CURRVAL('storefront_orders_id_seq');
```

### 2. Service Layer Testing

#### OrderService тесты:
```go
func TestOrderService_CreateOrder(t *testing.T) {
    // Тест создания заказа из корзины
    // Тест резервирования товаров
    // Тест расчета комиссий
    // Тест недостатка товара
}

func TestOrderService_ConfirmOrder(t *testing.T) {
    // Тест подтверждения и списания товара
    // Тест обновления статуса escrow
}
```

#### InventoryManager тесты:
```go
func TestInventoryManager_ReserveStock(t *testing.T) {
    // Тест резервирования
    // Тест недостатка товара
    // Тест истечения резервирований
}
```

### 3. API Testing

#### Postman/Thunder Client тесты:
```json
{
  "name": "Create Storefront Order",
  "method": "POST",
  "url": "{{base_url}}/api/v1/orders",
  "body": {
    "storefront_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ],
    "shipping_address": {
      "full_name": "Test User",
      "street": "Test Street 123",
      "city": "Belgrade",
      "postal_code": "11000",
      "country": "Serbia"
    },
    "shipping_method": "standard"
  }
}
```

### 4. Frontend Testing

#### Component тесты:
```typescript
// ShoppingCart.test.tsx
test('adds item to cart', async () => {
  // Тест добавления товара
});

test('updates item quantity', async () => {
  // Тест изменения количества
});
```

#### E2E тесты с Playwright:
```typescript
test('complete purchase flow', async ({ page }) => {
  // 1. Перейти на страницу витрины
  // 2. Добавить товар в корзину
  // 3. Перейти к checkout
  // 4. Заполнить форму
  // 5. Создать заказ
  // 6. Перейти к mock payment
  // 7. Завершить оплату
  // 8. Проверить статус заказа
});
```

---

## 📊 Метрики успеха

### Технические метрики:
- [ ] **Test Coverage**: >85% для service layer
- [ ] **API Response Time**: <500ms для создания заказа
- [ ] **Database Performance**: <100ms для запросов корзины
- [ ] **Frontend Load Time**: <2s для страницы checkout

### Бизнес метрики:
- [ ] **Cart Abandonment Rate**: отслеживание через analytics
- [ ] **Payment Success Rate**: >95%
- [ ] **Inventory Accuracy**: 100% точность резервирований
- [ ] **Order Processing Time**: автоматизация 90% заказов

---

## 🔄 Rollback Plan

В случае проблем:

1. **Database Rollback**:
   ```bash
   cd backend && go run ./cmd/migrate/main.go down
   ```

2. **Feature Flag**: Временно отключить UI кнопки "В корзину"

3. **API Versioning**: Сохранить совместимость с существующими marketplace API

---

## 📅 Timeline

**Общая оценка**: 11-13 дней

- **Фаза 1 (Backend)**: 3-4 дня
- **Фаза 2 (Backend Testing)**: 2 дня  
- **Фаза 3 (Frontend)**: 4-5 дня
- **Фаза 4 (Integration)**: 2 дня
- **Фаза 5 (E2E Testing)**: 2 дня

**Критический путь**: Backend Integration → Frontend Cart → Payment Integration → E2E Testing

**Риски и митигация**:
- **AllSecure integration**: Использовать mock до получения production credentials
- **Database performance**: Индексы уже включены в миграцию
- **Frontend complexity**: Переиспользовать существующие компоненты payment системы