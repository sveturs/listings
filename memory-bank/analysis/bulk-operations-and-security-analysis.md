# Анализ массовых операций и проблемы безопасности изображений

## Дата анализа: 2025-06-24
## Ветка: feature/improve-product-attributes-selection

## Оглавление
1. [Текущий статус задачи](#текущий-статус-задачи)
2. [Архитектура Backend](#архитектура-backend)
3. [Архитектура Frontend](#архитектура-frontend)
4. [База данных](#база-данных)
5. [Критическая проблема безопасности](#критическая-проблема-безопасности)
6. [План реализации массовых операций](#план-реализации-массовых-операций)
7. [План решения проблемы безопасности](#план-решения-проблемы-безопасности)

## Текущий статус задачи

### Что уже сделано:
- ✅ Исправлены ошибки загрузки изображений с Unsplash
- ✅ Исправлена проблема с async params в Next.js 13+
- ✅ Реализован UI для массового выбора товаров (чекбоксы) - **НО ЭТО БЫЛО ОШИБОЧНО**
- ✅ Создан компонент BulkActions для массовых операций - **НО ЭТО БЫЛО ОШИБОЧНО**
- ✅ Добавлены переводы для массовых операций
- ✅ Исправлены замечания по PR (хардкод в MarketplaceCard, перегенерация типов)

### Реальный статус:
**ВАЖНО**: При детальном изучении кода выяснилось, что функционал массовых операций НЕ реализован. Предыдущие отметки о реализации были ошибочными.

## Архитектура Backend

### Структура файлов

```
backend/
├── internal/proj/storefronts/
│   ├── handler/
│   │   ├── product_handler.go      # HTTP handlers для товаров
│   │   └── import_handler.go       # Handlers для импорта товаров
│   └── service/
│       ├── product_service.go      # Бизнес-логика товаров
│       └── import_service.go       # Логика импорта
├── internal/storage/
│   ├── postgres/
│   │   └── storefront_product.go   # Методы работы с БД
│   ├── minio/
│   │   └── minio.go               # Работа с S3-хранилищем
│   └── opensearch/
│       └── storefront_search.go    # Поисковая индексация
└── internal/domain/models/
    └── storefront_product.go       # Модели данных
```

### Существующие API endpoints

#### Product endpoints (в `product_handler.go`):

```go
// Публичные (без авторизации):
GET    /api/v1/storefronts/slug/:slug/products        // GetProducts() - список товаров
GET    /api/v1/storefronts/slug/:slug/products/:id    // GetProduct() - один товар

// Защищенные (требуют Bearer token):
POST   /api/v1/storefronts/slug/:slug/products        // CreateProduct() - создание
PUT    /api/v1/storefronts/slug/:slug/products/:id    // UpdateProduct() - обновление
DELETE /api/v1/storefronts/slug/:slug/products/:id    // DeleteProduct() - удаление
POST   /api/v1/storefronts/slug/:slug/products/:id/inventory // UpdateInventory()
GET    /api/v1/storefronts/slug/:slug/products/stats  // GetProductStats() - статистика
```

#### Import endpoints (в `import_handler.go`):

```go
// Через ID витрины:
POST /api/v1/storefronts/:id/import/url     // ImportFromURL()
POST /api/v1/storefronts/:id/import/file    // ImportFromFile()

// Через slug витрины:
POST /api/v1/storefronts/slug/:slug/import/url   // ImportFromURLBySlug()
POST /api/v1/storefronts/slug/:slug/import/file  // ImportFromFileBySlug()

// Управление импортом:
GET    /api/v1/import/jobs           // GetImportJobs() - список задач
GET    /api/v1/import/jobs/:jobId    // GetImportJob() - детали задачи
DELETE /api/v1/import/jobs/:jobId    // CancelImportJob() - отмена
```

### Модели данных (models/storefront_product.go)

```go
// Основная модель товара
type StorefrontProduct struct {
    ID            int                    `json:"id"`
    StorefrontID  int                    `json:"storefront_id"`
    Storefront    *Storefront           `json:"storefront,omitempty"`
    Name          string                `json:"name"`
    Description   *string               `json:"description,omitempty"`
    CategoryID    *int                  `json:"category_id,omitempty"`
    Category      *MarketplaceCategory  `json:"category,omitempty"`
    Images        []StorefrontProductImage `json:"images,omitempty"`
    Price         float64               `json:"price"`
    OldPrice      *float64              `json:"old_price,omitempty"`
    Currency      string                `json:"currency"`
    SKU           *string               `json:"sku,omitempty"`
    Barcode       *string               `json:"barcode,omitempty"`
    Stock         int                   `json:"stock"`
    StockStatus   string                `json:"stock_status"`
    Unit          *string               `json:"unit,omitempty"`
    Attributes    JSONB                 `json:"attributes,omitempty"`
    Variants      []StorefrontProductVariant `json:"variants,omitempty"`
    IsActive      bool                  `json:"is_active"`
    CreatedAt     time.Time             `json:"created_at"`
    UpdatedAt     time.Time             `json:"updated_at"`
}

// Изображения товара
type StorefrontProductImage struct {
    ID           int       `json:"id"`
    ProductID    int       `json:"product_id"`
    ImageURL     string    `json:"image_url"`    // ВАЖНО: хранит внешний URL!
    DisplayOrder int       `json:"display_order"`
    IsMain       bool      `json:"is_main"`
    CreatedAt    time.Time `json:"created_at"`
}

// Запросы для API
type CreateProductRequest struct {
    Name         string                 `json:"name" validate:"required,min=3,max=200"`
    Description  *string                `json:"description" validate:"omitempty,max=2000"`
    CategoryID   *int                   `json:"category_id" validate:"omitempty,min=1"`
    Images       []CreateProductImage   `json:"images" validate:"omitempty,dive"`
    Price        float64                `json:"price" validate:"required,min=0"`
    OldPrice     *float64               `json:"old_price" validate:"omitempty,min=0"`
    Currency     string                 `json:"currency" validate:"required,len=3"`
    SKU          *string                `json:"sku" validate:"omitempty,max=100"`
    Barcode      *string                `json:"barcode" validate:"omitempty,max=100"`
    Stock        int                    `json:"stock" validate:"min=0"`
    Unit         *string                `json:"unit" validate:"omitempty,max=20"`
    Attributes   map[string]interface{} `json:"attributes,omitempty"`
    IsActive     bool                   `json:"is_active"`
}
```

### Storage методы (postgres/storefront_product.go)

```go
// Текущие методы - только для единичных операций:
func (s *PostgresStorage) GetStorefrontProducts(ctx context.Context, storefrontID int, params storage.GetProductsParams) ([]models.StorefrontProduct, int, error)
func (s *PostgresStorage) GetStorefrontProduct(ctx context.Context, storefrontID, productID int) (*models.StorefrontProduct, error)
func (s *PostgresStorage) CreateStorefrontProduct(ctx context.Context, product *models.StorefrontProduct) error
func (s *PostgresStorage) UpdateStorefrontProduct(ctx context.Context, product *models.StorefrontProduct) error
func (s *PostgresStorage) DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error
func (s *PostgresStorage) UpdateProductInventory(ctx context.Context, storefrontID, productID int, stockChange int, movementType, reason string) error
func (s *PostgresStorage) GetProductStats(ctx context.Context, storefrontID int) (*models.ProductStats, error)

// Специальные методы для импорта (без проверки владельца):
func (s *PostgresStorage) CreateProductForImport(ctx context.Context, product *models.StorefrontProduct) error
func (s *PostgresStorage) UpdateProductForImport(ctx context.Context, product *models.StorefrontProduct) error
```

### Отсутствующие компоненты для массовых операций

1. **Модели для bulk операций** - НЕ СУЩЕСТВУЮТ
2. **Storage методы для bulk операций** - НЕ СУЩЕСТВУЮТ
3. **API endpoints для bulk операций** - НЕ СУЩЕСТВУЮТ
4. **Оптимизация для OpenSearch bulk индексации** - НЕ РЕАЛИЗОВАНА

## Архитектура Frontend

### Структура файлов

```
frontend/svetu/src/
├── app/[locale]/storefronts/[slug]/products/
│   ├── page.tsx                    # Страница списка товаров
│   ├── new/page.tsx               # Создание товара
│   └── import/page.tsx            # Импорт товаров
├── components/
│   ├── products/
│   │   ├── ProductWizard.tsx      # Визард создания товара
│   │   └── steps/                 # Шаги визарда
│   └── import/
│       ├── ImportManager.tsx      # Управление импортом
│       └── ImportWizard.tsx       # Визард импорта
├── store/slices/
│   ├── storefrontSlice.ts         # Redux для витрин
│   └── importSlice.ts             # Redux для импорта
└── types/
    └── storefront.ts              # TypeScript типы
```

### Компонент списка товаров (products/page.tsx)

```typescript
// Локальное состояние (НЕ Redux):
const [products, setProducts] = useState<StorefrontProduct[]>([]);
const [loading, setLoading] = useState(true);
const [search, setSearch] = useState('');
const [filterStatus, setFilterStatus] = useState<string>('all');
const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
const [page, setPage] = useState(1);
const [hasMore, setHasMore] = useState(true);

// Функции для работы с товарами:
const loadProducts = async (pageNum: number, append = false)
const handleEdit = (product: StorefrontProduct)
const handleDelete = async (product: StorefrontProduct)
const handleStatusChange = (value: string)
```

### Отсутствующие компоненты для массовых операций

1. **Redux slice для товаров** - НЕ СУЩЕСТВУЕТ
2. **Состояние для выбранных товаров** - НЕ СУЩЕСТВУЕТ
3. **Чекбоксы для выбора** - НЕ СУЩЕСТВУЮТ
4. **Компонент BulkActions** - НЕ СУЩЕСТВУЕТ
5. **API методы для bulk операций** - НЕ СУЩЕСТВУЮТ

### TypeScript типы (types/storefront.ts)

```typescript
export interface StorefrontProduct {
  id: number;
  storefront_id: number;
  name: string;
  description?: string;
  category_id?: number;
  images?: StorefrontProductImage[];
  price: number;
  old_price?: number;
  currency: string;
  sku?: string;
  barcode?: string;
  stock: number;
  stock_status: 'in_stock' | 'low_stock' | 'out_of_stock';
  unit?: string;
  attributes?: Record<string, any>;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface StorefrontProductImage {
  id: number;
  product_id: number;
  image_url: string;  // ВАЖНО: внешний URL!
  display_order: number;
  is_main: boolean;
}
```

## База данных

### Таблицы PostgreSQL

#### storefront_products
```sql
CREATE TABLE storefront_products (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES marketplace_categories(id),
    price DECIMAL(10,2) NOT NULL,
    old_price DECIMAL(10,2),
    currency VARCHAR(3) NOT NULL DEFAULT 'RSD',
    sku VARCHAR(100),
    barcode VARCHAR(100),
    stock INTEGER NOT NULL DEFAULT 0,
    stock_status VARCHAR(20) NOT NULL DEFAULT 'out_of_stock',
    unit VARCHAR(20),
    attributes JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Уникальные индексы
    UNIQUE(storefront_id, sku),
    UNIQUE(storefront_id, barcode)
);

-- Полнотекстовый поиск
CREATE INDEX idx_storefront_products_name_gin ON storefront_products 
USING gin(to_tsvector('simple', name));

-- Триггер для автоматического обновления stock_status
CREATE TRIGGER update_stock_status 
BEFORE UPDATE ON storefront_products
FOR EACH ROW EXECUTE FUNCTION update_stock_status_trigger();
```

#### storefront_product_images
```sql
CREATE TABLE storefront_product_images (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,  -- ВАЖНО: хранит внешние URL!
    display_order INTEGER DEFAULT 0,
    is_main BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Индексы
    INDEX idx_product_images_product_id (product_id),
    INDEX idx_product_images_order (product_id, display_order)
);
```

#### storefront_product_variants
```sql
CREATE TABLE storefront_product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    sku VARCHAR(100),
    price_modifier DECIMAL(10,2) DEFAULT 0,
    stock INTEGER DEFAULT 0,
    attributes JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### storefront_inventory_movements
```sql
CREATE TABLE storefront_inventory_movements (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    movement_type VARCHAR(50) NOT NULL, -- 'purchase', 'sale', 'adjustment', etc
    quantity INTEGER NOT NULL,
    reason TEXT,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Сравнение с marketplace_images (безопасная реализация)

```sql
-- marketplace_images - БЕЗОПАСНО (изображения в MinIO)
CREATE TABLE marketplace_images (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER REFERENCES marketplace_listings(id),
    file_name VARCHAR(255) NOT NULL,    -- Имя файла в MinIO
    file_size INTEGER,
    mime_type VARCHAR(100),
    public_url TEXT,                     -- URL из MinIO (наш сервер)
    thumbnail_url TEXT,
    display_order INTEGER DEFAULT 0,
    is_main BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Критическая проблема безопасности

### Суть проблемы

**Marketplace (объявления)** - БЕЗОПАСНАЯ реализация:
- При создании объявления изображения загружаются на сервер
- Сохраняются в MinIO (S3-совместимое хранилище)
- В БД хранятся ссылки на файлы в MinIO
- Полный контроль над контентом

**Storefronts (товары магазинов)** - НЕБЕЗОПАСНАЯ реализация:
- При импорте товаров из XML/CSV внешние URL изображений сохраняются напрямую в БД
- Изображения НЕ загружаются в MinIO
- Frontend отображает изображения напрямую с внешних серверов

### Риски

1. **Утечка данных пользователей**:
   - IP-адреса пользователей передаются внешним серверам
   - Возможность отслеживания активности пользователей
   - Нарушение GDPR и законов о конфиденциальности

2. **Безопасность контента**:
   - Изображения могут быть изменены после модерации
   - Возможна подмена на вредоносный контент
   - Отсутствие контроля над контентом

3. **Надежность**:
   - Зависимость от внешних серверов
   - Изображения могут стать недоступными
   - Внешние серверы могут быть медленными

### Пример проблемного кода

В `import_service.go`:
```go
func (s *ImportService) parseXMLProduct(xmlProduct XMLProduct) (*models.StorefrontProduct, error) {
    // ...
    
    // Изображения из XML
    if len(xmlProduct.Images) > 0 {
        for i, imgURL := range xmlProduct.Images {
            product.Images = append(product.Images, models.StorefrontProductImage{
                ImageURL:     imgURL,  // ПРОБЛЕМА: внешний URL сохраняется напрямую!
                DisplayOrder: i,
                IsMain:       i == 0,
            })
        }
    }
    
    return product, nil
}
```

### Решение уже существует в коде

В `minio/minio.go` есть готовый метод:
```go
func (s *MinioStorage) UploadImageFromURL(ctx context.Context, imageURL, bucketName, objectName string) (string, error)
```

Этот метод:
1. Скачивает изображение по URL
2. Проверяет тип файла
3. Загружает в MinIO
4. Возвращает публичный URL из MinIO

## План реализации массовых операций

### Backend

#### 1. Создать модели для bulk операций в `models/storefront_product.go`:

```go
// Массовое создание
type BulkCreateProductsRequest struct {
    Products []CreateProductRequest `json:"products" validate:"required,min=1,max=100,dive"`
}

type BulkCreateProductsResponse struct {
    Created []int                `json:"created"`    // ID созданных товаров
    Failed  []BulkOperationError `json:"failed"`     // Ошибки
}

// Массовое обновление
type BulkUpdateProductsRequest struct {
    Updates []BulkUpdateItem `json:"updates" validate:"required,min=1,max=100,dive"`
}

type BulkUpdateItem struct {
    ProductID int                  `json:"product_id" validate:"required"`
    Updates   UpdateProductRequest `json:"updates" validate:"required"`
}

// Массовое удаление
type BulkDeleteProductsRequest struct {
    ProductIDs []int `json:"product_ids" validate:"required,min=1,max=100"`
}

// Массовое обновление статуса
type BulkUpdateStatusRequest struct {
    ProductIDs []int `json:"product_ids" validate:"required,min=1,max=100"`
    IsActive   bool  `json:"is_active"`
}

// Ошибка операции
type BulkOperationError struct {
    Index   int    `json:"index"`
    ID      int    `json:"id,omitempty"`
    Error   string `json:"error"`
}
```

#### 2. Добавить storage методы в `postgres/storefront_product.go`:

```go
// Массовое создание с транзакцией
func (s *PostgresStorage) BulkCreateProducts(ctx context.Context, storefrontID int, products []models.StorefrontProduct) ([]int, []error) {
    // Использовать транзакцию
    // Batch INSERT для оптимизации
    // Вернуть ID созданных и ошибки
}

// Массовое обновление
func (s *PostgresStorage) BulkUpdateProducts(ctx context.Context, storefrontID int, updates map[int]models.UpdateProductRequest) ([]int, []error) {
    // Проверить принадлежность товаров витрине
    // Использовать транзакцию
    // Оптимизированные UPDATE запросы
}

// Массовое удаление
func (s *PostgresStorage) BulkDeleteProducts(ctx context.Context, storefrontID int, productIDs []int) ([]int, []error) {
    // Проверить принадлежность товаров
    // Soft delete или hard delete
    // Каскадное удаление связанных данных
}
```

#### 3. Создать API endpoints в `product_handler.go`:

```go
// Регистрация роутов
auth.POST("/products/bulk/create", h.BulkCreateProducts)
auth.PUT("/products/bulk/update", h.BulkUpdateProducts)
auth.DELETE("/products/bulk/delete", h.BulkDeleteProducts)
auth.PUT("/products/bulk/status", h.BulkUpdateStatus)

// Handlers
func (h *ProductHandler) BulkCreateProducts(c *fiber.Ctx) error
func (h *ProductHandler) BulkUpdateProducts(c *fiber.Ctx) error
func (h *ProductHandler) BulkDeleteProducts(c *fiber.Ctx) error
func (h *ProductHandler) BulkUpdateStatus(c *fiber.Ctx) error
```

#### 4. Оптимизировать OpenSearch индексацию:

```go
// В product_service.go добавить bulk индексацию
func (s *ProductService) BulkIndexProducts(products []models.StorefrontProduct) error {
    // Использовать OpenSearch Bulk API
    // Асинхронная обработка больших объемов
}
```

### Frontend

#### 1. Создать Redux slice для товаров:

```typescript
// store/slices/productSlice.ts
interface ProductState {
  products: StorefrontProduct[];
  selectedIds: number[];
  loading: boolean;
  bulkOperation: {
    isProcessing: boolean;
    progress: number;
    errors: BulkOperationError[];
  };
}

const productSlice = createSlice({
  name: 'products',
  initialState,
  reducers: {
    toggleProductSelection: (state, action) => {},
    selectAll: (state) => {},
    clearSelection: (state) => {},
    setBulkOperationProgress: (state, action) => {},
  },
});
```

#### 2. Добавить компонент BulkActions:

```typescript
// components/products/BulkActions.tsx
interface BulkActionsProps {
  selectedCount: number;
  onBulkDelete: () => void;
  onBulkStatusChange: (active: boolean) => void;
  onBulkExport: () => void;
}

export function BulkActions({ selectedCount, ...handlers }: BulkActionsProps) {
  return (
    <div className="flex items-center gap-4 p-4 bg-base-200 rounded-lg">
      <span>{selectedCount} товаров выбрано</span>
      <div className="flex gap-2">
        <button className="btn btn-sm" onClick={() => handlers.onBulkStatusChange(false)}>
          Деактивировать
        </button>
        <button className="btn btn-sm btn-error" onClick={handlers.onBulkDelete}>
          Удалить
        </button>
      </div>
    </div>
  );
}
```

#### 3. Обновить список товаров с чекбоксами:

```typescript
// В ProductCard добавить чекбокс
<input
  type="checkbox"
  className="checkbox"
  checked={isSelected}
  onChange={() => dispatch(toggleProductSelection(product.id))}
/>

// Header с "выбрать все"
<input
  type="checkbox"
  className="checkbox"
  checked={allSelected}
  onChange={() => dispatch(allSelected ? clearSelection() : selectAll())}
/>
```

#### 4. Создать API методы:

```typescript
// services/productApi.ts
export const productApi = {
  bulkCreate: async (storefrontId: string, products: CreateProductRequest[]) => {},
  bulkUpdate: async (storefrontId: string, updates: BulkUpdateItem[]) => {},
  bulkDelete: async (storefrontId: string, productIds: number[]) => {},
  bulkUpdateStatus: async (storefrontId: string, productIds: number[], isActive: boolean) => {},
};
```

## План решения проблемы безопасности

### 1. Модифицировать import_service.go:

```go
func (s *ImportService) processProductImages(ctx context.Context, product *models.StorefrontProduct) error {
    for i, img := range product.Images {
        if strings.HasPrefix(img.ImageURL, "http") {
            // Генерировать уникальное имя файла
            objectName := fmt.Sprintf("storefronts/%d/products/%s_%d.jpg", 
                product.StorefrontID, uuid.New().String(), i)
            
            // Загрузить изображение в MinIO
            publicURL, err := s.minioStorage.UploadImageFromURL(
                ctx, img.ImageURL, "listings", objectName)
            if err != nil {
                log.Printf("Failed to upload image: %v", err)
                continue
            }
            
            // Заменить внешний URL на URL из MinIO
            product.Images[i].ImageURL = publicURL
        }
    }
    return nil
}
```

### 2. Добавить миграцию для существующих товаров:

```go
// Создать отдельный скрипт или background job
func MigrateExternalImages() {
    // Получить все товары с внешними URL
    // Для каждого изображения:
    //   - Скачать с внешнего URL
    //   - Загрузить в MinIO
    //   - Обновить URL в БД
    // Логировать прогресс и ошибки
}
```

### 3. Обновить валидацию при создании товаров:

```go
// В product_handler.go
func validateProductImages(images []CreateProductImage) error {
    for _, img := range images {
        if strings.HasPrefix(img.ImageURL, "http") && 
           !strings.Contains(img.ImageURL, os.Getenv("MINIO_PUBLIC_URL")) {
            return errors.New("external image URLs are not allowed")
        }
    }
    return nil
}
```

### 4. Добавить настройку для режима безопасности:

```yaml
# В конфигурации
security:
  allow_external_images: false  # По умолчанию запретить
  auto_download_images: true    # Автоматически загружать при импорте
```

## Важные замечания

1. **Приоритеты**:
   - Сначала решить проблему безопасности (критично)
   - Затем реализовать массовые операции

2. **Совместимость**:
   - Сохранить обратную совместимость для существующих API
   - Добавить флаги для контроля поведения

3. **Производительность**:
   - Использовать горутины для параллельной загрузки изображений
   - Ограничить количество одновременных операций
   - Добавить прогресс-бары для длительных операций

4. **Мониторинг**:
   - Логировать все массовые операции
   - Отслеживать failed операции
   - Уведомлять администраторов о проблемах

## Тестирование

### Backend тесты:
- Unit тесты для bulk методов в storage
- Integration тесты для API endpoints
- Тесты производительности для больших объемов
- Тесты загрузки изображений из разных источников

### Frontend тесты:
- Компонентные тесты для BulkActions
- E2E тесты для сценариев массовых операций
- Тесты производительности при большом количестве товаров

## Метрики успеха

1. **Безопасность**:
   - 100% изображений товаров хранятся в MinIO
   - 0 внешних запросов при отображении товаров
   - Полное соответствие GDPR

2. **Производительность массовых операций**:
   - Обработка 1000 товаров < 10 секунд
   - Отзывчивый UI во время операций
   - Корректная обработка ошибок

3. **Пользовательский опыт**:
   - Интуитивный интерфейс выбора товаров
   - Понятные сообщения об ошибках
   - Возможность отмены операций