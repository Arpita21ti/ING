// The certification and academic achievements store for students.
// Uses StudentDocumentTable for storing documents related to certifications.
// The DocumentID is a foreign key referencing the StudentDocumentTable table.
package student

import (
	"time"
)

type StudentCertificationDetailsTable struct {
	// ID = Primary Key (unique identifier for each student).
	// To be used internally only. Not to be sent in responses.
	ID uint32 `gorm:"primaryKey;autoIncrement" json:"-" bson:"-"`

	// CertificationName = Name of the certification earned by the student.
	CertificationName string `gorm:"type:varchar(255);not null" json:"certificationName" bson:"certificationName" binding:"required"`

	// IssuingAuthority = Name of the authority or organization that issued the certification.
	IssuingAuthority string `gorm:"type:varchar(255);not null" json:"issuingAuthority" bson:"issuingAuthority" binding:"required"`

	// IssuingDate = Date when the certification was issued. Defaults to the 30th of the month.
	IssuingDate string `gorm:"not null;" json:"issuingDate" bson:"issuingDate" binding:"required"`

	// DocumentID = Foreign key referencing the Documents table for storing certification-related documents.
	DocumentID uint32 `gorm:"not null;unique;" json:"-" bson:"-"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Automatically update timestamp

}

// TableName overrides the default table name
func (StudentCertificationDetailsTable) TableName() string {
	return "student_schema.student_certification_and_achievements_details_table"
}
