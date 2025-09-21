'use client';

import { useState, useEffect } from 'react';
import {
  MagnifyingGlassIcon,
  TruckIcon,
  MapPinIcon,
  ClockIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  PhoneIcon,
  BuildingStorefrontIcon,
  UserIcon,
  DocumentTextIcon,
} from '@heroicons/react/24/outline';
import { useTranslations } from 'next-intl';
import configManager from '@/config';

interface TrackingEvent {
  id: number;
  event_time: string;
  status: string;
  location: string;
  description: string;
  raw_data?: Record<string, any>;
}

interface ShipmentInfo {
  id: number;
  provider_id: number;
  provider_name: string;
  provider_code: string;
  tracking_number: string;
  status: string;
  sender_info: {
    name: string;
    phone?: string;
    city: string;
    address?: string;
  };
  recipient_info: {
    name: string;
    phone?: string;
    city: string;
    address?: string;
  };
  package_info: {
    weight_kg: number;
    dimensions?: {
      length_cm: number;
      width_cm: number;
      height_cm: number;
    };
    is_fragile?: boolean;
    packaging_type?: string;
  };
  delivery_cost: number;
  cod_amount?: number;
  pickup_date?: string;
  estimated_delivery?: string;
  actual_delivery_date?: string;
  created_at: string;
  events: TrackingEvent[];
}

interface Props {
  initialTrackingNumber?: string;
  className?: string;
}

const STATUS_MAPPING: Record<
  string,
  { label: string; color: string; icon: any }
> = {
  pending: {
    label: 'Ожидает отправки',
    color: 'text-warning',
    icon: ClockIcon,
  },
  pickup_scheduled: {
    label: 'Забор запланирован',
    color: 'text-info',
    icon: ClockIcon,
  },
  picked_up: {
    label: 'Забрано у отправителя',
    color: 'text-primary',
    icon: TruckIcon,
  },
  in_transit: { label: 'В пути', color: 'text-primary', icon: TruckIcon },
  out_for_delivery: {
    label: 'Передано в доставку',
    color: 'text-secondary',
    icon: TruckIcon,
  },
  delivered: {
    label: 'Доставлено',
    color: 'text-success',
    icon: CheckCircleIcon,
  },
  failed_delivery: {
    label: 'Неудачная доставка',
    color: 'text-error',
    icon: ExclamationTriangleIcon,
  },
  returned: {
    label: 'Возвращено отправителю',
    color: 'text-warning',
    icon: ExclamationTriangleIcon,
  },
  cancelled: {
    label: 'Отменено',
    color: 'text-error',
    icon: ExclamationTriangleIcon,
  },
};

export default function TrackingPage({
  initialTrackingNumber,
  className = '',
}: Props) {
  const _t = useTranslations('delivery');
  const [trackingNumber, setTrackingNumber] = useState(
    initialTrackingNumber || ''
  );
  const [shipment, setShipment] = useState<ShipmentInfo | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (initialTrackingNumber) {
      trackShipment(initialTrackingNumber);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialTrackingNumber]);

  const trackShipment = async (number?: string) => {
    const trackingNum = number || trackingNumber;
    if (!trackingNum.trim()) return;

    setLoading(true);
    setError(null);

    try {
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/shipments/track/${encodeURIComponent(trackingNum)}`
      );
      const data = await response.json();

      if (data.success && data.data) {
        setShipment(data.data);
      } else {
        setError(data.message || 'Отправление не найдено');
        setShipment(null);
      }
    } catch (err) {
      setError('Ошибка при поиске отправления');
      console.error('Tracking error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    trackShipment();
  };

  const getStatusInfo = (status: string) => {
    return (
      STATUS_MAPPING[status] || {
        label: status,
        color: 'text-base-content',
        icon: InformationCircleIcon,
      }
    );
  };

  const formatDateTime = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    });
  };

  const getProgressPercentage = () => {
    if (!shipment) return 0;

    const statusOrder = [
      'pending',
      'pickup_scheduled',
      'picked_up',
      'in_transit',
      'out_for_delivery',
      'delivered',
    ];
    const currentIndex = statusOrder.indexOf(shipment.status);

    if (currentIndex === -1) return 0;
    if (shipment.status === 'delivered') return 100;

    return Math.round(((currentIndex + 1) / statusOrder.length) * 100);
  };

  return (
    <div className={`max-w-4xl mx-auto space-y-6 ${className}`}>
      {/* Search Form */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body p-6">
          <div className="text-center mb-6">
            <div className="p-3 bg-primary/10 rounded-full inline-flex mb-3">
              <MagnifyingGlassIcon className="w-8 h-8 text-primary" />
            </div>
            <h1 className="text-2xl font-bold mb-2">
              Отслеживание отправления
            </h1>
            <p className="text-base-content/60">
              Введите трек-номер для получения информации о статусе доставки
            </p>
          </div>

          <form onSubmit={handleSubmit} className="max-w-md mx-auto">
            <div className="join w-full">
              <input
                type="text"
                className="input input-bordered join-item flex-1 focus:input-primary"
                placeholder="Трек-номер (например: PE123456789)"
                value={trackingNumber}
                onChange={(e) => setTrackingNumber(e.target.value)}
                disabled={loading}
              />
              <button
                type="submit"
                className={`btn btn-primary join-item ${loading ? 'loading' : ''}`}
                disabled={loading || !trackingNumber.trim()}
              >
                {!loading && <MagnifyingGlassIcon className="w-4 h-4" />}
                Найти
              </button>
            </div>
          </form>

          {error && (
            <div className="alert alert-error mt-4">
              <ExclamationTriangleIcon className="w-5 h-5" />
              <span>{error}</span>
            </div>
          )}
        </div>
      </div>

      {/* Shipment Information */}
      {shipment && (
        <>
          {/* Status Summary */}
          <div className="card bg-gradient-to-r from-primary/5 to-secondary/5 shadow-lg">
            <div className="card-body p-6">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h2 className="text-xl font-bold mb-2">
                    Отправление {shipment.tracking_number}
                  </h2>
                  <div className="flex items-center gap-2">
                    <img
                      src={`/images/delivery/${shipment.provider_code}-logo.png`}
                      alt={shipment.provider_name}
                      className="w-6 h-6 object-contain"
                      onError={(e) => {
                        (e.target as HTMLImageElement).style.display = 'none';
                      }}
                    />
                    <span className="text-sm text-base-content/60">
                      {shipment.provider_name}
                    </span>
                  </div>
                </div>

                <div className="text-right">
                  <div
                    className={`text-lg font-semibold ${getStatusInfo(shipment.status).color}`}
                  >
                    {getStatusInfo(shipment.status).label}
                  </div>
                  {shipment.estimated_delivery && (
                    <div className="text-sm text-base-content/60">
                      Ожидается: {formatDate(shipment.estimated_delivery)}
                    </div>
                  )}
                </div>
              </div>

              {/* Progress Bar */}
              <div className="mb-4">
                <div className="flex justify-between text-xs text-base-content/60 mb-2">
                  <span>Обработка</span>
                  <span>В пути</span>
                  <span>Доставка</span>
                  <span>Получено</span>
                </div>
                <progress
                  className="progress progress-primary w-full"
                  value={getProgressPercentage()}
                  max="100"
                ></progress>
              </div>

              <div className="grid md:grid-cols-3 gap-4">
                <div className="text-center">
                  <div className="text-2xl font-bold">
                    {shipment.delivery_cost.toFixed(0)} RSD
                  </div>
                  <div className="text-xs text-base-content/60">
                    стоимость доставки
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">
                    {shipment.package_info.weight_kg.toFixed(1)} кг
                  </div>
                  <div className="text-xs text-base-content/60">
                    вес посылки
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">
                    {shipment.events.length}
                  </div>
                  <div className="text-xs text-base-content/60">событий</div>
                </div>
              </div>
            </div>
          </div>

          {/* Route Information */}
          <div className="grid md:grid-cols-2 gap-6">
            {/* Sender */}
            <div className="card bg-base-100 shadow-lg">
              <div className="card-body p-6">
                <div className="flex items-center gap-3 mb-4">
                  <div className="p-2 bg-info/10 rounded-lg">
                    <BuildingStorefrontIcon className="w-5 h-5 text-info" />
                  </div>
                  <h3 className="font-semibold">Отправитель</h3>
                </div>

                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    <UserIcon className="w-4 h-4 text-base-content/60" />
                    <span>{shipment.sender_info.name}</span>
                  </div>
                  {shipment.sender_info.phone && (
                    <div className="flex items-center gap-2">
                      <PhoneIcon className="w-4 h-4 text-base-content/60" />
                      <span>{shipment.sender_info.phone}</span>
                    </div>
                  )}
                  <div className="flex items-center gap-2">
                    <MapPinIcon className="w-4 h-4 text-base-content/60" />
                    <span>
                      {shipment.sender_info.city}
                      {shipment.sender_info.address &&
                        `, ${shipment.sender_info.address}`}
                    </span>
                  </div>
                  {shipment.pickup_date && (
                    <div className="flex items-center gap-2">
                      <ClockIcon className="w-4 h-4 text-base-content/60" />
                      <span>Забор: {formatDate(shipment.pickup_date)}</span>
                    </div>
                  )}
                </div>
              </div>
            </div>

            {/* Recipient */}
            <div className="card bg-base-100 shadow-lg">
              <div className="card-body p-6">
                <div className="flex items-center gap-3 mb-4">
                  <div className="p-2 bg-success/10 rounded-lg">
                    <UserIcon className="w-5 h-5 text-success" />
                  </div>
                  <h3 className="font-semibold">Получатель</h3>
                </div>

                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    <UserIcon className="w-4 h-4 text-base-content/60" />
                    <span>{shipment.recipient_info.name}</span>
                  </div>
                  {shipment.recipient_info.phone && (
                    <div className="flex items-center gap-2">
                      <PhoneIcon className="w-4 h-4 text-base-content/60" />
                      <span>{shipment.recipient_info.phone}</span>
                    </div>
                  )}
                  <div className="flex items-center gap-2">
                    <MapPinIcon className="w-4 h-4 text-base-content/60" />
                    <span>
                      {shipment.recipient_info.city}
                      {shipment.recipient_info.address &&
                        `, ${shipment.recipient_info.address}`}
                    </span>
                  </div>
                  {shipment.cod_amount && shipment.cod_amount > 0 && (
                    <div className="flex items-center gap-2">
                      <DocumentTextIcon className="w-4 h-4 text-warning" />
                      <span className="text-warning">
                        Наложенный платеж: {shipment.cod_amount.toFixed(0)} RSD
                      </span>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* Package Information */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body p-6">
              <div className="flex items-center gap-3 mb-4">
                <div className="p-2 bg-warning/10 rounded-lg">
                  <DocumentTextIcon className="w-5 h-5 text-warning" />
                </div>
                <h3 className="font-semibold">Информация о посылке</h3>
              </div>

              <div className="grid md:grid-cols-3 gap-4">
                <div>
                  <div className="text-sm text-base-content/60 mb-1">Вес</div>
                  <div className="font-medium">
                    {shipment.package_info.weight_kg} кг
                  </div>
                </div>
                {shipment.package_info.dimensions && (
                  <div>
                    <div className="text-sm text-base-content/60 mb-1">
                      Размеры
                    </div>
                    <div className="font-medium">
                      {shipment.package_info.dimensions.length_cm} ×{' '}
                      {shipment.package_info.dimensions.width_cm} ×{' '}
                      {shipment.package_info.dimensions.height_cm} см
                    </div>
                  </div>
                )}
                <div>
                  <div className="text-sm text-base-content/60 mb-1">
                    Упаковка
                  </div>
                  <div className="font-medium">
                    {shipment.package_info.packaging_type || 'Стандартная'}
                    {shipment.package_info.is_fragile && (
                      <span className="ml-2 badge badge-warning badge-xs">
                        Хрупкое
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Tracking History */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body p-6">
              <h3 className="text-lg font-semibold mb-6 flex items-center gap-2">
                <ClockIcon className="w-5 h-5" />
                История перемещения
              </h3>

              <div className="space-y-4">
                {shipment.events.length > 0 ? (
                  shipment.events.map((event, index) => {
                    const StatusIcon = getStatusInfo(event.status).icon;
                    const isLatest = index === 0;

                    return (
                      <div key={event.id} className="flex gap-4">
                        <div className="flex flex-col items-center">
                          <div
                            className={`
                            p-2 rounded-full
                            ${isLatest ? 'bg-primary/20' : 'bg-base-200'}
                          `}
                          >
                            <StatusIcon
                              className={`
                              w-4 h-4
                              ${isLatest ? 'text-primary' : 'text-base-content/60'}
                            `}
                            />
                          </div>
                          {index < shipment.events.length - 1 && (
                            <div className="w-px h-8 bg-base-200 mt-2"></div>
                          )}
                        </div>

                        <div className="flex-1 pb-4">
                          <div className="flex justify-between items-start">
                            <div>
                              <div
                                className={`font-medium ${isLatest ? 'text-primary' : ''}`}
                              >
                                {getStatusInfo(event.status).label}
                              </div>
                              <div className="text-sm text-base-content/70 mt-1">
                                {event.description}
                              </div>
                              {event.location && (
                                <div className="text-xs text-base-content/60 mt-1 flex items-center gap-1">
                                  <MapPinIcon className="w-3 h-3" />
                                  {event.location}
                                </div>
                              )}
                            </div>
                            <div className="text-sm text-base-content/60 text-right">
                              {formatDateTime(event.event_time)}
                            </div>
                          </div>
                        </div>
                      </div>
                    );
                  })
                ) : (
                  <div className="text-center py-8 text-base-content/60">
                    <ClockIcon className="w-12 h-12 mx-auto mb-2 opacity-50" />
                    <p>Пока нет событий отслеживания</p>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Help Information */}
          <div className="card bg-gradient-to-r from-info/5 to-info/10">
            <div className="card-body p-6">
              <div className="flex items-start gap-3">
                <InformationCircleIcon className="w-5 h-5 text-info flex-shrink-0 mt-0.5" />
                <div className="space-y-2 text-sm">
                  <div className="font-semibold">Полезная информация:</div>
                  <ul className="list-disc list-inside space-y-1 text-xs">
                    <li>
                      Информация обновляется автоматически каждые 30 минут
                    </li>
                    <li>SMS уведомления отправляются при изменении статуса</li>
                    <li>Посылка хранится в пункте выдачи до 5 рабочих дней</li>
                    <li>
                      При возникновении вопросов обращайтесь в службу поддержки{' '}
                      {shipment.provider_name}
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
