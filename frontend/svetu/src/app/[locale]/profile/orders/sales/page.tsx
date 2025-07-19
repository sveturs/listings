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

export default function MySalesPage() {
  const locale = useLocale();
  const t = useTranslations();
  const { isAuthenticated } = useAuth();

  const [orders, setOrders] = useState<MarketplaceOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [processingOrderId, setProcessingOrderId] = useState<number | null>(
    null
  );
  const [shippingModal, setShippingModal] = useState<{
    orderId: number | null;
    shippingMethod: string;
    trackingNumber: string;
  }>({
    orderId: null,
    shippingMethod: '',
    trackingNumber: '',
  });

  const limit = 20;

  useEffect(() => {
    if (isAuthenticated) {
      fetchOrders();
    }
  }, [isAuthenticated, page]); // eslint-disable-line react-hooks/exhaustive-deps

  const fetchOrders = async () => {
    try {
      setLoading(true);
      const response = await marketplaceOrdersService.getMySales(page, limit);
      setOrders(response.orders);
      setTotal(response.total);
    } catch (error) {
      console.error('Error fetching sales:', error);
      toast.error(t('orders.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const handleMarkAsShipped = async () => {
    if (
      !shippingModal.orderId ||
      !shippingModal.shippingMethod ||
      !shippingModal.trackingNumber
    ) {
      toast.error(t('orders.fillAllShippingFields'));
      return;
    }

    setProcessingOrderId(shippingModal.orderId);
    try {
      await marketplaceOrdersService.markAsShipped(
        shippingModal.orderId,
        shippingModal.shippingMethod,
        shippingModal.trackingNumber
      );
      toast.success(t('orders.markedAsShipped'));
      setShippingModal({
        orderId: null,
        shippingMethod: '',
        trackingNumber: '',
      });
      await fetchOrders(); // Обновляем список
    } catch (error) {
      console.error('Error marking as shipped:', error);
      toast.error(t('orders.shippingError'));
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
              <Link href={`/${locale}`}>{t('home.title')}</Link>
            </li>
            <li>
              <Link href={`/${locale}/profile`}>{t('profile.title')}</Link>
            </li>
            <li>{t('orders.mySales')}</li>
          </ul>
        </div>

        <h1 className="text-3xl font-bold mb-4">{t('orders.mySales')}</h1>

        {/* Табы */}
        <div className="tabs tabs-boxed mb-6">
          <Link href={`/${locale}/profile/orders/purchases`} className="tab">
            {t('orders.purchases')}
          </Link>
          <Link
            href={`/${locale}/profile/orders/sales`}
            className="tab tab-active"
          >
            {t('orders.sales')}
          </Link>
        </div>

        {/* Статистика */}
        <div className="stats shadow mb-6">
          <div className="stat">
            <div className="stat-title">{t('orders.totalSales')}</div>
            <div className="stat-value">{total}</div>
          </div>
          <div className="stat">
            <div className="stat-title">{t('orders.pendingShipment')}</div>
            <div className="stat-value text-warning">
              {orders.filter((o) => o.status === 'paid').length}
            </div>
          </div>
          <div className="stat">
            <div className="stat-title">{t('orders.inDispute')}</div>
            <div className="stat-value text-error">
              {orders.filter((o) => o.status === 'disputed').length}
            </div>
          </div>
        </div>
      </div>

      {/* Список заказов */}
      {orders.length === 0 ? (
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body text-center py-16">
            <h3 className="text-xl font-semibold mb-2">
              {t('orders.noSales')}
            </h3>
            <p className="text-base-content/70 mb-6">
              {t('orders.noSalesDesc')}
            </p>
            <Link
              href={`/${locale}/profile/listings`}
              className="btn btn-primary"
            >
              {t('orders.createListing')}
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
                          {order.listing?.title || t('orders.unknownItem')}
                        </h3>
                        <p className="text-sm text-base-content/70">
                          {t('orders.orderNumber')}: #{order.id}
                        </p>
                        <p className="text-sm text-base-content/70">
                          {t('orders.orderDate')}:{' '}
                          {new Date(order.created_at).toLocaleDateString(
                            locale
                          )}
                        </p>
                      </div>

                      {/* Статус и выплата */}
                      <div className="text-right">
                        <div
                          className={`badge ${marketplaceOrdersService.getStatusColor(order.status)} badge-lg mb-2`}
                        >
                          {marketplaceOrdersService.getStatusLabel(
                            order.status,
                            locale
                          )}
                        </div>
                        <p className="text-sm">
                          {t('orders.yourEarnings')}:{' '}
                          <span className="font-bold text-success">
                            {balanceService.formatAmount(
                              order.seller_payout_amount,
                              'RSD'
                            )}
                          </span>
                        </p>
                        <p className="text-xs text-base-content/70">
                          {t('orders.afterFee', {
                            fee: `${order.platform_fee_rate * 100}%`,
                          })}
                        </p>
                      </div>
                    </div>

                    {/* Детали */}
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                      <div>
                        <p className="text-sm text-base-content/70">
                          {t('orders.buyer')}
                        </p>
                        <p className="font-medium">
                          {order.buyer?.name || 'Unknown'}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-base-content/70">
                          {t('orders.price')}
                        </p>
                        <p className="font-medium">
                          {balanceService.formatAmount(order.item_price, 'RSD')}
                        </p>
                      </div>
                      {order.shipped_at && (
                        <div>
                          <p className="text-sm text-base-content/70">
                            {t('orders.shippedDate')}
                          </p>
                          <p className="font-medium">
                            {new Date(order.shipped_at).toLocaleDateString(
                              locale
                            )}
                          </p>
                        </div>
                      )}
                    </div>

                    {/* Предупреждение о необходимости отправки */}
                    {order.status === 'paid' && (
                      <div className="alert alert-warning mb-4">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="stroke-current shrink-0 h-6 w-6"
                          fill="none"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                          />
                        </svg>
                        <span>{t('orders.shipmentRequired')}</span>
                      </div>
                    )}

                    {/* Действия */}
                    <div className="flex flex-wrap gap-2">
                      <Link
                        href={`/${locale}/profile/orders/${order.id}`}
                        className="btn btn-sm btn-ghost"
                      >
                        {t('orders.viewDetails')}
                      </Link>

                      {marketplaceOrdersService.canShip(
                        order,
                        order.seller_id
                      ) && (
                        <button
                          onClick={() =>
                            setShippingModal({
                              orderId: order.id,
                              shippingMethod: '',
                              trackingNumber: '',
                            })
                          }
                          className="btn btn-sm btn-primary"
                        >
                          📦 {t('orders.markAsShipped')}
                        </button>
                      )}

                      {order.listing && (
                        <Link
                          href={`/${locale}/marketplace/${order.listing_id}`}
                          className="btn btn-sm btn-ghost"
                        >
                          {t('orders.viewListing')}
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
            {t('common.page')} {page} {t('common.of')} {totalPages}
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

      {/* Модальное окно для отправки */}
      {shippingModal.orderId && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg mb-4">
              {t('orders.shippingDetails')}
            </h3>

            <div className="form-control mb-4">
              <label className="label">
                <span className="label-text">{t('orders.shippingMethod')}</span>
              </label>
              <select
                className="select select-bordered w-full"
                value={shippingModal.shippingMethod}
                onChange={(e) =>
                  setShippingModal({
                    ...shippingModal,
                    shippingMethod: e.target.value,
                  })
                }
              >
                <option value="">{t('orders.selectShippingMethod')}</option>
                <option value="post">{t('orders.shippingPost')}</option>
                <option value="courier">{t('orders.shippingCourier')}</option>
                <option value="personal">{t('orders.shippingPersonal')}</option>
                <option value="other">{t('orders.shippingOther')}</option>
              </select>
            </div>

            <div className="form-control mb-6">
              <label className="label">
                <span className="label-text">{t('orders.trackingNumber')}</span>
              </label>
              <input
                type="text"
                className="input input-bordered w-full"
                placeholder={t('orders.trackingNumberPlaceholder')}
                value={shippingModal.trackingNumber}
                onChange={(e) =>
                  setShippingModal({
                    ...shippingModal,
                    trackingNumber: e.target.value,
                  })
                }
              />
            </div>

            <div className="modal-action">
              <button
                className="btn btn-ghost"
                onClick={() =>
                  setShippingModal({
                    orderId: null,
                    shippingMethod: '',
                    trackingNumber: '',
                  })
                }
                disabled={processingOrderId === shippingModal.orderId}
              >
                {t('common.cancel')}
              </button>
              <button
                className="btn btn-primary"
                onClick={handleMarkAsShipped}
                disabled={
                  processingOrderId === shippingModal.orderId ||
                  !shippingModal.shippingMethod ||
                  !shippingModal.trackingNumber
                }
              >
                {processingOrderId === shippingModal.orderId ? (
                  <span className="loading loading-spinner loading-sm"></span>
                ) : (
                  t('orders.confirmShipment')
                )}
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button
              onClick={() =>
                setShippingModal({
                  orderId: null,
                  shippingMethod: '',
                  trackingNumber: '',
                })
              }
            >
              close
            </button>
          </form>
        </dialog>
      )}
    </div>
  );
}
