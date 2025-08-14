'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';
import {
  translationAdminApi,
  SyncStatus,
  SyncConflict,
} from '@/services/translationAdminApi';
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
  skipped?: number;
  errors?: string[];
}

export default function SyncManager() {
  const t = useTranslations('admin');
  const [syncing, setSyncing] = useState(false);
  const [conflicts, setConflicts] = useState<SyncConflict[]>([]);
  const [syncStatus, setSyncStatus] = useState<SyncStatus | null>(null);
  const [lastSyncResult, setLastSyncResult] = useState<SyncResult | null>(null);
  const [_exportedData, _setExportedData] = useState<any>(null);

  useEffect(() => {
    fetchSyncStatus();
  }, []);

  const syncFrontendToDB = async () => {
    setSyncing(true);
    try {
      const response = await translationAdminApi.syncFrontendToDB();

      if (response.success && response.data) {
        setLastSyncResult({
          added: response.data.added || 0,
          updated: response.data.updated || 0,
          conflicts: response.data.conflicts || 0,
          total_items: response.data.total_items || response.data.total || 0,
          skipped: response.data.skipped || 0,
        });
        toast.success(t('translations.syncSuccess'));
        fetchSyncStatus();
      } else {
        throw new Error(response.error || 'Sync failed');
      }
    } catch (err) {
      console.error('Failed to sync from server', err);
      toast.error(t('translations.syncError'));
    } finally {
      setSyncing(false);
    }
  };

  const syncDBToFrontend = async () => {
    setSyncing(true);
    try {
      const response = await translationAdminApi.syncDBToFrontend();

      if (response.success && response.data) {
        setLastSyncResult({
          added: response.data.added || 0,
          updated: response.data.updated || 0,
          conflicts: response.data.conflicts || 0,
          total_items: response.data.total_items || response.data.total || 0,
          skipped: response.data.skipped || 0,
        });
        toast.success(t('translations.syncSuccess'));
        fetchSyncStatus();
      } else {
        throw new Error(response.error || 'Sync failed');
      }
    } catch (err) {
      console.error('Failed to sync from server', err);
      toast.error(t('translations.syncError'));
    } finally {
      setSyncing(false);
    }
  };

  const syncDBToOpenSearch = async () => {
    setSyncing(true);
    try {
      const response = await translationAdminApi.syncDBToOpenSearch();

      if (response.success) {
        toast.success(t('translations.syncSuccess'));
        fetchSyncStatus();
      } else {
        throw new Error(response.error || 'Sync failed');
      }
    } catch (err) {
      console.error('Failed to sync from server', err);
      toast.error(t('translations.syncError'));
    } finally {
      setSyncing(false);
    }
  };

  const exportTranslations = async () => {
    try {
      const response = await translationAdminApi.export({
        format: 'json',
        only_verified: false,
        include_metadata: true,
      });

      if (response.success && response.data) {
        // Create downloadable file
        const blob = new Blob([JSON.stringify(response.data, null, 2)], {
          type: 'application/json',
        });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `translations-export-${new Date().toISOString().split('T')[0]}.json`;
        a.click();
        URL.revokeObjectURL(url);

        toast.success(t('translations.exportSuccess'));
      } else {
        throw new Error(response.error || 'Export failed');
      }
    } catch (err) {
      console.error('Failed to export data', err);
      toast.error(t('translations.exportError'));
    }
  };

  const importTranslations = async (file: File) => {
    try {
      const fileContent = await file.text();
      const data = JSON.parse(fileContent);

      const response = await translationAdminApi.import({
        format: 'json',
        data,
        overwrite_existing: true,
        validate_only: false,
      });

      if (!response.success) {
        throw new Error(response.error || 'Import failed');
      }

      toast.success(
        t('translations.importSuccess', {
          count: response.data?.imported || response.data?.success || 0,
        })
      );
      fetchSyncStatus();
    } catch (err) {
      console.error('Failed to import data', err);
      toast.error(t('translations.importError'));
    }
  };

  const fetchSyncStatus = async () => {
    try {
      const response = await translationAdminApi.getSyncStatus();

      if (response.success && response.data) {
        setSyncStatus(response.data);
      }
    } catch (err) {
      console.error('Failed to fetch sync status', err);
    }
  };

  const fetchConflicts = async () => {
    try {
      const response = await translationAdminApi.getConflicts();

      if (response.success && response.data) {
        setConflicts(response.data);
      } else {
        setConflicts([]);
      }
    } catch (err) {
      console.error('Failed to fetch conflicts', err);
    }
  };

  const resolveConflict = async (
    conflictId: number,
    resolution: 'frontend' | 'database' | 'manual'
  ) => {
    try {
      const response = await translationAdminApi.resolveConflict(
        conflictId,
        resolution
      );

      if (response.success) {
        toast.success(t('translations.conflictResolved'));
        fetchConflicts();
      } else {
        throw new Error(response.error || 'Failed to resolve conflict');
      }
    } catch (err) {
      console.error('Failed to resolve conflict', err);
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
                  {syncStatus.last_sync?.completed_at
                    ? new Date(
                        syncStatus.last_sync.completed_at
                      ).toLocaleString()
                    : t('translations.never')}
                </span>
              </div>
              <div>
                <span className="text-base-content/60">
                  {t('translations.conflicts')}:
                </span>
                <span className="ml-2 font-semibold">
                  {syncStatus.conflicts_count || 0}
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
                        {conflict.entity_type} #{conflict.entity_id}
                      </span>
                      <span className="ml-2 text-sm text-base-content/60">
                        {conflict.field_name} ({conflict.language})
                      </span>
                    </div>
                    <span className="badge badge-warning">
                      {conflict.resolved ? 'Resolved' : 'Pending'}
                    </span>
                  </div>

                  <div className="grid grid-cols-2 gap-4 mb-4">
                    <div>
                      <p className="text-sm font-medium mb-1">
                        {t('translations.sourceValue')}:
                      </p>
                      <p className="text-sm bg-base-200 p-2 rounded">
                        {conflict.frontend_value || t('translations.empty')}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm font-medium mb-1">
                        {t('translations.targetValue')}:
                      </p>
                      <p className="text-sm bg-base-200 p-2 rounded">
                        {conflict.database_value || t('translations.empty')}
                      </p>
                    </div>
                  </div>

                  <div className="flex gap-2">
                    <button
                      className="btn btn-sm btn-success"
                      onClick={() => resolveConflict(conflict.id, 'frontend')}
                      disabled={conflict.resolved}
                    >
                      <CheckCircleIcon className="h-4 w-4 mr-1" />
                      {t('translations.useFrontend')}
                    </button>
                    <button
                      className="btn btn-sm btn-info"
                      onClick={() => resolveConflict(conflict.id, 'database')}
                      disabled={conflict.resolved}
                    >
                      <CheckCircleIcon className="h-4 w-4 mr-1" />
                      {t('translations.useDatabase')}
                    </button>
                    <button
                      className="btn btn-sm btn-warning"
                      onClick={() => {
                        const value = prompt('Enter manual value:');
                        if (value) resolveConflict(conflict.id, 'manual');
                      }}
                      disabled={conflict.resolved}
                    >
                      <XCircleIcon className="h-4 w-4 mr-1" />
                      {t('translations.manual')}
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
