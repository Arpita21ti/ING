// This is an independent table to store leaderboard related information.
package student

import (
	"time"
)

type StudentLeaderboardRecordTable struct {
	// LeaderboardRecordID = Unique identifier for each leaderboard record
	LeaderboardRecordID uint32 `gorm:"primaryKey;autoIncrement" json:"-" bson:"-"`

	// Rank = Position of the student in the leaderboard
	Rank int `gorm:"not null" json:"rank" bson:"rank" binding:"required"`

	// Score = Total score obtained by the student
	Score float64 `gorm:"not null" json:"score" bson:"score" binding:"required"` // Float allows partial marking

	// Domain = Category of the questions (e.g., Programming, Mathematics)
	Domain string `gorm:"type:varchar(100);not null" json:"domain" bson:"domain" binding:"required"`

	// SubDomain = Sub-category of the questions (e.g., Data Structures, Algebra)
	SubDomain string `gorm:"type:varchar(100);not null" json:"subDomain" bson:"subDomain" binding:"required"`

	// TimePeriod = Time duration for which the leaderboard is valid (e.g., weekly, monthly)
	TimePeriod string `gorm:"type:varchar(7);not null" json:"timePeriod" bson:"timePeriod" binding:"required"`

	// LastUpdated = Timestamp of the last update to the leaderboard record
	LastUpdated time.Time `gorm:"type:timestamp with time zone;autoUpdateTime" json:"lastUpdated" bson:"lastUpdated"`
}

// TableName returns the name of the table in the database
func (StudentLeaderboardRecordTable) TableName() string {
	return "student_schema.student_leaderboard_records_table"
}
