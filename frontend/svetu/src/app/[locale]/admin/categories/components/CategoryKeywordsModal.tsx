import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import {
  X,
  Loader2,
  Tags,
  TrendingUp,
  Languages,
  Trash2,
} from 'lucide-react';
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

interface CategoryKeywordsModalProps {
  categoryId: number;
  categoryName: string;
  isOpen: boolean;
  onClose: () => void;
}

export default function CategoryKeywordsModal({
  categoryId,
  categoryName,
  isOpen,
  onClose,
}: CategoryKeywordsModalProps) {
  const t = useTranslations('admin.categories.keywords');
  const [keywords, setKeywords] = useState<CategoryKeyword[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [newKeyword, setNewKeyword] = useState({
    keyword: '',
    language: 'ru',
    weight: 1.0,
    keyword_type: 'main',
    is_negative: false,
  });

  const loadKeywords = useCallback(async () => {
    if (!isOpen || !categoryId) return;

    try {
      setLoading(true);
      const response = await fetch(
        `/api/v1/admin/categories/${categoryId}/keywords`
      );

      if (response.ok) {
        const data = await response.json();
        // Обрабатываем разные форматы ответа API
        if (Array.isArray(data)) {
          setKeywords(data);
        } else if (data && Array.isArray(data.data)) {
          setKeywords(data.data);
        } else {
          setKeywords([]);
        }
      } else {
        throw new Error('Failed to load keywords');
      }
    } catch (error) {
      console.error('Error loading keywords:', error);
      toast.error(t('loadError'));
      setKeywords([]); // Устанавливаем пустой массив при ошибке
    } finally {
      setLoading(false);
    }
  }, [categoryId, isOpen, t]);

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
          body: JSON.stringify(newKeyword),
        }
      );

      if (response.ok) {
        await loadKeywords();
        setNewKeyword({
          keyword: '',
          language: 'ru',
          weight: 1.0,
          keyword_type: 'main',
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

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box w-11/12 max-w-6xl max-h-[90vh]">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <Tags className="w-6 h-6 text-primary" />
            <div>
              <h3 className="text-xl font-bold">{t('title')}</h3>
              <p className="text-sm text-base-content/60">{categoryName}</p>
            </div>
          </div>
          <button onClick={onClose} className="btn btn-ghost btn-sm btn-circle">
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Add keyword form */}
        <div className="bg-base-200 rounded-lg p-6 mb-6">
          <h4 className="text-lg font-medium mb-4">{t('addKeyword')}</h4>
          <div className="grid grid-cols-1 lg:grid-cols-6 gap-4">
            <div className="lg:col-span-2">
              <label className="label">
                <span className="label-text">{t('keywordPlaceholder')}</span>
              </label>
              <input
                type="text"
                placeholder={t('keywordPlaceholder')}
                className="input input-bordered w-full"
                value={newKeyword.keyword}
                onChange={(e) =>
                  setNewKeyword({ ...newKeyword, keyword: e.target.value })
                }
              />
            </div>

            <div>
              <label className="label">
                <span className="label-text">Язык</span>
              </label>
              <select
                className="select select-bordered w-full"
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
            </div>

            <div>
              <label className="label">
                <span className="label-text">Тип</span>
              </label>
              <select
                className="select select-bordered w-full"
                value={newKeyword.keyword_type}
                onChange={(e) =>
                  setNewKeyword({ ...newKeyword, keyword_type: e.target.value })
                }
              >
                <option value="main">{t('types.main')}</option>
                <option value="synonym">{t('types.synonym')}</option>
                <option value="brand">{t('types.brand')}</option>
                <option value="attribute">{t('types.attribute')}</option>
                <option value="context">{t('types.context')}</option>
                <option value="pattern">{t('types.pattern')}</option>
              </select>
            </div>

            <div>
              <label className="label">
                <span className="label-text">{t('weight')}</span>
              </label>
              <input
                type="number"
                className="input input-bordered w-full"
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
            </div>

            <div className="flex items-end">
              <button
                onClick={addKeyword}
                disabled={saving || !newKeyword.keyword.trim()}
                className="btn btn-primary w-full"
              >
                {saving && <Loader2 className="w-4 h-4 animate-spin mr-2" />}
                {t('add')}
              </button>
            </div>
          </div>

          <div className="form-control mt-4">
            <label className="label cursor-pointer justify-start gap-3">
              <input
                type="checkbox"
                className="checkbox"
                checked={newKeyword.is_negative}
                onChange={(e) =>
                  setNewKeyword({
                    ...newKeyword,
                    is_negative: e.target.checked,
                  })
                }
              />
              <span className="label-text">{t('negative')}</span>
            </label>
          </div>
        </div>

        {/* Keywords table */}
        <div className="overflow-x-auto">
          {loading ? (
            <div className="flex items-center justify-center py-12">
              <Loader2 className="w-8 h-8 animate-spin" />
            </div>
          ) : !Array.isArray(keywords) || keywords.length === 0 ? (
            <div className="text-center py-12 text-base-content/60">
              <Tags className="w-12 h-12 mx-auto mb-4 opacity-50" />
              <p>{t('noKeywords')}</p>
            </div>
          ) : (
            <table className="table table-zebra w-full">
              <thead>
                <tr>
                  <th>Ключевое слово</th>
                  <th>Язык</th>
                  <th>Тип</th>
                  <th>Вес</th>
                  <th>Статистика</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                {Array.isArray(keywords) && keywords.map((keyword) => (
                  <tr
                    key={keyword.id}
                    className={keyword.is_negative ? 'bg-error/5' : ''}
                  >
                    <td>
                      <div className="flex items-center gap-2">
                        <span className="font-medium">{keyword.keyword}</span>
                        {keyword.is_negative && (
                          <span className="badge badge-error badge-sm">
                            {t('negative')}
                          </span>
                        )}
                      </div>
                    </td>
                    <td>
                      <div className="flex items-center gap-2">
                        <Languages className="w-4 h-4" />
                        <span>
                          {keyword.language === '*'
                            ? t('allLanguages')
                            : keyword.language.toUpperCase()}
                        </span>
                      </div>
                    </td>
                    <td>
                      <span className="badge badge-outline">
                        {t(`types.${keyword.keyword_type}`)}
                      </span>
                    </td>
                    <td>
                      <input
                        type="number"
                        className="input input-bordered input-sm w-20"
                        value={keyword.weight}
                        min="0"
                        max="10"
                        step="0.1"
                        onChange={(e) =>
                          updateWeight(keyword.id, parseFloat(e.target.value))
                        }
                      />
                    </td>
                    <td>
                      <div className="flex flex-col gap-1">
                        <div className="flex items-center gap-2 text-sm">
                          <span>
                            {t('usage')}: {keyword.usage_count}
                          </span>
                        </div>
                        <div className="flex items-center gap-2 text-sm">
                          <TrendingUp className="w-4 h-4 text-success" />
                          <span>
                            {(keyword.success_rate * 100).toFixed(1)}%
                          </span>
                        </div>
                      </div>
                    </td>
                    <td>
                      <button
                        onClick={() => deleteKeyword(keyword.id)}
                        className="btn btn-ghost btn-sm text-error"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>

        {/* Info section */}
        <div className="mt-6 p-4 bg-info/10 rounded-lg">
          <h4 className="font-medium mb-2">{t('info.title')}</h4>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm text-base-content/80">
            <div>• {t('info.weights')}</div>
            <div>• {t('info.types')}</div>
            <div>• {t('info.negative')}</div>
            <div>• {t('info.languages')}</div>
            <div className="md:col-span-2">• {t('info.statistics')}</div>
          </div>
        </div>

        {/* Footer */}
        <div className="modal-action">
          <button onClick={onClose} className="btn">
            Закрыть
          </button>
        </div>
      </div>
    </div>
  );
}
