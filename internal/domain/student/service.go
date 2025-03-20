// internal/domain/student/service.go

package student

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Common errors
var (
	ErrStudentNotFound      = errors.New("student not found")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrEnrollmentIDExists   = errors.New("enrollment ID already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrInvalidEmail         = errors.New("invalid email format")
	ErrInvalidEnrollmentID  = errors.New("invalid enrollment ID format")
	ErrInvalidPasswordStrength = errors.New("password does not meet strength requirements")
)

// Service provides student-related operations
type Service struct {
	repo Repository
	// Add other necessary dependencies like event publisher, logger, etc.
	// eventPublisher eventbus.Publisher
	// logger         logger.Logger
}

// NewService creates a new instance of the student service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Register registers a new student in the system
func (s *Service) Register(ctx context.Context, student *Student, password string) (*Student, error) {
	// Validate email format
	if !isValidEmail(student.Email) {
		return nil, ErrInvalidEmail
	}

	// Validate enrollment ID format
	if !isValidEnrollmentID(student.EnrollmentID) {
		return nil, ErrInvalidEnrollmentID
	}

	// Check if email already exists
	existingStudent, err := s.repo.GetByEmail(ctx, student.Email)
	if err == nil && existingStudent != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Check if enrollment ID already exists
	existingStudent, err = s.repo.GetByEnrollmentID(ctx, student.EnrollmentID)
	if err == nil && existingStudent != nil {
		return nil, ErrEnrollmentIDExists
	}

	// Validate password strength
	if !isPasswordStrong(password) {
		return nil, ErrInvalidPasswordStrength
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default values
	if student.ID == uuid.Nil {
		student.ID = uuid.New()
	}
	student.PasswordHash = hashedPassword
	student.Status = StatusActive
	now := time.Now()
	student.CreatedAt = now
	student.UpdatedAt = now

	// Save student to repository
	if err := s.repo.Create(ctx, student); err != nil {
		return nil, fmt.Errorf("failed to create student: %w", err)
	}

	// Publish student created event
	// s.eventPublisher.Publish("student.created", student)

	return student, nil
}

// Authenticate authenticates a student with email and password
func (s *Service) Authenticate(ctx context.Context, email, password string) (*Student, error) {
	student, err := s.repo.GetByEmail(ctx, email)
	if err != nil || student == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if student is active
	if !student.IsActive() {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := verifyPassword(student.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login time
	now := time.Now()
	student.LastLoginAt = &now
	student.UpdatedAt = now
	if err := s.repo.Update(ctx, student); err != nil {
		return nil, fmt.Errorf("failed to update last login time: %w", err)
	}

	return student, nil
}

// GetByID retrieves a student by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Student, error) {
	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrStudentNotFound
	}
	return student, nil
}

// GetByEnrollmentID retrieves a student by enrollment ID
func (s *Service) GetByEnrollmentID(ctx context.Context, enrollmentID string) (*Student, error) {
	student, err := s.repo.GetByEnrollmentID(ctx, enrollmentID)
	if err != nil {
		return nil, ErrStudentNotFound
	}
	return student, nil
}

// Update updates a student's profile
func (s *Service) Update(ctx context.Context, student *Student) (*Student, error) {
	existingStudent, err := s.repo.GetByID(ctx, student.ID)
	if err != nil {
		return nil, ErrStudentNotFound
	}

	// Prevent updating critical fields
	student.EnrollmentID = existingStudent.EnrollmentID
	student.Email = existingStudent.Email
	student.PasswordHash = existingStudent.PasswordHash
	student.CreatedAt = existingStudent.CreatedAt
	student.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, student); err != nil {
		return nil, fmt.Errorf("failed to update student: %w", err)
	}

	return student, nil
}

// UpdateEmail updates a student's email
func (s *Service) UpdateEmail(ctx context.Context, studentID uuid.UUID, newEmail string) error {
	if !isValidEmail(newEmail) {
		return ErrInvalidEmail
	}

	// Check if email already exists
	existingStudent, err := s.repo.GetByEmail(ctx, newEmail)
	if err == nil && existingStudent != nil && existingStudent.ID != studentID {
		return ErrEmailAlreadyExists
	}

	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return ErrStudentNotFound
	}

	student.Email = newEmail
	student.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, student); err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	// Publish email updated event
	// s.eventPublisher.Publish("student.email_updated", student)

	return nil
}

// ChangePassword changes a student's password
func (s *Service) ChangePassword(ctx context.Context, studentID uuid.UUID, currentPassword, newPassword string) error {
	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return ErrStudentNotFound
	}

	// Verify current password
	if err := verifyPassword(student.PasswordHash, currentPassword); err != nil {
		return ErrInvalidCredentials
	}

	// Validate password strength
	if !isPasswordStrong(newPassword) {
		return ErrInvalidPasswordStrength
	}

	// Hash new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	student.PasswordHash = hashedPassword
	student.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, student); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all sessions for security
	// TODO: Implement session revocation logic

	return nil
}

// UpdateStatus updates a student's status
func (s *Service) UpdateStatus(ctx context.Context, studentID uuid.UUID, status StudentStatus) error {
	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return ErrStudentNotFound
	}

	// Validate status
	if !isValidStatus(status) {
		return ErrInvalidStatus
	}

	student.Status = status
	student.UpdatedAt = time.Now()

	if status == StatusDeactivated {
		now := time.Now()
		student.DeactivatedAt = &now
	}

	if err := s.repo.Update(ctx, student); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Publish status updated event
	// s.eventPublisher.Publish("student.status_updated", student)

	return nil
}

// ListStudents retrieves a list of students with filtering and pagination
func (s *Service) ListStudents(ctx context.Context, filter StudentFilter, pagination Pagination) ([]*Student, int64, error) {
	return s.repo.List(ctx, filter, pagination)
}

// EnrollCourse enrolls a student in a course
func (s *Service) EnrollCourse(ctx context.Context, studentID uuid.UUID, course Course) error {
	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return ErrStudentNotFound
	}

	// Check if student is already enrolled in the course
	for _, existingCourse := range student.Courses {
		if existingCourse.CourseCode == course.CourseCode {
			return errors.New("student already enrolled in this course")
		}
	}

	// Set default values
	course.ID = uuid.New()
	course.EnrolledAt = time.Now()

	return s.repo.AddCourse(ctx, studentID, course)
}

// UnenrollCourse removes a student from a course
func (s *Service) UnenrollCourse(ctx context.Context, studentID, courseID uuid.UUID) error {
	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return ErrStudentNotFound
	}

	// Check if student is enrolled in the course
	found := false
	for _, course := range student.Courses {
		if course.ID == courseID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("student not enrolled in this course")
	}

	return s.repo.RemoveCourse(ctx, studentID, courseID)
}

// UpdateCourseGrade updates a student's grade for a specific course
func (s *Service) UpdateCourseGrade(ctx context.Context, studentID, courseID uuid.UUID, grade string) error {
	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return ErrStudentNotFound
	}

	// Check if student is enrolled in the course
	found := false
	for _, course := range student.Courses {
		if course.ID == courseID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("student not enrolled in this course")
	}

	// Validate grade format
	if !isValidGrade(grade) {
		return errors.New("invalid grade format")
	}

	return s.repo.UpdateCourseGrade(ctx, studentID, courseID, grade)
}

// GetAttendance retrieves attendance records for a student
func (s *Service) GetAttendance(ctx context.Context, studentID uuid.UUID, courseID *uuid.UUID, semester *int) ([]StudentAttendance, error) {
	_, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return nil, ErrStudentNotFound
	}

	return s.repo.GetAttendance(ctx, studentID, courseID, semester)
}

// RecordAttendance records attendance for a student
func (s *Service) RecordAttendance(ctx context.Context, attendance *StudentAttendance) error {
	_, err := s.repo.GetByID(ctx, attendance.StudentID)
	if err != nil {
		return ErrStudentNotFound
	}

	// Validate attendance record
	if attendance.ID == uuid.Nil {
		attendance.ID = uuid.New()
	}
	attendance.CreatedAt = time.Now()
	attendance.UpdatedAt = time.Now()

	// Calculate statistics
	calculateAttendanceStatistics(attendance)

	return s.repo.RecordAttendance(ctx, attendance)
}

// GetPreferences retrieves preferences for a student
func (s *Service) GetPreferences(ctx context.Context, studentID uuid.UUID) (*StudentPreferences, error) {
	_, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return nil, ErrStudentNotFound
	}

	return s.repo.GetPreferences(ctx, studentID)
}

// UpdatePreferences updates preferences for a student
func (s *Service) UpdatePreferences(ctx context.Context, preferences *StudentPreferences) error {
	_, err := s.repo.GetByID(ctx, preferences.StudentID)
	if err != nil {
		return ErrStudentNotFound
	}

	preferences.UpdatedAt = time.Now()
	return s.repo.UpdatePreferences(ctx, preferences)
}

// AddAchievement adds an achievement to a student
func (s *Service) AddAchievement(ctx context.Context, studentID uuid.UUID, achievement Achievement) error {
	_, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return ErrStudentNotFound
	}

	// Set default values
	if achievement.ID == uuid.Nil {
		achievement.ID = uuid.New()
	}
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	return s.repo.AddAchievement(ctx, studentID, achievement)
}

// GetAchievements retrieves achievements for a student
func (s *Service) GetAchievements(ctx context.Context, studentID uuid.UUID) ([]Achievement, error) {
	_, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return nil, ErrStudentNotFound
	}

	return s.repo.GetAchievements(ctx, studentID)
}

// CreateSession creates a new session for a student
func (s *Service) CreateSession(ctx context.Context, studentID uuid.UUID, userAgent, ipAddress string) (*StudentSession, error) {
	student, err := s.repo.GetByID(ctx, studentID)
	if err != nil {
		return nil, ErrStudentNotFound
	}

	if !student.IsActive() {
		return nil, errors.New("inactive student cannot create session")
	}

	// Generate tokens
	token := generateToken(student.ID)
	refreshToken := generateRefreshToken()

	// Create session
	session := &StudentSession{
		ID:           uuid.New(),
		StudentID:    studentID,
		Token:        token,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(24 * time.Hour), // Token expires in 24 hours
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// RevokeSession revokes a session
func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.repo.RevokeSession(ctx, sessionID)
}

// Helper functions

// isValidEmail validates email format
func isValidEmail(email string) bool {
	// Simple validation, can be expanded with regex
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}


// isValidEnrollmentID validates enrollment ID format
func isValidEnrollmentID(enrollmentID string) bool {
	// This can be expanded with specific validation logic based on your university's enrollment ID format
	return len(enrollmentID) >= 5 && len(enrollmentID) <= 15
}

// isPasswordStrong validates password strength
func isPasswordStrong(password string) bool {
	// Password must be at least 8 characters, contain uppercase, lowercase, number, and special character
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:'\",.<>/?", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// isValidStatus validates student status
func isValidStatus(status StudentStatus) bool {
	validStatuses := []StudentStatus{
		StatusActive,
		StatusOnLeave,
		StatusGraduated,
		StatusSuspended,
		StatusDeactivated,
		StatusProvisional,
	}

	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

// isValidGrade validates grade format
func isValidGrade(grade string) bool {
	validGrades := []string{"A", "A-", "B+", "B", "B-", "C+", "C", "C-", "D+", "D", "F", "I", "W"}
	for _, g := range validGrades {
		if grade == g {
			return true
		}
	}
	return false
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword verifies a password against its hash
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// generateToken generates a JWT token for a student
func generateToken(studentID uuid.UUID) string {
	// This would typically use a JWT library to generate a proper token
	// Placeholder implementation
	return fmt.Sprintf("token_%s_%d", studentID.String(), time.Now().Unix())
}

// generateRefreshToken generates a refresh token
func generateRefreshToken() string {
	// This would typically generate a cryptographically secure random token
	// Placeholder implementation
	return fmt.Sprintf("refresh_%d", time.Now().UnixNano())
}

// calculateAttendanceStatistics calculates statistics for attendance records
func calculateAttendanceStatistics(attendance *StudentAttendance) {
	var present, absent, leave int
	
	for _, record := range attendance.Records {
		switch record.Status {
		case AttendancePresent:
			present++
		case AttendanceAbsent:
			absent++
		case AttendanceLeave:
			leave++
		}
	}
	
	totalClasses := present + absent + leave
	var attendanceRate float64
	if totalClasses > 0 {
		attendanceRate = float64(present) / float64(totalClasses) * 100
	}
	
	attendance.Statistics = AttendanceStatistics{
		TotalClasses:   totalClasses,
		Present:        present,
		Absent:         absent,
		Leave:          leave,
		AttendanceRate: attendanceRate,
	}
}