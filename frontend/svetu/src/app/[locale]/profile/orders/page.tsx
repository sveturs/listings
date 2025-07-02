'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import OrderHistory from '@/components/orders/OrderHistory';

export default function OrdersPage() {
  const t = useTranslations('orders');
  const [selectedStatus, setSelectedStatus] = useState<string>('');

  const statusOptions = [
    { value: '', label: t('allOrders') },
    { value: 'pending', label: t('status.pending') },
    { value: 'confirmed', label: t('status.confirmed') },
    { value: 'shipped', label: t('status.shipped') },
    { value: 'delivered', label: t('status.delivered') },
    { value: 'cancelled', label: t('status.cancelled') },
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
          <div>
            <h1 className="text-3xl font-bold">{t('myOrders')}</h1>
            <p className="text-base-content/70 mt-1">
              {t('ordersDescription')}
            </p>
          </div>

          {/* Status Filter */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">{t('filterByStatus')}</span>
            </label>
            <select
              className="select select-bordered w-full max-w-xs"
              value={selectedStatus}
              onChange={(e) => setSelectedStatus(e.target.value)}
            >
              {statusOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Orders List */}
        <OrderHistory status={selectedStatus || undefined} limit={10} />
      </div>
    </div>
  );
}
