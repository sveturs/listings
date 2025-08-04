'use client';

import { useState, useEffect } from 'react';
import { Category, adminApi } from '@/services/admin';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import IconPicker from '@/components/IconPicker';
import { TranslationStatus } from '@/components/attributes/TranslationStatus';

interface CategoryFormProps {
  category?: Category | null;
  categories: Category[];
  onSave: (data: Partial<Category>) => void;
  onCancel: () => void;
}

export default function CategoryForm({
  category,
  categories,
  onSave,
  onCancel,
}: CategoryFormProps) {
  const t = useTranslations('admin');
  const tCommon = useTranslations('admin');

  const [formData, setFormData] = useState<Partial<Category>>({
    name: '',
    slug: '',
    parent_id: undefined,
    icon: '',
    description: '',
    is_active: true,
    seo_title: '',
    seo_description: '',
    seo_keywords: '',
  });

  const [isTranslating, setIsTranslating] = useState(false);

  useEffect(() => {
    if (category) {
      setFormData({
        name: category.name || '',
        slug: category.slug || '',
        parent_id: category.parent_id,
        icon: category.icon || '',
        description: category.description || '',
        is_active: category.is_active !== false,
        seo_title: category.seo_title || '',
        seo_description: category.seo_description || '',
        seo_keywords: category.seo_keywords || '',
      });
    }
  }, [category]);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value, type } = e.target;

    setFormData((prev) => ({
      ...prev,
      [name]:
        type === 'checkbox' ? (e.target as HTMLInputElement).checked : value,
    }));

    // Auto-generate slug from name
    if (name === 'name' && !category) {
      const slug = value
        .toLowerCase()
        .replace(/[^a-z0-9-]/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '');
      setFormData((prev) => ({ ...prev, slug }));
    }
  };

  const handleTranslate = async () => {
    if (!formData.name) {
      toast.error('Введите название категории для перевода');
      return;
    }

    setIsTranslating(true);
    try {
      const translations = await adminApi.translate(formData.name);

      // Here we would typically save these translations to the database
      // For now, we'll just show a success message
      toast.success(
        'Переводы получены: ' + Object.values(translations).join(', ')
      );
    } catch (error) {
      toast.error('Ошибка при переводе');
      console.error('Translation error:', error);
    } finally {
      setIsTranslating(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name || !formData.slug) {
      toast.error('Заполните обязательные поля');
      return;
    }

    onSave(formData);
  };

  // Get available parent categories (exclude current category and its children)
  const getAvailableParents = () => {
    if (!category) return categories;

    const excludeIds = new Set<number>([category.id]);

    // Recursively find all children
    const findChildren = (parentId: number) => {
      categories
        .filter((c) => c.parent_id === parentId)
        .forEach((child) => {
          excludeIds.add(child.id);
          findChildren(child.id);
        });
    };

    findChildren(category.id);

    return categories.filter((c) => !excludeIds.has(c.id));
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('categories.categoryName')} *</span>
        </label>
        <input
          type="text"
          name="name"
          value={formData.name}
          onChange={handleChange}
          className="input input-bordered"
          required
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('categories.slug')} *</span>
        </label>
        <input
          type="text"
          name="slug"
          value={formData.slug}
          onChange={handleChange}
          className="input input-bordered"
          pattern="^[a-z0-9\-]+$"
          required
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('categories.parentCategory')}</span>
        </label>
        <select
          name="parent_id"
          value={formData.parent_id || ''}
          onChange={handleChange}
          className="select select-bordered"
        >
          <option value="">{t('categories.noParent')}</option>
          {getAvailableParents().map((cat) => (
            <option key={cat.id} value={cat.id}>
              {cat.name}
            </option>
          ))}
        </select>
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('categories.icon')}</span>
        </label>
        <IconPicker
          value={formData.icon || ''}
          onChange={(icon) => setFormData((prev) => ({ ...prev, icon }))}
          placeholder="Выберите иконку"
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('categories.description')}</span>
        </label>
        <textarea
          name="description"
          value={formData.description}
          onChange={handleChange}
          className="textarea textarea-bordered"
          rows={3}
        />
      </div>

      <div className="form-control">
        <label className="label cursor-pointer">
          <span className="label-text">{t('categories.isActive')}</span>
          <input
            type="checkbox"
            name="is_active"
            checked={formData.is_active}
            onChange={handleChange}
            className="checkbox checkbox-primary"
          />
        </label>
      </div>

      <div className="divider">SEO</div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('categories.seoTitle')}</span>
        </label>
        <input
          type="text"
          name="seo_title"
          value={formData.seo_title}
          onChange={handleChange}
          className="input input-bordered"
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('categories.seoDescription')}</span>
        </label>
        <textarea
          name="seo_description"
          value={formData.seo_description}
          onChange={handleChange}
          className="textarea textarea-bordered"
          rows={2}
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('categories.seoKeywords')}</span>
        </label>
        <input
          type="text"
          name="seo_keywords"
          value={formData.seo_keywords}
          onChange={handleChange}
          className="input input-bordered"
          placeholder="ключ1, ключ2, ключ3"
        />
      </div>

      {category?.id && (
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('categories.translationStatus')}</span>
          </label>
          <TranslationStatus
            entityType="category"
            entityId={category.id}
            onTranslateClick={async () => {
              try {
                await adminApi.translateCategory(category.id);
                toast.success(t('categories.translationSuccess'));
                // Обновляем форму после перевода
                if (onSave) {
                  onSave(formData);
                }
              } catch {
                toast.error(t('categories.translationError'));
              }
            }}
          />
        </div>
      )}

      {!category?.id && (
        <div className="flex gap-2">
          <button
            type="button"
            onClick={handleTranslate}
            className="btn btn-secondary"
            disabled={isTranslating || !formData.name}
          >
            {isTranslating ? (
              <>
                <span className="loading loading-spinner loading-sm"></span>
                {t('categories.translating')}
              </>
            ) : (
              <>🌍 {t('categories.translate')}</>
            )}
          </button>
        </div>
      )}

      <div className="flex gap-2 pt-4">
        <button type="submit" className="btn btn-primary flex-1">
          {tCommon('common.save')}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="btn btn-ghost flex-1"
        >
          {tCommon('common.cancel')}
        </button>
      </div>
    </form>
  );
}
