// The basic profile details for the students.
// Uses document IDs from the StudentDocumentTable table for Photograph and Resume.
// Depends of the StudentDocumentTable table.
package student

import (
	"time"
)

type StudentProfileDetailsTable struct {
	// ID = Primary Key (unique identifier for each student).
	// To be used internally only. Not to be sent in responses.
	ID uint32 `gorm:"primaryKey;autoIncrement;unique" json:"-" bson:"-"`

	// UserRole = Privelage check parameter, defaults to 'STU' (Student)
	// To be used internally only. Not to be sent in responses.
	UserRole string `gorm:"type:varchar(3);size:3;not null;default:'STU';check:user_role IN ('STU','VOL','COR','ADM');" bson:"-" json:"-"`

	// Full name of the student
	Name string `gorm:"type:varchar(100);not null;size:100" json:"name" bson:"name" binding:"required"`

	// Gender: 'M' (Male), 'F' (Female), 'O' (Other)
	Gender string `gorm:"type:varchar(1);not null;size:1;check:gender IN ('M','F','O')" json:"gender" bson:"gender" binding:"required,oneof=M F O"`

	// Category: 'GEN', 'EWS', 'OBC', 'SC', or 'ST'
	// Extend if needed.
	Category string `gorm:"type:varchar(3);not null;size:3;check:category IN ('GEN','EWS','OBC','SC','ST')" json:"category" bson:"category" binding:"required,oneof=GEN EWS OBC SC ST"`

	// Foreign Key: References `DocumentID` in the `Documents` table for Photograph
	// To be used internally only. Not to be sent in responses.
	PhotographID uint32 `gorm:"not null;unique" json:"-" bson:"-" binding:"required"`

	// Foreign Key: References `DocumentID` in the `Documents` table for Resume
	// To be used internally only. Not to be sent in responses.
	ResumeID uint32 `gorm:"not null;unique" json:"-" bson:"-" binding:"required"`

	CreatedAt time.Time `gorm:"autoCreateTime"` // Automatically set timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Automatically update timestamp

	// Foreign key relations
	// Relationship to Documents table for Photograph
	Photograph StudentDocumentTable `gorm:"foreignKey:photograph_id;references:document_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Relationship to Documents table for Resume
	Resume StudentDocumentTable `gorm:"foreignKey:resume_id;references:document_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (StudentProfileDetailsTable) TableName() string {
	return "student_schema.student_profile_details_table"
}
