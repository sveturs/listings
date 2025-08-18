'use client';

import { useState, useEffect } from 'react';
import {
  MapPinIcon,
  MagnifyingGlassIcon,
  ClockIcon,
  PhoneIcon,
  CheckIcon,
  BuildingStorefrontIcon,
  ArrowsPointingOutIcon,
  InformationCircleIcon,
} from '@heroicons/react/24/outline';
// import { useTranslations } from 'next-intl';

interface PostOffice {
  id: number;
  code: string;
  name: string;
  address: string;
  city: string;
  postal_code: string;
  phone?: string;
  latitude?: number;
  longitude?: number;
  working_hours: {
    monday?: string;
    tuesday?: string;
    wednesday?: string;
    thursday?: string;
    friday?: string;
    saturday?: string;
    sunday?: string;
  };
  services: string[];
  distance?: number; // км от пользователя
}

interface Props {
  selectedCity?: string;
  onOfficeSelect: (office: PostOffice) => void;
  selectedOffice?: PostOffice;
  className?: string;
}

export default function PostExpressOfficeSelector({
  selectedCity,
  onOfficeSelect,
  selectedOffice,
  className = '',
}: Props) {
  // const t = useTranslations('delivery');
  const [offices, setOffices] = useState<PostOffice[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<'distance' | 'name' | 'working_hours'>(
    'distance'
  );
  const [showMap, setShowMap] = useState(false);

  // Загрузка отделений
  useEffect(() => {
    if (selectedCity) {
      loadOffices();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedCity]);

  const loadOffices = async () => {
    if (!selectedCity) return;

    setLoading(true);
    try {
      const params = new URLSearchParams({
        city: selectedCity,
        limit: '50',
        sort: sortBy,
      });

      const response = await fetch(`/api/v1/postexpress/offices?${params}`);
      const data = await response.json();

      if (data.success) {
        setOffices(data.data || []);
      }
    } catch (error) {
      console.error('Ошибка загрузки отделений:', error);
    } finally {
      setLoading(false);
    }
  };

  // Фильтрация отделений по поиску
  const filteredOffices = offices.filter(
    (office) =>
      office.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      office.address.toLowerCase().includes(searchQuery.toLowerCase()) ||
      office.code.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleOfficeSelect = (office: PostOffice) => {
    onOfficeSelect(office);
  };

  const getWorkingHoursToday = (office: PostOffice) => {
    const dayNames = [
      'sunday',
      'monday',
      'tuesday',
      'wednesday',
      'thursday',
      'friday',
      'saturday',
    ];
    const todayName = dayNames[
      new Date().getDay()
    ] as keyof typeof office.working_hours;

    return office.working_hours[todayName] || 'Часы не указаны';
  };

  const isOpenNow = (office: PostOffice) => {
    const now = new Date();
    const currentHour = now.getHours();
    const todayHours = getWorkingHoursToday(office);

    if (todayHours === 'Закрыто' || todayHours === 'Часы не указаны') {
      return false;
    }

    // Простая проверка (может быть улучшена)
    const match = todayHours.match(/(\d{1,2}):(\d{2})-(\d{1,2}):(\d{2})/);
    if (match) {
      const [, startHour, startMin, endHour, endMin] = match;
      const start = parseInt(startHour) + parseInt(startMin) / 60;
      const end = parseInt(endHour) + parseInt(endMin) / 60;
      const current = currentHour + now.getMinutes() / 60;

      return current >= start && current <= end;
    }

    return false;
  };

  if (!selectedCity) {
    return (
      <div className={`text-center py-12 ${className}`}>
        <BuildingStorefrontIcon className="w-16 h-16 mx-auto text-base-content/30 mb-4" />
        <h3 className="text-lg font-semibold mb-2">Выберите город</h3>
        <p className="text-base-content/60">
          Сначала укажите город для отображения доступных отделений
        </p>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Заголовок и поиск */}
      <div className="space-y-4">
        <div className="text-center">
          <h3 className="text-xl font-bold mb-2">
            Выберите отделение Post Express
          </h3>
          <p className="text-base-content/70">
            {filteredOffices.length} отделений в городе {selectedCity}
          </p>
        </div>

        {/* Поиск и фильтры */}
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body p-4">
            <div className="flex flex-col sm:flex-row gap-4">
              {/* Поиск */}
              <div className="flex-1">
                <div className="relative">
                  <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
                  <input
                    type="text"
                    className="input input-bordered focus:input-primary w-full pl-11"
                    placeholder="Поиск по названию, адресу или коду..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>
              </div>

              {/* Сортировка */}
              <div className="flex gap-2">
                <select
                  className="select select-bordered focus:select-primary"
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value as any)}
                >
                  <option value="distance">По расстоянию</option>
                  <option value="name">По названию</option>
                  <option value="working_hours">По времени работы</option>
                </select>

                <button
                  className={`btn ${showMap ? 'btn-primary' : 'btn-outline'}`}
                  onClick={() => setShowMap(!showMap)}
                >
                  <ArrowsPointingOutIcon className="w-5 h-5" />
                  <span className="hidden sm:inline">Карта</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Список отделений */}
      {loading ? (
        <div className="text-center py-12">
          <span className="loading loading-spinner loading-lg"></span>
          <p className="mt-4 text-base-content/60">Загрузка отделений...</p>
        </div>
      ) : filteredOffices.length === 0 ? (
        <div className="text-center py-12">
          <BuildingStorefrontIcon className="w-16 h-16 mx-auto text-base-content/30 mb-4" />
          <h3 className="text-lg font-semibold mb-2">Отделения не найдены</h3>
          <p className="text-base-content/60">
            {searchQuery
              ? `По запросу "${searchQuery}" ничего не найдено`
              : `В городе ${selectedCity} нет доступных отделений`}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {filteredOffices.map((office) => {
            const isSelected = selectedOffice?.id === office.id;
            const isOpen = isOpenNow(office);

            return (
              <div
                key={office.id}
                className={`
                  card cursor-pointer transition-all duration-200 border-2
                  ${
                    isSelected
                      ? 'border-primary shadow-xl scale-[1.02] bg-primary/5'
                      : 'border-transparent hover:border-primary/30 hover:shadow-lg bg-base-100'
                  }
                `}
                onClick={() => handleOfficeSelect(office)}
              >
                <div className="card-body p-6">
                  {/* Заголовок с статусом */}
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h4 className="font-semibold text-lg">{office.name}</h4>
                        {isSelected && (
                          <div className="p-1 bg-success text-success-content rounded-full">
                            <CheckIcon className="w-4 h-4" />
                          </div>
                        )}
                      </div>
                      <div className="text-sm text-base-content/60 mb-2">
                        Код: {office.code}
                      </div>
                    </div>

                    <div
                      className={`badge ${isOpen ? 'badge-success' : 'badge-error'} badge-sm`}
                    >
                      {isOpen ? 'Открыто' : 'Закрыто'}
                    </div>
                  </div>

                  {/* Адрес */}
                  <div className="flex items-start gap-2 mb-4">
                    <MapPinIcon className="w-5 h-5 text-primary mt-0.5 flex-shrink-0" />
                    <div className="flex-1">
                      <div className="font-medium">{office.address}</div>
                      <div className="text-sm text-base-content/60">
                        {office.city}, {office.postal_code}
                      </div>
                    </div>
                    {office.distance && (
                      <div className="text-sm text-primary font-medium">
                        {office.distance} км
                      </div>
                    )}
                  </div>

                  {/* Время работы */}
                  <div className="flex items-center gap-2 mb-4 p-3 bg-base-200/50 rounded-lg">
                    <ClockIcon className="w-5 h-5 text-secondary flex-shrink-0" />
                    <div className="flex-1">
                      <div className="text-sm font-medium">Сегодня:</div>
                      <div className="text-sm text-base-content/70">
                        {getWorkingHoursToday(office)}
                      </div>
                    </div>
                  </div>

                  {/* Контакты */}
                  {office.phone && (
                    <div className="flex items-center gap-2 mb-4">
                      <PhoneIcon className="w-5 h-5 text-accent flex-shrink-0" />
                      <div className="text-sm">{office.phone}</div>
                    </div>
                  )}

                  {/* Услуги */}
                  {office.services && office.services.length > 0 && (
                    <div className="space-y-2">
                      <div className="text-sm font-medium">
                        Доступные услуги:
                      </div>
                      <div className="flex flex-wrap gap-1">
                        {office.services.slice(0, 3).map((service, index) => (
                          <div
                            key={index}
                            className="badge badge-outline badge-sm"
                          >
                            {service}
                          </div>
                        ))}
                        {office.services.length > 3 && (
                          <div className="badge badge-ghost badge-sm">
                            +{office.services.length - 3}
                          </div>
                        )}
                      </div>
                    </div>
                  )}

                  {/* Детали расписания */}
                  {isSelected && (
                    <div className="mt-4 pt-4 border-t space-y-2">
                      <div className="text-sm font-medium">
                        Полное расписание:
                      </div>
                      <div className="grid grid-cols-2 gap-2 text-sm">
                        {Object.entries(office.working_hours).map(
                          ([day, hours]) => (
                            <div key={day} className="flex justify-between">
                              <span className="capitalize text-base-content/70">
                                {day === 'monday' && 'Пн'}
                                {day === 'tuesday' && 'Вт'}
                                {day === 'wednesday' && 'Ср'}
                                {day === 'thursday' && 'Чт'}
                                {day === 'friday' && 'Пт'}
                                {day === 'saturday' && 'Сб'}
                                {day === 'sunday' && 'Вс'}
                              </span>
                              <span className="font-medium">{hours}</span>
                            </div>
                          )
                        )}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Выбранное отделение */}
      {selectedOffice && (
        <div className="alert alert-success">
          <CheckIcon className="w-5 h-5" />
          <div>
            <h4 className="font-semibold">
              Выбрано отделение: {selectedOffice.name}
            </h4>
            <p className="text-sm mt-1">
              {selectedOffice.address}, {selectedOffice.city}
              <br />
              Код отделения: {selectedOffice.code}
            </p>
          </div>
        </div>
      )}

      {/* Информация о времени хранения */}
      <div className="card bg-gradient-to-r from-info/5 to-info/10">
        <div className="card-body p-6">
          <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
            <InformationCircleIcon className="w-5 h-5 text-info" />
            Важная информация
          </h4>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
            <div className="space-y-2">
              <div className="font-medium">⏰ Время хранения:</div>
              <div className="text-base-content/70">
                Посылки хранятся в отделении до 5 рабочих дней бесплатно
              </div>
            </div>

            <div className="space-y-2">
              <div className="font-medium">📄 Документы для получения:</div>
              <div className="text-base-content/70">
                Личный документ (паспорт или ID карта)
              </div>
            </div>

            <div className="space-y-2">
              <div className="font-medium">📞 Уведомления:</div>
              <div className="text-base-content/70">
                SMS о прибытии посылки в отделение
              </div>
            </div>

            <div className="space-y-2">
              <div className="font-medium">💰 Оплата:</div>
              <div className="text-base-content/70">
                Наличными или картой при получении
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
