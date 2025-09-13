'use client';

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { MarketplaceService } from '@/services/marketplace';
import { toast } from '@/utils/toast';

// Icons as components
const ChevronRight = ({ size = 16 }: { size?: number }) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <polyline points="9 18 15 12 9 6"></polyline>
  </svg>
);

const ChevronDown = ({ size = 16 }: { size?: number }) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <polyline points="6 9 12 15 18 9"></polyline>
  </svg>
);

const Folder = ({ size = 16 }: { size?: number }) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
  </svg>
);

const File = ({ size = 16 }: { size?: number }) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
    <polyline points="13 2 13 9 20 9"></polyline>
  </svg>
);

const Search = ({ size = 16 }: { size?: number }) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <circle cx="11" cy="11" r="8"></circle>
    <path d="m21 21-4.35-4.35"></path>
  </svg>
);

const X = ({ size = 16 }: { size?: number }) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <line x1="18" y1="6" x2="6" y2="18"></line>
    <line x1="6" y1="6" x2="18" y2="18"></line>
  </svg>
);

interface Category {
  id: number;
  name: string;
  parent_id?: number | null;
  icon?: string;
  children?: Category[];
  count?: number;
  translations?: Record<string, string>;
}

interface CategoryTreeSelectorProps {
  value?: number | number[];
  onChange: (value: number | number[]) => void;
  multiple?: boolean;
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  showPath?: boolean;
  allowParentSelection?: boolean;
  maxDepth?: number;
  className?: string;
}

export function CategoryTreeSelector({
  value,
  onChange,
  multiple = false,
  placeholder,
  disabled = false,
  required = false,
  showPath = false,
  allowParentSelection = true,
  maxDepth,
  className = '',
}: CategoryTreeSelectorProps) {
  const t = useTranslations('marketplace');
  const locale = useLocale();
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<Set<number>>(new Set());
  const [searchQuery, setSearchQuery] = useState('');
  const [showDropdown, setShowDropdown] = useState(false);
  const [selectedCategories, setSelectedCategories] = useState<Set<number>>(
    new Set()
  );

  // Initialize selected categories from value prop
  useEffect(() => {
    if (value !== undefined) {
      setSelectedCategories(prevCategories => {
        const newSet = Array.isArray(value)
          ? new Set(value)
          : value
            ? new Set([value])
            : new Set();

        // Проверяем, изменилось ли значение
        const currentArray = Array.from(prevCategories);
        const newArray = Array.from(newSet);

        if (JSON.stringify(currentArray.sort()) !== JSON.stringify(newArray.sort())) {
          return newSet;
        }
        return prevCategories;
      });
    }
  }, [value]);

  // Load categories from API
  useEffect(() => {
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
        toast.error(t('categoriesLoadError'));
      } finally {
        setLoading(false);
      }
    };

    loadCategories();
  }, [locale, t]);

  // Build category path for display
  const getCategoryPath = useCallback(
    (categoryId: number): string => {
      const path: string[] = [];

      const findPath = (cats: Category[], targetId: number): boolean => {
        for (const cat of cats) {
          if (cat.id === targetId) {
            path.unshift(cat.name);
            return true;
          }
          if (cat.children && findPath(cat.children, targetId)) {
            path.unshift(cat.name);
            return true;
          }
        }
        return false;
      };

      findPath(categories, categoryId);
      return path.join(' > ');
    },
    [categories]
  );

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
  const handleSelect = (categoryId: number, hasChildren: boolean) => {
    if (!allowParentSelection && hasChildren) {
      // Only expand/collapse if parent selection not allowed
      toggleExpand(categoryId);
      return;
    }

    if (multiple) {
      const newSelection = new Set(selectedCategories);
      if (newSelection.has(categoryId)) {
        newSelection.delete(categoryId);
      } else {
        newSelection.add(categoryId);
      }
      setSelectedCategories(newSelection);
      onChange(Array.from(newSelection));
    } else {
      setSelectedCategories(new Set([categoryId]));
      onChange(categoryId);
      setShowDropdown(false);
    }
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
    if (maxDepth && level >= maxDepth) return [];

    return items.map((item) => {
      const hasChildren = item.children && item.children.length > 0;
      const isExpanded = expanded.has(item.id);
      const isSelected = selectedCategories.has(item.id);
      const isSelectable = allowParentSelection || !hasChildren;

      return (
        <div key={item.id}>
          <div
            className={`
              flex items-center gap-2 px-2 py-1 rounded
              ${isSelectable ? 'cursor-pointer hover:bg-base-200' : ''}
              ${isSelected ? 'bg-primary/10 text-primary' : ''}
              ${!isSelectable ? 'opacity-60' : ''}
            `}
            style={{ paddingLeft: `${level * 20 + 8}px` }}
            onClick={() =>
              isSelectable && handleSelect(item.id, hasChildren || false)
            }
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

            {multiple && (
              <input
                type="checkbox"
                checked={isSelected}
                onChange={() => {}}
                className="checkbox checkbox-sm"
                disabled={!isSelectable}
              />
            )}
          </div>

          {hasChildren && isExpanded && (
            <div>{renderTree(item.children!, level + 1)}</div>
          )}
        </div>
      );
    });
  };

  // Get display text for selected categories
  const getDisplayText = () => {
    if (selectedCategories.size === 0) {
      return placeholder || t('selectCategory');
    }

    if (showPath && selectedCategories.size === 1) {
      const categoryId = Array.from(selectedCategories)[0];
      return getCategoryPath(categoryId);
    }

    if (selectedCategories.size === 1) {
      const categoryId = Array.from(selectedCategories)[0];
      const findCategory = (cats: Category[]): Category | null => {
        for (const cat of cats) {
          if (cat.id === categoryId) return cat;
          if (cat.children) {
            const found = findCategory(cat.children);
            if (found) return found;
          }
        }
        return null;
      };
      const category = findCategory(categories);
      return category?.name || '';
    }

    return t('categoriesSelected', { count: selectedCategories.size });
  };

  if (loading) {
    return (
      <div className={`form-control ${className}`}>
        <div className="loading loading-spinner loading-sm"></div>
      </div>
    );
  }

  return (
    <div className={`form-control relative ${className}`}>
      {/* Main selector button */}
      <div
        className={`
          input input-bordered flex items-center justify-between cursor-pointer
          ${disabled ? 'input-disabled' : ''}
          ${required && selectedCategories.size === 0 ? 'input-error' : ''}
        `}
        onClick={() => !disabled && setShowDropdown(!showDropdown)}
      >
        <span
          className={
            selectedCategories.size === 0 ? 'text-base-content/50' : ''
          }
        >
          {getDisplayText()}
        </span>
        <div className="flex items-center gap-2">
          {selectedCategories.size > 0 && !disabled && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                setSelectedCategories(new Set());
                onChange(multiple ? [] : 0);
              }}
              className="btn btn-ghost btn-xs p-0 min-h-0 h-5 w-5"
            >
              <X size={16} />
            </button>
          )}
          <div className={showDropdown ? 'rotate-180' : ''}>
            <ChevronDown size={16} />
          </div>
        </div>
      </div>

      {/* Dropdown */}
      {showDropdown && !disabled && (
        <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-96 overflow-hidden">
          {/* Search input */}
          <div className="p-2 border-b border-base-300">
            <div className="input input-sm input-bordered flex items-center gap-2">
              <Search size={16} />
              <input
                type="text"
                placeholder={t('searchCategories')}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="flex-1 outline-none bg-transparent"
                onClick={(e) => e.stopPropagation()}
              />
            </div>
          </div>

          {/* Category tree */}
          <div className="overflow-auto max-h-80 p-2">
            {filteredCategories.length > 0 ? (
              renderTree(filteredCategories)
            ) : (
              <div className="text-center py-4 text-base-content/50">
                {t('noCategoriesFound')}
              </div>
            )}
          </div>

          {/* Action buttons for multiple selection */}
          {multiple && (
            <div className="p-2 border-t border-base-300 flex justify-end gap-2">
              <button
                onClick={() => setShowDropdown(false)}
                className="btn btn-sm btn-ghost"
              >
                {t('cancel')}
              </button>
              <button
                onClick={() => setShowDropdown(false)}
                className="btn btn-sm btn-primary"
              >
                {t('apply')}
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
