'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { adminApi, Category } from '@/services/admin';
import { toast } from '@/utils/toast';
import CategoryTree from './components/CategoryTree';
import CategoryForm from './components/CategoryForm';
import CategoryAttributes from './components/CategoryAttributes';

export default function CategoriesPage() {
  const t = useTranslations('admin');
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(
    null
  );
  const [showForm, setShowForm] = useState(false);
  const [showAttributes, setShowAttributes] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    loadCategories();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadCategories = async () => {
    console.log('=== loadCategories called ===');
    try {
      console.log('Setting loading to true');
      setLoading(true);
      console.log('Loading categories...');
      console.log('About to call adminApi.categories.getAll()');

      let data;
      try {
        data = await adminApi.categories.getAll();
        console.log('adminApi.categories.getAll() completed successfully');
      } catch (apiError) {
        console.error('ERROR in adminApi.categories.getAll():', apiError);
        throw apiError;
      }

      console.log('Raw API response:', data);
      console.log('Response type:', typeof data);
      console.log('Is array:', Array.isArray(data));
      console.log('Response keys:', Object.keys(data || {}));

      // Убеждаемся, что data это массив
      const categoriesArray = Array.isArray(data) ? data : [];
      console.log('Final categories array:', categoriesArray);
      console.log('Categories count:', categoriesArray.length);
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

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <CategoryTree
                categories={categories}
                onEdit={handleEditCategory}
                onDelete={handleDeleteCategory}
                onManageAttributes={handleManageAttributes}
                onReorder={handleReorderCategories}
                onMove={handleMoveCategory}
              />
            </div>
          </div>
        </div>

        <div className="lg:col-span-1">
          {showForm && (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">
                  {isEditing
                    ? t('categories.editCategory')
                    : t('categories.addCategory')}
                </h2>
                <CategoryForm
                  category={selectedCategory}
                  categories={categories}
                  onSave={handleSaveCategory}
                  onCancel={() => setShowForm(false)}
                />
              </div>
            </div>
          )}

          {showAttributes && selectedCategory && (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">
                  {t('sections.attributes')}: {selectedCategory.name}
                </h2>
                <CategoryAttributes
                  category={selectedCategory}
                  onUpdate={loadCategories}
                />
              </div>
            </div>
          )}

          {!showForm && !showAttributes && (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">{t('categories.title')}</h2>
                <p className="text-base-content/60">
                  {t('categories.description')}
                </p>
                <div className="stats stats-vertical shadow mt-4">
                  <div className="stat">
                    <div className="stat-title">{t('common.total')}</div>
                    <div className="stat-value">{categories.length}</div>
                    <div className="stat-desc">
                      {t('categories.totalCategories')}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
