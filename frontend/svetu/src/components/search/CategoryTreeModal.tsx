'use client';

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { MarketplaceService } from '@/services/c2c';
import { toast } from '@/utils/toast';
import {
  ChevronRight,
  ChevronDown,
  Folder,
  File,
  Search,
  X,
} from 'lucide-react';

interface Category {
  id: number;
  name: string;
  parent_id?: number | null;
  icon?: string;
  children?: Category[];
  count?: number;
  translations?: Record<string, string>;
}

interface CategoryTreeModalProps {
  isOpen: boolean;
  onClose: () => void;
  selectedCategoryId?: number;
  onCategorySelect: (categoryId: number) => void;
}

export const CategoryTreeModal: React.FC<CategoryTreeModalProps> = ({
  isOpen,
  onClose,
  selectedCategoryId,
  onCategorySelect,
}) => {
  const t = useTranslations('search');
  const locale = useLocale();
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<Set<number>>(new Set());
  const [searchQuery, setSearchQuery] = useState('');

  // Load categories from API
  useEffect(() => {
    if (!isOpen) return;

    const loadCategories = async () => {
      try {
        setLoading(true);
        const response = await MarketplaceService.getCategories(locale);

        // Build tree structure from flat array
        const categoryMap = new Map<number, Category>();
        const rootCategories: Category[] = [];

        // First pass: create all category objects
        response.data.forEach((cat) => {
          categoryMap.set(cat.id, {
            ...cat,
            children: [],
          });
        });

        // Second pass: build tree structure
        response.data.forEach((cat) => {
          const category = categoryMap.get(cat.id);
          if (!category) return;

          if (cat.parent_id) {
            const parent = categoryMap.get(cat.parent_id);
            if (parent) {
              parent.children = parent.children || [];
              parent.children.push(category);
            }
          } else {
            rootCategories.push(category);
          }
        });

        setCategories(rootCategories);
      } catch (error) {
        console.error('Failed to load categories:', error);
        toast.error(t('categoriesLoadError') || 'Ошибка загрузки категорий');
      } finally {
        setLoading(false);
      }
    };

    loadCategories();
  }, [locale, t, isOpen]);

  // Toggle category expansion
  const toggleExpand = (categoryId: number) => {
    setExpanded((prev) => {
      const next = new Set(prev);
      if (next.has(categoryId)) {
        next.delete(categoryId);
      } else {
        next.add(categoryId);
      }
      return next;
    });
  };

  // Handle category selection
  const handleSelect = (categoryId: number) => {
    onCategorySelect(categoryId);
    onClose();
  };

  // Filter categories based on search
  const filterCategories = useCallback(
    (cats: Category[], query: string): Category[] => {
      if (!query) return cats;

      const lowerQuery = query.toLowerCase();
      return cats.reduce((acc: Category[], cat) => {
        const matches = cat.name.toLowerCase().includes(lowerQuery);
        const filteredChildren = cat.children
          ? filterCategories(cat.children, query)
          : [];

        if (matches || filteredChildren.length > 0) {
          acc.push({
            ...cat,
            children: filteredChildren,
          });
        }

        return acc;
      }, []);
    },
    []
  );

  const filteredCategories = useMemo(
    () => filterCategories(categories, searchQuery),
    [categories, searchQuery, filterCategories]
  );

  // Render category tree recursively
  const renderTree = (items: Category[], level = 0): React.ReactElement[] => {
    return items.map((item) => {
      const hasChildren = item.children && item.children.length > 0;
      const isExpanded = expanded.has(item.id);
      const isSelected = selectedCategoryId === item.id;

      return (
        <div key={item.id}>
          <div
            className={`
              flex items-center gap-2 px-2 py-2 rounded cursor-pointer
              hover:bg-base-200 transition-colors
              ${isSelected ? 'bg-primary/10 text-primary' : ''}
            `}
            style={{ paddingLeft: `${level * 20 + 8}px` }}
            onClick={() => handleSelect(item.id)}
          >
            {hasChildren && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  toggleExpand(item.id);
                }}
                className="btn btn-xs btn-ghost p-0 min-h-0 h-5 w-5"
              >
                {isExpanded ? (
                  <ChevronDown size={16} />
                ) : (
                  <ChevronRight size={16} />
                )}
              </button>
            )}
            {!hasChildren && <div className="w-5" />}

            {hasChildren ? <Folder size={16} /> : <File size={16} />}

            <span className="flex-1">{item.name}</span>

            {item.count !== undefined && item.count > 0 && (
              <span className="badge badge-sm">{item.count}</span>
            )}
          </div>

          {hasChildren && isExpanded && (
            <div>{renderTree(item.children!, level + 1)}</div>
          )}
        </div>
      );
    });
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-2xl max-h-[80vh] p-0">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-base-300">
          <h3 className="text-lg font-semibold">{t('allCategories')}</h3>
          <button onClick={onClose} className="btn btn-ghost btn-sm btn-circle">
            <X size={20} />
          </button>
        </div>

        {/* Search */}
        <div className="p-4 border-b border-base-300">
          <div className="input input-bordered flex items-center gap-2">
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchCategories') || 'Поиск категорий...'}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="flex-1 outline-none bg-transparent"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                className="btn btn-ghost btn-xs p-0 min-h-0 h-5 w-5"
              >
                <X size={16} />
              </button>
            )}
          </div>
        </div>

        {/* Category tree */}
        <div
          className="overflow-auto flex-1 p-4"
          style={{ maxHeight: 'calc(80vh - 200px)' }}
        >
          {loading ? (
            <div className="flex justify-center items-center py-8">
              <div className="loading loading-spinner loading-md"></div>
            </div>
          ) : filteredCategories.length > 0 ? (
            <div className="space-y-1">{renderTree(filteredCategories)}</div>
          ) : (
            <div className="text-center py-8 text-base-content/50">
              {searchQuery
                ? t('noCategoriesFound') || 'Категории не найдены'
                : t('noCategories') || 'Нет категорий'}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="p-4 border-t border-base-300 flex justify-end gap-2">
          <button onClick={onClose} className="btn btn-ghost">
            {t('cancel') || 'Отмена'}
          </button>
        </div>
      </div>
      <div className="modal-backdrop" onClick={onClose}></div>
    </div>
  );
};
