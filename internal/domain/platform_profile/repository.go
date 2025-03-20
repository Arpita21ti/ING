package platform_profile

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the data access interface for platform profiles
type Repository interface {

	// Profile operations (Verification)
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)

	// Profile operations (CRUD)
	CreateProfile(ctx context.Context, profile *PlatformProfile) error
	GetProfileByID(ctx context.Context, id uuid.UUID) (*PlatformProfile, error)
	GetProfileByEmail(ctx context.Context, email string) (*PlatformProfile, error)
	GetProfileByUsername(ctx context.Context, username string) (*PlatformProfile, error)

	UpdateProfile(ctx context.Context, profile *PlatformProfile) error

	SoftDeleteProfile(ctx context.Context, id uuid.UUID) error

	// Status management
	VerifyProfile(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error

	// Password / Authentication  operations
	CreatePasswordResetToken(ctx context.Context, resetToken *PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, token string) error
	DeleteOtherPasswordResetTokens(ctx context.Context, profileID uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error

	// Recovery operations
	GetSoftDeletedProfileByUsername(ctx context.Context, username string) (*PlatformProfile, error)
	GetSoftDeletedProfileByEmail(ctx context.Context, email string) (*PlatformProfile, error)
	RestoreSoftDeletedProfile(ctx context.Context, id uuid.UUID) error

	// Hard / Permanaent delete operations
	HardDeleteProfile(ctx context.Context, id uuid.UUID) error

	// Role associations
	AssignRoleToProfile(ctx context.Context, profileID, roleID uuid.UUID) error
	GetProfileRoles(ctx context.Context, profileID uuid.UUID) ([]uuid.UUID, error)
	RemoveRoleFromProfile(ctx context.Context, profileID, roleID uuid.UUID) error

	// Preferences
	SavePreferences(ctx context.Context, prefs *ProfilePreference) error
	GetPreferences(ctx context.Context, profileID uuid.UUID) (*ProfilePreference, error)

	// Bulk/query operations (only for admin)
	GetProfiles(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*PlatformProfile, int, error)
	GetProfilesByRole(ctx context.Context, roleID uuid.UUID, page, pageSize int) ([]*PlatformProfile, int, error)
	GetProfilesByPreferences(ctx context.Context, preferences map[string]interface{}, page, pageSize int) ([]*PlatformProfile, int, error)
	DeleteProfiles(ctx context.Context, profileIDs []uuid.UUID, hardDelete bool) error

	// Login management
	RecordLogin(ctx context.Context, id uuid.UUID) error

	// Failed Login Attempt management
	IncrementFailedLoginAttempts(ctx context.Context, id uuid.UUID) error
	ResetFailedLoginAttempts(ctx context.Context, id uuid.UUID) error
}
