/**
 * Система регистрации кастомных UI компонентов для категорий и атрибутов
 * Позволяет динамически регистрировать и получать компоненты для отображения
 * специфичных для категорий форм и фильтров
 */

import React from 'react';

// Типы пропсов для компонентов категорий
export interface CategoryUiComponentProps {
  categoryId: number;
  values: Record<string, any>;
  onChange: (values: Record<string, any>) => void;
}

// Типы пропсов для компонентов атрибутов
export interface AttributeComponentProps {
  attribute: {
    id: number;
    name: string;
    display_name: string;
    attribute_type: string;
    options?: any;
    validation_rules?: any;
    custom_component?: string;
  };
  value: any;
  onChange: (value: any) => void;
}

/**
 * Класс для регистрации и получения кастомных компонентов
 */
class ComponentRegistry {
  // Реестр компонентов для категорий
  private static categoryComponents: Record<string, React.ComponentType<CategoryUiComponentProps>> = {};
  
  // Реестр компонентов для атрибутов
  private static attributeComponents: Record<string, React.ComponentType<AttributeComponentProps>> = {};

  /**
   * Регистрирует компонент для категории
   * @param name Уникальное имя компонента
   * @param component Компонент React
   */
  static registerCategoryComponent(name: string, component: React.ComponentType<CategoryUiComponentProps>) {
    this.categoryComponents[name] = component;
    console.log(`Registered category component: ${name}`);
  }

  /**
   * Регистрирует компонент для атрибута
   * @param name Уникальное имя компонента
   * @param component Компонент React
   */
  static registerAttributeComponent(name: string, component: React.ComponentType<AttributeComponentProps>) {
    this.attributeComponents[name] = component;
    console.log(`Registered attribute component: ${name}`);
  }

  /**
   * Получает компонент для категории по имени
   * @param name Имя компонента
   * @returns Компонент React или null, если компонент не найден
   */
  static getCategoryComponent(name: string): React.ComponentType<CategoryUiComponentProps> | null {
    return this.categoryComponents[name] || null;
  }

  /**
   * Получает компонент для атрибута по имени
   * @param name Имя компонента
   * @returns Компонент React или null, если компонент не найден
   */
  static getAttributeComponent(name: string): React.ComponentType<AttributeComponentProps> | null {
    return this.attributeComponents[name] || null;
  }

  /**
   * Возвращает все зарегистрированные компоненты для категорий
   */
  static getAllCategoryComponents(): Record<string, React.ComponentType<CategoryUiComponentProps>> {
    return { ...this.categoryComponents };
  }

  /**
   * Возвращает все зарегистрированные компоненты для атрибутов
   */
  static getAllAttributeComponents(): Record<string, React.ComponentType<AttributeComponentProps>> {
    return { ...this.attributeComponents };
  }

  /**
   * Проверяет, существует ли компонент для категории с указанным именем
   * @param name Имя компонента
   */
  static hasCategoryComponent(name: string): boolean {
    return !!this.categoryComponents[name];
  }

  /**
   * Проверяет, существует ли компонент для атрибута с указанным именем
   * @param name Имя компонента
   */
  static hasAttributeComponent(name: string): boolean {
    return !!this.attributeComponents[name];
  }
}

export default ComponentRegistry;