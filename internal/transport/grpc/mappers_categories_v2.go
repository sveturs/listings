package grpc

import (
	categoriesv2 "github.com/vondi-global/listings/api/proto/categories/v2"
	"github.com/vondi-global/listings/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DomainToProtoCategoryV2 converts domain.LocalizedCategory to proto CategoryV2
func DomainToProtoCategoryV2(cat *domain.LocalizedCategory) *categoriesv2.CategoryV2 {
	if cat == nil {
		return nil
	}

	pbCat := &categoriesv2.CategoryV2{
		Id:              cat.ID.String(),
		Slug:            cat.Slug,
		Level:           cat.Level,
		Path:            cat.Path,
		SortOrder:       cat.SortOrder,
		Name:            cat.Name,
		Description:     cat.Description,
		MetaTitle:       cat.MetaTitle,
		MetaDescription: cat.MetaDescription,
		MetaKeywords:    cat.MetaKeywords,
		IsActive:        cat.IsActive,
		CreatedAt:       timestamppb.New(cat.CreatedAt),
		UpdatedAt:       timestamppb.New(cat.UpdatedAt),
	}

	// Optional fields
	if cat.ParentID != nil {
		parentIDStr := cat.ParentID.String()
		pbCat.ParentId = &parentIDStr
	}

	if cat.Icon != nil {
		pbCat.Icon = cat.Icon
	}

	if cat.ImageURL != nil {
		pbCat.ImageUrl = cat.ImageURL
	}

	return pbCat
}

// DomainToProtoCategoryTreeV2 converts domain.CategoryTreeV2 to proto CategoryTreeV2
func DomainToProtoCategoryTreeV2(tree []*domain.CategoryTreeV2) []*categoriesv2.CategoryTreeV2 {
	if tree == nil {
		return nil
	}

	pbTrees := make([]*categoriesv2.CategoryTreeV2, 0, len(tree))
	for _, node := range tree {
		if node == nil {
			continue
		}

		pbNode := &categoriesv2.CategoryTreeV2{
			Category:      DomainToProtoCategoryV2(node.Category),
			Subcategories: DomainToProtoCategoryTreeV2(node.Subcategories),
		}
		pbTrees = append(pbTrees, pbNode)
	}

	return pbTrees
}

// DomainToProtoBreadcrumb converts domain.CategoryBreadcrumb to proto CategoryBreadcrumb
func DomainToProtoBreadcrumb(breadcrumbs []*domain.CategoryBreadcrumb) []*categoriesv2.CategoryBreadcrumb {
	if breadcrumbs == nil {
		return nil
	}

	pbBreadcrumbs := make([]*categoriesv2.CategoryBreadcrumb, 0, len(breadcrumbs))
	for _, bc := range breadcrumbs {
		if bc == nil {
			continue
		}

		pbBc := &categoriesv2.CategoryBreadcrumb{
			Id:    bc.ID.String(),
			Slug:  bc.Slug,
			Name:  bc.Name,
			Level: bc.Level,
		}
		pbBreadcrumbs = append(pbBreadcrumbs, pbBc)
	}

	return pbBreadcrumbs
}
