'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';

interface TrustSetupStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function TrustSetupStep({
  onNext,
  onBack,
}: TrustSetupStepProps) {
  const t = useTranslations();
  const [formData, setFormData] = useState({
    phoneVerified: false,
    preferredMeetingType: 'personal', // personal, pickup, delivery
    meetingLocations: [] as string[],
    availableHours: '',
    trustBadges: [] as string[],
  });
  const [isVerifyingPhone, setIsVerifyingPhone] = useState(false);

  const handlePhoneVerification = async () => {
    setIsVerifyingPhone(true);
    // TODO: Интегрировать реальную верификацию телефона через SMS/звонок
    // ВРЕМЕННОЕ РЕШЕНИЕ: Фейковая верификация для демонстрации
    // См. TODO #15: Интегрировать реальную верификацию телефона
    setTimeout(() => {
      setFormData((prev) => ({ ...prev, phoneVerified: true }));
      setIsVerifyingPhone(false);
    }, 2000);
  };

  const meetingTypes = [
    {
      id: 'personal',
      label: 'trust.meeting.personal',
      icon: '🤝',
      description: 'trust.meeting.personal_desc',
    },
    {
      id: 'pickup',
      label: 'trust.meeting.pickup',
      icon: '🏪',
      description: 'trust.meeting.pickup_desc',
    },
    {
      id: 'delivery',
      label: 'trust.meeting.delivery',
      icon: '🚚',
      description: 'trust.meeting.delivery_desc',
    },
  ];

  // TODO: Загружать из API безопасные места встреч по городам
  // ВРЕМЕННОЕ РЕШЕНИЕ: Хардкодные места для демонстрации
  // См. TODO #17: Создать API для безопасных мест встреч по городам
  const commonMeetingPlaces = [
    'Тржни центар',
    'Кафе',
    'Паркинг',
    'Аутобуска станица',
    'Трг',
  ];

  const canProceed = formData.phoneVerified && formData.preferredMeetingType;

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            🛡️ {t('trust.setup_title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('trust.setup_description')}
          </p>

          {/* Верификация телефона */}
          <div className="space-y-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  📞 {t('trust.phone_verification')}
                </span>
                <span className="label-text-alt text-error">*</span>
              </label>

              {!formData.phoneVerified ? (
                <div className="space-y-2">
                  <div className="alert alert-info">
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
                      <p>{t('trust.phone_verification_info')}</p>
                      <p className="text-xs mt-1">
                        {t('trust.local_operators_only')}
                      </p>
                    </div>
                  </div>

                  <button
                    className={`btn btn-primary ${isVerifyingPhone ? 'loading' : ''}`}
                    onClick={handlePhoneVerification}
                    disabled={isVerifyingPhone}
                  >
                    {isVerifyingPhone
                      ? t('trust.verifying')
                      : t('trust.verify_phone')}
                  </button>
                </div>
              ) : (
                <div className="alert alert-success">
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
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <span>{t('trust.phone_verified')}</span>
                </div>
              )}
            </div>

            {/* Тип встречи */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  🤝 {t('trust.preferred_meeting')}
                </span>
                <span className="label-text-alt text-error">*</span>
              </label>

              <div className="grid gap-3">
                {meetingTypes.map((type) => (
                  <label key={type.id} className="cursor-pointer">
                    <input
                      type="radio"
                      name="meetingType"
                      value={type.id}
                      checked={formData.preferredMeetingType === type.id}
                      onChange={(e) =>
                        setFormData((prev) => ({
                          ...prev,
                          preferredMeetingType: e.target.value,
                        }))
                      }
                      className="sr-only"
                    />
                    <div
                      className={`
                      card border-2 transition-all duration-200
                      ${
                        formData.preferredMeetingType === type.id
                          ? 'border-primary bg-primary/5'
                          : 'border-base-300 hover:border-primary/50'
                      }
                    `}
                    >
                      <div className="card-body p-4">
                        <div className="flex items-start gap-3">
                          <span className="text-2xl">{type.icon}</span>
                          <div className="flex-1">
                            <h3 className="font-medium">{t(type.label)}</h3>
                            <p className="text-sm text-base-content/60 mt-1">
                              {t(type.description)}
                            </p>
                          </div>
                          {formData.preferredMeetingType === type.id && (
                            <svg
                              className="w-6 h-6 text-primary"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                clipRule="evenodd"
                              />
                            </svg>
                          )}
                        </div>
                      </div>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            {/* Места встреч (если выбрана личная передача) */}
            {formData.preferredMeetingType === 'personal' && (
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    📍 {t('trust.meeting_places')}
                  </span>
                </label>

                <div className="flex flex-wrap gap-2 mb-3">
                  {commonMeetingPlaces.map((place) => (
                    <button
                      key={place}
                      type="button"
                      onClick={() => {
                        const newLocations = formData.meetingLocations.includes(
                          place
                        )
                          ? formData.meetingLocations.filter((l) => l !== place)
                          : [...formData.meetingLocations, place];
                        setFormData((prev) => ({
                          ...prev,
                          meetingLocations: newLocations,
                        }));
                      }}
                      className={`
                        btn btn-sm
                        ${
                          formData.meetingLocations.includes(place)
                            ? 'btn-primary'
                            : 'btn-outline'
                        }
                      `}
                    >
                      {place}
                    </button>
                  ))}
                </div>

                <p className="text-xs text-base-content/60">
                  {t('trust.meeting_places_hint')}
                </p>
              </div>
            )}

            {/* Доступные часы */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  ⏰ {t('trust.available_hours')}
                </span>
              </label>
              <input
                type="text"
                placeholder="нпр. 9:00 - 18:00, сваки дан"
                className="input input-bordered"
                value={formData.availableHours}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    availableHours: e.target.value,
                  }))
                }
              />
              <label className="label">
                <span className="label-text-alt">
                  {t('trust.available_hours_hint')}
                </span>
              </label>
            </div>
          </div>

          {/* Кнопки навигации */}
          <div className="card-actions justify-between mt-6">
            <button className="btn btn-outline" onClick={onBack}>
              ← {t('common.back')}
            </button>
            <button
              className={`btn btn-primary ${!canProceed ? 'btn-disabled' : ''}`}
              onClick={onNext}
              disabled={!canProceed}
            >
              {t('common.continue')} →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
