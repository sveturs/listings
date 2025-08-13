// import { apiClient } from './api-client';

// Types для новой системы переводов
export interface TranslationVersion {
  id: number;
  translation_id: number;
  entity_type: string;
  entity_id: number;
  field_name: string;
  language: string;
  translated_text: string;
  previous_text?: string;
  version: number;
  version_number?: number; // для обратной совместимости
  change_type: 'created' | 'updated' | 'deleted' | 'restored';
  changed_by?: number;
  changed_at: string;
  change_reason?: string;
  change_comment?: string; // для обратной совместимости
  metadata?: Record<string, any>;
}

export interface TranslationAuditLog {
  id: number;
  user_id?: number;
  action: string;
  entity_type?: string;
  entity_id?: number;
  old_value?: string;
  new_value?: string;
  ip_address?: string;
  user_agent?: string;
  created_at: string;
}

export interface BulkTranslateRequest {
  entity_type: string;
  entity_ids?: number[];
  source_language: string;
  target_languages: string[];
  provider_id?: number;
  auto_approve: boolean;
  overwrite_existing: boolean;
}

export interface BulkTranslateResult {
  total_processed: number;
  successful: number;
  failed: number;
  skipped: number;
  errors?: string[];
  details?: {
    successful_items?: Array<{
      entity_id: number;
      entity_name: string;
      languages: string[];
    }>;
    failed_items?: Array<{
      entity_id: number;
      entity_name: string;
      error: string;
      language?: string;
    }>;
    skipped_items?: Array<{
      entity_id: number;
      entity_name: string;
      reason: string;
      existing_languages?: string[];
    }>;
  };
  processing_time?: number;
  provider_used?: string;
}

export interface ExportRequest {
  format: 'json' | 'csv' | 'xliff';
  entity_type?: string;
  language?: string;
  module?: string;
  only_verified: boolean;
  include_metadata: boolean;
}

export interface ImportRequest {
  format: 'json' | 'csv' | 'xliff';
  data: any;
  overwrite_existing: boolean;
  validate_only: boolean;
  metadata?: Record<string, any>;
}

export interface VersionDiff {
  version1: TranslationVersion;
  version2: TranslationVersion;
  text_changes: TextChange[];
  metadata_changes: Record<string, any>;
}

export interface TextChange {
  type: 'addition' | 'deletion' | 'modification';
  position: number;
  old_text: string;
  new_text: string;
  length: number;
}

export interface AuditStatistics {
  total_actions: number;
  actions_by_type: Record<string, number>;
  actions_by_user: Record<number, number>;
  recent_actions: TranslationAuditLog[];
}

// Helper function to get auth headers
async function getAuthHeaders(): Promise<Record<string, string>> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (typeof window !== 'undefined') {
    try {
      // For demo pages, use a test token
      if (window.location.pathname.includes('/demo/')) {
        // This is a hardcoded test token for demo purposes only
        headers['Authorization'] = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwiZXhwIjoxNzU1MTUyNTE5LCJpYXQiOjE3NTUwNjYxMTksImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjoxfQ.Tlq5EwIkiIYmlvEJ-TQMMz0_WT06xudZArjHZ8tArjs`;
      } else {
        const { tokenManager } = await import('@/utils/tokenManager');
        const token = await tokenManager.getAccessToken();
        if (token) {
          headers['Authorization'] = `Bearer ${token}`;
        }
      }
    } catch {
      console.log('No auth token available');
    }
  }

  return headers;
}

// Translation Admin API Service
export const translationAdminApi = {
  // Version History
  versions: {
    async getByEntity(
      entityType: string,
      entityId: number
    ): Promise<any> {
      const headers = await getAuthHeaders();

      const response = await fetch(
        `/api/v1/admin/translations/versions/${entityType}/${entityId}`,
        {
          method: 'GET',
          headers,
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      // API возвращает объект с versions внутри data
      return data.data || { versions: [], current_version: 0, total_versions: 0 };
    },

    async getDiff(
      versionId1: number,
      versionId2: number
    ): Promise<VersionDiff> {
      const headers = await getAuthHeaders();

      const response = await fetch(
        `/api/v1/admin/translations/versions/diff?version_id1=${versionId1}&version_id2=${versionId2}`,
        {
          method: 'GET',
          headers,
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data.data;
    },

    async rollback(
      translationId: number,
      versionId: number,
      comment?: string
    ): Promise<void> {
      const headers = await getAuthHeaders();

      const response = await fetch(
        `/api/v1/admin/translations/versions/rollback`,
        {
          method: 'POST',
          headers,
          body: JSON.stringify({
            translation_id: translationId,
            version_id: versionId,
            comment,
          }),
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
    },
  },

  // Audit Logs
  audit: {
    async getLogs(
      page: number = 1,
      pageSize: number = 20,
      filters?: {
        user_id?: number;
        action?: string;
        entity_type?: string;
        start_date?: string;
        end_date?: string;
      }
    ): Promise<{
      data: TranslationAuditLog[];
      total: number;
      page: number;
      total_pages: number;
    }> {
      const headers = await getAuthHeaders();
      const params = new URLSearchParams({
        page: page.toString(),
        page_size: pageSize.toString(),
      });

      if (filters) {
        Object.entries(filters).forEach(([key, value]) => {
          if (value !== undefined) {
            params.append(key, value.toString());
          }
        });
      }

      const response = await fetch(
        `/api/v1/admin/translations/audit/logs?${params.toString()}`,
        {
          method: 'GET',
          headers,
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return await response.json();
    },

    async getStatistics(): Promise<AuditStatistics> {
      const headers = await getAuthHeaders();

      const response = await fetch(
        `/api/v1/admin/translations/audit/statistics`,
        {
          method: 'GET',
          headers,
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data.data;
    },
  },

  // Bulk Operations
  bulk: {
    async translate(
      request: BulkTranslateRequest
    ): Promise<BulkTranslateResult> {
      const headers = await getAuthHeaders();

      const response = await fetch(
        `/api/v1/admin/translations/bulk/translate`,
        {
          method: 'POST',
          headers,
          body: JSON.stringify(request),
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data.data;
    },
  },

  // Export/Import
  async export(request: ExportRequest): Promise<any> {
    const headers = await getAuthHeaders();

    const response = await fetch(`/api/v1/admin/translations/export/advanced`, {
      method: 'POST',
      headers,
      body: JSON.stringify(request),
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    // Для CSV и XLIFF возвращаем blob для скачивания
    if (request.format === 'csv' || request.format === 'xliff') {
      return await response.blob();
    }

    // Для JSON возвращаем данные
    return await response.json();
  },

  async import(request: ImportRequest): Promise<{
    success: number;
    failed: number;
    skipped: number;
    errors?: string[];
  }> {
    const headers = await getAuthHeaders();

    const response = await fetch(`/api/v1/admin/translations/import/advanced`, {
      method: 'POST',
      headers,
      body: JSON.stringify(request),
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data.data;
  },

  // Translation Providers
  providers: {
    async getAll(): Promise<any[]> {
      const headers = await getAuthHeaders();

      const response = await fetch(`/api/v1/admin/translations/providers`, {
        method: 'GET',
        headers,
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data.data || [];
    },

    async update(id: number, provider: any): Promise<void> {
      const headers = await getAuthHeaders();

      const response = await fetch(
        `/api/v1/admin/translations/providers/${id}`,
        {
          method: 'PUT',
          headers,
          body: JSON.stringify(provider),
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
    },
  },
};
