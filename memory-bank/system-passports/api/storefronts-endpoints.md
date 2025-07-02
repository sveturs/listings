# Паспорт API Endpoints: Storefronts (Витрины)

## 📋 Метаданные
- **Группа API**: Storefronts
- **Базовый путь**: `/api/v1/storefronts`
- **Handler**: `backend/internal/proj/storefronts/module.go`
- **Количество endpoints**: 43 (12 публичных, 31 защищенный)
- **Интеграции**: PostgreSQL, OpenSearch, MinIO, Redis

## 🎯 Назначение
Система витрин для бизнесов на платформе:
- Публичные страницы витрин с геолокацией
- Управление товарами и инвентарем
- Bulk операции для массового управления
- Импорт товаров из внешних источников
- Аналитика продаж и просмотров
- Карта витрин с геопространственным поиском

## 📡 Endpoints

### 🌐 Публичные (без авторизации)

#### Основная информация о витринах
```typescript
GET /api/v1/storefronts/
// Список всех публичных витрин с пагинацией
// Handler: m.storefrontHandler.ListStorefronts
// Query: page, limit, category, city, status=active
// Response: PublicStorefront[] с базовой информацией
```

```typescript
GET /api/v1/storefronts/search
// Поиск витрин через OpenSearch
// Handler: m.storefrontHandler.SearchOpenSearch
// Query: q, location, category, radius, open_now
// Response: SearchResults с геофильтрацией и релевантностью
```

```typescript
GET /api/v1/storefronts/nearby
// Ближайшие витрины по геолокации
// Handler: m.storefrontHandler.GetNearbyStorefronts
// Query: lat, lng, radius (км), limit
// Response: NearbyStorefront[] с расстояниями
```

```typescript
GET /api/v1/storefronts/map
// Данные для карты витрин (кластеризация)
// Handler: m.storefrontHandler.GetMapData
// Query: bounds (ne_lat, ne_lng, sw_lat, sw_lng), zoom
// Response: MapCluster[] для производительного отображения
```

```typescript
GET /api/v1/storefronts/building
// Все бизнесы в одном здании
// Handler: m.storefrontHandler.GetBusinessesInBuilding
// Query: building_id или lat,lng с точностью до метра
// Response: BusinessInBuilding[] с этажами/номерами
```

#### Просмотр витрин
```typescript
GET /api/v1/storefronts/slug/:slug
// Публичная страница витрины по slug
// Handler: m.storefrontHandler.GetStorefrontBySlug
// Response: DetailedStorefront с полной информацией
// Includes: часы работы, контакты, рейтинги, фото
```

```typescript
GET /api/v1/storefronts/:id
// Витрина по ID (для внутренних ссылок)
// Handler: m.storefrontHandler.GetStorefront
// Response: DetailedStorefront
```

```typescript
GET /api/v1/storefronts/:id/staff
// Персонал витрины (публичная информация)
// Handler: m.storefrontHandler.GetStaff
// Response: StaffMember[] с ролями и фото
```

```typescript
POST /api/v1/storefronts/:id/view
// Запись просмотра витрины (аналитика)
// Handler: m.storefrontHandler.RecordView
// Body: {visitor_id?, source?, utm_params?}
// Effect: Увеличение счетчика просмотров
```

#### Продукты витрин (публичное API)
```typescript
GET /api/v1/storefronts/slug/:slug/products
// Каталог товаров витрины
// Handler: m.getProductsBySlug
// Query: category, price_range, sort, in_stock
// Response: ProductCatalog с фильтрами и пагинацией
```

```typescript
GET /api/v1/storefronts/slug/:slug/products/:id
// Детальная информация о товаре
// Handler: m.getProductBySlug
// Response: DetailedProduct с характеристиками и фото
```

### 🔒 Защищенные (требуют авторизации)

#### Управление витринами
```typescript
GET /api/v1/storefronts/my
// Мои витрины (владелец или сотрудник)
// Handler: m.storefrontHandler.GetMyStorefronts
// Response: MyStorefront[] с правами доступа
```

```typescript
POST /api/v1/storefronts/
// Создание новой витрины
// Handler: m.storefrontHandler.CreateStorefront
// Body: CreateStorefrontRequest
// Response: CreatedStorefront с ID и slug
```

```typescript
PUT /api/v1/storefronts/:id
// Обновление витрины
// Handler: m.storefrontHandler.UpdateStorefront
// Security: Только владелец или админ сотрудник
// Body: UpdateStorefrontRequest (partial)
```

```typescript
DELETE /api/v1/storefronts/:id
// Удаление витрины (soft delete)
// Handler: m.storefrontHandler.DeleteStorefront
// Security: Только владелец
// Effect: Статус меняется на 'deleted'
```

#### Управление товарами
```typescript
POST /api/v1/storefronts/slug/:slug/products
// Создание нового товара
// Handler: m.createProductBySlug
// Body: CreateProductRequest
// Response: CreatedProduct с ID
```

```typescript
PUT /api/v1/storefronts/slug/:slug/products/:id
// Обновление товара
// Handler: m.updateProductBySlug
// Body: UpdateProductRequest (partial)
// Effect: Обновление в БД + переиндексация
```

```typescript
DELETE /api/v1/storefronts/slug/:slug/products/:id
// Удаление товара
// Handler: m.deleteProductBySlug
// Effect: Soft delete или архивация
```

```typescript
POST /api/v1/storefronts/slug/:slug/products/:id/inventory
// Обновление остатков товара
// Handler: m.updateInventoryBySlug
// Body: {stock_quantity: number, low_stock_threshold?: number}
// Real-time: WebSocket уведомления при низких остатках
```

```typescript
GET /api/v1/storefronts/slug/:slug/products/stats
// Статистика по товарам витрины
// Handler: m.getProductStatsBySlug
// Response: ProductStats с топ товарами, оборотом, остатками
```

#### Bulk операции
```typescript
POST /api/v1/storefronts/slug/:slug/products/bulk/create
// Массовое создание товаров
// Handler: m.bulkCreateProductsBySlug
// Body: CreateProductRequest[]
// Response: BulkOperationResult с успехами/ошибками
```

```typescript
PUT /api/v1/storefronts/slug/:slug/products/bulk/update
// Массовое обновление
// Handler: m.bulkUpdateProductsBySlug
// Body: {product_ids: string[], updates: ProductUpdates}
// Use-case: Изменение цен, категорий, статусов
```

```typescript
DELETE /api/v1/storefronts/slug/:slug/products/bulk/delete
// Массовое удаление
// Handler: m.bulkDeleteProductsBySlug
// Body: {product_ids: string[], reason?: string}
```

```typescript
PUT /api/v1/storefronts/slug/:slug/products/bulk/status
// Массовое изменение статусов
// Handler: m.bulkUpdateStatusBySlug
// Body: {product_ids: string[], status: ProductStatus}
// Use-case: Активация/деактивация, снятие с продажи
```

#### Система импорта
```typescript
POST /api/v1/storefronts/:id/import/url
// Импорт товаров из URL (API/каталог)
// Handler: m.importHandler.ImportFromURL
// Body: {url: string, mapping: FieldMapping, options: ImportOptions}
// Response: ImportJob с ID для отслеживания
```

```typescript
POST /api/v1/storefronts/:id/import/file
// Импорт из файла (CSV/Excel)
// Handler: m.importHandler.ImportFromFile
// Content-Type: multipart/form-data
// Files: CSV/XLSX с товарами
```

```typescript
POST /api/v1/storefronts/:id/import/validate
// Валидация файла перед импортом
// Handler: m.importHandler.ValidateImportFile
// Response: ValidationResult с ошибками и предпросмотром
```

```typescript
GET /api/v1/storefronts/:id/import/jobs
// Список задач импорта
// Handler: m.importHandler.GetJobs
// Response: ImportJob[] с прогрессом и статусами
```

```typescript
GET /api/v1/storefronts/:id/import/jobs/:jobId
// Детали конкретной задачи импорта
// Handler: m.importHandler.GetJobDetails
// Response: DetailedImportJob с логами и ошибками
```

```typescript
GET /api/v1/storefronts/:id/import/jobs/:jobId/status
// Статус задачи в реальном времени
// Handler: m.importHandler.GetJobStatus
// Response: JobStatus с прогрессом %
```

```typescript
POST /api/v1/storefronts/:id/import/jobs/:jobId/cancel
// Отмена импорта
// Handler: m.importHandler.CancelJob
// Effect: Остановка задачи + откат изменений
```

```typescript
POST /api/v1/storefronts/:id/import/jobs/:jobId/retry
// Повторный запуск неудачного импорта
// Handler: m.importHandler.RetryJob
// Body: {retry_failed_only: boolean}
```

## 🏗️ Структуры данных

### Основные модели
```typescript
interface Storefront {
  id: string;
  slug: string;                      // уникальный URL slug
  name: string;
  description: string;
  category: string;
  owner_id: string;
  
  // Локация
  location: {
    address: string;
    city: string;
    country: string;
    coordinates: [number, number];
    building_id?: string;
    floor?: string;
    unit?: string;
  };
  
  // Контакты
  contacts: {
    phone?: string;
    email?: string;
    website?: string;
    social_media?: SocialMediaLinks;
  };
  
  // Режим работы
  business_hours: BusinessHours;
  
  // Медиа
  logo_url?: string;
  cover_images: string[];
  
  // Настройки
  settings: StorefrontSettings;
  
  // Статистика
  stats: StorefrontStats;
  
  status: "draft" | "active" | "paused" | "closed" | "deleted";
  created_at: string;
  updated_at: string;
}

interface BusinessHours {
  monday: DaySchedule;
  tuesday: DaySchedule;
  wednesday: DaySchedule;
  thursday: DaySchedule;
  friday: DaySchedule;
  saturday: DaySchedule;
  sunday: DaySchedule;
  timezone: string;
  special_hours?: SpecialHour[];     // праздники, отпуска
}

interface DaySchedule {
  is_open: boolean;
  open_time?: string;                // "09:00"
  close_time?: string;               // "18:00"
  breaks?: TimeRange[];              // обеденные перерывы
}

interface Product {
  id: string;
  storefront_id: string;
  name: string;
  description: string;
  category: string;
  sku?: string;
  
  // Цена и инвентарь
  price: number;
  currency: "RSD" | "EUR";
  stock_quantity: number;
  low_stock_threshold: number;
  
  // Медиа
  images: ProductImage[];
  
  // Характеристики
  attributes: Record<string, any>;
  variations?: ProductVariation[];   // размеры, цвета и т.д.
  
  // SEO и метаданные
  tags: string[];
  meta_title?: string;
  meta_description?: string;
  
  status: "draft" | "active" | "out_of_stock" | "discontinued";
  created_at: string;
  updated_at: string;
}
```

### Импорт данных
```typescript
interface ImportJob {
  id: string;
  storefront_id: string;
  type: "url" | "file";
  source: string;                    // URL или filename
  status: ImportStatus;
  progress: {
    total_items: number;
    processed_items: number;
    successful_items: number;
    failed_items: number;
    percentage: number;
  };
  mapping: FieldMapping;
  options: ImportOptions;
  errors: ImportError[];
  created_at: string;
  started_at?: string;
  completed_at?: string;
}

type ImportStatus = 
  | "queued"
  | "validating" 
  | "processing"
  | "completed"
  | "failed"
  | "cancelled";

interface FieldMapping {
  name: string;                      // поле в источнике → Product.name
  price: string;
  description: string;
  sku?: string;
  category?: string;
  stock_quantity?: string;
  images?: string[];                 // URL изображений
  custom_fields?: Record<string, string>;
}

interface ImportOptions {
  update_existing: boolean;          // обновлять существующие товары
  skip_duplicates: boolean;          // пропускать дубликаты
  auto_categorize: boolean;          // автоматическая категоризация
  download_images: boolean;          // скачивать изображения
  validate_urls: boolean;            // проверять валидность URL
}
```

## 🔄 Интеграции

### OpenSearch Integration
- **Index**: `storefronts`
- **Geospatial**: Поддержка geo_point для карты
- **Filtering**: По категориям, локации, режиму работы
- **Autocomplete**: Поиск по названиям и описаниям

### Database Schema
```sql
-- Основная таблица витрин
storefronts (
  id, slug, name, description, category,
  owner_id, location_data, contacts_json,
  business_hours_json, logo_url, cover_images,
  settings_json, status, created_at, updated_at
);

-- Товары витрин
storefront_products (
  id, storefront_id, name, description,
  category, sku, price, currency,
  stock_quantity, low_stock_threshold,
  images_json, attributes_json, status
);

-- Персонал витрин
storefront_staff (
  storefront_id, user_id, role, permissions,
  created_at, updated_at
);

-- Аналитика
storefront_analytics (
  storefront_id, date, views, clicks,
  orders, revenue, visitors
);

-- Импорт задачи
import_jobs (
  id, storefront_id, type, source,
  status, progress_json, mapping_json,
  options_json, errors_json
);
```

### MinIO Integration
- **Bucket**: `storefronts`
- **Logos**: `/storefronts/{id}/logo.{ext}`
- **Covers**: `/storefronts/{id}/covers/{image_id}.{ext}`
- **Products**: `/storefronts/{id}/products/{product_id}/{image_id}.{ext}`

## 🎛️ Бизнес-логика

### Геопространственный поиск
```typescript
function findNearbyStorefronts(
  lat: number, 
  lng: number, 
  radius: number
): NearbyStorefront[] {
  // PostgreSQL + PostGIS query
  const query = `
    SELECT *, ST_Distance(
      location_point, 
      ST_SetSRID(ST_MakePoint($2, $1), 4326)
    ) as distance
    FROM storefronts 
    WHERE ST_DWithin(
      location_point, 
      ST_SetSRID(ST_MakePoint($2, $1), 4326),
      $3
    )
    AND status = 'active'
    ORDER BY distance
  `;
  
  return query(lat, lng, radiusInMeters);
}
```

### Система слагов
- Автогенерация из названия: "Мој дућан" → "moj-ducan"
- Проверка уникальности
- Поддержка кириллицы и латиницы
- История изменений для SEO

### Управление инвентарем
```typescript
interface InventoryUpdate {
  product_id: string;
  stock_quantity: number;
  reason: "sale" | "restock" | "adjustment" | "return";
  notes?: string;
}

function updateInventory(updates: InventoryUpdate[]): void {
  // Atomic операция с транзакциями
  // WebSocket уведомления при низких остатках
  // Автоматические уведомления поставщикам
  // Интеграция с системами учета
}
```

## 🛡️ Безопасность и валидация

### Права доступа
```typescript
enum StorefrontRole {
  OWNER = "owner",           // полный доступ
  MANAGER = "manager",       // управление товарами + аналитика
  STAFF = "staff",           // только просмотр заказов
  VIEWER = "viewer"          // только чтение
}

interface StorefrontPermissions {
  can_edit_info: boolean;
  can_manage_products: boolean;
  can_view_analytics: boolean;
  can_manage_staff: boolean;
  can_delete: boolean;
}
```

### Валидация импорта
- Проверка форматов файлов (CSV, XLSX)
- Лимит размера: 100MB, 50K строк
- Валидация данных: цены, URL изображений
- Дубликаты по SKU или названию
- Карантин подозрительных товаров

## ⚠️ Известные особенности

### Performance
- Кеширование витрин на CDN
- Lazy loading товаров и изображений
- Пагинация курсорами для больших каталогов
- Оптимизация геопространственных запросов

### SEO и Discovery
- Уникальные URL слаги
- Meta теги для социальных сетей
- Структурированные данные (JSON-LD)
- Sitemap с приоритизацией

### Analytics
- Трекинг просмотров и кликов
- Конверсионная воронка
- A/B тесты страниц витрин
- Интеграция с Google Analytics

## 🧪 Примеры использования

### Создание витрины
```bash
curl -X POST /api/v1/storefronts/ \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Техно маркет",
    "description": "Продажа компјутера и електронике",
    "category": "electronics",
    "location": {
      "address": "Кнез Михаилова 5",
      "city": "Београд",
      "coordinates": [44.8176, 20.4633]
    }
  }'
```

### Поиск витрин
```bash
curl "/api/v1/storefronts/search?q=elektronika&location=Belgrade&radius=10&open_now=true"
```

### Bulk импорт товаров
```bash
curl -X POST /api/v1/storefronts/123/import/file \
  -H "Authorization: Bearer <token>" \
  -F "file=@products.csv" \
  -F "mapping={\"name\":\"Naziv\",\"price\":\"Cena\"}" \
  -F "options={\"update_existing\":true}"
```