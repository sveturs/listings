'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { adminApi } from '@/services/admin';
import { toast } from '@/utils/toast';
import type { Attribute, VariantAttribute } from '@/services/admin';

interface AttributeMappingEditorProps {
  variantAttribute: VariantAttribute;
  onSave?: () => void;
  onCancel?: () => void;
}

interface LinkedAttribute {
  id: number;
  name: string;
  display_name: string;
  is_linked: boolean;
}

export default function AttributeMappingEditor({
  variantAttribute,
  onSave,
  onCancel,
}: AttributeMappingEditorProps) {
  const t = useTranslations('admin.variantAttributes');
  const [categoryAttributes, setCategoryAttributes] = useState<
    LinkedAttribute[]
  >([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedAttributes, setSelectedAttributes] = useState<Set<number>>(
    new Set()
  );

  useEffect(() => {
    loadAttributes();
  }, [variantAttribute.id]);

  const loadAttributes = async () => {
    try {
      setIsLoading(true);

      // Загружаем все атрибуты категорий через существующий метод
      const response = await adminApi.attributes.getAll(
        1, // page
        100, // limit
        '', // search
        '' // filterType
      );

      // Фильтруем только атрибуты с is_variant_compatible = true
      const variantCompatibleAttributes = response.data.filter(
        (attr: Attribute) => attr.is_variant_compatible === true
      );

      // Загружаем связи для текущего вариативного атрибута
      const mappings =
        await adminApi.variantAttributes.getVariantAttributeMappings(
          variantAttribute.id
        );
      const linkedIds = new Set(
        mappings.map((m: any) => m.category_attribute_id)
      );

      // Формируем список с отметками о связях
      const attributesWithLinks = variantCompatibleAttributes.map(
        (attr: Attribute) => ({
          id: attr.id,
          name: attr.name,
          display_name: attr.display_name,
          is_linked: linkedIds.has(attr.id),
        })
      );

      setCategoryAttributes(attributesWithLinks);
      setSelectedAttributes(linkedIds);
    } catch (error) {
      console.error('Failed to load attributes:', error);
      toast.error(t('loadMappingsError'));
    } finally {
      setIsLoading(false);
    }
  };

  const handleToggleAttribute = (attributeId: number) => {
    setSelectedAttributes((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(attributeId)) {
        newSet.delete(attributeId);
      } else {
        newSet.add(attributeId);
      }
      return newSet;
    });
  };

  const handleAutoDetect = async () => {
    // Автоматическое определение связей по похожим названиям
    const variantName = variantAttribute.name.toLowerCase();
    const autoDetected = new Set<number>();

    categoryAttributes.forEach((attr) => {
      const attrName = attr.name.toLowerCase();
      // Проверяем совпадение или вхождение
      if (
        attrName === variantName ||
        attrName.includes(variantName) ||
        variantName.includes(attrName)
      ) {
        autoDetected.add(attr.id);
      }
    });

    if (autoDetected.size > 0) {
      setSelectedAttributes(autoDetected);
      toast.success(t('autoDetectSuccess', { count: autoDetected.size }));
    } else {
      toast.info(t('autoDetectNoResults'));
    }
  };

  const handleSave = async () => {
    try {
      setIsSaving(true);

      // Отправляем обновленные связи
      await adminApi.variantAttributes.updateVariantAttributeMappings(
        variantAttribute.id,
        Array.from(selectedAttributes)
      );

      toast.success(t('saveMappingsSuccess'));
      if (onSave) onSave();
    } catch (error) {
      console.error('Failed to save mappings:', error);
      toast.error(t('saveMappingsError'));
    } finally {
      setIsSaving(false);
    }
  };

  const filteredAttributes = categoryAttributes.filter(
    (attr) =>
      searchQuery === '' ||
      attr.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      attr.display_name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div>
        <h3 className="text-lg font-semibold">
          {t('mappingTitle')}: {variantAttribute.display_name}
        </h3>
        <p className="text-sm text-base-content/70 mt-1">
          {t('mappingDescription')}
        </p>
      </div>

      {/* Actions Bar */}
      <div className="flex gap-2 items-center">
        <div className="flex-1">
          <input
            type="text"
            placeholder={t('searchAttributes')}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="input input-bordered input-sm w-full"
          />
        </div>
        <button
          onClick={handleAutoDetect}
          className="btn btn-sm btn-secondary"
          disabled={isSaving}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-4 w-4 mr-1"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
            />
          </svg>
          {t('autoDetect')}
        </button>
      </div>

      {/* Attributes List */}
      <div className="border rounded-lg max-h-96 overflow-y-auto">
        {filteredAttributes.length === 0 ? (
          <div className="p-8 text-center text-base-content/60">
            {t('noAttributesFound')}
          </div>
        ) : (
          <div className="divide-y">
            {filteredAttributes.map((attr) => (
              <label
                key={attr.id}
                className="flex items-center p-3 hover:bg-base-200 cursor-pointer transition-colors"
              >
                <input
                  type="checkbox"
                  checked={selectedAttributes.has(attr.id)}
                  onChange={() => handleToggleAttribute(attr.id)}
                  className="checkbox checkbox-primary mr-3"
                />
                <div className="flex-1">
                  <div className="font-medium">{attr.display_name}</div>
                  <div className="text-sm text-base-content/60">
                    {attr.name}
                  </div>
                </div>
                {attr.is_linked && !selectedAttributes.has(attr.id) && (
                  <span className="badge badge-ghost badge-sm">
                    {t('previouslyLinked')}
                  </span>
                )}
              </label>
            ))}
          </div>
        )}
      </div>

      {/* Summary */}
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
          {t('selectedCount', { count: selectedAttributes.size })}
        </span>
      </div>

      {/* Actions */}
      <div className="flex justify-end gap-2">
        {onCancel && (
          <button
            onClick={onCancel}
            className="btn btn-ghost"
            disabled={isSaving}
          >
            {t('cancel')}
          </button>
        )}
        <button
          onClick={handleSave}
          className="btn btn-primary"
          disabled={isSaving || selectedAttributes.size === 0}
        >
          {isSaving ? (
            <>
              <span className="loading loading-spinner loading-sm"></span>
              {t('saving')}
            </>
          ) : (
            t('saveMappings')
          )}
        </button>
      </div>
    </div>
  );
}
