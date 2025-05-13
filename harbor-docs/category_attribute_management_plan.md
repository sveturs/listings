# План системы управления категориями и атрибутами для маркетплейса

## 1. Текущее состояние системы

Система категорий и атрибутов имеет структуру:

- **Категории** хранятся в таблице `marketplace_categories` с поддержкой иерархий и многоязычности
- **Атрибуты** хранятся в `category_attributes` и поддерживают различные типы (text, number, select, boolean)
- **Связи** между категориями и атрибутами реализованы через `category_attribute_mapping`
- **Значения атрибутов** хранятся в `listing_attribute_values` с поддержкой различных типов данных
- **Переводы** реализованы через таблицу `translations` с поддержкой многоязычности

## 2. Административный интерфейс (Backend)

### 2.1. Новые эндпоинты API

```go
// Категории
POST /api/admin/categories                 // Создание категории
GET /api/admin/categories                  // Получение списка категорий
GET /api/admin/categories/:id              // Получение категории по ID
PUT /api/admin/categories/:id              // Обновление категории
DELETE /api/admin/categories/:id           // Удаление категории
POST /api/admin/categories/:id/reorder      // Изменение порядка категорий
PUT /api/admin/categories/:id/move         // Перемещение категории в иерархии

// Атрибуты
POST /api/admin/attributes                 // Создание атрибута
GET /api/admin/attributes                  // Получение списка атрибутов
GET /api/admin/attributes/:id              // Получение атрибута по ID
PUT /api/admin/attributes/:id              // Обновление атрибута
DELETE /api/admin/attributes/:id           // Удаление атрибута
POST /api/admin/attributes/bulk-update     // Массовое обновление атрибутов

// Связи категорий и атрибутов
POST /api/admin/categories/:id/attributes  // Привязка атрибута к категории
DELETE /api/admin/categories/:id/attributes/:attr_id  // Удаление связи
PUT /api/admin/categories/:id/attributes/:attr_id     // Обновление настроек связи
```

### 2.2. Обновление сервисного слоя

Добавить в интерфейс `MarketplaceServiceInterface` новые методы:

```go
type MarketplaceServiceInterface interface {
    // Существующие методы...
    
    // Новые методы для управления категориями
    CreateCategory(ctx context.Context, category *models.MarketplaceCategory) (int, error)
    UpdateCategory(ctx context.Context, category *models.MarketplaceCategory) error
    DeleteCategory(ctx context.Context, id int) error
    ReorderCategories(ctx context.Context, orderedIDs []int) error
    MoveCategory(ctx context.Context, id int, newParentID int) error
    
    // Новые методы для управления атрибутами
    CreateAttribute(ctx context.Context, attribute *models.CategoryAttribute) (int, error)
    UpdateAttribute(ctx context.Context, attribute *models.CategoryAttribute) error
    DeleteAttribute(ctx context.Context, id int) error
    GetAttributeByID(ctx context.Context, id int) (*models.CategoryAttribute, error)
    
    // Новые методы для управления связями
    AddAttributeToCategory(ctx context.Context, categoryID int, attributeID int, isRequired bool) error
    RemoveAttributeFromCategory(ctx context.Context, categoryID int, attributeID int) error
    UpdateAttributeCategory(ctx context.Context, categoryID int, attributeID int, isRequired bool, isEnabled bool) error
    InvalidateAttributeCache(ctx context.Context, categoryID int) error
}
```

### 2.3. Обработка специфики кастомных UI компонентов

Добавить поле для указания кастомного UI-компонента в модель `CategoryAttribute`:

```go
type CategoryAttribute struct {
    // Существующие поля...
    CustomComponent string `json:"custom_component,omitempty"`
}
```

И аналогично для категорий, указывая необходимость кастомного UI:

```go
type MarketplaceCategory struct {
    // Существующие поля...
    HasCustomUI bool `json:"has_custom_ui,omitempty"`
    CustomUIComponent string `json:"custom_ui_component,omitempty"`
}
```

## 3. Административный интерфейс (Frontend)

### 3.1. Новые страницы и компоненты

- **Страница управления категориями** (`/admin/categories`)
  - Древовидная структура категорий с возможностью drag-and-drop
  - Контекстное меню для редактирования и удаления
  - Модальное окно для создания/редактирования категории
  - Функция поиска по категориям

- **Страница управления атрибутами** (`/admin/attributes`)
  - Таблица со списком атрибутов и фильтрами
  - Модальное окно для создания/редактирования атрибута
  - Функция поиска по атрибутам

- **Страница связей категорий и атрибутов** (`/admin/category/:id/attributes`)
  - Список привязанных атрибутов к конкретной категории
  - Возможность добавления/удаления/настройки атрибутов
  - Наследование атрибутов от родительских категорий

### 3.2. Компоненты для редактирования атрибутов

```tsx
// Компонент для редактирования атрибута
const AttributeForm: React.FC<{attribute?: CategoryAttribute}> = ({ attribute }) => {
  // Логика редактирования атрибута
};

// Компонент для выбора типа атрибута
const AttributeTypeSelector: React.FC<{value: string, onChange: (type: string) => void}> = ({ value, onChange }) => {
  // Выбор типа атрибута с подсказками
};

// Редактор опций для select/multiselect
const AttributeOptionsEditor: React.FC<{options: any, onChange: (options: any) => void}> = ({ options, onChange }) => {
  // Редактирование опций выбора
};

// Компонент для настройки валидационных правил
const ValidationRulesEditor: React.FC<{rules: any, onChange: (rules: any) => void}> = ({ rules, onChange }) => {
  // Настройка правил валидации
};
```

### 3.3. Компоненты для управления переводами

```tsx
// Компонент для многоязычного редактирования
const TranslationEditor: React.FC<{
  entityType: string,
  entityId: number,
  fieldName: string,
  translations: Record<string, string>,
  onChange: (translations: Record<string, string>) => void
}> = (props) => {
  // Редактирование переводов для разных языков
};

// Компонент для переводов опций атрибутов
const OptionTranslationEditor: React.FC<{
  attributeId: number,
  options: string[],
  translations: Record<string, Record<string, string>>,
  onChange: (translations: Record<string, Record<string, string>>) => void
}> = (props) => {
  // Редактирование переводов для опций атрибутов
};
```

## 4. Реестр компонентов для кастомных UI

### 4.1. Система регистрации компонентов

```tsx
// types.ts
type CategoryUiComponentProps = {
  categoryId: number;
  values: Record<string, any>;
  onChange: (values: Record<string, any>) => void;
};

type AttributeComponentProps = {
  attribute: CategoryAttribute;
  value: any;
  onChange: (value: any) => void;
};

// registry.ts
class ComponentRegistry {
  private static categoryComponents: Record<string, React.ComponentType<CategoryUiComponentProps>> = {};
  private static attributeComponents: Record<string, React.ComponentType<AttributeComponentProps>> = {};

  static registerCategoryComponent(name: string, component: React.ComponentType<CategoryUiComponentProps>) {
    this.categoryComponents[name] = component;
  }

  static registerAttributeComponent(name: string, component: React.ComponentType<AttributeComponentProps>) {
    this.attributeComponents[name] = component;
  }

  static getCategoryComponent(name: string): React.ComponentType<CategoryUiComponentProps> | null {
    return this.categoryComponents[name] || null;
  }

  static getAttributeComponent(name: string): React.ComponentType<AttributeComponentProps> | null {
    return this.attributeComponents[name] || null;
  }
}
```

### 4.2. Фабрика компонентов

```tsx
// CategoryAttributeComponent.tsx
const CategoryAttributeComponent: React.FC<{
  attribute: CategoryAttribute;
  value: any;
  onChange: (value: any) => void;
}> = ({ attribute, value, onChange }) => {
  // Если указан кастомный компонент
  if (attribute.CustomComponent) {
    const CustomComponent = ComponentRegistry.getAttributeComponent(attribute.CustomComponent);
    if (CustomComponent) {
      return <CustomComponent attribute={attribute} value={value} onChange={onChange} />;
    }
  }

  // Стандартные компоненты по типу атрибута
  switch (attribute.AttributeType) {
    case 'text':
      return <TextAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    case 'number':
      return <NumberAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    case 'select':
      return <SelectAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    case 'multiselect':
      return <MultiSelectAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    case 'boolean':
      return <BooleanAttributeInput attribute={attribute} value={value} onChange={onChange} />;
    default:
      return <TextAttributeInput attribute={attribute} value={value} onChange={onChange} />;
  }
};
```

## 5. Интеграция с существующей системой

### 5.1. Обновление миграций БД

```sql
-- Добавление полей для управления кастомными компонентами
ALTER TABLE marketplace_categories 
ADD COLUMN has_custom_ui BOOLEAN DEFAULT FALSE,
ADD COLUMN custom_ui_component VARCHAR(255);

ALTER TABLE category_attributes 
ADD COLUMN custom_component VARCHAR(255);

-- Индексы для оптимизации запросов в админке
CREATE INDEX IF NOT EXISTS idx_category_attributes_name ON category_attributes(name);
CREATE INDEX IF NOT EXISTS idx_categories_path ON marketplace_categories USING GIN (category_path);
```

### 5.2. Обновление логики поиска и фильтрации

Модифицируем методы для учета новых полей при поиске и фильтрации. Ключевые функции:

- `GetCategoryAttributes` - обновим для возврата информации о кастомных компонентах
- `SearchListingsAdvanced` - обновим для поддержки фильтрации по кастомным атрибутам

## 6. Производительность и безопасность

### 6.1. Кеширование и пагинация

- Расширим существующую логику кеширования атрибутов
- Внедрим пагинацию для списков категорий и атрибутов в админке
- Оптимизируем запросы для больших наборов данных

### 6.2. Проверки доступа

```go
// Middleware для проверки прав администратора
func AdminAuthMiddleware(service users.ServiceInterface) func(c *fiber.Ctx) error {
    return func(c *fiber.Ctx) error {
        // Проверка прав администратора
        userID, ok := c.Locals("user_id").(int)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        }

        isAdmin, err := service.IsUserAdmin(c.Context(), userID)
        if err != nil || !isAdmin {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Admin access required",
            })
        }

        return c.Next()
    }
}
```

## 7. Поэтапное внедрение

### 7.1. Этап 1: Базовый интерфейс администрирования

- Создание CRUD для категорий
- Создание CRUD для атрибутов
- Базовое управление связями

### 7.2. Этап 2: Кастомные компоненты

- Внедрение реестра компонентов
- Создание базовых кастомных компонентов для автомобилей и недвижимости
- Интеграция с существующей системой фильтров

### 7.3. Этап 3: Расширенные возможности

- Перетаскивание для сортировки категорий и атрибутов
- Массовое редактирование
- Улучшенное управление переводами
- Статистика использования атрибутов

## 8. Интеграция с существующим UI фильтров

Создадим фабрику фильтров, которая будет динамически строить компоненты фильтрации на основе доступных атрибутов категории:

```tsx
// FilterFactory.tsx
type FilterFactoryProps = {
  categoryId: number;
  attributes: CategoryAttribute[];
  values: Record<string, any>;
  onChange: (values: Record<string, any>) => void;
};

const FilterFactory: React.FC<FilterFactoryProps> = ({ categoryId, attributes, values, onChange }) => {
  // Получаем информацию о категории
  const category = useCategory(categoryId);
  
  if (category?.HasCustomUI && category?.CustomUIComponent) {
    // Используем кастомный компонент фильтрации
    const CustomFilter = ComponentRegistry.getCategoryComponent(category.CustomUIComponent);
    if (CustomFilter) {
      return <CustomFilter categoryId={categoryId} values={values} onChange={onChange} />;
    }
  }
  
  // Стандартная логика генерации фильтров на основе атрибутов
  return (
    <div className="filters-container">
      {attributes.map(attr => 
        attr.IsFilterable ? (
          <div key={attr.ID} className="filter-item">
            <label>{attr.DisplayName}</label>
            <CategoryAttributeComponent 
              attribute={attr}
              value={values[attr.Name]} 
              onChange={(value) => {
                onChange({...values, [attr.Name]: value});
              }}
            />
          </div>
        ) : null
      )}
    </div>
  );
};
```
