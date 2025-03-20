package student

type StudentPracticeSessionLookupTable struct {
	EnrollmentNo      string `gorm:"type:varchar(12);size:12;not null" json:"enrollmentNo"`                                                               // FK to student
	PracticeSessionID uint32 `gorm:"not null;primaryKey" json:"practiceSessionID"`                                                                        // FK to practice session
	Status            string `gorm:"type:varchar(9);size:9;not null;default:'Active';check:status IN ('Submitted', 'Active', 'Force End')" json:"status"` // Status of the practice session

	// Foreign key relationships
	PracticeSessionRecord StudentPracticeSessionRecordTable `gorm:"foreignKey:PracticeSessionID;references:PracticeSessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (StudentPracticeSessionLookupTable) TableName() string {
	return "student_schema.student_practice_session_lookup_table"
}
