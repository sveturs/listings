# Индекс Common компонентов Sve Tu Platform

## Компоненты UI

### Базовые элементы
- **FormField** - Обертка для полей формы с лейблом и ошибками
- **GoogleIcon** - SVG иконка Google с официальными цветами

### Изображения
- **OptimizedImage** - Оптимизированные изображения с CDN и fallback
- **SafeImage** - Безопасные изображения с проверкой URL

### Навигация и управление
- **ViewToggle** - Переключатель сетка/список
- **LanguageSwitcher** - Переключатель языков RU/EN
- **InfiniteScrollTrigger** - Бесконечная прокрутка

### Выбор и ввод
- **IconPicker** - Выбор эмодзи из 240+ иконок по категориям

## Компоненты состояния

### Черновики
- **DraftStatus** - Статус сохранения черновика
- **DraftIndicator** - Индикатор несохраненных изменений  
- **OfflineIndicator** - Индикатор оффлайн режима
- **DraftsModal** - Управление черновиками

### Ошибки
- **ErrorBoundary/AuthErrorBoundary** - Обработчик ошибок React

## Системные компоненты

### Провайдеры
- **ReduxProvider** - Redux + React Query провайдер

### Менеджеры
- **WebSocketManager** - Управление WebSocket соединением
- **AuthStateManager** - Очистка данных при смене пользователя

### Защита
- **AdminGuard** - Защита админских роутов

## Основные зависимости

- **UI**: DaisyUI, Tailwind CSS, @heroicons/react
- **Состояние**: Redux Toolkit, React Query
- **Локализация**: next-intl, date-fns locales
- **Изображения**: next/image
- **Маршрутизация**: Next.js App Router

## Ключевые паттерны

1. Все компоненты - 'use client'
2. Строгая TypeScript типизация
3. Поддержка локализации
4. DaisyUI для стилей
5. Доступность (ARIA)
6. Обработка ошибок и fallback