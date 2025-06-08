'use client';

import { useState } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { useLocale } from 'next-intl';
import { format } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

interface SellerInfoProps {
  listing: {
    id: number;
    user_id: number;
    user?: {
      id: number;
      name: string;
      email: string;
      picture_url?: string;
      created_at: string;
    };
    seller_response_rate?: number;
    seller_response_time?: string;
    seller_total_listings?: number;
  };
  onChatClick: () => void;
}

export default function SellerInfo({ listing, onChatClick }: SellerInfoProps) {
  const locale = useLocale();
  const { user } = useAuth();
  const router = useRouter();
  const [showPhone, setShowPhone] = useState(false);
  const [rating] = useState(4.5); // TODO: Get from API
  const [reviewsCount] = useState(23); // TODO: Get from API

  const dateLocale = locale === 'ru' ? ru : enUS;
  const formatDate = (date: string) => {
    return format(new Date(date), 'dd MMMM yyyy', { locale: dateLocale });
  };

  const renderStars = (rating: number) => {
    return (
      <div className="flex items-center gap-1">
        {[1, 2, 3, 4, 5].map((star) => (
          <svg
            key={star}
            className={`w-4 h-4 ${
              star <= Math.floor(rating)
                ? 'text-warning fill-warning'
                : star <= Math.ceil(rating)
                  ? 'text-warning fill-warning opacity-50'
                  : 'text-base-300 fill-base-300'
            }`}
            viewBox="0 0 20 20"
          >
            <path d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z" />
          </svg>
        ))}
        <span className="text-sm ml-1">({reviewsCount})</span>
      </div>
    );
  };

  return (
    <div className="card bg-base-200">
      <div className="card-body">
        <h3 className="font-semibold mb-4">
          {locale === 'ru' ? 'Продавец' : 'Seller'}
        </h3>

        {/* Seller Avatar and Basic Info */}
        <div className="flex items-start gap-4 mb-4">
          <div className="avatar">
            <div className="w-20 h-20 rounded-full bg-base-300">
              {listing.user?.picture_url ? (
                <Image
                  src={listing.user.picture_url}
                  alt={listing.user.name}
                  width={80}
                  height={80}
                  className="rounded-full"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-3xl font-bold">
                  {listing.user?.name.charAt(0).toUpperCase()}
                </div>
              )}
            </div>
          </div>
          <div className="flex-1">
            <p className="font-medium text-lg">{listing.user?.name}</p>
            {renderStars(rating)}
            <p className="text-sm text-base-content/60 mt-1">
              {locale === 'ru' ? 'На сайте с' : 'Member since'}{' '}
              {listing.user && formatDate(listing.user.created_at)}
            </p>
          </div>
        </div>

        {/* Seller Stats */}
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="text-center p-3 bg-base-100 rounded-lg">
            <div className="text-2xl font-bold text-primary">
              {listing.seller_response_rate || 95}%
            </div>
            <div className="text-xs text-base-content/60">
              {locale === 'ru' ? 'Процент ответов' : 'Response rate'}
            </div>
          </div>
          <div className="text-center p-3 bg-base-100 rounded-lg">
            <div className="text-2xl font-bold text-primary">
              {listing.seller_response_time || '< 1ч'}
            </div>
            <div className="text-xs text-base-content/60">
              {locale === 'ru' ? 'Время ответа' : 'Response time'}
            </div>
          </div>
        </div>

        {/* Seller Badges */}
        <div className="flex flex-wrap gap-2 mb-6">
          <div className="badge badge-success gap-1">
            <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                clipRule="evenodd"
              />
            </svg>
            {locale === 'ru' ? 'Проверен' : 'Verified'}
          </div>
          {(listing.seller_total_listings || 0) > 10 && (
            <div className="badge badge-info gap-1">
              <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" />
                <path
                  fillRule="evenodd"
                  d="M4 5a2 2 0 012-2 1 1 0 000 2H6a2 2 0 100 4h2a2 2 0 100 4h-2a1 1 0 100 2 2 2 0 01-2-2V5z"
                  clipRule="evenodd"
                />
              </svg>
              {locale === 'ru' ? 'Опытный продавец' : 'Experienced'}
            </div>
          )}
        </div>

        {/* Action Buttons */}
        <div className="space-y-3">
          {user && user.id !== listing.user_id ? (
            <>
              <button onClick={onChatClick} className="btn btn-primary w-full">
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
                {locale === 'ru' ? 'Написать сообщение' : 'Send Message'}
              </button>

              <button
                onClick={() => setShowPhone(!showPhone)}
                className="btn btn-outline w-full"
              >
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
                    d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"
                  />
                </svg>
                {showPhone
                  ? '+381 69 123 4567'
                  : locale === 'ru'
                    ? 'Показать телефон'
                    : 'Show Phone'}
              </button>

              <Link
                href={`/${locale}/marketplace?user_id=${listing.user_id}`}
                className="btn btn-ghost w-full"
              >
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
                    d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                  />
                </svg>
                {locale === 'ru' ? 'Все товары продавца' : 'All seller items'}(
                {listing.seller_total_listings || 0})
              </Link>
            </>
          ) : user && user.id === listing.user_id ? (
            <div className="text-center">
              <p className="text-base-content/60 mb-3">
                {locale === 'ru'
                  ? 'Это ваше объявление'
                  : 'This is your listing'}
              </p>
              <Link
                href={`/${locale}/profile/listings/${listing.id}/edit`}
                className="btn btn-outline btn-sm"
              >
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
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                  />
                </svg>
                {locale === 'ru' ? 'Редактировать' : 'Edit'}
              </Link>
            </div>
          ) : (
            <div className="text-center">
              <p className="text-base-content/60 mb-3">
                {locale === 'ru'
                  ? 'Войдите, чтобы связаться с продавцом'
                  : 'Sign in to contact seller'}
              </p>
              <button
                onClick={() => router.push('/')}
                className="btn btn-primary btn-sm"
              >
                {locale === 'ru' ? 'Войти' : 'Sign In'}
              </button>
            </div>
          )}
        </div>

        {/* Trust & Safety */}
        <div className="divider"></div>
        <div className="text-xs text-base-content/60 space-y-1">
          <div className="flex items-start gap-2">
            <svg
              className="w-4 h-4 mt-0.5 text-success"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
              />
            </svg>
            <span>
              {locale === 'ru'
                ? 'Все сделки защищены правилами платформы'
                : 'All transactions are protected by platform rules'}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
