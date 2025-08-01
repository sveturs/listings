'use client';

import React, { useState, useEffect } from 'react';
import { ChevronRight, ChevronDown, Search, Package, X } from 'lucide-react';
import { renderCategoryIcon } from '@/utils/iconMapper';
import type { components } from '@/types/generated/api';

type MarketplaceCategory =
  components['schemas']['backend_internal_domain_models.MarketplaceCategory'];

interface CategorySelectorProps {
  categories: MarketplaceCategory[];
  selectedCategory: MarketplaceCategory | null;
  onCategorySelect: (category: MarketplaceCategory) => void;
  locale: string;
  compact?: boolean;
}

interface CategoryTreeItem extends MarketplaceCategory {
  children: CategoryTreeItem[];
  level: number;
}

export default function CategorySelector({
  categories,
  selectedCategory,
  onCategorySelect,
  locale: _locale,
  compact = false,
}: CategorySelectorProps) {
  const [searchTerm, setSearchTerm] = useState('');
  const [expandedCategories, setExpandedCategories] = useState<Set<number>>(
    new Set()
  );
  const [showModal, setShowModal] = useState(false);
  const [categoryTree, setCategoryTree] = useState<CategoryTreeItem[]>([]);

  // Строим дерево категорий
  useEffect(() => {
    const buildCategoryTree = (
      categories: MarketplaceCategory[]
    ): CategoryTreeItem[] => {
      const categoryMap = new Map<number, CategoryTreeItem>();
      const rootCategories: CategoryTreeItem[] = [];

      // Создаем все узлы
      categories.forEach((category) => {
        categoryMap.set(category.id || 0, {
          ...category,
          children: [],
          level: 0,
        });
      });

      // Строим дерево
      categories.forEach((category) => {
        const categoryItem = categoryMap.get(category.id || 0);
        if (!categoryItem) return;

        if (category.parent_id && category.parent_id !== 0) {
          const parent = categoryMap.get(category.parent_id);
          if (parent) {
            categoryItem.level = parent.level + 1;
            parent.children.push(categoryItem);
          }
        } else {
          rootCategories.push(categoryItem);
        }
      });

      // Сортируем по sort_order
      const sortByOrder = (items: CategoryTreeItem[]) => {
        items.sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0));
        items.forEach((item) => {
          if (item.children.length > 0) {
            sortByOrder(item.children);
          }
        });
      };

      sortByOrder(rootCategories);
      return rootCategories;
    };

    setCategoryTree(buildCategoryTree(categories));
  }, [categories]);

  // Автоматически разворачиваем путь к выбранной категории
  useEffect(() => {
    if (selectedCategory) {
      const findPath = (category: MarketplaceCategory): number[] => {
        const path: number[] = [];
        let current = category;

        while (current.parent_id && current.parent_id !== 0) {
          path.unshift(current.parent_id);
          const parent = categories.find((c) => c.id === current.parent_id);
          if (!parent) break;
          current = parent;
        }

        return path;
      };

      const pathToSelected = findPath(selectedCategory);
      setExpandedCategories(new Set(pathToSelected));
    }
  }, [selectedCategory, categories]);

  const toggleExpanded = (categoryId: number) => {
    const newExpanded = new Set(expandedCategories);
    if (newExpanded.has(categoryId)) {
      newExpanded.delete(categoryId);
    } else {
      newExpanded.add(categoryId);
    }
    setExpandedCategories(newExpanded);
  };

  const getCategoryName = (category: MarketplaceCategory) => {
    return category.translations?.name || category.name || '';
  };

  const renderCategoryItem = (category: CategoryTreeItem) => {
    const hasChildren = category.children.length > 0;
    const isExpanded = expandedCategories.has(category.id || 0);
    const isSelected = selectedCategory?.id === category.id;
    const categoryName = getCategoryName(category);

    // Фильтрация по поиску
    if (
      searchTerm &&
      !categoryName.toLowerCase().includes(searchTerm.toLowerCase())
    ) {
      // Проверяем, есть ли совпадения в дочерних категориях
      const hasMatchingChildren = category.children.some((child) =>
        getCategoryName(child).toLowerCase().includes(searchTerm.toLowerCase())
      );
      if (!hasMatchingChildren) return null;
    }

    return (
      <div key={category.id} className={`ml-${category.level * 4}`}>
        <div
          className={`flex items-center justify-between p-2 rounded-lg cursor-pointer transition-colors ${
            isSelected ? 'bg-primary text-primary-content' : 'hover:bg-base-200'
          }`}
        >
          <div
            className="flex items-center flex-1"
            onClick={() => onCategorySelect(category)}
          >
            {hasChildren && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  toggleExpanded(category.id || 0);
                }}
                className="mr-2 p-1 rounded"
              >
                {isExpanded ? (
                  <ChevronDown className="w-4 h-4" />
                ) : (
                  <ChevronRight className="w-4 h-4" />
                )}
              </button>
            )}

            {!hasChildren && <div className="w-6 mr-2" />}

            {renderCategoryIcon(category.icon, 'w-5 h-5 mr-2')}

            <span className="text-sm font-medium">{categoryName}</span>

            {category.listing_count !== undefined &&
              category.listing_count > 0 && (
                <span className="ml-2 text-xs opacity-60">
                  ({category.listing_count})
                </span>
              )}
          </div>
        </div>

        {hasChildren && isExpanded && (
          <div className="mt-1">
            {category.children.map((child) => renderCategoryItem(child))}
          </div>
        )}
      </div>
    );
  };

  const renderCompactView = () => (
    <div className="form-control">
      <label className="label">
        <span className="label-text font-semibold">Категория</span>
        <button
          onClick={() => setShowModal(true)}
          className="label-text-alt link link-primary"
        >
          Выбрать другую
        </button>
      </label>

      <button
        onClick={() => setShowModal(true)}
        className={`btn btn-outline justify-start ${
          selectedCategory ? 'btn-primary' : ''
        }`}
      >
        {selectedCategory ? (
          <>
            {renderCategoryIcon(selectedCategory.icon, 'w-4 h-4')}
            {getCategoryName(selectedCategory)}
          </>
        ) : (
          <>
            <Package className="w-4 h-4" />
            Выберите категорию
          </>
        )}
        <ChevronRight className="w-4 h-4 ml-auto" />
      </button>
    </div>
  );

  const renderFullView = () => (
    <div className="card bg-base-200">
      <div className="card-body">
        <h3 className="card-title text-base mb-4">
          <Package className="w-5 h-5" />
          Выберите категорию
        </h3>

        {/* Поиск */}
        <div className="form-control mb-4">
          <div className="input-group">
            <span>
              <Search className="w-4 h-4" />
            </span>
            <input
              type="text"
              placeholder="Поиск по категориям..."
              className="input input-bordered input-sm flex-1"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
            {searchTerm && (
              <button
                onClick={() => setSearchTerm('')}
                className="btn btn-ghost btn-sm"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>
        </div>

        {/* Выбранная категория */}
        {selectedCategory && (
          <div className="alert alert-success mb-4">
            <div className="flex items-center">
              {renderCategoryIcon(selectedCategory.icon, 'w-5 h-5 mr-2')}
              <div>
                <div className="font-bold">Выбрано:</div>
                <div className="text-sm">
                  {getCategoryName(selectedCategory)}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Дерево категорий */}
        <div className="max-h-96 overflow-y-auto space-y-1">
          {categoryTree.map((category) => renderCategoryItem(category))}
        </div>
      </div>
    </div>
  );

  return (
    <>
      {compact ? renderCompactView() : renderFullView()}

      {/* Модальное окно для компактного режима */}
      {showModal && (
        <div className="modal modal-open">
          <div className="modal-box max-w-2xl">
            <h3 className="font-bold text-lg mb-4">Выберите категорию</h3>

            {/* Поиск */}
            <div className="form-control mb-4">
              <div className="input-group">
                <span>
                  <Search className="w-4 h-4" />
                </span>
                <input
                  type="text"
                  placeholder="Поиск по категориям..."
                  className="input input-bordered input-sm flex-1"
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                />
                {searchTerm && (
                  <button
                    onClick={() => setSearchTerm('')}
                    className="btn btn-ghost btn-sm"
                  >
                    <X className="w-4 h-4" />
                  </button>
                )}
              </div>
            </div>

            {/* Дерево категорий */}
            <div className="max-h-96 overflow-y-auto space-y-1 mb-4">
              {categoryTree.map((category) => renderCategoryItem(category))}
            </div>

            <div className="modal-action">
              <button
                className="btn btn-ghost"
                onClick={() => setShowModal(false)}
              >
                Отмена
              </button>
              <button
                className="btn btn-primary"
                onClick={() => {
                  setShowModal(false);
                  setSearchTerm('');
                }}
                disabled={!selectedCategory}
              >
                Выбрать
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
