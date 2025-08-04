'use client';

import { useTranslations } from 'next-intl';

interface Step {
  id: string;
  label: string;
}

interface StepWizardProps {
  steps: Step[];
  currentStep: number;
  onStepClick: (step: number) => void;
}

export default function StepWizard({
  steps,
  currentStep,
  onStepClick,
}: StepWizardProps) {
  const t = useTranslations('create_listing');

  return (
    <div className="w-full mb-8">
      {/* Прогресс-бар */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex-1 bg-base-300 h-2 rounded-full overflow-hidden">
          <div
            className="h-full bg-primary transition-all duration-300 ease-out"
            style={{ width: `${((currentStep + 1) / steps.length) * 100}%` }}
          />
        </div>
        <span className="ml-4 text-sm font-medium text-base-content/70">
          {currentStep + 1} / {steps.length}
        </span>
      </div>

      {/* Шаги для десктопа */}
      <div className="hidden sm:block">
        <div className="flex justify-between items-center">
          {steps.map((step, index) => (
            <div key={step.id} className="flex flex-col items-center relative">
              {/* Линия между шагами */}
              {index < steps.length - 1 && (
                <div
                  className={`absolute top-4 left-8 w-full h-0.5 -z-10 ${
                    index < currentStep ? 'bg-primary' : 'bg-base-300'
                  }`}
                  style={{ width: 'calc(100vw / ' + steps.length + ' - 2rem)' }}
                />
              )}

              {/* Номер шага */}
              <button
                onClick={() => onStepClick(index)}
                disabled={index > currentStep}
                className={`
                  w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium
                  transition-all duration-200 border-2
                  ${
                    index <= currentStep
                      ? 'bg-primary text-primary-content border-primary'
                      : 'bg-base-100 text-base-content/50 border-base-300'
                  }
                  ${index < currentStep ? 'hover:bg-primary-focus cursor-pointer' : ''}
                  ${index > currentStep ? 'cursor-not-allowed' : ''}
                `}
              >
                {index < currentStep ? (
                  <svg
                    className="w-4 h-4"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clipRule="evenodd"
                    />
                  </svg>
                ) : (
                  index + 1
                )}
              </button>

              {/* Название шага */}
              <span
                className={`
                mt-2 text-xs text-center max-w-20
                ${index <= currentStep ? 'text-base-content font-medium' : 'text-base-content/50'}
              `}
              >
                {t(step.label)}
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* Мобильная версия - только текущий шаг */}
      <div className="sm:hidden">
        <div className="text-center">
          <div className="inline-flex items-center gap-2 bg-base-200 rounded-full px-4 py-2">
            <div
              className={`
              w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium
              bg-primary text-primary-content
            `}
            >
              {currentStep + 1}
            </div>
            <span className="text-sm font-medium">
              {t(steps[currentStep].label)}
            </span>
          </div>
        </div>
      </div>

      {/* Региональные подсказки */}
      <div className="mt-4 text-center">
        {currentStep === 0 && (
          <p className="text-xs text-base-content/60">
            💡 {t('category')}
          </p>
        )}
        {currentStep === 2 && (
          <p className="text-xs text-base-content/60">
            🤝 {t('trust')}
          </p>
        )}
        {currentStep === 6 && (
          <p className="text-xs text-base-content/60">
            📦 {t('delivery')}
          </p>
        )}
      </div>
    </div>
  );
}
