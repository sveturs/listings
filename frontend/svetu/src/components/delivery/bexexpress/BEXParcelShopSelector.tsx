'use client';

import { useState, useEffect } from 'react';
import {
  BuildingStorefrontIcon,
  MapPinIcon,
  ClockIcon,
  PhoneIcon,
  CheckCircleIcon,
  MagnifyingGlassIcon,
  MapIcon,
  ListBulletIcon,
  AdjustmentsHorizontalIcon,
  InformationCircleIcon,
} from '@heroicons/react/24/outline';
import { motion, AnimatePresence } from 'framer-motion';

interface ParcelShop {
  id: number;
  code: string;
  name: string;
  address: string;
  city: string;
  postal_code: string;
  phone?: string;
  latitude?: number;
  longitude?: number;
  working_hours: Record<string, string>;
  services: string[];
  distance?: number;
}

interface Props {
  selectedCity: string;
  onShopSelect: (shop: ParcelShop) => void;
  selectedShop?: ParcelShop;
  className?: string;
}

export default function BEXParcelShopSelector({
  selectedCity,
  onShopSelect,
  selectedShop,
  className = '',
}: Props) {
  const [shops, setShops] = useState<ParcelShop[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState<'list' | 'map'>('list');
  const [showFilters, setShowFilters] = useState(false);
  const [filters, setFilters] = useState({
    openNow: false,
    openWeekends: false,
    hasParking: false,
  });

  // Загрузка пунктов выдачи
  useEffect(() => {
    if (!selectedCity) return;

    const fetchShops = async () => {
      setLoading(true);
      try {
        const response = await fetch(
          `/api/v1/bex/parcel-shops?city=${encodeURIComponent(selectedCity)}`
        );
        if (response.ok) {
          const data = await response.json();
          if (data.success && data.data) {
            // Добавляем случайные демо-данные для примера
            const enrichedShops = data.data.map((shop: ParcelShop) => ({
              ...shop,
              distance: Math.random() * 10, // Случайное расстояние для демо
              services: ['cod', 'insurance', 'packaging'],
              working_hours: shop.working_hours || {
                monday: '08:00-20:00',
                tuesday: '08:00-20:00',
                wednesday: '08:00-20:00',
                thursday: '08:00-20:00',
                friday: '08:00-20:00',
                saturday: '09:00-15:00',
                sunday: 'Закрыто',
              },
            }));
            setShops(enrichedShops);
          }
        }
      } catch (error) {
        console.error('Failed to fetch parcel shops:', error);
        // Используем демо-данные при ошибке
        setShops(getDemoShops(selectedCity));
      } finally {
        setLoading(false);
      }
    };

    fetchShops();
  }, [selectedCity]);

  // Демо-данные для тестирования
  const getDemoShops = (city: string): ParcelShop[] => [
    {
      id: 1,
      code: 'BEX-NS-001',
      name: 'BEX Express Центр',
      address: 'Булевар ослобођења 45',
      city,
      postal_code: '21000',
      phone: '+381 21 123 456',
      latitude: 45.2671,
      longitude: 19.8335,
      working_hours: {
        monday: '08:00-20:00',
        tuesday: '08:00-20:00',
        wednesday: '08:00-20:00',
        thursday: '08:00-20:00',
        friday: '08:00-20:00',
        saturday: '09:00-15:00',
        sunday: 'Закрыто',
      },
      services: ['cod', 'insurance', 'packaging', 'parking'],
      distance: 1.2,
    },
    {
      id: 2,
      code: 'BEX-NS-002',
      name: 'BEX Express Лиман',
      address: 'Народног фронта 12',
      city,
      postal_code: '21000',
      phone: '+381 21 234 567',
      latitude: 45.2394,
      longitude: 19.8365,
      working_hours: {
        monday: '09:00-19:00',
        tuesday: '09:00-19:00',
        wednesday: '09:00-19:00',
        thursday: '09:00-19:00',
        friday: '09:00-19:00',
        saturday: '10:00-14:00',
        sunday: 'Закрыто',
      },
      services: ['cod', 'insurance'],
      distance: 2.5,
    },
    {
      id: 3,
      code: 'BEX-NS-003',
      name: 'BEX Express Детелинара',
      address: 'Корнелија Станковића 3',
      city,
      postal_code: '21000',
      phone: '+381 21 345 678',
      working_hours: {
        monday: '08:00-18:00',
        tuesday: '08:00-18:00',
        wednesday: '08:00-18:00',
        thursday: '08:00-18:00',
        friday: '08:00-18:00',
        saturday: 'Закрыто',
        sunday: 'Закрыто',
      },
      services: ['cod', 'parking'],
      distance: 3.8,
    },
  ];

  // Фильтрация магазинов
  const filteredShops = shops.filter((shop) => {
    // Поиск по названию или адресу
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      if (
        !shop.name.toLowerCase().includes(query) &&
        !shop.address.toLowerCase().includes(query) &&
        !shop.code.toLowerCase().includes(query)
      ) {
        return false;
      }
    }

    // Фильтр "Открыто сейчас"
    if (filters.openNow) {
      const now = new Date();
      const day = [
        'sunday',
        'monday',
        'tuesday',
        'wednesday',
        'thursday',
        'friday',
        'saturday',
      ][now.getDay()];
      const hours = shop.working_hours[day];
      if (hours === 'Закрыто') return false;
      // Здесь можно добавить более точную проверку времени
    }

    // Фильтр "Работает в выходные"
    if (filters.openWeekends) {
      if (
        shop.working_hours.saturday === 'Закрыто' &&
        shop.working_hours.sunday === 'Закрыто'
      ) {
        return false;
      }
    }

    // Фильтр "Есть парковка"
    if (filters.hasParking) {
      if (!shop.services.includes('parking')) {
        return false;
      }
    }

    return true;
  });

  // Сортировка по расстоянию
  const sortedShops = [...filteredShops].sort(
    (a, b) => (a.distance || 0) - (b.distance || 0)
  );

  // Получение текущего дня недели
  const getCurrentDay = () => {
    const days = [
      'sunday',
      'monday',
      'tuesday',
      'wednesday',
      'thursday',
      'friday',
      'saturday',
    ];
    return days[new Date().getDay()];
  };

  // Проверка, открыт ли магазин
  const isOpenNow = (shop: ParcelShop) => {
    const currentDay = getCurrentDay();
    const hours = shop.working_hours[currentDay];
    if (hours === 'Закрыто') return false;

    // Простая проверка для демо
    const now = new Date();
    const currentHour = now.getHours();
    const [openTime, closeTime] = hours
      .split('-')
      .map((t) => parseInt(t.split(':')[0]));
    return currentHour >= openTime && currentHour < closeTime;
  };

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Панель управления */}
      <div className="flex flex-col sm:flex-row gap-4">
        {/* Поиск */}
        <div className="flex-1">
          <div className="relative">
            <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
            <input
              type="text"
              className="input input-bordered w-full pl-11"
              placeholder="Поиск по названию, адресу или коду..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
        </div>

        {/* Кнопки управления */}
        <div className="flex gap-2">
          {/* Переключатель вида */}
          <div className="btn-group">
            <button
              className={`btn btn-sm ${viewMode === 'list' ? 'btn-active' : ''}`}
              onClick={() => setViewMode('list')}
            >
              <ListBulletIcon className="w-4 h-4" />
            </button>
            <button
              className={`btn btn-sm ${viewMode === 'map' ? 'btn-active' : ''}`}
              onClick={() => setViewMode('map')}
            >
              <MapIcon className="w-4 h-4" />
            </button>
          </div>

          {/* Фильтры */}
          <button
            className={`btn btn-sm ${showFilters ? 'btn-active' : ''}`}
            onClick={() => setShowFilters(!showFilters)}
          >
            <AdjustmentsHorizontalIcon className="w-4 h-4" />
            Фильтры
          </button>
        </div>
      </div>

      {/* Панель фильтров */}
      <AnimatePresence>
        {showFilters && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            className="card bg-base-200 overflow-hidden"
          >
            <div className="card-body p-4">
              <div className="flex flex-wrap gap-4">
                <label className="label cursor-pointer gap-2">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-sm"
                    checked={filters.openNow}
                    onChange={(e) =>
                      setFilters({ ...filters, openNow: e.target.checked })
                    }
                  />
                  <span className="label-text">Открыто сейчас</span>
                </label>
                <label className="label cursor-pointer gap-2">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-sm"
                    checked={filters.openWeekends}
                    onChange={(e) =>
                      setFilters({ ...filters, openWeekends: e.target.checked })
                    }
                  />
                  <span className="label-text">Работает в выходные</span>
                </label>
                <label className="label cursor-pointer gap-2">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-sm"
                    checked={filters.hasParking}
                    onChange={(e) =>
                      setFilters({ ...filters, hasParking: e.target.checked })
                    }
                  />
                  <span className="label-text">Есть парковка</span>
                </label>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Список или карта */}
      {loading ? (
        <div className="flex justify-center items-center h-64">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : viewMode === 'list' ? (
        <div className="space-y-3">
          {sortedShops.length === 0 ? (
            <div className="text-center py-8">
              <BuildingStorefrontIcon className="w-12 h-12 mx-auto text-base-content/30 mb-3" />
              <p className="text-base-content/60">
                {searchQuery ||
                filters.openNow ||
                filters.openWeekends ||
                filters.hasParking
                  ? 'Не найдено пунктов выдачи по заданным критериям'
                  : `Пункты выдачи в городе ${selectedCity} не найдены`}
              </p>
            </div>
          ) : (
            sortedShops.map((shop) => {
              const isSelected = selectedShop?.id === shop.id;
              const isOpen = isOpenNow(shop);

              return (
                <motion.div
                  key={shop.id}
                  whileHover={{ scale: 1.01 }}
                  onClick={() => onShopSelect(shop)}
                  className={`
                    card bg-base-100 shadow-sm cursor-pointer transition-all
                    ${isSelected ? 'ring-2 ring-primary shadow-lg' : 'hover:shadow-md'}
                  `}
                >
                  <div className="card-body p-4">
                    <div className="flex items-start gap-4">
                      {/* Иконка */}
                      <div
                        className={`
                        p-3 rounded-lg flex-shrink-0
                        ${isSelected ? 'bg-primary/10 text-primary' : 'bg-base-200'}
                      `}
                      >
                        <BuildingStorefrontIcon className="w-6 h-6" />
                      </div>

                      {/* Информация */}
                      <div className="flex-1">
                        <div className="flex items-start justify-between">
                          <div>
                            <h4 className="font-semibold text-lg">
                              {shop.name}
                            </h4>
                            <p className="text-sm text-base-content/60 mb-1">
                              {shop.code}
                            </p>
                          </div>
                          {isSelected && (
                            <CheckCircleIcon className="w-6 h-6 text-primary flex-shrink-0" />
                          )}
                        </div>

                        <div className="space-y-2 mt-3">
                          {/* Адрес */}
                          <div className="flex items-start gap-2">
                            <MapPinIcon className="w-4 h-4 text-base-content/60 mt-0.5 flex-shrink-0" />
                            <div className="text-sm">
                              <p>{shop.address}</p>
                              <p className="text-base-content/60">
                                {shop.postal_code} {shop.city}
                              </p>
                            </div>
                          </div>

                          {/* Расстояние */}
                          {shop.distance && (
                            <div className="flex items-center gap-2">
                              <MapIcon className="w-4 h-4 text-info flex-shrink-0" />
                              <span className="text-sm">
                                {shop.distance.toFixed(1)} км от центра города
                              </span>
                            </div>
                          )}

                          {/* Телефон */}
                          {shop.phone && (
                            <div className="flex items-center gap-2">
                              <PhoneIcon className="w-4 h-4 text-base-content/60 flex-shrink-0" />
                              <span className="text-sm">{shop.phone}</span>
                            </div>
                          )}

                          {/* Часы работы */}
                          <div className="flex items-start gap-2">
                            <ClockIcon className="w-4 h-4 text-base-content/60 mt-0.5 flex-shrink-0" />
                            <div className="text-sm">
                              <div className="flex items-center gap-2">
                                <span>
                                  Сегодня: {shop.working_hours[getCurrentDay()]}
                                </span>
                                {isOpen ? (
                                  <span className="badge badge-success badge-sm">
                                    Открыто
                                  </span>
                                ) : (
                                  <span className="badge badge-error badge-sm">
                                    Закрыто
                                  </span>
                                )}
                              </div>

                              {/* Расписание на неделю (свернуто) */}
                              <details className="mt-1">
                                <summary className="cursor-pointer text-primary hover:underline">
                                  Показать расписание
                                </summary>
                                <div className="mt-2 space-y-1">
                                  <div>Пн: {shop.working_hours.monday}</div>
                                  <div>Вт: {shop.working_hours.tuesday}</div>
                                  <div>Ср: {shop.working_hours.wednesday}</div>
                                  <div>Чт: {shop.working_hours.thursday}</div>
                                  <div>Пт: {shop.working_hours.friday}</div>
                                  <div className="font-medium">
                                    Сб: {shop.working_hours.saturday}
                                  </div>
                                  <div className="font-medium">
                                    Вс: {shop.working_hours.sunday}
                                  </div>
                                </div>
                              </details>
                            </div>
                          </div>

                          {/* Услуги */}
                          <div className="flex flex-wrap gap-2 mt-2">
                            {shop.services.includes('cod') && (
                              <span className="badge badge-sm">
                                Наложенный платеж
                              </span>
                            )}
                            {shop.services.includes('insurance') && (
                              <span className="badge badge-sm">
                                Страхование
                              </span>
                            )}
                            {shop.services.includes('packaging') && (
                              <span className="badge badge-sm">Упаковка</span>
                            )}
                            {shop.services.includes('parking') && (
                              <span className="badge badge-sm badge-success">
                                Парковка
                              </span>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </motion.div>
              );
            })
          )}
        </div>
      ) : (
        <div className="card bg-base-200 h-96">
          <div className="card-body flex items-center justify-center">
            <MapIcon className="w-12 h-12 text-base-content/30 mb-3" />
            <p className="text-base-content/60">
              Карта пунктов выдачи будет доступна в ближайшее время
            </p>
            <button
              className="btn btn-sm btn-primary mt-3"
              onClick={() => setViewMode('list')}
            >
              Вернуться к списку
            </button>
          </div>
        </div>
      )}

      {/* Информационная панель */}
      {selectedShop && (
        <div className="alert alert-info">
          <InformationCircleIcon className="w-5 h-5" />
          <div>
            <h4 className="font-semibold">Выбран пункт: {selectedShop.name}</h4>
            <p className="text-sm">
              Посылка будет доставлена по адресу: {selectedShop.address}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
