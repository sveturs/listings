'use client';

import { useTranslations } from 'next-intl';
import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../../../contexts/AuthContext';

interface UserContact {
  id: number;
  user_id: number;
  contact_user_id: number;
  status: 'pending' | 'accepted' | 'blocked';
  notes?: string;
  added_from_chat_id?: number;
  created_at: string;
  updated_at: string;
  contact_user?: {
    id: number;
    name: string;
    email: string;
  };
}

interface ContactsResponse {
  contacts: UserContact[];
  total_count: number;
  page: number;
  limit: number;
}

interface PrivacySettings {
  user_id: number;
  allow_contact_requests: boolean;
  allow_messages_from_contacts_only: boolean;
  created_at: string;
  updated_at: string;
}

export default function UserContactsPage() {
  const t = useTranslations('userContacts');
  const { user, isLoading: authLoading } = useAuth();
  const [mounted, setMounted] = useState(false);
  const [contacts, setContacts] = useState<UserContact[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<
    'all' | 'pending' | 'accepted' | 'blocked'
  >('all');
  const [privacySettings, setPrivacySettings] =
    useState<PrivacySettings | null>(null);
  const [showPrivacyModal, setShowPrivacyModal] = useState(false);

  const fetchContacts = useCallback(
    async (status: string = '') => {
      try {
        setLoading(true);
        let url = '/api/v1/contacts';
        if (status) {
          url += `?status=${encodeURIComponent(status)}`;
        }

        // Получаем JWT токен через tokenManager
        const { tokenManager } = await import('@/utils/tokenManager');
        const accessToken = await tokenManager.getAccessToken();

        const headers: HeadersInit = {
          'Content-Type': 'application/json',
        };

        if (accessToken) {
          headers['Authorization'] = `Bearer ${accessToken}`;
        }

        const response = await fetch(url, {
          credentials: 'include',
          headers,
        });

        if (!response.ok) {
          throw new Error('Failed to fetch contacts');
        }

        const result = await response.json();
        // Учитываем, что данные могут быть обернуты в поле data
        const data: ContactsResponse = result.data || result;
        setContacts(data.contacts || []);
        setError(''); // Очищаем ошибку при успешной загрузке
      } catch (err) {
        console.error('Failed to fetch contacts:', err);
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    },
    [] // Убираем зависимость от contacts
  );

  const fetchPrivacySettings = useCallback(async () => {
    try {
      // Получаем JWT токен через tokenManager
      const { tokenManager } = await import('@/utils/tokenManager');
      const accessToken = await tokenManager.getAccessToken();

      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (accessToken) {
        headers['Authorization'] = `Bearer ${accessToken}`;
      }

      const response = await fetch('/api/v1/contacts/privacy', {
        credentials: 'include',
        headers,
      });

      if (response.ok) {
        const result = await response.json();
        const data: PrivacySettings = result.data || result;
        setPrivacySettings(data);
      } else {
        // Если настройки приватности недоступны, используем дефолтные
        console.warn('Failed to fetch privacy settings, using defaults');
        setPrivacySettings({
          user_id: user!.id,
          allow_contact_requests: true,
          allow_messages_from_contacts_only: false,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        });
      }
    } catch (err) {
      console.error('Failed to fetch privacy settings:', err);
      // Используем дефолтные настройки при ошибке
      setPrivacySettings({
        user_id: user!.id,
        allow_contact_requests: true,
        allow_messages_from_contacts_only: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      });
    }
  }, [user]);

  useEffect(() => {
    setMounted(true);
  }, []);

  // Загружаем данные при изменении activeTab
  useEffect(() => {
    if (user && !authLoading && mounted) {
      const status = activeTab === 'all' ? '' : activeTab;
      fetchContacts(status);
    }
  }, [activeTab, user, authLoading, mounted, fetchContacts]);

  // Загружаем настройки приватности только один раз при монтировании
  useEffect(() => {
    if (user && !authLoading && mounted) {
      fetchPrivacySettings();
    }
  }, [user, authLoading, mounted, fetchPrivacySettings]);

  const updateContactStatus = async (
    contactUserID: number,
    status: 'accepted' | 'blocked',
    notes: string = ''
  ) => {
    try {
      // Получаем JWT токен через tokenManager
      const { tokenManager } = await import('@/utils/tokenManager');
      const accessToken = await tokenManager.getAccessToken();

      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (accessToken) {
        headers['Authorization'] = `Bearer ${accessToken}`;
      }

      // Добавляем CSRF токен для PUT запроса
      try {
        const { AuthService } = await import('@/services/auth');
        const csrfToken = await AuthService.getCsrfToken();
        if (csrfToken) {
          (headers as any)['X-CSRF-Token'] = csrfToken;
        }
      } catch (error) {
        console.warn('Failed to get CSRF token:', error);
      }

      const response = await fetch(`/api/v1/contacts/${contactUserID}/status`, {
        method: 'PUT',
        credentials: 'include',
        headers,
        body: JSON.stringify({ status, notes }),
      });

      if (!response.ok) {
        throw new Error('Failed to update contact status');
      }

      // Refresh contacts list
      const currentStatus = activeTab === 'all' ? '' : activeTab;
      fetchContacts(currentStatus);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    }
  };

  const removeContact = async (contactUserID: number) => {
    try {
      // Получаем JWT токен через tokenManager
      const { tokenManager } = await import('@/utils/tokenManager');
      const accessToken = await tokenManager.getAccessToken();

      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (accessToken) {
        headers['Authorization'] = `Bearer ${accessToken}`;
      }

      // Добавляем CSRF токен для DELETE запроса
      try {
        const { AuthService } = await import('@/services/auth');
        const csrfToken = await AuthService.getCsrfToken();
        if (csrfToken) {
          (headers as any)['X-CSRF-Token'] = csrfToken;
        }
      } catch (error) {
        console.warn('Failed to get CSRF token:', error);
      }

      const response = await fetch(`/api/v1/contacts/${contactUserID}`, {
        method: 'DELETE',
        credentials: 'include',
        headers,
      });

      if (!response.ok) {
        throw new Error('Failed to remove contact');
      }

      // Refresh contacts list
      const currentStatus = activeTab === 'all' ? '' : activeTab;
      fetchContacts(currentStatus);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    }
  };

  const updatePrivacySettings = async (settings: Partial<PrivacySettings>) => {
    try {
      // Получаем JWT токен через tokenManager
      const { tokenManager } = await import('@/utils/tokenManager');
      const accessToken = await tokenManager.getAccessToken();

      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (accessToken) {
        headers['Authorization'] = `Bearer ${accessToken}`;
      }

      // Добавляем CSRF токен для PUT запроса
      try {
        const { AuthService } = await import('@/services/auth');
        const csrfToken = await AuthService.getCsrfToken();
        if (csrfToken) {
          (headers as any)['X-CSRF-Token'] = csrfToken;
        }
      } catch (error) {
        console.warn('Failed to get CSRF token:', error);
      }

      const response = await fetch('/api/v1/contacts/privacy', {
        method: 'PUT',
        credentials: 'include',
        headers,
        body: JSON.stringify(settings),
      });

      if (!response.ok) {
        throw new Error('Failed to update privacy settings');
      }

      const result = await response.json();
      const data: PrivacySettings = result.data || result;
      setPrivacySettings(data);
      setShowPrivacyModal(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'accepted':
        return 'text-green-600 bg-green-100';
      case 'pending':
        return 'text-yellow-600 bg-yellow-100';
      case 'blocked':
        return 'text-red-600 bg-red-100';
      default:
        return 'text-gray-600 bg-gray-100';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'accepted':
        return t('status.accepted');
      case 'pending':
        return t('status.pending');
      case 'blocked':
        return t('status.blocked');
      default:
        return status;
    }
  };

  if (!mounted || authLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center min-h-[400px]">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">{t('auth.required')}</h1>
          <p>{t('auth.loginMessage')}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">{t('title')}</h1>
        <button
          onClick={() => setShowPrivacyModal(true)}
          className="btn btn-outline btn-sm"
        >
          {t('privacy.settings')}
        </button>
      </div>

      {/* Tabs */}
      <div className="tabs tabs-bordered mb-6">
        <button
          className={`tab ${activeTab === 'all' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('all')}
        >
          {t('tabs.all')}
        </button>
        <button
          className={`tab ${activeTab === 'pending' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('pending')}
        >
          {t('tabs.pending')}
        </button>
        <button
          className={`tab ${activeTab === 'accepted' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('accepted')}
        >
          {t('tabs.accepted')}
        </button>
        <button
          className={`tab ${activeTab === 'blocked' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('blocked')}
        >
          {t('tabs.blocked')}
        </button>
      </div>

      {error && (!contacts || contacts.length === 0) && (
        <div className="alert alert-error mb-4">
          <span>{error}</span>
        </div>
      )}

      {loading ? (
        <div className="flex justify-center py-8">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : (
        <div className="space-y-4">
          {!contacts || contacts.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <p>{t('empty.message')}</p>
            </div>
          ) : (
            contacts.map((contact) => (
              <div
                key={contact.id}
                className="card bg-base-100 shadow-sm border"
              >
                <div className="card-body">
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <h3 className="card-title">
                        {contact.contact_user?.name ||
                          `User #${contact.contact_user_id}`}
                      </h3>
                      <p className="text-gray-600">
                        {contact.contact_user?.email}
                      </p>
                      <div className="flex items-center gap-2 mt-2">
                        <span
                          className={`badge ${getStatusColor(contact.status)}`}
                        >
                          {getStatusText(contact.status)}
                        </span>
                        <span className="text-sm text-gray-500">
                          {new Date(contact.created_at).toLocaleDateString()}
                        </span>
                      </div>
                      {contact.notes && (
                        <p className="text-sm text-gray-600 mt-2">
                          {contact.notes}
                        </p>
                      )}
                    </div>
                    <div className="flex gap-2">
                      {contact.status === 'pending' && (
                        <>
                          {/* Показываем кнопки принять/заблокировать только если запрос отправлен НАМ */}
                          {contact.contact_user_id === user?.id && (
                            <>
                              <button
                                onClick={() =>
                                  updateContactStatus(
                                    contact.user_id,
                                    'accepted'
                                  )
                                }
                                className="btn btn-success btn-sm"
                              >
                                {t('actions.accept')}
                              </button>
                              <button
                                onClick={() =>
                                  updateContactStatus(
                                    contact.user_id,
                                    'blocked'
                                  )
                                }
                                className="btn btn-error btn-sm"
                              >
                                {t('actions.block')}
                              </button>
                            </>
                          )}
                        </>
                      )}
                      {contact.status === 'blocked' && (
                        <button
                          onClick={() => {
                            // Для разблокировки нужно определить правильные ID
                            const otherUserId =
                              contact.user_id === user?.id
                                ? contact.contact_user_id
                                : contact.user_id;
                            updateContactStatus(otherUserId, 'accepted');
                          }}
                          className="btn btn-success btn-sm"
                        >
                          {t('actions.unblock')}
                        </button>
                      )}
                      <button
                        onClick={() => {
                          // Определяем ID другого пользователя
                          const otherUserId =
                            contact.user_id === user?.id
                              ? contact.contact_user_id
                              : contact.user_id;
                          removeContact(otherUserId);
                        }}
                        className="btn btn-outline btn-error btn-sm"
                      >
                        {t('actions.remove')}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      )}

      {/* Privacy Settings Modal */}
      {showPrivacyModal && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg mb-4">{t('privacy.title')}</h3>

            <div className="space-y-4">
              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">
                    {t('privacy.allowRequests')}
                  </span>
                  <input
                    type="checkbox"
                    className="checkbox"
                    checked={privacySettings?.allow_contact_requests ?? true}
                    onChange={(e) =>
                      setPrivacySettings((prev) =>
                        prev
                          ? {
                              ...prev,
                              allow_contact_requests: e.target.checked,
                            }
                          : null
                      )
                    }
                  />
                </label>
              </div>

              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">
                    {t('privacy.contactsOnly')}
                  </span>
                  <input
                    type="checkbox"
                    className="checkbox"
                    checked={
                      privacySettings?.allow_messages_from_contacts_only ??
                      false
                    }
                    onChange={(e) =>
                      setPrivacySettings((prev) =>
                        prev
                          ? {
                              ...prev,
                              allow_messages_from_contacts_only:
                                e.target.checked,
                            }
                          : null
                      )
                    }
                  />
                </label>
              </div>
            </div>

            <div className="modal-action">
              <button
                onClick={() => setShowPrivacyModal(false)}
                className="btn btn-ghost"
              >
                {t('actions.cancel')}
              </button>
              <button
                onClick={() => updatePrivacySettings(privacySettings || {})}
                className="btn btn-primary"
              >
                {t('actions.save')}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
