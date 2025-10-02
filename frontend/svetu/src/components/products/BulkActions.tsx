'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import {
  FiTrash2,
  FiToggleLeft,
  FiToggleRight,
  FiDownload,
  FiEdit3,
  FiPackage,
  FiX,
  FiAlertTriangle,
} from 'react-icons/fi';

interface BulkActionsProps {
  selectedCount: number;
  onBulkDelete: () => void;
  onBulkStatusChange: (active: boolean) => void;
  onBulkExport: () => void;
  onBulkEdit?: () => void;
  onClearSelection: () => void;
  isProcessing?: boolean;
}

export function BulkActions({
  selectedCount,
  onBulkDelete,
  onBulkStatusChange,
  onBulkExport,
  onBulkEdit,
  onClearSelection,
  isProcessing = false,
}: BulkActionsProps) {
  const t = useTranslations('storefronts');
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const handleDelete = () => {
    if (!showDeleteConfirm) {
      setShowDeleteConfirm(true);
      setTimeout(() => setShowDeleteConfirm(false), 5000); // Автосброс через 5 секунд
    } else {
      onBulkDelete();
      setShowDeleteConfirm(false);
    }
  };

  return (
    <div className="sticky top-0 z-40 bg-base-100 border-b border-base-300 shadow-lg">
      <div className="px-4 py-3">
        <div className="flex items-center justify-between">
          {/* Левая часть - информация о выборе */}
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <FiPackage className="text-primary w-5 h-5" />
              <span className="text-lg font-semibold">
                {t('selected', { count: selectedCount })}
              </span>
            </div>

            <button
              onClick={onClearSelection}
              className="btn btn-ghost btn-sm"
              disabled={isProcessing}
            >
              <FiX className="w-4 h-4" />
              {t('clearSelection')}
            </button>
          </div>

          {/* Правая часть - действия */}
          <div className="flex items-center gap-2">
            {/* Активировать/Деактивировать */}
            <div className="dropdown dropdown-end">
              <button
                tabIndex={0}
                className="btn btn-sm btn-ghost gap-2"
                disabled={isProcessing}
              >
                <FiToggleLeft className="w-4 h-4" />
                {t('statusLabel')}
              </button>
              <ul
                tabIndex={0}
                className="dropdown-content z-[1] menu p-2 shadow-lg bg-base-100 rounded-box w-52 border border-base-300"
              >
                <li>
                  <button
                    onClick={() => onBulkStatusChange(true)}
                    className="gap-2"
                  >
                    <FiToggleRight className="w-4 h-4 text-success" />
                    {t('activate')}
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => onBulkStatusChange(false)}
                    className="gap-2"
                  >
                    <FiToggleLeft className="w-4 h-4 text-warning" />
                    {t('deactivate')}
                  </button>
                </li>
              </ul>
            </div>

            {/* Массовое редактирование */}
            {onBulkEdit && (
              <button
                onClick={onBulkEdit}
                className="btn btn-sm btn-ghost gap-2"
                disabled={isProcessing}
              >
                <FiEdit3 className="w-4 h-4" />
                {t('edit')}
              </button>
            )}

            {/* Экспорт */}
            <button
              onClick={onBulkExport}
              className="btn btn-sm btn-ghost gap-2"
              disabled={isProcessing}
            >
              <FiDownload className="w-4 h-4" />
              {t('export')}
            </button>

            <div className="divider divider-horizontal mx-1"></div>

            {/* Удалить */}
            <button
              onClick={handleDelete}
              className={`btn btn-sm gap-2 ${
                showDeleteConfirm ? 'btn-error' : 'btn-ghost text-error'
              }`}
              disabled={isProcessing}
            >
              {showDeleteConfirm ? (
                <>
                  <FiAlertTriangle className="w-4 h-4" />
                  {t('confirmDelete')}
                </>
              ) : (
                <>
                  <FiTrash2 className="w-4 h-4" />
                  {t('delete')}
                </>
              )}
            </button>
          </div>
        </div>

        {/* Прогресс-бар при обработке */}
        {isProcessing && (
          <div className="mt-3">
            <progress className="progress progress-primary w-full"></progress>
            <p className="text-sm text-base-content/70 mt-1">
              {t('processing')}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
