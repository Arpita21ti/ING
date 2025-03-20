// This is an independent table.
// Referenced in the EnrollmentDetailsTable.
package student

import (
	"time"
)

type StudentScholarshipDetailsTable struct {
	// ID = Primary Key (unique identifier for the scholarship record)
	ID uint32 `gorm:"primaryKey;autoIncrement" json:"-" bson:"-"`

	// Name of the scholarship received
	ScholarshipName string `gorm:"type:varchar(100);not null" json:"scholarshipName" bson:"scholarshipName" binding:"required"`

	// Organization or authority providing the scholarship
	ProvidedBy string `gorm:"type:varchar(100);not null" json:"providedBy" bson:"providedBy" binding:"required"`

	// Amount received for the scholarship
	AmountReceived int `gorm:"type:int;not null;" json:"amountReceived" bson:"amountReceived" binding:"required"`

	// Timestamps for record updates
	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Automatically update timestamp
}

// TableName overrides the default table name
func (StudentScholarshipDetailsTable) TableName() string {
	return "student_schema.student_scholarship_details_table"
}
