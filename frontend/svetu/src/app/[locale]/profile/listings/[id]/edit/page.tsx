'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { useState, useEffect, useCallback } from 'react';
import { useParams } from 'next/navigation';
import { useRouter } from '@/i18n/routing';
import { apiClient } from '@/services/api-client';
import { Link } from '@/i18n/routing';
import { toast } from 'react-hot-toast';

// Import components
import { BasicInfoSection } from './components/BasicInfoSection';
import { ImagesSection } from './components/ImagesSection';
import { AttributesSection } from './components/AttributesSection';
import { LocationSection } from './components/LocationSection';
import { SEOSection } from './components/SEOSection';
import { PreviewCard } from './components/PreviewCard';

interface AttributeValue {
  attribute_id: number;
  value: string | number | boolean | string[];
}

interface Listing {
  id: number;
  title: string;
  description: string;
  price: number;
  condition: string;
  status: string;
  city: string;
  country: string;
  address_city?: string;
  address_country?: string;
  location: string;
  latitude?: number;
  longitude?: number;
  show_on_map: boolean;
  category_id: number;
  user_id: number;
  category: {
    id: number;
    name: string;
    slug: string;
  };
  images: Array<{
    id: number;
    file_path: string;
    file_name: string;
    is_main: boolean;
    public_url: string;
  }>;
  attributes?: AttributeValue[];
  metadata?: {
    seo?: {
      keywords: string;
      slug: string;
    };
  };
}

interface ListingResponse {
  data: Listing;
  success: boolean;
}

type TabType =
  | 'basic'
  | 'images'
  | 'attributes'
  | 'location'
  | 'seo'
  | 'preview';

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
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isDirty, setIsDirty] = useState(false);
  const [activeTab, setActiveTab] = useState<TabType>('basic');
  const [previewMode, setPreviewMode] = useState<'card' | 'list' | 'mobile'>(
    'card'
  );

  const [formData, setFormData] = useState({
    title: '',
    description: '',
    price: 0,
    condition: 'used',
    city: '',
    country: '',
    location: '',
    latitude: undefined as number | undefined,
    longitude: undefined as number | undefined,
    show_on_map: true,
  });

  const [images, setImages] = useState<Listing['images']>([]);
  const [attributes, setAttributes] = useState<AttributeValue[]>([]);
  const [seoData, setSeoData] = useState({
    keywords: '',
    slug: '',
  });

  const fetchListing = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await apiClient.get<ListingResponse>(
        `/api/v1/marketplace/listings/${listingId}`
      );

      if (!response.error && response.data) {
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
          city: listingData.address_city || listingData.city || '',
          country: listingData.address_country || listingData.country || '',
          location: listingData.location || '',
          latitude: listingData.latitude,
          longitude: listingData.longitude,
          show_on_map: listingData.show_on_map,
        });
        setImages(listingData.images || []);
        setAttributes(
          Array.isArray(listingData.attributes) ? listingData.attributes : []
        );

        // Load SEO data from metadata
        if (listingData.metadata && listingData.metadata.seo) {
          setSeoData(listingData.metadata.seo);
        }
      } else {
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

  // Auto-save draft every 30 seconds
  useEffect(() => {
    if (!isDirty || !listing) return;

    const timer = setTimeout(() => {
      saveDraft();
    }, 30000);

    return () => clearTimeout(timer);
  }, [isDirty, formData, attributes]);

  const saveDraft = async () => {
    // TODO: Implement draft saving
    console.log('Saving draft...');
  };

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.title.trim()) {
      newErrors.title = t('errors.titleRequired');
    } else if (formData.title.length < 10) {
      newErrors.title = t('errors.titleTooShort');
    }

    if (!formData.description.trim()) {
      newErrors.description = t('errors.descriptionRequired');
    } else if (formData.description.length < 20) {
      newErrors.description = t('errors.descriptionTooShort');
    }

    if (formData.price <= 0) {
      newErrors.price = t('errors.priceRequired');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!listing || !validateForm()) return;

    try {
      setSaving(true);
      setError(null);

      const updateData: any = {
        title: formData.title.trim(),
        description: formData.description.trim(),
        price: formData.price,
        condition: formData.condition,
        address_city: formData.city.trim(),
        address_country: formData.country.trim(),
        location: formData.location.trim(),
        latitude: formData.latitude,
        longitude: formData.longitude,
        show_on_map: formData.show_on_map,
        category_id: listing.category_id,
        attributes: attributes,
        images: images.map((img, index) => ({
          ...img,
          display_order: index,
        })),
        metadata: {
          seo: seoData,
        },
      };

      const response = await apiClient.put(
        `/api/v1/marketplace/listings/${listingId}`,
        updateData
      );

      if (!response.error && response.data) {
        toast.success(t('updateSuccess'));
        setIsDirty(false);

        // Update SEO data with the new slug from backend
        if (
          response.data.metadata?.seo?.slug &&
          response.data.metadata.seo.slug !== seoData.slug
        ) {
          const newSlug = response.data.metadata.seo.slug;
          setSeoData((prev) => ({
            ...prev,
            slug: newSlug,
          }));
          toast(t('seo.slugChanged', { suggestion: newSlug }));
        }

        // Don't redirect immediately to allow user to see the updated slug
        // User can navigate away manually or save again
      } else {
        setError(t('updateError'));
      }
    } catch (err) {
      console.error('Error updating listing:', err);
      setError(t('updateError'));
    } finally {
      setSaving(false);
    }
  };

  const handleInputChange = (
    field: string,
    value: string | number | boolean
  ) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));
    setIsDirty(true);
  };

  const handleImagesChange = (newImages: Listing['images']) => {
    setImages(newImages);
    setIsDirty(true);
  };

  const handleAttributesChange = (newAttributes: AttributeValue[]) => {
    setAttributes(newAttributes);
    setIsDirty(true);
  };

  const handleSEOChange = (field: keyof typeof seoData, value: string) => {
    setSeoData((prev) => ({
      ...prev,
      [field]: value,
    }));
    setIsDirty(true);
  };

  // Show confirmation when leaving with unsaved changes
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (isDirty) {
        e.preventDefault();
        e.returnValue = '';
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [isDirty]);

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

  if (error && !listing) {
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
    <div className="min-h-screen bg-base-200">
      <div className="container mx-auto px-4 py-8">
        <form onSubmit={handleSubmit} className="max-w-7xl mx-auto">
          {/* Header */}
          <div className="mb-8">
            <div className="flex items-center justify-between mb-4">
              <Link
                href="/profile/listings"
                className="btn btn-ghost btn-sm gap-2"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                  className="w-4 h-4"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
                  />
                </svg>
                {t('back')}
              </Link>
              {isDirty && (
                <div className="badge badge-warning gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    className="w-4 h-4 stroke-current"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                    />
                  </svg>
                  {t('unsavedChanges')}
                </div>
              )}
            </div>
            <div className="space-y-2">
              <h1 className="text-4xl font-bold">{t('title')}</h1>
              <p className="text-base-content/70 text-lg">{t('description')}</p>
              {listing && (
                <div className="flex items-center gap-2 text-sm text-base-content/60">
                  <span>ID: #{listing.id}</span>
                  <span>•</span>
                  <span>{listing.category.name}</span>
                </div>
              )}
            </div>
          </div>

          {/* Tabs */}
          <div className="card bg-base-100 shadow-sm mb-6">
            <div className="card-body p-4">
              <div className="tabs tabs-boxed bg-base-200">
                <a
                  className={`tab gap-2 ${activeTab === 'basic' ? 'tab-active' : ''}`}
                  onClick={() => setActiveTab('basic')}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
                    />
                  </svg>
                  {t('tabs.basic')}
                </a>
                <a
                  className={`tab gap-2 ${activeTab === 'images' ? 'tab-active' : ''}`}
                  onClick={() => setActiveTab('images')}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"
                    />
                  </svg>
                  {t('tabs.images')}
                  {images.length > 0 && (
                    <span className="badge badge-primary badge-sm">
                      {images.length}
                    </span>
                  )}
                </a>
                <a
                  className={`tab gap-2 ${activeTab === 'attributes' ? 'tab-active' : ''}`}
                  onClick={() => setActiveTab('attributes')}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M6 6h.008v.008H6V6z"
                    />
                  </svg>
                  {t('tabs.attributes')}
                </a>
                <a
                  className={`tab gap-2 ${activeTab === 'location' ? 'tab-active' : ''}`}
                  onClick={() => setActiveTab('location')}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"
                    />
                  </svg>
                  {t('tabs.location')}
                </a>
                <a
                  className={`tab gap-2 ${activeTab === 'seo' ? 'tab-active' : ''}`}
                  onClick={() => setActiveTab('seo')}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
                    />
                  </svg>
                  {t('tabs.seo')}
                </a>
                <a
                  className={`tab gap-2 ${activeTab === 'preview' ? 'tab-active' : ''}`}
                  onClick={() => setActiveTab('preview')}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                  </svg>
                  {t('tabs.preview')}
                </a>
              </div>
            </div>
          </div>

          {/* Error message */}
          {error && (
            <div className="alert alert-error mb-6 shadow-lg">
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
                  d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span>{error}</span>
            </div>
          )}

          {/* Content */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 space-y-6">
              {/* Tab content */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  {activeTab === 'basic' && (
                    <BasicInfoSection
                      data={formData}
                      errors={errors}
                      onChange={handleInputChange}
                    />
                  )}

                  {activeTab === 'images' && (
                    <ImagesSection
                      listingId={listing.id}
                      images={images}
                      onImagesChange={handleImagesChange}
                    />
                  )}

                  {activeTab === 'attributes' && (
                    <AttributesSection
                      categoryId={listing.category_id}
                      values={attributes || []}
                      errors={errors}
                      onChange={handleAttributesChange}
                    />
                  )}

                  {activeTab === 'location' && (
                    <LocationSection
                      data={formData}
                      errors={errors}
                      onChange={handleInputChange}
                    />
                  )}

                  {activeTab === 'seo' && (
                    <SEOSection
                      data={seoData}
                      basicData={{
                        title: formData.title,
                        description: formData.description,
                      }}
                      errors={errors}
                      onChange={handleSEOChange}
                    />
                  )}

                  {activeTab === 'preview' && (
                    <div className="space-y-4">
                      <div className="flex gap-2 mb-4">
                        <button
                          type="button"
                          className={`btn btn-sm ${previewMode === 'card' ? 'btn-primary' : 'btn-ghost'}`}
                          onClick={() => setPreviewMode('card')}
                        >
                          {t('preview.card')}
                        </button>
                        <button
                          type="button"
                          className={`btn btn-sm ${previewMode === 'list' ? 'btn-primary' : 'btn-ghost'}`}
                          onClick={() => setPreviewMode('list')}
                        >
                          {t('preview.list')}
                        </button>
                        <button
                          type="button"
                          className={`btn btn-sm ${previewMode === 'mobile' ? 'btn-primary' : 'btn-ghost'}`}
                          onClick={() => setPreviewMode('mobile')}
                        >
                          {t('preview.mobile')}
                        </button>
                      </div>
                      <PreviewCard
                        data={{
                          ...formData,
                          category: listing.category,
                          images: images,
                        }}
                        viewMode={previewMode}
                      />
                    </div>
                  )}
                </div>
              </div>
            </div>

            {/* Sidebar */}
            <div className="lg:col-span-1">
              <div className="card bg-base-100 shadow-xl sticky top-24">
                <div className="card-body">
                  <h3 className="card-title text-lg">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      strokeWidth={1.5}
                      stroke="currentColor"
                      className="w-5 h-5"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"
                      />
                    </svg>
                    {t('sidebar.status')}
                  </h3>

                  {/* Completion progress */}
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>{t('sidebar.completion')}</span>
                      <span className="font-bold text-primary">85%</span>
                    </div>
                    <progress
                      className="progress progress-primary w-full"
                      value="85"
                      max="100"
                    ></progress>
                  </div>

                  {/* Quality score */}
                  <div className="divider"></div>
                  <div className="space-y-3">
                    <h4 className="font-medium flex items-center gap-2">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        strokeWidth={1.5}
                        stroke="currentColor"
                        className="w-4 h-4"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                      {t('sidebar.quality')}
                    </h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex items-center gap-3">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 24 24"
                          strokeWidth={2}
                          stroke="currentColor"
                          className="w-5 h-5 text-success"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                        <span>{t('sidebar.hasTitle')}</span>
                      </div>
                      <div className="flex items-center gap-3">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 24 24"
                          strokeWidth={2}
                          stroke="currentColor"
                          className="w-5 h-5 text-success"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                        <span>{t('sidebar.hasDescription')}</span>
                      </div>
                      <div className="flex items-center gap-3">
                        {images.length > 0 ? (
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={2}
                            stroke="currentColor"
                            className="w-5 h-5 text-success"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                          </svg>
                        ) : (
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={2}
                            stroke="currentColor"
                            className="w-5 h-5 text-warning"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"
                            />
                          </svg>
                        )}
                        <span>{t('sidebar.hasImages')}</span>
                      </div>
                      <div className="flex items-center gap-3">
                        {formData.location ? (
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={2}
                            stroke="currentColor"
                            className="w-5 h-5 text-success"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                          </svg>
                        ) : (
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={2}
                            stroke="currentColor"
                            className="w-5 h-5 text-warning"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"
                            />
                          </svg>
                        )}
                        <span>{t('sidebar.hasLocation')}</span>
                      </div>
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="divider"></div>
                  <div className="space-y-3">
                    <button
                      type="submit"
                      className={`btn btn-primary w-full gap-2 ${saving ? 'loading' : ''}`}
                      disabled={saving || !isDirty}
                    >
                      {!saving && (
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 24 24"
                          strokeWidth={1.5}
                          stroke="currentColor"
                          className="w-5 h-5"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M17.593 3.322c1.1.128 1.907 1.077 1.907 2.185V21L12 17.25 4.5 21V5.507c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0111.186 0z"
                          />
                        </svg>
                      )}
                      {saving ? t('saving') : t('saveChanges')}
                    </button>
                    <Link
                      href="/profile/listings"
                      className="btn btn-ghost w-full gap-2"
                      onClick={(e) => {
                        if (isDirty && !confirm(t('confirmLeave'))) {
                          e.preventDefault();
                        }
                      }}
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        strokeWidth={1.5}
                        stroke="currentColor"
                        className="w-5 h-5"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          d="M6 18L18 6M6 6l12 12"
                        />
                      </svg>
                      {t('cancel')}
                    </Link>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
