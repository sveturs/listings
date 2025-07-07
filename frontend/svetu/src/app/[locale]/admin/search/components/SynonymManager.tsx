'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';
import { tokenManager } from '@/utils/tokenManager';

interface Synonym {
  id: string;
  word: string;
  synonyms: string[];
  language: 'en' | 'ru' | 'sr';
  active: boolean;
}

export default function SynonymManager() {
  const t = useTranslations();
  const [synonyms, setSynonyms] = useState<Synonym[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [newWord, setNewWord] = useState('');
  const [newSynonyms, setNewSynonyms] = useState('');
  const [selectedLanguage, setSelectedLanguage] = useState<'en' | 'ru' | 'sr'>(
    'ru'
  );
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    fetchSynonyms();
  }, [selectedLanguage]);

  const fetchSynonyms = async () => {
    try {
      const accessToken = await tokenManager.getAccessToken();
      const response = await fetch(
        `/api/admin/search/config/synonyms?lang=${selectedLanguage}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        }
      );
      if (!response.ok) throw new Error('Failed to fetch synonyms');
      const data = await response.json();
      setSynonyms(data.synonyms || []);
    } catch (error) {
      console.error('Error fetching synonyms:', error);
      toast.error(t('admin.search.synonyms.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = async () => {
    if (!newWord.trim() || !newSynonyms.trim()) {
      toast.error(t('admin.search.synonyms.fillAllFields'));
      return;
    }

    try {
      const accessToken = await tokenManager.getAccessToken();
      const response = await fetch('/api/admin/search/synonyms', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          word: newWord.trim(),
          synonyms: newSynonyms
            .split(',')
            .map((s) => s.trim())
            .filter(Boolean),
          language: selectedLanguage,
          active: true,
        }),
      });

      if (!response.ok) throw new Error('Failed to add synonym');

      setNewWord('');
      setNewSynonyms('');
      fetchSynonyms();
      toast.success(t('admin.search.synonyms.addSuccess'));
    } catch (error) {
      console.error('Error adding synonym:', error);
      toast.error(t('admin.search.synonyms.addError'));
    }
  };

  const handleUpdate = async (synonym: Synonym) => {
    try {
      const accessToken = await tokenManager.getAccessToken();
      const response = await fetch(`/api/admin/search/synonyms/${synonym.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify(synonym),
      });

      if (!response.ok) throw new Error('Failed to update synonym');

      setEditingId(null);
      fetchSynonyms();
      toast.success(t('admin.search.synonyms.updateSuccess'));
    } catch (error) {
      console.error('Error updating synonym:', error);
      toast.error(t('admin.search.synonyms.updateError'));
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm(t('admin.search.synonyms.deleteConfirm'))) return;

    try {
      const accessToken = await tokenManager.getAccessToken();
      const response = await fetch(`/api/admin/search/synonyms/${id}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });

      if (!response.ok) throw new Error('Failed to delete synonym');

      fetchSynonyms();
      toast.success(t('admin.search.synonyms.deleteSuccess'));
    } catch (error) {
      console.error('Error deleting synonym:', error);
      toast.error(t('admin.search.synonyms.deleteError'));
    }
  };

  const toggleActive = async (synonym: Synonym) => {
    await handleUpdate({ ...synonym, active: !synonym.active });
  };

  const filteredSynonyms = synonyms.filter(
    (s) =>
      s.word.toLowerCase().includes(searchTerm.toLowerCase()) ||
      s.synonyms.some((syn) =>
        syn.toLowerCase().includes(searchTerm.toLowerCase())
      )
  );

  if (loading) {
    return <div className="loading loading-spinner loading-lg"></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex gap-4 items-end">
        <div className="form-control">
          <label className="label">
            <span className="label-text">
              {t('admin.search.synonyms.language')}
            </span>
          </label>
          <select
            className="select select-bordered"
            value={selectedLanguage}
            onChange={(e) =>
              setSelectedLanguage(e.target.value as 'en' | 'ru' | 'sr')
            }
          >
            <option value="ru">Русский</option>
            <option value="en">English</option>
            <option value="sr">Српски</option>
          </select>
        </div>

        <div className="form-control flex-1">
          <label className="label">
            <span className="label-text">
              {t('admin.search.synonyms.search')}
            </span>
          </label>
          <input
            type="text"
            className="input input-bordered"
            placeholder={t('admin.search.synonyms.searchPlaceholder')}
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
      </div>

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">{t('admin.search.synonyms.addNew')}</h3>
          <div className="flex gap-4 items-end">
            <div className="form-control flex-1">
              <label className="label">
                <span className="label-text">
                  {t('admin.search.synonyms.word')}
                </span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                placeholder={t('admin.search.synonyms.wordPlaceholder')}
                value={newWord}
                onChange={(e) => setNewWord(e.target.value)}
              />
            </div>
            <div className="form-control flex-1">
              <label className="label">
                <span className="label-text">
                  {t('admin.search.synonyms.synonymsList')}
                </span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                placeholder={t('admin.search.synonyms.synonymsPlaceholder')}
                value={newSynonyms}
                onChange={(e) => setNewSynonyms(e.target.value)}
              />
            </div>
            <button className="btn btn-primary" onClick={handleAdd}>
              {t('admin.search.synonyms.add')}
            </button>
          </div>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="table table-zebra">
          <thead>
            <tr>
              <th>{t('admin.search.synonyms.word')}</th>
              <th>{t('admin.search.synonyms.synonymsList')}</th>
              <th>{t('admin.search.synonyms.status')}</th>
              <th>{t('admin.search.synonyms.actions')}</th>
            </tr>
          </thead>
          <tbody>
            {filteredSynonyms.map((synonym) => (
              <tr key={synonym.id}>
                <td>
                  {editingId === synonym.id ? (
                    <input
                      type="text"
                      className="input input-bordered input-sm"
                      value={synonym.word}
                      onChange={(e) => {
                        const updated = synonyms.map((s) =>
                          s.id === synonym.id
                            ? { ...s, word: e.target.value }
                            : s
                        );
                        setSynonyms(updated);
                      }}
                    />
                  ) : (
                    <span className="font-mono">{synonym.word}</span>
                  )}
                </td>
                <td>
                  {editingId === synonym.id ? (
                    <input
                      type="text"
                      className="input input-bordered input-sm w-full"
                      value={synonym.synonyms.join(', ')}
                      onChange={(e) => {
                        const updated = synonyms.map((s) =>
                          s.id === synonym.id
                            ? {
                                ...s,
                                synonyms: e.target.value
                                  .split(',')
                                  .map((s) => s.trim())
                                  .filter(Boolean),
                              }
                            : s
                        );
                        setSynonyms(updated);
                      }}
                    />
                  ) : (
                    <div className="flex flex-wrap gap-1">
                      {synonym.synonyms.map((syn, idx) => (
                        <span key={idx} className="badge badge-outline">
                          {syn}
                        </span>
                      ))}
                    </div>
                  )}
                </td>
                <td>
                  <input
                    type="checkbox"
                    className="toggle toggle-success"
                    checked={synonym.active}
                    onChange={() => toggleActive(synonym)}
                    disabled={editingId === synonym.id}
                  />
                </td>
                <td>
                  <div className="flex gap-2">
                    {editingId === synonym.id ? (
                      <>
                        <button
                          className="btn btn-success btn-xs"
                          onClick={() => handleUpdate(synonym)}
                        >
                          {t('admin.search.synonyms.save')}
                        </button>
                        <button
                          className="btn btn-ghost btn-xs"
                          onClick={() => {
                            setEditingId(null);
                            fetchSynonyms();
                          }}
                        >
                          {t('admin.search.synonyms.cancel')}
                        </button>
                      </>
                    ) : (
                      <>
                        <button
                          className="btn btn-ghost btn-xs"
                          onClick={() => setEditingId(synonym.id)}
                        >
                          {t('admin.search.synonyms.edit')}
                        </button>
                        <button
                          className="btn btn-error btn-xs"
                          onClick={() => handleDelete(synonym.id)}
                        >
                          {t('admin.search.synonyms.delete')}
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
    </div>
  );
}
