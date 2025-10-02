'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { adminApi, Category } from '@/services/admin';
import { toast } from '@/utils/toast';
import CategoryTree from './components/CategoryTree';
import CategoryForm from './components/CategoryForm';
import CategoryAttributesOptimized from './components/CategoryAttributesOptimized';
import CategoryKeywordsModal from './components/CategoryKeywordsModal';
import CategoryVariantAttributesModal from './components/CategoryVariantAttributesModal';

export default function CategoriesPage() {
  const t = useTranslations('admin');
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(
    null
  );
  const [showForm, setShowForm] = useState(false);
  const [showAttributes, setShowAttributes] = useState(false);
  const [showKeywordsModal, setShowKeywordsModal] = useState(false);
  const [showVariantAttributesModal, setShowVariantAttributesModal] =
    useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    loadCategories();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadCategories = async () => {
    try {
      setLoading(true);
      const data = await adminApi.categories.getAll();

      // Убеждаемся, что data это массив
      const categoriesArray = Array.isArray(data) ? data : [];
      setCategories(categoriesArray);
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to load categories:', error);
      setCategories([]); // Устанавливаем пустой массив при ошибке
    } finally {
      setLoading(false);
    }
  };

  const handleAddCategory = () => {
    setSelectedCategory(null);
    setIsEditing(false);
    setShowForm(true);
    setShowAttributes(false);
  };

  const handleEditCategory = (category: Category) => {
    setSelectedCategory(category);
    setIsEditing(true);
    setShowForm(true);
    setShowAttributes(false);
  };

  const handleManageAttributes = (category: Category) => {
    setSelectedCategory(category);
    setShowForm(false);
    setShowAttributes(true);
  };

  const handleManageKeywords = (category: Category) => {
    setSelectedCategory(category);
    setShowKeywordsModal(true);
  };

  const handleManageVariantAttributes = (category: Category) => {
    setSelectedCategory(category);
    setShowVariantAttributesModal(true);
  };

  const handleDeleteCategory = async (category: Category) => {
    if (!confirm(t('common.confirmDelete'))) return;

    try {
      await adminApi.categories.delete(category.id);
      toast.success(t('common.deleteSuccess'));
      await loadCategories();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to delete category:', error);
    }
  };

  const handleSaveCategory = async (data: Partial<Category>) => {
    try {
      if (isEditing && selectedCategory) {
        await adminApi.categories.update(selectedCategory.id, data);
        toast.success(t('common.saveSuccess'));
      } else {
        await adminApi.categories.create(data);
        toast.success(t('common.saveSuccess'));
      }
      setShowForm(false);
      await loadCategories();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to save category:', error);
    }
  };

  const handleReorderCategories = async (orderedIds: number[]) => {
    try {
      await adminApi.categories.reorder(orderedIds);
      await loadCategories();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to reorder categories:', error);
    }
  };

  const handleMoveCategory = async (
    categoryId: number,
    newParentId: number
  ) => {
    try {
      await adminApi.categories.move(categoryId, newParentId);
      await loadCategories();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to move category:', error);
    }
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
        <h1 className="text-2xl font-bold">{t('categories.title')}</h1>
        <button className="btn btn-primary" onClick={handleAddCategory}>
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
          {t('categories.addCategory')}
        </button>
      </div>

      <div className="grid grid-cols-1 gap-6">
        <div className="col-span-1">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <CategoryTree
                categories={categories}
                onEdit={handleEditCategory}
                onDelete={handleDeleteCategory}
                onManageAttributes={handleManageAttributes}
                onManageKeywords={handleManageKeywords}
                onManageVariantAttributes={handleManageVariantAttributes}
                onReorder={handleReorderCategories}
                onMove={handleMoveCategory}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Modal for Category Form */}
      {showForm && (
        <div className="modal modal-open">
          <div className="modal-box w-11/12 max-w-4xl">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold">
                {isEditing
                  ? t('categories.editCategory')
                  : t('categories.addCategory')}
              </h2>
              <button
                className="btn btn-sm btn-circle btn-ghost"
                onClick={() => setShowForm(false)}
              >
                ✕
              </button>
            </div>
            <CategoryForm
              category={selectedCategory}
              categories={categories}
              onSave={handleSaveCategory}
              onCancel={() => setShowForm(false)}
            />
          </div>
          <div
            className="modal-backdrop"
            onClick={() => setShowForm(false)}
          ></div>
        </div>
      )}

      {/* Modal for Category Attributes */}
      {showAttributes && selectedCategory && (
        <div className="modal modal-open">
          <div className="modal-box w-11/12 max-w-6xl">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold">
                {t('sections.attributes')}: {selectedCategory.name}
              </h2>
              <button
                className="btn btn-sm btn-circle btn-ghost"
                onClick={() => setShowAttributes(false)}
              >
                ✕
              </button>
            </div>
            <CategoryAttributesOptimized
              category={selectedCategory}
              onUpdate={loadCategories}
            />
          </div>
          <div
            className="modal-backdrop"
            onClick={() => setShowAttributes(false)}
          ></div>
        </div>
      )}

      {/* Modal for managing keywords */}
      {selectedCategory && (
        <CategoryKeywordsModal
          categoryId={selectedCategory.id}
          categoryName={selectedCategory.name}
          isOpen={showKeywordsModal}
          onClose={() => setShowKeywordsModal(false)}
        />
      )}

      {/* Modal for managing variant attributes */}
      {selectedCategory && (
        <CategoryVariantAttributesModal
          categoryId={selectedCategory.id}
          categoryName={selectedCategory.name}
          isOpen={showVariantAttributesModal}
          onClose={() => setShowVariantAttributesModal(false)}
        />
      )}
    </div>
  );
}
