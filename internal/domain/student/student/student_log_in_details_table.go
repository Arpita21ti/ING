// This stores the log in details of the student
// This is an independent table
// Referenced in the EnrollmentMasterLookupTable.
package student

type StudentLogInDetailsTable struct {
	ID       uint32 `gorm:"primaryKey" json:"-" bson:"-"` // The primary Key for this table
	Email    string `gorm:"type:varchar(255);unique;not null" json:"email"  bson:"email" validate:"required,email"`
	Password string `gorm:"type:varchar(255);not null" json:"password"  bson:"password" validate:"required"`
	Phone    string `gorm:"type:varchar(15);not null" json:"phone"  bson:"phone" validate:"required,phone"`
}

func (StudentLogInDetailsTable) TableName() string {
	return "student_schema.student_login_details_table"
}
