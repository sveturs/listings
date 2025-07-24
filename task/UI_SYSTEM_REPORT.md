# Отчет о UI системе проекта Sve Tu Platform

## Обзор архитектуры UI

### Технологический стек
- **React 19** с TypeScript
- **Next.js 15.3.2** с App Router
- **Tailwind CSS v4** + **DaisyUI v5.0.42**
- **Redux Toolkit** для управления состоянием
- **next-intl** для интернационализации (en/ru)

### Структура компонентов

#### Основные категории компонентов:

1. **Базовые компоненты**
   - `FormField` - обертка для полей формы с лейблом и обработкой ошибок
   - `SafeImage` / `OptimizedImage` - компоненты для безопасной загрузки изображений
   - `ErrorBoundary` - обработка ошибок React
   - `ViewToggle` - переключатель режима отображения (grid/list)
   - `InfiniteScrollTrigger` - триггер бесконечной прокрутки

2. **Навигация и макет**
   - `Header` - главная навигация с поддержкой корзины, авторизации
   - `SearchBar` - поисковая строка с автодополнением
   - `LanguageSwitcher` - переключатель языков
   - `CategorySidebar` - боковая панель категорий

3. **Аутентификация**
   - `AuthButton` - кнопка входа/профиля
   - `LoginModal` - модальное окно авторизации
   - `LoginForm` / `RegisterForm` - формы входа и регистрации
   - `AuthStateManager` - управление состоянием авторизации

4. **Marketplace компоненты**
   - `MarketplaceCard` - карточка товара (поддерживает grid/list режимы)
   - `MarketplaceList` - список товаров с фильтрами
   - `MarketplaceFilters` / `ListingFilters` - система фильтров
   - `QuickFilters` / `SmartFilters` - быстрые и умные фильтры

5. **Карты и геолокация**
   - `InteractiveMap` - интерактивная карта на Mapbox
   - `LocationPicker` - выбор местоположения
   - `DistrictSelector` - выбор района
   - `RadiusSearchControl` - контроль радиуса поиска
   - `WalkingAccessibilityControl` - контроль пешей доступности
   - Кластеризация маркеров (`MapboxClusterLayer`, `PriceCluster`)

6. **Галереи и медиа**
   - `ImageGallery` (два варианта: `ui/ImageGallery` и `reviews/ImageGallery`)
   - Поддержка полноэкранного режима
   - Жесты свайпа на мобильных
   - Навигация клавиатурой и колесом мыши

7. **Чат система**
   - `ChatLayout` - макет чата
   - `ChatWindow` - окно чата
   - `MessageItem` - отдельное сообщение
   - `EmojiPicker` - выбор эмодзи
   - `AnimatedEmoji` - анимированные эмодзи
   - `ChatAttachments` - вложения в чате

8. **Витрины (Storefronts)**
   - `StorefrontHeader` - шапка витрины
   - `StorefrontProducts` - список товаров витрины
   - `StorefrontMap` - карта витрин
   - `ProductCard` - карточка товара витрины
   - Мастер создания витрины (multi-step wizard)

9. **Корзина и оплата**
   - `CartIcon` - иконка корзины с счетчиком
   - `ShoppingCartModal` - модальное окно корзины
   - `PaymentButton` / `PaymentMethodSelector` - оплата
   - `EscrowStatus` - статус эскроу

10. **Отзывы**
    - `ReviewsSection` - секция отзывов
    - `ReviewList` - список отзывов
    - `ReviewForm` - форма добавления отзыва
    - `RatingDisplay` / `RatingInput` - отображение и ввод рейтинга
    - `RatingStats` - статистика рейтингов

11. **Атрибуты и фильтры**
    - `AttributeFilters` - фильтры по атрибутам
    - `MultiSelectAttribute` - множественный выбор
    - `RangeAttribute` - диапазонные фильтры
    - `TranslationStatus` - статус переводов
    - `BatchTranslationModal` - массовый перевод

## Система стилей

### DaisyUI интеграция
- Используется DaisyUI v5 как основная UI библиотека
- Применяются встроенные классы DaisyUI:
  - `btn`, `btn-primary`, `btn-ghost` - кнопки
  - `card`, `card-body`, `card-title` - карточки
  - `form-control`, `label`, `input` - формы
  - `modal`, `drawer` - модальные окна
  - `navbar`, `dropdown` - навигация
  - `badge`, `alert`, `toast` - уведомления

### Кастомные стили
- Минимальное количество кастомных CSS
- Основные кастомные стили в:
  - `globals.css` - глобальные переменные и темы
  - `chat-patterns.css` - стили для чата
  - `chat-bubble.css` - пузыри сообщений
  - `markdown.css` - стили для markdown
  - `map-popup.css` - стили попапов карты

### Темы
- Поддержка светлой/темной темы через CSS переменные
- Переменные:
  - `--background` / `--foreground`
  - Автоматическое переключение по `prefers-color-scheme`

## Управление состоянием (Redux)

### Store структура:
```typescript
{
  chat: chatSlice,
  reviews: reviewsSlice,
  storefronts: storefrontSlice,
  import: importSlice,
  products: productSlice,
  payment: paymentSlice,
  cart: cartSlice
}
```

### Middleware:
- `websocketMiddleware` - управление WebSocket соединениями

## Хуки

### Основные кастомные хуки:
1. **Геолокация и карты**
   - `useGeolocation` - получение текущей геопозиции
   - `useGeoSearch` - поиск по геоданным
   - `useDistanceCalculation` - расчет расстояний
   - `useRadiusSearch` - поиск в радиусе
   - `useVisibleCities` - видимые города на карте
   - `useAddressGeocoding` - геокодирование адресов

2. **UI и взаимодействие**
   - `useDebounce` - дебаунс значений
   - `useInfiniteScroll` - бесконечная прокрутка
   - `useViewPreference` - предпочтения отображения
   - `useLocalStorage` - работа с localStorage
   - `usePaginatedData` - пагинация данных

3. **Бизнес-логика**
   - `useChat` - функционал чата
   - `useReviews` - работа с отзывами
   - `useStorefronts` - управление витринами
   - `useBalance` - баланс пользователя
   - `useAllSecurePayment` - безопасные платежи
   - `useBehaviorTracking` - отслеживание поведения
   - `useAnalytics` - аналитика

4. **Формы и валидация**
   - `useFormValidation` - валидация форм
   - `useAuthForm` - формы авторизации
   - `useListingDraft` - черновики объявлений
   - `useCategoryFilters` - фильтры категорий

## Мобильная адаптация

### Компоненты для мобильных:
- `MobileFiltersDrawer` - выдвижная панель фильтров
- Поддержка жестов в `ImageGallery`
- Адаптивная навигация в `Header`
- Оптимизированные карточки товаров для мобильных

### Breakpoints:
- Используются стандартные breakpoints Tailwind CSS
- Компоненты адаптируются через responsive классы

## Интернационализация

### Поддерживаемые языки:
- Английский (en)
- Русский (ru)

### Файлы переводов:
- `/messages/en.json`
- `/messages/ru.json`

### Использование:
```typescript
const t = useTranslations('header');
// t('nav.blog')
```

## Оптимизация производительности

1. **Ленивая загрузка изображений**
   - `SafeImage` с fallback
   - Next.js Image оптимизация

2. **Виртуализация списков**
   - `VirtualizedList` компонент
   - `AttributeListVirtualized` для больших списков

3. **Мемоизация и оптимизация рендеров**
   - Использование React.memo
   - useMemo/useCallback хуки

4. **Code splitting**
   - Динамический импорт для тяжелых компонентов
   - Разделение по маршрутам

## Рекомендации по улучшению

1. **Унификация компонентов**
   - Два варианта `ImageGallery` можно объединить
   - Стандартизировать обработку ошибок

2. **Типизация**
   - Добавить более строгую типизацию для Redux actions
   - Улучшить типы для API responses

3. **Тестирование**
   - Расширить покрытие unit тестами
   - Добавить интеграционные тесты для критических flow

4. **Документация**
   - Создать Storybook для компонентов
   - Документировать паттерны использования

5. **Производительность**
   - Оптимизировать bundle size
   - Реализовать более агрессивное кеширование

## Заключение

UI система проекта Sve Tu хорошо структурирована и использует современные подходы. Основные сильные стороны:
- Четкая компонентная архитектура
- Эффективное использование DaisyUI
- Хорошая поддержка мобильных устройств
- Продуманная система управления состоянием

Области для улучшения в основном касаются унификации похожих компонентов и расширения тестового покрытия.