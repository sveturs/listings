# База данных: Схема и миграция

### Библиотека микросервиса для монолита

| Пакет | Файлы | Назначение | Строк кода |
|-------|-------|------------|------------|
| **pkg/client** | `client.go`, `types.go`, `converter.go` | Низкоуровневый gRPC клиент | ~400 |
| **pkg/service** | `delivery.go`, `validator.go`, `retry.go`, `cache.go` | Высокоуровневая обертка | ~600 |

**Итого библиотека**: ~1000 строк

---

## 🗄️ База данных: Текущая vs Целевая

### ТЕКУЩЕЕ: Одна БД (svetubd)

```sql
-- PostgreSQL: svetubd (монолит)

-- Все таблицы вместе:
marketplace_listings
marketplace_categories
marketplace_orders
users
user_profiles
storefronts
storefront_products
payments
payment_transactions
delivery_shipments              ⚠️ → микросервис
delivery_providers              ⚠️ → микросервис
delivery_tracking_events        ⚠️ → микросервис
delivery_category_defaults      ⚠️ → микросервис
delivery_pricing_rules          ⚠️ → микросервис
delivery_zones                  ⚠️ → микросервис
chat_messages
notifications
```

### ЦЕЛЕВОЕ: Две БД

```sql
-- PostgreSQL: svetubd (монолит)
marketplace_listings
marketplace_categories
marketplace_orders
users
user_profiles
storefronts
storefront_products
payments
payment_transactions
chat_messages
notifications
-- delivery таблицы УДАЛЕНЫ ❌
```

