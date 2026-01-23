package grpc

import (
	"context"

	categoriespb "github.com/vondi-global/listings/api/proto/categories/v1"
)

// CategoryServiceServer is a separate handler for CategoryService gRPC interface
// This is needed because CategoryService and ListingsService have conflicting method names
type CategoryServiceServer struct {
	categoriespb.UnimplementedCategoryServiceServer
	server *Server // Delegate to main server
}

// NewCategoryServiceServer creates a new CategoryService handler
func NewCategoryServiceServer(mainServer *Server) *CategoryServiceServer {
	return &CategoryServiceServer{
		server: mainServer,
	}
}

// GetCategories delegates to handler
func (css *CategoryServiceServer) GetCategories(ctx context.Context, req *categoriespb.GetCategoriesRequest) (*categoriespb.GetCategoriesResponse, error) {
	return css.server.GetCategoriesForCategorySvc(ctx, req)
}

// GetCategoryTree delegates to handler
func (css *CategoryServiceServer) GetCategoryTree(ctx context.Context, req *categoriespb.GetCategoryTreeRequest) (*categoriespb.GetCategoryTreeResponse, error) {
	return css.server.GetCategoryTreeForCategorySvc(ctx, req)
}

// GetCategory delegates to handler
func (css *CategoryServiceServer) GetCategory(ctx context.Context, req *categoriespb.GetCategoryRequest) (*categoriespb.GetCategoryResponse, error) {
	return css.server.GetCategoryForCategorySvc(ctx, req)
}

// GetCategoryBySlug delegates to handler
func (css *CategoryServiceServer) GetCategoryBySlug(ctx context.Context, req *categoriespb.GetCategoryBySlugRequest) (*categoriespb.GetCategoryBySlugResponse, error) {
	return css.server.GetCategoryBySlugForCategorySvc(ctx, req)
}

// CreateCategory delegates to handler
func (css *CategoryServiceServer) CreateCategory(ctx context.Context, req *categoriespb.CreateCategoryRequest) (*categoriespb.CreateCategoryResponse, error) {
	return css.server.CreateCategoryForCategorySvc(ctx, req)
}

// UpdateCategory delegates to handler
func (css *CategoryServiceServer) UpdateCategory(ctx context.Context, req *categoriespb.UpdateCategoryRequest) (*categoriespb.UpdateCategoryResponse, error) {
	return css.server.UpdateCategoryForCategorySvc(ctx, req)
}

// DeleteCategory delegates to handler
func (css *CategoryServiceServer) DeleteCategory(ctx context.Context, req *categoriespb.DeleteCategoryRequest) (*categoriespb.DeleteCategoryResponse, error) {
	return css.server.DeleteCategoryForCategorySvc(ctx, req)
}
