'use client';

import { useState, useEffect } from 'react';
import dynamic from 'next/dynamic';
import {
  MapPinIcon,
  ClockIcon,
  PhoneIcon,
  CheckCircleIcon,
  MagnifyingGlassIcon,
  FunnelIcon,
  BuildingStorefrontIcon,
} from '@heroicons/react/24/outline';
import 'leaflet/dist/leaflet.css';

// Dynamic import for Leaflet components to avoid SSR issues
const MapContainer = dynamic(
  () => import('react-leaflet').then((mod) => mod.MapContainer),
  { ssr: false }
);
const TileLayer = dynamic(
  () => import('react-leaflet').then((mod) => mod.TileLayer),
  { ssr: false }
);
const Marker = dynamic(
  () => import('react-leaflet').then((mod) => mod.Marker),
  { ssr: false }
);
const Popup = dynamic(() => import('react-leaflet').then((mod) => mod.Popup), {
  ssr: false,
});

interface ParcelShop {
  id: number;
  name: string;
  address: string;
  city: string;
  distance: string;
  workHours: string;
  saturdayHours: string;
  phone: string;
  coordinates: { lat: number; lng: number };
  type: 'bex' | 'post-express' | 'yettel' | 'svetu-warehouse';
  company: string;
}

export default function SerbianParcelShopMap() {
  const [selectedShop, setSelectedShop] = useState<number | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterCity, setFilterCity] = useState('all');
  const [filterCompany, setFilterCompany] = useState('all');
  const [isClient, setIsClient] = useState(false);

  useEffect(() => {
    setIsClient(true);
  }, []);

  const parcelShops: ParcelShop[] = [
    {
      id: 1,
      name: 'BexExpress Централни депо',
      address: 'Милутина Миланковића 9б',
      city: 'Нови Београд',
      distance: '2.5 км',
      workHours: '08:00 - 20:00',
      saturdayHours: '09:00 - 16:00',
      phone: '011/555-7890',
      coordinates: { lat: 44.8097, lng: 20.4291 },
      type: 'bex',
      company: 'BexExpress',
    },
    {
      id: 2,
      name: 'Sve Tu Platforma - Главни магацин',
      address: 'Булевар Краља Александра 73',
      city: 'Београд',
      distance: '1.2 км',
      workHours: '08:00 - 18:00',
      saturdayHours: '10:00 - 14:00',
      phone: '011/555-1234',
      coordinates: { lat: 44.8036, lng: 20.4848 },
      type: 'svetu-warehouse',
      company: 'Sve Tu Platforma',
    },
    {
      id: 3,
      name: 'Post Express Кнез Михаилова',
      address: 'Кнез Михаилова 23',
      city: 'Београд',
      distance: '0.5 км',
      workHours: '08:00 - 22:00',
      saturdayHours: '10:00 - 20:00',
      phone: '011/333-4567',
      coordinates: { lat: 44.8176, lng: 20.4568 },
      type: 'post-express',
      company: 'Post Express',
    },
    {
      id: 4,
      name: 'Yettel Post Пакомат БГЦЦ',
      address: 'Омладинских бригада 86б',
      city: 'Нови Београд',
      distance: '3.2 км',
      workHours: '24/7',
      saturdayHours: '24/7',
      phone: '0800 100 100',
      coordinates: { lat: 44.8094, lng: 20.4322 },
      type: 'yettel',
      company: 'Yettel Post',
    },
    {
      id: 5,
      name: 'BexExpress Parcel Shop Земун',
      address: 'Главна 45',
      city: 'Земун',
      distance: '7.3 км',
      workHours: '09:00 - 19:00',
      saturdayHours: '10:00 - 16:00',
      phone: '011/316-7890',
      coordinates: { lat: 44.8433, lng: 20.4012 },
      type: 'bex',
      company: 'BexExpress',
    },
    {
      id: 6,
      name: 'Post Express Променада',
      address: 'Булевар Ослобођења 119',
      city: 'Нови Сад',
      distance: '85 км',
      workHours: '10:00 - 22:00',
      saturdayHours: '10:00 - 22:00',
      phone: '021/123-4567',
      coordinates: { lat: 45.2471, lng: 19.8335 },
      type: 'post-express',
      company: 'Post Express',
    },
    {
      id: 7,
      name: 'BexExpress Ниш Центар',
      address: 'Обреновићева 12',
      city: 'Ниш',
      distance: '240 км',
      workHours: '09:00 - 20:00',
      saturdayHours: '10:00 - 18:00',
      phone: '018/123-7890',
      coordinates: { lat: 43.3209, lng: 21.8958 },
      type: 'bex',
      company: 'BexExpress',
    },
  ];

  const filteredShops = parcelShops.filter((shop) => {
    const matchesSearch =
      shop.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      shop.address.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesCity = filterCity === 'all' || shop.city === filterCity;
    const matchesCompany =
      filterCompany === 'all' || shop.company === filterCompany;
    return matchesSearch && matchesCity && matchesCompany;
  });

  const getShopColor = (type: string) => {
    switch (type) {
      case 'bex':
        return 'bg-purple-500 text-white';
      case 'svetu-warehouse':
        return 'bg-green-500 text-white';
      case 'post-express':
        return 'bg-blue-500 text-white';
      case 'yettel':
        return 'bg-orange-500 text-white';
      default:
        return 'bg-base-100';
    }
  };

  const getMarkerIcon = (type: string) => {
    if (!isClient) return null;

    const L = require('leaflet');
    const iconHtml =
      {
        bex: `<div class="flex items-center justify-center w-8 h-8 bg-purple-600 rounded-full text-white shadow-lg"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></div>`,
        'svetu-warehouse': `<div class="flex items-center justify-center w-8 h-8 bg-green-600 rounded-full text-white shadow-lg"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path></svg></div>`,
        'post-express': `<div class="flex items-center justify-center w-8 h-8 bg-blue-600 rounded-full text-white shadow-lg"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></div>`,
        yettel: `<div class="flex items-center justify-center w-8 h-8 bg-orange-600 rounded-full text-white shadow-lg"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></div>`,
      }[type] || '';

    return L.divIcon({
      html: iconHtml,
      iconSize: [32, 32],
      iconAnchor: [16, 32],
      popupAnchor: [0, -32],
      className: 'custom-marker',
    });
  };

  return (
    <div className="space-y-6">
      {/* Search and Filter */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="flex gap-3 flex-wrap">
            <div className="flex-1 min-w-[200px]">
              <div className="relative">
                <MagnifyingGlassIcon className="w-5 h-5 absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" />
                <input
                  type="text"
                  placeholder="Претрага по називу или адреси"
                  className="input input-bordered w-full pl-10"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
              </div>
            </div>
            <select
              className="select select-bordered"
              value={filterCity}
              onChange={(e) => setFilterCity(e.target.value)}
            >
              <option value="all">Сви градови</option>
              <option value="Београд">Београд</option>
              <option value="Нови Београд">Нови Београд</option>
              <option value="Земун">Земун</option>
              <option value="Нови Сад">Нови Сад</option>
              <option value="Ниш">Ниш</option>
            </select>
            <select
              className="select select-bordered"
              value={filterCompany}
              onChange={(e) => setFilterCompany(e.target.value)}
            >
              <option value="all">Све службе</option>
              <option value="BexExpress">BexExpress</option>
              <option value="Sve Tu Platforma">Sve Tu Platforma</option>
              <option value="Post Express">Post Express</option>
              <option value="Yettel Post">Yettel Post</option>
            </select>
            <button className="btn btn-outline">
              <FunnelIcon className="w-5 h-5" />
              Филтери
            </button>
          </div>
        </div>
      </div>

      <div className="grid lg:grid-cols-2 gap-6">
        {/* Interactive Map */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body p-0">
            <div className="h-[500px] relative rounded-lg overflow-hidden">
              {isClient && (
                <MapContainer
                  center={[44.8176, 20.4568]}
                  zoom={11}
                  style={{ height: '100%', width: '100%' }}
                  className="rounded-lg"
                >
                  <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                  />
                  {filteredShops.map((shop) => (
                    <Marker
                      key={shop.id}
                      position={[shop.coordinates.lat, shop.coordinates.lng]}
                      icon={getMarkerIcon(shop.type)}
                      eventHandlers={{
                        click: () => setSelectedShop(shop.id),
                      }}
                    >
                      <Popup>
                        <div className="p-2">
                          <h4 className="font-bold text-sm">{shop.name}</h4>
                          <p className="text-xs">{shop.address}</p>
                          <p className="text-xs">{shop.workHours}</p>
                          <p className="text-xs font-semibold mt-1">
                            {shop.company}
                          </p>
                        </div>
                      </Popup>
                    </Marker>
                  ))}
                </MapContainer>
              )}
              {!isClient && (
                <div className="flex items-center justify-center h-full bg-base-200">
                  <div className="text-center">
                    <MapPinIcon className="w-12 h-12 mx-auto text-base-content/30 mb-2" />
                    <p className="text-base-content/50">Учитава се карта...</p>
                  </div>
                </div>
              )}

              {/* Legend */}
              <div className="absolute bottom-4 left-4 bg-base-100/90 backdrop-blur-sm p-3 rounded-lg z-[1000]">
                <h5 className="font-semibold text-sm mb-2">Легенда:</h5>
                <div className="flex items-center gap-2 text-sm mb-1">
                  <div className="w-4 h-4 bg-purple-600 rounded-full"></div>
                  <span>BexExpress</span>
                </div>
                <div className="flex items-center gap-2 text-sm mb-1">
                  <div className="w-4 h-4 bg-green-600 rounded-full"></div>
                  <span>Sve Tu Platforma</span>
                </div>
                <div className="flex items-center gap-2 text-sm mb-1">
                  <div className="w-4 h-4 bg-blue-600 rounded-full"></div>
                  <span>Post Express</span>
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <div className="w-4 h-4 bg-orange-600 rounded-full"></div>
                  <span>Yettel Post 24/7</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Shops List */}
        <div className="space-y-4 max-h-[500px] overflow-y-auto">
          {filteredShops.map((shop) => (
            <div
              key={shop.id}
              className={`
                card cursor-pointer transition-all
                ${
                  selectedShop === shop.id
                    ? 'ring-2 ring-primary shadow-xl'
                    : 'bg-base-100 hover:shadow-lg'
                }
              `}
              onClick={() => setSelectedShop(shop.id)}
            >
              <div className="card-body">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <h4 className="font-semibold">{shop.name}</h4>
                      <div
                        className={`badge badge-sm ${
                          shop.type === 'bex'
                            ? 'badge-primary'
                            : shop.type === 'svetu-warehouse'
                              ? 'badge-success'
                              : shop.type === 'post-express'
                                ? 'badge-info'
                                : 'badge-warning'
                        }`}
                      >
                        {shop.company}
                      </div>
                    </div>
                    <p className="text-sm text-base-content/60 mt-1">
                      {shop.address}
                    </p>
                    <p className="text-sm text-base-content/60">{shop.city}</p>

                    <div className="flex items-center gap-4 mt-3 text-sm">
                      <div className="flex items-center gap-1">
                        <MapPinIcon className="w-4 h-4 text-primary" />
                        <span className="font-medium">{shop.distance}</span>
                      </div>
                      <div className="flex items-center gap-1">
                        <ClockIcon className="w-4 h-4 text-base-content/40" />
                        <span>{shop.workHours}</span>
                      </div>
                    </div>

                    {shop.saturdayHours !== '24/7' && (
                      <div className="text-xs text-base-content/60 mt-2">
                        Субота: {shop.saturdayHours}
                      </div>
                    )}

                    <div className="flex items-center gap-1 mt-2">
                      <PhoneIcon className="w-4 h-4 text-base-content/40" />
                      <span className="text-sm">{shop.phone}</span>
                    </div>

                    {shop.type === 'yettel' && (
                      <div className="badge badge-accent badge-sm mt-2">
                        🕐 Доступан 24/7
                      </div>
                    )}
                    {shop.type === 'bex' && (
                      <div className="badge badge-primary badge-sm mt-2">
                        📡 API интеграција
                      </div>
                    )}
                    {shop.type === 'svetu-warehouse' && (
                      <div className="badge badge-success badge-sm mt-2">
                        <BuildingStorefrontIcon className="w-3 h-3 mr-1" />
                        Наш магацин
                      </div>
                    )}
                  </div>

                  {selectedShop === shop.id && (
                    <div className="p-2 bg-success text-success-content rounded-full">
                      <CheckCircleIcon className="w-5 h-5" />
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Selected Shop Actions */}
      {selectedShop && (
        <div className="alert alert-success">
          <CheckCircleIcon className="w-5 h-5" />
          <div className="flex-1">
            <h4 className="font-semibold">
              Изабран пункт:{' '}
              {parcelShops.find((s) => s.id === selectedShop)?.name}
            </h4>
            <p className="text-sm mt-1">
              Пошиљка ће бити достављена на овај пункт. Добићете СМС обавештење
              када будете могли да преузмете пакет.
            </p>
          </div>
          <button className="btn btn-sm">Потврди избор</button>
        </div>
      )}

      {/* Serbian Info */}
      <div className="bg-gradient-to-r from-purple-50 via-green-50 to-blue-50 p-4 rounded-lg">
        <h4 className="font-semibold text-purple-700 mb-2">
          🇷🇸 Важне информације
        </h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
          <div>
            <strong className="text-purple-600">BexExpress:</strong> API
            интеграција, 3 пункта
          </div>
          <div>
            <strong className="text-green-600">Sve Tu Platforma:</strong> Наш
            главни магацин
          </div>
          <div>
            <strong className="text-blue-600">Post Express:</strong> 500+
            пунктова, ради до 22:00
          </div>
          <div>
            <strong className="text-orange-600">Yettel Post:</strong> Пакомати
            24/7, QR код
          </div>
          <div className="md:col-span-2">
            <strong>За преузимање:</strong> Обавезна лична карта или пасош
          </div>
        </div>
      </div>
    </div>
  );
}
