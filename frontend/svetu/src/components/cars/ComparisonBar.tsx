'use client';

import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState, AppDispatch } from '@/store';
import {
  removeFromCompare,
  clearCompare,
  toggleComparePanel,
  initializeCompare,
} from '@/store/slices/compareSlice';
import { useTranslations } from 'next-intl';
import { X, Scale, ChevronUp, ChevronDown } from 'lucide-react';
import Image from 'next/image';
import { useRouter } from 'next/navigation';

export default function ComparisonBar() {
  const dispatch = useDispatch<AppDispatch>();
  const router = useRouter();
  const t = useTranslations('cars');

  const { items, isOpen, maxItems } = useSelector(
    (state: RootState) => state.compare
  );

  useEffect(() => {
    dispatch(initializeCompare());
  }, [dispatch]);

  if (items.length === 0) {
    return null;
  }

  const handleCompareClick = () => {
    if (items.length >= 2) {
      // Navigate to comparison page
      const ids = items.map((item) => item.id).join(',');
      router.push(`/cars/compare?ids=${ids}`);
    }
  };

  return (
    <>
      {/* Floating button for mobile */}
      <div className="fixed bottom-20 right-4 z-40 lg:hidden">
        <button
          onClick={() => dispatch(toggleComparePanel())}
          className="btn btn-primary btn-circle shadow-lg relative"
        >
          <Scale className="w-5 h-5" />
          {items.length > 0 && (
            <span className="absolute -top-2 -right-2 badge badge-error badge-sm">
              {items.length}
            </span>
          )}
        </button>
      </div>

      {/* Comparison bar */}
      <div
        className={`fixed bottom-0 left-0 right-0 bg-base-100 border-t border-base-300 shadow-2xl z-30 transition-transform duration-300 ${
          isOpen
            ? 'translate-y-0'
            : 'translate-y-full lg:translate-y-[calc(100%-4rem)]'
        }`}
      >
        {/* Header */}
        <div
          className="flex items-center justify-between p-4 border-b border-base-300 cursor-pointer lg:cursor-default"
          onClick={() => dispatch(toggleComparePanel())}
        >
          <div className="flex items-center gap-3">
            <Scale className="w-5 h-5 text-primary" />
            <span className="font-semibold">
              {t('compare.title')} ({items.length}/{maxItems})
            </span>
          </div>

          <div className="flex items-center gap-2">
            {items.length >= 2 && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  handleCompareClick();
                }}
                className="btn btn-primary btn-sm"
              >
                {t('compare.compareNow')}
              </button>
            )}

            <button
              onClick={(e) => {
                e.stopPropagation();
                dispatch(clearCompare());
              }}
              className="btn btn-ghost btn-sm"
            >
              {t('compare.clearAll')}
            </button>

            <button className="btn btn-ghost btn-sm btn-circle lg:hidden">
              {isOpen ? (
                <ChevronDown className="w-4 h-4" />
              ) : (
                <ChevronUp className="w-4 h-4" />
              )}
            </button>
          </div>
        </div>

        {/* Content */}
        {isOpen && (
          <div className="p-4 max-h-64 overflow-y-auto">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {items.map((car) => (
                <div
                  key={car.id}
                  className="flex items-center gap-3 p-3 bg-base-200 rounded-lg relative"
                >
                  {/* Remove button */}
                  <button
                    onClick={() => dispatch(removeFromCompare(car.id))}
                    className="absolute top-2 right-2 btn btn-ghost btn-xs btn-circle"
                  >
                    <X className="w-3 h-3" />
                  </button>

                  {/* Car image */}
                  {car.imageUrl && (
                    <div className="w-16 h-16 flex-shrink-0">
                      <Image
                        src={car.imageUrl}
                        alt={car.title}
                        width={64}
                        height={64}
                        className="w-full h-full object-cover rounded"
                      />
                    </div>
                  )}

                  {/* Car info */}
                  <div className="flex-1 min-w-0">
                    <h4 className="font-medium truncate">{car.title}</h4>
                    <p className="text-sm text-base-content/70">
                      {car.year} •{' '}
                      {car.mileage
                        ? `${car.mileage.toLocaleString()} ${t('common.km')}`
                        : t('common.noMileage')}
                    </p>
                    <p className="text-sm font-semibold text-primary">
                      €{car.price.toLocaleString()}
                    </p>
                  </div>
                </div>
              ))}

              {/* Add more cars prompt */}
              {items.length < maxItems && (
                <div className="flex items-center justify-center p-3 border-2 border-dashed border-base-300 rounded-lg min-h-[100px]">
                  <div className="text-center">
                    <p className="text-sm text-base-content/70">
                      {t('compare.addMore', { count: maxItems - items.length })}
                    </p>
                  </div>
                </div>
              )}
            </div>

            {/* Comparison requirements */}
            {items.length < 2 && (
              <div className="mt-4 text-center">
                <p className="text-sm text-warning">
                  {t('compare.minRequired')}
                </p>
              </div>
            )}
          </div>
        )}
      </div>
    </>
  );
}
