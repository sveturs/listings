'use client';

import { useState, useEffect, useCallback, DragEvent } from 'react';
import { useTranslations } from 'next-intl';
import { adminApi } from '@/services/admin';
import { toast } from '@/utils/toast';
import type { Attribute, VariantAttribute } from '@/services/admin';

interface DragDropMappingEditorProps {
  variantAttribute: VariantAttribute;
  onSave?: () => void;
  onCancel?: () => void;
}

export default function DragDropMappingEditor({
  variantAttribute,
  onSave,
  onCancel,
}: DragDropMappingEditorProps) {
  const t = useTranslations('admin');
  const [availableAttributes, setAvailableAttributes] = useState<Attribute[]>(
    []
  );
  const [linkedAttributes, setLinkedAttributes] = useState<Attribute[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [draggedAttribute, setDraggedAttribute] = useState<Attribute | null>(
    null
  );
  const [isDragOverLinked, setIsDragOverLinked] = useState(false);
  const [isDragOverAvailable, setIsDragOverAvailable] = useState(false);

  const loadAttributes = useCallback(async () => {
    try {
      setIsLoading(true);

      // Загружаем все атрибуты категорий
      const response = await adminApi.attributes.getAll(1, 100, '', '');

      // Фильтруем только атрибуты с is_variant_compatible = true
      const variantCompatibleAttributes = response.data.filter(
        (attr: Attribute) => attr.is_variant_compatible === true
      );

      // Загружаем связи для текущего вариативного атрибута
      const mappings =
        await adminApi.variantAttributes.getVariantAttributeMappings(
          variantAttribute.id
        );

      // Разделяем на связанные и доступные
      const linkedIds = new Set(mappings.map((m: any) => m.id));
      const linked = variantCompatibleAttributes.filter((attr) =>
        linkedIds.has(attr.id)
      );
      const available = variantCompatibleAttributes.filter(
        (attr) => !linkedIds.has(attr.id)
      );

      setLinkedAttributes(linked);
      setAvailableAttributes(available);
    } catch (error) {
      console.error('Failed to load attributes:', error);
      toast.error(t('variantAttributes.loadMappingsError'));
    } finally {
      setIsLoading(false);
    }
  }, [variantAttribute.id, t]);

  useEffect(() => {
    loadAttributes();
  }, [variantAttribute.id, loadAttributes]);

  const handleDragStart = (
    e: DragEvent<HTMLDivElement>,
    attribute: Attribute
  ) => {
    setDraggedAttribute(attribute);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragEnd = () => {
    setDraggedAttribute(null);
    setIsDragOverLinked(false);
    setIsDragOverAvailable(false);
  };

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  const handleDropOnLinked = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragOverLinked(false);

    if (!draggedAttribute) return;

    // Если атрибут уже в связанных, ничего не делаем
    if (linkedAttributes.find((attr) => attr.id === draggedAttribute.id))
      return;

    // Перемещаем из доступных в связанные
    setAvailableAttributes((prev) =>
      prev.filter((attr) => attr.id !== draggedAttribute.id)
    );
    setLinkedAttributes((prev) => [...prev, draggedAttribute]);
  };

  const handleDropOnAvailable = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragOverAvailable(false);

    if (!draggedAttribute) return;

    // Если атрибут уже в доступных, ничего не делаем
    if (availableAttributes.find((attr) => attr.id === draggedAttribute.id))
      return;

    // Перемещаем из связанных в доступные
    setLinkedAttributes((prev) =>
      prev.filter((attr) => attr.id !== draggedAttribute.id)
    );
    setAvailableAttributes((prev) => [...prev, draggedAttribute]);
  };

  const handleAutoDetect = async () => {
    try {
      // Ищем атрибуты с похожими названиями
      const variantName = variantAttribute.name.toLowerCase();
      const autoDetected = availableAttributes.filter((attr) => {
        const attrName = attr.name.toLowerCase();
        const displayName = attr.display_name.toLowerCase();

        // Точное совпадение
        if (attrName === variantName) return true;

        // Проверяем похожесть (например, color/colour, size/sizes)
        if (attrName.includes(variantName) || variantName.includes(attrName))
          return true;
        if (
          displayName.includes(variantName) ||
          variantName.includes(displayName)
        )
          return true;

        // Проверяем известные синонимы
        const synonyms: Record<string, string[]> = {
          color: ['colour', 'цвет', 'boja'],
          size: ['размер', 'veličina', 'величина'],
          memory: ['ram', 'память', 'memorija'],
          storage: ['disk', 'hdd', 'ssd', 'память'],
        };

        const variantSynonyms = synonyms[variantName] || [];
        if (
          variantSynonyms.some(
            (syn) => attrName.includes(syn) || displayName.includes(syn)
          )
        ) {
          return true;
        }

        return false;
      });

      if (autoDetected.length > 0) {
        // Перемещаем найденные атрибуты в связанные
        setAvailableAttributes((prev) =>
          prev.filter(
            (attr) => !autoDetected.find((detected) => detected.id === attr.id)
          )
        );
        setLinkedAttributes((prev) => [...prev, ...autoDetected]);

        toast.success(
          t('variantAttributes.autoDetectSuccess', {
            count: autoDetected.length,
          })
        );
      } else {
        toast.info(t('variantAttributes.autoDetectNoResults'));
      }
    } catch (error) {
      console.error('Auto-detect failed:', error);
      toast.error(t('variantAttributes.autoDetectError'));
    }
  };

  const handleSave = async () => {
    try {
      setIsSaving(true);

      // Отправляем только ID связанных атрибутов
      await adminApi.variantAttributes.updateVariantAttributeMappings(
        variantAttribute.id,
        linkedAttributes.map((attr) => attr.id)
      );

      toast.success(t('variantAttributes.mappingsUpdated'));
      onSave?.();
    } catch (error) {
      console.error('Failed to save mappings:', error);
      toast.error(t('variantAttributes.saveMappingsError'));
    } finally {
      setIsSaving(false);
    }
  };

  const filteredAvailable = availableAttributes.filter(
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
          />
        </svg>
        <div>
          <h3 className="font-bold">{t('variantAttributes.dragDropTitle')}</h3>
          <div className="text-xs">{t('variantAttributes.dragDropHint')}</div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        {/* Доступные атрибуты */}
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <h3 className="font-semibold">
              {t('variantAttributes.availableAttributes')}
            </h3>
            <button
              className="btn btn-sm btn-ghost"
              onClick={handleAutoDetect}
              title={t('variantAttributes.autoDetect')}
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
                  d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                />
              </svg>
              {t('variantAttributes.autoDetect')}
            </button>
          </div>

          <input
            type="text"
            placeholder={t('variantAttributes.searchAttributes')}
            className="input input-bordered input-sm w-full"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />

          <div
            className={`border-2 rounded-lg p-4 min-h-[400px] transition-colors ${
              isDragOverAvailable
                ? 'border-primary bg-base-200'
                : 'border-base-300'
            }`}
            onDragOver={handleDragOver}
            onDragEnter={() => setIsDragOverAvailable(true)}
            onDragLeave={() => setIsDragOverAvailable(false)}
            onDrop={handleDropOnAvailable}
          >
            {filteredAvailable.length === 0 ? (
              <p className="text-base-content/50 text-center">
                {searchQuery
                  ? t('variantAttributes.noMatchingAttributes')
                  : t('variantAttributes.noAvailableAttributes')}
              </p>
            ) : (
              <div className="space-y-2">
                {filteredAvailable.map((attr) => (
                  <div
                    key={attr.id}
                    draggable
                    onDragStart={(e) => handleDragStart(e, attr)}
                    onDragEnd={handleDragEnd}
                    className="card bg-base-100 shadow-sm cursor-move hover:shadow-md transition-shadow"
                  >
                    <div className="card-body p-3">
                      <div className="flex items-center justify-between">
                        <div>
                          <div className="font-medium">{attr.display_name}</div>
                          <code className="text-xs text-base-content/70">
                            {attr.name}
                          </code>
                        </div>
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-5 w-5 text-base-content/30"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M4 6h16M4 12h16M4 18h16"
                          />
                        </svg>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Связанные атрибуты */}
        <div className="space-y-2">
          <h3 className="font-semibold">
            {t('variantAttributes.linkedAttributes')}
          </h3>

          <div className="text-sm text-base-content/70">
            {t('variantAttributes.linkedWith')}:{' '}
            <strong>{variantAttribute.display_name}</strong>
          </div>

          <div
            className={`border-2 rounded-lg p-4 min-h-[400px] transition-colors ${
              isDragOverLinked
                ? 'border-success bg-success/10'
                : 'border-base-300'
            }`}
            onDragOver={handleDragOver}
            onDragEnter={() => setIsDragOverLinked(true)}
            onDragLeave={() => setIsDragOverLinked(false)}
            onDrop={handleDropOnLinked}
          >
            {linkedAttributes.length === 0 ? (
              <p className="text-base-content/50 text-center">
                {t('variantAttributes.noLinkedAttributes')}
              </p>
            ) : (
              <div className="space-y-2">
                {linkedAttributes.map((attr) => (
                  <div
                    key={attr.id}
                    draggable
                    onDragStart={(e) => handleDragStart(e, attr)}
                    onDragEnd={handleDragEnd}
                    className="card bg-success/10 shadow-sm cursor-move hover:shadow-md transition-shadow"
                  >
                    <div className="card-body p-3">
                      <div className="flex items-center justify-between">
                        <div>
                          <div className="font-medium">{attr.display_name}</div>
                          <code className="text-xs text-base-content/70">
                            {attr.name}
                          </code>
                        </div>
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-5 w-5 text-success"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
                          />
                        </svg>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="flex justify-end gap-2 pt-4">
        <button
          className="btn btn-ghost"
          onClick={onCancel}
          disabled={isSaving}
        >
          {t('common.cancel')}
        </button>
        <button
          className="btn btn-primary"
          onClick={handleSave}
          disabled={isSaving || isLoading}
        >
          {isSaving && (
            <span className="loading loading-spinner loading-sm mr-2"></span>
          )}
          {t('common.save')}
        </button>
      </div>
    </div>
  );
}
