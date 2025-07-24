# 🏷️ Руководство по реализации системы скидок

## 📋 Обзор существующей функциональности

### ✅ Что уже реализовано в системе:

1. **База данных**:
   - Таблица `price_history` автоматически отслеживает все изменения цен
   - Поле `metadata` в `marketplace_listings` хранит информацию о скидках
   - SQL функция `check_price_manipulation()` защищает от фейковых скидок
   - Триггеры автоматически вычисляют скидки при изменении цены

2. **Backend API**:
   - `GET /api/v1/marketplace/listings/{id}/price-history` - готовый эндпоинт
   - Модель `MarketplaceListing` уже имеет поля `old_price` и `has_discount`
   - Автоматическое вычисление скидок от максимальной цены за период

3. **Frontend**:
   - Базовое отображение старой цены и процента скидки
   - Поддержка в `MarketplaceCard` компоненте

## 🛠️ Пошаговая инструкция реализации

### Шаг 1: Создание компонента DiscountBadge

```tsx
// frontend/svetu/src/components/ui/DiscountBadge.tsx

import React from 'react';
import { TrendingDown } from 'lucide-react';

interface DiscountBadgeProps {
  oldPrice: number;
  currentPrice: number;
  onClick?: () => void;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const DiscountBadge: React.FC<DiscountBadgeProps> = ({
  oldPrice,
  currentPrice,
  onClick,
  size = 'md',
  className = ''
}) => {
  const discountPercent = Math.round(((oldPrice - currentPrice) / oldPrice) * 100);
  
  // Не показывать скидки менее 5%
  if (discountPercent < 5) return null;
  
  const sizeClasses = {
    sm: 'text-xs px-2 py-1',
    md: 'text-sm px-3 py-1.5',
    lg: 'text-base px-4 py-2'
  };
  
  return (
    <button
      onClick={onClick}
      className={`
        badge badge-error gap-1 cursor-pointer 
        hover:scale-105 transition-transform
        ${sizeClasses[size]} ${className}
      `}
      title="Нажмите, чтобы увидеть историю цены"
    >
      <TrendingDown className="w-3 h-3" />
      -{discountPercent}%
    </button>
  );
};
```

### Шаг 2: Создание модалки истории цен

```tsx
// frontend/svetu/src/components/marketplace/PriceHistoryModal.tsx

import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { X, AlertTriangle } from 'lucide-react';

interface PriceHistoryModalProps {
  listingId: number;
  isOpen: boolean;
  onClose: () => void;
}

export const PriceHistoryModal: React.FC<PriceHistoryModalProps> = ({
  listingId,
  isOpen,
  onClose
}) => {
  const [priceHistory, setPriceHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [hasManipulation, setHasManipulation] = useState(false);
  
  useEffect(() => {
    if (isOpen && listingId) {
      fetchPriceHistory();
    }
  }, [isOpen, listingId]);
  
  const fetchPriceHistory = async () => {
    try {
      const response = await fetch(`/api/v1/marketplace/listings/${listingId}/price-history`);
      const data = await response.json();
      
      setPriceHistory(data.data);
      // Проверка на манипуляции (резкий рост > 30% с последующим снижением)
      checkForManipulation(data.data);
      setLoading(false);
    } catch (error) {
      console.error('Error fetching price history:', error);
      setLoading(false);
    }
  };
  
  const checkForManipulation = (history) => {
    // Логика проверки манипуляций
    for (let i = 1; i < history.length; i++) {
      const prevPrice = history[i - 1].price;
      const currPrice = history[i].price;
      const changePercent = ((currPrice - prevPrice) / prevPrice) * 100;
      
      // Если цена выросла более чем на 30%
      if (changePercent > 30) {
        // Проверяем последующее снижение
        for (let j = i + 1; j < history.length; j++) {
          const futurePrice = history[j].price;
          if (futurePrice < prevPrice * 1.1) {
            setHasManipulation(true);
            return;
          }
        }
      }
    }
  };
  
  if (!isOpen) return null;
  
  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-3xl">
        <button
          className="btn btn-sm btn-circle absolute right-2 top-2"
          onClick={onClose}
        >
          <X />
        </button>
        
        <h3 className="font-bold text-lg mb-4">История цены товара</h3>
        
        {hasManipulation && (
          <div className="alert alert-warning mb-4">
            <AlertTriangle className="w-5 h-5" />
            <span>Обнаружены подозрительные изменения цены!</span>
          </div>
        )}
        
        {loading ? (
          <div className="flex justify-center py-8">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        ) : (
          <div className="h-64">
            {/* Здесь будет график с Chart.js */}
            <Line
              data={chartData}
              options={chartOptions}
            />
          </div>
        )}
        
        <div className="modal-action">
          <button className="btn" onClick={onClose}>Закрыть</button>
        </div>
      </div>
    </div>
  );
};
```

### Шаг 3: Интеграция в EnhancedListingCard

```tsx
// Добавить в EnhancedListingCard.tsx

import { DiscountBadge } from '@/components/ui/DiscountBadge';
import { PriceHistoryModal } from '@/components/marketplace/PriceHistoryModal';

// В компоненте:
const [showPriceHistory, setShowPriceHistory] = useState(false);

// В разметке где отображается цена:
<div className="flex items-center gap-2">
  <p className="text-xl font-bold">
    {formatPrice(item.price, item.currency || 'RSD')}
  </p>
  {item.has_discount && item.old_price && (
    <>
      <p className="text-sm line-through text-base-content/50">
        {formatPrice(item.old_price, item.currency || 'RSD')}
      </p>
      <DiscountBadge
        oldPrice={item.old_price}
        currentPrice={item.price}
        onClick={() => setShowPriceHistory(true)}
        size="sm"
      />
    </>
  )}
</div>

{showPriceHistory && (
  <PriceHistoryModal
    listingId={item.id}
    isOpen={showPriceHistory}
    onClose={() => setShowPriceHistory(false)}
  />
)}
```

### Шаг 4: Добавление Black Friday Badge для витрин

```tsx
// frontend/svetu/src/components/storefronts/BlackFridayBadge.tsx

import React from 'react';
import { Zap } from 'lucide-react';

interface BlackFridayBadgeProps {
  discountStats: {
    totalProducts: number;
    discountedProducts: number;
    averageDiscount: number;
  };
}

export const BlackFridayBadge: React.FC<BlackFridayBadgeProps> = ({ discountStats }) => {
  const discountedPercent = (discountStats.discountedProducts / discountStats.totalProducts) * 100;
  
  // Показывать только если >20% товаров со скидками >10%
  if (discountedPercent < 20 || discountStats.averageDiscount < 10) {
    return null;
  }
  
  return (
    <div className="badge badge-lg gap-2 bg-black text-white">
      <Zap className="w-4 h-4 text-yellow-400" />
      <span className="font-bold">BLACK FRIDAY</span>
      <span className="text-xs">
        {Math.round(discountedPercent)}% товаров со скидками
      </span>
    </div>
  );
};
```

### Шаг 5: Backend - добавление статистики скидок для витрин

```go
// backend/internal/proj/storefronts/handler/analytics.go

// Добавить в существующий handler аналитики витрины:
func (h *Handler) GetStorefrontDiscountStats(c *fiber.Ctx) error {
    storefrontID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidStorefrontID")
    }
    
    stats, err := h.storefrontService.GetDiscountStats(storefrontID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.failedToGetStats")
    }
    
    return utils.SuccessResponse(c, stats)
}

// В service добавить метод:
func (s *Service) GetDiscountStats(storefrontID int) (*models.DiscountStats, error) {
    // SQL запрос для получения статистики
    query := `
        SELECT 
            COUNT(*) as total_products,
            COUNT(CASE WHEN has_discount = true THEN 1 END) as discounted_products,
            AVG(CASE 
                WHEN has_discount = true AND old_price > 0 
                THEN ((old_price - price) / old_price * 100) 
                ELSE 0 
            END) as average_discount
        FROM storefront_products
        WHERE storefront_id = $1 AND status = 'active'
    `
    
    var stats models.DiscountStats
    err := s.db.QueryRow(query, storefrontID).Scan(
        &stats.TotalProducts,
        &stats.DiscountedProducts,
        &stats.AverageDiscount,
    )
    
    return &stats, err
}
```

### Шаг 6: Добавление фильтра по скидкам

```tsx
// В компоненте фильтров добавить:
<label className="label cursor-pointer">
  <span className="label-text">Только товары со скидками</span>
  <input
    type="checkbox"
    className="checkbox checkbox-primary"
    checked={filters.onlyDiscounted}
    onChange={(e) => setFilters({
      ...filters,
      onlyDiscounted: e.target.checked
    })}
  />
</label>

// В API запросе добавить параметр:
if (filters.onlyDiscounted) {
  params.append('only_discounted', 'true');
}
```

### Шаг 7: SQL миграция для оптимизации

```sql
-- migrations/000XX_add_discount_indexes.up.sql

-- Индекс для быстрого поиска товаров со скидками
CREATE INDEX IF NOT EXISTS idx_listings_has_discount 
ON marketplace_listings(has_discount) 
WHERE has_discount = true;

-- Индекс для сортировки по размеру скидки
CREATE INDEX IF NOT EXISTS idx_listings_discount_percent 
ON marketplace_listings((metadata->>'discount'->>'discount_percent')::numeric) 
WHERE has_discount = true;

-- Материализованное представление для статистики витрин
CREATE MATERIALIZED VIEW IF NOT EXISTS storefront_discount_stats AS
SELECT 
    s.id as storefront_id,
    COUNT(sp.id) as total_products,
    COUNT(CASE WHEN ml.has_discount = true THEN 1 END) as discounted_products,
    AVG(CASE 
        WHEN ml.has_discount = true AND ml.old_price > 0 
        THEN ((ml.old_price - ml.price) / ml.old_price * 100) 
        ELSE 0 
    END) as average_discount_percent,
    MAX((ml.metadata->'discount'->>'discount_percent')::numeric) as max_discount_percent
FROM storefronts s
LEFT JOIN storefront_products sp ON s.id = sp.storefront_id
LEFT JOIN marketplace_listings ml ON sp.listing_id = ml.id
WHERE sp.status = 'active'
GROUP BY s.id;

-- Обновлять каждые 15 минут
CREATE UNIQUE INDEX ON storefront_discount_stats(storefront_id);
```

## 📊 Метрики для отслеживания

```typescript
// frontend/svetu/src/utils/analytics.ts

// Трекинг клика на бейдж скидки
export const trackDiscountBadgeClick = (listingId: number, discountPercent: number) => {
  if (typeof window !== 'undefined' && window.gtag) {
    window.gtag('event', 'discount_badge_click', {
      event_category: 'engagement',
      event_label: `listing_${listingId}`,
      value: discountPercent
    });
  }
};

// Трекинг просмотра истории цен
export const trackPriceHistoryView = (listingId: number, hasManipulation: boolean) => {
  if (typeof window !== 'undefined' && window.gtag) {
    window.gtag('event', 'price_history_view', {
      event_category: 'engagement',
      event_label: `listing_${listingId}`,
      custom_parameter: hasManipulation ? 'manipulation_detected' : 'normal'
    });
  }
};
```

## ⚠️ Важные моменты

1. **Защита от манипуляций**: Система уже имеет встроенную защиту через `check_price_manipulation()`. НЕ отключайте эту проверку!

2. **Производительность**: Используйте индексы и материализованные представления для быстрой фильтрации по скидкам.

3. **Кэширование**: Кэшируйте статистику скидок витрин на 15 минут, чтобы не нагружать БД.

4. **SEO**: Добавьте structured data для товаров со скидками:
```json
{
  "@type": "Product",
  "offers": {
    "@type": "Offer",
    "price": "75.00",
    "priceCurrency": "RSD",
    "priceValidUntil": "2024-12-31",
    "itemCondition": "https://schema.org/NewCondition",
    "availability": "https://schema.org/InStock",
    "discount": {
      "@type": "Discount",
      "discountPercent": 25,
      "previousPrice": "100.00"
    }
  }
}
```

## 🚀 Последовательность внедрения

1. **День 1**: Создать компоненты DiscountBadge и BlackFridayBadge
2. **День 2**: Интегрировать в существующие карточки товаров
3. **День 3**: Реализовать модалку истории цен с графиком
4. **День 4**: Добавить фильтрацию и сортировку по скидкам
5. **День 5**: Внедрить аналитику и A/B тестирование

## 📝 Чек-лист готовности

- [ ] DiscountBadge отображается на всех карточках
- [ ] График истории цен работает и показывает манипуляции
- [ ] Black Friday badge автоматически появляется у витрин
- [ ] Фильтр "Только со скидками" работает
- [ ] Метрики отслеживаются в Google Analytics
- [ ] Нет проблем с производительностью при фильтрации
- [ ] SEO разметка добавлена для товаров со скидками