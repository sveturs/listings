'use client';

import { useState } from 'react';
import {
  MapPinIcon,
  ClockIcon,
  PhoneIcon,
  CheckCircleIcon,
  MagnifyingGlassIcon,
  FunnelIcon,
} from '@heroicons/react/24/outline';

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
  type: 'shop' | 'warehouse';
}

export default function ParcelShopMap() {
  const [selectedShop, setSelectedShop] = useState<number | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterCity, setFilterCity] = useState('all');

  const parcelShops: ParcelShop[] = [
    {
      id: 1,
      name: 'BEX Paket Shop Centar',
      address: 'Knez Mihailova 23',
      city: 'Белград',
      distance: '0.5 км',
      workHours: '09:00 - 20:00',
      saturdayHours: '10:00 - 18:00',
      phone: '+381 11 123 4567',
      coordinates: { lat: 44.8125, lng: 20.4612 },
      type: 'shop',
    },
    {
      id: 2,
      name: 'BEX Magacin Novi Beograd',
      address: 'Bulevar Mihajla Pupina 165',
      city: 'Белград',
      distance: '3.2 км',
      workHours: '08:00 - 19:00',
      saturdayHours: '09:00 - 15:00',
      phone: '+381 11 234 5678',
      coordinates: { lat: 44.8176, lng: 20.4633 },
      type: 'warehouse',
    },
    {
      id: 3,
      name: 'BEX Paket Shop Zemun',
      address: 'Glavna 45',
      city: 'Земун',
      distance: '5.8 км',
      workHours: '09:00 - 19:00',
      saturdayHours: '10:00 - 16:00',
      phone: '+381 11 345 6789',
      coordinates: { lat: 44.8433, lng: 20.4112 },
      type: 'shop',
    },
    {
      id: 4,
      name: 'BEX Paket Shop Novi Sad',
      address: 'Dunavska 15',
      city: 'Нови Сад',
      distance: '85 км',
      workHours: '09:00 - 20:00',
      saturdayHours: '10:00 - 17:00',
      phone: '+381 21 123 4567',
      coordinates: { lat: 45.2671, lng: 19.8335 },
      type: 'shop',
    },
  ];

  const filteredShops = parcelShops.filter((shop) => {
    const matchesSearch =
      shop.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      shop.address.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesCity = filterCity === 'all' || shop.city === filterCity;
    return matchesSearch && matchesCity;
  });

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
                  placeholder="Поиск по названию или адресу"
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
              <option value="all">Все города</option>
              <option value="Белград">Белград</option>
              <option value="Земун">Земун</option>
              <option value="Нови Сад">Нови Сад</option>
            </select>
            <button className="btn btn-outline">
              <FunnelIcon className="w-5 h-5" />
              Фильтры
            </button>
          </div>
        </div>
      </div>

      <div className="grid lg:grid-cols-2 gap-6">
        {/* Map Placeholder */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body p-0">
            <div className="bg-gradient-to-br from-primary/20 to-secondary/20 h-[500px] relative rounded-lg overflow-hidden">
              {/* Map Background Pattern */}
              <div className="absolute inset-0 opacity-30">
                <div className="grid grid-cols-8 grid-rows-8 h-full">
                  {[...Array(64)].map((_, i) => (
                    <div
                      key={i}
                      className="border border-base-content/10"
                    ></div>
                  ))}
                </div>
              </div>

              {/* Map Pins */}
              {filteredShops.map((shop) => (
                <div
                  key={shop.id}
                  className="absolute cursor-pointer transform -translate-x-1/2 -translate-y-1/2"
                  style={{
                    left: `${20 + (shop.coordinates.lng - 19.5) * 80}%`,
                    top: `${80 - (shop.coordinates.lat - 44.5) * 160}%`,
                  }}
                  onClick={() => setSelectedShop(shop.id)}
                >
                  <div
                    className={`
                    p-2 rounded-full shadow-lg transition-all
                    ${
                      selectedShop === shop.id
                        ? 'bg-primary text-primary-content scale-125'
                        : shop.type === 'warehouse'
                          ? 'bg-secondary text-secondary-content hover:scale-110'
                          : 'bg-base-100 hover:scale-110'
                    }
                  `}
                  >
                    <MapPinIcon className="w-6 h-6" />
                  </div>
                  {selectedShop === shop.id && (
                    <div className="absolute top-10 left-1/2 -translate-x-1/2 bg-base-100 p-2 rounded shadow-lg whitespace-nowrap z-10">
                      <div className="font-semibold text-sm">{shop.name}</div>
                      <div className="text-xs text-base-content/60">
                        {shop.distance}
                      </div>
                    </div>
                  )}
                </div>
              ))}

              {/* Legend */}
              <div className="absolute bottom-4 left-4 bg-base-100/90 backdrop-blur-sm p-3 rounded-lg">
                <div className="flex items-center gap-2 text-sm mb-2">
                  <div className="w-4 h-4 bg-base-100 rounded-full border-2 border-base-content"></div>
                  <span>Пункт выдачи</span>
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <div className="w-4 h-4 bg-secondary rounded-full"></div>
                  <span>Склад</span>
                </div>
              </div>

              {/* Map Controls */}
              <div className="absolute top-4 right-4 flex flex-col gap-2">
                <button className="btn btn-circle btn-sm bg-base-100">
                  <span className="text-lg">+</span>
                </button>
                <button className="btn btn-circle btn-sm bg-base-100">
                  <span className="text-lg">−</span>
                </button>
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
                      {shop.type === 'warehouse' && (
                        <div className="badge badge-secondary badge-sm">
                          Склад
                        </div>
                      )}
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

                    <div className="text-xs text-base-content/60 mt-2">
                      Суббота: {shop.saturdayHours}
                    </div>

                    <div className="flex items-center gap-1 mt-2">
                      <PhoneIcon className="w-4 h-4 text-base-content/40" />
                      <span className="text-sm">{shop.phone}</span>
                    </div>
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
              Выбран пункт:{' '}
              {parcelShops.find((s) => s.id === selectedShop)?.name}
            </h4>
            <p className="text-sm mt-1">
              Посылка будет доставлена в этот пункт выдачи. Вы получите СМС
              когда можно будет забрать.
            </p>
          </div>
          <button className="btn btn-sm">Подтвердить выбор</button>
        </div>
      )}
    </div>
  );
}
