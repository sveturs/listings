'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from '@/i18n/routing';
import { Link } from '@/i18n/routing';
import config from '@/config';
import { PageTransition } from '@/components/ui/PageTransition';

interface Contact {
  id: number;
  user_id: number;
  contact_user_id: number;
  status: 'pending' | 'accepted' | 'blocked';
  created_at: string;
  updated_at: string;
  notes?: string;
  contact_user?: {
    id: number;
    name: string;
    email: string;
    picture_url?: string;
    city?: string;
    country?: string;
  };
}

export default function ContactsPage() {
  const { user, isAuthenticated, isLoading } = useAuth();
  const t = useTranslations('profile');
  const tCommon = useTranslations('common');
  const router = useRouter();

  const [mounted, setMounted] = useState(false);
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [pendingRequests, setPendingRequests] = useState<Contact[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<
    'contacts' | 'pending' | 'blocked'
  >('contacts');
  const [searchQuery, setSearchQuery] = useState('');
  const [showAddModal, setShowAddModal] = useState(false);
  const [addContactEmail, setAddContactEmail] = useState('');
  const [addingContact, setAddingContact] = useState(false);
  const [updatingStatus, setUpdatingStatus] = useState<number | null>(null);
  const [removingContact, setRemovingContact] = useState<number | null>(null);

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
      fetchContacts();
    }
  }, [isAuthenticated, user, fetchContacts]);

  const fetchContacts = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams();
      if (activeTab === 'blocked') {
        params.append('status', 'blocked');
      } else if (activeTab === 'pending') {
        params.append('status', 'pending');
      } else {
        params.append('status', 'accepted');
      }

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/contacts?${params}`,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        if (activeTab === 'pending') {
          setPendingRequests(data.data?.contacts || []);
        } else {
          setContacts(data.data?.contacts || []);
        }
      } else {
        throw new Error('Failed to fetch contacts');
      }
    } catch (err) {
      console.error('Error fetching contacts:', err);
      setError(t('contacts.error.fetchFailed'));
    } finally {
      setLoading(false);
    }
  }, [activeTab, t]);

  const handleAddContact = async () => {
    if (!addContactEmail.trim()) return;

    try {
      setAddingContact(true);
      setError(null);

      // First find user by email
      const searchResponse = await fetch(
        `${config.getApiUrl()}/api/v1/users/search?email=${encodeURIComponent(addContactEmail)}`,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
          },
        }
      );

      if (!searchResponse.ok) {
        throw new Error(t('contacts.error.userNotFound'));
      }

      const searchData = await searchResponse.json();
      const targetUser = searchData.data;

      if (!targetUser || !targetUser.id) {
        throw new Error(t('contacts.error.userNotFound'));
      }

      // Add contact
      const response = await fetch(`${config.getApiUrl()}/api/v1/contacts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('access_token')}`,
        },
        body: JSON.stringify({
          contact_user_id: targetUser.id,
        }),
      });

      if (response.ok) {
        setShowAddModal(false);
        setAddContactEmail('');
        fetchContacts();
      } else {
        const errorData = await response.json();
        throw new Error(errorData.message || t('contacts.error.addFailed'));
      }
    } catch (err: any) {
      setError(err.message || t('contacts.error.addFailed'));
    } finally {
      setAddingContact(false);
    }
  };

  const handleUpdateStatus = async (
    contactId: number,
    newStatus: 'accepted' | 'blocked'
  ) => {
    try {
      setUpdatingStatus(contactId);

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/contacts/${contactId}/status`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
          },
          body: JSON.stringify({ status: newStatus }),
        }
      );

      if (response.ok) {
        fetchContacts();
      } else {
        throw new Error('Failed to update contact status');
      }
    } catch (err) {
      console.error('Error updating contact status:', err);
      setError(t('contacts.error.updateFailed'));
    } finally {
      setUpdatingStatus(null);
    }
  };

  const handleRemoveContact = async (contactId: number) => {
    if (!confirm(t('contacts.confirmRemove'))) return;

    try {
      setRemovingContact(contactId);

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/contacts/${contactId}`,
        {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
          },
        }
      );

      if (response.ok) {
        fetchContacts();
      } else {
        throw new Error('Failed to remove contact');
      }
    } catch (err) {
      console.error('Error removing contact:', err);
      setError(t('contacts.error.removeFailed'));
    } finally {
      setRemovingContact(null);
    }
  };

  const getFilteredContacts = () => {
    const list = activeTab === 'pending' ? pendingRequests : contacts;
    if (!searchQuery) return list;

    return list.filter(
      (contact) =>
        contact.contact_user?.name
          ?.toLowerCase()
          .includes(searchQuery.toLowerCase()) ||
        contact.contact_user?.email
          ?.toLowerCase()
          .includes(searchQuery.toLowerCase())
    );
  };

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
    return null;
  }

  const filteredContacts = getFilteredContacts();

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
              <li>{t('contacts.title')}</li>
            </ul>
          </div>

          <div className="flex justify-between items-center mb-8">
            <h1 className="text-3xl font-bold">{t('contacts.title')}</h1>
            <button
              className="btn btn-primary"
              onClick={() => setShowAddModal(true)}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-5 w-5 mr-2"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 4v16m8-8H4"
                />
              </svg>
              {t('contacts.addContact')}
            </button>
          </div>

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

          {/* Tabs */}
          <div className="tabs tabs-boxed mb-6">
            <a
              className={`tab ${activeTab === 'contacts' ? 'tab-active' : ''}`}
              onClick={() => setActiveTab('contacts')}
            >
              {t('contacts.tabs.contacts')}
              {contacts.length > 0 && (
                <span className="badge badge-sm ml-2">{contacts.length}</span>
              )}
            </a>
            <a
              className={`tab ${activeTab === 'pending' ? 'tab-active' : ''}`}
              onClick={() => setActiveTab('pending')}
            >
              {t('contacts.tabs.pending')}
              {pendingRequests.length > 0 && (
                <span className="badge badge-sm badge-warning ml-2">
                  {pendingRequests.length}
                </span>
              )}
            </a>
            <a
              className={`tab ${activeTab === 'blocked' ? 'tab-active' : ''}`}
              onClick={() => setActiveTab('blocked')}
            >
              {t('contacts.tabs.blocked')}
            </a>
          </div>

          {/* Search */}
          <div className="form-control mb-6">
            <input
              type="text"
              placeholder={t('contacts.searchPlaceholder')}
              className="input input-bordered"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>

          {/* Contacts List */}
          {loading ? (
            <div className="flex justify-center py-8">
              <span className="loading loading-spinner loading-lg"></span>
            </div>
          ) : filteredContacts.length === 0 ? (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body text-center py-12">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-24 w-24 mx-auto text-base-300 mb-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={1}
                    d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                  />
                </svg>
                <p className="text-lg text-base-content/60">
                  {activeTab === 'pending'
                    ? t('contacts.noPending')
                    : activeTab === 'blocked'
                      ? t('contacts.noBlocked')
                      : t('contacts.noContacts')}
                </p>
              </div>
            </div>
          ) : (
            <div className="grid gap-4">
              {filteredContacts.map((contact) => (
                <div key={contact.id} className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-4">
                        {contact.contact_user?.picture_url ? (
                          <div className="avatar">
                            <div className="w-12 rounded-full">
                              <img
                                src={contact.contact_user.picture_url}
                                alt={contact.contact_user.name}
                              />
                            </div>
                          </div>
                        ) : (
                          <div className="avatar placeholder">
                            <div className="bg-neutral-focus text-neutral-content rounded-full w-12">
                              <span className="text-xl">
                                {contact.contact_user?.name
                                  ?.charAt(0)
                                  .toUpperCase()}
                              </span>
                            </div>
                          </div>
                        )}
                        <div>
                          <h3 className="font-semibold">
                            {contact.contact_user?.name}
                          </h3>
                          <p className="text-sm text-base-content/60">
                            {contact.contact_user?.email}
                          </p>
                          {contact.contact_user?.city && (
                            <p className="text-sm text-base-content/50">
                              {contact.contact_user.city}
                              {contact.contact_user.country &&
                                `, ${contact.contact_user.country}`}
                            </p>
                          )}
                        </div>
                      </div>

                      <div className="flex gap-2">
                        {activeTab === 'pending' ? (
                          <>
                            <button
                              className={`btn btn-success btn-sm ${updatingStatus === contact.id ? 'loading' : ''}`}
                              onClick={() =>
                                handleUpdateStatus(contact.id, 'accepted')
                              }
                              disabled={updatingStatus === contact.id}
                            >
                              {t('contacts.accept')}
                            </button>
                            <button
                              className={`btn btn-error btn-sm ${updatingStatus === contact.id ? 'loading' : ''}`}
                              onClick={() =>
                                handleUpdateStatus(contact.id, 'blocked')
                              }
                              disabled={updatingStatus === contact.id}
                            >
                              {t('contacts.decline')}
                            </button>
                          </>
                        ) : activeTab === 'blocked' ? (
                          <button
                            className={`btn btn-sm ${updatingStatus === contact.id ? 'loading' : ''}`}
                            onClick={() =>
                              handleUpdateStatus(contact.id, 'accepted')
                            }
                            disabled={updatingStatus === contact.id}
                          >
                            {t('contacts.unblock')}
                          </button>
                        ) : (
                          <>
                            <Link
                              href={`/chat?user=${contact.contact_user_id}`}
                              className="btn btn-ghost btn-sm"
                            >
                              {t('contacts.message')}
                            </Link>
                            <button
                              className={`btn btn-ghost btn-sm text-error ${removingContact === contact.id ? 'loading' : ''}`}
                              onClick={() => handleRemoveContact(contact.id)}
                              disabled={removingContact === contact.id}
                            >
                              {t('contacts.remove')}
                            </button>
                          </>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Add Contact Modal */}
      {showAddModal && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg mb-4">
              {t('contacts.addContact')}
            </h3>

            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('contacts.enterEmail')}</span>
              </label>
              <input
                type="email"
                placeholder={t('contacts.emailPlaceholder')}
                className="input input-bordered"
                value={addContactEmail}
                onChange={(e) => setAddContactEmail(e.target.value)}
              />
            </div>

            <div className="modal-action">
              <button
                className="btn"
                onClick={() => {
                  setShowAddModal(false);
                  setAddContactEmail('');
                  setError(null);
                }}
                disabled={addingContact}
              >
                {tCommon('cancel')}
              </button>
              <button
                className={`btn btn-primary ${addingContact ? 'loading' : ''}`}
                onClick={handleAddContact}
                disabled={addingContact || !addContactEmail.trim()}
              >
                {addingContact ? t('contacts.adding') : tCommon('add')}
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button onClick={() => setShowAddModal(false)}>close</button>
          </form>
        </dialog>
      )}
    </PageTransition>
  );
}
