'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import VariantStockManager from './VariantStockManager';

interface CategoryAttribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type: string;
  is_required: boolean;
  options?: any;
  values?: Array<{
    value: string;
    display_name: string;
    color_hex?: string;
  }>;
}

interface CategoryVariantGeneratorProps {
  categoryId: number;
  categoryAttributes: CategoryAttribute[];
  selectedAttributes: Record<number, any>; // Добавляем выбранные атрибуты с шага Attributes
  basePrice: number;
  onGenerate: (variants: any[]) => void;
  onCancel: () => void;
}

export default function CategoryVariantGenerator({
  categoryId: _categoryId,
  categoryAttributes,
  selectedAttributes,
  basePrice,
  onGenerate,
  onCancel,
}: CategoryVariantGeneratorProps) {
  const t = useTranslations('storefronts.products');

  // Фильтруем атрибуты, подходящие для вариантов, которые были выбраны на шаге Attributes
  const stockAttributes = React.useMemo(() => {
    console.log('Debug CategoryVariantGenerator:');
    console.log('categoryAttributes:', categoryAttributes);
    console.log('selectedAttributes:', selectedAttributes);

    // Используем атрибуты, которые обычно влияют на варианты (цвет, размер, материал etc.)
    const variantLikeAttributes = categoryAttributes.filter((attr) => {
      const name = (attr.name || attr.display_name || '').toLowerCase();
      const keywords = [
        'color',
        'colour',
        'size',
        'memory',
        'storage',
        'material',
        'style',
        'цвет',
        'размер',
        'материал',
        'стиль',
        'brand',
        'model',
        'version',
        'type',
        'finish',
        'pattern',
      ];
      return keywords.some((keyword) => name.includes(keyword));
    });
    console.log('variant-like attributes:', variantLikeAttributes);

    // Фильтруем только те, которые выбраны и для которых выбрано несколько значений
    const selectedVariantAttrs = variantLikeAttributes.filter((attr) => {
      const value = selectedAttributes[attr.id];
      if (!value || value === '') return false;

      // Проверяем, что это multiselect с несколькими значениями
      if (Array.isArray(value)) {
        return value.length > 1; // Нужно минимум 2 значения для вариантов
      }

      // Если это строка с запятыми, считаем количество
      if (typeof value === 'string' && value.includes(',')) {
        return value.split(',').filter(Boolean).length > 1;
      }

      return false;
    });

    console.log('selected variant attributes:', selectedVariantAttrs);
    return selectedVariantAttrs;
  }, [categoryAttributes, selectedAttributes]);

  const [selectedVariantValues, setSelectedVariantValues] = useState<
    Record<string, string[]>
  >({});
  const [stockSettings, setStockSettings] = useState({
    defaultQuantity: 10,
    useIndividualQuantities: false,
  });
  const [individualQuantities, _setIndividualQuantities] = useState<
    Record<string, number>
  >({});
  const [generatedVariants, setGeneratedVariants] = useState<any[]>([]);
  const [showStockManager, setShowStockManager] = useState(false);

  // Используем выбранные значения атрибутов из шага Attributes
  const [attributeValues, setAttributeValues] = useState<
    Record<number, string[]>
  >({});
  const [_loading] = useState(false);

  // Преобразуем выбранные атрибуты в формат для работы с вариантами
  React.useEffect(() => {
    const prepareAttributeValues = () => {
      const values: Record<number, string[]> = {};
      const variantValues: Record<string, string[]> = {};

      stockAttributes.forEach((attr) => {
        // Находим соответствующий атрибут по имени из selectedAttributes
        let selectedValue = null;

        // Сначала пытаемся найти по ID
        selectedValue = selectedAttributes[attr.id];

        // Если не нашли по ID, ищем по имени атрибута
        if (!selectedValue) {
          Object.entries(selectedAttributes).forEach(([key, value]) => {
            // Находим categoryAttribute по id
            const catAttrId = parseInt(key);
            if (!isNaN(catAttrId) && value) {
              // Проверяем что это атрибут цвета или размера
              if (
                (attr.name === 'color' && key === '2004') ||
                (attr.name === 'size' && key === '2801')
              ) {
                selectedValue = value;
              }
            }
          });
        }

        if (selectedValue) {
          // Для multiselect атрибутов используем массив
          if (Array.isArray(selectedValue)) {
            values[attr.id] = selectedValue;
            variantValues[attr.display_name || attr.name] = selectedValue;
          } else {
            // Для одиночных значений создаем массив
            values[attr.id] = [String(selectedValue)];
            variantValues[attr.display_name || attr.name] = [
              String(selectedValue),
            ];
          }
        }
      });

      console.log('prepareAttributeValues - values:', values);
      console.log('prepareAttributeValues - variantValues:', variantValues);

      setAttributeValues(values);
      setSelectedVariantValues(variantValues);
    };

    prepareAttributeValues();
  }, [stockAttributes, selectedAttributes]);

  const toggleAttributeValue = (attributeName: string, value: string) => {
    setSelectedVariantValues((prev) => {
      const current = prev[attributeName] || [];
      if (current.includes(value)) {
        return {
          ...prev,
          [attributeName]: current.filter((v) => v !== value),
        };
      } else {
        return {
          ...prev,
          [attributeName]: [...current, value],
        };
      }
    });
  };

  const generateVariants = () => {
    const variants: any[] = [];
    const attributeNames = Object.keys(selectedVariantValues);

    if (attributeNames.length === 0) {
      return;
    }

    // Добавляем уникальный суффикс для избежания дубликатов SKU
    const timestamp = Date.now();
    const randomSuffix = Math.random()
      .toString(36)
      .substring(2, 6)
      .toUpperCase();

    // Генерируем все комбинации
    const generateCombinations = (
      index: number,
      current: Record<string, string>
    ) => {
      if (index === attributeNames.length) {
        const variantKey = Object.values(current).join('-');
        const sku = `${variantKey}-${timestamp}-${randomSuffix}`
          .toUpperCase()
          .replace(/\s+/g, '-');

        variants.push({
          sku: sku,
          price: basePrice,
          stock_quantity: stockSettings.useIndividualQuantities
            ? individualQuantities[variantKey] || stockSettings.defaultQuantity
            : stockSettings.defaultQuantity,
          variant_attributes: { ...current },
          is_default: variants.length === 0,
        });
        return;
      }

      const attrName = attributeNames[index];
      const values = selectedVariantValues[attrName] || [];

      for (const value of values) {
        current[attrName] = value;
        generateCombinations(index + 1, current);
      }
    };

    generateCombinations(0, {});
    setGeneratedVariants(variants);
    setShowStockManager(true);
  };

  const handleVariantsConfirm = (finalVariants: any[]) => {
    onGenerate(finalVariants);
  };

  const getTotalVariants = () => {
    return Object.values(selectedVariantValues).reduce(
      (total, values) => total * Math.max(values.length, 1),
      1
    );
  };

  if (_loading) {
    return (
      <div className="text-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (stockAttributes.length === 0) {
    return (
      <div className="text-center py-8">
        <div className="alert alert-warning mb-4">
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
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z"
            />
          </svg>
          <div>
            <h3 className="font-bold">{t('noVariantAttributes')}</h3>
            <div className="text-xs">{t('setupAttributesFirst')}</div>
          </div>
        </div>

        <div className="space-y-2 text-sm">
          <p className="text-base-content/70">Для создания вариантов нужно:</p>
          <ul className="list-disc list-inside text-left text-base-content/60 space-y-1">
            <li>Вернуться на шаг &quot;Атрибуты&quot;</li>
            <li>Выбрать атрибуты для вариантов (цвет, размер)</li>
            <li>Выбрать multiselect значения (несколько цветов/размеров)</li>
          </ul>

          <div className="mt-4">
            <p className="text-xs text-base-content/50">
              Доступно: {categoryAttributes.length} атрибутов, выбрано:{' '}
              {Object.keys(selectedAttributes).length}
            </p>
          </div>
        </div>
      </div>
    );
  }

  // Если показываем stock manager
  if (showStockManager) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-medium">
            {t('configureStockAndPrices')}
          </h3>
          <button
            onClick={() => setShowStockManager(false)}
            className="btn btn-outline btn-sm"
          >
            ← {t('backToSettings')}
          </button>
        </div>

        <VariantStockManager
          variants={generatedVariants}
          onVariantsChange={setGeneratedVariants}
        />

        <div className="flex justify-end space-x-2">
          <button
            onClick={() => setShowStockManager(false)}
            className="btn btn-outline"
          >
            {t('cancel')}
          </button>
          <button
            onClick={() => handleVariantsConfirm(generatedVariants)}
            className="btn btn-primary"
          >
            {t('confirmVariants')}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Выбор атрибутов */}
      <div className="space-y-4">
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
          <span>{t('attributesFromPreviousStep')}</span>
        </div>

        {stockAttributes.map((attr) => {
          const values = attributeValues[attr.id] || [];
          const selected =
            selectedVariantValues[attr.display_name || attr.name] || [];
          const attributeName = attr.display_name || attr.name;

          return (
            <div key={attr.id} className="space-y-2">
              <h4 className="font-medium text-base-content flex items-center gap-2">
                {attributeName}
                <div className="badge badge-primary badge-sm">
                  Значений: {values.length}
                </div>
                {values.length > 0 && (
                  <div className="badge badge-info badge-sm">
                    Выбрано: {selected.length}
                  </div>
                )}
              </h4>
              <div className="flex flex-wrap gap-2">
                {values.map((value: string) => (
                  <label
                    key={value}
                    className={`
                      btn btn-sm
                      ${
                        selected.includes(value) ? 'btn-primary' : 'btn-outline'
                      }
                    `}
                  >
                    <input
                      type="checkbox"
                      className="hidden"
                      checked={selected.includes(value)}
                      onChange={() =>
                        toggleAttributeValue(attributeName, value)
                      }
                    />
                    {(attr.attribute_type === 'color' ||
                      attr.name?.toLowerCase().includes('color') ||
                      attr.name?.toLowerCase().includes('цвет')) && (
                      <span
                        className="w-4 h-4 rounded-full mr-2 border border-base-300"
                        style={{ backgroundColor: value.toLowerCase() }}
                      />
                    )}
                    {value}
                  </label>
                ))}
              </div>
            </div>
          );
        })}
      </div>

      {/* Настройки остатков */}
      <div className="divider"></div>
      <div className="space-y-4">
        <h4 className="font-medium text-base-content">{t('stockSettings')}</h4>

        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('defaultStockQuantity')}</span>
          </label>
          <input
            type="number"
            min="0"
            value={stockSettings.defaultQuantity}
            onChange={(e) =>
              setStockSettings((prev) => ({
                ...prev,
                defaultQuantity: parseInt(e.target.value) || 0,
              }))
            }
            className="input input-bordered w-32"
          />
        </div>

        <div className="form-control">
          <label className="label cursor-pointer justify-start">
            <input
              type="checkbox"
              className="checkbox checkbox-primary mr-2"
              checked={stockSettings.useIndividualQuantities}
              onChange={(e) =>
                setStockSettings((prev) => ({
                  ...prev,
                  useIndividualQuantities: e.target.checked,
                }))
              }
            />
            <span className="label-text">{t('useIndividualQuantities')}</span>
          </label>
        </div>
      </div>

      {/* Информация о генерации */}
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
        <span>{t('willGenerateVariants', { count: getTotalVariants() })}</span>
      </div>

      {/* Кнопки */}
      <div className="flex justify-end space-x-2">
        <button onClick={onCancel} className="btn btn-outline">
          {t('cancel')}
        </button>
        <button
          onClick={generateVariants}
          disabled={getTotalVariants() === 0}
          className="btn btn-primary"
        >
          {t('generateVariants')}
        </button>
      </div>
    </div>
  );
}
