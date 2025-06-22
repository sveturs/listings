import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { ImportApi } from '@/services/importApi';
import type {
  ImportState,
  ImportJob,
  ImportJobStatus,
  UploadProgress,
  ImportRequest,
} from '@/types/import';

// Initial state
const initialState: ImportState = {
  jobs: [],
  currentJob: null,
  isUploading: false,
  uploadProgress: null,
  validationErrors: [],
  formats: null,

  // UI states
  isLoading: false,
  error: null,
  selectedFiles: [],
  importUrl: '',
  selectedFileType: '',
  updateMode: 'upsert',
  categoryMappingMode: 'auto',

  // Modal states
  isImportModalOpen: false,
  isJobDetailsModalOpen: false,
  isErrorsModalOpen: false,
};

// Async thunks
export const fetchImportFormats = createAsyncThunk(
  'import/fetchFormats',
  async () => {
    return await ImportApi.getFormats();
  }
);

export const fetchImportJobs = createAsyncThunk(
  'import/fetchJobs',
  async (params: { storefrontId: number; status?: string; limit?: number }) => {
    return await ImportApi.getJobs(params.storefrontId, params);
  }
);

export const importFromUrl = createAsyncThunk(
  'import/importFromUrl',
  async (params: {
    storefrontId: number;
    storefrontSlug?: string;
    request: Omit<ImportRequest, 'storefront_id'>;
  }) => {
    // Use slug-based API if slug is provided
    if (params.storefrontSlug && params.request.file_url) {
      return await ImportApi.importFromUrlBySlug(
        params.storefrontSlug,
        params.request.file_url,
        {
          file_type: params.request.file_type,
          update_mode: params.request.update_mode,
          category_mapping_mode: params.request.category_mapping_mode,
        }
      );
    }
    if (!params.request.file_url) {
      throw new Error('file_url is required for URL import');
    }
    return await ImportApi.importFromUrl(params.storefrontId, params.request);
  }
);

export const importFromFile = createAsyncThunk(
  'import/importFromFile',
  async (
    params: {
      storefrontId: number;
      storefrontSlug?: string;
      file: File;
      options: {
        file_type: 'xml' | 'csv' | 'zip';
        update_mode?: 'create_only' | 'update_only' | 'upsert';
        category_mapping_mode?: 'auto' | 'manual' | 'skip';
      };
    },
    { dispatch }
  ) => {
    dispatch(setIsUploading(true));
    dispatch(setUploadProgress({ loaded: 0, total: 0, percentage: 0 }));

    try {
      // Use slug-based API if slug is provided
      const result = params.storefrontSlug
        ? await ImportApi.importFromFileBySlug(
            params.storefrontSlug,
            params.file,
            params.options,
            (progress) => {
              dispatch(setUploadProgress(progress));
            }
          )
        : await ImportApi.importFromFile(
            params.storefrontId,
            params.file,
            params.options,
            (progress) => {
              dispatch(setUploadProgress(progress));
            }
          );
      return result;
    } finally {
      dispatch(setIsUploading(false));
      dispatch(setUploadProgress(null));
    }
  }
);

export const validateImportFile = createAsyncThunk(
  'import/validateFile',
  async (params: {
    storefrontId: number;
    file: File;
    fileType: 'xml' | 'csv' | 'zip';
  }) => {
    return await ImportApi.validateFile(
      params.storefrontId,
      params.file,
      params.fileType
    );
  }
);

export const fetchJobStatus = createAsyncThunk(
  'import/fetchJobStatus',
  async (jobId: number) => {
    return await ImportApi.getJobStatus(jobId);
  }
);

export const fetchJobDetails = createAsyncThunk(
  'import/fetchJobDetails',
  async (jobId: number) => {
    return await ImportApi.getJobDetails(jobId);
  }
);

export const downloadCsvTemplate = createAsyncThunk(
  'import/downloadCsvTemplate',
  async () => {
    const blob = await ImportApi.downloadCsvTemplate();
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = 'product_import_template.csv';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  }
);

export const cancelImportJob = createAsyncThunk(
  'import/cancelJob',
  async (jobId: number) => {
    await ImportApi.cancelJob(jobId);
    return jobId;
  }
);

export const retryImportJob = createAsyncThunk(
  'import/retryJob',
  async (jobId: number) => {
    return await ImportApi.retryJob(jobId);
  }
);

// Slice
const importSlice = createSlice({
  name: 'import',
  initialState,
  reducers: {
    // UI state reducers
    setIsUploading: (state, action: PayloadAction<boolean>) => {
      state.isUploading = action.payload;
    },
    setUploadProgress: (
      state,
      action: PayloadAction<UploadProgress | null>
    ) => {
      state.uploadProgress = action.payload;
    },
    setSelectedFiles: (state, action: PayloadAction<File[]>) => {
      state.selectedFiles = action.payload;
    },
    setImportUrl: (state, action: PayloadAction<string>) => {
      state.importUrl = action.payload;
    },
    setSelectedFileType: (
      state,
      action: PayloadAction<'xml' | 'csv' | 'zip' | ''>
    ) => {
      state.selectedFileType = action.payload;
    },
    setUpdateMode: (
      state,
      action: PayloadAction<'create_only' | 'update_only' | 'upsert'>
    ) => {
      state.updateMode = action.payload;
    },
    setCategoryMappingMode: (
      state,
      action: PayloadAction<'auto' | 'manual' | 'skip'>
    ) => {
      state.categoryMappingMode = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    clearError: (state) => {
      state.error = null;
    },

    // Modal state reducers
    setImportModalOpen: (state, action: PayloadAction<boolean>) => {
      state.isImportModalOpen = action.payload;
    },
    setJobDetailsModalOpen: (state, action: PayloadAction<boolean>) => {
      state.isJobDetailsModalOpen = action.payload;
    },
    setErrorsModalOpen: (state, action: PayloadAction<boolean>) => {
      state.isErrorsModalOpen = action.payload;
    },

    // Current job
    setCurrentJob: (state, action: PayloadAction<ImportJob | null>) => {
      state.currentJob = action.payload;
    },

    // Update job status in the list
    updateJobStatus: (
      state,
      action: PayloadAction<{ jobId: number; status: ImportJobStatus }>
    ) => {
      const { jobId, status } = action.payload;
      const job = state.jobs.find((j) => j.id === jobId);
      if (job) {
        job.status = status.status as any;
        job.processed_records = status.processed_records;
        job.successful_records = status.successful_records;
        job.failed_records = status.failed_records;
        job.total_records = status.total_records;
      }
    },

    // Reset form
    resetForm: (state) => {
      state.selectedFiles = [];
      state.importUrl = '';
      state.selectedFileType = '';
      state.updateMode = 'upsert';
      state.categoryMappingMode = 'auto';
      state.uploadProgress = null;
      state.validationErrors = [];
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch formats
      .addCase(fetchImportFormats.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchImportFormats.fulfilled, (state, action) => {
        state.isLoading = false;
        state.formats = action.payload;
      })
      .addCase(fetchImportFormats.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to fetch import formats';
      })

      // Fetch jobs
      .addCase(fetchImportJobs.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchImportJobs.fulfilled, (state, action) => {
        state.isLoading = false;
        if (action.payload && action.payload.jobs) {
          state.jobs = action.payload.jobs;
        } else {
          state.jobs = [];
        }
      })
      .addCase(fetchImportJobs.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to fetch import jobs';
      })

      // Import from URL
      .addCase(importFromUrl.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(importFromUrl.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentJob = action.payload;
        state.jobs.unshift(action.payload);
        state.isImportModalOpen = false;
        // Reset form after successful import
        importSlice.caseReducers.resetForm(state);
      })
      .addCase(importFromUrl.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to import from URL';
      })

      // Import from file
      .addCase(importFromFile.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(importFromFile.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentJob = action.payload;
        state.jobs.unshift(action.payload);
        state.isImportModalOpen = false;
        // Reset form after successful import
        importSlice.caseReducers.resetForm(state);
      })
      .addCase(importFromFile.rejected, (state, action) => {
        state.isLoading = false;
        state.isUploading = false;
        state.uploadProgress = null;
        state.error = action.error.message || 'Failed to import file';
      })

      // Validate file
      .addCase(validateImportFile.pending, (state) => {
        state.isLoading = true;
        state.error = null;
        state.validationErrors = [];
      })
      .addCase(validateImportFile.fulfilled, (state, action) => {
        state.isLoading = false;
        // Handle validation results
        if (action.payload.errors && action.payload.errors.length > 0) {
          state.validationErrors = action.payload.errors as any;
        }
      })
      .addCase(validateImportFile.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to validate file';
      })

      // Fetch job status
      .addCase(fetchJobStatus.fulfilled, (state, action) => {
        importSlice.caseReducers.updateJobStatus(state, {
          payload: { jobId: action.payload.id, status: action.payload },
          type: 'updateJobStatus',
        });
      })

      // Fetch job details
      .addCase(fetchJobDetails.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchJobDetails.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentJob = action.payload as ImportJob;
      })
      .addCase(fetchJobDetails.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to fetch job details';
      })

      // Cancel job
      .addCase(cancelImportJob.fulfilled, (state, action) => {
        const job = state.jobs.find((j) => j.id === action.payload);
        if (job) {
          job.status = 'failed';
        }
      })

      // Retry job
      .addCase(retryImportJob.fulfilled, (state, action) => {
        const jobIndex = state.jobs.findIndex(
          (j) => j.id === action.payload.id
        );
        if (jobIndex !== -1) {
          state.jobs[jobIndex] = action.payload;
        } else {
          state.jobs.unshift(action.payload);
        }
        state.currentJob = action.payload;
      });
  },
});

export const {
  setIsUploading,
  setUploadProgress,
  setSelectedFiles,
  setImportUrl,
  setSelectedFileType,
  setUpdateMode,
  setCategoryMappingMode,
  setError,
  clearError,
  setImportModalOpen,
  setJobDetailsModalOpen,
  setErrorsModalOpen,
  setCurrentJob,
  updateJobStatus,
  resetForm,
} = importSlice.actions;

export default importSlice.reducer;
