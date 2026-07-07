package rbac

import "time"

// Organization represents a multi-tenant isolation boundary.
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Role defines a named set of permissions.
type Role struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"orgId"` // Empty for system-wide roles
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Membership binds a User to an Organization with a specific Role.
type Membership struct {
	ID          string     `json:"id"`
	OrgID       string     `json:"orgId"`
	UserID      string     `json:"userId"`
	RoleID      string     `json:"roleId"`
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	JoinedAt    *time.Time `json:"joinedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type MembershipDetails struct {
	ID          string     `json:"id"`
	OrgID       string     `json:"orgId"`
	UserID      string     `json:"userId"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	RoleID      string     `json:"roleId"`
	RoleName    string     `json:"roleName"`
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	JoinedAt    *time.Time `json:"joinedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// ResourcePolicy maps a specific resource to its owning Organization.
type ResourcePolicy struct {
	ResourceID   string `json:"resourceId"`
	ResourceType string `json:"resourceType"` // e.g., "workspace", "project"
	OrgID        string `json:"orgId"`
}
