'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
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
import { CategoryTreeModal } from './CategoryTreeModal';
import { MarketplaceService } from '@/services/c2c';

interface CategorySelectorProps {
  selectedCategoryId?: number;
  onCategorySelect: (categoryId: number | undefined) => void;
  className?: string;
}

const popularCategories = [
  { id: 1003, icon: Car, labelKey: 'automotive' },
  { id: 1004, icon: Home, labelKey: 'realEstate' },
  { id: 1001, icon: Monitor, labelKey: 'electronics' },
  { id: 1002, icon: Shirt, labelKey: 'clothing' },
  { id: 1005, icon: Sofa, labelKey: 'furniture' },
  { id: 1007, icon: Wrench, labelKey: 'tools' },
];

export const CategorySelector: React.FC<CategorySelectorProps> = ({
  selectedCategoryId,
  onCategorySelect,
  className = '',
}) => {
  const t = useTranslations('search');
  const locale = useLocale();
  const [showCategoryModal, setShowCategoryModal] = useState(false);
  const [selectedCategoryName, setSelectedCategoryName] = useState<string>('');

  // Load category name when selectedCategoryId changes
  useEffect(() => {
    const loadCategoryName = async () => {
      if (!selectedCategoryId) {
        setSelectedCategoryName('');
        return;
      }

      // Check if it's a popular category first
      const popularCategory = popularCategories.find(
        (c) => c.id === selectedCategoryId
      );
      if (popularCategory) {
        setSelectedCategoryName(t(`categories.${popularCategory.labelKey}`));
        return;
      }

      // If not popular, load from API
      try {
        const response = await MarketplaceService.getCategories(locale);
        const category = response.data.find(
          (cat) => cat.id === selectedCategoryId
        );
        if (category) {
          setSelectedCategoryName(category.name);
        } else {
          setSelectedCategoryName(t('categories.other'));
        }
      } catch (error) {
        console.error('Failed to load category:', error);
        setSelectedCategoryName(t('categories.other'));
      }
    };

    loadCategoryName();
  }, [selectedCategoryId, locale, t]);

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
              category: selectedCategoryName || t('categories.other'),
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

      <button
        className="btn btn-ghost btn-sm w-full mt-3"
        onClick={() => setShowCategoryModal(true)}
      >
        <MoreHorizontal className="w-4 h-4" />
        {t('allCategories')}
      </button>

      <CategoryTreeModal
        isOpen={showCategoryModal}
        onClose={() => setShowCategoryModal(false)}
        selectedCategoryId={selectedCategoryId}
        onCategorySelect={onCategorySelect}
      />
    </div>
  );
};
