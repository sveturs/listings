'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateProduct } from '@/contexts/CreateProductContext';
import { apiClient } from '@/services/api-client';
import { toast } from '@/utils/toast';
import type { components } from '@/types/generated/api';

type CategoryAttribute =
  components['schemas']['backend_internal_domain_models.CategoryAttribute'];

interface AttributeGroup {
  id: string;
  name: string;
  icon: string;
  attributes: CategoryAttribute[];
}

interface AttributesStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function AttributesStep({
  onNext,
  onBack,
}: AttributesStepProps) {
  const t = useTranslations();
  const { state, setAttribute, setError, clearError } = useCreateProduct();
  const [attributes, setAttributes] = useState<CategoryAttribute[]>([]);
  const [loading, setLoading] = useState(true);
  const [formData, setFormData] = useState<Record<number, any>>(
    state.attributes || {}
  );

  const loadAttributes = async () => {
    if (!state.category) return;

    try {
      setLoading(true);
      const response = await apiClient.get(
        `/api/v1/marketplace/categories/${state.category.id}/attributes`
      );

      if (response.data) {
        const responseData = response.data.data || response.data;
        if (Array.isArray(responseData)) {
          setAttributes(responseData);
        }
      }
    } catch (error: any) {
      console.error('Failed to load attributes:', error);
      if (error.response?.status !== 404) {
        toast.error('Failed to load attributes');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (state.category) {
      loadAttributes();
    }
  }, [state.category]); // eslint-disable-line react-hooks/exhaustive-deps

  // Группируем атрибуты по логическим группам
  const groupAttributes = (): AttributeGroup[] => {
    const groups: AttributeGroup[] = [
      {
        id: 'basic',
        name: t('storefronts.products.attributeGroups.basic'),
        icon: '🏷️',
        attributes: [],
      },
      {
        id: 'technical',
        name: t('storefronts.products.attributeGroups.technical'),
        icon: '⚙️',
        attributes: [],
      },
      {
        id: 'condition',
        name: t('storefronts.products.attributeGroups.condition'),
        icon: '✨',
        attributes: [],
      },
      {
        id: 'accessories',
        name: t('storefronts.products.attributeGroups.accessories'),
        icon: '📦',
        attributes: [],
      },
      {
        id: 'other',
        name: t('storefronts.products.attributeGroups.other'),
        icon: '📋',
        attributes: [],
      },
    ];

    // Распределяем атрибуты по группам на основе их названий
    attributes.forEach((attr) => {
      const name = attr.name?.toLowerCase() || '';

      if (
        ['brand', 'model', 'manufacturer'].some((key) => name.includes(key))
      ) {
        groups[0].attributes.push(attr); // basic
      } else if (
        [
          'storage',
          'memory',
          'screen_size',
          'resolution',
          'processor',
          'ram',
        ].some((key) => name.includes(key))
      ) {
        groups[1].attributes.push(attr); // technical
      } else if (
        ['condition', 'device_condition', 'warranty', 'used', 'new'].some(
          (key) => name.includes(key)
        )
      ) {
        groups[2].attributes.push(attr); // condition
      } else if (
        ['accessories', 'included', 'box', 'charger', 'cable'].some((key) =>
          name.includes(key)
        )
      ) {
        groups[3].attributes.push(attr); // accessories
      } else {
        groups[4].attributes.push(attr); // other
      }
    });

    // Возвращаем только группы с атрибутами
    return groups.filter((group) => group.attributes.length > 0);
  };

  const handleAttributeChange = (attributeId: number, value: any) => {
    setFormData((prev) => ({ ...prev, [attributeId]: value }));
    setAttribute(attributeId, value);
    clearError(`attribute_${attributeId}`);
  };

  const validateRequiredAttributes = (): boolean => {
    let isValid = true;

    attributes.forEach((attr) => {
      if (
        attr.is_required &&
        (!formData[attr.id!] || formData[attr.id!] === '')
      ) {
        setError(
          `attribute_${attr.id}`,
          `${attr.display_name || attr.name} is required`
        );
        isValid = false;
      }
    });

    return isValid;
  };

  const handleNext = () => {
    if (validateRequiredAttributes()) {
      onNext();
    }
  };

  const renderAttribute = (attribute: CategoryAttribute) => {
    const attrOptions = attribute.options as any;
    const value = formData[attribute.id!] || '';
    const isRequired = attribute.is_required;
    const hasError = state.errors[`attribute_${attribute.id}`];

    return (
      <div key={attribute.id} className="form-control">
        <label className="label">
          <span className="label-text font-medium">
            {attribute.display_name || attribute.name}
            {isRequired && <span className="text-error ml-1">*</span>}
          </span>
        </label>

        {/* Text input */}
        {attribute.attribute_type === 'text' && (
          <input
            type="text"
            className={`input input-bordered ${hasError ? 'input-error' : ''}`}
            placeholder={`Enter ${attribute.display_name || attribute.name}`}
            value={value}
            onChange={(e) =>
              handleAttributeChange(attribute.id!, e.target.value)
            }
          />
        )}

        {/* Number input */}
        {attribute.attribute_type === 'number' && (
          <input
            type="number"
            className={`input input-bordered ${hasError ? 'input-error' : ''}`}
            placeholder={`Enter ${attribute.display_name || attribute.name}`}
            value={value}
            onChange={(e) =>
              handleAttributeChange(
                attribute.id!,
                parseFloat(e.target.value) || 0
              )
            }
          />
        )}

        {/* Select */}
        {attribute.attribute_type === 'select' && (
          <select
            className={`select select-bordered ${hasError ? 'select-error' : ''}`}
            value={value}
            onChange={(e) =>
              handleAttributeChange(attribute.id!, e.target.value)
            }
          >
            <option value="">
              Select {attribute.display_name || attribute.name}
            </option>
            {(() => {
              let options: string[] = [];

              if (Array.isArray(attrOptions)) {
                options = attrOptions.map(String);
              } else if (
                attrOptions?.values &&
                Array.isArray(attrOptions.values)
              ) {
                options = attrOptions.values.map(String);
              }

              return options.map((option) => (
                <option key={`${attribute.id}-${option}`} value={option}>
                  {option}
                </option>
              ));
            })()}
          </select>
        )}

        {/* Boolean toggle */}
        {attribute.attribute_type === 'boolean' && (
          <div className="flex items-center gap-3">
            <input
              type="checkbox"
              className="toggle toggle-primary"
              checked={!!value}
              onChange={(e) =>
                handleAttributeChange(attribute.id!, e.target.checked)
              }
            />
            <span className="text-sm text-base-content/70">
              {value ? 'Yes' : 'No'}
            </span>
          </div>
        )}

        {/* Multiselect checkboxes */}
        {attribute.attribute_type === 'multiselect' && attrOptions?.values && (
          <div className="grid grid-cols-2 gap-2">
            {attrOptions.values.map((option: string) => {
              const currentValues = (value as string[]) || [];
              const isChecked = currentValues.includes(option);

              return (
                <label
                  key={`${attribute.id}-multiselect-${option}`}
                  className="label cursor-pointer justify-start gap-3 bg-base-200 rounded-lg p-3 hover:bg-base-300 transition-colors"
                >
                  <input
                    type="checkbox"
                    className="checkbox checkbox-primary"
                    checked={isChecked}
                    onChange={(e) => {
                      if (e.target.checked) {
                        handleAttributeChange(attribute.id!, [
                          ...currentValues,
                          option,
                        ]);
                      } else {
                        handleAttributeChange(
                          attribute.id!,
                          currentValues.filter((v) => v !== option)
                        );
                      }
                    }}
                  />
                  <span className="label-text">{option}</span>
                </label>
              );
            })}
          </div>
        )}

        {/* Error message */}
        {hasError && (
          <label className="label">
            <span className="label-text-alt text-error">
              {state.errors[`attribute_${attribute.id}`]}
            </span>
          </label>
        )}
      </div>
    );
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  const attributeGroups = groupAttributes();

  return (
    <div className="max-w-6xl mx-auto">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold text-base-content mb-4">
          {t('storefronts.products.categoryAttributes')}
        </h2>
        <p className="text-lg text-base-content/70">
          {t('storefronts.products.categoryAttributesDescription')}
        </p>
        {state.category && (
          <div className="badge badge-primary badge-lg mt-2">
            {state.category.name}
          </div>
        )}
      </div>

      {attributeGroups.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">🎉</div>
          <h3 className="text-2xl font-bold text-base-content mb-2">
            {t('storefronts.products.noAttributesTitle')}
          </h3>
          <p className="text-lg text-base-content/70">
            {t('storefronts.products.noAttributesMessage')}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {attributeGroups.map((group) => (
            <div key={group.id} className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title text-xl mb-4 flex items-center gap-3">
                  <span className="text-2xl">{group.icon}</span>
                  {group.name}
                  <div className="badge badge-neutral">
                    {group.attributes.length}
                  </div>
                </h3>

                <div className="space-y-4">
                  {group.attributes.map(renderAttribute)}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Навигация */}
      <div className="flex justify-between items-center">
        <button onClick={onBack} className="btn btn-outline btn-lg px-8">
          <svg
            className="w-5 h-5 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 19l-7-7 7-7"
            />
          </svg>
          {t('common.back')}
        </button>

        <button onClick={handleNext} className="btn btn-primary btn-lg px-8">
          {t('common.next')}
          <svg
            className="w-5 h-5 ml-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 5l7 7-7 7"
            />
          </svg>
        </button>
      </div>
    </div>
  );
}
