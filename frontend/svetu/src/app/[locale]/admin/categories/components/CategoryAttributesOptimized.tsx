'use client';

import { useState, useEffect, useMemo, useCallback } from 'react';
import {
  Category,
  AttributeGroup,
  adminApi,
  // CategoryAttributeMapping, // Will be used in future updates
} from '@/services/admin';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
// TranslationStatus will be used in future updates
import { BatchTranslationModal } from '@/components/attributes/BatchTranslationModal';
import { AttributeListVirtualized } from '@/components/attributes/AttributeListVirtualized';
import { AttributeFiltersComponent } from '@/components/attributes/AttributeFilters';
import { useAttributesPagination } from '@/hooks/useAttributesPagination';

interface CategoryAttributesOptimizedProps {
  category: Category;
  onUpdate: () => void;
}

// CategoryAttribute interface - will be used for type checking in future updates

export default function CategoryAttributesOptimized({
  category,
  onUpdate,
}: CategoryAttributesOptimizedProps) {
  const t = useTranslations('admin');

  // State for UI
  const [allGroups, setAllGroups] = useState<AttributeGroup[]>([]);
  const [categoryGroups, setCategoryGroups] = useState<AttributeGroup[]>([]);
  const [groupsLoading, setGroupsLoading] = useState(true);
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
  const [bulkOperationLoading, setBulkOperationLoading] = useState(false);
  const [selectedBulkGroup, setSelectedBulkGroup] = useState<number | null>(
    null
  );
  const [showBatchTranslateModal, setShowBatchTranslateModal] = useState(false);

  // Use paginated attributes hook
  const {
    attributes: allAttributes,
    loading: attributesLoading,
    hasMore,
    filters,
    updateFilters,
    clearFilters,
    loadMore,
    refresh: refreshAttributes,
  } = useAttributesPagination({}, 50);

  // Category attributes (filtered from all attributes)
  const [categoryAttributeIds, setCategoryAttributeIds] = useState<Set<number>>(
    new Set()
  );

  const categoryAttributes = useMemo(() => {
    return allAttributes.filter((attr) => categoryAttributeIds.has(attr.id));
  }, [allAttributes, categoryAttributeIds]);

  useEffect(() => {
    loadInitialData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [category.id]);

  const loadInitialData = async () => {
    try {
      setGroupsLoading(true);

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

      // Load category attributes IDs
      try {
        const categoryAttrs = await adminApi.categories.getAttributes(
          category.id
        );
        setCategoryAttributeIds(new Set(categoryAttrs.map((attr) => attr.id)));
      } catch (error) {
        console.error('Failed to load category attributes:', error);
        setCategoryAttributeIds(new Set());
      }
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to load initial data:', error);
    } finally {
      setGroupsLoading(false);
    }
  };

  const handleAddAttribute = async () => {
    if (!selectedAttribute) return;

    try {
      await adminApi.categories.addAttribute(category.id, selectedAttribute);
      setCategoryAttributeIds((prev) => new Set([...prev, selectedAttribute]));
      setSelectedAttribute(null);
      toast.success(t('categories.attributeAdded'));
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to add attribute:', error);
    }
  };

  const handleRemoveAttribute = async (attributeId: number) => {
    try {
      await adminApi.categories.removeAttribute(category.id, attributeId);
      setCategoryAttributeIds((prev) => {
        const newSet = new Set(prev);
        newSet.delete(attributeId);
        return newSet;
      });
      setSelectedAttributeIds((prev) =>
        prev.filter((id) => id !== attributeId)
      );
      toast.success(t('categories.attributeRemoved'));
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to remove attribute:', error);
    }
  };

  const handleAddGroup = async () => {
    if (!selectedGroup) return;

    try {
      await adminApi.categories.attachGroup(category.id, selectedGroup);
      const group = allGroups.find((g) => g.id === selectedGroup);
      if (group) {
        setCategoryGroups((prev) => [...prev, group]);
      }
      setSelectedGroup(null);
      toast.success(t('categories.groupAdded'));
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to add group:', error);
    }
  };

  const handleRemoveGroup = async (groupId: number) => {
    try {
      await adminApi.categories.detachGroup(category.id, groupId);
      setCategoryGroups((prev) => prev.filter((g) => g.id !== groupId));
      toast.success(t('categories.groupRemoved'));
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to remove group:', error);
    }
  };

  // Handle selection
  const handleAttributeSelection = useCallback((attributeId: number) => {
    setSelectedAttributeIds((prev) =>
      prev.includes(attributeId)
        ? prev.filter((id) => id !== attributeId)
        : [...prev, attributeId]
    );
  }, []);

  const handleSelectAll = useCallback(() => {
    if (selectedAttributeIds.length === categoryAttributes.length) {
      setSelectedAttributeIds([]);
    } else {
      setSelectedAttributeIds(categoryAttributes.map((attr) => attr.id));
    }
  }, [selectedAttributeIds.length, categoryAttributes]);

  // Bulk operations
  const handleBulkDelete = async () => {
    if (selectedAttributeIds.length === 0) return;

    const confirmMsg = t('categories.confirmBulkDelete', {
      count: selectedAttributeIds.length,
    });
    if (!confirm(confirmMsg)) return;

    setBulkOperationLoading(true);
    try {
      await Promise.all(
        selectedAttributeIds.map((attrId) =>
          adminApi.categories.removeAttribute(category.id, attrId)
        )
      );

      setCategoryAttributeIds((prev) => {
        const newSet = new Set(prev);
        selectedAttributeIds.forEach((id) => newSet.delete(id));
        return newSet;
      });

      toast.success(t('categories.bulkDeleteSuccess'));
      setSelectedAttributeIds([]);
      onUpdate();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to bulk delete attributes:', error);
    } finally {
      setBulkOperationLoading(false);
    }
  };

  const handleBulkUpdateRequired = async (isRequired: boolean) => {
    if (selectedAttributeIds.length === 0) return;

    setBulkOperationLoading(true);
    try {
      await Promise.all(
        selectedAttributeIds.map((attrId) =>
          adminApi.categories.updateAttributeSettings(category.id, attrId, {
            is_required: isRequired,
          })
        )
      );
      toast.success(t('categories.bulkUpdateSuccess'));
      setSelectedAttributeIds([]);
      await loadInitialData();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to bulk update attributes:', error);
    } finally {
      setBulkOperationLoading(false);
    }
  };

  const handleBulkMoveToGroup = async () => {
    if (selectedAttributeIds.length === 0 || !selectedBulkGroup) return;

    setBulkOperationLoading(true);
    try {
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
      await loadInitialData();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to bulk move attributes:', error);
    } finally {
      setBulkOperationLoading(false);
    }
  };

  if (groupsLoading) {
    return (
      <div className="flex justify-center items-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
        <span className="ml-2">{t('common.loading')}</span>
      </div>
    );
  }

  const availableAttributes = allAttributes.filter(
    (attr) => !categoryAttributeIds.has(attr.id)
  );
  const availableGroups = allGroups.filter(
    (group) => !categoryGroups.some((cg) => cg.id === group.id)
  );

  return (
    <div className="space-y-6">
      {/* Tabs */}
      <div className="tabs tabs-boxed">
        <button
          className={`tab ${activeTab === 'attributes' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('attributes')}
        >
          {t('attributes')} ({categoryAttributes.length})
        </button>
        <button
          className={`tab ${activeTab === 'groups' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('groups')}
        >
          {t('attributeGroups')} ({categoryGroups.length})
        </button>
      </div>

      {activeTab === 'attributes' ? (
        <div className="space-y-6">
          {/* Filters */}
          <AttributeFiltersComponent
            filters={filters}
            onFiltersChange={updateFilters}
            onClearFilters={clearFilters}
            loading={attributesLoading}
          />

          {/* Add attribute section */}
          <div className="bg-base-100 p-4 rounded-lg border border-base-300">
            <h3 className="font-medium mb-4">{t('categories.addAttribute')}</h3>
            <div className="flex gap-2">
              <select
                className="select select-bordered flex-1"
                value={selectedAttribute || ''}
                onChange={(e) =>
                  setSelectedAttribute(Number(e.target.value) || null)
                }
              >
                <option value="">{t('categories.selectAttribute')}</option>
                {availableAttributes.map((attr) => (
                  <option key={attr.id} value={attr.id}>
                    {attr.display_name} ({attr.attribute_type})
                  </option>
                ))}
              </select>
              <button
                className="btn btn-primary"
                onClick={handleAddAttribute}
                disabled={!selectedAttribute}
              >
                {t('common.add')}
              </button>
            </div>
          </div>

          {/* Bulk actions */}
          {selectedAttributeIds.length > 0 && (
            <div className="bg-base-200 p-4 rounded-lg border">
              <div className="flex flex-wrap gap-2 items-center">
                <span className="text-sm font-medium">
                  {t('common.selectedCount', {
                    count: selectedAttributeIds.length,
                  })}
                </span>

                <button
                  className="btn btn-error btn-sm"
                  onClick={handleBulkDelete}
                  disabled={bulkOperationLoading}
                >
                  {t('common.delete')}
                </button>

                <button
                  className="btn btn-warning btn-sm"
                  onClick={() => handleBulkUpdateRequired(true)}
                  disabled={bulkOperationLoading}
                >
                  {t('attributes.makeRequired')}
                </button>

                <button
                  className="btn btn-info btn-sm"
                  onClick={() => handleBulkUpdateRequired(false)}
                  disabled={bulkOperationLoading}
                >
                  {t('attributes.makeOptional')}
                </button>

                <div className="flex gap-1">
                  <select
                    className="select select-bordered select-sm"
                    value={selectedBulkGroup || ''}
                    onChange={(e) =>
                      setSelectedBulkGroup(Number(e.target.value) || null)
                    }
                  >
                    <option value="">{t('categories.selectGroup')}</option>
                    {categoryGroups.map((group) => (
                      <option key={group.id} value={group.id}>
                        {group.name}
                      </option>
                    ))}
                  </select>
                  <button
                    className="btn btn-primary btn-sm"
                    onClick={handleBulkMoveToGroup}
                    disabled={!selectedBulkGroup || bulkOperationLoading}
                  >
                    {t('categories.moveToGroup')}
                  </button>
                </div>

                <button
                  className="btn btn-secondary btn-sm"
                  onClick={() => setShowBatchTranslateModal(true)}
                >
                  {t('translations.batchTranslate')}
                </button>
              </div>
            </div>
          )}

          {/* Virtualized attributes list */}
          <AttributeListVirtualized
            attributes={categoryAttributes}
            loading={attributesLoading}
            hasMore={hasMore}
            onLoadMore={loadMore}
            selectedIds={selectedAttributeIds}
            onSelectionChange={handleAttributeSelection}
            onSelectAll={handleSelectAll}
            onRemove={handleRemoveAttribute}
            showTranslationStatus={true}
            containerHeight={500}
          />
        </div>
      ) : (
        // Groups tab content (existing implementation)
        <div className="space-y-4">
          {/* Add group section */}
          <div className="bg-base-100 p-4 rounded-lg border border-base-300">
            <h3 className="font-medium mb-4">{t('categories.addGroup')}</h3>
            <div className="flex gap-2">
              <select
                className="select select-bordered flex-1"
                value={selectedGroup || ''}
                onChange={(e) =>
                  setSelectedGroup(Number(e.target.value) || null)
                }
              >
                <option value="">{t('categories.selectGroup')}</option>
                {availableGroups.map((group) => (
                  <option key={group.id} value={group.id}>
                    {group.name}
                  </option>
                ))}
              </select>
              <button
                className="btn btn-primary"
                onClick={handleAddGroup}
                disabled={!selectedGroup}
              >
                {t('common.add')}
              </button>
            </div>
          </div>

          {/* Category groups list */}
          {categoryGroups.length === 0 ? (
            <div className="text-center py-8 text-base-content/60">
              {t('categories.noGroups')}
            </div>
          ) : (
            <div className="space-y-2">
              {categoryGroups.map((group) => (
                <div
                  key={group.id}
                  className="flex items-center justify-between p-3 bg-base-200 rounded-lg"
                >
                  <div>
                    <span className="font-medium">{group.name}</span>
                    {group.description && (
                      <p className="text-sm text-base-content/70">
                        {group.description}
                      </p>
                    )}
                  </div>
                  <button
                    className="btn btn-error btn-sm"
                    onClick={() => handleRemoveGroup(group.id)}
                  >
                    {t('common.remove')}
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Batch translation modal */}
      {showBatchTranslateModal && (
        <BatchTranslationModal
          isOpen={showBatchTranslateModal}
          entityType="attribute"
          selectedIds={selectedAttributeIds}
          onClose={() => setShowBatchTranslateModal(false)}
          onComplete={() => {
            setShowBatchTranslateModal(false);
            setSelectedAttributeIds([]);
            refreshAttributes();
          }}
        />
      )}
    </div>
  );
}
