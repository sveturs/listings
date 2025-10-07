import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { CategoryProposalsApi } from '@/services/categoryProposalsApi';
import type {
  CategoryProposal,
  CategoryProposalStatus,
  CategoryProposalListResponse,
  CategoryProposalFilter,
  CategoryProposalApproveRequest,
  CategoryProposalRejectRequest,
  CategoryProposalApproveResponse,
} from '@/types/categoryProposals';

interface CategoryProposalsState {
  proposals: CategoryProposal[];
  currentProposal: CategoryProposal | null;
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  pendingCount: number;

  // Filters
  statusFilter: CategoryProposalStatus | 'all';
  storefrontIdFilter: number | null;

  // UI states
  isLoading: boolean;
  isApproving: boolean;
  isRejecting: boolean;
  isDeleting: boolean;
  error: string | null;

  // Modal states
  isApproveModalOpen: boolean;
  isRejectModalOpen: boolean;
  selectedProposalId: number | null;
}

const initialState: CategoryProposalsState = {
  proposals: [],
  currentProposal: null,
  total: 0,
  page: 1,
  pageSize: 20,
  totalPages: 0,
  pendingCount: 0,

  statusFilter: 'all',
  storefrontIdFilter: null,

  isLoading: false,
  isApproving: false,
  isRejecting: false,
  isDeleting: false,
  error: null,

  isApproveModalOpen: false,
  isRejectModalOpen: false,
  selectedProposalId: null,
};

// Async thunks
export const fetchCategoryProposals = createAsyncThunk(
  'categoryProposals/fetchList',
  async (filter?: CategoryProposalFilter) => {
    return await CategoryProposalsApi.listProposals(filter);
  }
);

export const fetchPendingCount = createAsyncThunk(
  'categoryProposals/fetchPendingCount',
  async () => {
    return await CategoryProposalsApi.getPendingCount();
  }
);

export const fetchCategoryProposal = createAsyncThunk(
  'categoryProposals/fetchOne',
  async (id: number) => {
    return await CategoryProposalsApi.getProposal(id);
  }
);

export const approveCategoryProposal = createAsyncThunk(
  'categoryProposals/approve',
  async (params: { id: number; request: CategoryProposalApproveRequest }) => {
    return await CategoryProposalsApi.approveProposal(params.id, params.request);
  }
);

export const rejectCategoryProposal = createAsyncThunk(
  'categoryProposals/reject',
  async (params: { id: number; request: CategoryProposalRejectRequest }) => {
    return await CategoryProposalsApi.rejectProposal(params.id, params.request);
  }
);

export const deleteCategoryProposal = createAsyncThunk(
  'categoryProposals/delete',
  async (id: number) => {
    await CategoryProposalsApi.deleteProposal(id);
    return id;
  }
);

const categoryProposalsSlice = createSlice({
  name: 'categoryProposals',
  initialState,
  reducers: {
    setStatusFilter: (state, action: PayloadAction<CategoryProposalStatus | 'all'>) => {
      state.statusFilter = action.payload;
      state.page = 1; // Reset page when filter changes
    },
    setStorefrontIdFilter: (state, action: PayloadAction<number | null>) => {
      state.storefrontIdFilter = action.payload;
      state.page = 1;
    },
    setPage: (state, action: PayloadAction<number>) => {
      state.page = action.payload;
    },
    setPageSize: (state, action: PayloadAction<number>) => {
      state.pageSize = action.payload;
      state.page = 1;
    },
    openApproveModal: (state, action: PayloadAction<number>) => {
      state.isApproveModalOpen = true;
      state.selectedProposalId = action.payload;
    },
    closeApproveModal: (state) => {
      state.isApproveModalOpen = false;
      state.selectedProposalId = null;
    },
    openRejectModal: (state, action: PayloadAction<number>) => {
      state.isRejectModalOpen = true;
      state.selectedProposalId = action.payload;
    },
    closeRejectModal: (state) => {
      state.isRejectModalOpen = false;
      state.selectedProposalId = null;
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    // Fetch list
    builder.addCase(fetchCategoryProposals.pending, (state) => {
      state.isLoading = true;
      state.error = null;
    });
    builder.addCase(fetchCategoryProposals.fulfilled, (state, action) => {
      state.isLoading = false;
      state.proposals = action.payload.proposals;
      state.total = action.payload.total;
      state.page = action.payload.page;
      state.pageSize = action.payload.page_size;
      state.totalPages = action.payload.total_pages;
    });
    builder.addCase(fetchCategoryProposals.rejected, (state, action) => {
      state.isLoading = false;
      state.error = action.error.message || 'Failed to fetch proposals';
    });

    // Fetch pending count
    builder.addCase(fetchPendingCount.fulfilled, (state, action) => {
      state.pendingCount = action.payload.count;
    });

    // Fetch single proposal
    builder.addCase(fetchCategoryProposal.pending, (state) => {
      state.isLoading = true;
      state.error = null;
    });
    builder.addCase(fetchCategoryProposal.fulfilled, (state, action) => {
      state.isLoading = false;
      state.currentProposal = action.payload;
    });
    builder.addCase(fetchCategoryProposal.rejected, (state, action) => {
      state.isLoading = false;
      state.error = action.error.message || 'Failed to fetch proposal';
    });

    // Approve
    builder.addCase(approveCategoryProposal.pending, (state) => {
      state.isApproving = true;
      state.error = null;
    });
    builder.addCase(approveCategoryProposal.fulfilled, (state, action) => {
      state.isApproving = false;
      state.isApproveModalOpen = false;
      state.selectedProposalId = null;

      // Update proposal in list
      const index = state.proposals.findIndex(p => p.id === action.payload.proposal.id);
      if (index !== -1) {
        state.proposals[index] = action.payload.proposal;
      }

      // Update current proposal if it's the same
      if (state.currentProposal?.id === action.payload.proposal.id) {
        state.currentProposal = action.payload.proposal;
      }

      // Decrement pending count if status changed to approved
      if (action.payload.proposal.status === 'approved' && state.pendingCount > 0) {
        state.pendingCount--;
      }
    });
    builder.addCase(approveCategoryProposal.rejected, (state, action) => {
      state.isApproving = false;
      state.error = action.error.message || 'Failed to approve proposal';
    });

    // Reject
    builder.addCase(rejectCategoryProposal.pending, (state) => {
      state.isRejecting = true;
      state.error = null;
    });
    builder.addCase(rejectCategoryProposal.fulfilled, (state, action) => {
      state.isRejecting = false;
      state.isRejectModalOpen = false;
      state.selectedProposalId = null;

      // Update proposal in list
      const index = state.proposals.findIndex(p => p.id === action.payload.id);
      if (index !== -1) {
        state.proposals[index] = action.payload;
      }

      // Update current proposal if it's the same
      if (state.currentProposal?.id === action.payload.id) {
        state.currentProposal = action.payload;
      }

      // Decrement pending count if status changed to rejected
      if (action.payload.status === 'rejected' && state.pendingCount > 0) {
        state.pendingCount--;
      }
    });
    builder.addCase(rejectCategoryProposal.rejected, (state, action) => {
      state.isRejecting = false;
      state.error = action.error.message || 'Failed to reject proposal';
    });

    // Delete
    builder.addCase(deleteCategoryProposal.pending, (state) => {
      state.isDeleting = true;
      state.error = null;
    });
    builder.addCase(deleteCategoryProposal.fulfilled, (state, action) => {
      state.isDeleting = false;

      // Remove from list
      state.proposals = state.proposals.filter(p => p.id !== action.payload);
      state.total--;

      // Clear current if it was deleted
      if (state.currentProposal?.id === action.payload) {
        state.currentProposal = null;
      }

      // Decrement pending count if it was a pending proposal
      const deletedProposal = state.proposals.find(p => p.id === action.payload);
      if (deletedProposal?.status === 'pending' && state.pendingCount > 0) {
        state.pendingCount--;
      }
    });
    builder.addCase(deleteCategoryProposal.rejected, (state, action) => {
      state.isDeleting = false;
      state.error = action.error.message || 'Failed to delete proposal';
    });
  },
});

export const {
  setStatusFilter,
  setStorefrontIdFilter,
  setPage,
  setPageSize,
  openApproveModal,
  closeApproveModal,
  openRejectModal,
  closeRejectModal,
  clearError,
} = categoryProposalsSlice.actions;

export default categoryProposalsSlice.reducer;
