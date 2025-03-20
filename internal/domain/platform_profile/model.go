// 1. TODO: Include gradually - line: 42
// 2. TODO: Include OAuth
// 3. TODO: Add models for soft delete recovery.

// Rest models are final.
// No changes needed.

package platform_profile

import (
	"time"

	"github.com/google/uuid"
)

// Status represents the current state of a platform profile
type Status string

const (
	StatusActivated   Status = "activated"
	StatusDeactivated Status = "deactivated"
	StatusSuspended   Status = "suspended"
	StatusPending     Status = "pending"
	StatusLocked      Status = "locked"
)

// PlatformProfile represents a user account on the platform
// This is separate from student-specific information
type PlatformProfile struct {
	ID                  uuid.UUID  `json:"id"`
	Username            string     `json:"username"`
	Email               string     `json:"email"`
	PasswordHash        string     `json:"-"`
	Status              Status     `json:"status"`
	VerifiedAt          *time.Time `json:"verified_at,omitempty"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	FailedLoginAttempts int        `json:"-"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedByUserAt     time.Time  `json:"updated_by_user_at"`
	UpdatedBySystemAt   time.Time  `json:"updated_by_system_at"`
}

// CreateProfileRequest represents data needed to create a new profile
type CreateProfileRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// GetProfileRequest represents data needed to get an existing profile
type GetProfileRequest struct {
	ID       *uuid.UUID `json:"id,omitempty" validate:"omitempty,uuid4"`
	Username *string    `json:"username,omitempty" validate:"omitempty,min=3,max=30"`
	Email    *string    `json:"email,omitempty" validate:"omitempty,email"`
}

// UpdateProfileRequest represents data that can be updated for a profile
type UpdateProfileRequest struct {
	// Used pointers to string (*string) to enable partial updates and to check incoming data explicitely
	// it allows for the distinction between "field not provided" and "field provided but empty"
	// needed because there's an important distinction between
	// "don't change this field" (nil) and "change this field to empty" (pointer to empty string).
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=30"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
}

// LoginRequest represents the credentials needed for login
type LoginRequest struct {
	// omitempty is used to give users the liberty to login using either Username or Email.
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password" validate:"required"`
}

// GetProfileRequest represents data needed to get an existing profile
type GetSoftDeletedProfileRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=30"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
}

// PasswordResetRequest represents data needed to request password reset
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordChangeRequest represents data needed for password change
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,nefield=CurrentPassword"`
}

// PasswordResetToken stores information for password reset functionality
type PasswordResetToken struct {
	ProfileID uuid.UUID `json:"-"`
	Token     string    `json:"-"`
	ExpiresAt time.Time `json:"-"` // Kept redundantly for admin tasks like auditing, lifetime extension, etc. and adding security.
	IsUsed    bool      `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// PasswordResetConfirmation represents data needed to confirm password reset
type PasswordResetConfirmation struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ProfilePreference stores user preferences for the platform
type ProfilePreference struct {
	ID                 uuid.UUID `json:"id"`
	ProfileID          uuid.UUID `json:"profile_id"`
	NotificationsEmail bool      `json:"notifications_email"`

	// TODO: 1.
	// NotificationsPush  bool      `json:"notifications_push"`
	// NotificationsSMS   bool      `json:"notifications_sms"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Store in frontend cache and preferences.
	// Language  string    `json:"language"`
	// DarkModeOn      bool      `json:"dark_mode_on"`
}

// ProfileRole associates a profile with one or more roles
type ProfileRole struct {
	ProfileID  uuid.UUID `json:"profile_id"`
	RoleID     uuid.UUID `json:"role_id"`
	AssignedAt time.Time `json:"assigned_at"`
}
