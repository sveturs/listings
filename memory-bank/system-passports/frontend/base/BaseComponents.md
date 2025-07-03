# Паспорт модуля: Base Components

## 📋 Метаданные
- **Название**: Base Components
- **Путь**: `frontend/svetu/src/app/` и основные компоненты
- **Роль**: Базовые компоненты архитектуры приложения
- **Уровень**: Основная структура приложения

## 🎯 Назначение
Модуль содержит фундаментальные компоненты, формирующие основную структуру и архитектуру приложения: layout'ы, навигацию, провайдеры и системные компоненты.

## 📊 Состав модуля: 5 компонентов

### 🏗️ Layout компоненты (2 компонента)

#### 1. RootLayout
**Путь**: `app/[locale]/layout.tsx`

```typescript
interface RootLayoutProps {
  children: React.ReactNode;
  params: { locale: 'en' | 'ru' };
}
```

**Назначение**: Корневой макет приложения с провайдерами и метаданными

**Ключевые особенности**:
- Интернационализация с next-intl
- Подключение всех провайдеров (Redux, Auth, WebSocket)
- SEO мета-теги и Open Graph
- Responsive viewport настройки
- DaisyUI theme integration
- Google Fonts загрузка

**Структура провайдеров**:
```jsx
<html lang={locale} data-theme="light">
  <body className="bg-base-100">
    <NextIntlClientProvider messages={messages}>
      <ReduxProvider>
        <AuthStateManager>
          <WebSocketManager>
            <ErrorBoundary>
              <Header />
              {children}
            </ErrorBoundary>
          </WebSocketManager>
        </AuthStateManager>
      </ReduxProvider>
    </NextIntlClientProvider>
  </body>
</html>
```

#### 2. AdminLayout
**Путь**: `app/[locale]/admin/layout.tsx`

```typescript
interface AdminLayoutProps {
  children: React.ReactNode;
}
```

**Назначение**: Административная панель с drawer навигацией

**Ключевые особенности**:
- AdminGuard для защиты доступа
- Drawer navigation с мобильной адаптацией
- Breadcrumbs навигация
- Sidebar с иконками меню
- Responsive layout

**Структура административной панели**:
```jsx
<AdminGuard>
  <div className="drawer lg:drawer-open">
    <input id="admin-drawer" type="checkbox" className="drawer-toggle" />
    
    <div className="drawer-content">
      <AdminHeader />
      <main className="p-6">
        <Breadcrumbs />
        {children}
      </main>
    </div>
    
    <div className="drawer-side">
      <AdminSidebar />
    </div>
  </div>
</AdminGuard>
```

### 🧭 Навигационные компоненты (1 компонент)

#### 3. Header
**Путь**: `components/Header.tsx`

```typescript
interface HeaderProps {
  className?: string;
}
```

**Назначение**: Главный заголовок с навигацией, поиском и аутентификацией

**Ключевые особенности**:
- Responsive дизайн с mobile drawer
- Интегрированная поисковая строка
- Аутентификация через Google OAuth
- Переключатель языков
- Уведомления и профиль пользователя
- Логотип с брендингом

**Навигационная структура**:
```jsx
<header className="navbar bg-base-100 border-b">
  <div className="navbar-start">
    <MobileMenuButton />
    <Logo />
  </div>
  
  <div className="navbar-center">
    <SearchBar className="w-full max-w-lg" />
  </div>
  
  <div className="navbar-end">
    <LanguageSwitcher />
    <NotificationsButton />
    <AuthButton />
  </div>
</header>
```

### 📄 Страницы (1 компонент)

#### 4. HomePage
**Путь**: `app/[locale]/page.tsx`

```typescript
interface HomePageProps {
  searchParams: { [key: string]: string | string[] | undefined };
}
```

**Назначение**: Главная страница маркетплейса

**Ключевые особенности**:
- MarketplaceList с товарами из OpenSearch
- Фильтрация по категориям
- Infinite scroll загрузка
- SEO оптимизация
- Responsive grid layout

**Структура главной страницы**:
```jsx
<main className="container mx-auto px-4 py-6">
  <section className="hero mb-8">
    <WelcomeBanner />
  </section>
  
  <section className="marketplace">
    <MarketplaceFilters />
    <MarketplaceList 
      searchParams={searchParams}
      showFilters={true}
    />
  </section>
</main>
```

### ⚙️ Системные провайдеры (1 компонент)

#### 5. AppProviders
**Путь**: Интегрирован в RootLayout

```typescript
interface AppProvidersProps {
  children: React.ReactNode;
  locale: string;
  messages: any;
}
```

**Назначение**: Объединение всех системных провайдеров

**Ключевые особенности**:
- Redux store с middleware
- React Query для кэширования
- Интернационализация
- WebSocket менеджер
- Аутентификация
- Error boundaries

**Цепочка провайдеров**:
```jsx
// 1. Локализация
<NextIntlClientProvider messages={messages}>
  
  // 2. Redux состояние
  <ReduxProvider>
    
    // 3. Аутентификация
    <AuthStateManager>
      
      // 4. WebSocket
      <WebSocketManager url={wsUrl}>
        
        // 5. Обработка ошибок
        <ErrorBoundary>
          {children}
        </ErrorBoundary>
        
      </WebSocketManager>
    </AuthStateManager>
  </ReduxProvider>
</NextIntlClientProvider>
```

## 🔗 Архитектурные зависимости

### Интернационализация
```typescript
// next-intl конфигурация
import { NextIntlClientProvider } from 'next-intl';
import { getMessages } from 'next-intl/server';

// Поддерживаемые локали
const locales = ['en', 'ru'] as const;
type Locale = typeof locales[number];
```

### Redux интеграция
```typescript
// Store с middleware
import { Provider } from 'react-redux';
import { store } from '@/store';

// RTK Query для API
import { setupListeners } from '@reduxjs/toolkit/query';
```

### WebSocket интеграция
```typescript
// Real-time обновления
interface WebSocketConfig {
  url: string;
  reconnectInterval: number;
  maxReconnectAttempts: number;
}
```

## 📱 Responsive дизайн

### Breakpoints (DaisyUI)
- **sm**: 640px+
- **md**: 768px+  
- **lg**: 1024px+
- **xl**: 1280px+

### Mobile-first подход
```css
/* Mobile (default) */
.navbar-center { display: none; }

/* Desktop */
@media (min-width: 768px) {
  .navbar-center { display: flex; }
}
```

### Drawer navigation
```jsx
// Mobile: Hamburger menu
<div className="lg:hidden">
  <label htmlFor="drawer-toggle" className="btn btn-square btn-ghost">
    <HamburgerIcon />
  </label>
</div>

// Desktop: Always visible sidebar
<div className="hidden lg:flex">
  <NavigationMenu />
</div>
```

## 🎨 Тематизация

### DaisyUI themes
```typescript
// Поддерживаемые темы
const themes = [
  'light',
  'dark', 
  'cupcake',
  'corporate'
] as const;

// Переключение темы
const toggleTheme = () => {
  document.documentElement.setAttribute('data-theme', newTheme);
};
```

### CSS переменные
```css
:root {
  --primary: 219 70% 50%;
  --secondary: 262 80% 50%;
  --accent: 321 70% 50%;
  --neutral: 222 13% 19%;
  --base-100: 0 0% 100%;
}
```

## 🔒 Безопасность

### AdminGuard защита
```typescript
const AdminGuard: FC<AdminGuardProps> = ({ children, requiredRole = 'admin' }) => {
  const { user } = useAuth();
  
  if (!user || user.role !== requiredRole) {
    return <UnauthorizedPage />;
  }
  
  return <>{children}</>;
};
```

### CSP Headers
```typescript
// Content Security Policy
const securityHeaders = [
  {
    key: 'Content-Security-Policy',
    value: contentSecurityPolicy.replace(/\s{2,}/g, ' ').trim()
  }
];
```

## 🌐 SEO оптимизация

### Metadata
```typescript
export const metadata: Metadata = {
  title: {
    default: 'Sve Tu - Marketplace',
    template: '%s | Sve Tu'
  },
  description: 'Sve Tu Platform - Serbian marketplace for local business',
  keywords: ['marketplace', 'serbia', 'local business'],
  authors: [{ name: 'Sve Tu Team' }],
  openGraph: {
    type: 'website',
    locale: 'sr_RS',
    url: 'https://svetu.rs',
    siteName: 'Sve Tu',
  },
  robots: {
    index: true,
    follow: true,
  }
};
```

### Structured data
```json
{
  "@context": "https://schema.org",
  "@type": "WebSite",
  "name": "Sve Tu",
  "url": "https://svetu.rs",
  "potentialAction": {
    "@type": "SearchAction",
    "target": "https://svetu.rs/search?q={search_term_string}",
    "query-input": "required name=search_term_string"
  }
}
```

## ⚡ Производительность

### Code splitting
```typescript
// Dynamic imports для больших компонентов
const AdminPanel = dynamic(() => import('@/components/AdminPanel'), {
  loading: () => <Skeleton />,
  ssr: false
});
```

### Image optimization
```jsx
<Image
  src="/logo.svg"
  alt="Sve Tu Logo"
  width={120}
  height={40}
  priority
  className="h-10 w-auto"
/>
```

### Bundle optimization
```typescript
// Tree shaking для иконок
import { SearchIcon, UserIcon } from '@heroicons/react/24/outline';

// Lazy loading для routes
const LazyRoute = lazy(() => import('./routes/LazyRoute'));
```

## 🎯 Примеры использования

### Создание новой страницы
```jsx
// app/[locale]/new-page/page.tsx
export default function NewPage() {
  return (
    <main className="container mx-auto px-4 py-6">
      <h1 className="text-2xl font-bold mb-4">New Page</h1>
      {/* Контент автоматически обернется в RootLayout */}
    </main>
  );
}
```

### Административная страница
```jsx
// app/[locale]/admin/new-admin-page/page.tsx
export default function NewAdminPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Admin Page</h1>
      {/* Автоматически защищено AdminGuard */}
    </div>
  );
}
```

### Добавление в Header
```jsx
// Новая кнопка в Header
<div className="navbar-end">
  <NewFeatureButton />
  <LanguageSwitcher />
  <AuthButton />
</div>
```

## 🐛 Известные особенности

1. **RootLayout**: Обязательно должен включать все провайдеры
2. **AdminLayout**: Требует роль 'admin' для доступа
3. **Header**: SearchBar скрывается на мобильных устройствах
4. **HomePage**: Зависит от OpenSearch для получения данных
5. **AppProviders**: Порядок провайдеров критичен для работы

## 🔄 Жизненный цикл приложения

1. **Initialization**: RootLayout → провайдеры → аутентификация
2. **Navigation**: Header → SearchBar/меню → роутинг
3. **Data Loading**: Redux → API → кэширование
4. **Real-time**: WebSocket → обновления → UI
5. **Error Handling**: ErrorBoundary → fallback → восстановление

## 🌍 Интернационализация

### Конфигурация локалей
```typescript
// i18n.ts
export const locales = ['en', 'ru'] as const;
export const defaultLocale = 'ru' as const;

// Переводы
const messages = {
  en: () => import('./messages/en.json'),
  ru: () => import('./messages/ru.json')
};
```

### Использование переводов
```jsx
import { useTranslations } from 'next-intl';

const Component = () => {
  const t = useTranslations('common');
  
  return (
    <h1>{t('welcome')}</h1>
  );
};
```