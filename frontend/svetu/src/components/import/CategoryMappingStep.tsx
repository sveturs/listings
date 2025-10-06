'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import type { CategoryMapping, MappingConfidence } from '@/types/import';

interface CategoryMappingStepProps {
  mappings: CategoryMapping[];
  onMappingChange: (
    externalCategory: string,
    internalCategoryId: number | null
  ) => void;
  onApproveMapping: (externalCategory: string) => void;
  onRequestNewCategory: (externalCategory: string, reasoning: string) => void;
  availableCategories: Array<{ id: number; name: string; parent?: string }>;
  isLoading?: boolean;
}

const getConfidenceBadgeColor = (
  confidence: MappingConfidence
): string => {
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

export default function CategoryMappingStep({
  mappings,
  onMappingChange,
  onApproveMapping,
  onRequestNewCategory,
  availableCategories,
  isLoading = false,
}: CategoryMappingStepProps) {
  const t = useTranslations('storefronts.import.categoryMapping');
  const [showNewCategoryForm, setShowNewCategoryForm] = useState<string | null>(
    null
  );
  const [newCategoryReasoning, setNewCategoryReasoning] = useState('');

  // Group mappings by confidence
  const groupedMappings = {
    high: mappings.filter((m) => m.confidence === 'high'),
    medium: mappings.filter((m) => m.confidence === 'medium'),
    low: mappings.filter((m) => m.confidence === 'low'),
    unmapped: mappings.filter((m) => m.suggested_internal_category_id === null),
  };

  const handleNewCategoryRequest = (externalCategory: string) => {
    if (!newCategoryReasoning.trim()) {
      alert(t('errors.reasoningRequired'));
      return;
    }
    onRequestNewCategory(externalCategory, newCategoryReasoning);
    setShowNewCategoryForm(null);
    setNewCategoryReasoning('');
  };

  const renderMappingRow = (mapping: CategoryMapping) => {
    const isEditing = showNewCategoryForm === mapping.external_category;

    return (
      <div
        key={mapping.external_category}
        className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
      >
        <div className="flex items-start justify-between">
          {/* External Category */}
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <span className="font-medium text-gray-900">
                {mapping.external_category}
              </span>
              <span
                className={`px-2 py-1 text-xs font-medium border rounded-full ${getConfidenceBadgeColor(mapping.confidence)}`}
              >
                {t(`confidence.${mapping.confidence}`)}
              </span>
              {mapping.is_approved && (
                <span className="px-2 py-1 text-xs font-medium bg-blue-100 text-blue-800 border border-blue-200 rounded-full">
                  ✓ {t('status.approved')}
                </span>
              )}
            </div>

            {/* AI Suggestion */}
            {mapping.suggested_internal_category_id && (
              <div className="mt-2 p-3 bg-blue-50 border border-blue-200 rounded-md">
                <div className="flex items-center gap-2 mb-1">
                  <svg
                    className="w-4 h-4 text-blue-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                    />
                  </svg>
                  <span className="text-sm font-medium text-blue-900">
                    {t('aiSuggestion')}:{' '}
                    {mapping.suggested_internal_category_name}
                  </span>
                </div>
                {mapping.reasoning && (
                  <p className="text-xs text-blue-700 mt-1">
                    {mapping.reasoning}
                  </p>
                )}
              </div>
            )}

            {/* Manual Category Selection */}
            <div className="mt-3">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                {t('selectCategory')}
              </label>
              <select
                value={mapping.suggested_internal_category_id || ''}
                onChange={(e) =>
                  onMappingChange(
                    mapping.external_category,
                    e.target.value ? parseInt(e.target.value) : null
                  )
                }
                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                disabled={isLoading}
              >
                <option value="">{t('selectPlaceholder')}</option>
                {availableCategories.map((cat) => (
                  <option key={cat.id} value={cat.id}>
                    {cat.parent ? `${cat.parent} > ${cat.name}` : cat.name}
                  </option>
                ))}
              </select>
            </div>

            {/* New Category Request Form */}
            {isEditing ? (
              <div className="mt-3 p-3 bg-gray-50 border border-gray-200 rounded-md">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t('newCategory.reasoning')}
                </label>
                <textarea
                  value={newCategoryReasoning}
                  onChange={(e) => setNewCategoryReasoning(e.target.value)}
                  placeholder={t('newCategory.reasoningPlaceholder')}
                  className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                  rows={3}
                />
                <div className="flex gap-2 mt-2">
                  <button
                    onClick={() =>
                      handleNewCategoryRequest(mapping.external_category)
                    }
                    className="px-3 py-1 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700"
                  >
                    {t('newCategory.submit')}
                  </button>
                  <button
                    onClick={() => setShowNewCategoryForm(null)}
                    className="px-3 py-1 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                  >
                    {t('newCategory.cancel')}
                  </button>
                </div>
              </div>
            ) : (
              <button
                onClick={() =>
                  setShowNewCategoryForm(mapping.external_category)
                }
                className="mt-2 text-sm text-blue-600 hover:text-blue-800 underline"
              >
                {t('newCategory.request')}
              </button>
            )}
          </div>

          {/* Approve Button */}
          {mapping.suggested_internal_category_id && !mapping.is_approved && (
            <button
              onClick={() => onApproveMapping(mapping.external_category)}
              className="ml-4 px-3 py-2 text-sm font-medium text-white bg-green-600 border border-transparent rounded-md hover:bg-green-700"
              disabled={isLoading}
            >
              {t('actions.approve')}
            </button>
          )}
        </div>
      </div>
    );
  };

  return (
    <div className="space-y-6">
      {/* Header with summary */}
      <div className="bg-white border border-gray-200 rounded-lg p-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">
          {t('title')}
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="text-center">
            <div className="text-2xl font-bold text-green-600">
              {groupedMappings.high.length}
            </div>
            <div className="text-sm text-gray-600">
              {t('summary.highConfidence')}
            </div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-yellow-600">
              {groupedMappings.medium.length}
            </div>
            <div className="text-sm text-gray-600">
              {t('summary.mediumConfidence')}
            </div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-red-600">
              {groupedMappings.low.length}
            </div>
            <div className="text-sm text-gray-600">
              {t('summary.lowConfidence')}
            </div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-gray-600">
              {groupedMappings.unmapped.length}
            </div>
            <div className="text-sm text-gray-600">
              {t('summary.unmapped')}
            </div>
          </div>
        </div>
      </div>

      {/* High Confidence Mappings (Auto-approved) */}
      {groupedMappings.high.length > 0 && (
        <div>
          <h4 className="text-md font-semibold text-green-700 mb-3 flex items-center gap-2">
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            {t('sections.highConfidence')}
          </h4>
          <div className="space-y-3">
            {groupedMappings.high.map(renderMappingRow)}
          </div>
        </div>
      )}

      {/* Medium Confidence Mappings (Review Recommended) */}
      {groupedMappings.medium.length > 0 && (
        <div>
          <h4 className="text-md font-semibold text-yellow-700 mb-3 flex items-center gap-2">
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
            {t('sections.mediumConfidence')}
          </h4>
          <div className="space-y-3">
            {groupedMappings.medium.map(renderMappingRow)}
          </div>
        </div>
      )}

      {/* Low Confidence Mappings (Manual Required) */}
      {groupedMappings.low.length > 0 && (
        <div>
          <h4 className="text-md font-semibold text-red-700 mb-3 flex items-center gap-2">
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            {t('sections.lowConfidence')}
          </h4>
          <div className="space-y-3">
            {groupedMappings.low.map(renderMappingRow)}
          </div>
        </div>
      )}

      {/* Unmapped Categories */}
      {groupedMappings.unmapped.length > 0 && (
        <div>
          <h4 className="text-md font-semibold text-gray-700 mb-3 flex items-center gap-2">
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            {t('sections.unmapped')}
          </h4>
          <div className="space-y-3">
            {groupedMappings.unmapped.map(renderMappingRow)}
          </div>
        </div>
      )}

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
              <li>{t('help.highConfidence')}</li>
              <li>{t('help.mediumConfidence')}</li>
              <li>{t('help.lowConfidence')}</li>
              <li>{t('help.newCategory')}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}
