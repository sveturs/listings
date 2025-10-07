/**
 * Category Proposal types matching backend API
 * Backend: backend/internal/domain/models/category_proposal.go
 */

export type CategoryProposalStatus = 'pending' | 'approved' | 'rejected';

export interface NameTranslations {
  en?: string;
  ru?: string;
  sr?: string;
}

export interface CategoryProposal {
  id: number;
  proposed_by_user_id: number;
  storefront_id?: number;
  name: string;
  name_translations: NameTranslations;
  parent_category_id?: number;
  description?: string;
  reasoning?: string; // AI reasoning for this category
  expected_products: number; // Number of products that would use this category
  external_category_source?: string; // Original external category from import
  similar_categories?: number[]; // Related category IDs
  tags?: string[];
  status: CategoryProposalStatus;
  reviewed_by_user_id?: number;
  reviewed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCategoryProposalRequest {
  storefront_id?: number;
  name: string;
  name_translations?: NameTranslations;
  parent_category_id?: number;
  description?: string;
  reasoning?: string;
  expected_products?: number;
  external_category_source?: string;
  similar_categories?: number[];
  tags?: string[];
}

export interface UpdateCategoryProposalRequest {
  name?: string;
  name_translations?: NameTranslations;
  parent_category_id?: number;
  description?: string;
  tags?: string[];
}

export interface CategoryProposalListResponse {
  proposals: CategoryProposal[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface CategoryProposalApproveRequest {
  create_category: boolean; // If true, creates the category in marketplace_categories
}

export interface CategoryProposalRejectRequest {
  reason?: string; // Optional rejection reason
}

export interface CategoryProposalApproveResponse {
  proposal: CategoryProposal;
  category?: {
    id: number;
    name: string;
    slug: string;
    parent_id?: number;
  } | null;
}

export interface CategoryProposalFilter {
  status?: CategoryProposalStatus;
  storefront_id?: number;
  page?: number;
  page_size?: number;
}
