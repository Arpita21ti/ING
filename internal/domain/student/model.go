// internal/domain/student/model.go

package student

import (
	"time"

	"github.com/google/uuid"
)

// Student represents the core student entity in the system
type Student struct {
	// ID            uuid.UUID     `json:"id"`
	EnrollmentID  string        `json:"enrollment_id"` // University enrollment ID/number
	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	Email         string        `json:"email"`
	PhoneNumber   string        `json:"phone_number"`
	DateOfBirth   time.Time     `json:"date_of_birth"`
	Gender        string        `json:"gender"`
	PresentAddress       Address       `json:"present_address"`
	PermanentAddress Address `json:"permanent_address"`

	Program       string        `json:"program"` // Degree program (e.g., "B.Tech Computer Science")
	Batch         string        `json:"batch"`   // Admission year/batch (e.g., "2023-2027")
	
	PresentSemester      int           `json:"semester"` // Set to 10 for pass outs
	Section       string        `json:"section,omitempty"` // Can be empty
	Status        StudentStatus `json:"status"`
	ProfileImageUrl  string        `json:"profile_image_url"`
	
	// Achievements  []Achievement `json:"achievements,omitempty"`
	// Courses       []Course      `json:"courses,omitempty"`

	RoleIDs       []uuid.UUID   `json:"role_ids,omitempty"` // Link to roles in role domain
	PasswordHash  string        `json:"-"`                  // Never exposed in JSON
	LastLoginAt   *time.Time    `json:"last_login_at,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeactivatedAt *time.Time    `json:"deactivated_at,omitempty"`
}

// Address represents a student's address
type Address struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
}

// Achievement represents student achievements, certifications, or awards
// type Achievement struct {
// 	ID          uuid.UUID `json:"id"`
// 	Title       string    `json:"title"`
// 	Description string    `json:"description,omitempty"`
// 	IssuedBy    string    `json:"issued_by,omitempty"`
// 	IssuedDate  time.Time `json:"issued_date"`
// 	ExpiryDate  *time.Time `json:"expiry_date,omitempty"`
// 	Certificate string    `json:"certificate,omitempty"` // File path or URL to certificate
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

// // Course represents a course that a student is enrolled in
// type Course struct {
// 	ID           uuid.UUID  `json:"id"`
// 	CourseCode   string     `json:"course_code"`
// 	Name         string     `json:"name"`
// 	Credits      float64    `json:"credits"`
// 	Semester     int        `json:"semester"`
// 	Grade        string     `json:"grade,omitempty"`
// 	AttendanceID *uuid.UUID `json:"attendance_id,omitempty"` // Link to attendance records
// 	EnrolledAt   time.Time  `json:"enrolled_at"`
// }

// StudentStatus represents the current status of a student
type StudentStatus string

// Possible student statuses
const (
	StatusActive       StudentStatus = "active"
	StatusOnLeave      StudentStatus = "on_leave"
	StatusGraduated    StudentStatus = "graduated"
	StatusSuspended    StudentStatus = "suspended"
	StatusDeactivated  StudentStatus = "deactivated"
	StatusProvisional  StudentStatus = "provisional"
)

// StudentAttendance represents attendance records for a student
type StudentAttendance struct {
	ID         uuid.UUID            `json:"id"`
	StudentID  uuid.UUID            `json:"student_id"`
	CourseID   uuid.UUID            `json:"course_id"`
	Semester   int                  `json:"semester"`
	Records    []AttendanceRecord   `json:"records"`
	Statistics AttendanceStatistics `json:"statistics"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

// AttendanceRecord represents a single attendance entry
type AttendanceRecord struct {
	Date     time.Time      `json:"date"`
	Status   AttendanceType `json:"status"`
	Remarks  string         `json:"remarks,omitempty"`
	RecordedBy uuid.UUID    `json:"recorded_by"` // ID of coordinator/faculty who marked attendance
}

// AttendanceType represents the status of attendance for a class
type AttendanceType string

// Possible attendance types
const (
	AttendancePresent AttendanceType = "present"
	AttendanceAbsent  AttendanceType = "absent"
	AttendanceLeave   AttendanceType = "leave"
)

// AttendanceStatistics represents calculated attendance metrics
type AttendanceStatistics struct {
	TotalClasses   int     `json:"total_classes"`
	Present        int     `json:"present"`
	Absent         int     `json:"absent"`
	Leave          int     `json:"leave"`
	AttendanceRate float64 `json:"attendance_rate"` // Percentage
}

// StudentPreferences represents a student's preferences and settings
type StudentPreferences struct {
	ID                 uuid.UUID `json:"id"`
	StudentID          uuid.UUID `json:"student_id"`
	Calendar           Calendar  `json:"calendar"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Calendar represents calendar sync preferences
type Calendar struct {
	Enabled    bool   `json:"enabled"`
	Provider   string `json:"provider,omitempty"` // e.g., "google", "outlook"
	ExternalID string `json:"external_id,omitempty"`
}

// StudentContact represents a student's emergency contact information
type StudentContact struct {
	ID         uuid.UUID `json:"id"`
	StudentID  uuid.UUID `json:"student_id"`
	Name       string    `json:"name"`
	Relation   string    `json:"relation"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email,omitempty"`
	IsEmergency bool      `json:"is_emergency"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}



// NewStudent creates a new student instance with default values
func NewStudent(enrollmentID, firstName, lastName, email, program, batch string, semester int) *Student {
	now := time.Now()
	return &Student{
		// ID:           uuid.New(),
		EnrollmentID: enrollmentID,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Program:      program,
		Batch:        batch,
		PresentSemester:     semester,
		Status:       StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// FullName returns the student's full name
func (s *Student) FullName() string {
	return s.FirstName + " " + s.LastName
}

// IsActive checks if the student is currently active
func (s *Student) IsActive() bool {
	return s.Status == StatusActive && s.DeactivatedAt == nil
}

// HasRole checks if the student has a specific role
func (s *Student) HasRole(roleID uuid.UUID) bool {
	for _, id := range s.RoleIDs {
		if id == roleID {
			return true
		}
	}
	return false
}

// // AddCourse adds a course to the student's enrollment
// func (s *Student) AddCourse(course Course) {
// 	s.Courses = append(s.Courses, course)
// 	s.UpdatedAt = time.Now()
// }

// GetAttendanceRate calculates the overall attendance rate across all courses
func (s *Student) GetAttendanceRate() float64 {
	// This would typically be calculated from attendance records
	// This is a simplified placeholder
	return 0.0
}

// Deactivate sets the student's status to deactivated
func (s *Student) Deactivate() {
	now := time.Now()
	s.Status = StatusDeactivated
	s.DeactivatedAt = &now
	s.UpdatedAt = now
}