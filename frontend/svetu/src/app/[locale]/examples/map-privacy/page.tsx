'use client';

import React, { useState } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

const MapPrivacy = () => {
  const [privacyLevel, setPrivacyLevel] = useState<'exact' | 'area' | 'city'>(
    'area'
  );
  const [showMyLocation, setShowMyLocation] = useState(true);
  const [selectedListing, setSelectedListing] = useState<any>(null);
  const [mapView, setMapView] = useState<'map' | 'satellite' | 'hybrid'>('map');
  const [showFilters, setShowFilters] = useState(true);

  const listings = [
    {
      id: 1,
      title: 'iPhone 14 Pro Max',
      price: 899,
      exactLocation: { lat: 44.8125, lng: 20.4612 },
      areaLocation: { lat: 44.81, lng: 20.46, radius: 500 },
      city: 'Белград',
      address: 'Кнез Михаилова 25',
      privacyLevel: 'area',
      image:
        '/api/minio/download?fileName=listings/0a47e66f-d8da-459f-a2ba-8e2b85ae0163/38ad29e6-7b07-4bfc-9db2-d965cb6b966f.jpg',
    },
    {
      id: 2,
      title: 'Квартира 2-комнатная',
      price: 650,
      exactLocation: { lat: 45.2671, lng: 19.8335 },
      areaLocation: { lat: 45.27, lng: 19.83, radius: 800 },
      city: 'Нови Сад',
      address: 'Данила Киша 15',
      privacyLevel: 'exact',
      image:
        '/api/minio/download?fileName=listings/0c1fc30d-5d84-485f-a86a-5c5dc37f8b97/4b8b8e48-ddd8-4c04-ad8e-00c4b4d10d26.jpg',
    },
    {
      id: 3,
      title: 'MacBook Pro M2',
      price: 1299,
      exactLocation: { lat: 43.3209, lng: 21.8954 },
      areaLocation: { lat: 43.32, lng: 21.89, radius: 1000 },
      city: 'Ниш',
      address: 'Обреновићева 8',
      privacyLevel: 'city',
      image:
        '/api/minio/download?fileName=listings/0c91d2f7-53f7-4bff-87fe-d7e82dc3e2f0/3b26f07f-c5d6-4ff7-ba56-06ec69bb7f4d.jpg',
    },
  ];

  const privacyOptions = [
    {
      level: 'exact',
      title: 'Точный адрес',
      description: 'Показывать точное местоположение',
      icon: '📍',
      color: 'error',
    },
    {
      level: 'area',
      title: 'Район',
      description: 'Показывать только район (±500м)',
      icon: '🏘️',
      color: 'warning',
    },
    {
      level: 'city',
      title: 'Только город',
      description: 'Скрыть адрес, показать город',
      icon: '🏙️',
      color: 'success',
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      {/* Header */}
      <div className="navbar bg-base-100 shadow-lg">
        <div className="navbar-start">
          <SveTuLogoStatic variant="gradient" width={120} height={40} />
        </div>
        <div className="navbar-center">
          <h1 className="text-xl font-bold">
            🗺️ Карта с настройками приватности
          </h1>
        </div>
        <div className="navbar-end">
          <div className="btn-group">
            <button
              className={`btn btn-sm ${mapView === 'map' ? 'btn-active' : ''}`}
              onClick={() => setMapView('map')}
            >
              Карта
            </button>
            <button
              className={`btn btn-sm ${mapView === 'satellite' ? 'btn-active' : ''}`}
              onClick={() => setMapView('satellite')}
            >
              Спутник
            </button>
            <button
              className={`btn btn-sm ${mapView === 'hybrid' ? 'btn-active' : ''}`}
              onClick={() => setMapView('hybrid')}
            >
              Гибрид
            </button>
          </div>
        </div>
      </div>

      <div className="flex h-[calc(100vh-64px)]">
        {/* Map Container */}
        <div className="flex-1 relative">
          {/* Mock Map */}
          <div className="w-full h-full bg-gradient-to-br from-blue-50 to-green-50 relative overflow-hidden">
            {/* Map Grid */}
            <div className="absolute inset-0 opacity-10">
              {[...Array(20)].map((_, i) => (
                <div
                  key={i}
                  className="absolute w-full border-t border-gray-400"
                  style={{ top: `${i * 5}%` }}
                ></div>
              ))}
              {[...Array(20)].map((_, i) => (
                <div
                  key={i}
                  className="absolute h-full border-l border-gray-400"
                  style={{ left: `${i * 5}%` }}
                ></div>
              ))}
            </div>

            {/* Map Markers */}
            {listings.map((listing, idx) => {
              const position = {
                left: `${20 + idx * 25}%`,
                top: `${30 + idx * 15}%`,
              };

              return (
                <div key={listing.id} className="absolute" style={position}>
                  {listing.privacyLevel === 'exact' && (
                    <div
                      className="cursor-pointer transform hover:scale-110 transition-transform"
                      onClick={() => setSelectedListing(listing)}
                    >
                      <div className="relative">
                        <div className="absolute -inset-1 bg-red-500 rounded-full animate-ping opacity-75"></div>
                        <div className="relative bg-red-500 text-white rounded-full w-10 h-10 flex items-center justify-center shadow-lg">
                          📍
                        </div>
                      </div>
                      <div className="absolute top-12 left-1/2 transform -translate-x-1/2 whitespace-nowrap">
                        <div className="bg-white rounded-lg shadow-lg px-3 py-1 text-sm">
                          <div className="font-semibold">{listing.title}</div>
                          <div className="text-xs text-gray-500">
                            {listing.address}
                          </div>
                        </div>
                      </div>
                    </div>
                  )}

                  {listing.privacyLevel === 'area' && (
                    <div
                      className="cursor-pointer"
                      onClick={() => setSelectedListing(listing)}
                    >
                      <div className="relative">
                        <div className="absolute -inset-16 bg-orange-200 rounded-full opacity-30"></div>
                        <div className="absolute -inset-8 bg-orange-300 rounded-full opacity-30"></div>
                        <div className="relative bg-orange-500 text-white rounded-full w-10 h-10 flex items-center justify-center shadow-lg z-10">
                          🏘️
                        </div>
                      </div>
                      <div className="absolute top-12 left-1/2 transform -translate-x-1/2 whitespace-nowrap">
                        <div className="bg-white rounded-lg shadow-lg px-3 py-1 text-sm">
                          <div className="font-semibold">{listing.title}</div>
                          <div className="text-xs text-gray-500">
                            Район ±500м
                          </div>
                        </div>
                      </div>
                    </div>
                  )}

                  {listing.privacyLevel === 'city' && (
                    <div
                      className="cursor-pointer"
                      onClick={() => setSelectedListing(listing)}
                    >
                      <div className="relative">
                        <div className="absolute -inset-24 bg-green-200 rounded-full opacity-20"></div>
                        <div className="absolute -inset-16 bg-green-300 rounded-full opacity-20"></div>
                        <div className="absolute -inset-8 bg-green-400 rounded-full opacity-20"></div>
                        <div className="relative bg-green-500 text-white rounded-full w-10 h-10 flex items-center justify-center shadow-lg z-10">
                          🏙️
                        </div>
                      </div>
                      <div className="absolute top-12 left-1/2 transform -translate-x-1/2 whitespace-nowrap">
                        <div className="bg-white rounded-lg shadow-lg px-3 py-1 text-sm">
                          <div className="font-semibold">{listing.title}</div>
                          <div className="text-xs text-gray-500">
                            {listing.city}
                          </div>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}

            {/* User Location */}
            {showMyLocation && (
              <div className="absolute" style={{ left: '50%', top: '50%' }}>
                <div className="relative">
                  <div className="absolute -inset-2 bg-blue-500 rounded-full animate-ping opacity-75"></div>
                  <div className="relative bg-blue-500 text-white rounded-full w-12 h-12 flex items-center justify-center shadow-lg border-4 border-white">
                    <svg
                      className="w-6 h-6"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm0-2a6 6 0 100-12 6 6 0 000 12z"
                        clipRule="evenodd"
                      />
                    </svg>
                  </div>
                </div>
                <div className="absolute top-14 left-1/2 transform -translate-x-1/2 whitespace-nowrap">
                  <div className="bg-blue-500 text-white rounded-lg px-3 py-1 text-sm shadow-lg">
                    Вы здесь
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Map Controls */}
          <div className="absolute top-4 left-4 space-y-2">
            <button className="btn btn-circle btn-sm bg-white shadow-lg">
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                />
              </svg>
            </button>
            <button className="btn btn-circle btn-sm bg-white shadow-lg">
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M20 12H4"
                />
              </svg>
            </button>
          </div>

          {/* Selected Listing Details */}
          {selectedListing && (
            <AnimatedSection
              animation="slideUp"
              className="absolute bottom-4 left-4 right-4"
            >
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <button
                    className="btn btn-sm btn-circle absolute right-2 top-2"
                    onClick={() => setSelectedListing(null)}
                  >
                    ✕
                  </button>
                  <div className="flex gap-4">
                    <figure className="w-24 h-24 flex-shrink-0">
                      <img
                        src={selectedListing.image}
                        alt={selectedListing.title}
                        className="w-full h-full object-cover rounded-lg"
                      />
                    </figure>
                    <div className="flex-1">
                      <h3 className="card-title">{selectedListing.title}</h3>
                      <p className="text-2xl font-bold text-primary">
                        €{selectedListing.price}
                      </p>
                      <div className="flex items-center gap-2 mt-2">
                        <span
                          className={`badge badge-${selectedListing.privacyLevel === 'exact' ? 'error' : selectedListing.privacyLevel === 'area' ? 'warning' : 'success'}`}
                        >
                          {
                            privacyOptions.find(
                              (o) => o.level === selectedListing.privacyLevel
                            )?.title
                          }
                        </span>
                        <span className="text-sm text-base-content/60">
                          {selectedListing.privacyLevel === 'exact'
                            ? selectedListing.address
                            : selectedListing.privacyLevel === 'area'
                              ? 'Район показан приблизительно'
                              : selectedListing.city}
                        </span>
                      </div>
                    </div>
                    <div className="card-actions">
                      <button className="btn btn-primary">Подробнее</button>
                    </div>
                  </div>
                </div>
              </div>
            </AnimatedSection>
          )}
        </div>

        {/* Sidebar */}
        <AnimatedSection
          animation="slideRight"
          className={`transition-all duration-300 ${showFilters ? 'w-96' : 'w-0'} overflow-hidden`}
        >
          <div className="w-96 h-full bg-base-100 shadow-xl p-6 overflow-y-auto">
            <h2 className="text-2xl font-bold mb-6">Настройки приватности</h2>

            {/* Privacy Level Selector */}
            <div className="mb-6">
              <h3 className="font-semibold mb-3">
                Уровень приватности для моих объявлений:
              </h3>
              <div className="space-y-3">
                {privacyOptions.map((option) => (
                  <label
                    key={option.level}
                    className={`card cursor-pointer transition-all hover:shadow-lg ${privacyLevel === option.level ? 'ring-2 ring-primary' : ''}`}
                  >
                    <div className="card-body p-4">
                      <div className="flex items-center gap-3">
                        <input
                          type="radio"
                          name="privacy"
                          className="radio radio-primary"
                          checked={privacyLevel === option.level}
                          onChange={() => setPrivacyLevel(option.level as any)}
                        />
                        <div className="text-2xl">{option.icon}</div>
                        <div className="flex-1">
                          <div className="font-semibold">{option.title}</div>
                          <div className="text-sm text-base-content/60">
                            {option.description}
                          </div>
                        </div>
                      </div>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            <div className="divider"></div>

            {/* Additional Settings */}
            <div className="space-y-4">
              <label className="flex items-center justify-between cursor-pointer">
                <span>Показывать мое местоположение</span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={showMyLocation}
                  onChange={(e) => setShowMyLocation(e.target.checked)}
                />
              </label>

              <label className="flex items-center justify-between cursor-pointer">
                <span>Скрывать адрес до подтверждения</span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  defaultChecked
                />
              </label>

              <label className="flex items-center justify-between cursor-pointer">
                <span>Размытие фото с адресом</span>
                <input type="checkbox" className="toggle toggle-primary" />
              </label>
            </div>

            <div className="divider"></div>

            {/* Info Card */}
            <div className="card bg-info/10 border border-info/20">
              <div className="card-body p-4">
                <h4 className="font-semibold flex items-center gap-2">
                  <svg
                    className="w-5 h-5 text-info"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                      clipRule="evenodd"
                    />
                  </svg>
                  Как это работает?
                </h4>
                <ul className="text-sm space-y-2 mt-2">
                  <li className="flex items-start gap-2">
                    <span className="text-primary">•</span>
                    <span>
                      <strong>Точный адрес:</strong> Покупатели видят точное
                      местоположение на карте
                    </span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-warning">•</span>
                    <span>
                      <strong>Район:</strong> Показывается область радиусом 500м
                    </span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-success">•</span>
                    <span>
                      <strong>Только город:</strong> Максимальная приватность,
                      только название города
                    </span>
                  </li>
                </ul>
              </div>
            </div>

            {/* Statistics */}
            <div className="mt-6">
              <h4 className="font-semibold mb-3">Статистика приватности:</h4>
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <span className="text-sm">Точный адрес</span>
                  <div className="flex items-center gap-2">
                    <progress
                      className="progress progress-error w-32"
                      value="25"
                      max="100"
                    ></progress>
                    <span className="text-sm font-semibold">25%</span>
                  </div>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Район</span>
                  <div className="flex items-center gap-2">
                    <progress
                      className="progress progress-warning w-32"
                      value="60"
                      max="100"
                    ></progress>
                    <span className="text-sm font-semibold">60%</span>
                  </div>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Только город</span>
                  <div className="flex items-center gap-2">
                    <progress
                      className="progress progress-success w-32"
                      value="15"
                      max="100"
                    ></progress>
                    <span className="text-sm font-semibold">15%</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </AnimatedSection>

        {/* Toggle Sidebar Button */}
        <button
          className="absolute right-0 top-1/2 transform -translate-y-1/2 bg-base-100 rounded-l-lg shadow-lg p-2"
          onClick={() => setShowFilters(!showFilters)}
        >
          <svg
            className={`w-6 h-6 transition-transform ${showFilters ? '' : 'rotate-180'}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 5l7 7-7 7"
            />
          </svg>
        </button>
      </div>
    </div>
  );
};

export default MapPrivacy;
