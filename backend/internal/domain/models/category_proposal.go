package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// CategoryProposalStatus represents the status of a category proposal
type CategoryProposalStatus string

const (
	CategoryProposalStatusPending  CategoryProposalStatus = "pending"
	CategoryProposalStatusApproved CategoryProposalStatus = "approved"
	CategoryProposalStatusRejected CategoryProposalStatus = "rejected"
)

// NameTranslations represents translations for category name
type NameTranslations map[string]string

// Value implements driver.Valuer for NameTranslations
func (nt NameTranslations) Value() (driver.Value, error) {
	return json.Marshal(nt)
}

// Scan implements sql.Scanner for NameTranslations
func (nt *NameTranslations) Scan(value interface{}) error {
	if value == nil {
		*nt = make(NameTranslations)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to scan NameTranslations: unsupported type")
	}

	return json.Unmarshal(bytes, nt)
}

// CategoryProposal represents a proposal for a new category
// AI suggests new categories based on import analysis
type CategoryProposal struct {
	ID                     int                    `json:"id" db:"id"`
	ProposedByUserID       int                    `json:"proposed_by_user_id" db:"proposed_by_user_id"`
	StorefrontID           *int                   `json:"storefront_id,omitempty" db:"storefront_id"`
	Name                   string                 `json:"name" db:"name"`
	NameTranslations       NameTranslations       `json:"name_translations" db:"name_translations"`
	ParentCategoryID       *int                   `json:"parent_category_id,omitempty" db:"parent_category_id"`
	Description            *string                `json:"description,omitempty" db:"description"`
	Reasoning              *string                `json:"reasoning,omitempty" db:"reasoning"`                               // AI reasoning
	ExpectedProducts       int                    `json:"expected_products" db:"expected_products"`                         // Number of products that would use this
	ExternalCategorySource *string                `json:"external_category_source,omitempty" db:"external_category_source"` // Original external category
	SimilarCategories      []int                  `json:"similar_categories,omitempty" db:"similar_categories"`             // Related category IDs
	Tags                   []string               `json:"tags,omitempty" db:"tags"`
	Status                 CategoryProposalStatus `json:"status" db:"status"`
	ReviewedByUserID       *int                   `json:"reviewed_by_user_id,omitempty" db:"reviewed_by_user_id"`
	ReviewedAt             *time.Time             `json:"reviewed_at,omitempty" db:"reviewed_at"`
	CreatedAt              time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at" db:"updated_at"`
}

// CreateCategoryProposalRequest represents request to create a category proposal
type CreateCategoryProposalRequest struct {
	StorefrontID           *int             `json:"storefront_id,omitempty"`
	Name                   string           `json:"name" validate:"required,min=2,max=255"`
	NameTranslations       NameTranslations `json:"name_translations,omitempty"`
	ParentCategoryID       *int             `json:"parent_category_id,omitempty"`
	Description            *string          `json:"description,omitempty"`
	Reasoning              *string          `json:"reasoning,omitempty"`
	ExpectedProducts       int              `json:"expected_products,omitempty"`
	ExternalCategorySource *string          `json:"external_category_source,omitempty"`
	SimilarCategories      []int            `json:"similar_categories,omitempty"`
	Tags                   []string         `json:"tags,omitempty"`
}

// UpdateCategoryProposalRequest represents request to update a category proposal
type UpdateCategoryProposalRequest struct {
	Name             *string          `json:"name,omitempty"`
	NameTranslations NameTranslations `json:"name_translations,omitempty"`
	ParentCategoryID *int             `json:"parent_category_id,omitempty"`
	Description      *string          `json:"description,omitempty"`
	Tags             []string         `json:"tags,omitempty"`
}

// CategoryProposalListResponse represents paginated list of proposals
type CategoryProposalListResponse struct {
	Proposals  []CategoryProposal `json:"proposals"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// CategoryProposalApproveRequest represents request to approve a proposal
type CategoryProposalApproveRequest struct {
	CreateCategory bool `json:"create_category"` // If true, creates the category in c2c_categories
}

// CategoryProposalRejectRequest represents request to reject a proposal
type CategoryProposalRejectRequest struct {
	Reason *string `json:"reason,omitempty"` // Optional rejection reason
}
