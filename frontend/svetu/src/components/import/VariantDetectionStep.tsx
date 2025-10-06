'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import type { VariantGroup, MappingConfidence } from '@/types/import';

interface VariantDetectionStepProps {
  variantGroups: VariantGroup[];
  approvedGroups: string[]; // base_names of approved groups
  onToggleGroup: (baseName: string) => void;
  onApproveAll: () => void;
  onRejectAll: () => void;
  isLoading?: boolean;
}

const getConfidenceBadgeColor = (confidence: MappingConfidence): string => {
  switch (confidence) {
    case 'high':
      return 'bg-green-100 text-green-800 border-green-200';
    case 'medium':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    case 'low':
      return 'bg-red-100 text-red-800 border-red-200';
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200';
  }
};

export default function VariantDetectionStep({
  variantGroups,
  approvedGroups,
  onToggleGroup,
  onApproveAll,
  onRejectAll,
  isLoading = false,
}: VariantDetectionStepProps) {
  const t = useTranslations('storefronts.import.variantDetection');
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(new Set());
  const [searchQuery, setSearchQuery] = useState('');
  const [filterConfidence, setFilterConfidence] = useState<string>('all');

  // Filter groups
  const filteredGroups = variantGroups.filter((group) => {
    const matchesSearch =
      searchQuery === '' ||
      group.base_name.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesConfidence =
      filterConfidence === 'all' || group.confidence === filterConfidence;

    return matchesSearch && matchesConfidence;
  });

  // Stats
  const stats = {
    totalGroups: variantGroups.length,
    totalVariants: variantGroups.reduce(
      (sum, g) => sum + g.variant_count,
      0
    ),
    approved: approvedGroups.length,
    byConfidence: {
      high: variantGroups.filter((g) => g.confidence === 'high').length,
      medium: variantGroups.filter((g) => g.confidence === 'medium').length,
      low: variantGroups.filter((g) => g.confidence === 'low').length,
    },
  };

  const toggleExpanded = (baseName: string) => {
    const newExpanded = new Set(expandedGroups);
    if (newExpanded.has(baseName)) {
      newExpanded.delete(baseName);
    } else {
      newExpanded.add(baseName);
    }
    setExpandedGroups(newExpanded);
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
              {stats.totalGroups}
            </div>
            <div className="text-sm text-gray-600">{t('stats.groups')}</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-purple-600">
              {stats.totalVariants}
            </div>
            <div className="text-sm text-gray-600">{t('stats.variants')}</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-green-600">
              {stats.approved}
            </div>
            <div className="text-sm text-gray-600">{t('stats.approved')}</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-green-600">
              {stats.byConfidence.high}
            </div>
            <div className="text-sm text-gray-600">{t('confidence.high')}</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-yellow-600">
              {stats.byConfidence.medium}
            </div>
            <div className="text-sm text-gray-600">
              {t('confidence.medium')}
            </div>
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white border border-gray-200 rounded-lg p-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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

          {/* Confidence Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {t('filters.confidence')}
            </label>
            <select
              value={filterConfidence}
              onChange={(e) => setFilterConfidence(e.target.value)}
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            >
              <option value="all">{t('filters.allConfidence')}</option>
              <option value="high">{t('confidence.high')}</option>
              <option value="medium">{t('confidence.medium')}</option>
              <option value="low">{t('confidence.low')}</option>
            </select>
          </div>
        </div>

        {/* Bulk Actions */}
        <div className="flex gap-2 mt-4">
          <button
            onClick={onApproveAll}
            disabled={isLoading}
            className="px-3 py-2 text-sm font-medium text-green-600 bg-green-50 border border-green-200 rounded-md hover:bg-green-100 disabled:opacity-50"
          >
            {t('actions.approveAll')}
          </button>
          <button
            onClick={onRejectAll}
            disabled={isLoading}
            className="px-3 py-2 text-sm font-medium text-gray-600 bg-gray-50 border border-gray-200 rounded-md hover:bg-gray-100 disabled:opacity-50"
          >
            {t('actions.rejectAll')}
          </button>
        </div>
      </div>

      {/* Variant Groups List */}
      <div className="space-y-3">
        {filteredGroups.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            {t('noGroupsFound')}
          </div>
        ) : (
          filteredGroups.map((group) => {
            const isExpanded = expandedGroups.has(group.base_name);
            const isApproved = approvedGroups.includes(group.base_name);

            return (
              <div
                key={group.base_name}
                className={`border rounded-lg transition-all ${
                  isApproved
                    ? 'border-green-500 bg-green-50'
                    : 'border-gray-200 bg-white hover:shadow-md'
                }`}
              >
                {/* Group Header */}
                <div className="p-4">
                  <div className="flex items-start gap-4">
                    {/* Checkbox */}
                    <div className="flex items-center h-5 mt-1">
                      <input
                        type="checkbox"
                        checked={isApproved}
                        onChange={() => onToggleGroup(group.base_name)}
                        disabled={isLoading}
                        className="w-5 h-5 rounded border-gray-300 text-green-600 shadow-sm focus:border-green-500 focus:ring-green-500"
                      />
                    </div>

                    {/* Group Info */}
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2">
                        <span className="font-medium text-gray-900">
                          {group.base_name}
                        </span>
                        <span
                          className={`px-2 py-1 text-xs font-medium border rounded-full ${getConfidenceBadgeColor(group.confidence)}`}
                        >
                          {t(`confidence.${group.confidence}`)}
                        </span>
                        {isApproved && (
                          <span className="px-2 py-1 text-xs font-medium bg-green-100 text-green-800 border border-green-200 rounded-full">
                            ✓ {t('status.approved')}
                          </span>
                        )}
                      </div>

                      {/* Variant Stats */}
                      <div className="flex items-center gap-4 text-sm text-gray-600">
                        <span>
                          <strong>{group.variant_count}</strong>{' '}
                          {t('variants')}
                        </span>
                        <span>•</span>
                        <span>
                          {t('attributes')}:{' '}
                          <span className="italic">
                            {group.variant_attributes.join(', ')}
                          </span>
                        </span>
                      </div>

                      {/* Expand/Collapse Button */}
                      <button
                        onClick={() => toggleExpanded(group.base_name)}
                        className="mt-2 text-sm text-blue-600 hover:text-blue-800 underline"
                      >
                        {isExpanded
                          ? t('actions.hideVariants')
                          : t('actions.showVariants')}
                      </button>
                    </div>
                  </div>
                </div>

                {/* Expanded Variants */}
                {isExpanded && (
                  <div className="border-t border-gray-200 p-4 bg-gray-50">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                      {group.products.map((product) => (
                        <div
                          key={product.sku}
                          className="bg-white border border-gray-200 rounded-md p-3"
                        >
                          <div className="flex gap-3">
                            {/* Product Image */}
                            {product.images && product.images[0] && (
                              <img
                                src={product.images[0]}
                                alt={product.name}
                                className="w-16 h-16 object-cover rounded border border-gray-300"
                              />
                            )}

                            {/* Product Info */}
                            <div className="flex-1 min-w-0">
                              <div className="text-sm font-medium text-gray-900 truncate">
                                {product.name}
                              </div>
                              <div className="text-xs text-gray-500 mt-1">
                                SKU: {product.sku}
                              </div>
                              <div className="text-sm font-semibold text-gray-900 mt-1">
                                ${product.price.toFixed(2)}
                              </div>
                              {/* Variant Values */}
                              <div className="flex flex-wrap gap-1 mt-2">
                                {Object.entries(product.variant_values).map(
                                  ([key, value]) => (
                                    <span
                                      key={key}
                                      className="px-2 py-0.5 text-xs bg-blue-100 text-blue-800 rounded"
                                    >
                                      {key}: {value}
                                    </span>
                                  )
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
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
              <li>{t('help.variantGroups')}</li>
              <li>{t('help.confidence')}</li>
              <li>{t('help.approve')}</li>
              <li>{t('help.expandDetails')}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}
