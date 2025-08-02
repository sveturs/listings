'use client';

import React, { useState, useEffect } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import { MarketplaceService } from '@/services/marketplace';
import { ChevronRight, Grid3X3, Filter } from 'lucide-react';
import type { components } from '@/types/generated/api';
import { CategoryModal } from './CategoryModal';
import { FilterModal } from './FilterModal';
import { DesktopFilters } from './DesktopFilters';
import { renderCategoryIcon } from '@/utils/iconMapper';

type MarketplaceCategory =
  components['schemas']['backend_internal_domain_models.MarketplaceCategory'];

interface CategoryTreeNode extends MarketplaceCategory {
  children: CategoryTreeNode[];
}

interface BentoGridCategoriesProps {
  onCategorySelect: (categoryId: number | null) => void;
  selectedCategoryId?: number | null;
  filters?: Record<string, any>;
  onFiltersChange?: (filters: Record<string, any>) => void;
}

export const BentoGridCategories: React.FC<BentoGridCategoriesProps> = ({
  onCategorySelect,
  selectedCategoryId,
  filters = {},
  onFiltersChange,
}) => {
  const locale = useLocale();
  const t = useTranslations('marketplace');
  const [categories, setCategories] = useState<CategoryTreeNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [expandedCategories, setExpandedCategories] = useState<Set<number>>(
    new Set()
  );
  const [showCategoryModal, setShowCategoryModal] = useState(false);
  const [showFilterModal, setShowFilterModal] = useState(false);

  // Построение дерева категорий
  const buildCategoryTree = (
    categories: MarketplaceCategory[]
  ): CategoryTreeNode[] => {
    const categoryMap = new Map<number, CategoryTreeNode>();
    const rootCategories: CategoryTreeNode[] = [];

    categories.forEach((category) => {
      categoryMap.set(category.id!, {
        ...category,
        children: [],
      });
    });

    categories.forEach((category) => {
      const node = categoryMap.get(category.id!)!;
      if (category.parent_id && categoryMap.has(category.parent_id)) {
        const parent = categoryMap.get(category.parent_id)!;
        parent.children.push(node);
      } else {
        rootCategories.push(node);
      }
    });

    const sortCategories = (cats: CategoryTreeNode[]) => {
      cats.sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0));
      cats.forEach((cat) => sortCategories(cat.children));
    };

    sortCategories(rootCategories);
    return rootCategories;
  };

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        setLoading(true);
        const response = await MarketplaceService.getCategories(locale);
        if (response.success && response.data) {
          const mappedCategories: MarketplaceCategory[] = response.data.map(
            (cat) => ({
              ...cat,
              parent_id: cat.parent_id === null ? undefined : cat.parent_id,
            })
          );
          const tree = buildCategoryTree(mappedCategories);
          setCategories(tree);
        }
      } catch (err) {
        console.error('Error fetching categories:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchCategories();
  }, [locale]);

  const toggleExpanded = (categoryId: number) => {
    const newExpanded = new Set(expandedCategories);
    if (newExpanded.has(categoryId)) {
      newExpanded.delete(categoryId);
    } else {
      newExpanded.add(categoryId);
    }
    setExpandedCategories(newExpanded);
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
          onClick={() => onCategorySelect(category.id!)}
        >
          {hasChildren && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleExpanded(category.id!);
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

  if (loading) {
    return (
      <div className="col-span-1 row-span-3 bg-base-100 rounded-2xl shadow-xl p-6">
        <div className="flex items-center gap-3 mb-6">
          <div className="p-3 bg-primary/10 rounded-xl">
            <Grid3X3 className="w-6 h-6 text-primary" />
          </div>
          <h3 className="text-xl font-bold">{t('categories')}</h3>
        </div>
        <div className="space-y-3">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="skeleton h-12 w-full rounded-lg" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <>
      {/* Мобильная версия - кнопки внизу экрана */}
      <div className="fixed bottom-0 left-0 right-0 z-40 lg:hidden bg-base-100 border-t border-base-200 p-4">
        <div className="grid grid-cols-2 gap-2">
          <button
            onClick={() => setShowCategoryModal(true)}
            className="btn btn-primary"
          >
            <Grid3X3 className="w-5 h-5" />
            {t('categories')}
            {selectedCategoryId && (
              <span className="badge badge-secondary badge-sm">1</span>
            )}
          </button>
          <button
            onClick={() => setShowFilterModal(true)}
            className="btn btn-outline"
          >
            <Filter className="w-5 h-5" />
            {t('filters.title')}
            {Object.keys(filters).filter((k) => filters[k]).length > 0 && (
              <span className="badge badge-primary badge-sm">
                {Object.keys(filters).filter((k) => filters[k]).length}
              </span>
            )}
          </button>
        </div>
      </div>

      {/* Десктопная версия - категории */}
      <div className="hidden lg:block col-span-1 row-span-3 bg-base-100 rounded-2xl shadow-xl p-2 sm:p-4 lg:p-6 overflow-hidden">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-primary/10 rounded-xl">
              <Grid3X3 className="w-6 h-6 text-primary" />
            </div>
            <h3 className="text-xl font-bold">{t('categories')}</h3>
          </div>
          {selectedCategoryId && (
            <button
              onClick={() => onCategorySelect(null)}
              className="btn btn-ghost btn-sm btn-circle"
              title={t('filters.clearFilter')}
            >
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          )}
        </div>

        <div className="overflow-y-auto max-h-[calc(100%-5rem)] space-y-1 pr-2 -mr-2">
          {categories.length === 0 ? (
            <p className="text-base-content/60 text-center py-8">
              {t('noCategoriesFound')}
            </p>
          ) : (
            categories.map((category) => renderCategory(category))
          )}
        </div>
      </div>

      {/* Десктопная версия - фильтры */}
      <DesktopFilters
        filters={filters}
        onFiltersChange={onFiltersChange!}
        selectedCategoryId={selectedCategoryId}
      />

      {/* Модальные окна */}
      <CategoryModal
        isOpen={showCategoryModal}
        onClose={() => setShowCategoryModal(false)}
        categories={categories}
        expandedCategories={expandedCategories}
        onToggleExpanded={toggleExpanded}
        selectedCategoryId={selectedCategoryId}
        onCategorySelect={onCategorySelect}
        title={t('categories')}
      />

      <FilterModal
        isOpen={showFilterModal}
        onClose={() => setShowFilterModal(false)}
        filters={filters}
        onFiltersChange={onFiltersChange!}
        selectedCategoryId={selectedCategoryId}
      />
    </>
  );
};
