# 📦 GIS модуль: Интеграция с сервисами доставки

**Версия**: 1.0  
**Дата**: 2025-01-10  
**Сервисы доставки**: D Express, Пошта Србије

## 🚚 Обзор интеграции доставки

### Поддерживаемые сервисы:

1. **D Express**
   - Экспресс-доставка по всей Сербии
   - API для отслеживания посылок
   - Расчет стоимости по зонам
   - Срок доставки: 1-2 дня

2. **Пошта Србије (Почта Сербии)**
   - Национальный почтовый оператор
   - Широкая сеть отделений
   - Бюджетный вариант доставки
   - Срок доставки: 2-5 дней

## 📋 Архитектура интеграции

### 1. Модель данных для служб доставки

```sql
-- Таблица служб доставки
CREATE TABLE delivery_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL, -- 'dexpress', 'posta_srbije'
    name VARCHAR(100) NOT NULL,
    name_cyrillic VARCHAR(100),
    api_endpoint VARCHAR(255),
    api_key_encrypted TEXT,
    is_active BOOLEAN DEFAULT true,
    capabilities JSONB DEFAULT '{}', -- tracking, cod, insurance, etc
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Зоны доставки для каждого провайдера
CREATE TABLE delivery_provider_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID REFERENCES delivery_providers(id),
    zone_code VARCHAR(50), -- внутренний код зоны провайдера
    zone_name VARCHAR(100),
    municipalities TEXT[], -- список општина
    postal_codes TEXT[], -- список почтовых индексов
    base_price DECIMAL(10,2),
    price_per_kg DECIMAL(10,2),
    estimated_days_min INT,
    estimated_days_max INT,
    metadata JSONB DEFAULT '{}'
);

-- Пункты выдачи (почтовые отделения, пункты D Express)
CREATE TABLE pickup_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID REFERENCES delivery_providers(id),
    external_id VARCHAR(100), -- ID в системе провайдера
    name VARCHAR(200),
    address VARCHAR(500),
    location GEOGRAPHY(POINT, 4326),
    city VARCHAR(100),
    postal_code VARCHAR(10),
    working_hours JSONB, -- {"mon": "08:00-17:00", ...}
    services JSONB, -- ["pickup", "drop_off", "cod"]
    is_active BOOLEAN DEFAULT true,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Отслеживание посылок
CREATE TABLE delivery_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    provider_id UUID REFERENCES delivery_providers(id),
    tracking_number VARCHAR(100) UNIQUE,
    status VARCHAR(50), -- 'pending', 'picked_up', 'in_transit', 'delivered'
    status_details JSONB,
    last_location VARCHAR(200),
    last_update TIMESTAMPTZ,
    estimated_delivery DATE,
    actual_delivery TIMESTAMPTZ,
    tracking_history JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 2. Интеграция с API служб доставки

#### D Express интеграция:
```go
// internal/proj/delivery/providers/dexpress.go
package providers

type DExpressClient struct {
    apiURL    string
    apiKey    string
    client    *http.Client
}

// Расчет стоимости доставки
func (d *DExpressClient) CalculatePrice(params DeliveryParams) (*DeliveryQuote, error) {
    request := DExpressCalculateRequest{
        FromPostalCode: params.FromPostalCode,
        ToPostalCode:   params.ToPostalCode,
        Weight:         params.WeightKg,
        CODAmount:      params.CODAmount,
    }
    
    // Вызов API D Express
    resp, err := d.client.Post(
        d.apiURL + "/calculate-price",
        "application/json",
        json.Marshal(request),
    )
    
    return &DeliveryQuote{
        Provider:      "D Express",
        Price:         resp.Price,
        EstimatedDays: resp.EstimatedDays,
        Services:      resp.AvailableServices,
    }, nil
}

// Создание заказа на доставку
func (d *DExpressClient) CreateShipment(order DeliveryOrder) (*Shipment, error) {
    // Маппинг данных для D Express API
    shipment := d.mapToExpressFormat(order)
    
    resp, err := d.client.Post(
        d.apiURL + "/create-shipment",
        "application/json",
        json.Marshal(shipment),
    )
    
    return &Shipment{
        TrackingNumber: resp.TrackingNumber,
        Label:          resp.LabelPDF,
        PickupDate:     resp.PickupDate,
    }, nil
}

// Отслеживание посылки
func (d *DExpressClient) TrackShipment(trackingNumber string) (*TrackingInfo, error) {
    resp, err := d.client.Get(
        fmt.Sprintf("%s/track/%s", d.apiURL, trackingNumber),
    )
    
    return &TrackingInfo{
        Status:        d.mapStatus(resp.Status),
        Location:      resp.CurrentLocation,
        LastUpdate:    resp.LastUpdate,
        History:       d.mapHistory(resp.Events),
        EstimatedDate: resp.EstimatedDelivery,
    }, nil
}
```

#### Пошта Србије интеграция:
```go
// internal/proj/delivery/providers/posta_srbije.go
package providers

type PostaSrbijeClient struct {
    apiURL   string
    username string
    password string
    client   *http.Client
}

// Расчет стоимости через таблицу тарифов
func (p *PostaSrbijeClient) CalculatePrice(params DeliveryParams) (*DeliveryQuote, error) {
    // Пошта Србије использует зоны на основе расстояния
    zone := p.determineZone(params.FromPostalCode, params.ToPostalCode)
    
    // Базовая цена + цена за вес
    basePrice := p.getZoneBasePrice(zone)
    weightPrice := p.calculateWeightPrice(params.WeightKg, zone)
    
    return &DeliveryQuote{
        Provider:      "Пошта Србије",
        Price:         basePrice + weightPrice,
        EstimatedDays: p.getZoneEstimatedDays(zone),
        Services:      []string{"standard", "registered", "cod"},
    }, nil
}

// Поиск ближайшего почтового отделения
func (p *PostaSrbijeClient) FindNearestPostOffice(lat, lng float64) (*PickupPoint, error) {
    // Запрос к базе данных почтовых отделений
    query := `
        SELECT 
            id, name, address, 
            ST_Distance(location, ST_Point($1, $2)::geography) as distance,
            working_hours, services
        FROM pickup_points
        WHERE 
            provider_id = (SELECT id FROM delivery_providers WHERE code = 'posta_srbije')
            AND is_active = true
        ORDER BY location <-> ST_Point($1, $2)::geography
        LIMIT 1
    `
    
    var office PickupPoint
    err := p.db.QueryRow(query, lng, lat).Scan(&office)
    
    return &office, err
}
```

### 3. Компоненты для выбора доставки

#### Выбор службы доставки:
```typescript
// src/components/Delivery/DeliveryServiceSelector.tsx
interface DeliveryOption {
    provider: 'dexpress' | 'posta_srbije';
    price: number;
    estimatedDays: string;
    services: string[];
}

export const DeliveryServiceSelector = ({ 
    fromAddress, 
    toAddress, 
    weight,
    onSelect 
}) => {
    const { data: options, isLoading } = useQuery({
        queryKey: ['delivery-options', fromAddress, toAddress, weight],
        queryFn: () => api.getDeliveryOptions({ fromAddress, toAddress, weight })
    });
    
    return (
        <div className="space-y-4">
            {options?.map(option => (
                <DeliveryOptionCard
                    key={option.provider}
                    provider={option.provider}
                    price={option.price}
                    estimatedDays={option.estimatedDays}
                    services={option.services}
                    onSelect={() => onSelect(option)}
                />
            ))}
        </div>
    );
};

// Карточка опции доставки
const DeliveryOptionCard = ({ provider, price, estimatedDays, services, onSelect }) => {
    const providerInfo = {
        dexpress: {
            name: 'D Express',
            logo: '/images/dexpress-logo.png',
            color: 'bg-red-500',
            description: 'Брза достава широм Србије'
        },
        posta_srbije: {
            name: 'Пошта Србије',
            logo: '/images/posta-srbije-logo.png', 
            color: 'bg-blue-600',
            description: 'Национална поштанска служба'
        }
    };
    
    const info = providerInfo[provider];
    
    return (
        <Card 
            className="cursor-pointer hover:shadow-lg transition-shadow"
            onClick={onSelect}
        >
            <CardContent className="flex items-center justify-between p-4">
                <div className="flex items-center space-x-4">
                    <img src={info.logo} alt={info.name} className="h-12 w-auto" />
                    <div>
                        <h3 className="font-semibold">{info.name}</h3>
                        <p className="text-sm text-gray-600">{info.description}</p>
                        <p className="text-sm mt-1">
                            Достава за {estimatedDays}
                        </p>
                    </div>
                </div>
                <div className="text-right">
                    <p className="text-2xl font-bold">{price} РСД</p>
                    {services.includes('cod') && (
                        <Badge variant="secondary">Плаћање поузећем</Badge>
                    )}
                </div>
            </CardContent>
        </Card>
    );
};
```

#### Карта пунктов выдачи:
```typescript
// src/components/Delivery/PickupPointsMap.tsx
export const PickupPointsMap = ({ provider, userLocation }) => {
    const { data: points } = useQuery({
        queryKey: ['pickup-points', provider],
        queryFn: () => api.getPickupPoints(provider)
    });
    
    return (
        <Map
            center={userLocation || BELGRADE_CENTER}
            zoom={12}
        >
            {/* Маркеры почтовых отделений / пунктов выдачи */}
            {points?.map(point => (
                <Marker
                    key={point.id}
                    position={[point.lat, point.lng]}
                    icon={getProviderIcon(provider)}
                >
                    <Popup>
                        <div className="p-2">
                            <h4 className="font-semibold">{point.name}</h4>
                            <p className="text-sm">{point.address}</p>
                            <WorkingHours hours={point.workingHours} />
                            <Button
                                size="sm"
                                className="mt-2 w-full"
                                onClick={() => selectPickupPoint(point)}
                            >
                                Изабери
                            </Button>
                        </div>
                    </Popup>
                </Marker>
            ))}
            
            {/* Круг радиуса доступности */}
            {userLocation && (
                <Circle
                    center={userLocation}
                    radius={5000} // 5км
                    pathOptions={{
                        color: 'blue',
                        fillOpacity: 0.1
                    }}
                />
            )}
        </Map>
    );
};
```

### 4. Отслеживание доставки

```typescript
// src/components/Delivery/DeliveryTracking.tsx
export const DeliveryTracking = ({ orderId }) => {
    const { data: tracking, isLoading } = useQuery({
        queryKey: ['delivery-tracking', orderId],
        queryFn: () => api.getDeliveryTracking(orderId),
        refetchInterval: 60000 // обновление каждую минуту
    });
    
    if (!tracking) return null;
    
    const statusSteps = {
        pending: { label: 'Чека преузимање', icon: Package, color: 'gray' },
        picked_up: { label: 'Преузето', icon: Truck, color: 'blue' },
        in_transit: { label: 'У транзиту', icon: Navigation, color: 'blue' },
        out_for_delivery: { label: 'Испорука у току', icon: MapPin, color: 'orange' },
        delivered: { label: 'Испоручено', icon: CheckCircle, color: 'green' }
    };
    
    return (
        <Card>
            <CardHeader>
                <CardTitle>Праћење пошиљке</CardTitle>
                <p className="text-sm text-gray-600">
                    Број: {tracking.trackingNumber}
                </p>
            </CardHeader>
            <CardContent>
                {/* Статус доставки */}
                <div className="space-y-4">
                    {Object.entries(statusSteps).map(([status, info]) => {
                        const isPast = getStatusOrder(status) <= getStatusOrder(tracking.status);
                        const isCurrent = status === tracking.status;
                        
                        return (
                            <div 
                                key={status}
                                className={`flex items-center space-x-3 ${
                                    isPast ? 'text-gray-900' : 'text-gray-400'
                                }`}
                            >
                                <info.icon 
                                    className={`h-6 w-6 ${
                                        isCurrent ? `text-${info.color}-600` : ''
                                    }`}
                                />
                                <div className="flex-1">
                                    <p className="font-medium">{info.label}</p>
                                    {isCurrent && tracking.lastLocation && (
                                        <p className="text-sm text-gray-600">
                                            {tracking.lastLocation}
                                        </p>
                                    )}
                                </div>
                                {isPast && (
                                    <CheckIcon className="h-5 w-5 text-green-600" />
                                )}
                            </div>
                        );
                    })}
                </div>
                
                {/* Процењено време */}
                {tracking.estimatedDelivery && (
                    <Alert className="mt-4">
                        <AlertDescription>
                            Процењена испорука: {formatDate(tracking.estimatedDelivery)}
                        </AlertDescription>
                    </Alert>
                )}
            </CardContent>
        </Card>
    );
};
```

### 5. Расчет зон доставки

```sql
-- Функция определения зоны доставки для Пошты Србије
CREATE OR REPLACE FUNCTION calculate_postal_zone(
    from_postal VARCHAR,
    to_postal VARCHAR
) RETURNS INT AS $$
DECLARE
    from_region VARCHAR;
    to_region VARCHAR;
BEGIN
    -- Определяем регионы по первым цифрам индекса
    from_region := CASE 
        WHEN from_postal LIKE '11%' THEN 'BEOGRAD'
        WHEN from_postal LIKE '21%' THEN 'NOVI_SAD'
        WHEN from_postal LIKE '18%' THEN 'NIS'
        WHEN from_postal LIKE '34%' THEN 'KRAGUJEVAC'
        ELSE 'OTHER'
    END;
    
    to_region := CASE 
        WHEN to_postal LIKE '11%' THEN 'BEOGRAD'
        WHEN to_postal LIKE '21%' THEN 'NOVI_SAD'
        WHEN to_postal LIKE '18%' THEN 'NIS'
        WHEN to_postal LIKE '34%' THEN 'KRAGUJEVAC'
        ELSE 'OTHER'
    END;
    
    -- Зона 1: внутри города
    IF from_region = to_region AND from_region != 'OTHER' THEN
        RETURN 1;
    -- Зона 2: между крупными городами
    ELSIF from_region != 'OTHER' AND to_region != 'OTHER' THEN
        RETURN 2;
    -- Зона 3: остальное
    ELSE
        RETURN 3;
    END IF;
END;
$$ LANGUAGE plpgsql;
```

## 📊 Мониторинг доставок

```typescript
// src/components/Admin/DeliveryDashboard.tsx
export const DeliveryDashboard = () => {
    const { data: stats } = useDeliveryStats();
    
    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* Статистика по провайдерам */}
            <Card>
                <CardHeader>
                    <CardTitle>По службама доставе</CardTitle>
                </CardHeader>
                <CardContent>
                    <PieChart
                        data={[
                            { name: 'D Express', value: stats.dexpress.count },
                            { name: 'Пошта Србије', value: stats.posta.count }
                        ]}
                    />
                </CardContent>
            </Card>
            
            {/* Средний срок доставки */}
            <Card>
                <CardHeader>
                    <CardTitle>Средњи рок испоруке</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-2">
                        <div className="flex justify-between">
                            <span>D Express:</span>
                            <span className="font-bold">{stats.dexpress.avgDays} дана</span>
                        </div>
                        <div className="flex justify-between">
                            <span>Пошта Србије:</span>
                            <span className="font-bold">{stats.posta.avgDays} дана</span>
                        </div>
                    </div>
                </CardContent>
            </Card>
            
            {/* Проблемные доставки */}
            <Card>
                <CardHeader>
                    <CardTitle>Проблемне испоруке</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-2">
                        <Alert variant="warning">
                            <AlertDescription>
                                {stats.delayed.count} кашњења (>{stats.delayed.threshold} дана)
                            </AlertDescription>
                        </Alert>
                        <Button 
                            variant="outline" 
                            size="sm"
                            onClick={() => navigate('/admin/deliveries/delayed')}
                        >
                            Прегледај
                        </Button>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};
```

## 🔧 Конфигурация

```yaml
# .env
# D Express
DEXPRESS_API_URL=https://api.dexpress.rs/v1
DEXPRESS_API_KEY=your_api_key_here
DEXPRESS_TEST_MODE=true

# Пошта Србије
POSTA_SRBIJE_API_URL=https://api.posta.rs/v2
POSTA_SRBIJE_USERNAME=your_username
POSTA_SRBIJE_PASSWORD=your_password

# Настройки доставки
DELIVERY_DEFAULT_PROVIDER=dexpress
DELIVERY_FALLBACK_PROVIDER=posta_srbije
DELIVERY_CACHE_TTL=3600
```

## 📈 KPI для доставки

| Метрика | Цель | Текущее |
|---------|------|---------|
| Средний срок доставки | < 3 дня | - |
| % вовремя доставленных | > 95% | - |
| Стоимость доставки / заказ | < 300 RSD | - |
| % успешных отслеживаний | > 98% | - |

---

Это реалистичная интеграция с D Express и Поштой Србије, без велокурьеров и прочей экзотики. Фокус на надежности и простоте использования.