'use client';

import { useState, useEffect } from 'react';
import {
  MapPinIcon,
  MagnifyingGlassIcon,
  CheckIcon,
  ExclamationTriangleIcon,
  HomeIcon,
  PhoneIcon,
  UserIcon,
  EnvelopeIcon,
} from '@heroicons/react/24/outline';
import { useTranslations } from 'next-intl';

interface Location {
  id: number;
  name: string;
  postal_code: string;
  region: string;
  country: string;
}

interface Props {
  onAddressChange: (address: any) => void;
  initialAddress?: any;
  deliveryMethod: string;
  className?: string;
}

export default function PostExpressAddressForm({
  onAddressChange,
  initialAddress,
  deliveryMethod,
  className = '',
}: Props) {
  // const t = useTranslations('delivery');
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<Location[]>([]);
  const [selectedLocation, setSelectedLocation] = useState<Location | null>(
    null
  );
  const [loading, setLoading] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);

  const [formData, setFormData] = useState({
    recipient_name: '',
    recipient_phone: '',
    recipient_email: '',
    street_address: '',
    apartment: '',
    floor: '',
    entrance: '',
    postal_code: '',
    city: '',
    note: '',
  });

  useEffect(() => {
    if (initialAddress) {
      setFormData(initialAddress);
      if (initialAddress.city) {
        setSelectedLocation({
          id: 0,
          name: initialAddress.city,
          postal_code: initialAddress.postal_code,
          region: '',
          country: 'Сербия',
        });
      }
    }
  }, [initialAddress]);

  // Поиск населенных пунктов
  const searchLocations = async (query: string) => {
    if (query.length < 2) {
      setSearchResults([]);
      return;
    }

    setLoading(true);
    try {
      const response = await fetch(
        `/api/v1/postexpress/locations/search?q=${encodeURIComponent(query)}`
      );
      const data = await response.json();

      if (data.success) {
        setSearchResults(data.data || []);
        setShowSuggestions(true);
      }
    } catch (error) {
      console.error('Ошибка поиска населенных пунктов:', error);
    } finally {
      setLoading(false);
    }
  };

  // Дебаунс для поиска
  useEffect(() => {
    const timer = setTimeout(() => {
      searchLocations(searchQuery);
    }, 300);

    return () => clearTimeout(timer);
  }, [searchQuery]);

  const handleLocationSelect = (location: Location) => {
    setSelectedLocation(location);
    setSearchQuery(location.name);
    setShowSuggestions(false);

    const updatedData = {
      ...formData,
      city: location.name,
      postal_code: location.postal_code,
    };
    setFormData(updatedData);
    onAddressChange(updatedData);
  };

  const handleInputChange = (field: string, value: string) => {
    const updatedData = { ...formData, [field]: value };
    setFormData(updatedData);
    onAddressChange(updatedData);
  };

  const isValidForm = () => {
    return (
      formData.recipient_name.trim() &&
      formData.recipient_phone.trim() &&
      selectedLocation &&
      (deliveryMethod === 'courier' ? formData.street_address.trim() : true)
    );
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Заголовок */}
      <div className="text-center">
        <h3 className="text-xl font-bold mb-2">
          {deliveryMethod === 'courier'
            ? 'Адрес доставки'
            : 'Данные получателя'}
        </h3>
        <p className="text-base-content/70">
          {deliveryMethod === 'courier'
            ? 'Укажите точный адрес для курьерской доставки'
            : 'Контактные данные для получения в отделении'}
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Основная информация */}
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body p-6">
            <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
              <UserIcon className="w-5 h-5 text-primary" />
              Контактные данные
            </h4>

            <div className="space-y-4">
              {/* Имя получателя */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    Имя и фамилия получателя *
                  </span>
                </label>
                <input
                  type="text"
                  className="input input-bordered focus:input-primary"
                  placeholder="Петар Петрович"
                  value={formData.recipient_name}
                  onChange={(e) =>
                    handleInputChange('recipient_name', e.target.value)
                  }
                />
              </div>

              {/* Телефон */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    Номер телефона *
                  </span>
                </label>
                <div className="relative">
                  <PhoneIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
                  <input
                    type="tel"
                    className="input input-bordered focus:input-primary pl-11"
                    placeholder="+381 60 123 4567"
                    value={formData.recipient_phone}
                    onChange={(e) =>
                      handleInputChange('recipient_phone', e.target.value)
                    }
                  />
                </div>
              </div>

              {/* Email (опционально) */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    Email (опционально)
                  </span>
                </label>
                <div className="relative">
                  <EnvelopeIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
                  <input
                    type="email"
                    className="input input-bordered focus:input-primary pl-11"
                    placeholder="email@example.com"
                    value={formData.recipient_email}
                    onChange={(e) =>
                      handleInputChange('recipient_email', e.target.value)
                    }
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Адрес доставки */}
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body p-6">
            <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
              <MapPinIcon className="w-5 h-5 text-primary" />
              {deliveryMethod === 'courier'
                ? 'Адрес доставки'
                : 'Населенный пункт'}
            </h4>

            <div className="space-y-4">
              {/* Поиск города */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">Город *</span>
                </label>
                <div className="relative">
                  <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
                  <input
                    type="text"
                    className={`input input-bordered focus:input-primary pl-11 ${
                      selectedLocation ? 'input-success' : ''
                    }`}
                    placeholder="Начните вводить название города..."
                    value={searchQuery}
                    onChange={(e) => {
                      setSearchQuery(e.target.value);
                      if (
                        selectedLocation &&
                        e.target.value !== selectedLocation.name
                      ) {
                        setSelectedLocation(null);
                      }
                    }}
                    onFocus={() => setShowSuggestions(searchResults.length > 0)}
                  />
                  {selectedLocation && (
                    <CheckIcon className="absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-success" />
                  )}
                  {loading && (
                    <span className="loading loading-spinner loading-sm absolute right-3 top-1/2 -translate-y-1/2"></span>
                  )}
                </div>

                {/* Результаты поиска */}
                {showSuggestions && searchResults.length > 0 && (
                  <div className="absolute z-50 w-full mt-1 card bg-base-100 shadow-xl max-h-60 overflow-auto">
                    <div className="card-body p-2">
                      {searchResults.map((location) => (
                        <div
                          key={location.id}
                          className="flex items-center gap-3 p-3 hover:bg-base-200 rounded-lg cursor-pointer transition-colors"
                          onClick={() => handleLocationSelect(location)}
                        >
                          <MapPinIcon className="w-5 h-5 text-primary flex-shrink-0" />
                          <div className="flex-1">
                            <div className="font-medium">{location.name}</div>
                            <div className="text-sm text-base-content/60">
                              {location.postal_code} • {location.region}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>

              {/* Адрес для курьерской доставки */}
              {deliveryMethod === 'courier' && (
                <>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text font-medium">
                        Улица и номер дома *
                      </span>
                    </label>
                    <div className="relative">
                      <HomeIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
                      <input
                        type="text"
                        className="input input-bordered focus:input-primary pl-11"
                        placeholder="Кнез Михаилова 42"
                        value={formData.street_address}
                        onChange={(e) =>
                          handleInputChange('street_address', e.target.value)
                        }
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-3">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text font-medium">Квартира</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered focus:input-primary"
                        placeholder="12"
                        value={formData.apartment}
                        onChange={(e) =>
                          handleInputChange('apartment', e.target.value)
                        }
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text font-medium">Этаж</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered focus:input-primary"
                        placeholder="3"
                        value={formData.floor}
                        onChange={(e) =>
                          handleInputChange('floor', e.target.value)
                        }
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text font-medium">Подъезд</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered focus:input-primary"
                        placeholder="А"
                        value={formData.entrance}
                        onChange={(e) =>
                          handleInputChange('entrance', e.target.value)
                        }
                      />
                    </div>
                  </div>
                </>
              )}

              {/* Комментарий */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    Комментарий для курьера
                  </span>
                </label>
                <textarea
                  className="textarea textarea-bordered focus:textarea-primary"
                  placeholder="Дополнительные указания для доставки..."
                  rows={3}
                  value={formData.note}
                  onChange={(e) => handleInputChange('note', e.target.value)}
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Статус валидации */}
      <div
        className={`alert ${isValidForm() ? 'alert-success' : 'alert-warning'}`}
      >
        {isValidForm() ? (
          <>
            <CheckIcon className="w-5 h-5" />
            <span>Все обязательные поля заполнены правильно</span>
          </>
        ) : (
          <>
            <ExclamationTriangleIcon className="w-5 h-5" />
            <span>
              Заполните обязательные поля: имя получателя, телефон, город
              {deliveryMethod === 'courier' && ', адрес доставки'}
            </span>
          </>
        )}
      </div>

      {/* Информация о выбранном населенном пункте */}
      {selectedLocation && (
        <div className="card bg-gradient-to-r from-primary/5 to-secondary/5">
          <div className="card-body p-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 rounded-lg">
                <MapPinIcon className="w-5 h-5 text-primary" />
              </div>
              <div>
                <div className="font-semibold">{selectedLocation.name}</div>
                <div className="text-sm text-base-content/70">
                  Почтовый индекс: {selectedLocation.postal_code}
                  {selectedLocation.region && ` • ${selectedLocation.region}`}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
