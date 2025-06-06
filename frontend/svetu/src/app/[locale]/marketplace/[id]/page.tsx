'use client';

import { use, useEffect, useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import { useAuth } from '@/contexts/AuthContext';
import config from '@/config';
import Image from 'next/image';
import Link from 'next/link';

interface Listing {
  id: number;
  title: string;
  description: string;
  price: number;
  user_id: number;
  images?: Array<{
    id: number;
    public_url: string;
  }>;
  location?: string;
  created_at: string;
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
      <div className="flex items-center justify-center h-screen">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (!listing) {
    return (
      <div className="container mx-auto p-4">
        <h1 className="text-2xl font-bold">
          {locale === 'ru' ? 'Объявление не найдено' : 'Listing not found'}
        </h1>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <div className="max-w-4xl mx-auto">
        {/* Images */}
        {listing.images && listing.images.length > 0 && (
          <div className="mb-6">
            <Image
              src={config.buildImageUrl(listing.images[0].public_url)}
              alt={listing.title}
              width={800}
              height={600}
              className="w-full h-96 object-cover rounded-lg"
            />
          </div>
        )}

        {/* Title and Price */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold mb-2">{listing.title}</h1>
          <p className="text-2xl font-semibold text-primary">
            {listing.price} $
          </p>
        </div>

        {/* Description */}
        <div className="mb-6">
          <h2 className="text-xl font-semibold mb-2">
            {locale === 'ru' ? 'Описание' : 'Description'}
          </h2>
          <p className="text-base-content/70">{listing.description}</p>
        </div>

        {/* Location */}
        {listing.location && (
          <div className="mb-6">
            <p className="flex items-center gap-2 text-base-content/70">
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
          </div>
        )}

        {/* Actions */}
        <div className="flex gap-4">
          {user && user.id !== listing.user_id && (
            <button onClick={handleChatClick} className="btn btn-primary">
              <svg
                className="w-5 h-5 mr-2"
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
              {locale === 'ru' ? 'Написать сообщение' : 'Send Message'}
            </button>
          )}

          <Link href={`/${locale}`} className="btn btn-ghost">
            {locale === 'ru' ? 'Назад к объявлениям' : 'Back to listings'}
          </Link>
        </div>
      </div>
    </div>
  );
}
