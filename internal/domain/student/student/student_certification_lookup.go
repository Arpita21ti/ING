package student

type StudentCertificationLookup struct {
	EnrollmentNo                  string `gorm:"type:varchar(12);size:12;not null" json:"enrollmentNo"`
	StudentCertificationDetailsID uint32 `gorm:"primaryKey;not null" json:"certificationDetailsID"`

	// References to the certification details
	CertificationsAndAchievementsDetails StudentCertificationDetailsTable `gorm:"foreignKey:StudentCertificationDetailsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (StudentCertificationLookup) TableName() string {
	return "student_schema.student_certification_lookup_table"
}
