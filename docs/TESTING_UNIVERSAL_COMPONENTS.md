# 🧪 Инструкция по тестированию универсальных компонентов

## 📋 Предварительная подготовка

### 1. Применение миграций БД
```bash
cd /data/hostel-booking-system/backend

# Проверить статус миграций
./migrator status

# Применить новую миграцию для универсальной истории просмотров
./migrator up

# Проверить, что миграция 000020 применена
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable" \
  -c "SELECT * FROM schema_migrations WHERE version = 20;"

# Проверить создание новых таблиц
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable" \
  -c "\dt user_view_history"

psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable" \
  -c "\dt view_statistics"
```

### 2. Проверка файловой структуры
```bash
# Проверить созданные компоненты
ls -la /data/hostel-booking-system/frontend/svetu/src/components/universal/
# Должны быть папки: cards/, filters/, calculators/, recommendations/

# Проверить Redux slices
ls -la /data/hostel-booking-system/frontend/svetu/src/store/slices/universalCompareSlice.ts

# Проверить миграции
ls -la /data/hostel-booking-system/backend/migrations/000020_*
```

## 🎨 Тестирование Frontend компонентов

### 1. Создание тестовой страницы

Создайте временную страницу для тестирования всех компонентов:

```typescript
// frontend/svetu/src/app/[locale]/test-universal/page.tsx
'use client';

import { useState } from 'react';
import UniversalListingCard from '@/components/universal/cards/UniversalListingCard';
import UniversalFilters from '@/components/universal/filters/UniversalFilters';
import UniversalCreditCalculator from '@/components/universal/calculators/UniversalCreditCalculator';
import RecommendationsEngine from '@/components/universal/recommendations/RecommendationsEngine';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { initializeCompare } from '@/store/slices/universalCompareSlice';

export default function TestUniversalPage() {
  const [filters, setFilters] = useState({});
  const [selectedCategory] = useState('cars');

  // Тестовые данные для карточки
  const testListing = {
    id: 1,
    title: 'Toyota Camry 2020',
    price: 25000,
    currency: '€',
    images: ['https://picsum.photos/400/300'],
    location: { city: 'Belgrade' },
    category: 'cars',
    createdAt: new Date().toISOString(),
    customFields: [
      { label: 'Year', value: '2020' },
      { label: 'Mileage', value: '45,000 km' },
      { label: 'Fuel', value: 'Gasoline' },
    ],
    badges: [
      { type: 'new', label: 'New' },
      { type: 'discount', label: '-15%', value: '15' }
    ],
    stats: {
      views: 234,
      favorites: 12,
    }
  };

  return (
    <div className="container mx-auto p-4 space-y-8">
      <h1 className="text-3xl font-bold mb-8">Тестирование универсальных компонентов</h1>

      {/* 1. UniversalListingCard */}
      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">1. UniversalListingCard</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <h3 className="text-lg mb-2">Grid Layout</h3>
            <UniversalListingCard
              data={testListing}
              type="cars"
              layout="grid"
              showBadges={true}
              showFavorite={true}
              showCompare={true}
              showStats={true}
            />
          </div>

          <div>
            <h3 className="text-lg mb-2">List Layout</h3>
            <UniversalListingCard
              data={testListing}
              type="cars"
              layout="list"
              showBadges={true}
              showFavorite={true}
              showCompare={true}
              showStats={true}
            />
          </div>
        </div>
      </section>

      {/* 2. UniversalFilters */}
      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">2. UniversalFilters</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <h3 className="text-lg mb-2">Vertical Layout</h3>
            <div className="bg-base-200 p-4 rounded">
              <UniversalFilters
                category={selectedCategory}
                filters={filters}
                onFiltersChange={setFilters}
                layout="vertical"
              />
            </div>
          </div>

          <div>
            <h3 className="text-lg mb-2">Current Filters</h3>
            <pre className="bg-base-300 p-4 rounded overflow-auto">
              {JSON.stringify(filters, null, 2)}
            </pre>
          </div>
        </div>
      </section>

      {/* 3. UniversalCreditCalculator */}
      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">3. UniversalCreditCalculator</h2>

        <div className="max-w-2xl">
          <UniversalCreditCalculator
            price={25000}
            category="cars"
            onApply={(calculation) => {
              console.log('Credit calculation:', calculation);
              alert(`Monthly payment: €${calculation.monthlyPayment.toFixed(2)}`);
            }}
          />
        </div>
      </section>

      {/* 4. RecommendationsEngine */}
      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">4. RecommendationsEngine</h2>

        <RecommendationsEngine
          type="similar"
          category="cars"
          currentItemId={1}
          limit={4}
          layout="grid"
          showTitle={true}
          showDescription={true}
        />

        <RecommendationsEngine
          type="trending"
          category="cars"
          limit={3}
          layout="carousel"
          showTitle={true}
        />
      </section>
    </div>
  );
}
```

### 2. Обновление Redux Store

Убедитесь, что новые slices добавлены в store:

```typescript
// frontend/svetu/src/store/index.ts
import universalCompareReducer from './slices/universalCompareSlice';

export const store = configureStore({
  reducer: {
    // ... другие reducers
    universalCompare: universalCompareReducer,
  },
});
```

### 3. Запуск и проверка

```bash
# Запустить frontend
cd /data/hostel-booking-system/frontend/svetu
yarn dev -p 3001

# Открыть в браузере
# http://localhost:3001/ru/test-universal
```

## ✅ Чек-лист проверки функционала

### UniversalListingCard
- [ ] Отображается корректно в grid режиме
- [ ] Отображается корректно в list режиме
- [ ] Работает кнопка добавления в избранное
- [ ] Работает кнопка добавления в сравнение
- [ ] Отображаются бейджи (new, discount)
- [ ] Показывается статистика (просмотры, избранное)
- [ ] Корректно форматируется цена
- [ ] Правильно отображается время публикации

### UniversalFilters
- [ ] Отображаются фильтры для выбранной категории
- [ ] Работают select фильтры
- [ ] Работают multiselect фильтры
- [ ] Работают range фильтры (слайдеры)
- [ ] Работает фильтр цены
- [ ] Сворачиваются/разворачиваются группы фильтров
- [ ] Показывается счетчик активных фильтров
- [ ] Работает кнопка очистки фильтров

### UniversalCreditCalculator
- [ ] Работает слайдер первоначального взноса
- [ ] Работает слайдер срока кредита
- [ ] Работает слайдер процентной ставки
- [ ] Корректно рассчитывается ежемесячный платеж
- [ ] Показывается общая сумма выплат
- [ ] Показывается сумма переплаты
- [ ] Работает выбор банка
- [ ] Отображаются дополнительные расходы
- [ ] Работает график платежей

### UniversalCompareSlice
- [ ] Добавляются элементы в сравнение
- [ ] Удаляются элементы из сравнения
- [ ] Соблюдается лимит элементов (3 для авто)
- [ ] Сохраняется в localStorage
- [ ] Восстанавливается после перезагрузки

### RecommendationsEngine
- [ ] Отображаются рекомендации в grid режиме
- [ ] Работает carousel режим
- [ ] Работает list режим
- [ ] Показываются разные типы рекомендаций
- [ ] Работает кнопка обновления рекомендаций

## 🗄️ Проверка базы данных

### Проверка таблиц
```sql
-- Подключение к БД
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"

-- Проверить структуру user_view_history
\d user_view_history

-- Проверить структуру view_statistics
\d view_statistics

-- Проверить, мигрированы ли данные из старой таблицы
SELECT COUNT(*) FROM user_view_history;

-- Проверить работу триггера
INSERT INTO user_view_history (user_id, listing_id, category_id, viewed_at)
VALUES (1, 1, 1, NOW());

SELECT is_return_visit FROM user_view_history ORDER BY id DESC LIMIT 1;

-- Проверить функцию агрегации
SELECT update_view_statistics(CURRENT_DATE);
SELECT * FROM view_statistics WHERE date = CURRENT_DATE;
```

## 🐛 Типичные проблемы и решения

### 1. Ошибка импорта компонентов
```bash
# Убедитесь, что пути правильные
find /data/hostel-booking-system/frontend/svetu/src/components -name "Universal*.tsx"
```

### 2. Ошибка в Redux store
```bash
# Проверьте, что slice добавлен в store
grep -r "universalCompare" /data/hostel-booking-system/frontend/svetu/src/store/
```

### 3. Ошибка миграции БД
```bash
# Откатить миграцию и применить заново
cd /data/hostel-booking-system/backend
./migrator down
./migrator up
```

### 4. TypeScript ошибки
```bash
# Проверить типы
cd /data/hostel-booking-system/frontend/svetu
yarn tsc --noEmit
```

## 📊 Проверка производительности

```bash
# 1. Открыть Chrome DevTools
# 2. Перейти на вкладку Performance
# 3. Записать загрузку страницы с компонентами
# 4. Проверить:
#    - Время рендера < 100ms
#    - Размер bundle < 200KB
#    - Нет утечек памяти
```

## 🚀 Интеграция в продакшен

После успешного тестирования:

1. **Удалить тестовую страницу**
```bash
rm /data/hostel-booking-system/frontend/svetu/src/app/[locale]/test-universal/page.tsx
```

2. **Интегрировать в реальные страницы**
- Заменить `CarListingCardEnhanced` на `UniversalListingCard`
- Заменить `CarFilters` на `UniversalFilters`
- Добавить `RecommendationsEngine` на страницы

3. **Создать API endpoints**
```bash
# Backend endpoints для:
- /api/v1/recommendations
- /api/v1/view-history
- /api/v1/credit/calculate
```

4. **Запустить полное тестирование**
```bash
cd /data/hostel-booking-system/frontend/svetu
yarn test
yarn build
```

---

*Документ создан: 27.09.2025*
*Для проверки универсальных компонентов маркетплейса*