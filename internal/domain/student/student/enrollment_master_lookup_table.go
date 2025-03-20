// Most Dependent Table
// References to all other tables and used to fetch all details of a student in one go.
// Receiving all the Traffic for lookups.
// The certification related information is stored separately in the Certifications lookup table.
// The leaderboard and practice session records must be looked in respective separate dedicated
// lookup tables.
package student

type EnrollmentMasterLookupTable struct {
	EnrollmentNo         string `gorm:"type:varchar(12);size:12;primaryKey;check:char_length(enrollment_no) = 12" json:"enrollmentNo"` // Enrollment number (Primary Key)
	LogInDetailsID       uint32 `gorm:"not null;unique" json:"-" bson:"-"`                                                             // Reference to login details table id for student (uint32)
	AcademicDetailsID    uint32 `gorm:"not null;unique" json:"-" bson:"-"`                                                             // Reference to academic details table id for student (uint32)
	FamilyDetailsID      uint32 `gorm:"not null;unique" json:"-" bson:"-"`                                                             // Reference to family details table id for student (uint32)
	ProfileDetailsID     uint32 `gorm:"not null;unique" json:"-" bson:"-"`                                                             // Reference to profile details table id for student (uint32)
	ScholarshipDetailsID uint32 `gorm:"not null;unique" json:"-" bson:"-"`                                                             // Reference to scholarship details table id for student (uint32)

	// Foreign keys to maintain referential integrity
	StudentLogInDetailsTable StudentLogInDetailsTable       `gorm:"foreignKey:LogInDetailsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-" bson:"-"`
	AcademicDetails          StudentAcademicDetailsTable    `gorm:"foreignKey:AcademicDetailsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-" bson:"-"`
	FamilyDetails            StudentFamilyDetailsTable      `gorm:"foreignKey:FamilyDetailsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-" bson:"-"`
	ProfileDetails           StudentProfileDetailsTable     `gorm:"foreignKey:ProfileDetailsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-" bson:"-"`
	ScholarshipDetails       StudentScholarshipDetailsTable `gorm:"foreignKey:ScholarshipDetailsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-" bson:"-"`
}

func (EnrollmentMasterLookupTable) TableName() string {
	return "public.enrollment_master_lookup_table"
}
