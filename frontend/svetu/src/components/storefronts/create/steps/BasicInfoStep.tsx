'use client';

import { useState, useEffect, useRef } from 'react';
import { useTranslations } from 'next-intl';
import Image from 'next/image';
import { useCreateStorefrontContext } from '@/contexts/CreateStorefrontContext';

interface BasicInfoStepProps {
  onNext: () => void;
}

export default function BasicInfoStep({ onNext }: BasicInfoStepProps) {
  const t = useTranslations('create_storefront');
  const tCommon = useTranslations('common');
  const { formData, updateFormData } = useCreateStorefrontContext();
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [logoPreview, setLogoPreview] = useState<string | null>(null);
  const [bannerPreview, setBannerPreview] = useState<string | null>(null);
  const logoInputRef = useRef<HTMLInputElement>(null);
  const bannerInputRef = useRef<HTMLInputElement>(null);

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

  // Handle logo upload
  const handleLogoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > 5 * 1024 * 1024) {
        setErrors({
          ...errors,
          logo: t('errors.logo_too_large'),
        });
        return;
      }
      updateFormData({ logoFile: file });
      const reader = new FileReader();
      reader.onloadend = () => {
        setLogoPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
      setErrors({ ...errors, logo: '' });
    }
  };

  // Handle banner upload
  const handleBannerChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > 10 * 1024 * 1024) {
        setErrors({
          ...errors,
          banner: t('errors.banner_too_large'),
        });
        return;
      }
      updateFormData({ bannerFile: file });
      const reader = new FileReader();
      reader.onloadend = () => {
        setBannerPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
      setErrors({ ...errors, banner: '' });
    }
  };

  const validate = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name || formData.name.length < 3) {
      newErrors.name = t('errors.name_required');
    }

    if (!formData.slug || formData.slug.length < 3) {
      newErrors.slug = t('errors.slug_required');
    } else if (!/^[a-z0-9-]+$/.test(formData.slug)) {
      newErrors.slug = t('errors.slug_invalid');
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
            {t('title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('subtitle')}
          </p>

          <div className="form-control w-full">
            <label className="label">
              <span className="label-text">
                {t('name')}
              </span>
              <span className="label-text-alt text-error">*</span>
            </label>
            <input
              type="text"
              placeholder={t('name_placeholder')}
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
                {t('slug')}
              </span>
              <span className="label-text-alt text-error">*</span>
            </label>
            <div className="flex items-center gap-2">
              <span className="text-base-content/70">svetu.rs/</span>
              <input
                type="text"
                placeholder={t('slug_placeholder')}
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
                {t('description')}
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
                {t('business_type')}
              </span>
            </label>
            <select
              className="select select-bordered w-full"
              value={formData.businessType}
              onChange={(e) => updateFormData({ businessType: e.target.value })}
            >
              <option value="retail">
                {t('business_types.retail')}
              </option>
              <option value="service">
                {t('business_types.service')}
              </option>
              <option value="restaurant">
                {t('business_types.restaurant')}
              </option>
              <option value="grocery">
                {t('business_types.grocery')}
              </option>
              <option value="other">
                {t('business_types.other')}
              </option>
            </select>
          </div>

          {/* Logo Upload */}
          <div className="form-control w-full mt-6">
            <label className="label">
              <span className="label-text">
                {t('logo')}
              </span>
              <span className="label-text-alt text-base-content/60">
                {t('logo_hint')}
              </span>
            </label>
            <div className="flex items-center gap-4">
              <div
                className="avatar cursor-pointer group"
                onClick={() => logoInputRef.current?.click()}
              >
                <div className="w-24 h-24 rounded-xl ring-2 ring-base-300 overflow-hidden bg-base-200 group-hover:ring-primary transition-all">
                  {logoPreview ? (
                    <Image
                      src={logoPreview}
                      alt="Logo preview"
                      width={96}
                      height={96}
                      className="object-cover"
                    />
                  ) : (
                    <div className="flex items-center justify-center h-full">
                      <svg
                        className="w-8 h-8 text-base-content/40"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"
                        />
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"
                        />
                      </svg>
                    </div>
                  )}
                </div>
              </div>
              <div>
                <input
                  ref={logoInputRef}
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  className="hidden"
                  onChange={handleLogoChange}
                />
                <button
                  type="button"
                  className="btn btn-sm btn-outline"
                  onClick={() => logoInputRef.current?.click()}
                >
                  {t('choose_logo')}
                </button>
                <p className="text-xs text-base-content/60 mt-1">
                  {t('logo_requirements')}
                </p>
              </div>
            </div>
            {errors.logo && (
              <label className="label">
                <span className="label-text-alt text-error">{errors.logo}</span>
              </label>
            )}
          </div>

          {/* Banner Upload */}
          <div className="form-control w-full mt-6">
            <label className="label">
              <span className="label-text">
                {t('banner')}
              </span>
              <span className="label-text-alt text-base-content/60">
                {t('banner_hint')}
              </span>
            </label>
            <div
              className="relative h-32 bg-base-200 rounded-xl overflow-hidden cursor-pointer group"
              onClick={() => bannerInputRef.current?.click()}
            >
              {bannerPreview ? (
                <Image
                  src={bannerPreview}
                  alt="Banner preview"
                  fill
                  className="object-cover"
                />
              ) : (
                <div className="flex flex-col items-center justify-center h-full">
                  <svg
                    className="w-12 h-12 text-base-content/40 mb-2"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                    />
                  </svg>
                  <span className="text-sm text-base-content/60">
                    {t('click_to_upload_banner')}
                  </span>
                </div>
              )}
              <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                <span className="text-white text-sm font-medium">
                  {t('change_banner')}
                </span>
              </div>
            </div>
            <input
              ref={bannerInputRef}
              type="file"
              accept="image/jpeg,image/png,image/webp"
              className="hidden"
              onChange={handleBannerChange}
            />
            <p className="text-xs text-base-content/60 mt-2">
              {t('banner_requirements')}
            </p>
            {errors.banner && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {errors.banner}
                </span>
              </label>
            )}
          </div>

          <div className="card-actions justify-end mt-6">
            <button className="btn btn-primary" onClick={handleNext}>
              {tCommon('next')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
