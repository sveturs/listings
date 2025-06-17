# План внедрения системы отзывов на фронтенде

## Дата: 17.06.2025
## Анализ требований

### Система агрегации оценок:
1. **Отзывы на товары** → оценка попадает продавцу/магазину
2. **Отзывы на магазин** → оценка попадает владельцу магазина  
3. **Отзывы на продавца** → личная оценка продавца
4. **Сохранение оценок** → даже после удаления товара оценка остается у владельца

## Архитектура решения

### 1. Структура данных для отзывов

```typescript
// Типы отзывов
enum ReviewEntityType {
  LISTING = 'listing',      // отзыв на товар
  STOREFRONT = 'storefront', // отзыв на магазин
  USER = 'user'             // отзыв на продавца
}

// Расширенная модель отзыва
interface Review {
  id: number;
  user_id: number;
  entity_type: ReviewEntityType;
  entity_id: number;
  
  // Для сохранения связи после удаления
  entity_origin_type: 'user' | 'storefront';
  entity_origin_id: number;
  
  rating: number;
  comment: string;
  photos: string[];
  is_verified_purchase: boolean;
  seller_confirmed: boolean;
  // ...
}

// Агрегированный рейтинг
interface AggregatedRating {
  average: number;
  total_reviews: number;
  distribution: Record<1|2|3|4|5, number>;
  verified_reviews: number;
  listing_reviews: number;    // отзывы через товары
  direct_reviews: number;      // прямые отзывы
  recent_trend: 'up' | 'down' | 'stable';
}
```

### 2. Компоненты для реализации

#### 2.1 Базовые компоненты отзывов
```
src/components/reviews/
├── ReviewsList.tsx          # Список отзывов с фильтрами
├── ReviewItem.tsx           # Отдельный отзыв
├── ReviewForm.tsx           # Форма добавления отзыва
├── ReviewStats.tsx          # Статистика и графики
├── RatingDisplay.tsx        # Отображение рейтинга (звезды)
├── RatingInput.tsx          # Ввод рейтинга
├── ReviewFilters.tsx        # Фильтры отзывов
├── ReviewPhotos.tsx         # Галерея фото отзыва
└── VerifiedBadge.tsx        # Бейдж верификации
```

#### 2.2 Специализированные компоненты
```
src/components/reviews/
├── listing/
│   └── ListingReviews.tsx  # Отзывы на товар
├── storefront/
│   └── StorefrontReviews.tsx # Отзывы на магазин
└── seller/
    └── SellerReviews.tsx    # Отзывы на продавца
```

### 3. Страницы для добавления

#### 3.1 Страница продавца/магазина
```
src/app/[locale]/seller/[id]/page.tsx
src/app/[locale]/storefront/[id]/page.tsx
```

#### 3.2 Модальные окна
```
src/components/modals/
├── ReviewModal.tsx          # Модалка для написания отзыва
├── ReviewPhotoUpload.tsx    # Загрузка фото
└── ReviewDispute.tsx        # Оспаривание отзыва
```

### 4. Интеграция в существующие страницы

#### 4.1 Страница товара (`marketplace/[id]/page.tsx`)
- Добавить секцию отзывов под описанием
- Показывать рейтинг товара И продавца
- Кнопка "Оставить отзыв" (если был чат)

#### 4.2 Компонент SellerInfo
- Заменить захардкоженный рейтинг на реальный
- Добавить ссылку на все отзывы продавца
- Показывать распределение оценок

#### 4.3 Карточка товара в списке
- Добавить отображение рейтинга
- Количество отзывов

### 5. Redux Store структура

```typescript
// store/slices/reviewsSlice.ts
interface ReviewsState {
  // Списки отзывов по типу и ID
  byEntity: {
    [entityType: string]: {
      [entityId: string]: {
        reviews: Review[];
        stats: ReviewStats;
        loading: boolean;
        hasMore: boolean;
        filters: ReviewFilters;
      }
    }
  };
  
  // Агрегированные рейтинги
  ratings: {
    [entityType: string]: {
      [entityId: string]: AggregatedRating;
    }
  };
  
  // Управление формой
  form: {
    isOpen: boolean;
    entityType: ReviewEntityType;
    entityId: number;
    draft: Partial<Review>;
  };
  
  // Можно ли оставить отзыв
  canReview: {
    [key: string]: boolean; // "type:id" -> true/false
  };
}
```

### 6. API интеграция

#### 6.1 Новые хуки
```typescript
// hooks/useReviews.ts
- useReviewsList(entityType, entityId, filters)
- useReviewStats(entityType, entityId)
- useCanLeaveReview(entityType, entityId)
- useCreateReview()
- useUpdateReview()
- useReviewVote()
```

#### 6.2 API клиент
```typescript
// services/api/reviews.ts
- getReviews(params)
- getReviewStats(entityType, entityId)
- createReview(data)
- updateReview(id, data)
- deleteReview(id)
- voteReview(id, voteType)
- uploadReviewPhotos(id, files)
- checkCanReview(entityType, entityId)
```

### 7. Этапы реализации

#### Этап 1: Базовая инфраструктура (3-4 дня)
1. Создать типы данных для отзывов
2. Добавить Redux slice для отзывов
3. Создать API сервис и базовые хуки
4. Создать базовые компоненты (RatingDisplay, RatingInput)

#### Этап 2: Отзывы на товары (3-4 дня)
1. Интегрировать в страницу товара
2. Создать ReviewsList и ReviewItem
3. Создать форму добавления отзыва
4. Добавить загрузку фото

#### Этап 3: Страницы продавцов/магазинов (4-5 дней)
1. Создать страницу продавца с отзывами
2. Создать страницу магазина
3. Реализовать агрегацию рейтингов
4. Добавить фильтры и сортировку

#### Этап 4: Расширенный функционал (3-4 дня)
1. Система подтверждения продавцом
2. Оспаривание отзывов
3. Модерация и жалобы
4. Email/push уведомления

#### Этап 5: Оптимизация и UX (2-3 дня)
1. Оптимизация загрузки (lazy loading)
2. Анимации и переходы
3. Мобильная адаптация
4. A/B тестирование

### 8. Технические решения

#### 8.1 Кеширование рейтингов
- Использовать React Query для кеширования
- Обновлять рейтинги через WebSocket
- Prefetch при наведении на карточку

#### 8.2 Оптимизация изображений
- Lazy loading для фото отзывов
- Превью в base64 для быстрой загрузки
- Сжатие на клиенте перед загрузкой

#### 8.3 SEO оптимизация
- Structured data для отзывов (schema.org)
- Серверный рендеринг рейтингов
- Rich snippets для поисковиков

### 9. UI/UX рекомендации

#### 9.1 Отображение рейтинга
- Всегда показывать количество отзывов
- Выделять верифицированные покупки
- График распределения оценок
- Фильтр "только с фото"

#### 9.2 Форма отзыва
- Пошаговый wizard для удобства
- Подсказки что писать
- Автосохранение черновика
- Превью перед отправкой

#### 9.3 Доверие и прозрачность
- Бейджи верификации
- История изменений отзыва
- Ответы продавца выделены
- Причины удаления отзывов

### 10. Метрики успеха

1. **Вовлеченность**
   - % пользователей оставивших отзыв
   - Среднее количество отзывов на товар
   - Время до первого отзыва

2. **Качество**
   - Средняя длина отзыва
   - % отзывов с фото
   - % верифицированных отзывов

3. **Доверие**
   - Конверсия после просмотра отзывов
   - Снижение возвратов
   - Рост повторных покупок

### 11. Риски и митигация

1. **Накрутка отзывов**
   - Алгоритмы определения подозрительной активности
   - Ограничения по IP/устройству
   - Обязательная верификация через чат

2. **Негативные отзывы**
   - Система быстрого реагирования для продавцов
   - Возможность публичного ответа
   - Медиация споров

3. **Производительность**
   - Пагинация и виртуализация списков
   - CDN для изображений отзывов
   - Фоновая загрузка статистики