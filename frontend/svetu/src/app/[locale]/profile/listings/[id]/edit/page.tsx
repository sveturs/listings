'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { useState, useEffect, useCallback } from 'react';
import { useParams } from 'next/navigation';
import { useRouter } from '@/i18n/routing';
import { apiClient } from '@/services/api-client';
import { Link } from '@/i18n/routing';

interface Listing {
  id: number;
  title: string;
  description: string;
  price: number;
  condition: string;
  status: string;
  city: string;
  country: string;
  category_id: number;
  user_id: number;
  category: {
    id: number;
    name: string;
    slug: string;
  };
}

interface ListingResponse {
  data: Listing;
  success: boolean;
}

export default function EditListingPage() {
  const { user, isAuthenticated, isLoading } = useAuth();
  const t = useTranslations('profile.listings.editListing');
  const router = useRouter();
  const params = useParams();
  const listingId = params?.id as string;

  const [mounted, setMounted] = useState(false);
  const [listing, setListing] = useState<Listing | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    price: 0,
    condition: 'used',
    city: '',
    country: '',
  });

  const fetchListing = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await apiClient.get<ListingResponse>(
        `/api/v1/marketplace/listings/${listingId}`
      );

      if (!response.error && response.data) {
        // Адаптируем под структуру ответа
        const listingData = response.data.data || response.data;

        // Check if user owns this listing
        if (listingData.user_id !== user?.id) {
          setError(t('permissionError'));
          return;
        }

        setListing(listingData);
        setFormData({
          title: listingData.title,
          description: listingData.description,
          price: listingData.price,
          condition: listingData.condition,
          city: listingData.city || '',
          country: listingData.country || '',
        });
      } else {
        console.error('Failed to load listing:', response);
        setError(t('listingNotFound'));
      }
    } catch (err) {
      console.error('Error fetching listing:', err);
      setError(t('fields.loadError'));
    } finally {
      setLoading(false);
    }
  }, [listingId, user?.id, t]);

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (mounted && !isLoading && !isAuthenticated) {
      router.push('/');
    }
  }, [mounted, isAuthenticated, isLoading, router]);

  useEffect(() => {
    if (mounted && isAuthenticated && listingId) {
      fetchListing();
    }
  }, [mounted, isAuthenticated, listingId, fetchListing]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!listing) return;

    try {
      setSaving(true);
      setError(null);

      const updateData = {
        title: formData.title.trim(),
        description: formData.description.trim(),
        price: formData.price,
        condition: formData.condition,
        city: formData.city.trim(),
        country: formData.country.trim(),
        category_id: listing.category_id,
      };

      const response = await apiClient.put(
        `/api/v1/marketplace/listings/${listingId}`,
        updateData
      );

      if (!response.error) {
        // Redirect back to listings or show success
        router.push('/profile/listings');
      } else {
        setError(t('updateError'));
      }
    } catch (err) {
      console.error('Error updating listing:', err);
      setError('Failed to update listing');
    } finally {
      setSaving(false);
    }
  };

  const handleInputChange = (field: string, value: string | number) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));
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

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          <div className="alert alert-error">
            <span>{error}</span>
          </div>
          <div className="mt-4">
            <Link href="/profile/listings" className="btn btn-ghost">
              ← {t('back')}
            </Link>
          </div>
        </div>
      </div>
    );
  }

  if (!listing) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          <div className="alert alert-warning">
            <span>{t('listingNotFound')}</span>
          </div>
          <div className="mt-4">
            <Link href="/profile/listings" className="btn btn-ghost">
              ← {t('back')}
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-4 mb-4">
            <Link href="/profile/listings" className="btn btn-ghost btn-sm">
              ← {t('back')}
            </Link>
          </div>
          <h1 className="text-3xl font-bold">{t('title')}</h1>
          <p className="text-base-content/70 mt-2">{t('description')}</p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Title */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {t('fields.title')} *
              </span>
            </label>
            <input
              type="text"
              className="input input-bordered w-full"
              value={formData.title}
              onChange={(e) => handleInputChange('title', e.target.value)}
              required
              maxLength={100}
            />
            <label className="label">
              <span className="label-text-alt">
                {formData.title.length}/100
              </span>
            </label>
          </div>

          {/* Description */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {t('fields.description')} *
              </span>
            </label>
            <textarea
              className="textarea textarea-bordered h-32"
              value={formData.description}
              onChange={(e) => handleInputChange('description', e.target.value)}
              required
              maxLength={1000}
            />
            <label className="label">
              <span className="label-text-alt">
                {formData.description.length}/1000
              </span>
            </label>
          </div>

          {/* Price */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {t('fields.price')} *
              </span>
            </label>
            <div className="input-group">
              <span>$</span>
              <input
                type="number"
                className="input input-bordered w-full"
                value={formData.price}
                onChange={(e) =>
                  handleInputChange('price', Number(e.target.value))
                }
                required
                min="0"
                step="0.01"
              />
            </div>
          </div>

          {/* Condition */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {t('fields.condition')} *
              </span>
            </label>
            <select
              className="select select-bordered w-full"
              value={formData.condition}
              onChange={(e) => handleInputChange('condition', e.target.value)}
              required
            >
              <option value="new">{t('condition.new')}</option>
              <option value="used">{t('condition.used')}</option>
              <option value="refurbished">{t('condition.refurbished')}</option>
            </select>
          </div>

          {/* Location */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  {t('fields.city')}
                </span>
              </label>
              <input
                type="text"
                className="input input-bordered w-full"
                value={formData.city}
                onChange={(e) => handleInputChange('city', e.target.value)}
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  {t('fields.country')}
                </span>
              </label>
              <input
                type="text"
                className="input input-bordered w-full"
                value={formData.country}
                onChange={(e) => handleInputChange('country', e.target.value)}
              />
            </div>
          </div>

          {/* Category (read-only for now) */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">
                {t('fields.category')}
              </span>
            </label>
            <input
              type="text"
              className="input input-bordered w-full"
              value={listing.category.name}
              disabled
            />
            <label className="label">
              <span className="label-text-alt">{t('fields.categoryHint')}</span>
            </label>
          </div>

          {/* Actions */}
          <div className="flex gap-4 pt-6">
            <Link href="/profile/listings" className="btn btn-ghost flex-1">
              {t('cancel')}
            </Link>
            <button
              type="submit"
              className={`btn btn-primary flex-1 ${saving ? 'loading' : ''}`}
              disabled={saving}
            >
              {saving ? t('saving') : t('saveChanges')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
