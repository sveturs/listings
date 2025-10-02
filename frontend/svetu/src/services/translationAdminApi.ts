import { apiClient } from './api-client';

// Types для новой системы переводов
export interface Translation {
  id: number;
  entity_type: string;
  entity_id: number;
  field_name: string;
  language: string;
  translated_text: string;
  is_verified?: boolean;
  created_at?: string;
  updated_at?: string;
}

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
  version_number?: number;
  change_type: 'created' | 'updated' | 'deleted' | 'restored';
  changed_by?: number;
  changed_at: string;
  change_reason?: string;
  change_comment?: string;
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
      translations?: Record<string, string>; // Добавляем поле для переведенных текстов
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
  line: number;
  content: string;
  old_content?: string;
}

export interface AuditStatistics {
  total_actions: number;
  actions_by_type: Record<string, number>;
  actions_by_user: Array<{
    user_id: number;
    user_name: string;
    action_count: number;
  }>;
  actions_by_date: Array<{
    date: string;
    count: number;
  }>;
  recent_actions: TranslationAuditLog[];
}

export interface TranslationProvider {
  id: number;
  name: string;
  type: 'openai' | 'google' | 'deepl' | 'claude';
  is_active: boolean;
  api_key?: string;
  settings?: Record<string, any>;
  usage_count?: number;
  last_used?: string;
}

export interface TranslationStatistics {
  total_translations: number;
  complete_translations: number;
  missing_translations: number;
  placeholder_count: number;
  language_stats: Record<
    string,
    {
      total: number;
      complete: number;
      machine_translated: number;
      verified: number;
      coverage: number;
    }
  >;
  module_stats: Array<{
    name: string;
    keys: number;
    complete: number;
    incomplete: number;
    placeholders: number;
    missing: number;
    languages: Record<
      string,
      {
        total: number;
        complete: number;
        incomplete: number;
        placeholders: number;
        missing: number;
      }
    >;
  }>;
  recent_changes: TranslationAuditLog[];
}

export interface SyncConflict {
  id: number;
  entity_type: string;
  entity_id: number;
  field_name: string;
  language: string;
  frontend_value: string;
  database_value: string;
  detected_at: string;
  resolved: boolean;
  resolution?: 'frontend' | 'database' | 'manual';
  resolved_value?: string;
  resolved_by?: number;
  resolved_at?: string;
}

export interface SyncStatus {
  in_progress: boolean;
  last_sync?: {
    type: 'frontend-to-db' | 'db-to-frontend' | 'db-to-opensearch';
    started_at: string;
    completed_at?: string;
    items_processed: number;
    items_synced: number;
    errors: number;
  };
  conflicts_count: number;
}

// API Response wrapper
interface ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

// Translation Admin API Client
// Используем BFF proxy для всех запросов

class TranslationAdminApi {
  private getBaseUrl(): string {
    return '/admin/translations';
  }

  private async request<T>(
    path: string,
    options?: RequestInit
  ): Promise<ApiResponse<T>> {
    try {
      const baseUrl = this.getBaseUrl();
      const method = (options?.method || 'GET').toLowerCase() as
        | 'get'
        | 'post'
        | 'put'
        | 'delete';
      const body = options?.body ? JSON.parse(options.body as string) : undefined;

      let response;
      if (method === 'get') {
        response = await apiClient.get(`${baseUrl}${path}`);
      } else if (method === 'post') {
        response = await apiClient.post(`${baseUrl}${path}`, body);
      } else if (method === 'put') {
        response = await apiClient.put(`${baseUrl}${path}`, body);
      } else if (method === 'delete') {
        response = await apiClient.delete(`${baseUrl}${path}`);
      }

      if (!response?.data) {
        throw new Error('API request failed');
      }

      return response.data;
    } catch (error) {
      console.error('API request failed:', error);
      return {
        success: false,
        error:
          error instanceof Error ? error.message : 'Unknown error occurred',
      };
    }
  }

  // Statistics
  async getStatistics(): Promise<ApiResponse<TranslationStatistics>> {
    return this.request<TranslationStatistics>('/stats/overview');
  }

  async getCoverage(language?: string): Promise<ApiResponse<any>> {
    const params = language ? `?language=${language}` : '';
    return this.request(`/stats/coverage${params}`);
  }

  async getQuality(): Promise<ApiResponse<any>> {
    return this.request('/stats/quality');
  }

  async getUsage(): Promise<ApiResponse<any>> {
    return this.request('/stats/usage');
  }

  // Providers
  async getProviders(): Promise<ApiResponse<TranslationProvider[]>> {
    return this.request<TranslationProvider[]>('/providers');
  }

  async updateProvider(
    id: number,
    data: Partial<TranslationProvider>
  ): Promise<ApiResponse<TranslationProvider>> {
    return this.request<TranslationProvider>(`/providers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  // Bulk operations
  async bulkTranslate(
    request: BulkTranslateRequest
  ): Promise<ApiResponse<BulkTranslateResult>> {
    return this.request<BulkTranslateResult>('/bulk/translate', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  // Version history
  async getVersionHistory(
    entityType: string,
    entityId: number
  ): Promise<ApiResponse<TranslationVersion[]>> {
    return this.request<TranslationVersion[]>(
      `/versions/${entityType}/${entityId}`
    );
  }

  async getVersionDiff(
    version1Id: number,
    version2Id: number
  ): Promise<ApiResponse<VersionDiff>> {
    return this.request<VersionDiff>(
      `/versions/diff?v1=${version1Id}&v2=${version2Id}`
    );
  }

  async rollbackVersion(
    versionId: number,
    reason?: string
  ): Promise<ApiResponse<TranslationVersion>> {
    return this.request<TranslationVersion>('/versions/rollback', {
      method: 'POST',
      body: JSON.stringify({ version_id: versionId, reason }),
    });
  }

  // Sync operations
  async syncFrontendToDB(): Promise<ApiResponse<any>> {
    return this.request('/sync/frontend-to-db', {
      method: 'POST',
    });
  }

  async syncDBToFrontend(): Promise<ApiResponse<any>> {
    return this.request('/sync/db-to-frontend', {
      method: 'POST',
    });
  }

  async syncDBToOpenSearch(): Promise<ApiResponse<any>> {
    return this.request('/sync/db-to-opensearch', {
      method: 'POST',
    });
  }

  async getSyncStatus(): Promise<ApiResponse<SyncStatus>> {
    return this.request<SyncStatus>('/sync/status');
  }

  async getConflicts(): Promise<ApiResponse<SyncConflict[]>> {
    return this.request<SyncConflict[]>('/sync/conflicts');
  }

  async resolveConflict(
    conflictId: number,
    resolution: 'frontend' | 'database' | 'manual',
    value?: string
  ): Promise<ApiResponse<SyncConflict>> {
    return this.request<SyncConflict>(`/sync/conflicts/${conflictId}/resolve`, {
      method: 'POST',
      body: JSON.stringify({ resolution, value }),
    });
  }

  async resolveConflictsBatch(
    conflictIds: number[],
    resolution: 'frontend' | 'database'
  ): Promise<ApiResponse<any>> {
    return this.request('/sync/conflicts/resolve', {
      method: 'POST',
      body: JSON.stringify({ conflict_ids: conflictIds, resolution }),
    });
  }

  // Export/Import
  async export(request: ExportRequest): Promise<ApiResponse<any>> {
    return this.request('/export/advanced', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async import(request: ImportRequest): Promise<ApiResponse<any>> {
    return this.request('/import/advanced', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async validateImport(data: any, format: string): Promise<ApiResponse<any>> {
    return this.request('/import/validate', {
      method: 'POST',
      body: JSON.stringify({ data, format }),
    });
  }

  // Audit
  async getAuditLogs(
    page = 1,
    limit = 50,
    filters?: {
      user_id?: number;
      action?: string;
      entity_type?: string;
      start_date?: string;
      end_date?: string;
    }
  ): Promise<
    ApiResponse<{
      logs: TranslationAuditLog[];
      total: number;
      page: number;
      pages: number;
    }>
  > {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
      ...(filters?.user_id && { user_id: filters.user_id.toString() }),
      ...(filters?.action && { action: filters.action }),
      ...(filters?.entity_type && { entity_type: filters.entity_type }),
      ...(filters?.start_date && { start_date: filters.start_date }),
      ...(filters?.end_date && { end_date: filters.end_date }),
    });

    return this.request(`/audit/logs?${params}`);
  }

  async getAuditStatistics(): Promise<ApiResponse<AuditStatistics>> {
    return this.request<AuditStatistics>('/audit/statistics');
  }

  // AI Translation
  async translateWithAI(
    text: string,
    targetLanguage: string,
    providerId?: number
  ): Promise<
    ApiResponse<{
      translated_text: string;
      provider_used: string;
      confidence?: number;
    }>
  > {
    return this.request('/ai/translate', {
      method: 'POST',
      body: JSON.stringify({
        text,
        target_language: targetLanguage,
        provider_id: providerId,
      }),
    });
  }

  async translateBatch(
    texts: string[],
    targetLanguages: string[],
    providerId?: number
  ): Promise<ApiResponse<any>> {
    return this.request('/ai/batch', {
      method: 'POST',
      body: JSON.stringify({
        texts,
        target_languages: targetLanguages,
        provider_id: providerId,
      }),
    });
  }

  // Frontend module translations
  async getFrontendModules(): Promise<ApiResponse<string[]>> {
    return this.request<string[]>('/frontend/modules');
  }

  async getModuleTranslations(
    moduleName: string
  ): Promise<ApiResponse<Record<string, any>>> {
    return this.request(`/frontend/module/${moduleName}`);
  }

  async updateModuleTranslations(
    moduleName: string,
    translations: Record<string, any>
  ): Promise<ApiResponse<any>> {
    return this.request(`/frontend/module/${moduleName}`, {
      method: 'PUT',
      body: JSON.stringify(translations),
    });
  }

  // Database translations
  async getDatabaseTranslations(
    entityType?: string,
    entityId?: number
  ): Promise<ApiResponse<any>> {
    const params = new URLSearchParams();
    if (entityType) params.append('entity_type', entityType);
    if (entityId) params.append('entity_id', entityId.toString());

    return this.request(`/database?${params}`);
  }

  async updateTranslation(
    id: number,
    data: {
      translated_text: string;
      is_verified?: boolean;
    }
  ): Promise<ApiResponse<any>> {
    return this.request(`/database/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteTranslation(id: number): Promise<ApiResponse<any>> {
    return this.request(`/database/${id}`, {
      method: 'DELETE',
    });
  }

  // Search translations
  async searchTranslations(params: {
    query: string;
    entity_type?: string;
    language?: string;
    limit?: number;
    offset?: number;
  }): Promise<ApiResponse<Translation[]>> {
    const searchParams = new URLSearchParams();
    searchParams.append('q', params.query);
    if (params.entity_type)
      searchParams.append('entity_type', params.entity_type);
    if (params.language) searchParams.append('language', params.language);
    if (params.limit) searchParams.append('limit', params.limit.toString());
    if (params.offset) searchParams.append('offset', params.offset.toString());

    return this.request<Translation[]>(`/search?${searchParams}`);
  }

  // Duplicate translateWithAI removed - using the one at line 475
}

// Export singleton instance
export const translationAdminApi = new TranslationAdminApi();
