'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';
import {
  MarketplaceService,
  CategoryAttributeMapping,
} from '@/services/marketplace';
import { getTranslatedAttribute } from '@/utils/translatedAttribute';
import { CarSelector } from '@/components/cars/CarSelector';
import type { CarSelection } from '@/types/cars';

interface AttributeFormData {
  attribute_id: number;
  attribute_name: string;
  display_name: string;
  attribute_type: string;
  text_value?: string;
  numeric_value?: number;
  boolean_value?: boolean;
  json_value?: any;
  display_value?: string;
  unit?: string;
}

interface AttributeGroup {
  id: string;
  name: string;
  icon: string;
  attributes: CategoryAttributeMapping[];
}

interface AttributesStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function AttributesStep({
  onNext,
  onBack,
}: AttributesStepProps) {
  const t = useTranslations('create_listing.attributes.groups');
  const tCommon = useTranslations('common');
  const tCreate_listing.attributes = useTranslations('create_listing.attributes');
  const locale = useLocale();
  const { state, dispatch } = useCreateListing();
  const [attributes, setAttributes] = useState<CategoryAttributeMapping[]>([]);
  const [formData, setFormData] = useState<Record<number, AttributeFormData>>(
    state.attributes || {}
  );
  const [loading, setLoading] = useState(true);
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(
    new Set(['basic', 'technical'])
  );
  const [carSelection, setCarSelection] = useState<CarSelection>({});
  const [isAutomotiveCategory, setIsAutomotiveCategory] = useState(false);

  useEffect(() => {
    const loadAttributes = async () => {
      if (!state.category) {
        setLoading(false);
        return;
      }

      try {
        setLoading(true);

        // Проверяем, является ли категория автомобильной (ID >= 10100 и < 10200)
        const isAuto = state.category.id >= 10100 && state.category.id < 10200;
        setIsAutomotiveCategory(isAuto);

        const response = await MarketplaceService.getCategoryAttributes(
          state.category.id
        );

        if (response.success && response.data) {
          // Преобразуем атрибуты в CategoryAttributeMapping формат
          const attributeMappings: CategoryAttributeMapping[] =
            response.data.map((attr) => ({
              category_id: state.category!.id,
              attribute_id: attr.id,
              is_enabled: true, // Все возвращенные атрибуты включены
              is_required: attr.is_required,
              sort_order: attr.sort_order,
              attribute: attr,
            }));

          // Фильтруем дубликаты по attribute_id и сортируем
          const uniqueAttributes = attributeMappings
            .filter(
              (attr, index, self) =>
                index ===
                self.findIndex((a) => a.attribute?.id === attr.attribute?.id)
            )
            .sort((a, b) => a.sort_order - b.sort_order);

          // Логируем атрибуты для отладки
          console.log('Loaded attributes:', uniqueAttributes);
          console.log('Current locale:', locale);
          uniqueAttributes.forEach((attr) => {
            if (attr.attribute) {
              console.log(`Attribute ${attr.attribute.name}:`, {
                display_name: attr.attribute.display_name,
                attribute_type: attr.attribute.attribute_type,
                options: attr.attribute.options,
                translations: attr.attribute.translations,
                option_translations: attr.attribute.option_translations,
              });
            }
          });

          setAttributes(uniqueAttributes);

          // Инициализируем formData для новых атрибутов
          const initialFormData: Record<number, AttributeFormData> = {};
          uniqueAttributes.forEach((mapping) => {
            if (mapping.attribute && !formData[mapping.attribute.id]) {
              initialFormData[mapping.attribute.id] = {
                attribute_id: mapping.attribute.id,
                attribute_name: mapping.attribute.name,
                display_name: mapping.attribute.display_name,
                attribute_type: mapping.attribute.attribute_type,
              };
            }
          });

          // Объединяем с существующими данными
          setFormData((prev) => ({ ...initialFormData, ...prev }));
        }
      } catch (error) {
        console.error('Error loading attributes:', error);
        // В случае ошибки показываем пустой список атрибутов
        setAttributes([]);
      } finally {
        setLoading(false);
      }
    };

    loadAttributes();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [state.category]);

  useEffect(() => {
    dispatch({ type: 'SET_ATTRIBUTES', payload: formData });
  }, [formData, dispatch]);

  // Обновляем атрибуты при изменении выбора автомобиля
  useEffect(() => {
    if (isAutomotiveCategory && carSelection.make) {
      const makeAttr = attributes.find(
        (a) => a.attribute?.name === 'car_make_id'
      );
      const modelAttr = attributes.find(
        (a) => a.attribute?.name === 'car_model_id'
      );
      const selectedMake = carSelection.make; // Сохраняем ссылку для TypeScript
      const selectedModel = carSelection.model;

      setFormData((prev) => ({
        ...prev,
        ...(makeAttr?.attribute?.id && {
          [makeAttr.attribute.id]: {
            attribute_id: makeAttr.attribute.id,
            attribute_name: 'car_make_id',
            display_name: 'Car Make ID',
            attribute_type: 'number',
            numeric_value: selectedMake.id,
            display_value: selectedMake.name,
          },
        }),
        ...(selectedModel &&
          modelAttr?.attribute?.id && {
            [modelAttr.attribute.id]: {
              attribute_id: modelAttr.attribute.id,
              attribute_name: 'car_model_id',
              display_name: 'Car Model ID',
              attribute_type: 'number',
              numeric_value: selectedModel.id,
              display_value: selectedModel.name,
            },
          }),
      }));
    }
  }, [carSelection, isAutomotiveCategory, attributes]);

  // Группируем атрибуты по логическим группам
  const groupAttributes = useCallback((): AttributeGroup[] => {
    const groupsMap = new Map<string, AttributeGroup>();

    // Предопределенные группы с иконками
    const predefinedGroups: Record<
      string,
      { name: string; icon: string; priority: number }
    > = {
      basic: {
        name: tCreate_listing.attributes('groups.basic'),
        icon: '🏷️',
        priority: 1,
      },
      technical: {
        name: tCreate_listing.attributes('groups.technical'),
        icon: '⚙️',
        priority: 2,
      },
      condition: {
        name: tCreate_listing.attributes('groups.condition'),
        icon: '✨',
        priority: 3,
      },
      accessories: {
        name: tCreate_listing.attributes('groups.accessories'),
        icon: '📦',
        priority: 4,
      },
      dimensions: {
        name: tCreate_listing.attributes('groups.dimensions'),
        icon: '📏',
        priority: 5,
      },
      other: {
        name: tCreate_listing.attributes('groups.other'),
        icon: '📋',
        priority: 6,
      },
    };

    attributes.forEach((mapping) => {
      const attr = mapping.attribute;
      if (!attr) return;

      let groupId = 'other';
      const name = attr.name.toLowerCase();

      // Определяем группу на основе имени атрибута
      if (
        ['brand', 'model', 'type', 'category', 'name', 'title'].some((key) =>
          name.includes(key)
        )
      ) {
        groupId = 'basic';
      } else if (
        [
          'year',
          'engine',
          'fuel',
          'transmission',
          'power',
          'volume',
          'memory',
          'storage',
          'display',
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

      groupsMap.get(groupId)!.attributes.push(mapping);
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
  }, [attributes, t]);

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
  }, [attributes, groupAttributes]);

  const handleInputChange = (
    attributeId: number,
    value: any,
    attributeType: string
  ) => {
    setFormData((prev) => {
      const mapping = attributes.find((a) => a.attribute?.id === attributeId);
      const attribute = mapping?.attribute;
      if (!attribute) return prev;

      const updatedAttribute: AttributeFormData = {
        ...prev[attributeId],
        attribute_id: attributeId,
        attribute_name: attribute.name,
        display_name: attribute.display_name,
        attribute_type: attributeType,
      };

      // Очищаем все значения
      delete updatedAttribute.text_value;
      delete updatedAttribute.numeric_value;
      delete updatedAttribute.boolean_value;
      delete updatedAttribute.json_value;

      // Устанавливаем значение в соответствующее поле
      switch (attributeType) {
        case 'text':
        case 'select':
          updatedAttribute.text_value = value;
          updatedAttribute.display_value = value;
          break;
        case 'number':
          updatedAttribute.numeric_value = value;
          updatedAttribute.display_value = value.toString();
          break;
        case 'boolean':
          updatedAttribute.boolean_value = value;
          updatedAttribute.display_value = value ? 'Да' : 'Не';
          break;
        case 'multiselect':
          updatedAttribute.json_value = value;
          updatedAttribute.display_value = Array.isArray(value)
            ? value.join(', ')
            : '';
          break;
      }

      return {
        ...prev,
        [attributeId]: updatedAttribute,
      };
    });
  };

  // Универсальная функция для получения значений опций
  const getOptionValues = (options: any): string[] => {
    if (!options) return [];

    // Если options это строка JSON, пытаемся распарсить
    if (typeof options === 'string') {
      try {
        const parsed = JSON.parse(options);
        return getOptionValues(parsed); // Рекурсивный вызов с распарсенным объектом
      } catch (e) {
        console.error('Failed to parse options:', e);
        return [];
      }
    }

    // Если options это массив напрямую
    if (Array.isArray(options)) {
      return options;
    }

    // Если options это объект с полем values
    if (options.values && Array.isArray(options.values)) {
      return options.values;
    }

    // Если options это объект с другой структурой, логируем для отладки
    console.warn('Unknown options structure:', options);
    return [];
  };

  const renderAttribute = (
    mapping: CategoryAttributeMapping,
    getOptionLabel: (option: string) => string
  ) => {
    const attribute = mapping.attribute;
    if (!attribute) return null;

    // Для автомобильных категорий используем CarSelector для марки и модели
    if (
      isAutomotiveCategory &&
      (attribute.name === 'car_make_id' || attribute.name === 'car_model_id')
    ) {
      // Скрываем отдельные поля для марки и модели, так как они управляются через CarSelector
      return null;
    }

    const formAttribute = formData[attribute.id];
    const value =
      formAttribute?.text_value ||
      formAttribute?.numeric_value ||
      formAttribute?.boolean_value ||
      formAttribute?.json_value ||
      '';

    switch (attribute.attribute_type) {
      case 'text':
        return (
          <input
            type="text"
            placeholder=""
            className="input input-bordered"
            value={value}
            onChange={(e) =>
              handleInputChange(
                attribute.id,
                e.target.value,
                attribute.attribute_type
              )
            }
          />
        );

      case 'number':
        return (
          <div className="flex items-center gap-2">
            <input
              type="number"
              placeholder="0"
              className="input input-bordered flex-1"
              value={value}
              onChange={(e) =>
                handleInputChange(
                  attribute.id,
                  parseFloat(e.target.value) || 0,
                  attribute.attribute_type
                )
              }
              min="0"
            />
            {attribute.options?.step && (
              <span className="text-sm text-base-content/60">
                шаг: {attribute.options.step}
              </span>
            )}
          </div>
        );

      case 'select':
        const selectOptions = getOptionValues(attribute.options);
        return (
          <select
            className="select select-bordered"
            value={value}
            onChange={(e) =>
              handleInputChange(
                attribute.id,
                e.target.value,
                attribute.attribute_type
              )
            }
          >
            <option value="">{tCommon('select')}</option>
            {selectOptions.map((option) => (
              <option key={option} value={option}>
                {getOptionLabel(option)}
              </option>
            ))}
          </select>
        );

      case 'boolean':
        return (
          <div className="form-control">
            <label className="label cursor-pointer justify-start gap-3">
              <input
                type="checkbox"
                className="checkbox checkbox-primary"
                checked={!!value}
                onChange={(e) =>
                  handleInputChange(
                    attribute.id,
                    e.target.checked,
                    attribute.attribute_type
                  )
                }
              />
              <span className="label-text">{tCommon('yes')}</span>
            </label>
          </div>
        );

      case 'multiselect':
        const multiselectOptions = getOptionValues(attribute.options);
        return (
          <div className="space-y-2">
            {multiselectOptions.map((option) => (
              <label
                key={option}
                className="label cursor-pointer justify-start gap-3"
              >
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm"
                  checked={Array.isArray(value) && value.includes(option)}
                  onChange={(e) => {
                    const currentArray = Array.isArray(value) ? value : [];
                    const newArray = e.target.checked
                      ? [...currentArray, option]
                      : currentArray.filter((item: string) => item !== option);
                    handleInputChange(
                      attribute.id,
                      newArray,
                      attribute.attribute_type
                    );
                  }}
                />
                <span className="label-text text-sm">
                  {getOptionLabel(option)}
                </span>
              </label>
            ))}
          </div>
        );

      default:
        return null;
    }
  };

  const requiredAttributesFilled = attributes
    .filter((mapping) => mapping.is_required && mapping.attribute)
    .every((mapping) => {
      const attr = mapping.attribute!;

      // Для автомобильных атрибутов проверяем carSelection
      if (isAutomotiveCategory) {
        if (attr.name === 'car_make_id') {
          return !!carSelection.make;
        }
        if (attr.name === 'car_model_id') {
          return !!carSelection.model;
        }
      }

      const formAttr = formData[attr.id];
      if (!formAttr) return false;

      // Проверяем, что хотя бы одно значение заполнено
      return (
        formAttr.text_value !== undefined ||
        formAttr.numeric_value !== undefined ||
        formAttr.boolean_value !== undefined ||
        (formAttr.json_value &&
          Array.isArray(formAttr.json_value) &&
          formAttr.json_value.length > 0)
      );
    });

  if (loading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            🏷️ {tCreate_listing.attributes('title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {tCreate_listing.attributes('description')}
          </p>

          {attributes.length === 0 ? (
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
              <span>{tCreate_listing.attributes('none_required')}</span>
            </div>
          ) : (
            <div className="space-y-6 mb-8">
              {/* Для автомобильных категорий показываем CarSelector */}
              {isAutomotiveCategory && (
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body">
                    <h3 className="card-title text-xl flex items-center gap-3">
                      <span className="text-2xl">🚗</span>
                      {tCreate_listing.attributes('groups.car_selection')}
                      <div className="badge badge-warning">
                        {tCommon('required')}
                      </div>
                    </h3>
                    <CarSelector
                      value={carSelection}
                      onChange={setCarSelection}
                      required={true}
                      className="mt-4"
                    />
                  </div>
                </div>
              )}

              {/* Сводка обязательных полей */}
              {attributes.some((mapping) => mapping.is_required) && (
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
                    />
                  </svg>
                  <div>
                    <h3 className="font-bold">
                      {tCommon('required')} атрибуты
                    </h3>
                    <div className="text-xs">
                      Обязательные поля автоматически развернуты для удобства
                      заполнения
                    </div>
                  </div>
                </div>
              )}

              {/* Группы атрибутов */}
              <div className="grid grid-cols-1 gap-4">
                {groupAttributes().map((group) => {
                  const isExpanded = expandedGroups.has(group.id);
                  const hasRequiredFields = group.attributes.some(
                    (mapping) => mapping.is_required
                  );
                  const filledRequiredFields = group.attributes
                    .filter((mapping) => mapping.is_required)
                    .every((mapping) => {
                      const attr = mapping.attribute!;
                      const formAttr = formData[attr.id];
                      return (
                        formAttr &&
                        (formAttr.text_value !== undefined ||
                          formAttr.numeric_value !== undefined ||
                          formAttr.boolean_value !== undefined ||
                          (formAttr.json_value &&
                            Array.isArray(formAttr.json_value) &&
                            formAttr.json_value.length > 0))
                      );
                    });

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
                                {filledRequiredFields
                                  ? '✓'
                                  : tCommon('required')}
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
                        <div className="card-body pt-0">
                          <div className="space-y-4">
                            {group.attributes.map((mapping) => {
                              const attribute = mapping.attribute;
                              if (!attribute) return null;

                              const { displayName, getOptionLabel } =
                                getTranslatedAttribute(attribute, locale);

                              return (
                                <div
                                  key={attribute.id}
                                  className="form-control"
                                >
                                  <label className="label">
                                    <span className="label-text font-medium">
                                      {displayName}
                                    </span>
                                    {mapping.is_required && (
                                      <span className="label-text-alt text-error">
                                        *
                                      </span>
                                    )}
                                  </label>
                                  {renderAttribute(mapping, getOptionLabel)}
                                </div>
                              );
                            })}
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          )}

          {/* Региональная подсказка */}
          <div className="alert alert-info mt-6">
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
            <div className="text-sm">
              <p className="font-medium">
                💡 {tCreate_listing.attributes('tip')}
              </p>
              <p className="text-xs mt-1">
                {tCreate_listing.attributes('tip_description')}
              </p>
            </div>
          </div>

          {/* Кнопки навигации */}
          <div className="card-actions justify-between mt-6">
            <button className="btn btn-outline" onClick={onBack}>
              ← {tCommon('back')}
            </button>
            <button
              className={`btn btn-primary ${!requiredAttributesFilled ? 'btn-disabled' : ''}`}
              onClick={onNext}
              disabled={!requiredAttributesFilled}
            >
              {tCommon('continue')} →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
