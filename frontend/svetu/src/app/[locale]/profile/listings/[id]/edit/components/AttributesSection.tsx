'use client';

import { useTranslations } from 'next-intl';
import { useEffect, useState, useCallback } from 'react';
import { apiClient } from '@/services/api-client';

interface CategoryAttribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type: 'text' | 'numeric' | 'select' | 'multiselect' | 'boolean';
  options?:
    | {
        values?: string[];
      }
    | Array<{ value: string; label: string }>;
  validation_rules?: {
    min?: number;
    max?: number;
    pattern?: string;
    required?: boolean;
  };
  is_required: boolean;
  unit?: string;
  translations?: {
    [key: string]: string;
  };
}

interface AttributeValue {
  attribute_id: number;
  value: string | number | boolean | string[];
}

interface AttributesSectionProps {
  categoryId: number;
  values: AttributeValue[];
  errors?: Record<string, string>;
  onChange: (values: AttributeValue[]) => void;
}

export function AttributesSection({
  categoryId,
  values,
  errors = {},
  onChange,
}: AttributesSectionProps) {
  const t = useTranslations('profile');
  const [attributes, setAttributes] = useState<CategoryAttribute[]>([]);
  const [loading, setLoading] = useState(true);

  const loadAttributes = useCallback(async () => {
    try {
      const response = await apiClient.get(
        `/api/v1/c2c/categories/${categoryId}/attributes`
      );
      if (response.data) {
        // Проверяем структуру ответа
        let attributesData = [];
        if (response.data.data && Array.isArray(response.data.data)) {
          attributesData = response.data.data;
        } else if (Array.isArray(response.data)) {
          attributesData = response.data;
        }
        console.log('Loaded attributes:', attributesData);
        setAttributes(attributesData);
      }
    } catch (error) {
      console.error('Error loading attributes:', error);
      setAttributes([]);
    } finally {
      setLoading(false);
    }
  }, [categoryId]);

  useEffect(() => {
    loadAttributes();
  }, [categoryId, loadAttributes]);

  const getValue = (attributeId: number) => {
    if (!Array.isArray(values)) return '';
    const value = values.find((v) => v.attribute_id === attributeId);
    return value?.value ?? '';
  };

  const updateValue = (attributeId: number, value: any) => {
    const currentValues = Array.isArray(values) ? values : [];
    const newValues = currentValues.filter(
      (v) => v.attribute_id !== attributeId
    );
    if (value !== '' && value !== null && value !== undefined) {
      newValues.push({ attribute_id: attributeId, value });
    }
    onChange(newValues);
  };

  const renderField = (attr: CategoryAttribute) => {
    const value = getValue(attr.id);
    const error = errors[`attribute_${attr.id}`];

    switch (attr.attribute_type) {
      case 'text':
        return (
          <div key={attr.id} className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {attr.display_name}
                {attr.is_required && ' *'}
              </span>
            </label>
            <input
              type="text"
              className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
              value={value as string}
              onChange={(e) => updateValue(attr.id, e.target.value)}
              required={attr.is_required}
            />
            {error && (
              <label className="label">
                <span className="label-text-alt text-error">{error}</span>
              </label>
            )}
          </div>
        );

      case 'numeric':
        return (
          <div key={attr.id} className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {attr.display_name}
                {attr.is_required && ' *'}
              </span>
            </label>
            <div className="join w-full">
              <input
                type="number"
                className={`input input-bordered join-item flex-1 ${error ? 'input-error' : ''}`}
                value={value as number}
                onChange={(e) => updateValue(attr.id, Number(e.target.value))}
                required={attr.is_required}
                min={attr.validation_rules?.min}
                max={attr.validation_rules?.max}
              />
              {attr.unit && (
                <span className="join-item btn btn-disabled">{attr.unit}</span>
              )}
            </div>
            {error && (
              <label className="label">
                <span className="label-text-alt text-error">{error}</span>
              </label>
            )}
          </div>
        );

      case 'select':
        const selectOptions = attr.options;
        let optionsList: Array<{ value: string; label: string }> = [];

        // Handle different option formats
        if (
          selectOptions &&
          'values' in selectOptions &&
          Array.isArray(selectOptions.values)
        ) {
          // Format: { values: ["option1", "option2"] }
          optionsList = selectOptions.values.map((val) => ({
            value: val,
            label:
              val.charAt(0).toUpperCase() + val.slice(1).replace(/_/g, ' '),
          }));
        } else if (Array.isArray(selectOptions)) {
          // Format: [{ value: "...", label: "..." }]
          optionsList = selectOptions;
        }

        return (
          <div key={attr.id} className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {attr.display_name}
                {attr.is_required && ' *'}
              </span>
            </label>
            <select
              className={`select select-bordered w-full ${error ? 'select-error' : ''}`}
              value={value as string}
              onChange={(e) => updateValue(attr.id, e.target.value)}
              required={attr.is_required}
            >
              <option value="">{t('attributes.selectOption')}</option>
              {optionsList.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
            {error && (
              <label className="label">
                <span className="label-text-alt text-error">{error}</span>
              </label>
            )}
          </div>
        );

      case 'multiselect':
        const multiselectOptions = attr.options;
        let multiOptionsList: Array<{ value: string; label: string }> = [];

        // Handle different option formats
        if (
          multiselectOptions &&
          'values' in multiselectOptions &&
          Array.isArray(multiselectOptions.values)
        ) {
          // Format: { values: ["option1", "option2"] }
          multiOptionsList = multiselectOptions.values.map((val) => ({
            value: val,
            label:
              val.charAt(0).toUpperCase() + val.slice(1).replace(/_/g, ' '),
          }));
        } else if (Array.isArray(multiselectOptions)) {
          // Format: [{ value: "...", label: "..." }]
          multiOptionsList = multiselectOptions;
        }

        return (
          <div key={attr.id} className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {attr.display_name}
                {attr.is_required && ' *'}
              </span>
            </label>
            <div className="border border-base-300 rounded-lg p-3 space-y-2 max-h-48 overflow-y-auto">
              {multiOptionsList.map((opt) => (
                <label
                  key={opt.value}
                  className="cursor-pointer flex items-center gap-2"
                >
                  <input
                    type="checkbox"
                    className="checkbox checkbox-sm"
                    checked={(value as string[])?.includes(opt.value) || false}
                    onChange={(e) => {
                      const currentValues = (value as string[]) || [];
                      if (e.target.checked) {
                        updateValue(attr.id, [...currentValues, opt.value]);
                      } else {
                        updateValue(
                          attr.id,
                          currentValues.filter((v) => v !== opt.value)
                        );
                      }
                    }}
                  />
                  <span className="text-sm">{opt.label}</span>
                </label>
              ))}
            </div>
            {error && (
              <label className="label">
                <span className="label-text-alt text-error">{error}</span>
              </label>
            )}
          </div>
        );

      case 'boolean':
        return (
          <div key={attr.id} className="form-control">
            <label className="label cursor-pointer">
              <span className="label-text font-medium">
                {attr.display_name}
                {attr.is_required && ' *'}
              </span>
              <input
                type="checkbox"
                className="checkbox"
                checked={(value as boolean) || false}
                onChange={(e) => updateValue(attr.id, e.target.checked)}
              />
            </label>
            {error && (
              <label className="label">
                <span className="label-text-alt text-error">{error}</span>
              </label>
            )}
          </div>
        );

      default:
        return null;
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (!attributes || attributes.length === 0) {
    return (
      <div className="text-center py-8 text-base-content/70">
        {t('attributes.noAttributes')}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="w-5 h-5 text-primary"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3z"
            />
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M6 6h.008v.008H6V6z"
            />
          </svg>
          {t('attributes.title')}
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {attributes.map((attr) => renderField(attr))}
        </div>
      </div>
    </div>
  );
}
