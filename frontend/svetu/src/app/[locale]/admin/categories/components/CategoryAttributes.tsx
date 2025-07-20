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
import { TranslationStatus } from '@/components/attributes/TranslationStatus';
import { BatchTranslationModal } from '@/components/attributes/BatchTranslationModal';

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
  const [allGroups, setAllGroups] = useState<AttributeGroup[]>([]);
  const [categoryGroups, setCategoryGroups] = useState<AttributeGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedAttribute, setSelectedAttribute] = useState<number | null>(
    null
  );
  const [selectedGroup, setSelectedGroup] = useState<number | null>(null);
  const [activeTab, setActiveTab] = useState<'attributes' | 'groups'>(
    'attributes'
  );
  const [selectedAttributeIds, setSelectedAttributeIds] = useState<number[]>(
    []
  );
  const [_bulkOperationLoading, setBulkOperationLoading] = useState(false);
  const [selectedBulkGroup, setSelectedBulkGroup] = useState<number | null>(
    null
  );
  const [showBatchTranslateModal, setShowBatchTranslateModal] = useState(false);

  useEffect(() => {
    loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [category.id]);

  const loadData = async () => {
    try {
      setLoading(true);

      // Load all attributes - увеличиваем лимит для поддержки большего количества атрибутов
      const attributesResponse = await adminApi.attributes.getAll(1, 1000);
      setAllAttributes(attributesResponse.data);

      // Load all attribute groups
      const attributeGroups = await adminApi.attributeGroups.getAll();
      setAllGroups(attributeGroups);

      // Load category groups
      try {
        const categoryGroupsData = await adminApi.categories.getGroups(
          category.id
        );
        setCategoryGroups(categoryGroupsData);
      } catch (error) {
        console.error('Failed to load category groups:', error);
        setCategoryGroups([]);
      }

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

  const handleAddGroup = async () => {
    if (!selectedGroup) {
      toast.error('Выберите группу атрибутов');
      return;
    }

    try {
      await adminApi.categories.attachGroup(category.id, selectedGroup, 0);
      toast.success(t('common.saveSuccess'));
      setSelectedGroup(null);
      await loadData();
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to add group:', error);
    }
  };

  const handleRemoveGroup = async (groupId: number) => {
    if (!confirm(t('common.confirmDelete'))) return;

    try {
      await adminApi.categories.detachGroup(category.id, groupId);
      toast.success(t('common.deleteSuccess'));
      await loadData();
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to remove group:', error);
    }
  };

  // Get available attributes (not already in category)
  const availableAttributes = allAttributes.filter(
    (attr) => !categoryAttributes.some((catAttr) => catAttr.id === attr.id)
  );

  // Get available groups (not already in category)
  const availableGroups = allGroups.filter(
    (group) => !categoryGroups.some((catGroup) => catGroup.id === group.id)
  );

  // Check if all attributes are selected
  const isAllSelected =
    categoryAttributes.length > 0 &&
    selectedAttributeIds.length === categoryAttributes.length;

  // Handle select all/none
  const _handleSelectAll = () => {
    if (isAllSelected) {
      setSelectedAttributeIds([]);
    } else {
      setSelectedAttributeIds(categoryAttributes.map((attr) => attr.id));
    }
  };

  // Handle individual attribute selection
  const _handleSelectAttribute = (attributeId: number) => {
    setSelectedAttributeIds((prev) =>
      prev.includes(attributeId)
        ? prev.filter((id) => id !== attributeId)
        : [...prev, attributeId]
    );
  };

  // Handle bulk delete
  const _handleBulkDelete = async () => {
    if (selectedAttributeIds.length === 0) return;

    const confirmMsg = t('categories.confirmBulkDelete', {
      count: selectedAttributeIds.length,
    });
    if (!confirm(confirmMsg)) return;

    setBulkOperationLoading(true);
    try {
      // Delete all selected attributes
      await Promise.all(
        selectedAttributeIds.map((attrId) =>
          adminApi.categories.removeAttribute(category.id, attrId)
        )
      );
      toast.success(t('categories.bulkDeleteSuccess'));
      setSelectedAttributeIds([]);
      await loadData();
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to bulk delete attributes:', error);
    } finally {
      setBulkOperationLoading(false);
    }
  };

  // Handle bulk update required status
  const _handleBulkUpdateRequired = async (isRequired: boolean) => {
    if (selectedAttributeIds.length === 0) return;

    setBulkOperationLoading(true);
    try {
      // Update all selected attributes
      await Promise.all(
        selectedAttributeIds.map((attrId) =>
          adminApi.categories.updateAttributeSettings(category.id, attrId, {
            is_required: isRequired,
          })
        )
      );
      toast.success(t('categories.bulkUpdateSuccess'));
      setSelectedAttributeIds([]);
      await loadData();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to bulk update attributes:', error);
    } finally {
      setBulkOperationLoading(false);
    }
  };

  // Handle bulk move to group
  const _handleBulkMoveToGroup = async () => {
    if (selectedAttributeIds.length === 0 || !selectedBulkGroup) return;

    setBulkOperationLoading(true);
    try {
      // Move all selected attributes to the group
      await Promise.all(
        selectedAttributeIds.map((attrId) =>
          adminApi.categories.updateAttributeSettings(category.id, attrId, {
            group_id: selectedBulkGroup,
          })
        )
      );
      toast.success(t('categories.bulkMoveSuccess'));
      setSelectedAttributeIds([]);
      setSelectedBulkGroup(null);
      await loadData();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to bulk move attributes:', error);
    } finally {
      setBulkOperationLoading(false);
    }
  };

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
              {categoryAttributes.length > 0 && (
                <button
                  className="btn btn-ghost btn-xs"
                  onClick={() => {
                    if (
                      selectedAttributeIds.length === categoryAttributes.length
                    ) {
                      setSelectedAttributeIds([]);
                    } else {
                      setSelectedAttributeIds(
                        categoryAttributes.map((a) => a.id)
                      );
                    }
                  }}
                >
                  {selectedAttributeIds.length === categoryAttributes.length
                    ? t('common.deselectAll')
                    : t('common.selectAll')}
                </button>
              )}
            </label>

            {/* Bulk operations toolbar */}
            {selectedAttributeIds.length > 0 && (
              <div className="alert alert-info mb-4">
                <div className="flex items-center justify-between w-full">
                  <span>
                    {t('common.selected')}: {selectedAttributeIds.length}
                  </span>
                  <div className="flex gap-2">
                    <button
                      className="btn btn-sm btn-primary"
                      onClick={() => setShowBatchTranslateModal(true)}
                    >
                      🌍 {t('translations.translateSelected')}
                    </button>
                    <button
                      className="btn btn-sm btn-error"
                      onClick={() => _handleBulkDelete()}
                    >
                      {t('common.delete')}
                    </button>
                  </div>
                </div>
              </div>
            )}

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
                    <input
                      type="checkbox"
                      className="checkbox checkbox-sm"
                      checked={selectedAttributeIds.includes(attr.id)}
                      onChange={(e) => {
                        if (e.target.checked) {
                          setSelectedAttributeIds([
                            ...selectedAttributeIds,
                            attr.id,
                          ]);
                        } else {
                          setSelectedAttributeIds(
                            selectedAttributeIds.filter((id) => id !== attr.id)
                          );
                        }
                      }}
                    />
                    <span className="flex-1 font-medium">
                      {attr.display_name}
                    </span>
                    <code className="text-xs">{attr.name}</code>

                    <TranslationStatus
                      entityType="attribute"
                      entityId={attr.id}
                      compact={true}
                      onTranslateClick={async () => {
                        try {
                          await adminApi.translateAttribute(attr.id);
                          toast.success(
                            t('admin.translations.translationSuccess')
                          );
                          loadData();
                        } catch {
                          toast.error(t('admin.translations.translationError'));
                        }
                      }}
                    />

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
        <>
          {/* Add group */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">
                {t('attributeGroups.availableGroups')}
              </span>
            </label>
            <div className="flex gap-2">
              <select
                value={selectedGroup || ''}
                onChange={(e) => setSelectedGroup(Number(e.target.value))}
                className="select select-bordered flex-1"
              >
                <option value="">Выберите группу атрибутов</option>
                {availableGroups.map((group) => (
                  <option key={group.id} value={group.id}>
                    {group.display_name} ({group.name})
                  </option>
                ))}
              </select>
              <button
                onClick={handleAddGroup}
                className="btn btn-primary"
                disabled={!selectedGroup}
              >
                {t('common.add')}
              </button>
            </div>
          </div>

          {/* Current groups */}
          <div>
            <label className="label">
              <span className="label-text">
                {t('attributeGroups.selectedGroups')}
              </span>
            </label>

            {categoryGroups.length === 0 ? (
              <p className="text-center text-base-content/60 py-4">
                Нет групп атрибутов в категории
              </p>
            ) : (
              <div className="space-y-2">
                {categoryGroups.map((group) => (
                  <div
                    key={group.id}
                    className="flex items-center gap-2 p-3 bg-base-200 rounded-lg"
                  >
                    <span className="flex-1 font-medium">
                      {group.display_name}
                    </span>
                    <code className="text-xs">{group.name}</code>
                    {group.description && (
                      <span className="text-xs text-base-content/60">
                        {group.description}
                      </span>
                    )}

                    <button
                      onClick={() => handleRemoveGroup(group.id)}
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

          {/* Info about groups */}
          <div className="mt-6">
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
                Группы атрибутов помогают организовать атрибуты для лучшего
                пользовательского опыта при создании объявлений
              </span>
            </div>
          </div>
        </>
      )}

      {/* Batch Translation Modal */}
      <BatchTranslationModal
        isOpen={showBatchTranslateModal}
        onClose={() => setShowBatchTranslateModal(false)}
        entityType="attribute"
        selectedIds={selectedAttributeIds}
        selectedNames={categoryAttributes
          .filter((attr) => selectedAttributeIds.includes(attr.id))
          .map((attr) => attr.display_name)}
        onComplete={() => {
          loadData();
          setSelectedAttributeIds([]);
        }}
      />
    </div>
  );
}
