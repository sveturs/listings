'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';

interface Shipment {
  id: number;
  tracking_number: string;
  provider: string;
  status: string;
  sender_name: string;
  recipient_name: string;
  sender_city: string;
  recipient_city: string;
  created_at: string;
  estimated_delivery: string;
  delivery_cost: number;
  weight: number;
  items_count?: number;
}

export default function DeliveryShipments() {
  const _t = useTranslations('admin.delivery');
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [filters, setFilters] = useState({
    status: '',
    provider: '',
    search: '',
  });

  const fetchShipments = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: currentPage.toString(),
        limit: '20',
        ...(filters.status && { status: filters.status }),
        ...(filters.provider && { provider: filters.provider }),
        ...(filters.search && { search: filters.search }),
      });

      const response = await apiClient.get(
        `/admin/delivery/shipments?${params}`
      );

      if (response.data) {
        setShipments(response.data.data?.shipments || []);
        const total = response.data.data?.total || 0;
        setTotalPages(Math.ceil(total / 20));
      }
    } catch (error) {
      console.error('Error fetching shipments:', error);
      // Use mock data for development
      setShipments(getMockShipments());
      setTotalPages(5);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchShipments();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentPage, filters]);

  const getMockShipments = (): Shipment[] => [
    {
      id: 1,
      tracking_number: 'PE-2024-001234',
      provider: 'Post Express',
      status: 'delivered',
      sender_name: 'Марко Петровић',
      recipient_name: 'Јована Николић',
      sender_city: 'Београд',
      recipient_city: 'Нови Сад',
      created_at: '2024-01-20T10:00:00Z',
      estimated_delivery: '2024-01-22T18:00:00Z',
      delivery_cost: 450,
      weight: 2.5,
      items_count: 3,
    },
    {
      id: 2,
      tracking_number: 'BEX-2024-005678',
      provider: 'BEX Express',
      status: 'in_transit',
      sender_name: 'Милан Јовановић',
      recipient_name: 'Ана Стојковић',
      sender_city: 'Ниш',
      recipient_city: 'Београд',
      created_at: '2024-01-21T14:30:00Z',
      estimated_delivery: '2024-01-23T12:00:00Z',
      delivery_cost: 380,
      weight: 1.8,
      items_count: 1,
    },
    {
      id: 3,
      tracking_number: 'AKS-2024-009012',
      provider: 'AKS Express',
      status: 'pending',
      sender_name: 'Драган Милић',
      recipient_name: 'Милица Павловић',
      sender_city: 'Крагујевац',
      recipient_city: 'Суботица',
      created_at: '2024-01-22T09:15:00Z',
      estimated_delivery: '2024-01-24T16:00:00Z',
      delivery_cost: 520,
      weight: 3.2,
      items_count: 2,
    },
  ];

  const getStatusBadge = (status: string) => {
    const badges: Record<string, string> = {
      pending: 'badge badge-warning',
      in_transit: 'badge badge-info',
      delivered: 'badge badge-success',
      failed: 'badge badge-error',
      cancelled: 'badge badge-ghost',
    };
    return badges[status] || 'badge';
  };

  const getStatusText = (status: string) => {
    const texts: Record<string, string> = {
      pending: 'Ожидает',
      in_transit: 'В пути',
      delivered: 'Доставлено',
      failed: 'Проблема',
      cancelled: 'Отменено',
    };
    return texts[status] || status;
  };

  const handleShipmentAction = async (id: number, action: string) => {
    try {
      const response = await apiClient.post(
        `/admin/delivery/shipments/${id}/action`,
        { action }
      );

      if (response.data) {
        fetchShipments();
      }
    } catch (error) {
      console.error('Error performing action:', error);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header and Filters */}
      <div className="bg-base-200 rounded-lg p-4">
        <div className="flex flex-wrap gap-4 items-center">
          <input
            type="text"
            placeholder="Поиск по трек-номеру..."
            className="input input-bordered"
            value={filters.search}
            onChange={(e) => setFilters({ ...filters, search: e.target.value })}
          />

          <select
            className="select select-bordered"
            value={filters.status}
            onChange={(e) => setFilters({ ...filters, status: e.target.value })}
          >
            <option value="">Все статусы</option>
            <option value="pending">Ожидает</option>
            <option value="in_transit">В пути</option>
            <option value="delivered">Доставлено</option>
            <option value="failed">Проблема</option>
          </select>

          <select
            className="select select-bordered"
            value={filters.provider}
            onChange={(e) =>
              setFilters({ ...filters, provider: e.target.value })
            }
          >
            <option value="">Все провайдеры</option>
            <option value="post_express">Post Express</option>
            <option value="bex_express">BEX Express</option>
            <option value="aks_express">AKS Express</option>
            <option value="d_express">D Express</option>
            <option value="city_express">City Express</option>
            <option value="dhl">DHL</option>
          </select>

          <button className="btn btn-primary ml-auto">Экспорт CSV</button>
        </div>
      </div>

      {/* Shipments Table */}
      <div className="bg-base-100 rounded-lg shadow-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="table table-zebra">
            <thead>
              <tr>
                <th>Трек-номер</th>
                <th>Провайдер</th>
                <th>Статус</th>
                <th>Отправитель</th>
                <th>Получатель</th>
                <th>Вес (кг)</th>
                <th>Стоимость</th>
                <th>Дата создания</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={9} className="text-center py-8">
                    <span className="loading loading-spinner loading-lg"></span>
                  </td>
                </tr>
              ) : shipments.length === 0 ? (
                <tr>
                  <td
                    colSpan={9}
                    className="text-center py-8 text-base-content/50"
                  >
                    Отправления не найдены
                  </td>
                </tr>
              ) : (
                shipments.map((shipment) => (
                  <tr key={shipment.id} className="hover">
                    <td className="font-mono text-sm">
                      {shipment.tracking_number}
                    </td>
                    <td>{shipment.provider}</td>
                    <td>
                      <span className={getStatusBadge(shipment.status)}>
                        {getStatusText(shipment.status)}
                      </span>
                    </td>
                    <td>
                      <div>
                        <div className="font-semibold">
                          {shipment.sender_name}
                        </div>
                        <div className="text-sm text-base-content/70">
                          {shipment.sender_city}
                        </div>
                      </div>
                    </td>
                    <td>
                      <div>
                        <div className="font-semibold">
                          {shipment.recipient_name}
                        </div>
                        <div className="text-sm text-base-content/70">
                          {shipment.recipient_city}
                        </div>
                      </div>
                    </td>
                    <td>{shipment.weight}</td>
                    <td>{shipment.delivery_cost} RSD</td>
                    <td>
                      {new Date(shipment.created_at).toLocaleDateString(
                        'sr-RS'
                      )}
                    </td>
                    <td>
                      <div className="dropdown dropdown-end">
                        <label tabIndex={0} className="btn btn-sm btn-ghost">
                          ⋮
                        </label>
                        <ul
                          tabIndex={0}
                          className="dropdown-content z-[1] menu p-2 shadow bg-base-200 rounded-box w-52"
                        >
                          <li>
                            <button
                              onClick={() =>
                                window.open(
                                  `/tracking/${shipment.tracking_number}`,
                                  '_blank'
                                )
                              }
                            >
                              Отследить
                            </button>
                          </li>
                          <li>
                            <button
                              onClick={() =>
                                handleShipmentAction(shipment.id, 'print_label')
                              }
                            >
                              Печать этикетки
                            </button>
                          </li>
                          {shipment.status === 'pending' && (
                            <li>
                              <button
                                onClick={() =>
                                  handleShipmentAction(shipment.id, 'confirm')
                                }
                              >
                                Подтвердить отправку
                              </button>
                            </li>
                          )}
                          {shipment.status === 'in_transit' && (
                            <li>
                              <button
                                onClick={() =>
                                  handleShipmentAction(
                                    shipment.id,
                                    'mark_delivered'
                                  )
                                }
                              >
                                Отметить доставленным
                              </button>
                            </li>
                          )}
                          <li>
                            <button
                              className="text-error"
                              onClick={() =>
                                handleShipmentAction(shipment.id, 'cancel')
                              }
                            >
                              Отменить
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

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-center items-center p-4 border-t">
            <div className="join">
              <button
                className="join-item btn"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage(currentPage - 1)}
              >
                «
              </button>
              {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                const page = currentPage > 3 ? currentPage - 2 + i : i + 1;
                if (page > totalPages) return null;
                return (
                  <button
                    key={page}
                    className={`join-item btn ${currentPage === page ? 'btn-active' : ''}`}
                    onClick={() => setCurrentPage(page)}
                  >
                    {page}
                  </button>
                );
              })}
              <button
                className="join-item btn"
                disabled={currentPage === totalPages}
                onClick={() => setCurrentPage(currentPage + 1)}
              >
                »
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="stat bg-base-200 rounded-lg">
          <div className="stat-title">Всего отправлений</div>
          <div className="stat-value text-primary">234</div>
          <div className="stat-desc">За последние 30 дней</div>
        </div>

        <div className="stat bg-base-200 rounded-lg">
          <div className="stat-title">В пути</div>
          <div className="stat-value text-info">42</div>
          <div className="stat-desc">Активные доставки</div>
        </div>

        <div className="stat bg-base-200 rounded-lg">
          <div className="stat-title">Доставлено</div>
          <div className="stat-value text-success">178</div>
          <div className="stat-desc">Успешных доставок</div>
        </div>

        <div className="stat bg-base-200 rounded-lg">
          <div className="stat-title">Проблемные</div>
          <div className="stat-value text-warning">14</div>
          <div className="stat-desc">Требуют внимания</div>
        </div>
      </div>
    </div>
  );
}
