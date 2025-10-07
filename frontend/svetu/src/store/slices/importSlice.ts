import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { ImportApi } from '@/services/importApi';
import type {
  ImportState,
  ImportJob,
  ImportJobStatus,
  UploadProgress,
  ImportRequest,
  CategoryAnalysisResponse,
  AttributeAnalysisResponse,
  VariantDetectionResponse,
  ClientCategoriesResponse,
  CategoryMapping,
  CategoryProposal,
} from '@/types/import';

// Initial state
const initialState: ImportState = {
  jobs: [],
  currentJob: null,
  isUploading: false,
  uploadProgress: null,
  validationErrors: [],
  formats: null,

  // Preview states
  previewData: null,
  isPreviewLoading: false,
  previewError: null,

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

  // Analysis states (Enhanced Import)
  analysisFile: null,
  analysisFileType: '',
  categoryAnalysis: null,
  attributeAnalysis: null,
  variantDetection: null,
  clientCategories: null,
  isAnalyzing: false,
  analysisError: null,
  analysisProgress: 0,
  currentAnalysisStep: 0,

  // User modifications for analysis
  approvedMappings: [],
  customMappings: {},
  selectedAttributes: [],
  approvedVariantGroups: [],
  categoryProposals: [],
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

export const previewImportFile = createAsyncThunk(
  'import/previewFile',
  async (params: {
    storefrontId?: number;
    storefrontSlug?: string;
    file: File;
    fileType: 'xml' | 'csv' | 'zip';
    previewLimit?: number;
  }) => {
    // Use slug-based API if slug is provided
    if (params.storefrontSlug) {
      return await ImportApi.previewFileBySlug(
        params.storefrontSlug,
        params.file,
        params.fileType,
        params.previewLimit
      );
    }

    if (!params.storefrontId) {
      throw new Error('Either storefrontId or storefrontSlug is required');
    }

    return await ImportApi.previewFile(
      params.storefrontId,
      params.file,
      params.fileType,
      params.previewLimit
    );
  }
);

export const fetchJobStatus = createAsyncThunk(
  'import/fetchJobStatus',
  async (params: { storefrontId: number; jobId: number }) => {
    return await ImportApi.getJobStatus(params.storefrontId, params.jobId);
  }
);

export const fetchJobDetails = createAsyncThunk(
  'import/fetchJobDetails',
  async (params: { storefrontId: number; jobId: number }) => {
    return await ImportApi.getJobDetails(params.storefrontId, params.jobId);
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

// Enhanced Import Analysis thunks
export const analyzeImportFile = createAsyncThunk(
  'import/analyzeFile',
  async (
    params: {
      storefrontId: number;
      file: File;
      fileType: 'xml' | 'csv' | 'zip';
    },
    { dispatch }
  ) => {
    // Step 1: Analyze categories (progress 0-33%)
    dispatch(setAnalysisProgress(10));
    const categoryAnalysis = await ImportApi.analyzeCategories(
      params.storefrontId,
      params.file,
      params.fileType
    );
    dispatch(setAnalysisProgress(33));

    // Step 2: Analyze attributes (progress 33-66%)
    const attributeAnalysis = await ImportApi.analyzeAttributes(
      params.storefrontId,
      params.file,
      params.fileType
    );
    dispatch(setAnalysisProgress(66));

    // Step 3: Detect variants (progress 66-100%)
    const variantDetection = await ImportApi.detectVariants(
      params.storefrontId,
      params.file,
      params.fileType
    );
    dispatch(setAnalysisProgress(100));

    return {
      categoryAnalysis,
      attributeAnalysis,
      variantDetection,
    };
  }
);

export const analyzeCategories = createAsyncThunk(
  'import/analyzeCategories',
  async (params: {
    storefrontId: number;
    file: File;
    fileType: 'xml' | 'csv' | 'zip';
  }) => {
    return await ImportApi.analyzeCategories(
      params.storefrontId,
      params.file,
      params.fileType
    );
  }
);

export const analyzeAttributes = createAsyncThunk(
  'import/analyzeAttributes',
  async (params: {
    storefrontId: number;
    file: File;
    fileType: 'xml' | 'csv' | 'zip';
  }) => {
    return await ImportApi.analyzeAttributes(
      params.storefrontId,
      params.file,
      params.fileType
    );
  }
);

export const detectVariants = createAsyncThunk(
  'import/detectVariants',
  async (params: {
    storefrontId: number;
    file: File;
    fileType: 'xml' | 'csv' | 'zip';
  }) => {
    return await ImportApi.detectVariants(
      params.storefrontId,
      params.file,
      params.fileType
    );
  }
);

export const analyzeClientCategories = createAsyncThunk(
  'import/analyzeClientCategories',
  async (params: {
    storefrontId: number;
    file: File;
    fileType: 'xml' | 'csv' | 'zip';
  }) => {
    return await ImportApi.analyzeClientCategories(
      params.storefrontId,
      params.file,
      params.fileType
    );
  }
);

export const fetchCategoryProposals = createAsyncThunk(
  'import/fetchCategoryProposals',
  async (params?: {
    status?: 'pending' | 'approved' | 'rejected';
    storefront_id?: number;
    limit?: number;
    offset?: number;
  }) => {
    return await ImportApi.getCategoryProposals(params);
  }
);

export const approveCategoryProposal = createAsyncThunk(
  'import/approveCategoryProposal',
  async (proposalId: number) => {
    await ImportApi.approveCategoryProposal(proposalId);
    return proposalId;
  }
);

export const rejectCategoryProposal = createAsyncThunk(
  'import/rejectCategoryProposal',
  async (params: { proposalId: number; reason?: string }) => {
    await ImportApi.rejectCategoryProposal(params.proposalId, params.reason);
    return params.proposalId;
  }
);

export const cancelImportJob = createAsyncThunk(
  'import/cancelJob',
  async (params: { storefrontId: number; jobId: number }) => {
    await ImportApi.cancelJob(params.storefrontId, params.jobId);
    return params.jobId;
  }
);

export const retryImportJob = createAsyncThunk(
  'import/retryJob',
  async (params: { storefrontId: number; jobId: number }) => {
    return await ImportApi.retryJob(params.storefrontId, params.jobId);
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
      state.previewData = null;
      state.previewError = null;
    },

    // Preview reducers
    clearPreview: (state) => {
      state.previewData = null;
      state.previewError = null;
    },

    // Analysis reducers
    setAnalysisFile: (state, action: PayloadAction<File | null>) => {
      state.analysisFile = action.payload;
    },
    setAnalysisFileType: (
      state,
      action: PayloadAction<'xml' | 'csv' | 'zip' | ''>
    ) => {
      state.analysisFileType = action.payload;
    },
    setAnalysisProgress: (state, action: PayloadAction<number>) => {
      state.analysisProgress = action.payload;
    },
    setCurrentAnalysisStep: (state, action: PayloadAction<number>) => {
      state.currentAnalysisStep = action.payload;
    },
    setApprovedMappings: (state, action: PayloadAction<CategoryMapping[]>) => {
      state.approvedMappings = action.payload;
    },
    addApprovedMapping: (state, action: PayloadAction<CategoryMapping>) => {
      const index = state.approvedMappings.findIndex(
        (m) => m.external_category === action.payload.external_category
      );
      if (index >= 0) {
        state.approvedMappings[index] = action.payload;
      } else {
        state.approvedMappings.push(action.payload);
      }
    },
    setCustomMapping: (
      state,
      action: PayloadAction<{ externalCategory: string; internalCategoryId: number }>
    ) => {
      state.customMappings[action.payload.externalCategory] =
        action.payload.internalCategoryId;
    },
    setSelectedAttributes: (state, action: PayloadAction<string[]>) => {
      state.selectedAttributes = action.payload;
    },
    toggleSelectedAttribute: (state, action: PayloadAction<string>) => {
      const index = state.selectedAttributes.indexOf(action.payload);
      if (index >= 0) {
        state.selectedAttributes.splice(index, 1);
      } else {
        state.selectedAttributes.push(action.payload);
      }
    },
    setApprovedVariantGroups: (state, action: PayloadAction<string[]>) => {
      state.approvedVariantGroups = action.payload;
    },
    toggleApprovedVariantGroup: (state, action: PayloadAction<string>) => {
      const index = state.approvedVariantGroups.indexOf(action.payload);
      if (index >= 0) {
        state.approvedVariantGroups.splice(index, 1);
      } else {
        state.approvedVariantGroups.push(action.payload);
      }
    },
    clearAnalysis: (state) => {
      state.analysisFile = null;
      state.analysisFileType = '';
      state.categoryAnalysis = null;
      state.attributeAnalysis = null;
      state.variantDetection = null;
      state.clientCategories = null;
      state.isAnalyzing = false;
      state.analysisError = null;
      state.analysisProgress = 0;
      state.currentAnalysisStep = 0;
      state.approvedMappings = [];
      state.customMappings = {};
      state.selectedAttributes = [];
      state.approvedVariantGroups = [];
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
      })

      // Preview file
      .addCase(previewImportFile.pending, (state) => {
        state.isPreviewLoading = true;
        state.previewError = null;
        state.previewData = null;
      })
      .addCase(previewImportFile.fulfilled, (state, action) => {
        state.isPreviewLoading = false;
        state.previewData = action.payload;
      })
      .addCase(previewImportFile.rejected, (state, action) => {
        state.isPreviewLoading = false;
        state.previewError =
          action.error.message || 'Failed to preview import file';
      })

      // Analyze import file (full analysis)
      .addCase(analyzeImportFile.pending, (state) => {
        state.isAnalyzing = true;
        state.analysisError = null;
        state.analysisProgress = 0;
      })
      .addCase(analyzeImportFile.fulfilled, (state, action) => {
        state.isAnalyzing = false;
        state.categoryAnalysis = action.payload.categoryAnalysis;
        state.attributeAnalysis = action.payload.attributeAnalysis;
        state.variantDetection = action.payload.variantDetection;
        state.analysisProgress = 100;
      })
      .addCase(analyzeImportFile.rejected, (state, action) => {
        state.isAnalyzing = false;
        state.analysisError = action.error.message || 'Failed to analyze import file';
        state.analysisProgress = 0;
      })

      // Analyze categories
      .addCase(analyzeCategories.pending, (state) => {
        state.isAnalyzing = true;
        state.analysisError = null;
      })
      .addCase(analyzeCategories.fulfilled, (state, action) => {
        state.isAnalyzing = false;
        state.categoryAnalysis = action.payload;
      })
      .addCase(analyzeCategories.rejected, (state, action) => {
        state.isAnalyzing = false;
        state.analysisError = action.error.message || 'Failed to analyze categories';
      })

      // Analyze attributes
      .addCase(analyzeAttributes.pending, (state) => {
        state.isAnalyzing = true;
        state.analysisError = null;
      })
      .addCase(analyzeAttributes.fulfilled, (state, action) => {
        state.isAnalyzing = false;
        state.attributeAnalysis = action.payload;
      })
      .addCase(analyzeAttributes.rejected, (state, action) => {
        state.isAnalyzing = false;
        state.analysisError = action.error.message || 'Failed to analyze attributes';
      })

      // Detect variants
      .addCase(detectVariants.pending, (state) => {
        state.isAnalyzing = true;
        state.analysisError = null;
      })
      .addCase(detectVariants.fulfilled, (state, action) => {
        state.isAnalyzing = false;
        state.variantDetection = action.payload;
      })
      .addCase(detectVariants.rejected, (state, action) => {
        state.isAnalyzing = false;
        state.analysisError = action.error.message || 'Failed to detect variants';
      })

      // Analyze client categories
      .addCase(analyzeClientCategories.pending, (state) => {
        state.isAnalyzing = true;
        state.analysisError = null;
      })
      .addCase(analyzeClientCategories.fulfilled, (state, action) => {
        state.isAnalyzing = false;
        state.clientCategories = action.payload;
      })
      .addCase(analyzeClientCategories.rejected, (state, action) => {
        state.isAnalyzing = false;
        state.analysisError =
          action.error.message || 'Failed to analyze client categories';
      })

      // Fetch category proposals
      .addCase(fetchCategoryProposals.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchCategoryProposals.fulfilled, (state, action) => {
        state.isLoading = false;
        state.categoryProposals = action.payload.proposals;
      })
      .addCase(fetchCategoryProposals.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to fetch category proposals';
      })

      // Approve category proposal
      .addCase(approveCategoryProposal.fulfilled, (state, action) => {
        const proposal = state.categoryProposals.find(
          (p) => p.id === action.payload
        );
        if (proposal) {
          proposal.status = 'approved';
        }
      })

      // Reject category proposal
      .addCase(rejectCategoryProposal.fulfilled, (state, action) => {
        const proposal = state.categoryProposals.find(
          (p) => p.id === action.payload
        );
        if (proposal) {
          proposal.status = 'rejected';
        }
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
  clearPreview,
  // Analysis actions
  setAnalysisFile,
  setAnalysisFileType,
  setAnalysisProgress,
  setCurrentAnalysisStep,
  setApprovedMappings,
  addApprovedMapping,
  setCustomMapping,
  setSelectedAttributes,
  toggleSelectedAttribute,
  setApprovedVariantGroups,
  toggleApprovedVariantGroup,
  clearAnalysis,
} = importSlice.actions;

export default importSlice.reducer;
