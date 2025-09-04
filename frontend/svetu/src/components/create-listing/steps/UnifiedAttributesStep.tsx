'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';
import { UnifiedAttributeField } from '@/components/shared/UnifiedAttributeField';
import { unifiedAttributeService } from '@/services/unifiedAttributeService';
import { CarSelector } from '@/components/cars/CarSelector';
import type { CarSelection } from '@/types/cars';
import type { components } from '@/types/generated/api';

type UnifiedAttribute =
  components['schemas']['backend_internal_domain_models.UnifiedAttribute'];
type UnifiedAttributeValue =
  components['schemas']['backend_internal_domain_models.UnifiedAttributeValue'];

interface AttributeGroup {
  id: string;
  name: string;
  icon: string;
  attributes: UnifiedAttribute[];
  priority: number;
}

interface UnifiedAttributesStepProps {
  onNext: () => void;
  onBack: () => void;
}

// Cache для атрибутов категорий
const attributesCache = new Map<number, UnifiedAttribute[]>();
const CACHE_TTL = 5 * 60 * 1000; // 5 минут
const cacheTimestamps = new Map<number, number>();

export default function UnifiedAttributesStep({
  onNext,
  onBack,
}: UnifiedAttributesStepProps) {
  const t = useTranslations('create_listing');
  const tCommon = useTranslations('common');
  const { state, dispatch } = useCreateListing();

  const [attributes, setAttributes] = useState<UnifiedAttribute[]>([]);
  const [attributeValues, setAttributeValues] = useState<
    Record<number, UnifiedAttributeValue>
  >({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(
    new Set(['basic', 'technical'])
  );
  const [carSelection, setCarSelection] = useState<CarSelection>({});
  const [validationErrors, setValidationErrors] = useState<
    Record<number, string>
  >({});

  // Проверяем является ли категория автомобильной
  const isAutomotiveCategory = useMemo(() => {
    return state.category
      ? state.category.id >= 10100 && state.category.id < 10200
      : false;
  }, [state.category]);

  // Загрузка атрибутов категории с кешированием
  useEffect(() => {
    const loadAttributes = async () => {
      if (!state.category) {
        setLoading(false);
        return;
      }

      const categoryId = state.category.id;
      setLoading(true);
      setError(null);

      try {
        // Проверяем кеш
        const cachedData = attributesCache.get(categoryId);
        const cachedTime = cacheTimestamps.get(categoryId);
        const now = Date.now();

        if (cachedData && cachedTime && now - cachedTime < CACHE_TTL) {
          console.log('Using cached attributes for category:', categoryId);
          setAttributes(cachedData);
          setLoading(false);
          return;
        }

        // Загружаем атрибуты через унифицированный сервис
        const response =
          await unifiedAttributeService.getCategoryAttributes(categoryId);

        if (response.success && response.data) {
          // Фильтруем только активные атрибуты
          const activeAttributes = response.data.filter(
            (attr) => attr.is_active !== false
          );

          // Сохраняем в кеш
          attributesCache.set(categoryId, activeAttributes);
          cacheTimestamps.set(categoryId, now);

          setAttributes(activeAttributes);

          // Инициализируем значения для обязательных атрибутов
          const initialValues: Record<number, UnifiedAttributeValue> = {};
          activeAttributes.forEach((attr) => {
            if (attr.id && attr.is_required) {
              initialValues[attr.id] = {
                attribute_id: attr.id,
                text_value: '',
                numeric_value: undefined,
                boolean_value: undefined,
                date_value: undefined,
                json_value: undefined,
              };
            }
          });

          // Объединяем с существующими значениями из контекста
          const existingValues = state.unifiedAttributes || {};
          setAttributeValues({ ...initialValues, ...existingValues });
        } else {
          throw new Error(response.error || 'Failed to load attributes');
        }
      } catch (err) {
        console.error('Error loading attributes:', err);
        setError(t('attributes.load_error'));

        // Пробуем fallback на v1 API
        try {
          unifiedAttributeService.useV1Api();
          const fallbackResponse =
            await unifiedAttributeService.getCategoryAttributes(categoryId);
          if (fallbackResponse.success && fallbackResponse.data) {
            setAttributes(fallbackResponse.data);
          }
        } catch (fallbackErr) {
          console.error('Fallback also failed:', fallbackErr);
        }
      } finally {
        setLoading(false);
      }
    };

    loadAttributes();
  }, [state.category, t, state.unifiedAttributes]);

  // Сохраняем значения атрибутов в контексте
  useEffect(() => {
    dispatch({ type: 'SET_UNIFIED_ATTRIBUTES', payload: attributeValues });
  }, [attributeValues, dispatch]);

  // Обновляем атрибуты при изменении выбора автомобиля
  useEffect(() => {
    if (isAutomotiveCategory && carSelection.make) {
      const makeAttr = attributes.find((a) => a.name === 'car_make_id');
      const modelAttr = attributes.find((a) => a.name === 'car_model_id');

      setAttributeValues((prev) => ({
        ...prev,
        ...(makeAttr?.id &&
          carSelection.make && {
            [makeAttr.id]: {
              attribute_id: makeAttr.id,
              numeric_value: carSelection.make.id,
              display_value: carSelection.make.name,
            },
          }),
        ...(carSelection.model &&
          modelAttr?.id && {
            [modelAttr.id]: {
              attribute_id: modelAttr.id,
              numeric_value: carSelection.model.id,
              display_value: carSelection.model.name,
            },
          }),
      }));
    }
  }, [carSelection, isAutomotiveCategory, attributes]);

  // Группировка атрибутов
  const groupedAttributes = useMemo((): AttributeGroup[] => {
    const groups = new Map<string, AttributeGroup>();

    const predefinedGroups: Record<
      string,
      { name: string; icon: string; priority: number }
    > = {
      basic: { name: t('attributes.groups.basic'), icon: '🏷️', priority: 1 },
      technical: {
        name: t('attributes.groups.technical'),
        icon: '⚙️',
        priority: 2,
      },
      condition: {
        name: t('attributes.groups.condition'),
        icon: '✨',
        priority: 3,
      },
      accessories: {
        name: t('attributes.groups.accessories'),
        icon: '📦',
        priority: 4,
      },
      dimensions: {
        name: t('attributes.groups.dimensions'),
        icon: '📏',
        priority: 5,
      },
      other: { name: t('attributes.groups.other'), icon: '📋', priority: 6 },
    };

    attributes.forEach((attr) => {
      if (!attr.name) return;

      let groupId = 'other';
      const name = attr.name.toLowerCase();

      // Определение группы по имени атрибута
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
        ['accessories', 'included', 'box', 'charger'].some((key) =>
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

      if (!groups.has(groupId)) {
        const groupInfo = predefinedGroups[groupId] || predefinedGroups.other;
        groups.set(groupId, {
          id: groupId,
          name: groupInfo.name,
          icon: groupInfo.icon,
          priority: groupInfo.priority,
          attributes: [],
        });
      }

      groups.get(groupId)!.attributes.push(attr);
    });

    // Сортировка групп и атрибутов
    return Array.from(groups.values())
      .sort((a, b) => a.priority - b.priority)
      .map((group) => ({
        ...group,
        attributes: group.attributes.sort(
          (a, b) => (a.sort_order || 0) - (b.sort_order || 0)
        ),
      }));
  }, [attributes, t]);

  // Автоматическое разворачивание групп с обязательными полями
  useEffect(() => {
    const groupsWithRequired = groupedAttributes
      .filter((group) => group.attributes.some((attr) => attr.is_required))
      .map((group) => group.id);

    setExpandedGroups(new Set(['basic', 'technical', ...groupsWithRequired]));
  }, [groupedAttributes]);

  // Обработчик изменения значения атрибута
  const handleAttributeChange = useCallback(
    (attributeId: number, value: UnifiedAttributeValue) => {
      setAttributeValues((prev) => ({
        ...prev,
        [attributeId]: value,
      }));

      // Очищаем ошибку валидации при изменении
      setValidationErrors((prev) => {
        const next = { ...prev };
        delete next[attributeId];
        return next;
      });
    },
    []
  );

  // Валидация обязательных полей
  const validateRequiredFields = useCallback(() => {
    const errors: Record<number, string> = {};
    let isValid = true;

    attributes.forEach((attr) => {
      if (attr.id && attr.is_required) {
        const value = attributeValues[attr.id];

        // Специальная проверка для автомобильных атрибутов
        if (isAutomotiveCategory) {
          if (attr.name === 'car_make_id' && !carSelection.make) {
            errors[attr.id] = t('attributes.required_field');
            isValid = false;
            return;
          }
          if (attr.name === 'car_model_id' && !carSelection.model) {
            errors[attr.id] = t('attributes.required_field');
            isValid = false;
            return;
          }
        }

        // Обычная проверка
        if (
          !value ||
          (!value.text_value &&
            value.numeric_value === undefined &&
            value.boolean_value === undefined &&
            !value.date_value &&
            (!value.json_value ||
              (Array.isArray(value.json_value) &&
                value.json_value.length === 0)))
        ) {
          errors[attr.id] = t('attributes.required_field');
          isValid = false;
        }
      }
    });

    setValidationErrors(errors);
    return isValid;
  }, [attributes, attributeValues, isAutomotiveCategory, carSelection, t]);

  // Обработчик кнопки "Далее"
  const handleNext = () => {
    if (validateRequiredFields()) {
      onNext();
    }
  };

  // Переключение группы
  const toggleGroup = (groupId: string) => {
    setExpandedGroups((prev) => {
      const next = new Set(prev);
      if (next.has(groupId)) {
        next.delete(groupId);
      } else {
        next.add(groupId);
      }
      return next;
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="stroke-current shrink-0 h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <span>{error}</span>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            🏷️ {t('attributes.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('attributes.description')}
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
                />
              </svg>
              <span>{t('attributes.none_required')}</span>
            </div>
          ) : (
            <div className="space-y-6 mb-8">
              {/* CarSelector для автомобильных категорий */}
              {isAutomotiveCategory && (
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body">
                    <h3 className="card-title text-xl flex items-center gap-3">
                      <span className="text-2xl">🚗</span>
                      {t('attributes.groups.car_selection')}
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

              {/* Информация об обязательных полях */}
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
                    />
                  </svg>
                  <div>
                    <h3 className="font-bold">
                      {t('attributes.required_info')}
                    </h3>
                    <div className="text-xs">
                      {t('attributes.required_expanded')}
                    </div>
                  </div>
                </div>
              )}

              {/* Группы атрибутов */}
              <div className="grid grid-cols-1 gap-4">
                {groupedAttributes.map((group) => {
                  const isExpanded = expandedGroups.has(group.id);
                  const hasRequired = group.attributes.some(
                    (attr) => attr.is_required
                  );
                  const allRequiredFilled = group.attributes
                    .filter((attr) => attr.is_required)
                    .every((attr) => {
                      if (!attr.id) return true;
                      const value = attributeValues[attr.id];
                      return (
                        value &&
                        (value.text_value ||
                          value.numeric_value !== undefined ||
                          value.boolean_value !== undefined ||
                          value.date_value ||
                          (value.json_value && value.json_value.length > 0))
                      );
                    });

                  return (
                    <div key={group.id} className="card bg-base-100 shadow-lg">
                      <div
                        className="card-body cursor-pointer select-none"
                        onClick={() => toggleGroup(group.id)}
                      >
                        <div className="flex items-center justify-between">
                          <h3 className="card-title text-xl flex items-center gap-3">
                            <span className="text-2xl">{group.icon}</span>
                            {group.name}
                            <div className="badge badge-neutral">
                              {group.attributes.length}
                            </div>
                            {hasRequired && (
                              <div
                                className={`badge ${allRequiredFilled ? 'badge-success' : 'badge-warning'}`}
                              >
                                {allRequiredFilled ? '✓' : tCommon('required')}
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
                            {group.attributes.map((attr) => {
                              // Пропускаем автомобильные атрибуты, управляемые через CarSelector
                              if (
                                isAutomotiveCategory &&
                                (attr.name === 'car_make_id' ||
                                  attr.name === 'car_model_id')
                              ) {
                                return null;
                              }

                              if (!attr.id) return null;

                              return (
                                <UnifiedAttributeField
                                  key={attr.id}
                                  attribute={attr}
                                  value={attributeValues[attr.id]}
                                  onChange={(value) =>
                                    handleAttributeChange(attr.id!, value)
                                  }
                                  error={validationErrors[attr.id]}
                                  required={attr.is_required}
                                />
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

          {/* Подсказка */}
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
              />
            </svg>
            <div className="text-sm">
              <p className="font-medium">💡 {t('attributes.tip')}</p>
              <p className="text-xs mt-1">{t('attributes.tip_description')}</p>
            </div>
          </div>

          {/* Навигационные кнопки */}
          <div className="card-actions justify-between mt-6">
            <button className="btn btn-outline" onClick={onBack}>
              ← {tCommon('back')}
            </button>
            <button className="btn btn-primary" onClick={handleNext}>
              {tCommon('continue')} →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
