# Архитектура редактирования объявлений - Полная спецификация

## 📋 Оглавление
1. [Обзор системы](#обзор-системы)
2. [База данных](#база-данных)
3. [Backend API](#backend-api)
4. [Frontend архитектура](#frontend-архитектура)
5. [Работа с изображениями](#работа-с-изображениями)
6. [Система атрибутов](#система-атрибутов)
7. [Поисковая индексация](#поисковая-индексация)
8. [План реализации продакшн версии](#план-реализации)

## 🎯 Обзор системы

### Текущее состояние
- **Путь**: `/profile/listings/[id]/edit`
- **Файл**: `frontend/svetu/src/app/[locale]/profile/listings/[id]/edit/page.tsx`
- **Статус**: Базовая версия без поддержки атрибутов, фото и расширенных функций

### Целевое состояние
Современная система редактирования с:
- Полной поддержкой всех полей объявления
- Управление фотографиями (загрузка, удаление, сортировка)
- Динамические атрибуты по категориям
- Умные подсказки и валидация
- SEO оптимизация
- Предпросмотр в реальном времени

## 💾 База данных

### Основные таблицы

#### marketplace_listings
```sql
- id: integer (PK)
- user_id: integer (FK -> users)
- category_id: integer (FK -> marketplace_categories)
- title: varchar(255)
- description: text
- price: numeric(12,2)
- condition: varchar(50) -- 'new', 'used', 'refurbished'
- status: varchar(20) -- 'active', 'inactive', 'sold'
- location: varchar(255)
- latitude: numeric(10,8)
- longitude: numeric(11,8)
- address_city: varchar(100)
- address_country: varchar(100)
- views_count: integer
- show_on_map: boolean
- original_language: varchar(10)
- storefront_id: integer (FK -> user_storefronts, nullable)
- metadata: jsonb
- created_at: timestamp
- updated_at: timestamp
```

#### listing_attribute_values
```sql
- id: integer (PK)
- listing_id: integer (FK -> marketplace_listings)
- attribute_id: integer (FK -> category_attributes)
- text_value: text
- numeric_value: numeric(20,5)
- boolean_value: boolean
- json_value: jsonb
- unit: varchar(20)
- value_type: varchar(20) -- 'text', 'numeric', 'boolean', 'json'
```

#### category_attributes
```sql
- id: integer (PK)
- name: varchar(100) -- системное имя (brand, model, year)
- display_name: varchar(255) -- отображаемое имя
- attribute_type: varchar(50) -- 'text', 'numeric', 'select', 'multiselect', 'boolean'
- options: jsonb -- для select/multiselect
- validation_rules: jsonb
- is_searchable: boolean
- is_filterable: boolean
- is_required: boolean
- sort_order: integer
```

#### marketplace_images
```sql
- id: integer (PK)
- listing_id: integer (FK -> marketplace_listings)
- file_path: varchar(255)
- file_name: varchar(255)
- file_size: integer
- content_type: varchar(100)
- is_main: boolean
- created_at: timestamp
- storage_type: varchar(50) -- 'minio'
- storage_bucket: varchar(100) -- 'listings'
- public_url: varchar(500)
```

## 🔌 Backend API

### Основные эндпоинты

#### GET /api/v1/marketplace/listings/{id}
Получение полной информации об объявлении включая:
- Основные данные
- Категорию с путем (category_path)
- Атрибуты (attributes)
- Изображения (images)
- Переводы (translations)

#### PUT /api/v1/marketplace/listings/{id}
```typescript
interface UpdateListingRequest {
  title: string;
  description: string;
  price: number;
  condition: 'new' | 'used' | 'refurbished';
  city?: string;
  country?: string;
  location?: string;
  latitude?: number;
  longitude?: number;
  show_on_map?: boolean;
  category_id: number;
  attributes?: Array<{
    attribute_id: number;
    value: string | number | boolean;
  }>;
}
```

#### POST /api/v1/marketplace/listings/{id}/images
- Multipart form-data
- Поле: `file` (множественная загрузка)
- Поле: `main_image_index` (индекс главного фото)

#### DELETE /api/v1/marketplace/listings/{id}/images/{imageId}
Удаление конкретного изображения

#### PUT /api/v1/marketplace/listings/{id}/images/reorder
```typescript
interface ReorderImagesRequest {
  image_ids: number[]; // новый порядок ID изображений
  main_image_id: number; // ID главного изображения
}
```

## 🎨 Frontend архитектура

### Компонентная структура

```
/profile/listings/[id]/edit/
├── page.tsx                    # Главная страница
├── components/
│   ├── BasicInfoSection.tsx    # Основная информация
│   ├── ImagesSection.tsx       # Управление фото
│   ├── AttributesSection.tsx   # Динамические атрибуты
│   ├── LocationSection.tsx     # Локация и карта
│   ├── SEOSection.tsx          # SEO и переводы
│   └── PreviewCard.tsx         # Превью объявления
```

### Ключевые хуки и сервисы

```typescript
// Хук для работы с объявлением
const useListingEdit = (listingId: string) => {
  const [listing, setListing] = useState<Listing | null>(null);
  const [isDirty, setIsDirty] = useState(false);
  const [errors, setErrors] = useState<ValidationErrors>({});
  
  // Загрузка, сохранение, валидация
  return { listing, save, validate, isDirty, errors };
};

// Хук для работы с атрибутами категории
const useCategoryAttributes = (categoryId: number) => {
  const [attributes, setAttributes] = useState<CategoryAttribute[]>([]);
  const [values, setValues] = useState<AttributeValues>({});
  
  return { attributes, values, updateValue };
};
```

## 📸 Работа с изображениями

### MinIO конфигурация
- **Bucket**: listings
- **Путь**: `{listing_id}/{timestamp}_{random}.{ext}`
- **Публичный URL**: `/listings/{path}`

### Процесс загрузки
1. Клиент отправляет файлы через multipart/form-data
2. Backend валидирует (тип, размер < 10MB)
3. Генерирует уникальное имя с timestamp
4. Загружает в MinIO
5. Сохраняет метаданные в БД
6. Возвращает публичные URL

### Оптимизация изображений
- Автоматическое сжатие JPEG (качество 85%)
- Генерация превью (300x300)
- Поддержка WebP для современных браузеров

## 🏷️ Система атрибутов

### Загрузка атрибутов категории
```sql
SELECT ca.*, cam.is_required as category_required
FROM category_attributes ca
JOIN category_attribute_mapping cam ON ca.id = cam.attribute_id
WHERE cam.category_id = $1
ORDER BY cam.sort_order, ca.sort_order;
```

### Типы атрибутов

#### text
```typescript
<input type="text" value={value} onChange={onChange} />
```

#### numeric
```typescript
<input type="number" value={value} onChange={onChange} />
{attribute.unit && <span>{attribute.unit}</span>}
```

#### select
```typescript
<select value={value} onChange={onChange}>
  {attribute.options.map(opt => (
    <option key={opt.value} value={opt.value}>{opt.label}</option>
  ))}
</select>
```

#### multiselect
```typescript
<MultiSelect 
  options={attribute.options}
  value={value as string[]}
  onChange={onChange}
/>
```

#### boolean
```typescript
<input type="checkbox" checked={value} onChange={onChange} />
```

## 🔍 Поисковая индексация

### Автоматическая переиндексация при:
- Изменении title, description, price
- Изменении атрибутов
- Изменении статуса (active/inactive)
- Добавлении/удалении фото

### Индексируемые поля
```typescript
interface SearchDocument {
  id: number;
  title: string;
  title_ngram: string; // для частичного поиска
  description: string;
  price: number;
  category_id: number;
  category_path: string[]; // для фасетного поиска
  attributes: {
    [key: string]: string | number | boolean;
  };
  location: {
    lat: number;
    lon: number;
  };
  city: string;
  country: string;
  images_count: number;
  has_images: boolean;
  created_at: string;
  updated_at: string;
}
```

## 📝 План реализации продакшн версии

### Фаза 1: Базовый функционал (2-3 дня)
1. **Рефакторинг текущей страницы**
   - Разделение на компоненты
   - Добавление валидации
   - Улучшение UX (loading states, error handling)

2. **Интеграция с изображениями**
   - Компонент загрузки фото
   - Drag & drop сортировка
   - Выбор главного фото
   - Удаление фото

3. **Поддержка атрибутов**
   - Загрузка атрибутов категории
   - Динамическая генерация полей
   - Валидация по правилам

### Фаза 2: Расширенный функционал (2-3 дня)
1. **Локация и карта**
   - Интеграция с LocationPicker
   - Геокодинг адреса
   - Опция "Скрыть на карте"

2. **SEO и оптимизация**
   - Счетчики символов
   - SEO-скор заголовка
   - Предложения по улучшению

3. **Предпросмотр**
   - Реалтайм превью карточки
   - Переключение вида (список/сетка)
   - Мобильный вид

### Фаза 3: Умные функции (3-4 дня)
1. **Анализ рынка**
   - Средняя цена для категории
   - Похожие объявления
   - Рекомендации по цене

2. **AI-помощник** (опционально)
   - Улучшение описания
   - Генерация заголовка
   - Перевод на другие языки

3. **История изменений**
   - Версионирование
   - Откат изменений
   - Сравнение версий

## 🚀 Быстрый старт для следующей сессии

```bash
# 1. Открыть основной файл редактирования
code /data/hostel-booking-system/frontend/svetu/src/app/[locale]/profile/listings/[id]/edit/page.tsx

# 2. Посмотреть примеры UI
http://localhost:3001/ru/examples/listing-edit-ux

# 3. Проверить API документацию
cd /data/hostel-booking-system/backend/docs && python3 -m http.server 8888
# Затем использовать JSON MCP для поиска PUT /api/v1/marketplace/listings/{id}

# 4. Начать с создания компонентов в папке
mkdir -p /data/hostel-booking-system/frontend/svetu/src/app/[locale]/profile/listings/[id]/edit/components
```

## ⚡ Критически важные моменты

1. **Проверка владельца** - обязательна на backend и frontend
2. **Валидация атрибутов** - использовать validation_rules из БД
3. **Оптимистичные обновления** - показывать изменения сразу, откатывать при ошибке
4. **Автосохранение** - сохранять черновики каждые 30 секунд
5. **Обработка ошибок** - использовать toast уведомления
6. **Мобильная версия** - адаптивный дизайн обязателен

## 📚 Дополнительные ресурсы

- **Текущая страница редактирования**: `/frontend/svetu/src/app/[locale]/profile/listings/[id]/edit/page.tsx`
- **API handler**: `/backend/internal/proj/marketplace/handler/listings.go` (функция UpdateListing)
- **Примеры UI**: `/frontend/svetu/src/app/[locale]/examples/listing-edit-ux/page.tsx`
- **Компонент LocationPicker**: см. документацию `LOCATION_PICKER_INSTRUCTION.md`
- **Система атрибутов**: см. документацию `CATEGORY_ATTRIBUTES_STATUS.md`