'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { useState, useEffect, useCallback } from 'react';
import { useRouter } from '@/i18n/routing';
import Image from 'next/image';
import {
  validateProfileData,
  hasFormChanges,
  type ValidationError,
} from '@/utils/validation';

export default function ProfilePage() {
  const {
    user,
    isAuthenticated,
    isLoading,
    isUpdatingProfile,
    error,
    updateProfile,
    clearError,
  } = useAuth();
  const t = useTranslations('profile');
  const router = useRouter();
  const [mounted, setMounted] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [updateError, setUpdateError] = useState<string | null>(null);
  const [updateSuccess, setUpdateSuccess] = useState(false);
  const [validationErrors, setValidationErrors] = useState<ValidationError[]>(
    []
  );
  const [formData, setFormData] = useState({
    name: '',
    phone: '',
    city: '',
    country: '',
  });

  // useCallback hooks must be before any conditional returns
  const getFieldError = useCallback(
    (field: string) => {
      return validationErrors.find((error) => error.field === field)?.message;
    },
    [validationErrors]
  );

  const hasChanges = useCallback(() => {
    if (!user) return false;
    return hasFormChanges(
      {
        name: user.name,
        phone: user.phone || '',
        city: user.city || '',
        country: user.country || '',
      },
      formData,
      ['name', 'phone', 'city', 'country']
    );
  }, [user, formData]);

  // Ensure component is mounted before rendering dynamic content
  useEffect(() => {
    setMounted(true);
  }, []);

  // Handle authentication redirect
  useEffect(() => {
    if (mounted && !isLoading && !isAuthenticated) {
      router.push('/');
    }
  }, [mounted, isAuthenticated, isLoading, router]);

  // Clear messages after delay
  useEffect(() => {
    if (updateSuccess) {
      const timer = setTimeout(() => setUpdateSuccess(false), 3000);
      return () => clearTimeout(timer);
    }
  }, [updateSuccess]);

  useEffect(() => {
    if (updateError) {
      const timer = setTimeout(() => setUpdateError(null), 5000);
      return () => clearTimeout(timer);
    }
  }, [updateError]);

  // Always show loading state during SSR and initial mount
  if (!mounted || isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null; // Redirect handled by useEffect
  }

  const handleEdit = () => {
    if (user) {
      setFormData({
        name: user.name,
        phone: user.phone || '',
        city: user.city || '',
        country: user.country || '',
      });
      setIsEditing(true);
    }
  };

  const handleSave = async () => {
    setUpdateError(null);
    setUpdateSuccess(false);
    setValidationErrors([]);

    // Validate form data
    const errors = validateProfileData(formData);
    if (errors.length > 0) {
      setValidationErrors(errors);
      return;
    }

    try {
      const updated = await updateProfile(formData);
      if (updated) {
        setIsEditing(false);
        setUpdateSuccess(true);
      } else {
        setUpdateError(t('errors.updateFailed'));
      }
    } catch {
      setUpdateError(t('errors.updateFailed'));
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
    setUpdateError(null);
    setValidationErrors([]);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">{t('title')}</h1>

        {/* Success/Error Messages */}
        {(error || updateError) && (
          <div className="alert alert-error mb-4">
            <span>{error || updateError}</span>
            <button
              onClick={() => {
                clearError();
                setUpdateError(null);
              }}
              className="btn btn-ghost btn-xs"
            >
              ✕
            </button>
          </div>
        )}

        {updateSuccess && (
          <div className="alert alert-success mb-4">
            <span>{t('messages.updateSuccess')}</span>
          </div>
        )}

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <div className="flex items-center mb-6">
              {user?.picture_url && (
                <div className="avatar mr-4">
                  <div className="w-24 h-24 rounded-full relative">
                    <Image
                      src={user.picture_url}
                      alt={`Profile picture of ${user.name}`}
                      fill
                      className="rounded-full object-cover"
                      sizes="96px"
                      priority={false}
                      placeholder="empty"
                    />
                  </div>
                </div>
              )}
              <div>
                <h2 className="card-title">{user?.name || ''}</h2>
                <p className="text-base-content/70">{user?.email || ''}</p>
              </div>
            </div>

            <div className="divider"></div>

            {isEditing ? (
              <div>
                <fieldset className="fieldset">
                  <legend className="fieldset-legend">{t('title')}</legend>

                  <label className="label" htmlFor="name">
                    {t('fields.name')} <span className="text-error">*</span>
                  </label>
                  <input
                    id="name"
                    type="text"
                    className={`input ${
                      getFieldError('name') ? 'input-error' : ''
                    }`}
                    value={formData.name}
                    onChange={(e) =>
                      setFormData({ ...formData, name: e.target.value })
                    }
                    required
                  />
                  {getFieldError('name') && (
                    <p className="label text-error text-sm">
                      {t(getFieldError('name')!)}
                    </p>
                  )}

                  <label className="label" htmlFor="email">
                    {t('fields.email')}
                  </label>
                  <input
                    id="email"
                    type="email"
                    className="input input-disabled"
                    value={user?.email || ''}
                    disabled
                    readOnly
                  />

                  <label className="label" htmlFor="phone">
                    {t('fields.phone')}
                  </label>
                  <input
                    id="phone"
                    type="tel"
                    className={`input ${
                      getFieldError('phone') ? 'input-error' : ''
                    }`}
                    value={formData.phone}
                    onChange={(e) =>
                      setFormData({ ...formData, phone: e.target.value })
                    }
                    placeholder="+1234567890"
                  />
                  {getFieldError('phone') && (
                    <p className="label text-error text-sm">
                      {t(getFieldError('phone')!)}
                    </p>
                  )}

                  <label className="label" htmlFor="city">
                    {t('fields.city')}
                  </label>
                  <input
                    id="city"
                    type="text"
                    className={`input ${
                      getFieldError('city') ? 'input-error' : ''
                    }`}
                    value={formData.city}
                    onChange={(e) =>
                      setFormData({ ...formData, city: e.target.value })
                    }
                  />
                  {getFieldError('city') && (
                    <p className="label text-error text-sm">
                      {t(getFieldError('city')!)}
                    </p>
                  )}

                  <label className="label" htmlFor="country">
                    {t('fields.country')}
                  </label>
                  <input
                    id="country"
                    type="text"
                    className={`input ${
                      getFieldError('country') ? 'input-error' : ''
                    }`}
                    value={formData.country}
                    onChange={(e) =>
                      setFormData({ ...formData, country: e.target.value })
                    }
                  />
                  {getFieldError('country') && (
                    <p className="label text-error text-sm">
                      {t(getFieldError('country')!)}
                    </p>
                  )}
                </fieldset>

                <div className="card-actions justify-end mt-6">
                  <button
                    className="btn btn-ghost"
                    onClick={handleCancel}
                    disabled={isUpdatingProfile}
                  >
                    {t('actions.cancel')}
                  </button>
                  <button
                    className={`btn btn-primary ${isUpdatingProfile ? 'loading' : ''}`}
                    onClick={handleSave}
                    disabled={isUpdatingProfile || !hasChanges()}
                  >
                    {isUpdatingProfile
                      ? t('actions.saving')
                      : t('actions.save')}
                  </button>
                </div>
              </div>
            ) : (
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-base-content/70">
                    {t('fields.name')}
                  </p>
                  <p className="font-medium">{user?.name || ''}</p>
                </div>

                <div>
                  <p className="text-sm text-base-content/70">
                    {t('fields.email')}
                  </p>
                  <p className="font-medium">{user?.email || ''}</p>
                </div>

                <div>
                  <p className="text-sm text-base-content/70">
                    {t('fields.phone')}
                  </p>
                  <p className="font-medium">
                    {user?.phone || t('notProvided')}
                  </p>
                </div>

                <div>
                  <p className="text-sm text-base-content/70">
                    {t('fields.city')}
                  </p>
                  <p className="font-medium">
                    {user?.city || t('notProvided')}
                  </p>
                </div>

                <div>
                  <p className="text-sm text-base-content/70">
                    {t('fields.country')}
                  </p>
                  <p className="font-medium">
                    {user?.country || t('notProvided')}
                  </p>
                </div>

                <div className="card-actions justify-end mt-6">
                  <button className="btn btn-primary" onClick={handleEdit}>
                    {t('actions.edit')}
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
