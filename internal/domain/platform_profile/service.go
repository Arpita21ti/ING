// 1. TODO: Implement VerifyProfile
// 2. TODO: Implement ActivateProfile
// 3. TODO: Reset Preferences Function.
// 4. TODO: Save Preferences Function.
// 5. TODO: List All Profiles Function.
// 6. TODO: Add Lock and Unlock mechanisms for multiple failed logins.
// 7. TODO: Implement ValidatePasswordResetToken
// 8. TODO: Implement SendPasswordResetEmail
// 9. TODO: Implement ExpireOldResetTokens
// 10. TODO: Implement ForcePasswordChange
// 11. TODO: Implement DeleteExpiredResetTokens

package platform_profile

import (
	"context"
	"server/pkg/logger"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"server/internal/common/errors"
	"server/internal/domain/role"
)

// Service defines the business logic for platform profiles
type Service interface {
	// Account Management (CRUD)
	RegisterProfile(ctx context.Context, req CreateProfileRequest) (*PlatformProfile, error)
	GetProfile(ctx context.Context, req GetProfileRequest) (*PlatformProfile, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, req UpdateProfileRequest) (*PlatformProfile, error)
	SoftDeleteProfile(ctx context.Context, id uuid.UUID) error

	// Account Verification and Status Management
	VerifyProfile(ctx context.Context, id uuid.UUID) error
	ActivateProfile(ctx context.Context, id uuid.UUID) error
	SuspendProfile(ctx context.Context, id uuid.UUID) error
	DeactivateProfile(ctx context.Context, id uuid.UUID) error

	// Account Restoration
	GetSoftDeletedProfile(ctx context.Context, req GetSoftDeletedProfileRequest) (*PlatformProfile, error)

	// Authentication
	AuthenticateProfile(ctx context.Context, req LoginRequest) (*PlatformProfile, error)
	ChangePassword(ctx context.Context, id uuid.UUID, req PasswordChangeRequest) error
	RequestPasswordReset(ctx context.Context, req PasswordResetRequest) error
	ConfirmPasswordReset(ctx context.Context, req PasswordResetConfirmation) error
	LockAccountAfterFailedAttempts(ctx context.Context, id uuid.UUID, maxAttempts int) error
	UnlockAccount(ctx context.Context, id uuid.UUID) error
	ValidatePasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error)
	SendPasswordResetEmail(ctx context.Context, email string, resetLink string) error
	ExpireOldResetTokens(ctx context.Context, profileID uuid.UUID) error
	ForcePasswordChange(ctx context.Context, id uuid.UUID) error
	DeleteExpiredResetTokens(ctx context.Context) error

	// Role Management
	AssignRole(ctx context.Context, profileID uuid.UUID, roleID uuid.UUID) error
	RemoveRole(ctx context.Context, profileID uuid.UUID, roleID uuid.UUID) error
	HasRole(ctx context.Context, profileID uuid.UUID, roleName string) (bool, error)

	// Preferences
	UpdatePreferences(ctx context.Context, profileID uuid.UUID, prefs ProfilePreference) error
	GetPreferences(ctx context.Context, profileID uuid.UUID) (*ProfilePreference, error)

	// Administrative functions
	SearchProfiles(ctx context.Context, query string, page, pageSize int) ([]*PlatformProfile, int, error)
	ListProfilesByRole(ctx context.Context, roleID uuid.UUID, page, pageSize int) ([]*PlatformProfile, int, error)
	ResetPreferencesToDefault(ctx context.Context, profileID uuid.UUID) (*ProfilePreference, error)
}

// The "service" struct is the concrete implementation of the "Service" interface.
type service struct {
	repo        Repository
	roleService role.Service

	// tokenGenerator utils.TokenGenerator
	logger logger.Logger
}

// NewService creates a new platform profile service
func NewService(
	repo Repository,
	roleService role.Service,
	// tokenGen utils.TokenGenerator,
	logger logger.Logger,
) Service {
	return &service{
		repo:        repo,
		roleService: roleService,
		// tokenGenerator: tokenGen,
		logger: logger,
	}
}

// RegisterProfile creates a new platform profile
func (s *service) RegisterProfile(ctx context.Context, req CreateProfileRequest) (*PlatformProfile, error) {
	s.logger.Debug("Starting profile registration", "username", req.Username, "email", req.Email)

	// Check if username already exists
	existingUser, err := s.repo.UsernameExists(ctx, req.Username)
	if err != nil && !errors.IsNotFoundErrorDomain(err) {
		s.logger.Error("Failed to check username existence", "username", req.Username, "error", err)
		return nil, errors.NewDatabaseError("fetching username", err)
	}
	if existingUser {
		s.logger.Warn("Username already exists", "username", req.Username)
		return nil, errors.NewConflictError("username", map[string]interface{}{"username": req.Username})
	}

	// Check if email already exists
	existingEmail, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil && !errors.IsNotFoundErrorDomain(err) {
		s.logger.Error("Failed to check email existence", "email", req.Email, "error", err)
		return nil, errors.NewDatabaseError("fetching email", err)
	}
	if existingEmail {
		s.logger.Warn("Email already exists", "email", req.Email)
		return nil, errors.NewConflictError("email", map[string]interface{}{"email": req.Email})
	}

	// Hash password securely
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "username", req.Username, "error", err)
		return nil, errors.NewBusinessError("PASSWORD_HASHING_FAILED", "password hashing failed", nil)
	}

	now := time.Now()
	profile := &PlatformProfile{
		ID:                  uuid.New(),
		Username:            req.Username,
		Email:               req.Email,
		PasswordHash:        string(passwordHash),
		Status:              StatusPending, // Default to 'pending' for new profiles
		FailedLoginAttempts: 0,
		CreatedAt:           now,
		UpdatedByUserAt:     now,
		UpdatedBySystemAt:   now,
	}

	// Save profile to database
	if err := s.repo.CreateProfile(ctx, profile); err != nil {
		s.logger.Error("Failed to create profile", "username", req.Username, "error", err)
		return nil, errors.NewBusinessError("PROFILE_CREATION_FAILED", "failed to create profile", nil)
	}

	// Create default preferences (not a blocker if it fails)
	prefs := &ProfilePreference{
		ID:                 uuid.New(),
		ProfileID:          profile.ID,
		NotificationsEmail: true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := s.repo.SavePreferences(ctx, prefs); err != nil {
		s.logger.Warn("Failed to create default preferences", "profile_id", profile.ID, "error", err)
	}

	// Assign default role if needed
	defaultRoleID, err := s.roleService.GetDefaultRoleID(ctx)
	if err != nil {
		s.logger.Warn("Failed to get default role for profile", "profile_id", profile.ID, "error", err)
	} else {
		err = s.repo.AssignRoleToProfile(ctx, profile.ID, defaultRoleID)
		if err != nil {
			s.logger.Warn("Failed to assign default role", "profile_id", profile.ID, "role_id", defaultRoleID, "error", err)
		}
	}

	if err != nil {
		s.logger.Error("Failed to commit transaction for profile registration", "error", err)
		return nil, errors.NewDatabaseError("transaction commit", err)
	}

	s.logger.Info("Profile registered successfully", "username", req.Username, "profile_id", profile.ID)

	// Remove password hash before returning response
	profile.PasswordHash = ""
	return profile, nil
}

// GetProfile retrieves a profile based on provided identifier (ID, username, or email)
func (s *service) GetProfile(ctx context.Context, req GetProfileRequest) (*PlatformProfile, error) {

	var profile *PlatformProfile
	var err error

	// Log which identifier is being used
	if req.ID != nil {
		s.logger.Debug(
			"Retrieving profile by ID",
			"profileID",
			*req.ID,
		)

		profile, err = s.repo.GetProfileByID(ctx, *req.ID)

		if err != nil {
			s.logger.Error(
				"Failed to retrieve profile by ID",
				"profileID", *req.ID,
				"error", err,
			)
			return nil, err
		}

		s.logger.Debug(
			"Profile retrieved successfully by ID",
			"profileID", *req.ID,
		)

	} else if req.Username != nil {
		s.logger.Debug(
			"Retrieving profile by username",
			"username",
			*req.Username,
		)

		profile, err = s.repo.GetProfileByUsername(ctx, *req.Username)

		if err != nil {
			s.logger.Error(
				"Failed to retrieve profile by username",
				"username", *req.Username,
				"error", err,
			)
			return nil, err
		}

		s.logger.Debug(
			"Profile retrieved successfully by username",
			"username", *req.Username,
			"profileID", profile.ID,
		)

	} else if req.Email != nil {
		s.logger.Debug(
			"Retrieving profile by email",
			"email",
			*req.Email,
		)

		profile, err = s.repo.GetProfileByEmail(ctx, *req.Email)

		if err != nil {
			s.logger.Error(
				"Failed to retrieve profile by email",
				"email", *req.Email,
				"error", err,
			)
			return nil, err
		}

		s.logger.Debug(
			"Profile retrieved successfully by email",
			"email", *req.Email,
			"profileID", profile.ID,
		)

	} else {
		errMsg := "no valid identifier provided for profile retrieval"
		s.logger.Error(errMsg)
		return nil, errors.NewValidationError(errMsg, nil)
	}

	// Don't expose password hash
	profile.PasswordHash = ""
	return profile, nil
}

// UpdateProfile updates a profile's information
func (s *service) UpdateProfile(ctx context.Context, id uuid.UUID, req UpdateProfileRequest) (*PlatformProfile, error) {

	s.logger.Info(
		"Updating profile",
		"profileID", id,
		"usernameChange", req.Username != nil,
		"emailChange", req.Email != nil,
	)

	profile, err := s.repo.GetProfileByID(ctx, id)

	if err != nil {
		s.logger.Error(
			"Failed to retrieve profile for update",
			"profileID", id,
			"error", err,
		)
		return nil, err
	}

	if req.Username != nil {
		// Check if new username is available
		if *req.Username != profile.Username {
			s.logger.Debug(
				"Checking username availability",
				"profileID", id,
				"newUsername", *req.Username,
			)

			existing, err := s.repo.GetProfileByUsername(ctx, *req.Username)

			if err == nil && existing != nil {
				s.logger.Warn(
					"Username already taken",
					"profileID", id,
					"requestedUsername", *req.Username,
				)

				return nil, errors.NewDomainError(
					"username already taken",
					errors.ValidationError,
					"USERNAME_TAKEN", // Add a unique error code
					nil,              // No additional details
					nil,              // No underlying error
				)
			}

			profile.Username = *req.Username

			s.logger.Debug(
				"Username updated",
				"profileID", id,
				"newUsername", *req.Username,
			)
		}
	}

	if req.Email != nil {
		// Check if new email is available
		if *req.Email != profile.Email {
			s.logger.Debug(
				"Checking email availability",
				"profileID", id,
				"newEmail", *req.Email,
			)

			existing, err := s.repo.GetProfileByEmail(ctx, *req.Email)

			if err == nil && existing != nil {
				s.logger.Warn(
					"Email already registered",
					"profileID", id,
					"requestedEmail", *req.Email,
				)

				return nil, errors.NewDomainError(
					"email already registered",
					errors.ValidationError,
					"EMAIL_REGISTERED", // Add a unique error code
					nil,                // No additional details
					nil,                // No underlying error
				)

			}
			profile.Email = *req.Email
			s.logger.Debug(
				"Email updated",
				"profileID", id,
				"newEmail", *req.Email,
			)
		}
	}

	// Set timestamps at the service level
	now := time.Now()

	profile.UpdatedByUserAt = now
	profile.UpdatedBySystemAt = now

	if err := s.repo.UpdateProfile(ctx, profile); err != nil {
		s.logger.Error(
			"Failed to update profile in database",
			"profileID", id,
			"error", err,
		)
		return nil, errors.NewDatabaseError("updating profile", err)
	}

	s.logger.Info("Profile updated successfully", "profileID", id)

	// Don't expose password hash
	profile.PasswordHash = ""
	return profile, nil
}

// DeleteProfile soft-deletes a profile
func (s *service) SoftDeleteProfile(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("Attempting to soft-delete profile", "id", id)

	err := s.repo.SoftDeleteProfile(ctx, id)
	if err != nil {
		s.logger.Error("Failed to soft-delete profile", "id", id, "error", err)
		return errors.NewDatabaseError("soft delete profile", err)
	}

	s.logger.Info("Successfully soft-deleted profile", "id", id)
	return nil
}

// VerifyProfile verifies the profile and updates the status to activated
func (s *service) VerifyProfile(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("Starting profile verification", "profile_id", id)

	// Mark the profile as verified
	err := s.repo.VerifyProfile(ctx, id)
	if err != nil {
		if errors.IsNotFoundErrorDomain(err) {
			s.logger.Warn("Profile not found for verification", "profile_id", id)
			return errors.NewNotFoundError("profile", map[string]interface{}{"id": id})
		}
		s.logger.Error("Failed to verify profile", "profile_id", id, "error", err)
		return errors.NewDatabaseError("profile verification", err)
	}

	// Update profile status to 'activated'
	err = s.repo.UpdateStatus(ctx, id, StatusActivated)
	if err != nil {
		s.logger.Error("Failed to update profile status after verification", "profile_id", id, "error", err)
		return errors.NewDatabaseError("status update", err)
	}

	s.logger.Info("Profile verified and activated successfully", "profile_id", id)
	return nil
}

// ActivateProfile sets a profile's status to activated
func (s *service) ActivateProfile(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("Activating profile", "profile_id", id)

	err := s.repo.UpdateStatus(ctx, id, StatusActivated)
	if err != nil {
		if errors.IsNotFoundErrorDomain(err) {
			s.logger.Warn("Profile not found for activation", "profile_id", id)
			return errors.NewNotFoundError("profile", map[string]interface{}{"id": id})
		}
		s.logger.Error("Failed to activate profile", "profile_id", id, "error", err)
		return errors.NewDatabaseError("profile activation", err)
	}

	s.logger.Info("Profile activated successfully", "profile_id", id)
	return nil
}

// SuspendProfile sets a profile's status to suspended
func (s *service) SuspendProfile(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("Suspending profile", "profile_id", id)

	err := s.repo.UpdateStatus(ctx, id, StatusSuspended)
	if err != nil {
		if errors.IsNotFoundErrorDomain(err) {
			s.logger.Warn("Profile not found for suspension", "profile_id", id)
			return errors.NewNotFoundError("profile", map[string]interface{}{"id": id})
		}
		s.logger.Error("Failed to suspend profile", "profile_id", id, "error", err)
		return errors.NewDatabaseError("profile suspension", err)
	}

	s.logger.Info("Profile suspended successfully", "profile_id", id)
	return nil
}

// DeactivateProfile sets a profile's status to inactive
func (s *service) DeactivateProfile(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("Deactivating profile", "profile_id", id)

	err := s.repo.UpdateStatus(ctx, id, StatusDeactivated)
	if err != nil {
		if errors.IsNotFoundErrorDomain(err) {
			s.logger.Warn("Profile not found for deactivation", "profile_id", id)
			return errors.NewNotFoundError("profile", map[string]interface{}{"id": id})
		}
		s.logger.Error("Failed to deactivate profile", "profile_id", id, "error", err)
		return errors.NewDatabaseError("profile deactivation", err)
	}

	s.logger.Info("Profile deactivated successfully", "profile_id", id)
	return nil
}

// AuthenticateProfile verifies login credentials and returns the profile if valid
func (s *service) AuthenticateProfile(ctx context.Context, req LoginRequest) (*PlatformProfile, error) {
	var profile *PlatformProfile
	var err error

	s.logger.Debug("Authenticating profile", "username", req.Username, "email", req.Email)

	// Lookup by username or email
	switch {
	case req.Username != "":
		profile, err = s.repo.GetProfileByUsername(ctx, req.Username)
	case req.Email != "":
		profile, err = s.repo.GetProfileByEmail(ctx, req.Email)
	default:
		s.logger.Warn("Username or email not provided")
		return nil, errors.NewValidationError("username or email is required", map[string]any{"field": "username_or_email"})
	}

	// Handle profile lookup errors
	if err != nil {
		if errors.IsNotFoundErrorDomain(err) {
			s.logger.Warn("Invalid credentials provided", "username", req.Username, "email", req.Email)
			return nil, errors.NewUnauthorizedError("invalid credentials")
		}
		s.logger.Error("Failed to fetch profile", "error", err)
		return nil, errors.NewDatabaseError("error fetching profile", err)
	}

	// Check profile status before login
	if profile.Status == StatusDeactivated || profile.Status == StatusSuspended {
		s.logger.Warn("Account is not active", "profile_id", profile.ID, "status", profile.Status)
		return nil, errors.NewUnauthorizedError("account is not active")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(profile.PasswordHash), []byte(req.Password))
	if err != nil {
		// Increment failed login attempts
		if errInc := s.repo.IncrementFailedLoginAttempts(ctx, profile.ID); errInc != nil {
			s.logger.Error("Failed to increment failed login attempts", "error", errInc)
		}

		// Check if account should be suspended due to too many failed attempts
		if profile.FailedLoginAttempts+1 >= 5 {
			if errSusp := s.repo.UpdateStatus(ctx, profile.ID, StatusSuspended); errSusp != nil {
				s.logger.Error("Failed to suspend profile after max failed attempts", "error", errSusp)
			}
			s.logger.Warn("Account suspended due to excessive failed login attempts", "profile_id", profile.ID)
			return nil, errors.NewUnauthorizedError("account suspended due to too many failed login attempts")
		}

		s.logger.Warn("Invalid credentials provided", "profile_id", profile.ID)
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	// Reset failed attempts on successful login
	if profile.FailedLoginAttempts > 0 {
		if err := s.repo.ResetFailedLoginAttempts(ctx, profile.ID); err != nil {
			s.logger.Warn("Failed to reset failed attempts", "profile_id", profile.ID, "error", err)
		}
	}

	// Update last login timestamp
	if err := s.repo.RecordLogin(ctx, profile.ID); err != nil {
		s.logger.Warn("Failed to record login", "profile_id", profile.ID, "error", err)
	}

	// Activate account if pending on first login
	if profile.Status == StatusPending {
		if err := s.repo.UpdateStatus(ctx, profile.ID, StatusActivated); err != nil {
			s.logger.Warn("Failed to activate pending account", "profile_id", profile.ID, "error", err)
		}
		profile.Status = StatusActivated
		s.logger.Info("Pending account activated upon first login", "profile_id", profile.ID)
	}

	// Remove password hash before returning profile
	profile.PasswordHash = ""
	s.logger.Info("Profile authenticated successfully", "profile_id", profile.ID)
	return profile, nil
}

// ChangePassword changes a user's password when previous password is provided/known
func (s *service) ChangePassword(ctx context.Context, id uuid.UUID, req PasswordChangeRequest) error {
	profile, err := s.repo.GetProfileByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(profile.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		return errors.NewUnauthorizedError("current password is incorrect")
	}

	// Ensure new password is different
	if req.CurrentPassword == req.NewPassword {
		return errors.NewValidationError(
			"new password must be different from current password",
			map[string]interface{}{"field": "new_password"},
		)
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return errors.NewBusinessError("PASSWORD_HASHING_FAILED", "failed to update password", nil)
	}

	// Update password in repository
	err = s.repo.UpdatePassword(ctx, id, string(passwordHash))
	if err != nil {
		s.logger.Error("failed to update password in database", "error", err)
		return errors.NewDatabaseError("updating password", err)
	}

	return nil
}

// RequestPasswordReset initiates the password reset process
func (s *service) RequestPasswordReset(ctx context.Context, req PasswordResetRequest) error {
	// Check if profile exists for the given email
	profile, err := s.repo.GetProfileByEmail(ctx, req.Email)
	if err != nil {
		// Do not reveal if email exists or not for security reasons
		s.logger.Info("password reset requested for non-existent email", "email", req.Email)
		return nil
	}

	// Generate reset token
	token, err := s.tokenGenerator.GenerateToken()
	if err != nil {
		s.logger.Error("failed to generate reset token", "error", err)
		return errors.NewBusinessError("RESET_TOKEN_GENERATION_FAILED", "failed to initiate password reset", nil)
	}

	// Store reset token with expiration (1 hours)
	expires := time.Now().Add(1 * time.Hour)

	passwordResetToken := PasswordResetToken{
		ProfileID: profile.ID,
		Token:     token,
		ExpiresAt: expires,
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreatePasswordResetToken(ctx, &passwordResetToken); err != nil {
		s.logger.Error("failed to save reset token", "error", err)
		return errors.NewBusinessError("RESET_TOKEN_SAVE_FAILED", "failed to initiate password reset", nil)
	}

	// Send reset link via email
	err = s.mailer.SendPasswordResetEmail(profile.Email, token)
	if err != nil {
		s.logger.Error("failed to send password reset email", "error", err)
		return errors.NewBusinessError("EMAIL_SEND_FAILED", "failed to send password reset email", nil)
	}

	s.logger.Info("password reset token generated and email sent", "profileID", profile.ID, "email", profile.Email)

	return nil
}

// ConfirmPasswordReset validates the reset token and updates the password
func (s *service) ConfirmPasswordReset(ctx context.Context, req PasswordResetConfirmation) error {
	// Validate token
	profile, err := s.repo.GetPasswordResetToken(ctx, req.Token)
	if err != nil {
		return errors.NewDatabaseError("fetching reset token", err)
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return errors.NewBusinessError(
			"PASSWORD_HASH_FAILED",
			"failed to update password",
			nil,
		)
	}

	// Update password
	if err := s.repo.UpdatePassword(ctx, profile.ProfileID, string(passwordHash)); err != nil {
		return errors.NewBusinessError(
			"PASSWORD_UPDATE_FAILED",
			"failed to update password",
			nil,
		)
	}

	// Mark token as used
	if err := s.repo.MarkPasswordResetTokenUsed(ctx, req.Token); err != nil {
		s.logger.Warn("failed to mark reset token as used", "error", err)
	}

	// Invalidate / Delete all other tokens
	if err := s.repo.DeleteOtherPasswordResetTokens(ctx, profile.ProfileID); err != nil {
		s.logger.Warn("failed to invalidate/delete other reset tokens", "error", err)
	}

	s.logger.Info("password successfully reset and token invalidated", "profileID", profile.ProfileID)
	return nil

}

func (s *service) GetSoftDeletedProfile(ctx context.Context, req GetSoftDeletedProfileRequest) (*PlatformProfile, error) {
	var profile *PlatformProfile
	var err error

	// Check if request contains an email
	if req.Email != nil {
		s.logger.Debug("Fetching soft-deleted profile by email", "email", *req.Email)
		profile, err = s.repo.GetSoftDeletedProfileByEmail(ctx, *req.Email)
		if err == nil && profile != nil {
			s.logger.Info("Soft-deleted profile found by email", "email", *req.Email)
			return profile, nil
		}
		if err != nil {
			s.logger.Warn("Error retrieving soft-deleted profile by email", "email", *req.Email, "error", err)
		}
	}

	// Check if request contains a username (only if email search failed or not provided)
	if req.Username != nil {
		s.logger.Debug("Fetching soft-deleted profile by username", "username", *req.Username)
		profile, err = s.repo.GetSoftDeletedProfileByUsername(ctx, *req.Username)
		if err == nil && profile != nil {
			s.logger.Info("Soft-deleted profile found by username", "username", *req.Username)
			return profile, nil
		}
		if err != nil {
			s.logger.Warn("Error retrieving soft-deleted profile by username", "username", *req.Username, "error", err)
		}
	}

	// Return error if no profile is found after all checks
	s.logger.Warn("Soft-deleted profile not found", "request", req)
	return nil, errors.NewNotFoundError("SoftDeletedProfile", map[string]any{
		"username": req.Username,
		"email":    req.Email,
	})
}

// RestoreProfile restores a soft-deleted profile based on the given identifier (ID, username, or email)
func (s *service) RestoreProfile(ctx context.Context, id uuid.UUID) error {
	var profile *PlatformProfile
	var err error

	s.logger.Debug("Restoring soft-deleted profile by ID", "id", id)

	err = s.repo.RestoreSoftDeletedProfile(ctx, id)

	// Handle errors if the profile is not found or another issue occurs
	if err != nil {
		if errors.IsNotFoundErrorDomain(err) {
			s.logger.Warn("Profile not found in deleted_profiles", "error", err)
			return errors.NewNotFoundError("profile", map[string]interface{}{"id": id})
		}
		s.logger.Error("Failed to get soft-deleted profile", "error", err)
		return errors.NewBusinessError(
			"PROFILE_RESTORE_FAILED",
			"failed to restore profile",
			nil,
		)
	}

	s.logger.Info("Profile restored successfully", "id", profile.ID)
	return nil
}

// HardDeleteProfile hard-deletes a profile
func (s *service) HardDeleteProfile(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("initiating hard delete of profile", "profileID", id)

	err := s.repo.HardDeleteProfile(ctx, id)
	if err != nil {
		s.logger.Error("failed to hard delete profile", "profileID", id, "error", err)
		return errors.NewBusinessError(
			"PROFILE_HARD_DELETE_FAILED",
			"failed to hard delete profile",
			map[string]interface{}{"profileID": id},
		)
	}

	s.logger.Info("profile hard deleted successfully", "profileID", id)
	return nil
}

// AssignRole assigns a role to a profile
func (s *service) AssignRole(ctx context.Context, profileID uuid.UUID, roleID uuid.UUID) error {
	// Verify profile exists
	_, err := s.repo.GetProfileByID(ctx, profileID)
	if err != nil {
		return err
	}

	// Verify role exists
	_, err = s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Check if already assigned
	hasRole, err := s.repo.HasRoleAssignment(ctx, profileID, roleID)
	if err != nil {
		return err
	}

	if hasRole {
		return errors.NewDomainError("role already assigned to profile", errors.ValidationError)
	}

	return s.repo.AssignRole(ctx, profileID, roleID)
}

// RemoveRole removes a role from a profile
func (s *service) RemoveRole(ctx context.Context, profileID uuid.UUID, roleID uuid.UUID) error {
	// Verify assignment exists
	hasRole, err := s.repo.HasRoleAssignment(ctx, profileID, roleID)
	if err != nil {
		return err
	}

	if !hasRole {
		return errors.NewDomainError("role not assigned to profile", errors.ValidationError)
	}

	return s.repo.RemoveRole(ctx, profileID, roleID)
}

// HasRole checks if a profile has a specific role by name
func (s *service) HasRole(ctx context.Context, profileID uuid.UUID, roleName string) (bool, error) {
	// Get role ID from name
	role, err := s.roleRepo.GetRoleByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return s.repo.HasRoleAssignment(ctx, profileID, role.ID)
}

// GetPreferences retrieves a profile's preferences
func (s *service) GetPreferences(ctx context.Context, profileID uuid.UUID) (*ProfilePreference, error) {

	// Fetch preferences from the repository
	prefs, err := s.repo.GetPreferences(ctx, profileID)
	if err != nil {
		s.logger.Error("failed to fetch preferences", "profileID", profileID, "error", err)
		return nil, errors.NewDatabaseError("fetching preferences", err)
	}

	s.logger.Info("preferences retrieved successfully", "profileID", profileID)
	return prefs, nil
}

// UpdatePreferences updates a profile's preferences
func (s *service) UpdatePreferences(ctx context.Context, profileID uuid.UUID, prefs ProfilePreference) error {

	s.logger.Info("updating existing preferences", "profileID", profileID)

	if err := s.repo.SavePreferences(ctx, &prefs); err != nil {
		s.logger.Error("failed to update preferences", "profileID", profileID, "error", err)
		return errors.NewBusinessError(
			"PREFERENCES_UPDATE_FAILED",
			"failed to update preferences",
			map[string]interface{}{"profileID": profileID},
		)
	}

	s.logger.Info("preferences updated successfully", "profileID", profileID)
	return nil
}

// SearchProfiles searches for profiles matching a query string
func (s *service) SearchProfiles(ctx context.Context, query string, page, pageSize int) ([]*PlatformProfile, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	// Search profiles
	profiles, total, err := s.repo.SearchProfiles(ctx, query, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Remove password hashes from results
	for _, profile := range profiles {
		profile.PasswordHash = ""
	}

	return profiles, total, nil
}

// ListProfilesByRole retrieves profiles that have a specific role
func (s *service) ListProfilesByRole(ctx context.Context, roleID uuid.UUID, page, pageSize int) ([]*PlatformProfile, int, error) {
	// Verify role exists
	_, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, 0, err
	}

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	// Get profiles by role
	profiles, total, err := s.repo.GetProfilesByRole(ctx, roleID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Remove password hashes from results
	for _, profile := range profiles {
		profile.PasswordHash = ""
	}

	return profiles, total, nil
}
