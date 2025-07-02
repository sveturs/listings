# Паспорт компонента: ImportManager

## Метаданные
- **Путь**: `/frontend/svetu/src/components/import/ImportManager.tsx`
- **Роль**: Центр управления импортом товаров
- **Уровень**: Компонент страницы
- **Категория**: Import

## Описание

ImportManager - главный контейнер системы импорта, объединяющий все компоненты импорта в единый интерфейс. Отображает статистику импортов, список задач, поддерживаемые форматы и управляет модальными окнами.

## Props

```typescript
interface ImportManagerProps {
  storefrontId: number;      // ID витрины
  storefrontSlug?: string;   // Slug витрины для API
}
```

## Состояние

### Redux State (useAppSelector)
```typescript
{
  jobs: ImportJob[];         // Список задач импорта
  formats: ImportFormats;    // Поддерживаемые форматы
}
```

### Вычисляемые значения
```typescript
const stats = {
  pending: number;     // Ожидающие задачи
  processing: number;  // Обрабатываемые задачи
  completed: number;   // Завершенные задачи
  failed: number;      // Неудачные задачи
};
```

## Основная логика

### Инициализация
```typescript
useEffect(() => {
  dispatch(fetchImportJobs({ storefrontId }));
  dispatch(fetchImportFormats());
}, [dispatch, storefrontId]);
```

### Обработчики событий
```typescript
// Открытие мастера импорта
const handleStartImport = () => {
  dispatch(setImportModalOpen(true));
};

// Обновление после успешного импорта
const handleImportSuccess = (jobId: number) => {
  dispatch(fetchImportJobs({ storefrontId }));
};

// Подсчет статистики
const getJobsStats = () => {
  const pending = jobs.filter(job => job.status === 'pending').length;
  const processing = jobs.filter(job => job.status === 'processing').length;
  const completed = jobs.filter(job => job.status === 'completed').length;
  const failed = jobs.filter(job => job.status === 'failed').length;
  
  return { pending, processing, completed, failed };
};
```

## UI структура

### 1. Заголовок
```tsx
<div className="bg-white shadow rounded-lg">
  <h2>Import Products</h2>
  <p>Import products from CSV, XML or API sources</p>
  <button onClick={handleStartImport}>
    Start Import
  </button>
</div>
```

### 2. Статистика (4 карточки)
- **Pending**: Желтая иконка часов
- **Processing**: Синяя анимированная иконка
- **Completed**: Зеленая галочка
- **Failed**: Красный крестик

### 3. Поддерживаемые форматы
```tsx
{formats && (
  <div className="grid grid-cols-3 gap-6">
    {Object.entries(formats.supported_formats).map(([format, info]) => (
      <div key={format}>
        <span>{format.toUpperCase()}</span>
        <p>{info.description}</p>
        <p>Extensions: {info.file_extensions.join(', ')}</p>
        {info.encoding && <p>Encoding: {info.encoding}</p>}
      </div>
    ))}
  </div>
)}
```

### 4. Список задач импорта
```tsx
<ImportJobsList
  storefrontId={storefrontId}
  autoRefresh={true}
  refreshInterval={5000}
/>
```

### 5. Модальные окна
```tsx
<ImportWizard
  storefrontId={storefrontId}
  storefrontSlug={storefrontSlug}
  onSuccess={handleImportSuccess}
/>
<ImportJobDetails />
<ImportErrorsModal />
```

## Дочерние компоненты

1. **ImportWizard**: Мастер создания новой задачи импорта
2. **ImportJobsList**: Таблица с историей импортов
3. **ImportJobDetails**: Детальная информация о задаче
4. **ImportErrorsModal**: Просмотр ошибок импорта

## Зависимости

### Redux
- `importSlice`: Управление состоянием
- Actions: `fetchImportJobs`, `fetchImportFormats`, `setImportModalOpen`

### Локализация
- `useTranslations('import')`: Переводы для секции импорта

## Пример использования

```tsx
// На странице управления товарами витрины
function StorefrontProductsPage({ storefront }) {
  return (
    <div className="container mx-auto">
      <h1>Products Management</h1>
      
      <ImportManager
        storefrontId={storefront.id}
        storefrontSlug={storefront.slug}
      />
      
      {/* Другие компоненты управления товарами */}
    </div>
  );
}
```

## Особенности

### Автообновление
- Список задач обновляется каждые 5 секунд
- Обновление только для активных задач (pending/processing)

### Адаптивный дизайн
- Grid layout для карточек статистики
- Responsive сетка для форматов (1-3 колонки)

### Условный рендеринг
- Статистика показывается только при наличии задач
- Форматы отображаются после загрузки

## Стилизация

### Компоненты DaisyUI
- Card компоненты для секций
- Badge для форматов файлов
- Button variants для действий

### Tailwind классы
- Shadow и rounded для карточек
- Grid система для layout
- Анимация для processing иконки

## Проблемы и TODO

1. **Отсутствие пагинации**: Список задач может стать длинным
2. **Нет фильтрации**: Невозможно отфильтровать задачи по дате
3. **Жесткий интервал**: 5 секунд не настраивается
4. **Нет экспорта**: Отсутствует возможность экспортировать историю

## Рекомендации

1. Добавить пагинацию для списка задач
2. Реализовать фильтры по дате и статусу
3. Сделать интервал обновления настраиваемым
4. Добавить возможность экспорта истории импортов
5. Показывать общую статистику (всего импортировано товаров)