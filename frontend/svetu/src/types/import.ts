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

  // Analysis states (Enhanced Import)
  analysisFile: File | null;
  analysisFileType: 'xml' | 'csv' | 'zip' | '';
  categoryAnalysis: CategoryAnalysisResponse | null;
  attributeAnalysis: AttributeAnalysisResponse | null;
  variantDetection: VariantDetectionResponse | null;
  clientCategories: ClientCategoriesResponse | null;
  isAnalyzing: boolean;
  analysisError: string | null;
  analysisProgress: number; // 0-100
  currentAnalysisStep: number;

  // User modifications for analysis
  approvedMappings: CategoryMapping[];
  customMappings: Record<string, number>; // external_category -> internal_category_id
  selectedAttributes: string[];
  approvedVariantGroups: string[]; // base_names
  categoryProposals: CategoryProposal[];
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

// ============================================
// Enhanced Import Analysis Types (Phase 1)
// ============================================

/**
 * Confidence level for AI mapping
 */
export type MappingConfidence = 'high' | 'medium' | 'low';

/**
 * Category mapping suggestion from AI
 */
export interface CategoryMapping {
  external_category: string; // Original category from import file
  suggested_internal_category_id: number | null; // Suggested marketplace category ID
  suggested_internal_category_name?: string; // Category name (for display)
  confidence: MappingConfidence;
  reasoning?: string; // AI explanation why this mapping was suggested
  is_approved: boolean; // User approved this mapping
  requires_new_category: boolean; // Needs admin to create new category
}

/**
 * Response from category analysis endpoint
 */
export interface CategoryAnalysisResponse {
  total_categories: number;
  mappings: CategoryMapping[];
  quality_summary: {
    high_confidence: number;
    medium_confidence: number;
    low_confidence: number;
    requires_new_category: number;
  };
  unmapped_categories: string[]; // Categories with no suggestions
}

/**
 * Attribute detected in import file
 */
export interface DetectedAttribute {
  name: string; // Original attribute name from file
  value_type: 'string' | 'number' | 'boolean' | 'enum'; // Detected type
  sample_values: string[]; // Example values
  frequency: number; // How many products have this attribute
  suggested_mapping?: string; // Suggested internal attribute name
  is_variant_defining: boolean; // Could be used for variants (color, size, etc.)
}

/**
 * Response from attribute analysis endpoint
 */
export interface AttributeAnalysisResponse {
  total_attributes: number;
  attributes: DetectedAttribute[];
  variant_defining_attributes: string[]; // Attributes that could define variants
}

/**
 * Detected product variant group
 */
export interface VariantGroup {
  base_name: string; // Product base name without variant info
  variant_count: number; // Number of variants in this group
  variant_attributes: string[]; // Attributes that define variants (color, size, etc.)
  products: VariantProduct[]; // Products in this group
  confidence: MappingConfidence; // How confident we are this is a variant group
}

/**
 * Product that is part of variant group
 */
export interface VariantProduct {
  sku: string;
  name: string;
  variant_values: Record<string, string>; // e.g., { color: "Red", size: "M" }
  price: number;
  images?: string[];
}

/**
 * Response from variant detection endpoint
 */
export interface VariantDetectionResponse {
  total_groups: number;
  variant_groups: VariantGroup[];
  ungrouped_products: number; // Products that don't belong to any group
  confidence_summary: {
    high: number;
    medium: number;
    low: number;
  };
}

/**
 * Proposal for new category (admin review required)
 */
export interface CategoryProposal {
  id?: number;
  external_category: string;
  suggested_name: string;
  suggested_parent_id: number | null;
  reasoning: string;
  expected_products: number;
  status: 'pending' | 'approved' | 'rejected';
  created_at?: string;
  reviewed_at?: string;
}

/**
 * Response from client categories analysis
 */
export interface ClientCategoriesResponse {
  unique_categories: string[];
  category_tree: CategoryTreeNode[];
  total_products_per_category: Record<string, number>;
  suggested_proposals: CategoryProposal[];
}

/**
 * Category tree node for hierarchical display
 */
export interface CategoryTreeNode {
  name: string;
  level: number; // 1, 2, or 3
  children?: CategoryTreeNode[];
  product_count: number;
}

/**
 * Import analysis state (for wizard)
 */
export interface ImportAnalysisState {
  // Current step in wizard
  currentStep: number;

  // File info
  file: File | null;
  fileType: 'xml' | 'csv' | 'zip' | '';

  // Analysis results
  categoryAnalysis: CategoryAnalysisResponse | null;
  attributeAnalysis: AttributeAnalysisResponse | null;
  variantDetection: VariantDetectionResponse | null;
  clientCategories: ClientCategoriesResponse | null;

  // User modifications
  approvedMappings: CategoryMapping[];
  customMappings: Record<string, number>; // external_category -> internal_category_id
  selectedAttributes: string[]; // Attributes user wants to import
  approvedVariantGroups: string[]; // base_names of approved groups

  // Loading states
  isAnalyzing: boolean;
  analysisError: string | null;

  // Progress
  analysisProgress: number; // 0-100
}
