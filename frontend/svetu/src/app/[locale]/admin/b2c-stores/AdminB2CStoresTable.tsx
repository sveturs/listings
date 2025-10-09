'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useSearchParams, useRouter } from 'next/navigation';
import { apiClient } from '@/services/api-client';

interface B2CStore {
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

  const [storefronts, setStorefronts] = useState<B2CStore[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [stats, setStats] = useState({
    total: 0,
    active: 0,
    verified: 0,
    totalProducts: 0,
    totalSales: 0,
  });
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

      const response = await apiClient.get(`/admin/b2c?${params}`);

      if (response.data) {
        let storefrontsList =
          response.data.b2c_stores || response.data.storefronts || response.data.data || [];

        // Применяем фильтры
        if (statusFilter) {
          if (statusFilter === 'active') {
            storefrontsList = storefrontsList.filter(
              (s: B2CStore) => s.is_active
            );
          } else if (statusFilter === 'inactive') {
            storefrontsList = storefrontsList.filter(
              (s: B2CStore) => !s.is_active
            );
          }
        }

        if (verifiedFilter) {
          if (verifiedFilter === 'verified') {
            storefrontsList = storefrontsList.filter(
              (s: B2CStore) => s.is_verified
            );
          } else if (verifiedFilter === 'unverified') {
            storefrontsList = storefrontsList.filter(
              (s: B2CStore) => !s.is_verified
            );
          }
        }

        // Сортировка
        storefrontsList.sort((a: B2CStore, b: B2CStore) => {
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
        setTotal(response.data.total || storefrontsList.length);

        // Вычисляем статистику из полученных данных
        const allStorefronts = response.data.b2c_stores || response.data.storefronts || response.data.data || [];
        setStats({
          total: allStorefronts.length,
          active: allStorefronts.filter((s: B2CStore) => s.is_active).length,
          verified: allStorefronts.filter((s: B2CStore) => s.is_verified).length,
          totalProducts: allStorefronts.reduce(
            (sum: number, s: B2CStore) => sum + (s.products_count || 0),
            0
          ),
          totalSales: allStorefronts.reduce(
            (sum: number, s: B2CStore) => sum + (s.sales_count || 0),
            0
          ),
        });
      }
    } catch (error) {
      console.error('Error fetching storefronts:', error);
      setStorefronts([]);
      setTotal(0);
      setStats({
        total: 0,
        active: 0,
        verified: 0,
        totalProducts: 0,
        totalSales: 0,
      });
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
      const url =
        type === 'hard' ? `/b2c/${id}?hard_delete=true` : `/b2c/${id}`;

      const response = await apiClient.delete(url);

      if (response.data) {
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
      const response = await apiClient.post(`/b2c/${id}/restore`);

      if (response.data) {
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
          ? `${t('b2c.confirmBulkDelete')} ${t('b2c.hardDeleteDesc')}`
          : `${t('b2c.confirmBulkDelete')} ${t('b2c.softDeleteDesc')}`
      );
      if (!confirmed) return;

      try {
        for (const id of selectedStorefronts) {
          const url =
            deleteType === 'hard'
              ? `/b2c/${id}?hard_delete=true`
              : `/b2c/${id}`;

          await apiClient.delete(url);
        }
        setSelectedStorefronts(new Set());
        fetchStorefronts();
      } catch (error) {
        console.error('Error performing bulk action:', error);
      }
    } else if (action === 'activate' || action === 'deactivate') {
      try {
        const isActive = action === 'activate';

        for (const id of selectedStorefronts) {
          await apiClient.put(`/b2c/${id}`, { is_active: isActive });
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
    <div>
      {/* Заголовок и кнопки */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold">{t('b2c.title')}</h1>
          <p className="text-base-content/60 mt-1">{t('b2c.description')}</p>
        </div>

        <div className="flex gap-2">
          <button className="btn btn-primary">
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
                d="M12 4v16m8-8H4"
              />
            </svg>
            {t('b2c.createNew')}
          </button>
          <button className="btn btn-outline">
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
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
              />
            </svg>
            {t('b2c.export')}
          </button>
        </div>
      </div>

      {/* Статистика */}
      <div className="stats shadow mb-6 w-full">
        <div className="stat">
          <div className="stat-figure text-primary">
            <svg
              className="w-8 h-8"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
              />
            </svg>
          </div>
          <div className="stat-title">{t('b2c.stats.total')}</div>
          <div className="stat-value">{stats.total}</div>
          <div className="stat-desc">{t('b2c.stats.totalDesc')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-success">
            <svg
              className="w-8 h-8"
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
          </div>
          <div className="stat-title">{t('b2c.stats.active')}</div>
          <div className="stat-value">{stats.active}</div>
          <div className="stat-desc">{t('b2c.stats.activeDesc')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-warning">
            <svg
              className="w-8 h-8"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4M7 12l2 2 4-4m0 0a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <div className="stat-title">{t('b2c.stats.verified')}</div>
          <div className="stat-value">{stats.verified}</div>
          <div className="stat-desc">{t('b2c.stats.verifiedDesc')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-info">
            <svg
              className="w-8 h-8"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
              />
            </svg>
          </div>
          <div className="stat-title">{t('b2c.stats.products')}</div>
          <div className="stat-value">{stats.totalProducts}</div>
          <div className="stat-desc">{t('b2c.stats.productsDesc')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-secondary">
            <svg
              className="w-8 h-8"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
          </div>
          <div className="stat-title">{t('b2c.stats.sales')}</div>
          <div className="stat-value">{stats.totalSales.toLocaleString()}</div>
          <div className="stat-desc">{t('b2c.stats.salesDesc')}</div>
        </div>
      </div>

      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          {/* Панель фильтров */}
          <div className="flex justify-between items-center mb-4">
          <div className="flex gap-2 flex-1">
            <input
              type="text"
              placeholder={t('b2c.searchPlaceholder')}
              className="input input-bordered flex-1 max-w-md"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />

            <select
              className="select select-bordered"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
            >
              <option value="">{t('b2c.filters.allStatus')}</option>
              <option value="active">{t('b2c.filters.active')}</option>
              <option value="inactive">{t('b2c.filters.inactive')}</option>
            </select>

            <select
              className="select select-bordered"
              value={verifiedFilter}
              onChange={(e) => setVerifiedFilter(e.target.value)}
            >
              <option value="">{t('b2c.filters.allVerification')}</option>
              <option value="verified">{t('b2c.filters.verified')}</option>
              <option value="unverified">{t('b2c.filters.unverified')}</option>
            </select>

            <select
              className="select select-bordered"
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
            >
              <option value="date_desc">{t('b2c.sort.newest')}</option>
              <option value="date_asc">{t('b2c.sort.oldest')}</option>
              <option value="name_asc">{t('b2c.sort.nameAsc')}</option>
              <option value="name_desc">{t('b2c.sort.nameDesc')}</option>
              <option value="products_desc">
                {t('b2c.sort.mostProducts')}
              </option>
              <option value="sales_desc">{t('b2c.sort.mostSales')}</option>
              <option value="rating_desc">{t('b2c.sort.bestRating')}</option>
            </select>
          </div>

          {/* Массовые действия */}
          {selectedStorefronts.size > 0 && (
            <div className="dropdown dropdown-end">
              <label tabIndex={0} className="btn btn-primary dropdown-toggle">
                {t('b2c.bulkActions')} ({selectedStorefronts.size})
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
                    {t('b2c.activateSelected')}
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => handleBulkAction('deactivate')}
                    className="text-warning"
                  >
                    {t('b2c.deactivateSelected')}
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => handleBulkAction('soft-delete')}
                    className="text-warning"
                  >
                    {t('b2c.softDelete')} ({t('b2c.deleteSelected')})
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => handleBulkAction('hard-delete')}
                    className="text-error"
                  >
                    {t('b2c.hardDelete')} ({t('b2c.deleteSelected')})
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
                <th>{t('b2c.table.name')}</th>
                <th>{t('b2c.table.owner')}</th>
                <th>{t('b2c.table.location')}</th>
                <th>{t('b2c.table.products')}</th>
                <th>{t('b2c.table.sales')}</th>
                <th>{t('b2c.table.rating')}</th>
                <th>{t('b2c.table.status')}</th>
                <th>{t('b2c.table.verified')}</th>
                <th>{t('b2c.table.actions')}</th>
              </tr>
            </thead>
            <tbody>
              {storefronts.length === 0 ? (
                <tr>
                  <td colSpan={10} className="text-center py-8">
                    <p className="text-base-content/60">{t('b2c.noData')}</p>
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
                          ? t('b2c.active')
                          : t('b2c.inactive')}
                      </span>
                    </td>
                    <td>
                      {storefront.is_verified ? (
                        <span className="badge badge-primary">
                          {t('b2c.verified')}
                        </span>
                      ) : (
                        <span className="badge badge-outline">
                          {t('b2c.unverified')}
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
                              href={`/b2c/${storefront.slug}`}
                              target="_blank"
                              rel="noopener noreferrer"
                            >
                              {t('b2c.view')}
                            </a>
                          </li>
                          <li>
                            <button
                              onClick={() =>
                                router.push(`/admin/b2c/${storefront.id}/edit`)
                              }
                            >
                              {t('b2c.edit')}
                            </button>
                          </li>
                          <li>
                            <button
                              onClick={() =>
                                router.push(
                                  `/admin/b2c-products?storefront=${storefront.id}`
                                )
                              }
                            >
                              {t('b2c.viewProducts')}
                            </button>
                          </li>
                          {!storefront.is_active && (
                            <li>
                              <button
                                onClick={() => handleRestore(storefront.id)}
                                className="text-success"
                              >
                                {t('b2c.restore')}
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
                              {t('b2c.delete')}
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
            <h3 className="font-bold text-lg">{t('b2c.confirmDeleteTitle')}</h3>
            <p className="py-4">{t('b2c.confirmDeleteMessage')}</p>

            {/* Выбор типа удаления для администратора */}
            <div className="form-control py-4">
              <label className="label">
                <span className="label-text font-semibold">
                  {t('b2c.deleteType')}
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
                      {t('b2c.softDelete')}
                    </span>
                    <div className="label-text-alt text-warning">
                      {t('b2c.softDeleteDesc')}
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
                      {t('b2c.hardDelete')}
                    </span>
                    <div className="label-text-alt text-error">
                      {t('b2c.hardDeleteDesc')}
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
                  ? t('b2c.hardDelete')
                  : t('b2c.softDelete')}
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
    </div>
  );
}
