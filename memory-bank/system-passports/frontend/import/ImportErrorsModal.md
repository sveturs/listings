# Паспорт компонента: ImportErrorsModal

## Метаданные
- **Путь**: `/frontend/svetu/src/components/import/ImportErrorsModal.tsx`
- **Роль**: Модальное окно для просмотра ошибок импорта
- **Уровень**: Компонент-модал
- **Категория**: Import

## Описание

ImportErrorsModal отображает детальный список ошибок, возникших при импорте. Поддерживает поиск, фильтрацию по полям, пагинацию и экспорт ошибок в CSV формат.

## Props

Компонент не принимает props, работает через Redux состояние.

## Состояние

### Redux State (useAppSelector)
```typescript
{
  isErrorsModalOpen: boolean;   // Видимость модального окна
  currentJob: ImportJob | null; // Текущая задача с ошибками
}
```

### Local State
```typescript
{
  errors: ImportError[];        // Список ошибок
  isLoading: boolean;          // Состояние загрузки
  currentPage: number;         // Текущая страница (1-based)
  searchTerm: string;          // Поисковый запрос
  selectedField: string;       // Фильтр по полю ('all' или имя поля)
}
```

### Константы
```typescript
const errorsPerPage = 10;  // Количество ошибок на странице
```

## Основная логика

### Загрузка ошибок
```typescript
const loadErrors = useCallback(async () => {
  if (!currentJob) return;
  
  setIsLoading(true);
  try {
    const jobDetails = await ImportApi.getJobDetails(currentJob.id);
    setErrors(jobDetails.errors || []);
  } catch (error) {
    console.error('Failed to load import errors:', error);
  } finally {
    setIsLoading(false);
  }
}, [currentJob]);

useEffect(() => {
  if (isErrorsModalOpen && currentJob) {
    loadErrors();
  }
}, [isErrorsModalOpen, currentJob, loadErrors]);
```

### Фильтрация ошибок
```typescript
const filteredErrors = errors.filter(error => {
  const matchesSearch = 
    searchTerm === '' ||
    error.error_message.toLowerCase().includes(searchTerm.toLowerCase()) ||
    error.raw_data.toLowerCase().includes(searchTerm.toLowerCase());
  
  const matchesField = 
    selectedField === 'all' || 
    error.field_name === selectedField;
  
  return matchesSearch && matchesField;
});
```

### Пагинация
```typescript
const totalPages = Math.ceil(filteredErrors.length / errorsPerPage);
const startIndex = (currentPage - 1) * errorsPerPage;
const endIndex = startIndex + errorsPerPage;
const currentErrors = filteredErrors.slice(startIndex, endIndex);
```

### Экспорт ошибок
```typescript
const handleExportErrors = async () => {
  if (!currentJob) return;
  
  try {
    const blob = await ImportApi.exportResults(currentJob.id);
    downloadFile(blob, `import_errors_${currentJob.id}.csv`);
  } catch (error) {
    console.error('Failed to export errors:', error);
  }
};
```

## UI структура

### 1. Заголовок с действиями
```tsx
<div className="flex items-center justify-between p-6 border-b">
  <div>
    <h2>Import Errors</h2>
    <p>Errors from {currentJob.file_name}</p>
  </div>
  
  <div className="flex items-center space-x-3">
    {errors.length > 0 && (
      <button onClick={handleExportErrors}>
        <DownloadIcon />
        Export
      </button>
    )}
    
    <button onClick={handleClose}>
      <CloseIcon />
    </button>
  </div>
</div>
```

### 2. Панель фильтров
```tsx
<div className="p-6 border-b bg-gray-50">
  <div className="flex space-x-3">
    {/* Поиск */}
    <input
      type="text"
      placeholder="Search errors..."
      value={searchTerm}
      onChange={(e) => setSearchTerm(e.target.value)}
      className="flex-1"
    />
    
    {/* Фильтр по полю */}
    <select value={selectedField} onChange={(e) => setSelectedField(e.target.value)}>
      <option value="all">All fields</option>
      {uniqueFields.map(field => (
        <option key={field} value={field}>{field}</option>
      ))}
    </select>
    
    {/* Счетчик */}
    <span className="text-sm text-gray-600">
      {filteredErrors.length} of {errors.length}
    </span>
  </div>
</div>
```

### 3. Список ошибок
```tsx
{currentErrors.map((error, index) => (
  <div key={`${error.id}-${index}`} className="p-6">
    {/* Метки */}
    <div className="flex items-center space-x-2 mb-2">
      <span className="badge badge-error">
        Line {error.line_number}
      </span>
      <span className="badge">
        {error.field_name}
      </span>
    </div>
    
    {/* Сообщение об ошибке */}
    <p className="text-sm text-red-800 mb-3">
      {error.error_message}
    </p>
    
    {/* Исходные данные */}
    {error.raw_data && (
      <div className="bg-gray-50 rounded-md p-3">
        <p className="text-xs font-medium text-gray-500 mb-1">
          Raw data:
        </p>
        <p className="text-sm text-gray-800 font-mono break-all">
          {error.raw_data.length > 200
            ? `${error.raw_data.substring(0, 200)}...`
            : error.raw_data}
        </p>
      </div>
    )}
  </div>
))}
```

### 4. Состояния загрузки и пустые состояния
```tsx
{/* Загрузка */}
{isLoading && (
  <div className="flex items-center justify-center h-64">
    <SpinnerIcon className="animate-spin" />
    <span>Loading errors...</span>
  </div>
)}

{/* Нет ошибок */}
{errors.length === 0 && (
  <div className="text-center py-12">
    <CheckIcon className="w-12 h-12 text-green-400" />
    <h3>No errors found</h3>
    <p>All records were imported successfully</p>
  </div>
)}

{/* Нет результатов поиска */}
{currentErrors.length === 0 && (
  <div className="text-center py-12">
    <SearchIcon className="w-12 h-12 text-gray-400" />
    <h3>No matching errors</h3>
    <p>Try adjusting your search criteria</p>
  </div>
)}
```

### 5. Пагинация
```tsx
{totalPages > 1 && (
  <div className="px-6 py-4 border-t bg-gray-50">
    <div className="flex items-center justify-between">
      <div className="text-sm text-gray-700">
        Showing {startIndex + 1} to {Math.min(endIndex, filteredErrors.length)} 
        of {filteredErrors.length} entries
      </div>
      
      <div className="flex items-center space-x-2">
        <button
          onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
          disabled={currentPage === 1}
        >
          <ChevronLeftIcon />
          Previous
        </button>
        
        <span>{currentPage} / {totalPages}</span>
        
        <button
          onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
          disabled={currentPage === totalPages}
        >
          Next
          <ChevronRightIcon />
        </button>
      </div>
    </div>
  </div>
)}
```

## Типы ошибок

```typescript
interface ImportError {
  id: number;
  line_number: number;    // Номер строки в файле
  field_name: string;     // Поле с ошибкой
  error_message: string;  // Описание ошибки
  raw_data: string;       // Исходные данные строки
}
```

## Визуальные элементы

### Badges
- **Line number**: Красный badge (badge-error)
- **Field name**: Серый badge по умолчанию

### Layout
- Максимальная ширина 6xl (1152px)
- Flex column для правильного скроллинга
- Фиксированные header и footer

## Зависимости

### Redux
- `importSlice`: Управление состоянием
- Actions: `setErrorsModalOpen`

### API
- `ImportApi.getJobDetails()`: Загрузка деталей с ошибками
- `ImportApi.exportResults()`: Экспорт в CSV
- `downloadFile()`: Утилита для скачивания файла

## Особенности

### Производительность
- Пагинация для больших списков ошибок
- Обрезка длинных raw_data до 200 символов
- useCallback для оптимизации loadErrors

### UX
- Поиск по тексту ошибки и исходным данным
- Фильтрация по полям
- Экспорт всех ошибок одним кликом
- Отзывчивый дизайн для мобильных

### Доступность
- Disabled состояния для кнопок пагинации
- Описательные тексты для пустых состояний
- Контрастные цвета для ошибок

## Проблемы и TODO

1. **console.error**: Ошибки логируются только в консоль
2. **Нет сортировки**: Невозможно отсортировать ошибки
3. **Ограничение raw_data**: Обрезка может скрыть важные данные
4. **Нет группировки**: Ошибки не группируются по типам

## Рекомендации

1. Добавить toast-уведомления для ошибок загрузки
2. Реализовать сортировку по номеру строки/полю
3. Добавить возможность просмотра полных raw_data
4. Группировать похожие ошибки
5. Показывать рекомендации по исправлению типичных ошибок