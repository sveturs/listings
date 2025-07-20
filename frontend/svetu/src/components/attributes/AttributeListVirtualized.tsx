import React, { useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { VirtualizedList } from '@/components/common/VirtualizedList';
import { TranslationStatus } from '@/components/attributes/TranslationStatus';
import { Attribute } from '@/services/admin';

export interface AttributeListVirtualizedProps {
  attributes: Attribute[];
  loading: boolean;
  hasMore: boolean;
  onLoadMore: () => void;
  selectedIds: number[];
  onSelectionChange: (attributeId: number) => void;
  onSelectAll: () => void;
  onEdit?: (attribute: Attribute) => void;
  onRemove?: (attributeId: number) => void;
  showTranslationStatus?: boolean;
  containerHeight?: number;
}

interface AttributeItemProps {
  attribute: Attribute;
  index: number;
  isSelected: boolean;
  onSelectionChange: (attributeId: number) => void;
  onEdit?: (attribute: Attribute) => void;
  onRemove?: (attributeId: number) => void;
  showTranslationStatus?: boolean;
}

const AttributeItem: React.FC<AttributeItemProps> = ({
  attribute,
  isSelected,
  onSelectionChange,
  onEdit,
  onRemove,
  showTranslationStatus = false,
}) => {
  const t = useTranslations('admin');

  return (
    <div className="flex items-center gap-3 p-3 bg-base-200 rounded-lg hover:bg-base-300 transition-colors w-full">
      <input
        type="checkbox"
        className="checkbox checkbox-sm"
        checked={isSelected}
        onChange={() => onSelectionChange(attribute.id)}
      />

      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="font-medium truncate">{attribute.display_name}</span>
          <span className="badge badge-sm badge-outline">
            {attribute.attribute_type}
          </span>
          {showTranslationStatus && (
            <TranslationStatus entityType="attribute" entityId={attribute.id} />
          )}
        </div>

        {/* Description field will be added in future updates */}
      </div>

      <div className="flex items-center gap-2 flex-shrink-0">
        {attribute.is_required && (
          <span className="badge badge-warning badge-xs">
            {t('attributes.required')}
          </span>
        )}

        <div className="dropdown dropdown-end">
          <div
            tabIndex={0}
            role="button"
            className="btn btn-ghost btn-sm btn-square"
          >
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
              />
            </svg>
          </div>
          <ul
            tabIndex={0}
            className="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow-lg border"
          >
            {onEdit && (
              <li>
                <button onClick={() => onEdit(attribute)}>
                  {t('common.edit')}
                </button>
              </li>
            )}
            {onRemove && (
              <li>
                <button
                  onClick={() => onRemove(attribute.id)}
                  className="text-error"
                >
                  {t('common.remove')}
                </button>
              </li>
            )}
          </ul>
        </div>
      </div>
    </div>
  );
};

export const AttributeListVirtualized: React.FC<
  AttributeListVirtualizedProps
> = ({
  attributes,
  loading,
  hasMore,
  onLoadMore,
  selectedIds,
  onSelectionChange,
  onSelectAll,
  onEdit,
  onRemove,
  showTranslationStatus = false,
  containerHeight = 400,
}) => {
  const t = useTranslations('admin');

  const renderAttributeItem = (attribute: Attribute, index: number) => (
    <AttributeItem
      key={attribute.id}
      attribute={attribute}
      index={index}
      isSelected={selectedIds.includes(attribute.id)}
      onSelectionChange={onSelectionChange}
      onEdit={onEdit}
      onRemove={onRemove}
      showTranslationStatus={showTranslationStatus}
    />
  );

  const isAllSelected = useMemo(() => {
    return attributes.length > 0 && selectedIds.length === attributes.length;
  }, [attributes.length, selectedIds.length]);

  const isPartiallySelected = useMemo(() => {
    return selectedIds.length > 0 && selectedIds.length < attributes.length;
  }, [attributes.length, selectedIds.length]);

  if (attributes.length === 0 && !loading) {
    return (
      <div className="text-center py-8 text-base-content/60">
        {t('attributes.noAttributesFound')}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header with bulk actions */}
      <div className="flex items-center justify-between">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            className="checkbox"
            checked={isAllSelected}
            ref={(el) => {
              if (el) el.indeterminate = isPartiallySelected;
            }}
            onChange={onSelectAll}
          />
          <span className="text-sm">
            {selectedIds.length > 0
              ? t('common.selectedCount', { count: selectedIds.length })
              : t('common.selectAll')}
          </span>
        </label>

        <div className="text-sm text-base-content/60">
          {t('common.total')}: {attributes.length}
          {hasMore && '+'}
        </div>
      </div>

      {/* Virtualized list */}
      <VirtualizedList
        items={attributes}
        itemHeight={80} // Approximate height of AttributeItem
        containerHeight={containerHeight}
        renderItem={renderAttributeItem}
        onLoadMore={onLoadMore}
        hasMore={hasMore}
        loading={loading}
        threshold={3}
        className="border border-base-300 rounded-lg"
      />
    </div>
  );
};
