/**
 * Category Proposals API Client
 * All requests go through BFF proxy /api/v2 (which maps to backend /api/v1)
 */

import { apiClient } from './api-client';
import type {
  CategoryProposal,
  CategoryProposalListResponse,
  CategoryProposalFilter,
  CreateCategoryProposalRequest,
  UpdateCategoryProposalRequest,
  CategoryProposalApproveRequest,
  CategoryProposalRejectRequest,
  CategoryProposalApproveResponse,
} from '@/types/categoryProposals';

export class CategoryProposalsApi {
  /**
   * List category proposals with filters and pagination
   */
  static async listProposals(
    filter?: CategoryProposalFilter
  ): Promise<CategoryProposalListResponse> {
    const params = new URLSearchParams();

    if (filter?.status) {
      params.append('status', filter.status);
    }
    if (filter?.storefront_id) {
      params.append('storefront_id', filter.storefront_id.toString());
    }
    if (filter?.page) {
      params.append('page', filter.page.toString());
    }
    if (filter?.page_size) {
      params.append('page_size', filter.page_size.toString());
    }

    const url = `/admin/category-proposals${params.toString() ? '?' + params.toString() : ''}`;
    const response = await apiClient.get(url);

    return response.data;
  }

  /**
   * Get pending proposals count
   */
  static async getPendingCount(storefrontId?: number): Promise<number> {
    const params = new URLSearchParams();
    if (storefrontId) {
      params.append('storefront_id', storefrontId.toString());
    }

    const url = `/admin/category-proposals/pending/count${params.toString() ? '?' + params.toString() : ''}`;
    const response = await apiClient.get(url);

    return response.data?.count || 0;
  }

  /**
   * Get single proposal by ID
   */
  static async getProposal(id: number): Promise<CategoryProposal> {
    const response = await apiClient.get(`/admin/category-proposals/${id}`);
    return response.data;
  }

  /**
   * Create new category proposal
   */
  static async createProposal(
    request: CreateCategoryProposalRequest
  ): Promise<CategoryProposal> {
    const response = await apiClient.post('/admin/category-proposals', request);
    return response.data;
  }

  /**
   * Update pending proposal
   */
  static async updateProposal(
    id: number,
    request: UpdateCategoryProposalRequest
  ): Promise<CategoryProposal> {
    const response = await apiClient.put(`/admin/category-proposals/${id}`, request);
    return response.data;
  }

  /**
   * Approve proposal (optionally create category)
   */
  static async approveProposal(
    id: number,
    request: CategoryProposalApproveRequest
  ): Promise<CategoryProposalApproveResponse> {
    const response = await apiClient.post(
      `/admin/category-proposals/${id}/approve`,
      request
    );
    return response.data;
  }

  /**
   * Reject proposal with optional reason
   */
  static async rejectProposal(
    id: number,
    reason?: string
  ): Promise<CategoryProposal> {
    const response = await apiClient.post(
      `/admin/category-proposals/${id}/reject`,
      { reason } as CategoryProposalRejectRequest
    );
    return response.data;
  }

  /**
   * Delete proposal
   */
  static async deleteProposal(id: number): Promise<void> {
    await apiClient.delete(`/admin/category-proposals/${id}`);
  }
}
