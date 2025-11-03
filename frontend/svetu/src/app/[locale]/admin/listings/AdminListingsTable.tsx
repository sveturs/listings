'use client';

import React, { useState, useCallback, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { format } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import { useLocale } from 'next-intl';
import { apiClient } from '@/services/api-client';

interface Listing {
  id: number;
  title: string;
  description: string;
  price: number;
  currency: string;
  status: string;
  condition: string;
  category_id: number;
  category_name?: string;
  category?: {
    id: number;
    name: string;
    slug: string;
  };
  user_id: number;
  user_email?: string;
  user_name?: string;
  user?: {
    id: number;
    name: string;
    email: string;
  };
  storefront_id?: number;
  storefront_name?: string;
  address_city?: string;
  address_country?: string;
  city?: string;
  district?: string;
  images?: Array<
    | {
        id: number;
        listing_id: number;
        file_path: string;
        file_name: string;
        public_url: string;
        is_main: boolean;
      }
    | string
  > | null;
  views_count: number;
  created_at: string;
  updated_at: string;
  is_active?: boolean;
  is_promoted?: boolean;
  promotion_end_date?: string;
}

export default function AdminListingsTable() {
  const t = useTranslations('admin');
  const locale = useLocale();
  const router = useRouter();
  const searchParams = useSearchParams();

  const [listings, setListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [selectedListings, setSelectedListings] = useState<Set<number>>(
    new Set()
  );
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [listingToDelete, setListingToDelete] = useState<number | null>(null);

  // Фильтры и пагинация
  const [searchQuery, setSearchQuery] = useState(
    searchParams.get('query') || ''
  );
  const [categoryFilter] = useState(searchParams.get('category') || '');
  const [statusFilter, setStatusFilter] = useState(
    searchParams.get('status') || ''
  );
  const [page, setPage] = useState(parseInt(searchParams.get('page') || '1'));
  const [limit] = useState(20);
  const [sortBy, setSortBy] = useState(searchParams.get('sort') || 'date_desc');

  // Загрузка данных
  const fetchListings = useCallback(async () => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      if (searchQuery) params.append('query', searchQuery);
      if (categoryFilter) params.append('category_id', categoryFilter);
      if (statusFilter) params.append('status', statusFilter);
      params.append('limit', limit.toString());
      params.append('offset', ((page - 1) * limit).toString());
      params.append('sort_by', sortBy);
      // Исключить товары витрин - показываем только P2P объявления
      params.append('exclude_storefronts', 'true');

      const response = await apiClient.get(`/marketplace/listings?${params}`);

      if (response.data?.success && response.data.data) {
        // API возвращает data.data для массива объявлений и data.meta для метаданных
        setListings(response.data.data.data || []);
        setTotal(response.data.data.meta?.total || 0);
      }
    } catch (error) {
      console.error('Error fetching listings:', error);
    } finally {
      setLoading(false);
    }
  }, [searchQuery, categoryFilter, statusFilter, page, limit, sortBy]);

  useEffect(() => {
    fetchListings();
  }, [fetchListings]);

  // Обновление URL при изменении фильтров
  useEffect(() => {
    const params = new URLSearchParams();
    if (searchQuery) params.set('query', searchQuery);
    if (categoryFilter) params.set('category', categoryFilter);
    if (statusFilter) params.set('status', statusFilter);
    if (page > 1) params.set('page', page.toString());
    if (sortBy !== 'date_desc') params.set('sort', sortBy);

    const newUrl = params.toString() ? `?${params}` : '';
    router.push(`/${locale}/admin/listings${newUrl}`, { scroll: false });
  }, [searchQuery, categoryFilter, statusFilter, page, sortBy, router, locale]);

  // Обработчики действий
  const handleDelete = async (id: number) => {
    try {
      const response = await apiClient.delete(`/marketplace/listings/${id}`);

      if (response.data) {
        await fetchListings();
        setDeleteModalOpen(false);
        setListingToDelete(null);
      }
    } catch (error) {
      console.error('Error deleting listing:', error);
      alert('Ошибка при удалении объявления');
    }
  };

  const handleBulkDelete = async () => {
    for (const id of selectedListings) {
      await handleDelete(id);
    }
    setSelectedListings(new Set());
  };

  const handleToggleActive = async (id: number, isActive: boolean) => {
    try {
      const response = await apiClient.put(`/marketplace/listings/${id}`, {
        is_active: !isActive,
      });

      if (response.data) {
        await fetchListings();
      }
    } catch (error) {
      console.error('Error toggling listing status:', error);
      alert('Ошибка при изменении статуса объявления');
    }
  };

  const handleSelectAll = () => {
    if (selectedListings.size === listings.length) {
      setSelectedListings(new Set());
    } else {
      setSelectedListings(new Set(listings.map((l) => l.id)));
    }
  };

  const handleSelectOne = (id: number) => {
    const newSelected = new Set(selectedListings);
    if (newSelected.has(id)) {
      newSelected.delete(id);
    } else {
      newSelected.add(id);
    }
    setSelectedListings(newSelected);
  };

  const totalPages = Math.ceil(total / limit);

  const getStatusBadge = (status: string, isActive?: boolean) => {
    if (isActive === false)
      return (
        <span className="badge badge-ghost">
          {t('listings.status.inactive')}
        </span>
      );

    switch (status) {
      case 'active':
        return (
          <span className="badge badge-success">
            {t('listings.status.active')}
          </span>
        );
      case 'pending':
        return (
          <span className="badge badge-warning">
            {t('listings.status.pending')}
          </span>
        );
      case 'sold':
        return (
          <span className="badge badge-info">{t('listings.status.sold')}</span>
        );
      default:
        return <span className="badge">{status}</span>;
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const dateLocale = locale === 'ru' ? ru : enUS;
    return format(date, 'dd MMM yyyy HH:mm', { locale: dateLocale });
  };

  if (loading && listings.length === 0) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Фильтры и поиск */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <div className="flex flex-wrap gap-4">
            <div className="form-control flex-1 min-w-[200px]">
              <input
                type="text"
                placeholder={t('listings.search.placeholder')}
                className="input input-bordered"
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value);
                  setPage(1);
                }}
              />
            </div>

            <select
              className="select select-bordered"
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value);
                setPage(1);
              }}
            >
              <option value="">{t('listings.filter.allStatuses')}</option>
              <option value="active">{t('listings.status.active')}</option>
              <option value="pending">{t('listings.status.pending')}</option>
              <option value="sold">{t('listings.status.sold')}</option>
              <option value="inactive">{t('listings.status.inactive')}</option>
            </select>

            <select
              className="select select-bordered"
              value={sortBy}
              onChange={(e) => {
                setSortBy(e.target.value);
                setPage(1);
              }}
            >
              <option value="date_desc">{t('listings.sort.dateDesc')}</option>
              <option value="date_asc">{t('listings.sort.dateAsc')}</option>
              <option value="price_desc">{t('listings.sort.priceDesc')}</option>
              <option value="price_asc">{t('listings.sort.priceAsc')}</option>
              <option value="views_desc">{t('listings.sort.viewsDesc')}</option>
            </select>

            {selectedListings.size > 0 && (
              <button
                className="btn btn-error"
                onClick={() => {
                  if (confirm(t('listings.confirmBulkDelete'))) {
                    handleBulkDelete();
                  }
                }}
              >
                {t('listings.deleteSelected')} ({selectedListings.size})
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Таблица объявлений */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body p-0">
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>
                    <input
                      type="checkbox"
                      className="checkbox"
                      checked={
                        selectedListings.size === listings.length &&
                        listings.length > 0
                      }
                      onChange={handleSelectAll}
                    />
                  </th>
                  <th>{t('listings.table.id')}</th>
                  <th>{t('listings.table.image')}</th>
                  <th>{t('listings.table.title')}</th>
                  <th>{t('listings.table.price')}</th>
                  <th>{t('listings.table.category')}</th>
                  <th>{t('listings.table.user')}</th>
                  <th>{t('listings.table.status')}</th>
                  <th>{t('listings.table.views')}</th>
                  <th>{t('listings.table.created')}</th>
                  <th>{t('listings.table.actions')}</th>
                </tr>
              </thead>
              <tbody>
                {listings.map((listing) => (
                  <tr key={listing.id}>
                    <td>
                      <input
                        type="checkbox"
                        className="checkbox"
                        checked={selectedListings.has(listing.id)}
                        onChange={() => handleSelectOne(listing.id)}
                      />
                    </td>
                    <td>{listing.id}</td>
                    <td>
                      {listing.images &&
                      listing.images.length > 0 &&
                      listing.images[0] ? (
                        <div className="avatar">
                          <div className="w-12 h-12 rounded bg-base-200">
                            <img
                              src={(() => {
                                const firstImage = listing.images![0];
                                // Проверяем, является ли изображение объектом или строкой
                                const imagePath =
                                  typeof firstImage === 'object' &&
                                  firstImage !== null
                                    ? (firstImage as any).public_url ||
                                      (firstImage as any).file_path
                                    : firstImage;

                                // Добавляем базовый URL если нужно
                                if (!imagePath) return '/placeholder.jpg';
                                if (typeof imagePath !== 'string')
                                  return '/placeholder.jpg';
                                return imagePath.startsWith('http')
                                  ? imagePath
                                  : imagePath.startsWith('/')
                                    ? imagePath
                                    : `/listings/${imagePath}`;
                              })()}
                              alt={listing.title}
                              className="w-full h-full object-cover rounded"
                              onError={(e) => {
                                const target = e.target as HTMLImageElement;
                                target.style.display = 'none';
                                if (target.parentElement) {
                                  target.parentElement.innerHTML = `
                                    <svg class="w-6 h-6 text-base-content/20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                    </svg>
                                  `;
                                }
                              }}
                            />
                          </div>
                        </div>
                      ) : (
                        <div className="w-12 h-12 bg-base-200 rounded flex items-center justify-center">
                          <svg
                            className="w-6 h-6 text-base-content/20"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                            />
                          </svg>
                        </div>
                      )}
                    </td>
                    <td>
                      <div className="flex items-center gap-2">
                        <div className="font-medium">{listing.title}</div>
                        {listing.storefront_id && (
                          <span className="badge badge-sm badge-info">
                            Storefront
                          </span>
                        )}
                      </div>
                      <div className="text-xs text-base-content/60 truncate max-w-xs">
                        {listing.description}
                      </div>
                    </td>
                    <td>
                      <div className="font-semibold">
                        {listing.price} {listing.currency || 'RSD'}
                      </div>
                    </td>
                    <td>
                      {listing.category?.name ||
                        listing.category_name ||
                        `ID: ${listing.category_id}`}
                    </td>
                    <td>
                      <div className="text-sm">
                        <div>
                          {listing.user?.name ||
                            listing.user_name ||
                            `User #${listing.user_id}`}
                        </div>
                        {(listing.user?.email || listing.user_email) && (
                          <div className="text-xs text-base-content/60">
                            {listing.user?.email || listing.user_email}
                          </div>
                        )}
                        {listing.storefront_name && (
                          <div className="text-xs text-info">
                            {listing.storefront_name}
                          </div>
                        )}
                      </div>
                    </td>
                    <td>{getStatusBadge(listing.status, listing.is_active)}</td>
                    <td>
                      <div className="flex items-center gap-1">
                        <svg
                          className="w-4 h-4"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                          />
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                          />
                        </svg>
                        {listing.views_count}
                      </div>
                    </td>
                    <td>
                      <div className="text-xs">
                        {formatDate(listing.created_at)}
                      </div>
                    </td>
                    <td>
                      <div className="flex items-center gap-1">
                        <Link
                          href={`/${locale}/c2c/${listing.id}`}
                          className="btn btn-ghost btn-xs"
                          target="_blank"
                        >
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                            />
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                            />
                          </svg>
                        </Link>

                        <Link
                          href={`/${locale}/admin/listings/${listing.id}/edit`}
                          className="btn btn-ghost btn-xs"
                        >
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                            />
                          </svg>
                        </Link>

                        <button
                          className="btn btn-ghost btn-xs"
                          onClick={() =>
                            handleToggleActive(
                              listing.id,
                              listing.is_active || false
                            )
                          }
                        >
                          {listing.is_active ? (
                            <svg
                              className="w-4 h-4 text-warning"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                              />
                            </svg>
                          ) : (
                            <svg
                              className="w-4 h-4 text-success"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                              />
                            </svg>
                          )}
                        </button>

                        <button
                          className="btn btn-ghost btn-xs text-error"
                          onClick={() => {
                            setListingToDelete(listing.id);
                            setDeleteModalOpen(true);
                          }}
                        >
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            />
                          </svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Пагинация */}
          {totalPages > 1 && (
            <div className="flex justify-center p-4">
              <div className="join">
                <button
                  className="join-item btn"
                  disabled={page === 1}
                  onClick={() => setPage(page - 1)}
                >
                  «
                </button>

                {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                  let pageNum;
                  if (totalPages <= 5) {
                    pageNum = i + 1;
                  } else if (page <= 3) {
                    pageNum = i + 1;
                  } else if (page >= totalPages - 2) {
                    pageNum = totalPages - 4 + i;
                  } else {
                    pageNum = page - 2 + i;
                  }

                  return (
                    <button
                      key={pageNum}
                      className={`join-item btn ${pageNum === page ? 'btn-active' : ''}`}
                      onClick={() => setPage(pageNum)}
                    >
                      {pageNum}
                    </button>
                  );
                })}

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

      {/* Модальное окно удаления */}
      {deleteModalOpen && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">
              {t('listings.deleteModal.title')}
            </h3>
            <p className="py-4">{t('listings.deleteModal.message')}</p>
            <div className="modal-action">
              <button
                className="btn"
                onClick={() => {
                  setDeleteModalOpen(false);
                  setListingToDelete(null);
                }}
              >
                {t('common.cancel')}
              </button>
              <button
                className="btn btn-error"
                onClick={() => {
                  if (listingToDelete) {
                    handleDelete(listingToDelete);
                  }
                }}
              >
                {t('common.delete')}
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button onClick={() => setDeleteModalOpen(false)}>close</button>
          </form>
        </dialog>
      )}
    </div>
  );
}
