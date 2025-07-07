'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';
import { tokenManager } from '@/utils/tokenManager';

interface SearchWeight {
  field: string;
  weight: number;
  description: string;
}

interface CategoryWeight {
  categoryId: string;
  categoryName: string;
  weights: {
    title: number;
    description: number;
    attributes: number;
  };
}

export default function SearchWeights() {
  const t = useTranslations();
  const [globalWeights, setGlobalWeights] = useState<SearchWeight[]>([]);
  const [categoryWeights, setCategoryWeights] = useState<CategoryWeight[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetchWeights();
  }, []);

  const fetchWeights = async () => {
    try {
      const accessToken = await tokenManager.getAccessToken();
      const response = await fetch('/api/admin/search/config/weights', {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });
      if (!response.ok) throw new Error('Failed to fetch weights');
      const data = await response.json();

      setGlobalWeights(
        data.globalWeights || [
          {
            field: 'title',
            weight: 10,
            description: 'Вес заголовка объявления',
          },
          { field: 'description', weight: 5, description: 'Вес описания' },
          { field: 'attributes', weight: 3, description: 'Вес атрибутов' },
          { field: 'category', weight: 8, description: 'Вес категории' },
          { field: 'tags', weight: 4, description: 'Вес тегов' },
        ]
      );

      setCategoryWeights(data.categoryWeights || []);
    } catch (error) {
      console.error('Error fetching weights:', error);
      toast.error(t('admin.search.weights.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const handleGlobalWeightChange = (index: number, value: string) => {
    const numValue = parseInt(value) || 0;
    if (numValue >= 0 && numValue <= 100) {
      const updated = [...globalWeights];
      updated[index].weight = numValue;
      setGlobalWeights(updated);
    }
  };

  const handleCategoryWeightChange = (
    categoryIndex: number,
    field: keyof CategoryWeight['weights'],
    value: string
  ) => {
    const numValue = parseInt(value) || 0;
    if (numValue >= 0 && numValue <= 100) {
      const updated = [...categoryWeights];
      updated[categoryIndex].weights[field] = numValue;
      setCategoryWeights(updated);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      const accessToken = await tokenManager.getAccessToken();
      const response = await fetch('/api/admin/search/weights', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({ globalWeights, categoryWeights }),
      });

      if (!response.ok) throw new Error('Failed to save weights');
      toast.success(t('admin.search.weights.saveSuccess'));
    } catch (error) {
      console.error('Error saving weights:', error);
      toast.error(t('admin.search.weights.saveError'));
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <div className="loading loading-spinner loading-lg"></div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold mb-4">
          {t('admin.search.weights.global')}
        </h2>
        <div className="overflow-x-auto">
          <table className="table table-zebra">
            <thead>
              <tr>
                <th>{t('admin.search.weights.field')}</th>
                <th>{t('admin.search.weights.weight')}</th>
                <th>{t('admin.search.weights.description')}</th>
              </tr>
            </thead>
            <tbody>
              {globalWeights.map((weight, index) => (
                <tr key={weight.field}>
                  <td className="font-mono">{weight.field}</td>
                  <td>
                    <input
                      type="range"
                      min="0"
                      max="100"
                      value={weight.weight}
                      onChange={(e) =>
                        handleGlobalWeightChange(index, e.target.value)
                      }
                      className="range range-primary range-sm"
                    />
                    <span className="ml-2 badge badge-primary">
                      {weight.weight}
                    </span>
                  </td>
                  <td>{t(`admin.search.weights.${weight.field}Desc`)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {categoryWeights.length > 0 && (
        <div>
          <h2 className="text-2xl font-semibold mb-4">
            {t('admin.search.weights.perCategory')}
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {categoryWeights.map((category, categoryIndex) => (
              <div
                key={category.categoryId}
                className="card bg-base-100 shadow-xl"
              >
                <div className="card-body">
                  <h3 className="card-title">{category.categoryName}</h3>
                  <div className="space-y-2">
                    {Object.entries(category.weights).map(([field, value]) => (
                      <div key={field} className="form-control">
                        <label className="label">
                          <span className="label-text">
                            {t(`admin.search.weights.${field}`)}
                          </span>
                        </label>
                        <input
                          type="range"
                          min="0"
                          max="100"
                          value={value}
                          onChange={(e) =>
                            handleCategoryWeightChange(
                              categoryIndex,
                              field as keyof CategoryWeight['weights'],
                              e.target.value
                            )
                          }
                          className="range range-secondary range-sm"
                        />
                        <span className="text-center badge badge-secondary">
                          {value}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="flex justify-end">
        <button
          className={`btn btn-primary ${saving ? 'loading' : ''}`}
          onClick={handleSave}
          disabled={saving}
        >
          {t('admin.search.weights.save')}
        </button>
      </div>
    </div>
  );
}
