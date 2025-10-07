'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  analyzeImportFile,
  setCurrentAnalysisStep,
  setAnalysisFile,
  setAnalysisFileType,
  addApprovedMapping,
  setCustomMapping,
  toggleSelectedAttribute,
  toggleApprovedVariantGroup,
  setSelectedAttributes,
  setApprovedVariantGroups,
  clearAnalysis,
} from '@/store/slices/importSlice';
import CategoryMappingStep from './CategoryMappingStep';
import AttributeMappingStep from './AttributeMappingStep';
import VariantDetectionStep from './VariantDetectionStep';
import type { CategoryMapping } from '@/types/import';

interface ImportAnalysisWizardProps {
  storefrontId: number;
  storefrontSlug?: string;
  onClose?: () => void;
  onSuccess?: (jobId: number) => void;
  onSwitchToClassic?: () => void;
}

const WIZARD_STEPS = [
  'upload', // Step 0: File upload
  'analyzing', // Step 1: Auto analysis (progress indicator)
  'categories', // Step 2: Category mapping
  'attributes', // Step 3: Attribute selection
  'variants', // Step 4: Variant detection
  'summary', // Step 5: Summary and confirm
] as const;

type WizardStep = (typeof WIZARD_STEPS)[number];

export default function ImportAnalysisWizard({
  storefrontId,
  storefrontSlug,
  onClose,
  onSuccess,
  onSwitchToClassic,
}: ImportAnalysisWizardProps) {
  const t = useTranslations('storefronts.import.wizard');
  const dispatch = useAppDispatch();

  const {
    analysisFile,
    analysisFileType,
    categoryAnalysis,
    attributeAnalysis,
    variantDetection,
    isAnalyzing,
    analysisError,
    analysisProgress,
    currentAnalysisStep,
    approvedMappings,
    customMappings,
    selectedAttributes,
    approvedVariantGroups,
  } = useAppSelector((state) => state.import);

  const [currentStep, setCurrentStep] = useState<number>(0);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);

  // Mock available categories - should be fetched from API
  const [availableCategories] = useState<
    Array<{ id: number; name: string; parent?: string }>
  >([
    { id: 1, name: 'Electronics', parent: undefined },
    { id: 2, name: 'Smartphones', parent: 'Electronics' },
    { id: 3, name: 'Laptops', parent: 'Electronics' },
    { id: 4, name: 'Clothing', parent: undefined },
    { id: 5, name: 'Men', parent: 'Clothing' },
    { id: 6, name: 'Women', parent: 'Clothing' },
  ]);

  useEffect(() => {
    return () => {
      // Cleanup on unmount
      dispatch(clearAnalysis());
    };
  }, [dispatch]);

  const handleFileSelect = (file: File) => {
    setSelectedFile(file);
    dispatch(setAnalysisFile(file));

    // Detect file type
    const extension = file.name.split('.').pop()?.toLowerCase();
    let fileType: 'xml' | 'csv' | 'zip' | '' = '';
    if (extension === 'xml') fileType = 'xml';
    else if (extension === 'csv') fileType = 'csv';
    else if (extension === 'zip') fileType = 'zip';
    dispatch(setAnalysisFileType(fileType));
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const files = Array.from(e.dataTransfer.files);
    if (files.length > 0) {
      handleFileSelect(files[0]);
    }
  };

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleFileSelect(files[0]);
    }
  };

  const handleStartAnalysis = async () => {
    if (!selectedFile || !analysisFileType) {
      alert(t('errors.fileRequired'));
      return;
    }

    setCurrentStep(1); // Move to analyzing step
    dispatch(setCurrentAnalysisStep(1));

    try {
      await dispatch(
        analyzeImportFile({
          storefrontId,
          file: selectedFile,
          fileType: analysisFileType as 'xml' | 'csv' | 'zip',
        })
      ).unwrap();

      // Analysis complete, move to category mapping
      setCurrentStep(2);
      dispatch(setCurrentAnalysisStep(2));
    } catch (error) {
      console.error('Analysis failed:', error);
      // Stay on analyzing step to show error
    }
  };

  const handleCategoryMappingChange = (
    externalCategory: string,
    internalCategoryId: number | null
  ) => {
    if (internalCategoryId !== null) {
      dispatch(setCustomMapping({ externalCategory, internalCategoryId }));
    }
  };

  const handleApproveMapping = (externalCategory: string) => {
    const mapping = categoryAnalysis?.mappings.find(
      (m) => m.external_category === externalCategory
    );
    if (mapping) {
      dispatch(
        addApprovedMapping({
          ...mapping,
          is_approved: true,
        })
      );
    }
  };

  const handleRequestNewCategory = (
    externalCategory: string,
    reasoning: string
  ) => {
    // TODO: Implement category proposal submission
    console.log('Request new category:', externalCategory, reasoning);
  };

  const handleAttributeToggle = (attributeName: string) => {
    dispatch(toggleSelectedAttribute(attributeName));
  };

  const handleVariantGroupToggle = (baseName: string) => {
    dispatch(toggleApprovedVariantGroup(baseName));
  };

  const handleNext = () => {
    if (currentStep < WIZARD_STEPS.length - 1) {
      const nextStep = currentStep + 1;
      setCurrentStep(nextStep);
      dispatch(setCurrentAnalysisStep(nextStep));
    }
  };

  const handleBack = () => {
    if (currentStep > 0) {
      const prevStep = currentStep - 1;
      setCurrentStep(prevStep);
      dispatch(setCurrentAnalysisStep(prevStep));
    }
  };

  const handleComplete = () => {
    // TODO: Start actual import with all the selected options
    console.log('Starting import with:', {
      file: selectedFile,
      approvedMappings,
      customMappings,
      selectedAttributes,
      approvedVariantGroups,
    });

    // For now, close the wizard. In the future, this will trigger actual import
    // and call onSuccess with job ID
    onClose?.();
  };

  const renderStepContent = () => {
    const stepName = WIZARD_STEPS[currentStep];

    switch (stepName) {
      case 'upload':
        return (
          <div className="space-y-4">
            <h2 className="text-2xl font-bold text-gray-900">
              {t('steps.upload.title')}
            </h2>
            <p className="text-gray-600">{t('steps.upload.description')}</p>

            {/* Drag and drop zone */}
            <div
              className={`
                border-2 border-dashed rounded-lg p-12 text-center
                transition-colors cursor-pointer
                ${
                  isDragging
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-300 hover:border-gray-400'
                }
              `}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              onClick={() => document.getElementById('fileInput')?.click()}
            >
              <input
                id="fileInput"
                type="file"
                accept=".xml,.csv,.zip"
                onChange={handleFileInputChange}
                className="hidden"
              />

              <div className="space-y-2">
                <div className="text-4xl">📁</div>
                <p className="text-lg font-medium text-gray-900">
                  {selectedFile
                    ? selectedFile.name
                    : t('steps.upload.dragDrop')}
                </p>
                <p className="text-sm text-gray-500">
                  {t('steps.upload.supportedFormats')}
                </p>
              </div>
            </div>

            {selectedFile && (
              <div className="flex items-center justify-between p-4 bg-green-50 border border-green-200 rounded-lg">
                <div className="flex items-center gap-2">
                  <span className="text-green-600">✓</span>
                  <span className="font-medium">{selectedFile.name}</span>
                  <span className="text-sm text-gray-500">
                    ({(selectedFile.size / 1024 / 1024).toFixed(2)} MB)
                  </span>
                </div>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setSelectedFile(null);
                    dispatch(setAnalysisFile(null));
                  }}
                  className="text-red-600 hover:text-red-800"
                >
                  {t('steps.upload.remove')}
                </button>
              </div>
            )}

            <div className="flex justify-between mt-8">
              <div className="flex gap-4">
                <button
                  onClick={onClose}
                  className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                >
                  {t('buttons.cancel')}
                </button>
                {onSwitchToClassic && (
                  <button
                    onClick={onSwitchToClassic}
                    className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                  >
                    Switch to Classic Import
                  </button>
                )}
              </div>
              <button
                onClick={handleStartAnalysis}
                disabled={!selectedFile}
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {t('buttons.startAnalysis')}
              </button>
            </div>
          </div>
        );

      case 'analyzing':
        return (
          <div className="space-y-6">
            <h2 className="text-2xl font-bold text-gray-900">
              {t('steps.analyzing.title')}
            </h2>
            <p className="text-gray-600">{t('steps.analyzing.description')}</p>

            {/* Progress bar */}
            <div className="space-y-2">
              <div className="flex justify-between text-sm text-gray-600">
                <span>{t('steps.analyzing.progress')}</span>
                <span>{analysisProgress}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                  style={{ width: `${analysisProgress}%` }}
                />
              </div>
            </div>

            {/* Analysis stages */}
            <div className="space-y-3">
              <div
                className={`flex items-center gap-3 ${analysisProgress >= 33 ? 'text-green-600' : 'text-gray-400'}`}
              >
                <span className="text-xl">
                  {analysisProgress >= 33 ? '✓' : '⏳'}
                </span>
                <span>{t('steps.analyzing.stages.categories')}</span>
              </div>
              <div
                className={`flex items-center gap-3 ${analysisProgress >= 66 ? 'text-green-600' : 'text-gray-400'}`}
              >
                <span className="text-xl">
                  {analysisProgress >= 66 ? '✓' : '⏳'}
                </span>
                <span>{t('steps.analyzing.stages.attributes')}</span>
              </div>
              <div
                className={`flex items-center gap-3 ${analysisProgress >= 100 ? 'text-green-600' : 'text-gray-400'}`}
              >
                <span className="text-xl">
                  {analysisProgress >= 100 ? '✓' : '⏳'}
                </span>
                <span>{t('steps.analyzing.stages.variants')}</span>
              </div>
            </div>

            {analysisError && (
              <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
                <p className="text-red-800">{analysisError}</p>
                <button
                  onClick={() => setCurrentStep(0)}
                  className="mt-2 text-red-600 hover:text-red-800 underline"
                >
                  {t('buttons.backToUpload')}
                </button>
              </div>
            )}
          </div>
        );

      case 'categories':
        return (
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">
                  {t('steps.categories.title')}
                </h2>
                <p className="text-gray-600">
                  {t('steps.categories.description')}
                </p>
              </div>
              {categoryAnalysis && (
                <div className="text-sm text-gray-600">
                  {categoryAnalysis.total_categories} {t('categories.total')}
                </div>
              )}
            </div>

            {categoryAnalysis && categoryAnalysis.mapping_quality && (
              <>
                {/* Quality summary */}
                <div className="grid grid-cols-4 gap-4 mb-6">
                  <div className="p-4 bg-green-50 border border-green-200 rounded-lg">
                    <div className="text-2xl font-bold text-green-700">
                      {categoryAnalysis.mapping_quality.high_confidence?.length || 0}
                    </div>
                    <div className="text-sm text-green-600">
                      {t('categories.highConfidence')}
                    </div>
                  </div>
                  <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
                    <div className="text-2xl font-bold text-yellow-700">
                      {categoryAnalysis.mapping_quality.medium_confidence?.length || 0}
                    </div>
                    <div className="text-sm text-yellow-600">
                      {t('categories.mediumConfidence')}
                    </div>
                  </div>
                  <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
                    <div className="text-2xl font-bold text-red-700">
                      {categoryAnalysis.mapping_quality.low_confidence?.length || 0}
                    </div>
                    <div className="text-sm text-red-600">
                      {t('categories.lowConfidence')}
                    </div>
                  </div>
                  <div className="p-4 bg-gray-50 border border-gray-200 rounded-lg">
                    <div className="text-2xl font-bold text-gray-700">
                      {categoryAnalysis.total_categories - (categoryAnalysis.mapping_quality.total_mapped || 0)}
                    </div>
                    <div className="text-sm text-gray-600">
                      {t('categories.unmapped')}
                    </div>
                  </div>
                </div>

                <CategoryMappingStep
                  mappings={
                    categoryAnalysis.mapping_suggestions
                      ? Object.values(categoryAnalysis.mapping_suggestions)
                      : []
                  }
                  onMappingChange={handleCategoryMappingChange}
                  onApproveMapping={handleApproveMapping}
                  onRequestNewCategory={handleRequestNewCategory}
                  availableCategories={availableCategories}
                  isLoading={isAnalyzing}
                />
              </>
            )}

            <div className="flex justify-between gap-4 mt-8">
              <button
                onClick={handleBack}
                className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                {t('buttons.back')}
              </button>
              <button
                onClick={handleNext}
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                {t('buttons.next')}
              </button>
            </div>
          </div>
        );

      case 'attributes':
        return (
          <div className="space-y-4">
            <div>
              <h2 className="text-2xl font-bold text-gray-900">
                {t('steps.attributes.title')}
              </h2>
              <p className="text-gray-600">
                {t('steps.attributes.description')}
              </p>
            </div>

            {attributeAnalysis && (
              <AttributeMappingStep
                attributes={attributeAnalysis.detected_attributes || []}
                selectedAttributes={selectedAttributes}
                onToggleAttribute={handleAttributeToggle}
                onSelectAll={() => {
                  dispatch(
                    setSelectedAttributes(
                      (attributeAnalysis.detected_attributes || []).map(
                        (a) => a.name
                      )
                    )
                  );
                }}
                onDeselectAll={() => {
                  dispatch(setSelectedAttributes([]));
                }}
                isLoading={isAnalyzing}
              />
            )}

            <div className="flex justify-between gap-4 mt-8">
              <button
                onClick={handleBack}
                className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                {t('buttons.back')}
              </button>
              <button
                onClick={handleNext}
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                {t('buttons.next')}
              </button>
            </div>
          </div>
        );

      case 'variants':
        return (
          <div className="space-y-4">
            <div>
              <h2 className="text-2xl font-bold text-gray-900">
                {t('steps.variants.title')}
              </h2>
              <p className="text-gray-600">{t('steps.variants.description')}</p>
            </div>

            {variantDetection && (
              <>
                {/* Summary */}
                <div className="grid grid-cols-3 gap-4 mb-6">
                  <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
                    <div className="text-2xl font-bold text-blue-700">
                      {variantDetection.total_groups}
                    </div>
                    <div className="text-sm text-blue-600">
                      {t('variants.totalGroups')}
                    </div>
                  </div>
                  <div className="p-4 bg-green-50 border border-green-200 rounded-lg">
                    <div className="text-2xl font-bold text-green-700">
                      {variantDetection.variant_groups.reduce(
                        (acc, g) => acc + g.variant_count,
                        0
                      )}
                    </div>
                    <div className="text-sm text-green-600">
                      {t('variants.totalVariants')}
                    </div>
                  </div>
                  <div className="p-4 bg-gray-50 border border-gray-200 rounded-lg">
                    <div className="text-2xl font-bold text-gray-700">
                      {variantDetection.ungrouped_products}
                    </div>
                    <div className="text-sm text-gray-600">
                      {t('variants.ungrouped')}
                    </div>
                  </div>
                </div>

                <VariantDetectionStep
                  variantGroups={variantDetection.variant_groups || []}
                  approvedGroups={approvedVariantGroups}
                  onToggleGroup={handleVariantGroupToggle}
                  onApproveAll={() => {
                    dispatch(
                      setApprovedVariantGroups(
                        (variantDetection.variant_groups || []).map(
                          (g) => g.base_name
                        )
                      )
                    );
                  }}
                  onRejectAll={() => {
                    dispatch(setApprovedVariantGroups([]));
                  }}
                  isLoading={isAnalyzing}
                />
              </>
            )}

            <div className="flex justify-between gap-4 mt-8">
              <button
                onClick={handleBack}
                className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                {t('buttons.back')}
              </button>
              <button
                onClick={handleNext}
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                {t('buttons.next')}
              </button>
            </div>
          </div>
        );

      case 'summary':
        return (
          <div className="space-y-6">
            <h2 className="text-2xl font-bold text-gray-900">
              {t('steps.summary.title')}
            </h2>
            <p className="text-gray-600">{t('steps.summary.description')}</p>

            {/* Summary cards */}
            <div className="space-y-4">
              <div className="p-6 border border-gray-200 rounded-lg">
                <h3 className="font-semibold text-lg mb-4">
                  {t('summary.file')}
                </h3>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-gray-600">
                      {t('summary.fileName')}:
                    </span>
                    <span className="font-medium">{selectedFile?.name}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">
                      {t('summary.fileSize')}:
                    </span>
                    <span className="font-medium">
                      {selectedFile
                        ? (selectedFile.size / 1024 / 1024).toFixed(2)
                        : 0}{' '}
                      MB
                    </span>
                  </div>
                </div>
              </div>

              <div className="p-6 border border-gray-200 rounded-lg">
                <h3 className="font-semibold text-lg mb-4">
                  {t('summary.categories')}
                </h3>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-gray-600">
                      {t('summary.totalCategories')}:
                    </span>
                    <span className="font-medium">
                      {categoryAnalysis?.total_categories || 0}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">
                      {t('summary.approvedMappings')}:
                    </span>
                    <span className="font-medium">
                      {approvedMappings.length}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">
                      {t('summary.customMappings')}:
                    </span>
                    <span className="font-medium">
                      {Object.keys(customMappings).length}
                    </span>
                  </div>
                </div>
              </div>

              <div className="p-6 border border-gray-200 rounded-lg">
                <h3 className="font-semibold text-lg mb-4">
                  {t('summary.attributes')}
                </h3>
                <div className="flex justify-between">
                  <span className="text-gray-600">
                    {t('summary.selectedAttributes')}:
                  </span>
                  <span className="font-medium">
                    {selectedAttributes.length}
                  </span>
                </div>
              </div>

              <div className="p-6 border border-gray-200 rounded-lg">
                <h3 className="font-semibold text-lg mb-4">
                  {t('summary.variants')}
                </h3>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-gray-600">
                      {t('summary.approvedGroups')}:
                    </span>
                    <span className="font-medium">
                      {approvedVariantGroups.length}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">
                      {t('summary.totalVariants')}:
                    </span>
                    <span className="font-medium">
                      {variantDetection?.variant_groups
                        .filter((g) =>
                          approvedVariantGroups.includes(g.base_name)
                        )
                        .reduce((acc, g) => acc + g.variant_count, 0) || 0}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <div className="flex justify-between gap-4 mt-8">
              <button
                onClick={handleBack}
                className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                {t('buttons.back')}
              </button>
              <button
                onClick={handleComplete}
                className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
              >
                {t('buttons.startImport')}
              </button>
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="max-w-6xl mx-auto p-6">
      {/* Progress indicator */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          {WIZARD_STEPS.map((step, index) => (
            <React.Fragment key={step}>
              <div
                className={`
                  flex items-center justify-center w-10 h-10 rounded-full
                  ${
                    index <= currentStep
                      ? 'bg-blue-600 text-white'
                      : 'bg-gray-200 text-gray-500'
                  }
                  ${index === currentStep ? 'ring-4 ring-blue-200' : ''}
                `}
              >
                {index < currentStep ? '✓' : index + 1}
              </div>
              {index < WIZARD_STEPS.length - 1 && (
                <div
                  className={`
                    flex-1 h-1 mx-2
                    ${index < currentStep ? 'bg-blue-600' : 'bg-gray-200'}
                  `}
                />
              )}
            </React.Fragment>
          ))}
        </div>
        <div className="flex items-center justify-between mt-2">
          {WIZARD_STEPS.map((step, index) => (
            <div
              key={step}
              className="text-xs text-gray-600"
              style={{ width: '10%', textAlign: 'center' }}
            >
              {t(`steps.${step}.name`)}
            </div>
          ))}
        </div>
      </div>

      {/* Step content */}
      <div className="bg-white rounded-lg shadow-lg p-8">
        {renderStepContent()}
      </div>
    </div>
  );
}
