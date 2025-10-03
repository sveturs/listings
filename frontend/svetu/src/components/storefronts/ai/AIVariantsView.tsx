'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateAIProduct } from '@/contexts/CreateAIProductContext';
import SimplifiedVariantGenerator from '@/components/products/SimplifiedVariantGenerator';

interface AIVariantsViewProps {
  storefrontId: number | null;
  storefrontSlug: string;
}

export default function AIVariantsView({
  storefrontId: _storefrontId,
  storefrontSlug: _storefrontSlug,
}: AIVariantsViewProps) {
  const t = useTranslations('storefronts');
  const { state, setView, setHasVariants, setVariants } = useCreateAIProduct();

  const [activeMode, setActiveMode] = useState<'none' | 'simple' | 'advanced'>(
    state.aiData.hasVariants
      ? state.aiData.variants.length > 0
        ? 'advanced'
        : 'simple'
      : 'none'
  );

  // State for attribute selection
  const [availableAttributes, setAvailableAttributes] = useState<any[]>([]);
  const [selectedAttributes, setSelectedAttributes] = useState<
    Record<number, any>
  >({});
  const [loadingAttributes, setLoadingAttributes] = useState(false);
  const [showAttributeSelector, setShowAttributeSelector] = useState(false);

  // AI suggestions
  const aiSuggestedVariants = state.aiData.suggestedVariants || [];
  const hasAISuggestions = aiSuggestedVariants.length > 0;

  // Load variant attributes for the selected category
  useEffect(() => {
    const loadAttributes = async () => {
      console.log(
        '[AIVariantsView] Loading attributes for category:',
        state.aiData.categoryId
      );
      setLoadingAttributes(true);
      try {
        // Сначала получаем slug категории по ID
        const categoriesResponse = await fetch(
          `/api/v2/marketplace/categories?page=1&limit=1000`
        );

        if (!categoriesResponse.ok) {
          throw new Error('Failed to load categories');
        }

        const categoriesData = await categoriesResponse.json();
        console.log(
          '[AIVariantsView] Categories loaded:',
          categoriesData.data?.length
        );

        const category = categoriesData.data?.find(
          (cat: any) => cat.id === state.aiData.categoryId
        );

        if (!category || !category.slug) {
          console.warn('[AIVariantsView] Category not found or has no slug:', {
            categoryId: state.aiData.categoryId,
            found: !!category,
            slug: category?.slug,
          });
          setAvailableAttributes([]);
          return;
        }

        console.log('[AIVariantsView] Category found:', category.slug);

        // Загружаем variant attributes для этой категории
        const variantAttrsResponse = await fetch(
          `/api/v2/marketplace/categories/${category.slug}/variant-attributes`
        );

        if (variantAttrsResponse.ok) {
          const variantAttrsData = await variantAttrsResponse.json();
          console.log('[AIVariantsView] Variant attributes loaded:', {
            slug: category.slug,
            count: variantAttrsData.data?.length || 0,
            attributes: variantAttrsData.data,
          });
          setAvailableAttributes(variantAttrsData.data || []);
        } else {
          console.warn(
            '[AIVariantsView] No variant attributes found for category:',
            {
              slug: category.slug,
              status: variantAttrsResponse.status,
            }
          );
          setAvailableAttributes([]);
        }
      } catch (error) {
        console.error(
          '[AIVariantsView] Failed to load variant attributes:',
          error
        );
        setAvailableAttributes([]);
      } finally {
        setLoadingAttributes(false);
      }
    };

    console.log('[AIVariantsView] useEffect triggered:', {
      hasVariants: state.aiData.hasVariants,
      activeMode,
      categoryId: state.aiData.categoryId,
      shouldLoad:
        state.aiData.hasVariants &&
        activeMode === 'simple' &&
        state.aiData.categoryId,
    });

    if (
      state.aiData.hasVariants &&
      activeMode === 'simple' &&
      state.aiData.categoryId
    ) {
      loadAttributes();
    }
  }, [state.aiData.hasVariants, activeMode, state.aiData.categoryId]);

  const handleAttributeSelect = (attrId: number, attr: any) => {
    if (selectedAttributes[attrId]) {
      // Deselect
      const newSelected = { ...selectedAttributes };
      delete newSelected[attrId];
      setSelectedAttributes(newSelected);
    } else {
      // Select with empty values
      setSelectedAttributes({
        ...selectedAttributes,
        [attrId]: {
          id: attrId,
          name: attr.name,
          display_name: attr.display_name,
          values: [],
        },
      });
    }
  };

  const handleAttributeValuesChange = (attrId: number, values: string[]) => {
    if (selectedAttributes[attrId]) {
      setSelectedAttributes({
        ...selectedAttributes,
        [attrId]: {
          ...selectedAttributes[attrId],
          values: values,
        },
      });
    }
  };

  const handleVariantToggle = (enabled: boolean) => {
    setHasVariants(enabled);
    if (!enabled) {
      setActiveMode('none');
      setVariants([]);
    } else {
      setActiveMode('simple');
    }
  };

  const handleApplyAISuggestions = () => {
    const variants = aiSuggestedVariants.map((suggestion, idx) => ({
      sku: suggestion.sku,
      price: suggestion.price || state.aiData.price,
      stock_quantity: suggestion.stockQuantity || 0,
      variant_attributes: suggestion.attributes,
      is_default: idx === 0,
    }));
    setVariants(variants);
    setHasVariants(true);
    setActiveMode('advanced');
  };

  const handleVariantsSave = (newVariants: any[]) => {
    const formattedVariants = newVariants.map((v) => ({
      sku: v.sku,
      barcode: v.barcode,
      price: v.price,
      compare_at_price: v.compare_at_price,
      cost_price: v.cost_price,
      stock_quantity: v.stock_quantity || 0,
      low_stock_threshold: v.low_stock_threshold,
      variant_attributes: v.variant_attributes || {},
      weight: v.weight,
      dimensions: v.dimensions,
      is_default: v.is_default || false,
    }));
    setVariants(formattedVariants);
  };

  const handleNext = () => {
    setView('publish');
  };

  const handleSkip = () => {
    setHasVariants(false);
    setVariants([]);
    setView('publish');
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold mb-2">
          {t('productVariants') || 'Product Variants'}
        </h2>
        <p className="text-base-content/70">
          {t('variantsDescription') ||
            'Create variants like sizes or colors (optional)'}
        </p>
      </div>

      {/* AI Suggestions Alert */}
      {hasAISuggestions && state.aiData.variants.length === 0 && (
        <div className="alert alert-info">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            className="stroke-current shrink-0 w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <div>
            <div className="font-semibold">
              {t('aiSuggestedVariants') || 'AI Suggested Variants'}
            </div>
            <div className="text-sm">
              {t('aiSuggestedVariantsDescription', {
                count: aiSuggestedVariants.length,
              }) ||
                `AI detected ${aiSuggestedVariants.length} potential variants. Apply them automatically?`}
            </div>
          </div>
          <button
            onClick={handleApplyAISuggestions}
            className="btn btn-sm btn-primary"
          >
            {t('applyAISuggestions') || 'Apply AI Suggestions'}
          </button>
        </div>
      )}

      {/* Variant Toggle */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <h3 className="text-lg font-medium">
                {t('enableVariants') || 'Enable Product Variants'}
              </h3>
              <p className="text-sm text-base-content/70">
                {t('enableVariantsDescription') ||
                  'Create variants if your product has multiple options like size or color'}
              </p>
            </div>
            <div className="form-control">
              <label className="label cursor-pointer">
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={state.aiData.hasVariants}
                  onChange={(e) => handleVariantToggle(e.target.checked)}
                />
              </label>
            </div>
          </div>
        </div>
      </div>

      {/* Variants Configuration */}
      {state.aiData.hasVariants && (
        <div className="space-y-6">
          {/* Mode Selection */}
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h3 className="text-lg font-medium mb-4">
                {t('variantMode') || 'Variant Creation Mode'}
              </h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Simple Mode */}
                <div
                  className={`
                    card border-2 cursor-pointer transition-colors
                    ${
                      activeMode === 'simple'
                        ? 'border-primary bg-primary/5'
                        : 'border-base-300 hover:border-base-400'
                    }
                  `}
                  onClick={() => setActiveMode('simple')}
                >
                  <div className="card-body p-4">
                    <div className="flex items-center space-x-3">
                      <input
                        type="radio"
                        name="variantMode"
                        className="radio radio-primary"
                        checked={activeMode === 'simple'}
                        readOnly
                      />
                      <div>
                        <h4 className="font-medium">
                          {t('simpleMode') || 'Simple Mode'}
                        </h4>
                        <p className="text-sm text-base-content/70">
                          {t('simpleModeDescription') ||
                            'Quick variant creation'}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Advanced Mode */}
                <div
                  className={`
                    card border-2 cursor-pointer transition-colors
                    ${
                      activeMode === 'advanced'
                        ? 'border-primary bg-primary/5'
                        : 'border-base-300 hover:border-base-400'
                    }
                  `}
                  onClick={() => setActiveMode('advanced')}
                >
                  <div className="card-body p-4">
                    <div className="flex items-center space-x-3">
                      <input
                        type="radio"
                        name="variantMode"
                        className="radio radio-primary"
                        checked={activeMode === 'advanced'}
                        readOnly
                      />
                      <div>
                        <h4 className="font-medium">
                          {t('advancedMode') || 'Advanced Mode'}
                        </h4>
                        <p className="text-sm text-base-content/70">
                          {t('advancedModeDescription') ||
                            'Full control over variants'}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Variant Components */}
          {activeMode === 'simple' && (
            <>
              {/* Attribute Selection */}
              {!showAttributeSelector ? (
                <div className="card bg-base-100 shadow-sm">
                  <div className="card-body">
                    <h3 className="text-lg font-medium mb-4">
                      {t('selectVariantAttributes') ||
                        'Select Variant Attributes'}
                    </h3>

                    {loadingAttributes ? (
                      <div className="flex justify-center py-8">
                        <span className="loading loading-spinner loading-lg"></span>
                      </div>
                    ) : availableAttributes.length === 0 ? (
                      <div className="alert alert-warning">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 24 24"
                          className="stroke-current shrink-0 w-6 h-6"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                          />
                        </svg>
                        <span>
                          {t('noVariantAttributesAvailable') ||
                            'No variant attributes available. Please contact administrator.'}
                        </span>
                      </div>
                    ) : (
                      <>
                        <div className="space-y-3">
                          {availableAttributes.map((attr) => (
                            <div
                              key={attr.id}
                              className="border border-base-300 rounded-lg p-4"
                            >
                              <div className="flex items-start gap-3">
                                <input
                                  type="checkbox"
                                  className="checkbox checkbox-primary mt-1"
                                  checked={!!selectedAttributes[attr.id]}
                                  onChange={() =>
                                    handleAttributeSelect(attr.id, attr)
                                  }
                                />
                                <div className="flex-1">
                                  <div className="font-medium">
                                    {attr.display_name || attr.name}
                                  </div>
                                  {attr.description && (
                                    <div className="text-sm text-base-content/70 mt-1">
                                      {attr.description}
                                    </div>
                                  )}

                                  {selectedAttributes[attr.id] && (
                                    <div className="mt-3">
                                      <label className="label">
                                        <span className="label-text">
                                          {t('enterValues') ||
                                            'Enter values (comma-separated)'}
                                        </span>
                                      </label>
                                      <input
                                        type="text"
                                        className="input input-bordered w-full"
                                        placeholder={
                                          t('exampleValues', {
                                            example:
                                              attr.name === 'Color'
                                                ? 'Red, Blue, Green'
                                                : 'S, M, L, XL',
                                          }) || 'e.g., Red, Blue, Green'
                                        }
                                        value={selectedAttributes[
                                          attr.id
                                        ].values.join(', ')}
                                        onChange={(e) => {
                                          // Разрешаем пользователю вводить любые символы
                                          // Разбиваем на массив только при наличии запятых
                                          const inputValue = e.target.value;
                                          const values = inputValue
                                            .split(',')
                                            .map((v) => v.trim());
                                          handleAttributeValuesChange(
                                            attr.id,
                                            values
                                          );
                                        }}
                                      />
                                    </div>
                                  )}
                                </div>
                              </div>
                            </div>
                          ))}
                        </div>

                        <div className="flex justify-end gap-3 mt-6">
                          <button
                            onClick={() => {
                              setActiveMode('none');
                              setHasVariants(false);
                              setSelectedAttributes({});
                            }}
                            className="btn btn-ghost"
                          >
                            {t('cancel') || 'Cancel'}
                          </button>
                          <button
                            onClick={() => {
                              const hasValues = Object.values(
                                selectedAttributes
                              ).every(
                                (attr: any) =>
                                  attr.values && attr.values.length > 0
                              );
                              if (
                                Object.keys(selectedAttributes).length > 0 &&
                                hasValues
                              ) {
                                setShowAttributeSelector(true);
                              }
                            }}
                            className="btn btn-primary"
                            disabled={
                              Object.keys(selectedAttributes).length === 0 ||
                              !Object.values(selectedAttributes).every(
                                (attr: any) =>
                                  attr.values && attr.values.length > 0
                              )
                            }
                          >
                            {t('generateVariants') || 'Generate Variants'}
                          </button>
                        </div>
                      </>
                    )}
                  </div>
                </div>
              ) : (
                <SimplifiedVariantGenerator
                  selectedAttributes={selectedAttributes}
                  categoryAttributes={Object.values(selectedAttributes)}
                  basePrice={state.aiData.price}
                  onGenerate={handleVariantsSave}
                  onCancel={() => {
                    setShowAttributeSelector(false);
                    setSelectedAttributes({});
                    setActiveMode('none');
                    setHasVariants(false);
                  }}
                  categoryId={state.aiData.categoryId}
                />
              )}
            </>
          )}

          {activeMode === 'advanced' && (
            <div className="alert alert-info">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                className="stroke-current shrink-0 w-6 h-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div>
                <div className="font-semibold">
                  {t('advancedMode') || 'Advanced Mode'}
                </div>
                <div className="text-sm">
                  {t('advancedModeNotAvailable') ||
                    'Advanced variant management will be available after product creation. For now, use Simple Mode to create basic variants.'}
                </div>
              </div>
              <button
                onClick={() => setActiveMode('simple')}
                className="btn btn-sm"
              >
                {t('switchToSimple') || 'Switch to Simple'}
              </button>
            </div>
          )}
        </div>
      )}

      {/* Actions */}
      <div className="flex justify-between gap-3">
        <button onClick={() => setView('enhance')} className="btn btn-outline">
          {t('back') || 'Back'}
        </button>

        <div className="flex gap-3">
          <button onClick={handleSkip} className="btn btn-ghost">
            {t('skipVariants') || 'Skip Variants'}
          </button>
          <button
            onClick={handleNext}
            className="btn btn-primary px-8"
            disabled={
              state.aiData.hasVariants && state.aiData.variants.length === 0
            }
          >
            {t('continueToPublish') || 'Continue to Publish'}
          </button>
        </div>
      </div>
    </div>
  );
}
