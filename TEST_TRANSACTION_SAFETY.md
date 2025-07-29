# 🧪 Инструкция по тестированию транзакционной безопасности

## Подготовка

1. **Убедитесь что сервисы запущены:**
   ```bash
   # Backend на порту 3000
   lsof -i :3000
   
   # Frontend на порту 3001  
   lsof -i :3001
   
   # Если не запущены:
   /home/dim/.local/bin/kill-port-3000.sh && screen -dmS backend-3000 bash -c 'cd /data/hostel-booking-system/backend && go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'
   /home/dim/.local/bin/start-frontend-screen.sh
   ```

2. **Откройте мониторинг логов в отдельном терминале:**
   ```bash
   tail -f /tmp/backend.log | grep -i "order\|transaction\|stock"
   ```

## Тест 1: Защита от overselling

### Сценарий:
Два пользователя пытаются купить последний товар одновременно.

### Шаги:

1. **Создайте товар с количеством = 1**
   - Войдите на http://localhost:3001
   - Создайте объявление
   - Установите количество = 1
   - Запомните ID товара

2. **Откройте два браузера**
   - Браузер A: Chrome обычный режим
   - Браузер B: Chrome инкогнито или Firefox

3. **В обоих браузерах:**
   - Найдите созданный товар
   - Добавьте в корзину
   - Перейдите к оформлению заказа
   - Заполните адрес доставки

4. **Одновременное оформление:**
   - НЕ НАЖИМАЙТЕ "Оформить заказ" сразу
   - Приготовьтесь нажать в обоих браузерах
   - Нажмите "Оформить заказ" почти одновременно (с разницей 1-2 сек)

### Ожидаемый результат:
- ✅ Первый заказ успешно создан
- ❌ Второй заказ отклонен с ошибкой "Недостаточно товара на складе"
- В логах видно ROLLBACK для второй транзакции

## Тест 2: Проверка резервирований

### Через базу данных:

```bash
# Подключитесь к БД
psql "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable"

# Проверьте резервирования
SELECT 
    r.id,
    r.product_id,
    r.quantity,
    r.order_id,
    r.status,
    r.expires_at,
    p.name as product_name
FROM inventory_reservations r
JOIN storefront_products p ON p.id = r.product_id
ORDER BY r.created_at DESC
LIMIT 10;

# Проверьте остатки товаров
SELECT 
    p.id,
    p.name,
    p.stock_quantity,
    COUNT(r.id) as active_reservations,
    COALESCE(SUM(r.quantity), 0) as reserved_quantity
FROM storefront_products p
LEFT JOIN inventory_reservations r ON r.product_id = p.id AND r.status = 'reserved'
WHERE p.stock_quantity < 5
GROUP BY p.id, p.name, p.stock_quantity;
```

## Тест 3: Проверка отката транзакции

### Сценарий:
Создание заказа с несуществующим товаром

### Шаги:

1. **Измените ID товара в корзине через консоль браузера:**
   ```javascript
   // Откройте DevTools (F12)
   // Найдите в Network запрос к API корзины
   // Измените product_id на несуществующий (999999)
   ```

2. **Попробуйте оформить заказ**

### Ожидаемый результат:
- Ошибка "Товар не найден"
- В БД нет частично созданных данных
- В логах видно ROLLBACK

## Тест 4: Мониторинг производительности

### Проверьте время выполнения транзакций:

```bash
# В логах backend
grep "Creating order with transaction" /tmp/backend.log -A 20 | grep -E "took|duration|ms"

# Через БД - активные транзакции
psql "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable" -c "
SELECT 
    pid,
    now() - pg_stat_activity.query_start AS duration,
    state,
    query 
FROM pg_stat_activity 
WHERE (now() - pg_stat_activity.query_start) > interval '5 seconds'
AND state != 'idle';"
```

## Тест 5: Параллельная нагрузка

### Используйте curl для создания множества заказов:

```bash
# Создайте скрипт test_concurrent_orders.sh
cat > test_concurrent_orders.sh << 'EOF'
#!/bin/bash
TOKEN="ваш_токен_авторизации"
PRODUCT_ID=1
STOREFRONT_ID=1

for i in {1..10}; do
  curl -X POST http://localhost:3000/api/v1/orders \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "storefront_id": '$STOREFRONT_ID',
      "items": [{
        "product_id": '$PRODUCT_ID',
        "quantity": 1
      }],
      "shipping_method": "pickup",
      "shipping_address": {
        "street": "Test Street",
        "city": "Belgrade"
      }
    }' &
done
wait
EOF

chmod +x test_concurrent_orders.sh
./test_concurrent_orders.sh
```

### Ожидаемый результат:
- Только один заказ успешно создан
- Остальные получили ошибку недостатка товара
- Нет deadlock'ов в БД

## Проверка результатов

### 1. Статистика транзакций:
```sql
-- Успешные vs откаченные заказы за последний час
SELECT 
    DATE_TRUNC('minute', created_at) as minute,
    COUNT(*) FILTER (WHERE status != 'cancelled') as successful_orders,
    COUNT(*) FILTER (WHERE status = 'cancelled') as failed_orders
FROM storefront_orders
WHERE created_at > NOW() - INTERVAL '1 hour'
GROUP BY minute
ORDER BY minute DESC;
```

### 2. Проверка консистентности данных:
```sql
-- Не должно быть заказов без позиций
SELECT o.id, o.created_at
FROM storefront_orders o
LEFT JOIN storefront_order_items oi ON oi.order_id = o.id
WHERE oi.id IS NULL
AND o.created_at > NOW() - INTERVAL '1 day';

-- Не должно быть резервирований без заказов
SELECT r.*
FROM inventory_reservations r
LEFT JOIN storefront_orders o ON o.id = r.order_id
WHERE o.id IS NULL;
```

## Отладка проблем

Если что-то пошло не так:

1. **Проверьте использует ли handler новый метод:**
   ```bash
   grep -n "CreateOrderWithTx" backend/internal/proj/orders/handler/order_handler.go
   ```

2. **Проверьте логи на ошибки:**
   ```bash
   grep ERROR /tmp/backend.log | tail -20
   ```

3. **Проверьте состояние БД:**
   ```sql
   -- Активные блокировки
   SELECT * FROM pg_locks WHERE NOT granted;
   
   -- Долгие транзакции
   SELECT * FROM pg_stat_activity 
   WHERE state != 'idle' 
   AND now() - query_start > interval '30 seconds';
   ```

## Успешные индикаторы

✅ Невозможно купить больше товара чем есть на складе
✅ При ошибке все изменения откатываются
✅ Нет частичных данных в БД
✅ Производительность не упала значительно
✅ Нет deadlock'ов при параллельных запросах