'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';

type B2COrder = components['schemas']['models.B2COrder'];

interface OrderDetailsProps {
  order: B2COrder;
}

export default function OrderDetails({ order }: OrderDetailsProps) {
  const t = useTranslations('orders');
  const tCommon = useTranslations('common');

  const formatAddress = (address: any) => {
    if (!address) return null;

    const parts = [];
    if (address.street) parts.push(address.street);
    if (address.city) parts.push(address.city);
    if (address.postal_code) parts.push(address.postal_code);
    if (address.country) parts.push(address.country);

    return parts.join(', ');
  };

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

  return (
    <div className="space-y-6">
      {/* Order Info */}
      <div className="card bg-base-200">
        <div className="card-body">
          <h3 className="card-title">{t('orderInfo')}</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-base-content/70">{t('orderNumber')}</p>
              <p className="font-medium">{order.order_number}</p>
            </div>
            <div>
              <p className="text-sm text-base-content/70">{t('orderDate')}</p>
              <p className="font-medium">{formatDate(order.created_at)}</p>
            </div>
            <div>
              <p className="text-sm text-base-content/70">
                {t('status.label')}
              </p>
              <p className="font-medium">{t(`status.${order.status}`)}</p>
            </div>
            <div>
              <p className="text-sm text-base-content/70">
                {t('paymentMethod')}
              </p>
              <p className="font-medium">
                {t(`payment.${order.payment_method}`)}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Addresses */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Pickup Address */}
        {(order as any).pickup_address && (
          <div className="card bg-base-200">
            <div className="card-body">
              <h3 className="card-title">{t('pickupAddress.title')}</h3>
              <p className="text-sm text-base-content/70 mb-2">
                {t('pickupAddress.description')}
              </p>
              <div className="space-y-2">
                {(order as any).pickup_address.name && (
                  <p className="font-medium">
                    {(order as any).pickup_address.name}
                  </p>
                )}
                {(order as any).pickup_address.formatted ? (
                  <p>{(order as any).pickup_address.formatted}</p>
                ) : (
                  <p>{formatAddress((order as any).pickup_address)}</p>
                )}
                {(order as any).pickup_address.phone && (
                  <p className="text-sm">
                    <span className="text-base-content/70">
                      {tCommon('phone')}:
                    </span>{' '}
                    {(order as any).pickup_address.phone}
                  </p>
                )}
                {(order as any).pickup_address.email && (
                  <p className="text-sm">
                    <span className="text-base-content/70">
                      {tCommon('email')}:
                    </span>{' '}
                    {(order as any).pickup_address.email}
                  </p>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Delivery Address */}
        {order.shipping_address && (
          <div className="card bg-base-200">
            <div className="card-body">
              <h3 className="card-title">{t('deliveryAddress.title')}</h3>
              <p className="text-sm text-base-content/70 mb-2">
                {t('deliveryAddress.description')}
              </p>
              <div className="space-y-2">
                <p>{formatAddress(order.shipping_address)}</p>
                {order.shipping_method && (
                  <p className="text-sm">
                    <span className="text-base-content/70">
                      {t('shippingMethod')}:
                    </span>{' '}
                    {t(`shipping.${order.shipping_method}`)}
                  </p>
                )}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Tracking Info */}
      {order.tracking_number && (
        <div className="card bg-base-200">
          <div className="card-body">
            <h3 className="card-title">{t('tracking.title')}</h3>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-base-content/70">
                  {t('tracking.number')}
                </p>
                <p className="font-medium">{order.tracking_number}</p>
              </div>
              {order.shipping_provider && (
                <div>
                  <p className="text-sm text-base-content/70">
                    {t('tracking.provider')}
                  </p>
                  <p className="font-medium">{order.shipping_provider}</p>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Order Items */}
      <div className="card bg-base-200">
        <div className="card-body">
          <h3 className="card-title">{t('items')}</h3>
          <div className="space-y-3">
            {order.items?.map((item) => (
              <div
                key={item.id}
                className="flex items-center gap-4 p-3 bg-base-100 rounded"
              >
                <div className="flex-1">
                  <p className="font-medium">{item.product_name}</p>
                  {item.variant_name && (
                    <p className="text-sm text-base-content/70">
                      {item.variant_name}
                    </p>
                  )}
                  <p className="text-sm">
                    {item.quantity} × {item.price_per_unit} {order.currency}
                  </p>
                </div>
                <div className="text-right">
                  <p className="font-semibold">
                    {item.total_price} {order.currency}
                  </p>
                </div>
              </div>
            ))}
          </div>

          {/* Totals */}
          <div className="divider"></div>
          <div className="space-y-2">
            <div className="flex justify-between">
              <span>{t('subtotal')}</span>
              <span>
                {order.subtotal} {order.currency}
              </span>
            </div>
            {(order.shipping ?? 0) > 0 && (
              <div className="flex justify-between">
                <span>{t('shipping')}</span>
                <span>
                  {order.shipping} {order.currency}
                </span>
              </div>
            )}
            {(order.tax ?? 0) > 0 && (
              <div className="flex justify-between">
                <span>{t('tax')}</span>
                <span>
                  {order.tax} {order.currency}
                </span>
              </div>
            )}
            {(order.discount ?? 0) > 0 && (
              <div className="flex justify-between text-success">
                <span>{t('discount')}</span>
                <span>
                  -{order.discount} {order.currency}
                </span>
              </div>
            )}
            <div className="divider my-2"></div>
            <div className="flex justify-between text-lg font-bold">
              <span>{t('total')}</span>
              <span className="text-primary">
                {order.total} {order.currency}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
