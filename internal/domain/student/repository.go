// internal/domain/student/repository.go

package student

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the data access contract for students
type Repository interface {
	// Create creates a new student
	Create(ctx context.Context, student *Student) error

	// GetByID retrieves a student by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Student, error)

	// GetByEnrollmentID retrieves a student by enrollment ID
	GetByEnrollmentID(ctx context.Context, enrollmentID string) (*Student, error)

	// GetByEmail retrieves a student by email
	GetByEmail(ctx context.Context, email string) (*Student, error)

	// Update updates an existing student
	Update(ctx context.Context, student *Student) error

	// Delete deletes a student by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves students with pagination
	List(ctx context.Context, filter StudentFilter, pagination Pagination) ([]*Student, int64, error)

	// AddCourse adds a course to a student
	// AddCourse(ctx context.Context, studentID uuid.UUID, course Course) error

	// RemoveCourse removes a course from a student
	RemoveCourse(ctx context.Context, studentID uuid.UUID, courseID uuid.UUID) error

	// UpdateCourseGrade updates a grade for a specific course
	UpdateCourseGrade(ctx context.Context, studentID, courseID uuid.UUID, grade string) error

	// GetAttendance retrieves attendance for a student
	GetAttendance(ctx context.Context, studentID uuid.UUID, courseID *uuid.UUID, semester *int) ([]StudentAttendance, error)

	// RecordAttendance records attendance for a student
	RecordAttendance(ctx context.Context, attendance *StudentAttendance) error

	// GetPreferences retrieves preferences for a student
	GetPreferences(ctx context.Context, studentID uuid.UUID) (*StudentPreferences, error)

	// UpdatePreferences updates preferences for a student
	UpdatePreferences(ctx context.Context, preferences *StudentPreferences) error

	// GetContacts retrieves contacts for a student
	GetContacts(ctx context.Context, studentID uuid.UUID) ([]StudentContact, error)

	// AddContact adds a contact for a student
	AddContact(ctx context.Context, contact *StudentContact) error

	// UpdateContact updates a contact for a student
	UpdateContact(ctx context.Context, contact *StudentContact) error

	// DeleteContact deletes a contact for a student
	DeleteContact(ctx context.Context, contactID uuid.UUID) error

	// CreateSession creates a new session for a student
	// CreateSession(ctx context.Context, session *StudentSession) error

	// GetActiveSession retrieves an active session by token
	// GetActiveSession(ctx context.Context, token string) (*StudentSession, error)

	// RevokeSession revokes a session
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error

	// AddAchievement adds an achievement to a student
	// AddAchievement(ctx context.Context, studentID uuid.UUID, achievement Achievement) error

	// GetAchievements retrieves achievements for a student
	// GetAchievements(ctx context.Context, studentID uuid.UUID) ([]Achievement, error)
}

// StudentFilter defines the filter options for listing students
type StudentFilter struct {
	Status      *StudentStatus
	Program     *string
	Batch       *string
	Semester    *int
	SearchQuery *string // Will match against name, email, or enrollment ID
}

// Pagination defines the pagination options
type Pagination struct {
	Page     int
	PageSize int
	SortBy   string
	SortDesc bool
}