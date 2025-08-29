import { apiClient } from './api-client';
import { tokenManager } from '@/utils/tokenManager';
import configManager from '@/config';
import type {
  ImportJob,
  ImportJobStatus,
  ImportRequest,
  ImportSummary,
  ImportFormats,
  UploadProgress,
} from '@/types/import';

export class ImportApi {
  /**
   * Imports products from a URL using storefront slug
   */
  static async importFromUrlBySlug(
    storefrontSlug: string,
    url: string,
    options: {
      file_type?: 'xml' | 'csv' | 'zip';
      update_mode?: 'create_only' | 'update_only' | 'upsert';
      category_mapping_mode?: 'auto' | 'manual' | 'skip';
    }
  ): Promise<ImportJob> {
    const response = await apiClient.post(
      `/api/v1/storefronts/slug/${storefrontSlug}/import/url`,
      {
        url,
        ...options,
      }
    );
    return response.data;
  }

  /**
   * Uploads and imports products from a file using storefront slug
   */
  static async importFromFileBySlug(
    storefrontSlug: string,
    file: File,
    options: {
      file_type: 'xml' | 'csv' | 'zip';
      update_mode?: 'create_only' | 'update_only' | 'upsert';
      category_mapping_mode?: 'auto' | 'manual' | 'skip';
    },
    onProgress?: (progress: UploadProgress) => void
  ): Promise<ImportJob> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('file_type', options.file_type);
    formData.append('update_mode', options.update_mode || 'upsert');
    formData.append(
      'category_mapping_mode',
      options.category_mapping_mode || 'auto'
    );

    // Use fetch directly for file upload with progress tracking
    const xhr = new XMLHttpRequest();

    return new Promise((resolve, reject) => {
      xhr.upload.addEventListener('progress', (event) => {
        if (onProgress && event.lengthComputable) {
          const progress: UploadProgress = {
            loaded: event.loaded,
            total: event.total,
            percentage: Math.round((event.loaded * 100) / event.total),
          };
          onProgress(progress);
        }
      });

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const data = JSON.parse(xhr.responseText);
            resolve(data);
          } catch {
            reject(new Error('Invalid JSON response'));
          }
        } else {
          reject(new Error(`HTTP Error: ${xhr.status}`));
        }
      });

      xhr.addEventListener('error', () => {
        reject(new Error('Network error'));
      });

      xhr.open(
        'POST',
        `${configManager.getApiUrl()}/api/v1/storefronts/slug/${storefrontSlug}/import/file`
      );
      xhr.withCredentials = true; // Include cookies

      // Add authorization header if token exists
      const accessToken = tokenManager.getAccessToken();
      if (accessToken) {
        xhr.setRequestHeader('Authorization', `Bearer ${accessToken}`);
      }

      xhr.send(formData);
    });
  }
  /**
   * Imports products from a URL
   */
  static async importFromUrl(
    storefrontId: number,
    request: Omit<ImportRequest, 'storefront_id'>
  ): Promise<ImportJob> {
    const response = await apiClient.post(
      `/api/v1/storefronts/${storefrontId}/import/url`,
      {
        ...request,
        storefront_id: storefrontId,
      }
    );
    return response.data;
  }

  /**
   * Uploads and imports products from a file
   */
  static async importFromFile(
    storefrontId: number,
    file: File,
    options: {
      file_type: 'xml' | 'csv' | 'zip';
      update_mode?: 'create_only' | 'update_only' | 'upsert';
      category_mapping_mode?: 'auto' | 'manual' | 'skip';
    },
    onProgress?: (progress: UploadProgress) => void
  ): Promise<ImportJob> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('file_type', options.file_type);
    formData.append('update_mode', options.update_mode || 'upsert');
    formData.append(
      'category_mapping_mode',
      options.category_mapping_mode || 'auto'
    );

    // Use fetch directly for file upload with progress tracking
    const xhr = new XMLHttpRequest();

    return new Promise((resolve, reject) => {
      xhr.upload.addEventListener('progress', (event) => {
        if (onProgress && event.lengthComputable) {
          const progress: UploadProgress = {
            loaded: event.loaded,
            total: event.total,
            percentage: Math.round((event.loaded * 100) / event.total),
          };
          onProgress(progress);
        }
      });

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const data = JSON.parse(xhr.responseText);
            resolve(data);
          } catch {
            reject(new Error('Invalid JSON response'));
          }
        } else {
          reject(new Error(`HTTP Error: ${xhr.status}`));
        }
      });

      xhr.addEventListener('error', () => {
        reject(new Error('Network error'));
      });

      xhr.open(
        'POST',
        `${configManager.getApiUrl()}/api/v1/storefronts/${storefrontId}/import/file`
      );
      xhr.withCredentials = true; // Include cookies

      // Add authorization header if token exists
      const accessToken = tokenManager.getAccessToken();
      if (accessToken) {
        xhr.setRequestHeader('Authorization', `Bearer ${accessToken}`);
      }

      xhr.send(formData);
    });
  }

  /**
   * Validates import file without importing
   */
  static async validateFile(
    storefrontId: number,
    file: File,
    fileType: 'xml' | 'csv' | 'zip'
  ): Promise<ImportJobStatus> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('file_type', fileType);

    const response = await apiClient.post(
      `/api/v1/storefronts/${storefrontId}/import/validate`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    );

    return response.data;
  }

  /**
   * Gets import job status
   */
  static async getJobStatus(jobId: number): Promise<ImportJobStatus> {
    const response = await apiClient.get(`/api/v1/import/jobs/${jobId}/status`);
    return response.data;
  }

  /**
   * Gets list of import jobs for a storefront
   */
  static async getJobs(
    storefrontId: number,
    params?: {
      status?: string;
      limit?: number;
      offset?: number;
    }
  ): Promise<{ jobs: ImportJob[]; total: number }> {
    const queryParams = new URLSearchParams();
    if (params?.status) queryParams.append('status', params.status);
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const queryString = queryParams.toString();
    const url = `/api/v1/storefronts/${storefrontId}/import/jobs${queryString ? `?${queryString}` : ''}`;

    const response = await apiClient.get(url);
    return response.data;
  }

  /**
   * Gets detailed job information including errors
   */
  static async getJobDetails(
    jobId: number
  ): Promise<ImportJob & { errors?: any[] }> {
    const response = await apiClient.get(`/api/v1/import/jobs/${jobId}`);
    return response.data;
  }

  /**
   * Downloads CSV template for product import
   */
  static async downloadCsvTemplate(): Promise<Blob> {
    const accessToken = tokenManager.getAccessToken();
    const headers: Record<string, string> = {};
    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
    }

    const response = await fetch(
      `${configManager.getApiUrl()}/api/v1/storefronts/import/csv-template`,
      {
        headers,
        credentials: 'include',
      }
    );
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.blob();
  }

  /**
   * Gets supported import formats information
   */
  static async getFormats(): Promise<ImportFormats> {
    const response = await apiClient.get('/api/v1/storefronts/import/formats');
    return response.data;
  }

  /**
   * Cancels an ongoing import job
   */
  static async cancelJob(jobId: number): Promise<void> {
    await apiClient.post(`/api/v1/import/jobs/${jobId}/cancel`);
  }

  /**
   * Retries a failed import job
   */
  static async retryJob(jobId: number): Promise<ImportJob> {
    const response = await apiClient.post(`/api/v1/import/jobs/${jobId}/retry`);
    return response.data;
  }

  /**
   * Exports import results as CSV
   */
  static async exportResults(jobId: number): Promise<Blob> {
    const accessToken = tokenManager.getAccessToken();
    const headers: Record<string, string> = {};
    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
    }

    const response = await fetch(
      `${configManager.getApiUrl()}/api/v1/import/jobs/${jobId}/export`,
      {
        headers,
        credentials: 'include',
      }
    );
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.blob();
  }

  /**
   * Gets import summary statistics
   */
  static async getSummary(
    storefrontId: number,
    params?: {
      start_date?: string;
      end_date?: string;
    }
  ): Promise<ImportSummary[]> {
    const queryParams = new URLSearchParams();
    if (params?.start_date) queryParams.append('start_date', params.start_date);
    if (params?.end_date) queryParams.append('end_date', params.end_date);

    const queryString = queryParams.toString();
    const url = `/api/v1/storefronts/${storefrontId}/import/summary${queryString ? `?${queryString}` : ''}`;

    const response = await apiClient.get(url);
    return response.data;
  }

  /**
   * Uploads file and gets preview of data without importing
   */
  static async previewFile(
    file: File,
    fileType: 'xml' | 'csv' | 'zip'
  ): Promise<{
    sample_data: any[];
    total_records: number;
    detected_fields: string[];
    validation_errors: any[];
  }> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('file_type', fileType);

    const response = await apiClient.post(
      '/api/v1/storefronts/import/preview',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    );

    return response.data;
  }

  /**
   * Gets category mapping suggestions for import
   */
  static async getCategoryMappings(
    storefrontId: number,
    importCategories: string[]
  ): Promise<{
    mappings: Array<{
      import_category: string;
      suggested_category_id: number;
      confidence: number;
    }>;
  }> {
    const response = await apiClient.post(
      `/api/v1/storefronts/${storefrontId}/import/category-mappings`,
      {
        categories: importCategories,
      }
    );
    return response.data;
  }

  /**
   * Creates custom category mapping
   */
  static async createCategoryMapping(
    storefrontId: number,
    mapping: {
      import_category1: string;
      import_category2?: string;
      import_category3?: string;
      local_category_id: number;
    }
  ): Promise<void> {
    await apiClient.post(
      `/api/v1/storefronts/${storefrontId}/import/category-mappings/create`,
      mapping
    );
  }

  /**
   * Downloads sample import file for given format
   */
  static async downloadSample(format: 'csv' | 'xml'): Promise<Blob> {
    const accessToken = tokenManager.getAccessToken();
    const headers: Record<string, string> = {};
    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
    }

    const response = await fetch(
      `${configManager.getApiUrl()}/api/v1/storefronts/import/sample/${format}`,
      {
        headers,
        credentials: 'include',
      }
    );
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.blob();
  }
}

// Helper functions for file handling
export const downloadFile = (blob: Blob, filename: string) => {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
};

export const validateFileType = (
  file: File,
  allowedTypes: string[]
): boolean => {
  return allowedTypes.includes(file.type);
};

export const validateFileSize = (file: File, maxSize: number): boolean => {
  return file.size <= maxSize;
};

export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

export const getFileTypeFromExtension = (
  filename: string
): 'xml' | 'csv' | 'zip' | '' => {
  const extension = filename.split('.').pop()?.toLowerCase();
  switch (extension) {
    case 'xml':
      return 'xml';
    case 'csv':
      return 'csv';
    case 'zip':
      return 'zip';
    default:
      return '';
  }
};
