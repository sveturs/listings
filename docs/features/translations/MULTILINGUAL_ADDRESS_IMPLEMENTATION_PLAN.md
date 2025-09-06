# План реализации мультиязычного хранения адресов

## Обзор
Реализация системы хранения адресов на трех языках (сербский, английский, русский) для обеспечения корректного отображения адресов в зависимости от выбранной локали без дополнительных запросов к API геокодирования.

## Архитектурное решение
Использование отдельных полей в БД для каждого языка как наиболее производительное и простое решение.

## Детальный план реализации

### 1. База данных (Backend)

#### 1.1 Создание миграции для новых полей
```sql
-- backend/migrations/00XXX_add_multilingual_addresses.up.sql
ALTER TABLE marketplace_listings 
ADD COLUMN location_sr VARCHAR(255),
ADD COLUMN location_en VARCHAR(255),
ADD COLUMN location_ru VARCHAR(255);

-- Создание индексов для поиска
CREATE INDEX idx_marketplace_listings_location_sr ON marketplace_listings(location_sr) WHERE location_sr IS NOT NULL;
CREATE INDEX idx_marketplace_listings_location_en ON marketplace_listings(location_en) WHERE location_en IS NOT NULL;
CREATE INDEX idx_marketplace_listings_location_ru ON marketplace_listings(location_ru) WHERE location_ru IS NOT NULL;

-- Миграция существующих данных (копируем в location_ru, т.к. большинство адресов на русском)
UPDATE marketplace_listings 
SET location_ru = location 
WHERE location IS NOT NULL;
```

```sql
-- backend/migrations/00XXX_add_multilingual_addresses.down.sql
DROP INDEX IF EXISTS idx_marketplace_listings_location_sr;
DROP INDEX IF EXISTS idx_marketplace_listings_location_en;
DROP INDEX IF EXISTS idx_marketplace_listings_location_ru;

ALTER TABLE marketplace_listings 
DROP COLUMN location_sr,
DROP COLUMN location_en,
DROP COLUMN location_ru;
```

#### 1.2 Обновление моделей
```go
// backend/internal/domain/marketplace.go
type MarketplaceListing struct {
    // ... существующие поля
    Location   string  `json:"location"`         // Оставляем для обратной совместимости
    LocationSr *string `json:"location_sr"`      // Адрес на сербском
    LocationEn *string `json:"location_en"`      // Адрес на английском
    LocationRu *string `json:"location_ru"`      // Адрес на русском
}
```

#### 1.3 Обновление репозитория
- Обновить SQL запросы в `GetListingByID` для выборки новых полей
- Обновить методы `CreateListing` и `UpdateListing` для сохранения адресов на всех языках
- Добавить логику выбора адреса на нужном языке в зависимости от заголовка `Accept-Language`

### 2. API (Backend)

#### 2.1 Создание эндпоинта для мультиязычного геокодирования
```go
// POST /api/v1/gis/geocode/multilingual
type MultilingualGeocodeRequest struct {
    Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
    Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
}

type MultilingualGeocodeResponse struct {
    AddressSr string `json:"address_sr"`
    AddressEn string `json:"address_en"`
    AddressRu string `json:"address_ru"`
}
```

#### 2.2 Реализация в GeocodingService
```go
func (s *GeocodingService) GetMultilingualAddress(ctx context.Context, lat, lng float64) (*MultilingualGeocodeResponse, error) {
    // Параллельные запросы для всех языков
    var wg sync.WaitGroup
    var addressSr, addressEn, addressRu string
    var errSr, errEn, errRu error
    
    wg.Add(3)
    
    go func() {
        defer wg.Done()
        result, err := s.ReverseGeocode(ctx, Point{Lat: lat, Lng: lng}, "sr")
        if err == nil && result != nil {
            addressSr = result.PlaceName
        }
        errSr = err
    }()
    
    // Аналогично для en и ru
    
    wg.Wait()
    
    return &MultilingualGeocodeResponse{
        AddressSr: addressSr,
        AddressEn: addressEn,
        AddressRu: addressRu,
    }, nil
}
```

### 3. Frontend

#### 3.1 Обновление процесса создания объявления

##### В `CreateListingContext`:
```typescript
interface ListingData {
  // ... существующие поля
  location: string;        // Для обратной совместимости
  location_sr?: string;
  location_en?: string;
  location_ru?: string;
}
```

##### В `LocationStep`:
```typescript
// При выборе адреса через карту или поиск
const handleLocationSelect = async (lat: number, lng: number) => {
  // Получаем адреса на всех языках
  const response = await fetch('/api/v1/gis/geocode/multilingual', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ latitude: lat, longitude: lng })
  });
  
  const { address_sr, address_en, address_ru } = await response.json();
  
  // Сохраняем все варианты
  updateListingData({
    latitude: lat,
    longitude: lng,
    location: address_sr, // По умолчанию сербский
    location_sr: address_sr,
    location_en: address_en,
    location_ru: address_ru
  });
};
```

#### 3.2 Обновление отображения объявлений

##### В `marketplace/[id]/page.tsx`:
```typescript
// Удаляем fetchLocalizedAddress и используем готовые адреса
const displayAddress = useMemo(() => {
  if (!listing) return '';
  
  // Выбираем адрес на нужном языке
  if (locale === 'ru' && listing.location_ru) {
    return listing.location_ru;
  } else if (locale === 'en' && listing.location_en) {
    return listing.location_en;
  }
  
  // По умолчанию сербский или fallback на location
  return listing.location_sr || listing.location;
}, [listing, locale]);

// Используем в рендере
{formatAddressWithPrivacy(displayAddress, listing.location_privacy)}
```

#### 3.3 Аналогичные изменения для:
- Компонентов карты (MarkerHoverPopup, MarkerClickPopup)
- Карточек товаров (MarketplaceCard)
- Страниц витрин

### 4. Миграция существующих данных

#### 4.1 Скрипт для обновления адресов
```typescript
// scripts/migrate-addresses.ts
async function migrateAddresses() {
  // Получаем все объявления с координатами
  const listings = await getListingsWithCoordinates();
  
  for (const listing of listings) {
    if (listing.latitude && listing.longitude) {
      try {
        // Получаем адреса на всех языках
        const addresses = await getMultilingualAddress(
          listing.latitude,
          listing.longitude
        );
        
        // Обновляем в БД
        await updateListingAddresses(listing.id, addresses);
      } catch (error) {
        console.error(`Failed to migrate listing ${listing.id}:`, error);
      }
    }
  }
}
```

### 5. Тестирование

#### 5.1 Unit тесты
- Тест мультиязычного геокодирования
- Тест выбора правильного языка при отображении
- Тест миграции данных

#### 5.2 Интеграционные тесты
- Создание объявления с автоматическим получением адресов
- Отображение адресов в разных локалях
- Поиск по адресам на разных языках

### 6. Оптимизации

#### 6.1 Кэширование геокодирования
- Использовать существующий кэш для хранения результатов по координатам
- Ключ кэша: `geocode:${lat}:${lng}:${lang}`

#### 6.2 Batch обработка
- При импорте товаров получать адреса пачками
- Ограничение rate limit для API геокодирования

### 7. Обратная совместимость

- Поле `location` остается для старых клиентов
- При отсутствии локализованных адресов используется `location`
- Постепенная миграция через фоновые задачи

## Преимущества решения

1. **Производительность**: Нет дополнительных запросов при отображении
2. **Консистентность**: Адреса всегда корректны для каждого языка
3. **SEO**: Правильные адреса для индексации на разных языках
4. **UX**: Мгновенное переключение языков без задержек

## Оценка времени

- Backend (миграция + API): 4-6 часов
- Frontend (создание + отображение): 3-4 часа
- Миграция данных: 2-3 часа
- Тестирование: 2-3 часа

**Итого**: 11-16 часов

## Риски и митигация

1. **Rate limiting геокодирования**
   - Митигация: Batch обработка, кэширование, throttling

2. **Неполные данные при миграции**
   - Митигация: Fallback на существующее поле location

3. **Увеличение размера БД**
   - Митигация: Индексы только на NOT NULL значениях