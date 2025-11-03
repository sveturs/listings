'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useCreateProduct } from '@/contexts/CreateProductContext';
import { apiClient } from '@/services/api-client';
import { toast } from '@/utils/toast';
import { getTranslatedAttribute } from '@/utils/translatedAttribute';
import MultiSelectAttribute from '@/components/attributes/MultiSelectAttribute';
import RangeAttribute from '@/components/attributes/RangeAttribute';
import type { components } from '@/types/generated/api';

type CategoryAttribute = components['schemas']['models.CategoryAttribute'];

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
  const _tCommon = useTranslations('common');
  const t = useTranslations('storefronts');
  const locale = useLocale();
  const { state, setAttribute, setError, clearError } = useCreateProduct();
  const [attributes, setAttributes] = useState<CategoryAttribute[]>([]);
  const [loading, setLoading] = useState(true);
  const [formData, setFormData] = useState<Record<number, any>>(
    state.attributes || {}
  );
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(
    new Set(['basic', 'technical'])
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

  // Автоматически разворачиваем группы с обязательными полями
  useEffect(() => {
    if (attributes.length > 0) {
      const attributeGroups = groupAttributes();
      const groupsWithRequired = attributeGroups
        .filter((group) => group.attributes.some((attr) => attr.is_required))
        .map((group) => group.id);

      // Объединяем базовые группы с группами, содержащими обязательные поля
      const autoExpandGroups = new Set([
        'basic',
        'technical',
        ...groupsWithRequired,
      ]);

      setExpandedGroups(autoExpandGroups);
    }
  }, [attributes]); // eslint-disable-line react-hooks/exhaustive-deps

  // Группируем атрибуты по логическим группам с учетом группы атрибутов
  const groupAttributes = (): AttributeGroup[] => {
    const groupsMap = new Map<string, AttributeGroup>();

    // Предопределенные группы с иконками
    const predefinedGroups: Record<
      string,
      { name: string; icon: string; priority: number }
    > = {
      basic: {
        name: t('products.attributeGroups.basic'),
        icon: '🏷️',
        priority: 1,
      },
      technical: {
        name: t('products.attributeGroups.technical'),
        icon: '⚙️',
        priority: 2,
      },
      condition: {
        name: t('products.attributeGroups.condition'),
        icon: '✨',
        priority: 3,
      },
      accessories: {
        name: t('products.attributeGroups.accessories'),
        icon: '📦',
        priority: 4,
      },
      dimensions: {
        name: t('products.attributeGroups.dimensions'),
        icon: '📏',
        priority: 5,
      },
      other: {
        name: t('products.attributeGroups.other'),
        icon: '📋',
        priority: 99,
      },
    };

    // Создаем группы на основе attribute_group_id или логики названий
    attributes.forEach((attr) => {
      let groupId = 'other';
      const name = attr.name?.toLowerCase() || '';

      // Группировка по названию атрибута
      if (
        ['brand', 'model', 'manufacturer', 'year'].some((key) =>
          name.includes(key)
        )
      ) {
        groupId = 'basic';
      } else if (
        [
          'storage',
          'memory',
          'screen',
          'resolution',
          'processor',
          'ram',
          'battery',
        ].some((key) => name.includes(key))
      ) {
        groupId = 'technical';
      } else if (
        ['condition', 'warranty', 'used', 'new'].some((key) =>
          name.includes(key)
        )
      ) {
        groupId = 'condition';
      } else if (
        ['accessories', 'included', 'box', 'charger', 'cable'].some((key) =>
          name.includes(key)
        )
      ) {
        groupId = 'accessories';
      } else if (
        ['width', 'height', 'length', 'weight', 'size'].some((key) =>
          name.includes(key)
        )
      ) {
        groupId = 'dimensions';
      }

      if (!groupsMap.has(groupId)) {
        const groupInfo = predefinedGroups[groupId] || predefinedGroups.other;
        groupsMap.set(groupId, {
          id: groupId,
          name: groupInfo.name,
          icon: groupInfo.icon,
          attributes: [],
        });
      }

      groupsMap.get(groupId)!.attributes.push(attr);
    });

    // Сортируем группы по приоритету и атрибуты внутри групп по sort_order
    const groups = Array.from(groupsMap.values()).sort((a, b) => {
      const priorityA = predefinedGroups[a.id]?.priority || 99;
      const priorityB = predefinedGroups[b.id]?.priority || 99;
      return priorityA - priorityB;
    });

    groups.forEach((group) => {
      group.attributes.sort(
        (a, b) => (a.sort_order || 0) - (b.sort_order || 0)
      );
    });

    return groups;
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

    // Получаем переведенные значения
    const { displayName, getOptionLabel } = attribute.id
      ? getTranslatedAttribute(attribute as any, locale)
      : {
          displayName: attribute.display_name || attribute.name || '',
          getOptionLabel: (v: string) => v,
        };

    return (
      <div key={attribute.id} className="form-control">
        <label className="label">
          <span className="label-text font-medium">
            {displayName}
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
          <div className="flex items-center gap-2">
            <input
              type="number"
              className={`input input-bordered flex-1 ${hasError ? 'input-error' : ''}`}
              placeholder="0"
              value={value}
              onChange={(e) =>
                handleAttributeChange(
                  attribute.id!,
                  parseFloat(e.target.value) || 0
                )
              }
              min="0"
              step={attrOptions?.step || 1}
            />
            {attrOptions?.unit && (
              <span className="text-sm text-base-content/60 min-w-fit">
                {attrOptions.unit}
              </span>
            )}
          </div>
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
              {t('select')} {displayName.toLowerCase()}
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
                  {getOptionLabel(option)}
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
              {value ? t('yes') : t('no')}
            </span>
          </div>
        )}

        {/* Multiselect */}
        {attribute.attribute_type === 'multiselect' && (
          <MultiSelectAttribute
            attribute={attribute as any}
            value={value}
            onChange={(values) => handleAttributeChange(attribute.id!, values)}
            error={
              hasError ? state.errors[`attribute_${attribute.id}`] : undefined
            }
            locale={locale}
          />
        )}

        {/* Range */}
        {attribute.attribute_type === 'range' && (
          <RangeAttribute
            attribute={attribute as any}
            value={value}
            onChange={(rangeValue) =>
              handleAttributeChange(attribute.id!, rangeValue)
            }
            error={
              hasError ? state.errors[`attribute_${attribute.id}`] : undefined
            }
            locale={locale}
          />
        )}

        {/* Error message for other attribute types */}
        {hasError &&
          attribute.attribute_type !== 'multiselect' &&
          attribute.attribute_type !== 'range' && (
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
          {t('products.categoryAttributes')}
        </h2>
        <p className="text-lg text-base-content/70">
          {t('products.categoryAttributesDescription')}
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
            {t('products.noAttributesTitle')}
          </h3>
          <p className="text-lg text-base-content/70">
            {t('products.noAttributesMessage')}
          </p>
        </div>
      ) : (
        <div className="space-y-6 mb-8">
          {/* Сводка обязательных полей */}
          {attributes.some((attr) => attr.is_required) && (
            <div className="alert alert-info">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                className="stroke-current shrink-0 w-6 h-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path>
              </svg>
              <span>{t('products.requiredFieldsInfo')}</span>
            </div>
          )}

          {/* Группы атрибутов */}
          <div className="grid grid-cols-1 gap-4">
            {attributeGroups.map((group) => {
              const isExpanded = expandedGroups.has(group.id);
              const hasRequiredFields = group.attributes.some(
                (attr) => attr.is_required
              );
              const filledRequiredFields = group.attributes
                .filter((attr) => attr.is_required)
                .every(
                  (attr) => formData[attr.id!] && formData[attr.id!] !== ''
                );

              return (
                <div key={group.id} className="card bg-base-100 shadow-lg">
                  <div
                    className="card-body cursor-pointer select-none"
                    onClick={() => {
                      const newExpanded = new Set(expandedGroups);
                      if (isExpanded) {
                        newExpanded.delete(group.id);
                      } else {
                        newExpanded.add(group.id);
                      }
                      setExpandedGroups(newExpanded);
                    }}
                  >
                    <div className="flex items-center justify-between">
                      <h3 className="card-title text-xl flex items-center gap-3">
                        <span className="text-2xl">{group.icon}</span>
                        {group.name}
                        <div className="badge badge-neutral">
                          {group.attributes.length}
                        </div>
                        {hasRequiredFields && (
                          <div
                            className={`badge ${filledRequiredFields ? 'badge-success' : 'badge-warning'}`}
                          >
                            {filledRequiredFields ? '✓' : t('required')}
                          </div>
                        )}
                      </h3>
                      <svg
                        className={`w-6 h-6 transition-transform ${isExpanded ? 'rotate-180' : ''}`}
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M19 9l-7 7-7-7"
                        />
                      </svg>
                    </div>
                  </div>

                  {isExpanded && (
                    <div className="px-8 pb-8">
                      <div className="space-y-4">
                        {group.attributes.map(renderAttribute)}
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Прогресс заполнения */}
      {attributes.length > 0 && (
        <div className="card bg-base-200 mb-6">
          <div className="card-body p-4">
            <div className="flex items-center justify-between mb-2">
              <h4 className="text-sm font-medium">{t('progress')}</h4>
              <span className="text-sm text-base-content/70">
                {
                  Object.keys(formData).filter(
                    (key) =>
                      formData[parseInt(key)] && formData[parseInt(key)] !== ''
                  ).length
                }{' '}
                / {attributes.length}
              </span>
            </div>
            <progress
              className="progress progress-primary"
              value={
                Object.keys(formData).filter(
                  (key) =>
                    formData[parseInt(key)] && formData[parseInt(key)] !== ''
                ).length
              }
              max={attributes.length}
            ></progress>
          </div>
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
          {t('back')}
        </button>

        <button onClick={handleNext} className="btn btn-primary btn-lg px-8">
          {t('next')}
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
