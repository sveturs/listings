package models

import (
	"time"
)

// Role represents a system role with specific permissions
type Role struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	DisplayName  string    `json:"display_name" db:"display_name"`
	Description  string    `json:"description" db:"description"`
	IsSystem     bool      `json:"is_system" db:"is_system"`
	IsAssignable bool      `json:"is_assignable" db:"is_assignable"`
	Priority     int       `json:"priority" db:"priority"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// Relations
	Permissions []Permission `json:"permissions,omitempty"`
}

// Permission represents a specific action that can be performed on a resource
type Permission struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UserRole represents the association between a user and a role
type UserRole struct {
	UserID     int       `json:"user_id" db:"user_id"`
	RoleID     int       `json:"role_id" db:"role_id"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
	AssignedBy *int      `json:"assigned_by" db:"assigned_by"`
}

// RolePermission represents the association between a role and a permission
type RolePermission struct {
	RoleID       int       `json:"role_id" db:"role_id"`
	PermissionID int       `json:"permission_id" db:"permission_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// RoleAuditLog represents an audit log entry for role changes
type RoleAuditLog struct {
	ID           int                    `json:"id" db:"id"`
	UserID       *int                   `json:"user_id" db:"user_id"`
	TargetUserID int                    `json:"target_user_id" db:"target_user_id"`
	Action       string                 `json:"action" db:"action"`
	OldRoleID    *int                   `json:"old_role_id" db:"old_role_id"`
	NewRoleID    *int                   `json:"new_role_id" db:"new_role_id"`
	Details      map[string]interface{} `json:"details" db:"details"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

// UserWithRole extends UserProfile with role information
type UserWithRole struct {
	UserProfile
	Role *Role `json:"role,omitempty"`
}

// RoleFilter represents filter criteria for roles
type RoleFilter struct {
	Name         string `json:"name,omitempty"`
	IsSystem     *bool  `json:"is_system,omitempty"`
	IsAssignable *bool  `json:"is_assignable,omitempty"`
}

// PermissionFilter represents filter criteria for permissions
type PermissionFilter struct {
	Resource string `json:"resource,omitempty"`
	Action   string `json:"action,omitempty"`
}

// AssignRoleRequest represents a request to assign a role to a user
type AssignRoleRequest struct {
	UserID int `json:"user_id" validate:"required"`
	RoleID int `json:"role_id" validate:"required"`
}

// UpdateUserRoleRequest represents a request to update a user's role
type UpdateUserRoleRequest struct {
	RoleID int `json:"role_id" validate:"required"`
}

// CheckPermissionRequest represents a request to check if a user has a permission
type CheckPermissionRequest struct {
	UserID     int    `json:"user_id" validate:"required"`
	Permission string `json:"permission" validate:"required"`
}

// RoleResponse represents a role in API responses
type RoleResponse struct {
	ID           int                  `json:"id"`
	Name         string               `json:"name"`
	DisplayName  string               `json:"display_name"`
	Description  string               `json:"description"`
	IsSystem     bool                 `json:"is_system"`
	IsAssignable bool                 `json:"is_assignable"`
	Priority     int                  `json:"priority"`
	Permissions  []PermissionResponse `json:"permissions,omitempty"`
	UserCount    int                  `json:"user_count,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

// PermissionResponse represents a permission in API responses
type PermissionResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// Common role constants
const (
	RoleSuperAdmin         = "super_admin"
	RoleAdmin              = "admin"
	RoleModerator          = "moderator"
	RoleContentModerator   = "content_moderator"
	RoleReviewModerator    = "review_moderator"
	RoleChatModerator      = "chat_moderator"
	RoleDisputeManager     = "dispute_manager"
	RoleVendorManager      = "vendor_manager"
	RoleCategoryManager    = "category_manager"
	RoleMarketingManager   = "marketing_manager"
	RoleFinancialManager   = "financial_manager"
	RoleWarehouseManager   = "warehouse_manager"
	RoleWarehouseWorker    = "warehouse_worker"
	RolePickupManager      = "pickup_manager"
	RolePickupWorker       = "pickup_worker"
	RoleCourier            = "courier"
	RoleSupportL1          = "support_l1"
	RoleSupportL2          = "support_l2"
	RoleSupportL3          = "support_l3"
	RoleLegalAdvisor       = "legal_advisor"
	RoleComplianceOfficer  = "compliance_officer"
	RoleProfessionalVendor = "professional_vendor"
	RoleVendor             = "vendor"
	RoleIndividualSeller   = "individual_seller"
	RoleStorefrontStaff    = "storefront_staff"
	RoleVerifiedBuyer      = "verified_buyer"
	RoleVIPCustomer        = "vip_customer"
	RoleUser               = "user"
	RoleDataAnalyst        = "data_analyst"
	RoleBusinessAnalyst    = "business_analyst"
)

// Common permission constants
const (
	// User permissions
	PermUsersView       = "users.view"
	PermUsersList       = "users.list"
	PermUsersEdit       = "users.edit"
	PermUsersDelete     = "users.delete"
	PermUsersBlock      = "users.block"
	PermUsersVerify     = "users.verify"
	PermUsersAssignRole = "users.assign_role"
	PermUsersExport     = "users.export"

	// Admin permissions
	PermAdminAccess       = "admin.access"
	PermAdminUsers        = "admin.users"
	PermAdminCategories   = "admin.categories"
	PermAdminAttributes   = "admin.attributes"
	PermAdminTranslations = "admin.translations"

	// Listing permissions
	PermListingsCreate    = "listings.create"
	PermListingsEditOwn   = "listings.edit_own"
	PermListingsEditAny   = "listings.edit_any"
	PermListingsDeleteOwn = "listings.delete_own"
	PermListingsDeleteAny = "listings.delete_any"
	PermListingsModerate  = "listings.moderate"
	PermListingsViewAll   = "listings.view_all"
	PermListingsApprove   = "listings.approve"
	PermListingsReject    = "listings.reject"

	// Order permissions
	PermOrdersViewAll = "orders.view_all"
	PermOrdersViewOwn = "orders.view_own"
	PermOrdersProcess = "orders.process"
	PermOrdersCancel  = "orders.cancel"
	PermOrdersRefund  = "orders.refund"
	PermOrdersExport  = "orders.export"

	// System permissions
	PermSystemViewLogs       = "system.view_logs"
	PermSystemManageSettings = "system.manage_settings"
	PermSystemManageRoles    = "system.manage_roles"
	PermSystemViewAudit      = "system.view_audit"
)

// HasPermission checks if a role has a specific permission
func (r *Role) HasPermission(permissionName string) bool {
	for _, p := range r.Permissions {
		if p.Name == permissionName {
			return true
		}
	}
	return false
}

// IsHigherPriority checks if this role has higher priority than another
func (r *Role) IsHigherPriority(other *Role) bool {
	return r.Priority < other.Priority
}
