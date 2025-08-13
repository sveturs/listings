'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import {
  translationAdminApi,
  TranslationVersion,
  VersionDiff,
} from '@/services/translationAdminApi';
import {
  ClockIcon,
  ArrowPathIcon,
  EyeIcon,
  DocumentTextIcon,
  UserIcon,
  CodeBracketIcon,
} from '@heroicons/react/24/outline';

interface VersionHistoryViewerProps {
  entityType: string;
  entityId: number;
  onClose: () => void;
}

export default function VersionHistoryViewer({
  entityType,
  entityId,
  onClose,
}: VersionHistoryViewerProps) {
  const t = useTranslations('admin.translations');

  const [versions, setVersions] = useState<TranslationVersion[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedVersions, setSelectedVersions] = useState<number[]>([]);
  const [diff, setDiff] = useState<VersionDiff | null>(null);
  const [isViewingDiff, setIsViewingDiff] = useState(false);
  const [isRollingBack, setIsRollingBack] = useState(false);

  useEffect(() => {
    loadVersionHistory();
  }, [entityType, entityId]);

  const loadVersionHistory = async () => {
    try {
      setIsLoading(true);
      const data = await translationAdminApi.versions.getByEntity(
        entityType,
        entityId
      );
      setVersions(data);
    } catch (error) {
      console.error('Failed to load version history:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleVersionSelect = (versionId: number) => {
    setSelectedVersions((prev) => {
      if (prev.includes(versionId)) {
        return prev.filter((id) => id !== versionId);
      } else if (prev.length < 2) {
        return [...prev, versionId];
      } else {
        // Replace the first selected version
        return [prev[1], versionId];
      }
    });
  };

  const viewDiff = async () => {
    if (selectedVersions.length !== 2) return;

    try {
      setIsLoading(true);
      const diffData = await translationAdminApi.versions.getDiff(
        selectedVersions[0],
        selectedVersions[1]
      );
      setDiff(diffData);
      setIsViewingDiff(true);
    } catch (error) {
      console.error('Failed to load diff:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const rollbackToVersion = async (versionId: number) => {
    const version = versions.find((v) => v.id === versionId);
    if (!version) return;

    const comment = prompt('Комментарий к откату (необязательно):');
    if (comment === null) return; // User cancelled

    try {
      setIsRollingBack(true);
      await translationAdminApi.versions.rollback(
        version.translation_id,
        versionId,
        comment || undefined
      );

      // Reload versions after rollback
      await loadVersionHistory();
      alert('Откат выполнен успешно');
    } catch (error) {
      console.error('Failed to rollback:', error);
      alert('Ошибка при выполнении отката');
    } finally {
      setIsRollingBack(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU');
  };

  const renderTextDiff = (changes: any[]) => {
    if (!changes || changes.length === 0) {
      return <p className="text-base-content/60">Нет изменений в тексте</p>;
    }

    return (
      <div className="space-y-2">
        {changes.map((change, index) => (
          <div key={index} className="bg-base-200 rounded-lg p-3">
            <div className="text-xs uppercase font-semibold mb-2">
              {change.type === 'modification'
                ? 'Изменение'
                : change.type === 'addition'
                  ? 'Добавление'
                  : 'Удаление'}
            </div>
            {change.old_text && (
              <div className="bg-error/10 p-2 rounded mb-2">
                <div className="text-xs text-error font-medium mb-1">Было:</div>
                <div className="text-sm">{change.old_text}</div>
              </div>
            )}
            {change.new_text && (
              <div className="bg-success/10 p-2 rounded">
                <div className="text-xs text-success font-medium mb-1">
                  Стало:
                </div>
                <div className="text-sm">{change.new_text}</div>
              </div>
            )}
          </div>
        ))}
      </div>
    );
  };

  if (isViewingDiff && diff) {
    return (
      <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div className="bg-base-100 rounded-lg max-w-4xl w-full max-h-[90vh] overflow-hidden flex flex-col">
          <div className="p-6 border-b border-base-300">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-bold">Сравнение версий</h2>
              <button
                className="btn btn-ghost btn-sm"
                onClick={() => setIsViewingDiff(false)}
              >
                ✕
              </button>
            </div>
            <div className="text-sm text-base-content/60 mt-2">
              Версия {diff.version1.version_number} vs Версия{' '}
              {diff.version2.version_number}
            </div>
          </div>

          <div className="flex-1 overflow-auto p-6 space-y-6">
            <div>
              <h3 className="font-semibold mb-3 flex items-center gap-2">
                <DocumentTextIcon className="h-4 w-4" />
                Изменения в тексте
              </h3>
              {renderTextDiff(diff.text_changes)}
            </div>

            {Object.keys(diff.metadata_changes).length > 0 && (
              <div>
                <h3 className="font-semibold mb-3 flex items-center gap-2">
                  <CodeBracketIcon className="h-4 w-4" />
                  Изменения в метаданных
                </h3>
                <div className="bg-base-200 rounded-lg p-4">
                  <pre className="text-sm overflow-auto">
                    {JSON.stringify(diff.metadata_changes, null, 2)}
                  </pre>
                </div>
              </div>
            )}
          </div>

          <div className="p-6 border-t border-base-300">
            <button
              className="btn btn-primary"
              onClick={() => setIsViewingDiff(false)}
            >
              Закрыть
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-base-100 rounded-lg max-w-4xl w-full max-h-[90vh] overflow-hidden flex flex-col">
        <div className="p-6 border-b border-base-300">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-xl font-bold">История версий</h2>
              <p className="text-base-content/60 mt-1">
                {entityType} ID: {entityId}
              </p>
            </div>
            <button className="btn btn-ghost btn-sm" onClick={onClose}>
              ✕
            </button>
          </div>

          {selectedVersions.length === 2 && (
            <div className="mt-4 flex gap-2">
              <button
                className="btn btn-primary btn-sm gap-2"
                onClick={viewDiff}
                disabled={isLoading}
              >
                <EyeIcon className="h-4 w-4" />
                Сравнить версии
              </button>
            </div>
          )}
        </div>

        <div className="flex-1 overflow-auto">
          {isLoading ? (
            <div className="flex justify-center items-center py-12">
              <span className="loading loading-spinner loading-lg"></span>
            </div>
          ) : versions.length === 0 ? (
            <div className="text-center py-12 text-base-content/60">
              История версий не найдена
            </div>
          ) : (
            <div className="p-6">
              <div className="text-sm text-base-content/60 mb-4">
                Найдено {versions.length} версий. Выберите до 2 версий для
                сравнения.
              </div>

              <div className="space-y-3">
                {versions.map((version) => (
                  <div
                    key={version.id}
                    className={`card bg-base-200/50 border-2 transition-colors cursor-pointer
                      ${
                        selectedVersions.includes(version.id)
                          ? 'border-primary bg-primary/10'
                          : 'border-transparent hover:border-base-300'
                      }`}
                    onClick={() => handleVersionSelect(version.id)}
                  >
                    <div className="card-body p-4">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                            <div className="badge badge-primary">
                              v{version.version_number}
                            </div>
                            <div className="text-xs text-base-content/60">
                              {version.language}
                            </div>
                            <div className="text-xs text-base-content/60">
                              {version.field_name}
                            </div>
                          </div>

                          <div className="text-sm mb-2 bg-base-100 p-2 rounded">
                            {version.translated_text || 'Пустой перевод'}
                          </div>

                          <div className="flex items-center gap-4 text-xs text-base-content/60">
                            <div className="flex items-center gap-1">
                              <ClockIcon className="h-3 w-3" />
                              {formatDate(version.changed_at)}
                            </div>
                            {version.changed_by && (
                              <div className="flex items-center gap-1">
                                <UserIcon className="h-3 w-3" />
                                Пользователь {version.changed_by}
                              </div>
                            )}
                          </div>

                          {version.change_comment && (
                            <div className="text-xs text-base-content/60 mt-2 italic">
                              "{version.change_comment}"
                            </div>
                          )}
                        </div>

                        <div className="flex items-center gap-2 ml-4">
                          {selectedVersions.includes(version.id) && (
                            <div className="badge badge-primary badge-sm">
                              Выбрано
                            </div>
                          )}

                          <div className="dropdown dropdown-end">
                            <label
                              tabIndex={0}
                              className="btn btn-ghost btn-xs"
                            >
                              •••
                            </label>
                            <ul
                              tabIndex={0}
                              className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52"
                            >
                              <li>
                                <button
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    rollbackToVersion(version.id);
                                  }}
                                  disabled={isRollingBack}
                                >
                                  <ArrowPathIcon className="h-4 w-4" />
                                  Откатиться к этой версии
                                </button>
                              </li>
                            </ul>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        <div className="p-6 border-t border-base-300">
          <div className="flex justify-between items-center">
            <div className="text-sm text-base-content/60">
              {selectedVersions.length > 0 &&
                `Выбрано версий: ${selectedVersions.length}/2`}
            </div>
            <button className="btn" onClick={onClose}>
              Закрыть
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
