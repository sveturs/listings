'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from '@/i18n/routing';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { fetchMyStorefronts } from '@/store/slices/storefrontSlice';
import { Link } from '@/i18n/routing';
import { Storefront } from '@/types/storefront';
import {
  MapPinIcon,
  ClockIcon,
  PlusIcon,
  ChartBarIcon,
  PencilIcon,
  ShoppingBagIcon,
  UserGroupIcon,
  CogIcon,
  EyeIcon,
  ChatBubbleLeftRightIcon,
  StarIcon,
  SparklesIcon,
  ArrowTrendingUpIcon,
  CurrencyDollarIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon as Clock8Icon,
} from '@heroicons/react/24/outline';
import { StarIcon as StarSolidIcon } from '@heroicons/react/24/solid';

export default function MyStorefrontsPage() {
  const t = useTranslations();
  const router = useRouter();
  const dispatch = useAppDispatch();
  const [selectedTab, setSelectedTab] = useState<'all' | 'active' | 'inactive'>(
    'all'
  );

  const { myStorefronts: userStorefronts, isLoading } = useAppSelector(
    (state) => state.storefronts
  );

  useEffect(() => {
    dispatch(fetchMyStorefronts());
  }, [dispatch]);

  const handleCreateStorefront = () => {
    router.push('/create-storefront');
  };

  const getStatusIcon = (isActive: boolean, isVerified: boolean) => {
    if (isActive && isVerified) {
      return <CheckCircleIcon className="w-5 h-5 text-success" />;
    } else if (isActive && !isVerified) {
      return <Clock8Icon className="w-5 h-5 text-info" />;
    } else {
      return <XCircleIcon className="w-5 h-5 text-base-content/50" />;
    }
  };

  const formatBusinessHours = (settings: any) => {
    // Business hours might be stored in settings JSON
    const hours = settings?.business_hours;
    if (!hours || typeof hours !== 'object') return t('storefronts.alwaysOpen');

    const today = new Date()
      .toLocaleDateString('en-US', { weekday: 'long' })
      .toLowerCase();
    const todayHours = hours[today];

    if (!todayHours || !todayHours.open) {
      return <span className="text-error">{t('storefronts.closedToday')}</span>;
    }

    return (
      <span className="text-success">
        {todayHours.from} - {todayHours.to}
      </span>
    );
  };

  const filteredStorefronts = userStorefronts.filter(
    (storefront: Storefront) => {
      if (selectedTab === 'all') return true;
      if (selectedTab === 'active') return storefront.is_active;
      if (selectedTab === 'inactive') return !storefront.is_active;
      return true;
    }
  );

  // Calculate summary stats
  const stats = {
    total: userStorefronts.length,
    active: userStorefronts.filter((s: Storefront) => s.is_active).length,
    totalViews: userStorefronts.reduce(
      (sum: number, s: Storefront) => sum + (s.views_count || 0),
      0
    ),
    totalProducts: userStorefronts.reduce(
      (sum: number, s: Storefront) => sum + (s.products_count || 0),
      0
    ),
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-base-200">
        <div className="container mx-auto px-4 py-8">
          <div className="flex justify-center items-center min-h-[600px]">
            <div className="text-center">
              <span className="loading loading-spinner loading-lg text-primary"></span>
              <p className="mt-4 text-base-content/60">{t('common.loading')}</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Hero Section */}
      <div className="bg-gradient-to-br from-primary/10 via-secondary/5 to-accent/10 backdrop-blur-sm">
        <div className="container mx-auto px-4 py-12">
          <div className="flex flex-col lg:flex-row justify-between items-start lg:items-center gap-6">
            <div>
              <h1 className="text-4xl lg:text-5xl font-bold mb-3 flex items-center gap-3">
                <SparklesIcon className="w-10 h-10 text-primary" />
                {t('storefronts.myStorefronts')}
              </h1>
              <p className="text-lg text-base-content/70 max-w-2xl">
                {t('storefronts.manageDescription')}
              </p>
            </div>
            <button
              onClick={handleCreateStorefront}
              className="btn btn-primary btn-lg shadow-lg hover:shadow-xl transition-shadow"
            >
              <PlusIcon className="w-6 h-6" />
              {t('storefronts.createNew')}
            </button>
          </div>

          {/* Summary Stats */}
          {userStorefronts.length > 0 && (
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mt-8">
              <div className="stat bg-base-100 rounded-2xl shadow-md">
                <div className="stat-figure text-primary">
                  <ShoppingBagIcon className="w-8 h-8" />
                </div>
                <div className="stat-title">
                  {t('storefronts.totalStorefronts')}
                </div>
                <div className="stat-value text-2xl">{stats.total}</div>
                <div className="stat-desc">
                  {stats.active} {t('storefronts.activeCount')}
                </div>
              </div>

              <div className="stat bg-base-100 rounded-2xl shadow-md">
                <div className="stat-figure text-secondary">
                  <EyeIcon className="w-8 h-8" />
                </div>
                <div className="stat-title">{t('storefronts.totalViews')}</div>
                <div className="stat-value text-2xl">
                  {stats.totalViews.toLocaleString()}
                </div>
                <div className="stat-desc">
                  <ArrowTrendingUpIcon className="w-4 h-4 inline text-success" />
                  <span className="text-success ml-1">+12%</span>
                </div>
              </div>

              <div className="stat bg-base-100 rounded-2xl shadow-md">
                <div className="stat-figure text-accent">
                  <ShoppingBagIcon className="w-8 h-8" />
                </div>
                <div className="stat-title">
                  {t('storefronts.totalProducts')}
                </div>
                <div className="stat-value text-2xl">{stats.totalProducts}</div>
                <div className="stat-desc">
                  {t('storefronts.acrossAllStorefronts')}
                </div>
              </div>

              <div className="stat bg-base-100 rounded-2xl shadow-md">
                <div className="stat-figure text-success">
                  <CurrencyDollarIcon className="w-8 h-8" />
                </div>
                <div className="stat-title">
                  {t('storefronts.monthlyRevenue')}
                </div>
                <div className="stat-value text-2xl">€0</div>
                <div className="stat-desc">{t('storefronts.thisMonth')}</div>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Tabs */}
        {userStorefronts.length > 0 && (
          <div className="tabs tabs-boxed bg-base-100 shadow-sm mb-8 inline-flex">
            <a
              className={`tab ${selectedTab === 'all' ? 'tab-active' : ''}`}
              onClick={() => setSelectedTab('all')}
            >
              {t('common.all')} ({stats.total})
            </a>
            <a
              className={`tab ${selectedTab === 'active' ? 'tab-active' : ''}`}
              onClick={() => setSelectedTab('active')}
            >
              {t('storefronts.active')} ({stats.active})
            </a>
            <a
              className={`tab ${selectedTab === 'inactive' ? 'tab-active' : ''}`}
              onClick={() => setSelectedTab('inactive')}
            >
              {t('storefronts.inactive')} ({stats.total - stats.active})
            </a>
          </div>
        )}

        {/* Storefronts Grid */}
        {filteredStorefronts.length === 0 ? (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body text-center py-20">
              <div className="max-w-md mx-auto">
                <div className="mb-6">
                  <ShoppingBagIcon className="w-24 h-24 mx-auto text-base-content/20" />
                </div>
                <h3 className="text-2xl font-bold mb-3">
                  {selectedTab === 'all'
                    ? t('storefronts.noStorefronts')
                    : t('storefronts.noStorefrontsInCategory')}
                </h3>
                <p className="text-base-content/60 mb-8 text-lg">
                  {t('storefronts.createFirstStorefront')}
                </p>
                <button
                  onClick={handleCreateStorefront}
                  className="btn btn-primary btn-lg shadow-lg"
                >
                  <PlusIcon className="w-6 h-6" />
                  {t('storefronts.createStorefront')}
                </button>
              </div>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
            {filteredStorefronts.map((storefront: Storefront) => (
              <div
                key={storefront.id}
                className="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 overflow-hidden group"
              >
                {/* Card Header with Banner */}
                <figure className="relative h-48 bg-gradient-to-br from-primary/20 to-secondary/20 overflow-hidden">
                  {storefront.banner_url ? (
                    <img
                      src={storefront.banner_url}
                      alt={storefront.name}
                      className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center">
                      <ShoppingBagIcon className="w-16 h-16 text-base-content/20" />
                    </div>
                  )}

                  {/* Status Badge Overlay */}
                  <div className="absolute top-4 right-4">
                    <div className="flex items-center gap-2 bg-base-100/90 backdrop-blur-sm px-3 py-1.5 rounded-full shadow-lg">
                      {getStatusIcon(
                        storefront.is_active || false,
                        storefront.is_verified || false
                      )}
                      <span className="text-sm font-medium">
                        {storefront.is_active
                          ? t('storefronts.status.active')
                          : t('storefronts.status.inactive')}
                      </span>
                    </div>
                  </div>

                  {/* Logo Overlay */}
                  <div className="absolute bottom-0 left-6 transform translate-y-1/2">
                    <div className="avatar shadow-xl">
                      <div className="w-24 rounded-2xl ring-4 ring-base-100 bg-base-100">
                        {storefront.logo_url ? (
                          <img
                            src={storefront.logo_url}
                            alt={`${storefront.name} logo`}
                          />
                        ) : (
                          <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-primary/20 to-secondary/20">
                            <ShoppingBagIcon className="w-10 h-10 text-primary" />
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                </figure>

                <div className="card-body pt-16">
                  {/* Name and Rating */}
                  <div className="mb-3">
                    <h2 className="card-title text-xl mb-1">
                      {storefront.name}
                    </h2>
                    <div className="flex items-center gap-2">
                      <div className="flex items-center gap-0.5">
                        {[...Array(5)].map((_, i) => (
                          <StarSolidIcon
                            key={i}
                            className={`w-4 h-4 ${i < Math.floor(storefront.rating || 0) ? 'text-warning' : 'text-base-300'}`}
                          />
                        ))}
                      </div>
                      <span className="text-sm text-base-content/60">
                        {storefront.rating || 0} (
                        {storefront.reviews_count || 0} {t('common.reviews')})
                      </span>
                    </div>
                  </div>

                  {/* Description */}
                  <p className="text-sm text-base-content/70 line-clamp-2 mb-4">
                    {storefront.description || t('storefronts.noDescription')}
                  </p>

                  {/* Info Grid */}
                  <div className="grid grid-cols-2 gap-3 mb-4">
                    <div className="flex items-center gap-2 text-sm">
                      <MapPinIcon className="w-4 h-4 text-base-content/50" />
                      <span className="text-base-content/70 truncate">
                        {storefront.location?.city || t('common.notSpecified')}
                      </span>
                    </div>
                    <div className="flex items-center gap-2 text-sm">
                      <ClockIcon className="w-4 h-4 text-base-content/50" />
                      {formatBusinessHours(storefront.settings)}
                    </div>
                  </div>

                  {/* Quick Stats */}
                  <div className="grid grid-cols-3 gap-2 mb-4">
                    <div className="text-center p-3 bg-base-200 rounded-xl">
                      <div className="text-2xl font-bold text-primary">
                        {storefront.products_count || 0}
                      </div>
                      <div className="text-xs text-base-content/60">
                        {t('storefronts.products')}
                      </div>
                    </div>
                    <div className="text-center p-3 bg-base-200 rounded-xl">
                      <div className="text-2xl font-bold text-secondary">
                        {storefront.views_count || 0}
                      </div>
                      <div className="text-xs text-base-content/60">
                        {t('storefronts.views')}
                      </div>
                    </div>
                    <div className="text-center p-3 bg-base-200 rounded-xl">
                      <div className="text-2xl font-bold text-accent">0</div>
                      <div className="text-xs text-base-content/60">
                        {t('storefronts.orders')}
                      </div>
                    </div>
                  </div>

                  {/* Action Buttons */}
                  <div className="card-actions">
                    <Link
                      href={`/storefronts/${storefront.slug}/dashboard`}
                      className="btn btn-primary flex-1"
                    >
                      <ChartBarIcon className="w-5 h-5" />
                      {t('storefronts.dashboard')}
                    </Link>

                    <div className="dropdown dropdown-end">
                      <label tabIndex={0} className="btn btn-ghost btn-square">
                        <CogIcon className="w-5 h-5" />
                      </label>
                      <ul
                        tabIndex={0}
                        className="dropdown-content z-[1] menu menu-sm p-2 shadow-lg bg-base-100 rounded-box w-56 mt-2"
                      >
                        <li>
                          <Link
                            href={`/storefronts/${storefront.slug}`}
                            className="gap-3"
                          >
                            <EyeIcon className="w-4 h-4" />
                            {t('storefronts.viewPublicPage')}
                          </Link>
                        </li>
                        <li>
                          <Link
                            href={`/storefronts/${storefront.slug}/edit`}
                            className="gap-3"
                          >
                            <PencilIcon className="w-4 h-4" />
                            {t('storefronts.editDetails')}
                          </Link>
                        </li>
                        <li>
                          <Link
                            href={`/storefronts/${storefront.slug}/products`}
                            className="gap-3"
                          >
                            <ShoppingBagIcon className="w-4 h-4" />
                            {t('storefronts.manageProducts')}
                          </Link>
                        </li>
                        <li>
                          <Link
                            href={`/storefronts/${storefront.slug}/staff`}
                            className="gap-3"
                          >
                            <UserGroupIcon className="w-4 h-4" />
                            {t('storefronts.manageStaff')}
                          </Link>
                        </li>
                        <li>
                          <Link
                            href={`/storefronts/${storefront.slug}/reviews`}
                            className="gap-3"
                          >
                            <StarIcon className="w-4 h-4" />
                            {t('storefronts.manageReviews')}
                          </Link>
                        </li>
                        <li>
                          <Link
                            href={`/storefronts/${storefront.slug}/messages`}
                            className="gap-3"
                          >
                            <ChatBubbleLeftRightIcon className="w-4 h-4" />
                            {t('storefronts.messages')}
                          </Link>
                        </li>
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
