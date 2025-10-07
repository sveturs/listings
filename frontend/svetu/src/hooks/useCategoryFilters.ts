import { useState, useEffect, useCallback } from 'react';
import type { components } from '@/types/generated/api';

type CategoryAttribute =
  components['schemas']['models.CategoryAttribute'];

interface UseCategoryFiltersOptions {
  lang?: string;
  onError?: (error: Error) => void;
}

interface UseCategoryFiltersResult {
  attributes: CategoryAttribute[];
  loading: boolean;
  error: Error | null;
  refetch: () => void;
}

export function useCategoryFilters(
  categoryId: number | null | undefined,
  options: UseCategoryFiltersOptions = {}
): UseCategoryFiltersResult {
  const [attributes, setAttributes] = useState<CategoryAttribute[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const { lang = 'sr', onError } = options;

  const fetchAttributes = useCallback(async () => {
    if (!categoryId) {
      setAttributes([]);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `/api/v1/marketplace/categories/${categoryId}/attributes?lang=${lang}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to fetch attributes: ${response.status}`);
      }

      const data = await response.json();

      if (data.data && Array.isArray(data.data)) {
        // Сортируем атрибуты по sort_order
        const sortedAttributes = data.data.sort(
          (a: CategoryAttribute, b: CategoryAttribute) =>
            (a.sort_order || 0) - (b.sort_order || 0)
        );
        setAttributes(sortedAttributes);
      } else {
        setAttributes([]);
      }
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Unknown error');
      setError(error);
      onError?.(error);
      console.error('Error fetching category attributes:', error);
    } finally {
      setLoading(false);
    }
  }, [categoryId, lang, onError]);

  useEffect(() => {
    fetchAttributes();
  }, [categoryId, lang, fetchAttributes]);

  return {
    attributes,
    loading,
    error,
    refetch: fetchAttributes,
  };
}
