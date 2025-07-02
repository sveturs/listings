# Паспорт компонента: ImportJobDetails

## Метаданные
- **Путь**: `/frontend/svetu/src/components/import/ImportJobDetails.tsx`
- **Роль**: Модальное окно с детальной информацией о задаче импорта
- **Уровень**: Компонент-модал
- **Категория**: Import

## Описание

ImportJobDetails отображает полную информацию о выбранной задаче импорта: статус, прогресс, статистику, временные метки и ошибки. Предоставляет действия для управления задачей.

## Props

Компонент не принимает props, работает через Redux состояние.

## Состояние

### Redux State (useAppSelector)
```typescript
{
  isJobDetailsModalOpen: boolean;  // Видимость модального окна
  currentJob: ImportJob | null;    // Текущая выбранная задача
  isLoading: boolean;              // Состояние загрузки
}
```

## Основная логика

### Загрузка деталей
```typescript
useEffect(() => {
  if (isJobDetailsModalOpen && currentJob) {
    dispatch(fetchJobDetails(currentJob.id));
  }
}, [isJobDetailsModalOpen, currentJob, dispatch]);
```

### Обработчики действий
```typescript
// Закрытие модального окна
const handleClose = () => {
  dispatch(setJobDetailsModalOpen(false));
};

// Показать ошибки
const handleShowErrors = () => {
  dispatch(setErrorsModalOpen(true));
};

// Отменить задачу
const handleCancelJob = async () => {
  if (currentJob && confirm(t('actions.confirmCancel'))) {
    await dispatch(cancelImportJob(currentJob.id));
    handleClose();
  }
};

// Повторить задачу
const handleRetryJob = async () => {
  if (currentJob) {
    await dispatch(retryImportJob(currentJob.id));
    handleClose();
  }
};
```

### Вычисления
```typescript
// Форматирование даты
const formatDate = (dateString?: string) => {
  if (!dateString) return t('jobs.details.notSet');
  return new Date(dateString).toLocaleString();
};

// Расчет прогресса
const calculateProgress = () => {
  if (currentJob.total_records === 0) return 0;
  return Math.round(
    (currentJob.processed_records / currentJob.total_records) * 100
  );
};

// Время выполнения
const getProcessingTime = () => {
  if (!currentJob.started_at) return null;
  
  const startTime = new Date(currentJob.started_at);
  const endTime = currentJob.completed_at 
    ? new Date(currentJob.completed_at) 
    : new Date();
  
  const diffMs = endTime.getTime() - startTime.getTime();
  const diffSeconds = Math.floor(diffMs / 1000);
  const diffMinutes = Math.floor(diffSeconds / 60);
  const diffHours = Math.floor(diffMinutes / 60);
  
  if (diffHours > 0) {
    return `${diffHours}h ${diffMinutes % 60}m ${diffSeconds % 60}s`;
  } else if (diffMinutes > 0) {
    return `${diffMinutes}m ${diffSeconds % 60}s`;
  } else {
    return `${diffSeconds}s`;
  }
};
```

## UI структура

### 1. Заголовок
```tsx
<div className="flex items-center justify-between p-6 border-b">
  <h2>Import Job Details</h2>
  <button onClick={handleClose}>
    <CloseIcon />
  </button>
</div>
```

### 2. Статус и прогресс
```tsx
<div className="bg-gray-50 rounded-lg p-4">
  <span className={IMPORT_STATUS_COLORS[currentJob.status]}>
    {getStatusIcon(currentJob.status)}
    {t(`status.${currentJob.status}`)}
  </span>
  
  {currentJob.status === 'processing' && (
    <>
      <span>{calculateProgress()}% complete</span>
      <div className="w-full bg-gray-200 rounded-full h-3">
        <div 
          className="bg-blue-600 h-3 rounded-full"
          style={{ width: `${calculateProgress()}%` }}
        />
      </div>
    </>
  )}
</div>
```

### 3. Информация о файле
```tsx
<dl className="grid grid-cols-2 gap-4">
  <div>
    <dt>File Name</dt>
    <dd>{currentJob.file_name}</dd>
  </div>
  <div>
    <dt>File Type</dt>
    <dd>{currentJob.file_type.toUpperCase()}</dd>
  </div>
  {currentJob.file_url && (
    <div className="col-span-2">
      <dt>File URL</dt>
      <dd>
        <a href={currentJob.file_url} target="_blank">
          {currentJob.file_url}
        </a>
      </dd>
    </div>
  )}
</dl>
```

### 4. Статистика обработки
```tsx
<div className="grid grid-cols-4 gap-4">
  <div className="bg-blue-50 rounded-lg p-3 text-center">
    <div className="text-2xl font-bold text-blue-600">
      {currentJob.total_records}
    </div>
    <div className="text-sm">Total</div>
  </div>
  
  <div className="bg-yellow-50 rounded-lg p-3 text-center">
    <div className="text-2xl font-bold text-yellow-600">
      {currentJob.processed_records}
    </div>
    <div className="text-sm">Processed</div>
  </div>
  
  <div className="bg-green-50 rounded-lg p-3 text-center">
    <div className="text-2xl font-bold text-green-600">
      {currentJob.successful_records}
    </div>
    <div className="text-sm">Successful</div>
  </div>
  
  <div className="bg-red-50 rounded-lg p-3 text-center">
    <div className="text-2xl font-bold text-red-600">
      {currentJob.failed_records}
    </div>
    <div className="text-sm">Failed</div>
  </div>
</div>
```

### 5. Временные метки
```tsx
<dl className="grid grid-cols-2 gap-4">
  <div>
    <dt>Created</dt>
    <dd>{formatDate(currentJob.created_at)}</dd>
  </div>
  <div>
    <dt>Started</dt>
    <dd>{formatDate(currentJob.started_at)}</dd>
  </div>
  <div>
    <dt>Completed</dt>
    <dd>{formatDate(currentJob.completed_at)}</dd>
  </div>
  {getProcessingTime() && (
    <div>
      <dt>Duration</dt>
      <dd>{getProcessingTime()}</dd>
    </div>
  )}
</dl>
```

### 6. Сообщение об ошибке
```tsx
{currentJob.status === 'failed' && currentJob.error_message && (
  <div>
    <h3>Error Details</h3>
    <div className="bg-red-50 border border-red-200 rounded-md p-4">
      <p className="text-sm text-red-800">
        {currentJob.error_message}
      </p>
    </div>
  </div>
)}
```

### 7. Действия
```tsx
<div className="flex items-center justify-between pt-6 border-t">
  <div className="flex space-x-3">
    {currentJob.failed_records > 0 && (
      <button onClick={handleShowErrors}>
        View Errors
      </button>
    )}
    
    {currentJob.status === 'processing' && (
      <button onClick={handleCancelJob}>
        Cancel
      </button>
    )}
    
    {currentJob.status === 'failed' && (
      <button onClick={handleRetryJob}>
        Retry
      </button>
    )}
  </div>
  
  <button onClick={handleClose}>
    Close
  </button>
</div>
```

## Визуальные элементы

### Цветовая схема статистики
- **Total**: Синий фон (blue-50)
- **Processed**: Желтый фон (yellow-50)
- **Successful**: Зеленый фон (green-50)
- **Failed**: Красный фон (red-50)

### Модальное окно
- Фиксированное позиционирование
- Затемненный фон (black/50)
- Максимальная ширина 2xl (672px)
- Скроллинг при переполнении

## Зависимости

### Redux
- `importSlice`: Управление состоянием
- Actions: `setJobDetailsModalOpen`, `fetchJobDetails`, `cancelImportJob`, `retryImportJob`

### Константы
- `IMPORT_STATUS_COLORS`: Цвета для статусов
- `IMPORT_STATUS_ICONS`: Иконки для статусов

## Особенности

### Условное отображение
- URL файла показывается только если есть
- Время выполнения только для начатых задач
- Действия зависят от статуса задачи

### Форматирование времени
- Автоматический выбор формата (часы/минуты/секунды)
- Локализованное отображение дат
- "Not set" для отсутствующих дат

## Проблемы и TODO

1. **confirm() dialog**: Используется нативный confirm
2. **Нет копирования**: Невозможно скопировать детали
3. **Отсутствие логов**: Не показываются логи обработки
4. **Нет экспорта**: Невозможно экспортировать детали

## Рекомендации

1. Заменить confirm() на кастомный диалог
2. Добавить кнопки копирования для важных данных
3. Показывать логи обработки в отдельной вкладке
4. Добавить экспорт деталей в JSON/CSV
5. Реализовать live-обновление для активных задач