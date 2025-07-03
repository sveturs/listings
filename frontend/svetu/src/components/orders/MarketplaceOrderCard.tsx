'use client';

import React from 'react';
import Link from 'next/link';
import { useLocale, useTranslations } from 'next-intl';
import { balanceService } from '@/services/balance';

interface MarketplaceOrder {
  id: number;
  buyer_id: number;
  seller_id: number;
  listing_id: number;
  item_price: number;
  status: string;
  created_at: string;
  updated_at: string;
  listing?: {
    id: number;
    title: string;
    description?: string;
    price: number;
    images?: Array<{
      id: number;
      url: string;
    }>;
  };
  seller?: {
    id: number;
    name?: string;
    email: string;
  };
  buyer?: {
    id: number;
    name?: string;
    email: string;
  };
}

interface MarketplaceOrderCardProps {
  order: MarketplaceOrder;
}

export default function MarketplaceOrderCard({
  order,
}: MarketplaceOrderCardProps) {
  const _t = useTranslations('orders');
  const locale = useLocale();

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'pending':
        return 'badge-warning';
      case 'paid':
        return 'badge-info';
      case 'shipped':
        return 'badge-primary';
      case 'delivered':
        return 'badge-success';
      case 'completed':
        return 'badge-success';
      case 'cancelled':
        return 'badge-error';
      case 'disputed':
        return 'badge-error';
      case 'refunded':
        return 'badge-ghost';
      default:
        return 'badge-ghost';
    }
  };

  const getStatusLabel = (status: string) => {
    const statusMap: Record<string, string> = {
      pending: 'В ожидании',
      paid: 'Оплачен',
      shipped: 'Отправлен',
      delivered: 'Доставлен',
      completed: 'Завершен',
      cancelled: 'Отменен',
      disputed: 'Спор',
      refunded: 'Возврат',
    };
    return statusMap[status] || status;
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat(locale, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
  };

  return (
    <div className="card bg-base-100 shadow hover:shadow-lg transition-shadow">
      <div className="card-body">
        <div className="flex flex-col sm:flex-row justify-between gap-4">
          {/* Order Info */}
          <div className="flex-1">
            <div className="flex items-start gap-4">
              {/* Product Image */}
              {order.listing?.images?.[0] && (
                <div className="avatar">
                  <div className="w-20 h-20 rounded">
                    <img
                      src={order.listing.images[0].url}
                      alt={order.listing.title}
                      className="object-cover"
                    />
                  </div>
                </div>
              )}

              {/* Order Details */}
              <div className="flex-1">
                <h3 className="font-semibold text-lg">
                  {order.listing?.title || `Заказ #${order.id}`}
                </h3>

                <div className="text-sm text-base-content/70 mt-1 space-y-1">
                  <p>Заказ #{order.id}</p>
                  <p>
                    Продавец:{' '}
                    {order.seller?.name || order.seller?.email || 'Неизвестно'}
                  </p>
                  <p>{formatDate(order.created_at)}</p>
                </div>

                {/* Price */}
                <div className="font-semibold text-lg mt-2">
                  {balanceService.formatAmount(order.item_price, 'RSD')}
                </div>
              </div>
            </div>
          </div>

          {/* Status and Actions */}
          <div className="flex flex-col items-end gap-2">
            <div className={`badge ${getStatusBadgeClass(order.status)}`}>
              {getStatusLabel(order.status)}
            </div>

            <Link
              href={`/${locale}/orders/${order.id}`}
              className="btn btn-sm btn-primary"
            >
              Подробнее
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
