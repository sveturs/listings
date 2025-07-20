'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';
import { MarketplaceService } from '@/services/marketplace';
import { toast } from '@/utils/toast';

interface Category {
  id: number;
  name: string;
  icon?: string;
  slug?: string;
  parent_id?: number | null;
  translations?: Record<string, string>;
  level?: number;
  count?: number;
}

interface CategorySelectionStepProps {
  onNext: () => void;
}

export default function CategorySelectionStep({
  onNext,
}: CategorySelectionStepProps) {
  const t = useTranslations();
  const locale = useLocale();
  const { state, setCategory } = useCreateListing();
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(
    state.category
      ? {
          id: state.category.id,
          name: state.category.name,
        }
      : null
  );
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadCategories = async () => {
      try {
        setLoading(true);
        setError(null);

        // Загружаем категории из API с текущей локалью
        const response = await MarketplaceService.getCategories(locale);

        // Преобразуем категории для отображения
        const processedCategories: Category[] = response.data.map((cat) => ({
          id: cat.id,
          name: cat.name,
          icon: cat.icon || '📦',
          slug: cat.slug,
          parent_id: cat.parent_id,
          translations: cat.translations,
          level: cat.level || 0,
          count: cat.count || 0,
        }));

        // Фильтруем только корневые категории для отображения
        const rootCategories = processedCategories.filter(
          (cat) => !cat.parent_id
        );

        setCategories(rootCategories);
      } catch (err) {
        console.error('Error loading categories:', err);
        setError(t('create_listing.category.error_loading'));
        toast.error(t('create_listing.category.error_loading'));
      } finally {
        setLoading(false);
      }
    };

    loadCategories();
  }, [t, locale]);

  const handleCategorySelect = (category: Category) => {
    console.log('CategorySelectionStep - Selecting category:', category);
    setSelectedCategory(category);
    setCategory({
      id: category.id,
      name: category.name,
      slug: category.slug || `category-${category.id}`,
    });
    console.log('CategorySelectionStep - Category set in context');
  };

  const handleNext = () => {
    if (selectedCategory) {
      onNext();
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            🏪 {t('create_listing.category.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('create_listing.category.description')}
          </p>

          {/* Популярные категории */}
          <div className="mb-6">
            <h3 className="text-lg font-medium mb-3 flex items-center gap-2">
              ⭐ {t('create_listing.category.popular')}
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              {categories
                .slice(0, 6) // Показываем первые 6 категорий как популярные
                .map((category) => (
                  <label
                    key={`popular-${category.id}`}
                    className="cursor-pointer"
                  >
                    <input
                      type="radio"
                      name="category"
                      value={category.id}
                      checked={selectedCategory?.id === category.id}
                      onChange={() => handleCategorySelect(category)}
                      className="sr-only"
                    />
                    <div
                      className={`
                    card border-2 transition-all duration-200 hover:scale-105
                    ${
                      selectedCategory?.id === category.id
                        ? 'border-primary bg-primary/5 shadow-lg'
                        : 'border-base-300 hover:border-primary/50'
                    }
                  `}
                    >
                      <div className="card-body p-4 text-center">
                        <div className="text-3xl mb-2">{category.icon}</div>
                        <h4 className="font-medium text-sm mb-1">
                          {category.name}
                        </h4>
                        {category.count && category.count > 0 && (
                          <p className="text-xs text-base-content/60">
                            {category.count}{' '}
                            {t('create_listing.category.listings')}
                          </p>
                        )}
                        {selectedCategory?.id === category.id && (
                          <div className="mt-2">
                            <svg
                              className="w-6 h-6 text-primary mx-auto"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                clipRule="evenodd"
                              />
                            </svg>
                          </div>
                        )}
                      </div>
                    </div>
                  </label>
                ))}
            </div>
          </div>

          {/* Остальные категории */}
          <div>
            <h3 className="text-lg font-medium mb-3 flex items-center gap-2">
              📂 {t('create_listing.category.all')}
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
              {categories
                .slice(6) // Остальные категории
                .map((category) => (
                  <label key={`all-${category.id}`} className="cursor-pointer">
                    <input
                      type="radio"
                      name="category"
                      value={category.id}
                      checked={selectedCategory?.id === category.id}
                      onChange={() => handleCategorySelect(category)}
                      className="sr-only"
                    />
                    <div
                      className={`
                    card border-2 transition-all duration-200
                    ${
                      selectedCategory?.id === category.id
                        ? 'border-primary bg-primary/5'
                        : 'border-base-300 hover:border-primary/50'
                    }
                  `}
                    >
                      <div className="card-body p-3 text-center">
                        <div className="text-2xl mb-1">{category.icon}</div>
                        <h4 className="font-medium text-xs mb-1">
                          {category.name}
                        </h4>
                        {category.count && category.count > 0 && (
                          <p className="text-xs text-base-content/60 hidden sm:block">
                            {category.count}
                          </p>
                        )}
                        {selectedCategory?.id === category.id && (
                          <div className="mt-1">
                            <svg
                              className="w-4 h-4 text-primary mx-auto"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                clipRule="evenodd"
                              />
                            </svg>
                          </div>
                        )}
                      </div>
                    </div>
                  </label>
                ))}
            </div>
          </div>

          {/* Показываем ошибку если есть */}
          {error && (
            <div className="alert alert-error mb-4">
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
              <span>{error}</span>
            </div>
          )}

          {/* Подсказка в региональном стиле */}
          <div className="alert alert-info mt-6">
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
              <p className="font-medium">
                💡 {t('create_listing.category.tip_title')}
              </p>
              <p className="text-xs mt-1">
                {t('create_listing.category.tip_description')}
              </p>
            </div>
          </div>

          {/* Кнопка продолжения */}
          <div className="card-actions justify-end mt-6">
            <button
              className={`btn btn-primary ${!selectedCategory ? 'btn-disabled' : ''}`}
              onClick={handleNext}
              disabled={!selectedCategory}
            >
              {t('common.continue')} →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
