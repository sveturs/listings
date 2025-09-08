'use client';

import { useState, useEffect } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import {
  marketplaceOrdersService,
  type MarketplaceOrder,
} from '@/services/marketplaceOrders';
import { balanceService } from '@/services/balance';
import SafeImage from '@/components/SafeImage';
import { toast } from 'react-hot-toast';

export default function MyPurchasesPage() {
  const locale = useLocale();
  const t = useTranslations('orders');
  const tHome = useTranslations('marketplace.homeContent');
  const tProfile = useTranslations('profile');
  const tCommon = useTranslations('common');
  const { isAuthenticated } = useAuth();

  const [orders, setOrders] = useState<MarketplaceOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [processingOrderId, setProcessingOrderId] = useState<number | null>(
    null
  );

  const limit = 20;

  useEffect(() => {
    if (isAuthenticated) {
      fetchOrders();
    }
  }, [isAuthenticated, page]); // eslint-disable-line react-hooks/exhaustive-deps

  const fetchOrders = async () => {
    try {
      setLoading(true);
      const response = await marketplaceOrdersService.getMyPurchases(
        page,
        limit
      );
      setOrders(response.orders);
      setTotal(response.total);
    } catch (error) {
      console.error('Error fetching purchases:', error);
      toast.error(t('fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const handleConfirmDelivery = async (orderId: number) => {
    if (!confirm(t('confirmDeliveryPrompt'))) {
      return;
    }

    setProcessingOrderId(orderId);
    try {
      await marketplaceOrdersService.confirmDelivery(orderId);
      toast.success(t('deliveryConfirmed'));
      await fetchOrders(); // Обновляем список
    } catch (error) {
      console.error('Error confirming delivery:', error);
      toast.error(t('confirmDeliveryError'));
    } finally {
      setProcessingOrderId(null);
    }
  };

  const handleOpenDispute = async (orderId: number) => {
    const reason = prompt(t('disputeReasonPrompt'));
    if (!reason || reason.length < 10) {
      toast.error(t('disputeReasonTooShort'));
      return;
    }

    setProcessingOrderId(orderId);
    try {
      await marketplaceOrdersService.openDispute(orderId, reason);
      toast.success(t('disputeOpened'));
      await fetchOrders(); // Обновляем список
    } catch (error) {
      console.error('Error opening dispute:', error);
      toast.error(t('disputeError'));
    } finally {
      setProcessingOrderId(null);
    }
  };

  const totalPages = Math.ceil(total / limit);

  if (loading && orders.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center min-h-[400px]">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      {/* Заголовок и навигация */}
      <div className="mb-8">
        <div className="breadcrumbs text-sm mb-4">
          <ul>
            <li>
              <Link href={`/${locale}`}>{tHome('title')}</Link>
            </li>
            <li>
              <Link href={`/${locale}/profile`}>{tProfile('title')}</Link>
            </li>
            <li>{t('myPurchases')}</li>
          </ul>
        </div>

        <h1 className="text-3xl font-bold mb-4">{t('myPurchases')}</h1>

        {/* Табы */}
        <div className="tabs tabs-boxed mb-6">
          <Link
            href={`/${locale}/profile/orders/purchases`}
            className="tab tab-active"
          >
            {t('purchases')}
          </Link>
          <Link href={`/${locale}/profile/orders/sales`} className="tab">
            {t('sales')}
          </Link>
        </div>
      </div>

      {/* Список заказов */}
      {orders.length === 0 ? (
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body text-center py-16">
            <h3 className="text-xl font-semibold mb-2">{t('noPurchases')}</h3>
            <p className="text-base-content/70 mb-6">{t('noPurchasesDesc')}</p>
            <Link href={`/${locale}/search`} className="btn btn-primary">
              {t('startShopping')}
            </Link>
          </div>
        </div>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <div key={order.id} className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <div className="flex flex-col lg:flex-row gap-6">
                  {/* Изображение товара */}
                  <figure className="relative w-full lg:w-32 h-32 flex-shrink-0 bg-base-200 rounded-lg overflow-hidden">
                    <SafeImage
                      src={order.listing?.images?.[0]?.public_url}
                      alt={order.listing?.title || 'Order item'}
                      fill
                      className="object-cover"
                    />
                  </figure>

                  {/* Информация о заказе */}
                  <div className="flex-grow">
                    <div className="flex justify-between items-start mb-4">
                      <div>
                        <h3 className="text-xl font-semibold mb-1">
                          {order.listing?.title || t('unknownItem')}
                        </h3>
                        <p className="text-sm text-base-content/70">
                          {t('orderNumber')}: #{order.id}
                        </p>
                        <p className="text-sm text-base-content/70">
                          {t('orderDate')}:{' '}
                          {new Date(order.created_at).toLocaleDateString(
                            locale
                          )}
                        </p>
                      </div>

                      {/* Статус */}
                      <div className="text-right">
                        <div
                          className={`badge ${marketplaceOrdersService.getStatusColor(order.status)} badge-lg`}
                        >
                          {marketplaceOrdersService.getStatusLabel(
                            order.status,
                            locale
                          )}
                        </div>
                        {order.protection_expires_at &&
                          order.status === 'delivered' && (
                            <p className="text-sm text-warning mt-2">
                              ⏱️ {t('protectionEnds')}:{' '}
                              {marketplaceOrdersService.getProtectionTimeLeft(
                                order
                              )}
                            </p>
                          )}
                      </div>
                    </div>

                    {/* Детали */}
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                      <div>
                        <p className="text-sm text-base-content/70">
                          {t('seller')}
                        </p>
                        <p className="font-medium">
                          {order.seller?.name || 'Unknown'}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-base-content/70">
                          {t('price')}
                        </p>
                        <p className="font-medium text-primary">
                          {balanceService.formatAmount(order.item_price, 'RSD')}
                        </p>
                      </div>
                      {order.tracking_number && (
                        <div>
                          <p className="text-sm text-base-content/70">
                            {t('trackingNumber')}
                          </p>
                          <p className="font-medium">{order.tracking_number}</p>
                        </div>
                      )}
                    </div>

                    {/* Действия */}
                    <div className="flex flex-wrap gap-2">
                      <Link
                        href={`/${locale}/profile/orders/${order.id}`}
                        className="btn btn-sm btn-ghost"
                      >
                        {t('viewDetails')}
                      </Link>

                      {marketplaceOrdersService.canConfirmDelivery(
                        order,
                        order.buyer_id
                      ) && (
                        <button
                          onClick={() => handleConfirmDelivery(order.id)}
                          disabled={processingOrderId === order.id}
                          className="btn btn-sm btn-success"
                        >
                          {processingOrderId === order.id ? (
                            <span className="loading loading-spinner loading-xs"></span>
                          ) : (
                            '✅ ' + t('confirmDelivery')
                          )}
                        </button>
                      )}

                      {marketplaceOrdersService.canOpenDispute(
                        order,
                        order.buyer_id
                      ) && (
                        <button
                          onClick={() => handleOpenDispute(order.id)}
                          disabled={processingOrderId === order.id}
                          className="btn btn-sm btn-error"
                        >
                          {processingOrderId === order.id ? (
                            <span className="loading loading-spinner loading-xs"></span>
                          ) : (
                            '⚠️ ' + t('openDispute')
                          )}
                        </button>
                      )}

                      {order.listing && (
                        <Link
                          href={`/${locale}/marketplace/${order.listing_id}`}
                          className="btn btn-sm btn-ghost"
                        >
                          {t('viewListing')}
                        </Link>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Пагинация */}
      {totalPages > 1 && (
        <div className="join mt-8 flex justify-center">
          <button
            className="join-item btn"
            onClick={() => setPage(Math.max(1, page - 1))}
            disabled={page === 1}
          >
            «
          </button>
          <button className="join-item btn btn-active">
            {tCommon('page')} {page} {tCommon('of')} {totalPages}
          </button>
          <button
            className="join-item btn"
            onClick={() => setPage(Math.min(totalPages, page + 1))}
            disabled={page === totalPages}
          >
            »
          </button>
        </div>
      )}
    </div>
  );
}
