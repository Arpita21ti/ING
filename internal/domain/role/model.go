package role

import (
	"time"
)

// Role represents a role in the system with associated permissions
type Role struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Permission represents an action that can be performed in the system
type Permission struct {
	ID          uint      `json:"id"`
	Resource    string    `json:"resource"` // e.g., "student", "quiz", "event"
	Action      string    `json:"action"`   // e.g., "create", "read", "update", "delete"
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRole represents the relationship between users and roles
type UserRole struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	RoleID    uint      `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RolePermission represents the relationship between roles and permissions
type RolePermission struct {
	ID           uint      `json:"id"`
	RoleID       uint      `json:"role_id"`
	PermissionID uint      `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Common predefined roles
const (
	RoleAdmin       = "admin"
	RoleCoordinator = "coordinator"
	RoleStudent     = "student"
	RoleFaculty     = "faculty"
	RoleGuest       = "guest"
)

// Permission actions
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionManage = "manage" // Special permission for all actions
)
