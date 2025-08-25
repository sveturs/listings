'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { tokenManager } from '@/utils/tokenManager';
import config from '@/config';

interface IndexInfo {
  index_name: string;
  document_count: number;
  size_in_bytes: number;
  size_formatted: string;
  created_at?: string;
  last_updated?: string;
  health: string;
  status: string;
  number_of_shards: number;
  mappings: Record<string, any>;
  settings: Record<string, any>;
}

interface IndexStatistics {
  total_documents: number;
  listings_count: number;
  products_count: number;
  last_reindexed?: string;
  documents_by_category: Record<string, number>;
  documents_by_status: Record<string, number>;
  index_health: string;
  searchable_fields: string[];
}

interface IndexedDocument {
  id: string;
  type: string;
  title: string;
  category_id: number;
  category_name: string;
  user_id: number;
  storefront_id?: number;
  indexed_at: string;
  last_modified: string;
  status: string;
  searchable_fields: Record<string, any>;
}

export default function SearchIndexManager() {
  const t = useTranslations('admin');

  const [activeTab, setActiveTab] = useState<
    'overview' | 'documents' | 'mappings'
  >('overview');
  const [indexInfo, setIndexInfo] = useState<IndexInfo | null>(null);
  const [statistics, setStatistics] = useState<IndexStatistics | null>(null);
  const [documents, setDocuments] = useState<IndexedDocument[]>([]);
  const [documentsTotal, setDocumentsTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [reindexing, setReindexing] = useState(false);

  // Параметры поиска документов
  const [searchQuery, setSearchQuery] = useState('');
  const [docType, setDocType] = useState('');
  const [categoryId, setCategoryId] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);

  // Загрузка информации об индексе
  const fetchIndexInfo = async () => {
    try {
      const token = tokenManager.getAccessToken();
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/admin/search/index/info`,
        {
          headers,
          credentials: 'include',
        }
      );

      if (response.ok) {
        const data = await response.json();
        setIndexInfo(data.data);
      }
    } catch (error) {
      console.error('Failed to fetch index info:', error);
    }
  };

  // Загрузка статистики
  const fetchStatistics = async () => {
    try {
      const token = tokenManager.getAccessToken();
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/admin/search/index/statistics`,
        {
          headers,
          credentials: 'include',
        }
      );

      if (response.ok) {
        const data = await response.json();
        setStatistics(data.data);
      }
    } catch (error) {
      console.error('Failed to fetch statistics:', error);
    }
  };

  // Поиск документов
  const searchDocuments = async () => {
    setLoading(true);
    try {
      const token = tokenManager.getAccessToken();
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const params = new URLSearchParams({
        page: currentPage.toString(),
        limit: pageSize.toString(),
      });

      if (searchQuery) params.append('query', searchQuery);
      if (docType) params.append('type', docType);
      if (categoryId) params.append('category_id', categoryId);

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/admin/search/index/documents?${params}`,
        {
          headers,
          credentials: 'include',
        }
      );

      if (response.ok) {
        const data = await response.json();
        // API возвращает массив документов прямо в data.data
        setDocuments(data.data || []);
        setDocumentsTotal(data.data?.length || 0);
      }
    } catch (error) {
      console.error('Failed to search documents:', error);
    } finally {
      setLoading(false);
    }
  };

  // Запуск переиндексации
  const handleReindex = async () => {
    if (!confirm(t('search.index.confirmReindex'))) {
      return;
    }

    setReindexing(true);
    try {
      const token = tokenManager.getAccessToken();
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/admin/search/index/reindex`,
        {
          method: 'POST',
          headers,
          credentials: 'include',
        }
      );

      if (response.ok) {
        alert(t('search.index.reindexStarted'));
        // Обновляем статистику через 5 секунд
        setTimeout(() => {
          fetchStatistics();
        }, 5000);
      }
    } catch (error) {
      console.error('Failed to start reindexing:', error);
      alert(t('search.index.reindexError'));
    } finally {
      setReindexing(false);
    }
  };

  useEffect(() => {
    fetchIndexInfo();
    fetchStatistics();
  }, []);

  useEffect(() => {
    if (activeTab === 'documents') {
      searchDocuments();
    }
  }, [activeTab, currentPage, searchQuery, docType, categoryId]);

  const getHealthBadgeClass = (health: string) => {
    switch (health) {
      case 'green':
        return 'badge-success';
      case 'yellow':
        return 'badge-warning';
      case 'red':
        return 'badge-error';
      default:
        return 'badge-ghost';
    }
  };

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'active':
        return 'badge-success';
      case 'inactive':
        return 'badge-error';
      case 'pending':
        return 'badge-warning';
      default:
        return 'badge-ghost';
    }
  };

  return (
    <div className="space-y-6">
      {/* Заголовок и действия */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold">{t('search.index.title')}</h2>
          <p className="text-base-content/70">
            {t('search.index.description')}
          </p>
        </div>
        <button
          className={`btn btn-primary ${reindexing ? 'loading' : ''}`}
          onClick={handleReindex}
          disabled={reindexing}
        >
          {reindexing
            ? t('search.index.reindexing')
            : t('search.index.reindex')}
        </button>
      </div>

      {/* Табы */}
      <div className="tabs tabs-boxed">
        <button
          className={`tab ${activeTab === 'overview' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('overview')}
        >
          {t('search.index.tabs.overview')}
        </button>
        <button
          className={`tab ${activeTab === 'documents' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('documents')}
        >
          {t('search.index.tabs.documents')}
        </button>
        <button
          className={`tab ${activeTab === 'mappings' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('mappings')}
        >
          {t('search.index.tabs.mappings')}
        </button>
      </div>

      {/* Контент */}
      {activeTab === 'overview' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Основная информация */}
          <div className="card bg-base-100">
            <div className="card-body">
              <h3 className="card-title">{t('search.index.info')}</h3>
              {indexInfo && (
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-base-content/70">
                      {t('search.index.name')}
                    </span>
                    <span className="font-mono">{indexInfo.index_name}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">
                      {t('search.index.documents')}
                    </span>
                    <span className="font-bold">
                      {indexInfo.document_count.toLocaleString()}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">
                      {t('search.index.size')}
                    </span>
                    <span>{indexInfo.size_formatted}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">
                      {t('search.index.health')}
                    </span>
                    <span
                      className={`badge ${getHealthBadgeClass(indexInfo.health)}`}
                    >
                      {indexInfo.health}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">
                      {t('search.index.status')}
                    </span>
                    <span className="badge badge-success">
                      {indexInfo.status}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">
                      {t('search.index.shards')}
                    </span>
                    <span>{indexInfo.number_of_shards}</span>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Статистика */}
          <div className="card bg-base-100">
            <div className="card-body">
              <h3 className="card-title">{t('search.index.statistics')}</h3>
              {statistics && (
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-base-content/70">
                      {t('search.index.totalDocs')}
                    </span>
                    <span className="font-bold">
                      {statistics.total_documents.toLocaleString()}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">
                      {t('search.index.listings')}
                    </span>
                    <span>{statistics.listings_count.toLocaleString()}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">
                      {t('search.index.products')}
                    </span>
                    <span>{statistics.products_count.toLocaleString()}</span>
                  </div>
                  {statistics.last_reindexed && (
                    <div className="flex justify-between">
                      <span className="text-base-content/70">
                        {t('search.index.lastReindexed')}
                      </span>
                      <span className="text-sm">
                        {new Date(statistics.last_reindexed).toLocaleString()}
                      </span>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>

          {/* Документы по категориям */}
          {statistics &&
            Object.keys(statistics.documents_by_category).length > 0 && (
              <div className="card bg-base-100 lg:col-span-2">
                <div className="card-body">
                  <h3 className="card-title">{t('search.index.byCategory')}</h3>
                  <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
                    {Object.entries(statistics.documents_by_category)
                      .sort(([, a], [, b]) => b - a)
                      .slice(0, 12)
                      .map(([category, count]) => (
                        <div
                          key={category}
                          className="flex justify-between p-2 bg-base-200 rounded"
                        >
                          <span className="text-sm truncate">{category}</span>
                          <span className="badge badge-sm">{count}</span>
                        </div>
                      ))}
                  </div>
                </div>
              </div>
            )}

          {/* Поля для поиска */}
          {statistics && (
            <div className="card bg-base-100 lg:col-span-2">
              <div className="card-body">
                <h3 className="card-title">
                  {t('search.index.searchableFields')}
                </h3>
                <div className="flex flex-wrap gap-2">
                  {statistics.searchable_fields.map((field) => (
                    <span key={field} className="badge badge-outline">
                      {field}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {activeTab === 'documents' && (
        <div className="space-y-4">
          {/* Фильтры */}
          <div className="card bg-base-100">
            <div className="card-body">
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <input
                  type="text"
                  placeholder={t('search.index.searchPlaceholder')}
                  className="input input-bordered"
                  value={searchQuery}
                  onChange={(e) => {
                    setSearchQuery(e.target.value);
                    setCurrentPage(1);
                  }}
                />
                <select
                  className="select select-bordered"
                  value={docType}
                  onChange={(e) => {
                    setDocType(e.target.value);
                    setCurrentPage(1);
                  }}
                >
                  <option value="">{t('search.index.allTypes')}</option>
                  <option value="listing">
                    {t('search.index.typeListings')}
                  </option>
                  <option value="product">
                    {t('search.index.typeProducts')}
                  </option>
                </select>
                <input
                  type="number"
                  placeholder={t('search.index.categoryId')}
                  className="input input-bordered"
                  value={categoryId}
                  onChange={(e) => {
                    setCategoryId(e.target.value);
                    setCurrentPage(1);
                  }}
                />
                <button
                  className={`btn btn-primary ${loading ? 'loading' : ''}`}
                  onClick={() => searchDocuments()}
                  disabled={loading}
                >
                  {t('search.index.search')}
                </button>
              </div>
            </div>
          </div>

          {/* Список документов */}
          <div className="card bg-base-100">
            <div className="card-body">
              {loading ? (
                <div className="flex justify-center py-8">
                  <span className="loading loading-spinner loading-lg"></span>
                </div>
              ) : documents.length > 0 ? (
                <div className="overflow-x-auto">
                  <table className="table table-sm">
                    <thead>
                      <tr>
                        <th>{t('search.index.docId')}</th>
                        <th>{t('search.index.docType')}</th>
                        <th>{t('search.index.docTitle')}</th>
                        <th>{t('search.index.docCategory')}</th>
                        <th>{t('search.index.docStatus')}</th>
                        <th>{t('search.index.docIndexedAt')}</th>
                        <th>{t('search.index.docSearchFields')}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {documents.map((doc) => (
                        <tr key={doc.id}>
                          <td className="font-mono text-xs">{doc.id}</td>
                          <td>
                            <span className="badge badge-sm">{doc.type}</span>
                          </td>
                          <td className="max-w-xs truncate">{doc.title}</td>
                          <td>{doc.category_name}</td>
                          <td>
                            <span
                              className={`badge badge-sm ${getStatusBadgeClass(doc.status)}`}
                            >
                              {doc.status}
                            </span>
                          </td>
                          <td className="text-xs">
                            {new Date(doc.indexed_at).toLocaleString()}
                          </td>
                          <td>
                            <button
                              className="btn btn-xs btn-ghost"
                              onClick={() => {
                                alert(
                                  JSON.stringify(doc.searchable_fields, null, 2)
                                );
                              }}
                            >
                              {t('search.index.viewFields')}
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className="text-center py-8 text-base-content/50">
                  {t('search.index.noDocuments')}
                </div>
              )}

              {/* Пагинация */}
              {documentsTotal > pageSize && (
                <div className="flex justify-center mt-4">
                  <div className="join">
                    <button
                      className="join-item btn"
                      disabled={currentPage === 1}
                      onClick={() => setCurrentPage(currentPage - 1)}
                    >
                      «
                    </button>
                    <button className="join-item btn">
                      {t('search.index.pageInfo', {
                        current: currentPage,
                        total: Math.ceil(documentsTotal / pageSize),
                      })}
                    </button>
                    <button
                      className="join-item btn"
                      disabled={
                        currentPage >= Math.ceil(documentsTotal / pageSize)
                      }
                      onClick={() => setCurrentPage(currentPage + 1)}
                    >
                      »
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {activeTab === 'mappings' && indexInfo && (
        <div className="card bg-base-100">
          <div className="card-body">
            <h3 className="card-title">{t('search.index.mappingsTitle')}</h3>
            <div className="mockup-code">
              <pre className="text-xs">
                <code>{JSON.stringify(indexInfo.mappings, null, 2)}</code>
              </pre>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
