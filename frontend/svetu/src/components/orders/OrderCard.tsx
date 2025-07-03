'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { ordersService } from '@/services/orders';
import SafeImage from '@/components/SafeImage';
import type { components } from '@/types/generated/api';

type StorefrontOrder = components['schemas']['models.StorefrontOrder'];
type StorefrontOrderItem = components['schemas']['models.StorefrontOrderItem'];

interface OrderCardProps {
  order: StorefrontOrder;
  onOrderUpdate?: (order: StorefrontOrder) => void;
}

export default function OrderCard({ order, onOrderUpdate }: OrderCardProps) {
  const t = useTranslations('orders');
  const [loading, setLoading] = useState(false);
  const [showDetails, setShowDetails] = useState(false);

  const formatDate = (dateString?: string) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'confirmed':
        return 'badge-success';
      case 'shipped':
        return 'badge-info';
      case 'delivered':
        return 'badge-success';
      case 'cancelled':
        return 'badge-error';
      case 'pending':
      default:
        return 'badge-warning';
    }
  };

  const getPaymentStatusColor = (status?: string) => {
    switch (status) {
      case 'paid':
      case 'completed':
        return 'badge-success';
      case 'failed':
      case 'cancelled':
        return 'badge-error';
      case 'pending':
      default:
        return 'badge-warning';
    }
  };

  const canCancelOrder = () => {
    return order.status === 'pending' || order.status === 'confirmed';
  };

  const handleCancelOrder = async () => {
    if (!order.id || !canCancelOrder()) return;

    const confirmed = window.confirm(t('confirmCancel'));
    if (!confirmed) return;

    try {
      setLoading(true);
      const reason = prompt(t('cancelReason'));
      const updatedOrder = await ordersService.cancelOrder(
        order.id,
        reason || undefined
      );
      onOrderUpdate?.(updatedOrder);
    } catch (error) {
      console.error('Failed to cancel order:', error);
      alert(t('cancelError'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card bg-base-100 shadow-lg">
      <div className="card-body">
        {/* Header */}
        <div className="flex justify-between items-start mb-4">
          <div>
            <h3 className="text-lg font-semibold">
              {t('orderNumber', { number: order.order_number })}
            </h3>
            <p className="text-sm text-base-content/70">
              {t('ordered')} {formatDate(order.created_at)}
            </p>
            {order.storefront && (
              <p className="text-sm text-base-content/70">
                {t('from')}{' '}
                <span className="font-medium">{order.storefront.name}</span>
              </p>
            )}
          </div>

          <div className="text-right">
            <div className={`badge ${getStatusColor(order.status)}`}>
              {t(`status.${order.status || 'pending'}`)}
            </div>
            {order.payment_status && (
              <div
                className={`badge ${getPaymentStatusColor(order.payment_status)} mt-1`}
              >
                {t(`paymentStatus.${order.payment_status}`)}
              </div>
            )}
          </div>
        </div>

        {/* Items Preview */}
        <div className="space-y-2">
          {order.items?.slice(0, 2).map((item: StorefrontOrderItem) => (
            <div
              key={item.id}
              className="flex items-center gap-3 p-2 bg-base-200 rounded"
            >
              {item.product?.images?.[0] && (
                <div className="w-12 h-12 rounded overflow-hidden flex-shrink-0">
                  <SafeImage
                    src={
                      item.product.images[0].thumbnail_url ||
                      item.product.images[0].image_url ||
                      ''
                    }
                    alt={item.product.name || 'Product'}
                    width={48}
                    height={48}
                    className="object-cover"
                  />
                </div>
              )}

              <div className="flex-1 min-w-0">
                <p className="font-medium truncate">{item.product?.name}</p>
                <p className="text-sm text-base-content/70">
                  {t('quantity')}: {item.quantity} × {item.price_per_unit}{' '}
                  {order.currency}
                </p>
              </div>

              <div className="text-right">
                <p className="font-medium">
                  {item.total_price} {order.currency}
                </p>
              </div>
            </div>
          ))}

          {order.items && order.items.length > 2 && (
            <button
              onClick={() => setShowDetails(!showDetails)}
              className="text-sm text-primary hover:underline"
            >
              {showDetails
                ? t('hideDetails')
                : t('showMore', { count: order.items.length - 2 })}
            </button>
          )}
        </div>

        {/* Expanded Details */}
        {showDetails && order.items && order.items.length > 2 && (
          <div className="space-y-2 mt-2">
            {order.items.slice(2).map((item: StorefrontOrderItem) => (
              <div
                key={item.id}
                className="flex items-center gap-3 p-2 bg-base-200 rounded"
              >
                {item.product?.images?.[0] && (
                  <div className="w-12 h-12 rounded overflow-hidden flex-shrink-0">
                    <SafeImage
                      src={
                        item.product.images[0].thumbnail_url ||
                        item.product.images[0].image_url ||
                        ''
                      }
                      alt={item.product.name || 'Product'}
                      width={48}
                      height={48}
                      className="object-cover"
                    />
                  </div>
                )}

                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{item.product?.name}</p>
                  <p className="text-sm text-base-content/70">
                    {t('quantity')}: {item.quantity} × {item.price_per_unit}{' '}
                    {order.currency}
                  </p>
                </div>

                <div className="text-right">
                  <p className="font-medium">
                    {item.total_price} {order.currency}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Total */}
        <div className="divider"></div>
        <div className="flex justify-between items-center">
          <div>
            <p className="text-sm text-base-content/70">
              {t('subtotal')}: {order.subtotal} {order.currency}
            </p>
            {order.shipping_amount && order.shipping_amount > 0 && (
              <p className="text-sm text-base-content/70">
                {t('shipping')}: {order.shipping_amount} {order.currency}
              </p>
            )}
            {order.tax_amount && order.tax_amount > 0 && (
              <p className="text-sm text-base-content/70">
                {t('tax')}: {order.tax_amount} {order.currency}
              </p>
            )}
          </div>

          <div className="text-right">
            <p className="text-lg font-bold">
              {t('total')}: {order.total_amount} {order.currency}
            </p>
          </div>
        </div>

        {/* Actions */}
        {canCancelOrder() && (
          <div className="card-actions justify-end mt-4">
            <button
              onClick={handleCancelOrder}
              className="btn btn-outline btn-sm btn-error"
              disabled={loading}
            >
              {loading && (
                <span className="loading loading-spinner loading-sm mr-2" />
              )}
              {t('cancel')}
            </button>
          </div>
        )}

        {/* Tracking Info */}
        {order.tracking_number && (
          <div className="mt-4 p-3 bg-info/10 rounded">
            <p className="text-sm">
              <span className="font-medium">{t('trackingNumber')}:</span>{' '}
              {order.tracking_number}
            </p>
            {order.shipping_provider && (
              <p className="text-sm">
                <span className="font-medium">{t('carrier')}:</span>{' '}
                {order.shipping_provider}
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
