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

// RegisterQuizRoutes sets up all quiz-related routes
func RegisterQuizRoutes(r *gin.RouterGroup, db *sqlx.DB, log *logger.Logger, cfg *config.Config) {
	// Create repositories
	quizRepo := repository.NewQuizRepository(db)
	
	// Create services
	quizService := service.NewQuizService(quizRepo, log)
	
	// Create handlers
	quizHandler := handlers.NewQuizHandler(quizService, log)
	
	// Define routes
	quizzes := r.Group("/quizzes")
	{
		quizzes.GET("", quizHandler.GetAllQuizzes)
		quizzes.GET("/:id", quizHandler.GetQuizByID)
		quizzes.POST("", quizHandler.CreateQuiz)
		quizzes.PUT("/:id", quizHandler.UpdateQuiz)
		quizzes.DELETE("/:id", quizHandler.DeleteQuiz)
		
		// Questions routes
		questions := quizzes.Group("/:quizId/questions")
		{
			questions.GET("", quizHandler.GetQuizQuestions)
			questions.POST("", quizHandler.AddQuizQuestion)
			questions.PUT("/:questionId", quizHandler.UpdateQuizQuestion)
			questions.DELETE("/:questionId", quizHandler.DeleteQuizQuestion)
		}
		
		// Add more quiz-related routes as needed
	}
}