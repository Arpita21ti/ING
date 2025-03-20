package role

import (
	"context"
)

// Repository defines the interface for role data access
type Repository interface {
	// Role operations
	CreateRole(ctx context.Context, role *Role) error
	GetRoleByID(ctx context.Context, id uint) (*Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	ListRoles(ctx context.Context) ([]Role, error)
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id uint) error

	// Permission operations
	CreatePermission(ctx context.Context, permission *Permission) error
	GetPermissionByID(ctx context.Context, id uint) (*Permission, error)
	ListPermissions(ctx context.Context) ([]Permission, error)
	UpdatePermission(ctx context.Context, permission *Permission) error
	DeletePermission(ctx context.Context, id uint) error

	// Role-Permission operations
	AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]Permission, error)

	// User-Role operations
	AssignRoleToUser(ctx context.Context, userID, roleID uint) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]Role, error)
	GetUsersWithRole(ctx context.Context, roleID uint) ([]uint, error) // Returns user IDs

	// Special queries
	HasPermission(ctx context.Context, userID uint, resource, action string) (bool, error)
}
