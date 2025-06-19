'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTranslations, useLocale } from 'next-intl';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { fetchStorefrontBySlug } from '@/store/slices/storefrontSlice';
import { useAuth } from '@/contexts/AuthContext';
import Link from 'next/link';
import {
  ArrowLeftIcon,
  ChartBarIcon,
  UserGroupIcon,
  CalendarIcon,
  CurrencyDollarIcon,
} from '@heroicons/react/24/outline';

export default function StorefrontAnalyticsPage() {
  const t = useTranslations();
  const locale = useLocale();
  const router = useRouter();
  const params = useParams();
  const dispatch = useAppDispatch();
  const slug = params.slug as string;
  const { user } = useAuth();

  const { currentStorefront, isLoading } = useAppSelector(
    (state) => state.storefronts
  );

  const [accessDenied, setAccessDenied] = useState(false);

  useEffect(() => {
    if (slug) {
      dispatch(fetchStorefrontBySlug(slug));
    }
  }, [dispatch, slug]);

  // Проверка доступа
  useEffect(() => {
    if (!isLoading && currentStorefront) {
      // Если пользователь не авторизован
      if (!user) {
        setAccessDenied(true);
        router.push(`/${locale}/storefronts/${slug}`);
        return;
      }

      // Если пользователь не владелец
      if (currentStorefront.user_id !== user.id) {
        setAccessDenied(true);
        router.push(`/${locale}/storefronts/${slug}`);
      }
    }
  }, [currentStorefront, user, isLoading, router, slug, locale]);

  if (accessDenied) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">🔒</div>
          <h1 className="text-2xl font-bold mb-2">
            {t('common.accessDenied')}
          </h1>
          <p className="text-base-content/60">
            {t('storefronts.analyticsAccessDenied')}
          </p>
        </div>
      </div>
    );
  }

  if (isLoading || !currentStorefront) {
    return (
      <div className="min-h-screen bg-base-200">
        <div className="container mx-auto px-4 py-8">
          <div className="flex justify-center items-center min-h-[600px]">
            <div className="text-center">
              <span className="loading loading-spinner loading-lg text-primary"></span>
              <p className="mt-4 text-base-content/60">
                {t('admin.common.loading')}
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-base-100 shadow-sm border-b border-base-300">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Link
                href={`/${locale}/storefronts/${slug}`}
                className="btn btn-ghost btn-sm btn-square"
              >
                <ArrowLeftIcon className="w-5 h-5" />
              </Link>
              <div>
                <h1 className="text-2xl font-bold">{currentStorefront.name}</h1>
                <p className="text-sm text-base-content/60">
                  {t('storefronts.analyticsTitle')}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
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
          <span>{t('storefronts.analyticsComingSoon')}</span>
        </div>

        <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Placeholder cards for future analytics */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <ChartBarIcon className="w-6 h-6" />
                {t('storefronts.conversionRate')}
              </h2>
              <p className="text-base-content/60">{t('common.comingSoon')}</p>
            </div>
          </div>

          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <UserGroupIcon className="w-6 h-6" />
                {t('storefronts.customerDemographics')}
              </h2>
              <p className="text-base-content/60">{t('common.comingSoon')}</p>
            </div>
          </div>

          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <CalendarIcon className="w-6 h-6" />
                {t('storefronts.seasonalTrends')}
              </h2>
              <p className="text-base-content/60">{t('common.comingSoon')}</p>
            </div>
          </div>

          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <CurrencyDollarIcon className="w-6 h-6" />
                {t('storefronts.revenueForecasting')}
              </h2>
              <p className="text-base-content/60">{t('common.comingSoon')}</p>
            </div>
          </div>
        </div>

        <div className="mt-8 text-center">
          <Link
            href={`/${locale}/storefronts/${slug}/dashboard`}
            className="btn btn-primary"
          >
            <ArrowLeftIcon className="w-4 h-4" />
            {t('storefronts.backToDashboard')}
          </Link>
        </div>
      </div>
    </div>
  );
}
