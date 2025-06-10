# Анализ процесса выбора атрибутов при создании товаров

## Обзор

Процесс создания товара в маркетплейсе включает этап выбора и заполнения атрибутов категории. Система атрибутов позволяет структурированно хранить дополнительную информацию о товарах.

## Frontend (React/Next.js)

### 1. Компонент AttributesStep

**Расположение**: `/frontend/svetu/src/components/create-listing/steps/AttributesStep.tsx`

#### Основные функции:

1. **Загрузка атрибутов категории**:
```typescript
const response = await MarketplaceService.getCategoryAttributes(state.category.id);
```

2. **Типы атрибутов**:
- `text` - текстовое поле
- `number` - числовое поле
- `select` - выпадающий список
- `boolean` - чекбокс
- `multiselect` - множественный выбор

3. **Структура данных атрибута**:
```typescript
interface AttributeFormData {
  attribute_id: number;
  attribute_name: string;
  display_name: string;
  attribute_type: string;
  text_value?: string;
  numeric_value?: number;
  boolean_value?: boolean;
  json_value?: any;
  display_value?: string;
  unit?: string;
}
```

4. **Валидация**:
- Проверка обязательных атрибутов
- Блокировка кнопки "Продолжить" при незаполненных обязательных полях

5. **UX особенности**:
- Автоматическая сортировка атрибутов по `sort_order`
- Фильтрация дубликатов
- Отображение обязательных полей с красной звездочкой
- Региональные подсказки

### 2. Контекст CreateListingContext

**Расположение**: `/frontend/svetu/src/contexts/CreateListingContext.tsx`

- Управление состоянием атрибутов через Redux-подобный reducer
- Action `SET_ATTRIBUTES` для обновления атрибутов
- Атрибуты хранятся в `state.attributes` как объект с ключами по ID атрибута

### 3. Сервис ListingsService

**Расположение**: `/frontend/svetu/src/services/listings.ts`

При создании товара атрибуты преобразуются и отправляются на backend:

```typescript
request.attributes = Object.entries(data.attributes).map(
  ([_key, value]) => {
    const attributeData = value as any;
    return {
      attribute_id: attributeData.attribute_id,
      attribute_name: attributeData.attribute_name,
      display_name: attributeData.display_name,
      attribute_type: attributeData.attribute_type,
      text_value: attributeData.text_value,
      numeric_value: attributeData.numeric_value,
      boolean_value: attributeData.boolean_value,
      json_value: attributeData.json_value,
      unit: attributeData.unit,
    };
  }
);
```

## Backend (Go)

### 1. API Endpoint

**Расположение**: `/backend/internal/proj/marketplace/handler/listings.go`

**Endpoint**: `POST /api/v1/marketplace/listings`

#### Обработка атрибутов:

1. **Парсинг атрибутов из запроса**:
```go
func processAttributesFromRequest(requestBody map[string]interface{}, listing *models.MarketplaceListing)
```

2. **Преобразование типов**:
- Поддержка различных форматов входных данных
- Автоматическое определение типа значения
- Обработка multiselect как JSON

3. **Сохранение атрибутов**:
```go
if listing.Attributes != nil && len(listing.Attributes) > 0 {
    // Фильтрация дубликатов
    uniqueAttrs := make(map[int]models.ListingAttributeValue)
    for _, attr := range listing.Attributes {
        uniqueAttrs[attr.AttributeID] = attr
    }
    
    // Сохранение в БД
    if err := s.SaveListingAttributes(ctx, listingID, filteredAttrs); err != nil {
        log.Printf("Error saving attributes for listing %d: %v", listingID, err)
    }
}
```

### 2. Модели данных

**Расположение**: `/backend/internal/domain/models/category_attributes.go`

#### ListingAttributeValue:
```go
type ListingAttributeValue struct {
    ListingID     int
    AttributeID   int
    AttributeName string
    DisplayName   string
    AttributeType string
    TextValue     *string
    NumericValue  *float64
    BooleanValue  *bool
    JSONValue     json.RawMessage
    DisplayValue  string
    Unit          string
    Translations  map[string]string
    OptionTranslations map[string]map[string]string
}
```

### 3. Сохранение в БД

**Расположение**: `/backend/internal/proj/marketplace/service/attribute_admin.go`

Метод `SaveListingAttributes`:
1. Начинает транзакцию
2. Удаляет старые атрибуты товара
3. Сохраняет новые атрибуты с правильными типами значений
4. Поддерживает различные типы: text, numeric, boolean, json

## База данных

### Таблица listing_attribute_values

```sql
CREATE TABLE listing_attribute_values (
    listing_id integer,
    attribute_id integer,
    value_type varchar,
    text_value text,
    numeric_value numeric,
    boolean_value boolean,
    json_value jsonb,
    unit varchar
);
```

## Процесс создания товара с атрибутами

1. **Выбор категории** → загрузка атрибутов категории
2. **Заполнение атрибутов** → валидация обязательных полей
3. **Отправка на сервер** → преобразование данных
4. **Сохранение в БД** → транзакционное сохранение
5. **Индексация в OpenSearch** → для поиска по атрибутам

## Проблемы и рекомендации

### Выявленные проблемы:

1. **Отсутствие кеширования атрибутов** - каждый раз загружаются с сервера
2. **Нет валидации на backend** - только проверка обязательных полей на frontend
3. **Дублирование атрибутов** - требуется фильтрация на уровне сервиса
4. **Отсутствие динамической валидации** - не используются validation_rules из БД

### Рекомендации:

1. **Добавить кеширование атрибутов категории**
2. **Реализовать валидацию на backend согласно validation_rules**
3. **Оптимизировать загрузку атрибутов** - использовать batch запросы
4. **Добавить поддержку custom_component** для специальных UI компонентов
5. **Улучшить UX** - добавить подсказки для каждого атрибута
6. **Реализовать автозаполнение** для часто используемых значений

## Примеры кода

### Frontend - рендеринг атрибута:
```typescript
case 'select':
  const selectOptions = getOptionValues(attribute.options);
  return (
    <select
      className="select select-bordered"
      value={value}
      onChange={(e) => handleInputChange(
        attribute.id,
        e.target.value,
        attribute.attribute_type
      )}
    >
      <option value="">{t('common.select')}</option>
      {selectOptions.map((option) => (
        <option key={option} value={option}>
          {option}
        </option>
      ))}
    </select>
  );
```

### Backend - обработка multiselect:
```go
case "multiselect":
    if jsonValues, ok := attrMap["json_value"]; ok {
        jsonBytes, err := json.Marshal(jsonValues)
        if err == nil {
            attr.JSONValue = jsonBytes
        }
    }
```