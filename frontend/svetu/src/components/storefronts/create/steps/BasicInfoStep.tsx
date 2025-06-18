'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateStorefrontContext } from '@/contexts/CreateStorefrontContext';

interface BasicInfoStepProps {
  onNext: () => void;
}

export default function BasicInfoStep({ onNext }: BasicInfoStepProps) {
  const t = useTranslations();
  const { formData, updateFormData } = useCreateStorefrontContext();
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Auto-generate slug from name
  useEffect(() => {
    if (formData.name && !formData.slug) {
      const slug = formData.name
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/^-+|-+$/g, '');
      updateFormData({ slug });
    }
  }, [formData.name, formData.slug, updateFormData]);

  const validate = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name || formData.name.length < 3) {
      newErrors.name = t('create_storefront.errors.name_required');
    }

    if (!formData.slug || formData.slug.length < 3) {
      newErrors.slug = t('create_storefront.errors.slug_required');
    } else if (!/^[a-z0-9-]+$/.test(formData.slug)) {
      newErrors.slug = t('create_storefront.errors.slug_invalid');
    }

    if (!formData.description || formData.description.length < 20) {
      newErrors.description = t(
        'create_storefront.errors.description_required'
      );
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
            {t('create_storefront.basic_info.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('create_storefront.basic_info.subtitle')}
          </p>

          <div className="form-control w-full">
            <label className="label">
              <span className="label-text">
                {t('create_storefront.basic_info.name')}
              </span>
              <span className="label-text-alt text-error">*</span>
            </label>
            <input
              type="text"
              placeholder={t('create_storefront.basic_info.name_placeholder')}
              className={`input input-bordered w-full ${errors.name ? 'input-error' : ''}`}
              value={formData.name}
              onChange={(e) => updateFormData({ name: e.target.value })}
            />
            {errors.name && (
              <label className="label">
                <span className="label-text-alt text-error">{errors.name}</span>
              </label>
            )}
          </div>

          <div className="form-control w-full mt-4">
            <label className="label">
              <span className="label-text">
                {t('create_storefront.basic_info.slug')}
              </span>
              <span className="label-text-alt text-error">*</span>
            </label>
            <div className="flex items-center gap-2">
              <span className="text-base-content/70">svetu.rs/</span>
              <input
                type="text"
                placeholder={t('create_storefront.basic_info.slug_placeholder')}
                className={`input input-bordered flex-1 ${errors.slug ? 'input-error' : ''}`}
                value={formData.slug}
                onChange={(e) =>
                  updateFormData({ slug: e.target.value.toLowerCase() })
                }
              />
            </div>
            {errors.slug && (
              <label className="label">
                <span className="label-text-alt text-error">{errors.slug}</span>
              </label>
            )}
          </div>

          <div className="form-control w-full mt-4">
            <label className="label">
              <span className="label-text">
                {t('create_storefront.basic_info.description')}
              </span>
              <span className="label-text-alt text-error">*</span>
            </label>
            <textarea
              className={`textarea textarea-bordered h-24 ${errors.description ? 'textarea-error' : ''}`}
              placeholder={t(
                'create_storefront.basic_info.description_placeholder'
              )}
              value={formData.description}
              onChange={(e) => updateFormData({ description: e.target.value })}
            />
            {errors.description && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {errors.description}
                </span>
              </label>
            )}
          </div>

          <div className="form-control w-full mt-4">
            <label className="label">
              <span className="label-text">
                {t('create_storefront.basic_info.business_type')}
              </span>
            </label>
            <select
              className="select select-bordered w-full"
              value={formData.businessType}
              onChange={(e) => updateFormData({ businessType: e.target.value })}
            >
              <option value="retail">
                {t('create_storefront.business_types.retail')}
              </option>
              <option value="service">
                {t('create_storefront.business_types.service')}
              </option>
              <option value="restaurant">
                {t('create_storefront.business_types.restaurant')}
              </option>
              <option value="grocery">
                {t('create_storefront.business_types.grocery')}
              </option>
              <option value="other">
                {t('create_storefront.business_types.other')}
              </option>
            </select>
          </div>

          <div className="card-actions justify-end mt-6">
            <button className="btn btn-primary" onClick={handleNext}>
              {t('common.next')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
