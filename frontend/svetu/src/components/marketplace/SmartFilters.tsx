'use client';

import React, { useState, useEffect, lazy, Suspense } from 'react';
import { useTranslations } from 'next-intl';
import { useCategoryFilters } from '@/hooks/useCategoryFilters';
import type { components } from '@/types/generated/api';

// Lazy load CarFilters для автомобильных категорий
const CarFilters = lazy(() =>
  import('./CarFilters').then((module) => ({ default: module.CarFilters }))
);

type CategoryAttribute =
  components['schemas']['backend_internal_domain_models.CategoryAttribute'];

interface FilterValue {
  [attributeId: string]: any;
}

interface SmartFiltersProps {
  categoryId: number | null;
  onChange: (filters: FilterValue) => void;
  lang?: string;
  className?: string;
}

// Функция для проверки автомобильной категории
const isAutomotiveCategory = (categoryId: number | null): boolean => {
  if (!categoryId) return false;

  // Проверяем основные автомобильные категории
  const automotiveCategories = [
    1003, // automotive - основная категория автомобилей
    1301, // cars - легковые автомобили
    1303, // auto-parts - автозапчасти
  ];

  // Проверяем основные категории
  if (automotiveCategories.includes(categoryId)) {
    return true;
  }

  // Также проверяем диапазон специальных категорий (домашнее производство и т.д.)
  if (categoryId >= 10100 && categoryId <= 10199) {
    return true;
  }

  // Проверяем специальные подкатегории автомобилей
  if (categoryId >= 10170 && categoryId <= 10179) {
    return true;
  }

  return false;
};

export function SmartFilters({
  categoryId,
  onChange,
  lang = 'sr',
  className = '',
}: SmartFiltersProps) {
  const t = useTranslations('marketplace');
  const { attributes, loading, error } = useCategoryFilters(categoryId, {
    lang,
  });
  const [filterValues, setFilterValues] = useState<FilterValue>({});

  // Сброс фильтров при смене категории
  useEffect(() => {
    // Сохраняем предыдущее состояние перед сбросом
    const hadFilters = Object.keys(filterValues).length > 0;
    setFilterValues({});
    // Вызываем onChange только если были активные фильтры
    if (hadFilters) {
      onChange({});
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categoryId]); // onChange не включаем чтобы избежать бесконечного цикла

  const handleFilterChange = (attributeId: string, value: any) => {
    const newFilters = {
      ...filterValues,
      [attributeId]: value,
    };

    // Удаляем пустые значения
    if (
      value === '' ||
      value === null ||
      value === undefined ||
      (Array.isArray(value) && value.length === 0)
    ) {
      delete newFilters[attributeId];
    }

    setFilterValues(newFilters);
    onChange(newFilters);
  };

  const renderFilter = (attribute: CategoryAttribute) => {
    if (!attribute.id) return null;
    const value = filterValues[attribute.id] || '';

    switch (attribute.attribute_type) {
      case 'text':
        return (
          <input
            type="text"
            value={value}
            onChange={(e) =>
              handleFilterChange(attribute.id!.toString(), e.target.value)
            }
            placeholder={t('filters.enterValue')}
            className="input input-bordered input-sm w-full"
          />
        );

      case 'number':
        return (
          <input
            type="number"
            value={value}
            onChange={(e) =>
              handleFilterChange(
                attribute.id!.toString(),
                parseFloat(e.target.value) || ''
              )
            }
            placeholder={t('filters.enterNumber')}
            className="input input-bordered input-sm w-full"
          />
        );

      case 'range':
        const rangeValue = value || { min: '', max: '' };
        return (
          <div className="flex gap-2">
            <input
              type="number"
              value={rangeValue.min || ''}
              onChange={(e) =>
                handleFilterChange(attribute.id!.toString(), {
                  ...rangeValue,
                  min: parseFloat(e.target.value) || '',
                })
              }
              placeholder={t('filters.min')}
              className="input input-bordered input-sm flex-1"
            />
            <span className="self-center">-</span>
            <input
              type="number"
              value={rangeValue.max || ''}
              onChange={(e) =>
                handleFilterChange(attribute.id!.toString(), {
                  ...rangeValue,
                  max: parseFloat(e.target.value) || '',
                })
              }
              placeholder={t('filters.max')}
              className="input input-bordered input-sm flex-1"
            />
          </div>
        );

      case 'select':
        let options: Array<{ value: string; label: string }> = [];

        // Проверяем translated_options
        if (attribute.translated_options) {
          try {
            options = JSON.parse(attribute.translated_options.toString());
          } catch {}
        }

        // Если нет translated_options, используем options как массив чисел
        if (
          options.length === 0 &&
          attribute.options &&
          Array.isArray(attribute.options)
        ) {
          // options это массив чисел, используем их как значения
          options = attribute.options.map((val) => ({
            value: val.toString(),
            label: val.toString(),
          }));
        }

        return (
          <select
            value={value}
            onChange={(e) =>
              handleFilterChange(attribute.id!.toString(), e.target.value)
            }
            className="select select-bordered select-sm w-full"
          >
            <option value="">{t('filters.selectOption')}</option>
            {options.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        );

      case 'multiselect':
        let multiOptions: Array<{ value: string; label: string }> = [];

        // Аналогично select
        if (attribute.translated_options) {
          try {
            multiOptions = JSON.parse(attribute.translated_options.toString());
          } catch {}
        }

        if (
          multiOptions.length === 0 &&
          attribute.options &&
          Array.isArray(attribute.options)
        ) {
          // options это массив чисел, используем их как значения
          multiOptions = attribute.options.map((val) => ({
            value: val.toString(),
            label: val.toString(),
          }));
        }

        const selectedValues = Array.isArray(value) ? value : [];

        return (
          <div className="space-y-2">
            {multiOptions.map((option) => (
              <label
                key={option.value}
                className="flex items-center gap-2 cursor-pointer"
              >
                <input
                  type="checkbox"
                  checked={selectedValues.includes(option.value)}
                  onChange={(e) => {
                    if (e.target.checked) {
                      handleFilterChange(attribute.id!.toString(), [
                        ...selectedValues,
                        option.value,
                      ]);
                    } else {
                      handleFilterChange(
                        attribute.id!.toString(),
                        selectedValues.filter((v) => v !== option.value)
                      );
                    }
                  }}
                  className="checkbox checkbox-sm"
                />
                <span className="text-sm">{option.label}</span>
              </label>
            ))}
          </div>
        );

      case 'boolean':
        return (
          <select
            value={value === '' ? '' : value.toString()}
            onChange={(e) => {
              const val = e.target.value;
              handleFilterChange(
                attribute.id!.toString(),
                val === '' ? '' : val === 'true'
              );
            }}
            className="select select-bordered select-sm w-full"
          >
            <option value="">{t('filters.any')}</option>
            <option value="true">{t('filters.yes')}</option>
            <option value="false">{t('filters.no')}</option>
          </select>
        );

      default:
        return null;
    }
  };

  if (!categoryId) {
    return (
      <div className={`p-4 text-center text-base-content/60 ${className}`}>
        <p>{t('filters.selectCategory')}</p>
      </div>
    );
  }

  // Для автомобильных категорий используем специальный компонент CarFilters
  if (isAutomotiveCategory(categoryId)) {
    return (
      <Suspense
        fallback={
          <div className={`p-4 ${className}`}>
            <div className="flex justify-center">
              <span className="loading loading-spinner loading-sm"></span>
            </div>
          </div>
        }
      >
        <CarFilters onFiltersChange={onChange} className={className} />
      </Suspense>
    );
  }

  if (loading) {
    return (
      <div className={`p-4 ${className}`}>
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-sm"></span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`p-4 ${className}`}>
        <div className="alert alert-error">
          <span>{t('filters.loadError')}</span>
        </div>
      </div>
    );
  }

  if (attributes.length === 0) {
    return (
      <div className={`p-4 text-center text-base-content/60 ${className}`}>
        <p>{t('filters.noFilters')}</p>
      </div>
    );
  }

  // Группируем атрибуты по типу для лучшей организации
  const filterableAttributes = attributes.filter((attr) => attr.is_filterable);

  return (
    <div className={`space-y-4 ${className}`}>
      {filterableAttributes.map((attribute) => {
        if (!attribute.id) return null;
        return (
          <div key={attribute.id} className="form-control">
            <label className="label">
              <span className="label-text font-medium flex items-center gap-2">
                {attribute.icon && <span>{attribute.icon}</span>}
                {attribute.display_name}
                {attribute.is_required && <span className="text-error">*</span>}
              </span>
            </label>
            {renderFilter(attribute)}
          </div>
        );
      })}

      {filterValues && Object.keys(filterValues).length > 0 && (
        <button
          onClick={() => {
            setFilterValues({});
            onChange({});
          }}
          className="btn btn-sm btn-ghost w-full"
        >
          {t('filters.clearAll')}
        </button>
      )}
    </div>
  );
}
