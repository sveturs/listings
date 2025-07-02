# Паспорта базовых компонентов Frontend Sve Tu Platform

## 1. RootLayout (Main Layout)

**Файл:** `/frontend/svetu/src/app/[locale]/layout.tsx`

### Назначение
Корневой макет приложения, обеспечивающий базовую структуру, интернационализацию и инициализацию провайдеров.

### Основные компоненты и props
```typescript
interface Props {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}
```

### Зависимости
- **Шрифты:** Geist, Geist_Mono от Google Fonts
- **Интернационализация:** next-intl (NextIntlClientProvider, getMessages, getTranslations)
- **Роутинг:** next/navigation, @/i18n/routing
- **Провайдеры:** ReduxProvider, AuthProvider
- **Компоненты:** Header, AuthStateManager, WebSocketManager
- **Runtime env:** next-runtime-env (PublicEnvScript)

### Ключевая логика работы
1. **Валидация локали:** Проверяет входящую локаль против routing.locales
2. **Генерация метаданных:** Асинхронно генерирует title и description через next-intl
3. **SSG:** generateStaticParams для всех поддерживаемых локалей  
4. **Структура провайдеров:** NextIntlClientProvider > ReduxProvider > AuthProvider

### UI структура
```html
<html lang={locale} data-theme="cupcake">
  <head>
    <PublicEnvScript />
  </head>
  <body className="geist-fonts antialiased">
    <NextIntlClientProvider>
      <ReduxProvider>
        <AuthProvider>
          <AuthStateManager />
          <WebSocketManager />
          <Header />
          <main className="min-h-screen pt-28 lg:pt-16">
            {children}
          </main>
        </AuthProvider>
      </ReduxProvider>
    </NextIntlClientProvider>
  </body>
</html>
```

### Особенности реализации
- **Тема:** Использует DaisyUI тему "cupcake"
- **Отступы:** Адаптивные отступы для header (pt-28 на мобильных, pt-16 на десктопе)
- **SSR/SSG:** Полная поддержка серверного рендеринга
- **Error handling:** notFound() для невалидных локалей

---

## 2. AdminLayout

**Файл:** `/frontend/svetu/src/app/[locale]/admin/layout.tsx`

### Назначение
Layout для административной панели с боковой навигацией и защитой доступа.

### Основные компоненты и props
```typescript
interface Props {
  children: React.ReactNode;
}
```

### Зависимости
- **Интернационализация:** next-intl (useTranslations)
- **Навигация:** @/i18n/routing (Link), next/navigation (usePathname)
- **Защита:** AdminGuard компонент

### Ключевая логика работы
1. **Route detection:** isActive функция для подсветки активных пунктов меню
2. **Responsive navigation:** Drawer pattern для мобильных устройств
3. **Access control:** Обернут в AdminGuard для проверки прав доступа

### UI структура
- **Desktop:** Постоянная боковая панель (drawer-open)
- **Mobile:** Скрытая боковая панель с hamburger menu
- **Navigation sections:**
  - Dashboard
  - Catalog (Categories, Attributes, Attribute Groups)
  - Content (Listings, Users)

### Особенности реализации
- **DaisyUI Drawer:** drawer-mobile lg:drawer-open паттерн
- **SVG Icons:** Встроенные иконки для всех пунктов меню
- **Локализация:** Полная поддержка переводов через admin namespace
- **Active states:** Динамическое выделение активных пунктов

---

## 3. Header

**Файл:** `/frontend/svetu/src/components/Header.tsx`

### Назначение
Главный заголовок приложения с навигацией, поиском и аутентификацией.

### Основные компоненты и props
Не принимает props, использует контексты и хуки.

### Зависимости
- **State:** useState, useEffect для локального состояния
- **Интернационализация:** next-intl (useTranslations)
- **Навигация:** @/i18n/routing (Link), next/navigation (usePathname)
- **Компоненты:** LanguageSwitcher, AuthButton, LoginModal, SearchBar
- **Контекст:** useAuthContext для получения пользователя

### Ключевая логика работы
1. **Responsive behavior:** Адаптивное отображение навигации и поиска
2. **Auth integration:** Интеграция с системой аутентификации
3. **Modal management:** Управление модальным окном входа
4. **Search conditional rendering:** Скрытие мобильного поиска на странице поиска

### UI структура
```
Header (fixed top)
├── Mobile menu (dropdown)
├── Logo (SveTu)
├── Desktop navigation (hidden lg:flex)
├── Search bar (desktop только)
└── Right section
    ├── Create listing button (if authenticated)
    ├── LanguageSwitcher
    └── AuthButton
```

### Особенности реализации
- **Fixed positioning:** Фиксированный заголовок с z-index 100
- **Conditional search:** Отдельная мобильная поисковая строка ниже header
- **Adaptive button:** Кнопка создания объявления только для авторизованных
- **Responsive layout:** Скрытие/показ элементов в зависимости от размера экрана

---

## 4. HomePage (Main Page)

**Файл:** `/frontend/svetu/src/app/[locale]/page.tsx`

### Назначение
Главная страница маркетплейса с отображением товаров и основной функциональностью.

### Основные компоненты и props
```typescript
interface Props {
  params: Promise<{ locale: string }>;
}
```

### Зависимости
- **Services:** MarketplaceService для загрузки данных
- **Компоненты:** MarketplaceList для отображения товаров
- **Config:** configManager для feature flags
- **Интернационализация:** next-intl (getTranslations)

### Ключевая логика работы
1. **SSR conditional:** Отключение SSR в development режиме
2. **Error handling:** Graceful handling ошибок загрузки данных
3. **Feature flags:** Условное отображение функций (payments)
4. **Empty state:** Обработка случая отсутствия товаров

### UI структура
```
HomePage
├── Hero section (gradient background)
├── Feature alerts (if enabled)
├── Error alerts (if error)
├── MarketplaceList (if data available)
├── Empty state (if no data)
└── Floating action button (create listing)
```

### Особенности реализации
- **SSR bypass:** Специальная логика для development окружения
- **Network resilience:** Timeout и error handling для API запросов
- **Feature toggles:** Интеграция с системой feature flags
- **Floating FAB:** Кнопка создания объявления в правом нужнем углу

---

## 5. ReduxProvider

**Файл:** `/frontend/svetu/src/components/ReduxProvider.tsx`

### Назначение
Провайдер состояния приложения, объединяющий Redux и React Query.

### Основные компоненты и props
```typescript
interface Props {
  children: React.ReactNode;
}
```

### Зависимости
- **Redux:** react-redux (Provider)
- **Queries:** @tanstack/react-query (QueryClient, QueryClientProvider)
- **Store:** @/store (Redux store)

### Ключевая логика работы
1. **QueryClient configuration:** Настройка React Query с кэшированием
2. **Provider composition:** Комбинирование Redux и React Query провайдеров
3. **Memoization:** useState для предотвращения пересоздания QueryClient

### Особенности реализации
- **Stale time:** 1 минута для запросов
- **Window focus:** Отключен refetch при фокусе окна
- **Provider nesting:** Redux > React Query для правильного порядка

---

## 6. AuthStateManager

**Файл:** `/frontend/svetu/src/components/AuthStateManager.tsx`

### Назначение
Менеджер состояния аутентификации, очищающий данные при смене пользователя.

### Зависимости
- **Контексты:** useAuth для получения пользователя
- **Hooks:** useChat для очистки чат данных
- **Storage:** localStorage, sessionStorage для управления хранилищем

### Ключевая логика работы
1. **User change detection:** Отслеживание смены пользователя через useRef
2. **Selective cleanup:** Очистка только чат данных, сохранение токенов
3. **Storage management:** Умная очистка sessionStorage и localStorage

### Особенности реализации
- **Invisible component:** Не рендерит UI, только логика
- **Storage filtering:** Сохраняет важные ключи при очистке
- **Error handling:** Try-catch для операций с storage

---

## 7. WebSocketManager

**Файл:** `/frontend/svetu/src/components/WebSocketManager.tsx`

### Назначение
Менеджер WebSocket соединений для реального времени (чаты).

### Зависимости
- **Контексты:** useAuth для статуса аутентификации
- **Hooks:** useChat для управления WebSocket
- **Utils:** tokenManager для работы с токенами

### Ключевая логика работы
1. **Connection lifecycle:** Инициализация при входе, закрытие при выходе
2. **Token verification:** Проверка наличия токена перед подключением
3. **Initialization guard:** Предотвращение множественных инициализаций
4. **Delayed initialization:** Таймаут для ожидания токена

### Особенности реализации
- **Timeout handling:** 500мс задержка для стабильности
- **Cleanup management:** Правильная очистка таймаутов и соединений
- **Invisible component:** Не рендерит UI, только управление WebSocket

---

## 8. SearchBar

**Файл:** `/frontend/svetu/src/components/SearchBar/SearchBar.tsx`

### Назначение
Универсальная поисковая строка с автодополнением и историей поиска.

### Основные компоненты и props
```typescript
interface SearchBarProps {
  placeholder?: string;
  initialQuery?: string;
  onSearch?: (query: string) => void;
  className?: string;
  variant?: 'default' | 'hero' | 'minimal';
  showTrending?: boolean;
}
```

### Зависимости
- **Services:** UnifiedSearchService для поиска и предложений
- **Hooks:** useDebounce для оптимизации запросов
- **Icons:** Собственные иконки из ./icons
- **Навигация:** next/navigation (useRouter)

### Ключевая логика работы
1. **Debounced search:** 300мс задержка для оптимизации API запросов
2. **Keyboard navigation:** Полная поддержка навигации стрелками
3. **Suggestion types:** Различные типы предложений (продукты, категории)
4. **History management:** Сохранение и отображение истории поиска

### UI структура
```
SearchBar
├── Input field (with icons)
├── Suggestions dropdown
│   ├── Search results
│   ├── Search history
│   └── Trending searches
└── Action buttons (clear, search)
```

### Особенности реализации
- **Variant system:** Три варианта оформления для разных контекстов
- **Click outside:** Автоматическое закрытие предложений
- **Route-aware:** Умное поведение на странице поиска
- **Accessibility:** Полная поддержка ARIA и клавиатурной навигации

---

## 9. LanguageSwitcher

**Файл:** `/frontend/svetu/src/components/LanguageSwitcher.tsx`

### Назначение
Переключатель языков интерфейса.

### Зависимости
- **Интернационализация:** next-intl (useLocale)
- **Навигация:** @/i18n/routing для смены локали
- **Transitions:** useTransition для плавного переключения

### Ключевая логика работы
1. **Locale switching:** Смена языка с сохранением текущего маршрута
2. **Transition handling:** Плавное переключение с loading состоянием
3. **Current locale display:** Показ текущего языка в кнопке

### Особенности реализации
- **Dropdown pattern:** DaisyUI dropdown для выбора языка
- **Disabled state:** Блокировка во время перехода
- **Route preservation:** Сохранение текущего пути при смене языка

---

## 10. AuthButton

**Файл:** `/frontend/svetu/src/components/AuthButton.tsx`

### Назначение
Кнопка аутентификации с меню пользователя и индикацией чатов.

### Основные компоненты и props
```typescript
interface AuthButtonProps {
  onLoginClick?: () => void;
}
```

### Зависимости
- **Контексты:** useAuth для состояния аутентификации
- **Hooks:** useChat для счетчика непрочитанных сообщений
- **Image:** next/image для аватаров
- **Интернационализация:** next-intl

### Ключевая логика работы
1. **State management:** Различные состояния (loading, error, authenticated)
2. **Image fallback:** Обработка ошибок загрузки аватаров
3. **Unread indicator:** Пульсирующий индикатор непрочитанных чатов
4. **Menu structure:** Контекстное меню для авторизованных пользователей

### UI структура (authenticated)
```
AuthButton
├── Chat button (with unread indicator)
└── User dropdown
    ├── Avatar (with fallback)
    ├── User info
    ├── Profile links
    ├── Admin links (if admin)
    └── Logout button
```

### Особенности реализации
- **SSR hydration:** Правильная обработка гидратации
- **Error auto-dismiss:** Автоматическое скрытие ошибок через 5 секунд
- **Image optimization:** Lazy loading и оптимизация аватаров
- **Admin privileges:** Условное отображение административных ссылок

---

## 11. ErrorBoundary

**Файл:** `/frontend/svetu/src/components/ErrorBoundary.tsx`  

### Назначение
Компонент для перехвата и обработки ошибок React с локализованными сообщениями.

### Основные компоненты и props
```typescript
interface ErrorMessages {
  title: string;
  description: string;
  details: string;  
  reload: string;
}

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  messages: ErrorMessages;
}
```

### Зависимости
- **React:** Component, ErrorInfo для class компонента
- **Интернационализация:** next-intl для переводов

### Ключевая логика работы
1. **Error catching:** getDerivedStateFromError для перехвата ошибок
2. **Error logging:** componentDidCatch для логирования
3. **Fallback UI:** Красивый интерфейс с деталями ошибки
4. **Recovery mechanism:** Кнопка перезагрузки страницы

### UI структура (error state)
```
ErrorBoundary
└── Card layout
    ├── Error title
    ├── Description
    ├── Collapsible details
    └── Reload button
```

### Особенности реализации
- **Class component:** Использует class syntax для error boundaries
- **Локализация:** Специализированный AuthErrorBoundary с переводами
- **Accessibility:** Proper ARIA roles и semantic HTML
- **Details disclosure:** Раскрывающиеся технические подробности

---

## Общие особенности базовых компонентов

### Архитектурные принципы
1. **Composition over inheritance:** Компоненты легко комбинируются
2. **Single responsibility:** Каждый компонент имеет четкую задачу  
3. **Error resilience:** Graceful handling ошибок на всех уровнях
4. **Accessibility first:** ARIA и семантический HTML

### Технологический стек
- **React 19** с новыми возможностями
- **Next.js 15.3.2** с App Router
- **TypeScript** для типобезопасности
- **DaisyUI + Tailwind** для стилизации
- **next-intl** для интернационализации

### Интеграции
- **Redux Toolkit** для state management
- **React Query** для server state
- **WebSocket** для real-time функций
- **Google OAuth** для аутентификации

### Performance оптимизации
- **SSR/SSG** где возможно
- **Code splitting** для крупных компонентов
- **Lazy loading** для изображений
- **Debouncing** для поисковых запросов
- **Memoization** для предотвращения лишних ререндеров