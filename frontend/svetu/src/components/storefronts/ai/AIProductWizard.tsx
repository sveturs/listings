'use client';

import React, { useEffect, useState } from 'react';
import { useCreateAIProduct } from '@/contexts/CreateAIProductContext';
import { useRouter } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import { tokenManager } from '@/utils/tokenManager';
import UploadView from './UploadView';
import ProcessView from './ProcessView';
import EnhanceView from './EnhanceView';
import PublishView from './PublishView';

interface AIProductWizardProps {
  storefrontSlug: string;
}

export default function AIProductWizard({
  storefrontSlug,
}: AIProductWizardProps) {
  const { state } = useCreateAIProduct();
  const [storefrontId, setStorefrontId] = useState<number | null>(null);

  // Загрузка storefront ID по slug
  useEffect(() => {
    const fetchStorefrontId = async () => {
      console.log(
        '[AIProductWizard] Fetching storefront for slug:',
        storefrontSlug
      );
      try {
        const token = tokenManager.getAccessToken();
        console.log(
          '[AIProductWizard] Token available:',
          token ? 'yes' : 'no'
        );

        const headers: HeadersInit = {
          'Content-Type': 'application/json',
        };
        if (token) {
          headers['Authorization'] = `Bearer ${token}`;
        }

        const url = `/api/v1/storefronts/slug/${storefrontSlug}`;
        console.log('[AIProductWizard] Fetching URL:', url);

        const response = await fetch(url, { headers });
        console.log('[AIProductWizard] Response status:', response.status);

        if (response.ok) {
          const data = await response.json();
          console.log('[AIProductWizard] Fetched storefront:', data);
          // API возвращает storefront напрямую, а не обернутый в data
          const id = data.id || null;
          console.log('[AIProductWizard] Extracted ID:', id);
          setStorefrontId(id);
        } else {
          const errorText = await response.text();
          console.error(
            '[AIProductWizard] Failed to fetch storefront:',
            response.status,
            errorText
          );
        }
      } catch (error) {
        console.error('[AIProductWizard] Error fetching storefront ID:', error);
      }
    };

    if (storefrontSlug) {
      fetchStorefrontId();
    }
  }, [storefrontSlug]);

  // Progress indicator
  const t = useTranslations('storefronts');
  const steps = [
    { id: 'upload', label: t('stepUpload'), icon: '📤' },
    { id: 'process', label: t('stepAIProcessing'), icon: '🤖' },
    { id: 'enhance', label: t('stepEnhance'), icon: '✨' },
    { id: 'publish', label: t('stepPublish'), icon: '🚀' },
  ];

  const currentStepIndex = steps.findIndex(
    (step) => step.id === state.currentView
  );

  return (
    <div className="max-w-6xl mx-auto">
      {/* Progress Bar */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          {steps.map((step, index) => (
            <React.Fragment key={step.id}>
              <div className="flex flex-col items-center">
                <div
                  className={`w-12 h-12 rounded-full flex items-center justify-center text-xl transition-all ${
                    index <= currentStepIndex
                      ? 'bg-primary text-primary-content'
                      : 'bg-base-300 text-base-content/50'
                  }`}
                >
                  {step.icon}
                </div>
                <span
                  className={`text-sm mt-2 ${
                    index <= currentStepIndex
                      ? 'text-primary font-semibold'
                      : 'text-base-content/50'
                  }`}
                >
                  {step.label}
                </span>
              </div>
              {index < steps.length - 1 && (
                <div className="flex-1 h-1 mx-4 bg-base-300 rounded">
                  <div
                    className={`h-full bg-primary rounded transition-all duration-500 ${
                      index < currentStepIndex ? 'w-full' : 'w-0'
                    }`}
                  />
                </div>
              )}
            </React.Fragment>
          ))}
        </div>
      </div>

      {/* Error Display */}
      {state.error && (
        <div className="alert alert-error mb-6">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{state.error}</span>
        </div>
      )}

      {/* Views */}
      <div className="bg-base-100 rounded-lg shadow-lg p-6">
        {state.currentView === 'upload' && <UploadView />}
        {state.currentView === 'process' && (
          <ProcessView storefrontSlug={storefrontSlug} />
        )}
        {state.currentView === 'enhance' && (
          <EnhanceView
            storefrontId={storefrontId}
            storefrontSlug={storefrontSlug}
          />
        )}
        {state.currentView === 'publish' && (
          <PublishView
            storefrontId={storefrontId}
            storefrontSlug={storefrontSlug}
          />
        )}
      </div>
    </div>
  );
}
