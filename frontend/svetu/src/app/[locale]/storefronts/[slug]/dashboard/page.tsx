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
  ShoppingBagIcon,
  BellIcon,
  ChatBubbleLeftRightIcon,
  Cog6ToothIcon,
  PlusIcon,
  TagIcon,
  TruckIcon,
  UsersIcon,
  ClipboardDocumentListIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  ClockIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';

export default function StorefrontDashboardPage() {
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
      if (!user) {
        setAccessDenied(true);
        router.push(`/${locale}/storefronts/${slug}`);
        return;
      }

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
            {t('storefronts.dashboardAccessDenied')}
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
                href="/profile/storefronts"
                className="btn btn-ghost btn-sm btn-square"
              >
                <ArrowLeftIcon className="w-5 h-5" />
              </Link>
              <div>
                <h1 className="text-2xl font-bold">{currentStorefront.name}</h1>
                <p className="text-sm text-base-content/60">
                  {t('storefronts.managementDashboard')}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-2">
              {/* Analytics Button */}
              <Link
                href={`/${locale}/storefronts/${currentStorefront.slug}/analytics`}
                className="btn btn-outline btn-sm"
              >
                <ChartBarIcon className="w-4 h-4" />
                {t('storefronts.viewAnalytics')}
              </Link>

              {/* Settings Button */}
              <Link
                href={`/${locale}/storefronts/${currentStorefront.slug}/settings`}
                className="btn btn-ghost btn-sm btn-square"
              >
                <Cog6ToothIcon className="w-5 h-5" />
              </Link>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Quick Stats Bar */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
          <div className="stat bg-base-100 rounded-lg shadow-sm">
            <div className="stat-figure text-primary">
              <ShoppingBagIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('storefronts.activeProducts')}</div>
            <div className="stat-value text-primary">48</div>
            <div className="stat-desc">
              {t('storefronts.totalProducts', { count: 52 })}
            </div>
          </div>

          <div className="stat bg-base-100 rounded-lg shadow-sm">
            <div className="stat-figure text-secondary">
              <ClipboardDocumentListIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('storefronts.pendingOrders')}</div>
            <div className="stat-value text-secondary">12</div>
            <div className="stat-desc">{t('storefronts.requiresAction')}</div>
          </div>

          <div className="stat bg-base-100 rounded-lg shadow-sm">
            <div className="stat-figure text-accent">
              <ChatBubbleLeftRightIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('storefronts.unreadMessages')}</div>
            <div className="stat-value text-accent">7</div>
            <div className="stat-desc">{t('storefronts.fromCustomers')}</div>
          </div>

          <div className="stat bg-base-100 rounded-lg shadow-sm">
            <div className="stat-figure text-warning">
              <ExclamationTriangleIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('storefronts.lowStock')}</div>
            <div className="stat-value text-warning">3</div>
            <div className="stat-desc">{t('storefronts.productsLowStock')}</div>
          </div>
        </div>

        {/* Main Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column - Main Actions */}
          <div className="lg:col-span-2 space-y-6">
            {/* Quick Actions */}
            <div className="card bg-base-100 shadow-md">
              <div className="card-body">
                <h2 className="card-title text-lg mb-4">
                  {t('storefronts.quickActions')}
                </h2>
                <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                  <Link
                    href={`/${locale}/storefronts/${currentStorefront.slug}/products/new`}
                    className="btn btn-primary"
                  >
                    <PlusIcon className="w-5 h-5" />
                    {t('storefronts.addProduct')}
                  </Link>
                  <Link
                    href={`/${locale}/storefronts/${currentStorefront.slug}/orders`}
                    className="btn btn-outline"
                  >
                    <ClipboardDocumentListIcon className="w-5 h-5" />
                    {t('storefronts.manageOrders')}
                  </Link>
                  <Link
                    href={`/${locale}/storefronts/${currentStorefront.slug}/products`}
                    className="btn btn-outline"
                  >
                    <ShoppingBagIcon className="w-5 h-5" />
                    {t('storefronts.products')}
                  </Link>
                  <Link
                    href={`/${locale}/storefronts/${currentStorefront.slug}/customers`}
                    className="btn btn-outline"
                  >
                    <UsersIcon className="w-5 h-5" />
                    {t('storefronts.customers')}
                  </Link>
                  <Link
                    href={`/${locale}/storefronts/${currentStorefront.slug}/promotions`}
                    className="btn btn-outline"
                  >
                    <TagIcon className="w-5 h-5" />
                    {t('storefronts.promotions')}
                  </Link>
                  <Link
                    href={`/${locale}/storefronts/${currentStorefront.slug}/delivery`}
                    className="btn btn-outline"
                  >
                    <TruckIcon className="w-5 h-5" />
                    {t('storefronts.delivery')}
                  </Link>
                </div>
              </div>
            </div>

            {/* Recent Orders */}
            <div className="card bg-base-100 shadow-md">
              <div className="card-body">
                <div className="flex justify-between items-center mb-4">
                  <h2 className="card-title text-lg">
                    {t('storefronts.recentOrders')}
                  </h2>
                  <Link
                    href={`/${locale}/storefronts/${currentStorefront.slug}/orders`}
                    className="btn btn-ghost btn-sm"
                  >
                    {t('common.viewAll')}
                  </Link>
                </div>

                <div className="overflow-x-auto">
                  <table className="table table-zebra">
                    <thead>
                      <tr>
                        <th>{t('storefronts.orderId')}</th>
                        <th>{t('storefronts.customer')}</th>
                        <th>{t('storefronts.products')}</th>
                        <th>{t('storefronts.total')}</th>
                        <th>{t('storefronts.status')}</th>
                        <th>{t('common.actions')}</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr>
                        <td>#1234</td>
                        <td>John Doe</td>
                        <td>2 items</td>
                        <td>$89.99</td>
                        <td>
                          <span className="badge badge-warning badge-sm">
                            <ClockIcon className="w-3 h-3 mr-1" />
                            {t('storefronts.pending')}
                          </span>
                        </td>
                        <td>
                          <button className="btn btn-ghost btn-xs">
                            {t('common.view')}
                          </button>
                        </td>
                      </tr>
                      <tr>
                        <td>#1233</td>
                        <td>Jane Smith</td>
                        <td>1 item</td>
                        <td>$45.00</td>
                        <td>
                          <span className="badge badge-info badge-sm">
                            <TruckIcon className="w-3 h-3 mr-1" />
                            {t('storefronts.shipping')}
                          </span>
                        </td>
                        <td>
                          <button className="btn btn-ghost btn-xs">
                            {t('common.view')}
                          </button>
                        </td>
                      </tr>
                      <tr>
                        <td>#1232</td>
                        <td>Bob Wilson</td>
                        <td>3 items</td>
                        <td>$125.50</td>
                        <td>
                          <span className="badge badge-success badge-sm">
                            <CheckCircleIcon className="w-3 h-3 mr-1" />
                            {t('storefronts.completed')}
                          </span>
                        </td>
                        <td>
                          <button className="btn btn-ghost btn-xs">
                            {t('common.view')}
                          </button>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            {/* Low Stock Alert */}
            <div className="card bg-warning/10 border-warning shadow-md">
              <div className="card-body">
                <div className="flex items-center gap-3 mb-4">
                  <ExclamationTriangleIcon className="w-6 h-6 text-warning" />
                  <h2 className="card-title text-lg">
                    {t('storefronts.lowStockAlert')}
                  </h2>
                </div>

                <div className="space-y-2">
                  <div className="flex justify-between items-center p-2 bg-base-100 rounded">
                    <span>iPhone 15 Pro Max</span>
                    <span className="badge badge-warning">
                      {t('storefronts.stockLeft', { count: 2 })}
                    </span>
                  </div>
                  <div className="flex justify-between items-center p-2 bg-base-100 rounded">
                    <span>MacBook Pro M3</span>
                    <span className="badge badge-warning">
                      {t('storefronts.stockLeft', { count: 1 })}
                    </span>
                  </div>
                  <div className="flex justify-between items-center p-2 bg-base-100 rounded">
                    <span>Sony WH-1000XM5</span>
                    <span className="badge badge-error">
                      {t('storefronts.outOfStock')}
                    </span>
                  </div>
                </div>

                <Link
                  href={`/${locale}/storefronts/${currentStorefront.slug}/inventory`}
                  className="btn btn-warning btn-sm w-full mt-4"
                >
                  {t('storefronts.manageInventory')}
                </Link>
              </div>
            </div>
          </div>

          {/* Right Column - Notifications & Messages */}
          <div className="space-y-6">
            {/* Notifications */}
            <div className="card bg-base-100 shadow-md">
              <div className="card-body">
                <div className="flex justify-between items-center mb-4">
                  <h2 className="card-title text-lg">
                    <BellIcon className="w-5 h-5" />
                    {t('storefronts.notifications')}
                  </h2>
                  <span className="badge badge-primary">
                    5 {t('common.new')}
                  </span>
                </div>

                <div className="space-y-3">
                  <div className="p-3 bg-base-200 rounded-lg">
                    <div className="flex items-start gap-3">
                      <div className="p-1 bg-primary/10 rounded">
                        <ShoppingBagIcon className="w-4 h-4 text-primary" />
                      </div>
                      <div className="flex-1">
                        <p className="text-sm font-medium">
                          {t('storefronts.newOrderReceived')}
                        </p>
                        <p className="text-xs text-base-content/60">
                          5 {t('common.minutesAgo')}
                        </p>
                      </div>
                    </div>
                  </div>

                  <div className="p-3 bg-base-200 rounded-lg">
                    <div className="flex items-start gap-3">
                      <div className="p-1 bg-accent/10 rounded">
                        <ChatBubbleLeftRightIcon className="w-4 h-4 text-accent" />
                      </div>
                      <div className="flex-1">
                        <p className="text-sm font-medium">
                          {t('storefronts.newMessage')}
                        </p>
                        <p className="text-xs text-base-content/60">
                          15 {t('common.minutesAgo')}
                        </p>
                      </div>
                    </div>
                  </div>

                  <div className="p-3 bg-base-200 rounded-lg">
                    <div className="flex items-start gap-3">
                      <div className="p-1 bg-warning/10 rounded">
                        <ExclamationTriangleIcon className="w-4 h-4 text-warning" />
                      </div>
                      <div className="flex-1">
                        <p className="text-sm font-medium">
                          {t('storefronts.productLowStock')}
                        </p>
                        <p className="text-xs text-base-content/60">
                          1 {t('common.hourAgo')}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>

                <button className="btn btn-ghost btn-sm w-full mt-4">
                  {t('storefronts.viewAllNotifications')}
                </button>
              </div>
            </div>

            {/* Recent Messages */}
            <div className="card bg-base-100 shadow-md">
              <div className="card-body">
                <div className="flex justify-between items-center mb-4">
                  <h2 className="card-title text-lg">
                    <ChatBubbleLeftRightIcon className="w-5 h-5" />
                    {t('storefronts.messages')}
                  </h2>
                  <span className="badge badge-accent">
                    7 {t('storefronts.unread')}
                  </span>
                </div>

                <div className="space-y-3">
                  <div className="flex items-start gap-3 p-2 hover:bg-base-200 rounded cursor-pointer">
                    <div className="avatar">
                      <div className="w-10 h-10 rounded-full bg-primary/10">
                        <span className="text-lg font-bold flex items-center justify-center h-full">
                          JD
                        </span>
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex justify-between">
                        <p className="text-sm font-medium">John Doe</p>
                        <p className="text-xs text-base-content/60">5m</p>
                      </div>
                      <p className="text-sm text-base-content/70 truncate">
                        {t('storefronts.messagePreview')}
                      </p>
                    </div>
                    <div className="badge badge-accent badge-xs"></div>
                  </div>

                  <div className="flex items-start gap-3 p-2 hover:bg-base-200 rounded cursor-pointer">
                    <div className="avatar">
                      <div className="w-10 h-10 rounded-full bg-secondary/10">
                        <span className="text-lg font-bold flex items-center justify-center h-full">
                          JS
                        </span>
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex justify-between">
                        <p className="text-sm font-medium">Jane Smith</p>
                        <p className="text-xs text-base-content/60">1h</p>
                      </div>
                      <p className="text-sm text-base-content/70 truncate">
                        {t('storefronts.messagePreview2')}
                      </p>
                    </div>
                  </div>
                </div>

                <Link
                  href={`/${locale}/chat`}
                  className="btn btn-ghost btn-sm w-full mt-4"
                >
                  {t('storefronts.goToMessages')}
                </Link>
              </div>
            </div>

            {/* Store Status */}
            <div className="card bg-base-100 shadow-md">
              <div className="card-body">
                <h2 className="card-title text-lg mb-4">
                  {t('storefronts.storeStatus')}
                </h2>

                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-sm">
                      {t('storefronts.storeOpen')}
                    </span>
                    <input
                      type="checkbox"
                      className="toggle toggle-success"
                      defaultChecked
                    />
                  </div>

                  <div className="divider my-2"></div>

                  <div className="flex justify-between items-center">
                    <span className="text-sm">
                      {t('storefronts.acceptingOrders')}
                    </span>
                    <input
                      type="checkbox"
                      className="toggle toggle-primary"
                      defaultChecked
                    />
                  </div>

                  <div className="flex justify-between items-center">
                    <span className="text-sm">
                      {t('storefronts.vacationMode')}
                    </span>
                    <input type="checkbox" className="toggle" />
                  </div>
                </div>

                <Link
                  href={`/${locale}/storefronts/${currentStorefront.slug}/settings`}
                  className="btn btn-outline btn-sm w-full mt-4"
                >
                  <Cog6ToothIcon className="w-4 h-4" />
                  {t('storefronts.storeSettings')}
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
