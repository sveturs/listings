'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useSearchParams, useRouter } from 'next/navigation';
import { tokenManager } from '@/utils/tokenManager';

interface Storefront {
  id: number;
  user_id: number;
  slug: string;
  name: string;
  description: string;
  logo_url?: string;
  banner_url?: string;
  phone?: string;
  email?: string;
  website?: string;
  address?: string;
  city?: string;
  postal_code?: string;
  country?: string;
  latitude?: number;
  longitude?: number;
  is_active: boolean;
  is_verified: boolean;
  rating: number;
  reviews_count: number;
  products_count: number;
  sales_count: number;
  views_count: number;
  subscription_plan: string;
  commission_rate: number;
  created_at: string;
  updated_at: string;
}

export default function AdminStorefrontsTable() {
  const t = useTranslations('admin');
  const searchParams = useSearchParams();
  const router = useRouter();

  const [storefronts, setStorefronts] = useState<Storefront[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [selectedStorefronts, setSelectedStorefronts] = useState<Set<number>>(
    new Set()
  );
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [storefrontToDelete, setStorefrontToDelete] = useState<number | null>(
    null
  );
  const [deleteType, setDeleteType] = useState<'soft' | 'hard'>('soft'); // soft = мягкое, hard = жесткое удаление
  const [, setBulkAction] = useState('');

  // Фильтры и пагинация
  const [searchQuery, setSearchQuery] = useState(
    searchParams.get('query') || ''
  );
  const [statusFilter, setStatusFilter] = useState(
    searchParams.get('status') || ''
  );
  const [verifiedFilter, setVerifiedFilter] = useState(
    searchParams.get('verified') || ''
  );
  const [page, setPage] = useState(parseInt(searchParams.get('page') || '1'));
  const [limit] = useState(20);
  const [sortBy, setSortBy] = useState(searchParams.get('sort') || 'date_desc');

  // Загрузка данных
  const fetchStorefronts = useCallback(async () => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      if (searchQuery) params.append('query', searchQuery);
      params.append('limit', limit.toString());
      params.append('offset', ((page - 1) * limit).toString());
      // Добавляем параметр для получения всех витрин, включая неактивные
      params.append('include_inactive', 'true');

      const response = await fetch(
        `http://localhost:3000/api/v1/storefronts?${params}`,
        {
          credentials: 'include',
        }
      );

      if (response.ok) {
        const data = await response.json();
        let storefrontsList = data.storefronts || data.data || [];

        // Применяем фильтры
        if (statusFilter) {
          if (statusFilter === 'active') {
            storefrontsList = storefrontsList.filter(
              (s: Storefront) => s.is_active
            );
          } else if (statusFilter === 'inactive') {
            storefrontsList = storefrontsList.filter(
              (s: Storefront) => !s.is_active
            );
          }
        }

        if (verifiedFilter) {
          if (verifiedFilter === 'verified') {
            storefrontsList = storefrontsList.filter(
              (s: Storefront) => s.is_verified
            );
          } else if (verifiedFilter === 'unverified') {
            storefrontsList = storefrontsList.filter(
              (s: Storefront) => !s.is_verified
            );
          }
        }

        // Сортировка
        storefrontsList.sort((a: Storefront, b: Storefront) => {
          switch (sortBy) {
            case 'name_asc':
              return a.name.localeCompare(b.name);
            case 'name_desc':
              return b.name.localeCompare(a.name);
            case 'products_desc':
              return b.products_count - a.products_count;
            case 'sales_desc':
              return b.sales_count - a.sales_count;
            case 'rating_desc':
              return b.rating - a.rating;
            case 'date_asc':
              return (
                new Date(a.created_at).getTime() -
                new Date(b.created_at).getTime()
              );
            case 'date_desc':
            default:
              return (
                new Date(b.created_at).getTime() -
                new Date(a.created_at).getTime()
              );
          }
        });

        setStorefronts(storefrontsList);
        setTotal(data.total || storefrontsList.length);
      }
    } catch (error) {
      console.error('Error fetching storefronts:', error);
      setStorefronts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [searchQuery, statusFilter, verifiedFilter, page, limit, sortBy]);

  useEffect(() => {
    fetchStorefronts();
  }, [fetchStorefronts]);

  // Обработка удаления
  const handleDelete = async (id: number, type: 'soft' | 'hard' = 'soft') => {
    try {
      const token = tokenManager.getAccessToken();
      const url =
        type === 'hard'
          ? `http://localhost:3000/api/v1/storefronts/${id}?hard_delete=true`
          : `http://localhost:3000/api/v1/storefronts/${id}`;

      const response = await fetch(url, {
        method: 'DELETE',
        credentials: 'include',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        fetchStorefronts();
        setDeleteModalOpen(false);
        setStorefrontToDelete(null);
        setDeleteType('soft'); // Сброс к значению по умолчанию
      }
    } catch (error) {
      console.error('Error deleting storefront:', error);
    }
  };

  // Обработка восстановления
  const handleRestore = async (id: number) => {
    try {
      const token = tokenManager.getAccessToken();
      const response = await fetch(
        `http://localhost:3000/api/v1/storefronts/${id}/restore`,
        {
          method: 'POST',
          credentials: 'include',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.ok) {
        fetchStorefronts();
      }
    } catch (error) {
      console.error('Error restoring storefront:', error);
    }
  };

  // Массовые действия
  const handleBulkAction = async (action: string) => {
    if (selectedStorefronts.size === 0) return;

    if (action === 'soft-delete' || action === 'hard-delete') {
      const deleteType = action === 'hard-delete' ? 'hard' : 'soft';
      const confirmed = window.confirm(
        deleteType === 'hard'
          ? `${t('storefronts.confirmBulkDelete')} ${t('storefronts.hardDeleteDesc')}`
          : `${t('storefronts.confirmBulkDelete')} ${t('storefronts.softDeleteDesc')}`
      );
      if (!confirmed) return;

      try {
        const token = tokenManager.getAccessToken();
        for (const id of selectedStorefronts) {
          const url =
            deleteType === 'hard'
              ? `http://localhost:3000/api/v1/storefronts/${id}?hard_delete=true`
              : `http://localhost:3000/api/v1/storefronts/${id}`;

          await fetch(url, {
            method: 'DELETE',
            credentials: 'include',
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
        }
        setSelectedStorefronts(new Set());
        fetchStorefronts();
      } catch (error) {
        console.error('Error performing bulk action:', error);
      }
    } else if (action === 'activate' || action === 'deactivate') {
      try {
        const token = tokenManager.getAccessToken();
        const isActive = action === 'activate';

        for (const id of selectedStorefronts) {
          await fetch(`http://localhost:3000/api/v1/storefronts/${id}`, {
            method: 'PUT',
            credentials: 'include',
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({ is_active: isActive }),
          });
        }
        setSelectedStorefronts(new Set());
        fetchStorefronts();
      } catch (error) {
        console.error('Error performing bulk action:', error);
      }
    }
    setBulkAction('');
  };

  // Переключение выбора
  const toggleSelection = (id: number) => {
    const newSelection = new Set(selectedStorefronts);
    if (newSelection.has(id)) {
      newSelection.delete(id);
    } else {
      newSelection.add(id);
    }
    setSelectedStorefronts(newSelection);
  };

  const selectAll = () => {
    if (selectedStorefronts.size === storefronts.length) {
      setSelectedStorefronts(new Set());
    } else {
      setSelectedStorefronts(new Set(storefronts.map((s) => s.id)));
    }
  };

  const totalPages = Math.ceil(total / limit);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="card bg-base-100 shadow-lg">
      <div className="card-body">
        {/* Панель фильтров */}
        <div className="flex justify-between items-center mb-4">
          <div className="flex gap-2 flex-1">
            <input
              type="text"
              placeholder={t('storefronts.searchPlaceholder')}
              className="input input-bordered flex-1 max-w-md"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />

            <select
              className="select select-bordered"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
            >
              <option value="">{t('storefronts.filters.allStatus')}</option>
              <option value="active">{t('storefronts.filters.active')}</option>
              <option value="inactive">
                {t('storefronts.filters.inactive')}
              </option>
            </select>

            <select
              className="select select-bordered"
              value={verifiedFilter}
              onChange={(e) => setVerifiedFilter(e.target.value)}
            >
              <option value="">
                {t('storefronts.filters.allVerification')}
              </option>
              <option value="verified">
                {t('storefronts.filters.verified')}
              </option>
              <option value="unverified">
                {t('storefronts.filters.unverified')}
              </option>
            </select>

            <select
              className="select select-bordered"
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
            >
              <option value="date_desc">{t('storefronts.sort.newest')}</option>
              <option value="date_asc">{t('storefronts.sort.oldest')}</option>
              <option value="name_asc">{t('storefronts.sort.nameAsc')}</option>
              <option value="name_desc">
                {t('storefronts.sort.nameDesc')}
              </option>
              <option value="products_desc">
                {t('storefronts.sort.mostProducts')}
              </option>
              <option value="sales_desc">
                {t('storefronts.sort.mostSales')}
              </option>
              <option value="rating_desc">
                {t('storefronts.sort.bestRating')}
              </option>
            </select>
          </div>

          {/* Массовые действия */}
          {selectedStorefronts.size > 0 && (
            <div className="dropdown dropdown-end">
              <label tabIndex={0} className="btn btn-primary dropdown-toggle">
                {t('storefronts.bulkActions')} ({selectedStorefronts.size})
              </label>
              <ul
                tabIndex={0}
                className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52"
              >
                <li>
                  <button
                    onClick={() => handleBulkAction('activate')}
                    className="text-success"
                  >
                    {t('storefronts.activateSelected')}
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => handleBulkAction('deactivate')}
                    className="text-warning"
                  >
                    {t('storefronts.deactivateSelected')}
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => handleBulkAction('soft-delete')}
                    className="text-warning"
                  >
                    {t('storefronts.softDelete')} (
                    {t('storefronts.deleteSelected')})
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => handleBulkAction('hard-delete')}
                    className="text-error"
                  >
                    {t('storefronts.hardDelete')} (
                    {t('storefronts.deleteSelected')})
                  </button>
                </li>
              </ul>
            </div>
          )}
        </div>

        {/* Таблица */}
        <div className="overflow-x-auto">
          <table className="table table-zebra">
            <thead>
              <tr>
                <th>
                  <label>
                    <input
                      type="checkbox"
                      className="checkbox"
                      checked={
                        selectedStorefronts.size === storefronts.length &&
                        storefronts.length > 0
                      }
                      onChange={selectAll}
                    />
                  </label>
                </th>
                <th>{t('storefronts.table.name')}</th>
                <th>{t('storefronts.table.owner')}</th>
                <th>{t('storefronts.table.location')}</th>
                <th>{t('storefronts.table.products')}</th>
                <th>{t('storefronts.table.sales')}</th>
                <th>{t('storefronts.table.rating')}</th>
                <th>{t('storefronts.table.status')}</th>
                <th>{t('storefronts.table.verified')}</th>
                <th>{t('storefronts.table.actions')}</th>
              </tr>
            </thead>
            <tbody>
              {storefronts.length === 0 ? (
                <tr>
                  <td colSpan={10} className="text-center py-8">
                    <p className="text-base-content/60">
                      {t('storefronts.noData')}
                    </p>
                  </td>
                </tr>
              ) : (
                storefronts.map((storefront) => (
                  <tr key={storefront.id}>
                    <td>
                      <label>
                        <input
                          type="checkbox"
                          className="checkbox"
                          checked={selectedStorefronts.has(storefront.id)}
                          onChange={() => toggleSelection(storefront.id)}
                        />
                      </label>
                    </td>
                    <td>
                      <div className="flex items-center gap-3">
                        {storefront.logo_url && (
                          <div className="avatar">
                            <div className="mask mask-squircle w-12 h-12">
                              <img
                                src={storefront.logo_url}
                                alt={storefront.name}
                              />
                            </div>
                          </div>
                        )}
                        <div>
                          <div className="font-bold">{storefront.name}</div>
                          <div className="text-sm opacity-50">
                            @{storefront.slug}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div className="text-sm">ID: {storefront.user_id}</div>
                    </td>
                    <td>
                      <div className="text-sm">
                        {storefront.city || storefront.address || '-'}
                      </div>
                    </td>
                    <td>
                      <div className="badge badge-ghost">
                        {storefront.products_count}
                      </div>
                    </td>
                    <td>
                      <div className="badge badge-ghost">
                        {storefront.sales_count}
                      </div>
                    </td>
                    <td>
                      <div className="flex items-center gap-1">
                        <span className="text-warning">★</span>
                        <span>{storefront.rating.toFixed(1)}</span>
                        <span className="text-xs opacity-50">
                          ({storefront.reviews_count})
                        </span>
                      </div>
                    </td>
                    <td>
                      <span
                        className={`badge ${
                          storefront.is_active ? 'badge-success' : 'badge-error'
                        }`}
                      >
                        {storefront.is_active
                          ? t('storefronts.active')
                          : t('storefronts.inactive')}
                      </span>
                    </td>
                    <td>
                      {storefront.is_verified ? (
                        <span className="badge badge-primary">
                          {t('storefronts.verified')}
                        </span>
                      ) : (
                        <span className="badge badge-outline">
                          {t('storefronts.unverified')}
                        </span>
                      )}
                    </td>
                    <td>
                      <div className="dropdown dropdown-end">
                        <label tabIndex={0} className="btn btn-ghost btn-sm">
                          <svg
                            className="w-5 h-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
                            />
                          </svg>
                        </label>
                        <ul
                          tabIndex={0}
                          className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52"
                        >
                          <li>
                            <a
                              href={`/${storefront.slug}`}
                              target="_blank"
                              rel="noopener noreferrer"
                            >
                              {t('storefronts.view')}
                            </a>
                          </li>
                          <li>
                            <button
                              onClick={() =>
                                router.push(
                                  `/admin/storefronts/${storefront.id}/edit`
                                )
                              }
                            >
                              {t('storefronts.edit')}
                            </button>
                          </li>
                          <li>
                            <button
                              onClick={() =>
                                router.push(
                                  `/admin/storefront-products?storefront=${storefront.id}`
                                )
                              }
                            >
                              {t('storefronts.viewProducts')}
                            </button>
                          </li>
                          {!storefront.is_active && (
                            <li>
                              <button
                                onClick={() => handleRestore(storefront.id)}
                                className="text-success"
                              >
                                {t('storefronts.restore')}
                              </button>
                            </li>
                          )}
                          <li>
                            <button
                              onClick={() => {
                                setStorefrontToDelete(storefront.id);
                                setDeleteModalOpen(true);
                              }}
                              className="text-error"
                            >
                              {t('storefronts.delete')}
                            </button>
                          </li>
                        </ul>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Пагинация */}
        {totalPages > 1 && (
          <div className="flex justify-center mt-4">
            <div className="join">
              <button
                className="join-item btn"
                disabled={page === 1}
                onClick={() => setPage(page - 1)}
              >
                «
              </button>
              {[...Array(Math.min(5, totalPages))].map((_, i) => {
                const pageNumber = i + 1;
                return (
                  <button
                    key={pageNumber}
                    className={`join-item btn ${
                      page === pageNumber ? 'btn-active' : ''
                    }`}
                    onClick={() => setPage(pageNumber)}
                  >
                    {pageNumber}
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

      {/* Модальное окно удаления */}
      {deleteModalOpen && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">
              {t('storefronts.confirmDeleteTitle')}
            </h3>
            <p className="py-4">{t('storefronts.confirmDeleteMessage')}</p>

            {/* Выбор типа удаления для администратора */}
            <div className="form-control py-4">
              <label className="label">
                <span className="label-text font-semibold">
                  {t('storefronts.deleteType')}
                </span>
              </label>
              <div className="space-y-2">
                <label className="label cursor-pointer justify-start gap-4">
                  <input
                    type="radio"
                    name="deleteType"
                    value="soft"
                    checked={deleteType === 'soft'}
                    onChange={(e) =>
                      setDeleteType(e.target.value as 'soft' | 'hard')
                    }
                    className="radio radio-primary"
                  />
                  <div>
                    <span className="label-text font-medium">
                      {t('storefronts.softDelete')}
                    </span>
                    <div className="label-text-alt text-warning">
                      {t('storefronts.softDeleteDesc')}
                    </div>
                  </div>
                </label>
                <label className="label cursor-pointer justify-start gap-4">
                  <input
                    type="radio"
                    name="deleteType"
                    value="hard"
                    checked={deleteType === 'hard'}
                    onChange={(e) =>
                      setDeleteType(e.target.value as 'soft' | 'hard')
                    }
                    className="radio radio-error"
                  />
                  <div>
                    <span className="label-text font-medium text-error">
                      {t('storefronts.hardDelete')}
                    </span>
                    <div className="label-text-alt text-error">
                      {t('storefronts.hardDeleteDesc')}
                    </div>
                  </div>
                </label>
              </div>
            </div>

            <div className="modal-action">
              <button
                className="btn"
                onClick={() => {
                  setDeleteModalOpen(false);
                  setStorefrontToDelete(null);
                  setDeleteType('soft');
                }}
              >
                {t('common.cancel')}
              </button>
              <button
                className={`btn ${deleteType === 'hard' ? 'btn-error' : 'btn-warning'}`}
                onClick={() => {
                  if (storefrontToDelete) {
                    handleDelete(storefrontToDelete, deleteType);
                  }
                }}
              >
                {deleteType === 'hard'
                  ? t('storefronts.hardDelete')
                  : t('storefronts.softDelete')}
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button
              onClick={() => {
                setDeleteModalOpen(false);
                setStorefrontToDelete(null);
                setDeleteType('soft');
              }}
            >
              close
            </button>
          </form>
        </dialog>
      )}
    </div>
  );
}
