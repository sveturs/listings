'use client';

import { useState, useEffect } from 'react';
import {
  TruckIcon,
  CogIcon,
  ChartBarIcon,
  BuildingStorefrontIcon,
  ClockIcon,
  CheckCircleIcon,
  XCircleIcon,
  ExclamationTriangleIcon,
  ArrowPathIcon,
  EyeIcon,
  PrinterIcon,
  PlusIcon,
} from '@heroicons/react/24/outline';
import { PageTransition } from '@/components/ui/PageTransition';
import { configManager } from '@/config';

interface Shipment {
  id: number;
  tracking_number: string;
  status: string;
  created_at: string;
  recipient_name: string;
  recipient_phone: string;
  delivery_address?: string;
  office_code?: string;
  office_name?: string;
  cod_amount?: number;
  weight: number;
  declared_value: number;
  order_id: number;
  order_number: string;
  estimated_delivery?: string;
  last_status_update: string;
}

interface ShipmentStats {
  total: number;
  pending: number;
  in_transit: number;
  delivered: number;
  cancelled: number;
  total_value: number;
  total_cod: number;
}

const tabs = [
  { id: 'overview', label: 'Обзор', icon: ChartBarIcon },
  { id: 'shipments', label: 'Отправления', icon: TruckIcon },
  { id: 'warehouse', label: 'Склад', icon: BuildingStorefrontIcon },
  { id: 'settings', label: 'Настройки', icon: CogIcon },
] as const;

type TabId = (typeof tabs)[number]['id'];

export default function PostExpressAdminPage() {
  // const t = useTranslations('admin');
  const [activeTab, setActiveTab] = useState<TabId>('overview');
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [stats, setStats] = useState<ShipmentStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [selectedShipments, setSelectedShipments] = useState<number[]>([]);
  const [filters, setFilters] = useState({
    status: '',
    dateFrom: '',
    dateTo: '',
    search: '',
  });

  useEffect(() => {
    loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab, filters]);

  const loadData = async () => {
    setLoading(true);
    try {
      if (activeTab === 'overview') {
        await loadStats();
      } else if (activeTab === 'shipments') {
        await loadShipments();
      }
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadStats = async () => {
    const apiUrl = configManager.get('api.url');
    const response = await fetch(`${apiUrl}/api/v1/postexpress/shipments/stats`);
    const data = await response.json();
    if (data.success) {
      setStats(data.data);
    }
  };

  const loadShipments = async () => {
    const params = new URLSearchParams();
    if (filters.status) params.append('status', filters.status);
    if (filters.dateFrom) params.append('date_from', filters.dateFrom);
    if (filters.dateTo) params.append('date_to', filters.dateTo);
    if (filters.search) params.append('search', filters.search);
    params.append('limit', '50');

    const apiUrl = configManager.get('api.url');
    const response = await fetch(`${apiUrl}/api/v1/postexpress/shipments?${params}`);
    const data = await response.json();
    if (data.success) {
      setShipments(data.data || []);
    }
  };

  const handleBulkAction = async (action: string) => {
    if (selectedShipments.length === 0) return;

    setLoading(true);
    try {
      const apiUrl = configManager.get('api.url');
      let endpoint = '';
      let method = 'POST';
      const body: any = { shipment_ids: selectedShipments };

      switch (action) {
        case 'sync':
          endpoint = `${apiUrl}/api/v1/postexpress/track/sync`;
          break;
        case 'cancel':
          endpoint = `${apiUrl}/api/v1/postexpress/shipments/bulk-cancel`;
          break;
        case 'print':
          endpoint = `${apiUrl}/api/v1/postexpress/shipments/bulk-labels`;
          method = 'GET';
          break;
      }

      const response = await fetch(endpoint, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: method === 'POST' ? JSON.stringify(body) : undefined,
      });

      const data = await response.json();
      if (data.success) {
        if (action === 'print') {
          // Handle PDF download
          window.open(data.data.pdf_url, '_blank');
        } else {
          await loadShipments();
        }
        setSelectedShipments([]);
      }
    } catch (error) {
      console.error(`Failed to ${action}:`, error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'delivered':
        return <CheckCircleIcon className="w-5 h-5 text-success" />;
      case 'in_transit':
        return <TruckIcon className="w-5 h-5 text-primary" />;
      case 'at_office':
        return <BuildingStorefrontIcon className="w-5 h-5 text-warning" />;
      case 'pending':
        return <ClockIcon className="w-5 h-5 text-info" />;
      case 'cancelled':
        return <XCircleIcon className="w-5 h-5 text-error" />;
      default:
        return (
          <ExclamationTriangleIcon className="w-5 h-5 text-base-content/60" />
        );
    }
  };

  const getStatusText = (status: string) => {
    switch (status.toLowerCase()) {
      case 'delivered':
        return 'Доставлено';
      case 'in_transit':
        return 'В пути';
      case 'at_office':
        return 'В отделении';
      case 'pending':
        return 'Ожидает отправки';
      case 'cancelled':
        return 'Отменено';
      default:
        return status;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <PageTransition>
      <div className="min-h-screen bg-base-100 pt-24">
        <div className="container mx-auto px-4 py-8">
          {/* Header */}
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
            <div>
              <h1 className="text-3xl font-bold flex items-center gap-3">
                <TruckIcon className="w-8 h-8 text-primary" />
                Управление Post Express
              </h1>
              <p className="text-base-content/70 mt-2">
                Управление отправлениями, складом и настройками доставки
              </p>
            </div>

            <div className="flex gap-2">
              <button className="btn btn-primary btn-sm">
                <PlusIcon className="w-4 h-4" />
                Создать отправление
              </button>
              <button
                className="btn btn-outline btn-sm"
                onClick={() => loadData()}
              >
                <ArrowPathIcon className="w-4 h-4" />
                Обновить
              </button>
            </div>
          </div>

          {/* Tabs */}
          <div className="tabs tabs-boxed w-full mb-8">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                className={`tab ${activeTab === tab.id ? 'tab-active' : ''}`}
                onClick={() => setActiveTab(tab.id)}
              >
                <tab.icon className="w-4 h-4 mr-2" />
                {tab.label}
              </button>
            ))}
          </div>

          {/* Content */}
          <div className="space-y-6">
            {/* Overview Tab */}
            {activeTab === 'overview' && (
              <div className="space-y-6">
                {/* Stats Cards */}
                {stats && (
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <div className="card bg-base-100 shadow-lg">
                      <div className="card-body">
                        <div className="flex items-center justify-between">
                          <div>
                            <p className="text-base-content/60 text-sm">
                              Всего отправлений
                            </p>
                            <p className="text-3xl font-bold">{stats.total}</p>
                          </div>
                          <TruckIcon className="w-12 h-12 text-primary/30" />
                        </div>
                      </div>
                    </div>

                    <div className="card bg-base-100 shadow-lg">
                      <div className="card-body">
                        <div className="flex items-center justify-between">
                          <div>
                            <p className="text-base-content/60 text-sm">
                              В пути
                            </p>
                            <p className="text-3xl font-bold text-primary">
                              {stats.in_transit}
                            </p>
                          </div>
                          <ClockIcon className="w-12 h-12 text-primary/30" />
                        </div>
                      </div>
                    </div>

                    <div className="card bg-base-100 shadow-lg">
                      <div className="card-body">
                        <div className="flex items-center justify-between">
                          <div>
                            <p className="text-base-content/60 text-sm">
                              Доставлено
                            </p>
                            <p className="text-3xl font-bold text-success">
                              {stats.delivered}
                            </p>
                          </div>
                          <CheckCircleIcon className="w-12 h-12 text-success/30" />
                        </div>
                      </div>
                    </div>

                    <div className="card bg-base-100 shadow-lg">
                      <div className="card-body">
                        <div className="flex items-center justify-between">
                          <div>
                            <p className="text-base-content/60 text-sm">
                              Общая стоимость
                            </p>
                            <p className="text-3xl font-bold">
                              {stats.total_value.toFixed(0)} RSD
                            </p>
                          </div>
                          <ChartBarIcon className="w-12 h-12 text-secondary/30" />
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {/* Recent Shipments */}
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body">
                    <h3 className="card-title mb-4">Последние отправления</h3>
                    {/* This would show recent shipments in a condensed format */}
                    <div className="text-center py-8 text-base-content/60">
                      Недавние отправления будут отображены здесь
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Shipments Tab */}
            {activeTab === 'shipments' && (
              <div className="space-y-6">
                {/* Filters */}
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body p-4">
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                      <div className="form-control">
                        <input
                          type="text"
                          placeholder="Поиск по номеру, получателю..."
                          className="input input-bordered input-sm"
                          value={filters.search}
                          onChange={(e) =>
                            setFilters((prev) => ({
                              ...prev,
                              search: e.target.value,
                            }))
                          }
                        />
                      </div>

                      <div className="form-control">
                        <select
                          className="select select-bordered select-sm"
                          value={filters.status}
                          onChange={(e) =>
                            setFilters((prev) => ({
                              ...prev,
                              status: e.target.value,
                            }))
                          }
                        >
                          <option value="">Все статусы</option>
                          <option value="pending">Ожидает отправки</option>
                          <option value="in_transit">В пути</option>
                          <option value="at_office">В отделении</option>
                          <option value="delivered">Доставлено</option>
                          <option value="cancelled">Отменено</option>
                        </select>
                      </div>

                      <div className="form-control">
                        <input
                          type="date"
                          className="input input-bordered input-sm"
                          value={filters.dateFrom}
                          onChange={(e) =>
                            setFilters((prev) => ({
                              ...prev,
                              dateFrom: e.target.value,
                            }))
                          }
                        />
                      </div>

                      <div className="form-control">
                        <input
                          type="date"
                          className="input input-bordered input-sm"
                          value={filters.dateTo}
                          onChange={(e) =>
                            setFilters((prev) => ({
                              ...prev,
                              dateTo: e.target.value,
                            }))
                          }
                        />
                      </div>
                    </div>
                  </div>
                </div>

                {/* Bulk Actions */}
                {selectedShipments.length > 0 && (
                  <div className="card bg-primary/5 border border-primary/20">
                    <div className="card-body p-4">
                      <div className="flex flex-wrap items-center justify-between gap-4">
                        <span className="text-sm">
                          Выбрано: {selectedShipments.length} отправлений
                        </span>

                        <div className="flex gap-2">
                          <button
                            className="btn btn-sm btn-outline"
                            onClick={() => handleBulkAction('sync')}
                            disabled={loading}
                          >
                            <ArrowPathIcon className="w-4 h-4" />
                            Синхронизировать
                          </button>
                          <button
                            className="btn btn-sm btn-outline"
                            onClick={() => handleBulkAction('print')}
                            disabled={loading}
                          >
                            <PrinterIcon className="w-4 h-4" />
                            Печать этикеток
                          </button>
                          <button
                            className="btn btn-sm btn-error btn-outline"
                            onClick={() => handleBulkAction('cancel')}
                            disabled={loading}
                          >
                            <XCircleIcon className="w-4 h-4" />
                            Отменить
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {/* Shipments Table */}
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body p-0">
                    <div className="overflow-x-auto">
                      <table className="table table-sm">
                        <thead>
                          <tr>
                            <th>
                              <input
                                type="checkbox"
                                className="checkbox checkbox-sm"
                                checked={
                                  selectedShipments.length ===
                                    shipments.length && shipments.length > 0
                                }
                                onChange={(e) => {
                                  if (e.target.checked) {
                                    setSelectedShipments(
                                      shipments.map((s) => s.id)
                                    );
                                  } else {
                                    setSelectedShipments([]);
                                  }
                                }}
                              />
                            </th>
                            <th>Трек-номер</th>
                            <th>Статус</th>
                            <th>Заказ</th>
                            <th>Получатель</th>
                            <th>Адрес/Отделение</th>
                            <th>COD</th>
                            <th>Создан</th>
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
                                className="text-center py-8 text-base-content/60"
                              >
                                Отправления не найдены
                              </td>
                            </tr>
                          ) : (
                            shipments.map((shipment) => (
                              <tr key={shipment.id} className="hover">
                                <td>
                                  <input
                                    type="checkbox"
                                    className="checkbox checkbox-sm"
                                    checked={selectedShipments.includes(
                                      shipment.id
                                    )}
                                    onChange={(e) => {
                                      if (e.target.checked) {
                                        setSelectedShipments((prev) => [
                                          ...prev,
                                          shipment.id,
                                        ]);
                                      } else {
                                        setSelectedShipments((prev) =>
                                          prev.filter(
                                            (id) => id !== shipment.id
                                          )
                                        );
                                      }
                                    }}
                                  />
                                </td>
                                <td>
                                  <div className="font-mono text-sm">
                                    {shipment.tracking_number}
                                  </div>
                                </td>
                                <td>
                                  <div className="flex items-center gap-2">
                                    {getStatusIcon(shipment.status)}
                                    <span className="text-sm">
                                      {getStatusText(shipment.status)}
                                    </span>
                                  </div>
                                </td>
                                <td>
                                  <div className="text-sm">
                                    <div className="font-medium">
                                      {shipment.order_number}
                                    </div>
                                    <div className="text-base-content/60">
                                      #{shipment.order_id}
                                    </div>
                                  </div>
                                </td>
                                <td>
                                  <div className="text-sm">
                                    <div className="font-medium">
                                      {shipment.recipient_name}
                                    </div>
                                    <div className="text-base-content/60">
                                      {shipment.recipient_phone}
                                    </div>
                                  </div>
                                </td>
                                <td>
                                  <div className="text-sm">
                                    {shipment.delivery_address ||
                                      shipment.office_name ||
                                      'Склад Sve Tu'}
                                    {shipment.office_code && (
                                      <div className="text-base-content/60">
                                        Код: {shipment.office_code}
                                      </div>
                                    )}
                                  </div>
                                </td>
                                <td>
                                  {shipment.cod_amount ? (
                                    <span className="text-warning font-medium">
                                      {shipment.cod_amount} RSD
                                    </span>
                                  ) : (
                                    <span className="text-base-content/40">
                                      —
                                    </span>
                                  )}
                                </td>
                                <td>
                                  <div className="text-sm">
                                    {formatDate(shipment.created_at)}
                                  </div>
                                </td>
                                <td>
                                  <div className="flex gap-1">
                                    <button className="btn btn-ghost btn-xs">
                                      <EyeIcon className="w-4 h-4" />
                                    </button>
                                    <button className="btn btn-ghost btn-xs">
                                      <PrinterIcon className="w-4 h-4" />
                                    </button>
                                  </div>
                                </td>
                              </tr>
                            ))
                          )}
                        </tbody>
                      </table>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Warehouse Tab */}
            {activeTab === 'warehouse' && (
              <div className="space-y-6">
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body">
                    <h3 className="card-title mb-4">Управление складом</h3>
                    <div className="text-center py-8 text-base-content/60">
                      Функции управления складом будут добавлены здесь
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Settings Tab */}
            {activeTab === 'settings' && (
              <div className="space-y-6">
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body">
                    <h3 className="card-title mb-4">Настройки Post Express</h3>

                    <div className="space-y-6">
                      {/* API Settings */}
                      <div>
                        <h4 className="font-semibold mb-3">API подключение</h4>
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">API Username</span>
                            </label>
                            <input
                              type="text"
                              className="input input-bordered"
                              placeholder="SVE_TU_API"
                              disabled
                            />
                          </div>
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">API Endpoint</span>
                            </label>
                            <input
                              type="text"
                              className="input input-bordered"
                              placeholder="https://wsp.postexpress.rs/api/Transakcija"
                              disabled
                            />
                          </div>
                        </div>
                      </div>

                      {/* Default Settings */}
                      <div>
                        <h4 className="font-semibold mb-3">
                          Настройки по умолчанию
                        </h4>
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">Отправитель</span>
                            </label>
                            <input
                              type="text"
                              className="input input-bordered"
                              placeholder="Sve Tu d.o.o."
                            />
                          </div>
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">
                                Адрес отправителя
                              </span>
                            </label>
                            <input
                              type="text"
                              className="input input-bordered"
                              placeholder="Микија Манојловића 53, 21000 Нови Сад"
                            />
                          </div>
                        </div>
                      </div>

                      {/* Sync Settings */}
                      <div>
                        <h4 className="font-semibold mb-3">Синхронизация</h4>
                        <div className="space-y-4">
                          <div className="form-control">
                            <label className="label cursor-pointer">
                              <span className="label-text">
                                Автоматическая синхронизация статусов
                              </span>
                              <input
                                type="checkbox"
                                className="toggle toggle-primary"
                                defaultChecked
                              />
                            </label>
                          </div>
                          <div className="form-control">
                            <label className="label">
                              <span className="label-text">
                                Интервал синхронизации (минуты)
                              </span>
                            </label>
                            <input
                              type="number"
                              className="input input-bordered w-32"
                              defaultValue={30}
                              min={5}
                              max={1440}
                            />
                          </div>
                        </div>
                      </div>

                      <div className="card-actions justify-end">
                        <button className="btn btn-primary">
                          Сохранить настройки
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </PageTransition>
  );
}
