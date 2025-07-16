'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useAuthContext } from '@/contexts/AuthContext';
import StepWizard from '@/components/create-listing/StepWizard';
import CategorySelectionStep from '@/components/create-listing/steps/CategorySelectionStep';
import BasicInfoStep from '@/components/create-listing/steps/BasicInfoStep';
import TrustSetupStep from '@/components/create-listing/steps/TrustSetupStep';
import AttributesStep from '@/components/create-listing/steps/AttributesStep';
import PhotosStep from '@/components/create-listing/steps/PhotosStep';
import LocationStep from '@/components/create-listing/steps/LocationStep';
import PaymentDeliveryStep from '@/components/create-listing/steps/PaymentDeliveryStep';
import PreviewPublishStep from '@/components/create-listing/steps/PreviewPublishStep';
import { CreateListingProvider } from '@/contexts/CreateListingContext';
import {
  DraftStatus,
  DraftIndicator,
  OfflineIndicator,
} from '@/components/DraftStatus';
import { DraftsModal } from '@/components/DraftsModal';
import { toast } from '@/utils/toast';

const steps = [
  { id: 'category', label: 'create_listing.steps.category' },
  { id: 'basic_info', label: 'create_listing.steps.basic_info' },
  { id: 'trust_setup', label: 'create_listing.steps.trust_setup' }, // Система доверия
  { id: 'attributes', label: 'create_listing.steps.attributes' },
  { id: 'photos', label: 'create_listing.steps.photos' },
  { id: 'location', label: 'create_listing.steps.location' },
  { id: 'payment_delivery', label: 'create_listing.steps.payment_delivery' }, // Лична предаја + наложенный платеж
  { id: 'preview', label: 'create_listing.steps.preview' },
];

export default function CreateListingPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const t = useTranslations();
  const { user, isLoading: authLoading } = useAuthContext();
  const [currentStep, setCurrentStep] = useState(0);
  const [showDraftsModal, setShowDraftsModal] = useState(false);
  const [isClient, setIsClient] = useState(false);

  // Получаем ID черновика из URL параметров
  const draftId = searchParams?.get('draft') || undefined;

  useEffect(() => {
    setIsClient(true);
  }, []);

  useEffect(() => {
    if (isClient && !authLoading && !user) {
      toast.error(t('create_listing.auth_required'));
      router.push('/');
    }
  }, [user, authLoading, router, t, isClient]);

  const handleStepChange = (newStep: number) => {
    setCurrentStep(newStep);
    if (typeof window !== 'undefined') {
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
  };

  const handleComplete = () => {
    toast.success(t('create_listing.success'));
    router.push('/profile');
  };

  if (!isClient || authLoading) {
    return (
      <div className="container mx-auto px-2 sm:px-4 py-4 sm:py-8 min-h-screen">
        <div className="text-center mb-6 sm:mb-8">
          <div className="loading loading-spinner loading-lg"></div>
        </div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  const renderStep = () => {
    switch (currentStep) {
      case 0:
        return <CategorySelectionStep onNext={() => handleStepChange(1)} />;
      case 1:
        return (
          <BasicInfoStep
            onNext={() => handleStepChange(2)}
            onBack={() => handleStepChange(0)}
          />
        );
      case 2:
        return (
          <TrustSetupStep
            onNext={() => handleStepChange(3)}
            onBack={() => handleStepChange(1)}
          />
        );
      case 3:
        return (
          <AttributesStep
            onNext={() => handleStepChange(4)}
            onBack={() => handleStepChange(2)}
          />
        );
      case 4:
        return (
          <PhotosStep
            onNext={() => handleStepChange(5)}
            onBack={() => handleStepChange(3)}
          />
        );
      case 5:
        return (
          <LocationStep
            onNext={() => handleStepChange(6)}
            onBack={() => handleStepChange(4)}
          />
        );
      case 6:
        return (
          <PaymentDeliveryStep
            onNext={() => handleStepChange(7)}
            onBack={() => handleStepChange(5)}
          />
        );
      case 7:
        return (
          <PreviewPublishStep
            onBack={() => handleStepChange(6)}
            onComplete={handleComplete}
          />
        );
      default:
        return null;
    }
  };

  return (
    <CreateListingProvider draftId={draftId}>
      <div className="container mx-auto px-2 sm:px-4 py-4 sm:py-8 min-h-screen">
        {/* Оффлайн индикатор */}
        <OfflineIndicator />

        {/* Региональный заголовок с традиционным стилем */}
        <div className="text-center mb-6 sm:mb-8">
          <div className="flex items-center justify-between mb-4">
            <div className="flex-1" />
            <h1 className="text-2xl sm:text-3xl font-bold text-primary">
              {t('create_listing.title')}
            </h1>
            <div className="flex items-center gap-2">
              <DraftStatus />
              <DraftIndicator onClick={() => setShowDraftsModal(true)} />
            </div>
          </div>
          <p className="text-sm text-base-content/70">
            {t('create_listing.subtitle_regional')}
          </p>
        </div>

        {/* Индикатор экономии трафика */}
        <div className="alert alert-info mb-4 sm:hidden">
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
            ></path>
          </svg>
          <span className="text-xs">
            {t('create_listing.data_saving_mode')}
          </span>
        </div>

        <StepWizard
          steps={steps}
          currentStep={currentStep}
          onStepClick={handleStepChange}
        />

        <div className="mt-4 sm:mt-8">{renderStep()}</div>

        {/* Модальное окно черновиков */}
        <DraftsModal
          isOpen={showDraftsModal}
          onClose={() => setShowDraftsModal(false)}
        />
      </div>
    </CreateListingProvider>
  );
}
