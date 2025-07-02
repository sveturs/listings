# Паспорта Common (переиспользуемых) компонентов Sve Tu Platform

## 1. InfiniteScrollTrigger

**Путь:** `/frontend/svetu/src/components/common/InfiniteScrollTrigger.tsx`

### Назначение
Компонент для реализации бесконечной прокрутки с визуальным триггером и fallback-кнопкой.

### Props
```typescript
interface InfiniteScrollTriggerProps {
  loading: boolean;              // Индикатор загрузки
  hasMore: boolean;              // Есть ли еще данные
  onLoadMore: () => void;        // Колбэк загрузки
  showButton?: boolean;          // Показывать кнопку (default: true)
  className?: string;            // Дополнительные классы
  loadMoreText?: string;         // Текст кнопки (default: 'Load more')
}
```

### Особенности реализации
- Использует `forwardRef` для передачи ref к элементу-триггеру
- Триггер-элемент высотой 20px с отступом 8px сверху
- Встроенный спиннер DaisyUI при загрузке
- Fallback-кнопка с иконкой стрелки вниз
- Доступность: `aria-hidden` для триггера, `aria-label` для кнопки

### UI структура
1. Невидимый div-триггер для IntersectionObserver
2. Спиннер загрузки внутри триггера
3. Опциональная кнопка "Load more"

### Интеграция
- Работает с IntersectionObserver API
- Совместим с любыми списками и сетками
- Поддерживает кастомизацию через className

---

## 2. ViewToggle

**Путь:** `/frontend/svetu/src/components/common/ViewToggle.tsx`

### Назначение
Переключатель режимов отображения между сеткой и списком.

### Props
```typescript
export type ViewMode = 'grid' | 'list';

interface ViewToggleProps {
  currentView: ViewMode;
  onViewChange: (view: ViewMode) => void;
}
```

### Зависимости
- `next-intl` для локализации
- `@heroicons/react` для иконок

### UI структура
- Контейнер с фоном base-200 и скругленными углами
- Две кнопки с иконками Squares2X2Icon и ListBulletIcon
- Активная кнопка имеет класс btn-primary
- Tooltips с локализованными подсказками

### Особенности
- Использует DaisyUI классы
- Полная поддержка локализации
- Компактный дизайн с gap-1

---

## 3. ErrorBoundary (AuthErrorBoundary)

**Путь:** `/frontend/svetu/src/components/ErrorBoundary.tsx`

### Назначение
Обработчик ошибок React с локализованными сообщениями и красивым UI.

### Props
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

### Архитектура
- Класс-компонент `ErrorBoundaryClass` для перехвата ошибок
- Функциональная обертка `AuthErrorBoundary` для локализации
- Компонент `DefaultErrorFallback` для UI ошибки

### Особенности
- Логирование ошибок в консоль
- Складываемые детали ошибки (details)
- Кнопка перезагрузки страницы
- Центрированная карточка DaisyUI
- Мемоизация сообщений через `useMemo`

### UI структура
- Полноэкранный контейнер с центрированием
- Карточка 384px шириной
- Заголовок с классом text-error
- Раскрываемый блок с деталями ошибки
- Кнопка перезагрузки

---

## 4. FormField

**Путь:** `/frontend/svetu/src/components/FormField.tsx`

### Назначение
Универсальная обертка для полей формы с лейблом и обработкой ошибок.

### Props
```typescript
interface FormFieldProps {
  label: string;
  required?: boolean;
  children: React.ReactNode;
  error?: string;
  className?: string;
}
```

### UI структура
1. Контейнер с классом form-control
2. Label с текстом и звездочкой для обязательных полей
3. Слот для дочернего элемента (input/select/textarea)
4. Опциональное сообщение об ошибке

### Особенности
- Использует DaisyUI классы form-control
- Красная звездочка для required полей
- Красный текст ошибки с классом label-text-alt
- Полная кастомизация через className

---

## 5. OptimizedImage

**Путь:** `/frontend/svetu/src/components/OptimizedImage.tsx`

### Назначение
Оптимизированный компонент изображения с поддержкой CDN и fallback.

### Props
```typescript
interface OptimizedImageProps {
  src: string;
  alt: string;
  width: number;
  height: number;
  className?: string;
  priority?: boolean;
}
```

### Зависимости
- `next/image` для оптимизации
- `configManager` для построения URL

### Особенности
- Автоматическое построение полного URL через configManager
- Fallback на placeholder при ошибке
- Blur placeholder для улучшения UX
- Встроенный base64 blur data URL
- Поддержка priority загрузки

### Обработка ошибок
- State для отслеживания ошибок
- Автоматический переход на `/placeholder-listing.jpg`

---

## 6. SafeImage

**Путь:** `/frontend/svetu/src/components/SafeImage.tsx`

### Назначение
Безопасный компонент изображения с проверкой URL и продвинутыми состояниями загрузки.

### Props
```typescript
interface SafeImageProps extends Omit<ImageProps, 'src'> {
  src: string | null | undefined;
  fallback?: React.ReactNode;
}
```

### Зависимости
- `next/image` для рендеринга
- `getSafeImageUrl` из imageUtils

### Особенности безопасности
- Проверка URL через белый список доменов
- Поддержка относительных путей
- Блокировка внешних небезопасных URL

### Состояния
1. **Loading** - анимированный placeholder
2. **Error** - кастомный fallback или SVG иконка
3. **Success** - плавное появление через opacity transition

### Debug функции
- Специальная отладка для объявления 177
- Логирование состояний в консоль

### UI fallback
- Серый контейнер с иконкой изображения
- Адаптивные размеры на основе props
- Минимальная высота 100px

---

## 7. LanguageSwitcher

**Путь:** `/frontend/svetu/src/components/LanguageSwitcher.tsx`

### Назначение
Переключатель языков с поддержкой маршрутизации и анимаций.

### Зависимости
- `next-intl` для текущей локали
- Кастомный `useRouter` из i18n/routing
- React `useTransition` для анимаций

### Особенности
- Dropdown меню DaisyUI
- Поддержка RU/EN локалей
- Отключение во время переключения
- Сохранение текущего пути при смене языка
- Плавная анимация через startTransition

### UI структура
- Кнопка с текущей локалью в верхнем регистре
- Иконка стрелки вниз
- Выпадающее меню шириной 96px
- Активный язык подсвечен

---

## 8. DraftStatus

**Путь:** `/frontend/svetu/src/components/DraftStatus.tsx`

### Назначение
Набор компонентов для отображения статуса черновиков и работы оффлайн.

### Компоненты

#### DraftStatus
- Показывает статус сохранения черновика
- Три состояния: сохранение, сохранено, есть изменения
- Использует date-fns для форматирования времени
- Поддержка SR и EN локалей

#### DraftIndicator
- Кнопка с иконкой документа
- Индикатор несохраненных изменений (желтая точка)
- Tooltip с локализованным текстом

#### OfflineIndicator
- Alert компонент для оффлайн режима
- Проверка navigator.onLine
- Предупреждение с иконкой

### Зависимости
- `CreateListingContext` для состояния
- `date-fns` с поддержкой локалей
- `next-intl` для переводов

---

## 9. DraftsModal

**Путь:** `/frontend/svetu/src/components/DraftsModal.tsx`

### Назначение
Модальное окно для управления черновиками объявлений.

### Props
```typescript
interface DraftsModalProps {
  isOpen: boolean;
  onClose: () => void;
}
```

### Функционал
- Список всех черновиков пользователя
- Открытие черновика для редактирования
- Удаление черновиков
- Экспорт в JSON
- Отображение времени обновления и истечения

### Зависимости
- `useListingDrafts` hook
- `draftService` для операций
- `toast` для уведомлений
- `date-fns` для дат

### UI структура
- Модальное окно max-w-4xl
- Карточки черновиков с информацией
- Dropdown меню с действиями
- Состояния: загрузка, пустой список, список

### Особенности
- Автоматическое обновление после удаления
- Индикация полноты черновика
- Экспорт с timestamp в имени файла

---

## 10. IconPicker

**Путь:** `/frontend/svetu/src/components/IconPicker.tsx`

### Назначение
Выбор эмодзи-иконок для категорий и атрибутов.

### Props
```typescript
interface IconPickerProps {
  value: string;
  onChange: (icon: string) => void;
  placeholder?: string;
}
```

### Категории иконок
- Транспорт (20 иконок)
- Электроника (20 иконок)
- Дом и быт (20 иконок)
- Одежда (20 иконок)
- Еда и напитки (20 иконок)
- Спорт и отдых (20 иконок)
- Красота и здоровье (20 иконок)
- Книги и обучение (20 иконок)
- Природа и животные (20 иконок)
- Инструменты (20 иконок)
- Числа и символы (20 иконок)
- Атрибуты (20 иконок)

### UI структура
- Input для ручного ввода
- Кнопка с текущей иконкой
- Popup с табами категорий
- Сетка иконок 8 колонок
- Скролл для длинных списков

### Особенности
- 240+ эмодзи иконок
- Группировка по категориям
- Поддержка ручного ввода
- Подсветка выбранной иконки
- Абсолютное позиционирование popup

---

## 11. GoogleIcon

**Путь:** `/frontend/svetu/src/components/GoogleIcon.tsx`

### Назначение
SVG иконка Google с официальными цветами бренда.

### Props
```typescript
interface GoogleIconProps {
  className?: string;  // default: 'w-5 h-5'
}
```

### Особенности
- Официальные цвета Google (#4285F4, #34A853, #FBBC05, #EA4335)
- ViewBox 24x24
- Полностью векторная графика
- Кастомизируемый размер через className

---

## 12. WebSocketManager

**Путь:** `/frontend/svetu/src/components/WebSocketManager.tsx`

### Назначение
Управляющий компонент для WebSocket соединения и чатов.

### Зависимости
- `AuthContext` для пользователя
- `useChat` hook для WebSocket
- `tokenManager` для токенов

### Логика работы
1. Инициализация только для авторизованных
2. Задержка 500мс для ожидания токена
3. Автоматическое переподключение
4. Очистка при logout или unmount

### Особенности
- Не рендерит UI (return null)
- Использует ref для предотвращения дублирования
- Логирование всех этапов жизненного цикла
- Правильная очистка таймеров

### Интеграция
- Передает getUserId функцию в WebSocket
- Автоматически загружает чаты после подключения

---

## 13. ReduxProvider

**Путь:** `/frontend/svetu/src/components/ReduxProvider.tsx`

### Назначение
Провайдер для Redux store и React Query.

### Зависимости
- `react-redux` Provider
- `@tanstack/react-query`
- Redux store из `/store`

### Конфигурация Query Client
```typescript
{
  queries: {
    staleTime: 60 * 1000,        // 1 минута
    refetchOnWindowFocus: false   // Отключен рефетч при фокусе
  }
}
```

### Особенности
- Создание QueryClient через useState для SSR
- Два провайдера в одном компоненте
- Client-side only компонент

---

## 14. AdminGuard

**Путь:** `/frontend/svetu/src/components/AdminGuard.tsx`

### Назначение
Защита админских роутов с проверкой прав.

### Props
```typescript
interface AdminGuardProps {
  children: React.ReactNode;
  loading?: React.ReactNode;
}
```

### Логика защиты
1. Проверка монтирования компонента
2. Проверка загрузки auth
3. Проверка авторизации
4. Проверка флага is_admin
5. Редирект на главную при отсутствии прав

### Состояния
- **Unmounted/Loading** - спиннер
- **Unauthorized** - редирект + спиннер
- **Authorized Admin** - рендер children

### Особенности
- Предотвращение hydration mismatch
- Кастомный loading компонент
- Использование isMounted паттерна

---

## 15. AuthStateManager

**Путь:** `/frontend/svetu/src/components/AuthStateManager.tsx`

### Назначение
Управление состоянием при смене пользователя.

### Зависимости
- `AuthContext` для отслеживания user
- `useChat` для очистки данных чата

### Логика работы
1. Отслеживание изменения userId через ref
2. При смене пользователя:
   - Очистка данных чата
   - Очистка sessionStorage (кроме токенов)
   - Очистка localStorage (кроме локали)
3. Сохранение важных данных

### Сохраняемые ключи
- `svetu_access_token`
- `svetu_user`
- `client_id`
- `NEXT_LOCALE`

### Особенности
- Не рендерит UI (return null)
- Безопасная очистка с try/catch
- Логирование действий

---

## Общие паттерны и best practices

1. **Client-side компоненты** - все помечены 'use client'
2. **Локализация** - использование next-intl во всех UI компонентах
3. **DaisyUI** - последовательное использование классов
4. **Доступность** - aria атрибуты, роли, семантичная разметка
5. **Обработка ошибок** - fallback UI, логирование
6. **Оптимизация** - lazy loading, мемоизация, переиспользование
7. **TypeScript** - строгая типизация всех props и состояний