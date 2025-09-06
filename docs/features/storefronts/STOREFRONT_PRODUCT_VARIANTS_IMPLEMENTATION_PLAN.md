# 📋 СИСТЕМА ВАРИАНТОВ ТОВАРОВ ДЛЯ ВИТРИН - АКТУАЛЬНОЕ СОСТОЯНИЕ

## ✅ СТАТУС РЕАЛИЗАЦИИ

**Общий прогресс:** 95% завершено  
**Последнее обновление:** 30.07.2025  

### 🎯 Реализованные функции:
- ✅ **Полная система вариантов** - создание, редактирование, удаление
- ✅ **Цветовые плитки** - современный UI для выбора цветов (как на Avito)
- ✅ **Корзина с вариантами** - добавление конкретных вариантов в корзину
- ✅ **Модальные окна** - выбор вариантов при добавлении в корзину
- ✅ **Marketplace интеграция** - товары витрин корректно отображаются в маркетплейсе
- ✅ **Резервирование товаров** - система inventory_reservations
- ✅ **Управление остатками** - автоматический подсчет доступного количества
- ✅ **API endpoints** - полный набор публичных и приватных методов
- ✅ **Валидация уникальности** - проверка комбинаций атрибутов
- ✅ **Транзакционная безопасность** - все операции в транзакциях

## 📊 Анализ текущего состояния системы

### 🗄️ База данных

#### Таблица: `storefront_product_variants`
Расположение: `/data/hostel-booking-system/backend/migrations/000074_create_table_storefront_product_variants.up.sql`

```sql
CREATE TABLE public.storefront_product_variants (
    id integer NOT NULL,
    product_id integer NOT NULL,
    sku character varying(100),
    barcode character varying(100),
    price numeric(15,2),
    compare_at_price numeric(15,2),
    cost_price numeric(15,2),
    stock_quantity integer DEFAULT 0 NOT NULL,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,  -- добавлено в миграции 000179
    available_quantity INTEGER GENERATED ALWAYS AS (stock_quantity - reserved_quantity) STORED,
    stock_status character varying(20) DEFAULT 'in_stock'::character varying NOT NULL,
    low_stock_threshold integer DEFAULT 5,
    variant_attributes jsonb DEFAULT '{}'::jsonb NOT NULL,
    weight numeric(10,3),
    dimensions jsonb,
    is_active boolean DEFAULT true NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
    view_count integer DEFAULT 0 NOT NULL,
    sold_count integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

**Важные индексы:**
- `idx_storefront_product_variants_unique_attributes` - уникальность комбинации атрибутов
- `idx_storefront_product_variants_low_stock` - поиск товаров с низкими остатками
- `idx_storefront_product_variants_default_unique` - только один вариант может быть дефолтным

**Триггеры:**
- `trigger_update_variant_stock_status` - автоматическое обновление stock_status
- `trigger_update_storefront_product_variants_updated_at` - обновление updated_at

#### Связанные таблицы:
- `product_variant_attributes` (ID: 2000+) - глобальные атрибуты вариантов
- `product_variant_attribute_values` - возможные значения атрибутов
- `storefront_product_attributes` - настройки атрибутов для конкретного товара
- `inventory_reservations` - резервирование товаров

### 🔧 Backend структура

#### 1. Handlers (HTTP слой)

**Файл:** `/data/hostel-booking-system/backend/internal/proj/storefront/handler/variant_handler.go`
- `GetProductVariants(c *fiber.Ctx)` - получение вариантов товара
- `GetVariantByID(c *fiber.Ctx)` - получение конкретного варианта
- `CreateVariant(c *fiber.Ctx)` - создание варианта
- `UpdateVariant(c *fiber.Ctx)` - обновление варианта
- `DeleteVariant(c *fiber.Ctx)` - удаление варианта
- `BulkCreateVariants(c *fiber.Ctx)` - массовое создание
- `GenerateVariants(c *fiber.Ctx)` - автогенерация вариантов
- `ImportVariantsCSV(c *fiber.Ctx)` - импорт из CSV
- `ExportVariantsCSV(c *fiber.Ctx)` - экспорт в CSV

**Файл:** `/data/hostel-booking-system/backend/internal/proj/storefront/handler/public_variant_handler.go`
- Публичные методы для неавторизованных пользователей

**Файл:** `/data/hostel-booking-system/backend/internal/proj/storefronts/handler/product_handler.go`
- `CreateProduct(c *fiber.Ctx)` - создание товара (НЕ поддерживает варианты!)
- `UpdateProduct(c *fiber.Ctx)` - обновление товара
- `BulkCreateProducts(c *fiber.Ctx)` - массовое создание

#### 2. Services (бизнес-логика)

**Файл:** `/data/hostel-booking-system/backend/internal/proj/storefront/service/variant_service.go`
- `CreateVariant(ctx, productID, variant)` - создание с валидацией
- `UpdateVariant(ctx, variantID, updates)` - обновление с проверками
- `GenerateVariants(ctx, productID, matrix)` - генерация комбинаций
- `ValidateVariantUniqueness(ctx, productID, attributes)` - проверка уникальности

**Файл:** `/data/hostel-booking-system/backend/internal/proj/storefronts/service/product_service.go`
- `CreateProduct(ctx, storefrontID, userID, req)` - создание товара
- `validateCreateRequest(req)` - валидация запроса

#### 3. Repository (слой данных)

**Файл:** `/data/hostel-booking-system/backend/internal/proj/storefront/repository/variant_repository.go`
- `CreateVariant(ctx, variant)` - SQL вставка
- `GetVariantsByProductID(ctx, productID)` - получение всех вариантов
- `UpdateVariantStock(ctx, variantID, quantity)` - обновление остатков
- `BulkCreateVariants(ctx, variants)` - массовая вставка

#### 4. Models (структуры данных)

**Файл:** `/data/hostel-booking-system/backend/internal/proj/storefront/types/variant_types.go`
```go
type ProductVariant struct {
    ID                int                    `json:"id" db:"id"`
    ProductID         int                    `json:"product_id" db:"product_id"`
    SKU               *string                `json:"sku,omitempty" db:"sku"`
    Price             *float64               `json:"price,omitempty" db:"price"`
    StockQuantity     int                    `json:"stock_quantity" db:"stock_quantity"`
    ReservedQuantity  int                    `json:"reserved_quantity" db:"reserved_quantity"`
    AvailableQuantity int                    `json:"available_quantity" db:"available_quantity"`
    VariantAttributes map[string]interface{} `json:"variant_attributes" db:"variant_attributes"`
    // ... другие поля
}

type CreateVariantRequest struct {
    ProductID         int                    `json:"product_id" validate:"required"`
    SKU               *string                `json:"sku,omitempty"`
    Price             *float64               `json:"price,omitempty"`
    StockQuantity     int                    `json:"stock_quantity" validate:"min=0"`
    VariantAttributes map[string]interface{} `json:"variant_attributes" validate:"required"`
    IsDefault         bool                   `json:"is_default"`
}

type GenerateVariantsRequest struct {
    ProductID         int                    `json:"product_id" validate:"required"`
    AttributeMatrix   map[string][]string    `json:"attribute_matrix" validate:"required"`
    PriceModifiers    map[string]float64     `json:"price_modifiers,omitempty"`
    StockQuantities   map[string]int         `json:"stock_quantities,omitempty"`
}
```

**Файл:** `/data/hostel-booking-system/backend/internal/domain/models/storefront_product.go`
```go
type CreateProductRequest struct {
    Name          string                 `json:"name" validate:"required,min=3,max=255"`
    Description   string                 `json:"description" validate:"required,min=10"`
    Price         float64                `json:"price" validate:"required,min=0"`
    CategoryID    int                    `json:"category_id" validate:"required"`
    StockQuantity int                    `json:"stock_quantity" validate:"min=0"`
    // НЕТ поддержки вариантов!
}
```

#### 5. Routes (маршруты)

**Файл:** `/data/hostel-booking-system/backend/internal/proj/storefronts/module.go`
```go
// Защищенные маршруты для вариантов
variants := protected.Group("/products/:product_id/variants")
variants.Get("/", m.variantHandler.GetProductVariants)
variants.Post("/", m.variantHandler.CreateVariant)
variants.Post("/bulk", m.variantHandler.BulkCreateVariants)
variants.Post("/generate", m.variantHandler.GenerateVariants)

// Публичные маршруты
publicVariants := api.Group("/public")
publicVariants.Get("/storefronts/:slug/products/:product_id/variants", m.publicVariantHandler.GetProductVariantsPublic)
```

### 🎨 Frontend структура

#### 1. Components

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/components/products/ProductWizard.tsx`
- 7-шаговый wizard для создания товара
- Шаг 5: `VariantsStep` - управление вариантами

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/components/products/steps/VariantsStep.tsx`
```typescript
// Состояние компонента
const [activeMode, setActiveMode] = useState<'none' | 'simple' | 'advanced'>('none');
const [hasVariants, setHasVariants] = useState(false);
const [variants, setVariants] = useState<any[]>([]);

// НЕ интегрирован с CreateProductContext!
// Данные о вариантах не сохраняются в контексте
```

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/components/Storefront/ProductVariants/VariantManager.tsx`
- Полноценный UI для управления вариантами
- Таблица вариантов с редактированием
- Массовые операции

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/components/Storefront/ProductVariants/VariantGenerator.tsx`
- Генератор вариантов по матрице атрибутов
- НЕ имеет полей для указания stock quantity!

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/components/Storefront/ProductVariants/AttributeSetup.tsx`
- Настройка атрибутов для товара
- Выбор из глобальных атрибутов категории

#### 2. Context

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/contexts/CreateProductContext.tsx`
```typescript
interface ProductFormData {
    // Базовая информация
    name: string;
    description: string;
    price: number;
    categoryId: number | null;
    // НЕТ полей для вариантов!
}
```

#### 3. Services/API

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/services/storefrontProducts.ts`
- `createProduct(storefrontSlug, data)` - создание товара
- НЕ поддерживает варианты в запросе!

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/services/productApi.ts`
- Методы для работы с вариантами существующих товаров
- `getProductVariants(productId)`
- `createVariant(productId, data)`
- `generateVariants(productId, matrix)`

#### 4. Types

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/types/storefront.ts`
```typescript
interface ProductVariant {
    id: number;
    product_id: number;
    sku?: string;
    price?: number;
    stock_quantity: number;
    reserved_quantity: number;
    available_quantity: number;
    variant_attributes: Record<string, any>;
    is_active: boolean;
    is_default: boolean;
}
```

## 🚨 ВЫЯВЛЕННЫЕ ПРОБЛЕМЫ

### 1. Отсутствие интеграции между созданием товара и вариантами
- `CreateProductRequest` НЕ содержит поля для вариантов
- `ProductService.CreateProduct` НЕ создает варианты
- Frontend `CreateProductContext` НЕ хранит данные о вариантах
- `VariantsStep` НЕ интегрирован с контекстом

### 2. Невозможность указать количество для вариантов при создании
- `VariantGenerator` НЕ имеет полей для stock quantity
- Нет UI для массового указания количества
- Приходится редактировать каждый вариант отдельно после создания

### 3. Разрыв между атрибутами категории и вариантами
- Атрибуты категории (brand, color, storage) не связаны автоматически с вариантами
- Нет валидации соответствия атрибутов варианта категории товара

## 🛠️ ДЕТАЛЬНЫЙ ПЛАН РЕАЛИЗАЦИИ

### Этап 1: Модификация Backend для поддержки вариантов при создании товара

#### 1.1 Расширение модели CreateProductRequest

**Файл для изменения:** `/data/hostel-booking-system/backend/internal/domain/models/storefront_product.go`

```go
// CreateProductRequest - добавить поля для вариантов
type CreateProductRequest struct {
    // ... существующие поля ...
    
    // Новые поля для вариантов
    HasVariants      bool                      `json:"has_variants"`
    Variants         []CreateVariantInline     `json:"variants,omitempty" validate:"omitempty,dive"`
    VariantSettings  *VariantSettings          `json:"variant_settings,omitempty"`
}

// CreateVariantInline - структура для создания варианта вместе с товаром
type CreateVariantInline struct {
    SKU               *string                `json:"sku,omitempty"`
    Barcode           *string                `json:"barcode,omitempty"`
    Price             *float64               `json:"price,omitempty"`
    CompareAtPrice    *float64               `json:"compare_at_price,omitempty"`
    CostPrice         *float64               `json:"cost_price,omitempty"`
    StockQuantity     int                    `json:"stock_quantity" validate:"min=0"`
    LowStockThreshold *int                   `json:"low_stock_threshold,omitempty"`
    VariantAttributes map[string]interface{} `json:"variant_attributes" validate:"required"`
    Weight            *float64               `json:"weight,omitempty"`
    Dimensions        map[string]interface{} `json:"dimensions,omitempty"`
    IsDefault         bool                   `json:"is_default"`
}

// VariantSettings - настройки для вариантов
type VariantSettings struct {
    TrackInventory    bool     `json:"track_inventory"`
    ContinueSelling   bool     `json:"continue_selling"` // продолжать продажи при 0 stock
    RequireShipping   bool     `json:"require_shipping"`
    TaxableProduct    bool     `json:"taxable_product"`
    WeightUnit        string   `json:"weight_unit,omitempty"` // kg, g, lb, oz
    SelectedAttributes []string `json:"selected_attributes"` // какие атрибуты используются
}
```

#### 1.2 Модификация ProductService

**Файл для изменения:** `/data/hostel-booking-system/backend/internal/proj/storefronts/service/product_service.go`

```go
// CreateProduct - добавить поддержку создания вариантов
func (s *ProductService) CreateProduct(ctx context.Context, storefrontID, userID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
    // ... существующая валидация ...
    
    // Новая валидация для вариантов
    if req.HasVariants {
        if err := s.validateVariants(req); err != nil {
            return nil, fmt.Errorf("invalid variants: %w", err)
        }
    }
    
    // Начинаем транзакцию
    tx, err := s.storage.BeginTx(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %w", err)
    }
    defer tx.Rollback()
    
    // Создаем основной товар
    product, err := s.storage.CreateStorefrontProduct(ctx, storefrontID, req)
    if err != nil {
        return nil, err
    }
    
    // Создаем варианты если указаны
    if req.HasVariants && len(req.Variants) > 0 {
        variantRequests := make([]types.CreateVariantRequest, len(req.Variants))
        for i, v := range req.Variants {
            variantRequests[i] = types.CreateVariantRequest{
                ProductID:         product.ID,
                SKU:               v.SKU,
                Barcode:           v.Barcode,
                Price:             v.Price,
                CompareAtPrice:    v.CompareAtPrice,
                CostPrice:         v.CostPrice,
                StockQuantity:     v.StockQuantity,
                LowStockThreshold: v.LowStockThreshold,
                VariantAttributes: v.VariantAttributes,
                Weight:            v.Weight,
                Dimensions:        v.Dimensions,
                IsDefault:         v.IsDefault,
            }
        }
        
        createdVariants, err := s.variantService.BulkCreateVariants(ctx, product.ID, variantRequests)
        if err != nil {
            return nil, fmt.Errorf("failed to create variants: %w", err)
        }
        
        product.Variants = createdVariants
    }
    
    // Коммитим транзакцию
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    // Индексируем в OpenSearch
    go s.indexProductWithVariants(product)
    
    return product, nil
}

// validateVariants - валидация вариантов перед созданием
func (s *ProductService) validateVariants(req *models.CreateProductRequest) error {
    if len(req.Variants) == 0 {
        return errors.New("at least one variant is required when has_variants is true")
    }
    
    // Проверяем уникальность комбинаций атрибутов
    seen := make(map[string]bool)
    defaultCount := 0
    
    for i, v := range req.Variants {
        // Создаем ключ из атрибутов для проверки уникальности
        attrKey := generateAttributeKey(v.VariantAttributes)
        if seen[attrKey] {
            return fmt.Errorf("duplicate variant attributes at index %d", i)
        }
        seen[attrKey] = true
        
        // Проверяем что только один вариант дефолтный
        if v.IsDefault {
            defaultCount++
        }
        
        // Валидация SKU если указан
        if v.SKU != nil && *v.SKU != "" {
            if err := s.validateSKU(*v.SKU); err != nil {
                return fmt.Errorf("invalid SKU at index %d: %w", i, err)
            }
        }
    }
    
    if defaultCount > 1 {
        return errors.New("only one variant can be default")
    }
    
    if defaultCount == 0 && len(req.Variants) > 0 {
        // Делаем первый вариант дефолтным
        req.Variants[0].IsDefault = true
    }
    
    return nil
}
```

#### 1.3 Создание нового эндпоинта API

**Файл для изменения:** `/data/hostel-booking-system/backend/internal/proj/storefronts/handler/product_handler.go`

Модифицировать существующий `CreateProduct` для поддержки вариантов или создать новый метод `CreateProductWithVariants`.

### Этап 2: Интеграция Frontend

#### 2.1 Расширение CreateProductContext

**Файл для изменения:** `/data/hostel-booking-system/frontend/svetu/src/contexts/CreateProductContext.tsx`

```typescript
// Расширяем интерфейс ProductFormData
interface ProductFormData {
    // ... существующие поля ...
    
    // Новые поля для вариантов
    hasVariants: boolean;
    variants: ProductVariantCreate[];
    variantSettings: {
        trackInventory: boolean;
        continueSelling: boolean;
        selectedAttributes: string[];
    };
}

interface ProductVariantCreate {
    sku?: string;
    price?: number;
    stockQuantity: number;
    variantAttributes: Record<string, any>;
    isDefault: boolean;
}

// Добавить в initialFormData
const initialFormData: ProductFormData = {
    // ... существующие поля ...
    hasVariants: false,
    variants: [],
    variantSettings: {
        trackInventory: true,
        continueSelling: false,
        selectedAttributes: [],
    },
};

// Добавить методы для управления вариантами
const updateVariants = (variants: ProductVariantCreate[]) => {
    setFormData(prev => ({ ...prev, variants }));
};

const toggleHasVariants = (enabled: boolean) => {
    setFormData(prev => ({ 
        ...prev, 
        hasVariants: enabled,
        variants: enabled ? prev.variants : []
    }));
};
```

#### 2.2 Модификация VariantsStep

**Файл для изменения:** `/data/hostel-booking-system/frontend/svetu/src/components/products/steps/VariantsStep.tsx`

```typescript
export default function VariantsStep({ onNext, onBack }: VariantsStepProps) {
    const t = useTranslations('storefronts.products');
    const { state, updateFormData } = useCreateProduct();
    
    // Используем данные из контекста вместо локального состояния
    const hasVariants = state.formData.hasVariants;
    const variants = state.formData.variants;
    
    const handleVariantToggle = (enabled: boolean) => {
        updateFormData({ hasVariants: enabled });
        if (!enabled) {
            updateFormData({ variants: [] });
        }
    };
    
    const handleVariantsSave = (newVariants: any[]) => {
        // Преобразуем варианты в формат для API
        const formattedVariants = newVariants.map(v => ({
            sku: v.sku,
            price: v.price,
            stockQuantity: v.stock_quantity || 0,
            variantAttributes: v.variant_attributes,
            isDefault: v.is_default || false,
        }));
        
        updateFormData({ variants: formattedVariants });
    };
    
    // ... остальной код компонента
}
```

#### 2.3 Модификация VariantGenerator для поддержки stock

**Файл для изменения:** `/data/hostel-booking-system/frontend/svetu/src/components/Storefront/ProductVariants/VariantGenerator.tsx`

```typescript
// Добавить в состояние компонента
const [stockSettings, setStockSettings] = useState({
    defaultQuantity: 0,
    useIndividualQuantities: false,
    quantities: {} as Record<string, number>,
});

// Добавить UI для управления stock
<div className="mt-6">
    <h4 className="text-sm font-medium mb-4">Stock Quantities</h4>
    
    <div className="mb-4">
        <label className="flex items-center gap-2">
            <input
                type="checkbox"
                checked={stockSettings.useIndividualQuantities}
                onChange={(e) => setStockSettings(prev => ({
                    ...prev,
                    useIndividualQuantities: e.target.checked
                }))}
                className="checkbox checkbox-sm"
            />
            <span className="text-sm">Set individual quantities for each variant</span>
        </label>
    </div>
    
    {!stockSettings.useIndividualQuantities ? (
        <div className="form-control">
            <label className="label">
                <span className="label-text">Default quantity for all variants</span>
            </label>
            <input
                type="number"
                min="0"
                value={stockSettings.defaultQuantity}
                onChange={(e) => setStockSettings(prev => ({
                    ...prev,
                    defaultQuantity: parseInt(e.target.value) || 0
                }))}
                className="input input-bordered"
            />
        </div>
    ) : (
        <div className="space-y-2 max-h-60 overflow-y-auto">
            {generatedVariants.map((variant, index) => {
                const key = generateVariantKey(variant.variant_attributes);
                return (
                    <div key={index} className="flex items-center gap-4">
                        <span className="text-sm flex-1">
                            {Object.entries(variant.variant_attributes)
                                .map(([k, v]) => `${k}: ${v}`)
                                .join(', ')}
                        </span>
                        <input
                            type="number"
                            min="0"
                            value={stockSettings.quantities[key] || 0}
                            onChange={(e) => setStockSettings(prev => ({
                                ...prev,
                                quantities: {
                                    ...prev.quantities,
                                    [key]: parseInt(e.target.value) || 0
                                }
                            }))}
                            className="input input-bordered input-sm w-24"
                            placeholder="Qty"
                        />
                    </div>
                );
            })}
        </div>
    )}
</div>

// Модифицировать handleGenerate для включения stock
const handleGenerate = () => {
    const variantsWithStock = generatedVariants.map(variant => {
        const key = generateVariantKey(variant.variant_attributes);
        return {
            ...variant,
            stock_quantity: stockSettings.useIndividualQuantities 
                ? (stockSettings.quantities[key] || 0)
                : stockSettings.defaultQuantity
        };
    });
    
    onGenerate(variantsWithStock);
};
```

#### 2.4 Обновление сервиса создания товара

**Файл для изменения:** `/data/hostel-booking-system/frontend/svetu/src/services/storefrontProducts.ts`

```typescript
interface CreateProductWithVariantsData extends CreateProductData {
    has_variants: boolean;
    variants?: Array<{
        sku?: string;
        price?: number;
        stock_quantity: number;
        variant_attributes: Record<string, any>;
        is_default: boolean;
    }>;
    variant_settings?: {
        track_inventory: boolean;
        continue_selling: boolean;
        selected_attributes: string[];
    };
}

export const createProductWithVariants = async (
    storefrontSlug: string,
    data: CreateProductWithVariantsData
): Promise<StorefrontProduct> => {
    const response = await apiClient.post(
        `/storefronts/slug/${storefrontSlug}/products`,
        data
    );
    return response.data;
};
```

#### 2.5 Модификация PreviewStep для отображения вариантов

**Файл для изменения:** `/data/hostel-booking-system/frontend/svetu/src/components/products/steps/PreviewStep.tsx`

```typescript
// Добавить секцию для отображения вариантов
{state.formData.hasVariants && state.formData.variants.length > 0 && (
    <div className="mt-6">
        <h3 className="text-lg font-semibold mb-4">Product Variants</h3>
        <div className="overflow-x-auto">
            <table className="table table-sm">
                <thead>
                    <tr>
                        <th>Variant</th>
                        <th>SKU</th>
                        <th>Price</th>
                        <th>Stock</th>
                        <th>Default</th>
                    </tr>
                </thead>
                <tbody>
                    {state.formData.variants.map((variant, index) => (
                        <tr key={index}>
                            <td>
                                {Object.entries(variant.variantAttributes)
                                    .map(([k, v]) => `${k}: ${v}`)
                                    .join(', ')}
                            </td>
                            <td>{variant.sku || '-'}</td>
                            <td>
                                {variant.price 
                                    ? formatPrice(variant.price, state.formData.currency)
                                    : 'Base price'}
                            </td>
                            <td>{variant.stockQuantity}</td>
                            <td>
                                {variant.isDefault && (
                                    <span className="badge badge-primary badge-sm">Default</span>
                                )}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
        <div className="text-sm text-base-content/60 mt-2">
            Total variants: {state.formData.variants.length}
        </div>
    </div>
)}
```

### Этап 3: Интеграция с атрибутами категории

#### 3.1 Связь атрибутов категории с вариантами

**Backend изменения:**

1. Добавить валидацию соответствия атрибутов варианта атрибутам категории
2. Автоматически предлагать атрибуты категории для вариантов

**Frontend изменения:**

1. В `AttributesStep` сохранять выбранные атрибуты для использования в вариантах
2. В `VariantGenerator` использовать только атрибуты категории

### Этап 4: Система управления остатками

#### 4.1 Real-time обновления stock

**Backend:**
- WebSocket endpoint для подписки на изменения stock
- События при изменении остатков

**Frontend:**
- WebSocket подключение для real-time обновлений
- UI индикаторы изменения остатков

#### 4.2 Резервирование товаров

Использовать существующую таблицу `inventory_reservations`:
- При добавлении в корзину создавать резерв
- TTL для автоматического освобождения
- Учитывать резервы при отображении доступного количества

## 🔄 Порядок выполнения работ

### Фаза 1: MVP (критически важно)
1. **Backend**: Модифицировать `CreateProductRequest` для поддержки вариантов
2. **Backend**: Обновить `ProductService.CreateProduct` для создания вариантов в транзакции
3. **Frontend**: Интегрировать `VariantsStep` с `CreateProductContext`
4. **Frontend**: Добавить поля stock в `VariantGenerator`
5. **Frontend**: Обновить сервис создания товара для отправки вариантов
6. **Testing**: E2E тест создания товара с вариантами

### Фаза 2: Улучшения
1. Bulk операции с вариантами
2. Импорт/экспорт CSV с вариантами
3. Валидация атрибутов по категории
4. Real-time обновления stock

### Фаза 3: Расширенный функционал
1. Analytics по вариантам
2. A/B тестирование цен
3. Автоматические скидки
4. Интеграция с внешними системами

## 📝 Чек-лист для тестирования

### Backend тесты:
- [ ] Создание товара без вариантов работает как раньше
- [ ] Создание товара с вариантами создает все записи
- [ ] Транзакция откатывается при ошибке
- [ ] Валидация уникальности атрибутов работает
- [ ] Только один вариант может быть дефолтным
- [ ] Stock status автоматически обновляется

### Frontend тесты:
- [ ] VariantsStep сохраняет данные в контекст
- [ ] Можно указать stock для каждого варианта
- [ ] Preview отображает все варианты
- [ ] Форма отправляется с вариантами
- [ ] Ошибки валидации отображаются корректно

### E2E тесты:
- [ ] Полный flow создания товара с вариантами
- [ ] Варианты отображаются в каталоге
- [ ] Покупатель может выбрать вариант
- [ ] Stock корректно уменьшается при заказе

## 🚀 Команды для разработки

### Backend:
```bash
# Запуск тестов
cd backend && go test ./internal/proj/storefronts/...

# Проверка линтера
make lint

# Генерация Swagger документации
make generate-swagger
```

### Frontend:
```bash
# Запуск dev сервера
cd frontend/svetu && yarn dev -p 3001

# Проверка типов
yarn tsc --noEmit

# Запуск тестов
yarn test
```

### База данных:
```bash
# Применение миграций
cd backend && ./migrator up

# Откат миграций
./migrator down

# Подключение к БД
psql "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable"
```

## 📊 Метрики успеха

1. **Время создания товара с вариантами** < 30 секунд
2. **Конверсия в покупку** увеличивается на 15% для товаров с вариантами
3. **Ошибки при создании** < 1%
4. **Производительность API** < 200ms для создания товара с 10 вариантами

## 🔗 Связанная документация

- `/data/hostel-booking-system/PRODUCT_VARIANTS_DETAILED_TECHNICAL_SPEC.md` - техническая спецификация
- `/data/hostel-booking-system/backend/docs/swagger.json` - API документация
- `/data/hostel-booking-system/frontend/svetu/README.md` - Frontend документация

---

**Автор:** Claude Assistant  
**Дата создания:** 30.07.2025  
**Версия:** 1.0