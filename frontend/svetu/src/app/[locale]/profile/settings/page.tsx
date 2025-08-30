'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { useState, useEffect, useCallback } from 'react';
import { useRouter } from '@/i18n/routing';
import { Link } from '@/i18n/routing';
import config from '@/config';
import { PageTransition } from '@/components/ui/PageTransition';

interface PrivacySettings {
  allow_contact_requests: boolean;
  allow_messages_from_contacts_only: boolean;
}

export default function ProfileSettingsPage() {
  const { user, isAuthenticated, isLoading } = useAuth();
  const t = useTranslations('profile');
  const router = useRouter();
  const [mounted, setMounted] = useState(false);
  const [privacySettings, setPrivacySettings] = useState<PrivacySettings>({
    allow_contact_requests: true,
    allow_messages_from_contacts_only: false,
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (mounted && !isLoading && !isAuthenticated) {
      router.push('/');
    }
  }, [mounted, isAuthenticated, isLoading, router]);

  useEffect(() => {
    if (isAuthenticated && user) {
      fetchPrivacySettings();
    }
  }, [isAuthenticated, user, fetchPrivacySettings]);

  useEffect(() => {
    if (success) {
      const timer = setTimeout(() => setSuccess(false), 3000);
      return () => clearTimeout(timer);
    }
  }, [success]);

  const fetchPrivacySettings = useCallback(async () => {
    try {
      setLoading(true);
      const response = await fetch(
        `${config.getApiUrl()}/api/v1/users/privacy-settings`,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        if (data.data) {
          setPrivacySettings(data.data);
        }
      }
    } catch (err) {
      console.error('Error fetching privacy settings:', err);
      setError(t('settings.privacy.error.fetchFailed'));
    } finally {
      setLoading(false);
    }
  }, [t]);

  const savePrivacySettings = async () => {
    try {
      setSaving(true);
      setError(null);

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/users/privacy-settings`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
          },
          body: JSON.stringify(privacySettings),
        }
      );

      if (response.ok) {
        setSuccess(true);
      } else {
        throw new Error('Failed to save settings');
      }
    } catch (err) {
      console.error('Error saving privacy settings:', err);
      setError(t('settings.privacy.error.saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handlePrivacyChange = (field: keyof PrivacySettings) => {
    setPrivacySettings((prev) => ({
      ...prev,
      [field]: !prev[field],
    }));
  };

  if (!mounted || isLoading || loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <PageTransition mode="slideUp">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          {/* Breadcrumbs */}
          <div className="breadcrumbs mb-6">
            <ul>
              <li>
                <Link href="/profile">{t('title')}</Link>
              </li>
              <li>{t('settings.title')}</li>
            </ul>
          </div>

          <h1 className="text-3xl font-bold mb-8">{t('settings.title')}</h1>

          {/* Success/Error Messages */}
          {error && (
            <div className="alert alert-error mb-4">
              <span>{error}</span>
              <button
                onClick={() => setError(null)}
                className="btn btn-ghost btn-xs"
              >
                ✕
              </button>
            </div>
          )}

          {success && (
            <div className="alert alert-success mb-4">
              <span>{t('settings.messages.saveSuccess')}</span>
            </div>
          )}

          {/* Settings Tabs */}
          <div className="tabs tabs-boxed mb-6">
            <a className="tab tab-active">{t('settings.tabs.privacy')}</a>
            <a className="tab" onClick={() => alert(t('settings.comingSoon'))}>
              {t('settings.tabs.notifications')}
            </a>
            <a className="tab" onClick={() => alert(t('settings.comingSoon'))}>
              {t('settings.tabs.security')}
            </a>
          </div>

          {/* Privacy Settings */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title mb-4">{t('settings.privacy.title')}</h2>
              <p className="text-base-content/70 mb-6">
                {t('settings.privacy.description')}
              </p>

              <div className="space-y-6">
                {/* Contact Requests */}
                <div className="form-control">
                  <label className="label cursor-pointer justify-start gap-4">
                    <input
                      type="checkbox"
                      className="toggle toggle-primary"
                      checked={privacySettings.allow_contact_requests}
                      onChange={() =>
                        handlePrivacyChange('allow_contact_requests')
                      }
                    />
                    <div className="flex-1">
                      <span className="label-text font-medium">
                        {t('settings.privacy.allowContactRequests')}
                      </span>
                      <p className="text-sm text-base-content/60 mt-1">
                        {t('settings.privacy.allowContactRequestsDesc')}
                      </p>
                    </div>
                  </label>
                </div>

                {/* Messages from Contacts Only */}
                <div className="form-control">
                  <label className="label cursor-pointer justify-start gap-4">
                    <input
                      type="checkbox"
                      className="toggle toggle-primary"
                      checked={
                        privacySettings.allow_messages_from_contacts_only
                      }
                      onChange={() =>
                        handlePrivacyChange('allow_messages_from_contacts_only')
                      }
                    />
                    <div className="flex-1">
                      <span className="label-text font-medium">
                        {t('settings.privacy.messagesFromContactsOnly')}
                      </span>
                      <p className="text-sm text-base-content/60 mt-1">
                        {t('settings.privacy.messagesFromContactsOnlyDesc')}
                      </p>
                    </div>
                  </label>
                </div>

                <div className="divider"></div>

                {/* Location Privacy */}
                <div className="space-y-4">
                  <h3 className="font-semibold">
                    {t('settings.privacy.locationPrivacy')}
                  </h3>
                  <p className="text-sm text-base-content/60">
                    {t('settings.privacy.locationPrivacyDesc')}
                  </p>

                  <div className="form-control">
                    <label className="label cursor-pointer">
                      <span className="label-text">
                        {t('settings.privacy.locationExact')}
                      </span>
                      <input
                        type="radio"
                        name="location-privacy"
                        className="radio"
                        disabled
                      />
                    </label>
                  </div>

                  <div className="form-control">
                    <label className="label cursor-pointer">
                      <span className="label-text">
                        {t('settings.privacy.locationArea')}
                      </span>
                      <input
                        type="radio"
                        name="location-privacy"
                        className="radio"
                        defaultChecked
                        disabled
                      />
                    </label>
                  </div>

                  <div className="form-control">
                    <label className="label cursor-pointer">
                      <span className="label-text">
                        {t('settings.privacy.locationCity')}
                      </span>
                      <input
                        type="radio"
                        name="location-privacy"
                        className="radio"
                        disabled
                      />
                    </label>
                  </div>

                  <div className="form-control">
                    <label className="label cursor-pointer">
                      <span className="label-text">
                        {t('settings.privacy.locationHidden')}
                      </span>
                      <input
                        type="radio"
                        name="location-privacy"
                        className="radio"
                        disabled
                      />
                    </label>
                  </div>

                  <div className="alert alert-info">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      className="stroke-info shrink-0 w-6 h-6"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      ></path>
                    </svg>
                    <span className="text-sm">
                      {t('settings.privacy.locationNote')}
                    </span>
                  </div>
                </div>
              </div>

              <div className="card-actions justify-end mt-6">
                <button
                  className={`btn btn-primary ${saving ? 'loading' : ''}`}
                  onClick={savePrivacySettings}
                  disabled={saving}
                >
                  {saving
                    ? t('settings.actions.saving')
                    : t('settings.actions.save')}
                </button>
              </div>
            </div>
          </div>

          {/* Account Management */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h2 className="card-title mb-4 text-error">
                {t('settings.danger.title')}
              </h2>
              <p className="text-base-content/70 mb-6">
                {t('settings.danger.description')}
              </p>

              <div className="space-y-4">
                <button className="btn btn-outline btn-warning" disabled>
                  {t('settings.danger.deactivateAccount')}
                </button>
                <button className="btn btn-outline btn-error" disabled>
                  {t('settings.danger.deleteAccount')}
                </button>
              </div>

              <div className="alert alert-warning mt-4">
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
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                  />
                </svg>
                <span className="text-sm">{t('settings.danger.warning')}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </PageTransition>
  );
}
