'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { DollarSign, MapPin, Package, TrendingUp } from 'lucide-react';

interface BaseFiltersProps {
  onFiltersChange: (filters: Record<string, any>) => void;
  className?: string;
}

export const BaseFilters: React.FC<BaseFiltersProps> = ({
  onFiltersChange,
  className = '',
}) => {
  const t = useTranslations('search');

  const [priceMin, setPriceMin] = useState<string>('');
  const [priceMax, setPriceMax] = useState<string>('');
  const [condition, setCondition] = useState<string>('');
  const [sortBy, setSortBy] = useState<string>('relevance');
  const [location, setLocation] = useState<string>('');
  const [distance, setDistance] = useState<number>(0);

  useEffect(() => {
    const filters: Record<string, any> = {};

    if (priceMin) filters.price_min = parseInt(priceMin);
    if (priceMax) filters.price_max = parseInt(priceMax);
    if (condition) filters.condition = condition;
    if (sortBy !== 'relevance') filters.sort_by = sortBy;
    if (location) filters.city = location;
    if (distance > 0) filters.distance = distance;

    onFiltersChange(filters);
  }, [
    priceMin,
    priceMax,
    condition,
    sortBy,
    location,
    distance,
    onFiltersChange,
  ]);

  return (
    <div className={`space-y-4 ${className}`}>
      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <DollarSign className="w-3 h-3" />
          {t('priceRange')}
        </label>
        <div className="flex gap-2">
          <input
            type="number"
            placeholder={t('from')}
            value={priceMin}
            onChange={(e) => setPriceMin(e.target.value)}
            className="input input-bordered input-sm w-full"
          />
          <input
            type="number"
            placeholder={t('to')}
            value={priceMax}
            onChange={(e) => setPriceMax(e.target.value)}
            className="input input-bordered input-sm w-full"
          />
        </div>
      </div>

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <Package className="w-3 h-3" />
          {t('condition')}
        </label>
        <select
          value={condition}
          onChange={(e) => setCondition(e.target.value)}
          className="select select-bordered select-sm w-full"
        >
          <option value="">{t('allConditions')}</option>
          <option value="new">{t('conditions.new')}</option>
          <option value="like_new">{t('conditions.likeNew')}</option>
          <option value="good">{t('conditions.good')}</option>
          <option value="used">{t('conditions.used')}</option>
        </select>
      </div>

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <MapPin className="w-3 h-3" />
          {t('location')}
        </label>
        <input
          type="text"
          placeholder={t('enterCityName')}
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          className="input input-bordered input-sm w-full mb-2"
        />
        {location && (
          <div>
            <label className="text-xs text-base-content/60">
              {t('searchRadius')}: {distance} km
            </label>
            <input
              type="range"
              min="0"
              max="100"
              value={distance}
              onChange={(e) => setDistance(parseInt(e.target.value))}
              className="range range-xs range-primary"
            />
          </div>
        )}
      </div>

      <div>
        <label className="text-xs font-medium text-base-content/70 mb-2 flex items-center gap-1">
          <TrendingUp className="w-3 h-3" />
          {t('sortBy')}
        </label>
        <select
          value={sortBy}
          onChange={(e) => setSortBy(e.target.value)}
          className="select select-bordered select-sm w-full"
        >
          <option value="relevance">{t('sortOptions.relevance')}</option>
          <option value="price_asc">{t('sortOptions.priceAsc')}</option>
          <option value="price_desc">{t('sortOptions.priceDesc')}</option>
          <option value="date">{t('sortOptions.date')}</option>
          <option value="popularity">{t('sortOptions.popularity')}</option>
        </select>
      </div>
    </div>
  );
};
