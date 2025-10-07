'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { AutocompleteAttributeField } from './AutocompleteAttributeField';
import type { components } from '@/types/generated/api';

type UnifiedAttribute =
  components['schemas']['models.UnifiedAttribute'];
type UnifiedAttributeValue =
  components['schemas']['models.UnifiedAttributeValue'];

interface UnifiedAttributeFieldProps {
  attribute: UnifiedAttribute;
  value?: UnifiedAttributeValue;
  onChange: (value: UnifiedAttributeValue) => void;
  error?: string;
  disabled?: boolean;
  required?: boolean;
  className?: string;
  enableAutocomplete?: boolean; // Новый проп для включения автокомплита
}

// Extended type for attribute with all needed properties
interface ExtendedAttribute extends UnifiedAttribute {
  placeholder?: string;
  description?: string;
  help_text?: string;
  unit?: string;
  translations?: Record<string, string>;
  validation?: {
    min?: number;
    max?: number;
    step?: number | string;
    minLength?: number;
    maxLength?: number;
    pattern?: string;
    rows?: number;
  };
}

export function UnifiedAttributeField({
  attribute: rawAttribute,
  value,
  onChange,
  error,
  disabled = false,
  required = false,
  className = '',
  enableAutocomplete = false,
}: UnifiedAttributeFieldProps) {
  // Cast to extended type for easier access
  const attribute = rawAttribute as ExtendedAttribute;
  const t = useTranslations('common');
  const locale = useLocale();

  // Локальное состояние для значений
  const [localValue, setLocalValue] = useState<UnifiedAttributeValue>({
    attribute_id: attribute.id,
    text_value: value?.text_value || '',
    numeric_value: value?.numeric_value,
    boolean_value: value?.boolean_value,
    date_value: value?.date_value,
    json_value: value?.json_value,
  });

  useEffect(() => {
    if (value) {
      setLocalValue(value);
    }
  }, [value]);

  const handleChange = (newValue: Partial<UnifiedAttributeValue>) => {
    const updatedValue = {
      ...localValue,
      ...newValue,
      attribute_id: attribute.id,
    };
    setLocalValue(updatedValue);
    onChange(updatedValue);
  };

  // Получаем локализованное имя атрибута
  const getLocalizedName = () => {
    if (attribute.translations && attribute.translations[locale]) {
      return attribute.translations[locale];
    }
    return attribute.display_name || attribute.name || '';
  };

  // Получаем локализованные опции для select
  const getLocalizedOptions = () => {
    if (!attribute.options || attribute.options.length === 0) {
      return [];
    }

    return attribute.options.map((option) => {
      if (typeof option === 'string' || typeof option === 'number') {
        return { value: String(option), label: String(option) };
      }
      if (typeof option === 'object' && option !== null) {
        const optionObj = option as any;
        if (optionObj.translations && optionObj.translations[locale]) {
          return {
            value: String(optionObj.value),
            label: optionObj.translations[locale],
          };
        }
        return {
          value: String(optionObj.value || option),
          label: String(optionObj.label || optionObj.value || option),
        };
      }
      return { value: String(option), label: String(option) };
    });
  };

  const renderField = () => {
    switch (attribute.attribute_type) {
      case 'text':
        // Используем автокомплит если он включен и атрибут не отключен
        if (enableAutocomplete && !disabled && attribute.id) {
          return (
            <AutocompleteAttributeField
              attribute={attribute}
              value={localValue}
              onChange={onChange}
              className={error ? 'has-error' : ''}
            />
          );
        }

        return (
          <input
            type="text"
            value={localValue.text_value || ''}
            onChange={(e) => handleChange({ text_value: e.target.value })}
            disabled={disabled}
            className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
            placeholder={attribute.placeholder || ''}
          />
        );

      case 'number':
        return (
          <div className="flex items-center gap-2">
            <input
              type="number"
              value={localValue.numeric_value || ''}
              onChange={(e) =>
                handleChange({
                  numeric_value: e.target.value
                    ? parseFloat(e.target.value)
                    : undefined,
                })
              }
              disabled={disabled}
              min={attribute.validation?.min}
              max={attribute.validation?.max}
              step={attribute.validation?.step || 'any'}
              className={`input input-bordered flex-1 ${error ? 'input-error' : ''}`}
              placeholder={attribute.placeholder || ''}
            />
            {attribute.unit && (
              <span className="text-sm text-base-content/70">
                {attribute.unit}
              </span>
            )}
          </div>
        );

      case 'select':
        const options = getLocalizedOptions();
        return (
          <select
            value={localValue.text_value || ''}
            onChange={(e) => handleChange({ text_value: e.target.value })}
            disabled={disabled}
            className={`select select-bordered w-full ${error ? 'select-error' : ''}`}
          >
            <option value="">{t('select_option')}</option>
            {options.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        );

      case 'multiselect':
        const multiOptions = getLocalizedOptions();
        let selectedValues: string[] = [];

        // Parse selected values from text_value (JSON string) or json_value (number array)
        if (localValue.text_value) {
          try {
            const parsed = JSON.parse(localValue.text_value);
            if (Array.isArray(parsed)) {
              selectedValues = parsed.map((v) => String(v));
            }
          } catch {
            // If not valid JSON, treat as comma-separated string
            selectedValues = localValue.text_value
              .split(',')
              .map((v) => v.trim());
          }
        } else if (Array.isArray(localValue.json_value)) {
          selectedValues = localValue.json_value.map((v) => String(v));
        }

        return (
          <div className="space-y-2">
            {multiOptions.map((option) => (
              <label
                key={option.value}
                className="flex items-center gap-2 cursor-pointer"
              >
                <input
                  type="checkbox"
                  checked={selectedValues.includes(option.value)}
                  onChange={(e) => {
                    const newValues = e.target.checked
                      ? [...selectedValues, option.value]
                      : selectedValues.filter(
                          (v: string) => v !== option.value
                        );
                    // Store as JSON string in text_value
                    handleChange({ text_value: JSON.stringify(newValues) });
                  }}
                  disabled={disabled}
                  className="checkbox checkbox-primary"
                />
                <span className="text-sm">{option.label}</span>
              </label>
            ))}
          </div>
        );

      case 'boolean':
        return (
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={localValue.boolean_value || false}
              onChange={(e) =>
                handleChange({ boolean_value: e.target.checked })
              }
              disabled={disabled}
              className="checkbox checkbox-primary"
            />
            <span>{attribute.placeholder || t('yes')}</span>
          </label>
        );

      case 'date':
        return (
          <input
            type="date"
            value={localValue.date_value || ''}
            onChange={(e) => handleChange({ date_value: e.target.value })}
            disabled={disabled}
            className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
          />
        );

      case 'datetime':
        return (
          <input
            type="datetime-local"
            value={localValue.date_value || ''}
            onChange={(e) => handleChange({ date_value: e.target.value })}
            disabled={disabled}
            className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
          />
        );

      case 'textarea':
        return (
          <textarea
            value={localValue.text_value || ''}
            onChange={(e) => handleChange({ text_value: e.target.value })}
            disabled={disabled}
            rows={attribute.validation?.rows || 4}
            className={`textarea textarea-bordered w-full ${error ? 'textarea-error' : ''}`}
            placeholder={attribute.placeholder || ''}
          />
        );

      case 'radio':
        const radioOptions = getLocalizedOptions();
        return (
          <div className="space-y-2">
            {radioOptions.map((option) => (
              <label
                key={option.value}
                className="flex items-center gap-2 cursor-pointer"
              >
                <input
                  type="radio"
                  name={`attribute-${attribute.id}`}
                  value={option.value}
                  checked={localValue.text_value === option.value}
                  onChange={(e) => handleChange({ text_value: e.target.value })}
                  disabled={disabled}
                  className="radio radio-primary"
                />
                <span className="text-sm">{option.label}</span>
              </label>
            ))}
          </div>
        );

      case 'range':
        return (
          <div className="space-y-2">
            <input
              type="range"
              value={localValue.numeric_value || attribute.validation?.min || 0}
              onChange={(e) =>
                handleChange({
                  numeric_value: parseFloat(e.target.value),
                })
              }
              disabled={disabled}
              min={attribute.validation?.min}
              max={attribute.validation?.max}
              step={attribute.validation?.step}
              className="range range-primary"
            />
            <div className="flex justify-between text-xs text-base-content/70">
              <span>{attribute.validation?.min || 0}</span>
              <span className="font-semibold">
                {localValue.numeric_value || 0}
              </span>
              <span>{attribute.validation?.max || 100}</span>
            </div>
          </div>
        );

      case 'email':
        return (
          <input
            type="email"
            value={localValue.text_value || ''}
            onChange={(e) => handleChange({ text_value: e.target.value })}
            disabled={disabled}
            className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
            placeholder={attribute.placeholder || 'email@example.com'}
          />
        );

      case 'tel':
        return (
          <input
            type="tel"
            value={localValue.text_value || ''}
            onChange={(e) => handleChange({ text_value: e.target.value })}
            disabled={disabled}
            className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
            placeholder={attribute.placeholder || '+1234567890'}
          />
        );

      case 'url':
        return (
          <input
            type="url"
            value={localValue.text_value || ''}
            onChange={(e) => handleChange({ text_value: e.target.value })}
            disabled={disabled}
            className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
            placeholder={attribute.placeholder || 'https://example.com'}
          />
        );

      case 'color':
        return (
          <div className="flex items-center gap-2">
            <input
              type="color"
              value={localValue.text_value || '#000000'}
              onChange={(e) => handleChange({ text_value: e.target.value })}
              disabled={disabled}
              className="w-16 h-10 border border-base-300 rounded cursor-pointer"
            />
            <input
              type="text"
              value={localValue.text_value || ''}
              onChange={(e) => handleChange({ text_value: e.target.value })}
              disabled={disabled}
              className={`input input-bordered flex-1 ${error ? 'input-error' : ''}`}
              placeholder="#000000"
            />
          </div>
        );

      default:
        // Fallback для неизвестных типов
        return (
          <input
            type="text"
            value={localValue.text_value || ''}
            onChange={(e) => handleChange({ text_value: e.target.value })}
            disabled={disabled}
            className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
            placeholder={attribute.placeholder || ''}
          />
        );
    }
  };

  return (
    <div className={`form-control ${className}`}>
      <label className="label">
        <span className="label-text">
          {getLocalizedName()}
          {required && <span className="text-error ml-1">*</span>}
        </span>
        {attribute.description && (
          <span className="label-text-alt text-base-content/60">
            {attribute.description}
          </span>
        )}
      </label>

      {renderField()}

      {error && (
        <label className="label">
          <span className="label-text-alt text-error">{error}</span>
        </label>
      )}

      {attribute.help_text && !error && (
        <label className="label">
          <span className="label-text-alt text-base-content/60">
            {attribute.help_text}
          </span>
        </label>
      )}
    </div>
  );
}
