'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import {
  Car,
  Home,
  ShoppingBag,
  Wrench,
  Shirt,
  Monitor,
  Sofa,
  MoreHorizontal,
} from 'lucide-react';

interface CategorySelectorProps {
  selectedCategoryId?: number;
  onCategorySelect: (categoryId: number | undefined) => void;
  className?: string;
}

const popularCategories = [
  { id: 1003, icon: Car, labelKey: 'automotive' },
  { id: 1004, icon: Home, labelKey: 'realEstate' },
  { id: 1002, icon: Monitor, labelKey: 'electronics' },
  { id: 1001, icon: Shirt, labelKey: 'clothing' },
  { id: 1005, icon: Sofa, labelKey: 'furniture' },
  { id: 1006, icon: Wrench, labelKey: 'tools' },
];

export const CategorySelector: React.FC<CategorySelectorProps> = ({
  selectedCategoryId,
  onCategorySelect,
  className = '',
}) => {
  const t = useTranslations('search');

  const handleCategoryClick = (categoryId: number) => {
    if (selectedCategoryId === categoryId) {
      onCategorySelect(undefined);
    } else {
      onCategorySelect(categoryId);
    }
  };

  return (
    <div className={`${className}`}>
      <div className="mb-4">
        <h3 className="text-sm font-semibold text-base-content/70 mb-3 flex items-center gap-2">
          <ShoppingBag className="w-4 h-4" />
          {t('popularCategories')}
        </h3>
        <div className="grid grid-cols-3 gap-2">
          {popularCategories.map(({ id, icon: Icon, labelKey }) => (
            <button
              key={id}
              onClick={() => handleCategoryClick(id)}
              className={`
                btn btn-sm flex flex-col gap-1 h-auto py-3 transition-all duration-200
                hover:scale-105 hover:shadow-lg
                ${
                  selectedCategoryId === id
                    ? 'btn-primary shadow-lg scale-105'
                    : 'btn-ghost border border-base-300 hover:border-primary'
                }
              `}
            >
              <Icon
                className={`w-5 h-5 ${selectedCategoryId === id ? 'animate-pulse' : ''}`}
              />
              <span className="text-xs">{t(`categories.${labelKey}`)}</span>
            </button>
          ))}
        </div>
      </div>

      {selectedCategoryId && (
        <div className="flex items-center justify-between p-2 bg-primary/10 rounded-lg animate-fadeIn">
          <span className="text-sm font-medium">
            {t('categorySelected', {
              category: t(
                `categories.${
                  popularCategories.find((c) => c.id === selectedCategoryId)
                    ?.labelKey
                }`
              ),
            })}
          </span>
          <button
            onClick={() => onCategorySelect(undefined)}
            className="btn btn-ghost btn-xs transition-transform hover:scale-110"
          >
            {t('clear')}
          </button>
        </div>
      )}

      <button className="btn btn-ghost btn-sm w-full mt-3">
        <MoreHorizontal className="w-4 h-4" />
        {t('allCategories')}
      </button>
    </div>
  );
};
