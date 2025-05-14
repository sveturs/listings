/**
 * Индексный файл для регистрации компонентов
 */
import ComponentRegistry from './ComponentRegistry';
import AttributeComponentFactory from './AttributeComponentFactory';
import CategoryComponentFactory from './CategoryComponentFactory';

// Импортируем кастомные компоненты
import AutoCategoryComponent from './custom/AutoCategoryComponent';
import ColorPickerAttribute from './custom/ColorPickerAttribute';

// Регистрируем кастомные компоненты для категорий
ComponentRegistry.registerCategoryComponent('AutoCategoryUI', AutoCategoryComponent);

// Регистрируем кастомные компоненты для атрибутов
ComponentRegistry.registerAttributeComponent('ColorPicker', ColorPickerAttribute);

// Экспортируем все необходимые компоненты и классы
export {
  ComponentRegistry,
  AttributeComponentFactory,
  CategoryComponentFactory,
};

// Экспортируем типы для использования в проекте
export type { 
  CategoryUiComponentProps, 
  AttributeComponentProps 
} from './ComponentRegistry';