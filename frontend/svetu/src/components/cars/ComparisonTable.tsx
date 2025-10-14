'use client';

import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '@/store';
import { useTranslations } from 'next-intl';
import Image from 'next/image';
import Link from 'next/link';
import { Check, X, Minus } from 'lucide-react';

interface ComparisonTableProps {
  locale: string;
}

export default function ComparisonTable({ locale }: ComparisonTableProps) {
  const t = useTranslations('cars');
  // Получаем элементы для категории 'cars' из universalCompare
  const items = useSelector(
    (state: RootState) => state.universalCompare.itemsByCategory['cars'] || []
  );

  if (items.length < 2) {
    return (
      <div className="container mx-auto px-4 py-12">
        <div className="alert alert-info">
          <p>{t('compare.minRequiredForTable')}</p>
          <Link
            href={`/${locale}/cars`}
            className="btn btn-primary btn-sm mt-2"
          >
            {t('compare.backToCars')}
          </Link>
        </div>
      </div>
    );
  }

  // Определяем все характеристики для сравнения
  const specifications = [
    { key: 'price', label: t('specs.price'), format: 'price' },
    { key: 'year', label: t('specs.year'), format: 'number' },
    { key: 'mileage', label: t('specs.mileage'), format: 'mileage' },
    { key: 'fuelType', label: t('specs.fuelType'), format: 'text' },
    { key: 'transmission', label: t('specs.transmission'), format: 'text' },
    { key: 'engineSize', label: t('specs.engineSize'), format: 'text' },
    { key: 'power', label: t('specs.power'), format: 'text' },
    { key: 'bodyType', label: t('specs.bodyType'), format: 'text' },
    { key: 'color', label: t('specs.color'), format: 'text' },
    { key: 'driveType', label: t('specs.driveType'), format: 'text' },
    { key: 'doors', label: t('specs.doors'), format: 'number' },
    { key: 'seats', label: t('specs.seats'), format: 'number' },
    { key: 'condition', label: t('specs.condition'), format: 'text' },
    {
      key: 'previousOwners',
      label: t('specs.previousOwners'),
      format: 'number',
    },
    { key: 'warranty', label: t('specs.warranty'), format: 'text' },
    {
      key: 'firstRegistration',
      label: t('specs.firstRegistration'),
      format: 'date',
    },
    {
      key: 'technicalInspection',
      label: t('specs.technicalInspection'),
      format: 'date',
    },
    { key: 'vin', label: t('specs.vin'), format: 'text' },
  ];

  // Форматирование значений
  const formatValue = (value: any, format: string) => {
    if (value === null || value === undefined || value === '') {
      return <Minus className="w-4 h-4 text-base-content/40" />;
    }

    switch (format) {
      case 'price':
        return (
          <span className="font-bold text-primary">
            €{value.toLocaleString()}
          </span>
        );
      case 'mileage':
        return `${value.toLocaleString()} ${t('common.km')}`;
      case 'number':
        return value.toLocaleString();
      case 'date':
        return new Date(value).toLocaleDateString();
      case 'boolean':
        return value ? (
          <Check className="w-5 h-5 text-success" />
        ) : (
          <X className="w-5 h-5 text-error" />
        );
      default:
        return value;
    }
  };

  // Проверка, есть ли различия в характеристике
  const hasDifference = (key: string) => {
    const values = items.map((item) => {
      // price это прямое свойство, остальное в attributes
      if (key === 'price') return item.price;
      return item.attributes?.[key];
    });
    return new Set(values).size > 1;
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">{t('compare.pageTitle')}</h1>

      <div className="overflow-x-auto">
        <table className="table table-zebra w-full">
          <thead>
            <tr>
              <th className="sticky left-0 bg-base-200 z-10">
                {t('compare.specification')}
              </th>
              {items.map((car) => (
                <th key={car.id} className="min-w-[250px]">
                  <div className="space-y-2">
                    {/* Car Image */}
                    {car.image && (
                      <div className="relative h-32 w-full">
                        <Image
                          src={car.image}
                          alt={car.title}
                          fill
                          className="object-cover rounded"
                        />
                      </div>
                    )}

                    {/* Car Title */}
                    <h3 className="font-semibold">{car.title}</h3>

                    {/* View Details Link */}
                    <Link
                      href={`/${locale}/listing/${car.id}`}
                      className="btn btn-primary btn-xs"
                    >
                      {t('actions.viewDetails')}
                    </Link>
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {specifications.map((spec) => {
              const isDifferent = hasDifference(spec.key);
              return (
                <tr
                  key={spec.key}
                  className={isDifferent ? 'bg-warning/10' : ''}
                >
                  <td className="sticky left-0 bg-base-100 font-medium">
                    {spec.label}
                    {isDifferent && (
                      <span className="ml-2 badge badge-warning badge-xs">
                        {t('compare.different')}
                      </span>
                    )}
                  </td>
                  {items.map((car) => (
                    <td key={`${car.id}-${spec.key}`}>
                      {formatValue(
                        spec.key === 'price'
                          ? car.price
                          : car.attributes?.[spec.key],
                        spec.format
                      )}
                    </td>
                  ))}
                </tr>
              );
            })}

            {/* Features Section */}
            <tr>
              <td className="sticky left-0 bg-base-100 font-medium">
                {t('specs.features')}
              </td>
              {items.map((car) => (
                <td key={`${car.id}-features`}>
                  {car.attributes?.features &&
                  Array.isArray(car.attributes.features) &&
                  car.attributes.features.length > 0 ? (
                    <ul className="text-sm space-y-1">
                      {car.attributes.features.map(
                        (feature: string, idx: number) => (
                          <li key={idx} className="flex items-center gap-1">
                            <Check className="w-3 h-3 text-success flex-shrink-0" />
                            {feature}
                          </li>
                        )
                      )}
                    </ul>
                  ) : (
                    <Minus className="w-4 h-4 text-base-content/40" />
                  )}
                </td>
              ))}
            </tr>
          </tbody>
        </table>
      </div>

      {/* Actions */}
      <div className="flex gap-4 mt-8 justify-center">
        <button onClick={() => window.print()} className="btn btn-outline">
          {t('compare.print')}
        </button>
        <Link href={`/${locale}/cars`} className="btn btn-primary">
          {t('compare.backToCars')}
        </Link>
      </div>

      {/* Print Styles */}
      <style jsx global>{`
        @media print {
          .navbar,
          .footer,
          button {
            display: none !important;
          }

          table {
            font-size: 12px;
          }

          .bg-warning\\/10 {
            background-color: #fffbeb !important;
          }
        }
      `}</style>
    </div>
  );
}
