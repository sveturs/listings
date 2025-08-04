'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';

interface RangeAttributeProps {
  attribute: {
    id: number;
    name: string;
    display_name: string;
    unit?: string;
    min_value?: number;
    max_value?: number;
    translations?: Record<string, string>;
    is_required?: boolean;
  };
  value?: { min?: number; max?: number } | string;
  onChange: (value: { min?: number; max?: number }) => void;
  error?: string;
  locale?: string;
}

export default function RangeAttribute({
  attribute,
  value,
  onChange,
  error,
  locale = 'ru',
}: RangeAttributeProps) {
  const t = useTranslations('marketplace');
  const [minValue, setMinValue] = useState<string>('');
  const [maxValue, setMaxValue] = useState<string>('');

  useEffect(() => {
    if (value) {
      if (typeof value === 'string') {
        try {
          const parsed = JSON.parse(value);
          setMinValue(parsed.min?.toString() || '');
          setMaxValue(parsed.max?.toString() || '');
        } catch {
          // Handle invalid JSON
        }
      } else if (typeof value === 'object') {
        setMinValue(value.min?.toString() || '');
        setMaxValue(value.max?.toString() || '');
      }
    }
  }, [value]);

  const handleMinChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setMinValue(val);

    const numVal = val ? Number(val) : undefined;
    const maxNum = maxValue ? Number(maxValue) : undefined;

    onChange({
      min: numVal,
      max: maxNum,
    });
  };

  const handleMaxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setMaxValue(val);

    const minNum = minValue ? Number(minValue) : undefined;
    const numVal = val ? Number(val) : undefined;

    onChange({
      min: minNum,
      max: numVal,
    });
  };

  const getDisplayName = () => {
    if (attribute.translations?.[locale]) {
      return attribute.translations[locale];
    }
    return attribute.display_name;
  };

  const hasValidationError = () => {
    if (minValue && maxValue) {
      const min = Number(minValue);
      const max = Number(maxValue);
      return min > max;
    }
    return false;
  };

  return (
    <div className="form-control">
      <label className="label">
        <span className="label-text">
          {getDisplayName()}
          {attribute.is_required && <span className="text-error"> *</span>}
        </span>
      </label>

      <div className="flex gap-2 items-center">
        <div className="form-control flex-1">
          <label className="label py-0">
            <span className="label-text text-xs">{t('create.from')}</span>
          </label>
          <input
            type="number"
            value={minValue}
            onChange={handleMinChange}
            className={`input input-bordered ${error || hasValidationError() ? 'input-error' : ''}`}
            placeholder={t('create.min')}
            min={attribute.min_value}
            max={attribute.max_value}
          />
        </div>

        <span className="mt-6">—</span>

        <div className="form-control flex-1">
          <label className="label py-0">
            <span className="label-text text-xs">{t('create.to')}</span>
          </label>
          <input
            type="number"
            value={maxValue}
            onChange={handleMaxChange}
            className={`input input-bordered ${error || hasValidationError() ? 'input-error' : ''}`}
            placeholder={t('create.max')}
            min={attribute.min_value}
            max={attribute.max_value}
          />
        </div>

        {attribute.unit && (
          <span className="mt-6 text-base-content/70 min-w-[3rem]">
            {attribute.unit}
          </span>
        )}
      </div>

      {hasValidationError() && (
        <label className="label">
          <span className="label-text-alt text-error">
            {t('create.minGreaterThanMax')}
          </span>
        </label>
      )}

      {error && !hasValidationError() && (
        <label className="label">
          <span className="label-text-alt text-error">{error}</span>
        </label>
      )}

      {attribute.min_value !== undefined ||
      attribute.max_value !== undefined ? (
        <label className="label">
          <span className="label-text-alt text-base-content/50">
            {attribute.min_value !== undefined &&
            attribute.max_value !== undefined
              ? `${t('create.allowedRange')}: ${attribute.min_value} - ${attribute.max_value} ${attribute.unit || ''}`
              : attribute.min_value !== undefined
                ? `${t('create.minAllowed')}: ${attribute.min_value} ${attribute.unit || ''}`
                : `${t('create.maxAllowed')}: ${attribute.max_value} ${attribute.unit || ''}`}
          </span>
        </label>
      ) : null}
    </div>
  );
}
