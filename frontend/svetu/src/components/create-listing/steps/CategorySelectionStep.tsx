'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';
import { MarketplaceService } from '@/services/c2c';
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

interface Category {
  id: number;
  name: string;
  parent_id?: number | null;
  icon?: string;
  children?: Category[];
  count?: number;
  translations?: Record<string, string>;
}

interface CategorySelectionStepProps {
  onNext: () => void;
}

export default function CategorySelectionStep({
  onNext,
}: CategorySelectionStepProps) {
  const t = useTranslations('create_listing');
  const tMarketplace = useTranslations('marketplace');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const { state, setCategory } = useCreateListing();
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState<Set<number>>(new Set());
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(
    state.category?.id || null
  );

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
        toast.error(tMarketplace('categoriesLoadError'));
      } finally {
        setLoading(false);
      }
    };

    loadCategories();
  }, [locale, t, tMarketplace]);

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
  const handleSelect = (
    categoryId: number,
    categoryName: string,
    hasChildren: boolean
  ) => {
    if (hasChildren) {
      // Only expand/collapse if it has children
      toggleExpand(categoryId);
      return;
    }

    setSelectedCategoryId(categoryId);
    setCategory({
      id: categoryId,
      name: categoryName,
      slug: `category-${categoryId}`,
    });
  };

  const handleNext = () => {
    if (selectedCategoryId) {
      onNext();
    }
  };

  // Filter categories based on search
  const filterCategories = (cats: Category[], query: string): Category[] => {
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
  };

  const filteredCategories = searchQuery
    ? filterCategories(categories, searchQuery)
    : categories;

  // Render category tree recursively
  const renderTree = (items: Category[], level = 0): React.ReactElement[] => {
    return items.map((item) => {
      const hasChildren = item.children && item.children.length > 0;
      const isExpanded = expanded.has(item.id);
      const isSelected = item.id === selectedCategoryId;
      const isSelectable = !hasChildren;

      return (
        <div key={item.id}>
          <div
            className={`
              flex items-center gap-2 px-3 py-2 rounded-lg transition-all
              ${isSelectable ? 'cursor-pointer hover:bg-base-200' : ''}
              ${isSelected ? 'bg-primary/10 text-primary font-medium' : ''}
              ${!isSelectable ? 'opacity-75' : ''}
            `}
            style={{ paddingLeft: `${level * 24 + 12}px` }}
            onClick={() =>
              handleSelect(item.id, item.name, hasChildren || false)
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

            {hasChildren ? <Folder size={18} /> : <File size={18} />}

            <span className="flex-1">{item.name}</span>

            {item.count !== undefined && item.count > 0 && (
              <span className="badge badge-sm badge-ghost">{item.count}</span>
            )}

            {isSelected && (
              <svg
                className="w-5 h-5 text-primary"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                  clipRule="evenodd"
                />
              </svg>
            )}
          </div>

          {hasChildren && isExpanded && (
            <div>{renderTree(item.children!, level + 1)}</div>
          )}
        </div>
      );
    });
  };

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto">
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body">
            <div className="flex items-center justify-center py-16">
              <div className="loading loading-spinner loading-lg"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            🏪 {t('title')}
          </h2>
          <p className="text-base-content/70 mb-6">{t('description')}</p>

          {/* Search input */}
          <div className="mb-4">
            <div className="input input-bordered flex items-center gap-2">
              <Search size={20} />
              <input
                type="text"
                placeholder={tMarketplace('searchCategories')}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="flex-1 outline-none bg-transparent"
              />
            </div>
          </div>

          {/* Category tree */}
          <div className="border border-base-300 rounded-lg p-4 max-h-96 overflow-y-auto mb-6">
            {filteredCategories.length > 0 ? (
              renderTree(filteredCategories)
            ) : (
              <div className="text-center py-8 text-base-content/50">
                {tMarketplace('noCategoriesFound')}
              </div>
            )}
          </div>

          {/* Selected category display */}
          {selectedCategoryId && (
            <div className="alert alert-success mb-6">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="stroke-current shrink-0 h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div>
                <div className="font-medium">{t('selected')}:</div>
                <div className="text-sm">{state.category?.name}</div>
              </div>
            </div>
          )}

          {/* Подсказка в региональном стиле */}
          <div className="alert alert-info">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="stroke-current shrink-0 w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
            <div className="text-sm">
              <p className="font-medium">💡 {t('tip_title')}</p>
              <p className="text-xs mt-1">{t('tip_description')}</p>
            </div>
          </div>

          {/* Кнопка продолжения */}
          <div className="card-actions justify-end mt-6">
            <button
              className={`btn btn-primary ${!selectedCategoryId ? 'btn-disabled' : ''}`}
              onClick={handleNext}
              disabled={!selectedCategoryId}
            >
              {tCommon('continue')} →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
