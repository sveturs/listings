'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { adminApi, Attribute } from '@/services/admin';
import { toast } from '@/utils/toast';
import AttributeForm from './components/AttributeForm';

export default function AttributesPage() {
  const t = useTranslations('admin');
  const [attributes, setAttributes] = useState<Attribute[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedAttribute, setSelectedAttribute] = useState<Attribute | null>(
    null
  );
  const [showForm, setShowForm] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState('');
  const [isInitialized, setIsInitialized] = useState(false);

  useEffect(() => {
    // Ждем инициализации авторизации
    const initAuth = async () => {
      try {
        // Проверяем, есть ли токен
        const { tokenManager } = await import('@/utils/tokenManager');

        // Даем время на обновление токена если нужно
        const token = await tokenManager.getAccessToken();
        if (!token) {
          // Попробуем обновить токен
          try {
            await tokenManager.refreshAccessToken();
          } catch (error) {
            console.log('Failed to refresh token:', error);
          }
        }

        setIsInitialized(true);
      } catch (error) {
        console.error('Auth initialization error:', error);
        setIsInitialized(true); // Все равно пытаемся загрузить
      }
    };

    initAuth();
  }, []);

  useEffect(() => {
    if (isInitialized) {
      loadAttributes();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isInitialized]);

  const loadAttributes = async () => {
    try {
      setLoading(true);
      const data = await adminApi.attributes.getAll();
      setAttributes(data);
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to load attributes:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddAttribute = () => {
    setSelectedAttribute(null);
    setIsEditing(false);
    setShowForm(true);
  };

  const handleEditAttribute = (attribute: Attribute) => {
    setSelectedAttribute(attribute);
    setIsEditing(true);
    setShowForm(true);
  };

  const handleDeleteAttribute = async (attribute: Attribute) => {
    if (!confirm(t('common.confirmDelete'))) return;

    try {
      await adminApi.attributes.delete(attribute.id);
      toast.success(t('common.deleteSuccess'));
      await loadAttributes();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to delete attribute:', error);
    }
  };

  const handleSaveAttribute = async (data: Partial<Attribute>) => {
    try {
      if (isEditing && selectedAttribute) {
        await adminApi.attributes.update(selectedAttribute.id, data);
        toast.success(t('common.saveSuccess'));
      } else {
        await adminApi.attributes.create(data);
        toast.success(t('common.saveSuccess'));
      }
      setShowForm(false);
      await loadAttributes();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to save attribute:', error);
    }
  };

  // Filter attributes based on search and type
  const filteredAttributes = attributes.filter((attr) => {
    const matchesSearch =
      attr.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      attr.display_name.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesType = !filterType || attr.attribute_type === filterType;
    return matchesSearch && matchesType;
  });

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">{t('attributes.title')}</h1>
        <button className="btn btn-primary" onClick={handleAddAttribute}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5 mr-2"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 4v16m8-8H4"
            />
          </svg>
          {t('attributes.addAttribute')}
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className={showForm ? 'lg:col-span-2' : 'lg:col-span-3'}>
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              {/* Filters */}
              <div className="flex gap-4 mb-4">
                <div className="form-control flex-1">
                  <input
                    type="text"
                    placeholder={t('common.search')}
                    className="input input-bordered"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                  />
                </div>
                <div className="form-control">
                  <select
                    className="select select-bordered"
                    value={filterType}
                    onChange={(e) => setFilterType(e.target.value)}
                  >
                    <option value="">{t('attributes.allTypes')}</option>
                    <option value="text">{t('attributes.types.text')}</option>
                    <option value="number">
                      {t('attributes.types.number')}
                    </option>
                    <option value="select">
                      {t('attributes.types.select')}
                    </option>
                    <option value="boolean">
                      {t('attributes.types.boolean')}
                    </option>
                    <option value="date">{t('attributes.types.date')}</option>
                    <option value="range">{t('attributes.types.range')}</option>
                    <option value="location">
                      {t('attributes.types.location')}
                    </option>
                    <option value="file">{t('attributes.types.file')}</option>
                    <option value="gallery">
                      {t('attributes.types.gallery')}
                    </option>
                  </select>
                </div>
              </div>

              {/* Attributes Table */}
              <div className="overflow-x-auto">
                <table className="table table-zebra">
                  <thead>
                    <tr>
                      <th>{t('attributes.systemName')}</th>
                      <th>{t('attributes.displayName')}</th>
                      <th>{t('attributes.type')}</th>
                      <th>{t('attributes.usedInCategories')}</th>
                      <th className="text-center">{t('common.actions')}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {filteredAttributes.length === 0 ? (
                      <tr>
                        <td colSpan={5} className="text-center">
                          {t('common.noData')}
                        </td>
                      </tr>
                    ) : (
                      filteredAttributes.map((attr) => (
                        <tr key={attr.id}>
                          <td>
                            <code className="text-sm">{attr.name}</code>
                          </td>
                          <td>{attr.display_name}</td>
                          <td>
                            <span className="badge badge-outline">
                              {t(`attributes.types.${attr.attribute_type}`)}
                            </span>
                          </td>
                          <td>
                            <div className="flex gap-1">
                              {attr.is_searchable && (
                                <span
                                  className="badge badge-sm badge-info"
                                  title={t('attributes.isSearchable')}
                                >
                                  🔍
                                </span>
                              )}
                              {attr.is_filterable && (
                                <span
                                  className="badge badge-sm badge-warning"
                                  title={t('attributes.isFilterable')}
                                >
                                  🔧
                                </span>
                              )}
                              {attr.is_required && (
                                <span
                                  className="badge badge-sm badge-error"
                                  title={t('attributes.isRequired')}
                                >
                                  *
                                </span>
                              )}
                            </div>
                          </td>
                          <td className="text-center">
                            <div className="dropdown dropdown-end">
                              <label
                                tabIndex={0}
                                className="btn btn-ghost btn-xs"
                              >
                                <svg
                                  xmlns="http://www.w3.org/2000/svg"
                                  className="h-4 w-4"
                                  fill="none"
                                  viewBox="0 0 24 24"
                                  stroke="currentColor"
                                >
                                  <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z"
                                  />
                                </svg>
                              </label>
                              <ul
                                tabIndex={0}
                                className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52"
                              >
                                <li>
                                  <a onClick={() => handleEditAttribute(attr)}>
                                    <svg
                                      xmlns="http://www.w3.org/2000/svg"
                                      className="h-4 w-4"
                                      fill="none"
                                      viewBox="0 0 24 24"
                                      stroke="currentColor"
                                    >
                                      <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                                      />
                                    </svg>
                                    {t('common.edit')}
                                  </a>
                                </li>
                                <li>
                                  <a
                                    onClick={() => handleDeleteAttribute(attr)}
                                    className="text-error"
                                  >
                                    <svg
                                      xmlns="http://www.w3.org/2000/svg"
                                      className="h-4 w-4"
                                      fill="none"
                                      viewBox="0 0 24 24"
                                      stroke="currentColor"
                                    >
                                      <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                                      />
                                    </svg>
                                    {t('common.delete')}
                                  </a>
                                </li>
                              </ul>
                            </div>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        {showForm && (
          <div className="lg:col-span-1">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">
                  {isEditing
                    ? t('attributes.editAttribute')
                    : t('attributes.addAttribute')}
                </h2>
                <AttributeForm
                  attribute={selectedAttribute}
                  onSave={handleSaveAttribute}
                  onCancel={() => setShowForm(false)}
                />
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
