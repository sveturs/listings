'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import { useRouter } from '@/i18n/routing';
import Link from 'next/link';
import { getTranslatedAttribute } from '@/utils/translatedAttribute';
import { ReviewsSection } from '@/components/reviews';
import Image from 'next/image';
import ContactSellerButton from '@/components/ContactSellerButton';

interface Listing {
  id: number;
  title: string;
  description: string;
  price: number;
  condition: string;
  status: string;
  location: string;
  show_on_map: boolean;
  latitude?: number;
  longitude?: number;
  views_count: number;
  created_at: string;
  updated_at: string;
  user: {
    id: number;
    name: string;
    email: string;
    picture_url?: string;
  };
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
  attributes?: Array<{
    attribute_id: number;
    attribute_name: string;
    display_name: string;
    attribute_type: string;
    value?: string | number | boolean;
    text_value?: string;
    numeric_value?: number;
    boolean_value?: boolean;
    display_value: string;
    unit?: string;
    is_required?: boolean;
    show_in_card?: boolean;
    translations?: any;
    option_translations?: any;
  }>;
}

export default function ListingDetailBySlugPage() {
  const params = useParams();
  const slug = params?.slug as string;
  const locale = params?.locale as string;
  const t = useTranslations('marketplace');
  const _router = useRouter();

  const [listing, setListing] = useState<Listing | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);

  useEffect(() => {
    const fetchListing = async () => {
      if (!slug) return;

      try {
        setLoading(true);
        setError(null);

        const response = await apiClient.get<{ data: Listing }>(
          `/api/v1/c2c/listings/slug/${slug}`
        );

        if (!response.error && response.data) {
          setListing(response.data.data || response.data);
        } else {
          setError(t('listing.notFound'));
        }
      } catch (err) {
        console.error('Error fetching listing:', err);
        setError(t('listing.loadError'));
      } finally {
        setLoading(false);
      }
    };

    fetchListing();
  }, [slug, t]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error || !listing) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">{t('listing.notFound')}</h1>
          <Link href="/" className="btn btn-primary">
            {t('listing.backToHome')}
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Images */}
            {listing.images.length > 0 && (
              <div className="card bg-base-100 shadow-xl">
                <figure className="relative">
                  <div className="relative w-full h-96">
                    <Image
                      src={`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000'}${listing.images[selectedImageIndex].public_url}`}
                      alt={listing.title}
                      fill
                      className="object-cover cursor-pointer"
                      onClick={() => setSelectedImageIndex(selectedImageIndex)}
                    />
                  </div>
                  {listing.images.length > 1 && (
                    <div className="absolute bottom-4 left-4 right-4 flex gap-2 overflow-x-auto">
                      {listing.images.map((image, index) => (
                        <div
                          key={image.id}
                          className={`relative w-20 h-20 rounded cursor-pointer ${
                            index === selectedImageIndex
                              ? 'ring-2 ring-primary'
                              : 'opacity-70 hover:opacity-100'
                          }`}
                          onClick={() => setSelectedImageIndex(index)}
                        >
                          <Image
                            src={`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000'}${image.public_url}`}
                            alt={`${listing.title} ${index + 1}`}
                            fill
                            className="object-cover rounded"
                          />
                        </div>
                      ))}
                    </div>
                  )}
                </figure>
              </div>
            )}

            {/* Title and Description */}
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h1 className="card-title text-3xl">{listing.title}</h1>
                <div className="flex items-center gap-4 text-sm text-base-content/70">
                  <span>
                    {t('listing.views', { count: listing.views_count })}
                  </span>
                  <span>•</span>
                  <span>
                    {new Date(listing.created_at).toLocaleDateString()}
                  </span>
                </div>
                <div className="divider"></div>
                <p className="whitespace-pre-wrap">{listing.description}</p>
              </div>
            </div>

            {/* Attributes */}
            {listing.attributes &&
              listing.attributes.length > 0 &&
              (() => {
                // Фильтруем атрибуты для отображения
                const visibleAttributes = listing.attributes.filter(
                  (attr) => attr.show_in_card !== false
                );

                // Разделяем на обязательные и дополнительные
                const requiredAttributes = visibleAttributes.filter(
                  (attr) => attr.is_required
                );
                const optionalAttributes = visibleAttributes.filter(
                  (attr) => !attr.is_required
                );

                if (visibleAttributes.length === 0) return null;

                return (
                  <div className="card bg-base-100 shadow-xl">
                    <div className="card-body">
                      {/* Основные характеристики */}
                      {requiredAttributes.length > 0 && (
                        <>
                          <h2 className="card-title">
                            {locale === 'ru'
                              ? 'Основные характеристики'
                              : 'Main Specifications'}
                          </h2>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                            {requiredAttributes.map((attr) => {
                              const translated = getTranslatedAttribute(
                                {
                                  id: attr.attribute_id,
                                  name: attr.attribute_name,
                                  display_name: attr.display_name,
                                  attribute_type: attr.attribute_type,
                                  translations: attr.translations,
                                  option_translations: attr.option_translations,
                                },
                                locale
                              );

                              return (
                                <div
                                  key={attr.attribute_id}
                                  className="flex justify-between py-2 border-b border-base-300"
                                >
                                  <span className="text-base-content/70 font-medium">
                                    {translated.displayName}:
                                  </span>
                                  <span className="font-semibold">
                                    {attr.display_value}
                                    {attr.unit && ` ${attr.unit}`}
                                  </span>
                                </div>
                              );
                            })}
                          </div>
                        </>
                      )}

                      {/* Дополнительные характеристики */}
                      {optionalAttributes.length > 0 && (
                        <>
                          {requiredAttributes.length > 0 && (
                            <div className="divider"></div>
                          )}
                          <h3 className="text-lg font-semibold">
                            {locale === 'ru'
                              ? 'Дополнительные характеристики'
                              : 'Additional Specifications'}
                          </h3>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                            {optionalAttributes.map((attr) => {
                              const translated = getTranslatedAttribute(
                                {
                                  id: attr.attribute_id,
                                  name: attr.attribute_name,
                                  display_name: attr.display_name,
                                  attribute_type: attr.attribute_type,
                                  translations: attr.translations,
                                  option_translations: attr.option_translations,
                                },
                                locale
                              );

                              return (
                                <div
                                  key={attr.attribute_id}
                                  className="flex justify-between py-2 border-b border-base-300"
                                >
                                  <span className="text-base-content/60">
                                    {translated.displayName}:
                                  </span>
                                  <span className="font-medium">
                                    {attr.display_value}
                                    {attr.unit && ` ${attr.unit}`}
                                  </span>
                                </div>
                              );
                            })}
                          </div>
                        </>
                      )}
                    </div>
                  </div>
                );
              })()}

            {/* Reviews Section */}
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <ReviewsSection
                  entityType="listing"
                  entityId={listing.id}
                  sellerId={listing.user.id}
                />
              </div>
            </div>
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-1 space-y-6">
            {/* Price and Actions */}
            <div className="card bg-base-100 shadow-xl sticky top-24">
              <div className="card-body">
                <div className="text-3xl font-bold text-primary">
                  {listing.price.toLocaleString()} RSD
                </div>
                <div className="divider"></div>

                {/* Seller Info */}
                <div className="flex items-center gap-3 mb-4">
                  {listing.user.picture_url ? (
                    <div className="relative w-12 h-12">
                      <Image
                        src={listing.user.picture_url}
                        alt={listing.user.name}
                        fill
                        className="rounded-full object-cover"
                      />
                    </div>
                  ) : (
                    <div className="w-12 h-12 rounded-full bg-primary flex items-center justify-center text-white font-bold">
                      {listing.user.name.charAt(0).toUpperCase()}
                    </div>
                  )}
                  <div>
                    <p className="font-semibold">{listing.user.name}</p>
                    <p className="text-sm text-base-content/70">
                      {t('listing.memberSince', {
                        date: new Date(listing.created_at).getFullYear(),
                      })}
                    </p>
                  </div>
                </div>

                <ContactSellerButton
                  sellerId={listing.user.id}
                  listingId={listing.id}
                  className="w-full"
                  size="lg"
                />

                {/* Location */}
                {listing.location && (
                  <div className="mt-4">
                    <h3 className="font-semibold mb-2 flex items-center gap-2">
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
                          d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"
                        />
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"
                        />
                      </svg>
                      {t('listing.location')}
                    </h3>
                    <p className="text-base-content/70">{listing.location}</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
