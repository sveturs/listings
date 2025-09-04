'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { Filter, Tag } from 'lucide-react';
import api from '@/services/api';

interface CategoryAttribute {
  id: number;
  name: string;
  display_name: string;
  attribute_type: string;
  options?: any;
  is_required: boolean;
  is_filterable: boolean;
}

interface GenericCategoryFiltersProps {
  categoryId: number;
  onFiltersChange: (filters: Record<string, any>) => void;
  className?: string;
}

export const GenericCategoryFilters: React.FC<GenericCategoryFiltersProps> = ({
  categoryId,
  onFiltersChange,
  className = '',
}) => {
  const t = useTranslations('search');
  const [attributes, setAttributes] = useState<CategoryAttribute[]>([]);
  const [loading, setLoading] = useState(true);
  const [filterValues, setFilterValues] = useState<Record<string, any>>({});

  useEffect(() => {
    loadCategoryAttributes();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categoryId]);

  useEffect(() => {
    onFiltersChange(filterValues);
  }, [filterValues, onFiltersChange]);

  const loadCategoryAttributes = async () => {
    setLoading(true);
    try {
      const response = await api.get(`/categories/${categoryId}/attributes`);
      if (response.data?.data) {
        const filterableAttributes = response.data.data.filter(
          (attr: CategoryAttribute) => attr.is_filterable
        );
        setAttributes(filterableAttributes);
      }
    } catch (error) {
      console.error('Error loading category attributes:', error);
      setAttributes([]);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (attributeName: string, value: any) => {
    setFilterValues((prev) => {
      const updated = { ...prev };
      if (value === '' || value === null || value === undefined) {
        delete updated[attributeName];
      } else {
        updated[attributeName] = value;
      }
      return updated;
    });
  };

  const renderAttributeFilter = (attribute: CategoryAttribute) => {
    const { id, name, display_name, attribute_type, options } = attribute;

    switch (attribute_type) {
      case 'select':
        return (
          <div key={id}>
            <label className="text-xs font-medium text-base-content/70 mb-2">
              {display_name || name}
            </label>
            <select
              value={filterValues[name] || ''}
              onChange={(e) => handleFilterChange(name, e.target.value)}
              className="select select-bordered select-sm w-full"
            >
              <option value="">{t('all')}</option>
              {options?.values?.map((opt: string) => (
                <option key={opt} value={opt}>
                  {opt}
                </option>
              ))}
            </select>
          </div>
        );

      case 'multiselect':
        return (
          <div key={id}>
            <label className="text-xs font-medium text-base-content/70 mb-2">
              {display_name || name}
            </label>
            <div className="space-y-2 max-h-32 overflow-y-auto">
              {options?.values?.map((opt: string) => (
                <label
                  key={opt}
                  className="flex items-center gap-2 cursor-pointer"
                >
                  <input
                    type="checkbox"
                    checked={filterValues[name]?.includes(opt) || false}
                    onChange={(e) => {
                      const currentValues = filterValues[name] || [];
                      if (e.target.checked) {
                        handleFilterChange(name, [...currentValues, opt]);
                      } else {
                        handleFilterChange(
                          name,
                          currentValues.filter((v: string) => v !== opt)
                        );
                      }
                    }}
                    className="checkbox checkbox-sm checkbox-primary"
                  />
                  <span className="text-sm">{opt}</span>
                </label>
              ))}
            </div>
          </div>
        );

      case 'number':
        return (
          <div key={id}>
            <label className="text-xs font-medium text-base-content/70 mb-2">
              {display_name || name}
            </label>
            <div className="flex gap-2">
              <input
                type="number"
                placeholder={t('from')}
                value={filterValues[`${name}_min`] || ''}
                onChange={(e) =>
                  handleFilterChange(`${name}_min`, e.target.value)
                }
                className="input input-bordered input-sm w-full"
              />
              <input
                type="number"
                placeholder={t('to')}
                value={filterValues[`${name}_max`] || ''}
                onChange={(e) =>
                  handleFilterChange(`${name}_max`, e.target.value)
                }
                className="input input-bordered input-sm w-full"
              />
            </div>
          </div>
        );

      case 'boolean':
        return (
          <div key={id}>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={filterValues[name] || false}
                onChange={(e) => handleFilterChange(name, e.target.checked)}
                className="checkbox checkbox-sm checkbox-primary"
              />
              <span className="text-sm">{display_name || name}</span>
            </label>
          </div>
        );

      case 'text':
        return (
          <div key={id}>
            <label className="text-xs font-medium text-base-content/70 mb-2">
              {display_name || name}
            </label>
            <input
              type="text"
              placeholder={t('enter', { field: display_name || name })}
              value={filterValues[name] || ''}
              onChange={(e) => handleFilterChange(name, e.target.value)}
              className="input input-bordered input-sm w-full"
            />
          </div>
        );

      default:
        return null;
    }
  };

  if (loading) {
    return (
      <div className={`space-y-4 ${className}`}>
        <div className="flex items-center justify-center py-8">
          <span className="loading loading-spinner loading-sm"></span>
          <span className="ml-2 text-sm text-base-content/60">
            {t('loadingFilters')}
          </span>
        </div>
      </div>
    );
  }

  if (attributes.length === 0) {
    return (
      <div className={`space-y-4 ${className}`}>
        <div className="text-center py-4 text-sm text-base-content/60">
          <Tag className="w-8 h-8 mx-auto mb-2 opacity-50" />
          {t('noCategoryFilters')}
        </div>
      </div>
    );
  }

  return (
    <div className={`space-y-4 ${className}`}>
      <div className="flex items-center gap-2 mb-3">
        <Filter className="w-4 h-4 text-base-content/70" />
        <span className="text-xs font-semibold text-base-content/70">
          {t('categorySpecificFilters')}
        </span>
      </div>
      {attributes.map(renderAttributeFilter)}
    </div>
  );
};
