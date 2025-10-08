'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTranslations, useLocale } from 'next-intl';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchStorefrontBySlug,
  fetchDashboardStats,
  fetchRecentOrders,
  fetchLowStockProducts,
  fetchDashboardNotifications,
  selectDashboardStats,
  selectRecentOrders,
  selectLowStockProducts,
  selectDashboardNotifications,
  selectIsLoadingDashboard,
} from '@/store/slices/b2cStoreSlice';
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
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const tAdmin = useTranslations('admin');
  const tStorefronts = useTranslations('storefronts');
  const tNotifications = useTranslations('notifications');
  const locale = useLocale();
  const router = useRouter();
  const params = useParams();
  const dispatch = useAppDispatch();
  const slug = params?.slug as string;
  const { user } = useAuth();

  const { currentStorefront, isLoading } = useAppSelector(
    (state) => state.b2cStores
  );

  const dashboardStats = useAppSelector(selectDashboardStats);
  const recentOrders = useAppSelector(selectRecentOrders);
  const lowStockProducts = useAppSelector(selectLowStockProducts);
  const notifications = useAppSelector(selectDashboardNotifications);
  const isLoadingDashboard = useAppSelector(selectIsLoadingDashboard);

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
        router.push(`/${locale}/b2c/${slug}`);
        return;
      }

      if (currentStorefront.user_id !== user.id) {
        setAccessDenied(true);
        router.push(`/${locale}/b2c/${slug}`);
      }
    }
  }, [currentStorefront, user, isLoading, router, slug, locale]);

  // Загрузка данных dashboard
  useEffect(() => {
    if (currentStorefront && slug && !accessDenied) {
      dispatch(fetchDashboardStats(slug));
      dispatch(fetchRecentOrders({ slug, limit: 5 }));
      dispatch(fetchLowStockProducts(slug));
      dispatch(fetchDashboardNotifications({ slug, limit: 10 }));
    }
  }, [dispatch, currentStorefront, slug, accessDenied]);

  if (accessDenied) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">🔒</div>
          <h1 className="text-2xl font-bold mb-2">{tCommon('accessDenied')}</h1>
          <p className="text-base-content/60">{t('dashboardAccessDenied')}</p>
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
                {tAdmin('common.loading')}
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
                href={`/${locale}/profile/b2c`}
                className="btn btn-ghost btn-sm btn-square"
              >
                <ArrowLeftIcon className="w-5 h-5" />
              </Link>
              <div>
                <h1 className="text-2xl font-bold">{currentStorefront.name}</h1>
                <p className="text-sm text-base-content/60">
                  {t('managementDashboard')}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-2">
              {/* Analytics Button */}
              <Link
                href={`/${locale}/b2c/${currentStorefront.slug}/analytics`}
                className="btn btn-outline btn-sm"
              >
                <ChartBarIcon className="w-4 h-4" />
                {t('viewAnalytics')}
              </Link>

              {/* Settings Button */}
              <Link
                href={`/${locale}/b2c/${currentStorefront.slug}/settings`}
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
            <div className="stat-title">{t('activeProducts')}</div>
            <div className="stat-value text-primary">
              {isLoadingDashboard ? (
                <span className="loading loading-dots loading-sm"></span>
              ) : (
                dashboardStats?.activeProducts || 0
              )}
            </div>
            <div className="stat-desc">
              {t('totalProducts', {
                count: dashboardStats?.totalProducts || 0,
              })}
            </div>
          </div>

          <div className="stat bg-base-100 rounded-lg shadow-sm">
            <div className="stat-figure text-secondary">
              <ClipboardDocumentListIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('pendingOrders')}</div>
            <div className="stat-value text-secondary">
              {isLoadingDashboard ? (
                <span className="loading loading-dots loading-sm"></span>
              ) : (
                dashboardStats?.pendingOrders || 0
              )}
            </div>
            <div className="stat-desc">{t('requiresAction')}</div>
          </div>

          <div className="stat bg-base-100 rounded-lg shadow-sm">
            <div className="stat-figure text-accent">
              <ChatBubbleLeftRightIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('unreadMessages')}</div>
            <div className="stat-value text-accent">
              {isLoadingDashboard ? (
                <span className="loading loading-dots loading-sm"></span>
              ) : (
                dashboardStats?.unreadMessages || 0
              )}
            </div>
            <div className="stat-desc">{t('fromCustomers')}</div>
          </div>

          <div className="stat bg-base-100 rounded-lg shadow-sm">
            <div className="stat-figure text-warning">
              <ExclamationTriangleIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('lowStock')}</div>
            <div className="stat-value text-warning">
              {isLoadingDashboard ? (
                <span className="loading loading-dots loading-sm"></span>
              ) : (
                dashboardStats?.lowStockProducts || 0
              )}
            </div>
            <div className="stat-desc">{t('productsLowStock')}</div>
          </div>
        </div>

        {/* Main Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column - Main Actions */}
          <div className="lg:col-span-2 space-y-6">
            {/* Quick Actions */}
            <div className="card bg-base-100 shadow-md">
              <div className="card-body">
                <h2 className="card-title text-lg mb-4">{t('quickActions')}</h2>
                <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                  <Link
                    href={`/${locale}/b2c/${currentStorefront.slug}/products/new`}
                    className="btn btn-primary"
                  >
                    <PlusIcon className="w-5 h-5" />
                    {t('addProduct')}
                  </Link>
                  <Link
                    href={`/${locale}/b2c/${currentStorefront.slug}/orders`}
                    className="btn btn-outline"
                  >
                    <ClipboardDocumentListIcon className="w-5 h-5" />
                    {t('manageOrders')}
                  </Link>
                  <Link
                    href={`/${locale}/b2c/${currentStorefront.slug}/products`}
                    className="btn btn-outline"
                  >
                    <ShoppingBagIcon className="w-5 h-5" />
                    {tStorefronts('products.title')}
                  </Link>
                  <Link
                    href={`/${locale}/b2c/${currentStorefront.slug}/customers`}
                    className="btn btn-outline"
                  >
                    <UsersIcon className="w-5 h-5" />
                    {t('customers')}
                  </Link>
                  <Link
                    href={`/${locale}/b2c/${currentStorefront.slug}/promotions`}
                    className="btn btn-outline"
                  >
                    <TagIcon className="w-5 h-5" />
                    {t('promotions')}
                  </Link>
                  <Link
                    href={`/${locale}/b2c/${currentStorefront.slug}/delivery`}
                    className="btn btn-outline"
                  >
                    <TruckIcon className="w-5 h-5" />
                    {tStorefronts('delivery.title') ||
                      tStorefronts('delivery.shipping')}
                  </Link>
                </div>
              </div>
            </div>

            {/* Recent Orders */}
            <div className="card bg-base-100 shadow-md">
              <div className="card-body">
                <div className="flex justify-between items-center mb-4">
                  <h2 className="card-title text-lg">{t('recentOrders')}</h2>
                  <Link
                    href={`/${locale}/b2c/${currentStorefront.slug}/orders`}
                    className="btn btn-ghost btn-sm"
                  >
                    {tCommon('viewAll')}
                  </Link>
                </div>

                {!recentOrders || recentOrders.length === 0 ? (
                  <div className="text-center py-8">
                    <p className="text-base-content/60">
                      {t('noRecentOrders')}
                    </p>
                  </div>
                ) : (
                  <div className="overflow-x-auto">
                    <table className="table table-zebra">
                      <thead>
                        <tr>
                          <th>{t('orderId')}</th>
                          <th>{t('customer')}</th>
                          <th>{tStorefronts('products.title')}</th>
                          <th>{t('total')}</th>
                          <th>{tCommon('status')}</th>
                          <th>{tCommon('actions')}</th>
                        </tr>
                      </thead>
                      <tbody>
                        {recentOrders &&
                          recentOrders.map((order) => (
                            <tr key={order.id}>
                              <td>{order.order_id}</td>
                              <td>{order.customer}</td>
                              <td>
                                {t('itemsCount', {
                                  count: order.items_count,
                                })}
                              </td>
                              <td>
                                {order.currency} {order.total.toFixed(2)}
                              </td>
                              <td>
                                {order.status === 'pending' && (
                                  <span className="badge badge-warning badge-sm">
                                    <ClockIcon className="w-3 h-3 mr-1" />
                                    {t('pending')}
                                  </span>
                                )}
                                {order.status === 'shipping' && (
                                  <span className="badge badge-info badge-sm">
                                    <TruckIcon className="w-3 h-3 mr-1" />
                                    {t('shipping')}
                                  </span>
                                )}
                                {order.status === 'completed' && (
                                  <span className="badge badge-success badge-sm">
                                    <CheckCircleIcon className="w-3 h-3 mr-1" />
                                    {t('completed')}
                                  </span>
                                )}
                                {order.status === 'cancelled' && (
                                  <span className="badge badge-error badge-sm">
                                    {t('cancelled')}
                                  </span>
                                )}
                              </td>
                              <td>
                                <Link
                                  href={`/${locale}/b2c/${currentStorefront.slug}/orders/${order.id}`}
                                  className="btn btn-ghost btn-xs"
                                >
                                  {tCommon('view')}
                                </Link>
                              </td>
                            </tr>
                          ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>
            </div>

            {/* Low Stock Alert */}
            <div className="card bg-warning/10 border-warning shadow-md">
              <div className="card-body">
                <div className="flex items-center gap-3 mb-4">
                  <ExclamationTriangleIcon className="w-6 h-6 text-warning" />
                  <h2 className="card-title text-lg">{t('lowStockAlert')}</h2>
                </div>

                {!lowStockProducts || lowStockProducts.length === 0 ? (
                  <div className="text-center py-4">
                    <p className="text-base-content/60">
                      {t('allProductsInStock')}
                    </p>
                  </div>
                ) : (
                  <div className="space-y-2">
                    {lowStockProducts &&
                      lowStockProducts.map((product) => (
                        <div
                          key={product.id}
                          className="flex justify-between items-center p-2 bg-base-100 rounded"
                        >
                          <span className="truncate flex-1 mr-2">
                            {product.name}
                          </span>
                          {product.stock_quantity === 0 ? (
                            <span className="badge badge-error">
                              {t('outOfStock')}
                            </span>
                          ) : (
                            <span className="badge badge-warning">
                              {t('stockLeft', {
                                count: product.stock_quantity,
                              })}
                            </span>
                          )}
                        </div>
                      ))}
                  </div>
                )}

                <Link
                  href={`/${locale}/b2c/${currentStorefront.slug}/inventory`}
                  className="btn btn-warning btn-sm w-full mt-4"
                >
                  {t('manageInventory')}
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
                    {t('notifications')}
                  </h2>
                  {notifications &&
                    notifications.filter((n) => !n.is_read).length > 0 && (
                      <span className="badge badge-primary">
                        {notifications.filter((n) => !n.is_read).length}{' '}
                        {tCommon('new')}
                      </span>
                    )}
                </div>

                {!notifications || notifications.length === 0 ? (
                  <div className="text-center py-4">
                    <p className="text-base-content/60">
                      {t('noNotifications')}
                    </p>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {notifications &&
                      notifications.slice(0, 5).map((notification) => {
                        const getNotificationIcon = () => {
                          switch (notification.type) {
                            case 'order':
                              return (
                                <ShoppingBagIcon className="w-4 h-4 text-primary" />
                              );
                            case 'message':
                              return (
                                <ChatBubbleLeftRightIcon className="w-4 h-4 text-accent" />
                              );
                            case 'stock':
                              return (
                                <ExclamationTriangleIcon className="w-4 h-4 text-warning" />
                              );
                            default:
                              return (
                                <BellIcon className="w-4 h-4 text-base-content" />
                              );
                          }
                        };

                        const getIconBgColor = () => {
                          switch (notification.type) {
                            case 'order':
                              return 'bg-primary/10';
                            case 'message':
                              return 'bg-accent/10';
                            case 'stock':
                              return 'bg-warning/10';
                            default:
                              return 'bg-base-300';
                          }
                        };

                        return (
                          <div
                            key={notification.id}
                            className={`p-3 bg-base-200 rounded-lg ${!notification.is_read ? 'border-l-4 border-primary' : ''}`}
                          >
                            <div className="flex items-start gap-3">
                              <div
                                className={`p-1 ${getIconBgColor()} rounded`}
                              >
                                {getNotificationIcon()}
                              </div>
                              <div className="flex-1">
                                <p className="text-sm font-medium">
                                  {notification.title.startsWith(
                                    'notifications.'
                                  )
                                    ? tNotifications(
                                        notification.title.replace(
                                          'notifications.',
                                          ''
                                        )
                                      )
                                    : notification.title}
                                </p>
                                {notification.message && (
                                  <p className="text-xs text-base-content/70 mt-1">
                                    {notification.message.startsWith(
                                      'notifications.'
                                    )
                                      ? tNotifications(
                                          notification.message.replace(
                                            'notifications.',
                                            ''
                                          )
                                        )
                                      : notification.message}
                                  </p>
                                )}
                                <p className="text-xs text-base-content/60 mt-1">
                                  {new Date(
                                    notification.created_at
                                  ).toLocaleString(locale)}
                                </p>
                              </div>
                            </div>
                          </div>
                        );
                      })}
                  </div>
                )}

                <button className="btn btn-ghost btn-sm w-full mt-4">
                  {t('viewAllNotifications')}
                </button>
              </div>
            </div>

            {/* Recent Messages */}
            <div className="card bg-base-100 shadow-md">
              <div className="card-body">
                <div className="flex justify-between items-center mb-4">
                  <h2 className="card-title text-lg">
                    <ChatBubbleLeftRightIcon className="w-5 h-5" />
                    {t('messages')}
                  </h2>
                  <span className="badge badge-accent">7 {t('unread')}</span>
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
                        {t('messagePreview')}
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
                        {t('messagePreview2')}
                      </p>
                    </div>
                  </div>
                </div>

                <Link
                  href={`/${locale}/chat`}
                  className="btn btn-ghost btn-sm w-full mt-4"
                >
                  {t('goToMessages')}
                </Link>
              </div>
            </div>

            {/* Store Status */}
            <div className="card bg-base-100 shadow-md">
              <div className="card-body">
                <h2 className="card-title text-lg mb-4">{t('storeStatus')}</h2>

                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-sm">{t('storeOpen')}</span>
                    <input
                      type="checkbox"
                      className="toggle toggle-success"
                      defaultChecked
                    />
                  </div>

                  <div className="divider my-2"></div>

                  <div className="flex justify-between items-center">
                    <span className="text-sm">{t('acceptingOrders')}</span>
                    <input
                      type="checkbox"
                      className="toggle toggle-primary"
                      defaultChecked
                    />
                  </div>

                  <div className="flex justify-between items-center">
                    <span className="text-sm">{t('vacationMode')}</span>
                    <input type="checkbox" className="toggle" />
                  </div>
                </div>

                <Link
                  href={`/${locale}/b2c/${currentStorefront.slug}/settings`}
                  className="btn btn-outline btn-sm w-full mt-4"
                >
                  <Cog6ToothIcon className="w-4 h-4" />
                  {t('storeSettings')}
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
