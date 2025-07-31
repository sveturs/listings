'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import VariantStockManager from './VariantStockManager';
import { getVariantAttributesForCategory } from '@/lib/variant-attributes-config';

interface ProductVariantAttribute {
  id: number;
  name: string;
  display_name: string;
  type: string;
  affects_stock: boolean;
}

interface EnhancedVariantGeneratorProps {
  categorySlug: string;
  categoryAttributes: any[];
  basePrice: number;
  onGenerate: (variants: any[]) => void;
  onCancel: () => void;
}

export default function EnhancedVariantGenerator({
  categorySlug,
  basePrice,
  onGenerate,
  onCancel,
}: EnhancedVariantGeneratorProps) {
  const t = useTranslations('storefronts.products');
  const [loading, setLoading] = useState(true);
  const [variantAttributes, setVariantAttributes] = useState<
    ProductVariantAttribute[]
  >([]);
  const [selectedValues, setSelectedValues] = useState<
    Record<string, string[]>
  >({});
  const [generatedVariants, setGeneratedVariants] = useState<any[]>([]);
  const [showStockManager, setShowStockManager] = useState(false);

  // Загружаем вариативные атрибуты
  useEffect(() => {
    const loadVariantAttributes = async () => {
      try {
        // Получаем список атрибутов для категории
        const attributeNames = getVariantAttributesForCategory(categorySlug);

        // Загружаем информацию о вариативных атрибутах
        const response = await fetch('/api/v1/product-variant-attributes');
        if (response.ok) {
          const data = await response.json();
          const allAttributes: ProductVariantAttribute[] =
            data.data || data || [];

          // Фильтруем только те, которые подходят для этой категории
          const categoryVariants = allAttributes.filter((attr) =>
            attributeNames.includes(attr.name)
          );

          setVariantAttributes(categoryVariants);

          // Инициализируем выбранные значения
          const initialValues: Record<string, string[]> = {};
          categoryVariants.forEach((attr) => {
            initialValues[attr.name] = [];
          });
          setSelectedValues(initialValues);
        }
      } catch (error) {
        console.error('Failed to load variant attributes:', error);
      } finally {
        setLoading(false);
      }
    };

    loadVariantAttributes();
  }, [categorySlug]);

  const toggleValue = (attributeName: string, value: string) => {
    setSelectedValues((prev) => {
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
    const attributesWithValues = variantAttributes.filter(
      (attr) => selectedValues[attr.name]?.length > 0
    );

    if (attributesWithValues.length === 0) return;

    // Уникальный суффикс для SKU
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
      if (index === attributesWithValues.length) {
        const variantKey = Object.values(current).join('-');
        const sku = `${variantKey}-${timestamp}-${randomSuffix}`
          .toUpperCase()
          .replace(/\s+/g, '-');

        variants.push({
          sku: sku,
          price: basePrice,
          stock_quantity: 10,
          variant_attributes: { ...current },
          is_default: variants.length === 0,
        });
        return;
      }

      const attr = attributesWithValues[index];
      const values = selectedValues[attr.name] || [];

      values.forEach((value) => {
        current[attr.display_name || attr.name] = value;
        generateCombinations(index + 1, { ...current });
      });
    };

    generateCombinations(0, {});
    setGeneratedVariants(variants);
    setShowStockManager(true);
  };

  const handleVariantsConfirm = (finalVariants: any[]) => {
    onGenerate(finalVariants);
  };

  const getTotalVariants = () => {
    return Object.values(selectedValues).reduce(
      (total, values) => total * Math.max(values.length, 1),
      1
    );
  };

  // Функция для получения возможных значений атрибута
  const getAttributeValues = (attribute: ProductVariantAttribute): string[] => {
    // Определяем возможные значения в зависимости от типа атрибута
    switch (attribute.name) {
      case 'color':
        return [
          'Black',
          'White',
          'Red',
          'Blue',
          'Green',
          'Yellow',
          'Purple',
          'Gray',
          'Pink',
          'Orange',
        ];
      case 'size':
        return ['XS', 'S', 'M', 'L', 'XL', 'XXL'];
      case 'memory':
        return ['4GB', '6GB', '8GB', '12GB', '16GB'];
      case 'storage':
        return ['64GB', '128GB', '256GB', '512GB', '1TB'];
      case 'material':
        return ['Cotton', 'Polyester', 'Wool', 'Leather', 'Synthetic'];
      case 'capacity':
        return ['0.5L', '1L', '1.5L', '2L', '3L', '5L'];
      case 'power':
        return ['500W', '750W', '1000W', '1500W', '2000W'];
      case 'connectivity':
        return ['USB-C', 'Lightning', 'Bluetooth', 'WiFi', '3.5mm Jack'];
      case 'style':
        return ['Classic', 'Modern', 'Vintage', 'Casual', 'Formal'];
      case 'pattern':
        return ['Solid', 'Striped', 'Checkered', 'Floral', 'Abstract'];
      case 'weight':
        return ['0.5kg', '1kg', '2kg', '5kg', '10kg'];
      case 'bundle':
        return ['Basic', 'Standard', 'Premium', 'Deluxe'];
      default:
        return [];
    }
  };

  if (loading) {
    return (
      <div className="text-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (variantAttributes.length === 0) {
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
            <h3 className="font-bold">{t('noVariantAttributesForCategory')}</h3>
            <div className="text-xs">{t('categoryDoesNotSupportVariants')}</div>
          </div>
        </div>
      </div>
    );
  }

  // Если показываем менеджер остатков
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
      {/* Информация о доступных атрибутах */}
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
        <span>{t('selectValuesForVariants')}</span>
      </div>

      {/* Выбор значений для каждого атрибута */}
      <div className="space-y-4">
        {variantAttributes.map((attr) => {
          const values = getAttributeValues(attr);
          const selected = selectedValues[attr.name] || [];

          return (
            <div key={attr.id} className="space-y-2">
              <h4 className="font-medium text-base-content flex items-center gap-2">
                {attr.display_name}
                {attr.affects_stock && (
                  <div className="badge badge-info badge-sm">
                    {t('affectsStock')}
                  </div>
                )}
              </h4>
              <div className="flex flex-wrap gap-2">
                {values.map((value) => (
                  <label
                    key={value}
                    className={`
                      btn btn-sm
                      ${selected.includes(value) ? 'btn-primary' : 'btn-outline'}
                    `}
                  >
                    <input
                      type="checkbox"
                      className="hidden"
                      checked={selected.includes(value)}
                      onChange={() => toggleValue(attr.name, value)}
                    />
                    {attr.name === 'color' && (
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
