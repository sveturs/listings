'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import { useDebounce } from '@/hooks/useDebounce';
import { adminApi, VariantAttribute } from '@/services/admin';
import VariantAttributeList from './components/VariantAttributeList';
import VariantAttributeForm from './components/VariantAttributeForm';
import AttributeMappingEditor from './components/AttributeMappingEditor';

// VariantAttribute type imported from @/services/admin

export default function VariantAttributesPage() {
  const t = useTranslations('admin');
  const [attributes, setAttributes] = useState<VariantAttribute[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedAttribute, setSelectedAttribute] =
    useState<VariantAttribute | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState('');
  const [isInitialized, setIsInitialized] = useState(false);
  const [showMappingEditor, setShowMappingEditor] = useState(false);
  const [mappingAttribute, setMappingAttribute] =
    useState<VariantAttribute | null>(null);

  // Пагинация
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);
  const [pageSize] = useState(20);

  // Используем debounce для поиска
  const debouncedSearchTerm = useDebounce(searchTerm, 500);

  useEffect(() => {
    // Ждем инициализации авторизации
    const initAuth = async () => {
      try {
        const { tokenManager } = await import('@/utils/tokenManager');
        const token = await tokenManager.getAccessToken();
        if (!token) {
          try {
            await tokenManager.refreshAccessToken();
          } catch (error) {
            console.log('Failed to refresh token:', error);
          }
        }
        setIsInitialized(true);
      } catch (error) {
        console.error('Auth initialization error:', error);
        setIsInitialized(true);
      }
    };

    initAuth();
  }, []);

  const loadAttributes = useCallback(async () => {
    try {
      setLoading(true);

      const response = await adminApi.variantAttributes.getAll(
        currentPage,
        pageSize,
        debouncedSearchTerm,
        filterType
      );

      setAttributes(response.data);
      setTotalPages(response.total_pages || 0);
      setTotalItems(response.total || 0);
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to load variant attributes:', error);
    } finally {
      setLoading(false);
    }
  }, [currentPage, pageSize, debouncedSearchTerm, filterType, t]);

  useEffect(() => {
    if (isInitialized) {
      loadAttributes();
    }
  }, [
    isInitialized,
    currentPage,
    debouncedSearchTerm,
    filterType,
    loadAttributes,
  ]);

  // Сбрасываем на первую страницу при изменении поиска или фильтра
  useEffect(() => {
    if (isInitialized && currentPage !== 1) {
      setCurrentPage(1);
    }
  }, [debouncedSearchTerm, filterType, isInitialized, currentPage]);

  const handleAddAttribute = () => {
    setSelectedAttribute(null);
    setIsEditing(false);
    setShowForm(true);
  };

  const handleEditAttribute = async (attribute: VariantAttribute) => {
    try {
      // Загружаем полные данные атрибута
      const fullAttribute = await adminApi.variantAttributes.getById(
        attribute.id
      );
      setSelectedAttribute(fullAttribute);
      setIsEditing(true);
      setShowForm(true);
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to load variant attribute details:', error);
    }
  };

  const handleDeleteAttribute = async (attribute: VariantAttribute) => {
    if (!confirm(t('common.confirmDelete'))) return;

    try {
      await adminApi.variantAttributes.delete(attribute.id);
      toast.success(t('common.deleteSuccess'));
      await loadAttributes();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to delete variant attribute:', error);
    }
  };

  const handleSaveAttribute = async (data: Partial<VariantAttribute>) => {
    try {
      if (isEditing && selectedAttribute) {
        await adminApi.variantAttributes.update(selectedAttribute.id, data);
        toast.success(t('common.saveSuccess'));
      } else {
        await adminApi.variantAttributes.create(data);
        toast.success(t('common.saveSuccess'));
      }
      setShowForm(false);
      await loadAttributes();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to save variant attribute:', error);
    }
  };

  const handleManageLinks = (attribute: VariantAttribute) => {
    setMappingAttribute(attribute);
    setShowMappingEditor(true);
  };

  const handleCloseMappingEditor = () => {
    setShowMappingEditor(false);
    setMappingAttribute(null);
    loadAttributes(); // Перезагружаем список для обновления данных
  };

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
        <h1 className="text-2xl font-bold">{t('variantAttributes.title')}</h1>
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
          {t('variantAttributes.addAttribute')}
        </button>
      </div>

      <div className="grid grid-cols-1 gap-6">
        <div className="col-span-1">
          <VariantAttributeList
            attributes={attributes}
            searchTerm={searchTerm}
            filterType={filterType}
            currentPage={currentPage}
            totalPages={totalPages}
            totalItems={totalItems}
            pageSize={pageSize}
            onSearchChange={setSearchTerm}
            onFilterChange={setFilterType}
            onPageChange={setCurrentPage}
            onEdit={handleEditAttribute}
            onDelete={handleDeleteAttribute}
            onManageLinks={handleManageLinks}
          />
        </div>

        {/* Modal for Variant Attribute Form */}
        {showForm && (
          <div className="modal modal-open">
            <div className="modal-box w-11/12 max-w-2xl">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-bold">
                  {isEditing
                    ? t('variantAttributes.editAttribute')
                    : t('variantAttributes.addAttribute')}
                </h2>
                <button
                  className="btn btn-sm btn-circle btn-ghost"
                  onClick={() => setShowForm(false)}
                >
                  ✕
                </button>
              </div>
              <VariantAttributeForm
                attribute={selectedAttribute}
                onSave={handleSaveAttribute}
                onCancel={() => setShowForm(false)}
              />
            </div>
            <div
              className="modal-backdrop"
              onClick={() => setShowForm(false)}
            ></div>
          </div>
        )}

        {/* Modal for Attribute Mapping Editor */}
        {showMappingEditor && mappingAttribute && (
          <div className="modal modal-open">
            <div className="modal-box w-11/12 max-w-3xl">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-bold">
                  {t('variantAttributes.manageMappings')}
                </h2>
                <button
                  className="btn btn-sm btn-circle btn-ghost"
                  onClick={handleCloseMappingEditor}
                >
                  ✕
                </button>
              </div>
              <AttributeMappingEditor
                variantAttribute={mappingAttribute}
                onSave={handleCloseMappingEditor}
                onCancel={handleCloseMappingEditor}
              />
            </div>
            <div
              className="modal-backdrop"
              onClick={handleCloseMappingEditor}
            ></div>
          </div>
        )}
      </div>
    </div>
  );
}
