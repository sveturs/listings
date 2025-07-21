'use client';

import React, { useState, useEffect } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import { MarketplaceService } from '@/services/marketplace';
import type { components } from '@/types/generated/api';

type MarketplaceCategory =
  components['schemas']['backend_internal_domain_models.MarketplaceCategory'];

interface CategorySidebarProps {
  onCategorySelect: (categoryId: number | null) => void;
  selectedCategoryId?: number | null;
  className?: string;
}

interface CategoryTreeNode extends MarketplaceCategory {
  children: CategoryTreeNode[];
}

export default function CategorySidebar({
  onCategorySelect,
  selectedCategoryId,
  className = '',
}: CategorySidebarProps) {
  const locale = useLocale();
  const t = useTranslations('marketplace');
  const [categories, setCategories] = useState<CategoryTreeNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [expandedCategories, setExpandedCategories] = useState<Set<number>>(
    new Set()
  );

  // Построение дерева категорий
  const buildCategoryTree = (
    categories: MarketplaceCategory[]
  ): CategoryTreeNode[] => {
    const categoryMap = new Map<number, CategoryTreeNode>();
    const rootCategories: CategoryTreeNode[] = [];

    // Создаем узлы для всех категорий
    categories.forEach((category) => {
      categoryMap.set(category.id!, {
        ...category,
        children: [],
      });
    });

    // Строим дерево
    categories.forEach((category) => {
      const node = categoryMap.get(category.id!)!;

      if (category.parent_id && categoryMap.has(category.parent_id)) {
        const parent = categoryMap.get(category.parent_id)!;
        parent.children.push(node);
      } else {
        rootCategories.push(node);
      }
    });

    // Сортируем по sort_order
    const sortCategories = (cats: CategoryTreeNode[]) => {
      cats.sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0));
      cats.forEach((cat) => sortCategories(cat.children));
    };

    sortCategories(rootCategories);
    return rootCategories;
  };

  // Загрузка категорий
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await MarketplaceService.getCategories(locale);

        if (response.success && response.data) {
          // TODO: Фильтруем только активные категории когда будет добавлено поле is_active
          // const activeCategories = response.data.filter((cat) => cat.is_active);

          // Приводим типы к MarketplaceCategory
          const mappedCategories: MarketplaceCategory[] = response.data.map(
            (cat) => ({
              ...cat,
              parent_id: cat.parent_id === null ? undefined : cat.parent_id,
            })
          );

          const tree = buildCategoryTree(mappedCategories);
          setCategories(tree);
        } else {
          setError(t('categoriesLoadError'));
        }
      } catch (err) {
        console.error('Error fetching categories:', err);
        setError('Ошибка загрузки категорий');
      } finally {
        setLoading(false);
      }
    };

    fetchCategories();
  }, [locale, t]);

  // Переключение раскрытия категории
  const toggleExpanded = (categoryId: number) => {
    const newExpanded = new Set(expandedCategories);
    if (newExpanded.has(categoryId)) {
      newExpanded.delete(categoryId);
    } else {
      newExpanded.add(categoryId);
    }
    setExpandedCategories(newExpanded);
  };

  // Рендер категории
  const renderCategory = (category: CategoryTreeNode, level: number = 0) => {
    const isExpanded = expandedCategories.has(category.id!);
    const hasChildren = category.children.length > 0;
    const isSelected = selectedCategoryId === category.id;

    return (
      <div key={category.id} className="w-full">
        <div
          className={`flex items-center gap-2 p-2 rounded-lg cursor-pointer transition-colors hover:bg-base-200 ${
            isSelected ? 'bg-primary/10 text-primary' : 'text-base-content'
          }`}
          style={{ paddingLeft: `${level * 12 + 8}px` }}
          onClick={() => onCategorySelect(category.id!)}
        >
          {hasChildren && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleExpanded(category.id!);
              }}
              className="btn btn-ghost btn-xs p-0 min-h-0 h-4 w-4"
            >
              <svg
                className={`w-3 h-3 transition-transform ${isExpanded ? 'rotate-90' : ''}`}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 5l7 7-7 7"
                />
              </svg>
            </button>
          )}

          {!hasChildren && <div className="w-4" />}

          {category.icon && (
            <span className="text-base-content/70 text-sm">
              {category.icon}
            </span>
          )}

          <span className="text-sm font-medium truncate">{category.name}</span>

          {category.listing_count !== undefined &&
            category.listing_count > 0 && (
              <span className="badge badge-sm badge-ghost ml-auto">
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
      <div
        className={`bg-base-100 p-4 rounded-lg shadow-sm border ${className}`}
      >
        <div className="flex items-center gap-2 mb-4">
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 11H5m14-7l2 2m0 0l2 2m-2-2l-2 2m-2-2l-2 2M5 20l14 0"
            />
          </svg>
          <h3 className="font-semibold text-base-content">{t('categories')}</h3>
        </div>

        <div className="space-y-2">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="animate-pulse">
              <div className="h-8 bg-base-200 rounded-lg"></div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div
        className={`bg-base-100 p-4 rounded-lg shadow-sm border ${className}`}
      >
        <div className="alert alert-error">
          <svg
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span className="text-sm">{error}</span>
        </div>
      </div>
    );
  }

  return (
    <div className={`bg-base-100 p-4 rounded-lg shadow-sm border ${className}`}>
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 11H5m14-7l2 2m0 0l2 2m-2-2l-2 2m-2-2l-2 2M5 20l14 0"
            />
          </svg>
          <h3 className="font-semibold text-base-content">{t('categories')}</h3>
        </div>

        {selectedCategoryId && (
          <button
            onClick={() => onCategorySelect(null)}
            className="btn btn-ghost btn-xs"
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

      <div className="space-y-1 max-h-96 overflow-y-auto">
        {categories.length === 0 ? (
          <p className="text-sm text-base-content/60 text-center py-4">
            {t('noCategoriesFound')}
          </p>
        ) : (
          categories.map((category) => renderCategory(category))
        )}
      </div>
    </div>
  );
}
