'use client';

import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useLocale, useTranslations } from 'next-intl';
import { balanceService } from '@/services/balance';

interface C2COrder {
  id: number;
  buyer_id: number;
  seller_id: number;
  listing_id: number;
  total: number | string;
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

interface C2COrderCardProps {
  order: C2COrder;
}

export default function C2COrderCard({
  order,
}: C2COrderCardProps) {
  const t = useTranslations('orders');
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
    return t(`status.${status}` as any) || status;
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
                    <Image
                      src={order.listing.images[0].url}
                      alt={order.listing.title}
                      fill
                      className="object-cover"
                    />
                  </div>
                </div>
              )}

              {/* Order Details */}
              <div className="flex-1">
                <h3 className="font-semibold text-lg">
                  {order.listing?.title || `${t('order')} #${order.id}`}
                </h3>

                <div className="text-sm text-base-content/70 mt-1 space-y-1">
                  <p>
                    {t('order')} #{order.id}
                  </p>
                  <p>
                    {t('seller')}:{' '}
                    {order.seller?.name || order.seller?.email || t('unknown')}
                  </p>
                  <p>{formatDate(order.created_at)}</p>
                </div>

                {/* Price */}
                <div className="font-semibold text-lg mt-2">
                  {balanceService.formatAmount(order.total, 'RSD')}
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
              {t('viewDetails')}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
