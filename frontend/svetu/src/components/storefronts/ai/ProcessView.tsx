'use client';

import React, { useEffect, useState } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useCreateAIProduct } from '@/contexts/CreateAIProductContext';
import { storefrontAI } from '@/services/ai/storefronts.service';
import { extractLocationFromImages } from '@/utils/exifUtils';
import { useAddressGeocoding } from '@/hooks/useAddressGeocoding';

interface ProcessViewProps {
  storefrontSlug: string;
}

interface ProcessStep {
  id: string;
  label: string;
  status: 'pending' | 'processing' | 'completed' | 'error';
  message?: string;
}

export default function ProcessView({
  storefrontSlug: _storefrontSlug,
}: ProcessViewProps) {
  const t = useTranslations('storefronts');
  const locale = useLocale();
  const { state, setAIData, setView, setError, setProcessing } =
    useCreateAIProduct();

  // Геокодирование для преобразования координат в адрес
  const geocoding = useAddressGeocoding({
    country: ['rs'],
    language: locale,
  });

  const [steps, setSteps] = useState<ProcessStep[]>([
    { id: 'analyze', label: t('analyzingImages'), status: 'pending' },
    {
      id: 'exif',
      label: t('detectingLocationFromExif') || 'Определение адреса из EXIF',
      status: 'pending',
    },
    { id: 'category', label: t('detectingCategory'), status: 'pending' },
    { id: 'titles', label: t('generatingTitleVariants'), status: 'pending' },
    { id: 'translations', label: t('creatingTranslations'), status: 'pending' },
  ]);

  const updateStepStatus = (
    id: string,
    status: ProcessStep['status'],
    message?: string
  ) => {
    setSteps((prev) =>
      prev.map((step) => (step.id === id ? { ...step, status, message } : step))
    );
  };

  useEffect(() => {
    let isCancelled = false;

    const processImages = async () => {
      setProcessing(true);
      setError(null);

      try {
        // AI сервис теперь использует BFF proxy с credentials: 'include'
        // Токен передается автоматически через httpOnly cookies

        // Step 1: Analyze main image
        updateStepStatus('analyze', 'processing');

        const mainImageFile = state.imageFiles[0];
        if (!mainImageFile) {
          throw new Error('No image file found');
        }

        const base64Image = await storefrontAI.imageToBase64(mainImageFile);
        // AI всегда анализирует на английском для консистентности
        const analysisResult = await storefrontAI.analyzeProductImage(
          base64Image,
          'en'
        );

        if (isCancelled) return;

        updateStepStatus('analyze', 'completed', 'Image analyzed successfully');

        // Step 2: Extract location from EXIF
        updateStepStatus('exif', 'processing');

        let locationData = null;
        const locationTranslations: Record<
          string,
          { address: string; city: string; region: string }
        > = {};
        try {
          const exifLocation = await extractLocationFromImages(
            state.imageFiles
          );
          if (exifLocation) {
            console.log('[ProcessView] EXIF location found:', exifLocation);

            // Геокодируем координаты в адрес для текущей локали
            const geocodedAddress = await geocoding.reverseGeocode(
              exifLocation.latitude,
              exifLocation.longitude
            );

            console.log('[ProcessView] Geocoded address:', geocodedAddress);

            if (geocodedAddress) {
              locationData = {
                latitude: exifLocation.latitude,
                longitude: exifLocation.longitude,
                address:
                  (geocodedAddress as any).place_name ||
                  (geocodedAddress as any).display_name ||
                  '',
                city: (geocodedAddress as any).address?.city || '',
                region: (geocodedAddress as any).address?.state || '',
                source: 'exif' as const,
              };

              // Геокодируем для всех языков (en, ru, sr)
              const languages = ['en', 'ru', 'sr'];
              for (const lang of languages) {
                try {
                  const response = await fetch(
                    `/api/v1/gis/geocode/reverse?lat=${exifLocation.latitude}&lng=${exifLocation.longitude}&language=${lang}`
                  );
                  if (response.ok) {
                    const geocoded = await response.json();
                    const data = geocoded.data || geocoded;
                    locationTranslations[lang] = {
                      address: data.place_name || data.display_name || '',
                      city:
                        data.address?.city ||
                        data.address_components?.city ||
                        '',
                      region:
                        data.address?.state ||
                        data.address_components?.district ||
                        '',
                    };
                  }
                } catch (err) {
                  console.error(`Failed to geocode for ${lang}:`, err);
                }
              }

              updateStepStatus(
                'exif',
                'completed',
                `${locationData.city || 'Location'} detected`
              );
            } else {
              updateStepStatus('exif', 'completed', 'No address found');
            }
          } else {
            updateStepStatus('exif', 'completed', 'No GPS data in photos');
          }
        } catch (exifError) {
          console.error('[ProcessView] EXIF extraction error:', exifError);
          updateStepStatus('exif', 'completed', 'No location data available');
        }

        if (isCancelled) return;

        // Step 3: Detect category
        updateStepStatus('category', 'processing');

        // Определение категории также на английском
        const categoryResult = await storefrontAI.detectCategory(
          analysisResult.title,
          analysisResult.description,
          analysisResult.categoryHints,
          'en'
        );

        if (isCancelled) return;

        updateStepStatus(
          'category',
          'completed',
          `Category: ${categoryResult.categoryName}`
        );

        // Step 3: A/B test titles
        updateStepStatus('titles', 'processing');

        const titleResult = await storefrontAI.abTestTitles(
          analysisResult.titleVariants
        );

        if (isCancelled) return;

        updateStepStatus(
          'titles',
          'completed',
          `${analysisResult.titleVariants.length} variants generated`
        );

        // Step 4: Translate content
        updateStepStatus('translations', 'processing');

        // AI анализирует на английском, переводим на остальные языки
        const targetLanguages = ['ru', 'sr']; // Исключаем 'en' так как это исходный язык
        const translationResult = await storefrontAI.translateProductContent(
          {
            title: analysisResult.title,
            description: analysisResult.description,
          },
          targetLanguages,
          'en' // Исходный язык - английский
        );

        if (isCancelled) return;

        updateStepStatus(
          'translations',
          'completed',
          `Translated to ${targetLanguages.length} languages`
        );

        // Debug logging
        console.log('[ProcessView] Translation result:', translationResult);
        console.log(
          '[ProcessView] Translations:',
          translationResult.translations
        );

        // Backend возвращает переводы напрямую без обёртки "translations"
        const translations =
          translationResult.translations || translationResult || {};

        // Добавляем английский оригинал к переводам
        translations.en = {
          title: analysisResult.title,
          description: analysisResult.description,
        };

        // Добавляем переводы адресов к каждому языку
        if (Object.keys(locationTranslations).length > 0) {
          Object.keys(translations).forEach((lang) => {
            if (locationTranslations[lang]) {
              (translations as any)[lang] = {
                ...translations[lang],
                address: locationTranslations[lang].address,
                city: locationTranslations[lang].city,
                region: locationTranslations[lang].region,
              };
            }
          });
        }

        console.log('[ProcessView] Final translations:', translations);
        console.log('[ProcessView] Final locationData to save:', locationData);
        console.log(
          '[ProcessView] analysisResult.location:',
          analysisResult.location
        );

        // Update AI data in context
        const aiDataToSave = {
          title: analysisResult.title,
          titleVariants: analysisResult.titleVariants,
          selectedTitleIndex: titleResult.bestVariantIndex || 0,
          description: analysisResult.description,
          category: categoryResult.categoryName,
          categoryId: categoryResult.categoryId,
          categoryProbabilities: [], // TODO: Map alternativeIds to full category info
          price: analysisResult.price || 0,
          priceRange: analysisResult.priceRange || { min: 0, max: 0 },
          currency: 'RSD',
          attributes: analysisResult.attributes || {},
          suggestedVariants:
            analysisResult.suggestedVariants?.map((v) => ({
              ...v,
              stockQuantity: v.price ? 1 : 0, // Default stock for variants
            })) || [],
          stockQuantity: analysisResult.stockEstimate || 1,
          condition: analysisResult.condition || 'new',
          keywords: analysisResult.keywords || [],
          translations: translations,
          location: locationData || analysisResult.location || null,
        };

        console.log(
          '[ProcessView] Saving aiData with location:',
          aiDataToSave.location
        );
        setAIData(aiDataToSave);

        // Small delay to show completion
        setTimeout(() => {
          if (!isCancelled) {
            setView('enhance');
          }
        }, 1000);
      } catch (error: any) {
        console.error('AI processing error:', error);

        const failedStep = steps.find((s) => s.status === 'processing');
        if (failedStep) {
          updateStepStatus(
            failedStep.id,
            'error',
            error.message || 'Processing failed'
          );
        }

        setError(error.message || 'Failed to process images with AI');
      } finally {
        setProcessing(false);
      }
    };

    processImages();

    return () => {
      isCancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Run once on mount

  const getStepIcon = (status: ProcessStep['status']) => {
    switch (status) {
      case 'completed':
        return (
          <svg
            className="w-6 h-6 text-success"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fillRule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
              clipRule="evenodd"
            />
          </svg>
        );
      case 'processing':
        return (
          <span className="loading loading-spinner loading-md text-primary"></span>
        );
      case 'error':
        return (
          <svg
            className="w-6 h-6 text-error"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fillRule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
              clipRule="evenodd"
            />
          </svg>
        );
      default:
        return (
          <div className="w-6 h-6 rounded-full border-2 border-base-300 bg-base-200"></div>
        );
    }
  };

  const allStepsCompleted = steps.every((s) => s.status === 'completed');
  const hasError = steps.some((s) => s.status === 'error');

  return (
    <div className="space-y-6">
      <div className="text-center">
        <h2 className="text-2xl font-bold text-base-content mb-2">
          {hasError
            ? t('aiProcessingError') || 'Processing Error'
            : allStepsCompleted
              ? t('aiProcessingComplete') || 'Processing Complete'
              : t('aiProcessingImages') || 'AI is analyzing your product'}
        </h2>
        <p className="text-base-content/70">
          {hasError
            ? t('aiProcessingErrorDesc') ||
              'Something went wrong during processing'
            : allStepsCompleted
              ? t('aiProcessingCompleteDesc') || 'Redirecting to enhancement...'
              : t('aiProcessingDesc') || 'This usually takes 10-15 seconds'}
        </p>
      </div>

      {/* Processing Steps */}
      <div className="space-y-4 max-w-2xl mx-auto">
        {steps.map((step, index) => (
          <div
            key={step.id}
            className="flex items-start gap-4 p-4 bg-base-200 rounded-lg"
          >
            <div className="flex-shrink-0 mt-1">{getStepIcon(step.status)}</div>

            <div className="flex-1 min-w-0">
              <div className="flex items-center justify-between mb-1">
                <h3 className="font-semibold text-base-content">
                  {index + 1}. {step.label}
                </h3>
                <span
                  className={`text-xs px-2 py-1 rounded-full ${
                    step.status === 'completed'
                      ? 'bg-success/20 text-success'
                      : step.status === 'processing'
                        ? 'bg-primary/20 text-primary'
                        : step.status === 'error'
                          ? 'bg-error/20 text-error'
                          : 'bg-base-300 text-base-content/50'
                  }`}
                >
                  {step.status === 'completed'
                    ? t('statusDone')
                    : step.status === 'processing'
                      ? t('statusProcessing')
                      : step.status === 'error'
                        ? t('statusFailed')
                        : t('statusWaiting')}
                </span>
              </div>

              {step.message && (
                <p className="text-sm text-base-content/60">{step.message}</p>
              )}
            </div>
          </div>
        ))}
      </div>

      {/* AI Animation */}
      {!allStepsCompleted && !hasError && (
        <div className="flex justify-center py-8">
          <div className="relative">
            <div className="w-24 h-24 rounded-full bg-primary/10 flex items-center justify-center animate-pulse">
              <svg
                className="w-12 h-12 text-primary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                />
              </svg>
            </div>
          </div>
        </div>
      )}

      {/* Retry button if error */}
      {hasError && (
        <div className="flex justify-center gap-3">
          <button onClick={() => setView('upload')} className="btn btn-outline">
            {t('backToUpload') || 'Back to Upload'}
          </button>
          <button
            onClick={() => window.location.reload()}
            className="btn btn-primary"
          >
            {t('retryProcessing') || 'Retry Processing'}
          </button>
        </div>
      )}
    </div>
  );
}
