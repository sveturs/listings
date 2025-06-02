'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { useState, useEffect, useCallback } from 'react';
import { useRouter } from '@/i18n/routing';
import Image from 'next/image';
import FormField from '@/components/FormField';
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

        {/* User Profile Header */}
        <div className="flex items-center mb-8">
          {user?.picture_url && (
            <div className="avatar mr-6">
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
            <h2 className="text-2xl font-semibold">{user?.name || ''}</h2>
            <p className="text-base-content/70">{user?.email || ''}</p>
          </div>
        </div>

        {isEditing ? (
          <form className="space-y-6">
            <FormField
              label={t('fields.name')}
              required
              error={
                getFieldError('name') ? t(getFieldError('name')!) : undefined
              }
            >
              <input
                id="name"
                type="text"
                className={`input input-bordered w-full ${
                  getFieldError('name') ? 'input-error' : ''
                }`}
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
                required
              />
            </FormField>

            <div className="form-control">
              <label className="label" htmlFor="email">
                <span className="label-text">{t('fields.email')}</span>
              </label>
              <input
                id="email"
                type="email"
                className="input input-bordered w-full"
                value={user?.email || ''}
                disabled
                readOnly
              />
            </div>

            <div className="form-control">
              <label className="label" htmlFor="phone">
                <span className="label-text">{t('fields.phone')}</span>
              </label>
              <input
                id="phone"
                type="tel"
                className={`input input-bordered w-full ${
                  getFieldError('phone') ? 'input-error' : ''
                }`}
                value={formData.phone}
                onChange={(e) =>
                  setFormData({ ...formData, phone: e.target.value })
                }
                placeholder="+1234567890"
              />
              {getFieldError('phone') && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {t(getFieldError('phone')!)}
                  </span>
                </label>
              )}
            </div>

            <div className="form-control">
              <label className="label" htmlFor="city">
                <span className="label-text">{t('fields.city')}</span>
              </label>
              <input
                id="city"
                type="text"
                className={`input input-bordered w-full ${
                  getFieldError('city') ? 'input-error' : ''
                }`}
                value={formData.city}
                onChange={(e) =>
                  setFormData({ ...formData, city: e.target.value })
                }
              />
              {getFieldError('city') && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {t(getFieldError('city')!)}
                  </span>
                </label>
              )}
            </div>

            <div className="form-control">
              <label className="label" htmlFor="country">
                <span className="label-text">{t('fields.country')}</span>
              </label>
              <input
                id="country"
                type="text"
                className={`input input-bordered w-full ${
                  getFieldError('country') ? 'input-error' : ''
                }`}
                value={formData.country}
                onChange={(e) =>
                  setFormData({ ...formData, country: e.target.value })
                }
              />
              {getFieldError('country') && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {t(getFieldError('country')!)}
                  </span>
                </label>
              )}
            </div>

            <div className="flex gap-4 pt-4">
              <button
                type="button"
                className="btn btn-ghost"
                onClick={handleCancel}
                disabled={isUpdatingProfile}
              >
                {t('actions.cancel')}
              </button>
              <button
                type="button"
                className={`btn btn-primary ${isUpdatingProfile ? 'loading' : ''}`}
                onClick={handleSave}
                disabled={isUpdatingProfile || !hasChanges()}
              >
                {isUpdatingProfile ? t('actions.saving') : t('actions.save')}
              </button>
            </div>
          </form>
        ) : (
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      {t('fields.name')}
                    </span>
                  </label>
                  <p className="text-lg">{user?.name || ''}</p>
                </div>
              </div>

              <div>
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      {t('fields.email')}
                    </span>
                  </label>
                  <p className="text-lg">{user?.email || ''}</p>
                </div>
              </div>

              <div>
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      {t('fields.phone')}
                    </span>
                  </label>
                  <p className="text-lg text-base-content/70">
                    {user?.phone || t('notProvided')}
                  </p>
                </div>
              </div>

              <div>
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      {t('fields.city')}
                    </span>
                  </label>
                  <p className="text-lg text-base-content/70">
                    {user?.city || t('notProvided')}
                  </p>
                </div>
              </div>

              <div>
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      {t('fields.country')}
                    </span>
                  </label>
                  <p className="text-lg text-base-content/70">
                    {user?.country || t('notProvided')}
                  </p>
                </div>
              </div>
            </div>

            <div className="flex pt-6">
              <button className="btn btn-primary" onClick={handleEdit}>
                {t('actions.edit')}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
