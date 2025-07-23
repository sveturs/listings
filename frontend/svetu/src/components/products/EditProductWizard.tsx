'use client';

import { useTranslations } from 'next-intl';
import { useEditProduct } from '@/contexts/EditProductContext';
import EditCategoryStep from './steps/EditCategoryStep';
import EditBasicInfoStep from './steps/EditBasicInfoStep';
import EditLocationStep from './steps/EditLocationStep';
import EditAttributesStep from './steps/EditAttributesStep';
import EditPhotosStep from './steps/EditPhotosStep';
import EditPreviewStep from './steps/EditPreviewStep';

interface EditProductWizardProps {
  storefrontSlug: string;
  productId: number;
}

const STEPS = [
  { id: 'category', icon: '🏷️' },
  { id: 'basic', icon: '📝' },
  { id: 'location', icon: '📍' },
  { id: 'attributes', icon: '⚙️' },
  { id: 'photos', icon: '📸' },
  { id: 'preview', icon: '👁️' },
];

export default function EditProductWizard({
  storefrontSlug,
  productId,
}: EditProductWizardProps) {
  const t = useTranslations();
  const { state, goToStep, nextStep, prevStep } = useEditProduct();

  const renderStepContent = () => {
    switch (state.currentStep) {
      case 0:
        return <EditCategoryStep onNext={nextStep} />;
      case 1:
        return <EditBasicInfoStep onNext={nextStep} onBack={prevStep} />;
      case 2:
        return <EditLocationStep onNext={nextStep} onBack={prevStep} />;
      case 3:
        return <EditAttributesStep onNext={nextStep} onBack={prevStep} />;
      case 4:
        return <EditPhotosStep onNext={nextStep} onBack={prevStep} />;
      case 5:
        return (
          <EditPreviewStep
            onBack={prevStep}
            storefrontSlug={storefrontSlug}
            productId={productId}
          />
        );
      default:
        return <EditCategoryStep onNext={nextStep} />;
    }
  };

  if (state.isLoading) {
    return (
      <div className="container mx-auto px-4 py-8 flex justify-center">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="w-full">
      {/* Индикатор несохраненных изменений */}
      {state.hasUnsavedChanges && (
        <div className="bg-warning/10 border border-warning rounded-2xl p-4 mb-6">
          <div className="flex items-center gap-3">
            <div className="w-2 h-2 bg-warning rounded-full animate-pulse" />
            <p className="text-warning font-medium">
              {t('storefronts.products.unsavedChanges')}
            </p>
          </div>
        </div>
      )}

      {/* Мобильная версия - компактный заголовок */}
      <div className="lg:hidden mb-6">
        <div className="bg-base-100 rounded-2xl shadow-lg p-4">
          <div className="flex items-center gap-4">
            <div
              className={`
                w-12 h-12 rounded-full flex items-center justify-center text-xl border-2
                bg-primary text-primary-content border-primary
              `}
            >
              {STEPS[state.currentStep].icon}
            </div>
            <div className="flex-1">
              <h2 className="text-lg font-bold text-base-content">
                {t(`storefronts.products.steps.${STEPS[state.currentStep].id}`)}
              </h2>
              <div className="flex items-center gap-2 mt-1">
                <div className="flex-1 bg-base-300 rounded-full h-1">
                  <div
                    className="bg-primary h-1 rounded-full transition-all duration-300"
                    style={{
                      width: `${((state.currentStep + 1) / STEPS.length) * 100}%`,
                    }}
                  />
                </div>
                <span className="text-xs text-base-content/60 min-w-0">
                  {state.currentStep + 1}/{STEPS.length}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Десктопная версия */}
      <div className="hidden lg:block mb-8">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-2xl font-bold text-base-content">
            {t(`storefronts.products.steps.${STEPS[state.currentStep].id}`)}
          </h2>
          <div className="text-sm text-base-content/60">
            {state.currentStep + 1} / {STEPS.length}
          </div>
        </div>

        {/* Прогресс бар */}
        <div className="w-full bg-base-300 rounded-full h-2 mb-6">
          <div
            className="bg-primary h-2 rounded-full transition-all duration-300"
            style={{
              width: `${((state.currentStep + 1) / STEPS.length) * 100}%`,
            }}
          />
        </div>

        {/* Шаги */}
        <div className="flex justify-between">
          {STEPS.map((step, index) => {
            const isCompleted = state.completedSteps.has(index);
            const isCurrent = state.currentStep === index;
            const canAccess =
              index <= Math.max(...state.completedSteps, state.currentStep);

            return (
              <button
                key={step.id}
                onClick={() => canAccess && goToStep(index)}
                disabled={!canAccess}
                className={`
                  flex flex-col items-center gap-2 p-3 rounded-lg transition-all
                  ${isCurrent ? 'text-primary' : ''}
                  ${isCompleted ? 'text-success' : ''}
                  ${!canAccess ? 'text-base-content/30 cursor-not-allowed' : 'hover:bg-base-200'}
                `}
              >
                <div
                  className={`
                    w-12 h-12 rounded-full flex items-center justify-center text-lg border-2
                    ${
                      isCompleted
                        ? 'bg-success text-success-content border-success'
                        : isCurrent
                          ? 'bg-primary text-primary-content border-primary'
                          : canAccess
                            ? 'bg-base-100 border-base-300'
                            : 'bg-base-200 border-base-300'
                    }
                  `}
                >
                  {isCompleted ? '✓' : step.icon}
                </div>
                <span className="text-sm font-medium">
                  {t(`storefronts.products.steps.${step.id}`)}
                </span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Содержимое шага */}
      <div className="bg-base-100 rounded-2xl shadow-xl p-4 lg:p-6">
        {renderStepContent()}
      </div>
    </div>
  );
}
