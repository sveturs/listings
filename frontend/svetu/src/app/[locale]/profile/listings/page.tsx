'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useState, useEffect, useCallback } from 'react';
import { useRouter } from '@/i18n/routing';
import { Link } from '@/i18n/routing';
import Image from 'next/image';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import { toast } from '@/utils/toast';

interface UserListing {
  id: number;
  title: string;
  description: string;
  price: number;
  condition: string;
  status: string;
  city: string;
  country: string;
  views_count: number;
  created_at: string;
  updated_at: string;
  category: {
    name: string;
    slug: string;
  };
  images?: Array<{
    id: number;
    file_path: string;
    file_name: string;
    file_size: number;
    content_type: string;
    is_main: boolean;
    storage_type: string;
    storage_bucket?: string;
    public_url?: string;
  }>;
}

// Backend возвращает обернутый ответ через utils.SuccessResponse
interface BackendResponse {
  success: boolean;
  data: {
    success: boolean;
    data: UserListing[];
    meta: {
      total: number;
      page: number;
      limit: number;
    };
  };
}

export default function MyListingsPage() {
  const { user, isAuthenticated, isLoading } = useAuth();
  const t = useTranslations('common');
  const router = useRouter();
  const [mounted, setMounted] = useState(false);
  const [listings, setListings] = useState<UserListing[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [updatingId, setUpdatingId] = useState<number | null>(null);

  const fetchMyListings = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await apiClient.get<BackendResponse>(
        `/api/v1/marketplace/listings?user_id=${user?.id}`
      );

      if (!response.error && response.data) {
        // Backend возвращает вложенную структуру через utils.SuccessResponse
        const backendData = response.data;

        if (backendData.success && backendData.data) {
          const listingsData = backendData.data;

          if (listingsData.success && Array.isArray(listingsData.data)) {
            console.log('Listings data:', listingsData.data);
            // Логируем изображения для отладки
            listingsData.data.forEach((listing, index) => {
              if (listing.images && listing.images.length > 0) {
                console.log(`Listing ${index} images:`, listing.images);
              }
            });
            setListings(listingsData.data);
          } else {
            console.error('Unexpected listings structure:', listingsData);
            setListings([]);
          }
        } else {
          console.error('Backend response not successful:', backendData);
          setListings([]);
        }
      } else {
        setError(response.error?.message || 'Failed to load listings');
      }
    } catch (err) {
      console.error('Error fetching listings:', err);
      setError('Failed to load your listings');
    } finally {
      setLoading(false);
    }
  }, [user?.id]);

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (mounted && !isLoading && !isAuthenticated) {
      router.push('/');
    }
  }, [mounted, isAuthenticated, isLoading, router]);

  useEffect(() => {
    if (mounted && isAuthenticated && user?.id) {
      fetchMyListings();
    }
  }, [mounted, isAuthenticated, user?.id, fetchMyListings]);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'active':
        return <span className="badge badge-success">Active</span>;
      case 'inactive':
        return <span className="badge badge-warning">Inactive</span>;
      case 'sold':
        return <span className="badge badge-neutral">Sold</span>;
      default:
        return <span className="badge badge-ghost">{status}</span>;
    }
  };

  const handleDelete = async (listingId: number, listingTitle: string) => {
    // Use browser confirm dialog
    if (window.confirm(`${t('confirmDelete')}: "${listingTitle}"`)) {
      await performDelete(listingId, listingTitle);
    }
  };

  const performDelete = async (listingId: number, _listingTitle: string) => {
    try {
      setDeletingId(listingId);

      const response = await apiClient.delete(
        `/api/v1/marketplace/listings/${listingId}`
      );

      if (!response.error) {
        toast.success(t('deleteSuccess'));
        // Remove the listing from the local state
        setListings(listings.filter((listing) => listing.id !== listingId));
      } else {
        toast.error(response.error.message || t('error'));
      }
    } catch (err) {
      console.error('Error deleting listing:', err);
      toast.error(t('error'));
    } finally {
      setDeletingId(null);
    }
  };

  const handleUpdateStatus = async (listingId: number, newStatus: string) => {
    try {
      setUpdatingId(listingId);

      const response = await apiClient.patch(
        `/api/v1/marketplace/listings/${listingId}`,
        { status: newStatus }
      );

      if (!response.error) {
        toast.success(`Listing marked as ${newStatus}`);
        // Update the listing status in local state
        setListings(
          listings.map((listing) =>
            listing.id === listingId
              ? { ...listing, status: newStatus }
              : listing
          )
        );
      } else {
        toast.error(
          response.error.message || 'Failed to update listing status'
        );
      }
    } catch (err) {
      console.error('Error updating listing status:', err);
      toast.error('Failed to update listing status. Please try again.');
    } finally {
      setUpdatingId(null);
    }
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

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold">My Listings</h1>
            <p className="text-base-content/70 mt-2">
              Manage your marketplace listings
            </p>
          </div>
          <Link href="/create-listing" className="btn btn-primary">
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
            Create New Listing
          </Link>
        </div>

        {/* Error State */}
        {error && (
          <div className="alert alert-error mb-6">
            <span>{error}</span>
            <button
              onClick={() => setError(null)}
              className="btn btn-ghost btn-xs"
            >
              ✕
            </button>
          </div>
        )}

        {/* Loading State */}
        {loading ? (
          <div className="flex justify-center py-12">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        ) : (
          <>
            {/* Listings Grid */}
            {listings.length === 0 ? (
              <div className="text-center py-12">
                <div className="mb-4">
                  <svg
                    className="mx-auto h-12 w-12 text-base-content/30"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                    />
                  </svg>
                </div>
                <h3 className="text-lg font-medium text-base-content/70 mb-2">
                  No listings yet
                </h3>
                <p className="text-base-content/50 mb-6">
                  Create your first listing to start selling
                </p>
                <Link href="/create-listing" className="btn btn-primary">
                  Create Your First Listing
                </Link>
              </div>
            ) : (
              <div className="space-y-4">
                {listings.map((listing) => (
                  <div
                    key={listing.id}
                    className={`card card-side bg-base-100 shadow-xl transition-opacity ${
                      deletingId === listing.id || updatingId === listing.id
                        ? 'opacity-50'
                        : ''
                    }`}
                  >
                    {/* Image */}
                    <figure className="w-48 h-36">
                      {listing.images &&
                      listing.images.length > 0 &&
                      listing.images[0].public_url ? (
                        <Image
                          src={listing.images[0].public_url}
                          alt={listing.title}
                          width={192}
                          height={144}
                          className="object-cover w-full h-full"
                        />
                      ) : (
                        <div className="w-full h-full bg-base-200 flex items-center justify-center">
                          <span className="text-base-content/30">No image</span>
                        </div>
                      )}
                    </figure>

                    {/* Content */}
                    <div className="card-body flex-1">
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <h2 className="card-title">
                            {listing.title}
                            {getStatusBadge(listing.status)}
                          </h2>
                          <p className="text-base-content/70 line-clamp-2">
                            {listing.description}
                          </p>

                          <div className="flex items-center gap-4 mt-2 text-sm text-base-content/60">
                            <span>${listing.price}</span>
                            <span>•</span>
                            <span>{listing.category.name}</span>
                            <span>•</span>
                            <span>{listing.condition}</span>
                            {listing.city && (
                              <>
                                <span>•</span>
                                <span>{listing.city}</span>
                              </>
                            )}
                          </div>

                          <div className="flex items-center gap-4 mt-2 text-xs text-base-content/50">
                            <span>{listing.views_count} views</span>
                            <span>•</span>
                            <span>
                              Created {formatDate(listing.created_at)}
                            </span>
                            {listing.updated_at !== listing.created_at && (
                              <>
                                <span>•</span>
                                <span>
                                  Updated {formatDate(listing.updated_at)}
                                </span>
                              </>
                            )}
                          </div>
                        </div>

                        {/* Actions */}
                        <div className="card-actions">
                          <div className="dropdown dropdown-end">
                            <div
                              tabIndex={0}
                              role="button"
                              className="btn btn-ghost btn-sm"
                            >
                              <svg
                                xmlns="http://www.w3.org/2000/svg"
                                className="h-4 w-4"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M12 5v.01M12 12v.01M12 19v.01"
                                />
                              </svg>
                            </div>
                            <ul
                              tabIndex={0}
                              className="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow"
                            >
                              <li>
                                <Link href={`/marketplace/${listing.id}`}>
                                  View Listing
                                </Link>
                              </li>
                              <li>
                                <Link
                                  href={`/profile/listings/${listing.id}/edit`}
                                >
                                  Edit Listing
                                </Link>
                              </li>
                              {listing.status !== 'sold' && (
                                <li>
                                  <button
                                    onClick={() =>
                                      handleUpdateStatus(
                                        listing.id,
                                        listing.status === 'active'
                                          ? 'inactive'
                                          : 'active'
                                      )
                                    }
                                    disabled={updatingId === listing.id}
                                  >
                                    {updatingId === listing.id ? (
                                      <span className="loading loading-spinner loading-xs"></span>
                                    ) : listing.status === 'active' ? (
                                      'Deactivate'
                                    ) : (
                                      'Activate'
                                    )}
                                  </button>
                                </li>
                              )}
                              <li>
                                <button
                                  className="text-warning"
                                  onClick={() =>
                                    handleUpdateStatus(listing.id, 'sold')
                                  }
                                  disabled={
                                    updatingId === listing.id ||
                                    listing.status === 'sold'
                                  }
                                >
                                  {updatingId === listing.id ? (
                                    <span className="loading loading-spinner loading-xs"></span>
                                  ) : (
                                    'Mark as Sold'
                                  )}
                                </button>
                              </li>
                              <li className="divider my-0"></li>
                              <li>
                                <button
                                  className="text-error"
                                  onClick={() =>
                                    handleDelete(listing.id, listing.title)
                                  }
                                  disabled={deletingId === listing.id}
                                >
                                  {deletingId === listing.id ? (
                                    <span className="loading loading-spinner loading-xs"></span>
                                  ) : (
                                    'Delete Listing'
                                  )}
                                </button>
                              </li>
                            </ul>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
