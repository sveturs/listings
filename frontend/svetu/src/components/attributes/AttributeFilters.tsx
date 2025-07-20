import React, { useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { AttributeFilters } from '@/hooks/useAttributesPagination';

export interface AttributeFiltersProps {
  filters: AttributeFilters;
  onFiltersChange: (filters: Partial<AttributeFilters>) => void;
  onClearFilters: () => void;
  loading?: boolean;
}

const ATTRIBUTE_TYPES = [
  'text',
  'number',
  'select',
  'multiselect',
  'boolean',
  'date',
  'range',
] as const;

export const AttributeFiltersComponent: React.FC<AttributeFiltersProps> = ({
  filters,
  onFiltersChange,
  onClearFilters,
  loading = false,
}) => {
  const t = useTranslations('admin');
  const [searchInput, setSearchInput] = useState(filters.search || '');

  // Debounced search
  const [searchTimeout, setSearchTimeout] = useState<NodeJS.Timeout | null>(
    null
  );

  const handleSearchChange = useCallback(
    (value: string) => {
      setSearchInput(value);

      if (searchTimeout) {
        clearTimeout(searchTimeout);
      }

      const timeout = setTimeout(() => {
        onFiltersChange({ search: value || undefined });
      }, 500);

      setSearchTimeout(timeout);
    },
    [onFiltersChange, searchTimeout]
  );

  const handleTypeChange = useCallback(
    (type: string) => {
      onFiltersChange({ type: type || undefined });
    },
    [onFiltersChange]
  );

  const hasActiveFilters = !!(filters.search || filters.type);

  return (
    <div className="bg-base-100 border border-base-300 rounded-lg p-4 space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium">{t('attributes.filters')}</h3>
        {hasActiveFilters && (
          <button
            className="btn btn-ghost btn-sm"
            onClick={onClearFilters}
            disabled={loading}
          >
            {t('common.clearFilters')}
          </button>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {/* Search input */}
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('common.search')}</span>
          </label>
          <div className="relative">
            <input
              type="text"
              placeholder={t('attributes.searchPlaceholder')}
              className="input input-bordered w-full pr-10"
              value={searchInput}
              onChange={(e) => handleSearchChange(e.target.value)}
              disabled={loading}
            />
            <div className="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
              <svg
                className="w-4 h-4 text-base-content/50"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
          </div>
        </div>

        {/* Type filter */}
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('attributes.type')}</span>
          </label>
          <select
            className="select select-bordered w-full"
            value={filters.type || ''}
            onChange={(e) => handleTypeChange(e.target.value)}
            disabled={loading}
          >
            <option value="">{t('common.all')}</option>
            {ATTRIBUTE_TYPES.map((type) => (
              <option key={type} value={type}>
                {t(`attributes.types.${type}`)}
              </option>
            ))}
          </select>
        </div>

        {/* Quick stats */}
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('common.status')}</span>
          </label>
          <div className="flex flex-wrap gap-2">
            <div className="badge badge-info badge-sm">
              {t('attributes.totalFound')}
            </div>
            {hasActiveFilters && (
              <div className="badge badge-warning badge-sm">
                {t('common.filtered')}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Active filters display */}
      {hasActiveFilters && (
        <div className="flex flex-wrap gap-2 pt-2 border-t border-base-300">
          <span className="text-sm text-base-content/70">
            {t('common.activeFilters')}:
          </span>

          {filters.search && (
            <div className="badge badge-primary gap-1">
              <span>
                {t('common.search')}: &quot;{filters.search}&quot;
              </span>
              <button
                className="btn btn-ghost btn-xs"
                onClick={() => {
                  setSearchInput('');
                  onFiltersChange({ search: undefined });
                }}
              >
                ✕
              </button>
            </div>
          )}

          {filters.type && (
            <div className="badge badge-secondary gap-1">
              <span>
                {t('attributes.type')}: {t(`attributes.types.${filters.type}`)}
              </span>
              <button
                className="btn btn-ghost btn-xs"
                onClick={() => onFiltersChange({ type: undefined })}
              >
                ✕
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
};
