'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import {
  getLocalizedColors,
  getColorHex,
  getLocalizedColorName,
} from '@/utils/colorLocalization';
import type { components } from '@/types/generated/api';

type UnifiedAttributeValue =
  components['schemas']['models.UnifiedAttributeValue'];

interface ColorPickerAttributeProps {
  attributeId: number;
  value?: UnifiedAttributeValue;
  onChange: (value: UnifiedAttributeValue) => void;
  className?: string;
}

export function ColorPickerAttribute({
  attributeId,
  value,
  onChange,
  className = '',
}: ColorPickerAttributeProps) {
  const t = useTranslations('common');
  const locale = useLocale();

  const [selectedColor, setSelectedColor] = useState<string>('');
  const [customColor, setCustomColor] = useState<string>('#000000');
  const [showCustomPicker, setShowCustomPicker] = useState(false);

  // Получаем локализованные цвета
  const localizedColors = getLocalizedColors(locale);

  // Инициализация из value
  useEffect(() => {
    if (value?.text_value) {
      const colorValue = value.text_value;
      setSelectedColor(colorValue);

      // Если это не предустановленный цвет, показываем пользовательский выбор
      const isPresetColor = localizedColors.some(
        (color) =>
          color.hex.toLowerCase() === colorValue.toLowerCase() ||
          color.name === colorValue ||
          color.originalName === colorValue
      );

      if (!isPresetColor && colorValue.startsWith('#')) {
        setCustomColor(colorValue);
        setShowCustomPicker(true);
      }
    }
  }, [value, localizedColors]);

  // Выбор предустановленного цвета
  const handlePresetColorSelect = (color: {
    hex: string;
    name: string;
    originalName: string;
  }) => {
    setSelectedColor(color.name);
    setShowCustomPicker(false);

    const attributeValue: UnifiedAttributeValue = {
      attribute_id: attributeId,
      text_value: color.name, // Сохраняем локализованное название
    };

    onChange(attributeValue);
  };

  // Выбор пользовательского цвета
  const handleCustomColorChange = (hex: string) => {
    setCustomColor(hex);
    setSelectedColor(hex);

    const attributeValue: UnifiedAttributeValue = {
      attribute_id: attributeId,
      text_value: hex,
    };

    onChange(attributeValue);
  };

  // Получение локализованного названия для отображения
  const getDisplayColorName = (colorValue: string): string => {
    return getLocalizedColorName(colorValue, locale);
  };

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Выбранный цвет */}
      {selectedColor && (
        <div className="flex items-center gap-3 p-3 bg-base-200 rounded-lg">
          <div
            className="w-8 h-8 rounded-full border-2 border-base-300 shadow-sm"
            style={{ backgroundColor: getColorHex(selectedColor) }}
          />
          <div>
            <div className="font-medium">
              {getDisplayColorName(selectedColor)}
            </div>
            <div className="text-sm text-base-content/70">
              {getColorHex(selectedColor)}
            </div>
          </div>
          <button
            onClick={() => {
              setSelectedColor('');
              onChange({ attribute_id: attributeId, text_value: '' });
            }}
            className="btn btn-xs btn-ghost ml-auto"
            title={t('clear')}
          >
            ✕
          </button>
        </div>
      )}

      {/* Предустановленные цвета */}
      <div>
        <h4 className="text-sm font-medium mb-3 flex items-center gap-2">
          🎨 {t('colors.popular_colors')}
        </h4>
        <div className="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 gap-2">
          {localizedColors.slice(0, 16).map((color) => (
            <button
              key={color.hex}
              onClick={() => handlePresetColorSelect(color)}
              className={`group relative w-12 h-12 rounded-lg border-2 transition-all hover:scale-110 hover:shadow-lg ${
                selectedColor === color.name
                  ? 'border-primary ring-2 ring-primary/20'
                  : 'border-base-300'
              }`}
              style={{ backgroundColor: color.hex }}
              title={color.name}
            >
              {/* Отметка выбора */}
              {selectedColor === color.name && (
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="w-4 h-4 bg-white rounded-full border border-gray-300 flex items-center justify-center">
                    <div className="w-2 h-2 bg-primary rounded-full" />
                  </div>
                </div>
              )}

              {/* Всплывающая подсказка */}
              <div className="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-neutral text-neutral-content text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap z-10">
                {color.name}
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Пользовательский цвет */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h4 className="text-sm font-medium flex items-center gap-2">
            🎯 {t('colors.custom_color')}
          </h4>
          <button
            onClick={() => setShowCustomPicker(!showCustomPicker)}
            className="btn btn-xs btn-outline"
          >
            {showCustomPicker ? t('hide') : t('show')}
          </button>
        </div>

        {showCustomPicker && (
          <div className="p-4 border border-base-300 rounded-lg bg-base-100">
            <div className="flex items-center gap-3">
              <input
                type="color"
                value={customColor}
                onChange={(e) => handleCustomColorChange(e.target.value)}
                className="w-12 h-12 rounded border-2 border-base-300 cursor-pointer"
              />
              <div className="flex-1">
                <input
                  type="text"
                  value={customColor}
                  onChange={(e) => {
                    const value = e.target.value;
                    if (value.match(/^#[0-9A-Fa-f]{6}$/)) {
                      handleCustomColorChange(value);
                    } else {
                      setCustomColor(value);
                    }
                  }}
                  placeholder="#000000"
                  className="input input-sm input-bordered w-full font-mono"
                />
              </div>
              <button
                onClick={() => handleCustomColorChange(customColor)}
                disabled={!customColor.match(/^#[0-9A-Fa-f]{6}$/)}
                className="btn btn-sm btn-primary"
              >
                {t('apply')}
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Быстрая очистка */}
      {selectedColor && (
        <div className="flex justify-end">
          <button
            onClick={() => {
              setSelectedColor('');
              setCustomColor('#000000');
              setShowCustomPicker(false);
              onChange({ attribute_id: attributeId, text_value: '' });
            }}
            className="btn btn-sm btn-outline btn-error"
          >
            {t('colors.clear_selection')}
          </button>
        </div>
      )}
    </div>
  );
}
