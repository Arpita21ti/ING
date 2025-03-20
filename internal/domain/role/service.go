package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrRoleNotFound        = errors.New("role not found")
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrDuplicateRole       = errors.New("role already exists")
	ErrDuplicatePermission = errors.New("permission already exists")
	ErrInvalidRole         = errors.New("invalid role data")
	ErrInvalidPermission   = errors.New("invalid permission data")
	ErrUnauthorized        = errors.New("user does not have required permission")
)

// Service provides role management operations
type Service interface {
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
	GetDefaultRoleID(ctx context.Context) (uuid.UUID, error)

	// User-Role operations
	AssignRoleToUser(ctx context.Context, userID, roleID uint) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]Role, error)

	// Permission checking
	HasPermission(ctx context.Context, userID uint, resource, action string) (bool, error)
	Authorize(ctx context.Context, userID uint, resource, action string) error

	// Bulk operations
	SyncUserRoles(ctx context.Context, userID uint, roleIDs []uint) error
	SyncRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new role service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetDefaultRoleID returns the ID of the default role for new profiles
func (s *service) GetDefaultRoleID(ctx context.Context) (uint, error) {
	// Get the role with the default role name (e.g., "guest")
	defaultRole, err := s.repo.GetRoleByName(ctx, "guest") // Or whatever your default role is called
	if err != nil {
		return 0, fmt.Errorf("failed to fetch default role: %w", err)
	}

	if defaultRole == nil {
		return 0, errors.New("default role does not exist")
	}

	return defaultRole.ID, nil
}

// Implementation of Service interface methods

func (s *service) CreateRole(ctx context.Context, role *Role) error {
	// Validate role
	if role.Name == "" {
		return ErrInvalidRole
	}

	// Check for duplicate
	existing, err := s.repo.GetRoleByName(ctx, role.Name)
	if err == nil && existing != nil {
		return ErrDuplicateRole
	}

	return s.repo.CreateRole(ctx, role)
}

func (s *service) GetRoleByID(ctx context.Context, id uint) (*Role, error) {
	return s.repo.GetRoleByID(ctx, id)
}

func (s *service) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	return s.repo.GetRoleByName(ctx, name)
}

func (s *service) ListRoles(ctx context.Context) ([]Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *service) UpdateRole(ctx context.Context, role *Role) error {
	// Validate role
	if role.Name == "" {
		return ErrInvalidRole
	}

	// Ensure role exists
	existing, err := s.repo.GetRoleByID(ctx, role.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRoleNotFound
	}

	return s.repo.UpdateRole(ctx, role)
}

func (s *service) DeleteRole(ctx context.Context, id uint) error {
	// Ensure role exists
	existing, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRoleNotFound
	}

	return s.repo.DeleteRole(ctx, id)
}

// Authorization logic
func (s *service) HasPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	return s.repo.HasPermission(ctx, userID, resource, action)
}

func (s *service) Authorize(ctx context.Context, userID uint, resource, action string) error {
	hasPermission, err := s.HasPermission(ctx, userID, resource, action)
	if err != nil {
		return err
	}
	if !hasPermission {
		return ErrUnauthorized
	}
	return nil
}

// Implement other service methods...
