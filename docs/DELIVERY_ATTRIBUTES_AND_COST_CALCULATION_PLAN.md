# 📦 ПЛАН ИНТЕГРАЦИИ АТРИБУТОВ ДОСТАВКИ И АВТОМАТИЧЕСКОГО РАСЧЕТА СТОИМОСТИ

## 📋 Обзор

План интеграции системы атрибутов доставки (габариты, вес, хрупкость) в маркетплейс и реализация автоматического расчета стоимости доставки на основе характеристик товара и маршрута доставки.

---

## 🎯 Цели

1. Добавить атрибуты доставки ко всем товарам (B2C витрины и C2C объявления)
2. Реализовать автоматический расчет стоимости доставки
3. Интегрировать расчеты с существующей системой логистики
4. Обеспечить прозрачность стоимости доставки для покупателей

---

## ⚙️ ФАЗА 1: ПОДГОТОВКА АТРИБУТОВ (2-3 дня)

### 1.1 Создание универсальных атрибутов доставки

#### Миграция базы данных:
```sql
-- 000019_delivery_attributes.up.sql

-- Добавление атрибутов доставки для всех товаров
ALTER TABLE marketplace_listings ADD COLUMN IF NOT EXISTS delivery_attributes JSONB DEFAULT '{}';
ALTER TABLE storefront_products ADD COLUMN IF NOT EXISTS delivery_attributes JSONB DEFAULT '{}';

-- Структура delivery_attributes:
-- {
--   "weight_kg": 0.5,
--   "dimensions": {
--     "length_cm": 30,
--     "width_cm": 20,
--     "height_cm": 10
--   },
--   "volume_m3": 0.006,
--   "is_fragile": false,
--   "requires_special_handling": false,
--   "stackable": true,
--   "max_stack_weight_kg": 50,
--   "packaging_type": "box", // box, envelope, pallet, custom
--   "hazmat_class": null
-- }

-- Создание таблицы предустановленных категорий с атрибутами
CREATE TABLE delivery_category_defaults (
    id SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES marketplace_categories(id),
    default_weight_kg DECIMAL(10,3),
    default_length_cm DECIMAL(10,2),
    default_width_cm DECIMAL(10,2),
    default_height_cm DECIMAL(10,2),
    default_packaging_type VARCHAR(50),
    is_typically_fragile BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_delivery_category_defaults_category ON delivery_category_defaults(category_id);
```

### 1.2 Создание системы расчета стоимости

```sql
-- Таблица формул расчета доставки
CREATE TABLE delivery_pricing_rules (
    id SERIAL PRIMARY KEY,
    provider_id INTEGER REFERENCES delivery_providers(id),
    rule_type VARCHAR(50) NOT NULL, -- 'weight_based', 'volume_based', 'zone_based', 'combined'

    -- Весовые правила
    weight_ranges JSONB, -- [{"from": 0, "to": 1, "price_per_kg": 5}, ...]

    -- Объемные правила
    volume_ranges JSONB, -- [{"from": 0, "to": 0.01, "price_per_m3": 100}, ...]

    -- Зональные правила
    zone_multipliers JSONB, -- {"local": 1.0, "regional": 1.5, "national": 2.0, "international": 3.5}

    -- Дополнительные сборы
    fragile_surcharge DECIMAL(10,2) DEFAULT 0,
    oversized_surcharge DECIMAL(10,2) DEFAULT 0, -- если любая сторона > 100cm
    special_handling_surcharge DECIMAL(10,2) DEFAULT 0,

    -- Минимальная и максимальная стоимость
    min_price DECIMAL(10,2),
    max_price DECIMAL(10,2),

    -- Формула расчета (для сложных случаев)
    custom_formula TEXT, -- PostgreSQL функция или выражение

    priority INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Таблица зон доставки
CREATE TABLE delivery_zones (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'local', 'regional', 'national', 'international'

    -- Географические границы
    countries TEXT[], -- массив кодов стран
    regions TEXT[], -- массив регионов
    cities TEXT[], -- массив городов
    postal_codes TEXT[], -- массив почтовых индексов

    -- Полигон для точного определения (GIS)
    boundary GEOMETRY(POLYGON, 4326),

    -- Расстояние от центра (для радиусных зон)
    center_point GEOMETRY(POINT, 4326),
    radius_km DECIMAL(10,2),

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_delivery_zones_boundary ON delivery_zones USING GIST(boundary);
CREATE INDEX idx_delivery_zones_center ON delivery_zones USING GIST(center_point);
```

---

## 🔧 ФАЗА 2: BACKEND РЕАЛИЗАЦИЯ (3-4 дня)

### 2.1 Сервис расчета стоимости доставки

```go
// backend/internal/proj/delivery/calculator/service.go

type DeliveryCalculator struct {
    db              *sql.DB
    providerManager *ProviderManager
    zoneService     *ZoneService
}

type CalculationRequest struct {
    // Откуда и куда
    FromLocation    Location `json:"from_location"`
    ToLocation      Location `json:"to_location"`

    // Характеристики товара
    Weight          float64  `json:"weight_kg"`
    Dimensions      Dims     `json:"dimensions"`
    IsFragile       bool     `json:"is_fragile"`
    SpecialHandling bool     `json:"special_handling"`

    // Опции
    ProviderId      *int     `json:"provider_id,omitempty"`
    InsuranceValue  float64  `json:"insurance_value,omitempty"`
}

type CalculationResponse struct {
    Providers []ProviderQuote `json:"providers"`
    Cheapest  *ProviderQuote  `json:"cheapest"`
    Fastest   *ProviderQuote  `json:"fastest"`
}

type ProviderQuote struct {
    ProviderId      int      `json:"provider_id"`
    ProviderName    string   `json:"provider_name"`
    BasePrice       float64  `json:"base_price"`
    Surcharges      []Charge `json:"surcharges"`
    TotalPrice      float64  `json:"total_price"`
    EstimatedDays   [2]int   `json:"estimated_days"` // [min, max]
    Restrictions    []string `json:"restrictions,omitempty"`
}
```

### 2.2 API эндпоинты

```go
// Расчет стоимости для конкретного товара
POST /api/v1/delivery/calculate-product
{
    "product_type": "listing|storefront_product",
    "product_id": 123,
    "to_address": {...}
}

// Расчет стоимости с произвольными параметрами
POST /api/v1/delivery/calculate-custom
{
    "from_location": {...},
    "to_location": {...},
    "items": [
        {
            "weight_kg": 0.5,
            "dimensions": {"length_cm": 30, "width_cm": 20, "height_cm": 10},
            "is_fragile": false,
            "quantity": 2
        }
    ]
}

// Получение атрибутов доставки для товара
GET /api/v1/products/{id}/delivery-attributes

// Обновление атрибутов доставки
PUT /api/v1/products/{id}/delivery-attributes
{
    "weight_kg": 0.5,
    "dimensions": {...},
    "is_fragile": true
}
```

---

## 🎨 ФАЗА 3: FRONTEND ИНТЕГРАЦИЯ (3-4 дня)

### 3.1 Форма добавления атрибутов доставки

#### Компонент для C2C объявлений:
```tsx
// frontend/svetu/src/components/marketplace/DeliveryAttributesForm.tsx

interface DeliveryAttributesFormProps {
    categoryId: number;
    onAttributesChange: (attrs: DeliveryAttributes) => void;
}

export function DeliveryAttributesForm({ categoryId, onAttributesChange }: DeliveryAttributesFormProps) {
    const [attributes, setAttributes] = useState<DeliveryAttributes>({
        weight_kg: 0,
        dimensions: { length_cm: 0, width_cm: 0, height_cm: 0 },
        is_fragile: false,
        packaging_type: 'box'
    });

    // Загрузка дефолтных значений для категории
    useEffect(() => {
        fetchCategoryDefaults(categoryId).then(setAttributes);
    }, [categoryId]);

    return (
        <div className="space-y-4">
            <h3>Параметры доставки</h3>

            {/* Вес */}
            <div>
                <label>Вес товара (кг)</label>
                <input
                    type="number"
                    step="0.1"
                    value={attributes.weight_kg}
                    onChange={(e) => updateAttribute('weight_kg', e.target.value)}
                />
            </div>

            {/* Габариты */}
            <div className="grid grid-cols-3 gap-2">
                <input placeholder="Длина (см)" type="number" />
                <input placeholder="Ширина (см)" type="number" />
                <input placeholder="Высота (см)" type="number" />
            </div>

            {/* Особые условия */}
            <div className="space-y-2">
                <label className="flex items-center">
                    <input type="checkbox" /> Хрупкий товар
                </label>
                <label className="flex items-center">
                    <input type="checkbox" /> Требует специальной обработки
                </label>
            </div>

            {/* Калькулятор стоимости */}
            <DeliveryCostPreview attributes={attributes} />
        </div>
    );
}
```

### 3.2 Компонент расчета стоимости в корзине

```tsx
// frontend/svetu/src/components/cart/DeliveryCalculator.tsx

export function DeliveryCalculator({ items, deliveryAddress }) {
    const [quotes, setQuotes] = useState<ProviderQuote[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (deliveryAddress) {
            calculateDelivery();
        }
    }, [items, deliveryAddress]);

    const calculateDelivery = async () => {
        setLoading(true);
        const response = await api.calculateDelivery({
            items: items.map(item => ({
                product_id: item.id,
                product_type: item.type,
                quantity: item.quantity
            })),
            to_address: deliveryAddress
        });
        setQuotes(response.providers);
        setLoading(false);
    };

    return (
        <div className="delivery-options">
            <h3>Варианты доставки</h3>
            {loading ? (
                <Skeleton />
            ) : (
                quotes.map(quote => (
                    <DeliveryOption
                        key={quote.provider_id}
                        quote={quote}
                        onSelect={() => selectDelivery(quote)}
                    />
                ))
            )}
        </div>
    );
}
```

---

## 📊 ФАЗА 4: ИНТЕГРАЦИЯ С ПРОВАЙДЕРАМИ (2-3 дня)

### 4.1 Адаптеры для расчета через API провайдеров

```go
// backend/internal/proj/delivery/providers/postexpress/calculator.go

func (p *PostExpressProvider) CalculateRate(req CalculationRequest) (*ProviderQuote, error) {
    // Преобразование в формат Post Express
    peRequest := &postexpress.RateRequest{
        FromZip: req.FromLocation.PostalCode,
        ToZip:   req.ToLocation.PostalCode,
        Weight:  req.Weight,
        Length:  req.Dimensions.Length,
        Width:   req.Dimensions.Width,
        Height:  req.Dimensions.Height,
    }

    // Вызов API (или mock в dev режиме)
    rate, err := p.client.GetRate(peRequest)
    if err != nil {
        return nil, err
    }

    // Добавление надбавок
    totalPrice := rate.BasePrice
    surcharges := []Charge{}

    if req.IsFragile {
        fragileCharge := rate.BasePrice * 0.15 // 15% за хрупкость
        surcharges = append(surcharges, Charge{
            Type: "fragile",
            Amount: fragileCharge,
        })
        totalPrice += fragileCharge
    }

    return &ProviderQuote{
        ProviderId:    p.ID,
        ProviderName:  "Post Express",
        BasePrice:     rate.BasePrice,
        Surcharges:    surcharges,
        TotalPrice:    totalPrice,
        EstimatedDays: [2]int{rate.MinDays, rate.MaxDays},
    }, nil
}
```

---

## 🧪 ФАЗА 5: ТЕСТИРОВАНИЕ И ОПТИМИЗАЦИЯ (2 дня)

### 5.1 Создание тестовых данных

```sql
-- Заполнение дефолтных атрибутов для популярных категорий
INSERT INTO delivery_category_defaults (category_id, default_weight_kg, default_length_cm, default_width_cm, default_height_cm, default_packaging_type, is_typically_fragile) VALUES
-- Электроника
(1, 0.5, 20, 15, 5, 'box', true),
-- Одежда
(2, 0.3, 30, 25, 5, 'envelope', false),
-- Мебель
(3, 15.0, 120, 60, 80, 'custom', false),
-- Книги
(4, 0.4, 20, 15, 3, 'envelope', false);
```

### 5.2 Unit тесты

```go
func TestDeliveryCalculator(t *testing.T) {
    calc := NewDeliveryCalculator(db)

    t.Run("Calculate for small package", func(t *testing.T) {
        req := CalculationRequest{
            Weight: 0.5,
            Dimensions: Dims{30, 20, 10},
            FromLocation: Location{City: "Belgrade"},
            ToLocation: Location{City: "Novi Sad"},
        }

        resp, err := calc.Calculate(req)
        assert.NoError(t, err)
        assert.NotEmpty(t, resp.Providers)
        assert.NotNil(t, resp.Cheapest)
    })
}
```

---

## 🔄 ФАЗА 6: МИГРАЦИЯ СУЩЕСТВУЮЩИХ ДАННЫХ (1-2 дня)

### 6.1 Скрипт миграции для существующих товаров

```sql
-- 000020_populate_delivery_attributes.up.sql

-- Заполнение атрибутов на основе категорий
UPDATE marketplace_listings ml
SET delivery_attributes = jsonb_build_object(
    'weight_kg', COALESCE(dcd.default_weight_kg, 1.0),
    'dimensions', jsonb_build_object(
        'length_cm', COALESCE(dcd.default_length_cm, 30),
        'width_cm', COALESCE(dcd.default_width_cm, 20),
        'height_cm', COALESCE(dcd.default_height_cm, 10)
    ),
    'is_fragile', COALESCE(dcd.is_typically_fragile, false),
    'packaging_type', COALESCE(dcd.default_packaging_type, 'box')
)
FROM delivery_category_defaults dcd
WHERE ml.category_id = dcd.category_id
  AND ml.delivery_attributes = '{}';

-- Аналогично для storefront_products
UPDATE storefront_products sp
SET delivery_attributes = jsonb_build_object(
    'weight_kg', COALESCE(dcd.default_weight_kg, 1.0),
    'dimensions', jsonb_build_object(
        'length_cm', COALESCE(dcd.default_length_cm, 30),
        'width_cm', COALESCE(dcd.default_width_cm, 20),
        'height_cm', COALESCE(dcd.default_height_cm, 10)
    ),
    'is_fragile', COALESCE(dcd.is_typically_fragile, false),
    'packaging_type', COALESCE(dcd.default_packaging_type, 'box')
)
FROM delivery_category_defaults dcd
WHERE sp.category_id = dcd.category_id
  AND sp.delivery_attributes = '{}';
```

---

## 📈 МЕТРИКИ УСПЕХА

- **Точность расчета**: отклонение от реальной стоимости < 10%
- **Скорость расчета**: < 500ms для корзины из 5 товаров
- **Заполненность атрибутов**: > 90% товаров имеют атрибуты доставки
- **Конверсия**: увеличение конверсии корзина → заказ на 10-15%
- **Поддержка**: снижение запросов о стоимости доставки на 50%

---

## 🚀 ДОПОЛНИТЕЛЬНЫЕ ВОЗМОЖНОСТИ

### Фаза 2 (после MVP):

1. **Умные подсказки**:
   - AI определение габаритов по фото
   - Автозаполнение веса на основе похожих товаров

2. **Оптимизация упаковки**:
   - Группировка товаров в посылки
   - Рекомендации по упаковке

3. **Динамическое ценообразование**:
   - Скидки на доставку при покупке нескольких товаров
   - Промо-кампании провайдеров

4. **Интеграция с весами и сканерами**:
   - API для подключения торговых весов
   - Сканирование штрих-кодов для автозаполнения

---

## ⚠️ ВАЖНЫЕ ЗАМЕЧАНИЯ

1. **Приватность**: не показывать точный адрес продавца при расчете
2. **Кеширование**: кешировать расчеты на 15 минут для одинаковых параметров
3. **Fallback**: если атрибуты не заполнены, использовать дефолтные для категории
4. **Валидация**: проверять реалистичность введенных габаритов
5. **Обратная связь**: собирать данные о реальной стоимости для улучшения алгоритмов

---

## 🗓️ TIMELINE

| Фаза | Длительность | Приоритет |
|------|--------------|-----------|
| Подготовка атрибутов | 2-3 дня | Критический |
| Backend реализация | 3-4 дня | Критический |
| Frontend интеграция | 3-4 дня | Высокий |
| Интеграция с провайдерами | 2-3 дня | Средний |
| Тестирование | 2 дня | Высокий |
| Миграция данных | 1-2 дня | Средний |

**Общее время**: 13-18 дней

---

*Документ подготовлен: 2025-01-20*
*Версия: 1.0*