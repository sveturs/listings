'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { useState } from 'react';
import { redirect } from 'next/navigation';
import Image from 'next/image';

export default function ProfilePage() {
  const { user, isAuthenticated, isLoading, updateProfile } = useAuth();
  const t = useTranslations('profile');
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    phone: '',
    city: '',
    country: '',
  });

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    redirect('/');
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
    const updated = await updateProfile(formData);
    if (updated) {
      setIsEditing(false);
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">{t('title')}</h1>

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <div className="flex items-center mb-6">
              {user?.picture_url && (
                <div className="avatar mr-4">
                  <div className="w-24 h-24 rounded-full relative">
                    <Image
                      src={user.picture_url}
                      alt={user.name}
                      fill
                      className="rounded-full object-cover"
                    />
                  </div>
                </div>
              )}
              <div>
                <h2 className="card-title">{user?.name}</h2>
                <p className="text-base-content/70">{user?.email}</p>
              </div>
            </div>

            <div className="divider"></div>

            {isEditing ? (
              <div className="space-y-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">{t('fields.name')}</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={formData.name}
                    onChange={(e) =>
                      setFormData({ ...formData, name: e.target.value })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">{t('fields.phone')}</span>
                  </label>
                  <input
                    type="tel"
                    className="input input-bordered"
                    value={formData.phone}
                    onChange={(e) =>
                      setFormData({ ...formData, phone: e.target.value })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">{t('fields.city')}</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={formData.city}
                    onChange={(e) =>
                      setFormData({ ...formData, city: e.target.value })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">{t('fields.country')}</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={formData.country}
                    onChange={(e) =>
                      setFormData({ ...formData, country: e.target.value })
                    }
                  />
                </div>

                <div className="card-actions justify-end mt-6">
                  <button className="btn btn-ghost" onClick={handleCancel}>
                    {t('actions.cancel')}
                  </button>
                  <button className="btn btn-primary" onClick={handleSave}>
                    {t('actions.save')}
                  </button>
                </div>
              </div>
            ) : (
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-base-content/70">
                    {t('fields.name')}
                  </p>
                  <p className="font-medium">{user?.name}</p>
                </div>

                <div>
                  <p className="text-sm text-base-content/70">
                    {t('fields.email')}
                  </p>
                  <p className="font-medium">{user?.email}</p>
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
