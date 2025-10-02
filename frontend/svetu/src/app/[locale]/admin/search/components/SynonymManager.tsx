'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { apiClientAuth } from '@/lib/api-client-auth';

interface Synonym {
  id: number;
  term: string;
  synonym: string;
  language: 'en' | 'ru' | 'sr';
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export default function SynonymManager() {
  const t = useTranslations('admin');
  const [synonyms, setSynonyms] = useState<Synonym[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedLanguage, setSelectedLanguage] = useState<'en' | 'ru' | 'sr'>(
    'ru'
  );
  const [searchTerm, setSearchTerm] = useState('');
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  // Состояние для создания/редактирования
  const [isEditing, setIsEditing] = useState<number | null>(null);
  const [editForm, setEditForm] = useState({
    term: '',
    synonym: '',
    language: 'ru' as 'en' | 'ru' | 'sr',
    is_active: true,
  });
  const [showCreateForm, setShowCreateForm] = useState(false);

  const fetchSynonyms = async () => {
    try {
      setLoading(true);

      const params = new URLSearchParams({
        page: page.toString(),
        limit: '20',
      });

      if (selectedLanguage) {
        params.append('language', selectedLanguage);
      }

      if (searchTerm.trim()) {
        params.append('search', searchTerm.trim());
      }

      const response = await apiClientAuth.get(
        `/api/v1/admin/search/synonyms?${params}`
      );

      if (response && response.success) {
        setSynonyms(response.data?.data || []);
        setTotal(response.data?.total || 0);
        setTotalPages(response.data?.total_pages || 1);
      }
    } catch (error: any) {
      console.error('Error fetching synonyms:', error);
      if (error.status === 401) {
        window.location.href = '/ru/login';
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSynonyms();
  }, [selectedLanguage, searchTerm, page]); // eslint-disable-line react-hooks/exhaustive-deps

  const handleCreate = async () => {
    if (!editForm.term.trim() || !editForm.synonym.trim()) {
      return;
    }

    try {
      await apiClientAuth.post('/api/v1/admin/search/synonyms', editForm);

      setEditForm({ term: '', synonym: '', language: 'ru', is_active: true });
      setShowCreateForm(false);
      fetchSynonyms();
    } catch (error: any) {
      console.error('Error creating synonym:', error);
      if (error.status === 401) {
        window.location.href = '/ru/login';
      }
    }
  };

  const handleUpdate = async (synonym: Synonym) => {
    try {
      await apiClientAuth.put(`/api/v1/admin/search/synonyms/${synonym.id}`, {
        term: synonym.term,
        synonym: synonym.synonym,
        language: synonym.language,
        is_active: synonym.is_active,
      });

      setIsEditing(null);
      fetchSynonyms();
    } catch (error: any) {
      console.error('Error updating synonym:', error);
      if (error.status === 401) {
        window.location.href = '/ru/login';
      }
    }
  };

  const handleDelete = async (synonymId: number) => {
    if (!confirm(t('confirmDelete'))) {
      return;
    }

    try {
      await apiClientAuth.delete(`/api/v1/admin/search/synonyms/${synonymId}`);
      fetchSynonyms();
    } catch (error: any) {
      console.error('Error deleting synonym:', error);
      if (error.status === 401) {
        window.location.href = '/ru/login';
      }
    }
  };

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
    setPage(1); // Сброс на первую страницу при поиске
  };

  const handleLanguageChange = (language: 'en' | 'ru' | 'sr') => {
    setSelectedLanguage(language);
    setPage(1); // Сброс на первую страницу при смене языка
  };

  const getLanguageLabel = (lang: string) => {
    switch (lang) {
      case 'en':
        return 'English';
      case 'ru':
        return 'Русский';
      case 'sr':
        return 'Српски';
      default:
        return lang;
    }
  };

  return (
    <div className="space-y-6">
      {/* Заголовок */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-base-content">{t('title')}</h2>
          <p className="text-base-content/70 mt-1">{t('description')}</p>
        </div>
        <button
          className="btn btn-primary"
          onClick={() => setShowCreateForm(true)}
        >
          {t('addSynonym')}
        </button>
      </div>

      {/* Фильтры */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <div className="flex flex-wrap gap-4 items-end">
            {/* Поиск */}
            <div className="form-control flex-1 min-w-[200px]">
              <label className="label">
                <span className="label-text">{t('searchLabel')}</span>
              </label>
              <input
                type="text"
                placeholder={t('searchPlaceholder')}
                className="input input-bordered"
                value={searchTerm}
                onChange={handleSearchChange}
              />
            </div>

            {/* Выбор языка */}
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('language')}</span>
              </label>
              <div className="tabs tabs-boxed">
                {(['ru', 'en', 'sr'] as const).map((lang) => (
                  <button
                    key={lang}
                    className={`tab ${selectedLanguage === lang ? 'tab-active' : ''}`}
                    onClick={() => handleLanguageChange(lang)}
                  >
                    {getLanguageLabel(lang)}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Форма создания */}
      {showCreateForm && (
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body">
            <h3 className="card-title">{t('addSynonym')}</h3>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('term')}</span>
                </label>
                <input
                  type="text"
                  className="input input-bordered"
                  value={editForm.term}
                  onChange={(e) =>
                    setEditForm({ ...editForm, term: e.target.value })
                  }
                  placeholder={t('termPlaceholder')}
                />
              </div>
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('synonym')}</span>
                </label>
                <input
                  type="text"
                  className="input input-bordered"
                  value={editForm.synonym}
                  onChange={(e) =>
                    setEditForm({ ...editForm, synonym: e.target.value })
                  }
                  placeholder={t('synonymPlaceholder')}
                />
              </div>
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('language')}</span>
                </label>
                <select
                  className="select select-bordered"
                  value={editForm.language}
                  onChange={(e) =>
                    setEditForm({
                      ...editForm,
                      language: e.target.value as 'en' | 'ru' | 'sr',
                    })
                  }
                >
                  <option value="ru">Русский</option>
                  <option value="en">English</option>
                  <option value="sr">Српски</option>
                </select>
              </div>
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('active')}</span>
                </label>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={editForm.is_active}
                  onChange={(e) =>
                    setEditForm({ ...editForm, is_active: e.target.checked })
                  }
                />
              </div>
            </div>
            <div className="card-actions justify-end mt-4">
              <button
                className="btn btn-ghost"
                onClick={() => setShowCreateForm(false)}
              >
                {t('cancel')}
              </button>
              <button className="btn btn-primary" onClick={handleCreate}>
                {t('create')}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Таблица синонимов */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h3 className="card-title">{t('synonymsList')}</h3>
            <div className="text-sm text-base-content/70">
              {t('total')}: {total}
            </div>
          </div>

          {loading ? (
            <div className="flex justify-center py-8">
              <span className="loading loading-spinner loading-lg"></span>
            </div>
          ) : synonyms.length === 0 ? (
            <div className="text-center py-8 text-base-content/50">
              {t('noSynonyms')}
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="table table-zebra w-full">
                <thead>
                  <tr>
                    <th>{t('term')}</th>
                    <th>{t('synonym')}</th>
                    <th>{t('language')}</th>
                    <th>{t('statusLabel')}</th>
                    <th>{t('created')}</th>
                    <th>{t('actions')}</th>
                  </tr>
                </thead>
                <tbody>
                  {synonyms.map((synonym) => (
                    <tr key={synonym.id}>
                      <td>
                        {isEditing === synonym.id ? (
                          <input
                            type="text"
                            className="input input-sm input-bordered"
                            value={synonym.term}
                            onChange={(e) => {
                              const updatedSynonyms = synonyms.map((s) =>
                                s.id === synonym.id
                                  ? { ...s, term: e.target.value }
                                  : s
                              );
                              setSynonyms(updatedSynonyms);
                            }}
                          />
                        ) : (
                          <span className="font-medium">{synonym.term}</span>
                        )}
                      </td>
                      <td>
                        {isEditing === synonym.id ? (
                          <input
                            type="text"
                            className="input input-sm input-bordered"
                            value={synonym.synonym}
                            onChange={(e) => {
                              const updatedSynonyms = synonyms.map((s) =>
                                s.id === synonym.id
                                  ? { ...s, synonym: e.target.value }
                                  : s
                              );
                              setSynonyms(updatedSynonyms);
                            }}
                          />
                        ) : (
                          synonym.synonym
                        )}
                      </td>
                      <td>
                        <span className="badge badge-outline">
                          {getLanguageLabel(synonym.language)}
                        </span>
                      </td>
                      <td>
                        {isEditing === synonym.id ? (
                          <input
                            type="checkbox"
                            className="toggle toggle-sm toggle-primary"
                            checked={synonym.is_active}
                            onChange={(e) => {
                              const updatedSynonyms = synonyms.map((s) =>
                                s.id === synonym.id
                                  ? { ...s, is_active: e.target.checked }
                                  : s
                              );
                              setSynonyms(updatedSynonyms);
                            }}
                          />
                        ) : (
                          <span
                            className={`badge ${synonym.is_active ? 'badge-success' : 'badge-error'}`}
                          >
                            {synonym.is_active ? t('active') : t('inactive')}
                          </span>
                        )}
                      </td>
                      <td>
                        <span className="text-sm text-base-content/70">
                          {new Date(synonym.created_at).toLocaleDateString()}
                        </span>
                      </td>
                      <td>
                        <div className="flex gap-2">
                          {isEditing === synonym.id ? (
                            <>
                              <button
                                className="btn btn-sm btn-primary"
                                onClick={() => handleUpdate(synonym)}
                              >
                                {t('save')}
                              </button>
                              <button
                                className="btn btn-sm btn-ghost"
                                onClick={() => {
                                  setIsEditing(null);
                                  fetchSynonyms(); // Перезагрузить данные
                                }}
                              >
                                {t('cancel')}
                              </button>
                            </>
                          ) : (
                            <>
                              <button
                                className="btn btn-sm btn-ghost"
                                onClick={() => setIsEditing(synonym.id)}
                              >
                                {t('edit')}
                              </button>
                              <button
                                className="btn btn-sm btn-error"
                                onClick={() => handleDelete(synonym.id)}
                              >
                                {t('delete')}
                              </button>
                            </>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {/* Пагинация */}
          {totalPages > 1 && (
            <div className="flex justify-center mt-6">
              <div className="join">
                <button
                  className="join-item btn"
                  disabled={page <= 1}
                  onClick={() => setPage(page - 1)}
                >
                  «
                </button>
                <button className="join-item btn btn-active">
                  {page} / {totalPages}
                </button>
                <button
                  className="join-item btn"
                  disabled={page >= totalPages}
                  onClick={() => setPage(page + 1)}
                >
                  »
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
