'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { CategoryService, Category } from '@/services/category';
import { renderCategoryIcon } from '@/utils/iconMapper';
import { Package } from 'lucide-react';

interface CategoryFilterProps {
  selectedCategories: number[];
  onCategoryChange: (categories: number[]) => void;
}

interface CategoryTreeItem extends Category {
  children: CategoryTreeItem[];
  level: number;
}

export default function CategoryFilter({
  selectedCategories,
  onCategoryChange,
}: CategoryFilterProps) {
  const t = useTranslations('search');
  const locale = useLocale();
  const [categoryTree, setCategoryTree] = useState<CategoryTreeItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<Record<number, boolean>>({});

  const loadCategories = useCallback(async () => {
    try {
      const data = await CategoryService.getCategories();

      // Построение дерева категорий
      const tree = buildCategoryTree(data);
      setCategoryTree(tree);
    } catch (error) {
      console.error('Failed to load categories:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadCategories();
  }, [loadCategories]);

  const buildCategoryTree = (categories: Category[]): CategoryTreeItem[] => {
    const categoryMap = new Map<number, CategoryTreeItem>();
    const rootCategories: CategoryTreeItem[] = [];

    // Создаем все узлы
    categories.forEach((category) => {
      categoryMap.set(category.id, {
        ...category,
        children: [],
        level: 0,
      });
    });

    // Строим дерево
    categories.forEach((category) => {
      const categoryItem = categoryMap.get(category.id);
      if (!categoryItem) return;

      if (category.parent_id && category.parent_id !== 0) {
        const parent = categoryMap.get(category.parent_id);
        if (parent) {
          categoryItem.level = parent.level + 1;
          parent.children.push(categoryItem);
        } else {
          // Если родитель не найден, добавляем как корневую
          rootCategories.push(categoryItem);
        }
      } else {
        rootCategories.push(categoryItem);
      }
    });

    // Сортируем по sort_order или имени
    const sortByOrder = (items: CategoryTreeItem[]) => {
      items.sort((a, b) => {
        if (a.sort_order !== undefined && b.sort_order !== undefined) {
          return a.sort_order - b.sort_order;
        }
        return a.name.localeCompare(b.name);
      });
      items.forEach((item) => {
        if (item.children.length > 0) {
          sortByOrder(item.children);
        }
      });
    };

    sortByOrder(rootCategories);
    return rootCategories;
  };

  const toggleCategory = (categoryId: number) => {
    const newSelection = selectedCategories.includes(categoryId)
      ? selectedCategories.filter((id) => id !== categoryId)
      : [...selectedCategories, categoryId];
    onCategoryChange(newSelection);
  };

  const toggleExpand = (categoryId: number) => {
    setExpanded((prev) => ({
      ...prev,
      [categoryId]: !prev[categoryId],
    }));
  };

  const renderCategory = (category: CategoryTreeItem) => {
    const hasChildren = category.children && category.children.length > 0;
    const isExpanded = expanded[category.id];
    const isSelected = selectedCategories.includes(category.id);
    const paddingLeft = category.level * 16;

    // Подсчет выбранных дочерних категорий
    const countSelectedChildren = (cat: CategoryTreeItem): number => {
      let count = 0;
      if (selectedCategories.includes(cat.id)) count++;
      cat.children.forEach((child) => {
        count += countSelectedChildren(child);
      });
      return count;
    };

    const selectedChildrenCount = hasChildren
      ? countSelectedChildren(category) - (isSelected ? 1 : 0)
      : 0;

    return (
      <div key={category.id}>
        <div
          className={`flex items-center gap-2 py-2 rounded-lg px-2 transition-all duration-200 ${
            isSelected
              ? 'bg-primary/10 border-l-2 border-primary'
              : 'hover:bg-base-200 border-l-2 border-transparent'
          }`}
          style={{ marginLeft: `${paddingLeft}px` }}
        >
          {hasChildren ? (
            <button
              onClick={(e) => {
                e.preventDefault();
                toggleExpand(category.id);
              }}
              className="btn btn-ghost btn-xs p-0 min-h-0 h-5 w-5 hover:bg-base-300"
              aria-label={`${isExpanded ? 'Collapse' : 'Expand'} ${category.translations?.[locale] || category.name}`}
            >
              <svg
                className={`w-4 h-4 transition-transform duration-200 ${isExpanded ? 'rotate-90' : ''}`}
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
          ) : (
            <div className="w-5 flex items-center justify-center">
              <div className="w-1.5 h-1.5 rounded-full bg-base-300" />
            </div>
          )}

          <label className="cursor-pointer flex items-center gap-2 flex-1">
            <input
              type="checkbox"
              className="checkbox checkbox-sm checkbox-primary"
              checked={isSelected}
              onChange={() => toggleCategory(category.id)}
              aria-label={category.translations?.[locale] || category.name}
            />
            {category.icon ? (
              renderCategoryIcon(category.icon, 'w-4 h-4 text-primary')
            ) : (
              <Package className="w-4 h-4 text-base-content/40" />
            )}
            <span className="text-sm font-medium flex-1">
              {category.translations?.[locale] || category.name}
            </span>

            {selectedChildrenCount > 0 && !isExpanded && (
              <span className="badge badge-primary badge-xs">
                {selectedChildrenCount}
              </span>
            )}

            {category.listing_count ||
            category.count ||
            category.product_count ? (
              <span className="text-xs text-base-content/60">
                (
                {category.listing_count ||
                  category.count ||
                  category.product_count}
                )
              </span>
            ) : null}
          </label>
        </div>

        {hasChildren && isExpanded && (
          <div>{category.children.map((child) => renderCategory(child))}</div>
        )}
      </div>
    );
  };

  if (loading) {
    return (
      <div className="space-y-2">
        {[1, 2, 3].map((i) => (
          <div
            key={i}
            className="skeleton h-6 w-full"
            data-testid="skeleton"
          ></div>
        ))}
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-2">
        <label className="label p-0">
          <span className="label-text font-medium flex items-center gap-2">
            <svg
              className="w-4 h-4 text-primary"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
              />
            </svg>
            {t('categories')}
          </span>
        </label>

        <div className="flex items-center gap-1">
          <button
            className="btn btn-ghost btn-xs"
            onClick={() => {
              const allCategoryIds: Record<number, boolean> = {};
              const collectIds = (cats: CategoryTreeItem[]) => {
                cats.forEach((cat) => {
                  if (cat.children.length > 0) {
                    allCategoryIds[cat.id] = true;
                    collectIds(cat.children);
                  }
                });
              };
              collectIds(categoryTree);
              setExpanded(allCategoryIds);
            }}
            title={t('expandAll') || 'Expand all'}
          >
            <svg
              className="w-3 h-3"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </button>

          <button
            className="btn btn-ghost btn-xs"
            onClick={() => setExpanded({})}
            title={t('collapseAll') || 'Collapse all'}
          >
            <svg
              className="w-3 h-3"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 15l7-7 7 7"
              />
            </svg>
          </button>

          {selectedCategories.length > 0 && (
            <>
              <div className="w-px h-4 bg-base-300 mx-1" />
              <button
                className="btn btn-ghost btn-xs text-primary"
                onClick={() => onCategoryChange([])}
              >
                {t('clear')} ({selectedCategories.length})
              </button>
            </>
          )}
        </div>
      </div>

      <div
        className="max-h-96 overflow-y-auto space-y-0.5 border border-base-300 rounded-lg p-2 
        scrollbar-thin scrollbar-thumb-base-300 scrollbar-track-transparent hover:scrollbar-thumb-base-400"
      >
        {categoryTree.length > 0 ? (
          categoryTree.map((category) => renderCategory(category))
        ) : (
          <div className="text-sm text-base-content/60 text-center py-4">
            {t('noCategories') || 'No categories available'}
          </div>
        )}
      </div>
    </div>
  );
}
