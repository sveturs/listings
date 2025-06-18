'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useAuthContext } from '@/contexts/AuthContext';
import StepWizard from '@/components/create-listing/StepWizard';
import BasicInfoStep from '@/components/storefronts/create/steps/BasicInfoStep';
import BusinessDetailsStep from '@/components/storefronts/create/steps/BusinessDetailsStep';
import LocationStep from '@/components/storefronts/create/steps/LocationStep';
import BusinessHoursStep from '@/components/storefronts/create/steps/BusinessHoursStep';
import PaymentDeliveryStep from '@/components/storefronts/create/steps/PaymentDeliveryStep';
import StaffSetupStep from '@/components/storefronts/create/steps/StaffSetupStep';
import PreviewPublishStep from '@/components/storefronts/create/steps/PreviewPublishStep';
import { CreateStorefrontProvider } from '@/contexts/CreateStorefrontContext';
import { toast } from '@/utils/toast';

const steps = [
  { id: 'basic_info', label: 'create_storefront.steps.basic_info' },
  { id: 'business_details', label: 'create_storefront.steps.business_details' },
  { id: 'location', label: 'create_storefront.steps.location' },
  { id: 'business_hours', label: 'create_storefront.steps.business_hours' },
  { id: 'payment_delivery', label: 'create_storefront.steps.payment_delivery' },
  { id: 'staff_setup', label: 'create_storefront.steps.staff_setup' },
  { id: 'preview', label: 'create_storefront.steps.preview' },
];

export default function CreateStorefrontPage() {
  const router = useRouter();
  const t = useTranslations();
  const { user, isLoading: authLoading } = useAuthContext();
  const [currentStep, setCurrentStep] = useState(0);

  useEffect(() => {
    if (!authLoading && !user) {
      toast.error(t('create_storefront.auth_required'));
      router.push('/');
    }
  }, [user, authLoading, router, t]);

  const handleStepChange = (newStep: number) => {
    setCurrentStep(newStep);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleComplete = () => {
    toast.success(t('create_storefront.success'));
    router.push('/storefronts/my');
  };

  if (authLoading || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  const renderStep = () => {
    switch (currentStep) {
      case 0:
        return <BasicInfoStep onNext={() => handleStepChange(1)} />;
      case 1:
        return (
          <BusinessDetailsStep
            onNext={() => handleStepChange(2)}
            onBack={() => handleStepChange(0)}
          />
        );
      case 2:
        return (
          <LocationStep
            onNext={() => handleStepChange(3)}
            onBack={() => handleStepChange(1)}
          />
        );
      case 3:
        return (
          <BusinessHoursStep
            onNext={() => handleStepChange(4)}
            onBack={() => handleStepChange(2)}
          />
        );
      case 4:
        return (
          <PaymentDeliveryStep
            onNext={() => handleStepChange(5)}
            onBack={() => handleStepChange(3)}
          />
        );
      case 5:
        return (
          <StaffSetupStep
            onNext={() => handleStepChange(6)}
            onBack={() => handleStepChange(4)}
          />
        );
      case 6:
        return (
          <PreviewPublishStep
            onBack={() => handleStepChange(5)}
            onComplete={handleComplete}
          />
        );
      default:
        return null;
    }
  };

  return (
    <CreateStorefrontProvider>
      <div className="container mx-auto px-2 sm:px-4 py-4 sm:py-8 min-h-screen">
        <div className="text-center mb-6 sm:mb-8">
          <h1 className="text-2xl sm:text-3xl font-bold text-primary">
            {t('create_storefront.title')}
          </h1>
          <p className="text-sm text-base-content/70 mt-2">
            {t('create_storefront.subtitle')}
          </p>
        </div>

        <StepWizard
          steps={steps}
          currentStep={currentStep}
          onStepClick={handleStepChange}
        />

        <div className="mt-4 sm:mt-8">{renderStep()}</div>
      </div>
    </CreateStorefrontProvider>
  );
}
