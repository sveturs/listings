package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	categoriespb "github.com/vondi-global/listings/api/proto/categories/v1"
	"github.com/vondi-global/listings/internal/domain"
)

// GetCategories retrieves a list of categories with pagination
func (s *Server) GetCategoriesForCategorySvc(ctx context.Context, req *categoriespb.GetCategoriesRequest) (*categoriespb.GetCategoriesResponse, error) {
	s.logger.Debug().
		Interface("parent_id", req.ParentId).
		Interface("is_active", req.IsActive).
		Int32("page", req.Page).
		Int32("page_size", req.PageSize).
		Msg("GetCategories called")

	// Validate pagination
	if req.Page <= 0 {
		return nil, status.Error(codes.InvalidArgument, "page must be greater than 0")
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		return nil, status.Error(codes.InvalidArgument, "page_size must be between 1 and 100")
	}

	// Convert page-based pagination to offset-based
	offset := (req.Page - 1) * req.PageSize
	limit := req.PageSize

	// Convert proto optional fields to pointers
	var parentID *string
	if req.ParentId != nil {
		parentID = req.ParentId
	}

	var isActive *bool
	if req.IsActive != nil {
		isActive = req.IsActive
	}

	// Call service
	categories, total, err := s.categoryService.GetCategories(ctx, parentID, isActive, limit, offset)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get categories")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get categories: %v", err))
	}

	// Convert to proto
	pbCategories := make([]*categoriespb.Category, len(categories))
	for i, cat := range categories {
		pbCategories[i] = DomainToCategoryServiceProtoCategory(cat)
	}

	s.logger.Debug().
		Int("count", len(categories)).
		Int32("total", total).
		Msg("categories retrieved")

	return &categoriespb.GetCategoriesResponse{
		Categories: pbCategories,
		TotalCount: total,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

// GetCategoryTree retrieves a tree structure of categories (CategoryService implementation)
func (s *Server) GetCategoryTreeForCategorySvc(ctx context.Context, req *categoriespb.GetCategoryTreeRequest) (*categoriespb.GetCategoryTreeResponse, error) {
	s.logger.Debug().
		Interface("root_id", req.RootId).
		Interface("active_only", req.ActiveOnly).
		Msg("GetCategoryTree called")

	// Determine root category ID (empty string means entire tree)
	rootID := ""
	if req.RootId != nil {
		rootID = *req.RootId
	}

	// Get category tree from service
	tree, err := s.categoryService.GetCategoryTree(ctx, rootID)
	if err != nil {
		s.logger.Error().Err(err).Str("root_id", rootID).Msg("failed to get category tree")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get category tree: %v", err))
	}

	// Convert to proto - if root is specified, return single tree, otherwise return forest
	var pbTrees []*categoriespb.CategoryTree
	if tree != nil {
		// Filter by active status if requested
		if req.ActiveOnly != nil && *req.ActiveOnly {
			tree = filterActiveCategoriesInTree(tree)
		}

		if rootID != "" {
			// Single tree for specific root
			pbTrees = []*categoriespb.CategoryTree{DomainToCategoryServiceProtoCategoryTree(tree)}
		} else {
			// Forest of root categories
			if len(tree.Children) > 0 {
				pbTrees = make([]*categoriespb.CategoryTree, len(tree.Children))
				for i := range tree.Children {
					pbTrees[i] = DomainToCategoryServiceProtoCategoryTree(&tree.Children[i])
				}
			} else {
				// Single root category without children
				pbTrees = []*categoriespb.CategoryTree{DomainToCategoryServiceProtoCategoryTree(tree)}
			}
		}
	}

	s.logger.Debug().
		Int("tree_count", len(pbTrees)).
		Msg("category tree retrieved")

	return &categoriespb.GetCategoryTreeResponse{
		Tree: pbTrees,
	}, nil
}

// GetCategory retrieves a single category by ID (CategoryService implementation)
func (s *Server) GetCategoryForCategorySvc(ctx context.Context, req *categoriespb.GetCategoryRequest) (*categoriespb.GetCategoryResponse, error) {
	s.logger.Debug().Str("category_id", req.Id).Msg("GetCategory called")

	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "category ID must not be empty")
	}

	// Get category from service
	category, err := s.categoryService.GetCategory(ctx, req.Id)
	if err != nil {
		s.logger.Error().Err(err).Str("category_id", req.Id).Msg("failed to get category")
		if contains(err.Error(), "not found") || contains(err.Error(), "no rows") {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get category: %v", err))
	}

	// Convert to proto
	pbCategory := DomainToCategoryServiceProtoCategory(category)

	s.logger.Debug().Str("category_id", req.Id).Msg("category retrieved")

	return &categoriespb.GetCategoryResponse{
		Category: pbCategory,
	}, nil
}

// GetCategoryBySlug retrieves a single category by slug
func (s *Server) GetCategoryBySlugForCategorySvc(ctx context.Context, req *categoriespb.GetCategoryBySlugRequest) (*categoriespb.GetCategoryBySlugResponse, error) {
	s.logger.Debug().Str("slug", req.Slug).Msg("GetCategoryBySlug called")

	// Validate request
	if req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "slug is required")
	}

	// Get category from service
	category, err := s.categoryService.GetCategoryBySlug(ctx, req.Slug)
	if err != nil {
		s.logger.Error().Err(err).Str("slug", req.Slug).Msg("failed to get category by slug")
		if contains(err.Error(), "not found") || contains(err.Error(), "no rows") {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get category: %v", err))
	}

	// Convert to proto
	pbCategory := DomainToCategoryServiceProtoCategory(category)

	s.logger.Debug().Str("slug", req.Slug).Msg("category retrieved")

	return &categoriespb.GetCategoryBySlugResponse{
		Category: pbCategory,
	}, nil
}

// CreateCategory creates a new category (Admin endpoint)
func (s *Server) CreateCategoryForCategorySvc(ctx context.Context, req *categoriespb.CreateCategoryRequest) (*categoriespb.CreateCategoryResponse, error) {
	s.logger.Debug().
		Str("name", req.Name).
		Str("slug", req.Slug).
		Msg("CreateCategory called")

	// Validate request
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "slug is required")
	}
	if len(req.Name) < 2 {
		return nil, status.Error(codes.InvalidArgument, "name must be at least 2 characters")
	}
	if len(req.Slug) < 2 {
		return nil, status.Error(codes.InvalidArgument, "slug must be at least 2 characters")
	}

	// Convert proto to domain
	category := ProtoToCategoryServiceCreateDomain(req)

	// Create category via service
	createdCategory, err := s.categoryService.CreateCategory(ctx, category)
	if err != nil {
		s.logger.Error().Err(err).Str("name", req.Name).Msg("failed to create category")

		// Map database errors
		errMsg := err.Error()
		if contains(errMsg, "duplicate") || contains(errMsg, "unique constraint") {
			if contains(errMsg, "slug") {
				return nil, status.Error(codes.AlreadyExists, "category with this slug already exists")
			}
			return nil, status.Error(codes.AlreadyExists, "category already exists")
		}
		if contains(errMsg, "foreign key") || contains(errMsg, "parent") {
			return nil, status.Error(codes.InvalidArgument, "parent category not found")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create category: %v", err))
	}

	// Convert to proto
	pbCategory := DomainToCategoryServiceProtoCategory(createdCategory)

	s.logger.Info().
		Str("category_id", createdCategory.ID).
		Str("name", createdCategory.Name).
		Msg("category created successfully")

	return &categoriespb.CreateCategoryResponse{
		Category: pbCategory,
	}, nil
}

// UpdateCategory updates an existing category (Admin endpoint)
func (s *Server) UpdateCategoryForCategorySvc(ctx context.Context, req *categoriespb.UpdateCategoryRequest) (*categoriespb.UpdateCategoryResponse, error) {
	s.logger.Debug().
		Str("category_id", req.Id).
		Msg("UpdateCategory called")

	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "category ID must not be empty")
	}

	// At least one field must be set for update
	if req.Name == nil && req.Slug == nil && req.ParentId == nil && req.Icon == nil &&
		req.Description == nil && req.CustomUiComponent == nil && req.SortOrder == nil &&
		req.IsActive == nil && req.SeoTitle == nil && req.SeoDescription == nil &&
		req.SeoKeywords == nil && req.TitleEn == nil && req.TitleRu == nil && req.TitleSr == nil {
		return nil, status.Error(codes.InvalidArgument, "at least one field must be provided for update")
	}

	// Validate individual fields if present
	if req.Name != nil && len(*req.Name) < 2 {
		return nil, status.Error(codes.InvalidArgument, "name must be at least 2 characters")
	}
	if req.Slug != nil && len(*req.Slug) < 2 {
		return nil, status.Error(codes.InvalidArgument, "slug must be at least 2 characters")
	}

	// Convert proto to domain
	category := ProtoToCategoryServiceUpdateDomain(req)

	// Update category via service
	updatedCategory, err := s.categoryService.UpdateCategory(ctx, category)
	if err != nil {
		s.logger.Error().Err(err).Str("category_id", req.Id).Msg("failed to update category")

		// Map database errors
		errMsg := err.Error()
		if contains(errMsg, "not found") || contains(errMsg, "no rows") {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		if contains(errMsg, "duplicate") || contains(errMsg, "unique constraint") {
			if contains(errMsg, "slug") {
				return nil, status.Error(codes.AlreadyExists, "category with this slug already exists")
			}
			return nil, status.Error(codes.AlreadyExists, "category already exists")
		}
		if contains(errMsg, "foreign key") || contains(errMsg, "parent") {
			return nil, status.Error(codes.InvalidArgument, "parent category not found")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update category: %v", err))
	}

	// Convert to proto
	pbCategory := DomainToCategoryServiceProtoCategory(updatedCategory)

	s.logger.Info().
		Str("category_id", updatedCategory.ID).
		Msg("category updated successfully")

	return &categoriespb.UpdateCategoryResponse{
		Category: pbCategory,
	}, nil
}

// DeleteCategory deletes a category (Admin endpoint)
func (s *Server) DeleteCategoryForCategorySvc(ctx context.Context, req *categoriespb.DeleteCategoryRequest) (*categoriespb.DeleteCategoryResponse, error) {
	s.logger.Debug().
		Str("category_id", req.Id).
		Msg("DeleteCategory called")

	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "category ID must not be empty")
	}

	// Delete category via service
	err := s.categoryService.DeleteCategory(ctx, req.Id)
	if err != nil {
		s.logger.Error().Err(err).Str("category_id", req.Id).Msg("failed to delete category")

		// Map errors
		errMsg := err.Error()
		if contains(errMsg, "not found") || contains(errMsg, "no rows") {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		if contains(errMsg, "foreign key") || contains(errMsg, "has subcategories") || contains(errMsg, "has listings") {
			return nil, status.Error(codes.FailedPrecondition, "cannot delete category with subcategories or listings")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete category: %v", err))
	}

	s.logger.Info().
		Str("category_id", req.Id).
		Msg("category deleted successfully")

	return &categoriespb.DeleteCategoryResponse{
		Success: true,
		Message: "category deleted successfully",
	}, nil
}

// Helper function to filter active categories in tree
// Note: CategoryTreeNode doesn't have IsActive field, so we skip filtering for now
// In production, you may need to fetch IsActive from category service
func filterActiveCategoriesInTree(node *domain.CategoryTreeNode) *domain.CategoryTreeNode {
	if node == nil {
		return nil
	}

	// CategoryTreeNode doesn't have IsActive field
	// We would need to enhance the domain model or fetch from service
	// For now, return all nodes

	// Filter children recursively
	if len(node.Children) > 0 {
		activeChildren := make([]domain.CategoryTreeNode, 0, len(node.Children))
		for i := range node.Children {
			if filteredChild := filterActiveCategoriesInTree(&node.Children[i]); filteredChild != nil {
				activeChildren = append(activeChildren, *filteredChild)
			}
		}
		node.Children = activeChildren
	}

	return node
}
