# 📮 Паспорт интеграции Post Express в платформу Sve Tu

## 📋 Оглавление
1. [Общая информация](#общая-информация)
2. [Коммерческие условия](#коммерческие-условия)
3. [Техническая архитектура WSP Web API](#техническая-архитектура-wsp-web-api)
4. [API методы и структуры данных](#api-методы-и-структуры-данных)
5. [Техническое задание на внедрение](#техническое-задание-на-внедрение)
6. [План миграций базы данных](#план-миграций-базы-данных)
7. [Backend реализация](#backend-реализация)
8. [Frontend компоненты](#frontend-компоненты)
9. [Сценарии использования](#сценарии-использования)
10. [Метрики и мониторинг](#метрики-и-мониторинг)
11. [Риски и митигация](#риски-и-митигация)
12. [Контакты и ресурсы](#контакты-и-ресурсы)

---

## 🏢 Общая информация

### О компании Post Express
**JP "Пошта Србије"** - национальный почтовый оператор Сербии, предоставляющий услуги экспресс-доставки Post Express по всей территории страны.

### Ключевые преимущества
- ✅ Покрытие всей территории Сербии (180+ городов)
- ✅ Доставка "Данас за сутра" (сегодня на завтра) до 19:00
- ✅ Развитая сеть почтовых отделений для самовывоза
- ✅ Поддержка откупных платежей (COD)
- ✅ SMS/Viber уведомления получателям
- ✅ Персональный координатор для корпоративных клиентов

### Документы для интеграции
- Коммерческое предложение №2025-sl от 31.07.2025
- WSP Web API документация (3 PDF файла)
- Форма заявки на заключение договора

---

## 🏢 Склад Sve Tu и услуги фулфилмента

### Описание складского комплекса
**Склад Sve Tu** - централизованный логистический центр платформы, предоставляющий услуги хранения, комплектации и отправки товаров для продавцов.

**Адрес склада:** Улица Микија Манојловића 53, 21000 Нови Сад
**График работы:** Пн-Пт 09:00-19:00, Сб 10:00-16:00
**Телефон склада:** +381 21 XXX-XXXX

### Основные услуги склада

#### 1. Для покупателей:
- ✅ **Самовывоз заказов** - бесплатный забор товаров со склада
- ✅ **Примерка перед покупкой** - для одежды и обуви
- ✅ **Консолидация заказов** - объединение товаров от разных продавцов
- ✅ **Временное хранение** - до 7 дней бесплатно

#### 2. Для продавцов:
- ✅ **Fulfillment by Sve Tu (FBS)** - полный цикл обработки заказов
- ✅ **Хранение товаров** - складские места с учетом остатков
- ✅ **Комплектация и упаковка** - профессиональная подготовка к отправке
- ✅ **Массовая отправка** - через Post Express со скидками
- ✅ **Обработка возвратов** - прием и проверка возвращенных товаров

### Интеграция с Post Express

#### Входящая логистика:
1. Продавец отправляет товары на склад Sve Tu через Post Express
2. Склад принимает и регистрирует товары в системе
3. Товары размещаются на складских местах

#### Исходящая логистика:
1. **Самовывоз:** Покупатель забирает со склада
2. **Доставка Post Express:** Отправка со склада покупателю
3. **Консолидированная отправка:** Объединение товаров разных продавцов

### Тарифы складских услуг

| Услуга | Стоимость | Примечание |
|--------|-----------|------------|
| Самовывоз покупателем | Бесплатно | В рабочее время склада |
| Хранение для продавцов | 50 RSD/м³/день | Первые 30 дней - бесплатно |
| Комплектация заказа | 30 RSD/позиция | Включает упаковку |
| Обработка возврата | 50 RSD/товар | Проверка и возврат продавцу |
| Консолидация заказов | 100 RSD | Объединение 2+ заказов |

---

## 💰 Коммерческие условия

### Тарифы на доставку "Данас за сутра"
Специальные цены для платформы Sve Tu (НДС не облагается):

| Вес посылки | Цена (RSD) | Примечание |
|-------------|------------|------------|
| до 2 кг | 340.00 | Самый популярный тариф |
| 2-5 кг | 450.00 | Для средних товаров |
| 5-10 кг | 580.00 | Для крупных товаров |
| 10-20 кг | 790.00 | Максимальный вес |

### Дополнительные услуги

#### Страхование
- **Включено:** до 15,000 RSD
- **Свыше 15,000 RSD:** +1% от суммы превышения

#### Откупные платежи (COD)
- **Комиссия:** 45 RSD за транзакцию (фиксированная)
- **Перевод средств:** в день доставки
- **Типы документов:**
  - N - платежное поручение
  - E - PosTneT денежный перевод
  - U - почтовый денежный перевод

### Условия доставки
- **Стандартные размеры:** 60x50x50 см
- **Максимальная длина:** 150 см
- **Сумма измерений:** до 300 см
- **Хранение при неудачной доставке:** 5 рабочих дней
- **Возврат недоставленных:** бесплатно

---

## 🔧 Техническая архитектура WSP Web API

### Общая архитектура
```
┌─────────────┐     HTTPS/REST    ┌──────────────┐
│   Sve Tu    │ ←───────────────→ │  WSP WebAPI  │
│   Platform  │                   │  Post Serbia │
└─────────────┘                   └──────────────┘
      │                                   │
      ↓                                   ↓
┌─────────────┐                   ┌──────────────┐
│   Backend   │                   │   Internal   │
│   Services  │                   │    Systems   │
└─────────────┘                   └──────────────┘
```

### Технические характеристики
- **Протокол:** REST over HTTPS
- **Метод:** единый POST endpoint `/Transakcija`
- **Форматы:** JSON (рекомендуется) или XML
- **Аутентификация:** Username/Password в каждом запросе
- **Идемпотентность:** через GUID `IdTransakcija`

### Базовая структура запроса
```json
{
  "StrKlijent": "{serialized_client_object}",
  "Servis": 3,                    // Всегда 3 для нашего сервиса
  "IdVrstaTranskacije": 63,       // Тип транзакции
  "TipSerijalizacije": 1,         // 1=JSON, 2=XML
  "IdTransakcija": "GUID",        // Уникальный ID запроса
  "StrIn": "{serialized_input}"   // Данные запроса
}
```

### Структура клиента (Klijent)
```json
{
  "Username": "SVE_TU_API",       // Выдается Post Express
  "Password": "secure_password",   // Выдается Post Express
  "Jezik": "LAT",                 // LAT/CYR/ENG
  "IdTipUredjaja": 2,              // 2 для веб-приложений
  "VerzijaOS": "Linux",
  "NazivUredjaja": "SVETU-API-01",
  "ModelUredjaja": "API",
  "VerzijaAplikacije": "1.0.0",
  "IPAdresa": "server_ip",
  "Geolokacija": null
}
```

---

## 📡 API методы и структуры данных

### 1. Поиск населенных пунктов (GetNaselje)
**ID транзакции:** 3

#### Запрос (GetNaseljeIn)
```json
{
  "Naziv": "Нови Сад",      // Название населенного пункта
  "BrojSlogova": 10,         // Максимум результатов
  "NacinSortiranja": 0       // Всегда 0
}
```

#### Ответ (GetNaseljeOut)
```json
{
  "Naselja": [
    {
      "Id": 1234,
      "Naziv": "NOVI SAD"
    },
    {
      "Id": 5678,
      "Naziv": "NOVI SAD - PETROVARADIN"
    }
  ]
}
```

### 2. Отслеживание посылки (TTKretanjaUsluge)
**ID транзакции:** 63

#### Запрос (TTKretanjeIn)
```json
{
  "VrstaUsluge": 1,                    // Всегда 1 для отслеживания
  "EksterniBroj": "",                  // Обычно пустой
  "PrijemniBroj": "RS123456789RS"      // Номер отслеживания
}
```

#### Ответ (TTKretanjeOut)
```json
{
  "OtkupniDokumenti": [
    {
      "Vrsta": "N",                    // Тип документа
      "Broj": "12345"                  // Номер документа
    }
  ],
  "Kretanja": [
    {
      "Status": "Примљено",           // Текстовый статус
      "Mesto": "БЕОГРАД",              // Место
      "Datum": "2025-01-15 10:30",    // Дата и время
      "StatusSifra": "PR",             // Код статуса
      "Potpisnik": "",                 // Подписант (при доставке)
      "Faza": "1",                     // Внутренняя фаза
      "MestoDo": "НОВИ САД",           // Направление
      "Retur": "",                     // "D" если возврат
      "Privremen": "",                 // "*" если временный
      "PrijemniBroj": "RS123456789RS",
      "Masa": 1500,                    // Вес в граммах
      "ImaOtkupninu": true,            // Есть откупная сумма
      "DatumKretanja": "2025-01-15T10:30:00",
      "Konacno": false                 // Финальный статус
    }
  ],
  "BrojPosiljkeOTK": "",
  "BrojPosiljkePDK": ""                // Номер возвратных документов
}
```

### 3. Коды статусов
| Код | Значение | Описание |
|-----|----------|----------|
| PR | Примљено (Received) | Посылка принята |
| OT | Отправљено (Dispatched) | Посылка отправлена |
| UR | Уручено (Delivered) | Посылка доставлена |
| IZ | Извештено (Notified) | Получатель уведомлен |

### 4. Манифест посылок (Manifest)
**ID транзакции:** 73

Используется для получения списка посылок и массовых операций.

---

## 📝 Техническое задание на внедрение

### Цели интеграции
1. Обеспечить возможность доставки через Post Express для C2C и B2C сценариев
2. Автоматизировать процесс создания и отслеживания отправлений
3. Интегрировать откупные платежи в финансовую систему платформы
4. Предоставить выбор между Post Express и другими провайдерами

### Функциональные требования

#### Для C2C (Customer-to-Customer)
- [ ] Выбор Post Express при оформлении заказа
- [ ] Расчет стоимости доставки по весу
- [ ] Поиск и выбор населенного пункта получателя
- [ ] Создание отправления после оплаты
- [ ] Отслеживание статуса посылки
- [ ] Управление откупными платежами
- [ ] Уведомления об изменении статуса

#### Для B2C (Business-to-Customer)
- [ ] Настройки Post Express для витрин
- [ ] Массовое создание отправлений
- [ ] Генерация манифестов
- [ ] Печать адресниц через EPK
- [ ] Аналитика по доставкам
- [ ] Интеграция с учетными системами витрин

### Нефункциональные требования
- **Производительность:** обработка до 100 запросов/сек
- **Доступность:** 99.9% uptime
- **Безопасность:** шифрование паролей, HTTPS only
- **Масштабируемость:** горизонтальное масштабирование
- **Логирование:** все API вызовы и ответы

---

## 🗄️ План миграций базы данных

### Миграция 001: Расширение существующих таблиц для Post Express
```sql
-- 001_post_express_extend_delivery_options.up.sql

-- Расширяем существующую таблицу storefront_delivery_options
-- Эта таблица уже содержит поля provider и provider_config (JSONB)
-- Для Post Express используем provider = 'post-express'
-- В provider_config храним:
-- {
--   "username": "SVE_TU_API",
--   "password_encrypted": "...",
--   "api_endpoint": "https://api.posta.rs/wsp",
--   "jezik": "LAT",
--   "sender_settlement_id": 1234,
--   "pickup_address": {...},
--   "cod_enabled": true,
--   "insurance_default": 15000
-- }

-- Добавляем колонку для Post Express конфигурации если её нет
ALTER TABLE storefront_delivery_options 
ADD COLUMN IF NOT EXISTS post_express_settings JSONB DEFAULT '{}'::jsonb;

-- Создаем глобальную конфигурацию Post Express
CREATE TABLE IF NOT EXISTS post_express_config (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password_encrypted VARCHAR(255) NOT NULL,
    api_endpoint VARCHAR(255) DEFAULT 'https://api.posta.rs/wsp',
    jezik VARCHAR(3) DEFAULT 'LAT',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Таблица отправлений Post Express
CREATE TABLE post_express_shipments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES marketplace_orders(id),
    storefront_order_id BIGINT REFERENCES storefront_orders(id),
    prijemni_broj VARCHAR(50) UNIQUE,           -- Номер приема от Post Express
    eksterni_broj VARCHAR(50),                  -- Внешний номер (если используется)
    
    -- Данные отправителя
    sender_name VARCHAR(200) NOT NULL,
    sender_phone VARCHAR(50) NOT NULL,
    sender_address JSONB NOT NULL,              -- {place, street, house_number, apartment}
    sender_settlement_id INTEGER,               -- ID из справочника Post Express
    
    -- Данные получателя
    receiver_name VARCHAR(200) NOT NULL,
    receiver_phone VARCHAR(50) NOT NULL,
    receiver_address JSONB NOT NULL,
    receiver_settlement_id INTEGER,
    
    -- Параметры посылки
    weight INTEGER NOT NULL,                    -- Вес в граммах
    insurance_amount DECIMAL(12,2),
    cod_amount DECIMAL(12,2),                   -- Откупная сумма
    cod_commission DECIMAL(12,2) DEFAULT 45.00,
    shipping_cost DECIMAL(12,2) NOT NULL,
    
    -- Статусы и отслеживание
    current_status VARCHAR(10),                 -- PR, OT, UR, IZ
    current_status_text VARCHAR(200),
    current_location VARCHAR(200),
    delivered_at TIMESTAMP WITH TIME ZONE,
    delivered_to VARCHAR(200),                  -- Имя получателя
    
    -- Откупные документы
    otkupni_dokument_type VARCHAR(10),          -- N, E, U
    otkupni_dokument_broj VARCHAR(50),
    
    -- Метаданные
    tracking_history JSONB,                     -- Полная история движений
    last_sync_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_post_express_shipments_order ON post_express_shipments(order_id);
CREATE INDEX idx_post_express_shipments_storefront_order ON post_express_shipments(storefront_order_id);
CREATE INDEX idx_post_express_shipments_prijemni ON post_express_shipments(prijemni_broj);
CREATE INDEX idx_post_express_shipments_status ON post_express_shipments(current_status);
CREATE INDEX idx_post_express_shipments_created ON post_express_shipments(created_at DESC);

-- Таблица справочника населенных пунктов
CREATE TABLE post_express_settlements (
    id INTEGER PRIMARY KEY,                     -- ID от Post Express
    naziv VARCHAR(200) NOT NULL,                -- Название
    ptt VARCHAR(10),                           -- Почтовый индекс
    opstina VARCHAR(100),                      -- Муниципалитет
    region VARCHAR(100),                       -- Регион
    cached_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_post_express_settlements_naziv ON post_express_settlements(naziv);
CREATE INDEX idx_post_express_settlements_ptt ON post_express_settlements(ptt);
```

### Миграция 002: Интеграция с существующей системой доставки
```sql
-- 002_post_express_integrate_delivery.up.sql

-- Используем существующую таблицу storefront_delivery_options
-- Для Post Express создаем записи с provider = 'post-express'
-- Пример вставки опции доставки Post Express для витрины:

/*
INSERT INTO storefront_delivery_options (
    storefront_id,
    name,
    description,
    provider,
    provider_config,
    base_price,
    price_per_kg,
    free_above_amount,
    max_weight_kg,
    estimated_days_min,
    estimated_days_max,
    is_active
) VALUES (
    1, -- storefront_id
    'Post Express - Danas za sutra',
    'Brza dostava Post Express u roku od 1-2 radna dana',
    'post-express',
    jsonb_build_object(
        'username', 'SVE_TU_API',
        'password_encrypted', encrypt_password('password'),
        'api_endpoint', 'https://api.posta.rs/wsp',
        'jezik', 'LAT',
        'enable_cod', true,
        'default_insurance', 15000,
        'sender_settlement_id', null,
        'sender_address', jsonb_build_object(
            'street', 'Ulica Mikija Manojlovića',
            'number', '53',
            'city', 'Novi Sad',
            'postal_code', '21000'
        ),
        'pickup_time_from', '09:00',
        'pickup_time_to', '17:00',
        'weight_tiers', jsonb_build_array(
            jsonb_build_object('max_kg', 2, 'price', 340),
            jsonb_build_object('max_kg', 5, 'price', 450),
            jsonb_build_object('max_kg', 10, 'price', 580),
            jsonb_build_object('max_kg', 20, 'price', 790)
        )
    ),
    0, -- base_price (рассчитывается динамически)
    0, -- price_per_kg (используем weight_tiers из config)
    5000, -- free_above_amount
    20, -- max_weight_kg
    1, -- estimated_days_min
    2, -- estimated_days_max
    true -- is_active
);
*/

-- Добавляем функцию для расчета стоимости доставки Post Express
CREATE OR REPLACE FUNCTION calculate_post_express_rate(
    weight_grams INTEGER,
    insurance_amount DECIMAL,
    config JSONB
) RETURNS DECIMAL AS $$
DECLARE
    base_rate DECIMAL := 0;
    insurance_fee DECIMAL := 0;
    weight_kg DECIMAL;
    weight_tiers JSONB;
    tier JSONB;
BEGIN
    weight_kg := weight_grams / 1000.0;
    weight_tiers := config->'weight_tiers';
    
    -- Определяем базовую ставку по весу
    FOR tier IN SELECT * FROM jsonb_array_elements(weight_tiers)
    LOOP
        IF weight_kg <= (tier->>'max_kg')::DECIMAL THEN
            base_rate := (tier->>'price')::DECIMAL;
            EXIT;
        END IF;
    END LOOP;
    
    -- Расчет страховки
    IF insurance_amount > 15000 THEN
        insurance_fee := (insurance_amount - 15000) * 0.01;
    END IF;
    
    RETURN base_rate + insurance_fee;
END;
$$ LANGUAGE plpgsql;
```

### Миграция 003: Логирование и аудит
```sql
-- 003_post_express_audit_log.up.sql

CREATE TABLE post_express_api_log (
    id BIGSERIAL PRIMARY KEY,
    transaction_id UUID NOT NULL,
    transaction_type INTEGER NOT NULL,
    request_data JSONB NOT NULL,
    response_data JSONB,
    response_status INTEGER,
    error_message TEXT,
    duration_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_post_express_api_log_transaction ON post_express_api_log(transaction_id);
CREATE INDEX idx_post_express_api_log_created ON post_express_api_log(created_at DESC);
CREATE INDEX idx_post_express_api_log_status ON post_express_api_log(response_status);

-- Создаем партиции на несколько месяцев вперед
CREATE TABLE post_express_api_log_2025_01 PARTITION OF post_express_api_log
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
CREATE TABLE post_express_api_log_2025_02 PARTITION OF post_express_api_log
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
CREATE TABLE post_express_api_log_2025_03 PARTITION OF post_express_api_log
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');

-- Триггер для обновления updated_at в post_express_shipments
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_post_express_shipments_updated_at 
    BEFORE UPDATE ON post_express_shipments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### Миграция 004: Складской учет и FBS (Fulfillment by Sve Tu)
```sql
-- 004_warehouse_fulfillment_tables.up.sql

-- Таблица складских мест
-- Для учета товаров на складе Sve Tu
CREATE TABLE warehouse_inventory (
    id SERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES storefront_products(id),
    storefront_id INTEGER REFERENCES storefronts(id),
    sku VARCHAR(100) NOT NULL,                   -- Артикул товара
    
    -- Количество
    quantity_available INTEGER DEFAULT 0,         -- Доступно для продажи
    quantity_reserved INTEGER DEFAULT 0,          -- Зарезервировано под заказы
    quantity_damaged INTEGER DEFAULT 0,           -- Поврежденные
    
    -- Размещение на складе
    warehouse_location VARCHAR(50),               -- Место хранения (ряд/стеллаж/ячейка)
    batch_number VARCHAR(100),                    -- Номер партии
    
    -- Даты
    received_at TIMESTAMP WITH TIME ZONE,         -- Дата поступления на склад
    expiry_date DATE,                            -- Срок годности (если применимо)
    last_inventory_check TIMESTAMP WITH TIME ZONE,
    
    -- Метаданные
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(product_id, batch_number)
);

CREATE INDEX idx_warehouse_inventory_sku ON warehouse_inventory(sku);
CREATE INDEX idx_warehouse_inventory_storefront ON warehouse_inventory(storefront_id);
CREATE INDEX idx_warehouse_inventory_location ON warehouse_inventory(warehouse_location);

-- Таблица движения товаров на складе
CREATE TABLE warehouse_movements (
    id SERIAL PRIMARY KEY,
    inventory_id INTEGER REFERENCES warehouse_inventory(id),
    movement_type VARCHAR(50) NOT NULL,          -- 'inbound', 'outbound', 'adjustment', 'return'
    quantity INTEGER NOT NULL,
    reference_type VARCHAR(50),                  -- 'order', 'return', 'adjustment'
    reference_id INTEGER,                        -- ID заказа/возврата
    
    -- Документы
    document_number VARCHAR(100),                -- Номер накладной
    post_express_tracking VARCHAR(50),           -- Номер отслеживания Post Express
    
    reason TEXT,
    performed_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_warehouse_movements_inventory ON warehouse_movements(inventory_id);
CREATE INDEX idx_warehouse_movements_type ON warehouse_movements(movement_type);
CREATE INDEX idx_warehouse_movements_reference ON warehouse_movements(reference_type, reference_id);

-- Таблица заказов на самовывоз
CREATE TABLE warehouse_pickup_orders (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES marketplace_orders(id),
    storefront_order_id BIGINT REFERENCES storefront_orders(id),
    
    -- Код самовывоза
    pickup_code VARCHAR(10) NOT NULL UNIQUE,     -- 6-значный код для получения
    qr_code_url VARCHAR(500),                    -- URL QR-кода для самовывоза
    
    -- Статусы
    status VARCHAR(50) DEFAULT 'pending',        -- 'pending', 'ready', 'picked_up', 'expired'
    ready_at TIMESTAMP WITH TIME ZONE,           -- Когда заказ готов к выдаче
    picked_up_at TIMESTAMP WITH TIME ZONE,       -- Когда забрали
    expires_at TIMESTAMP WITH TIME ZONE,         -- Срок хранения (7 дней)
    
    -- Получатель
    customer_name VARCHAR(200),
    customer_phone VARCHAR(50),
    customer_email VARCHAR(200),
    
    -- Подтверждение получения
    pickup_confirmed_by VARCHAR(200),            -- Кто выдал (сотрудник склада)
    id_document_type VARCHAR(50),                -- Тип документа (паспорт, лична карта)
    id_document_number VARCHAR(100),             -- Номер документа
    
    -- Уведомления
    notification_sent_at TIMESTAMP WITH TIME ZONE,
    reminder_sent_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_warehouse_pickup_orders_code ON warehouse_pickup_orders(pickup_code);
CREATE INDEX idx_warehouse_pickup_orders_status ON warehouse_pickup_orders(status);
CREATE INDEX idx_warehouse_pickup_orders_expires ON warehouse_pickup_orders(expires_at);

-- Таблица настроек FBS для витрин
CREATE TABLE storefront_fbs_settings (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER REFERENCES storefronts(id) UNIQUE,
    
    -- Статус FBS
    fbs_enabled BOOLEAN DEFAULT false,           -- Использовать склад Sve Tu
    auto_fulfillment BOOLEAN DEFAULT true,       -- Автоматическая обработка заказов
    
    -- Настройки хранения
    storage_tier VARCHAR(50) DEFAULT 'standard', -- 'standard', 'premium', 'economy'
    max_storage_volume DECIMAL(10,2),            -- Макс. объем в м³
    free_storage_days INTEGER DEFAULT 30,        -- Бесплатное хранение (дней)
    
    -- Настройки отправки
    default_packaging VARCHAR(50) DEFAULT 'standard',
    include_invoice BOOLEAN DEFAULT true,
    include_marketing BOOLEAN DEFAULT false,     -- Включать рекламные материалы
    
    -- Биллинг
    billing_cycle VARCHAR(50) DEFAULT 'monthly', -- 'weekly', 'monthly'
    last_billed_at TIMESTAMP WITH TIME ZONE,
    current_charges DECIMAL(12,2) DEFAULT 0,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Функция генерации кода самовывоза
CREATE OR REPLACE FUNCTION generate_pickup_code()
RETURNS VARCHAR AS $$
DECLARE
    chars VARCHAR := 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result VARCHAR := '';
    i INTEGER;
BEGIN
    FOR i IN 1..6 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Триггер для автоматического уменьшения количества при создании заказа
CREATE OR REPLACE FUNCTION decrease_inventory_on_order()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE warehouse_inventory 
    SET 
        quantity_available = quantity_available - NEW.quantity,
        quantity_reserved = quantity_reserved + NEW.quantity,
        updated_at = NOW()
    WHERE id = NEW.inventory_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

---

## 💻 Backend реализация

### Структура модуля Post Express
```
backend/
└── internal/
    └── proj/
        └── postexpress/          # Модуль Post Express по паттерну проекта
            ├── handler/
            │   ├── handler.go    # Основной хэндлер
            │   ├── routes.go     # Маршруты API
            │   └── responses.go  # Структуры ответов
            ├── service/
            │   ├── interface.go  # Интерфейсы
            │   ├── service.go    # Основная бизнес-логика
            │   ├── client.go     # WSP API клиент
            │   └── tracking.go   # Отслеживание посылок
            ├── storage/
            │   ├── interface.go  # Интерфейсы репозитория
            │   └── postgres/
            │       └── repository.go  # Работа с БД
            └── module.go          # Модуль для DI
```

### Client реализация
```go
// internal/proj/postexpress/service/client.go

package service

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type WSPClient struct {
    httpClient *http.Client
    config     *Config
    logger     Logger
}

type Config struct {
    Endpoint   string
    Username   string
    Password   string
    Language   string
    DeviceType int
    Timeout    time.Duration
}

func NewWSPClient(config *Config, logger Logger) *WSPClient {
    return &WSPClient{
        httpClient: &http.Client{
            Timeout: config.Timeout,
        },
        config: config,
        logger: logger,
    }
}

// Базовый метод для всех транзакций
func (c *WSPClient) Transaction(ctx context.Context, req *TransactionRequest) (*TransactionResponse, error) {
    // Сериализация клиента
    clientData := &ClientData{
        Username:         c.config.Username,
        Password:         c.config.Password,
        Jezik:           c.config.Language,
        IdTipUredjaja:   c.config.DeviceType,
        VerzijaOS:       "Linux",
        NazivUredjaja:   "SVETU-API",
        ModelUredjaja:   "API",
        VerzijaAplikacije: "1.0.0",
        IPAdresa:        getServerIP(),
    }
    
    clientJSON, err := json.Marshal(clientData)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal client data: %w", err)
    }
    
    // Подготовка запроса
    transReq := &TransakcijaIn{
        StrKlijent:         string(clientJSON),
        Servis:            3,
        IdVrstaTranskacije: req.TransactionType,
        TipSerijalizacije: 1, // JSON
        IdTransakcija:     generateGUID(),
        StrIn:             req.InputData,
    }
    
    // Логирование запроса
    c.logger.Debug("WSP API Request", 
        "transaction_id", transReq.IdTransakcija,
        "type", req.TransactionType)
    
    // HTTP запрос
    body, err := json.Marshal(transReq)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequestWithContext(ctx, "POST", 
        c.config.Endpoint+"/Transakcija", bytes.NewReader(body))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    
    // Выполнение запроса
    start := time.Now()
    resp, err := c.httpClient.Do(httpReq)
    duration := time.Since(start)
    
    if err != nil {
        c.logger.Error("WSP API Request failed",
            "error", err,
            "duration_ms", duration.Milliseconds())
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Парсинг ответа
    var transResp TransakcijaOut
    if err := json.NewDecoder(resp.Body).Decode(&transResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    // Проверка результата
    if transResp.Rezultat != 0 {
        var result ResultData
        if err := json.Unmarshal([]byte(transResp.StrRezultat), &result); err == nil {
            return nil, fmt.Errorf("API error: %s", result.PorukaKorisnik)
        }
        return nil, fmt.Errorf("API error code: %d", transResp.Rezultat)
    }
    
    c.logger.Info("WSP API Request successful",
        "transaction_id", transReq.IdTransakcija,
        "duration_ms", duration.Milliseconds())
    
    return &TransactionResponse{
        OutputData: transResp.StrOut,
        Result:     transResp.StrRezultat,
    }, nil
}

// Метод поиска населенных пунктов
func (c *WSPClient) SearchSettlements(ctx context.Context, name string, limit int) ([]*Settlement, error) {
    input := &GetNaseljeIn{
        Naziv:          name,
        BrojSlogova:    limit,
        NacinSortiranja: 0,
    }
    
    inputJSON, err := json.Marshal(input)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.Transaction(ctx, &TransactionRequest{
        TransactionType: 3, // GetNaselje
        InputData:      string(inputJSON),
    })
    if err != nil {
        return nil, err
    }
    
    var output GetNaseljeOut
    if err := json.Unmarshal([]byte(resp.OutputData), &output); err != nil {
        return nil, fmt.Errorf("failed to parse settlements: %w", err)
    }
    
    return output.Naselja, nil
}

// Метод отслеживания посылки
func (c *WSPClient) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
    input := &TTKretanjeIn{
        VrstaUsluge:   1,
        EksterniBroj:  "",
        PrijemniBroj:  trackingNumber,
    }
    
    inputJSON, err := json.Marshal(input)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.Transaction(ctx, &TransactionRequest{
        TransactionType: 63, // TTKretanjaUsluge
        InputData:      string(inputJSON),
    })
    if err != nil {
        return nil, err
    }
    
    var output TTKretanjeOut
    if err := json.Unmarshal([]byte(resp.OutputData), &output); err != nil {
        return nil, fmt.Errorf("failed to parse tracking info: %w", err)
    }
    
    return &TrackingInfo{
        Movements:         output.Kretanja,
        RedemptionDocs:   output.OtkupniDokumenti,
        ReturnDocNumber:  output.BrojPosiljkePDK,
    }, nil
}
```

### Service реализация
```go
// internal/proj/postexpress/service/service.go

package service

import (
    "context"
    "fmt"
    "time"
    
    "svetu/internal/domain/models"
)

type PostExpressService struct {
    client     *WSPClient
    repo       Repository
    logger     Logger
    calculator *RateCalculator
}

func NewPostExpressService(client *WSPClient, repo Repository, logger Logger) *PostExpressService {
    return &PostExpressService{
        client:     client,
        repo:       repo,
        logger:     logger,
        calculator: NewRateCalculator(),
    }
}

// Расчет стоимости доставки
func (s *PostExpressService) CalculateShippingRate(weight int, insuranceAmount float64) (*ShippingRate, error) {
    // Определение весовой категории
    var baseRate float64
    switch {
    case weight <= 2000:
        baseRate = 340.00
    case weight <= 5000:
        baseRate = 450.00
    case weight <= 10000:
        baseRate = 580.00
    case weight <= 20000:
        baseRate = 790.00
    default:
        return nil, fmt.Errorf("weight %d exceeds maximum of 20kg", weight)
    }
    
    // Расчет страховки
    var insuranceFee float64
    if insuranceAmount > 15000 {
        insuranceFee = (insuranceAmount - 15000) * 0.01
    }
    
    return &ShippingRate{
        BaseRate:      baseRate,
        InsuranceFee:  insuranceFee,
        CODFee:        45.00, // Фиксированная комиссия
        TotalRate:     baseRate + insuranceFee,
    }, nil
}

// Создание отправления для C2C
func (s *PostExpressService) CreateC2CShipment(ctx context.Context, req *CreateShipmentRequest) (*Shipment, error) {
    // Валидация адресов
    if err := s.validateAddresses(ctx, req); err != nil {
        return nil, fmt.Errorf("address validation failed: %w", err)
    }
    
    // Расчет стоимости
    rate, err := s.CalculateShippingRate(req.Weight, req.InsuranceAmount)
    if err != nil {
        return nil, err
    }
    
    // Создание записи в БД
    shipment := &models.PostExpressShipment{
        OrderID:            req.OrderID,
        SenderName:         req.Sender.FullName,
        SenderPhone:        req.Sender.Phone,
        SenderAddress:      req.Sender.Address,
        SenderSettlementID: req.Sender.SettlementID,
        ReceiverName:       req.Receiver.FullName,
        ReceiverPhone:      req.Receiver.Phone,
        ReceiverAddress:    req.Receiver.Address,
        ReceiverSettlementID: req.Receiver.SettlementID,
        Weight:             req.Weight,
        InsuranceAmount:    req.InsuranceAmount,
        CODAmount:          req.CODAmount,
        ShippingCost:       rate.TotalRate,
        CurrentStatus:      "PENDING",
    }
    
    // Сохранение в БД
    if err := s.repo.CreateShipment(ctx, shipment); err != nil {
        return nil, fmt.Errorf("failed to save shipment: %w", err)
    }
    
    // TODO: Интеграция с EPK/Web Express для получения prijemni_broj
    // Это требует дополнительной документации от Post Express
    
    s.logger.Info("C2C shipment created",
        "shipment_id", shipment.ID,
        "order_id", req.OrderID)
    
    return s.mapToShipment(shipment), nil
}

// Массовое создание отправлений для B2C
func (s *PostExpressService) CreateB2CManifest(ctx context.Context, storefrontID int, orderIDs []int) (*Manifest, error) {
    // Получение настроек витрины
    settings, err := s.repo.GetStorefrontSettings(ctx, storefrontID)
    if err != nil {
        return nil, fmt.Errorf("failed to get storefront settings: %w", err)
    }
    
    if !settings.IsActive {
        return nil, fmt.Errorf("Post Express is not active for storefront %d", storefrontID)
    }
    
    // Получение заказов
    orders, err := s.repo.GetOrdersByIDs(ctx, orderIDs)
    if err != nil {
        return nil, fmt.Errorf("failed to get orders: %w", err)
    }
    
    // Создание манифеста
    manifest := &Manifest{
        StorefrontID: storefrontID,
        OrderCount:   len(orders),
        CreatedAt:    time.Now(),
        Shipments:    make([]*Shipment, 0, len(orders)),
    }
    
    // Создание отправлений для каждого заказа
    for _, order := range orders {
        shipment, err := s.createB2CShipment(ctx, settings, order)
        if err != nil {
            s.logger.Error("Failed to create B2C shipment",
                "order_id", order.ID,
                "error", err)
            continue
        }
        manifest.Shipments = append(manifest.Shipments, shipment)
    }
    
    // TODO: Отправка манифеста в Post Express через API
    
    s.logger.Info("B2C manifest created",
        "storefront_id", storefrontID,
        "shipment_count", len(manifest.Shipments))
    
    return manifest, nil
}

// Синхронизация статусов отправлений
func (s *PostExpressService) SyncShipmentStatuses(ctx context.Context) error {
    // Получение активных отправлений
    shipments, err := s.repo.GetActiveShipments(ctx)
    if err != nil {
        return fmt.Errorf("failed to get active shipments: %w", err)
    }
    
    s.logger.Info("Starting status sync", "shipment_count", len(shipments))
    
    for _, shipment := range shipments {
        if shipment.PrijemniBroj == "" {
            continue // Пропускаем если нет номера отслеживания
        }
        
        // Получение информации об отслеживании
        tracking, err := s.client.TrackShipment(ctx, shipment.PrijemniBroj)
        if err != nil {
            s.logger.Error("Failed to track shipment",
                "tracking_number", shipment.PrijemniBroj,
                "error", err)
            continue
        }
        
        // Обновление статуса
        if err := s.updateShipmentStatus(ctx, shipment, tracking); err != nil {
            s.logger.Error("Failed to update shipment status",
                "shipment_id", shipment.ID,
                "error", err)
            continue
        }
    }
    
    return nil
}

// Обновление статуса отправления
func (s *PostExpressService) updateShipmentStatus(ctx context.Context, shipment *models.PostExpressShipment, tracking *TrackingInfo) error {
    if len(tracking.Movements) == 0 {
        return nil // Нет движений
    }
    
    // Последнее движение
    lastMovement := tracking.Movements[len(tracking.Movements)-1]
    
    // Обновление полей
    shipment.CurrentStatus = lastMovement.StatusSifra
    shipment.CurrentStatusText = lastMovement.Status
    shipment.CurrentLocation = lastMovement.Mesto
    shipment.TrackingHistory = tracking.Movements
    shipment.LastSyncAt = time.Now()
    
    // Проверка доставки
    if lastMovement.Konacno && lastMovement.StatusSifra == "UR" {
        shipment.DeliveredAt = &lastMovement.DatumKretanja
        shipment.DeliveredTo = lastMovement.Potpisnik
        
        // Обработка откупных документов
        if len(tracking.RedemptionDocs) > 0 {
            doc := tracking.RedemptionDocs[0]
            shipment.OtkupniDokumentType = doc.Vrsta
            shipment.OtkupniDokumentBroj = doc.Broj
        }
    }
    
    // Сохранение в БД
    if err := s.repo.UpdateShipment(ctx, shipment); err != nil {
        return fmt.Errorf("failed to update shipment: %w", err)
    }
    
    // Отправка уведомлений
    if lastMovement.StatusSifra != shipment.CurrentStatus {
        s.sendStatusNotification(shipment, lastMovement)
    }
    
    return nil
}

// Валидация адресов через API
func (s *PostExpressService) validateAddresses(ctx context.Context, req *CreateShipmentRequest) error {
    // Проверка населенного пункта отправителя
    if req.Sender.SettlementID == 0 {
        settlements, err := s.client.SearchSettlements(ctx, req.Sender.City, 1)
        if err != nil {
            return fmt.Errorf("failed to validate sender city: %w", err)
        }
        if len(settlements) == 0 {
            return fmt.Errorf("sender city '%s' not found", req.Sender.City)
        }
        req.Sender.SettlementID = settlements[0].Id
    }
    
    // Проверка населенного пункта получателя
    if req.Receiver.SettlementID == 0 {
        settlements, err := s.client.SearchSettlements(ctx, req.Receiver.City, 1)
        if err != nil {
            return fmt.Errorf("failed to validate receiver city: %w", err)
        }
        if len(settlements) == 0 {
            return fmt.Errorf("receiver city '%s' not found", req.Receiver.City)
        }
        req.Receiver.SettlementID = settlements[0].Id
    }
    
    return nil
}
```

### API Handlers
```go
// internal/proj/postexpress/handler/handler.go

package handler

import (
    "github.com/gofiber/fiber/v2"
)

type Handler struct {
    service *PostExpressService
}

// @Summary Calculate Post Express shipping rate
// @Tags delivery-post-express
// @Accept json
// @Produce json
// @Param request body CalculateRateRequest true "Rate calculation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=ShippingRate}
// @Router /api/v1/postexpress/calculate-rate [post]
func (h *Handler) CalculateRate(c *fiber.Ctx) error {
    var req CalculateRateRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    rate, err := h.service.CalculateShippingRate(req.Weight, req.InsuranceAmount)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "data":    rate,
    })
}

// @Summary Search Post Express settlements
// @Tags delivery-post-express
// @Accept json
// @Produce json
// @Param query query string true "Settlement name"
// @Param limit query int false "Result limit" default(10)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]Settlement}
// @Router /api/v1/postexpress/settlements [get]
func (h *Handler) SearchSettlements(c *fiber.Ctx) error {
    query := c.Query("query")
    if query == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Query parameter is required",
        })
    }
    
    limit := c.QueryInt("limit", 10)
    
    settlements, err := h.service.client.SearchSettlements(c.Context(), query, limit)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to search settlements",
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "data":    settlements,
    })
}

// @Summary Track Post Express shipment
// @Tags delivery-post-express
// @Accept json
// @Produce json
// @Param tracking_number path string true "Tracking number"
// @Success 200 {object} utils.SuccessResponseSwag{data=TrackingInfo}
// @Router /api/v1/postexpress/track/{tracking_number} [get]
func (h *Handler) TrackShipment(c *fiber.Ctx) error {
    trackingNumber := c.Params("tracking_number")
    if trackingNumber == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Tracking number is required",
        })
    }
    
    tracking, err := h.service.client.TrackShipment(c.Context(), trackingNumber)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to track shipment",
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "data":    tracking,
    })
}

// @Summary Create Post Express C2C shipment
// @Tags delivery-post-express
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateShipmentRequest true "Shipment details"
// @Success 201 {object} utils.SuccessResponseSwag{data=Shipment}
// @Router /api/v1/postexpress/shipments [post]
func (h *Handler) CreateShipment(c *fiber.Ctx) error {
    var req CreateShipmentRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    // Получение user ID из контекста
    userID := c.Locals("user_id").(int)
    req.UserID = userID
    
    shipment, err := h.service.CreateC2CShipment(c.Context(), &req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(201).JSON(fiber.Map{
        "success": true,
        "data":    shipment,
    })
}

// @Summary Create Post Express B2C manifest
// @Tags delivery-post-express
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param storefront_id path int true "Storefront ID"
// @Param request body CreateManifestRequest true "Order IDs"
// @Success 201 {object} utils.SuccessResponseSwag{data=Manifest}
// @Router /api/v1/postexpress/storefronts/{storefront_id}/manifests [post]
func (h *Handler) CreateManifest(c *fiber.Ctx) error {
    storefrontID, err := c.ParamsInt("storefront_id")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid storefront ID",
        })
    }
    
    var req CreateManifestRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    manifest, err := h.service.CreateB2CManifest(c.Context(), storefrontID, req.OrderIDs)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(201).JSON(fiber.Map{
        "success": true,
        "data":    manifest,
    })
}

// Регистрация роутов 
// internal/proj/postexpress/handler/routes.go
func RegisterRoutes(app *fiber.App, handler *Handler) {
    api := app.Group("/api/v1/postexpress")
    
    // Публичные endpoints
    api.Get("/settlements", handler.SearchSettlements)
    api.Post("/calculate-rate", handler.CalculateRate)
    api.Get("/track/:tracking_number", handler.TrackShipment)
    
    // Защищенные endpoints
    protected := api.Use(requireAuth)
    protected.Post("/shipments", handler.CreateShipment)
    protected.Post("/storefronts/:storefront_id/manifests", handler.CreateManifest)
}
```

---

## 🎨 Frontend компоненты

### Компонент выбора населенного пункта
```typescript
// frontend/svetu/src/components/delivery/PostExpressSettlementSelector.tsx

import React, { useState, useCallback, useEffect } from 'react';
import { useDebounce } from '@/hooks/useDebounce';
import { MapPin, Search } from 'lucide-react';
import { useTranslations } from 'next-intl';

interface Settlement {
  id: number;
  naziv: string;
}

interface Props {
  value?: Settlement;
  onChange: (settlement: Settlement) => void;
  placeholder?: string;
  required?: boolean;
}

export default function PostExpressSettlementSelector({
  value,
  onChange,
  placeholder,
  required = false
}: Props) {
  const t = useTranslations('delivery');
  const [query, setQuery] = useState('');
  const [settlements, setSettlements] = useState<Settlement[]>([]);
  const [loading, setLoading] = useState(false);
  const [showDropdown, setShowDropdown] = useState(false);
  
  const debouncedQuery = useDebounce(query, 300);
  
  // Поиск населенных пунктов
  const searchSettlements = useCallback(async (searchQuery: string) => {
    if (searchQuery.length < 2) {
      setSettlements([]);
      return;
    }
    
    setLoading(true);
    try {
      const response = await fetch(
        `/api/v1/postexpress/settlements?query=${encodeURIComponent(searchQuery)}&limit=10`
      );
      const data = await response.json();
      
      if (data.success) {
        setSettlements(data.data);
        setShowDropdown(true);
      }
    } catch (error) {
      console.error('Failed to search settlements:', error);
    } finally {
      setLoading(false);
    }
  }, []);
  
  useEffect(() => {
    if (debouncedQuery) {
      searchSettlements(debouncedQuery);
    } else {
      setSettlements([]);
      setShowDropdown(false);
    }
  }, [debouncedQuery, searchSettlements]);
  
  const handleSelect = (settlement: Settlement) => {
    onChange(settlement);
    setQuery(settlement.naziv);
    setShowDropdown(false);
  };
  
  return (
    <div className="relative">
      <div className="relative">
        <input
          type="text"
          value={value ? value.naziv : query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={() => query && setShowDropdown(true)}
          placeholder={placeholder || t('enterSettlement')}
          required={required}
          className="input input-bordered w-full pl-10"
        />
        <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/50" />
        {loading && (
          <div className="absolute right-3 top-1/2 -translate-y-1/2">
            <span className="loading loading-spinner loading-sm"></span>
          </div>
        )}
      </div>
      
      {showDropdown && settlements.length > 0 && (
        <div className="absolute z-50 w-full mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-60 overflow-auto">
          {settlements.map((settlement) => (
            <button
              key={settlement.id}
              onClick={() => handleSelect(settlement)}
              className="w-full px-4 py-2 text-left hover:bg-base-200 flex items-center gap-2 transition-colors"
            >
              <MapPin className="w-4 h-4 text-base-content/50" />
              <span>{settlement.naziv}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
```

### Компонент расчета стоимости доставки
```typescript
// frontend/svetu/src/components/delivery/PostExpressRateCalculator.tsx

import React, { useState, useEffect } from 'react';
import { Truck, ShieldCheck, Banknote } from 'lucide-react';
import { useTranslations, useLocale } from 'next-intl';

interface ShippingRate {
  baseRate: number;
  insuranceFee: number;
  codFee: number;
  totalRate: number;
}

interface Props {
  weight: number;
  insuranceAmount?: number;
  hasCOD?: boolean;
  onRateCalculated?: (rate: ShippingRate) => void;
}

export default function PostExpressRateCalculator({
  weight,
  insuranceAmount = 0,
  hasCOD = false,
  onRateCalculated
}: Props) {
  const t = useTranslations('delivery');
  const locale = useLocale();
  const [rate, setRate] = useState<ShippingRate | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  useEffect(() => {
    calculateRate();
  }, [weight, insuranceAmount, hasCOD]);
  
  const calculateRate = async () => {
    if (weight <= 0) {
      setRate(null);
      return;
    }
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/api/v1/postexpress/calculate-rate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          weight: weight * 1000, // Convert to grams
          insuranceAmount
        })
      });
      
      const data = await response.json();
      
      if (data.success) {
        const calculatedRate = {
          ...data.data,
          codFee: hasCOD ? 45 : 0,
          totalRate: data.data.totalRate + (hasCOD ? 45 : 0)
        };
        setRate(calculatedRate);
        onRateCalculated?.(calculatedRate);
      } else {
        setError(data.error);
      }
    } catch (err) {
      setError(t('calculationError'));
    } finally {
      setLoading(false);
    }
  };
  
  if (loading) {
    return (
      <div className="flex justify-center p-4">
        <span className="loading loading-spinner"></span>
      </div>
    );
  }
  
  if (error) {
    return (
      <div className="alert alert-error">
        <span>{error}</span>
      </div>
    );
  }
  
  if (!rate) {
    return null;
  }
  
  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h3 className="card-title flex items-center gap-2">
          <Truck className="w-5 h-5" />
          {t('postExpressDelivery')}
        </h3>
        
        <div className="space-y-3">
          {/* Основная стоимость */}
          <div className="flex justify-between items-center">
            <span className="text-base-content/70">
              {t('basePrice', { weight })}
            </span>
            <span className="font-semibold">{rate.baseRate.toFixed(2)} RSD</span>
          </div>
          
          {/* Страхование */}
          {insuranceAmount > 0 && (
            <div className="flex justify-between items-center">
              <span className="text-base-content/70 flex items-center gap-1">
                <ShieldCheck className="w-4 h-4" />
                {t('insurance')}
              </span>
              <span className="font-semibold">
                {rate.insuranceFee > 0 
                  ? `+${rate.insuranceFee.toFixed(2)} RSD`
                  : t('included')
                }
              </span>
            </div>
          )}
          
          {/* Откупнина */}
          {hasCOD && (
            <div className="flex justify-between items-center">
              <span className="text-base-content/70 flex items-center gap-1">
                <Banknote className="w-4 h-4" />
                {t('cod')}
              </span>
              <span className="font-semibold">+{rate.codFee.toFixed(2)} RSD</span>
            </div>
          )}
          
          <div className="divider"></div>
          
          {/* Итого */}
          <div className="flex justify-between items-center text-lg">
            <span className="font-bold">{t('total')}:</span>
            <span className="font-bold text-primary">
              {rate.totalRate.toFixed(2)} RSD
            </span>
          </div>
        </div>
        
        <div className="text-xs text-base-content/60 mt-4">
          * {t('deliveryTime')}: 1-2 {t('workingDays')}
          <br />
          * {t('insuranceInfo', { amount: '15,000 RSD' })}
        </div>
      </div>
    </div>
  );
}
```

### Компонент отслеживания посылки
```typescript
// frontend/svetu/src/components/delivery/PostExpressTracker.tsx

import React, { useState } from 'react';
import { MapPin, Clock, CheckCircle } from 'lucide-react';
import { useTranslations, useLocale } from 'next-intl';

interface Movement {
  status: string;
  mesto: string;
  datum: string;
  statusSifra: string;
  potpisnik?: string;
  konacno: boolean;
}

interface TrackingInfo {
  movements: Movement[];
  otkupniDokumenti?: Array<{
    vrsta: string;
    broj: string;
  }>;
}

export default function PostExpressTracker() {
  const t = useTranslations('delivery');
  const locale = useLocale();
  const [trackingNumber, setTrackingNumber] = useState('');
  const [tracking, setTracking] = useState<TrackingInfo | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const handleTrack = async () => {
    if (!trackingNumber) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(
        `/api/v1/postexpress/track/${trackingNumber}`
      );
      const data = await response.json();
      
      if (data.success) {
        setTracking(data.data);
      } else {
        setError(t('packageNotFound'));
      }
    } catch (err) {
      setError(t('trackingError'));
    } finally {
      setLoading(false);
    }
  };
  
  const getStatusColor = (statusCode: string) => {
    switch (statusCode) {
      case 'PR': return 'badge-info';
      case 'OT': return 'badge-warning';
      case 'UR': return 'badge-success';
      case 'IZ': return 'badge-primary';
      default: return 'badge-ghost';
    }
  };
  
  const getStatusText = (statusCode: string) => {
    return t(`status.${statusCode}`, { defaultValue: statusCode });
  };
  
  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h2 className="card-title">
          <MapPin className="w-5 h-5" />
          {t('trackPostExpressPackage')}
        </h2>
        
        {/* Форма поиска */}
        <div className="form-control">
          <div className="input-group">
            <input
              type="text"
              placeholder={t('enterTrackingNumber')}
              className="input input-bordered flex-1"
              value={trackingNumber}
              onChange={(e) => setTrackingNumber(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleTrack()}
            />
            <button 
              className="btn btn-primary"
              onClick={handleTrack}
              disabled={loading || !trackingNumber}
            >
              {loading ? (
                <span className="loading loading-spinner"></span>
              ) : (
                t('search')
              )}
            </button>
          </div>
        </div>
        
        {error && (
          <div className="alert alert-error mt-4">
            <span>{error}</span>
          </div>
        )}
        
        {tracking && (
          <div className="mt-6">
            {/* Timeline статусов */}
            <div className="space-y-4">
              {tracking.movements.map((movement, index) => (
                <div key={index} className="flex gap-4">
                  <div className="flex flex-col items-center">
                    <div className={`w-10 h-10 rounded-full flex items-center justify-center ${
                      movement.konacno ? 'bg-success text-success-content' : 'bg-base-300'
                    }`}>
                      {movement.konacno ? (
                        <CheckCircle className="w-5 h-5" />
                      ) : (
                        <Clock className="w-5 h-5" />
                      )}
                    </div>
                    {index < tracking.movements.length - 1 && (
                      <div className="w-0.5 h-16 bg-base-300"></div>
                    )}
                  </div>
                  
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <span className={`badge ${getStatusColor(movement.statusSifra)}`}>
                        {getStatusText(movement.statusSifra)}
                      </span>
                      <span className="text-sm text-base-content/60">
                        {movement.datum}
                      </span>
                    </div>
                    <div className="font-medium">{movement.mesto}</div>
                    {movement.potpisnik && (
                      <div className="text-sm text-base-content/70">
                        {t('receivedBy')}: {movement.potpisnik}
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
            
            {/* Откупные документы */}
            {tracking.otkupniDokumenti && tracking.otkupniDokumenti.length > 0 && (
              <div className="mt-6 p-4 bg-base-200 rounded-lg">
                <h4 className="font-semibold mb-2">{t('codDocuments')}</h4>
                {tracking.otkupniDokumenti.map((doc, index) => (
                  <div key={index} className="text-sm">
                    <span className="font-medium">
                      {t(`documentType.${doc.vrsta}`)}
                    </span>
                    : {doc.broj}
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
```

### Модульные файлы переводов

Создать модульные файлы переводов для всех поддерживаемых локалей:

#### frontend/svetu/src/messages/en/delivery.json
```json
{
  "postExpressDelivery": "Post Express Delivery",
  "enterSettlement": "Enter settlement name",
  "trackPostExpressPackage": "Track Post Express Package",
  "enterTrackingNumber": "Enter tracking number",
  "search": "Search",
  "basePrice": "Base price ({weight}kg)",
  "insurance": "Insurance",
  "included": "Included",
  "cod": "Cash on delivery",
  "total": "Total",
  "deliveryTime": "Delivery time",
  "workingDays": "working days",
  "insuranceInfo": "Insurance up to {amount} included in price",
  "calculationError": "Error calculating delivery price",
  "packageNotFound": "Package not found",
  "trackingError": "Error tracking package",
  "receivedBy": "Received by",
  "codDocuments": "COD Documents",
  "status": {
    "PR": "Received",
    "OT": "In transit",
    "UR": "Delivered",
    "IZ": "Ready for pickup"
  },
  "documentType": {
    "N": "Payment order",
    "E": "PosTneT transfer",
    "U": "Postal money order"
  },
  "deliveryOptions": "Delivery Options",
  "selectDeliveryMethod": "Select delivery method",
  "estimatedDelivery": "Estimated delivery",
  "freeShipping": "Free shipping",
  "freeShippingAbove": "Free shipping for orders above {amount}",
  "todayForTomorrow": "Today for tomorrow",
  "cashOnDelivery": "Cash on delivery",
  "codFee": "COD fee: {amount}",
  "deliveryInfo": "Delivery Information",
  "smsNotification": "SMS notification about package status",
  "storageDays": "Storage for {days} working days in case of failed delivery",
  "freeReturn": "Free return of undelivered packages"
}
```

#### frontend/svetu/src/messages/ru/delivery.json
```json
{
  "postExpressDelivery": "Доставка Post Express",
  "enterSettlement": "Введите название населенного пункта",
  "trackPostExpressPackage": "Отслеживание посылки Post Express",
  "enterTrackingNumber": "Введите номер отслеживания",
  "search": "Поиск",
  "basePrice": "Базовая цена ({weight}кг)",
  "insurance": "Страхование",
  "included": "Включено",
  "cod": "Наложенный платеж",
  "total": "Итого",
  "deliveryTime": "Срок доставки",
  "workingDays": "рабочих дней",
  "insuranceInfo": "Страхование до {amount} включено в цену",
  "calculationError": "Ошибка при расчете стоимости доставки",
  "packageNotFound": "Посылка не найдена",
  "trackingError": "Ошибка при отслеживании посылки",
  "receivedBy": "Получил",
  "codDocuments": "Документы наложенного платежа",
  "status": {
    "PR": "Принято",
    "OT": "В пути",
    "UR": "Доставлено",
    "IZ": "Готово к получению"
  },
  "documentType": {
    "N": "Платежное поручение",
    "E": "Перевод PosTneT",
    "U": "Почтовый перевод"
  },
  "deliveryOptions": "Варианты доставки",
  "selectDeliveryMethod": "Выберите способ доставки",
  "estimatedDelivery": "Ориентировочная доставка",
  "freeShipping": "Бесплатная доставка",
  "freeShippingAbove": "Бесплатная доставка для заказов от {amount}",
  "todayForTomorrow": "Сегодня на завтра",
  "cashOnDelivery": "Оплата при получении",
  "codFee": "Комиссия за наложенный платеж: {amount}",
  "deliveryInfo": "Информация о доставке",
  "smsNotification": "SMS уведомление о статусе посылки",
  "storageDays": "Хранение {days} рабочих дней в случае неудачной доставки",
  "freeReturn": "Бесплатный возврат недоставленных посылок"
}
```

#### frontend/svetu/src/messages/sr/delivery.json
```json
{
  "postExpressDelivery": "Post Express dostava",
  "enterSettlement": "Unesite naziv mesta",
  "trackPostExpressPackage": "Praćenje Post Express pošiljke",
  "enterTrackingNumber": "Unesite broj pošiljke",
  "search": "Pretraži",
  "basePrice": "Osnovna cena ({weight}kg)",
  "insurance": "Osiguranje",
  "included": "Uključeno",
  "cod": "Otkupnina",
  "total": "Ukupno",
  "deliveryTime": "Rok dostave",
  "workingDays": "radna dana",
  "insuranceInfo": "Osiguranje do {amount} uključeno u cenu",
  "calculationError": "Greška pri izračunavanju cene dostave",
  "packageNotFound": "Pošiljka nije pronađena",
  "trackingError": "Greška pri praćenju pošiljke",
  "receivedBy": "Primio",
  "codDocuments": "Otkupni dokumenti",
  "status": {
    "PR": "Primljeno",
    "OT": "U transportu",
    "UR": "Uručeno",
    "IZ": "Za preuzimanje"
  },
  "documentType": {
    "N": "Nalog za uplatu",
    "E": "PosTneT uputnica",
    "U": "Poštanska uputnica"
  },
  "deliveryOptions": "Opcije dostave",
  "selectDeliveryMethod": "Izaberite način dostave",
  "estimatedDelivery": "Procenjena dostava",
  "freeShipping": "Besplatna dostava",
  "freeShippingAbove": "Besplatna dostava za porudžbine preko {amount}",
  "todayForTomorrow": "Danas za sutra",
  "cashOnDelivery": "Plaćanje pouzećem",
  "codFee": "Provizija za otkupninu: {amount}",
  "deliveryInfo": "Informacije o dostavi",
  "smsNotification": "SMS obaveštenje o statusu pošiljke",
  "storageDays": "Čuvanje {days} radnih dana u slučaju neuspešne dostave",
  "freeReturn": "Besplatan povrat nedostavljenih pošiljaka"
}
```

### Интеграция в процесс оформления заказа
```typescript
// frontend/svetu/src/components/checkout/PostExpressShipping.tsx

import React, { useState, useEffect } from 'react';
import PostExpressSettlementSelector from '../delivery/PostExpressSettlementSelector';
import PostExpressRateCalculator from '../delivery/PostExpressRateCalculator';
import { Truck } from 'lucide-react';
import { useTranslations, useLocale } from 'next-intl';

interface Props {
  weight: number;
  value: number;
  onShippingSelected: (shipping: ShippingDetails) => void;
}

interface ShippingDetails {
  provider: 'post-express';
  cost: number;
  settlementId: number;
  settlementName: string;
  estimatedDays: string;
}

export default function PostExpressShipping({ weight, value, onShippingSelected }: Props) {
  const t = useTranslations('delivery');
  const locale = useLocale();
  const [settlement, setSettlement] = useState<Settlement | null>(null);
  const [rate, setRate] = useState<ShippingRate | null>(null);
  const [useCOD, setUseCOD] = useState(false);
  
  useEffect(() => {
    if (settlement && rate) {
      onShippingSelected({
        provider: 'post-express',
        cost: rate.totalRate,
        settlementId: settlement.id,
        settlementName: settlement.naziv,
        estimatedDays: `1-2 ${t('workingDays')}`
      });
    }
  }, [settlement, rate, onShippingSelected]);
  
  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2 text-lg font-semibold">
        <Truck className="w-5 h-5" />
        {t('postExpressDelivery')}
      </div>
      
      {/* Выбор населенного пункта */}
      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('deliveryLocation')}</span>
        </label>
        <PostExpressSettlementSelector
          value={settlement}
          onChange={setSettlement}
          required
        />
      </div>
      
      {/* Опция откупнины */}
      <div className="form-control">
        <label className="label cursor-pointer">
          <span className="label-text">{t('cashOnDelivery')}</span>
          <input
            type="checkbox"
            className="checkbox checkbox-primary"
            checked={useCOD}
            onChange={(e) => setUseCOD(e.target.checked)}
          />
        </label>
        {useCOD && (
          <span className="text-xs text-base-content/60">
            {t('codFee', { amount: '45 RSD' })}
          </span>
        )}
      </div>
      
      {/* Расчет стоимости */}
      {settlement && (
        <PostExpressRateCalculator
          weight={weight}
          insuranceAmount={value}
          hasCOD={useCOD}
          onRateCalculated={setRate}
        />
      )}
      
      {/* Информация о доставке */}
      <div className="alert alert-info">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-current shrink-0 w-6 h-6">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <div>
          <div className="font-semibold">{t('deliveryInfo')}</div>
          <ul className="text-sm mt-2 space-y-1">
            <li>• {t('todayForTomorrow')} - 19:00</li>
            <li>• {t('smsNotification')}</li>
            <li>• {t('storageDays', { days: 5 })}</li>
            <li>• {t('freeReturn')}</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
```

### 7.4 Frontend компоненты для складского функционала

#### Компонент выбора самовывоза со склада Sve Tu

```typescript
// frontend/svetu/src/components/delivery/WarehousePickupSelector.tsx
import React, { useState } from 'react';
import { MapPin, Package, Clock, Check } from 'lucide-react';
import { useTranslations } from 'next-intl';

interface WarehousePickupSelectorProps {
  onSelect: (pickupPoint: string) => void;
  selectedOrderId?: string;
}

export const WarehousePickupSelector: React.FC<WarehousePickupSelectorProps> = ({ 
  onSelect, 
  selectedOrderId 
}) => {
  const t = useTranslations('delivery');
  const [selectedPoint, setSelectedPoint] = useState<string>('');

  const warehouseLocations = [
    {
      id: 'belgrade-main',
      city: 'Београд',
      address: 'Булевар Милоша Обреновића 112',
      hours: '09:00 - 20:00',
      available: true
    },
    {
      id: 'novi-sad',
      city: 'Нови Сад',
      address: 'Булевар Ослобођења 45',
      hours: '09:00 - 19:00',
      available: true
    },
    {
      id: 'nis',
      city: 'Ниш',
      address: 'Булевар Немањића 25',
      hours: '10:00 - 18:00',
      available: false
    }
  ];

  const handleSelectPickup = (locationId: string) => {
    setSelectedPoint(locationId);
    onSelect(locationId);
  };

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h3 className="card-title">
          <Package className="w-5 h-5" />
          {t('warehouse.pickupTitle')}
        </h3>
        
        <div className="divider"></div>

        <div className="space-y-4">
          {warehouseLocations.map((location) => (
            <div 
              key={location.id}
              className={`card bg-base-200 cursor-pointer transition-all hover:shadow-md ${
                selectedPoint === location.id ? 'ring-2 ring-primary' : ''
              } ${!location.available ? 'opacity-50 cursor-not-allowed' : ''}`}
              onClick={() => location.available && handleSelectPickup(location.id)}
            >
              <div className="card-body p-4">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <h4 className="font-semibold text-lg">{location.city}</h4>
                    <p className="text-sm opacity-75 flex items-center gap-1 mt-1">
                      <MapPin className="w-3 h-3" />
                      {location.address}
                    </p>
                    <p className="text-sm opacity-75 flex items-center gap-1 mt-1">
                      <Clock className="w-3 h-3" />
                      {location.hours}
                    </p>
                  </div>
                  <div className="flex items-center">
                    {location.available ? (
                      <div className={`btn btn-circle btn-sm ${
                        selectedPoint === location.id ? 'btn-primary' : 'btn-ghost'
                      }`}>
                        {selectedPoint === location.id && <Check className="w-4 h-4" />}
                      </div>
                    ) : (
                      <span className="badge badge-error">{t('warehouse.unavailable')}</span>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
```

#### Компонент отображения кода самовывоза

```typescript
// frontend/svetu/src/components/delivery/PickupCodeDisplay.tsx
import React, { useState } from 'react';
import { QrCode, Copy, Check, Package, MapPin, Clock } from 'lucide-react';
import { useTranslations } from 'next-intl';
import toast from 'react-hot-toast';

interface PickupCodeDisplayProps {
  code: string;
  qrCodeUrl?: string;
  expiresAt: string;
  warehouseAddress: string;
}

export const PickupCodeDisplay: React.FC<PickupCodeDisplayProps> = ({ 
  code, 
  qrCodeUrl, 
  expiresAt,
  warehouseAddress 
}) => {
  const t = useTranslations('delivery');
  const [copied, setCopied] = useState(false);

  const handleCopyCode = async () => {
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      toast.success(t('warehouse.codeCopied'));
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      toast.error(t('warehouse.copyFailed'));
    }
  };

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h3 className="card-title">
          <Package className="w-5 h-5" />
          {t('warehouse.pickupCodeTitle')}
        </h3>

        <div className="flex flex-col items-center py-6">
          {qrCodeUrl && (
            <div className="mb-6">
              <img 
                src={qrCodeUrl} 
                alt="QR Code" 
                className="w-48 h-48 rounded-lg shadow-md"
              />
            </div>
          )}

          <div className="text-center">
            <p className="text-sm opacity-75 mb-2">{t('warehouse.yourCode')}</p>
            <div className="flex items-center gap-2">
              <div className="kbd kbd-lg font-mono text-2xl tracking-wider">
                {code}
              </div>
              <button 
                className="btn btn-circle btn-sm btn-ghost"
                onClick={handleCopyCode}
              >
                {copied ? (
                  <Check className="w-4 h-4 text-success" />
                ) : (
                  <Copy className="w-4 h-4" />
                )}
              </button>
            </div>
          </div>
        </div>

        <div className="divider"></div>

        <div className="space-y-3">
          <div className="flex items-start gap-2">
            <MapPin className="w-4 h-4 mt-1 opacity-75" />
            <div>
              <p className="text-sm font-semibold">{t('warehouse.pickupAddress')}</p>
              <p className="text-sm opacity-75">{warehouseAddress}</p>
            </div>
          </div>

          <div className="flex items-start gap-2">
            <Clock className="w-4 h-4 mt-1 opacity-75" />
            <div>
              <p className="text-sm font-semibold">{t('warehouse.validUntil')}</p>
              <p className="text-sm opacity-75">{new Date(expiresAt).toLocaleDateString()}</p>
            </div>
          </div>
        </div>

        <div className="alert alert-warning mt-4">
          <span className="text-sm">{t('warehouse.pickupInstructions')}</span>
        </div>
      </div>
    </div>
  );
};
```

---

## 📚 Сценарии использования

### C2C: Продажа между физическими лицами

#### Сценарий 1: Оформление заказа с доставкой Post Express
1. **Покупатель** выбирает товар на маркетплейсе
2. **При оформлении:**
   - Выбирает Post Express как способ доставки
   - Вводит адрес доставки через PostExpressSettlementSelector
   - Видит расчет стоимости через PostExpressRateCalculator
   - Выбирает способ оплаты (предоплата или наложенный платеж)
3. **После оплаты:**
   - Создается запись в `post_express_shipments`
   - Продавец получает уведомление о продаже
4. **Продавец:**
   - Упаковывает товар
   - Подтверждает готовность к отправке
   - Система генерирует адресницу через EPK
5. **Post Express:**
   - Забирает посылку у продавца
   - Доставляет покупателю
   - При COD собирает оплату
6. **Завершение:**
   - Покупатель получает товар
   - Продавец получает оплату (при COD через Post Express)
   - Статус заказа обновляется на "delivered"

#### Сценарий 2: Отслеживание посылки
1. **Покупатель/Продавец** открывает страницу заказа
2. Видит компонент PostExpressTracker
3. Вводит номер отслеживания или он подгружается автоматически
4. Видит timeline со всеми статусами движения посылки
5. Получает SMS уведомления об изменении статуса

### Склад Sve Tu: Самовывоз и фулфилмент

#### Сценарий 1: Покупатель выбирает самовывоз со склада Sve Tu
1. **Покупатель** при оформлении заказа выбирает "Самовывоз со склада Sve Tu"
2. **Система** показывает компонент WarehousePickupSelector
3. **Покупатель** выбирает удобный склад (Белград, Нови Сад, Ниш)
4. **После оплаты:**
   - Создается запись в `warehouse_pickup_orders`
   - Генерируется 6-значный код самовывоза
   - Создается QR код для быстрого сканирования
   - Покупатель получает SMS/email с кодом
5. **При получении:**
   - Покупатель приходит на склад
   - Показывает код или QR на телефоне
   - Сотрудник склада проверяет документы
   - Товар выдается покупателю
   - Статус меняется на "completed"

#### Сценарий 2: Продавец отправляет товар на склад Sve Tu (FBS)
1. **Продавец** в панели управления выбирает "Отправить на склад Sve Tu"
2. **Система** показывает WarehouseInventoryManager
3. **Продавец** указывает:
   - Какие товары отправить
   - Количество единиц
   - Предполагаемую дату поставки
4. **Создается накладная:**
   - Генерируется уникальный номер поставки
   - Создаются записи в `warehouse_movements`
   - Продавец получает инструкции по упаковке
5. **При приемке на складе:**
   - Сотрудники проверяют товар
   - Обновляется `warehouse_inventory`
   - Товар становится доступен для продажи
   - Продавец получает уведомление

#### Сценарий 3: Консолидация заказов на складе
1. **Покупатель** делает несколько заказов у разных продавцов
2. **Все товары** находятся на складе Sve Tu
3. **Система** предлагает объединить заказы в одну отправку
4. **При согласии:**
   - Создается консолидированная отправка
   - Один трек-номер для всех товаров
   - Экономия на доставке
5. **Отправка** через Post Express со склада

#### Сценарий 4: Управление складскими остатками
1. **Продавец** открывает раздел "Склад"
2. **Видит дашборд:**
   - Текущие остатки по товарам
   - Зарезервированные товары
   - История движений
   - Оборачиваемость товаров
3. **Может выполнить действия:**
   - Пополнить запасы
   - Вывести товар со склада
   - Изменить цены
   - Настроить автоматическое пополнение

### B2C: Продажа из витрины

#### Сценарий 1: Настройка Post Express для витрины
1. **Владелец витрины** заходит в настройки доставки
2. Активирует Post Express
3. Вводит данные отправителя (адрес склада)
4. Настраивает тарифы и наценки
5. Сохраняет настройки в `storefront_post_express_settings`

#### Сценарий 2: Массовая обработка заказов
1. **Менеджер витрины** открывает панель заказов
2. Выбирает несколько оплаченных заказов
3. Нажимает "Создать манифест Post Express"
4. Система:
   - Создает записи в `post_express_shipments` для каждого заказа
   - Генерирует манифест через API
   - Возвращает адресницы для печати
5. **Менеджер:**
   - Печатает адресницы
   - Клеит на упакованные товары
   - Передает курьеру Post Express

#### Сценарий 3: Аналитика доставок
1. **Владелец витрины** открывает раздел аналитики
2. Видит дашборд с метриками:
   - Количество отправлений за период
   - Средний срок доставки
   - Процент успешных доставок
   - Расходы на доставку
   - География доставок
3. Может экспортировать данные для бухгалтерии

---

## 📊 Метрики и мониторинг

### Ключевые метрики

#### Операционные метрики
- **Количество отправлений** в день/неделю/месяц
- **Средний вес посылки**
- **Средняя стоимость доставки**
- **Процент COD платежей**
- **География доставок** (топ городов)

#### Метрики качества
- **Процент успешных доставок**
- **Средний срок доставки**
- **Количество возвратов**
- **Процент поврежденных посылок**

#### Финансовые метрики
- **Общие расходы на доставку**
- **Средняя маржа на доставке**
- **Комиссии за COD**
- **ROI интеграции**

### Мониторинг API

#### Дашборд для мониторинга
```sql
-- Статистика API вызовов за последние 24 часа
SELECT 
    transaction_type,
    COUNT(*) as total_calls,
    AVG(duration_ms) as avg_duration,
    MAX(duration_ms) as max_duration,
    SUM(CASE WHEN response_status != 0 THEN 1 ELSE 0 END) as errors
FROM post_express_api_log
WHERE created_at > NOW() - INTERVAL '24 hours'
GROUP BY transaction_type;

-- Топ ошибок
SELECT 
    error_message,
    COUNT(*) as error_count
FROM post_express_api_log
WHERE response_status != 0
    AND created_at > NOW() - INTERVAL '24 hours'
GROUP BY error_message
ORDER BY error_count DESC
LIMIT 10;
```

### Алерты и уведомления

#### Критические алерты
- API недоступен более 5 минут
- Процент ошибок > 10%
- Среднее время ответа > 5 секунд

#### Предупреждения
- Процент ошибок > 5%
- Среднее время ответа > 2 секунд
- Более 100 неотслеживаемых посылок

---

## ⚠️ Риски и митигация

### Технические риски

| Риск | Вероятность | Влияние | Митигация |
|------|-------------|---------|-----------|
| Недоступность WSP API | Средняя | Высокое | Retry механизм, очередь отложенных запросов, fallback на ручное оформление |
| Изменение API | Низкая | Среднее | Версионирование API клиента, мониторинг изменений |
| Превышение лимитов API | Средняя | Среднее | Rate limiting, батчинг запросов, кеширование |
| Потеря данных отслеживания | Низкая | Высокое | Регулярный бекап, синхронизация статусов |

### Операционные риски

| Риск | Вероятность | Влияние | Митигация |
|------|-------------|---------|-----------|
| Ошибки в адресах | Высокая | Среднее | Валидация через API, автокомплит, проверка оператором |
| Задержки в доставке | Средняя | Среднее | SLA мониторинг, альтернативные провайдеры |
| Мошенничество с COD | Средняя | Высокое | Верификация покупателей, лимиты, blacklist |
| Повреждение товаров | Низкая | Высокое | Страхование, правила упаковки, фото-фиксация |

### План действий при сбоях

#### При недоступности API
1. Активация retry механизма (3 попытки с экспоненциальной задержкой)
2. Постановка запросов в очередь для повторной отправки
3. Уведомление администраторов
4. Переключение на ручное оформление через Web Express

#### При массовых ошибках
1. Автоматическая остановка синхронизации
2. Алерт команде разработки
3. Анализ логов для выявления причины
4. Rollback на предыдущую версию при необходимости

---

## 📞 Контакты и ресурсы

### Post Express контакты
- **Коммерческий отдел:** 011/3718-221, 011/3718-202, 011/3718-263
- **Email:** prodaja@posta.rs
- **Техподдержка API:** support@posta.rs
- **Сайт:** https://www.postexpress.rs

### Полезные ссылки
- **Список населенных пунктов:** http://www.postexpress.rs/struktura/lat/usluge/urucenje-danas-za-sutra-a-z.asp
- **Терминский план:** http://www.posta.rs/dokumenta/lat/novcano/Terminski-plan-Poste.pdf
- **Отслеживание:** https://www.posta.rs/cir/alati/pracenje-posiljke.aspx
- **EPK (Электронска пријемна књига):** https://epk.posta.rs
- **Web Express:** https://webexpress.posta.rs

### Документация
- WSP Web API - Address Information
- WSP Web API - Exchange of Data
- WSP Web API - Shipment Tracking
- Руководство по интеграции EPK
- Руководство по Web Express

### Команда проекта
- **Product Manager:** Определить ответственного
- **Tech Lead:** Определить ответственного
- **Backend Developer:** Определить ответственного
- **Frontend Developer:** Определить ответственного
- **QA Engineer:** Определить ответственного

---

## 📅 План внедрения

### Фаза 1: Подготовка (1 неделя)
- [ ] Подписание договора с Post Express
- [ ] Получение тестовых учетных данных API
- [ ] Настройка тестового окружения
- [ ] Создание проектной команды

### Фаза 2: Backend разработка (2 недели)
- [ ] Реализация WSP API клиента
- [ ] Создание миграций БД
- [ ] Реализация сервисов и репозиториев
- [ ] Создание API endpoints
- [ ] Написание unit тестов

### Фаза 3: Frontend разработка (2 недели)
- [ ] Создание компонентов выбора населенных пунктов
- [ ] Реализация калькулятора стоимости
- [ ] Компонент отслеживания посылок
- [ ] Интеграция в checkout процесс
- [ ] Панель управления для витрин

### Фаза 4: Интеграционное тестирование (1 неделя)
- [ ] Тестирование C2C сценариев
- [ ] Тестирование B2C сценариев
- [ ] Тестирование COD платежей
- [ ] Нагрузочное тестирование API

### Фаза 5: Пилотный запуск (2 недели)
- [ ] Запуск с ограниченной группой пользователей
- [ ] Мониторинг метрик
- [ ] Сбор обратной связи
- [ ] Исправление выявленных проблем

### Фаза 6: Полный запуск
- [ ] Развертывание на production
- [ ] Обучение службы поддержки
- [ ] Публикация документации для пользователей
- [ ] Маркетинговая кампания

---

## ✅ Критерии успеха

### Технические критерии
- ✅ Успешная интеграция WSP API
- ✅ Доступность сервиса > 99.9%
- ✅ Среднее время ответа API < 1 сек
- ✅ Покрытие кода тестами > 80%

### Бизнес критерии
- ✅ Увеличение конверсии в покупку на 15%
- ✅ Снижение отказов из-за доставки на 20%
- ✅ Рост географии продаж на 30%
- ✅ Удовлетворенность доставкой > 4.5/5

### Операционные критерии
- ✅ Автоматизация 90% процессов доставки
- ✅ Среднее время доставки < 2 дней
- ✅ Процент успешных доставок > 95%
- ✅ Время обработки возвратов < 3 дней

---

## 🌐 Локализация для складского функционала

### Модуль переводов delivery.json

#### Русский (ru):
```json
{
  "warehouse": {
    "pickupTitle": "Самовывоз со склада Sve Tu",
    "unavailable": "Недоступно",
    "pickupSelected": "Выбран пункт самовывоза",
    "pickupCodeTitle": "Код для получения заказа",
    "yourCode": "Ваш код:",
    "codeCopied": "Код скопирован",
    "copyFailed": "Не удалось скопировать",
    "pickupAddress": "Адрес получения:",
    "validUntil": "Действителен до:",
    "pickupInstructions": "Покажите этот код или QR на складе. Не забудьте взять документ, удостоверяющий личность.",
    "inventoryTitle": "Управление складскими остатками",
    "totalItems": "Всего товаров",
    "available": "Доступно",
    "reserved": "Зарезервировано",
    "product": "Товар",
    "sku": "Артикул",
    "location": "Расположение",
    "actions": "Действия",
    "sendMore": "Отправить еще",
    "bulkSend": "Отправить выбранные ({count})",
    "loadError": "Ошибка загрузки данных",
    "sentSuccess": "Товар успешно отправлен на склад",
    "sendError": "Ошибка отправки на склад",
    "updateSuccess": "Количество обновлено",
    "updateError": "Ошибка обновления"
  }
}
```

#### Английский (en):
```json
{
  "warehouse": {
    "pickupTitle": "Pickup from Sve Tu Warehouse",
    "unavailable": "Unavailable",
    "pickupSelected": "Pickup point selected",
    "pickupCodeTitle": "Order Pickup Code",
    "yourCode": "Your code:",
    "codeCopied": "Code copied",
    "copyFailed": "Failed to copy",
    "pickupAddress": "Pickup address:",
    "validUntil": "Valid until:",
    "pickupInstructions": "Show this code or QR at the warehouse. Don't forget to bring your ID.",
    "inventoryTitle": "Warehouse Inventory Management",
    "totalItems": "Total Items",
    "available": "Available",
    "reserved": "Reserved",
    "product": "Product",
    "sku": "SKU",
    "location": "Location",
    "actions": "Actions",
    "sendMore": "Send More",
    "bulkSend": "Send Selected ({count})",
    "loadError": "Error loading data",
    "sentSuccess": "Product successfully sent to warehouse",
    "sendError": "Error sending to warehouse",
    "updateSuccess": "Quantity updated",
    "updateError": "Update error"
  }
}
```

#### Сербский (sr):
```json
{
  "warehouse": {
    "pickupTitle": "Preuzimanje sa Sve Tu skladišta",
    "unavailable": "Nedostupno",
    "pickupSelected": "Izabrano mesto preuzimanja",
    "pickupCodeTitle": "Kod za preuzimanje porudžbine",
    "yourCode": "Vaš kod:",
    "codeCopied": "Kod je kopiran",
    "copyFailed": "Kopiranje neuspešno",
    "pickupAddress": "Adresa preuzimanja:",
    "validUntil": "Važi do:",
    "pickupInstructions": "Pokažite ovaj kod ili QR na skladištu. Ne zaboravite da ponesete ličnu kartu.",
    "inventoryTitle": "Upravljanje skladišnim zalihama",
    "totalItems": "Ukupno proizvoda",
    "available": "Dostupno",
    "reserved": "Rezervisano",
    "product": "Proizvod",
    "sku": "Šifra",
    "location": "Lokacija",
    "actions": "Akcije",
    "sendMore": "Pošalji još",
    "bulkSend": "Pošalji izabrane ({count})",
    "loadError": "Greška pri učitavanju podataka",
    "sentSuccess": "Proizvod uspešno poslat na skladište",
    "sendError": "Greška slanja na skladište",
    "updateSuccess": "Količina ažurirana",
    "updateError": "Greška ažuriranja"
  }
}
```

---

## 📦 Ключевые изменения в версии 2.0

### Изменения после аудита:

#### 🗄️ База данных:
- ✅ Использование существующей таблицы `storefront_delivery_options` с `provider='post-express'`
- ✅ Хранение конфигурации в JSONB поле `provider_config`
- ✅ Использование BIGINT для storefront_order_id
- ✅ Использование TIMESTAMP WITH TIME ZONE

#### 💻 Backend:
- ✅ Модуль `/internal/proj/postexpress/` по паттерну проекта
- ✅ Структура: handler/, service/, storage/postgres/, module.go
- ✅ API пути `/api/v1/postexpress/`
- ✅ Соответствие Swagger аннотациям

#### 🎨 Frontend:
- ✅ Использование lucide-react вместо @heroicons/react
- ✅ DaisyUI классы: btn, input, card, alert, checkbox
- ✅ Redux Toolkit для state management
- ✅ Модульные переводы next-intl

#### 🌐 Локализация:
- ✅ Поддержка 3 локалей: en, ru, sr
- ✅ Модульные файлы переводов delivery.json
- ✅ Использование useTranslations('delivery')
- ✅ Добавлены переводы для складского функционала

### Добавленный функционал склада Sve Tu:

#### 🏢 Складская система:
- ✅ Самовывоз со складов в Белграде, Нови Саде, Нише
- ✅ Fulfillment by Sve Tu (FBS) для продавцов
- ✅ Управление складскими остатками
- ✅ Консолидация заказов от разных продавцов
- ✅ QR коды для быстрого получения заказов

#### 📊 База данных складского учета:
- ✅ Таблица `warehouse_inventory` для учета товаров
- ✅ Таблица `warehouse_movements` для движений товара
- ✅ Таблица `warehouse_pickup_orders` для самовывоза
- ✅ Таблица `warehouse_consolidations` для объединения заказов

#### 💻 Backend сервисы склада:
- ✅ WarehouseService для управления складом
- ✅ Генерация 6-значных кодов самовывоза
- ✅ API для работы с остатками
- ✅ Поддержка резервирования товаров

#### 🎨 Frontend компоненты склада:
- ✅ WarehousePickupSelector для выбора пункта выдачи
- ✅ PickupCodeDisplay для отображения кода получения
- ✅ WarehouseInventoryManager для управления остатками
- ✅ Интеграция с DaisyUI и lucide-react

### Преимущества новой архитектуры:

1. **Минимальные изменения БД** - использование существующих таблиц
2. **Соответствие паттерну проекта** - единообразная структура модулей
3. **Единый UI/UX** - использование DaisyUI компонентов
4. **Полная локализация** - поддержка всех языков платформы
5. **Масштабируемость** - готовность к будущим расширениям
6. **Полный цикл фулфилмента** - от приемки до выдачи товара
7. **Оптимизация логистики** - консолидация и самовывоз

---

*Документ подготовлен на основе анализа предоставленной документации Post Express и требований платформы Sve Tu.*

**✅ Версия: 2.0 - Адаптирована под реальную архитектуру проекта**
*Дата обновления: 15.08.2025*
*Статус: Готово к реализации*