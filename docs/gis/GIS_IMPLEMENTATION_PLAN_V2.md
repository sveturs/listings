# 🗺️ План реализации GIS модуля v2.0 - С учетом витрин

**Версия**: 2.0  
**Дата**: 2025-01-10  
**Изменения**: Добавлена поддержка витрин с гибким отображением товаров

## 📋 Оглавление

1. [Обновленные требования](#обновленные-требования)
2. [Архитектурные изменения](#архитектурные-изменения)
3. [База данных - дополнения](#база-данных-дополнения)
4. [API - новые endpoints](#api-новые-endpoints)
5. [Frontend компоненты](#frontend-компоненты)
6. [Актуализированный план реализации](#актуализированный-план-реализации)

---

## 1. Обновленные требования

### 1.1 Типы отображения на карте

1. **Обычные объявления** - стандартные маркеры с точной/размытой геолокацией
2. **Витрины** - специальные маркеры по адресу магазина
3. **Товары витрин** - могут отображаться:
   - Внутри витрины (не показываются отдельно)
   - На карте с собственной геолокацией (недвижимость, прокат)

### 1.2 Уровни приватности

- **exact** - точные координаты
- **district** - размытие ~500м (показывается район)
- **area** - размытие ~2км (показывается область)

### 1.3 Бизнес-сценарии

1. **Магазин** - витрина в одной точке, товары не показываются отдельно
2. **Агентство недвижимости** - офис + объекты по городу
3. **Сервис проката** - точка выдачи + локации товаров
4. **Частный продавец** - размытая локация для безопасности

---

## 2. Архитектурные изменения

### 2.1 Структура данных

```typescript
// Расширенная модель витрины
interface StorefrontGeoSettings {
  // Настройки отображения витрины
  displayLocation: 'exact' | 'district' | 'area';
  showOnMap: boolean;
  
  // Настройки товаров
  productsDisplay: {
    showOnMap: boolean;           // показывать ли товары отдельно
    inheritLocation: boolean;      // наследовать локацию витрины
    defaultPrivacy: 'exact' | 'district' | 'area';
    allowCustomLocation: boolean;  // разрешить товарам свою локацию
  };
  
  // Визуальные настройки
  mapIcon?: string;               // кастомная иконка
  showProductCount: boolean;
  showRating: boolean;
  showWorkingStatus: boolean;
}

// Расширенная модель товара
interface ProductGeoData {
  location?: {
    lat: number;
    lng: number;
    privacy: 'exact' | 'district' | 'area';
    address?: string;
    customAddress?: string;  // "Район Земун" вместо точного адреса
  };
  inheritFromStorefront: boolean;
}
```

### 2.2 Сервисная архитектура

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Map Service   │     │Storefront Service│    │ Privacy Service │
│                 │────▶│                  │────▶│                 │
│ - Clustering    │     │ - Settings       │     │ - Blur coords   │
│ - Filtering     │     │ - Product rules  │     │ - Address mask  │
│ - Rendering     │     │ - Display logic  │     │ - Access control│
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

---

## 3. База данных - дополнения

### 3.1 Новые таблицы

```sql
-- Настройки отображения витрин на карте
CREATE TABLE storefront_map_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    storefront_id UUID NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    display_location VARCHAR(20) DEFAULT 'exact',
    show_on_map BOOLEAN DEFAULT true,
    products_show_on_map BOOLEAN DEFAULT false,
    products_inherit_location BOOLEAN DEFAULT true,
    products_default_privacy VARCHAR(20) DEFAULT 'district',
    products_allow_custom_location BOOLEAN DEFAULT false,
    map_icon VARCHAR(100),
    show_product_count BOOLEAN DEFAULT true,
    show_rating BOOLEAN DEFAULT true,
    show_working_status BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(storefront_id)
);

-- Геоданные товаров витрин
CREATE TABLE storefront_product_geo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    storefront_id UUID NOT NULL REFERENCES storefronts(id),
    location GEOGRAPHY(POINT, 4326),
    location_privacy VARCHAR(20) DEFAULT 'district',
    custom_address TEXT,
    inherit_from_storefront BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(product_id)
);

-- Индексы
CREATE INDEX idx_storefront_product_geo_location 
    ON storefront_product_geo USING GIST(location);
CREATE INDEX idx_storefront_product_geo_storefront 
    ON storefront_product_geo(storefront_id);

-- Расширение listings_geo для поддержки приватности
ALTER TABLE listings_geo 
    ADD COLUMN location_privacy VARCHAR(20) DEFAULT 'exact',
    ADD COLUMN blurred_location GEOGRAPHY(POINT, 4326),
    ADD COLUMN custom_address TEXT;

-- Функция для размытия координат
CREATE OR REPLACE FUNCTION blur_coordinates(
    p_location GEOGRAPHY,
    p_privacy VARCHAR
) RETURNS GEOGRAPHY AS $$
DECLARE
    lat FLOAT;
    lng FLOAT;
    radius FLOAT;
    angle FLOAT;
    distance FLOAT;
BEGIN
    IF p_privacy = 'exact' THEN
        RETURN p_location;
    END IF;
    
    lat := ST_Y(p_location::geometry);
    lng := ST_X(p_location::geometry);
    
    -- Радиус размытия
    radius := CASE p_privacy
        WHEN 'district' THEN 0.005  -- ~500м
        WHEN 'area' THEN 0.02       -- ~2км
        ELSE 0
    END;
    
    -- Случайное смещение
    angle := random() * 2 * PI();
    distance := random() * radius;
    
    lat := lat + distance * cos(angle);
    lng := lng + distance * sin(angle);
    
    RETURN ST_SetSRID(ST_MakePoint(lng, lat), 4326)::geography;
END;
$$ LANGUAGE plpgsql;

-- Триггер для автоматического размытия
CREATE OR REPLACE FUNCTION update_blurred_location()
RETURNS TRIGGER AS $$
BEGIN
    NEW.blurred_location := blur_coordinates(NEW.location, NEW.location_privacy);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_blur_listing_location
    BEFORE INSERT OR UPDATE ON listings_geo
    FOR EACH ROW
    EXECUTE FUNCTION update_blurred_location();
```

### 3.2 Представления для карты

```sql
-- Объединенное представление для карты
CREATE OR REPLACE VIEW map_items AS
SELECT 
    'listing' as item_type,
    l.id,
    l.title,
    CASE 
        WHEN lg.location_privacy = 'exact' THEN lg.location
        ELSE lg.blurred_location
    END as display_location,
    lg.location_privacy,
    lg.custom_address,
    l.price,
    l.category_id,
    NULL::UUID as storefront_id,
    l.created_at
FROM marketplace_listings l
JOIN listings_geo lg ON l.id = lg.listing_id
WHERE l.is_active = true

UNION ALL

SELECT 
    'storefront' as item_type,
    s.id,
    s.name as title,
    CASE 
        WHEN sms.display_location = 'exact' THEN 
            ST_SetSRID(ST_MakePoint(s.longitude, s.latitude), 4326)::geography
        ELSE 
            blur_coordinates(
                ST_SetSRID(ST_MakePoint(s.longitude, s.latitude), 4326)::geography,
                sms.display_location
            )
    END as display_location,
    sms.display_location as location_privacy,
    s.address as custom_address,
    NULL::DECIMAL as price,
    NULL::INT as category_id,
    s.id as storefront_id,
    s.created_at
FROM storefronts s
LEFT JOIN storefront_map_settings sms ON s.id = sms.storefront_id
WHERE s.is_active = true 
    AND COALESCE(sms.show_on_map, true) = true

UNION ALL

SELECT 
    'storefront_product' as item_type,
    sp.id,
    sp.name as title,
    CASE 
        WHEN spg.inherit_from_storefront THEN
            CASE 
                WHEN sms.display_location = 'exact' THEN 
                    ST_SetSRID(ST_MakePoint(s.longitude, s.latitude), 4326)::geography
                ELSE 
                    blur_coordinates(
                        ST_SetSRID(ST_MakePoint(s.longitude, s.latitude), 4326)::geography,
                        sms.display_location
                    )
            END
        ELSE
            CASE 
                WHEN spg.location_privacy = 'exact' THEN spg.location
                ELSE blur_coordinates(spg.location, spg.location_privacy)
            END
    END as display_location,
    COALESCE(spg.location_privacy, sms.products_default_privacy) as location_privacy,
    spg.custom_address,
    sp.price,
    sp.category_id,
    sp.storefront_id,
    sp.created_at
FROM storefront_products sp
JOIN storefronts s ON sp.storefront_id = s.id
LEFT JOIN storefront_product_geo spg ON sp.id = spg.product_id
LEFT JOIN storefront_map_settings sms ON s.id = sms.storefront_id
WHERE sp.is_active = true 
    AND s.is_active = true
    AND COALESCE(sms.products_show_on_map, false) = true;
```

---

## 4. API - новые endpoints

### 4.1 Настройки карты для витрин

```yaml
# Получить/обновить настройки отображения витрины
GET/PUT /api/v1/storefronts/{id}/map-settings:
  requestBody:
    schema:
      $ref: '#/components/schemas/StorefrontMapSettings'

# Массовое обновление геоданных товаров
POST /api/v1/storefronts/{id}/products/geo:
  requestBody:
    schema:
      type: object
      properties:
        action:
          type: string
          enum: [set_inherit, set_privacy, update_locations]
        privacy:
          type: string
          enum: [exact, district, area]
        products:
          type: array
          items:
            type: object
            properties:
              product_id: string
              location: 
                lat: number
                lng: number
              custom_address: string

# Получить элементы для карты с учетом приватности
POST /api/v1/gis/map/items:
  requestBody:
    schema:
      type: object
      properties:
        bounds:
          $ref: '#/components/schemas/Bounds'
        types:
          type: array
          items:
            enum: [listings, storefronts, storefront_products]
        filters:
          type: object
          properties:
            categories: array
            price_range: object
            only_exact_location: boolean
            storefront_id: string
```

### 4.2 Расширенный поиск

```yaml
# Поиск с учетом типов элементов
POST /api/v1/gis/search/unified:
  requestBody:
    schema:
      type: object
      properties:
        center:
          $ref: '#/components/schemas/Point'
        radius: integer
        types:
          type: array
          items:
            enum: [listings, storefronts, storefront_products]
        group_by_storefront: boolean
        include_blurred: boolean
```

---

## 5. Frontend компоненты

### 5.1 Компонент настроек карты для витрины

```typescript
// src/components/storefronts/settings/MapSettings.tsx
import { useState } from 'react';
import { Switch, RadioGroup, Alert } from '@/components/ui';

interface MapSettingsProps {
  storefront: Storefront;
  onSave: (settings: StorefrontMapSettings) => void;
}

export const MapSettings: React.FC<MapSettingsProps> = ({ 
  storefront, 
  onSave 
}) => {
  const [settings, setSettings] = useState<StorefrontMapSettings>(
    storefront.mapSettings || defaultSettings
  );

  return (
    <div className="space-y-6">
      <div className="bg-white p-6 rounded-lg shadow">
        <h3 className="text-lg font-semibold mb-4">
          Отображение витрины на карте
        </h3>
        
        <div className="space-y-4">
          <Switch
            label="Показывать витрину на карте"
            checked={settings.showOnMap}
            onChange={(showOnMap) => setSettings({...settings, showOnMap})}
          />
          
          {settings.showOnMap && (
            <RadioGroup
              label="Точность отображения адреса"
              value={settings.displayLocation}
              onChange={(displayLocation) => 
                setSettings({...settings, displayLocation})
              }
              options={[
                { value: 'exact', label: 'Точный адрес' },
                { value: 'district', label: 'Только район (~500м)' },
                { value: 'area', label: 'Область города (~2км)' }
              ]}
            />
          )}
        </div>
      </div>

      <div className="bg-white p-6 rounded-lg shadow">
        <h3 className="text-lg font-semibold mb-4">
          Отображение товаров на карте
        </h3>
        
        <Alert type="info" className="mb-4">
          Подходит для агентств недвижимости, прокатов и сервисов 
          с распределенными товарами
        </Alert>
        
        <div className="space-y-4">
          <Switch
            label="Показывать товары отдельно на карте"
            checked={settings.productsDisplay.showOnMap}
            onChange={(showOnMap) => 
              setSettings({
                ...settings, 
                productsDisplay: {...settings.productsDisplay, showOnMap}
              })
            }
          />
          
          {settings.productsDisplay.showOnMap && (
            <>
              <Switch
                label="Разрешить товарам собственную локацию"
                checked={settings.productsDisplay.allowCustomLocation}
                help="Если выключено, все товары будут показаны по адресу витрины"
              />
              
              <RadioGroup
                label="Приватность локации товаров по умолчанию"
                value={settings.productsDisplay.defaultPrivacy}
                options={[
                  { value: 'exact', label: 'Точный адрес' },
                  { value: 'district', label: 'Район' },
                  { value: 'area', label: 'Область' }
                ]}
              />
            </>
          )}
        </div>
      </div>

      <div className="bg-white p-6 rounded-lg shadow">
        <h3 className="text-lg font-semibold mb-4">
          Дополнительные настройки
        </h3>
        
        <div className="space-y-3">
          <Switch
            label="Показывать количество товаров"
            checked={settings.showProductCount}
          />
          <Switch
            label="Показывать рейтинг"
            checked={settings.showRating}
          />
          <Switch
            label="Показывать статус работы (открыто/закрыто)"
            checked={settings.showWorkingStatus}
          />
        </div>
      </div>

      <button
        onClick={() => onSave(settings)}
        className="w-full btn btn-primary"
      >
        Сохранить настройки
      </button>
    </div>
  );
};
```

### 5.2 Расширенный компонент карты

```typescript
// src/components/GIS/Map/UnifiedMap.tsx
import { useEffect, useState } from 'react';
import mapboxgl from 'mapbox-gl';
import { useMapItems } from '@/hooks/gis/useMapItems';

interface UnifiedMapProps {
  types?: ('listings' | 'storefronts' | 'storefront_products')[];
  onItemClick?: (item: MapItem) => void;
}

export const UnifiedMap: React.FC<UnifiedMapProps> = ({
  types = ['listings', 'storefronts', 'storefront_products'],
  onItemClick
}) => {
  const { items, isLoading } = useMapItems({ types });
  
  // Различные стили маркеров
  const createMarker = (item: MapItem) => {
    const el = document.createElement('div');
    
    switch (item.type) {
      case 'storefront':
        el.className = 'marker-storefront';
        el.innerHTML = `
          <div class="storefront-marker">
            ${item.logo ? `<img src="${item.logo}" />` : '🏪'}
            <span class="count">${item.productCount || 0}</span>
          </div>
        `;
        break;
        
      case 'storefront_product':
        el.className = 'marker-product';
        // Размытая локация - полупрозрачный круг
        if (item.locationPrivacy !== 'exact') {
          el.classList.add('blurred');
          el.style.opacity = '0.7';
        }
        break;
        
      case 'listing':
        el.className = 'marker-listing';
        if (item.locationPrivacy !== 'exact') {
          el.classList.add('blurred');
        }
        break;
    }
    
    return new mapboxgl.Marker(el)
      .setLngLat([item.location.lng, item.location.lat]);
  };
  
  // Группировка товаров по витринам при необходимости
  const groupItemsByStorefront = (items: MapItem[]) => {
    const groups = new Map<string, MapItem[]>();
    
    items.forEach(item => {
      if (item.storefrontId) {
        const group = groups.get(item.storefrontId) || [];
        group.push(item);
        groups.set(item.storefrontId, group);
      }
    });
    
    return groups;
  };
  
  // ... остальная логика карты
};
```

---

## 6. Актуализированный план реализации

### Фаза 1: Базовая инфраструктура (2 недели)

**Неделя 1: База данных и API**
- [ ] Создать таблицы для настроек витрин и геоданных товаров
- [ ] Реализовать функции размытия координат
- [ ] Создать представление map_items
- [ ] Базовые API endpoints для настроек

**Неделя 2: Интеграция с существующей системой**
- [ ] Обновить модели данных в Go
- [ ] Расширить OpenSearch маппинги
- [ ] Добавить поддержку приватности в существующие endpoints
- [ ] Написать тесты

### Фаза 2: Frontend базовый функционал (2 недели)

**Неделя 3: Компоненты настроек**
- [ ] UI для настроек отображения витрины
- [ ] Форма управления геоданными товаров
- [ ] Интеграция с API

**Неделя 4: Обновление карты**
- [ ] Поддержка разных типов маркеров
- [ ] Визуализация размытых локаций
- [ ] Фильтры по типам элементов

### Фаза 3: Расширенный функционал (3 недели)

**Неделя 5: Умная группировка**
- [ ] Кластеризация с учетом витрин
- [ ] Связывание товаров с витринами визуально
- [ ] Оптимизация отображения

**Неделя 6: Поиск и фильтрация**
- [ ] Унифицированный поиск по всем типам
- [ ] Фильтры с учетом приватности
- [ ] Поиск "только с точным адресом"

**Неделя 7: Аналитика и оптимизация**
- [ ] Статистика использования карты
- [ ] A/B тестирование разных режимов
- [ ] Оптимизация производительности

### Фаза 4: Интеграция с доставкой (2 недели)

**Неделя 8: Зоны доставки для витрин**
- [ ] Отображение зон с учетом типа витрины
- [ ] Расчет доставки для распределенных товаров

**Неделя 9: Финализация**
- [ ] Документация для продавцов
- [ ] Обучающие материалы
- [ ] Запуск в продакшн

---

## 📊 Метрики успеха

### Технические
- Время загрузки карты < 2 сек
- Поддержка 10k+ маркеров
- Точность геопоиска > 95%

### Продуктовые  
- Использование карты витринами > 60%
- Конверсия из карты +25%
- Удовлетворенность продавцов > 4.5/5

### Безопасность
- 0 утечек точных адресов при размытии
- 100% соблюдение настроек приватности

---

**План готов к реализации!** 

Основные изменения:
1. Гибкая система отображения для разных типов бизнеса
2. Трехуровневая система приватности
3. Раздельное управление витринами и их товарами
4. Визуальное различие точных и размытых локаций