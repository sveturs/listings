'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';
import type { components } from '@/types/generated/api';

type UnifiedAttribute =
  components['schemas']['backend_internal_domain_models.UnifiedAttribute'];

interface ProductVariant {
  id: string;
  attributes: Record<string, string>;
  price?: number;
  stock?: number;
  sku?: string;
  image?: string;
}

interface VariantsStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function VariantsStep({ onNext, onBack }: VariantsStepProps) {
  const t = useTranslations('create_listing');
  const tCommon = useTranslations('common');
  const { state, dispatch } = useCreateListing();

  const [variants, setVariants] = useState<ProductVariant[]>([]);
  const [selectedVariantAttributes, setSelectedVariantAttributes] = useState<
    string[]
  >([]);
  const [loading, setLoading] = useState(false);
  const [attributeValues, setAttributeValues] = useState<
    Record<string, string[]>
  >({});

  // Получаем вариантные атрибуты из унифицированных атрибутов
  const variantCompatibleAttributes = useMemo(() => {
    if (!state.unifiedAttributes) return [];

    const attributes: UnifiedAttribute[] = [];
    Object.values(state.unifiedAttributes).forEach((value) => {
      // Здесь нужно получить информацию об атрибуте и проверить is_variant_compatible
      // Для демо используем базовые атрибуты
      if (value.attribute_id) {
        // Заглушка - в реальности получать из API
        const attr: UnifiedAttribute = {
          id: value.attribute_id,
          name: `attribute_${value.attribute_id}`,
          display_name:
            value.display_value || `Attribute ${value.attribute_id}`,
          is_variant_compatible: true,
          input_type: 'select',
          is_required: false,
          is_active: true,
        };
        attributes.push(attr);
      }
    });

    return attributes;
  }, [state.unifiedAttributes]);

  // Генерация всех возможных комбинаций вариантов
  const generateVariantCombinations = useCallback(() => {
    if (selectedVariantAttributes.length === 0) {
      setVariants([]);
      return;
    }

    const combinations: ProductVariant[] = [];
    const attributeValueArrays = selectedVariantAttributes.map(
      (attrName) => attributeValues[attrName] || []
    );

    // Функция для генерации всех комбинаций (декартово произведение)
    const generateCombinations = (
      arrays: string[][],
      current: string[] = [],
      index = 0
    ): void => {
      if (index === arrays.length) {
        const variantId = current.join('_').toLowerCase().replace(/\s+/g, '_');
        const attributes: Record<string, string> = {};
        selectedVariantAttributes.forEach((attrName, idx) => {
          attributes[attrName] = current[idx];
        });

        combinations.push({
          id: variantId,
          attributes,
          price: state.price || 0,
          stock: 0,
          sku: `${state.title?.substring(0, 3).toUpperCase() || 'PROD'}-${variantId.substring(0, 8).toUpperCase()}`,
        });
        return;
      }

      const array = arrays[index];
      for (const value of array) {
        generateCombinations(arrays, [...current, value], index + 1);
      }
    };

    generateCombinations(attributeValueArrays);
    setVariants(combinations);
  }, [selectedVariantAttributes, attributeValues, state.price, state.title]);

  // Обновление вариантов при изменении параметров
  useEffect(() => {
    generateVariantCombinations();
  }, [generateVariantCombinations]);

  // Добавление нового значения атрибута
  const addAttributeValue = (attributeName: string, value: string) => {
    if (!value.trim()) return;

    setAttributeValues((prev) => ({
      ...prev,
      [attributeName]: [...(prev[attributeName] || []), value.trim()],
    }));
  };

  // Удаление значения атрибута
  const removeAttributeValue = (attributeName: string, value: string) => {
    setAttributeValues((prev) => ({
      ...prev,
      [attributeName]: (prev[attributeName] || []).filter((v) => v !== value),
    }));
  };

  // Добавление атрибута к вариантам
  const toggleVariantAttribute = (attributeName: string) => {
    setSelectedVariantAttributes((prev) => {
      const isSelected = prev.includes(attributeName);
      if (isSelected) {
        // Удаляем атрибут
        const filtered = prev.filter((name) => name !== attributeName);
        // Также удаляем его значения
        setAttributeValues((current) => {
          const newValues = { ...current };
          delete newValues[attributeName];
          return newValues;
        });
        return filtered;
      } else {
        // Добавляем атрибут
        return [...prev, attributeName];
      }
    });
  };

  // Обновление варианта
  const updateVariant = (
    variantId: string,
    field: keyof ProductVariant,
    value: any
  ) => {
    setVariants((prev) =>
      prev.map((variant) =>
        variant.id === variantId ? { ...variant, [field]: value } : variant
      )
    );
  };

  // Сохранение в контексте и переход дальше
  const handleNext = () => {
    dispatch({ type: 'SET_PRODUCT_VARIANTS', payload: variants });
    onNext();
  };

  // Применить цену ко всем вариантам
  const applyPriceToAll = () => {
    if (!state.price) return;

    setVariants((prev) =>
      prev.map((variant) => ({ ...variant, price: state.price || 0 }))
    );
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            🎨 {t('variants.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('variants.description')}
          </p>

          {variantCompatibleAttributes.length === 0 ? (
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
              <span>{t('variants.no_variant_attributes')}</span>
            </div>
          ) : (
            <div className="space-y-6">
              {/* Выбор атрибутов для вариантов */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="card-title text-lg">
                    📋 {t('variants.select_attributes')}
                  </h3>
                  <p className="text-sm text-base-content/70 mb-4">
                    {t('variants.select_attributes_description')}
                  </p>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {variantCompatibleAttributes.map((attr) => (
                      <div
                        key={attr.name}
                        className={`card cursor-pointer transition-all ${
                          selectedVariantAttributes.includes(attr.name!)
                            ? 'bg-primary text-primary-content'
                            : 'bg-base-100 hover:bg-base-300'
                        }`}
                        onClick={() => toggleVariantAttribute(attr.name!)}
                      >
                        <div className="card-body p-4">
                          <div className="flex items-center justify-between">
                            <span className="font-medium">
                              {attr.display_name}
                            </span>
                            <input
                              type="checkbox"
                              className="checkbox"
                              checked={selectedVariantAttributes.includes(
                                attr.name!
                              )}
                              readOnly
                            />
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Управление значениями атрибутов */}
              {selectedVariantAttributes.length > 0 && (
                <div className="card bg-base-200">
                  <div className="card-body">
                    <h3 className="card-title text-lg">
                      🎯 {t('variants.attribute_values')}
                    </h3>
                    <p className="text-sm text-base-content/70 mb-4">
                      {t('variants.attribute_values_description')}
                    </p>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      {selectedVariantAttributes.map((attrName) => {
                        const attr = variantCompatibleAttributes.find(
                          (a) => a.name === attrName
                        );
                        return (
                          <div key={attrName} className="space-y-3">
                            <h4 className="font-medium">
                              {attr?.display_name || attrName}
                            </h4>

                            {/* Добавление нового значения */}
                            <div className="flex gap-2">
                              <input
                                type="text"
                                className="input input-sm input-bordered flex-1"
                                placeholder={t('variants.add_value')}
                                onKeyDown={(e) => {
                                  if (e.key === 'Enter') {
                                    const input = e.target as HTMLInputElement;
                                    addAttributeValue(attrName, input.value);
                                    input.value = '';
                                  }
                                }}
                              />
                              <button
                                className="btn btn-sm btn-primary"
                                onClick={(e) => {
                                  const input = (
                                    e.target as HTMLElement
                                  ).parentElement?.querySelector('input');
                                  if (input?.value) {
                                    addAttributeValue(attrName, input.value);
                                    input.value = '';
                                  }
                                }}
                              >
                                +
                              </button>
                            </div>

                            {/* Список значений */}
                            <div className="flex flex-wrap gap-2">
                              {(attributeValues[attrName] || []).map(
                                (value, index) => (
                                  <div
                                    key={index}
                                    className="badge badge-outline gap-2"
                                  >
                                    {value}
                                    <button
                                      className="btn btn-ghost btn-xs"
                                      onClick={() =>
                                        removeAttributeValue(attrName, value)
                                      }
                                    >
                                      ×
                                    </button>
                                  </div>
                                )
                              )}
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                </div>
              )}

              {/* Матрица вариантов */}
              {variants.length > 0 && (
                <div className="card bg-base-100">
                  <div className="card-body">
                    <div className="flex items-center justify-between mb-4">
                      <h3 className="card-title text-lg">
                        🔀 {t('variants.generated_variants')} ({variants.length}
                        )
                      </h3>
                      <button
                        className="btn btn-sm btn-outline"
                        onClick={applyPriceToAll}
                      >
                        {t('variants.apply_base_price')} ({state.price || 0}{' '}
                        RSD)
                      </button>
                    </div>

                    <div className="overflow-x-auto">
                      <table className="table table-zebra">
                        <thead>
                          <tr>
                            <th>SKU</th>
                            {selectedVariantAttributes.map((attr) => (
                              <th key={attr}>
                                {variantCompatibleAttributes.find(
                                  (a) => a.name === attr
                                )?.display_name || attr}
                              </th>
                            ))}
                            <th>{t('variants.price')}</th>
                            <th>{t('variants.stock')}</th>
                          </tr>
                        </thead>
                        <tbody>
                          {variants.map((variant) => (
                            <tr key={variant.id}>
                              <td>
                                <input
                                  type="text"
                                  className="input input-xs input-bordered w-full"
                                  value={variant.sku || ''}
                                  onChange={(e) =>
                                    updateVariant(
                                      variant.id,
                                      'sku',
                                      e.target.value
                                    )
                                  }
                                />
                              </td>
                              {selectedVariantAttributes.map((attr) => (
                                <td key={attr} className="font-medium">
                                  {variant.attributes[attr]}
                                </td>
                              ))}
                              <td>
                                <div className="flex items-center gap-1">
                                  <input
                                    type="number"
                                    className="input input-xs input-bordered w-20"
                                    value={variant.price || 0}
                                    onChange={(e) =>
                                      updateVariant(
                                        variant.id,
                                        'price',
                                        Number(e.target.value)
                                      )
                                    }
                                  />
                                  <span className="text-xs">RSD</span>
                                </div>
                              </td>
                              <td>
                                <input
                                  type="number"
                                  className="input input-xs input-bordered w-16"
                                  value={variant.stock || 0}
                                  onChange={(e) =>
                                    updateVariant(
                                      variant.id,
                                      'stock',
                                      Number(e.target.value)
                                    )
                                  }
                                />
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>

                    <div className="alert alert-info mt-4">
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
                        <p className="font-medium">💡 {t('variants.tip')}</p>
                        <p>{t('variants.tip_description')}</p>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Кнопки навигации */}
          <div className="card-actions justify-between mt-8">
            <button className="btn btn-outline" onClick={onBack}>
              ← {tCommon('back')}
            </button>
            <button className="btn btn-primary" onClick={handleNext}>
              {variants.length > 0
                ? `${tCommon('continue')} (${variants.length} ${t('variants.variants')})`
                : tCommon('skip')}{' '}
              →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
