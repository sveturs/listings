'use client';

import { use, useState, useEffect } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import {
  marketplaceOrdersService,
  type MarketplaceOrder,
} from '@/services/marketplaceOrders';
import { balanceService } from '@/services/balance';
import SafeImage from '@/components/SafeImage';
import { toast } from 'react-hot-toast';

interface Props {
  params: Promise<{ id: string }>;
}

export default function OrderDetailsPage({ params }: Props) {
  const { id } = use(params);
  const locale = useLocale();
  const t = useTranslations();
  const router = useRouter();
  const { user, isAuthenticated } = useAuth();

  const [order, setOrder] = useState<MarketplaceOrder | null>(null);
  const [loading, setLoading] = useState(true);
  const [messageText, setMessageText] = useState('');
  const [sendingMessage, setSendingMessage] = useState(false);

  const orderId = parseInt(id);

  useEffect(() => {
    if (isAuthenticated && orderId) {
      fetchOrderDetails();
    }
  }, [isAuthenticated, orderId]);

  const fetchOrderDetails = async () => {
    try {
      setLoading(true);
      const orderData = await marketplaceOrdersService.getOrderDetails(orderId);
      setOrder(orderData);
    } catch (error) {
      console.error('Error fetching order details:', error);
      toast.error(t('orders.fetchError'));
      router.push(`/${locale}/profile/orders/purchases`);
    } finally {
      setLoading(false);
    }
  };

  const handleSendMessage = async () => {
    if (!messageText.trim()) return;

    setSendingMessage(true);
    try {
      await marketplaceOrdersService.addMessage(orderId, messageText);
      toast.success(t('orders.messageSent'));
      setMessageText('');
      await fetchOrderDetails(); // Обновляем данные заказа
    } catch (error) {
      console.error('Error sending message:', error);
      toast.error(t('orders.messageError'));
    } finally {
      setSendingMessage(false);
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center min-h-[400px]">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  if (!order) {
    return null;
  }

  const isBuyer = user?.id === order.buyer_id;
  const isSeller = user?.id === order.seller_id;

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      {/* Заголовок и навигация */}
      <div className="breadcrumbs text-sm mb-6">
        <ul>
          <li>
            <Link href={`/${locale}`}>{t('home.title')}</Link>
          </li>
          <li>
            <Link href={`/${locale}/profile`}>{t('profile.title')}</Link>
          </li>
          <li>
            <Link
              href={`/${locale}/profile/orders/${isBuyer ? 'purchases' : 'sales'}`}
            >
              {isBuyer ? t('orders.myPurchases') : t('orders.mySales')}
            </Link>
          </li>
          <li>
            {t('orders.orderNumber')}: #{order.id}
          </li>
        </ul>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Левая колонка - основная информация */}
        <div className="lg:col-span-2 space-y-6">
          {/* Карточка заказа */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <div className="flex justify-between items-start mb-6">
                <h1 className="card-title text-2xl">
                  {t('orders.orderNumber')}: #{order.id}
                </h1>
                <div
                  className={`badge ${marketplaceOrdersService.getStatusColor(order.status)} badge-lg`}
                >
                  {marketplaceOrdersService.getStatusLabel(
                    order.status,
                    locale
                  )}
                </div>
              </div>

              {/* Товар */}
              <div className="divider">{t('orders.itemDetails')}</div>
              <div className="flex gap-4 mb-6">
                <figure className="relative w-32 h-32 flex-shrink-0 bg-base-200 rounded-lg overflow-hidden">
                  <SafeImage
                    src={order.listing?.images?.[0]?.public_url}
                    alt={order.listing?.title || 'Order item'}
                    fill
                    className="object-cover"
                  />
                </figure>
                <div className="flex-1">
                  <h3 className="text-xl font-semibold mb-2">
                    {order.listing?.title || t('orders.unknownItem')}
                  </h3>
                  {order.listing?.description && (
                    <p className="text-base-content/70 line-clamp-2">
                      {order.listing.description}
                    </p>
                  )}
                  <Link
                    href={`/${locale}/marketplace/${order.listing_id}`}
                    className="link link-primary text-sm mt-2 inline-block"
                  >
                    {t('orders.viewListing')} →
                  </Link>
                </div>
              </div>

              {/* Участники сделки */}
              <div className="divider">{t('orders.participants')}</div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                <div>
                  <p className="text-sm text-base-content/70 mb-1">
                    {t('orders.buyer')}
                  </p>
                  <p className="font-medium">
                    {order.buyer?.name || 'Unknown'}
                    {isBuyer && ' (' + t('common.you') + ')'}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-base-content/70 mb-1">
                    {t('orders.seller')}
                  </p>
                  <p className="font-medium">
                    {order.seller?.name || 'Unknown'}
                    {isSeller && ' (' + t('common.you') + ')'}
                  </p>
                </div>
              </div>

              {/* Доставка */}
              {(order.shipping_method || order.tracking_number) && (
                <>
                  <div className="divider">{t('orders.shippingInfo')}</div>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                    {order.shipping_method && (
                      <div>
                        <p className="text-sm text-base-content/70 mb-1">
                          {t('orders.shippingMethod')}
                        </p>
                        <p className="font-medium">{order.shipping_method}</p>
                      </div>
                    )}
                    {order.tracking_number && (
                      <div>
                        <p className="text-sm text-base-content/70 mb-1">
                          {t('orders.trackingNumber')}
                        </p>
                        <p className="font-medium">{order.tracking_number}</p>
                      </div>
                    )}
                  </div>
                </>
              )}

              {/* История статусов */}
              {order.status_history && order.status_history.length > 0 && (
                <>
                  <div className="divider">{t('orders.statusHistory')}</div>
                  <div className="space-y-2">
                    {order.status_history.map((history) => (
                      <div
                        key={history.id}
                        className="flex items-center gap-3 text-sm"
                      >
                        <span className="text-base-content/50">
                          {new Date(history.created_at).toLocaleString(locale)}
                        </span>
                        <span>→</span>
                        <span className="font-medium">
                          {marketplaceOrdersService.getStatusLabel(
                            history.new_status as any,
                            locale
                          )}
                        </span>
                        {history.reason && (
                          <span className="text-base-content/70">
                            ({history.reason})
                          </span>
                        )}
                      </div>
                    ))}
                  </div>
                </>
              )}
            </div>
          </div>

          {/* Сообщения */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title mb-4">{t('orders.messages')}</h2>

              {order.messages && order.messages.length > 0 ? (
                <div className="space-y-3 mb-4 max-h-96 overflow-y-auto">
                  {order.messages.map((message) => (
                    <div
                      key={message.id}
                      className={`chat ${message.sender_id === user?.id ? 'chat-end' : 'chat-start'}`}
                    >
                      <div className="chat-header">
                        {message.sender?.name || 'Unknown'}
                        <time className="text-xs opacity-50 ml-2">
                          {new Date(message.created_at).toLocaleString(locale)}
                        </time>
                      </div>
                      <div className="chat-bubble">{message.content}</div>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-base-content/70 text-center py-8">
                  {t('orders.noMessages')}
                </p>
              )}

              {/* Форма отправки сообщения */}
              <div className="form-control">
                <div className="input-group">
                  <input
                    type="text"
                    placeholder={t('orders.typeMessage')}
                    className="input input-bordered flex-1"
                    value={messageText}
                    onChange={(e) => setMessageText(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
                    disabled={sendingMessage}
                  />
                  <button
                    className="btn btn-primary"
                    onClick={handleSendMessage}
                    disabled={sendingMessage || !messageText.trim()}
                  >
                    {sendingMessage ? (
                      <span className="loading loading-spinner loading-sm"></span>
                    ) : (
                      t('common.send')
                    )}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Правая колонка - финансы и действия */}
        <div className="space-y-6">
          {/* Финансовая информация */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h3 className="card-title text-lg mb-4">
                {t('orders.financialDetails')}
              </h3>

              <div className="space-y-3">
                <div className="flex justify-between">
                  <span>{t('orders.itemPrice')}</span>
                  <span className="font-medium">
                    {balanceService.formatAmount(order.item_price, 'RSD')}
                  </span>
                </div>

                <div className="flex justify-between text-sm text-base-content/70">
                  <span>
                    {t('orders.platformFee')} ({order.platform_fee_rate * 100}%)
                  </span>
                  <span>
                    -
                    {balanceService.formatAmount(
                      order.platform_fee_amount,
                      'RSD'
                    )}
                  </span>
                </div>

                <div className="divider my-2"></div>

                <div className="flex justify-between font-bold">
                  <span>
                    {isSeller
                      ? t('orders.yourEarnings')
                      : t('orders.sellerReceives')}
                  </span>
                  <span className="text-success">
                    {balanceService.formatAmount(
                      order.seller_payout_amount,
                      'RSD'
                    )}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Защита покупателя */}
          {order.protection_expires_at &&
            ['shipped', 'delivered'].includes(order.status) && (
              <div className="card bg-warning/10 border-warning shadow-xl">
                <div className="card-body">
                  <h3 className="card-title text-lg mb-2">
                    🛡️ {t('orders.buyerProtection')}
                  </h3>
                  <p className="text-sm mb-2">{t('orders.protectionActive')}</p>
                  <p className="font-bold">
                    ⏱️ {t('orders.timeRemaining')}:{' '}
                    {marketplaceOrdersService.getProtectionTimeLeft(order) ||
                      t('orders.expired')}
                  </p>
                </div>
              </div>
            )}

          {/* Действия */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h3 className="card-title text-lg mb-4">{t('orders.actions')}</h3>

              <div className="space-y-2">
                {isBuyer &&
                  marketplaceOrdersService.canConfirmDelivery(
                    order,
                    user!.id
                  ) && (
                    <button className="btn btn-success btn-block">
                      ✅ {t('orders.confirmDelivery')}
                    </button>
                  )}

                {isSeller &&
                  marketplaceOrdersService.canShip(order, user!.id) && (
                    <button className="btn btn-primary btn-block">
                      📦 {t('orders.markAsShipped')}
                    </button>
                  )}

                {marketplaceOrdersService.canOpenDispute(order, user!.id) && (
                  <button className="btn btn-error btn-outline btn-block">
                    ⚠️ {t('orders.openDispute')}
                  </button>
                )}

                <Link
                  href={`/${locale}/support`}
                  className="btn btn-ghost btn-block"
                >
                  {t('orders.contactSupport')}
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
