export interface Translation {
  id: number;
  entity_type: string;
  entity_id: number;
  field_name: string;
  language: string;
  key: string;
  value_en?: string;
  value_ru?: string;
  value_sr?: string;
  context?: string;
  is_verified?: boolean;
  created_at: string;
  updated_at: string;
}

export interface TranslationVersion {
  id: number;
  entity_type: string;
  entity_id: number;
  field_name: string;
  language: string;
  translated_text: string;
  version: number;
  is_verified: boolean;
  created_by: number;
  created_at: string;
  changes?: string;
  previous_version_id?: number;
}

export interface TranslationProvider {
  id: number;
  name: string;
  provider_type: string;
  api_key?: string;
  api_endpoint?: string;
  is_active: boolean;
  priority: number;
  settings?: Record<string, any>;
  rate_limit?: number;
  cost_per_character?: number;
  supported_languages?: string[];
  created_at: string;
  updated_at: string;
}

export interface TranslationAuditLog {
  id: number;
  user_id: number;
  action: string;
  entity_type: string;
  entity_id: number;
  field_name?: string;
  old_value?: string;
  new_value?: string;
  metadata?: Record<string, any>;
  ip_address?: string;
  user_agent?: string;
  created_at: string;
}

export interface VersionDiff {
  version1: TranslationVersion;
  version2: TranslationVersion;
  changes: Array<{
    field: string;
    old_value: any;
    new_value: any;
  }>;
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