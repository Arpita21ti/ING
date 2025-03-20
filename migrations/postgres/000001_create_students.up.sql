CREATE TABLE student_schema.student_leaderboard_records_table (
	LeaderboardRecordID SERIAL PRIMARY KEY,
	Rank INT NOT NULL,
	Score FLOAT NOT NULL,
	Domain VARCHAR(100) NOT NULL,
	SubDomain VARCHAR(100) NOT NULL,
	TimePeriod VARCHAR(7) NOT NULL,
	LastUpdated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE student_schema.student_academic_details_table (
	ID SERIAL PRIMARY KEY,
	Branch VARCHAR(7) NOT NULL,
	YearOfEnrollment INT CHECK (YearOfEnrollment BETWEEN 1990 AND 2100) NOT NULL,
	CGPA REAL,
	PreviousSemSGPA REAL,
	SchoolForClassTen VARCHAR(255) NOT NULL,
	ClassTenPercentage REAL,
	ClassTenMarksheetID INT NOT NULL UNIQUE,
	SchoolForClassTwelve VARCHAR(255) NOT NULL,
	ClassTwelvePercentage REAL,
	ClassTwelveMarksheetID INT NOT NULL UNIQUE,
	UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE student_schema.student_login_details_table (
	ID SERIAL PRIMARY KEY,
	Email VARCHAR(255) UNIQUE NOT NULL,
	Password VARCHAR(255) NOT NULL,
	Phone VARCHAR(15) NOT NULL
);