'use client';

import { use, useEffect, useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import { useAuth } from '@/contexts/AuthContext';
import config from '@/config';
import Link from 'next/link';
import { formatDistanceToNow } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import ImageGallery from '@/components/marketplace/listing/ImageGallery';
import SellerInfo from '@/components/marketplace/listing/SellerInfo';
import ListingActions from '@/components/marketplace/listing/ListingActions';
import SimilarListings from '@/components/marketplace/listing/SimilarListings';
import { getTranslatedAttribute } from '@/utils/translatedAttribute';
import { ReviewsSection } from '@/components/reviews';

interface User {
  id: number;
  name: string;
  email: string;
  picture_url?: string;
  created_at: string;
}

interface Category {
  id: number;
  name: string;
  slug: string;
}

interface Attribute {
  listing_id: number;
  attribute_id: number;
  attribute_name: string;
  display_name: string;
  attribute_type: string;
  text_value?: string;
  numeric_value?: number;
  display_value: string;
  translations?: Record<string, string>;
  option_translations?: Record<string, Record<string, string>>;
  is_required: boolean;
  show_in_card: boolean;
  show_in_list: boolean;
}

interface ListingImage {
  id: number;
  public_url: string;
  is_video?: boolean;
}

interface Listing {
  id: number;
  user_id: number;
  category_id: number;
  title: string;
  description: string;
  price: number;
  old_price?: number;
  condition?: string;
  status: string;
  location?: string;
  latitude?: number;
  longitude?: number;
  city?: string;
  country?: string;
  views_count: number;
  favorites_count?: number;
  created_at: string;
  updated_at: string;
  user?: User;
  category?: Category;
  category_path_names?: string[];
  category_path_ids?: number[];
  category_path_slugs?: string[];
  attributes?: Attribute[];
  images?: ListingImage[];
  is_favorite?: boolean;
  show_on_map?: boolean;
  has_discount?: boolean;
  discount_percentage?: number;
  seller_response_rate?: number;
  seller_response_time?: string;
  seller_total_listings?: number;
  helpful_votes?: number;
  not_helpful_votes?: number;
}

type Props = {
  params: Promise<{ id: string }>;
};

export default function ListingPage({ params }: Props) {
  const { id } = use(params);
  const locale = useLocale();
  const router = useRouter();
  const { user } = useAuth();
  const [listing, setListing] = useState<Listing | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [showFullDescription, setShowFullDescription] = useState(false);

  const fetchListing = useCallback(async () => {
    try {
      const response = await fetch(
        `${config.getApiUrl()}/api/v1/marketplace/listings/${id}`
      );
      if (!response.ok) throw new Error('Failed to fetch listing');
      const result = await response.json();
      // Проверяем обертку ответа
      if (result.data) {
        setListing(result.data);
      } else {
        setListing(result);
      }
    } catch (error) {
      console.error('Error fetching listing:', error);
    } finally {
      setIsLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchListing();
  }, [fetchListing]);

  const handleChatClick = () => {
    if (!user) {
      // Redirect to login if not authenticated
      router.push('/');
      return;
    }

    if (listing && user.id !== listing.user_id) {
      // Navigate to chat with listing_id and seller_id as params
      router.push(
        `/${locale}/chat?listing_id=${listing.id}&seller_id=${listing.user_id}`
      );
    }
  };

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="animate-pulse">
          <div className="h-4 bg-base-300 rounded w-1/4 mb-8"></div>
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-6">
              <div className="aspect-[4/3] bg-base-300 rounded-2xl"></div>
              <div className="space-y-4">
                <div className="h-8 bg-base-300 rounded w-3/4"></div>
                <div className="h-4 bg-base-300 rounded w-full"></div>
                <div className="h-4 bg-base-300 rounded w-5/6"></div>
              </div>
            </div>
            <div className="space-y-6">
              <div className="h-32 bg-base-300 rounded-2xl"></div>
              <div className="h-64 bg-base-300 rounded-2xl"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!listing) {
    return (
      <div className="container mx-auto p-4">
        <div className="text-center py-16">
          <svg
            className="w-24 h-24 mx-auto mb-4 text-base-content/20"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <h1 className="text-2xl font-bold mb-2">
            {locale === 'ru' ? 'Объявление не найдено' : 'Listing not found'}
          </h1>
          <p className="text-base-content/60 mb-4">
            {locale === 'ru'
              ? 'Возможно, объявление было удалено или перемещено'
              : 'The listing may have been deleted or moved'}
          </p>
          <Link href={`/${locale}`} className="btn btn-primary">
            {locale === 'ru' ? 'Вернуться к объявлениям' : 'Back to listings'}
          </Link>
        </div>
      </div>
    );
  }

  const dateLocale = locale === 'ru' ? ru : enUS;

  const formatRelativeDate = (date: string) => {
    return formatDistanceToNow(new Date(date), {
      addSuffix: true,
      locale: dateLocale,
    });
  };

  // Generate placeholder images if none exist
  const images =
    listing.images && listing.images.length > 0
      ? listing.images
      : [{ id: 0, public_url: '/placeholder-listing.jpg' }];

  // Check if description is long
  const isLongDescription = listing.description.length > 300;
  const displayDescription =
    showFullDescription || !isLongDescription
      ? listing.description
      : listing.description.slice(0, 300) + '...';

  return (
    <div className="min-h-screen bg-base-100">
      {/* Breadcrumbs */}
      <div className="container mx-auto px-4 py-4">
        <div className="text-sm breadcrumbs">
          <ul>
            <li>
              <Link href={`/${locale}`}>
                {locale === 'ru' ? 'Главная' : 'Home'}
              </Link>
            </li>
            {listing.category_path_names?.map((name, index) => (
              <li key={index}>
                <Link
                  href={`/${locale}/marketplace?category=${listing.category_path_ids?.[index]}`}
                >
                  {name}
                </Link>
              </li>
            ))}
            <li className="font-semibold">{listing.title}</li>
          </ul>
        </div>
      </div>

      <div className="container mx-auto px-4 pb-12">
        {/* Title and Actions Row */}
        <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-4 mb-6">
          <div className="flex-1">
            <h1 className="text-3xl font-bold mb-2">{listing.title}</h1>
            <div className="flex flex-wrap items-center gap-3 text-sm text-base-content/60">
              <span className="flex items-center gap-1">
                <svg
                  className="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                {formatRelativeDate(listing.created_at)}
              </span>
              <span className="flex items-center gap-1">
                <svg
                  className="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                  />
                </svg>
                {listing.views_count} {locale === 'ru' ? 'просмотров' : 'views'}
              </span>
              {listing.favorites_count !== undefined &&
                listing.favorites_count > 0 && (
                  <span className="flex items-center gap-1">
                    <svg
                      className="w-4 h-4"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                      />
                    </svg>
                    {listing.favorites_count}{' '}
                    {locale === 'ru' ? 'в избранном' : 'favorites'}
                  </span>
                )}
              <span>ID: {listing.id}</span>
            </div>
          </div>
          <ListingActions listing={listing} />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Left Column - Images and Description */}
          <div className="lg:col-span-2 space-y-6">
            {/* Image Gallery */}
            <ImageGallery images={images} title={listing.title} />

            {/* Price and Main Info */}
            <div className="card bg-base-200">
              <div className="card-body">
                <div className="flex flex-wrap items-baseline gap-4">
                  <div className="flex items-baseline gap-2">
                    <h2 className="text-3xl font-bold text-primary">
                      {listing.price} $
                    </h2>
                    {listing.old_price && listing.old_price > listing.price && (
                      <>
                        <span className="text-lg line-through text-base-content/50">
                          {listing.old_price} $
                        </span>
                        <span className="badge badge-error">
                          -
                          {Math.round(
                            ((listing.old_price - listing.price) /
                              listing.old_price) *
                              100
                          )}
                          %
                        </span>
                      </>
                    )}
                  </div>
                  {listing.condition && (
                    <span className="badge badge-lg">
                      {listing.condition === 'new'
                        ? locale === 'ru'
                          ? 'Новое'
                          : 'New'
                        : locale === 'ru'
                          ? 'Б/У'
                          : 'Used'}
                    </span>
                  )}
                </div>

                {listing.location && (
                  <p className="flex items-center gap-2 mt-4 text-base-content/70">
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
                        d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                      />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                      />
                    </svg>
                    {listing.location}
                  </p>
                )}
              </div>
            </div>

            {/* Description Card */}
            <div className="card bg-base-200">
              <div className="card-body">
                <h2 className="card-title text-2xl mb-4">
                  {locale === 'ru' ? 'Описание' : 'Description'}
                </h2>
                <p className="text-base-content/80 whitespace-pre-wrap">
                  {displayDescription}
                </p>
                {isLongDescription && (
                  <button
                    onClick={() => setShowFullDescription(!showFullDescription)}
                    className="btn btn-link btn-sm p-0 h-auto min-h-0 mt-2"
                  >
                    {showFullDescription
                      ? locale === 'ru'
                        ? 'Скрыть'
                        : 'Show less'
                      : locale === 'ru'
                        ? 'Показать полностью'
                        : 'Show more'}
                  </button>
                )}

                {/* Attributes */}
                {listing.attributes &&
                  listing.attributes.length > 0 &&
                  (() => {
                    // Фильтруем атрибуты для отображения в карточке
                    const cardAttributes = listing.attributes.filter(
                      (attr) => attr.show_in_card
                    );

                    // Разделяем на обязательные и дополнительные
                    const requiredAttributes = cardAttributes.filter(
                      (attr) => attr.is_required
                    );
                    const optionalAttributes = cardAttributes.filter(
                      (attr) => !attr.is_required
                    );

                    if (cardAttributes.length === 0) return null;

                    return (
                      <>
                        {/* Обязательные характеристики */}
                        {requiredAttributes.length > 0 && (
                          <div className="mt-6">
                            <h3 className="text-lg font-semibold mb-3">
                              {locale === 'ru'
                                ? 'Основные характеристики'
                                : 'Main Specifications'}
                            </h3>
                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                              {requiredAttributes.map((attr) => {
                                const translated = getTranslatedAttribute(
                                  {
                                    id: attr.attribute_id,
                                    name: attr.attribute_name,
                                    display_name: attr.display_name,
                                    attribute_type: attr.attribute_type,
                                    translations: attr.translations,
                                    option_translations:
                                      attr.option_translations,
                                  },
                                  locale
                                );

                                const displayValue =
                                  attr.attribute_type === 'select' &&
                                  attr.text_value
                                    ? translated.getOptionLabel(attr.text_value)
                                    : attr.display_value;

                                return (
                                  <div
                                    key={attr.attribute_id}
                                    className="flex justify-between py-2 border-b border-base-300"
                                  >
                                    <span className="text-base-content/70 font-medium">
                                      {translated.displayName}:
                                    </span>
                                    <span className="font-semibold">
                                      {displayValue}
                                    </span>
                                  </div>
                                );
                              })}
                            </div>
                          </div>
                        )}

                        {/* Дополнительные характеристики */}
                        {optionalAttributes.length > 0 && (
                          <div className="mt-6">
                            <h3 className="text-lg font-semibold mb-3">
                              {locale === 'ru'
                                ? 'Дополнительные характеристики'
                                : 'Additional Specifications'}
                            </h3>
                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                              {optionalAttributes.map((attr) => {
                                const translated = getTranslatedAttribute(
                                  {
                                    id: attr.attribute_id,
                                    name: attr.attribute_name,
                                    display_name: attr.display_name,
                                    attribute_type: attr.attribute_type,
                                    translations: attr.translations,
                                    option_translations:
                                      attr.option_translations,
                                  },
                                  locale
                                );

                                const displayValue =
                                  attr.attribute_type === 'select' &&
                                  attr.text_value
                                    ? translated.getOptionLabel(attr.text_value)
                                    : attr.display_value;

                                return (
                                  <div
                                    key={attr.attribute_id}
                                    className="flex justify-between py-2 border-b border-base-300"
                                  >
                                    <span className="text-base-content/60">
                                      {translated.displayName}:
                                    </span>
                                    <span className="font-medium">
                                      {displayValue}
                                    </span>
                                  </div>
                                );
                              })}
                            </div>
                          </div>
                        )}
                      </>
                    );
                  })()}
              </div>
            </div>

            {/* Reviews Section */}
            <div className="card bg-base-200">
              <div className="card-body">
                <ReviewsSection
                  entityType="listing"
                  entityId={listing.id}
                  sellerId={listing.user_id}
                />
              </div>
            </div>

            {/* Location Map (placeholder for now) */}
            {listing.show_on_map && listing.latitude && listing.longitude && (
              <div className="card bg-base-200">
                <div className="card-body">
                  <h2 className="card-title mb-4">
                    {locale === 'ru' ? 'Местоположение' : 'Location'}
                  </h2>
                  <div className="h-64 bg-base-300 rounded-lg flex items-center justify-center">
                    <div className="text-center">
                      <svg
                        className="w-12 h-12 mx-auto mb-2 text-base-content/40"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                        />
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                        />
                      </svg>
                      <p className="text-base-content/60">
                        {listing.city}, {listing.country}
                      </p>
                    </div>
                  </div>
                  {/* TODO: Add real map component */}
                </div>
              </div>
            )}
          </div>

          {/* Right Column - Seller Info and Actions */}
          <div className="space-y-6 lg:sticky lg:top-4">
            {/* Seller Info */}
            <SellerInfo listing={listing} onChatClick={handleChatClick} />

            {/* Safety Tips */}
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
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div>
                <h4 className="font-semibold">
                  {locale === 'ru' ? 'Советы по безопасности' : 'Safety Tips'}
                </h4>
                <ul className="text-sm space-y-1 mt-1">
                  <li>
                    •{' '}
                    {locale === 'ru'
                      ? 'Встречайтесь в публичных местах'
                      : 'Meet in public places'}
                  </li>
                  <li>
                    •{' '}
                    {locale === 'ru'
                      ? 'Проверяйте товар перед покупкой'
                      : 'Inspect items before purchase'}
                  </li>
                  <li>
                    •{' '}
                    {locale === 'ru'
                      ? 'Используйте безопасные способы оплаты'
                      : 'Use secure payment methods'}
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        {/* Similar Listings Section */}
        <SimilarListings listingId={listing.id} />
      </div>

      {/* Mobile Sticky Bottom Bar */}
      <div className="fixed bottom-0 left-0 right-0 bg-base-100 border-t border-base-300 p-4 lg:hidden">
        <div className="flex items-center justify-between gap-4">
          <div>
            <div className="text-2xl font-bold text-primary">
              {listing.price} $
            </div>
            <div className="text-sm text-base-content/60">
              {listing.condition === 'new'
                ? locale === 'ru'
                  ? 'Новое'
                  : 'New'
                : locale === 'ru'
                  ? 'Б/У'
                  : 'Used'}
            </div>
          </div>
          {user && user.id !== listing.user_id && (
            <button onClick={handleChatClick} className="btn btn-primary">
              <svg
                className="w-5 h-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
              </svg>
              {locale === 'ru' ? 'Написать' : 'Message'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
