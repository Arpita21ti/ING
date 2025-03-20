// This is an independent table that stores the family details of the student.
// Referenced by the `EnrollmentMasterLookupTable` table.
package student

import (
	"time"
)

type StudentFamilyDetailsTable struct {
	// ID = Primary Key (unique identifier for each student).
	// To be used internally only. Not to be sent in responses.
	ID uint32 `gorm:"primaryKey;autoIncrement" json:"-" bson:"-"`

	// Father's details
	FatherName          string `gorm:"type:varchar(100);not null" json:"fatherName" bson:"fatherName" binding:"required"`
	FatherQualification string `gorm:"type:varchar(100);not null" json:"fatherQualification" bson:"fatherQualification" binding:"required"`
	FatherProfession    string `gorm:"type:varchar(100);not null" json:"fatherProfession" bson:"fatherProfession" binding:"required"`

	// Mother's details
	MotherName          string `gorm:"type:varchar(100);not null" json:"motherName" bson:"motherName" binding:"required"`
	MotherQualification string `gorm:"type:varchar(100);not null" json:"motherQualification" bson:"motherQualification" binding:"required"`
	MotherProfession    string `gorm:"type:varchar(100);not null" json:"motherProfession" bson:"motherProfession" binding:"required"`

	// Number of siblings in the family
	NoOfSiblings int `gorm:"type:int;not null;" json:"noOfSiblings" bson:"noOfSiblings" binding:"required"`

	// Total family income in INR
	TotalFamilyIncome int `gorm:"type:int;not null;" json:"totalFamilyIncome" bson:"totalFamilyIncome" binding:"required"`

	// Timestamps for record updates
	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Automatically update timestamp

}

// TableName overrides the default table name
func (StudentFamilyDetailsTable) TableName() string {
	return "student_schema.student_family_details_table"
}
