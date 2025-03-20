// This file contains the model for the
// student_academic_details_table table in the student_schema of the PostgreSQL database.
// Contains all the details related to the academic performance of a student.
// Holds reference to the documents table for the marksheet documents.
package student

import (
	"time"
)

type StudentAcademicDetailsTable struct {
	// ID = Primary Key (unique identifier for each student).
	// To be used internally only. Not to be sent in responses.
	ID uint32 `gorm:"primaryKey;autoIncrement" json:"-" bson:"-"`

	// Branch = Short name for the student's branch, e.g., "CSE", "ECE".
	Branch string `gorm:"type:varchar(7);not null" json:"branch" binding:"required" bson:"branch"`

	// Year of Enrollment = Year the student enrolled, restricted to the range [1990-2100].
	YearOfEnrollment int `gorm:"check:year_of_enrollment BETWEEN 1990 AND 2100;not null" json:"yearOfEnrollment" binding:"required" bson:"yearOfEnrollment"`

	// CGPA = Cumulative Grade Point Average of the student.
	CGPA float32 `gorm:"type:real" json:"cgpa" binding:"gte=0,lte=10" bson:"cgpa"`

	// PreviousSemSGPA = SGPA of the last semester completed.
	PreviousSemSGPA float32 `gorm:"type:real" json:"previousSemSgpa" binding:"gte=0,lte=10" bson:"previousSemSgpa"`

	// School for Class Ten
	SchoolForClassTen string `gorm:"type:varchar(255)" json:"schoolForClassTen" binding:"required" bson:"schoolForClassTen"`

	// Class Ten Percentage
	ClassTenPercentage float32 `gorm:"type:real" json:"classTenPercentage" binding:"gte=0,lte=100" bson:"classTenPercentage"`

	// Class Ten Marksheet Document ID = Foreign key to Documents table.
	ClassTenMarksheetID uint32 `gorm:"not null;unique;" json:"-" bson:"-"`

	// School for Class Twelve
	SchoolForClassTwelve string `gorm:"type:varchar(255)" json:"schoolForClassTwelve" binding:"required" bson:"schoolForClassTwelve"`

	// Class Twelve Percentage
	ClassTwelvePercentage float32 `gorm:"type:real" json:"classTwelvePercentage" binding:"gte=0,lte=100" bson:"classTwelvePercentage"`

	// Class Twelve Marksheet Document ID = Foreign key to Documents table.
	ClassTwelveMarksheetID uint32 `gorm:"not null;unique;" json:"-" bson:"-"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Automatically update timestamp

	// Relationships
	// Class Ten Marksheet Document Relationship
	ClassTenMarksheet StudentDocumentTable `gorm:"foreignKey:class_ten_marksheet_id;references:document_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Class Twelve Marksheet Document Relationship
	ClassTwelveMarksheet StudentDocumentTable `gorm:"foreignKey:class_twelve_marksheet_id;references:document_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName overrides the default table name
func (StudentAcademicDetailsTable) TableName() string {
	return "student_schema.student_academic_details_table"
}
