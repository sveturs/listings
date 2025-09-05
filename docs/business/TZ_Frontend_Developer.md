# Техническое задание для Frontend-разработчика (верстальщика) платформы Sve Tu

## Общее описание проекта

**Sve Tu** — современная маркетплейс-платформа, построенная на передовых технологиях. Ваша задача — реализовать pixel-perfect верстку на основе дизайн-макетов, обеспечив высокую производительность, доступность и кроссбраузерность.

## Технологический стек

### Обязательные технологии
- **Next.js 15.3.2** (App Router)
- **React 19** с TypeScript
- **Tailwind CSS v4**
- **DaisyUI** (компонентная библиотека)
- **next-intl** для интернационализации

### Дополнительные библиотеки
- **Redux Toolkit** для state management
- **React Hook Form** для форм
- **Framer Motion** для анимаций
- **React Query** для работы с API
- **date-fns** для работы с датами

## Структура проекта

```
frontend/svetu/
├── src/
│   ├── app/[locale]/          # App Router страницы
│   ├── components/            # Переиспользуемые компоненты
│   ├── hooks/                 # Custom React hooks
│   ├── store/                 # Redux store и slices
│   ├── services/              # API сервисы
│   ├── types/                 # TypeScript типы
│   ├── utils/                 # Утилиты
│   └── messages/              # Файлы локализации
```

## Основные задачи

### 1. Компонентная система

Реализовать библиотеку переиспользуемых компонентов:

#### Базовые компоненты
```typescript
// Button с вариантами
<Button variant="primary|secondary|ghost" size="sm|md|lg" loading={boolean}>
  
// Input с валидацией
<Input 
  type="text|email|password|number" 
  error={string}
  icon={ReactNode}
/>

// Select с поиском
<SearchableSelect 
  options={Array} 
  multiple={boolean}
  async={boolean}
/>

// Card для товаров
<ProductCard 
  product={Product}
  variant="grid|list|compact"
  onQuickView={() => {}}
/>
```

#### Сложные компоненты

**ImageGallery**
- Поддержка zoom при наведении
- Свайпы на мобильных
- Lazy loading
- Оптимизация через Next.js Image

```typescript
<ImageGallery 
  images={string[]}
  showThumbnails={boolean}
  enableZoom={boolean}
/>
```

**SmartFilters**
- Динамические фильтры на основе категории
- Сохранение состояния в URL
- Адаптивный layout
- Счетчики результатов

```typescript
<SmartFilters
  category={Category}
  onFilterChange={(filters) => {}}
  resultsCount={number}
/>
```

**LocationPicker**
- Интеграция с картами
- Автокомплит адресов
- Настройка приватности
- Геолокация пользователя

```typescript
<LocationPicker
  value={Location}
  privacyLevel="exact|approximate|city_only|hidden"
  onChange={(location) => {}}
/>
```

### 2. Страницы для верстки

#### Главная страница (`/[locale]/page.tsx`)
- Hero секция с анимированным поиском
- Bento Grid с категориями и статистикой
- Интерактивная карта
- Секция рекомендаций
- Infinite scroll для товаров

**Особенности:**
- SSR для SEO
- Skeleton loading
- Оптимизация первой загрузки

#### Поиск и каталог (`/[locale]/search`)
- Адаптивная сетка товаров
- Переключение видов (grid/list/map)
- Sticky фильтры на desktop
- Bottom sheet фильтры на mobile
- Virtual scrolling для больших списков

#### Карточка товара (`/[locale]/marketplace/[id]`)
- Структурированные данные (JSON-LD)
- Progressive enhancement галереи
- Динамический import для тяжелых компонентов
- Оптимистичные обновления

#### Создание объявления

**AI Mode** (`/[locale]/create-listing-ai`)
```typescript
// Drag & drop с preview
<ImageUploader
  onImagesAnalyzed={(analysis) => {}}
  maxFiles={10}
  maxSize={10MB}
/>

// Стриминг ответов от AI
<AIStreamingResponse
  prompt={string}
  onComplete={(result) => {}}
/>
```

**Smart Mode** (`/[locale]/create-listing-smart`)
- Multi-step wizard
- Прогресс сохранения
- Валидация на каждом шаге
- Автосохранение черновиков

#### Dashboard витрины (`/[locale]/storefronts/[slug]/dashboard`)
- Responsive таблицы
- Интерактивные графики (Chart.js)
- Real-time уведомления (WebSocket)
- Drag & drop для управления товарами

#### Чат (`/[locale]/chat`)
- WebSocket подключение
- Оптимистичная отправка сообщений
- Виртуализация списка сообщений
- Анимированные эмодзи
- Push-уведомления

### 3. Оптимизация производительности

#### Изображения
```typescript
// Используйте Next.js Image везде
import Image from 'next/image'

<Image
  src={url}
  alt={description}
  width={300}
  height={200}
  placeholder="blur"
  blurDataURL={base64}
  sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
/>
```

#### Code Splitting
```typescript
// Динамический импорт тяжелых компонентов
const MapComponent = dynamic(() => import('@/components/Map'), {
  ssr: false,
  loading: () => <MapSkeleton />
})
```

#### Оптимизация bundle
- Tree shaking неиспользуемого кода
- Минификация CSS
- Компрессия ассетов
- Правильные chunk стратегии

### 4. Адаптивность и доступность

#### Breakpoints
```css
/* Tailwind CSS v4 breakpoints */
sm: 640px   /* Mobile landscape */
md: 768px   /* Tablet */
lg: 1024px  /* Desktop */
xl: 1280px  /* Wide desktop */
2xl: 1536px /* Extra wide */
```

#### Mobile-specific
- Touch targets минимум 44x44px
- Свайп-жесты для навигации
- Bottom navigation bar
- Pull-to-refresh
- Оптимизация для slow 3G

#### Accessibility
```typescript
// ARIA атрибуты
<button
  aria-label="Добавить в избранное"
  aria-pressed={isFavorite}
  role="button"
>

// Keyboard navigation
const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' || e.key === ' ') {
    // Handle action
  }
}

// Focus management
const trapFocus = (element: HTMLElement) => {
  // Implement focus trap for modals
}
```

### 5. Интернационализация

#### Настройка локалей
```typescript
// Поддержка 3 языков
const locales = ['en', 'ru', 'sr'] as const

// Использование next-intl
import { useTranslations } from 'next-intl'

const t = useTranslations('marketplace')
<h1>{t('title')}</h1>
```

#### RTL поддержка
- Подготовка layouts для RTL языков
- Правильные направления для свайпов
- Зеркальные иконки где необходимо

### 6. Формы и валидация

```typescript
// React Hook Form с Zod
import { useForm } from 'react-hook-form'
import { z } from 'zod'

const schema = z.object({
  title: z.string().min(5).max(100),
  price: z.number().positive(),
  category: z.string().uuid()
})

const { register, handleSubmit, formState: { errors } } = useForm({
  resolver: zodResolver(schema)
})
```

### 7. State Management

#### Redux Toolkit setup
```typescript
// Slices для основных модулей
- authSlice (user, session, tokens)
- marketplaceSlice (listings, filters, search)
- storefrontSlice (products, orders, analytics)
- chatSlice (messages, conversations, online users)
- cartSlice (items, totals, shipping)
```

#### Оптимистичные обновления
```typescript
// Пример для избранного
const toggleFavorite = async (listingId: string) => {
  // Оптимистично обновляем UI
  dispatch(updateFavorite({ listingId, isFavorite: !current }))
  
  try {
    await api.toggleFavorite(listingId)
  } catch (error) {
    // Откатываем если ошибка
    dispatch(updateFavorite({ listingId, isFavorite: current }))
  }
}
```

### 8. SEO оптимизация

#### Metadata API
```typescript
export async function generateMetadata({ params }): Promise<Metadata> {
  const product = await getProduct(params.id)
  
  return {
    title: product.title,
    description: product.description,
    openGraph: {
      images: [product.images[0]],
    },
  }
}
```

#### Структурированные данные
```typescript
<script
  type="application/ld+json"
  dangerouslySetInnerHTML={{
    __html: JSON.stringify({
      '@context': 'https://schema.org',
      '@type': 'Product',
      name: product.title,
      image: product.images,
      offers: {
        '@type': 'Offer',
        price: product.price,
        priceCurrency: product.currency,
      },
    }),
  }}
/>
```

### 9. Тестирование

#### Unit тесты компонентов
```typescript
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

test('ProductCard displays correctly', async () => {
  render(<ProductCard product={mockProduct} />)
  
  expect(screen.getByText(mockProduct.title)).toBeInTheDocument()
  expect(screen.getByText(`${mockProduct.price} RSD`)).toBeInTheDocument()
})
```

#### E2E тесты критических путей
- Создание объявления
- Покупка товара
- Регистрация витрины
- Отправка сообщения

### 10. Особые требования

#### PWA функциональность
```typescript
// next.config.js
const withPWA = require('next-pwa')({
  dest: 'public',
  register: true,
  skipWaiting: true,
})
```

#### Офлайн поддержка
- Service Worker для кеширования
- IndexedDB для черновиков
- Синхронизация при восстановлении связи

#### Push-уведомления
```typescript
// Запрос разрешения
const permission = await Notification.requestPermission()

// Подписка на уведомления
if (permission === 'granted') {
  const subscription = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: urlB64ToUint8Array(publicVapidKey)
  })
}
```

## Требования к коду

### Код стайл
- ESLint + Prettier настроены в проекте
- Использовать TypeScript strict mode
- Документировать сложные компоненты
- Meaningful имена переменных и функций

### Git workflow
- Conventional commits
- Отдельная ветка для каждой фичи
- Code review обязателен
- Тесты должны проходить перед merge

### Performance metrics
- Lighthouse score > 90
- First Contentful Paint < 1.5s
- Time to Interactive < 3.5s
- Cumulative Layout Shift < 0.1

## Дедлайны и приоритеты

### Фаза 1 (1 неделя)
- Базовые компоненты
- Главная страница
- Поиск и каталог
- Mobile адаптация

### Фаза 2 (1 неделя)
- Карточка товара
- Создание объявления
- Личный кабинет
- Чат

### Фаза 3 (1 неделя)
- Dashboard витрины
- Checkout процесс
- PWA функциональность
- Оптимизация

### Фаза 4 (3 дня)
- Тестирование
- Баг-фиксы
- Документация
- Performance tuning

## Критерии приемки

1. **Pixel-perfect** верстка согласно макетам
2. **Responsive** на всех устройствах
3. **Кроссбраузерность** (Chrome, Firefox, Safari, Edge)
4. **Performance** метрики в зеленой зоне
5. **Accessibility** WCAG 2.1 Level AA
6. **Код покрыт** тестами минимум на 70%
7. **Документация** для всех компонентов

## Ресурсы

- Figma макеты: [будет предоставлено]
- API документация: `/backend/docs/swagger.json`
- Storybook компонентов: [будет создан]
- Тестовые данные: `/frontend/svetu/src/utils/mockDataGenerator.ts`

---

*При возникновении вопросов обращайтесь к техническому руководителю. Успешной разработки!*