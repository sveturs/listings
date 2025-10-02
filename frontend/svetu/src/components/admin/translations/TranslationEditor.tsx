'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';
import {
  PencilIcon,
  CheckIcon,
  XMarkIcon,
  ExclamationTriangleIcon,
  MagnifyingGlassIcon,
  ArrowUpIcon,
  ArrowDownIcon,
} from '@heroicons/react/24/outline';
import { debounce } from 'lodash';
import { apiClient } from '@/services/api-client';

interface Translation {
  module: string;
  key: string;
  path: string;
  translations: Record<string, string>;
  status: 'complete' | 'incomplete' | 'placeholder' | 'missing';
  metadata?: any;
}

interface TranslationEditorProps {
  module: string;
}

export default function TranslationEditor({ module }: TranslationEditorProps) {
  const t = useTranslations('admin');
  const [translations, setTranslations] = useState<Translation[]>([]);
  const [filteredTranslations, setFilteredTranslations] = useState<
    Translation[]
  >([]);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [sortBy, setSortBy] = useState<'key' | 'status'>('key');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [editingKey, setEditingKey] = useState<string | null>(null);
  const [editValues, setEditValues] = useState<Record<string, string>>({});
  const [hasChanges, setHasChanges] = useState(false);

  const languages = ['sr', 'en', 'ru'];

  useEffect(() => {
    if (module) {
      fetchTranslations();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [module]);

  useEffect(() => {
    filterAndSortTranslations();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [translations, searchTerm, statusFilter, sortBy, sortOrder]);

  const fetchTranslations = async () => {
    setLoading(true);
    try {
      const response = await apiClient.get(
        `/admin/translations/frontend/module/${module}`
      );

      if (!response.data) {
        throw new Error('Failed to fetch translations');
      }

      setTranslations(response.data.data || []);
    } catch {
      toast.error(t('translations.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const filterAndSortTranslations = () => {
    let filtered = [...translations];

    // Search filter
    if (searchTerm) {
      filtered = filtered.filter(
        (trans) =>
          trans.key.toLowerCase().includes(searchTerm.toLowerCase()) ||
          Object.values(trans.translations).some((val) =>
            val.toLowerCase().includes(searchTerm.toLowerCase())
          )
      );
    }

    // Status filter
    if (statusFilter !== 'all') {
      filtered = filtered.filter((trans) => trans.status === statusFilter);
    }

    // Sort
    filtered.sort((a, b) => {
      let comparison = 0;
      if (sortBy === 'key') {
        comparison = a.key.localeCompare(b.key);
      } else if (sortBy === 'status') {
        comparison = a.status.localeCompare(b.status);
      }
      return sortOrder === 'asc' ? comparison : -comparison;
    });

    setFilteredTranslations(filtered);
  };

  const startEditing = (key: string, translations: Record<string, string>) => {
    setEditingKey(key);
    setEditValues(translations);
  };

  const cancelEditing = () => {
    setEditingKey(null);
    setEditValues({});
  };

  const saveTranslation = async (translationKey: string) => {
    setSaving(true);
    try {
      const translationToUpdate = translations.find(
        (t) => t.key === translationKey
      );
      if (!translationToUpdate) return;

      const updatedTranslation = {
        ...translationToUpdate,
        translations: editValues,
      };

      const response = await apiClient.put(
        `/admin/translations/frontend/module/${module}`,
        [updatedTranslation]
      );

      if (!response.data) {
        throw new Error('Failed to save translation');
      }

      // Update local state
      setTranslations((prev) =>
        prev.map((t) =>
          t.key === translationKey
            ? { ...t, translations: editValues, status: getStatus(editValues) }
            : t
        )
      );

      setEditingKey(null);
      setEditValues({});
      setHasChanges(true);
      toast.success(t('translations.saved'));
    } catch {
      toast.error(t('translations.saveError'));
    } finally {
      setSaving(false);
    }
  };

  const getStatus = (
    translations: Record<string, string>
  ): Translation['status'] => {
    const values = Object.values(translations);
    if (
      values.some(
        (v) => v.includes('[RU]') || v.includes('[EN]') || v.includes('[SR]')
      )
    ) {
      return 'placeholder';
    }
    if (values.some((v) => !v)) {
      return 'missing';
    }
    if (values.length < languages.length) {
      return 'incomplete';
    }
    return 'complete';
  };

  const getStatusBadge = (status: Translation['status']) => {
    switch (status) {
      case 'complete':
        return (
          <span className="badge badge-success badge-sm">
            {t('translations.complete')}
          </span>
        );
      case 'incomplete':
        return (
          <span className="badge badge-warning badge-sm">
            {t('translations.incomplete')}
          </span>
        );
      case 'placeholder':
        return (
          <span className="badge badge-error badge-sm">
            {t('translations.placeholder')}
          </span>
        );
      case 'missing':
        return (
          <span className="badge badge-ghost badge-sm">
            {t('translations.missing')}
          </span>
        );
    }
  };

  const debouncedSearch = useCallback(
    (value: string) => {
      const debounced = debounce((val: string) => setSearchTerm(val), 300);
      debounced(value);
    },

    []
  );

  if (loading) {
    return (
      <div className="bg-base-100 rounded-lg p-8">
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-base-100 rounded-lg p-4">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold">
          {t('translations.editing')} {module}
        </h3>
        {hasChanges && (
          <span className="text-sm text-success">
            {t('translations.changesSaved')}
          </span>
        )}
      </div>

      {/* Filters and Search */}
      <div className="flex flex-col lg:flex-row gap-4 mb-4">
        <div className="flex-1">
          <div className="input-group">
            <span>
              <MagnifyingGlassIcon className="h-5 w-5" />
            </span>
            <input
              type="text"
              placeholder={t('translations.searchKeys')}
              className="input input-bordered w-full"
              onChange={(e) => debouncedSearch(e.target.value)}
            />
          </div>
        </div>

        <select
          className="select select-bordered"
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value)}
        >
          <option value="all">{t('translations.allStatuses')}</option>
          <option value="complete">{t('translations.complete')}</option>
          <option value="incomplete">{t('translations.incomplete')}</option>
          <option value="placeholder">{t('translations.placeholder')}</option>
          <option value="missing">{t('translations.missing')}</option>
        </select>

        <div className="btn-group">
          <button
            className={`btn btn-sm ${sortBy === 'key' ? 'btn-active' : ''}`}
            onClick={() => setSortBy('key')}
          >
            {t('translations.sortByKey')}
          </button>
          <button
            className={`btn btn-sm ${sortBy === 'status' ? 'btn-active' : ''}`}
            onClick={() => setSortBy('status')}
          >
            {t('translations.sortByStatus')}
          </button>
          <button
            className="btn btn-sm"
            onClick={() => setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')}
          >
            {sortOrder === 'asc' ? (
              <ArrowUpIcon className="h-4 w-4" />
            ) : (
              <ArrowDownIcon className="h-4 w-4" />
            )}
          </button>
        </div>
      </div>

      {/* Translations Table */}
      <div className="overflow-x-auto">
        <table className="table table-sm">
          <thead>
            <tr>
              <th className="w-1/4">{t('translations.key')}</th>
              {languages.map((lang) => (
                <th key={lang} className="w-1/4">
                  {lang.toUpperCase()}
                </th>
              ))}
              <th className="w-32">{t('translations.status')}</th>
              <th className="w-20">{t('translations.actions')}</th>
            </tr>
          </thead>
          <tbody>
            {filteredTranslations.map((trans) => (
              <tr key={trans.key} className="hover">
                <td className="font-mono text-xs">{trans.key}</td>
                {languages.map((lang) => (
                  <td key={lang}>
                    {editingKey === trans.key ? (
                      <input
                        type="text"
                        className="input input-bordered input-sm w-full"
                        value={editValues[lang] || ''}
                        onChange={(e) =>
                          setEditValues((prev) => ({
                            ...prev,
                            [lang]: e.target.value,
                          }))
                        }
                      />
                    ) : (
                      <div className="text-sm">
                        {trans.translations[lang] || (
                          <span className="text-base-content/40">
                            {t('translations.empty')}
                          </span>
                        )}
                        {trans.translations[lang]?.includes(
                          `[${lang.toUpperCase()}]`
                        ) && (
                          <ExclamationTriangleIcon className="h-4 w-4 text-warning inline ml-1" />
                        )}
                      </div>
                    )}
                  </td>
                ))}
                <td>{getStatusBadge(trans.status)}</td>
                <td>
                  {editingKey === trans.key ? (
                    <div className="flex gap-1">
                      <button
                        className="btn btn-ghost btn-xs"
                        onClick={() => saveTranslation(trans.key)}
                        disabled={saving}
                      >
                        {saving ? (
                          <span className="loading loading-spinner loading-xs"></span>
                        ) : (
                          <CheckIcon className="h-4 w-4 text-success" />
                        )}
                      </button>
                      <button
                        className="btn btn-ghost btn-xs"
                        onClick={cancelEditing}
                        disabled={saving}
                      >
                        <XMarkIcon className="h-4 w-4 text-error" />
                      </button>
                    </div>
                  ) : (
                    <button
                      className="btn btn-ghost btn-xs"
                      onClick={() =>
                        startEditing(trans.key, trans.translations)
                      }
                    >
                      <PencilIcon className="h-4 w-4" />
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {filteredTranslations.length === 0 && (
        <div className="text-center py-8 text-base-content/60">
          {t('translations.noResults')}
        </div>
      )}
    </div>
  );
}
