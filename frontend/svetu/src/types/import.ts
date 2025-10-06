export interface ImportJob {
  id: number;
  storefront_id: number;
  user_id: number;
  file_name: string;
  file_type: 'xml' | 'csv' | 'zip';
  file_url?: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  total_records: number;
  processed_records: number;
  successful_records: number;
  failed_records: number;
  error_message?: string;
  started_at?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface ImportError {
  id: number;
  job_id: number;
  line_number: number;
  field_name: string;
  error_message: string;
  raw_data: string;
  created_at: string;
}

export interface ImportJobStatus {
  id: number;
  status: string;
  progress: number; // 0-100
  total_records: number;
  processed_records: number;
  successful_records: number;
  failed_records: number;
  errors?: ImportError[];
  started_at?: string;
  completed_at?: string;
  estimated_time_left?: number; // seconds
}

export interface ImportRequest {
  storefront_id: number;
  file_type: 'xml' | 'csv' | 'zip';
  file_url?: string;
  file_name?: string;
  update_mode: 'create_only' | 'update_only' | 'upsert';
  category_mapping_mode: 'auto' | 'manual' | 'skip';
}

export interface ImportSummary {
  job_id: number;
  total_records: number;
  successful_records: number;
  failed_records: number;
  new_products: number;
  updated_products: number;
  skipped_products: number;
  processing_time: string;
  completed_at: string;
}

export interface ImportFormats {
  supported_formats: {
    [key: string]: {
      description: string;
      file_extensions: string[];
      required_headers?: string[];
      optional_headers?: string[];
      sample_structure?: any;
      encoding?: string;
      delimiter?: string;
      note?: string;
    };
  };
  update_modes: string[];
  category_mapping_modes: string[];
  max_file_size: string;
  max_products_per_import: number;
}

export interface UploadProgress {
  loaded: number;
  total: number;
  percentage: number;
}

export interface ImportValidationError {
  field: string;
  message: string;
  value?: any;
}

// Preview types
export interface ImportPreviewRow {
  line_number: number;
  data: Record<string, any>;
  errors?: ImportValidationError[];
  is_valid: boolean;
}

export interface ImportPreviewResponse {
  file_type: string;
  total_rows: number;
  preview_rows: ImportPreviewRow[];
  headers?: string[]; // For CSV files
  validation_ok: boolean;
  error_summary?: string;
}

// UI State types
export interface ImportState {
  jobs: ImportJob[];
  currentJob: ImportJob | null;
  isUploading: boolean;
  uploadProgress: UploadProgress | null;
  validationErrors: ImportValidationError[];
  formats: ImportFormats | null;

  // Preview states
  previewData: ImportPreviewResponse | null;
  isPreviewLoading: boolean;
  previewError: string | null;

  // UI states
  isLoading: boolean;
  error: string | null;
  selectedFiles: File[];
  importUrl: string;
  selectedFileType: 'xml' | 'csv' | 'zip' | '';
  updateMode: 'create_only' | 'update_only' | 'upsert';
  categoryMappingMode: 'auto' | 'manual' | 'skip';

  // Modal states
  isImportModalOpen: boolean;
  isJobDetailsModalOpen: boolean;
  isErrorsModalOpen: boolean;
}

export interface FileUploadConfig {
  maxFileSize: number; // in bytes
  allowedTypes: string[];
  multiple: boolean;
}

export const IMPORT_FILE_CONFIG: FileUploadConfig = {
  maxFileSize: 100 * 1024 * 1024, // 100MB
  allowedTypes: [
    'text/csv',
    'application/csv',
    'text/xml',
    'application/xml',
    'application/zip',
    'application/x-zip-compressed',
  ],
  multiple: false,
};

export const IMPORT_STATUS_COLORS = {
  pending: 'bg-yellow-100 text-yellow-800',
  processing: 'bg-blue-100 text-blue-800',
  completed: 'bg-green-100 text-green-800',
  failed: 'bg-red-100 text-red-800',
} as const;

export const IMPORT_STATUS_ICONS = {
  pending: 'clock',
  processing: 'refresh',
  completed: 'check-circle',
  failed: 'x-circle',
} as const;
