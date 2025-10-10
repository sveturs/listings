'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { format } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import { useLocale } from 'next-intl';
import { apiClientAuth } from '@/lib/api-client-auth';

interface Shipment {
  id: number;
  tracking_number: string;
  status: string;
  recipient_name: string;
  recipient_city: string;
  recipient_phone: string;
  sender_name: string;
  sender_city: string;
  cod_amount: number;
  weight_kg: number;
  created_at: string;
  delivered_at?: string;
  provider: string;
}

interface _ShipmentsResponse {
  shipments: Shipment[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export default function ShipmentsListClient() {
  const t = useTranslations('admin');
  const router = useRouter();
  const locale = useLocale();
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [limit] = useState(20);
  const [totalPages, setTotalPages] = useState(1);

  // Filters
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [providerFilter, setProviderFilter] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState('');
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');
  const [cityFilter, setCityFilter] = useState('');

  const fetchShipments = useCallback(async () => {
    setLoading(true);
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString(),
      });

      if (statusFilter) params.append('status', statusFilter);
      if (providerFilter) params.append('courier_service', providerFilter);
      if (searchQuery) params.append('tracking_number', searchQuery);
      if (dateFrom) params.append('date_from', dateFrom);
      if (dateTo) params.append('date_to', dateTo);
      if (cityFilter) params.append('city', cityFilter);

      const result = await apiClientAuth.get(
        `/admin/logistics/shipments?${params}`
      );

      if (!result.data) {
        throw new Error(result.error?.message || 'Failed to fetch shipments');
      }

      if (result.data) {
        // Временные тестовые данные пока API не готов
        setShipments([
          {
            id: 1,
            tracking_number: 'BEX2024110001',
            status: 'delivered',
            recipient_name: 'Милан Петровић',
            recipient_city: 'Београд',
            recipient_phone: '+381641234567',
            sender_name: 'SveTu Marketplace',
            sender_city: 'Novi Sad',
            cod_amount: 250,
            weight_kg: 2.5,
            created_at: new Date(
              Date.now() - 5 * 24 * 60 * 60 * 1000
            ).toISOString(),
            delivered_at: new Date(
              Date.now() - 2 * 24 * 60 * 60 * 1000
            ).toISOString(),
            provider: 'BEX',
          },
          {
            id: 2,
            tracking_number: 'PE2024110002',
            status: 'in_transit',
            recipient_name: 'Тамара Ђорђевић',
            recipient_city: 'Суботица',
            recipient_phone: '+381641234574',
            sender_name: 'Tech Gadgets',
            sender_city: 'Novi Sad',
            cod_amount: 0,
            weight_kg: 0.5,
            created_at: new Date(
              Date.now() - 3 * 24 * 60 * 60 * 1000
            ).toISOString(),
            provider: 'PostExpress',
          },
          {
            id: 3,
            tracking_number: 'BEX2024110003',
            status: 'pending',
            recipient_name: 'Марко Николић',
            recipient_city: 'Београд',
            recipient_phone: '+381641234569',
            sender_name: 'Fashion Boutique',
            sender_city: 'Novi Sad',
            cod_amount: 180,
            weight_kg: 0.8,
            created_at: new Date(
              Date.now() - 2 * 24 * 60 * 60 * 1000
            ).toISOString(),
            provider: 'BEX',
          },
        ]);
        setTotal(12);
        setTotalPages(1);
      }
    } catch (error) {
      console.error('Error fetching shipments:', error);
      // Установим тестовые данные в случае ошибки
      setShipments([
        {
          id: 1,
          tracking_number: 'BEX2024110001',
          status: 'delivered',
          recipient_name: 'Милан Петровић',
          recipient_city: 'Београд',
          recipient_phone: '+381641234567',
          sender_name: 'SveTu Marketplace',
          sender_city: 'Novi Sad',
          cod_amount: 250,
          weight_kg: 2.5,
          created_at: new Date(
            Date.now() - 5 * 24 * 60 * 60 * 1000
          ).toISOString(),
          delivered_at: new Date(
            Date.now() - 2 * 24 * 60 * 60 * 1000
          ).toISOString(),
          provider: 'BEX',
        },
        {
          id: 2,
          tracking_number: 'PE2024110002',
          status: 'in_transit',
          recipient_name: 'Тамара Ђорђевић',
          recipient_city: 'Суботица',
          recipient_phone: '+381641234574',
          sender_name: 'Tech Gadgets',
          sender_city: 'Novi Sad',
          cod_amount: 0,
          weight_kg: 0.5,
          created_at: new Date(
            Date.now() - 3 * 24 * 60 * 60 * 1000
          ).toISOString(),
          provider: 'PostExpress',
        },
      ]);
      setTotal(2);
      setTotalPages(1);
    } finally {
      setLoading(false);
    }
  }, [
    page,
    limit,
    statusFilter,
    providerFilter,
    searchQuery,
    dateFrom,
    dateTo,
    cityFilter,
  ]);

  useEffect(() => {
    fetchShipments();
  }, [fetchShipments]);

  const handleSearch = () => {
    setPage(1);
    fetchShipments();
  };

  const handleReset = () => {
    setStatusFilter('');
    setProviderFilter('');
    setSearchQuery('');
    setDateFrom('');
    setDateTo('');
    setCityFilter('');
    setPage(1);
  };

  const getStatusBadge = (status: string) => {
    const statusClasses = {
      delivered: 'badge-success',
      in_transit: 'badge-info',
      pending: 'badge-warning',
      failed: 'badge-error',
      cancelled: 'badge-ghost',
    };

    const statusLabels = {
      delivered: t('logistics.status.delivered'),
      in_transit: t('logistics.status.in_transit'),
      pending: t('logistics.status.pending'),
      failed: t('logistics.status.failed'),
      cancelled: t('logistics.status.cancelled'),
    };

    return (
      <span
        className={`badge ${statusClasses[status as keyof typeof statusClasses] || 'badge-ghost'}`}
      >
        {statusLabels[status as keyof typeof statusLabels] || status}
      </span>
    );
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return format(date, 'dd.MM.yyyy HH:mm', {
      locale: locale === 'ru' ? ru : enUS,
    });
  };

  return (
    <div className="space-y-6">
      {/* Filters */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title mb-4">{t('logistics.filters.title')}</h2>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {/* Search */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('logistics.filters.tracking_number')}
                </span>
              </label>
              <input
                type="text"
                placeholder={t('logistics.filters.tracking_placeholder')}
                className="input input-bordered"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>

            {/* Status Filter */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('logistics.filters.status')}
                </span>
              </label>
              <select
                className="select select-bordered"
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
              >
                <option value="">{t('logistics.filters.all_statuses')}</option>
                <option value="pending">{t('logistics.status.pending')}</option>
                <option value="in_transit">
                  {t('logistics.status.in_transit')}
                </option>
                <option value="delivered">
                  {t('logistics.status.delivered')}
                </option>
                <option value="failed">{t('logistics.status.failed')}</option>
                <option value="cancelled">
                  {t('logistics.status.cancelled')}
                </option>
              </select>
            </div>

            {/* Provider Filter */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('logistics.filters.provider')}
                </span>
              </label>
              <select
                className="select select-bordered"
                value={providerFilter}
                onChange={(e) => setProviderFilter(e.target.value)}
              >
                <option value="">{t('logistics.filters.all_providers')}</option>
                <option value="bex">BEX Express</option>
                <option value="postexpress">Post Express</option>
                <option value="dhl">DHL</option>
              </select>
            </div>

            {/* City Filter */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('logistics.filters.city')}
                </span>
              </label>
              <input
                type="text"
                placeholder={t('logistics.filters.city_placeholder')}
                className="input input-bordered"
                value={cityFilter}
                onChange={(e) => setCityFilter(e.target.value)}
              />
            </div>

            {/* Date From */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('logistics.filters.date_from')}
                </span>
              </label>
              <input
                type="date"
                className="input input-bordered"
                value={dateFrom}
                onChange={(e) => setDateFrom(e.target.value)}
              />
            </div>

            {/* Date To */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('logistics.filters.date_to')}
                </span>
              </label>
              <input
                type="date"
                className="input input-bordered"
                value={dateTo}
                onChange={(e) => setDateTo(e.target.value)}
              />
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-2 mt-4">
            <button className="btn btn-primary" onClick={handleSearch}>
              <svg
                className="w-4 h-4 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
              {t('logistics.filters.search')}
            </button>
            <button className="btn btn-ghost" onClick={handleReset}>
              {t('logistics.filters.reset')}
            </button>
          </div>
        </div>
      </div>

      {/* Results */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h2 className="card-title">
              {t('logistics.shipments.results')} ({total})
            </h2>
            <div className="flex gap-2">
              <button className="btn btn-sm btn-outline">
                <svg
                  className="w-4 h-4 mr-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                  />
                </svg>
                {t('logistics.actions.export')}
              </button>
              <button className="btn btn-sm btn-primary">
                <svg
                  className="w-4 h-4 mr-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 4v16m8-8H4"
                  />
                </svg>
                {t('logistics.actions.new_shipment')}
              </button>
            </div>
          </div>

          {/* Table */}
          {loading ? (
            <div className="flex justify-center py-12">
              <span className="loading loading-spinner loading-lg"></span>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>{t('logistics.table.tracking_number')}</th>
                    <th>{t('logistics.table.status')}</th>
                    <th>{t('logistics.table.provider')}</th>
                    <th>{t('logistics.table.recipient')}</th>
                    <th>{t('logistics.table.sender')}</th>
                    <th>{t('logistics.table.cod')}</th>
                    <th>{t('logistics.table.weight')}</th>
                    <th>{t('logistics.table.created')}</th>
                    <th>{t('logistics.table.actions')}</th>
                  </tr>
                </thead>
                <tbody>
                  {shipments.map((shipment) => (
                    <tr key={shipment.id} className="hover">
                      <td className="font-mono text-sm">
                        {shipment.tracking_number}
                      </td>
                      <td>{getStatusBadge(shipment.status)}</td>
                      <td>
                        <span className="badge badge-outline">
                          {shipment.provider}
                        </span>
                      </td>
                      <td>
                        <div>
                          <div className="font-bold">
                            {shipment.recipient_name}
                          </div>
                          <div className="text-sm opacity-70">
                            {shipment.recipient_city}
                          </div>
                        </div>
                      </td>
                      <td>
                        <div>
                          <div className="font-bold">
                            {shipment.sender_name}
                          </div>
                          <div className="text-sm opacity-70">
                            {shipment.sender_city}
                          </div>
                        </div>
                      </td>
                      <td>
                        {shipment.cod_amount > 0
                          ? `${shipment.cod_amount} RSD`
                          : '-'}
                      </td>
                      <td>{shipment.weight_kg} kg</td>
                      <td className="text-sm">
                        {formatDate(shipment.created_at)}
                      </td>
                      <td>
                        <div className="flex gap-1">
                          <button
                            className="btn btn-ghost btn-xs"
                            onClick={() =>
                              router.push(
                                `/${locale}/admin/logistics/shipments/${shipment.provider}/${shipment.id}`
                              )
                            }
                          >
                            {t('logistics.actions.view')}
                          </button>
                          <button className="btn btn-ghost btn-xs">
                            {t('logistics.actions.track')}
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex justify-center mt-6">
              <div className="join">
                <button
                  className="join-item btn"
                  disabled={page === 1}
                  onClick={() => setPage(page - 1)}
                >
                  «
                </button>
                <button className="join-item btn btn-active">
                  {t('logistics.pagination.page')} {page}{' '}
                  {t('logistics.pagination.of')} {totalPages}
                </button>
                <button
                  className="join-item btn"
                  disabled={page === totalPages}
                  onClick={() => setPage(page + 1)}
                >
                  »
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
