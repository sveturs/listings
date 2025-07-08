# Отчет о несоответствиях в именовании полей базы данных

## Дата анализа: 2025-01-08

## 1. Поля с префиксом address_, которые должны быть без префикса

### Таблица: marketplace_listings
- **address_city** → должно быть **city**
- **address_country** → должно быть **country**

## 2. Поля с использованием camelCase вместо snake_case

Не обнаружено полей с camelCase - все поля используют корректный snake_case формат.

## 3. Несогласованные id поля

Все id поля используют согласованный формат с суффиксом _id (например, user_id, listing_id, category_id).

## 4. Дополнительные замечания по согласованности

### Несогласованность в полях адреса между таблицами:

#### marketplace_listings:
- address_city
- address_country
- location (отдельное поле)
- latitude
- longitude

#### storefronts:
- address (полный адрес в одном поле)
- city (без префикса)
- country (без префикса)
- postal_code
- latitude
- longitude

#### users:
- city (без префикса)
- country (без префикса)

#### user_storefronts:
- address
- city (без префикса)
- country (без префикса)
- latitude
- longitude

### Рекомендации:
1. В таблице `marketplace_listings` следует переименовать:
   - `address_city` → `city`
   - `address_country` → `country`

2. Это обеспечит согласованность с остальными таблицами (storefronts, users, user_storefronts), где используются поля `city` и `country` без префикса `address_`.

## 5. Проверка специально запрошенных таблиц

### users ✅
Все поля корректно именованы в snake_case без излишних префиксов.

### storefronts ✅
Все поля корректно именованы в snake_case без излишних префиксов.

### storefront_orders ✅
Все поля корректно именованы в snake_case без излишних префиксов.

### payment_transactions ✅
Все поля корректно именованы в snake_case без излишних префиксов.

### marketplace_categories ✅
Все поля корректно именованы в snake_case без излишних префиксов.

### attribute_groups ✅
Все поля корректно именованы в snake_case без излишних префиксов.

### category_attributes ✅
Все поля корректно именованы в snake_case без излишних префиксов.

### listing_attribute_values ✅
Все поля корректно именованы в snake_case без излишних префиксов.

### reviews ✅
Все поля корректно именованы в snake_case без излишних префиксов.

### search_behavior_tracking
Таблица не найдена в схеме. Возможно, используется таблица `user_behavior_events` вместо нее.

### search_synonyms
Таблица не найдена в схеме. Вместо нее используется таблица `search_synonyms_config`.

## Итог

Основная проблема с именованием обнаружена только в таблице `marketplace_listings`, где поля `address_city` и `address_country` используют префикс `address_`, в то время как во всех остальных таблицах аналогичные поля называются просто `city` и `country`.

Для исправления необходимо создать миграцию, которая переименует эти поля для обеспечения согласованности во всей базе данных.