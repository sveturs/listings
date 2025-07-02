# Паспорт компонента: ImportWizard

## Метаданные
- **Путь**: `/frontend/svetu/src/components/import/ImportWizard.tsx`
- **Роль**: Пошаговый мастер импорта товаров из файлов или URL
- **Уровень**: Компонент страницы
- **Категория**: Import

## Описание

ImportWizard - главный компонент системы импорта товаров. Предоставляет интерфейс для загрузки файлов через drag-and-drop, импорта по URL, валидации данных и отображения прогресса загрузки. Поддерживает форматы CSV, XML и ZIP.

## Props

```typescript
interface ImportWizardProps {
  storefrontId: number;      // ID витрины для импорта
  storefrontSlug?: string;   // Slug витрины (для slug-based API)
  onSuccess?: (jobId: number) => void;  // Callback после успешного импорта
  onClose?: () => void;      // Callback при закрытии
}
```

## Состояние

### Redux State (useAppSelector)
```typescript
{
  isImportModalOpen: boolean;      // Видимость модального окна
  selectedFiles: File[];           // Выбранные файлы
  importUrl: string;               // URL для импорта
  selectedFileType: string;        // Тип файла (csv/xml/zip)
  updateMode: 'upsert' | 'create_only' | 'update_only';
  categoryMappingMode: 'auto' | 'manual' | 'skip';
  isLoading: boolean;              // Загрузка
  isUploading: boolean;            // Процесс загрузки файла
  uploadProgress: UploadProgress;  // Прогресс загрузки
  error: string | null;            // Ошибка
  formats: ImportFormats;          // Поддерживаемые форматы
  validationErrors: ValidationError[];  // Ошибки валидации
}
```

### Local State
```typescript
{
  activeTab: 'file' | 'url';    // Активная вкладка
  dragActive: boolean;          // Активен drag-and-drop
  validationStep: boolean;      // Шаг валидации
}
```

## Основная логика

### Загрузка файлов
```typescript
// Обработка drag-and-drop
const handleDrop = (e: React.DragEvent) => {
  e.preventDefault();
  const files = Array.from(e.dataTransfer.files);
  handleFileSelection(files);
};

// Валидация и выбор файла
const handleFileSelection = (files: File[]) => {
  const file = files[0];
  
  // Валидация типа файла
  if (!validateFileType(file, IMPORT_FILE_CONFIG.allowedTypes)) {
    alert(t('errors.invalidFileType'));
    return;
  }
  
  // Валидация размера
  if (!validateFileSize(file, IMPORT_FILE_CONFIG.maxFileSize)) {
    alert(t('errors.fileTooLarge'));
    return;
  }
  
  // Автоопределение типа
  const detectedType = getFileTypeFromExtension(file.name);
  if (detectedType) {
    dispatch(setSelectedFileType(detectedType));
  }
  
  dispatch(setSelectedFiles([file]));
};
```

### Импорт из файла
```typescript
const handleImportFromFile = async () => {
  const result = await dispatch(
    importFromFile({
      storefrontId,
      storefrontSlug,
      file: selectedFiles[0],
      options: {
        file_type: selectedFileType,
        update_mode: updateMode,
        category_mapping_mode: categoryMappingMode,
      },
    })
  ).unwrap();
  
  onSuccess?.(result.id);
  handleClose();
};
```

### Импорт из URL
```typescript
const handleImportFromUrl = async () => {
  const result = await dispatch(
    importFromUrl({
      storefrontId,
      storefrontSlug,
      request: {
        file_url: importUrl,
        file_type: selectedFileType,
        update_mode: updateMode,
        category_mapping_mode: categoryMappingMode,
      },
    })
  ).unwrap();
  
  onSuccess?.(result.id);
  handleClose();
};
```

## UI/UX

### Модальное окно
- Полноэкранный overlay с затемнением
- Максимальная ширина 4xl (896px)
- Скроллинг при переполнении

### Вкладки
- **Upload File**: Drag-and-drop зона, выбор файла
- **Import from URL**: Поле ввода URL

### Drag-and-Drop зона
- Визуальная индикация при перетаскивании
- Отображение выбранного файла с размером
- Кнопка удаления файла

### Опции импорта
- **Update Mode**: 
  - upsert (создать или обновить)
  - create_only (только создать)
  - update_only (только обновить)
- **Category Mapping**:
  - auto (автоматическое сопоставление)
  - manual (ручное сопоставление)
  - skip (пропустить)

### Прогресс загрузки
- Процентный индикатор
- Анимированный progress bar

### Валидация
- Отдельный шаг после выбора файла
- Отображение предупреждений
- Возможность продолжить с предупреждениями

## Зависимости

### Redux
- `importSlice`: Управление состоянием импорта
- Actions: `setImportModalOpen`, `setSelectedFiles`, `importFromFile`, `importFromUrl`

### API
- `ImportApi`: HTTP-запросы для импорта
- `validateFileType`, `validateFileSize`: Утилиты валидации

### Компоненты
- Модальное окно с вкладками
- Drag-and-drop зона
- Progress bar
- Форма с select-полями

## Пример использования

```tsx
// В компоненте управления витриной
function StorefrontProducts() {
  const [showImport, setShowImport] = useState(false);
  
  const handleImportSuccess = (jobId: number) => {
    console.log('Import started with job ID:', jobId);
    // Обновить список задач импорта
  };
  
  return (
    <>
      <button onClick={() => setShowImport(true)}>
        Import Products
      </button>
      
      {showImport && (
        <ImportWizard
          storefrontId={123}
          storefrontSlug="my-store"
          onSuccess={handleImportSuccess}
          onClose={() => setShowImport(false)}
        />
      )}
    </>
  );
}
```

## Особенности

### Поддержка форматов
- CSV - текстовый формат с разделителями
- XML - структурированный формат
- ZIP - архив с несколькими файлами

### Ограничения
- Максимальный размер файла из конфигурации
- Поддержка только одного файла за раз
- Валидация перед импортом опциональна

### Безопасность
- Валидация типов файлов на клиенте
- Проверка размера файла
- CSRF-защита через Redux

### Производительность
- Отслеживание прогресса загрузки
- Отложенная загрузка форматов
- Очистка URL объектов после использования

## Проблемы и TODO

1. **alert() вместо toast**: Используются нативные alert для ошибок
2. **Отсутствие retry**: Нет возможности повторить неудачную загрузку
3. **Нет preview**: Отсутствует предпросмотр данных перед импортом
4. **Жесткая локализация**: Все тексты через t() без fallback

## Рекомендации

1. Заменить alert() на toast-уведомления
2. Добавить предпросмотр первых строк файла
3. Реализовать пакетную загрузку файлов
4. Добавить историю импортов в UI
5. Улучшить обработку ошибок с retry-логикой