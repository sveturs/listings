'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import {
  Calendar,
  Gauge,
  Fuel,
  Car,
  Cog,
  Palette,
  Package,
  Users,
} from 'lucide-react';

interface CarQuickInfoProps {
  year?: string | number;
  mileage?: string | number;
  fuelType?: string;
  transmission?: string;
  engineSize?: string | number;
  color?: string;
  bodyType?: string;
  seats?: string | number;
  compact?: boolean;
}

export const CarQuickInfo: React.FC<CarQuickInfoProps> = ({
  year,
  mileage,
  fuelType,
  transmission,
  engineSize,
  color,
  bodyType,
  seats,
  compact = false,
}) => {
  const t = useTranslations('cars');

  const formatMileage = (value: string | number) => {
    const num = typeof value === 'string' ? parseInt(value) : value;
    if (isNaN(num)) return value;
    return num.toLocaleString() + ' km';
  };

  const formatFuel = (fuel: string) => {
    const fuelTypes: Record<string, string> = {
      petrol: t('filters.petrol'),
      diesel: t('filters.diesel'),
      electric: t('filters.electric'),
      hybrid: t('filters.hybrid'),
      lpg: t('filters.lpg'),
      gas: t('filters.gas'),
    };
    return fuelTypes[fuel.toLowerCase()] || fuel;
  };

  const formatTransmission = (trans: string) => {
    const types: Record<string, string> = {
      manual: t('filters.manual'),
      automatic: t('filters.automatic'),
      semiautomatic: t('filters.semiAutomatic'),
      cvt: 'CVT',
    };
    return types[trans.toLowerCase()] || trans;
  };

  const formatEngineSize = (size: string | number) => {
    if (!size) return null;
    return typeof size === 'number' ? `${size}L` : size;
  };

  const items = [];

  if (year) {
    items.push({
      icon: Calendar,
      value: year,
      label: t('info.year'),
      color: 'text-primary',
    });
  }

  if (mileage) {
    items.push({
      icon: Gauge,
      value: formatMileage(mileage),
      label: t('info.mileage'),
      color: 'text-info',
    });
  }

  if (fuelType) {
    items.push({
      icon: Fuel,
      value: formatFuel(fuelType),
      label: t('info.fuel'),
      color: 'text-success',
    });
  }

  if (transmission) {
    items.push({
      icon: Car,
      value: formatTransmission(transmission),
      label: t('info.transmission'),
      color: 'text-warning',
    });
  }

  if (engineSize) {
    items.push({
      icon: Cog,
      value: formatEngineSize(engineSize),
      label: t('info.engine'),
      color: 'text-error',
    });
  }

  if (color) {
    items.push({
      icon: Palette,
      value: color,
      label: t('info.color'),
      color: 'text-secondary',
    });
  }

  if (bodyType) {
    items.push({
      icon: Package,
      value: bodyType,
      label: t('info.bodyType'),
      color: 'text-accent',
    });
  }

  if (seats) {
    items.push({
      icon: Users,
      value: seats,
      label: t('info.seats'),
      color: 'text-neutral',
    });
  }

  if (items.length === 0) return null;

  if (compact) {
    return (
      <div className="flex flex-wrap gap-3">
        {items.slice(0, 4).map((item, index) => {
          const Icon = item.icon;
          return (
            <div key={index} className="flex items-center gap-1.5">
              <Icon className={`w-4 h-4 ${item.color}`} />
              <span className="text-sm font-medium">{item.value}</span>
            </div>
          );
        })}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
      {items.map((item, index) => {
        const Icon = item.icon;
        return (
          <div
            key={index}
            className="flex flex-col items-center p-2 rounded-lg bg-base-200/50"
          >
            <Icon className={`w-5 h-5 ${item.color} mb-1`} />
            <span className="text-xs text-base-content/60">{item.label}</span>
            <span className="text-sm font-semibold">{item.value}</span>
          </div>
        );
      })}
    </div>
  );
};

export default CarQuickInfo;
