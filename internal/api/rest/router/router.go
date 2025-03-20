package router

import (
	"server/internal/config"
	"server/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(r *gin.Engine, db *sqlx.DB, log *logger.Logger, cfg *config.Config) {
	// API versioning
	v1 := r.Group("/api/v1")

	// Register all route groups
	RegisterStudentRoutes(v1, db, log, cfg)
	RegisterQuizRoutes(v1, db, log, cfg)
	
	// Add more route groups as needed
}