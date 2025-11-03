import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { Plus, X, Loader2, Tags, TrendingUp, Languages } from 'lucide-react';
import { toast } from '@/utils/toast';

interface CategoryKeyword {
  id: number;
  keyword: string;
  language: string;
  weight: number;
  keyword_type: string;
  is_negative: boolean;
  usage_count: number;
  success_rate: number;
}

interface CategoryKeywordsProps {
  categoryId: number;
  categoryName: string;
}

export default function CategoryKeywords({
  categoryId,
  categoryName,
}: CategoryKeywordsProps) {
  const t = useTranslations('admin');
  const [keywords, setKeywords] = useState<CategoryKeyword[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [newKeyword, setNewKeyword] = useState({
    keyword: '',
    language: 'ru',
    weight: 1.0,
    keyword_type: 'general',
    is_negative: false,
  });

  const loadKeywords = useCallback(async () => {
    try {
      setLoading(true);
      const response = await fetch(
        `/api/v1/marketplace/categories/${categoryId}/keywords`
      );
      if (response.ok) {
        const data = await response.json();
        setKeywords(data.data || []);
      }
    } catch (error) {
      console.error('Error loading keywords:', error);
      toast.error(t('loadError'));
    } finally {
      setLoading(false);
    }
  }, [categoryId, t]);

  useEffect(() => {
    loadKeywords();
  }, [loadKeywords]);

  const addKeyword = async () => {
    if (!newKeyword.keyword.trim()) {
      toast.error(t('keywordRequired'));
      return;
    }

    try {
      setSaving(true);
      const response = await fetch(
        `/api/v1/admin/categories/${categoryId}/keywords`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            ...newKeyword,
            category_id: categoryId,
          }),
        }
      );

      if (response.ok) {
        await loadKeywords();
        setNewKeyword({
          keyword: '',
          language: 'ru',
          weight: 1.0,
          keyword_type: 'general',
          is_negative: false,
        });
        toast.success(t('addSuccess'));
      } else {
        throw new Error('Failed to add keyword');
      }
    } catch (error) {
      console.error('Error adding keyword:', error);
      toast.error(t('addError'));
    } finally {
      setSaving(false);
    }
  };

  const deleteKeyword = async (keywordId: number) => {
    if (!confirm(t('deleteConfirm'))) return;

    try {
      const response = await fetch(
        `/api/v1/admin/categories/keywords/${keywordId}`,
        {
          method: 'DELETE',
        }
      );

      if (response.ok) {
        await loadKeywords();
        toast.success(t('deleteSuccess'));
      } else {
        throw new Error('Failed to delete keyword');
      }
    } catch (error) {
      console.error('Error deleting keyword:', error);
      toast.error(t('deleteError'));
    }
  };

  const updateWeight = async (keywordId: number, newWeight: number) => {
    try {
      const response = await fetch(
        `/api/v1/admin/categories/keywords/${keywordId}`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ weight: newWeight }),
        }
      );

      if (response.ok) {
        await loadKeywords();
        toast.success(t('updateSuccess'));
      }
    } catch (error) {
      console.error('Error updating weight:', error);
      toast.error(t('updateError'));
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <Loader2 className="w-6 h-6 animate-spin" />
      </div>
    );
  }

  return (
    <div className="bg-base-100 rounded-lg p-6">
      <div className="flex items-center gap-2 mb-6">
        <Tags className="w-5 h-5 text-primary" />
        <h3 className="text-lg font-semibold">{t('title')}</h3>
        <span className="text-sm text-base-content/60">({categoryName})</span>
      </div>

      {/* Add new keyword form */}
      <div className="bg-base-200 rounded-lg p-4 mb-6">
        <h4 className="font-medium mb-3">{t('addKeyword')}</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-3">
          <input
            type="text"
            placeholder={t('keywordPlaceholder')}
            className="input input-bordered input-sm"
            value={newKeyword.keyword}
            onChange={(e) =>
              setNewKeyword({ ...newKeyword, keyword: e.target.value })
            }
          />

          <select
            className="select select-bordered select-sm"
            value={newKeyword.language}
            onChange={(e) =>
              setNewKeyword({ ...newKeyword, language: e.target.value })
            }
          >
            <option value="ru">Русский</option>
            <option value="en">English</option>
            <option value="sr">Српски</option>
            <option value="*">{t('allLanguages')}</option>
          </select>

          <select
            className="select select-bordered select-sm"
            value={newKeyword.keyword_type}
            onChange={(e) =>
              setNewKeyword({ ...newKeyword, keyword_type: e.target.value })
            }
          >
            <option value="general">{t('types.general')}</option>
            <option value="main">{t('types.main')}</option>
            <option value="synonym">{t('types.synonym')}</option>
            <option value="brand">{t('types.brand')}</option>
            <option value="attribute">{t('types.attribute')}</option>
          </select>

          <input
            type="number"
            placeholder={t('weight')}
            className="input input-bordered input-sm"
            value={newKeyword.weight}
            min="0"
            max="10"
            step="0.1"
            onChange={(e) =>
              setNewKeyword({
                ...newKeyword,
                weight: parseFloat(e.target.value),
              })
            }
          />

          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              className="checkbox checkbox-sm"
              checked={newKeyword.is_negative}
              onChange={(e) =>
                setNewKeyword({ ...newKeyword, is_negative: e.target.checked })
              }
            />
            <span className="text-sm">{t('negative')}</span>
          </label>

          <button
            onClick={addKeyword}
            disabled={saving || !newKeyword.keyword.trim()}
            className="btn btn-primary btn-sm"
          >
            {saving ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Plus className="w-4 h-4" />
            )}
            {t('add')}
          </button>
        </div>
      </div>

      {/* Keywords list */}
      <div className="space-y-2">
        {keywords.length === 0 ? (
          <div className="text-center py-8 text-base-content/60">
            {t('noKeywords')}
          </div>
        ) : (
          keywords.map((keyword) => (
            <div
              key={keyword.id}
              className={`flex items-center gap-3 p-3 rounded-lg ${
                keyword.is_negative ? 'bg-error/10' : 'bg-base-200'
              }`}
            >
              <div className="flex-1">
                <span className="font-medium">{keyword.keyword}</span>
                <div className="flex items-center gap-4 text-sm text-base-content/60 mt-1">
                  <span className="flex items-center gap-1">
                    <Languages className="w-3 h-3" />
                    {keyword.language === '*'
                      ? t('allLanguages')
                      : keyword.language.toUpperCase()}
                  </span>
                  <span className="badge badge-sm">
                    {t(`types.${keyword.keyword_type}`)}
                  </span>
                  {keyword.is_negative && (
                    <span className="badge badge-error badge-sm">
                      {t('negative')}
                    </span>
                  )}
                  <span className="flex items-center gap-1">
                    {t('usage')}: {keyword.usage_count}
                  </span>
                  <span className="flex items-center gap-1">
                    <TrendingUp className="w-3 h-3" />
                    {(keyword.success_rate * 100).toFixed(1)}%
                  </span>
                </div>
              </div>

              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1">
                  <span className="text-sm">{t('weight')}:</span>
                  <input
                    type="number"
                    className="input input-bordered input-xs w-16"
                    value={keyword.weight}
                    min="0"
                    max="10"
                    step="0.1"
                    onChange={(e) =>
                      updateWeight(keyword.id, parseFloat(e.target.value))
                    }
                  />
                </div>

                <button
                  onClick={() => deleteKeyword(keyword.id)}
                  className="btn btn-ghost btn-sm text-error"
                >
                  <X className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Info section */}
      <div className="mt-6 p-4 bg-info/10 rounded-lg">
        <h4 className="font-medium mb-2">{t('info.title')}</h4>
        <ul className="text-sm space-y-1 text-base-content/80">
          <li>• {t('info.weights')}</li>
          <li>• {t('info.types')}</li>
          <li>• {t('info.negative')}</li>
          <li>• {t('info.languages')}</li>
          <li>• {t('info.statistics')}</li>
        </ul>
      </div>
    </div>
  );
}
