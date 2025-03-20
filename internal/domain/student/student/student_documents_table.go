// Description: This file includes the model for the student_documents_table table in the
// student_schema database.
// The central store for all the documents related to a student.
// Completely Independent of other tables.
// Referenced in other tables as a foreign key.
package student

import "time"

type StudentDocumentTable struct {
	DocumentID   uint32    `gorm:"primaryKey;not null;autoIncrement" json:"-" bson:"-"`
	StoredIn     string    `gorm:"varchar(255);size:255;not null" json:"-" bson:"-"`
	CreatedAt    time.Time `gorm:"type:timestamp with time zone;default:current_timestamp" json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `gorm:"type:timestamp with time zone;default:current_timestamp" json:"updatedAt" bson:"updatedAt"`
	DocumentType string    `gorm:"varchar(255);size:255" json:"documentType" bson:"documentType"`
	URL          string    `gorm:"varchar(255);size:255" json:"documentURL" bson:"documentURL"`

	// Relationships
	// Foreign key relationship with Documents table
	CertificationDocument StudentCertificationDetailsTable `gorm:"foreignKey:DocumentID;references:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

// TableName overrides the default table name
func (StudentDocumentTable) TableName() string {
	return "student_schema.student_documents_table"
}
