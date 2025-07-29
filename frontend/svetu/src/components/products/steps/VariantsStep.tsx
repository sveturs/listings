'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateProduct } from '@/contexts/CreateProductContext';
import VariantManager from '@/components/Storefront/ProductVariants/VariantManager';
import AttributeSetup from '@/components/Storefront/ProductVariants/AttributeSetup';
import VariantGenerator from '@/components/Storefront/ProductVariants/VariantGenerator';

interface VariantsStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function VariantsStep({ onNext, onBack }: VariantsStepProps) {
  const t = useTranslations('storefronts.products');
  const { state, setState } = useCreateProduct();

  const [activeMode, setActiveMode] = useState<'none' | 'simple' | 'advanced'>(
    'none'
  );
  const [hasVariants, setHasVariants] = useState(false);
  const [variants, setVariants] = useState<any[]>([]);

  // Initialize variants state from context
  useEffect(() => {
    if (state.formData.variants && state.formData.variants.length > 0) {
      setVariants(state.formData.variants);
      setHasVariants(true);
      setActiveMode('advanced');
    }
  }, [state.formData.variants]);

  const handleVariantToggle = (enabled: boolean) => {
    setHasVariants(enabled);
    if (!enabled) {
      setActiveMode('none');
      setVariants([]);
      // Clear variants from form data
      setState((prev) => ({
        ...prev,
        formData: {
          ...prev.formData,
          variants: [],
        },
      }));
    } else {
      setActiveMode('simple');
    }
  };

  const handleModeChange = (mode: 'simple' | 'advanced') => {
    setActiveMode(mode);
  };

  const handleVariantsSave = (newVariants: any[]) => {
    setVariants(newVariants);
    // Save to context
    setState((prev) => ({
      ...prev,
      formData: {
        ...prev.formData,
        variants: newVariants,
      },
    }));
  };

  const handleNext = () => {
    // Mark step as completed
    setState((prev) => ({
      ...prev,
      completedSteps: new Set([...prev.completedSteps, prev.currentStep]),
    }));
    onNext();
  };

  const handleSkipVariants = () => {
    setHasVariants(false);
    setVariants([]);
    setState((prev) => ({
      ...prev,
      formData: {
        ...prev.formData,
        variants: [],
      },
      completedSteps: new Set([...prev.completedSteps, prev.currentStep]),
    }));
    onNext();
  };

  // Mock product data for the variant manager
  const mockProduct = {
    id: 0, // Will be set after product creation
    category_id: state.formData.category?.id || 0,
    title: state.formData.title || 'New Product',
    price: state.formData.price || 0,
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
                  checked={hasVariants}
                  onChange={(e) => handleVariantToggle(e.target.checked)}
                />
              </label>
            </div>
          </div>
        </div>
      </div>

      {/* Variant Configuration */}
      {hasVariants && (
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
              productData={mockProduct}
              variants={variants}
              onVariantsChange={handleVariantsSave}
              t={t}
            />
          )}

          {activeMode === 'advanced' && (
            <AdvancedVariantConfig
              productData={mockProduct}
              variants={variants}
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
          {hasVariants && variants.length === 0 && (
            <button onClick={handleSkipVariants} className="btn btn-outline">
              {t('skipVariants')}
            </button>
          )}

          <button
            onClick={handleNext}
            className="btn btn-primary"
            disabled={hasVariants && variants.length === 0}
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
  t: (key: string, params?: any) => string;
}

function SimpleVariantConfig({
  productData,
  variants,
  onVariantsChange,
  t,
}: SimpleVariantConfigProps) {
  const [basePrice, setBasePrice] = useState(productData.price || 0);

  const handleGenerateVariants = (generatedVariants: any[]) => {
    onVariantsChange(generatedVariants);
  };

  return (
    <div className="card bg-base-100 shadow-sm">
      <div className="card-body">
        <h3 className="text-lg font-medium text-base-content mb-4">
          {t('generateVariants')}
        </h3>

        <VariantGenerator
          productId={0} // Temporary ID for new products
          basePrice={basePrice}
          onGenerate={handleGenerateVariants}
          onCancel={() => {}}
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
  onVariantsChange,
  t,
}: AdvancedVariantConfigProps) {
  const [showAttributeSetup, setShowAttributeSetup] = useState(false);

  const handleVariantManagerSave = () => {
    // This will be called when the user saves from VariantManager
    // The variants are already saved to the parent component
  };

  const handleAttributeSetupSave = (attributes: any[]) => {
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
