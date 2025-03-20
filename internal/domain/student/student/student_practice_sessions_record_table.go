// This is an independent table that stores the practice session records of the students.

package student

import (
	"time"
)

type StudentPracticeSessionRecordTable struct {
	// PracticeSessionID = Unique identifier for each practice session
	PracticeSessionID uint32 `gorm:"primaryKey;not null" json:"sessionId" bson:"sessionId"`

	// DomainID = Category of the questions (e.g., Programming, Mathematics)
	DomainID uint32 `gorm:"not null" json:"domainID" bson:"domainID" binding:"required"`

	// SubDomainID = Sub-category of the questions (e.g., Data Structures, Algebra)
	SubDomainID uint32 `gorm:"not null" json:"subCategoryID" bson:"subCategoryID" binding:"required"`

	// DifficultyLevelID = DifficultyLevelID level of the session (e.g., Easy, Medium, Hard)
	DifficultyLevelID uint32 `gorm:"not null" json:"difficultyID" bson:"difficultyID" binding:"required"`

	// QuestionsAttempted = Number of questions attempted during the session
	QuestionsAttempted int `gorm:"not null" json:"questionsAttempted" bson:"questionsAttempted" binding:"required"`

	// QuestionsCorrect = Number of questions answered correctly
	QuestionsCorrect int `gorm:"not null" json:"questionsCorrect" bson:"questionsCorrect" binding:"required"`

	// ScoreEarned = Total score earned in the session
	ScoreEarned float64 `gorm:"not null" json:"scoreEarned" bson:"scoreEarned" binding:"required"`

	// StartTime = The time when the session started
	StartTime time.Time `gorm:"type:timestamp with time zone;not null" json:"startTime" bson:"startTime" binding:"required"`

	// EndTime = The time when the session ended
	EndTime time.Time `gorm:"type:timestamp with time zone;not null" json:"endTime" bson:"endTime" binding:"required"`

	// Feedbacks = Optional field for any feedback related to the session
	Feedbacks string `gorm:"type:varchar(255);default:''" json:"feedbacks" bson:"feedbacks"`
}

// TableName returns the name of the table in the database
func (StudentPracticeSessionRecordTable) TableName() string {
	return "student_schema.student_practice_session_records"
}
