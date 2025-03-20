package router

import (
	"server/internal/api/rest/handlers"
	"server/internal/config"
	"server/internal/repository"
	"server/internal/service"
	"server/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// RegisterStudentRoutes sets up all student-related routes
func RegisterStudentRoutes(r *gin.RouterGroup, db *sqlx.DB, log *logger.Logger, cfg *config.Config) {
	// Create repositories
	studentRepo := repository.NewStudentRepository(db)
	
	// Create services
	studentService := service.NewStudentService(studentRepo, log)
	
	// Create handlers
	studentHandler := handlers.NewStudentHandler(studentService, log)
	
	// Define routes
	students := r.Group("/students")
	{
		students.GET("", studentHandler.GetAllStudents)
		students.GET("/:id", studentHandler.GetStudentByID)
		students.POST("", studentHandler.CreateStudent)
		students.PUT("/:id", studentHandler.UpdateStudent)
		students.DELETE("/:id", studentHandler.DeleteStudent)
		
		// Add more student-related routes as needed
	}
}

// package router

// import (
// 	"server/internal/api/rest/handlers/student"
// 	"server/internal/config"
// 	"server/internal/repository"
// 	"server/internal/service"
// 	"server/internal/middleware"
// 	"server/pkg/logger"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jmoiron/sqlx"
// )

// // RegisterStudentRoutes sets up all student-related routes
// func RegisterStudentRoutes(r *gin.RouterGroup, db *sqlx.DB, log logger.Logger, cfg *config.Config) {
// 	// Create repositories
// 	studentRepo := repository.NewStudentRepository(db)
// 	profileRepo := repository.NewProfileRepository(db)
	
// 	// Create services
// 	authService := service.NewAuthService(studentRepo, log, cfg)
// 	profileService := service.NewProfileService(profileRepo, studentRepo, log, cfg)
	
// 	// Create handlers
// 	authHandler := student.NewAuthHandler(authService, log)
// 	profileHandler := student.NewProfileHandler(profileService, log)
	
// 	// Auth routes (no authentication required)
// 	auth := r.Group("/auth")
// 	{
// 		auth.POST("/register", authHandler.Register)
// 		auth.POST("/login", authHandler.Login)
// 		auth.GET("/verify-email", authHandler.VerifyEmail)
// 		auth.POST("/forgot-password", authHandler.ForgotPassword)
// 		auth.POST("/reset-password", authHandler.ResetPassword)
// 	}
	
// 	// Student profile routes (authentication required)
// 	students := r.Group("/students")
// 	students.Use(middleware.Authenticate(authService))
// 	{
// 		students.GET("/profile", profileHandler.GetProfile)
// 		students.PUT("/profile", profileHandler.UpdateProfile)
// 		students.POST("/change-password", profileHandler.ChangePassword)
// 		students.POST("/profile-picture", profileHandler.UploadProfilePicture)
// 		students.GET("/academic-records", profileHandler.GetAcademicRecords)
// 		students.GET("/quiz-history", profileHandler.GetQuizHistory)
// 		students.POST("/logout", authHandler.Logout)
// 	}
// }