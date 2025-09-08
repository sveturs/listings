'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { debounce } from 'lodash';
import {
  UserIcon,
  PhoneIcon,
  EnvelopeIcon,
  MapPinIcon,
  HomeIcon,
  BuildingOfficeIcon,
  ChatBubbleBottomCenterTextIcon,
  CheckCircleIcon,
} from '@heroicons/react/24/outline';
import { configManager } from '@/config';

// Схема валидации адреса
const addressSchema = z.object({
  recipient_name: z.string().min(2, 'Имя получателя обязательно'),
  recipient_phone: z
    .string()
    .min(10, 'Телефон должен содержать минимум 10 цифр')
    .regex(/^(\+381|0)[0-9]{8,11}$/, 'Неверный формат телефона'),
  recipient_email: z
    .string()
    .email('Неверный email')
    .optional()
    .or(z.literal('')),
  city: z.string().min(2, 'Город обязателен'),
  postal_code: z
    .string()
    .regex(/^[0-9]{5}$/, 'Почтовый индекс должен содержать 5 цифр'),
  street_address: z.string().optional(),
  street_number: z.string().optional(),
  apartment: z.string().optional(),
  floor: z.string().optional(),
  entrance: z.string().optional(),
  note: z.string().optional(),
  // BEX specific fields
  municipality_id: z.number().optional(),
  place_id: z.number().optional(),
  street_id: z.number().optional(),
});

export type BEXAddressData = z.infer<typeof addressSchema>;

interface AddressSuggestion {
  id: number;
  name: string;
  type: 'municipality' | 'place' | 'street';
  postal_code?: string;
  municipality_name?: string;
  place_name?: string;
}

interface Props {
  onAddressChange: (address: BEXAddressData) => void;
  initialAddress?: Partial<BEXAddressData>;
  deliveryMethod?: string;
  className?: string;
}

export default function BEXAddressForm({
  onAddressChange,
  initialAddress,
  deliveryMethod = 'courier',
  className = '',
}: Props) {
  const [addressSuggestions, setAddressSuggestions] = useState<
    AddressSuggestion[]
  >([]);
  const [isSearching, setIsSearching] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedCity, setSelectedCity] = useState<string>(
    initialAddress?.city || ''
  );
  const [_selectedPlace, setSelectedPlace] = useState<AddressSuggestion | null>(
    null
  );

  const form = useForm<BEXAddressData>({
    resolver: zodResolver(addressSchema),
    defaultValues: {
      recipient_name: '',
      recipient_phone: '',
      recipient_email: '',
      city: '',
      postal_code: '',
      street_address: '',
      street_number: '',
      apartment: '',
      floor: '',
      entrance: '',
      note: '',
      ...initialAddress,
    },
    mode: 'onChange',
  });

  // Поиск адреса через BEX API
  const searchAddressBase = useCallback(
    async (query: string, city?: string) => {
      if (query.length < 2) {
        setAddressSuggestions([]);
        return;
      }

      setIsSearching(true);
      try {
        const apiUrl = configManager.get('api.url');
        const response = await fetch(`${apiUrl}/api/v1/bex/search-address`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ query, city, limit: 10 }),
        });

        if (response.ok) {
          const data = await response.json();
          if (data.success && data.data) {
            setAddressSuggestions(data.data);
            setShowSuggestions(true);
          }
        }
      } catch (error) {
        console.error('Address search failed:', error);
      } finally {
        setIsSearching(false);
      }
    },
    []
  );

  // Debounced версия функции поиска
  const searchAddress = useMemo(
    () => debounce(searchAddressBase, 300),
    [searchAddressBase]
  );

  // Обработчик изменения города
  const handleCityChange = (value: string) => {
    setSelectedCity(value);
    form.setValue('city', value);
    searchAddress(value);
  };

  // Обработчик изменения улицы
  const handleStreetChange = (value: string) => {
    form.setValue('street_address', value);
    if (selectedCity) {
      searchAddress(value, selectedCity);
    }
  };

  // Выбор подсказки адреса
  const selectSuggestion = (suggestion: AddressSuggestion) => {
    if (suggestion.type === 'place') {
      setSelectedPlace(suggestion);
      form.setValue('city', suggestion.name);
      form.setValue('postal_code', suggestion.postal_code || '');
      form.setValue('place_id', suggestion.id);
      setSelectedCity(suggestion.name);
    } else if (suggestion.type === 'street') {
      form.setValue('street_address', suggestion.name);
      form.setValue('street_id', suggestion.id);
    } else if (suggestion.type === 'municipality') {
      form.setValue('municipality_id', suggestion.id);
    }

    setShowSuggestions(false);
  };

  // Отправка изменений в родительский компонент
  useEffect(() => {
    const subscription = form.watch((data) => {
      if (form.formState.isValid) {
        onAddressChange(data as BEXAddressData);
      }
    });
    return () => subscription.unsubscribe();
  }, [form, onAddressChange]);

  const isAddressRequired = deliveryMethod === 'courier';

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Информация о получателе */}
      <div>
        <h4 className="font-semibold mb-4 flex items-center gap-2">
          <UserIcon className="w-5 h-5 text-primary" />
          Данные получателя
        </h4>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="form-control">
            <label className="label">
              <span className="label-text">ФИО получателя *</span>
            </label>
            <div className="relative">
              <UserIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
              <input
                type="text"
                className={`input input-bordered w-full pl-11 ${
                  form.formState.errors.recipient_name ? 'input-error' : ''
                }`}
                placeholder="Петар Петровић"
                {...form.register('recipient_name')}
              />
            </div>
            {form.formState.errors.recipient_name && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {form.formState.errors.recipient_name.message}
                </span>
              </label>
            )}
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Телефон *</span>
            </label>
            <div className="relative">
              <PhoneIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
              <input
                type="tel"
                className={`input input-bordered w-full pl-11 ${
                  form.formState.errors.recipient_phone ? 'input-error' : ''
                }`}
                placeholder="+381 XX XXX XXXX"
                {...form.register('recipient_phone')}
              />
            </div>
            {form.formState.errors.recipient_phone && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {form.formState.errors.recipient_phone.message}
                </span>
              </label>
            )}
          </div>

          <div className="form-control md:col-span-2">
            <label className="label">
              <span className="label-text">Email</span>
              <span className="label-text-alt">Необязательно</span>
            </label>
            <div className="relative">
              <EnvelopeIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
              <input
                type="email"
                className={`input input-bordered w-full pl-11 ${
                  form.formState.errors.recipient_email ? 'input-error' : ''
                }`}
                placeholder="email@example.com"
                {...form.register('recipient_email')}
              />
            </div>
            {form.formState.errors.recipient_email && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {form.formState.errors.recipient_email.message}
                </span>
              </label>
            )}
          </div>
        </div>
      </div>

      {/* Адрес доставки */}
      {isAddressRequired && (
        <div>
          <h4 className="font-semibold mb-4 flex items-center gap-2">
            <MapPinIcon className="w-5 h-5 text-primary" />
            Адрес доставки
          </h4>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Город с автодополнением */}
            <div className="form-control relative">
              <label className="label">
                <span className="label-text">Город *</span>
              </label>
              <div className="relative">
                <MapPinIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40 z-10" />
                <input
                  type="text"
                  className={`input input-bordered w-full pl-11 pr-10 ${
                    form.formState.errors.city ? 'input-error' : ''
                  }`}
                  placeholder="Начните вводить название города"
                  value={form.watch('city')}
                  onChange={(e) => handleCityChange(e.target.value)}
                  onFocus={() => setShowSuggestions(true)}
                />
                {isSearching && (
                  <span className="absolute right-3 top-1/2 -translate-y-1/2 loading loading-spinner loading-sm"></span>
                )}
              </div>

              {/* Подсказки городов */}
              {showSuggestions && addressSuggestions.length > 0 && (
                <div className="absolute top-full left-0 right-0 z-50 mt-1 max-h-60 overflow-auto bg-base-100 border border-base-300 rounded-lg shadow-lg">
                  {addressSuggestions
                    .filter((s) => s.type === 'place')
                    .map((suggestion) => (
                      <button
                        key={suggestion.id}
                        type="button"
                        className="w-full px-4 py-2 text-left hover:bg-base-200 flex items-center justify-between"
                        onClick={() => selectSuggestion(suggestion)}
                      >
                        <div>
                          <div className="font-medium">{suggestion.name}</div>
                          {suggestion.postal_code && (
                            <div className="text-sm text-base-content/60">
                              {suggestion.postal_code} •{' '}
                              {suggestion.municipality_name}
                            </div>
                          )}
                        </div>
                        <CheckCircleIcon className="w-5 h-5 text-success opacity-0 hover:opacity-100" />
                      </button>
                    ))}
                </div>
              )}

              {form.formState.errors.city && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {form.formState.errors.city.message}
                  </span>
                </label>
              )}
            </div>

            {/* Почтовый индекс */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">Почтовый индекс *</span>
              </label>
              <input
                type="text"
                className={`input input-bordered w-full ${
                  form.formState.errors.postal_code ? 'input-error' : ''
                }`}
                placeholder="21000"
                maxLength={5}
                {...form.register('postal_code')}
              />
              {form.formState.errors.postal_code && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {form.formState.errors.postal_code.message}
                  </span>
                </label>
              )}
            </div>

            {/* Улица с автодополнением */}
            <div className="form-control relative">
              <label className="label">
                <span className="label-text">Улица *</span>
              </label>
              <div className="relative">
                <HomeIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40 z-10" />
                <input
                  type="text"
                  className={`input input-bordered w-full pl-11 ${
                    form.formState.errors.street_address ? 'input-error' : ''
                  }`}
                  placeholder="Начните вводить название улицы"
                  value={form.watch('street_address')}
                  onChange={(e) => handleStreetChange(e.target.value)}
                  disabled={!selectedCity}
                />
              </div>

              {/* Подсказки улиц */}
              {showSuggestions &&
                selectedCity &&
                addressSuggestions.length > 0 && (
                  <div className="absolute top-full left-0 right-0 z-50 mt-1 max-h-60 overflow-auto bg-base-100 border border-base-300 rounded-lg shadow-lg">
                    {addressSuggestions
                      .filter((s) => s.type === 'street')
                      .map((suggestion) => (
                        <button
                          key={suggestion.id}
                          type="button"
                          className="w-full px-4 py-2 text-left hover:bg-base-200"
                          onClick={() => selectSuggestion(suggestion)}
                        >
                          <div className="font-medium">{suggestion.name}</div>
                          {suggestion.place_name && (
                            <div className="text-sm text-base-content/60">
                              {suggestion.place_name}
                            </div>
                          )}
                        </button>
                      ))}
                  </div>
                )}
            </div>

            {/* Номер дома */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">Номер дома</span>
              </label>
              <input
                type="text"
                className="input input-bordered w-full"
                placeholder="15а"
                {...form.register('street_number')}
              />
            </div>

            {/* Дополнительные поля */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">Квартира/Офис</span>
              </label>
              <div className="relative">
                <BuildingOfficeIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
                <input
                  type="text"
                  className="input input-bordered w-full pl-11"
                  placeholder="Кв. 25"
                  {...form.register('apartment')}
                />
              </div>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Этаж</span>
              </label>
              <input
                type="text"
                className="input input-bordered w-full"
                placeholder="3"
                {...form.register('floor')}
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Подъезд</span>
              </label>
              <input
                type="text"
                className="input input-bordered w-full"
                placeholder="2"
                {...form.register('entrance')}
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Комментарий курьеру</span>
              </label>
              <div className="relative">
                <ChatBubbleBottomCenterTextIcon className="absolute left-3 top-3 w-5 h-5 text-base-content/40" />
                <textarea
                  className="textarea textarea-bordered w-full pl-11"
                  rows={2}
                  placeholder="Позвонить за 30 минут, домофон не работает..."
                  {...form.register('note')}
                />
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Статус валидации */}
      {form.formState.isValid && (
        <div className="alert alert-success">
          <CheckCircleIcon className="w-5 h-5" />
          <span>Все данные заполнены корректно</span>
        </div>
      )}
    </div>
  );
}
