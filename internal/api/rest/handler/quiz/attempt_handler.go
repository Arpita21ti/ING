package quiz

import (
	"net/http"
	"strconv"
	"time"

	"server/internal/model"
	"server/internal/service"
	"server/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AttemptHandler handles HTTP requests related to quiz attempts
type AttemptHandler struct {
	attemptService service.AttemptService
	logger         logger.Logger
}

// NewAttemptHandler creates a new AttemptHandler instance
func NewAttemptHandler(attemptService service.AttemptService, logger logger.Logger) *AttemptHandler {
	return &AttemptHandler{
		attemptService: attemptService,
		logger:         logger,
	}
}

// StartQuizAttempt creates a new quiz attempt
func (h *AttemptHandler) StartQuizAttempt(c *gin.Context) {
	quizID, err := strconv.ParseInt(c.Param("quizId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	attempt := &model.QuizAttempt{
		QuizID:     quizID,
		StudentID:  studentID.(int64),
		StartTime:  time.Now(),
		Status:     "in_progress",
	}

	attemptID, err := h.attemptService.StartQuizAttempt(c.Request.Context(), attempt)
	if err != nil {
		h.logger.Error("Failed to start quiz attempt", "quizID", quizID, "studentID", studentID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start quiz attempt"})
		return
	}

	attempt.ID = attemptID
	c.JSON(http.StatusCreated, attempt)
}

// SubmitQuizAttempt submits answers for a quiz attempt
func (h *AttemptHandler) SubmitQuizAttempt(c *gin.Context) {
	attemptID, err := strconv.ParseInt(c.Param("attemptId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attempt ID"})
		return
	}

	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var answers []model.QuizAnswer
	if err := c.ShouldBindJSON(&answers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the attempt belongs to the student
	attempt, err := h.attemptService.GetAttemptByID(c.Request.Context(), attemptID)
	if err != nil {
		h.logger.Error("Failed to retrieve attempt", "attemptID", attemptID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit quiz attempt"})
		return
	}

	if attempt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attempt not found"})
		return
	}

	if attempt.StudentID != studentID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to submit this attempt"})
		return
	}

	result, err := h.attemptService.SubmitQuizAttempt(c.Request.Context(), attemptID, answers)
	if err != nil {
		h.logger.Error("Failed to submit quiz attempt", "attemptID", attemptID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit quiz attempt"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetStudentAttempts retrieves all attempts for a student
func (h *AttemptHandler) GetStudentAttempts(c *gin.Context) {
	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	attempts, err := h.attemptService.GetStudentAttempts(c.Request.Context(), studentID.(int64))
	if err != nil {
		h.logger.Error("Failed to get student attempts", "studentID", studentID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attempts"})
		return
	}

	c.JSON(http.StatusOK, attempts)
}

// GetAttemptDetails retrieves details of a specific attempt
func (h *AttemptHandler) GetAttemptDetails(c *gin.Context) {
	attemptID, err := strconv.ParseInt(c.Param("attemptId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attempt ID"})
		return
	}

	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	details, err := h.attemptService.GetAttemptDetails(c.Request.Context(), attemptID)
	if err != nil {
		h.logger.Error("Failed to get attempt details", "attemptID", attemptID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attempt details"})
		return
	}

	if details == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attempt not found"})
		return
	}

	if details.StudentID != studentID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view this attempt"})
		return
	}

	c.JSON(http.StatusOK, details)
}