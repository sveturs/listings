'use client';

import { use, useEffect, useState, useCallback } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import dynamic from 'next/dynamic';
import { useAuth } from '@/contexts/AuthContext';
import { storefrontApi } from '@/services/storefrontApi';
import type { Storefront } from '@/types/storefront';
import { ReviewsSection } from '@/components/reviews';
import { ImageGallery } from '@/components/reviews/ImageGallery';
import { usePageViewTracking } from '@/hooks/useAnalytics';
import StorefrontHeader from '@/components/storefronts/public/StorefrontHeader';
import StorefrontInfo from '@/components/storefronts/public/StorefrontInfo';
import StorefrontActions from '@/components/storefronts/public/StorefrontActions';
import StorefrontProducts from '@/components/storefronts/public/StorefrontProducts';

// Dynamic import for map to avoid SSR issues
const StorefrontLocationMap = dynamic(
  () => import('@/components/storefronts/public/StorefrontLocationMap'),
  { ssr: false }
);

type Props = {
  params: Promise<{ slug: string }>;
};

export default function StorefrontPage({ params }: Props) {
  const { slug } = use(params);
  const locale = useLocale();
  const t = useTranslations('storefronts');
  const { user } = useAuth();
  const [storefront, setStorefront] = useState<Storefront | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [galleryOpen, setGalleryOpen] = useState(false);
  const [galleryImages, setGalleryImages] = useState<string[]>([]);
  const [galleryIndex, setGalleryIndex] = useState(0);

  // Отслеживание просмотра публичной страницы витрины
  usePageViewTracking(storefront?.id, 'public');

  const fetchStorefront = useCallback(async () => {
    try {
      const data = await storefrontApi.getStorefrontBySlug(slug);
      setStorefront(data);
    } catch (error) {
      console.error('Error fetching storefront:', error);
    } finally {
      setIsLoading(false);
    }
  }, [slug]);

  useEffect(() => {
    fetchStorefront();
  }, [fetchStorefront]);

  const openGallery = (images: string[], index: number = 0) => {
    setGalleryImages(images);
    setGalleryIndex(index);
    setGalleryOpen(true);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-base-100">
        <div className="animate-pulse">
          {/* Header skeleton */}
          <div className="h-80 bg-base-300"></div>
          {/* Content skeleton */}
          <div className="container mx-auto px-4 py-8">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
              <div className="lg:col-span-2 space-y-6">
                <div className="h-32 bg-base-300 rounded-2xl"></div>
                <div className="h-64 bg-base-300 rounded-2xl"></div>
              </div>
              <div className="space-y-6">
                <div className="h-48 bg-base-300 rounded-2xl"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!storefront) {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">🏪</div>
          <h1 className="text-2xl font-bold mb-2">{t('notFound')}</h1>
          <p className="text-base-content/60 mb-6">
            {t('notFoundDescription')}
          </p>
          <Link href={`/${locale}/storefronts`} className="btn btn-primary">
            {t('browseStorefronts')}
          </Link>
        </div>
      </div>
    );
  }

  const isOwner = user?.id === storefront.user_id;

  return (
    <div className="min-h-screen bg-base-100">
      {/* Hero Header */}
      <StorefrontHeader
        storefront={storefront}
        isOwner={isOwner}
        onImageClick={openGallery}
      />

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Left Column - Main Info */}
          <div className="lg:col-span-2 space-y-8">
            {/* About Section */}
            <div className="card bg-base-200 shadow-xl">
              <div className="card-body">
                <h2 className="card-title text-2xl mb-4">{t('about')}</h2>
                <p className="text-base-content/80 whitespace-pre-wrap">
                  {storefront.description}
                </p>

                {/* Business Details */}
                {(storefront.registration_number ||
                  storefront.tax_number ||
                  storefront.vat_number) && (
                  <div className="mt-6 pt-6 border-t border-base-300">
                    <h3 className="font-semibold mb-3">
                      {t('businessDetails')}
                    </h3>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm">
                      {storefront.registration_number && (
                        <div>
                          <span className="text-base-content/60">
                            {t('registrationNumber')}:
                          </span>
                          <span className="ml-2 font-medium">
                            {storefront.registration_number}
                          </span>
                        </div>
                      )}
                      {storefront.tax_number && (
                        <div>
                          <span className="text-base-content/60">
                            {t('taxNumber')}:
                          </span>
                          <span className="ml-2 font-medium">
                            {storefront.tax_number}
                          </span>
                        </div>
                      )}
                      {storefront.vat_number && (
                        <div>
                          <span className="text-base-content/60">
                            {t('vatNumber')}:
                          </span>
                          <span className="ml-2 font-medium">
                            {storefront.vat_number}
                          </span>
                        </div>
                      )}
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Products Section */}
            <StorefrontProducts storefrontSlug={slug} />

            {/* Location Map */}
            {storefront.location &&
              storefront.location.user_lat &&
              storefront.location.user_lng &&
              storefront.location.full_address && (
                <div className="card bg-base-200 shadow-xl overflow-hidden">
                  <div className="card-body">
                    <h2 className="card-title text-2xl mb-4">
                      {t('location')}
                    </h2>
                    <div className="h-96 -mx-8 -mb-8">
                      <StorefrontLocationMap
                        location={{
                          user_lat: storefront.location.user_lat,
                          user_lng: storefront.location.user_lng,
                          full_address: storefront.location.full_address,
                        }}
                        storefront={storefront}
                      />
                    </div>
                  </div>
                </div>
              )}

            {/* Reviews Section */}
            <div className="card bg-base-200 shadow-xl">
              <div className="card-body">
                <ReviewsSection
                  entityType="storefront"
                  entityId={storefront.id || 0}
                  sellerId={storefront.user_id || 0}
                />
              </div>
            </div>
          </div>

          {/* Right Column - Info & Actions */}
          <div className="space-y-6">
            {/* Store Info Card */}
            <StorefrontInfo storefront={storefront} />

            {/* Actions Card */}
            <StorefrontActions storefront={storefront} isOwner={isOwner} />

            {/* Trust & Safety */}
            <div className="alert alert-info">
              <svg
                className="w-5 h-5 flex-shrink-0"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div>
                <h4 className="font-semibold">{t('trustSafety')}</h4>
                <ul className="text-sm space-y-1 mt-1">
                  {storefront.is_verified && <li>• {t('verifiedBusiness')}</li>}
                  <li>• {t('securePayments')}</li>
                  <li>• {t('buyerProtection')}</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Image Gallery Modal */}
      <ImageGallery
        images={galleryImages}
        initialIndex={galleryIndex}
        isOpen={galleryOpen}
        onClose={() => setGalleryOpen(false)}
      />
    </div>
  );
}
