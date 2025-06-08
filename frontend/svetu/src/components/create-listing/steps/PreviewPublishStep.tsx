'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';
import { ListingsService } from '@/services/listings';
import { toast } from '@/utils/toast';
import { useRouter } from '@/i18n/routing';
import Image from 'next/image';

interface PreviewPublishStepProps {
  onBack: () => void;
  onComplete: () => void;
}

export default function PreviewPublishStep({
  onBack,
  onComplete: _onComplete,
}: PreviewPublishStepProps) {
  const t = useTranslations();
  const router = useRouter();
  const { state, saveDraft, publish } = useCreateListing();
  const [isPublishing, setIsPublishing] = useState(false);
  const [isSavingDraft, setIsSavingDraft] = useState(false);
  const [uploadingImages, setUploadingImages] = useState(false);

  const handlePublish = async () => {
    setIsPublishing(true);
    try {
      // Создаем объявление
      const response = await ListingsService.createListing(state);

      // Если есть изображения, загружаем их
      if (state.images && state.images.length > 0) {
        setUploadingImages(true);
        try {
          // Преобразуем base64 изображения в File объекты
          const files = await Promise.all(
            state.images.map(async (imageUrl, index) => {
              const res = await fetch(imageUrl);
              const blob = await res.blob();
              return new File([blob], `image_${index}.jpg`, {
                type: 'image/jpeg',
              });
            })
          );

          await ListingsService.uploadImages(
            response.id,
            files,
            state.mainImageIndex
          );
        } catch (imageError) {
          console.error('Error uploading images:', imageError);
          toast.error(t('create_listing.errors.image_upload_failed'));
        } finally {
          setUploadingImages(false);
        }
      }

      // Удаляем черновик
      await ListingsService.deleteDraft(state.category?.id);

      // Обновляем состояние
      publish();

      // Показываем успешное сообщение
      toast.success(t('create_listing.success'));

      // Перенаправляем на страницу объявления
      setTimeout(() => {
        router.push(`/marketplace/${response.id}`);
      }, 1000);
    } catch (error) {
      console.error('Error publishing:', error);
      toast.error(t('create_listing.errors.publish_failed'));
    } finally {
      setIsPublishing(false);
    }
  };

  const handleSaveDraft = async () => {
    setIsSavingDraft(true);
    try {
      // Сохраняем черновик в localStorage
      await ListingsService.saveDraft(state);
      saveDraft();
      toast.success(t('create_listing.draft_saved'));
    } catch (error) {
      console.error('Error saving draft:', error);
      toast.error(t('create_listing.errors.draft_save_failed'));
    } finally {
      setIsSavingDraft(false);
    }
  };

  const formatCurrency = (amount: number, currency: string) => {
    const symbols: { [key: string]: string } = {
      RSD: 'РСД',
      EUR: '€',
      HRK: 'kn',
      MKD: 'ден',
    };
    return `${amount.toLocaleString()} ${symbols[currency] || currency}`;
  };

  const getConditionLabel = (condition: string) => {
    const labels: { [key: string]: string } = {
      new: 'Ново',
      used: 'Половно',
      refurbished: 'Обновљено',
    };
    return labels[condition] || condition;
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            👁️ {t('create_listing.preview.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('create_listing.preview.description')}
          </p>

          {/* Превью объявления */}
          <div className="card border border-base-200 bg-base-50">
            <div className="card-body p-4">
              {/* Основная информация */}
              <div className="flex items-start gap-4">
                {/* Главное изображение */}
                {state.images.length > 0 && (
                  <div className="avatar">
                    <div className="w-20 h-20 rounded-lg">
                      <Image
                        src={state.images[state.mainImageIndex]}
                        alt="Main"
                        width={80}
                        height={80}
                        className="object-cover"
                      />
                    </div>
                  </div>
                )}

                <div className="flex-1">
                  <h3 className="card-title text-lg">{state.title}</h3>
                  <p className="text-primary font-bold text-xl">
                    {formatCurrency(state.price, state.currency)}
                  </p>
                  <div className="flex items-center gap-2 mt-2">
                    <span className="badge badge-outline badge-sm">
                      {getConditionLabel(state.condition)}
                    </span>
                    {state.category && (
                      <span className="badge badge-primary badge-sm">
                        {state.category.name}
                      </span>
                    )}
                  </div>
                </div>
              </div>

              {/* Описание */}
              <div className="mt-4">
                <p className="text-sm text-base-content/80 line-clamp-3">
                  {state.description}
                </p>
              </div>

              {/* Галерея изображений */}
              {state.images.length > 1 && (
                <div className="mt-4">
                  <div className="flex gap-2 overflow-x-auto">
                    {state.images.map((image, index) => (
                      <div key={index} className="avatar">
                        <div className="w-12 h-12 rounded">
                          <Image
                            src={image}
                            alt={`Image ${index + 1}`}
                            width={48}
                            height={48}
                            className="object-cover"
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Местоположение */}
              {state.location && (
                <div className="mt-4 flex items-center gap-2 text-sm text-base-content/60">
                  <span>📍</span>
                  <span>
                    {state.location.city}, {state.location.region}
                  </span>
                </div>
              )}

              {/* Региональные особенности */}
              <div className="mt-4 space-y-2">
                {/* Система доверия */}
                {state.trust.phoneVerified && (
                  <div className="flex items-center gap-2 text-sm">
                    <span className="badge badge-success badge-sm">
                      ✅ Верифован телефон
                    </span>
                  </div>
                )}

                {/* Способы оплаты */}
                <div className="flex items-center gap-2 text-sm">
                  <span>💳</span>
                  <span className="text-base-content/60">
                    {state.payment.methods.includes('cod') && 'Наложен платеж'}
                    {state.payment.personalMeeting && ' • Лична предаја'}
                  </span>
                </div>

                {/* Тип встречи */}
                {state.trust.preferredMeetingType && (
                  <div className="flex items-center gap-2 text-sm">
                    <span>🤝</span>
                    <span className="text-base-content/60">
                      {state.trust.preferredMeetingType === 'personal' &&
                        'Лична предаја'}
                      {state.trust.preferredMeetingType === 'pickup' &&
                        'Преузимање'}
                      {state.trust.preferredMeetingType === 'delivery' &&
                        'Достава'}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Сводка регионального функционала */}
          <div className="alert alert-info mt-6">
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
            <div className="text-sm">
              <p className="font-medium">
                🏪 {t('create_listing.preview.regional_summary')}
              </p>
              <ul className="text-xs mt-2 space-y-1">
                <li>
                  • {t('create_listing.preview.script')}:{' '}
                  {state.localization.script === 'cyrillic'
                    ? 'Ћирилица'
                    : 'Latinica'}
                </li>
                <li>
                  • {t('create_listing.preview.trust')}:{' '}
                  {state.trust.phoneVerified ? 'Верификован' : 'Неверификован'}
                </li>
                <li>
                  • {t('create_listing.preview.payment')}:{' '}
                  {state.payment.codEnabled
                    ? 'Наложен платеж омогућен'
                    : 'Само готовина'}
                </li>
                <li>
                  • {t('create_listing.preview.meeting')}:{' '}
                  {state.trust.preferredMeetingType}
                </li>
              </ul>
            </div>
          </div>

          {/* Региональные правила публикации */}
          <div className="alert alert-warning mt-4">
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
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.982 16.5c-.77.833.192 2.5 1.732 2.5z"
              ></path>
            </svg>
            <div className="text-sm">
              <p className="font-medium">
                ⚖️ {t('create_listing.preview.rules.title')}
              </p>
              <ul className="text-xs mt-2 space-y-1">
                <li>
                  • {t('create_listing.preview.rules.honest_description')}
                </li>
                <li>• {t('create_listing.preview.rules.fair_pricing')}</li>
                <li>
                  • {t('create_listing.preview.rules.respectful_communication')}
                </li>
                <li>• {t('create_listing.preview.rules.safe_meetings')}</li>
              </ul>
            </div>
          </div>

          {/* Кнопки действий */}
          <div className="card-actions justify-between mt-8">
            <button className="btn btn-outline" onClick={onBack}>
              ← {t('common.back')}
            </button>

            <div className="flex gap-2">
              <button
                className={`btn btn-outline ${isSavingDraft ? 'loading' : ''}`}
                onClick={handleSaveDraft}
                disabled={isSavingDraft || isPublishing}
              >
                {isSavingDraft ? '' : '💾'} {t('create_listing.save_draft')}
              </button>

              <button
                className={`btn btn-primary ${isPublishing ? 'loading' : ''}`}
                onClick={handlePublish}
                disabled={isPublishing || isSavingDraft}
              >
                {isPublishing ? '' : '🚀'}
                {uploadingImages
                  ? t('create_listing.uploading_images')
                  : t('create_listing.publish')}
              </button>
            </div>
          </div>

          {/* Локальная подсказка о публикации */}
          <div className="text-center mt-4">
            <p className="text-xs text-base-content/60">
              {t('create_listing.preview.publish_info')}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
