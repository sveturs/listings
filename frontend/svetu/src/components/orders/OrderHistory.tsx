'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { ordersService } from '@/services/orders';
import { useAuthContext } from '@/contexts/AuthContext';
import OrderCard from './OrderCard';
import type { components } from '@/types/generated/api';

type StorefrontOrder = components['schemas']['models.StorefrontOrder'];

interface OrderHistoryProps {
  status?: string;
  limit?: number;
}

export default function OrderHistory({
  status,
  limit = 10,
}: OrderHistoryProps) {
  const t = useTranslations('orders');
  const { isAuthenticated } = useAuthContext();

  const [orders, setOrders] = useState<StorefrontOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [hasMore, setHasMore] = useState(false);
  const [total, setTotal] = useState(0);

  const fetchOrders = async (page = 1) => {
    if (!isAuthenticated) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const response = await ordersService.getUserOrders({
        status,
        limit,
        offset: (page - 1) * limit,
      });

      if (page === 1) {
        setOrders(response.orders);
      } else {
        setOrders((prev) => [...prev, ...response.orders]);
      }

      setTotal(response.total);
      setHasMore(
        response.orders.length === limit &&
          orders.length + response.orders.length < response.total
      );
    } catch (err) {
      console.error('Failed to fetch orders:', err);
      setError(err instanceof Error ? err.message : 'Failed to fetch orders');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders(1);
  }, [isAuthenticated, status]); // eslint-disable-line react-hooks/exhaustive-deps

  const handleLoadMore = () => {
    if (hasMore && !loading) {
      const nextPage = currentPage + 1;
      setCurrentPage(nextPage);
      fetchOrders(nextPage);
    }
  };

  if (!isAuthenticated) {
    return (
      <div className="text-center py-8">
        <p className="text-base-content/70">{t('loginRequired')}</p>
      </div>
    );
  }

  if (loading && currentPage === 1) {
    return (
      <div className="space-y-4">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="card bg-base-100 shadow animate-pulse">
            <div className="card-body">
              <div className="flex justify-between items-start">
                <div className="space-y-2 flex-1">
                  <div className="h-4 bg-base-content/10 rounded w-1/4"></div>
                  <div className="h-6 bg-base-content/10 rounded w-1/2"></div>
                  <div className="h-4 bg-base-content/10 rounded w-1/3"></div>
                </div>
                <div className="h-8 bg-base-content/10 rounded w-20"></div>
              </div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <span>{error}</span>
        <button onClick={() => fetchOrders(1)} className="btn btn-sm">
          {t('tryAgain')}
        </button>
      </div>
    );
  }

  if (orders.length === 0) {
    return (
      <div className="text-center py-12">
        <svg
          className="w-16 h-16 mx-auto mb-4 text-base-content/20"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
          />
        </svg>
        <h3 className="text-lg font-medium mb-2">{t('noOrders')}</h3>
        <p className="text-base-content/70">{t('noOrdersDescription')}</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">{t('history')}</h2>
        {total > 0 && (
          <span className="text-sm text-base-content/70">
            {t('totalOrders', { count: total })}
          </span>
        )}
      </div>

      {/* Orders List */}
      <div className="space-y-4">
        {orders.map((order) => (
          <OrderCard key={order.id} order={order} />
        ))}
      </div>

      {/* Load More */}
      {hasMore && (
        <div className="text-center py-4">
          <button
            onClick={handleLoadMore}
            className="btn btn-outline"
            disabled={loading}
          >
            {loading && (
              <span className="loading loading-spinner loading-sm mr-2" />
            )}
            {t('loadMore')}
          </button>
        </div>
      )}
    </div>
  );
}
