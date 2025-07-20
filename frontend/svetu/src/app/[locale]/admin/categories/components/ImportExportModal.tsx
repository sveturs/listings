'use client';

import { useState, useRef } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import { Category, Attribute, AttributeGroup } from '@/services/admin';

interface ImportExportModalProps {
  isOpen: boolean;
  onClose: () => void;
  category: Category;
  categoryAttributes: Attribute[];
  categoryGroups: AttributeGroup[];
  onImport: (data: ExportData) => Promise<void>;
}

export interface ExportData {
  categoryId: number;
  categoryName: string;
  exportDate: string;
  version: string;
  attributes: {
    id: number;
    name: string;
    display_name: string;
    attribute_type: string;
    is_required?: boolean;
    sort_order?: number;
    validation_rules?: Record<string, unknown>;
    custom_component?: string;
  }[];
  groups: {
    id: number;
    name: string;
    display_name: string;
    sort_order?: number;
  }[];
}

export default function ImportExportModal({
  isOpen,
  onClose,
  category,
  categoryAttributes,
  categoryGroups,
  onImport,
}: ImportExportModalProps) {
  const t = useTranslations('admin');
  const [importData, setImportData] = useState<ExportData | null>(null);
  const [showPreview, setShowPreview] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  if (!isOpen) return null;

  const handleExport = () => {
    const exportData: ExportData = {
      categoryId: category.id,
      categoryName: category.name,
      exportDate: new Date().toISOString(),
      version: '1.0',
      attributes: categoryAttributes.map((attr) => ({
        id: attr.id,
        name: attr.name,
        display_name: attr.display_name,
        attribute_type: attr.attribute_type,
        is_required: attr.is_required,
        sort_order: attr.sort_order,
        validation_rules: attr.validation_rules,
        custom_component: attr.custom_component,
      })),
      groups: categoryGroups.map((group) => ({
        id: group.id,
        name: group.name,
        display_name: group.display_name,
        sort_order: group.sort_order,
      })),
    };

    const blob = new Blob([JSON.stringify(exportData, null, 2)], {
      type: 'application/json',
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `category_${category.slug}_attributes_${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);

    toast.success(t('categories.exportSuccess'));
  };

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const data = JSON.parse(e.target?.result as string);

        // Validate data structure
        if (!data.attributes || !Array.isArray(data.attributes)) {
          throw new Error(t('categories.invalidFileFormat'));
        }

        setImportData(data);
        setShowPreview(true);
      } catch (error) {
        toast.error(t('categories.invalidFileFormat'));
        console.error('Failed to parse import file:', error);
      }
    };
    reader.readAsText(file);
  };

  const handleImport = async () => {
    if (!importData) return;

    try {
      await onImport(importData);
      toast.success(t('categories.importSuccess'));
      onClose();
    } catch (error) {
      toast.error(t('common.error'));
      console.error('Failed to import attributes:', error);
    }
  };

  const handleReset = () => {
    setImportData(null);
    setShowPreview(false);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-3xl">
        <h3 className="font-bold text-lg mb-4">
          {t('categories.importExportAttributes')}
        </h3>

        {!showPreview ? (
          <div className="space-y-6">
            {/* Export section */}
            <div className="border border-base-300 rounded-lg p-4">
              <h4 className="font-semibold mb-2">
                {t('categories.exportAttributes')}
              </h4>
              <p className="text-sm text-base-content/70 mb-4">
                {t('categories.exportDescription')}
              </p>
              <button onClick={handleExport} className="btn btn-primary">
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
                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                  />
                </svg>
                {t('common.export')}
              </button>
            </div>

            {/* Import section */}
            <div className="border border-base-300 rounded-lg p-4">
              <h4 className="font-semibold mb-2">
                {t('categories.importAttributes')}
              </h4>
              <p className="text-sm text-base-content/70 mb-4">
                {t('categories.importDescription')}
              </p>
              <input
                ref={fileInputRef}
                type="file"
                accept=".json"
                onChange={handleFileSelect}
                className="file-input file-input-bordered w-full"
              />
            </div>
          </div>
        ) : (
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
                ></path>
              </svg>
              <div>
                <div className="font-semibold">
                  {t('categories.importPreview')}
                </div>
                <div className="text-sm">
                  {t('categories.importFrom', {
                    category: importData?.categoryName || '',
                    date: new Date(
                      importData?.exportDate || ''
                    ).toLocaleDateString(),
                  })}
                </div>
              </div>
            </div>

            <div>
              <h4 className="font-semibold mb-2">
                {t('sections.attributes')} ({importData?.attributes.length || 0}
                )
              </h4>
              <div className="max-h-60 overflow-y-auto space-y-1">
                {importData?.attributes.map((attr) => (
                  <div
                    key={attr.id}
                    className="p-2 bg-base-200 rounded text-sm"
                  >
                    {attr.display_name} ({attr.name})
                    {attr.is_required && (
                      <span className="badge badge-sm badge-warning ml-2">
                        {t('attributes.isRequired')}
                      </span>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {importData?.groups && importData.groups.length > 0 && (
              <div>
                <h4 className="font-semibold mb-2">
                  {t('sections.attributeGroups')} ({importData.groups.length})
                </h4>
                <div className="max-h-40 overflow-y-auto space-y-1">
                  {importData.groups.map((group) => (
                    <div
                      key={group.id}
                      className="p-2 bg-base-200 rounded text-sm"
                    >
                      {group.display_name} ({group.name})
                    </div>
                  ))}
                </div>
              </div>
            )}

            <div className="flex gap-2">
              <button onClick={handleImport} className="btn btn-primary">
                {t('common.import')}
              </button>
              <button onClick={handleReset} className="btn btn-ghost">
                {t('common.cancel')}
              </button>
            </div>
          </div>
        )}

        <div className="modal-action">
          <button className="btn" onClick={onClose}>
            {t('common.close')}
          </button>
        </div>
      </div>
    </div>
  );
}
