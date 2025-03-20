// 2. TODO: Update the Soft and Hard Delete implementations to include `shifting
// deleted profiles to different recovery table` which is cleared by workers periodically
// and implement Restore Profile Function.
// 3. TODO: Update the schemas and table locations for tables not to be in public schema like profile_preferences, password_reset_tokens, recovery table, etc.
// 4. TODO: Update to include logging in all the functions.
// 5. TODO: Update if ProfilePreference model is changed to include other notifications.
//

// Soft Delete Profile archives the profile to deleted_profiles and removes it from platform_profiles
// The deleted_profiles table has deleted_at timestamp to track when the profile was deleted

package repositories

import (
	"context"
	"errors"
	"fmt"
	"server/internal/domain/platform_profile"
	"server/pkg/logger"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrProfileNotFound       = errors.New("profile not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrResetTokenNotFound    = errors.New("reset token not found")
	ErrResetTokenExpired     = errors.New("reset token expired")
	ErrResetTokenUsed        = errors.New("reset token already used")
	ErrRoleNotFound          = errors.New("role not found")
	ErrRoleNotAssigned       = errors.New("role not assigned or already removed")
)

// PostgresProfileRepository implements the platform_profile.Repository interface
type PostgresProfileRepository struct {
	pool   *pgxpool.Pool
	logger *logger.Logger // Add the logger
}

// NewPostgresProfileRepository creates a new PostgreSQL-backed profile repository
func NewPostgresProfileRepository(pool *pgxpool.Pool, logger *logger.Logger) platform_profile.Repository {
	return &PostgresProfileRepository{
		pool:   pool,
		logger: logger,
	}
}

// UsernameExists checks if a platform profile with the same username exists in the database
func (r *PostgresProfileRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	r.logger.Debug("Checking if username exists", "username", username)

	query := "SELECT EXISTS(SELECT 1 FROM platform_profiles WHERE username = $1)"
	var exists bool
	err := r.pool.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check username existence", "username", username, "error", err)
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	if exists {
		r.logger.Info("Username exists", "username", username)
	} else {
		r.logger.Debug("Username does not exist", "username", username)
	}
	return exists, nil
}

// EmailExists checks if a profile with the given email already exists
func (r *PostgresProfileRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	r.logger.Debug("Checking if email exists", "email", email)

	query := "SELECT EXISTS(SELECT 1 FROM platform_profiles WHERE email = $1)"
	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check email existence", "email", email, "error", err)
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	if exists {
		r.logger.Info("Email exists", "email", email)
	} else {
		r.logger.Debug("Email does not exist", "email", email)
	}
	return exists, nil
}

// CreateProfile stores a new platform profile in the database
func (r *PostgresProfileRepository) CreateProfile(ctx context.Context, profile *platform_profile.PlatformProfile) error {

	r.logger.Debug(
		"Creating new profile",
		"username", profile.Username,
		"email", profile.Email,
	)

	query := `
	INSERT INTO platform_profiles (
		id, username, email, password_hash, status, verified_at, 
		last_login_at, failed_login_attempts, created_at, updated_by_user_at,
		updated_by_system_at		
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	) RETURNING id`

	// Execute query
	var returnedID uuid.UUID
	err := r.pool.QueryRow(
		ctx,
		query,
		profile.ID,
		profile.Username,
		profile.Email,
		profile.PasswordHash,
		profile.Status,
		profile.VerifiedAt,
		profile.LastLoginAt,
		profile.FailedLoginAttempts,
		profile.CreatedAt,
		profile.UpdatedByUserAt,
		profile.UpdatedBySystemAt,
	).Scan(&returnedID)

	if err != nil {
		// Check for unique constraint violations
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "platform_profiles_username_key":
				r.logger.Warn("Username already exists",
					"username", profile.Username)
				return ErrUsernameAlreadyExists
			case "platform_profiles_email_key":
				r.logger.Warn("Email already exists",
					"email", profile.Email)
				return ErrEmailAlreadyExists
			}
		}
		r.logger.Error("Failed to create profile",
			"username", profile.Username,
			"error", err)
		return fmt.Errorf("failed to create profile: %w", err)
	}

	r.logger.Info(
		"Profile created successfully",
		"profileID", returnedID,
		"username", profile.Username,
	)

	profile.ID = returnedID
	return nil
}

// GetProfileByID retrieves a profile by its ID
func (r *PostgresProfileRepository) GetProfileByID(ctx context.Context, id uuid.UUID) (*platform_profile.PlatformProfile, error) {

	query := `
	SELECT 
		id, username, email, password_hash, status, verified_at, 
		last_login_at, failed_login_attempts, created_at, updated_by_user_at,
		updated_by_system_at
	FROM platform_profiles 
	WHERE id = $1`

	profile := &platform_profile.PlatformProfile{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.PasswordHash,
		&profile.Status,
		&profile.VerifiedAt,
		&profile.LastLoginAt,
		&profile.FailedLoginAttempts,
		&profile.CreatedAt,
		&profile.UpdatedByUserAt,
		&profile.UpdatedBySystemAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("failed to get profile by ID: %w", err)
	}

	return profile, nil
}

// GetProfileByEmail retrieves a profile by its email
func (r *PostgresProfileRepository) GetProfileByEmail(ctx context.Context, email string) (*platform_profile.PlatformProfile, error) {

	query := `
	SELECT 
		id, username, email, password_hash, status, verified_at, 
		last_login_at, failed_login_attempts, created_at, updated_by_user_at,
		updated_by_system_at
	FROM platform_profiles 
	WHERE email = $1`

	profile := &platform_profile.PlatformProfile{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.PasswordHash,
		&profile.Status,
		&profile.VerifiedAt,
		&profile.LastLoginAt,
		&profile.FailedLoginAttempts,
		&profile.CreatedAt,
		&profile.UpdatedByUserAt,
		&profile.UpdatedBySystemAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("failed to get profile by email: %w", err)
	}

	return profile, nil
}

// GetProfileByUsername retrieves a profile by its username
func (r *PostgresProfileRepository) GetProfileByUsername(ctx context.Context, username string) (*platform_profile.PlatformProfile, error) {

	query := `
	SELECT 
		id, username, email, password_hash, status, verified_at, 
		last_login_at, failed_login_attempts, created_at, updated_by_user_at,
		updated_by_system_at
	FROM platform_profiles 
	WHERE username = $1`

	profile := &platform_profile.PlatformProfile{}
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.PasswordHash,
		&profile.Status,
		&profile.VerifiedAt,
		&profile.LastLoginAt,
		&profile.FailedLoginAttempts,
		&profile.CreatedAt,
		&profile.UpdatedByUserAt,
		&profile.UpdatedBySystemAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("failed to get profile by username: %w", err)
	}

	return profile, nil
}

// UpdateProfile updates an existing platform profile
func (r *PostgresProfileRepository) UpdateProfile(ctx context.Context, profile *platform_profile.PlatformProfile) error {
	r.logger.Debug(
		"Updating profile",
		"id", profile.ID,
		"username", profile.Username,
		"email", profile.Email,
	)

	query := `
	   UPDATE platform_profiles SET
			username = $1,
			email = $2,			
			updated_by_user_at = $3,
			updated_by_system_at = $4
	   WHERE id = $5`

	commandTag, err := r.pool.Exec(
		ctx,
		query,
		profile.Username,
		profile.Email,
		profile.UpdatedByUserAt,
		profile.UpdatedBySystemAt,
		profile.ID,
	)

	if err != nil {
		// Check for unique constraint violations
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "platform_profiles_username_key":
				r.logger.Warn("Username already exists", "username", profile.Username)
				return ErrUsernameAlreadyExists
			case "platform_profiles_email_key":
				r.logger.Warn("Email already exists", "email", profile.Email)
				return ErrEmailAlreadyExists
			}
		}
		r.logger.Error("Failed to update profile", "id", profile.ID, "error", err)
		return fmt.Errorf("failed to update profile: %w", err)
	}

	// Check if no rows were affected
	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile not found for update", "id", profile.ID)
		return ErrProfileNotFound
	}

	r.logger.Info("Profile updated successfully", "id", profile.ID)
	return nil
}

// TODO: 2, 3

// SoftDeleteProfile archives the profile to deleted_profiles and removes it from platform_profiles
func (r *PostgresProfileRepository) SoftDeleteProfile(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Starting profile soft delete process", "id", id)

	// Begin a transaction to ensure atomicity
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.Error(
			"Failed to start transaction for soft delete",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("failed to start transaction for soft delete: %w", err)
	}
	defer tx.Rollback(ctx)

	// Archive the profile to the deleted_profiles table
	archiveQuery := `
	INSERT INTO profile_schema.deleted_profiles (
		id, username, email, password_hash, status, verified_at, 
		last_login_at, failed_login_attempts, created_at, updated_by_user_at,
		updated_by_system_at, deleted_at
	)
	SELECT 
		id, username, email, password_hash, status, verified_at, 
		last_login_at, failed_login_attempts, created_at, updated_by_user_at,
		updated_by_system_at, NOW()
	FROM platform_profiles
	WHERE id = $1
	`

	commandTag, err := tx.Exec(ctx, archiveQuery, id)
	if err != nil {
		r.logger.Error(
			"Failed to archive profile before soft delete",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("failed to archive profile before soft delete: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile not found for soft delete", "id", id)
		return ErrProfileNotFound
	}

	r.logger.Info("Profile archived successfully", "id", id)

	// Permanently remove the profile from platform_profiles
	deleteQuery := `
	DELETE FROM platform_profiles WHERE id = $1
	`

	commandTag, err = tx.Exec(ctx, deleteQuery, id)
	if err != nil {
		r.logger.Error("Failed to delete profile after archiving", "id", id, "error", err)
		return fmt.Errorf("failed to delete profile after archiving: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile already deleted or not found", "id", id)
		return ErrProfileNotFound
	}

	// Commit the transaction after successful archiving and deletion
	if err = tx.Commit(ctx); err != nil {
		r.logger.Error("Failed to commit transaction for soft delete", "id", id, "error", err)
		return fmt.Errorf("failed to commit transaction for soft delete: %w", err)
	}

	r.logger.Info("Profile archived and removed from platform_profiles successfully", "id", id)
	return nil
}

// VerifyProfile marks a profile as verified by setting VerifiedAt timestamp
func (r *PostgresProfileRepository) VerifyProfile(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Starting profile verification process", "id", id)

	query := `
	UPDATE platform_profiles SET
		status = 'activated',
		verified_at = $1,
		updated_by_user_at = $1,
		updated_by_system_at = $1
	WHERE id = $2
	`

	now := time.Now()

	// Log query execution
	r.logger.Debug("Executing profile verification query", "id", id, "query", query)

	// Execute query
	commandTag, err := r.pool.Exec(ctx, query, now, id)
	if err != nil {
		r.logger.Error(
			"Failed to verify profile",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("failed to verify profile: %w", err)
	}

	// Check if the profile was found and updated
	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile not found or already verified", "id", id)
		return ErrProfileNotFound
	}

	r.logger.Info("Profile verified successfully", "id", id)
	return nil
}

// UpdateStatus updates the status of a profile
func (r *PostgresProfileRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status platform_profile.Status) error {
	r.logger.Debug("Updating profile status", "id", id, "status", status)

	query := `
	UPDATE platform_profiles SET
		status = $1,
		updated_by_user_at = $2,
		updated_by_system_at = $2
	WHERE id = $3`

	// Execute the update query
	commandTag, err := r.pool.Exec(ctx, query, status, time.Now(), id)
	if err != nil {
		r.logger.Error("Failed to update profile status", "id", id, "status", status, "error", err)
		return fmt.Errorf("failed to update profile status: %w", err)
	}

	// Check if no rows were affected (i.e., profile not found)
	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile not found while updating status", "id", id)
		return ErrProfileNotFound
	}

	r.logger.Info("Profile status updated successfully", "id", id, "status", status)
	return nil
}

// CreatePasswordResetToken stores a new password reset token
func (r *PostgresProfileRepository) CreatePasswordResetToken(ctx context.Context, resetToken *platform_profile.PasswordResetToken) error {
	r.logger.Debug(
		"Starting to create password reset token",
		"profile_id", resetToken.ProfileID,
	)

	query := `
	INSERT INTO profile_schema.password_reset_tokens (
		profile_id, token, expires_at, is_used, created_at
	) VALUES (
		$1, $2, $3, $4, $5
	)`

	// Execute the insert query
	_, err := r.pool.Exec(
		ctx,
		query,
		resetToken.ProfileID,
		resetToken.Token,
		resetToken.ExpiresAt,
		resetToken.IsUsed,
		resetToken.CreatedAt,
	)

	if err != nil {
		r.logger.Error(
			"Failed to create password reset token",
			"profile_id", resetToken.ProfileID,
			"error", err,
		)
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	r.logger.Info(
		"Password reset token created successfully",
		"profile_id", resetToken.ProfileID,
		"expires_at", resetToken.ExpiresAt,
	)
	return nil
}

// GetPasswordResetToken retrieves a password reset token
func (r *PostgresProfileRepository) GetPasswordResetToken(ctx context.Context, token string) (*platform_profile.PasswordResetToken, error) {
	r.logger.Debug("Fetching password reset token", "token", token)

	query := `
	SELECT 
		profile_id, token, expires_at, is_used, created_at
	FROM profile_schema.password_reset_tokens
	WHERE token = $1`

	resetToken := &platform_profile.PasswordResetToken{}
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&resetToken.ProfileID,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.IsUsed,
		&resetToken.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("Password reset token not found", "token", token)
			return nil, ErrResetTokenNotFound
		}
		r.logger.Error("Failed to fetch password reset token", "token", token, "error", err)
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}

	r.logger.Info("Password reset token fetched successfully", "token", token, "profile_id", resetToken.ProfileID)
	return resetToken, nil
}

// MarkPasswordResetTokenUsed marks a token as used
func (r *PostgresProfileRepository) MarkPasswordResetTokenUsed(ctx context.Context, token string) error {
	r.logger.Debug("Marking password reset token as used", "token", token)

	query := `
	UPDATE profile_schema.password_reset_tokens SET
		is_used = true,
		updated_at = NOW()
	WHERE token = $1 AND is_used = false`

	commandTag, err := r.pool.Exec(ctx, query, token)
	if err != nil {
		r.logger.Error("Failed to mark password reset token as used", "token", token, "error", err)
		return fmt.Errorf("failed to mark password reset token as used: %w", err)
	}

	// Check if any rows were affected
	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Password reset token not found or already used", "token", token)
		return ErrResetTokenNotFound
	}

	r.logger.Info("Password reset token marked as used successfully", "token", token)
	return nil
}

// UpdatePassword updates a profile's password hash
func (r *PostgresProfileRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	r.logger.Debug("Starting password update process", "id", id)

	query := `
	UPDATE platform_profiles SET
		password_hash = $1,
		updated_by_user_at = $2,
		updated_by_system_at = $2
	WHERE id = $3
	`

	now := time.Now()

	// Log query execution
	r.logger.Debug("Executing password update query", "id", id, "query", query)

	// Execute query
	commandTag, err := r.pool.Exec(ctx, query, passwordHash, now, id)
	if err != nil {
		r.logger.Error(
			"Failed to update password",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Check if the profile was found and updated
	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile not found for password update", "id", id)
		return ErrProfileNotFound
	}

	r.logger.Info("Password updated successfully", "id", id)
	return nil
}

// DeleteOtherResetTokens deletes all other reset tokens for the profile except the used one
func (r *PostgresProfileRepository) DeleteOtherPasswordResetTokens(ctx context.Context, profileID uuid.UUID) error {
	r.logger.Debug("Starting to delete other password reset tokens", "profile_id", profileID)

	query := `
	DELETE FROM profile_schema.password_reset_tokens
	WHERE profile_id = $1
	`

	commandTag, err := r.pool.Exec(ctx, query, profileID)
	if err != nil {
		r.logger.Error(
			"Failed to delete other password reset tokens",
			"profile_id", profileID,
			"error", err,
		)
		return fmt.Errorf("failed to delete other password reset tokens: %w", err)
	}

	r.logger.Info(
		"Other password reset tokens deleted successfully",
		"profile_id", profileID,
		"tokens_deleted", commandTag.RowsAffected(),
	)

	return nil
}

// GetSoftDeletedProfileByEmail retrieves a soft-deleted profile by email
func (r *PostgresProfileRepository) GetSoftDeletedProfileByEmail(ctx context.Context, email string) (*platform_profile.PlatformProfile, error) {
	r.logger.Debug("Fetching soft-deleted profile by email", "email", email)

	query := `
	SELECT id, username, email, password_hash, status, verified_at, 
	       last_login_at, failed_login_attempts, created_at, updated_by_user_at,
	       updated_by_system_at
	FROM profile_schema.deleted_profiles
	WHERE email = $1
	`

	var profile platform_profile.PlatformProfile
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.PasswordHash,
		&profile.Status,
		&profile.VerifiedAt,
		&profile.LastLoginAt,
		&profile.FailedLoginAttempts,
		&profile.CreatedAt,
		&profile.UpdatedByUserAt,
		&profile.UpdatedBySystemAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("No soft-deleted profile found with the given email", "email", email)
			return nil, ErrProfileNotFound
		}
		r.logger.Error("Failed to fetch soft-deleted profile by email", "email", email, "error", err)
		return nil, ErrProfileNotFound
	}

	r.logger.Info("Soft-deleted profile fetched successfully by email", "email", email)
	return &profile, nil
}

// GetSoftDeletedProfileByUsername retrieves a soft-deleted profile by username
func (r *PostgresProfileRepository) GetSoftDeletedProfileByUsername(ctx context.Context, username string) (*platform_profile.PlatformProfile, error) {
	r.logger.Debug("Fetching soft-deleted profile by username", "username", username)

	query := `
	SELECT id, username, email, password_hash, status, verified_at, 
	       last_login_at, failed_login_attempts, created_at, updated_by_user_at,
	       updated_by_system_at
	FROM profile_schema.deleted_profiles
	WHERE username = $1
	`

	var profile platform_profile.PlatformProfile
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.PasswordHash,
		&profile.Status,
		&profile.VerifiedAt,
		&profile.LastLoginAt,
		&profile.FailedLoginAttempts,
		&profile.CreatedAt,
		&profile.UpdatedByUserAt,
		&profile.UpdatedBySystemAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("No soft-deleted profile found with the given username", "username", username)
			return nil, ErrProfileNotFound
		}
		r.logger.Error("Failed to fetch soft-deleted profile by username", "username", username, "error", err)
		return nil, ErrProfileNotFound
	}

	r.logger.Info("Soft-deleted profile fetched successfully by username", "username", username)
	return &profile, nil
}

// RestoreSoftDeletedProfile restores a soft-deleted profile by id
func (r *PostgresProfileRepository) RestoreSoftDeletedProfile(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Starting profile restoration process", "id", id)

	// Begin a transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.Error(
			"Failed to start transaction for profile restoration",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("failed to start transaction for profile restoration: %w", err)
	}
	defer tx.Rollback(ctx)

	// Restore profile from deleted_profiles to platform_profiles
	restoreQuery := `
	INSERT INTO platform_profiles (
		id, username, email, password_hash, status, verified_at, 
		last_login_at, failed_login_attempts, created_at, updated_by_user_at,
		updated_by_system_at
	)
	SELECT 
		id, username, email, password_hash, status, verified_at, 
		last_login_at, failed_login_attempts, created_at, updated_by_user_at,
		updated_by_system_at
	FROM profile_schema.deleted_profiles
	WHERE id = $1
	ON CONFLICT (id) DO NOTHING
	RETURNING id
	`

	var restoredID uuid.UUID
	err = tx.QueryRow(ctx, restoreQuery, id).Scan(&restoredID)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("No soft-deleted profile found for restoration", "id", id)
			return ErrProfileNotFound
		}
		r.logger.Error("Failed to restore profile from deleted_profiles", "id", id, "error", err)
		return fmt.Errorf("failed to restore profile: %w", err)
	}

	r.logger.Info("Profile restored successfully", "id", id, "restoredID", restoredID)

	// Delete the restored profile from deleted_profiles
	deleteQuery := `
	DELETE FROM profile_schema.deleted_profiles WHERE id = $1
	`

	commandTag, err := tx.Exec(ctx, deleteQuery, id)
	if err != nil {
		r.logger.Error(
			"Failed to delete restored profile from deleted_profiles",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("failed to delete restored profile from deleted_profiles: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile was not found or already restored", "id", id)
		return ErrProfileNotFound
	}

	// Commit transaction after successful restoration and deletion
	if err = tx.Commit(ctx); err != nil {
		r.logger.Error(
			"Failed to commit transaction for profile restoration",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("failed to commit transaction for profile restoration: %w", err)
	}

	r.logger.Info("Profile restoration completed successfully", "id", id)
	return nil
}

// HardDeleteProfile removes a profile from the deleted_profiles table in profile_schema
func (r *PostgresProfileRepository) HardDeleteProfile(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Starting hard delete process for profile", "id", id)

	// Delete query to remove the profile from deleted_profiles table
	query := `DELETE FROM profile_schema.deleted_profiles WHERE id = $1`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete profile from deleted_profiles", "id", id, "error", err)
		return fmt.Errorf("failed to delete profile from deleted_profiles: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile not found in deleted_profiles for hard delete", "id", id)
		return ErrProfileNotFound
	}

	r.logger.Info("Profile successfully deleted from deleted_profiles", "id", id)
	return nil
}

// AssignRoleToProfile assigns a role to a profile in the profile_roles table
func (r *PostgresProfileRepository) AssignRoleToProfile(ctx context.Context, profileID, roleID uuid.UUID) error {

	r.logger.Debug(
		"Starting to assign role to profile",
		"profile_id",
		profileID,
		"role_id",
		roleID,
	)

	query := `
	INSERT INTO profile_schema.profile_roles (
		profile_id, role_id, created_at
	) VALUES (
		$1, $2, $3
	) ON CONFLICT (profile_id, role_id) DO NOTHING
	`

	// Execute the query to assign the role
	_, err := r.pool.Exec(ctx, query, profileID, roleID, time.Now())
	if err != nil {
		r.logger.Error(
			"Failed to assign role to profile",
			"profile_id", profileID,
			"role_id", roleID,
			"error", err,
		)
		return fmt.Errorf("failed to assign role to profile: %w", err)
	}

	r.logger.Info("Role assigned to profile successfully", "profile_id", profileID, "role_id", roleID)
	return nil
}

// GetProfileRoles retrieves all roles assigned to a profile
func (r *PostgresProfileRepository) GetProfileRoles(ctx context.Context, profileID uuid.UUID) ([]uuid.UUID, error) {
	r.logger.Debug("Fetching roles assigned to profile", "profile_id", profileID)

	query := `
	SELECT role_id
	FROM profile_schema.profile_roles
	WHERE profile_id = $1`

	// Query to get roles
	rows, err := r.pool.Query(ctx, query, profileID)
	if err != nil {
		r.logger.Error("Failed to query profile roles", "profile_id", profileID, "error", err)
		return nil, fmt.Errorf("failed to get profile roles: %w", err)
	}
	defer rows.Close()

	var roleIDs []uuid.UUID

	// Iterate over the result set
	for rows.Next() {
		var roleID uuid.UUID
		if err := rows.Scan(&roleID); err != nil {
			r.logger.Error("Failed to scan role ID", "profile_id", profileID, "error", err)
			return nil, fmt.Errorf("failed to scan role ID: %w", err)
		}
		roleIDs = append(roleIDs, roleID)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over role rows", "profile_id", profileID, "error", err)
		return nil, fmt.Errorf("error iterating over role rows: %w", err)
	}

	// Log if no roles are found
	if len(roleIDs) == 0 {
		r.logger.Warn("No roles found for profile", "profile_id", profileID)
	}

	r.logger.Info("Roles retrieved successfully for profile", "profile_id", profileID, "role_count", len(roleIDs))
	return roleIDs, nil
}

// RemoveRoleFromProfile removes a role from a profile
func (r *PostgresProfileRepository) RemoveRoleFromProfile(ctx context.Context, profileID, roleID uuid.UUID) error {
	r.logger.Debug("Starting to remove role from profile", "profile_id", profileID, "role_id", roleID)

	query := `
	DELETE FROM profile_schema.profile_roles
	WHERE profile_id = $1 AND role_id = $2`

	// Execute the query to remove the role
	commandTag, err := r.pool.Exec(ctx, query, profileID, roleID)
	if err != nil {
		r.logger.Error(
			"Failed to remove role from profile",
			"profile_id", profileID,
			"role_id", roleID,
			"error", err,
		)
		return fmt.Errorf("failed to remove role from profile: %w", err)
	}

	// Check if any rows were affected
	if commandTag.RowsAffected() == 0 {
		r.logger.Warn(
			"Role not found for profile or already removed",
			"profile_id", profileID,
			"role_id", roleID,
		)
		return ErrRoleNotAssigned
	}

	r.logger.Info("Role removed successfully from profile", "profile_id", profileID, "role_id", roleID)
	return nil
}

// TODO: 5
// SavePreferences saves profile preferences
func (r *PostgresProfileRepository) SavePreferences(ctx context.Context, prefs *platform_profile.ProfilePreference) error {

	r.logger.Debug(
		"Saving profile preferences",
		"profileID", prefs.ProfileID,
	)

	// Check if preference already exists
	existingPrefs, err := r.GetPreferences(ctx, prefs.ProfileID)

	if err != nil && !errors.Is(err, ErrProfileNotFound) {
		r.logger.Error(
			"Failed to check existing preferences",
			"profileID", prefs.ProfileID,
			"error", err,
		)
		return err
	}

	var query string
	if existingPrefs != nil {
		// Update existing preferences
		query = `
		UPDATE profile_schema.profile_preferences SET 
			notifications_email = $1, 			
			updated_at = $2 
		WHERE profile_id = $3`

		_, err = r.pool.Exec(
			ctx,
			query,
			prefs.NotificationsEmail,
			time.Now(),
			prefs.ProfileID,
		)
	} else {
		// Insert new preferences
		query = `
		INSERT INTO profile_schema.profile_preferences ( 
			id, profile_id, notifications_email, created_at, updated_at 
		) VALUES ( 
			$1, $2, $3, $4, $5
		)`

		now := time.Now()

		// Generate a new UUID if ID is not provided
		if prefs.ID == uuid.Nil {
			prefs.ID = uuid.New()
		}

		_, err = r.pool.Exec(
			ctx,
			query,
			prefs.ID,
			prefs.ProfileID,
			prefs.NotificationsEmail,
			now,
			now,
		)
	}

	if err != nil {
		r.logger.Error(
			"Failed to save profile preferences",
			"profileID", prefs.ProfileID,
			"error", err,
		)
		return fmt.Errorf("failed to save profile preferences: %w", err)
	}

	r.logger.Info(
		"Profile preferences saved successfully",
		"profileID", prefs.ProfileID,
	)

	return nil
}

// TODO: 5
// GetPreferences retrieves profile preferences
func (r *PostgresProfileRepository) GetPreferences(ctx context.Context, profileID uuid.UUID) (*platform_profile.ProfilePreference, error) {

	r.logger.Debug(
		"Getting profile preferences",
		"profileID", profileID,
	)

	query := `
	SELECT  
		id,
		profile_id, 
		notifications_email,
		created_at, 
		updated_at 
	FROM profile_schema.profile_preferences  
	WHERE profile_id = $1`

	prefs := &platform_profile.ProfilePreference{}
	err := r.pool.QueryRow(ctx, query, profileID).Scan(
		&prefs.ID,
		&prefs.ProfileID,
		&prefs.NotificationsEmail,
		&prefs.CreatedAt,
		&prefs.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug(
				"Profile preferences not found",
				"profileID", profileID,
			)
			return nil, ErrProfileNotFound
		}
		r.logger.Error(
			"Failed to get profile preferences",
			"profileID", profileID,
			"error", err,
		)
		return nil, fmt.Errorf("failed to get profile preferences: %w", err)
	}

	r.logger.Debug(
		"Profile preferences retrieved successfully",
		"profileID", profileID,
	)

	return prefs, nil
}

// GetProfiles retrieves profiles with pagination and filtering
func (r *PostgresProfileRepository) GetProfiles(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*platform_profile.PlatformProfile, int, error) {
	r.logger.Debug("Fetching profiles with pagination and filters", "offset", offset, "limit", limit, "filters", filters)

	// Base queries for profiles and count
	baseQuery := `
	SELECT 
		id, username, email, password_hash, first_name, last_name, 
		display_name, status, created_at, updated_by_user_at, 
		last_login_at, failed_login_attempts, verified_at
	FROM profile_schema.platform_profiles`

	countQuery := `SELECT COUNT(*) FROM profile_schema.platform_profiles`

	// Build WHERE clause for filters
	whereClause := ""
	args := []interface{}{}
	paramIndex := 1

	if len(filters) > 0 {
		whereClause = " WHERE "
		conditions := []string{}

		for key, value := range filters {
			switch key {
			case "status":
				conditions = append(conditions, fmt.Sprintf("status = $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			case "verified":
				conditions = append(conditions, fmt.Sprintf("verified_at IS NOT NULL AND verified_at <= $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			case "search":
				search := fmt.Sprintf("%%%s%%", value)
				conditions = append(conditions, fmt.Sprintf("(username ILIKE $%d OR email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)",
					paramIndex, paramIndex+1, paramIndex+2, paramIndex+3))
				args = append(args, search, search, search, search)
				paramIndex += 4
			case "created_after":
				conditions = append(conditions, fmt.Sprintf("created_at > $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			case "created_before":
				conditions = append(conditions, fmt.Sprintf("created_at < $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			}
		}

		whereClause += strings.Join(conditions, " AND ")
	}

	// Build final queries
	dataQuery := baseQuery + whereClause + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	countQueryFinal := countQuery + whereClause

	// Add pagination parameters
	args = append(args, limit, offset)

	// Get total count of records
	var total int
	err := r.pool.QueryRow(ctx, countQueryFinal, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		r.logger.Error("Failed to get profiles count", "error", err)
		return nil, 0, fmt.Errorf("failed to get profiles count: %w", err)
	}

	// Execute data query
	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		r.logger.Error("Failed to get profiles", "error", err)
		return nil, 0, fmt.Errorf("failed to get profiles: %w", err)
	}
	defer rows.Close()

	// Process query results
	profiles := []*platform_profile.PlatformProfile{}
	for rows.Next() {
		profile := &platform_profile.PlatformProfile{}
		err := rows.Scan(
			&profile.ID,
			&profile.Username,
			&profile.Email,
			&profile.PasswordHash,
			&profile.Status,
			&profile.VerifiedAt,
			&profile.LastLoginAt,
			&profile.FailedLoginAttempts,
			&profile.CreatedAt,
			&profile.UpdatedByUserAt,
			&profile.UpdatedBySystemAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan profile", "error", err)
			return nil, 0, fmt.Errorf("failed to scan profile: %w", err)
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over profile rows", "error", err)
		return nil, 0, fmt.Errorf("error iterating over rows: %w", err)
	}

	r.logger.Info("Profiles fetched successfully", "total_profiles", len(profiles))
	return profiles, total, nil
}

// GetProfilesByRole retrieves profiles assigned to a specific role with pagination
func (r *PostgresProfileRepository) GetProfilesByRole(ctx context.Context, roleID uuid.UUID, page, pageSize int) ([]*platform_profile.PlatformProfile, int, error) {
	r.logger.Debug("Fetching profiles by role", "role_id", roleID, "page", page, "page_size", pageSize)

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Base query to get profiles by role
	baseQuery := `
	SELECT 
		p.id, p.username, p.email, p.password_hash, p.first_name, 
		p.last_name, p.display_name, p.status, p.created_at, 
		p.updated_by_user_at, p.last_login_at, p.failed_login_attempts, p.verified_at
	FROM profile_schema.platform_profiles AS p
	JOIN profile_schema.profile_roles AS pr ON p.id = pr.profile_id
	WHERE pr.role_id = $1
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3`

	// Query to count total profiles with the given role
	countQuery := `
	SELECT COUNT(*)
	FROM profile_schema.platform_profiles AS p
	JOIN profile_schema.profile_roles AS pr ON p.id = pr.profile_id
	WHERE pr.role_id = $1`

	// Get the total count of profiles with the given role
	var total int
	err := r.pool.QueryRow(ctx, countQuery, roleID).Scan(&total)
	if err != nil {
		r.logger.Error("Failed to get profile count by role", "role_id", roleID, "error", err)
		return nil, 0, fmt.Errorf("failed to get profile count by role: %w", err)
	}

	// Fetch profiles with the specified role
	rows, err := r.pool.Query(ctx, baseQuery, roleID, pageSize, offset)
	if err != nil {
		r.logger.Error("Failed to get profiles by role", "role_id", roleID, "error", err)
		return nil, 0, fmt.Errorf("failed to get profiles by role: %w", err)
	}
	defer rows.Close()

	// Process the query results
	profiles := []*platform_profile.PlatformProfile{}
	for rows.Next() {
		profile := &platform_profile.PlatformProfile{}
		err := rows.Scan(
			&profile.ID,
			&profile.Username,
			&profile.Email,
			&profile.PasswordHash,
			&profile.Status,
			&profile.VerifiedAt,
			&profile.LastLoginAt,
			&profile.FailedLoginAttempts,
			&profile.CreatedAt,
			&profile.UpdatedByUserAt,
			&profile.UpdatedBySystemAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan profile by role", "role_id", roleID, "error", err)
			return nil, 0, fmt.Errorf("failed to scan profile by role: %w", err)
		}
		profiles = append(profiles, profile)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over profile rows by role", "role_id", roleID, "error", err)
		return nil, 0, fmt.Errorf("error iterating over profile rows by role: %w", err)
	}

	r.logger.Info("Profiles fetched successfully by role", "role_id", roleID, "total_profiles", len(profiles))
	return profiles, total, nil
}

// GetProfilesByPreferences retrieves profiles that match specific preferences with pagination
func (r *PostgresProfileRepository) GetProfilesByPreferences(ctx context.Context, preferences map[string]interface{}, page, pageSize int) ([]*platform_profile.PlatformProfile, int, error) {
	r.logger.Debug("Fetching profiles by preferences", "preferences", preferences, "page", page, "page_size", pageSize)

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Base query to get profiles matching preferences
	baseQuery := `
	SELECT 
		p.id, p.username, p.email, p.password_hash, p.first_name, p.last_name,
		p.display_name, p.status, p.created_at, p.updated_by_user_at, 
		p.last_login_at, p.failed_login_attempts, p.verified_at
	FROM profile_schema.platform_profiles p
	JOIN profile_schema.profile_preferences pref ON p.id = pref.profile_id`

	// Count query to get the total number of matching profiles
	countQuery := `
	SELECT COUNT(*)
	FROM profile_schema.platform_profiles p
	JOIN profile_schema.profile_preferences pref ON p.id = pref.profile_id`

	// Prepare WHERE clause for preferences
	whereClause := ""
	args := []interface{}{}
	paramIndex := 1

	if len(preferences) > 0 {
		whereClause = " WHERE "
		conditions := []string{}

		for key, value := range preferences {
			switch key {
			case "notifications_email":
				conditions = append(conditions, fmt.Sprintf("pref.notifications_email = $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			case "notifications_push":
				conditions = append(conditions, fmt.Sprintf("pref.notifications_push = $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			case "notifications_sms":
				conditions = append(conditions, fmt.Sprintf("pref.notifications_sms = $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			case "language":
				conditions = append(conditions, fmt.Sprintf("pref.language ILIKE $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			case "dark_mode_on":
				conditions = append(conditions, fmt.Sprintf("pref.dark_mode_on = $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			case "status":
				conditions = append(conditions, fmt.Sprintf("p.status = $%d", paramIndex))
				args = append(args, value)
				paramIndex++
			case "verified":
				notStr := ""
				if value.(bool) {
					notStr = "NOT"
				}
				conditions = append(conditions, fmt.Sprintf("p.verified_at IS %s NULL", notStr))
			case "search":
				search := fmt.Sprintf("%%%s%%", value)
				conditions = append(conditions, fmt.Sprintf("(p.username ILIKE $%d OR p.email ILIKE $%d OR p.first_name ILIKE $%d OR p.last_name ILIKE $%d)",
					paramIndex, paramIndex+1, paramIndex+2, paramIndex+3))
				args = append(args, search, search, search, search)
				paramIndex += 4
			}
		}

		whereClause += strings.Join(conditions, " AND ")
	}

	// Complete data and count queries
	dataQuery := baseQuery + whereClause + fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	countQueryFinal := countQuery + whereClause

	// Add pagination parameters
	args = append(args, pageSize, offset)

	// Get total count of matching profiles
	var total int
	err := r.pool.QueryRow(ctx, countQueryFinal, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		r.logger.Error("Failed to get profile count by preferences", "preferences", preferences, "error", err)
		return nil, 0, fmt.Errorf("failed to get profile count by preferences: %w", err)
	}

	// Execute query to fetch profiles
	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		r.logger.Error("Failed to get profiles by preferences", "preferences", preferences, "error", err)
		return nil, 0, fmt.Errorf("failed to get profiles by preferences: %w", err)
	}
	defer rows.Close()

	// Process query results
	profiles := []*platform_profile.PlatformProfile{}
	for rows.Next() {
		profile := &platform_profile.PlatformProfile{}
		err := rows.Scan(
			&profile.ID,
			&profile.Username,
			&profile.Email,
			&profile.PasswordHash,
			&profile.Status,
			&profile.VerifiedAt,
			&profile.LastLoginAt,
			&profile.FailedLoginAttempts,
			&profile.CreatedAt,
			&profile.UpdatedByUserAt,
			&profile.UpdatedBySystemAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan profile by preferences", "preferences", preferences, "error", err)
			return nil, 0, fmt.Errorf("failed to scan profile by preferences: %w", err)
		}
		profiles = append(profiles, profile)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over profile rows by preferences", "preferences", preferences, "error", err)
		return nil, 0, fmt.Errorf("error iterating over profile rows by preferences: %w", err)
	}

	r.logger.Info("Profiles fetched successfully by preferences", "total_profiles", len(profiles))
	return profiles, total, nil
}

// DeleteProfiles deletes multiple profiles either as a hard delete or a soft delete
func (r *PostgresProfileRepository) DeleteProfiles(ctx context.Context, profileIDs []uuid.UUID, hardDelete bool) error {
	r.logger.Debug("Starting profile deletion process", "hard_delete", hardDelete, "profile_ids", profileIDs)

	if len(profileIDs) == 0 {
		r.logger.Warn("No profile IDs provided for deletion")
		return nil
	}

	// Prepare query placeholders for IN clause
	placeholders := make([]string, len(profileIDs))
	args := make([]interface{}, len(profileIDs))
	for i, id := range profileIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	idList := strings.Join(placeholders, ", ")

	var query string
	var err error

	if hardDelete {
		// Hard delete: Permanently remove profiles and related records
		query = fmt.Sprintf(`
		DELETE FROM profile_schema.platform_profiles
		WHERE id IN (%s)`, idList)

		_, err = r.pool.Exec(ctx, query, args...)
		if err != nil {
			r.logger.Error("Failed to hard delete profiles", "error", err, "profile_ids", profileIDs)
			return fmt.Errorf("failed to hard delete profiles: %w", err)
		}
		r.logger.Info("Profiles permanently deleted successfully", "profile_ids", profileIDs)
	} else {
		// Soft delete: Mark profiles as deleted and store them in deleted_profiles
		tx, err := r.pool.Begin(ctx)
		if err != nil {
			r.logger.Error("Failed to start transaction for soft delete", "error", err)
			return fmt.Errorf("failed to start transaction for soft delete: %w", err)
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback(ctx)
			}
		}()

		// Move profiles to deleted_profiles table before marking them as deleted
		moveQuery := fmt.Sprintf(`
		INSERT INTO profile_schema.deleted_profiles (
			id, username, email, password_hash, first_name, last_name, 
			display_name, status, created_at, updated_by_user_at, 
			last_login_at, failed_login_attempts, verified_at, 
			deleted_at
		)
		SELECT 
			id, username, email, password_hash, first_name, last_name, 
			display_name, status, created_at, updated_by_user_at, 
			last_login_at, failed_login_attempts, verified_at, 
			NOW()
		FROM profile_schema.platform_profiles
		WHERE id IN (%s)`, idList)

		_, err = tx.Exec(ctx, moveQuery, args...)
		if err != nil {
			r.logger.Error("Failed to move profiles to deleted_profiles for soft delete", "error", err, "profile_ids", profileIDs)
			return fmt.Errorf("failed to move profiles to deleted_profiles: %w", err)
		}

		// Commit the transaction
		err = tx.Commit(ctx)
		if err != nil {
			r.logger.Error("Failed to commit transaction for soft delete", "error", err)
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		r.logger.Info("Profiles soft deleted successfully and moved to deleted_profiles", "profile_ids", profileIDs)
	}

	return nil
}

// RecordLogin updates the last login timestamp and resets failed login attempts
func (r *PostgresProfileRepository) RecordLogin(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Recording login for profile", "profile_id", id)

	query := `
	UPDATE profile_schema.platform_profiles
	SET
		last_login_at = $1,
		updated_by_system_at = $1,
		failed_login_attempts = 0
	WHERE id = $2
	`

	now := time.Now()

	// Execute the query to update login information
	commandTag, err := r.pool.Exec(ctx, query, now, id)
	if err != nil {
		r.logger.Error("Failed to record login", "profile_id", id, "error", err)
		return fmt.Errorf("failed to record login: %w", err)
	}

	// Check if any rows were affected (handle profile not found case)
	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile not found or already deleted", "profile_id", id)
		return ErrProfileNotFound
	}

	r.logger.Info("Login recorded successfully", "profile_id", id, "last_login_at", now)
	return nil
}

// IncrementFailedLoginAttempts increments the failed login attempts counter for a profile
func (r *PostgresProfileRepository) IncrementFailedLoginAttempts(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Incrementing failed login attempts", "profile_id", id)

	query := `
	UPDATE profile_schema.platform_profiles
	SET
		failed_login_attempts = failed_login_attempts + 1,
		updated_by_system_at = $1
	WHERE id = $2 `

	now := time.Now()

	// Execute the query to increment failed login attempts
	commandTag, err := r.pool.Exec(ctx, query, now, id)
	if err != nil {
		r.logger.Error("Failed to increment failed login attempts", "profile_id", id, "error", err)
		return fmt.Errorf("failed to increment failed login attempts: %w", err)
	}

	// Check if any rows were affected (profile not found or deleted)
	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile not found or already deleted, failed attempts not incremented", "profile_id", id)
		return ErrProfileNotFound
	}

	r.logger.Info("Failed login attempts incremented successfully", "profile_id", id)
	return nil
}

// ResetFailedLoginAttempts resets the failed login attempts counter for a profile
func (r *PostgresProfileRepository) ResetFailedLoginAttempts(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Resetting failed login attempts", "profile_id", id)

	query := `
	UPDATE profile_schema.platform_profiles
	SET
		failed_login_attempts = 0,
		updated_by_system_at = $1
	WHERE id = $2
	`

	now := time.Now()

	// Execute the query to reset failed login attempts
	commandTag, err := r.pool.Exec(ctx, query, now, id)
	if err != nil {
		r.logger.Error("Failed to reset failed login attempts", "profile_id", id, "error", err)
		return fmt.Errorf("failed to reset failed login attempts: %w", err)
	}

	// Check if any rows were affected (profile not found or deleted)
	if commandTag.RowsAffected() == 0 {
		r.logger.Warn("Profile not found or already deleted, failed attempts not reset", "profile_id", id)
		return ErrProfileNotFound
	}

	r.logger.Info("Failed login attempts reset successfully", "profile_id", id)
	return nil
}
