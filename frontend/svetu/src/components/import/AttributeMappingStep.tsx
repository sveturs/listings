'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import type { DetectedAttribute } from '@/types/import';

interface AttributeMappingStepProps {
  attributes: DetectedAttribute[];
  selectedAttributes: string[];
  onToggleAttribute: (attributeName: string) => void;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  isLoading?: boolean;
}

const getValueTypeBadgeColor = (valueType: string): string => {
  switch (valueType) {
    case 'string':
      return 'bg-blue-100 text-blue-800 border-blue-200';
    case 'number':
      return 'bg-green-100 text-green-800 border-green-200';
    case 'boolean':
      return 'bg-purple-100 text-purple-800 border-purple-200';
    case 'enum':
      return 'bg-orange-100 text-orange-800 border-orange-200';
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200';
  }
};

export default function AttributeMappingStep({
  attributes,
  selectedAttributes,
  onToggleAttribute,
  onSelectAll,
  onDeselectAll,
  isLoading = false,
}: AttributeMappingStepProps) {
  const t = useTranslations('storefronts.import.attributeMapping');
  const [searchQuery, setSearchQuery] = useState('');
  const [filterType, setFilterType] = useState<string>('all');
  const [showVariantDefiningOnly, setShowVariantDefiningOnly] = useState(false);

  // Filter attributes
  const filteredAttributes = attributes.filter((attr) => {
    const matchesSearch =
      searchQuery === '' ||
      attr.name.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesType =
      filterType === 'all' || attr.value_type === filterType;
    const matchesVariantFilter =
      !showVariantDefiningOnly || attr.is_variant_defining;

    return matchesSearch && matchesType && matchesVariantFilter;
  });

  // Stats
  const stats = {
    total: attributes.length,
    selected: selectedAttributes.length,
    variantDefining: attributes.filter((a) => a.is_variant_defining).length,
    byType: {
      string: attributes.filter((a) => a.value_type === 'string').length,
      number: attributes.filter((a) => a.value_type === 'number').length,
      boolean: attributes.filter((a) => a.value_type === 'boolean').length,
      enum: attributes.filter((a) => a.value_type === 'enum').length,
    },
  };

  return (
    <div className="space-y-6">
      {/* Header with Stats */}
      <div className="bg-white border border-gray-200 rounded-lg p-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">
          {t('title')}
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
          <div className="text-center">
            <div className="text-2xl font-bold text-blue-600">
              {stats.total}
            </div>
            <div className="text-sm text-gray-600">{t('stats.total')}</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-green-600">
              {stats.selected}
            </div>
            <div className="text-sm text-gray-600">{t('stats.selected')}</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-purple-600">
              {stats.variantDefining}
            </div>
            <div className="text-sm text-gray-600">
              {t('stats.variantDefining')}
            </div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-orange-600">
              {stats.byType.string}
            </div>
            <div className="text-sm text-gray-600">{t('types.string')}</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-teal-600">
              {stats.byType.number}
            </div>
            <div className="text-sm text-gray-600">{t('types.number')}</div>
          </div>
        </div>
      </div>

      {/* Filters and Search */}
      <div className="bg-white border border-gray-200 rounded-lg p-4">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Search */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('filters.search')}
            </label>
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder={t('filters.searchPlaceholder')}
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          {/* Type Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('filters.type')}
            </label>
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value)}
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            >
              <option value="all">{t('filters.allTypes')}</option>
              <option value="string">{t('types.string')}</option>
              <option value="number">{t('types.number')}</option>
              <option value="boolean">{t('types.boolean')}</option>
              <option value="enum">{t('types.enum')}</option>
            </select>
          </div>

          {/* Variant Defining Filter */}
          <div className="flex items-end">
            <label className="flex items-center cursor-pointer">
              <input
                type="checkbox"
                checked={showVariantDefiningOnly}
                onChange={(e) => setShowVariantDefiningOnly(e.target.checked)}
                className="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              />
              <span className="ml-2 text-sm text-gray-700">
                {t('filters.variantDefiningOnly')}
              </span>
            </label>
          </div>
        </div>

        {/* Bulk Actions */}
        <div className="flex gap-2 mt-4">
          <button
            onClick={onSelectAll}
            disabled={isLoading}
            className="px-3 py-2 text-sm font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-md hover:bg-blue-100 disabled:opacity-50"
          >
            {t('actions.selectAll')}
          </button>
          <button
            onClick={onDeselectAll}
            disabled={isLoading}
            className="px-3 py-2 text-sm font-medium text-gray-600 bg-gray-50 border border-gray-200 rounded-md hover:bg-gray-100 disabled:opacity-50"
          >
            {t('actions.deselectAll')}
          </button>
        </div>
      </div>

      {/* Attributes List */}
      <div className="space-y-2">
        {filteredAttributes.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            {t('noAttributesFound')}
          </div>
        ) : (
          filteredAttributes.map((attr) => {
            const isSelected = selectedAttributes.includes(attr.name);

            return (
              <div
                key={attr.name}
                className={`border rounded-lg p-4 transition-all ${
                  isSelected
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 bg-white hover:shadow-md'
                }`}
              >
                <div className="flex items-start gap-4">
                  {/* Checkbox */}
                  <div className="flex items-center h-5 mt-1">
                    <input
                      type="checkbox"
                      checked={isSelected}
                      onChange={() => onToggleAttribute(attr.name)}
                      disabled={isLoading}
                      className="w-5 h-5 rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                    />
                  </div>

                  {/* Attribute Info */}
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <span className="font-medium text-gray-900">
                        {attr.name}
                      </span>
                      <span
                        className={`px-2 py-1 text-xs font-medium border rounded-full ${getValueTypeBadgeColor(attr.value_type)}`}
                      >
                        {attr.value_type}
                      </span>
                      {attr.is_variant_defining && (
                        <span className="px-2 py-1 text-xs font-medium bg-purple-100 text-purple-800 border border-purple-200 rounded-full">
                          🎨 {t('badges.variantDefining')}
                        </span>
                      )}
                    </div>

                    {/* Suggested Mapping */}
                    {attr.suggested_mapping && (
                      <div className="mb-2 text-sm text-gray-600">
                        <span className="font-medium">
                          {t('suggestedMapping')}:
                        </span>{' '}
                        {attr.suggested_mapping}
                      </div>
                    )}

                    {/* Stats */}
                    <div className="flex items-center gap-4 text-sm text-gray-600">
                      <span>
                        {t('frequency')}: <strong>{attr.frequency}</strong>{' '}
                        {t('products')}
                      </span>
                      <span>•</span>
                      <span>
                        {t('sampleValues')}:{' '}
                        <span className="italic">
                          {attr.sample_values?.slice(0, 3).join(', ') || 'N/A'}
                          {attr.sample_values && attr.sample_values.length > 3 && '...'}
                        </span>
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            );
          })
        )}
      </div>

      {/* Help Text */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <div className="flex gap-2">
          <svg
            className="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <div>
            <p className="text-sm text-blue-900 font-medium mb-1">
              {t('help.title')}
            </p>
            <ul className="text-sm text-blue-800 space-y-1 list-disc list-inside">
              <li>{t('help.selectAttributes')}</li>
              <li>{t('help.variantDefining')}</li>
              <li>{t('help.frequency')}</li>
              <li>{t('help.sampleValues')}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}
