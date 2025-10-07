'use client';

import React from 'react';
import { X, ChevronRight, LayoutGrid } from 'lucide-react';
import { useLocale, useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';
import { renderCategoryIcon } from '@/utils/iconMapper';

type MarketplaceCategory =
  components['schemas']['models.MarketplaceCategory'];

interface CategoryTreeNode extends MarketplaceCategory {
  children: CategoryTreeNode[];
}

interface CategoryModalProps {
  isOpen: boolean;
  onClose: () => void;
  categories: CategoryTreeNode[];
  expandedCategories: Set<number>;
  onToggleExpanded: (categoryId: number) => void;
  selectedCategoryId?: number | null;
  onCategorySelect: (categoryId: number | null) => void;
  title: string;
}

export const CategoryModal: React.FC<CategoryModalProps> = ({
  isOpen,
  onClose,
  categories,
  expandedCategories,
  onToggleExpanded,
  selectedCategoryId,
  onCategorySelect,
  title,
}) => {
  const locale = useLocale();
  const t = useTranslations('marketplace');

  if (!isOpen) return null;

  // Функция для подсчета общего количества товаров во всех категориях
  const getTotalListingsCount = (cats: CategoryTreeNode[]): number => {
    return cats.reduce((total, cat) => {
      const catCount = cat.listing_count || 0;
      const childrenCount = cat.children
        ? getTotalListingsCount(cat.children)
        : 0;
      return total + catCount + childrenCount;
    }, 0);
  };

  const handleCategoryClick = (categoryId: number) => {
    onCategorySelect(categoryId);
    onClose();
  };

  const renderCategory = (category: CategoryTreeNode, level: number = 0) => {
    const isExpanded = expandedCategories.has(category.id!);
    const hasChildren = category.children.length > 0;
    const isSelected = selectedCategoryId === category.id;

    return (
      <div key={category.id} className="w-full">
        <div
          className={`flex items-center gap-2 p-3 rounded-lg cursor-pointer transition-all hover:bg-base-200 ${
            isSelected ? 'bg-primary/10 text-primary font-medium' : ''
          }`}
          style={{ paddingLeft: `${level * 16 + 12}px` }}
          onClick={() => handleCategoryClick(category.id!)}
        >
          {hasChildren && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                onToggleExpanded(category.id!);
              }}
              className="btn btn-ghost btn-xs p-0 min-h-0 h-5 w-5"
            >
              <ChevronRight
                className={`w-4 h-4 transition-transform ${
                  isExpanded ? 'rotate-90' : ''
                }`}
              />
            </button>
          )}

          {!hasChildren && <div className="w-5" />}

          {renderCategoryIcon(category.icon, 'w-5 h-5 text-base-content/70')}

          <span className="flex-1">
            {category.translations && category.translations[locale]
              ? category.translations[locale]
              : category.name}
          </span>

          {category.listing_count !== undefined &&
            category.listing_count > 0 && (
              <span className="badge badge-ghost badge-sm">
                {category.listing_count}
              </span>
            )}
        </div>

        {hasChildren && isExpanded && (
          <div className="mt-1">
            {category.children.map((child) => renderCategory(child, level + 1))}
          </div>
        )}
      </div>
    );
  };

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50 z-50 lg:hidden"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="fixed inset-x-0 bottom-0 z-50 lg:hidden">
        <div className="bg-base-100 rounded-t-3xl shadow-xl max-h-[85vh] flex flex-col">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b border-base-200">
            <h3 className="text-lg font-bold">{title}</h3>
            <button
              onClick={onClose}
              className="btn btn-ghost btn-sm btn-circle"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          {/* Content */}
          <div className="flex-1 overflow-y-auto p-4">
            <div className="space-y-1">
              {/* Кнопка "Все категории" */}
              <div
                className={`flex items-center gap-2 p-3 rounded-lg cursor-pointer transition-all hover:bg-base-200 ${
                  !selectedCategoryId
                    ? 'bg-primary/10 text-primary font-medium'
                    : ''
                }`}
                onClick={() => {
                  onCategorySelect(null);
                  onClose();
                }}
              >
                <div className="btn btn-ghost btn-xs p-0 min-h-0 h-5 w-5 invisible" />
                <div className="w-5" />
                <LayoutGrid className="w-5 h-5 text-base-content/70" />
                <span className="flex-1 font-semibold">
                  {t('allCategories')}
                </span>
                <span className="badge badge-primary badge-sm">
                  {getTotalListingsCount(categories)}
                </span>
              </div>

              <div className="divider my-2" />

              {categories.length === 0 ? (
                <p className="text-base-content/60 text-center py-8">
                  Нет доступных категорий
                </p>
              ) : (
                categories.map((category) => renderCategory(category))
              )}
            </div>
          </div>
        </div>
      </div>
    </>
  );
};
