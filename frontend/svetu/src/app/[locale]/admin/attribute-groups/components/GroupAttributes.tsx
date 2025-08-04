'use client';

import { useState, useEffect } from 'react';
import {
  AttributeGroup,
  Attribute,
  AttributeGroupItem,
  adminApi,
} from '@/services/admin';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';

interface GroupAttributesProps {
  group: AttributeGroup;
  onUpdate: () => void;
  onClose: () => void;
}

export default function GroupAttributes({
  group,
  onUpdate,
  onClose,
}: GroupAttributesProps) {
  const t = useTranslations('admin');
  const tCommon = useTranslations('admin');

  const [allAttributes, setAllAttributes] = useState<Attribute[]>([]);
  const [groupItems, setGroupItems] = useState<AttributeGroupItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedAttribute, setSelectedAttribute] = useState<number | null>(
    null
  );

  useEffect(() => {
    loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [group.id]);

  const loadData = async () => {
    try {
      setLoading(true);

      // Load all attributes - для групп загружаем все атрибуты без пагинации
      const attributesResponse = await adminApi.attributes.getAll(1, 100);
      setAllAttributes(attributesResponse.data);

      // Load group items
      const { items } = await adminApi.attributeGroups.getWithItems(group.id);
      setGroupItems(items || []);
    } catch (error) {
      toast.error(tCommon('common.error'));
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
      await adminApi.attributeGroups.addItem(group.id, {
        attribute_id: selectedAttribute,
        sort_order: groupItems.length,
      });

      toast.success(tCommon('common.saveSuccess'));
      setSelectedAttribute(null);
      await loadData();
      onUpdate();
    } catch (error) {
      toast.error(tCommon('common.error'));
      console.error('Failed to add attribute:', error);
    }
  };

  const handleRemoveAttribute = async (attributeId: number) => {
    if (!confirm(tCommon('common.confirmDelete'))) return;

    try {
      await adminApi.attributeGroups.removeItem(group.id, attributeId);
      toast.success(tCommon('common.deleteSuccess'));
      await loadData();
      onUpdate();
    } catch (error) {
      toast.error(tCommon('common.error'));
      console.error('Failed to remove attribute:', error);
    }
  };

  // Get available attributes (not already in group)
  const availableAttributes = allAttributes.filter(
    (attr) => !groupItems.some((item) => item.attribute_id === attr.id)
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
      {/* Add attribute */}
      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('attributeGroups.availableAttributes')}</span>
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
            {tCommon('common.add')}
          </button>
        </div>
      </div>

      {/* Current attributes */}
      <div>
        <label className="label">
          <span className="label-text">{t('attributeGroups.selectedAttributes')}</span>
        </label>

        {groupItems.length === 0 ? (
          <p className="text-center text-base-content/60 py-4">
            Нет атрибутов в группе
          </p>
        ) : (
          <div className="space-y-2">
            {groupItems.map((item) => {
              const attribute = allAttributes.find(
                (a) => a.id === item.attribute_id
              );
              if (!attribute) return null;

              return (
                <div
                  key={item.id}
                  className="flex items-center gap-2 p-2 bg-base-200 rounded-lg"
                >
                  {item.icon && <span className="text-lg">{item.icon}</span>}
                  <span className="flex-1">
                    {item.custom_display_name || attribute.display_name}
                  </span>
                  <code className="text-xs">{attribute.name}</code>
                  <button
                    onClick={() => handleRemoveAttribute(item.attribute_id)}
                    className="btn btn-ghost btn-xs text-error"
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
              );
            })}
          </div>
        )}
      </div>

      <div className="text-sm text-base-content/60">
        <p>{t('attributeGroups.dragToReorder')}</p>
      </div>

      <div className="pt-4">
        <button onClick={onClose} className="btn btn-ghost w-full">
          {tCommon('common.close')}
        </button>
      </div>
    </div>
  );
}
