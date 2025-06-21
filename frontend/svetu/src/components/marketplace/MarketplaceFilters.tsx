'use client';

import { useState, useEffect } from 'react';
import { MarketplaceService } from '@/services/marketplace';
import { useConfig } from '@/hooks/useConfig';
import type { components } from '@/types/generated/api';

interface MarketplaceFiltersProps {
  onFilterChange: (filters: { category?: string }) => void;
}

type Category = components['schemas']['models.MarketplaceCategory'];

export function MarketplaceFilters({
  onFilterChange,
}: MarketplaceFiltersProps) {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const config = useConfig();

  // Загрузка категорий
  useEffect(() => {
    const loadCategories = async () => {
      try {
        setLoading(true);
        const response = await MarketplaceService.getCategories();
        if (response.data) {
          setCategories(response.data as Category[]);
        }
      } catch (err) {
        console.error('Failed to load categories:', err);
        setError(true);
      } finally {
        setLoading(false);
      }
    };

    loadCategories();
  }, []);

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h3 className="card-title text-lg">Filters</h3>

        {/* Debug info в development */}
        {config.env.isDevelopment && (
          <div className="text-xs text-base-content/50 mb-2">
            API: {config.api.url}
          </div>
        )}

        {loading ? (
          <div className="skeleton h-12 w-full"></div>
        ) : error ? (
          <div className="alert alert-error">
            <span>Failed to load categories</span>
          </div>
        ) : (
          <select
            onChange={(e) => onFilterChange({ category: e.target.value })}
            className="select select-bordered w-full"
          >
            <option value="">All Categories</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        )}
      </div>
    </div>
  );
}
