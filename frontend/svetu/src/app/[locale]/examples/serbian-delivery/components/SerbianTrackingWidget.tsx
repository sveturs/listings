'use client';

import { useState } from 'react';
import {
  MagnifyingGlassIcon,
  TruckIcon,
  MapPinIcon,
  CheckCircleIcon,
  ClockIcon,
} from '@heroicons/react/24/outline';

interface TrackingEvent {
  id: string;
  status: string;
  description: string;
  location: string;
  timestamp: string;
  completed: boolean;
  courier?: string;
}

export default function SerbianTrackingWidget() {
  const [trackingNumber, setTrackingNumber] = useState('BEX170123456');
  const [isLoading, setIsLoading] = useState(false);
  const [trackingData, setTrackingData] = useState<TrackingEvent[] | null>(
    null
  );

  const sampleTracking: TrackingEvent[] = [
    {
      id: '1',
      status: 'Пошиљка примљена',
      description: 'Пошиљка је примљена у BexExpress депоу',
      location: 'Београд, Централни депо',
      timestamp: '2024-08-13 14:30',
      completed: true,
      courier: 'BexExpress',
    },
    {
      id: '2',
      status: 'У транспорту',
      description: 'Пошиљка је у транспорту према одредишту',
      location: 'Београд - Нови Сад',
      timestamp: '2024-08-13 18:45',
      completed: true,
      courier: 'BexExpress',
    },
    {
      id: '3',
      status: 'Стигла у град',
      description: 'Пошиљка је стигла у одредишни град',
      location: 'Нови Сад, Централни депо',
      timestamp: '2024-08-14 08:15',
      completed: true,
      courier: 'BexExpress',
    },
    {
      id: '4',
      status: 'Спремна за доставу',
      description: 'Пошиљка је додељена куриру за финалну доставу',
      location: 'Нови Сад, Центар',
      timestamp: '2024-08-14 10:30',
      completed: false,
      courier: 'Милан Петровић (курир)',
    },
  ];

  const handleTrack = async () => {
    if (!trackingNumber.trim()) return;

    setIsLoading(true);

    // Simulate API call
    setTimeout(() => {
      setTrackingData(sampleTracking);
      setIsLoading(false);
    }, 1000);
  };

  const getStatusIcon = (status: string, completed: boolean) => {
    if (completed) {
      return <CheckCircleIcon className="w-6 h-6 text-success" />;
    }

    switch (status) {
      case 'У транспорту':
        return <TruckIcon className="w-6 h-6 text-info" />;
      case 'Спремна за доставу':
        return <MapPinIcon className="w-6 h-6 text-warning" />;
      default:
        return <ClockIcon className="w-6 h-6 text-base-content/50" />;
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Tracking Input */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">
            <TruckIcon className="w-8 h-8 text-primary" />
            Праћење пошиљке - српске курирске службе
          </h2>

          <div className="flex gap-2">
            <div className="form-control flex-1">
              <input
                type="text"
                value={trackingNumber}
                onChange={(e) => setTrackingNumber(e.target.value)}
                placeholder="Унесите број пошиљке (нпр. BEX170123456)"
                className="input input-bordered w-full"
              />
              <div className="label">
                <span className="label-text-alt">
                  Подржани формати: BEX, AKS, PE, CE, YP
                </span>
              </div>
            </div>
            <button
              onClick={handleTrack}
              disabled={isLoading || !trackingNumber.trim()}
              className={`btn btn-primary gap-2 ${isLoading ? 'loading' : ''}`}
            >
              {!isLoading && <MagnifyingGlassIcon className="w-5 h-5" />}
              Прати
            </button>
          </div>

          {/* Sample Tracking Numbers */}
          <div className="mt-4">
            <p className="text-sm text-base-content/70 mb-2">
              Примери бројева:
            </p>
            <div className="flex flex-wrap gap-2">
              {[
                'BEX170123456',
                'AKS789456123',
                'PE2024080145',
                'CE187456329',
                'YP0814567',
              ].map((sample) => (
                <button
                  key={sample}
                  onClick={() => setTrackingNumber(sample)}
                  className="btn btn-ghost btn-sm"
                >
                  {sample}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Tracking Results */}
      {trackingData && (
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold">
                Статус пошиљке #{trackingNumber}
              </h3>
              <div className="badge badge-success badge-lg">BexExpress</div>
            </div>

            {/* Current Status */}
            <div className="alert alert-info mb-6">
              <MapPinIcon className="w-6 h-6" />
              <div>
                <h4 className="font-semibold">Тренутни статус</h4>
                <p>
                  {trackingData[trackingData.length - 1]?.description} -{' '}
                  {trackingData[trackingData.length - 1]?.location}
                </p>
                <p className="text-sm opacity-75">
                  Очекивана достава: <strong>данас до 18:00</strong>
                </p>
              </div>
            </div>

            {/* Timeline */}
            <div className="space-y-4">
              <h4 className="font-semibold text-lg mb-4">Историја кретања</h4>

              {trackingData.map((event, index) => (
                <div key={event.id} className="flex gap-4">
                  <div className="flex flex-col items-center">
                    {getStatusIcon(event.status, event.completed)}
                    {index < trackingData.length - 1 && (
                      <div
                        className={`w-0.5 h-12 mt-2 ${
                          event.completed ? 'bg-success' : 'bg-base-300'
                        }`}
                      />
                    )}
                  </div>

                  <div className="flex-1 pb-6">
                    <div className="flex items-start justify-between">
                      <div>
                        <h5
                          className={`font-semibold ${
                            event.completed
                              ? 'text-base-content'
                              : 'text-warning'
                          }`}
                        >
                          {event.status}
                        </h5>
                        <p className="text-sm text-base-content/70 mt-1">
                          {event.description}
                        </p>
                        <div className="flex items-center gap-4 mt-2 text-xs text-base-content/50">
                          <span>📍 {event.location}</span>
                          <span>🕒 {event.timestamp}</span>
                          {event.courier && <span>👤 {event.courier}</span>}
                        </div>
                      </div>

                      {!event.completed && (
                        <div className="badge badge-warning badge-sm">
                          У току
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Courier Contact */}
            <div className="mt-6 p-4 bg-base-200 rounded-lg">
              <h5 className="font-semibold mb-2">📞 Контакт курира</h5>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div>
                  <strong>Курир:</strong> Милан Петровић
                </div>
                <div>
                  <strong>Телефон:</strong>
                  <button className="btn btn-ghost btn-sm ml-2">
                    066/123-456
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Courier Services Info */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {[
          {
            name: 'BexExpress',
            color: 'bg-purple-600',
            info: 'API интеграција',
            phone: '011/555-7890',
          },
          {
            name: 'Post Express',
            color: 'bg-blue-500',
            info: '500+ пунктова',
            phone: '011/555-0456',
          },
          {
            name: 'City Express',
            color: 'bg-green-500',
            info: '150+ градова',
            phone: '011/555-0789',
          },
          {
            name: 'Yettel Post',
            color: 'bg-purple-500',
            info: 'Пакомати 24/7',
            phone: '011/555-0321',
          },
        ].map((courier) => (
          <div key={courier.name} className="card bg-base-100 shadow-lg">
            <div className="card-body p-4">
              <div
                className={`w-12 h-12 ${courier.color} rounded-lg flex items-center justify-center text-white font-bold text-lg mb-3`}
              >
                {courier.name.charAt(0)}
              </div>
              <h4 className="font-semibold">{courier.name}</h4>
              <p className="text-sm text-base-content/70">{courier.info}</p>
              <p className="text-xs text-base-content/50 mt-2">
                {courier.phone}
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
