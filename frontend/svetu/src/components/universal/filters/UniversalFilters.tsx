'use client';

import { FC, useState, useEffect, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { FaFilter, FaTimes, FaChevronDown, FaChevronUp } from 'react-icons/fa';

// Типы для универсальных фильтров
export interface FilterOption {
  value: string | number;
  label: string;
  count?: number;
}

export interface FilterConfig {
  id: string;
  type:
    | 'select'
    | 'multiselect'
    | 'range'
    | 'boolean'
    | 'text'
    | 'date'
    | 'custom';
  label: string;
  placeholder?: string;
  options?: FilterOption[];
  min?: number;
  max?: number;
  step?: number;
  unit?: string;
  icon?: React.ComponentType<any>;
  defaultValue?: any;
  // Для зависимых фильтров
  dependsOn?: string;
  loadOptions?: (parentValue: any) => Promise<FilterOption[]>;
  // Для кастомных фильтров
  component?: React.ComponentType<any>;
  componentProps?: Record<string, any>;
  // Группировка
  group?: string;
  order?: number;
}

export interface FilterGroup {
  id: string;
  label: string;
  icon?: React.ComponentType<any>;
  collapsible?: boolean;
  defaultExpanded?: boolean;
  order?: number;
}

interface UniversalFiltersProps {
  category: string;
  filters: Record<string, any>;
  onFiltersChange: (filters: Record<string, any>) => void;
  config?: {
    showPriceRange?: boolean;
    showCondition?: boolean;
    showLocation?: boolean;
    showSort?: boolean;
    customFilters?: FilterConfig[];
    groups?: FilterGroup[];
  };
  layout?: 'vertical' | 'horizontal';
  className?: string;
}

// Предустановленные конфигурации для категорий
const CATEGORY_FILTERS: Record<string, FilterConfig[]> = {
  cars: [
    {
      id: 'make',
      type: 'select',
      label: 'Make',
      placeholder: 'Select make',
      order: 1,
      group: 'main',
    },
    {
      id: 'model',
      type: 'select',
      label: 'Model',
      placeholder: 'Select model',
      dependsOn: 'make',
      order: 2,
      group: 'main',
    },
    {
      id: 'year',
      type: 'range',
      label: 'Year',
      min: 1990,
      max: new Date().getFullYear() + 1,
      step: 1,
      order: 3,
      group: 'main',
    },
    {
      id: 'mileage',
      type: 'range',
      label: 'Mileage',
      min: 0,
      max: 500000,
      step: 5000,
      unit: 'km',
      order: 4,
      group: 'main',
    },
    {
      id: 'fuelType',
      type: 'multiselect',
      label: 'Fuel Type',
      options: [
        { value: 'gasoline', label: 'Gasoline' },
        { value: 'diesel', label: 'Diesel' },
        { value: 'electric', label: 'Electric' },
        { value: 'hybrid', label: 'Hybrid' },
      ],
      order: 5,
      group: 'technical',
    },
    {
      id: 'transmission',
      type: 'multiselect',
      label: 'Transmission',
      options: [
        { value: 'manual', label: 'Manual' },
        { value: 'automatic', label: 'Automatic' },
        { value: 'cvt', label: 'CVT' },
        { value: 'dct', label: 'DCT' },
      ],
      order: 6,
      group: 'technical',
    },
  ],
  real_estate: [
    {
      id: 'propertyType',
      type: 'multiselect',
      label: 'Property Type',
      options: [
        { value: 'apartment', label: 'Apartment' },
        { value: 'house', label: 'House' },
        { value: 'studio', label: 'Studio' },
        { value: 'villa', label: 'Villa' },
        { value: 'land', label: 'Land' },
      ],
      order: 1,
      group: 'main',
    },
    {
      id: 'rooms',
      type: 'multiselect',
      label: 'Rooms',
      options: [
        { value: '1', label: '1 Room' },
        { value: '2', label: '2 Rooms' },
        { value: '3', label: '3 Rooms' },
        { value: '4', label: '4 Rooms' },
        { value: '5+', label: '5+ Rooms' },
      ],
      order: 2,
      group: 'main',
    },
    {
      id: 'area',
      type: 'range',
      label: 'Area',
      min: 0,
      max: 500,
      step: 10,
      unit: 'm²',
      order: 3,
      group: 'main',
    },
    {
      id: 'floor',
      type: 'range',
      label: 'Floor',
      min: 1,
      max: 30,
      step: 1,
      order: 4,
      group: 'details',
    },
  ],
  electronics: [
    {
      id: 'brand',
      type: 'multiselect',
      label: 'Brand',
      order: 1,
      group: 'main',
    },
    {
      id: 'category',
      type: 'select',
      label: 'Category',
      options: [
        { value: 'phones', label: 'Phones' },
        { value: 'laptops', label: 'Laptops' },
        { value: 'tablets', label: 'Tablets' },
        { value: 'accessories', label: 'Accessories' },
      ],
      order: 2,
      group: 'main',
    },
  ],
};

// Группы фильтров по умолчанию
const DEFAULT_GROUPS: Record<string, FilterGroup[]> = {
  cars: [
    { id: 'main', label: 'Main Filters', order: 1, defaultExpanded: true },
    { id: 'technical', label: 'Technical', order: 2, collapsible: true },
    { id: 'features', label: 'Features', order: 3, collapsible: true },
  ],
  real_estate: [
    { id: 'main', label: 'Main Filters', order: 1, defaultExpanded: true },
    { id: 'details', label: 'Details', order: 2, collapsible: true },
    { id: 'amenities', label: 'Amenities', order: 3, collapsible: true },
  ],
  electronics: [
    { id: 'main', label: 'Main Filters', order: 1, defaultExpanded: true },
    { id: 'specs', label: 'Specifications', order: 2, collapsible: true },
  ],
};

const UniversalFilters: FC<UniversalFiltersProps> = ({
  category,
  filters,
  onFiltersChange,
  config = {},
  layout = 'vertical',
  className = '',
}) => {
  const t = useTranslations('filters');
  const [expandedGroups, setExpandedGroups] = useState<Record<string, boolean>>(
    {}
  );
  const [dependentOptions, setDependentOptions] = useState<
    Record<string, FilterOption[]>
  >({});

  // Получаем фильтры для категории
  const categoryFilters = useMemo(() => {
    const baseFilters = CATEGORY_FILTERS[category] || [];
    const customFilters = config.customFilters || [];
    return [...baseFilters, ...customFilters].sort(
      (a, b) => (a.order || 999) - (b.order || 999)
    );
  }, [category, config.customFilters]);

  // Получаем группы для категории
  const filterGroups = useMemo(() => {
    const baseGroups = DEFAULT_GROUPS[category] || [];
    const customGroups = config.groups || [];
    return [...baseGroups, ...customGroups].sort(
      (a, b) => (a.order || 999) - (b.order || 999)
    );
  }, [category, config.groups]);

  // Инициализация состояния групп
  useEffect(() => {
    const initial: Record<string, boolean> = {};
    filterGroups.forEach((group) => {
      initial[group.id] = group.defaultExpanded !== false;
    });
    setExpandedGroups(initial);
  }, [filterGroups]);

  // Загрузка зависимых опций
  useEffect(() => {
    categoryFilters.forEach(async (filter) => {
      if (filter.dependsOn && filter.loadOptions) {
        const parentValue = filters[filter.dependsOn];
        if (parentValue) {
          const options = await filter.loadOptions(parentValue);
          setDependentOptions((prev) => ({ ...prev, [filter.id]: options }));
        } else {
          setDependentOptions((prev) => ({ ...prev, [filter.id]: [] }));
        }
      }
    });
  }, [filters, categoryFilters]);

  const handleFilterChange = (filterId: string, value: any) => {
    const newFilters = { ...filters };

    if (value === undefined || value === null || value === '') {
      delete newFilters[filterId];
    } else {
      newFilters[filterId] = value;
    }

    // Сброс зависимых фильтров
    categoryFilters.forEach((filter) => {
      if (filter.dependsOn === filterId) {
        delete newFilters[filter.id];
      }
    });

    onFiltersChange(newFilters);
  };

  const toggleGroup = (groupId: string) => {
    setExpandedGroups((prev) => ({ ...prev, [groupId]: !prev[groupId] }));
  };

  const renderFilter = (filter: FilterConfig) => {
    const value = filters[filter.id];

    switch (filter.type) {
      case 'select':
        const selectOptions =
          dependentOptions[filter.id] || filter.options || [];
        return (
          <select
            className="select select-bordered w-full"
            value={value || ''}
            onChange={(e) => handleFilterChange(filter.id, e.target.value)}
            disabled={Boolean(filter.dependsOn && !filters[filter.dependsOn])}
          >
            <option value="">{filter.placeholder || t('selectOption')}</option>
            {selectOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
                {option.count !== undefined && ` (${option.count})`}
              </option>
            ))}
          </select>
        );

      case 'multiselect':
        const multiselectOptions =
          dependentOptions[filter.id] || filter.options || [];
        const selectedValues = value || [];
        return (
          <div className="space-y-2">
            {multiselectOptions.map((option) => (
              <label
                key={option.value}
                className="flex items-center gap-2 cursor-pointer"
              >
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm"
                  checked={selectedValues.includes(option.value)}
                  onChange={(e) => {
                    const newValues = e.target.checked
                      ? [...selectedValues, option.value]
                      : selectedValues.filter((v: any) => v !== option.value);
                    handleFilterChange(
                      filter.id,
                      newValues.length > 0 ? newValues : undefined
                    );
                  }}
                />
                <span className="text-sm">
                  {option.label}
                  {option.count !== undefined && (
                    <span className="text-base-content/60 ml-1">
                      ({option.count})
                    </span>
                  )}
                </span>
              </label>
            ))}
          </div>
        );

      case 'range':
        const rangeValue = value || [filter.min || 0, filter.max || 100];
        return (
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <input
                type="number"
                className="input input-bordered input-sm w-24"
                value={rangeValue[0]}
                min={filter.min}
                max={filter.max}
                onChange={(e) => {
                  const newMin = Number(e.target.value);
                  handleFilterChange(filter.id, [newMin, rangeValue[1]]);
                }}
              />
              <span className="text-base-content/60">—</span>
              <input
                type="number"
                className="input input-bordered input-sm w-24"
                value={rangeValue[1]}
                min={filter.min}
                max={filter.max}
                onChange={(e) => {
                  const newMax = Number(e.target.value);
                  handleFilterChange(filter.id, [rangeValue[0], newMax]);
                }}
              />
              {filter.unit && (
                <span className="text-sm text-base-content/60">
                  {filter.unit}
                </span>
              )}
            </div>
            <input
              type="range"
              className="range range-sm"
              min={filter.min}
              max={filter.max}
              step={filter.step}
              value={rangeValue[1]}
              onChange={(e) => {
                const newMax = Number(e.target.value);
                handleFilterChange(filter.id, [rangeValue[0], newMax]);
              }}
            />
          </div>
        );

      case 'boolean':
        return (
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              className="toggle toggle-sm"
              checked={value || false}
              onChange={(e) => handleFilterChange(filter.id, e.target.checked)}
            />
            <span className="text-sm">{filter.label}</span>
          </label>
        );

      case 'text':
        return (
          <input
            type="text"
            className="input input-bordered w-full"
            placeholder={filter.placeholder}
            value={value || ''}
            onChange={(e) => handleFilterChange(filter.id, e.target.value)}
          />
        );

      case 'date':
        return (
          <input
            type="date"
            className="input input-bordered w-full"
            value={value || ''}
            onChange={(e) => handleFilterChange(filter.id, e.target.value)}
          />
        );

      case 'custom':
        if (filter.component) {
          const CustomComponent = filter.component;
          return (
            <CustomComponent
              value={value}
              onChange={(newValue: any) =>
                handleFilterChange(filter.id, newValue)
              }
              {...(filter.componentProps || {})}
            />
          );
        }
        return null;

      default:
        return null;
    }
  };

  // Группировка фильтров
  const filtersByGroup = useMemo(() => {
    const grouped: Record<string, FilterConfig[]> = {};

    categoryFilters.forEach((filter) => {
      const groupId = filter.group || 'other';
      if (!grouped[groupId]) {
        grouped[groupId] = [];
      }
      grouped[groupId].push(filter);
    });

    return grouped;
  }, [categoryFilters]);

  // Подсчет активных фильтров
  const activeFilterCount = Object.keys(filters).filter((key) => {
    const value = filters[key];
    return (
      value !== undefined &&
      value !== null &&
      value !== '' &&
      (!Array.isArray(value) || value.length > 0)
    );
  }).length;

  if (layout === 'horizontal') {
    return (
      <div
        className={`flex flex-wrap gap-4 p-4 bg-base-200 rounded-lg ${className}`}
      >
        {categoryFilters.map((filter) => (
          <div key={filter.id} className="min-w-[200px] flex-1">
            <label className="label label-text text-sm font-medium">
              {filter.icon && <filter.icon className="w-4 h-4 mr-1" />}
              {t(filter.label) || filter.label}
            </label>
            {renderFilter(filter)}
          </div>
        ))}

        {activeFilterCount > 0 && (
          <button
            className="btn btn-ghost btn-sm"
            onClick={() => onFiltersChange({})}
          >
            <FaTimes className="w-4 h-4 mr-1" />
            {t('clearFilters')} ({activeFilterCount})
          </button>
        )}
      </div>
    );
  }

  // Vertical layout с группами
  return (
    <div className={`space-y-4 ${className}`}>
      {/* Заголовок с счетчиком */}
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-semibold text-lg flex items-center gap-2">
          <FaFilter className="w-4 h-4" />
          {t('filters')}
          {activeFilterCount > 0 && (
            <span className="badge badge-primary badge-sm">
              {activeFilterCount}
            </span>
          )}
        </h3>

        {activeFilterCount > 0 && (
          <button
            className="btn btn-ghost btn-xs"
            onClick={() => onFiltersChange({})}
          >
            {t('clearAll')}
          </button>
        )}
      </div>

      {/* Общие фильтры (всегда видимые) */}
      {config.showPriceRange !== false && (
        <div className="pb-4 border-b">
          <label className="label label-text font-medium">
            {t('priceRange')}
          </label>
          <div className="flex items-center gap-2">
            <input
              type="number"
              className="input input-bordered input-sm w-full"
              placeholder={t('min')}
              value={filters.priceMin || ''}
              onChange={(e) =>
                handleFilterChange(
                  'priceMin',
                  e.target.value ? Number(e.target.value) : undefined
                )
              }
            />
            <span>—</span>
            <input
              type="number"
              className="input input-bordered input-sm w-full"
              placeholder={t('max')}
              value={filters.priceMax || ''}
              onChange={(e) =>
                handleFilterChange(
                  'priceMax',
                  e.target.value ? Number(e.target.value) : undefined
                )
              }
            />
          </div>
        </div>
      )}

      {/* Группированные фильтры */}
      {filterGroups.map((group) => {
        const groupFilters = filtersByGroup[group.id];
        if (!groupFilters || groupFilters.length === 0) return null;

        const isExpanded = expandedGroups[group.id];

        return (
          <div key={group.id} className="pb-4 border-b">
            {group.collapsible ? (
              <button
                className="w-full flex items-center justify-between mb-3 hover:text-primary transition-colors"
                onClick={() => toggleGroup(group.id)}
              >
                <span className="font-medium flex items-center gap-2">
                  {group.icon && <group.icon className="w-4 h-4" />}
                  {t(group.label) || group.label}
                </span>
                {isExpanded ? <FaChevronUp /> : <FaChevronDown />}
              </button>
            ) : (
              <h4 className="font-medium mb-3 flex items-center gap-2">
                {group.icon && <group.icon className="w-4 h-4" />}
                {t(group.label) || group.label}
              </h4>
            )}

            {(!group.collapsible || isExpanded) && (
              <div className="space-y-4">
                {groupFilters.map((filter) => (
                  <div key={filter.id}>
                    {filter.type !== 'boolean' && (
                      <label className="label label-text text-sm font-medium">
                        {filter.icon && (
                          <filter.icon className="w-4 h-4 mr-1" />
                        )}
                        {t(filter.label) || filter.label}
                      </label>
                    )}
                    {renderFilter(filter)}
                  </div>
                ))}
              </div>
            )}
          </div>
        );
      })}

      {/* Фильтры без группы */}
      {filtersByGroup.other && filtersByGroup.other.length > 0 && (
        <div className="space-y-4">
          {filtersByGroup.other.map((filter) => (
            <div key={filter.id}>
              {filter.type !== 'boolean' && (
                <label className="label label-text text-sm font-medium">
                  {filter.icon && <filter.icon className="w-4 h-4 mr-1" />}
                  {t(filter.label) || filter.label}
                </label>
              )}
              {renderFilter(filter)}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default UniversalFilters;
