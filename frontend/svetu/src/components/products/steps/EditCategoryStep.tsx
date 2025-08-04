'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useEditProduct } from '@/contexts/EditProductContext';
import { apiClient } from '@/services/api-client';
import { toast } from '@/utils/toast';
import type { components } from '@/types/generated/api';

type MarketplaceCategory =
  components['schemas']['backend_internal_domain_models.MarketplaceCategory'];

interface EditCategoryStepProps {
  onNext: () => void;
}

export default function EditCategoryStep({ onNext }: EditCategoryStepProps) {
  const t = useTranslations('storefronts.products');
  const tCommon = useTranslations('common');
  const { state, setCategory, setError, clearError } = useEditProduct();
  const [allCategories, setAllCategories] = useState<MarketplaceCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] =
    useState<MarketplaceCategory | null>(state.category || null);

  // Навигация по иерархии
  const [currentParentId, setCurrentParentId] = useState<number | null>(null);
  const [breadcrumbs, setBreadcrumbs] = useState<MarketplaceCategory[]>([]);

  const buildBreadcrumbs = useCallback(
    (category: MarketplaceCategory, categories: MarketplaceCategory[]) => {
      const breadcrumbPath: MarketplaceCategory[] = [];
      let current = category;

      while (current.parent_id) {
        const parent = categories.find((c) => c.id === current.parent_id);
        if (parent) {
          breadcrumbPath.unshift(parent);
          current = parent;
        } else {
          break;
        }
      }

      setBreadcrumbs(breadcrumbPath);
      setCurrentParentId(category.parent_id || null);
    },
    []
  );

  const loadCategories = useCallback(async () => {
    try {
      setLoading(true);
      const response = await apiClient.get('/api/v1/marketplace/category-tree');

      if (response.data) {
        const responseData = response.data.data || response.data;

        if (Array.isArray(responseData)) {
          // Создаем плоский список всех категорий, сохраняя иерархию
          const flattenCategories = (cats: any[]): MarketplaceCategory[] => {
            const result: MarketplaceCategory[] = [];
            const seenIds = new Set<number>();

            const flatten = (categories: any[]) => {
              categories.forEach((cat) => {
                // Пропускаем дубликаты
                if (seenIds.has(cat.id)) {
                  console.warn(
                    `EditCategoryStep - Duplicate category ID found: ${cat.id}`
                  );
                  return;
                }

                seenIds.add(cat.id);
                result.push({
                  id: cat.id,
                  name: cat.name,
                  slug: cat.slug,
                  parent_id: cat.parent_id,
                  icon: cat.icon,
                  listing_count: cat.listing_count,
                });

                if (cat.children && cat.children.length > 0) {
                  flatten(cat.children);
                }
              });
            };

            flatten(cats);
            return result;
          };

          const flattened = flattenCategories(responseData);
          setAllCategories(flattened);

          // Если у нас уже есть выбранная категория, восстановим навигацию
          if (state.category) {
            const category = flattened.find((c) => c.id === state.category?.id);
            if (category) {
              setSelectedCategory(category);
              buildBreadcrumbs(category, flattened);
            }
          }
        }
      }
    } catch (error: any) {
      console.error('Failed to load categories:', error);
      toast.error(t('errorLoadingCategories'));
      setError('category', t('errorLoadingCategories'));
    } finally {
      setLoading(false);
    }
  }, [state.category, t, setError, setAllCategories, buildBreadcrumbs]);

  useEffect(() => {
    loadCategories();
  }, [loadCategories]);

  // Получить текущие категории для отображения
  const getCurrentCategories = () => {
    const filtered = allCategories.filter((cat) => {
      if (currentParentId === null) {
        return !cat.parent_id; // Корневые категории
      } else {
        return cat.parent_id === currentParentId; // Подкатегории
      }
    });
    return filtered;
  };

  // Получить дочерние категории
  const getChildCategories = (parentId: number) => {
    return allCategories.filter((cat) => cat.parent_id === parentId);
  };

  // Проверить, есть ли дочерние категории
  const hasChildren = (categoryId: number) => {
    return getChildCategories(categoryId).length > 0;
  };

  // Навигация вглубь
  const navigateToCategory = (category: MarketplaceCategory) => {
    if (category.id && hasChildren(category.id)) {
      setCurrentParentId(category.id);
      setBreadcrumbs([...breadcrumbs, category]);
    } else {
      // Это конечная категория - выбираем ее
      setSelectedCategory(category);
      setCategory(category);
      clearError('category');
    }
  };

  // Навигация назад
  const navigateBack = () => {
    if (breadcrumbs.length > 0) {
      const newBreadcrumbs = [...breadcrumbs];
      newBreadcrumbs.pop();
      setBreadcrumbs(newBreadcrumbs);

      if (newBreadcrumbs.length === 0) {
        setCurrentParentId(null);
      } else {
        setCurrentParentId(
          newBreadcrumbs[newBreadcrumbs.length - 1].id || null
        );
      }
    }
  };

  // Навигация к корню
  const navigateToRoot = () => {
    setCurrentParentId(null);
    setBreadcrumbs([]);
  };

  const handleNext = () => {
    if (!selectedCategory) {
      setError('category', t('categoryRequired'));
      return;
    }
    onNext();
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  const currentCategories = getCurrentCategories();

  return (
    <div className="w-full">
      <div className="text-center mb-6 lg:mb-8">
        <h2 className="text-xl lg:text-3xl font-bold text-base-content mb-2 lg:mb-4">
          {t('selectCategory')}
        </h2>
        <p className="text-sm lg:text-lg text-base-content/70">
          {t('categorySelectionDescription')}
        </p>
      </div>

      {/* Хлебные крошки */}
      {breadcrumbs.length > 0 && (
        <div className="mb-6">
          <div className="flex items-center gap-2 text-sm">
            <button
              onClick={navigateToRoot}
              className="text-primary hover:text-primary-focus transition-colors"
            >
              {t('allCategories')}
            </button>
            {breadcrumbs.map((crumb, index) => (
              <div key={crumb.id} className="flex items-center gap-2">
                <span className="text-base-content/40">›</span>
                <button
                  onClick={() => {
                    const newBreadcrumbs = breadcrumbs.slice(0, index + 1);
                    setBreadcrumbs(newBreadcrumbs);
                    setCurrentParentId(crumb.id || null);
                  }}
                  className="text-primary hover:text-primary-focus transition-colors"
                >
                  {crumb.name}
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Кнопка назад */}
      {breadcrumbs.length > 0 && (
        <div className="mb-4">
          <button onClick={navigateBack} className="btn btn-sm btn-ghost gap-2">
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
                d="M15 19l-7-7 7-7"
              />
            </svg>
            {tCommon('back')}
          </button>
        </div>
      )}

      {/* Текущие категории */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 lg:gap-4 mb-6 lg:mb-8">
        {currentCategories.map((category) => {
          const hasSubcategories = category.id
            ? hasChildren(category.id)
            : false;
          const isSelected = selectedCategory?.id === category.id;

          return (
            <div
              key={`category-${category.id}`}
              onClick={() => navigateToCategory(category)}
              className={`
                p-3 sm:p-4 lg:p-6 rounded-lg lg:rounded-xl border-2 cursor-pointer 
                transition-all duration-200 hover:shadow-lg active:scale-95
                ${
                  isSelected
                    ? 'border-primary bg-primary/10 shadow-lg sm:scale-105'
                    : 'border-base-300 bg-base-100 hover:border-primary/50'
                }
              `}
            >
              <div className="flex items-center gap-3 sm:gap-4">
                {/* Иконка категории */}
                <div
                  className={`
                  w-10 h-10 sm:w-12 sm:h-12 rounded-lg flex items-center justify-center 
                  text-xl sm:text-2xl flex-shrink-0
                  ${
                    isSelected
                      ? 'bg-primary text-primary-content'
                      : 'bg-base-200 text-base-content'
                  }
                `}
                >
                  {category.icon || '📦'}
                </div>

                {/* Название категории */}
                <div className="flex-1 min-w-0">
                  <h3
                    className={`
                    font-semibold text-sm sm:text-base lg:text-lg leading-tight truncate
                    ${isSelected ? 'text-primary' : 'text-base-content'}
                  `}
                  >
                    {category.name}
                  </h3>

                  {/* Количество товаров */}
                  {category.listing_count !== undefined && (
                    <p className="text-xs sm:text-sm text-base-content/60 mt-1 truncate">
                      {category.listing_count}{' '}
                      {t('productsCount')}
                    </p>
                  )}

                  {/* Индикатор подкатегорий */}
                  {hasSubcategories && (
                    <p className="text-xs text-primary mt-1">
                      {t('hasSubcategories')}
                    </p>
                  )}
                </div>

                {/* Индикатор выбора или стрелка */}
                <div className="flex-shrink-0">
                  {isSelected ? (
                    <div className="w-5 h-5 sm:w-6 sm:h-6 rounded-full bg-primary flex items-center justify-center">
                      <svg
                        className="w-3 h-3 sm:w-4 sm:h-4 text-primary-content"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                          clipRule="evenodd"
                        />
                      </svg>
                    </div>
                  ) : hasSubcategories ? (
                    <svg
                      className="w-5 h-5 text-base-content/40"
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
                  ) : null}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Ошибка */}
      {state.errors.category && (
        <div className="alert alert-error mb-6">
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
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{state.errors.category}</span>
        </div>
      )}

      {/* Выбранная категория */}
      {selectedCategory && (
        <div className="bg-success/10 border border-success/20 rounded-lg lg:rounded-xl p-4 lg:p-6 mb-6">
          <div className="flex items-center gap-3 lg:gap-4">
            <div className="w-12 h-12 lg:w-16 lg:h-16 rounded-lg lg:rounded-xl bg-success/20 flex items-center justify-center text-2xl lg:text-3xl flex-shrink-0">
              {selectedCategory.icon || '📦'}
            </div>
            <div className="min-w-0 flex-1">
              <h3 className="font-bold text-base lg:text-xl text-success">
                {t('selectedCategory')}
              </h3>
              <p className="text-sm lg:text-lg text-base-content truncate">
                {selectedCategory.name}
              </p>
              {/* Показать полный путь */}
              <div className="text-xs text-base-content/60 mt-1">
                {breadcrumbs.map((crumb) => crumb.name).join(' › ')}
                {breadcrumbs.length > 0 && ' › '}
                {selectedCategory.name}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Кнопка продолжить */}
      <div className="flex justify-end">
        <button
          onClick={handleNext}
          disabled={!selectedCategory}
          className={`
            btn btn-md lg:btn-lg px-6 lg:px-8 min-w-0
            ${selectedCategory ? 'btn-primary' : 'btn-disabled'}
          `}
        >
          <span className="text-sm lg:text-base">{tCommon('next')}</span>
          <svg
            className="w-4 h-4 lg:w-5 lg:h-5 ml-1 lg:ml-2"
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
      </div>
    </div>
  );
}
