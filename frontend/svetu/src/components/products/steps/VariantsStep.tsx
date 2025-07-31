'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateProduct } from '@/contexts/CreateProductContext';
import VariantManager from '@/components/Storefront/ProductVariants/VariantManager';
import AttributeSetup from '@/components/Storefront/ProductVariants/AttributeSetup';
import SimplifiedVariantGenerator from '@/components/products/SimplifiedVariantGenerator';

interface VariantsStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function VariantsStep({ onNext, onBack }: VariantsStepProps) {
  const t = useTranslations('storefronts.products');
  const { state, setHasVariants, setVariants, completeStep } =
    useCreateProduct();

  const [activeMode, setActiveMode] = useState<'none' | 'simple' | 'advanced'>(
    state.hasVariants
      ? state.variants.length > 0
        ? 'advanced'
        : 'simple'
      : 'none'
  );

  // Initialize from context
  useEffect(() => {
    if (state.hasVariants && state.variants.length > 0) {
      setActiveMode('advanced');
    }
  }, [state.hasVariants, state.variants.length]);

  const handleVariantToggle = (enabled: boolean) => {
    setHasVariants(enabled);
    if (!enabled) {
      setActiveMode('none');
      setVariants([]);
    } else {
      setActiveMode('simple');
    }
  };

  const handleModeChange = (mode: 'simple' | 'advanced') => {
    setActiveMode(mode);
  };

  const handleVariantsSave = (newVariants: any[]) => {
    // Convert to the format expected by the context
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
    completeStep(state.currentStep);
    onNext();
  };

  const handleSkipVariants = () => {
    setHasVariants(false);
    setVariants([]);
    completeStep(state.currentStep);
    onNext();
  };

  // Product data from context
  const productData = {
    id: 0, // Will be set after product creation
    category_id: state.productData.category_id || 0,
    title: state.productData.name || 'New Product',
    price: state.productData.price || 0,
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold text-base-content mb-2">
          {t('steps.variants')}
        </h2>
        <p className="text-base-content/70">{t('variantsStepDescription')}</p>
      </div>

      {/* Variant Toggle */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <h3 className="text-lg font-medium text-base-content">
                {t('enableVariants')}
              </h3>
              <p className="text-sm text-base-content/70">
                {t('enableVariantsDescription')}
              </p>
            </div>
            <div className="form-control">
              <label className="label cursor-pointer">
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={state.hasVariants}
                  onChange={(e) => handleVariantToggle(e.target.checked)}
                />
              </label>
            </div>
          </div>
        </div>
      </div>

      {/* Variant Configuration */}
      {state.hasVariants && (
        <div className="space-y-6">
          {/* Mode Selection */}
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h3 className="text-lg font-medium text-base-content mb-4">
                {t('variantMode')}
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
                  onClick={() => handleModeChange('simple')}
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
                        <h4 className="font-medium">{t('simpleVariants')}</h4>
                        <p className="text-sm text-base-content/70">
                          {t('simpleVariantsDescription')}
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
                  onClick={() => handleModeChange('advanced')}
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
                        <h4 className="font-medium">{t('advancedVariants')}</h4>
                        <p className="text-sm text-base-content/70">
                          {t('advancedVariantsDescription')}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Variant Configuration Content */}
          {activeMode === 'simple' && (
            <SimpleVariantConfig
              productData={productData}
              variants={state.variants}
              onVariantsChange={handleVariantsSave}
              selectedAttributes={state.attributes || {}}
              t={t}
            />
          )}

          {activeMode === 'advanced' && (
            <AdvancedVariantConfig
              productData={productData}
              variants={state.variants}
              onVariantsChange={handleVariantsSave}
              t={t}
            />
          )}
        </div>
      )}

      {/* Navigation */}
      <div className="flex justify-between items-center pt-6 border-t border-base-300">
        <button onClick={onBack} className="btn btn-outline">
          {t('back')}
        </button>

        <div className="flex space-x-2">
          {state.hasVariants && state.variants.length === 0 && (
            <button onClick={handleSkipVariants} className="btn btn-outline">
              {t('skipVariants')}
            </button>
          )}

          <button
            onClick={handleNext}
            className="btn btn-primary"
            disabled={state.hasVariants && state.variants.length === 0}
          >
            {t('continue')}
          </button>
        </div>
      </div>
    </div>
  );
}

// Simple Variant Configuration
interface SimpleVariantConfigProps {
  productData: any;
  variants: any[];
  onVariantsChange: (variants: any[]) => void;
  selectedAttributes: Record<number, any>;
  t: (key: string, params?: any) => string;
}

function SimpleVariantConfig({
  productData,
  variants: _variants,
  onVariantsChange,
  selectedAttributes,
  t,
}: SimpleVariantConfigProps) {
  const { state } = useCreateProduct();
  const [_basePrice] = useState(productData.price || 0);
  const [categoryAttributes, setCategoryAttributes] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  // Загружаем атрибуты категории
  useEffect(() => {
    const loadCategoryAttributes = async () => {
      if (!productData.category_id) {
        setLoading(false);
        return;
      }

      try {
        const response = await fetch(
          `/api/v1/marketplace/categories/${productData.category_id}/attributes`
        );
        if (response.ok) {
          const data = await response.json();
          setCategoryAttributes(data.data || data);
        }
      } catch (error) {
        console.error('Failed to load category attributes:', error);
      } finally {
        setLoading(false);
      }
    };

    loadCategoryAttributes();
  }, [productData.category_id]);

  const handleGenerateVariants = (generatedVariants: any[]) => {
    onVariantsChange(generatedVariants);
  };

  if (loading) {
    return (
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <div className="flex justify-center">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        </div>
      </div>
    );
  }

  if (categoryAttributes.length === 0) {
    return (
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <div className="text-center py-8">
            <p className="text-base-content/70">
              {t('noAttributesForCategory')}
            </p>
            <p className="text-sm text-base-content/50 mt-2">
              {t('selectDifferentCategory')}
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="card bg-base-100 shadow-sm">
      <div className="card-body">
        <h3 className="text-lg font-medium text-base-content mb-4">
          {t('generateVariants')}
        </h3>

        <SimplifiedVariantGenerator
          selectedAttributes={selectedAttributes}
          categoryAttributes={categoryAttributes}
          basePrice={productData.price || 0}
          onGenerate={handleGenerateVariants}
          onCancel={() => {}}
          categoryId={productData.category_id}
          categorySlug={state.category?.slug}
        />
      </div>
    </div>
  );
}

// Advanced Variant Configuration
interface AdvancedVariantConfigProps {
  productData: any;
  variants: any[];
  onVariantsChange: (variants: any[]) => void;
  t: (key: string, params?: any) => string;
}

function AdvancedVariantConfig({
  productData,
  variants,
  onVariantsChange: _onVariantsChange,
  t,
}: AdvancedVariantConfigProps) {
  const [showAttributeSetup, setShowAttributeSetup] = useState(false);

  const handleVariantManagerSave = () => {
    // This will be called when the user saves from VariantManager
    // The variants are already saved to the parent component
  };

  const handleAttributeSetupSave = (_attributes: any[]) => {
    setShowAttributeSetup(false);
    // Attributes are saved to the product
  };

  return (
    <div className="space-y-6">
      {/* Attribute Setup */}
      {!showAttributeSetup && variants.length === 0 && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <div className="text-center py-8">
              <div className="text-gray-500">
                <svg
                  className="mx-auto h-12 w-12 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10"
                  />
                </svg>
                <h3 className="mt-2 text-sm font-medium text-gray-900">
                  {t('noVariantsConfigured')}
                </h3>
                <p className="mt-1 text-sm text-gray-500">
                  {t('setupAttributesFirst')}
                </p>
              </div>
              <button
                onClick={() => setShowAttributeSetup(true)}
                className="btn btn-primary mt-4"
              >
                {t('setupAttributes')}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Attribute Setup Modal */}
      {showAttributeSetup && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <AttributeSetup
              productId={0} // Temporary ID for new products
              categoryId={productData.category_id}
              onSave={handleAttributeSetupSave}
              onCancel={() => setShowAttributeSetup(false)}
            />
          </div>
        </div>
      )}

      {/* Variant Manager */}
      {variants.length > 0 && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <VariantManager
              productId={0} // Temporary ID for new products
              storefrontId={0} // Will be determined from context
              onSave={handleVariantManagerSave}
              onCancel={() => {}}
            />
          </div>
        </div>
      )}
    </div>
  );
}
