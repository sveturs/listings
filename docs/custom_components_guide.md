# Руководство по кастомным компонентам для категорий и атрибутов

Данное руководство описывает систему кастомных компонентов для категорий и атрибутов в маркетплейсе Sve Tu Platform.

## Содержание

1. [Введение](#введение)
2. [Миграции и модели данных](#миграции-и-модели-данных)
3. [Серверный API](#серверный-api)
4. [Система регистрации компонентов](#система-регистрации-компонентов)
5. [Создание кастомных компонентов атрибутов](#создание-кастомных-компонентов-атрибутов)
6. [Создание кастомных компонентов категорий](#создание-кастомных-компонентов-категорий)
7. [Интеграция с системой фильтрации](#интеграция-с-системой-фильтрации)

## Введение

Система кастомных компонентов позволяет создавать специализированные интерфейсы для различных категорий товаров и их атрибутов. Это полезно для категорий с уникальными требованиями к вводу данных или отображению, например:

- Автомобили с зависимостью между маркой и моделью
- Недвижимость с интерактивными фильтрами по местоположению
- Товары со специфическим выбором цвета, размеров и т.д.

## Миграции и модели данных

### Таблицы и поля

Для поддержки кастомных компонентов в существующие таблицы добавлены новые поля:

```sql
-- В таблице категорий
ALTER TABLE marketplace_categories 
ADD COLUMN has_custom_ui BOOLEAN DEFAULT FALSE,
ADD COLUMN custom_ui_component VARCHAR(255);

-- В таблице атрибутов
ALTER TABLE category_attributes 
ADD COLUMN custom_component VARCHAR(255);
```

### Модели данных

Соответствующие поля добавлены в модели:

```go
// В структуре MarketplaceCategory
type MarketplaceCategory struct {
    // Существующие поля...
    HasCustomUI       bool   `json:"has_custom_ui,omitempty"`
    CustomUIComponent string `json:"custom_ui_component,omitempty"`
}

// В структуре CategoryAttribute
type CategoryAttribute struct {
    // Существующие поля...
    CustomComponent string `json:"custom_component,omitempty"`
}
```

## Серверный API

Для управления категориями и атрибутами с кастомными компонентами реализован административный API:

### API категорий

```
POST   /api/admin/categories                 // Создание категории
GET    /api/admin/categories                 // Получение списка категорий
GET    /api/admin/categories/:id             // Получение категории по ID
PUT    /api/admin/categories/:id             // Обновление категории
DELETE /api/admin/categories/:id             // Удаление категории
POST   /api/admin/categories/:id/reorder     // Изменение порядка категорий
PUT    /api/admin/categories/:id/move        // Перемещение категории в иерархии
```

### API атрибутов

```
POST   /api/admin/attributes                 // Создание атрибута
GET    /api/admin/attributes                 // Получение списка атрибутов
GET    /api/admin/attributes/:id             // Получение атрибута по ID
PUT    /api/admin/attributes/:id             // Обновление атрибута
DELETE /api/admin/attributes/:id             // Удаление атрибута
POST   /api/admin/attributes/bulk-update     // Массовое обновление атрибутов
```

### API связей

```
POST   /api/admin/categories/:id/attributes         // Привязка атрибута к категории
DELETE /api/admin/categories/:id/attributes/:attr_id // Удаление связи
PUT    /api/admin/categories/:id/attributes/:attr_id // Обновление настроек связи
```

## Система регистрации компонентов

Для управления кастомными компонентами создана система регистрации на клиенте:

```typescript
// Типы пропсов для компонентов
export interface CategoryUiComponentProps {
  categoryId: number;
  values: Record<string, any>;
  onChange: (values: Record<string, any>) => void;
}

export interface AttributeComponentProps {
  attribute: {
    id: number;
    name: string;
    display_name: string;
    attribute_type: string;
    options?: any;
    custom_component?: string;
    // ... другие свойства
  };
  value: any;
  onChange: (value: any) => void;
}

// Регистрация компонентов
ComponentRegistry.registerCategoryComponent('AutoCategoryUI', AutoCategoryComponent);
ComponentRegistry.registerAttributeComponent('ColorPicker', ColorPickerAttribute);

// Получение компонентов
const CustomComponent = ComponentRegistry.getCategoryComponent('AutoCategoryUI');
const ColorPickerComponent = ComponentRegistry.getAttributeComponent('ColorPicker');
```

## Создание кастомных компонентов атрибутов

Кастомные компоненты атрибутов должны реализовывать интерфейс `AttributeComponentProps`:

```tsx
import { AttributeComponentProps } from '../registry/ComponentRegistry';

const ColorPickerAttribute: React.FC<AttributeComponentProps> = ({ 
  attribute, 
  value, 
  onChange 
}) => {
  // Реализация компонента
  return (
    <div>
      {/* Ваш UI для выбора цвета */}
      <Button 
        onClick={() => onChange('red')}
        style={{ backgroundColor: 'red' }}
      >
        Красный
      </Button>
    </div>
  );
};

export default ColorPickerAttribute;
```

## Создание кастомных компонентов категорий

Кастомные компоненты категорий должны реализовывать интерфейс `CategoryUiComponentProps`:

```tsx
import { CategoryUiComponentProps } from '../registry/ComponentRegistry';

const AutoCategoryComponent: React.FC<CategoryUiComponentProps> = ({ 
  categoryId, 
  values, 
  onChange 
}) => {
  // Реализация компонента для категории автомобилей
  return (
    <div>
      <TextField 
        label="Марка"
        value={values.make || ''}
        onChange={(e) => onChange({ ...values, make: e.target.value })}
      />
      
      {/* Другие поля формы */}
    </div>
  );
};

export default AutoCategoryComponent;
```

## Интеграция с системой фильтрации

Система фильтрации автоматически использует кастомные компоненты для категорий и атрибутов:

```tsx
// FilterFactory выбирает кастомные компоненты, если они указаны
const FilterFactory: React.FC<FilterFactoryProps> = ({ 
  categoryId, 
  values, 
  onChange 
}) => {
  // Получение данных о категории
  const category = useCategory(categoryId);
  
  // Проверка наличия кастомного компонента
  if (category?.HasCustomUI && category?.CustomUIComponent) {
    // Используем кастомный компонент для категории
    const CustomFilter = ComponentRegistry.getCategoryComponent(category.CustomUIComponent);
    if (CustomFilter) {
      return <CustomFilter categoryId={categoryId} values={values} onChange={onChange} />;
    }
  }
  
  // Иначе используем стандартные компоненты для атрибутов
  return (
    <div>
      {attributes.map(attr => (
        <AttributeComponentFactory 
          attribute={attr}
          value={values[attr.name]} 
          onChange={value => handleChange(attr.name, value)}
        />
      ))}
    </div>
  );
};
```

## Примеры готовых кастомных компонентов

В системе уже реализованы следующие кастомные компоненты:

1. **AutoCategoryComponent** - интерфейс для категории автомобилей с зависимостью марка-модель
2. **ColorPickerAttribute** - визуальный выбор цвета для атрибутов типа "цвет"

Для добавления новых компонентов:

1. Создайте новый компонент в директории `components/marketplace/registry/custom`
2. Реализуйте необходимый интерфейс (CategoryUiComponentProps или AttributeComponentProps)
3. Зарегистрируйте компонент в `components/marketplace/registry/index.ts`
4. Укажите имя компонента в админке при редактировании категории или атрибута

## Полезные советы

1. При разработке сложных компонентов соблюдайте принцип чистых функций
2. Используйте контекст категории для зависимых атрибутов
3. Для кэширования данных используйте React.memo и useMemo
4. Создавайте компоненты с хорошей мобильной адаптацией

## Заключение

Система кастомных компонентов позволяет расширять функциональность категорий и атрибутов без изменения основного кода. Это обеспечивает гибкость при добавлении новых типов товаров и услуг на маркетплейсе.