# План рефакторинга системы транзакций для товаров с вариантами

## Текущая проблема

### Описание проблемы
При создании товара с вариантами через API витрин возникает ошибка 500 из-за конфликта транзакций.

### Причина проблемы
1. Метод `ProductService.CreateProduct` создает транзакцию для атомарного создания товара
2. Внутри этой транзакции вызывается `VariantService.BulkCreateVariants`
3. `BulkCreateVariants` вызывает `VariantRepository.CreateVariant` для каждого варианта
4. `CreateVariant` пытается создать свою собственную транзакцию через `BeginTxx`
5. Возникает конфликт - нельзя создать вложенную транзакцию в PostgreSQL

### Текущее временное решение
Транзакции в `ProductService.CreateProduct` были отключены (закомментированы), что позволяет создавать товары с вариантами, но нарушает атомарность операции.

## Цель рефакторинга

### Основные цели
1. **Обеспечить атомарность** - создание товара и всех его вариантов должно происходить в одной транзакции
2. **Сохранить целостность данных** - при ошибке создания любого варианта должен откатываться весь товар
3. **Поддержать масштабируемость** - решение должно работать для bulk операций
4. **Минимизировать изменения API** - публичные интерфейсы должны остаться неизменными

## Детальное техническое решение

### 1. Создание методов с поддержкой транзакций

#### 1.1 Обновление интерфейса VariantRepository
```go
// internal/proj/storefront/repository/variant_repository.go

type VariantRepository interface {
    // Существующие методы
    CreateVariant(ctx context.Context, req *types.CreateVariantRequest) (*types.ProductVariant, error)
    
    // Новые методы с поддержкой транзакций
    CreateVariantTx(ctx context.Context, tx *sqlx.Tx, req *types.CreateVariantRequest) (*types.ProductVariant, error)
    BulkCreateVariantsTx(ctx context.Context, tx *sqlx.Tx, productID int, variants []types.CreateVariantRequest) ([]*types.ProductVariant, error)
}
```

#### 1.2 Реализация методов с транзакциями
```go
func (r *VariantRepository) CreateVariantTx(ctx context.Context, tx *sqlx.Tx, req *types.CreateVariantRequest) (*types.ProductVariant, error) {
    // Вся логика CreateVariant, но используя переданную транзакцию tx
    // вместо создания новой
    
    // If this is set as default, unset other defaults for this product
    if req.IsDefault {
        _, err := tx.ExecContext(ctx,
            "UPDATE storefront_product_variants SET is_default = false WHERE product_id = $1",
            req.ProductID)
        if err != nil {
            return nil, err
        }
    }
    
    // ... остальная логика создания варианта
}
```

### 2. Обновление VariantService

#### 2.1 Добавление методов с поддержкой транзакций
```go
// internal/proj/storefronts/service/variant_service.go

type VariantService interface {
    // Существующие методы
    BulkCreateVariants(ctx context.Context, productID int, variants []variantTypes.CreateVariantRequest) ([]*variantTypes.ProductVariant, error)
    
    // Новые методы
    BulkCreateVariantsTx(ctx context.Context, tx *sqlx.Tx, productID int, variants []variantTypes.CreateVariantRequest) ([]*variantTypes.ProductVariant, error)
}

func (s *VariantServiceImpl) BulkCreateVariantsTx(ctx context.Context, tx *sqlx.Tx, productID int, variants []variantTypes.CreateVariantRequest) ([]*variantTypes.ProductVariant, error) {
    // Используем транзакционные методы репозитория
    return s.variantRepo.BulkCreateVariantsTx(ctx, tx, productID, variants)
}
```

### 3. Обновление Storage интерфейса

#### 3.1 Добавление транзакционных методов
```go
// internal/proj/storefronts/service/storage.go

type Storage interface {
    // Существующие методы
    CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error)
    BeginTx(ctx context.Context) (Transaction, error)
    
    // Новые методы
    CreateStorefrontProductTx(ctx context.Context, tx *sqlx.Tx, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error)
}
```

### 4. Обновление ProductService

#### 4.1 Использование единой транзакции
```go
func (s *ProductService) CreateProduct(ctx context.Context, storefrontID, userID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
    // Validate ownership
    if err := s.ValidateStorefrontOwnership(ctx, storefrontID, userID); err != nil {
        return nil, fmt.Errorf("ownership validation failed: %w", err)
    }

    // Validate request
    if err := s.validateCreateRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    // Validate variants if provided
    if req.HasVariants {
        if err := s.validateVariants(req); err != nil {
            return nil, fmt.Errorf("invalid variants: %w", err)
        }
    }

    // Start transaction
    tx, err := s.storage.BeginTx(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %w", err)
    }
    defer func() { _ = tx.Rollback() }()

    // Create product using transaction
    product, err := s.storage.CreateStorefrontProductTx(ctx, tx, storefrontID, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create product: %w", err)
    }

    // Create variants if provided
    if req.HasVariants && len(req.Variants) > 0 && s.variantService != nil {
        // Преобразуем варианты в нужный формат
        variantRequests := make([]variantTypes.CreateVariantRequest, len(req.Variants))
        // ... преобразование ...

        // Используем транзакционный метод
        createdVariants, err := s.variantService.BulkCreateVariantsTx(ctx, tx, product.ID, variantRequests)
        if err != nil {
            return nil, fmt.Errorf("failed to create variants: %w", err)
        }

        // Добавляем варианты к продукту
        product.Variants = convertVariantsToProductVariants(createdVariants)
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    // Index product in OpenSearch (after transaction is committed)
    if s.searchRepo != nil {
        go s.indexProductWithVariants(ctx, product)
    }

    return product, nil
}
```

### 5. Поэтапная миграция

#### Фаза 1: Добавление новых методов
1. Создать новые методы с суффиксом `Tx` не нарушая существующие
2. Протестировать новые методы изолированно

#### Фаза 2: Миграция критических путей
1. Обновить `CreateProduct` для использования транзакционных методов
2. Обновить bulk операции для поддержки транзакций

#### Фаза 3: Полная миграция
1. Обновить все оставшиеся методы
2. Пометить старые методы как deprecated
3. Постепенно удалить старые методы

### 6. Обработка ошибок и откат

```go
// Пример обработки ошибок с детальным логированием
if err := s.variantService.BulkCreateVariantsTx(ctx, tx, product.ID, variantRequests); err != nil {
    logger.Error().
        Err(err).
        Int("product_id", product.ID).
        Int("variant_count", len(variantRequests)).
        Msg("Failed to create variants, rolling back entire product creation")
    
    // Транзакция автоматически откатится через defer
    return nil, fmt.Errorf("failed to create variants: %w", err)
}
```

## Тестирование

### Unit тесты
1. Тест создания товара с вариантами в одной транзакции
2. Тест отката при ошибке создания варианта
3. Тест bulk операций с транзакциями

### Интеграционные тесты
1. Тест создания товара с 10+ вариантами
2. Тест параллельного создания товаров
3. Тест восстановления после сбоя

## Риски и митигация

### Риск 1: Длительные транзакции
**Митигация**: Добавить таймауты и мониторинг длительности транзакций

### Риск 2: Deadlocks при параллельных операциях
**Митигация**: Использовать правильный порядок блокировок и retry логику

### Риск 3: Увеличение сложности кода
**Митигация**: Хорошая документация и примеры использования

## Временная шкала

1. **Неделя 1**: Создание транзакционных методов в репозиториях
2. **Неделя 2**: Обновление сервисов и тестирование
3. **Неделя 3**: Миграция критических путей
4. **Неделя 4**: Полное тестирование и развертывание

## Заключение

Данный план позволит восстановить атомарность операций создания товаров с вариантами, сохранив при этом производительность и масштабируемость системы. Поэтапный подход минимизирует риски и позволяет откатиться на любом этапе.