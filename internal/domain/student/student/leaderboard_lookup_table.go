package student

type StudentLeaderboardLookupTable struct {
	EnrollmentNo        string `gorm:"type:varchar(12);size:12;not null" json:"enrollmentNo"` // FK to student
	LeaderboardRecordID uint32 `gorm:"primaryKey;not null" json:"leaderboardRecordID"`        // FK to leaderboard record

	// Foreign key relationships
	LeaderboardRecord StudentLeaderboardRecordTable `gorm:"foreignKey:LeaderboardRecordID;references:LeaderboardRecordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (StudentLeaderboardLookupTable) TableName() string {
	return "student_schema.student_leaderboard_lookup_table"
}
