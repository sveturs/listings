'use client';

import { useMemo, memo } from 'react';
import { UnifiedAttributeField } from './UnifiedAttributeField';
import { ColorPickerAttribute } from './attribute-inputs/ColorPickerAttribute';
import { SizePickerAttribute } from './attribute-inputs/SizePickerAttribute';
import type { components } from '@/types/generated/api';

type UnifiedAttribute =
  components['schemas']['backend_internal_domain_models.UnifiedAttribute'];
type UnifiedAttributeValue =
  components['schemas']['backend_internal_domain_models.UnifiedAttributeValue'];

interface IntuitivePart {
  component: 'color' | 'size' | 'default';
  props?: {
    sizeType?: 'clothing' | 'shoes' | 'generic';
  };
}

interface IntuitiveMappingRule {
  name_patterns: RegExp[];
  attribute_types: string[];
  component: IntuitivePart['component'];
  props?: IntuitivePart['props'];
  priority: number;
}

// Правила для автоматического определения интуитивных компонентов
const INTUITIVE_MAPPING_RULES: IntuitiveMappingRule[] = [
  // Цвета - высокий приоритет
  {
    name_patterns: [/color|colour|цвет|боја/i],
    attribute_types: ['text', 'select'],
    component: 'color',
    priority: 10,
  },

  // Размеры одежды
  {
    name_patterns: [/size|размер|величина/i, /clothing|одежд|облач/i],
    attribute_types: ['text', 'select'],
    component: 'size',
    props: { sizeType: 'clothing' },
    priority: 9,
  },

  // Размеры обуви
  {
    name_patterns: [/shoe.*size|размер.*обув|величина.*ципел/i],
    attribute_types: ['text', 'select'],
    component: 'size',
    props: { sizeType: 'shoes' },
    priority: 10,
  },

  // Общие размеры
  {
    name_patterns: [/size|размер|величина/i],
    attribute_types: ['text', 'select'],
    component: 'size',
    props: { sizeType: 'generic' },
    priority: 7,
  },
];

interface IntuitiverAttributeFieldProps {
  attribute: UnifiedAttribute;
  value?: UnifiedAttributeValue;
  onChange: (value: UnifiedAttributeValue) => void;
  error?: string;
  disabled?: boolean;
  required?: boolean;
  className?: string;
  enableAutocomplete?: boolean;
}

function IntuitiverAttributeFieldComponent(
  props: IntuitiverAttributeFieldProps
) {
  const { attribute } = props;

  // Определение подходящего интуитивного компонента
  const intuitiveComponent = useMemo((): IntuitivePart => {
    if (!attribute.name || !attribute.attribute_type) {
      return { component: 'default' };
    }

    const attributeName = attribute.name.toLowerCase();
    const attributeType = attribute.attribute_type;

    // Поиск подходящего правила
    const matchedRule = INTUITIVE_MAPPING_RULES.filter((rule) => {
      // Проверяем тип атрибута
      const typeMatches = rule.attribute_types.includes(attributeType);
      if (!typeMatches) return false;

      // Проверяем паттерны названий
      const nameMatches = rule.name_patterns.some(
        (pattern) =>
          pattern.test(attributeName) ||
          pattern.test(attribute.display_name?.toLowerCase() || '') ||
          pattern.test(attribute.description?.toLowerCase() || '')
      );

      return nameMatches;
    }).sort((a, b) => b.priority - a.priority)[0]; // Выбираем с наивысшим приоритетом

    if (matchedRule) {
      return {
        component: matchedRule.component,
        props: matchedRule.props,
      };
    }

    return { component: 'default' };
  }, [attribute]);

  // Рендеринг соответствующего компонента
  const renderIntuitiveComponent = () => {
    if (!attribute.id) {
      return <UnifiedAttributeField {...props} />;
    }

    switch (intuitiveComponent.component) {
      case 'color':
        return (
          <div className="space-y-2">
            {/* Заголовок атрибута */}
            <label className="label">
              <span className="label-text font-medium">
                {attribute.display_name || attribute.name}
                {attribute.is_required && (
                  <span className="text-error"> *</span>
                )}
              </span>
            </label>

            <ColorPickerAttribute
              attributeId={attribute.id}
              value={props.value}
              onChange={props.onChange}
              className={props.className}
            />

            {/* Описание атрибута */}
            {attribute.description && (
              <label className="label">
                <span className="label-text-alt opacity-70">
                  {attribute.description}
                </span>
              </label>
            )}

            {/* Ошибка */}
            {props.error && (
              <label className="label">
                <span className="label-text-alt text-error">{props.error}</span>
              </label>
            )}
          </div>
        );

      case 'size':
        return (
          <div className="space-y-2">
            {/* Заголовок атрибута */}
            <label className="label">
              <span className="label-text font-medium">
                {attribute.display_name || attribute.name}
                {attribute.is_required && (
                  <span className="text-error"> *</span>
                )}
              </span>
            </label>

            <SizePickerAttribute
              attributeId={attribute.id}
              value={props.value}
              onChange={props.onChange}
              className={props.className}
              sizeType={intuitiveComponent.props?.sizeType}
            />

            {/* Описание атрибута */}
            {attribute.description && (
              <label className="label">
                <span className="label-text-alt opacity-70">
                  {attribute.description}
                </span>
              </label>
            )}

            {/* Ошибка */}
            {props.error && (
              <label className="label">
                <span className="label-text-alt text-error">{props.error}</span>
              </label>
            )}
          </div>
        );

      default:
        // Используем стандартный компонент с автокомплитом для текстовых полей
        return (
          <UnifiedAttributeField
            {...props}
            enableAutocomplete={
              props.enableAutocomplete || attribute.attribute_type === 'text'
            }
          />
        );
    }
  };

  return renderIntuitiveComponent();
}

// Мемоизированный экспорт для оптимизации производительности
export const IntuitiverAttributeField = memo(
  IntuitiverAttributeFieldComponent,
  (prevProps, nextProps) => {
    // Глубокое сравнение для атрибута и значения
    const attributeEquals =
      prevProps.attribute.id === nextProps.attribute.id &&
      prevProps.attribute.name === nextProps.attribute.name &&
      prevProps.attribute.attribute_type ===
        nextProps.attribute.attribute_type &&
      prevProps.attribute.display_name === nextProps.attribute.display_name;

    const valueEquals =
      JSON.stringify(prevProps.value) === JSON.stringify(nextProps.value);

    return (
      attributeEquals &&
      valueEquals &&
      prevProps.disabled === nextProps.disabled &&
      prevProps.required === nextProps.required &&
      prevProps.className === nextProps.className &&
      prevProps.enableAutocomplete === nextProps.enableAutocomplete &&
      prevProps.error === nextProps.error &&
      prevProps.onChange === nextProps.onChange
    );
  }
);

// Вспомогательная функция для добавления новых правил (для администраторов)
export function addIntuitiveRule(rule: Omit<IntuitiveMappingRule, 'priority'>) {
  INTUITIVE_MAPPING_RULES.push({
    ...rule,
    priority: 5, // Средний приоритет для пользовательских правил
  });
}

// Функция для получения статистики использования интуитивных компонентов (для аналитики)
export function getIntuitiveComponentUsage(attributes: UnifiedAttribute[]) {
  const usage = {
    color: 0,
    size: 0,
    default: 0,
    total: attributes.length,
  };

  attributes.forEach((attr) => {
    if (!attr.name || !attr.attribute_type) {
      usage.default++;
      return;
    }

    const attributeName = attr.name.toLowerCase();
    const attributeType = attr.attribute_type;

    const matchedRule = INTUITIVE_MAPPING_RULES.filter((rule) => {
      const typeMatches = rule.attribute_types.includes(attributeType);
      if (!typeMatches) return false;

      const nameMatches = rule.name_patterns.some(
        (pattern) =>
          pattern.test(attributeName) ||
          pattern.test(attr.display_name?.toLowerCase() || '') ||
          pattern.test(attr.description?.toLowerCase() || '')
      );

      return nameMatches;
    }).sort((a, b) => b.priority - a.priority)[0];

    if (matchedRule) {
      usage[matchedRule.component]++;
    } else {
      usage.default++;
    }
  });

  return usage;
}
