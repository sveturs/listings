'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';

type UnifiedAttributeValue =
  components['schemas']['models.UnifiedAttributeValue'];

interface SizePickerAttributeProps {
  attributeId: number;
  value?: UnifiedAttributeValue;
  onChange: (value: UnifiedAttributeValue) => void;
  className?: string;
  sizeType?: 'clothing' | 'shoes' | 'generic';
}

// Размеры одежды
const CLOTHING_SIZES = [
  { value: 'XXS', label: 'XXS', numeric: 38 },
  { value: 'XS', label: 'XS', numeric: 40 },
  { value: 'S', label: 'S', numeric: 42 },
  { value: 'M', label: 'M', numeric: 44 },
  { value: 'L', label: 'L', numeric: 46 },
  { value: 'XL', label: 'XL', numeric: 48 },
  { value: 'XXL', label: 'XXL', numeric: 50 },
  { value: 'XXXL', label: 'XXXL', numeric: 52 },
];

// Размеры обуви (европейские)
const SHOE_SIZES = [
  { value: '35', label: '35', category: 'EU' },
  { value: '36', label: '36', category: 'EU' },
  { value: '37', label: '37', category: 'EU' },
  { value: '38', label: '38', category: 'EU' },
  { value: '39', label: '39', category: 'EU' },
  { value: '40', label: '40', category: 'EU' },
  { value: '41', label: '41', category: 'EU' },
  { value: '42', label: '42', category: 'EU' },
  { value: '43', label: '43', category: 'EU' },
  { value: '44', label: '44', category: 'EU' },
  { value: '45', label: '45', category: 'EU' },
  { value: '46', label: '46', category: 'EU' },
  { value: '47', label: '47', category: 'EU' },
  { value: '48', label: '48', category: 'EU' },
];

// Общие размеры
const GENERIC_SIZES = [
  { value: 'XS', label: 'XS - Очень маленький', description: 'Extra Small' },
  { value: 'S', label: 'S - Маленький', description: 'Small' },
  { value: 'M', label: 'M - Средний', description: 'Medium' },
  { value: 'L', label: 'L - Большой', description: 'Large' },
  { value: 'XL', label: 'XL - Очень большой', description: 'Extra Large' },
  { value: 'XXL', label: 'XXL - Огромный', description: 'Extra Extra Large' },
];

export function SizePickerAttribute({
  attributeId,
  value,
  onChange,
  className = '',
  sizeType = 'generic',
}: SizePickerAttributeProps) {
  const t = useTranslations('common');

  const [selectedSize, setSelectedSize] = useState<string>('');
  const [customSize, setCustomSize] = useState<string>('');
  const [showCustomInput, setShowCustomInput] = useState(false);

  // Получение подходящих размеров
  const getSizeOptions = () => {
    switch (sizeType) {
      case 'clothing':
        return CLOTHING_SIZES;
      case 'shoes':
        return SHOE_SIZES;
      default:
        return GENERIC_SIZES;
    }
  };

  const sizes = getSizeOptions();

  // Инициализация из value
  useEffect(() => {
    if (value?.text_value) {
      setSelectedSize(value.text_value);

      // Проверяем, есть ли размер в предустановленных
      const isPresetSize = sizes.some(
        (size) => size.value === value.text_value
      );
      if (!isPresetSize) {
        setCustomSize(value.text_value);
        setShowCustomInput(true);
      }
    }
  }, [value, sizes]);

  // Выбор предустановленного размера
  const handleSizeSelect = (size: (typeof sizes)[0]) => {
    setSelectedSize(size.value);
    setShowCustomInput(false);

    const attributeValue: UnifiedAttributeValue = {
      attribute_id: attributeId,
      text_value: size.value,
    };

    onChange(attributeValue);
  };

  // Ввод пользовательского размера
  const handleCustomSizeSubmit = () => {
    if (!customSize.trim()) return;

    setSelectedSize(customSize.trim());

    const attributeValue: UnifiedAttributeValue = {
      attribute_id: attributeId,
      text_value: customSize.trim(),
    };

    onChange(attributeValue);
  };

  // Очистка выбора
  const clearSelection = () => {
    setSelectedSize('');
    setCustomSize('');
    setShowCustomInput(false);
    onChange({ attribute_id: attributeId, text_value: '' });
  };

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Выбранный размер */}
      {selectedSize && (
        <div className="flex items-center gap-3 p-3 bg-base-200 rounded-lg">
          <div className="flex items-center justify-center w-12 h-12 bg-primary text-primary-content rounded-lg font-bold text-lg">
            {selectedSize}
          </div>
          <div>
            <div className="font-medium">
              {t('selected_size')}: {selectedSize}
            </div>
            {sizeType === 'clothing' && (
              <div className="text-sm text-base-content/70">
                {CLOTHING_SIZES.find((s) => s.value === selectedSize)?.numeric}
              </div>
            )}
          </div>
          <button
            onClick={clearSelection}
            className="btn btn-xs btn-ghost ml-auto"
            title={t('clear')}
          >
            ✕
          </button>
        </div>
      )}

      {/* Стандартные размеры */}
      <div>
        <h4 className="text-sm font-medium mb-3 flex items-center gap-2">
          📏 {t('sizes.standard_sizes')}
        </h4>

        <div className="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 gap-2">
          {sizes.map((size) => (
            <button
              key={size.value}
              onClick={() => handleSizeSelect(size)}
              className={`btn btn-sm h-12 transition-all ${
                selectedSize === size.value
                  ? 'btn-primary'
                  : 'btn-outline hover:btn-primary'
              }`}
              title={size.label}
            >
              <div className="text-center">
                <div className="font-bold">{size.value}</div>
                {sizeType === 'clothing' &&
                  'numeric' in size &&
                  size.numeric && (
                    <div className="text-xs opacity-70">{size.numeric}</div>
                  )}
              </div>
            </button>
          ))}
        </div>

        {/* Дополнительная информация для размеров одежды */}
        {sizeType === 'clothing' && (
          <div className="mt-3 text-sm text-base-content/70">
            💡 {t('sizes.clothing_hint')}
          </div>
        )}

        {/* Дополнительная информация для размеров обуви */}
        {sizeType === 'shoes' && (
          <div className="mt-3 text-sm text-base-content/70">
            👟 {t('sizes.shoe_hint')}
          </div>
        )}
      </div>

      {/* Пользовательский размер */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h4 className="text-sm font-medium flex items-center gap-2">
            ✏️ {t('sizes.custom_size')}
          </h4>
          <button
            onClick={() => setShowCustomInput(!showCustomInput)}
            className="btn btn-xs btn-outline"
          >
            {showCustomInput ? t('hide') : t('show')}
          </button>
        </div>

        {showCustomInput && (
          <div className="p-4 border border-base-300 rounded-lg bg-base-100">
            <div className="flex items-center gap-3">
              <input
                type="text"
                value={customSize}
                onChange={(e) => setCustomSize(e.target.value)}
                placeholder={t('sizes.custom_placeholder')}
                className="input input-sm input-bordered flex-1"
                onKeyPress={(e) =>
                  e.key === 'Enter' && handleCustomSizeSubmit()
                }
              />
              <button
                onClick={handleCustomSizeSubmit}
                disabled={!customSize.trim()}
                className="btn btn-sm btn-primary"
              >
                {t('apply')}
              </button>
            </div>
            <div className="text-xs text-base-content/60 mt-2">
              {t('sizes.custom_examples')}
            </div>
          </div>
        )}
      </div>

      {/* Быстрая очистка */}
      {selectedSize && (
        <div className="flex justify-end">
          <button
            onClick={clearSelection}
            className="btn btn-sm btn-outline btn-error"
          >
            {t('colors.clear_selection')}
          </button>
        </div>
      )}
    </div>
  );
}
