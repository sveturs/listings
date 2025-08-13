'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';
import { tokenManager } from '@/utils/tokenManager';
import {
  ArrowPathIcon,
  CloudArrowUpIcon,
  CloudArrowDownIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  XCircleIcon,
  DocumentArrowDownIcon,
  DocumentArrowUpIcon,
} from '@heroicons/react/24/outline';

interface SyncResult {
  added: number;
  updated: number;
  conflicts: number;
  total_items: number;
}

export default function SyncManager() {
  const t = useTranslations('admin');
  const [syncing, setSyncing] = useState(false);
  const [conflicts, setConflicts] = useState<any[]>([]);
  const [syncStatus, setSyncStatus] = useState<any>(null);
  const [lastSyncResult, setLastSyncResult] = useState<SyncResult | null>(null);
  const [_exportedData, _setExportedData] = useState<any>(null);

  useEffect(() => {
    fetchSyncStatus();
  }, []);

  const syncFrontendToDB = async () => {
    setSyncing(true);
    try {
      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const response = await fetch(
        '/api/v1/admin/translations/sync/frontend-to-db',
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Sync failed');
      }

      const data = await response.json();
      setLastSyncResult(data.data);
      toast.success(t('translations.syncSuccess'));
      fetchSyncStatus();
    } catch {
      toast.error(t('translations.syncError'));
    } finally {
      setSyncing(false);
    }
  };

  const syncDBToFrontend = async () => {
    setSyncing(true);
    try {
      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const response = await fetch(
        '/api/v1/admin/translations/sync/db-to-frontend',
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Sync failed');
      }

      const data = await response.json();
      setLastSyncResult(data.data);
      toast.success(t('translations.syncSuccess'));
      fetchSyncStatus();
    } catch {
      toast.error(t('translations.syncError'));
    } finally {
      setSyncing(false);
    }
  };

  const syncDBToOpenSearch = async () => {
    setSyncing(true);
    try {
      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const response = await fetch(
        '/api/v1/admin/translations/sync/db-to-opensearch',
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Sync failed');
      }

      toast.success(t('translations.syncSuccess'));
      fetchSyncStatus();
    } catch {
      toast.error(t('translations.syncError'));
    } finally {
      setSyncing(false);
    }
  };

  const exportTranslations = async () => {
    try {
      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const response = await fetch('/api/v1/admin/translations/export', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Export failed');
      }

      const data = await response.json();
      // setExportedData(data.data); // TODO: Implement export data handling

      // Create downloadable file
      const blob = new Blob([JSON.stringify(data.data, null, 2)], {
        type: 'application/json',
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `translations-export-${new Date().toISOString().split('T')[0]}.json`;
      a.click();
      URL.revokeObjectURL(url);

      toast.success(t('translations.exportSuccess'));
    } catch {
      toast.error(t('translations.exportError'));
    }
  };

  const importTranslations = async (file: File) => {
    try {
      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const fileContent = await file.text();
      const translations = JSON.parse(fileContent);

      const response = await fetch('/api/v1/admin/translations/import', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          translations,
          overwrite_existing: true,
        }),
      });

      if (!response.ok) {
        throw new Error('Import failed');
      }

      const data = await response.json();
      toast.success(
        t('translations.importSuccess', {
          count: data.data?.success || 0,
        })
      );
      fetchSyncStatus();
    } catch {
      toast.error(t('translations.importError'));
    }
  };

  const fetchSyncStatus = async () => {
    try {
      const token = tokenManager.getAccessToken();
      if (!token) return;

      const response = await fetch('/api/v1/admin/translations/sync/status', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        setSyncStatus(data.data);
      }
    } catch {
      console.error('Failed to fetch sync status', _err);
    }
  };

  const fetchConflicts = async () => {
    try {
      const token = tokenManager.getAccessToken();
      if (!token) return;

      const response = await fetch(
        '/api/v1/admin/translations/sync/conflicts',
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        setConflicts(data.data || []);
      }
    } catch {
      console.error('Failed to fetch conflicts', _err);
    }
  };

  const resolveConflict = async (conflictId: number, resolution: string) => {
    try {
      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const response = await fetch(
        `/api/v1/admin/translations/sync/conflicts/${conflictId}/resolve`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ resolution }),
        }
      );

      if (!response.ok) {
        throw new Error('Failed to resolve conflict');
      }

      toast.success(t('translations.conflictResolved'));
      fetchConflicts();
    } catch {
      toast.error(t('translations.conflictResolveError'));
    }
  };

  return (
    <div className="space-y-6">
      {/* Sync Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div className="card bg-base-100">
          <div className="card-body">
            <h3 className="card-title">
              <CloudArrowUpIcon className="h-6 w-6" />
              {t('translations.syncFrontendToDB')}
            </h3>
            <p className="text-base-content/60">
              {t('translations.syncFrontendToDBDescription')}
            </p>
            <div className="card-actions justify-end mt-4">
              <button
                className="btn btn-primary"
                onClick={syncFrontendToDB}
                disabled={syncing}
              >
                {syncing ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    {t('translations.syncing')}
                  </>
                ) : (
                  <>
                    <ArrowPathIcon className="h-5 w-5 mr-2" />
                    {t('translations.startSync')}
                  </>
                )}
              </button>
            </div>
          </div>
        </div>

        <div className="card bg-base-100">
          <div className="card-body">
            <h3 className="card-title">
              <CloudArrowDownIcon className="h-6 w-6" />
              {t('translations.syncDBToFrontend')}
            </h3>
            <p className="text-base-content/60">
              {t('translations.syncDBToFrontendDescription')}
            </p>
            <div className="card-actions justify-end mt-4">
              <button
                className="btn btn-primary"
                onClick={syncDBToFrontend}
                disabled={syncing}
              >
                {syncing ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    {t('translations.syncing')}
                  </>
                ) : (
                  <>
                    <ArrowPathIcon className="h-5 w-5 mr-2" />
                    {t('translations.startSync')}
                  </>
                )}
              </button>
            </div>
          </div>
        </div>

        <div className="card bg-base-100">
          <div className="card-body">
            <h3 className="card-title">
              <CloudArrowUpIcon className="h-6 w-6" />
              {t('translations.syncDBToOpenSearch')}
            </h3>
            <p className="text-base-content/60">
              {t('translations.syncDBToOpenSearchDescription')}
            </p>
            <div className="card-actions justify-end mt-4">
              <button
                className="btn btn-primary"
                onClick={syncDBToOpenSearch}
                disabled={syncing}
              >
                {syncing ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    {t('translations.syncing')}
                  </>
                ) : (
                  <>
                    <ArrowPathIcon className="h-5 w-5 mr-2" />
                    {t('translations.startSync')}
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Export/Import Section */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="card bg-base-100">
          <div className="card-body">
            <h3 className="card-title">
              <DocumentArrowDownIcon className="h-6 w-6" />
              {t('translations.exportTranslations')}
            </h3>
            <p className="text-base-content/60">
              {t('translations.exportDescription')}
            </p>
            <div className="card-actions justify-end mt-4">
              <button
                className="btn btn-secondary"
                onClick={exportTranslations}
              >
                <DocumentArrowDownIcon className="h-5 w-5 mr-2" />
                {t('translations.export')}
              </button>
            </div>
          </div>
        </div>

        <div className="card bg-base-100">
          <div className="card-body">
            <h3 className="card-title">
              <DocumentArrowUpIcon className="h-6 w-6" />
              {t('translations.importTranslations')}
            </h3>
            <p className="text-base-content/60">
              {t('translations.importDescription')}
            </p>
            <div className="card-actions justify-end mt-4">
              <input
                type="file"
                id="import-file"
                accept=".json"
                className="hidden"
                onChange={(e) => {
                  const file = e.target.files?.[0];
                  if (file) {
                    importTranslations(file);
                  }
                }}
              />
              <label htmlFor="import-file" className="btn btn-secondary">
                <DocumentArrowUpIcon className="h-5 w-5 mr-2" />
                {t('translations.import')}
              </label>
            </div>
          </div>
        </div>
      </div>

      {/* Last Sync Result */}
      {lastSyncResult && (
        <div className="card bg-base-100">
          <div className="card-body">
            <h3 className="card-title">{t('translations.lastSyncResult')}</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="stat">
                <div className="stat-title">{t('translations.added')}</div>
                <div className="stat-value text-success">
                  {lastSyncResult.added}
                </div>
              </div>
              <div className="stat">
                <div className="stat-title">{t('translations.updated')}</div>
                <div className="stat-value text-info">
                  {lastSyncResult.updated}
                </div>
              </div>
              <div className="stat">
                <div className="stat-title">{t('translations.conflicts')}</div>
                <div className="stat-value text-warning">
                  {lastSyncResult.conflicts}
                </div>
              </div>
              <div className="stat">
                <div className="stat-title">{t('translations.totalItems')}</div>
                <div className="stat-value">{lastSyncResult.total_items}</div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Sync Status */}
      {syncStatus && (
        <div className="card bg-base-100">
          <div className="card-body">
            <h3 className="card-title">{t('translations.syncStatus')}</h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <span className="text-base-content/60">
                  {t('translations.lastSync')}:
                </span>
                <span className="ml-2">
                  {syncStatus.last_sync
                    ? new Date(syncStatus.last_sync).toLocaleString()
                    : t('translations.never')}
                </span>
              </div>
              <div>
                <span className="text-base-content/60">
                  {t('translations.conflicts')}:
                </span>
                <span className="ml-2 font-semibold">
                  {syncStatus.conflicts || 0}
                </span>
              </div>
            </div>
            {syncStatus.in_progress && (
              <div className="alert alert-info mt-4">
                <span className="loading loading-spinner"></span>
                <span>{t('translations.syncInProgress')}</span>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Conflicts */}
      {conflicts.length > 0 && (
        <div className="card bg-base-100">
          <div className="card-body">
            <h3 className="card-title">
              <ExclamationTriangleIcon className="h-6 w-6 text-warning" />
              {t('translations.conflicts')} ({conflicts.length})
            </h3>
            <div className="space-y-4">
              {conflicts.map((conflict) => (
                <div key={conflict.id} className="border rounded-lg p-4">
                  <div className="flex justify-between items-start mb-2">
                    <div>
                      <span className="font-medium">
                        {conflict.entity_identifier}
                      </span>
                      <span className="ml-2 text-sm text-base-content/60">
                        {conflict.source_type} → {conflict.target_type}
                      </span>
                    </div>
                    <span className="badge badge-warning">
                      {conflict.conflict_type}
                    </span>
                  </div>

                  <div className="grid grid-cols-2 gap-4 mb-4">
                    <div>
                      <p className="text-sm font-medium mb-1">
                        {t('translations.sourceValue')}:
                      </p>
                      <p className="text-sm bg-base-200 p-2 rounded">
                        {conflict.source_value || t('translations.empty')}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm font-medium mb-1">
                        {t('translations.targetValue')}:
                      </p>
                      <p className="text-sm bg-base-200 p-2 rounded">
                        {conflict.target_value || t('translations.empty')}
                      </p>
                    </div>
                  </div>

                  <div className="flex gap-2">
                    <button
                      className="btn btn-sm btn-success"
                      onClick={() => resolveConflict(conflict.id, 'use_source')}
                    >
                      <CheckCircleIcon className="h-4 w-4 mr-1" />
                      {t('translations.useSource')}
                    </button>
                    <button
                      className="btn btn-sm btn-info"
                      onClick={() => resolveConflict(conflict.id, 'use_target')}
                    >
                      <CheckCircleIcon className="h-4 w-4 mr-1" />
                      {t('translations.useTarget')}
                    </button>
                    <button
                      className="btn btn-sm btn-error"
                      onClick={() => resolveConflict(conflict.id, 'skip')}
                    >
                      <XCircleIcon className="h-4 w-4 mr-1" />
                      {t('translations.skip')}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Load Actions */}
      <div className="flex gap-2">
        <button className="btn btn-outline" onClick={fetchSyncStatus}>
          {t('translations.checkStatus')}
        </button>
        <button className="btn btn-outline" onClick={fetchConflicts}>
          {t('translations.loadConflicts')}
        </button>
      </div>
    </div>
  );
}
