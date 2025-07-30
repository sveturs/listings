'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import VariantStockManager from './VariantStockManager';

interface SimplifiedVariantGeneratorProps {
  selectedAttributes: Record<number, any>;
  categoryAttributes: any[];
  basePrice: number;
  onGenerate: (variants: any[]) => void;
  onCancel: () => void;
}

export default function SimplifiedVariantGenerator({
  selectedAttributes,
  categoryAttributes,
  basePrice,
  onGenerate,
  onCancel,
}: SimplifiedVariantGeneratorProps) {
  const t = useTranslations('storefronts.products');
  const [generatedVariants, setGeneratedVariants] = useState<any[]>([]);
  const [showStockManager, setShowStockManager] = useState(false);

  // Фильтруем атрибуты которые могут быть использованы для вариантов
  const variantAttributes = React.useMemo(() => {
    console.log(
      'SimplifiedVariantGenerator - selectedAttributes:',
      selectedAttributes
    );
    console.log(
      'SimplifiedVariantGenerator - categoryAttributes:',
      categoryAttributes
    );

    return categoryAttributes.filter((attr) => {
      const value = selectedAttributes[attr.id];
      if (!value) return false;

      // Проверяем что это подходящий атрибут для вариантов
      const name = (attr.name || '').toLowerCase();
      const isVariantAttribute = [
        'color',
        'size',
        'цвет',
        'размер',
        'boja',
        'veličina',
      ].some((keyword) => name.includes(keyword));

      if (!isVariantAttribute) return false;

      // Для multiselect атрибутов
      if (attr.attribute_type === 'multiselect') {
        if (Array.isArray(value)) return value.length > 0;
        if (typeof value === 'string' && value.includes(',')) {
          return value.split(',').filter(Boolean).length > 0;
        }
      }

      return true;
    });
  }, [categoryAttributes, selectedAttributes]);

  console.log(
    'SimplifiedVariantGenerator - variantAttributes:',
    variantAttributes
  );

  const generateVariants = () => {
    const variants: any[] = [];
    const attributesWithValues: any[] = [];

    // Подготавливаем атрибуты и их значения
    variantAttributes.forEach((attr) => {
      const value = selectedAttributes[attr.id];
      let values: string[] = [];

      if (Array.isArray(value)) {
        values = value;
      } else if (typeof value === 'string') {
        if (value.includes(',')) {
          values = value.split(',').filter(Boolean);
        } else {
          values = [value];
        }
      }

      if (values.length > 0) {
        attributesWithValues.push({
          name: attr.display_name || attr.name,
          values: values,
        });
      }
    });

    if (attributesWithValues.length === 0) return;

    // Добавляем уникальный суффикс для избежания дубликатов SKU
    const timestamp = Date.now();
    const randomSuffix = Math.random()
      .toString(36)
      .substring(2, 6)
      .toUpperCase();

    // Генерируем комбинации
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
      attr.values.forEach((value: string) => {
        current[attr.name] = value;
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
    return variantAttributes.reduce((total, attr) => {
      const value = selectedAttributes[attr.id];
      let count = 0;

      if (Array.isArray(value)) {
        count = value.length;
      } else if (typeof value === 'string' && value.includes(',')) {
        count = value.split(',').filter(Boolean).length;
      } else if (value) {
        count = 1;
      }

      return total * Math.max(count, 1);
    }, 1);
  };

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
            <h3 className="font-bold">{t('noVariantAttributes')}</h3>
            <div className="text-xs">{t('setupAttributesFirst')}</div>
          </div>
        </div>

        <div className="space-y-2 text-sm">
          <p className="text-base-content/70">Для создания вариантов нужно:</p>
          <ul className="list-disc list-inside text-left text-base-content/60 space-y-1">
            <li>Вернуться на шаг &quot;Атрибуты&quot;</li>
            <li>Выбрать атрибуты Color/Цвет и Size/Размер</li>
            <li>Для multiselect атрибутов выбрать несколько значений</li>
          </ul>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
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
        <span>Используются атрибуты, выбранные на предыдущем шаге</span>
      </div>

      {/* Отображаем выбранные атрибуты */}
      <div className="space-y-4">
        {variantAttributes.map((attr) => {
          const value = selectedAttributes[attr.id];
          let values: string[] = [];

          if (Array.isArray(value)) {
            values = value;
          } else if (typeof value === 'string') {
            if (value.includes(',')) {
              values = value.split(',').filter(Boolean);
            } else {
              values = [value];
            }
          }

          return (
            <div key={attr.id} className="card bg-base-100 shadow-sm">
              <div className="card-body p-4">
                <h4 className="font-medium text-base-content">
                  {attr.display_name || attr.name}
                </h4>
                <div className="flex flex-wrap gap-2">
                  {values.map((val: string) => (
                    <div key={val} className="badge badge-primary">
                      {val}
                    </div>
                  ))}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Информация о генерации */}
      <div className="alert alert-success">
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
            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
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
