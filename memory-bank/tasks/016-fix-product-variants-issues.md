# Исправление проблем с вариантами товаров

## Дата: 2025-07-30

## Проблемы
1. Дублирование SKU при создании вариантов товара
2. Неправильный URL в VariantSelectionModal
3. Отсутствующие переводы для UI
4. Ошибка TypeScript в product-utils.ts

## Решения

### 1. Исправление дублирования SKU
- **Файлы**: 
  - `frontend/svetu/src/components/products/CategoryVariantGenerator.tsx`
  - `frontend/svetu/src/components/products/SimplifiedVariantGenerator.tsx`
- **Изменения**: Добавлен timestamp и случайный суффикс к SKU для обеспечения уникальности
```typescript
const timestamp = Date.now();
const randomSuffix = Math.random().toString(36).substring(2, 6).toUpperCase();
const sku = `${variantKey}-${timestamp}-${randomSuffix}`.toUpperCase().replace(/\s+/g, '-');
```

### 2. Исправление URL в VariantSelectionModal
- **Файл**: `frontend/svetu/src/components/cart/VariantSelectionModal.tsx`
- **Изменение**: Исправлен путь API с `/api/v1/storefronts/${storefrontSlug}/products/${productId}` на `/api/v1/storefronts/slug/${storefrontSlug}/products/${productId}`

### 3. Добавление недостающих переводов
- **Файлы**: 
  - `frontend/svetu/src/messages/en.json`
  - `frontend/svetu/src/messages/ru.json`
  - `frontend/svetu/src/messages/sr.json`
- **Добавлены переводы**:
  - `common.stockStatus` - статусы наличия товара
  - `storefronts.products.stock` - статусы для компонента ProductCard
  - `storefronts.products.default` - значение по умолчанию

### 4. Исправление ошибки TypeScript
- **Файл**: `frontend/svetu/src/utils/product-utils.ts`
- **Изменение**: Удалена проверка `product.storefront?.ownerId`, так как поле `ownerId` не существует в интерфейсе `UnifiedStorefrontInfo`

## Результат
- Товар с вариантами успешно создается (ID: 88 с 9 вариантами)
- Все переводы работают корректно
- Проект успешно собирается без ошибок
- SKU теперь уникальны благодаря добавлению временной метки и случайного суффикса

## Проверки выполнены
- ✅ yarn lint
- ✅ yarn format
- ✅ yarn build