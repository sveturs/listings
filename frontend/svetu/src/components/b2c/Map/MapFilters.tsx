'use client';

import { useState, useCallback } from 'react';
// import { useTranslations } from 'next-intl';
import type { StorefrontFilters, PaymentMethodType } from '@/types/b2c';

interface MapFiltersProps {
  filters: StorefrontFilters;
  onFiltersChange: (filters: Partial<StorefrontFilters>) => void;
  onResetFilters: () => void;
  isLoading?: boolean;
  className?: string;
}

const MapFilters: React.FC<MapFiltersProps> = ({
  filters,
  onFiltersChange,
  onResetFilters,
  isLoading = false,
  className = '',
}) => {
  // const t = useTranslations('storefronts');
  const [isExpanded, setIsExpanded] = useState(false);

  const handleInputChange = useCallback(
    (field: keyof StorefrontFilters, value: any) => {
      onFiltersChange({ [field]: value });
    },
    [onFiltersChange]
  );

  const handlePaymentMethodToggle = useCallback(
    (method: PaymentMethodType) => {
      const currentMethods = filters.paymentMethods || [];
      const newMethods = currentMethods.includes(method)
        ? currentMethods.filter((m) => m !== method)
        : [...currentMethods, method];

      onFiltersChange({ paymentMethods: newMethods });
    },
    [filters.paymentMethods, onFiltersChange]
  );

  const activeFiltersCount = Object.values(filters).filter(
    (value) =>
      value !== undefined &&
      value !== null &&
      value !== '' &&
      (Array.isArray(value) ? value.length > 0 : true)
  ).length;

  return (
    <div className={`bg-white rounded-lg shadow-lg border ${className}`}>
      {/* Заголовок с кнопкой развертывания */}
      <div className="flex items-center justify-between p-4 border-b">
        <div className="flex items-center gap-2">
          <h3 className="font-semibold">Фильтры</h3>
          {activeFiltersCount > 0 && (
            <span className="badge badge-primary">{activeFiltersCount}</span>
          )}
        </div>

        <div className="flex items-center gap-2">
          {activeFiltersCount > 0 && (
            <button
              className="btn btn-ghost btn-sm"
              onClick={onResetFilters}
              disabled={isLoading}
            >
              Сбросить
            </button>
          )}

          <button
            className="btn btn-ghost btn-sm"
            onClick={() => setIsExpanded(!isExpanded)}
          >
            {isExpanded ? '▲' : '▼'}
          </button>
        </div>
      </div>

      {/* Содержимое фильтров */}
      {isExpanded && (
        <div className="p-4 space-y-4">
          {/* Поиск по тексту */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">Поиск</span>
            </label>
            <input
              type="text"
              placeholder="Название витрины..."
              className="input input-bordered"
              value={filters.search || ''}
              onChange={(e) => handleInputChange('search', e.target.value)}
              disabled={isLoading}
            />
          </div>

          {/* Город */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">Город</span>
            </label>
            <input
              type="text"
              placeholder="Название города..."
              className="input input-bordered"
              value={filters.city || ''}
              onChange={(e) => handleInputChange('city', e.target.value)}
              disabled={isLoading}
            />
          </div>

          {/* Рейтинг */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">Минимальный рейтинг</span>
            </label>
            <select
              className="select select-bordered"
              value={filters.minRating || ''}
              onChange={(e) =>
                handleInputChange(
                  'minRating',
                  e.target.value ? Number(e.target.value) : undefined
                )
              }
              disabled={isLoading}
            >
              <option value="">Любой</option>
              <option value="4.5">4.5+ ⭐</option>
              <option value="4.0">4.0+ ⭐</option>
              <option value="3.5">3.5+ ⭐</option>
              <option value="3.0">3.0+ ⭐</option>
            </select>
          </div>

          {/* Статус */}
          <div className="space-y-2">
            <label className="label">
              <span className="label-text">Статус</span>
            </label>

            <div className="flex flex-wrap gap-2">
              <label className="label cursor-pointer">
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm"
                  checked={filters.isActive || false}
                  onChange={(e) =>
                    handleInputChange(
                      'isActive',
                      e.target.checked ? true : undefined
                    )
                  }
                  disabled={isLoading}
                />
                <span className="label-text ml-2">Активные</span>
              </label>

              <label className="label cursor-pointer">
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm"
                  checked={filters.isVerified || false}
                  onChange={(e) =>
                    handleInputChange(
                      'isVerified',
                      e.target.checked ? true : undefined
                    )
                  }
                  disabled={isLoading}
                />
                <span className="label-text ml-2">Проверенные</span>
              </label>

              <label className="label cursor-pointer">
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm"
                  checked={filters.isOpenNow || false}
                  onChange={(e) =>
                    handleInputChange(
                      'isOpenNow',
                      e.target.checked ? true : undefined
                    )
                  }
                  disabled={isLoading}
                />
                <span className="label-text ml-2">Открыто сейчас</span>
              </label>
            </div>
          </div>

          {/* Способы оплаты */}
          <div className="space-y-2">
            <label className="label">
              <span className="label-text">Способы оплаты</span>
            </label>

            <div className="flex flex-wrap gap-2">
              {(
                [
                  'cash',
                  'card',
                  'bank_transfer',
                  'cod',
                  'postanska',
                ] as PaymentMethodType[]
              ).map((method) => (
                <label key={method} className="label cursor-pointer">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-sm"
                    checked={(filters.paymentMethods || []).includes(method)}
                    onChange={() => handlePaymentMethodToggle(method)}
                    disabled={isLoading}
                  />
                  <span className="label-text ml-2">
                    {method === 'cash' && 'Наличные'}
                    {method === 'card' && 'Карта'}
                    {method === 'bank_transfer' && 'Перевод'}
                    {method === 'cod' && 'При доставке'}
                    {method === 'postanska' && 'Pošta'}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {/* Доставка */}
          <div className="space-y-2">
            <label className="label">
              <span className="label-text">Доставка</span>
            </label>

            <div className="flex flex-wrap gap-2">
              <label className="label cursor-pointer">
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm"
                  checked={filters.hasDelivery || false}
                  onChange={(e) =>
                    handleInputChange(
                      'hasDelivery',
                      e.target.checked ? true : undefined
                    )
                  }
                  disabled={isLoading}
                />
                <span className="label-text ml-2">Есть доставка</span>
              </label>

              <label className="label cursor-pointer">
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm"
                  checked={filters.hasSelfPickup || false}
                  onChange={(e) =>
                    handleInputChange(
                      'hasSelfPickup',
                      e.target.checked ? true : undefined
                    )
                  }
                  disabled={isLoading}
                />
                <span className="label-text ml-2">Самовывоз</span>
              </label>
            </div>
          </div>

          {/* Радиус поиска */}
          {filters.latitude && filters.longitude && (
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  Радиус поиска: {filters.radiusKm || 10} км
                </span>
              </label>
              <input
                type="range"
                min="1"
                max="50"
                value={filters.radiusKm || 10}
                className="range range-primary"
                onChange={(e) =>
                  handleInputChange('radiusKm', Number(e.target.value))
                }
                disabled={isLoading}
              />
              <div className="w-full flex justify-between text-xs px-2">
                <span>1км</span>
                <span>25км</span>
                <span>50км</span>
              </div>
            </div>
          )}

          {/* Сортировка */}
          <div className="grid grid-cols-2 gap-2">
            <div className="form-control">
              <label className="label">
                <span className="label-text">Сортировать по</span>
              </label>
              <select
                className="select select-bordered select-sm"
                value={filters.sortBy || ''}
                onChange={(e) =>
                  handleInputChange('sortBy', e.target.value || undefined)
                }
                disabled={isLoading}
              >
                <option value="">По умолчанию</option>
                <option value="rating">Рейтинг</option>
                <option value="distance">Расстояние</option>
                <option value="created_at">Дата создания</option>
                <option value="products_count">Кол-во товаров</option>
              </select>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Порядок</span>
              </label>
              <select
                className="select select-bordered select-sm"
                value={filters.sortOrder || 'desc'}
                onChange={(e) =>
                  handleInputChange(
                    'sortOrder',
                    e.target.value as 'asc' | 'desc'
                  )
                }
                disabled={isLoading}
              >
                <option value="desc">По убыванию</option>
                <option value="asc">По возрастанию</option>
              </select>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default MapFilters;
