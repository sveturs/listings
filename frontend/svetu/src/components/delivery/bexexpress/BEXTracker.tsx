'use client';

import { useState, useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  TruckIcon,
  MapPinIcon,
  ClockIcon,
  CheckCircleIcon,
  ExclamationCircleIcon,
  ArrowPathIcon,
  DocumentDuplicateIcon,
  PhoneIcon,
  MapIcon,
  ListBulletIcon,
} from '@heroicons/react/24/outline';
import dynamic from 'next/dynamic';

// Динамический импорт карты для избежания SSR проблем
const Map = dynamic(() => import('./BEXMap'), {
  ssr: false,
  loading: () => (
    <div className="h-96 bg-base-200 rounded-lg flex items-center justify-center">
      <span className="loading loading-spinner loading-lg"></span>
    </div>
  ),
});

export interface BEXTrackingEvent {
  id: number;
  status: string;
  status_text: string;
  description: string;
  location?: string;
  timestamp: string;
  is_current: boolean;
}

export interface BEXShipment {
  id: number;
  tracking_number: string;
  status: number;
  status_text: string;
  sender_name: string;
  sender_address: string;
  sender_city: string;
  recipient_name: string;
  recipient_address: string;
  recipient_city: string;
  recipient_phone?: string;
  weight_kg: number;
  total_packages: number;
  cod_amount?: number;
  insurance_amount?: number;
  registered_at?: string;
  picked_up_at?: string;
  in_transit_at?: string;
  delivered_at?: string;
  failed_at?: string;
  returned_at?: string;
  failed_reason?: string;
  current_location?: {
    city: string;
    address?: string;
    latitude?: number;
    longitude?: number;
  };
  delivery_location?: {
    city: string;
    address: string;
    latitude?: number;
    longitude?: number;
  };
  events?: BEXTrackingEvent[];
  estimated_delivery?: string;
}

interface Props {
  initialTrackingNumber?: string;
  onTrackingUpdate?: (shipment: BEXShipment) => void;
  className?: string;
}

export default function BEXTracker({
  initialTrackingNumber = '',
  onTrackingUpdate,
  className = '',
}: Props) {
  const [trackingNumber, setTrackingNumber] = useState(initialTrackingNumber);
  const [shipment, setShipment] = useState<BEXShipment | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [viewMode, setViewMode] = useState<'timeline' | 'map'>('timeline');
  const [autoRefresh, setAutoRefresh] = useState(false);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  // Загрузка информации о посылке
  const fetchShipment = async (tracking: string) => {
    if (!tracking) return;

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `/api/v1/bex/track/${encodeURIComponent(tracking)}`
      );
      const data = await response.json();

      if (data.success && data.data) {
        // Добавляем демо-данные для отображения
        const enrichedShipment: BEXShipment = {
          ...data.data,
          current_location: data.data.current_location || {
            city: 'Белград',
            address: 'Сортировочный центр BEX',
            latitude: 44.8176,
            longitude: 20.4633,
          },
          delivery_location: data.data.delivery_location || {
            city: data.data.recipient_city,
            address: data.data.recipient_address,
            latitude: 45.2671,
            longitude: 19.8335,
          },
          events: data.data.events || getDemoEvents(data.data.status),
          estimated_delivery:
            data.data.estimated_delivery ||
            new Date(Date.now() + 2 * 24 * 60 * 60 * 1000).toISOString(),
        };

        setShipment(enrichedShipment);
        onTrackingUpdate?.(enrichedShipment);
      } else {
        setError(data.message || 'Посылка не найдена');
      }
    } catch (err) {
      console.error('Tracking error:', err);
      setError('Ошибка при отслеживании посылки');
      // Используем демо-данные при ошибке
      const demoShipment = getDemoShipment(tracking);
      setShipment(demoShipment);
      onTrackingUpdate?.(demoShipment);
    } finally {
      setLoading(false);
    }
  };

  // Демо-данные для тестирования
  const getDemoEvents = (status: number): BEXTrackingEvent[] => {
    const events: BEXTrackingEvent[] = [
      {
        id: 1,
        status: '0',
        status_text: 'Зарегистрировано',
        description: 'Посылка зарегистрирована в системе',
        location: 'Нови Сад, Склад отправителя',
        timestamp: new Date(Date.now() - 48 * 60 * 60 * 1000).toISOString(),
        is_current: status === 0,
      },
    ];

    if (status >= 4) {
      events.push({
        id: 2,
        status: '4',
        status_text: 'Забрано курьером',
        description: 'Посылка забрана курьером BEX',
        location: 'Нови Сад',
        timestamp: new Date(Date.now() - 36 * 60 * 60 * 1000).toISOString(),
        is_current: status === 4,
      });
    }

    if (status >= 1) {
      events.push(
        {
          id: 3,
          status: '1',
          status_text: 'В пути',
          description: 'Посылка прибыла в сортировочный центр',
          location: 'Белград, Сортировочный центр',
          timestamp: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
          is_current: false,
        },
        {
          id: 4,
          status: '1',
          status_text: 'В пути',
          description: 'Посылка отправлена в город назначения',
          location: 'Белград → Нови Сад',
          timestamp: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
          is_current: status === 1,
        }
      );
    }

    if (status >= 2) {
      events.push({
        id: 5,
        status: '2',
        status_text: 'Доставлено',
        description: 'Посылка успешно доставлена получателю',
        location: 'Нови Сад, Адрес получателя',
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
        is_current: status === 2,
      });
    }

    if (status === 3) {
      events.push({
        id: 6,
        status: '3',
        status_text: 'Возвращено',
        description: 'Посылка возвращена отправителю',
        location: 'Нови Сад, Склад отправителя',
        timestamp: new Date().toISOString(),
        is_current: true,
      });
    }

    return events.reverse(); // Новые события сверху
  };

  const getDemoShipment = (tracking: string): BEXShipment => ({
    id: 1,
    tracking_number: tracking,
    status: 1,
    status_text: 'В пути',
    sender_name: 'Sve Tu d.o.o.',
    sender_address: 'Мике Манојловића 53',
    sender_city: 'Нови Сад',
    recipient_name: 'Петар Петровић',
    recipient_address: 'Булевар ослобођења 100',
    recipient_city: 'Нови Сад',
    recipient_phone: '+381 21 123 456',
    weight_kg: 2.5,
    total_packages: 1,
    cod_amount: 3500,
    registered_at: new Date(Date.now() - 48 * 60 * 60 * 1000).toISOString(),
    picked_up_at: new Date(Date.now() - 36 * 60 * 60 * 1000).toISOString(),
    in_transit_at: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
    current_location: {
      city: 'Белград',
      address: 'Сортировочный центр BEX',
      latitude: 44.8176,
      longitude: 20.4633,
    },
    delivery_location: {
      city: 'Нови Сад',
      address: 'Булевар ослобођења 100',
      latitude: 45.2671,
      longitude: 19.8335,
    },
    events: getDemoEvents(1),
    estimated_delivery: new Date(
      Date.now() + 24 * 60 * 60 * 1000
    ).toISOString(),
  });

  // Автообновление
  useEffect(() => {
    if (autoRefresh && trackingNumber) {
      intervalRef.current = setInterval(() => {
        fetchShipment(trackingNumber);
      }, 30000); // Каждые 30 секунд

      return () => {
        if (intervalRef.current) {
          clearInterval(intervalRef.current);
          intervalRef.current = null;
        }
      };
    }
  }, [autoRefresh, trackingNumber]);

  // Начальная загрузка
  useEffect(() => {
    if (initialTrackingNumber) {
      fetchShipment(initialTrackingNumber);
    }
  }, [initialTrackingNumber]);

  const handleSearch = () => {
    if (trackingNumber) {
      fetchShipment(trackingNumber);
    }
  };

  const getStatusColor = (status: number) => {
    switch (status) {
      case 0:
        return 'text-warning';
      case 1:
        return 'text-info';
      case 2:
        return 'text-success';
      case 3:
        return 'text-error';
      case 4:
        return 'text-primary';
      case 5:
        return 'text-base-content/50';
      default:
        return 'text-base-content';
    }
  };

  const getStatusIcon = (status: number) => {
    switch (status) {
      case 0:
        return ClockIcon;
      case 1:
        return TruckIcon;
      case 2:
        return CheckCircleIcon;
      case 3:
        return ExclamationCircleIcon;
      case 4:
        return MapPinIcon;
      case 5:
        return ExclamationCircleIcon;
      default:
        return ClockIcon;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Поиск */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body p-6">
          <h3 className="font-semibold text-lg mb-4 flex items-center gap-2">
            <TruckIcon className="w-5 h-5 text-primary" />
            Отслеживание посылки BEX Express
          </h3>

          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex-1">
              <input
                type="text"
                className="input input-bordered w-full"
                placeholder="Введите трек-номер BEX (например: BEX123456789RS)"
                value={trackingNumber}
                onChange={(e) => setTrackingNumber(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              />
            </div>

            <div className="flex gap-2">
              <button
                className={`btn btn-primary ${loading ? 'loading' : ''}`}
                onClick={handleSearch}
                disabled={loading || !trackingNumber}
              >
                {!loading && <MapPinIcon className="w-5 h-5" />}
                Отследить
              </button>

              {trackingNumber && (
                <button
                  className="btn btn-ghost"
                  onClick={() => navigator.clipboard.writeText(trackingNumber)}
                  title="Скопировать трек-номер"
                >
                  <DocumentDuplicateIcon className="w-5 h-5" />
                </button>
              )}
            </div>
          </div>

          {error && (
            <div className="alert alert-error mt-4">
              <ExclamationCircleIcon className="w-5 h-5" />
              <span>{error}</span>
            </div>
          )}
        </div>
      </div>

      {/* Результаты отслеживания */}
      <AnimatePresence mode="wait">
        {shipment && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="space-y-6"
          >
            {/* Основная информация */}
            <div className="card bg-base-100 shadow-lg">
              <div className="card-body p-6">
                <div className="flex flex-col sm:flex-row justify-between items-start gap-4">
                  <div>
                    <h3 className="text-xl font-bold mb-2">
                      Посылка {shipment.tracking_number}
                    </h3>
                    <div
                      className={`flex items-center gap-2 ${getStatusColor(shipment.status)}`}
                    >
                      {(() => {
                        const Icon = getStatusIcon(shipment.status);
                        return <Icon className="w-5 h-5" />;
                      })()}
                      <span className="font-semibold text-lg">
                        {shipment.status_text}
                      </span>
                    </div>
                  </div>

                  <div className="flex gap-2">
                    <label className="label cursor-pointer gap-2">
                      <input
                        type="checkbox"
                        className="checkbox checkbox-sm"
                        checked={autoRefresh}
                        onChange={(e) => setAutoRefresh(e.target.checked)}
                      />
                      <span className="label-text">Автообновление</span>
                    </label>

                    <button
                      className="btn btn-ghost btn-sm"
                      onClick={() => fetchShipment(trackingNumber)}
                    >
                      <ArrowPathIcon className="w-4 h-4" />
                    </button>
                  </div>
                </div>

                {/* Прогресс доставки */}
                <div className="mt-6">
                  <ul className="steps steps-vertical sm:steps-horizontal w-full">
                    <li
                      className={`step ${shipment.status >= 0 ? 'step-primary' : ''}`}
                    >
                      <span className="hidden sm:inline">Зарегистрировано</span>
                    </li>
                    <li
                      className={`step ${shipment.status >= 4 ? 'step-primary' : ''}`}
                    >
                      <span className="hidden sm:inline">Забрано</span>
                    </li>
                    <li
                      className={`step ${shipment.status >= 1 ? 'step-primary' : ''}`}
                    >
                      <span className="hidden sm:inline">В пути</span>
                    </li>
                    <li
                      className={`step ${shipment.status >= 2 ? 'step-primary' : ''}`}
                    >
                      <span className="hidden sm:inline">Доставлено</span>
                    </li>
                  </ul>
                </div>

                {/* Детали отправления */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">
                  {/* Отправитель */}
                  <div>
                    <h4 className="font-semibold mb-3">Отправитель</h4>
                    <div className="space-y-2 text-sm">
                      <div>{shipment.sender_name}</div>
                      <div className="flex items-start gap-2">
                        <MapPinIcon className="w-4 h-4 text-base-content/60 mt-0.5 flex-shrink-0" />
                        <div>
                          <div>{shipment.sender_address}</div>
                          <div className="text-base-content/60">
                            {shipment.sender_city}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Получатель */}
                  <div>
                    <h4 className="font-semibold mb-3">Получатель</h4>
                    <div className="space-y-2 text-sm">
                      <div>{shipment.recipient_name}</div>
                      <div className="flex items-start gap-2">
                        <MapPinIcon className="w-4 h-4 text-base-content/60 mt-0.5 flex-shrink-0" />
                        <div>
                          <div>{shipment.recipient_address}</div>
                          <div className="text-base-content/60">
                            {shipment.recipient_city}
                          </div>
                        </div>
                      </div>
                      {shipment.recipient_phone && (
                        <div className="flex items-center gap-2">
                          <PhoneIcon className="w-4 h-4 text-base-content/60" />
                          <span>{shipment.recipient_phone}</span>
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Параметры посылки */}
                  <div>
                    <h4 className="font-semibold mb-3">Параметры</h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-base-content/60">Вес:</span>
                        <span>{shipment.weight_kg} кг</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-base-content/60">
                          Количество мест:
                        </span>
                        <span>{shipment.total_packages}</span>
                      </div>
                      {shipment.cod_amount && shipment.cod_amount > 0 && (
                        <div className="flex justify-between">
                          <span className="text-base-content/60">
                            Наложенный платеж:
                          </span>
                          <span className="font-medium">
                            {shipment.cod_amount} RSD
                          </span>
                        </div>
                      )}
                      {shipment.insurance_amount &&
                        shipment.insurance_amount > 0 && (
                          <div className="flex justify-between">
                            <span className="text-base-content/60">
                              Страховка:
                            </span>
                            <span>{shipment.insurance_amount} RSD</span>
                          </div>
                        )}
                    </div>
                  </div>

                  {/* Даты */}
                  <div>
                    <h4 className="font-semibold mb-3">Временные метки</h4>
                    <div className="space-y-2 text-sm">
                      {shipment.registered_at && (
                        <div className="flex justify-between">
                          <span className="text-base-content/60">
                            Зарегистрировано:
                          </span>
                          <span>{formatDate(shipment.registered_at)}</span>
                        </div>
                      )}
                      {shipment.picked_up_at && (
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Забрано:</span>
                          <span>{formatDate(shipment.picked_up_at)}</span>
                        </div>
                      )}
                      {shipment.delivered_at && (
                        <div className="flex justify-between">
                          <span className="text-base-content/60">
                            Доставлено:
                          </span>
                          <span className="text-success">
                            {formatDate(shipment.delivered_at)}
                          </span>
                        </div>
                      )}
                      {shipment.estimated_delivery && shipment.status < 2 && (
                        <div className="flex justify-between">
                          <span className="text-base-content/60">
                            Ожидаемая доставка:
                          </span>
                          <span className="font-medium text-primary">
                            {formatDate(shipment.estimated_delivery)}
                          </span>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Переключатель вида */}
            <div className="flex justify-center">
              <div className="btn-group">
                <button
                  className={`btn btn-sm ${viewMode === 'timeline' ? 'btn-active' : ''}`}
                  onClick={() => setViewMode('timeline')}
                >
                  <ListBulletIcon className="w-4 h-4" />
                  История
                </button>
                <button
                  className={`btn btn-sm ${viewMode === 'map' ? 'btn-active' : ''}`}
                  onClick={() => setViewMode('map')}
                >
                  <MapIcon className="w-4 h-4" />
                  Карта
                </button>
              </div>
            </div>

            {/* История событий или карта */}
            {viewMode === 'timeline' ? (
              <div className="card bg-base-100 shadow-lg">
                <div className="card-body p-6">
                  <h3 className="font-semibold text-lg mb-4">
                    История движения
                  </h3>

                  {shipment.events && shipment.events.length > 0 ? (
                    <div className="space-y-4">
                      {shipment.events.map((event, index) => (
                        <motion.div
                          key={event.id}
                          initial={{ opacity: 0, x: -20 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: index * 0.1 }}
                          className={`
                            flex gap-4 p-4 rounded-lg
                            ${
                              event.is_current
                                ? 'bg-primary/10 border-2 border-primary'
                                : 'bg-base-200'
                            }
                          `}
                        >
                          <div
                            className={`
                            p-2 rounded-full flex-shrink-0
                            ${event.is_current ? 'bg-primary text-primary-content' : 'bg-base-300'}
                          `}
                          >
                            {(() => {
                              const Icon = getStatusIcon(
                                parseInt(event.status)
                              );
                              return <Icon className="w-5 h-5" />;
                            })()}
                          </div>

                          <div className="flex-1">
                            <div className="flex items-start justify-between">
                              <div>
                                <h4 className="font-semibold">
                                  {event.status_text}
                                </h4>
                                <p className="text-sm text-base-content/70 mt-1">
                                  {event.description}
                                </p>
                                {event.location && (
                                  <div className="flex items-center gap-2 mt-2">
                                    <MapPinIcon className="w-4 h-4 text-base-content/60" />
                                    <span className="text-sm">
                                      {event.location}
                                    </span>
                                  </div>
                                )}
                              </div>
                              <div className="text-sm text-base-content/60 text-right">
                                {formatDate(event.timestamp)}
                              </div>
                            </div>
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8 text-base-content/60">
                      История событий пока недоступна
                    </div>
                  )}
                </div>
              </div>
            ) : (
              <div className="card bg-base-100 shadow-lg">
                <div className="card-body p-6">
                  <h3 className="font-semibold text-lg mb-4">
                    Маршрут доставки
                  </h3>
                  <Map
                    currentLocation={shipment.current_location}
                    deliveryLocation={shipment.delivery_location}
                    events={shipment.events}
                  />

                  {/* Легенда карты */}
                  <div className="flex flex-wrap gap-4 mt-4 text-sm">
                    <div className="flex items-center gap-2">
                      <div className="w-3 h-3 rounded-full bg-primary"></div>
                      <span>Текущее местоположение</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="w-3 h-3 rounded-full bg-success"></div>
                      <span>Пункт назначения</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="w-3 h-3 rounded-full bg-base-300"></div>
                      <span>Пройденный маршрут</span>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
