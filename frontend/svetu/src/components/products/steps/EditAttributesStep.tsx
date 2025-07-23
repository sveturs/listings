'use client';

import { useTranslations } from 'next-intl';
import { useEditProduct } from '@/contexts/EditProductContext';

interface EditAttributesStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function EditAttributesStep({
  onNext,
  onBack,
}: EditAttributesStepProps) {
  const t = useTranslations();
  const { state } = useEditProduct();

  return (
    <div className="space-y-6">
      {/* Заголовок */}
      <div className="text-center">
        <div className="w-16 h-16 bg-primary/20 rounded-full flex items-center justify-center mx-auto mb-4">
          <span className="text-2xl">⚙️</span>
        </div>
        <h3 className="text-2xl font-bold text-base-content mb-2">
          {t('storefronts.products.steps.attributes')}
        </h3>
        <p className="text-base-content/70">
          {t('storefronts.products.categoryAttributesDescription')}
        </p>
      </div>

      {/* Сообщение */}
      <div className="text-center py-8">
        <div className="bg-base-200 rounded-2xl p-6">
          <h4 className="text-lg font-semibold text-base-content mb-2">
            {t('storefronts.products.noAttributesTitle')}
          </h4>
          <p className="text-base-content/70">
            {state.category?.name
              ? t('storefronts.products.noAttributesForCategory')
              : t('storefronts.products.noAttributesMessage')}
          </p>
        </div>
      </div>

      {/* Кнопки навигации */}
      <div className="flex justify-between">
        <button onClick={onBack} className="btn btn-outline btn-lg">
          {t('common.back')}
        </button>
        <button onClick={onNext} className="btn btn-primary btn-lg">
          {t('common.continue')}
        </button>
      </div>
    </div>
  );
}
