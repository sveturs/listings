'use client';

import React, { useState, useEffect } from 'react';
import { Package, Info, Star, ChevronDown } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';

type CategoryAttribute =
  components['schemas']['backend_internal_domain_models.CategoryAttribute'];
type MarketplaceCategory =
  components['schemas']['backend_internal_domain_models.MarketplaceCategory'];

interface AttributeValue {
  attribute_id: number;
  attribute_name: string;
  display_name: string;
  attribute_type: string;
  text_value?: string;
  numeric_value?: number;
  boolean_value?: boolean;
  unit?: string;
}

interface CategoryAttributesProps {
  selectedCategory: MarketplaceCategory | null;
  attributes: Record<string, AttributeValue>;
  onAttributeChange: (attributeId: number, value: AttributeValue) => void;
  locale: string;
}

export default function CategoryAttributes({
  selectedCategory,
  attributes,
  onAttributeChange,
  locale,
}: CategoryAttributesProps) {
  const t = useTranslations('common');
  const tListing = useTranslations('listing');
  const [categoryAttributes, setCategoryAttributes] = useState<
    CategoryAttribute[]
  >([]);
  const [loading, setLoading] = useState(false);
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(
    new Set(['required'])
  );

  // Загружаем атрибуты при изменении категории
  useEffect(() => {
    const fetchCategoryAttributes = async () => {
      if (!selectedCategory?.id) {
        setCategoryAttributes([]);
        return;
      }

      setLoading(true);
      try {
        const response = await fetch(
          `/api/v1/marketplace/categories/${selectedCategory.id}/attributes?lang=${locale}`
        );

        if (response.ok) {
          const data = await response.json();
          if (data.data && Array.isArray(data.data)) {
            setCategoryAttributes(data.data);
          }
        } else {
          console.error('Failed to fetch category attributes');
          setCategoryAttributes([]);
        }
      } catch (error) {
        console.error('Error fetching category attributes:', error);
        setCategoryAttributes([]);
      } finally {
        setLoading(false);
      }
    };

    fetchCategoryAttributes();
  }, [selectedCategory?.id, locale]);

  const getAttributeName = (attribute: CategoryAttribute) => {
    return (
      attribute.translations?.name ||
      attribute.display_name ||
      attribute.name ||
      ''
    );
  };

  const getAttributeOptions = (
    attribute: CategoryAttribute
  ): Array<{ value: string; label: string }> => {
    // Проверяем новый формат API где options это объект с values
    if (
      attribute.options &&
      typeof attribute.options === 'object' &&
      'values' in attribute.options
    ) {
      const values = (attribute.options as any).values;
      if (Array.isArray(values)) {
        return values.map((value) => {
          // Сначала проверяем переводы атрибута
          let label = value;

          // Проверяем есть ли переводы для конкретного атрибута
          if (
            attribute.option_translations &&
            typeof attribute.option_translations === 'object'
          ) {
            const translation =
              attribute.option_translations?.[locale]?.[value];
            if (translation) {
              label = translation;
            }
          }

          // Если нет специфичного перевода, пробуем использовать системные переводы
          if (label === value) {
            // Пробуем найти перевод в системе через ключи
            // Для состояний товара
            if (['new', 'used', 'refurbished', 'damaged'].includes(value)) {
              const translationKey = `condition.${value}`;
              try {
                const translated = t(translationKey);
                if (translated !== translationKey) {
                  label = translated;
                }
              } catch {
                // Игнорируем ошибку если ключ не найден
              }
            }
            // Для общих значений
            else if (['yes', 'no'].includes(value)) {
              const translationKey = `common.${value}`;
              try {
                const translated = t(translationKey);
                if (translated !== translationKey) {
                  label = translated;
                }
              } catch {
                // Игнорируем ошибку если ключ не найден
              }
            }
          }

          return { value, label };
        });
      }
    }

    // Старый формат где options это массив чисел
    if (Array.isArray(attribute.options) && attribute.options.length > 0) {
      return attribute.options
        .map((optionId) => {
          const optionKey = optionId.toString();
          const translation =
            attribute.option_translations?.[locale]?.[optionKey];
          return {
            value: optionKey,
            label: translation || optionKey,
          };
        })
        .filter(Boolean);
    }

    return [];
  };

  const handleAttributeChange = (attribute: CategoryAttribute, value: any) => {
    const attributeValue: AttributeValue = {
      attribute_id: attribute.id || 0,
      attribute_name: attribute.name || '',
      display_name: getAttributeName(attribute),
      attribute_type: attribute.attribute_type || 'text',
    };

    switch (attribute.attribute_type) {
      case 'numeric':
      case 'number':
        attributeValue.numeric_value = parseFloat(value) || 0;
        break;
      case 'boolean':
        attributeValue.boolean_value = Boolean(value);
        break;
      case 'date':
      case 'text':
      case 'select':
      case 'multiselect':
      default:
        attributeValue.text_value = value;
        break;
    }

    onAttributeChange(attribute.id || 0, attributeValue);
  };

  const renderAttributeInput = (attribute: CategoryAttribute) => {
    const attributeName = getAttributeName(attribute);
    const currentValue = attributes[attribute.id || 0];
    const options = getAttributeOptions(attribute);

    // Если есть кастомный компонент
    if (attribute.custom_component) {
      return (
        <div className="alert alert-info">
          <Info className="w-4 h-4" />
          <span className="text-sm">
            Для этого атрибута требуется специальный компонент:{' '}
            {attribute.custom_component}
          </span>
        </div>
      );
    }

    // Multiselect с чекбоксами
    if (attribute.attribute_type === 'multiselect' && options.length > 0) {
      const selectedValues = currentValue?.text_value
        ? currentValue.text_value.split(',')
        : [];

      return (
        <div className="dropdown dropdown-end w-full">
          <label
            tabIndex={0}
            className="input input-bordered input-sm w-full cursor-pointer flex items-center justify-between"
          >
            <span className="text-sm truncate">
              {selectedValues.length > 0
                ? t('selectedCount', { count: selectedValues.length })
                : t('selectOptions')}
            </span>
            <ChevronDown className="w-4 h-4 flex-shrink-0" />
          </label>
          <div
            tabIndex={0}
            className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-full max-h-60 overflow-y-auto"
          >
            {options.map((option, index) => (
              <label
                key={index}
                className="label cursor-pointer justify-start hover:bg-base-200 rounded px-2"
              >
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm mr-2"
                  checked={selectedValues.includes(option.value)}
                  onChange={(e) => {
                    let newValues = [...selectedValues];
                    if (e.target.checked) {
                      newValues.push(option.value);
                    } else {
                      newValues = newValues.filter((v) => v !== option.value);
                    }
                    handleAttributeChange(attribute, newValues.join(','));
                  }}
                />
                <span className="label-text">{option.label}</span>
              </label>
            ))}
          </div>
        </div>
      );
    }

    // Обычный селект с опциями
    if (options.length > 0) {
      return (
        <select
          className="select select-bordered select-sm w-full"
          value={currentValue?.text_value || ''}
          onChange={(e) => handleAttributeChange(attribute, e.target.value)}
        >
          <option value="">{t('select')}</option>
          {options.map((option, index) => (
            <option key={index} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      );
    }

    // Числовой ввод
    if (
      attribute.attribute_type === 'numeric' ||
      attribute.attribute_type === 'number'
    ) {
      return (
        <div className="flex items-center space-x-2">
          <input
            type="number"
            className="input input-bordered input-sm flex-1"
            placeholder={t('enterValue', {
              field: attributeName.toLowerCase(),
            })}
            value={currentValue?.numeric_value || ''}
            onChange={(e) => handleAttributeChange(attribute, e.target.value)}
          />
          {attribute.validation_rules && (
            <span className="text-xs text-base-content/60">
              {/* Здесь можно добавить отображение единиц измерения */}
            </span>
          )}
        </div>
      );
    }

    // Boolean чекбокс
    if (attribute.attribute_type === 'boolean') {
      return (
        <label className="label cursor-pointer justify-start">
          <input
            type="checkbox"
            className="checkbox checkbox-sm mr-2"
            checked={currentValue?.boolean_value || false}
            onChange={(e) => handleAttributeChange(attribute, e.target.checked)}
          />
          <span className="label-text">{t('yes')}</span>
        </label>
      );
    }

    // Date picker
    if (attribute.attribute_type === 'date') {
      return (
        <input
          type="date"
          className="input input-bordered input-sm w-full"
          value={currentValue?.text_value || ''}
          onChange={(e) => handleAttributeChange(attribute, e.target.value)}
        />
      );
    }

    // Текстовый ввод (по умолчанию)
    return (
      <input
        type="text"
        className="input input-bordered input-sm w-full"
        placeholder={`Введите ${attributeName.toLowerCase()}`}
        value={currentValue?.text_value || ''}
        onChange={(e) => handleAttributeChange(attribute, e.target.value)}
      />
    );
  };

  const toggleGroup = (groupName: string) => {
    const newExpanded = new Set(expandedGroups);
    if (newExpanded.has(groupName)) {
      newExpanded.delete(groupName);
    } else {
      newExpanded.add(groupName);
    }
    setExpandedGroups(newExpanded);
  };

  // Группируем атрибуты
  const groupedAttributes = {
    required: categoryAttributes.filter((attr) => attr.is_required),
    optional: categoryAttributes.filter((attr) => !attr.is_required),
  };

  if (!selectedCategory) {
    return null;
  }

  if (loading) {
    return (
      <div className="card bg-base-200 animate-pulse">
        <div className="card-body">
          <div className="flex items-center space-x-2">
            <div className="w-5 h-5 bg-base-300 rounded"></div>
            <div className="w-32 h-4 bg-base-300 rounded"></div>
          </div>
          <div className="space-y-3 mt-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="space-y-2">
                <div className="w-24 h-3 bg-base-300 rounded"></div>
                <div className="w-full h-8 bg-base-300 rounded"></div>
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (categoryAttributes.length === 0) {
    return (
      <div className="card bg-base-200">
        <div className="card-body">
          <h3 className="card-title text-base">
            <Package className="w-5 h-5" />
            {tListing('additionalInfo')}
          </h3>
          <div className="alert alert-info">
            <Info className="w-4 h-4" />
            <span className="text-sm">
              {t('noAttributesForCategory', {
                category:
                  selectedCategory.translations?.name ||
                  selectedCategory.name ||
                  '',
              })}
            </span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="card bg-base-200"
    >
      <div className="card-body">
        <h3 className="card-title text-base">
          <Package className="w-5 h-5" />
          {tListing('additionalInfo')}
          <div className="badge badge-primary badge-sm">
            {selectedCategory.translations?.name || selectedCategory.name}
          </div>
        </h3>

        <div className="space-y-4">
          {/* Обязательные атрибуты */}
          {groupedAttributes.required.length > 0 && (
            <div>
              <div
                className="flex items-center justify-between cursor-pointer mb-3"
                onClick={() => toggleGroup('required')}
              >
                <h4 className="font-semibold text-sm flex items-center">
                  <Star className="w-4 h-4 mr-1 text-error" />
                  {t('requiredFields')}
                  <span className="badge badge-error badge-xs ml-2">
                    {groupedAttributes.required.length}
                  </span>
                </h4>
                <ChevronDown
                  className={`w-4 h-4 transition-transform ${
                    expandedGroups.has('required') ? 'rotate-180' : ''
                  }`}
                />
              </div>

              <AnimatePresence>
                {expandedGroups.has('required') && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    exit={{ opacity: 0, height: 0 }}
                    className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-4"
                  >
                    {groupedAttributes.required.map((attr) => (
                      <div key={attr.id} className="form-control">
                        <label className="label">
                          <span className="label-text">
                            {getAttributeName(attr)}
                            <span className="text-error ml-1">*</span>
                          </span>
                          {attr.icon && (
                            <span className="label-text-alt">{attr.icon}</span>
                          )}
                        </label>
                        {renderAttributeInput(attr)}
                      </div>
                    ))}
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          )}

          {/* Дополнительные атрибуты */}
          {groupedAttributes.optional.length > 0 && (
            <div>
              <div
                className="flex items-center justify-between cursor-pointer mb-3"
                onClick={() => toggleGroup('optional')}
              >
                <h4 className="font-semibold text-sm">
                  {t('optionalFields')}
                  <span className="badge badge-neutral badge-xs ml-2">
                    {groupedAttributes.optional.length}
                  </span>
                </h4>
                <ChevronDown
                  className={`w-4 h-4 transition-transform ${
                    expandedGroups.has('optional') ? 'rotate-180' : ''
                  }`}
                />
              </div>

              <AnimatePresence>
                {expandedGroups.has('optional') && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    exit={{ opacity: 0, height: 0 }}
                    className="grid grid-cols-1 lg:grid-cols-2 gap-4"
                  >
                    {groupedAttributes.optional.map((attr) => (
                      <div key={attr.id} className="form-control">
                        <label className="label">
                          <span className="label-text">
                            {getAttributeName(attr)}
                          </span>
                          {attr.icon && (
                            <span className="label-text-alt">{attr.icon}</span>
                          )}
                        </label>
                        {renderAttributeInput(attr)}
                      </div>
                    ))}
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          )}
        </div>

        {/* Подсказка о влиянии на поиск */}
        <div className="alert alert-info mt-4">
          <Info className="w-4 h-4" />
          <div className="text-sm">
            <div className="font-semibold">
              {tListing('whyAttributesImportant')}
            </div>
            <div>{tListing('attributesHelpBuyers')}</div>
          </div>
        </div>
      </div>
    </motion.div>
  );
}
