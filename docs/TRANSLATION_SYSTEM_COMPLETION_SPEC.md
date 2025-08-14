# Техническое задание на доработку системы переводов

## 1. Текущее состояние системы

### 1.1 Реализованные компоненты

#### Backend (полностью функционален):
- ✅ REST API эндпоинты для всех операций с переводами
- ✅ Версионирование переводов с историей изменений
- ✅ Аудит всех операций
- ✅ Массовые операции перевода
- ✅ Интеграция с AI провайдерами (OpenAI, Google, DeepL, Claude)
- ✅ Экспорт/импорт в форматах JSON, CSV, XLIFF
- ✅ Синхронизация между БД, Frontend и OpenSearch
- ✅ Статистика и аналитика переводов

#### Frontend (UI реализован, интеграция отсутствует):
- ✅ 15 компонентов системы переводов
- ✅ UI для всех основных операций
- ❌ Отсутствует рабочая интеграция с Backend API
- ❌ Проблемы с авторизацией и доступом

### 1.2 Основные проблемы

1. **Авторизация**: AdminGuard редиректит на несуществующую страницу `/auth/login`
2. **API интеграция**: JWT токен не передается в запросы автоматически
3. **Функциональность**: Многие компоненты работают только с демо-данными

## 2. Требования к доработке

### 2.1 Приоритет 1: Критические исправления

#### 2.1.1 Исправление системы авторизации

**Задача**: Интегрировать модальное окно авторизации с AdminGuard

**Требования**:
- AdminGuard должен показывать модальное окно входа вместо редиректа
- При успешной авторизации автоматически обновлять состояние
- Сохранять JWT токен в localStorage и cookie
- Автоматически обновлять токен при истечении

**Файлы для изменения**:
```
frontend/svetu/src/components/AdminGuard.tsx
frontend/svetu/src/components/AuthModal.tsx (создать)
frontend/svetu/src/contexts/AuthContext.tsx
```

**Примерная реализация**:
```typescript
// AdminGuard.tsx
export default function AdminGuard({ children }: AdminGuardProps) {
  const { user, isLoading, showAuthModal } = useAuth();
  const [showLoginModal, setShowLoginModal] = useState(false);

  useEffect(() => {
    if (!isLoading && !user) {
      setShowLoginModal(true);
    }
  }, [user, isLoading]);

  if (showLoginModal) {
    return <AuthModal onClose={() => setShowLoginModal(false)} />;
  }

  if (!user?.is_admin) {
    return <UnauthorizedMessage />;
  }

  return <>{children}</>;
}
```

#### 2.1.2 Реализация API клиента с автоматической авторизацией

**Задача**: Создать полноценный API клиент для системы переводов

**Требования**:
- Автоматически добавлять JWT токен во все запросы
- Обрабатывать ошибки авторизации (401)
- Автоматически обновлять токен при необходимости
- Типизированные методы для всех эндпоинтов

**Файл для создания**:
```
frontend/svetu/src/services/translationAdminApi.ts
```

**Структура API клиента**:
```typescript
class TranslationAdminApi {
  private baseUrl = '/api/v1/admin/translations';
  
  private async request<T>(
    path: string, 
    options?: RequestInit
  ): Promise<T> {
    const token = localStorage.getItem('access_token');
    
    const response = await fetch(`${this.baseUrl}${path}`, {
      ...options,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });
    
    if (response.status === 401) {
      await this.refreshToken();
      // Повторить запрос
    }
    
    return response.json();
  }

  // Методы API
  async getStatistics() { ... }
  async getProviders() { ... }
  async bulkTranslate(request: BulkTranslateRequest) { ... }
  // ... остальные методы
}

export const translationAdminApi = new TranslationAdminApi();
```

### 2.2 Приоритет 2: Основная функциональность

#### 2.2.1 Массовый перевод (BulkTranslationManager)

**Требования**:
- Загрузка реальных категорий/атрибутов/объявлений из API
- Выбор провайдера перевода
- Прогресс-бар для длительных операций
- Логирование результатов операции
- Возможность отмены операции

**API эндпоинты**:
```
POST /api/v1/admin/translations/bulk/translate
GET /api/v1/admin/categories/all
GET /api/v1/admin/attributes/all
GET /api/v1/marketplace/listings (с пагинацией)
```

#### 2.2.2 История версий (VersionHistoryViewer)

**Требования**:
- Загрузка реальной истории из API
- Визуальное сравнение версий (diff)
- Возможность отката к предыдущей версии
- Фильтрация по дате, пользователю, типу изменения

**API эндпоинты**:
```
GET /api/v1/admin/translations/versions/{entity}/{id}
GET /api/v1/admin/translations/versions/diff?v1={id1}&v2={id2}
POST /api/v1/admin/translations/versions/rollback
```

#### 2.2.3 Синхронизация (SyncManager)

**Требования**:
- Синхронизация Frontend JSON ↔ База данных
- Обнаружение и разрешение конфликтов
- Предпросмотр изменений перед применением
- Резервное копирование перед синхронизацией

**API эндпоинты**:
```
POST /api/v1/admin/translations/sync/frontend-to-db
POST /api/v1/admin/translations/sync/db-to-frontend
GET /api/v1/admin/translations/sync/conflicts
POST /api/v1/admin/translations/sync/conflicts/{id}/resolve
```

#### 2.2.4 Экспорт/Импорт (ExportImportManager)

**Требования**:
- Экспорт в форматы: JSON, CSV, XLIFF
- Импорт с валидацией данных
- Предпросмотр импортируемых данных
- Отчет о результатах импорта

**API эндпоинты**:
```
POST /api/v1/admin/translations/export/advanced
POST /api/v1/admin/translations/import/advanced
POST /api/v1/admin/translations/import/validate
```

### 2.3 Приоритет 3: Новая функциональность

#### 2.3.1 Редактор переводов (TranslationEditor)

**Новый компонент для редактирования отдельных переводов**

**Требования**:
- Inline редактирование в таблице
- Автосохранение через debounce
- Поддержка Markdown и переменных
- Предпросмотр результата
- История изменений конкретного перевода

**Макет интерфейса**:
```
┌─────────────────────────────────────────────────┐
│ Редактор перевода: category.electronics.name    │
├─────────────────────────────────────────────────┤
│ [EN] Electronics                                 │
│ [RU] [___________________________] [Сохранить]  │
│ [SR] [___________________________] [Сохранить]  │
├─────────────────────────────────────────────────┤
│ Переменные: {count}, {name}                     │
│ Markdown: **жирный**, *курсив*, [ссылка](url)   │
├─────────────────────────────────────────────────┤
│ История изменений:                              │
│ • 2025-08-13 - Изменено с "Электроника"        │
│ • 2025-08-10 - Создан перевод                  │
└─────────────────────────────────────────────────┘
```

#### 2.3.2 Поиск и фильтрация переводов

**Новый компонент TranslationSearch**

**Требования**:
- Полнотекстовый поиск по всем переводам
- Фильтры:
  - По языку
  - По модулю
  - По статусу (проверенные/непроверенные)
  - По типу (ручной/машинный)
  - По дате изменения
- Экспорт результатов поиска
- Массовые операции с найденными переводами

**API эндпоинт**:
```
GET /api/v1/admin/translations/search?q={query}&filters={...}
```

#### 2.3.3 Визуальный Diff Viewer

**Компонент для сравнения версий переводов**

**Требования**:
- Построчное сравнение текста
- Подсветка добавленных/удаленных/измененных частей
- Возможность выбора версий для сравнения
- Экспорт diff в текстовый файл

**Пример визуализации**:
```
Версия 1 (2025-08-10)          |  Версия 2 (2025-08-13)
--------------------------------|--------------------------------
Электроника и гаджеты          |  Электроника и [+цифровые] гаджеты
Лучшие [-предложения] товары   |  Лучшие товары
                                |  [+Новая строка с описанием]
```

#### 2.3.4 Dashboard с метриками

**Компонент TranslationMetrics**

**Требования**:
- Графики:
  - Прогресс переводов по времени
  - Распределение по языкам
  - Активность переводчиков
  - Использование AI провайдеров
- KPI метрики:
  - Процент завершенности
  - Скорость перевода
  - Качество (процент проверенных)
  - Экономия от использования AI
- Экспорт отчетов в PDF

### 2.4 Приоритет 4: UX улучшения

#### 2.4.1 Система уведомлений

**Требования**:
- Toast уведомления для всех операций
- Различные типы: success, error, warning, info
- Возможность отмены операции из уведомления
- История уведомлений

**Реализация**:
```typescript
// hooks/useTranslationNotifications.ts
export function useTranslationNotifications() {
  const showSuccess = (message: string, action?: () => void) => {
    toast.success(message, {
      action: action ? { label: 'Отменить', onClick: action } : undefined
    });
  };
  
  return { showSuccess, showError, showWarning };
}
```

#### 2.4.2 Подтверждения опасных операций

**Требования**:
- Модальные окна подтверждения для:
  - Удаления переводов
  - Отката версий
  - Массовых операций
  - Импорта с перезаписью
- Информация о последствиях операции
- Возможность создания резервной копии

#### 2.4.3 Индикаторы загрузки

**Требования**:
- Skeleton loaders для таблиц
- Progress bar для длительных операций
- Отображение этапа операции
- Возможность отмены

#### 2.4.4 Горячие клавиши

**Требования**:
- Ctrl+S - сохранить перевод
- Ctrl+Z - отменить изменение
- Ctrl+F - поиск
- Ctrl+E - экспорт
- Ctrl+I - импорт
- Показ подсказок по горячим клавишам

### 2.5 Приоритет 5: Тестирование

#### 2.5.1 Unit тесты

**Покрытие минимум 80% для**:
- API клиента
- Утилит валидации
- Компонентов форм
- Хуков

**Технологии**:
- Jest
- React Testing Library
- MSW для мокирования API

#### 2.5.2 E2E тесты

**Основные сценарии**:
1. Авторизация администратора
2. Массовый перевод категорий
3. Экспорт и импорт переводов
4. Откат версии перевода
5. Синхронизация с разрешением конфликтов

**Технология**: Playwright

#### 2.5.3 Тестирование производительности

**Требования**:
- Загрузка страницы < 2 секунд
- Отклик на действия < 200мс
- Работа с 10000+ переводов без зависаний
- Оптимизация bundle size

## 3. План реализации

### Этап 1: Критические исправления (3-4 дня)
- [ ] Исправление AdminGuard
- [ ] Реализация API клиента
- [ ] Интеграция JWT токенов

### Этап 2: Основная функциональность (5-7 дней)
- [ ] Массовый перевод
- [ ] История версий
- [ ] Синхронизация
- [ ] Экспорт/Импорт

### Этап 3: Новая функциональность (5-7 дней)
- [ ] Редактор переводов
- [ ] Поиск и фильтрация
- [ ] Визуальный Diff
- [ ] Dashboard с метриками

### Этап 4: UX улучшения (3-4 дня)
- [ ] Система уведомлений
- [ ] Подтверждения операций
- [ ] Индикаторы загрузки
- [ ] Горячие клавиши

### Этап 5: Тестирование (3-4 дня)
- [ ] Unit тесты
- [ ] E2E тесты
- [ ] Тестирование производительности
- [ ] Исправление найденных багов

**Общее время реализации: 19-26 дней**

## 4. Критерии приемки

### 4.1 Функциональные требования
- ✅ Все компоненты работают с реальными данными из API
- ✅ Авторизация работает через модальное окно
- ✅ JWT токен автоматически передается во все запросы
- ✅ Все основные операции с переводами доступны
- ✅ Система обрабатывает ошибки корректно

### 4.2 Нефункциональные требования
- ✅ Производительность: загрузка страницы < 2 сек
- ✅ Доступность: WCAG 2.1 Level AA
- ✅ Совместимость: Chrome, Firefox, Safari, Edge
- ✅ Безопасность: защита от XSS, CSRF
- ✅ Масштабируемость: работа с 100k+ переводов

### 4.3 Документация
- ✅ Техническая документация API
- ✅ Руководство пользователя
- ✅ Инструкция по развертыванию
- ✅ Примеры использования

## 5. Риски и митигация

### 5.1 Технические риски
| Риск | Вероятность | Влияние | Митигация |
|------|------------|---------|-----------|
| Проблемы с производительностью при большом объеме данных | Средняя | Высокое | Пагинация, виртуальная прокрутка, кеширование |
| Конфликты при одновременном редактировании | Низкая | Среднее | Оптимистичные блокировки, WebSocket уведомления |
| Потеря данных при синхронизации | Низкая | Высокое | Резервное копирование, транзакции, подтверждения |

### 5.2 Организационные риски
| Риск | Вероятность | Влияние | Митигация |
|------|------------|---------|-----------|
| Изменение требований | Средняя | Среднее | Agile подход, регулярные демо |
| Недостаток времени | Средняя | Высокое | Приоритизация, MVP подход |
| Недостаток экспертизы | Низкая | Среднее | Консультации, код-ревью |

## 6. Определение готовности (Definition of Done)

Задача считается выполненной когда:
1. ✅ Код написан и соответствует стандартам
2. ✅ Написаны unit тесты (покрытие > 80%)
3. ✅ Пройдено код-ревью
4. ✅ Обновлена документация
5. ✅ Проведено ручное тестирование
6. ✅ Исправлены найденные баги
7. ✅ Задеплоено на staging
8. ✅ Получено подтверждение от QA

## 7. Технический стек

### Frontend
- Next.js 15.3.2
- React 19
- TypeScript 5.x
- Tailwind CSS v4
- DaisyUI 5.x
- Redux Toolkit
- React Query
- React Hook Form
- Zod (валидация)

### Backend
- Go 1.23
- Fiber v2
- PostgreSQL 15
- Redis (кеширование)
- OpenSearch
- MinIO (S3)

### Инструменты
- ESLint + Prettier
- Jest + React Testing Library
- Playwright (E2E)
- Docker + Docker Compose
- GitHub Actions (CI/CD)

## 8. Контакты и ресурсы

### Команда
- Frontend Lead: [контакт]
- Backend Lead: [контакт]
- QA Lead: [контакт]
- Product Owner: [контакт]

### Ресурсы
- [Figma дизайны](https://figma.com/...)
- [API документация](http://localhost:3000/swagger)
- [Confluence](https://...)
- [Jira проект](https://...)

## 9. Приложения

### Приложение A: Структура базы данных

```sql
-- Основная таблица переводов
CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    field_name VARCHAR(100) NOT NULL,
    language VARCHAR(5) NOT NULL,
    translated_text TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    translated_by INTEGER,
    verified_by INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(entity_type, entity_id, field_name, language)
);

-- Версионирование
CREATE TABLE translation_versions (
    id SERIAL PRIMARY KEY,
    translation_id INTEGER REFERENCES translations(id),
    version INTEGER NOT NULL,
    translated_text TEXT,
    changed_by INTEGER,
    changed_at TIMESTAMP DEFAULT NOW(),
    change_reason TEXT
);

-- Аудит
CREATE TABLE translation_audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    action VARCHAR(50),
    entity_type VARCHAR(50),
    entity_id INTEGER,
    old_value TEXT,
    new_value TEXT,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Приложение B: API эндпоинты

```yaml
/api/v1/admin/translations:
  /stats/overview:
    GET: Получить общую статистику
  /providers:
    GET: Список провайдеров переводов
    PUT: Обновить настройки провайдера
  /bulk/translate:
    POST: Массовый перевод
  /versions/{entity}/{id}:
    GET: История версий
  /versions/rollback:
    POST: Откат версии
  /sync/frontend-to-db:
    POST: Синхронизация Frontend -> БД
  /sync/db-to-frontend:
    POST: Синхронизация БД -> Frontend
  /sync/conflicts:
    GET: Список конфликтов
  /sync/conflicts/{id}/resolve:
    POST: Разрешить конфликт
  /export/advanced:
    POST: Экспорт с параметрами
  /import/advanced:
    POST: Импорт с валидацией
  /audit/logs:
    GET: Журнал аудита
  /search:
    GET: Поиск переводов
```

### Приложение C: Примеры использования

#### Пример 1: Массовый перевод категорий

```typescript
// Массовый перевод всех категорий на русский и сербский
const result = await translationAdminApi.bulkTranslate({
  entity_type: 'category',
  entity_ids: [1, 2, 3, 4, 5],
  source_language: 'en',
  target_languages: ['ru', 'sr'],
  provider_id: 1, // OpenAI
  auto_approve: false,
  overwrite_existing: false
});

console.log(`Переведено: ${result.successful} из ${result.total_processed}`);
```

#### Пример 2: Экспорт переводов для внешнего переводчика

```typescript
// Экспорт непроверенных переводов в XLIFF для профессионального перевода
const exportData = await translationAdminApi.export({
  format: 'xliff',
  language: 'ru',
  only_verified: false,
  include_metadata: true
});

// Сохранить файл
downloadFile(exportData, 'translations_for_review.xliff');
```

#### Пример 3: Откат версии перевода

```typescript
// Откатить перевод категории к предыдущей версии
const versions = await translationAdminApi.getVersionHistory('category', 123);
const previousVersion = versions[1]; // Предыдущая версия

await translationAdminApi.rollbackVersion({
  version_id: previousVersion.id,
  reason: 'Некорректный перевод от AI'
});
```

---

**Документ подготовлен**: 13.08.2025
**Версия**: 1.0
**Статус**: На согласовании