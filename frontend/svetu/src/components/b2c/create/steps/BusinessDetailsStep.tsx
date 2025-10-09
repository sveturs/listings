'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateStorefrontContext } from '@/contexts/CreateB2CStoreContext';

interface BusinessDetailsStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function BusinessDetailsStep({
  onNext,
  onBack,
}: BusinessDetailsStepProps) {
  const t = useTranslations('create_storefront');
  const tCommon = useTranslations('common');
  const { formData, updateFormData } = useCreateStorefrontContext();
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validate = () => {
    const newErrors: Record<string, string> = {};

    // Email validation if provided
    if (formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = t('errors.email_invalid');
    }

    // Phone validation if provided
    if (formData.phone && !/^[\d\s\-\+\(\)]+$/.test(formData.phone)) {
      newErrors.phone = t('errors.phone_invalid');
    }

    // Website validation if provided
    if (formData.website && !/^https?:\/\/.+\..+/.test(formData.website)) {
      newErrors.website = t('errors.website_invalid');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleNext = () => {
    if (validate()) {
      onNext();
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">
            {t('business_details.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('business_details.subtitle')}
          </p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control w-full">
              <label className="label">
                <span className="label-text">
                  {t('business_details.registration_number')}
                </span>
              </label>
              <input
                type="text"
                placeholder={t(
                  'business_details.registration_number_placeholder'
                )}
                className="input input-bordered w-full"
                value={formData.registrationNumber || ''}
                onChange={(e) =>
                  updateFormData({ registrationNumber: e.target.value })
                }
              />
            </div>

            <div className="form-control w-full">
              <label className="label">
                <span className="label-text">
                  {t('business_details.tax_number')}
                </span>
              </label>
              <input
                type="text"
                placeholder={t('business_details.tax_number_placeholder')}
                className="input input-bordered w-full"
                value={formData.taxNumber || ''}
                onChange={(e) => updateFormData({ taxNumber: e.target.value })}
              />
            </div>

            <div className="form-control w-full">
              <label className="label">
                <span className="label-text">
                  {t('business_details.vat_number')}
                </span>
              </label>
              <input
                type="text"
                placeholder={t('business_details.vat_number_placeholder')}
                className="input input-bordered w-full"
                value={formData.vatNumber || ''}
                onChange={(e) => updateFormData({ vatNumber: e.target.value })}
              />
            </div>

            <div className="form-control w-full">
              <label className="label">
                <span className="label-text">
                  {t('business_details.phone')}
                </span>
              </label>
              <input
                type="tel"
                placeholder={t('business_details.phone_placeholder')}
                className={`input input-bordered w-full ${errors.phone ? 'input-error' : ''}`}
                value={formData.phone || ''}
                onChange={(e) => updateFormData({ phone: e.target.value })}
              />
              {errors.phone && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {errors.phone}
                  </span>
                </label>
              )}
            </div>

            <div className="form-control w-full">
              <label className="label">
                <span className="label-text">
                  {t('business_details.email')}
                </span>
              </label>
              <input
                type="email"
                placeholder={t('business_details.email_placeholder')}
                className={`input input-bordered w-full ${errors.email ? 'input-error' : ''}`}
                value={formData.email || ''}
                onChange={(e) => updateFormData({ email: e.target.value })}
              />
              {errors.email && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {errors.email}
                  </span>
                </label>
              )}
            </div>

            <div className="form-control w-full">
              <label className="label">
                <span className="label-text">
                  {t('business_details.website')}
                </span>
              </label>
              <input
                type="url"
                placeholder={t('business_details.website_placeholder')}
                className={`input input-bordered w-full ${errors.website ? 'input-error' : ''}`}
                value={formData.website || ''}
                onChange={(e) => updateFormData({ website: e.target.value })}
              />
              {errors.website && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {errors.website}
                  </span>
                </label>
              )}
            </div>
          </div>

          <div className="card-actions justify-between mt-6">
            <button className="btn btn-ghost" onClick={onBack}>
              {tCommon('back')}
            </button>
            <button className="btn btn-primary" onClick={handleNext}>
              {tCommon('next')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
