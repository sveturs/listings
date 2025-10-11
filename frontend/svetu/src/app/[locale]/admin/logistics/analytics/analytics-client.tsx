'use client';

import { useState, useEffect, useCallback } from 'react';
import {
  FiTrendingUp,
  FiTrendingDown,
  FiDollarSign,
  FiTruck,
  FiClock,
  FiTarget,
  FiDownload,
  FiMap,
} from 'react-icons/fi';
import { apiClientAuth } from '@/lib/api-client-auth';
import DeliveryMap from '../components/DeliveryMap';

interface PerformanceMetrics {
  avg_delivery_time: number;
  success_rate: number;
  total_shipments: number;
  problem_rate: number;
  return_rate: number;
  on_time_delivery_rate: number;
}

interface FinancialReport {
  total_shipping_cost: number;
  total_cod_collected: number;
  avg_shipping_cost: number;
  revenue_per_shipment: number;
  cost_savings: number;
  return_costs: number;
}

interface CourierComparison {
  name: string;
  shipments: number;
  success_rate: number;
  avg_delivery_time: number;
  cost_per_shipment: number;
  problems: number;
}

interface MapMarker {
  id: string;
  lat: number;
  lng: number;
  type: 'sender' | 'receiver' | 'in_transit' | 'delivered' | 'problem';
  shipment_id: number;
  tracking_number: string;
  status: string;
  city: string;
  address: string;
  delivery_time?: string;
  problem_description?: string;
}

export default function AnalyticsClient() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [performanceMetrics, setPerformanceMetrics] =
    useState<PerformanceMetrics | null>(null);
  const [financialReport, setFinancialReport] =
    useState<FinancialReport | null>(null);
  const [courierComparison, setCourierComparison] = useState<
    CourierComparison[]
  >([]);
  const [mapMarkers, setMapMarkers] = useState<MapMarker[]>([]);
  const [selectedPeriod, setSelectedPeriod] = useState('30'); // days
  const [exportLoading, setExportLoading] = useState(false);
  const [showMap, setShowMap] = useState(false);

  const loadAnalytics = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      // Всегда загружаем данные для карты (тестовые данные)
      loadMapData();

      // Пытаемся загрузить аналитические данные (могут не работать)
      try {
        const [performanceRes, financialRes, courierRes] = await Promise.all([
          apiClientAuth.get(
            `/admin/logistics/analytics/performance?period=${selectedPeriod}`
          ),
          apiClientAuth.get(
            `/admin/logistics/analytics/financial?period=${selectedPeriod}`
          ),
          apiClientAuth.get(
            `/admin/logistics/analytics/couriers?period=${selectedPeriod}`
          ),
        ]);

        if (performanceRes.data) {
          setPerformanceMetrics(performanceRes.data);
        }

        if (financialRes.data) {
          setFinancialReport(financialRes.data);
        }

        if (courierRes.data) {
          setCourierComparison(courierRes.data || []);
        }
      } catch (apiErr) {
        console.warn('Analytics API endpoints not implemented yet:', apiErr);
        // Не считаем это критической ошибкой, показываем только карту
      }
    } catch (err) {
      console.error('Error loading analytics:', err);
      setError('Ошибка загрузки аналитики');
    } finally {
      setLoading(false);
    }
  }, [selectedPeriod]);

  useEffect(() => {
    loadAnalytics();
  }, [loadAnalytics]);

  const loadMapData = () => {
    // Тестовые данные для карты (города Сербии)
    const testMarkers: MapMarker[] = [
      {
        id: '1',
        lat: 44.7866,
        lng: 20.4489, // Белград
        type: 'in_transit',
        shipment_id: 1001,
        tracking_number: 'BEX2024110001',
        status: 'in_transit',
        city: 'Белград',
        address: 'Кнез Михаилова 5',
      },
      {
        id: '2',
        lat: 45.2671,
        lng: 19.8335, // Нови Сад
        type: 'delivered',
        shipment_id: 1002,
        tracking_number: 'BEX2024110002',
        status: 'delivered',
        city: 'Нови Сад',
        address: 'Змај Јовина 15',
        delivery_time: '2 часа назад',
      },
      {
        id: '3',
        lat: 43.3209,
        lng: 21.8958, // Ниш
        type: 'problem',
        shipment_id: 1003,
        tracking_number: 'PE2024110003',
        status: 'problem',
        city: 'Ниш',
        address: 'Обренова 25',
        problem_description: 'Неверный адрес доставки',
      },
      {
        id: '4',
        lat: 44.0165,
        lng: 21.0059, // Крагуевац
        type: 'in_transit',
        shipment_id: 1004,
        tracking_number: 'BEX2024110004',
        status: 'in_transit',
        city: 'Крагуевац',
        address: 'Светозара Марковића 36',
      },
      {
        id: '5',
        lat: 44.0147,
        lng: 20.9106, // Смедерево
        type: 'delivered',
        shipment_id: 1005,
        tracking_number: 'PE2024110005',
        status: 'delivered',
        city: 'Смедерево',
        address: 'Карађорђева 8',
        delivery_time: '1 день назад',
      },
      {
        id: '6',
        lat: 44.2206,
        lng: 21.4031, // Пожаревац
        type: 'problem',
        shipment_id: 1006,
        tracking_number: 'BEX2024110006',
        status: 'problem',
        city: 'Пожаревац',
        address: 'Дубравска 12',
        problem_description: 'Клиент не отвечает на звонки',
      },
    ];

    setMapMarkers(testMarkers);
  };

  const handleExportReport = async () => {
    try {
      setExportLoading(true);

      // Используем BFF proxy для скачивания файла
      const url = `/api/v2/admin/logistics/analytics/export?period=${selectedPeriod}&format=csv&report_type=performance`;

      const response = await fetch(url, {
        method: 'GET',
        credentials: 'include', // Включаем cookies для авторизации через BFF
      });

      if (response.ok) {
        const blob = await response.blob();
        const downloadUrl = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = `logistics-report-${selectedPeriod}d.csv`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(downloadUrl);
      } else {
        throw new Error('Ошибка при скачивании файла');
      }
    } catch (err) {
      console.error('Error exporting report:', err);
      setError('Ошибка экспорта отчета');
    } finally {
      setExportLoading(false);
    }
  };

  const formatDuration = (hours: number) => {
    if (hours < 24) {
      return `${hours.toFixed(1)} ч`;
    } else {
      const days = Math.floor(hours / 24);
      const remainingHours = hours % 24;
      return remainingHours > 0
        ? `${days}д ${remainingHours.toFixed(1)}ч`
        : `${days}д`;
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('sr-RS', {
      style: 'currency',
      currency: 'RSD',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <span>{error}</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Фильтры и экспорт */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <div className="flex justify-between items-center">
            <div className="flex items-center gap-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Период:</span>
                </label>
                <select
                  className="select select-bordered"
                  value={selectedPeriod}
                  onChange={(e) => setSelectedPeriod(e.target.value)}
                >
                  <option value="7">Последние 7 дней</option>
                  <option value="30">Последние 30 дней</option>
                  <option value="90">Последние 90 дней</option>
                  <option value="365">Последний год</option>
                </select>
              </div>
            </div>

            <div className="flex gap-2">
              <button
                className="btn btn-primary"
                onClick={handleExportReport}
                disabled={exportLoading}
              >
                {exportLoading && (
                  <span className="loading loading-spinner loading-sm"></span>
                )}
                <FiDownload />
                Экспорт CSV
              </button>

              <button
                className={`btn ${showMap ? 'btn-secondary' : 'btn-outline'}`}
                onClick={() => setShowMap(!showMap)}
              >
                <FiMap />
                {showMap ? 'Скрыть карту' : 'Показать карту'}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Метрики производительности */}
      {performanceMetrics && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h2 className="card-title mb-4">
              <FiTarget className="w-6 h-6" />
              Показатели производительности
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-figure text-primary">
                  <FiClock className="w-8 h-8" />
                </div>
                <div className="stat-title">Среднее время доставки</div>
                <div className="stat-value text-primary">
                  {formatDuration(performanceMetrics.avg_delivery_time)}
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-figure text-success">
                  <FiTrendingUp className="w-8 h-8" />
                </div>
                <div className="stat-title">Процент успешных доставок</div>
                <div className="stat-value text-success">
                  {(performanceMetrics.success_rate * 100).toFixed(1)}%
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-figure text-info">
                  <FiTruck className="w-8 h-8" />
                </div>
                <div className="stat-title">Всего отправлений</div>
                <div className="stat-value text-info">
                  {performanceMetrics.total_shipments.toLocaleString()}
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-figure text-warning">
                  <FiTrendingDown className="w-8 h-8" />
                </div>
                <div className="stat-title">Процент проблемных</div>
                <div className="stat-value text-warning">
                  {(performanceMetrics.problem_rate * 100).toFixed(1)}%
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-figure text-error">
                  <FiTrendingDown className="w-8 h-8" />
                </div>
                <div className="stat-title">Процент возвратов</div>
                <div className="stat-value text-error">
                  {(performanceMetrics.return_rate * 100).toFixed(1)}%
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-figure text-success">
                  <FiTarget className="w-8 h-8" />
                </div>
                <div className="stat-title">Доставка в срок</div>
                <div className="stat-value text-success">
                  {(performanceMetrics.on_time_delivery_rate * 100).toFixed(1)}%
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Финансовые показатели */}
      {financialReport && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h2 className="card-title mb-4">
              <FiDollarSign className="w-6 h-6" />
              Финансовые показатели
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-title">Общая стоимость доставок</div>
                <div className="stat-value text-primary">
                  {formatCurrency(financialReport.total_shipping_cost)}
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-title">Собрано COD платежей</div>
                <div className="stat-value text-success">
                  {formatCurrency(financialReport.total_cod_collected)}
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-title">Средняя стоимость доставки</div>
                <div className="stat-value text-info">
                  {formatCurrency(financialReport.avg_shipping_cost)}
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-title">Доход с отправления</div>
                <div className="stat-value text-success">
                  {formatCurrency(financialReport.revenue_per_shipment)}
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-title">Экономия</div>
                <div className="stat-value text-success">
                  {formatCurrency(financialReport.cost_savings)}
                </div>
              </div>

              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-title">Затраты на возвраты</div>
                <div className="stat-value text-error">
                  {formatCurrency(financialReport.return_costs)}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Сравнение курьерских служб */}
      {courierComparison.length > 0 && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h2 className="card-title mb-4">
              <FiTruck className="w-6 h-6" />
              Сравнение курьерских служб
            </h2>

            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>Курьерская служба</th>
                    <th>Отправлений</th>
                    <th>Успешность</th>
                    <th>Среднее время</th>
                    <th>Стоимость за отправление</th>
                    <th>Проблем</th>
                    <th>Рейтинг</th>
                  </tr>
                </thead>
                <tbody>
                  {courierComparison.map((courier, index) => (
                    <tr key={courier.name}>
                      <td>
                        <div className="font-medium">{courier.name}</div>
                      </td>
                      <td>
                        <div className="font-mono">
                          {courier.shipments.toLocaleString()}
                        </div>
                      </td>
                      <td>
                        <div
                          className={`badge ${
                            courier.success_rate > 0.95
                              ? 'badge-success'
                              : courier.success_rate > 0.9
                                ? 'badge-warning'
                                : 'badge-error'
                          }`}
                        >
                          {(courier.success_rate * 100).toFixed(1)}%
                        </div>
                      </td>
                      <td>
                        <div className="font-mono">
                          {formatDuration(courier.avg_delivery_time)}
                        </div>
                      </td>
                      <td>
                        <div className="font-mono">
                          {formatCurrency(courier.cost_per_shipment)}
                        </div>
                      </td>
                      <td>
                        <div
                          className={`badge ${
                            courier.problems === 0
                              ? 'badge-success'
                              : courier.problems < 5
                                ? 'badge-warning'
                                : 'badge-error'
                          }`}
                        >
                          {courier.problems}
                        </div>
                      </td>
                      <td>
                        <div className="rating rating-sm">
                          {[1, 2, 3, 4, 5].map((star) => (
                            <input
                              key={star}
                              type="radio"
                              name={`rating-${index}`}
                              className={`mask mask-star-2 ${
                                star <= Math.round(courier.success_rate * 5)
                                  ? 'bg-orange-400'
                                  : 'bg-gray-300'
                              }`}
                              disabled
                            />
                          ))}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* Интерактивная карта */}
      {showMap && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h2 className="card-title mb-4">
              <FiMap className="w-6 h-6" />
              Карта активных доставок
            </h2>

            <DeliveryMap
              markers={mapMarkers}
              height="500px"
              onMarkerClick={(marker) => {
                console.log('Clicked marker:', marker);
                // Здесь можно добавить логику для отображения детальной информации
                // или перехода на страницу отправления
              }}
            />
          </div>
        </div>
      )}
    </div>
  );
}
