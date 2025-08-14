'use client';

import { useState } from 'react';
import {
  MagnifyingGlassIcon,
  TruckIcon,
  ClockIcon,
  MapPinIcon,
  PhoneIcon,
  DocumentTextIcon,
  ExclamationCircleIcon,
} from '@heroicons/react/24/outline';

interface TrackingEvent {
  date: string;
  time: string;
  status: string;
  location: string;
  description: string;
  icon: React.ComponentType<any>;
  current?: boolean;
}

export default function TrackingWidget() {
  const [trackingNumber, setTrackingNumber] = useState('');
  const [isTracking, setIsTracking] = useState(false);
  const [showResults, setShowResults] = useState(false);

  const mockTrackingData = {
    number: '170123456',
    status: 'in_transit',
    statusText: 'В пути',
    sender: {
      name: 'Марко Петрович',
      city: 'Белград',
      phone: '+381 64 123 4567',
    },
    recipient: {
      name: 'Ана Йованович',
      city: 'Нови Сад',
      address: 'Пункт выдачи BEX #42',
      phone: '+381 65 987 6543',
    },
    details: {
      weight: '2.5 кг',
      dimensions: '40x30x20 см',
      value: '5,500 RSD',
      cod: '5,500 RSD',
    },
    estimatedDelivery: '29.01.2025',
    events: [
      {
        date: '29.01.2025',
        time: '14:00',
        status: 'Ожидается доставка',
        location: 'Нови Сад',
        description: 'Посылка будет доставлена сегодня',
        icon: ClockIcon,
      },
      {
        date: '29.01.2025',
        time: '08:30',
        status: 'Передана на доставку',
        location: 'Нови Сад, сортировочный центр',
        description: 'Посылка передана курьеру для доставки',
        icon: TruckIcon,
        current: true,
      },
      {
        date: '28.01.2025',
        time: '22:15',
        status: 'Прибыла в город назначения',
        location: 'Нови Сад, сортировочный центр',
        description: 'Посылка прибыла в ваш город',
        icon: MapPinIcon,
      },
      {
        date: '28.01.2025',
        time: '14:30',
        status: 'В пути',
        location: 'Между Белградом и Нови Садом',
        description: 'Посылка в пути к месту назначения',
        icon: TruckIcon,
      },
      {
        date: '27.01.2025',
        time: '16:45',
        status: 'Отправлена',
        location: 'Белград, сортировочный центр',
        description: 'Посылка отправлена в город назначения',
        icon: TruckIcon,
      },
      {
        date: '27.01.2025',
        time: '12:30',
        status: 'Принята к отправке',
        location: 'Белград, пункт приема',
        description: 'Посылка принята курьерской службой',
        icon: DocumentTextIcon,
      },
    ] as TrackingEvent[],
  };

  const handleTrack = () => {
    if (!trackingNumber) return;

    setIsTracking(true);
    setTimeout(() => {
      setIsTracking(false);
      setShowResults(true);
    }, 1500);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'delivered':
        return 'text-success';
      case 'in_transit':
        return 'text-info';
      case 'pending':
        return 'text-warning';
      case 'problem':
        return 'text-error';
      default:
        return 'text-base-content';
    }
  };

  const getProgressPercentage = () => {
    const totalSteps = 5;
    const currentStep = 3; // В пути
    return (currentStep / totalSteps) * 100;
  };

  return (
    <div className="max-w-4xl mx-auto space-y-4 sm:space-y-6">
      {/* Search Form */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body p-4 sm:p-6">
          <h3 className="card-title text-lg sm:text-xl mb-3 sm:mb-4">
            Отследить посылку
          </h3>

          <div className="flex flex-col sm:flex-row gap-2 sm:gap-3">
            <input
              type="text"
              placeholder="Номер отслеживания"
              className="input input-bordered flex-1 text-sm sm:text-base"
              value={trackingNumber}
              onChange={(e) => setTrackingNumber(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleTrack()}
            />
            <button
              className={`btn btn-primary ${isTracking ? 'loading' : ''}`}
              onClick={handleTrack}
              disabled={isTracking}
            >
              {!isTracking && (
                <MagnifyingGlassIcon className="w-4 h-4 sm:w-5 sm:h-5" />
              )}
              <span className="hidden sm:inline">Отследить</span>
              <span className="sm:hidden">Найти</span>
            </button>
          </div>

          <div className="flex gap-2 mt-2">
            <button
              className="btn btn-ghost btn-xs"
              onClick={() => {
                setTrackingNumber('170123456');
                setShowResults(false);
              }}
            >
              Пример: 170123456
            </button>
            <button
              className="btn btn-ghost btn-xs"
              onClick={() => {
                setTrackingNumber('170789012');
                setShowResults(false);
              }}
            >
              Пример: 170789012
            </button>
          </div>
        </div>
      </div>

      {/* Tracking Results */}
      {showResults && (
        <>
          {/* Status Overview */}
          <div className="card bg-gradient-to-r from-primary/10 to-secondary/10">
            <div className="card-body p-4 sm:p-6">
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 sm:gap-4">
                <div>
                  <div className="text-xs sm:text-sm opacity-60">
                    Номер отслеживания
                  </div>
                  <div className="text-xl sm:text-2xl font-bold">
                    {mockTrackingData.number}
                  </div>
                </div>
                <div
                  className={`badge badge-md sm:badge-lg ${getStatusColor(mockTrackingData.status)}`}
                >
                  <TruckIcon className="w-3 h-3 sm:w-4 sm:h-4 mr-1" />
                  {mockTrackingData.statusText}
                </div>
              </div>

              {/* Progress Bar */}
              <div className="mt-4 sm:mt-6">
                <div className="flex justify-between text-xs mb-2">
                  <span>Принята</span>
                  <span>В пути</span>
                  <span>Доставлена</span>
                </div>
                <progress
                  className="progress progress-primary w-full"
                  value={getProgressPercentage()}
                  max="100"
                />
              </div>

              {/* Estimated Delivery */}
              <div className="alert alert-info mt-3 sm:mt-4">
                <ClockIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                <div>
                  <div className="font-semibold text-sm sm:text-base">
                    Ожидаемая доставка
                  </div>
                  <div className="text-xs sm:text-sm">
                    {mockTrackingData.estimatedDelivery}
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Shipment Details */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6">
            {/* Sender & Recipient */}
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body p-4 sm:p-6">
                <h4 className="font-semibold text-base sm:text-lg mb-3 sm:mb-4">
                  Информация об отправлении
                </h4>

                <div className="space-y-3 sm:space-y-4">
                  <div>
                    <div className="text-xs sm:text-sm text-base-content/60 mb-1">
                      Отправитель
                    </div>
                    <div className="font-medium text-sm sm:text-base">
                      {mockTrackingData.sender.name}
                    </div>
                    <div className="text-xs sm:text-sm">
                      {mockTrackingData.sender.city}
                    </div>
                    <div className="text-xs sm:text-sm flex items-center gap-1 mt-1">
                      <PhoneIcon className="w-3 h-3 sm:w-4 sm:h-4" />
                      {mockTrackingData.sender.phone}
                    </div>
                  </div>

                  <div className="divider my-2"></div>

                  <div>
                    <div className="text-xs sm:text-sm text-base-content/60 mb-1">
                      Получатель
                    </div>
                    <div className="font-medium text-sm sm:text-base">
                      {mockTrackingData.recipient.name}
                    </div>
                    <div className="text-xs sm:text-sm">
                      {mockTrackingData.recipient.city}
                    </div>
                    <div className="text-xs sm:text-sm">
                      {mockTrackingData.recipient.address}
                    </div>
                    <div className="text-xs sm:text-sm flex items-center gap-1 mt-1">
                      <PhoneIcon className="w-3 h-3 sm:w-4 sm:h-4" />
                      {mockTrackingData.recipient.phone}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Package Details */}
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body p-4 sm:p-6">
                <h4 className="font-semibold text-base sm:text-lg mb-3 sm:mb-4">
                  Детали посылки
                </h4>

                <div className="space-y-2 sm:space-y-3">
                  <div className="flex justify-between">
                    <span className="text-xs sm:text-sm text-base-content/60">
                      Вес
                    </span>
                    <span className="font-medium text-sm sm:text-base">
                      {mockTrackingData.details.weight}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-xs sm:text-sm text-base-content/60">
                      Размеры
                    </span>
                    <span className="font-medium text-sm sm:text-base">
                      {mockTrackingData.details.dimensions}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-xs sm:text-sm text-base-content/60">
                      Объявленная стоимость
                    </span>
                    <span className="font-medium text-sm sm:text-base">
                      {mockTrackingData.details.value}
                    </span>
                  </div>
                  <div className="divider my-2"></div>
                  <div className="flex justify-between">
                    <span className="text-xs sm:text-sm text-base-content/60">
                      Наложенный платеж
                    </span>
                    <span className="font-bold text-sm sm:text-base text-primary">
                      {mockTrackingData.details.cod}
                    </span>
                  </div>
                </div>

                <div className="alert alert-warning mt-3 sm:mt-4">
                  <ExclamationCircleIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                  <span className="text-xs sm:text-sm">
                    Оплата при получении: {mockTrackingData.details.cod}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Tracking Timeline */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body p-4 sm:p-6">
              <h4 className="font-semibold text-base sm:text-lg mb-3 sm:mb-4">
                История перемещения
              </h4>

              <div className="space-y-3 sm:space-y-4">
                {mockTrackingData.events.map((event, index) => (
                  <div
                    key={index}
                    className={`
                      flex gap-3 sm:gap-4 p-3 sm:p-4 rounded-lg transition-all
                      ${
                        event.current
                          ? 'bg-primary/10 border-2 border-primary'
                          : 'bg-base-200 hover:bg-base-300'
                      }
                    `}
                  >
                    <div
                      className={`
                      p-2 sm:p-3 rounded-full flex-shrink-0
                      ${event.current ? 'bg-primary text-primary-content' : 'bg-base-300'}
                    `}
                    >
                      <event.icon className="w-4 h-4 sm:w-5 sm:h-5" />
                    </div>

                    <div className="flex-1">
                      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-2 sm:gap-4">
                        <div className="flex-1">
                          <div className="font-semibold text-sm sm:text-base">
                            {event.status}
                            {event.current && (
                              <span className="badge badge-primary badge-xs sm:badge-sm ml-2">
                                Текущий статус
                              </span>
                            )}
                          </div>
                          <div className="text-xs sm:text-sm text-base-content/60 mt-1">
                            {event.description}
                          </div>
                          <div className="text-xs sm:text-sm text-base-content/60 mt-2">
                            <MapPinIcon className="w-3 h-3 sm:w-4 sm:h-4 inline mr-1" />
                            {event.location}
                          </div>
                        </div>
                        <div className="sm:text-right text-xs sm:text-sm text-base-content/60">
                          <div>{event.date}</div>
                          <div>{event.time}</div>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex flex-col sm:flex-row gap-2 sm:gap-3 justify-center">
            <button className="btn btn-outline btn-sm sm:btn-md">
              <PhoneIcon className="w-4 h-4 sm:w-5 sm:h-5" />
              <span className="text-xs sm:text-base">
                Связаться с поддержкой
              </span>
            </button>
            <button className="btn btn-outline btn-sm sm:btn-md">
              <DocumentTextIcon className="w-4 h-4 sm:w-5 sm:h-5" />
              <span className="text-xs sm:text-base">Скачать квитанцию</span>
            </button>
          </div>
        </>
      )}
    </div>
  );
}
