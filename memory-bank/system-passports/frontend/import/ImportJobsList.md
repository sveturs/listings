# Паспорт компонента: ImportJobsList

## Метаданные
- **Путь**: `/frontend/svetu/src/components/import/ImportJobsList.tsx`
- **Роль**: Список задач импорта с автообновлением
- **Уровень**: Компонент-виджет
- **Категория**: Import

## Описание

ImportJobsList отображает таблицу всех задач импорта с их статусами, прогрессом и статистикой. Поддерживает автоматическое обновление активных задач, фильтрацию по статусу и действия над задачами.

## Props

```typescript
interface ImportJobsListProps {
  storefrontId: number;        // ID витрины
  autoRefresh?: boolean;       // Автообновление (default: true)
  refreshInterval?: number;    // Интервал обновления в мс (default: 5000)
}
```

## Состояние

### Redux State (useAppSelector)
```typescript
{
  jobs: ImportJob[];     // Список задач импорта
  isLoading: boolean;    // Состояние загрузки
}
```

### Local State
```typescript
{
  selectedStatus: string;  // Фильтр по статусу ('all' | 'pending' | 'processing' | 'completed' | 'failed')
}
```

## Основная логика

### Автообновление
```typescript
useEffect(() => {
  if (!autoRefresh) return;
  
  const hasActiveJobs = jobs.some(
    job => job.status === 'pending' || job.status === 'processing'
  );
  
  if (!hasActiveJobs) return;
  
  const interval = setInterval(() => {
    jobs.forEach(job => {
      if (job.status === 'pending' || job.status === 'processing') {
        dispatch(fetchJobStatus(job.id));
      }
    });
  }, refreshInterval);
  
  return () => clearInterval(interval);
}, [jobs, autoRefresh, refreshInterval, dispatch]);
```

### Обработчики действий
```typescript
// Показать детали задачи
const handleJobClick = (job: ImportJob) => {
  dispatch(setCurrentJob(job));
  dispatch(setJobDetailsModalOpen(true));
};

// Показать ошибки
const handleShowErrors = (job: ImportJob) => {
  dispatch(setCurrentJob(job));
  dispatch(setErrorsModalOpen(true));
};

// Отменить задачу
const handleCancelJob = async (jobId: number, e: React.MouseEvent) => {
  e.stopPropagation();
  if (confirm(t('actions.confirmCancel'))) {
    await dispatch(cancelImportJob(jobId));
  }
};

// Повторить задачу
const handleRetryJob = async (jobId: number, e: React.MouseEvent) => {
  e.stopPropagation();
  await dispatch(retryImportJob(jobId));
};
```

### Вспомогательные функции
```typescript
// Расчет прогресса
const calculateProgress = (job: ImportJob) => {
  if (job.total_records === 0) return 0;
  return Math.round((job.processed_records / job.total_records) * 100);
};

// Форматирование даты
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString();
};

// Получение иконки статуса
const getStatusIcon = (status: string) => {
  // Возвращает соответствующую SVG иконку
};
```

## UI структура

### Заголовок с фильтрами
```tsx
<div className="flex items-center justify-between">
  <h3>Import Jobs</h3>
  
  <div className="flex items-center space-x-2">
    <select value={selectedStatus} onChange={...}>
      <option value="all">All</option>
      <option value="pending">Pending</option>
      <option value="processing">Processing</option>
      <option value="completed">Completed</option>
      <option value="failed">Failed</option>
    </select>
    
    <button onClick={refreshJobs}>
      <RefreshIcon className={isLoading ? 'animate-spin' : ''} />
      Refresh
    </button>
  </div>
</div>
```

### Список задач
```tsx
{filteredJobs.map(job => (
  <li key={job.id} onClick={() => handleJobClick(job)}>
    {/* Статус с иконкой */}
    <span className={IMPORT_STATUS_COLORS[job.status]}>
      {getStatusIcon(job.status)}
      {t(`status.${job.status}`)}
    </span>
    
    {/* Информация о файле */}
    <div>
      <p>{job.file_name}</p>
      <span>{job.file_type.toUpperCase()}</span>
      <p>{formatDate(job.created_at)}</p>
    </div>
    
    {/* Прогресс для активных задач */}
    {job.status === 'processing' && (
      <div className="w-24 bg-gray-200 rounded-full h-2">
        <div 
          className="bg-blue-600 h-2 rounded-full"
          style={{ width: `${calculateProgress(job)}%` }}
        />
      </div>
    )}
    
    {/* Статистика */}
    <div className="text-sm">
      <span className="text-green-600">✓ {job.successful_records}</span>
      {job.failed_records > 0 && (
        <span className="text-red-600">✗ {job.failed_records}</span>
      )}
      <span>/ {job.total_records}</span>
    </div>
    
    {/* Действия */}
    <div className="flex space-x-2">
      <button onClick={() => handleJobClick(job)}>View Details</button>
      
      {job.failed_records > 0 && (
        <button onClick={() => handleShowErrors(job)}>View Errors</button>
      )}
      
      {job.status === 'processing' && (
        <button onClick={(e) => handleCancelJob(job.id, e)}>Cancel</button>
      )}
      
      {job.status === 'failed' && (
        <button onClick={(e) => handleRetryJob(job.id, e)}>Retry</button>
      )}
    </div>
  </li>
))}
```

### Пустое состояние
```tsx
{filteredJobs.length === 0 && (
  <div className="text-center py-12">
    <FileIcon className="w-12 h-12 text-gray-400 mx-auto" />
    <h3>No import jobs found</h3>
    <p>Start your first import to see it here</p>
  </div>
)}
```

## Визуальные индикаторы

### Статусы задач
- **pending**: Желтый badge с иконкой часов
- **processing**: Синий badge с анимированной иконкой
- **completed**: Зеленый badge с галочкой
- **failed**: Красный badge с крестиком

### Прогресс
- Показывается только для processing задач
- Синий progress bar с процентами
- Плавная анимация изменений

## Зависимости

### Redux
- `importSlice`: Управление состоянием
- Actions: `fetchImportJobs`, `fetchJobStatus`, `cancelImportJob`, `retryImportJob`

### Константы
- `IMPORT_STATUS_COLORS`: Цвета для статусов
- `IMPORT_STATUS_ICONS`: Иконки для статусов

## Особенности

### Производительность
- Обновляются только активные задачи
- Интервал очищается при отсутствии активных задач
- stopPropagation для предотвращения всплытия

### UX
- Hover эффект на строках
- Курсор pointer для кликабельных элементов
- Confirm dialog для отмены задач
- Анимация загрузки на кнопке refresh

### Адаптивность
- Flexbox layout для строк
- Отзывчивые размеры шрифтов
- Скрытие некритичных элементов на мобильных

## Проблемы и TODO

1. **Нет виртуализации**: Большие списки могут тормозить
2. **confirm() dialog**: Используется нативный confirm
3. **Отсутствие сортировки**: Нельзя отсортировать по дате/статусу
4. **Нет группировки**: Задачи не группируются по дням

## Рекомендации

1. Добавить виртуальный скроллинг для больших списков
2. Заменить confirm() на кастомный модальный диалог
3. Реализовать сортировку по колонкам
4. Добавить группировку задач по датам
5. Показывать время выполнения для завершенных задач