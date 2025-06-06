'use client';

import { useState, useEffect } from 'react';
import {
  Category,
  Attribute,
  AttributeGroup,
  adminApi,
  CategoryAttributeMapping,
} from '@/services/admin';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';

interface CategoryAttributesProps {
  category: Category;
  onUpdate: () => void;
}

interface CategoryAttribute extends Attribute {
  mapping?: CategoryAttributeMapping;
}

export default function CategoryAttributes({
  category,
  onUpdate,
}: CategoryAttributesProps) {
  const t = useTranslations('admin');
  const [allAttributes, setAllAttributes] = useState<Attribute[]>([]);
  const [categoryAttributes, setCategoryAttributes] = useState<
    CategoryAttribute[]
  >([]);
  const [_groups, _setGroups] = useState<AttributeGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedAttribute, setSelectedAttribute] = useState<number | null>(
    null
  );
  const [activeTab, setActiveTab] = useState<'attributes' | 'groups'>(
    'attributes'
  );

  useEffect(() => {
    loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [category.id]);

  const loadData = async () => {
    try {
      setLoading(true);

      // Load all attributes
      const attributes = await adminApi.attributes.getAll();
      setAllAttributes(attributes);

      // Load attribute groups
      const attributeGroups = await adminApi.attributeGroups.getAll();
      _setGroups(attributeGroups);

      // Load category attributes
      try {
        const categoryAttrs = await adminApi.categories.getAttributes(
          category.id
        );
        setCategoryAttributes(categoryAttrs);
      } catch (error) {
        console.error('Failed to load category attributes:', error);
        // If the endpoint doesn't exist or fails, use empty array
        setCategoryAttributes([]);
      }
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddAttribute = async () => {
    if (!selectedAttribute) {
      toast.error('Выберите атрибут');
      return;
    }

    try {
      await adminApi.categories.addAttribute(
        category.id,
        selectedAttribute,
        false
      );
      toast.success(t('common.saveSuccess'));
      setSelectedAttribute(null);
      await loadData();
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to add attribute:', error);
    }
  };

  const handleRemoveAttribute = async (attributeId: number) => {
    if (!confirm(t('common.confirmDelete'))) return;

    try {
      await adminApi.categories.removeAttribute(category.id, attributeId);
      toast.success(t('common.deleteSuccess'));
      await loadData();
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to remove attribute:', error);
    }
  };

  const handleUpdateAttributeSettings = async (
    attributeId: number,
    settings: Partial<CategoryAttributeMapping>
  ) => {
    try {
      await adminApi.categories.updateAttributeSettings(
        category.id,
        attributeId,
        settings
      );
      toast.success(t('common.saveSuccess'));
      await loadData();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to update attribute settings:', error);
    }
  };

  // Get available attributes (not already in category)
  const availableAttributes = allAttributes.filter(
    (attr) => !categoryAttributes.some((catAttr) => catAttr.id === attr.id)
  );

  if (loading) {
    return (
      <div className="flex justify-center items-center h-32">
        <span className="loading loading-spinner loading-md"></span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Tabs */}
      <div className="tabs tabs-boxed">
        <a
          className={`tab ${activeTab === 'attributes' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('attributes')}
        >
          {t('sections.attributes')}
        </a>
        <a
          className={`tab ${activeTab === 'groups' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('groups')}
        >
          {t('sections.attributeGroups')}
        </a>
      </div>

      {activeTab === 'attributes' && (
        <>
          {/* Add attribute */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">
                {t('attributeGroups.availableAttributes')}
              </span>
            </label>
            <div className="flex gap-2">
              <select
                value={selectedAttribute || ''}
                onChange={(e) => setSelectedAttribute(Number(e.target.value))}
                className="select select-bordered flex-1"
              >
                <option value="">Выберите атрибут</option>
                {availableAttributes.map((attr) => (
                  <option key={attr.id} value={attr.id}>
                    {attr.display_name} ({attr.name})
                  </option>
                ))}
              </select>
              <button
                onClick={handleAddAttribute}
                className="btn btn-primary"
                disabled={!selectedAttribute}
              >
                {t('common.add')}
              </button>
            </div>
          </div>

          {/* Current attributes */}
          <div>
            <label className="label">
              <span className="label-text">
                {t('attributeGroups.selectedAttributes')}
              </span>
            </label>

            {categoryAttributes.length === 0 ? (
              <p className="text-center text-base-content/60 py-4">
                Нет атрибутов в категории
              </p>
            ) : (
              <div className="space-y-2">
                {categoryAttributes.map((attr) => (
                  <div
                    key={attr.id}
                    className="flex items-center gap-2 p-3 bg-base-200 rounded-lg"
                  >
                    <span className="flex-1 font-medium">
                      {attr.display_name}
                    </span>
                    <code className="text-xs">{attr.name}</code>

                    <label className="label cursor-pointer gap-2">
                      <span className="label-text text-xs">
                        {t('attributes.isRequired')}
                      </span>
                      <input
                        type="checkbox"
                        className="checkbox checkbox-sm"
                        checked={attr.mapping?.is_required || false}
                        onChange={(e) =>
                          handleUpdateAttributeSettings(attr.id, {
                            is_required: e.target.checked,
                          })
                        }
                      />
                    </label>

                    <button
                      onClick={() => handleRemoveAttribute(attr.id)}
                      className="btn btn-ghost btn-sm text-error"
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
                          d="M6 18L18 6M6 6l12 12"
                        />
                      </svg>
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Inherited attributes */}
          {category.parent_id && (
            <div className="mt-6">
              <label className="label">
                <span className="label-text">Унаследованные атрибуты</span>
              </label>
              <div className="alert alert-info">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  className="stroke-current shrink-0 w-6 h-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  ></path>
                </svg>
                <span>
                  Атрибуты родительских категорий наследуются автоматически
                </span>
              </div>
            </div>
          )}
        </>
      )}

      {activeTab === 'groups' && (
        <div>
          <p className="text-center text-base-content/60 py-8">
            Управление группами атрибутов будет доступно в следующей версии
          </p>
        </div>
      )}
    </div>
  );
}
