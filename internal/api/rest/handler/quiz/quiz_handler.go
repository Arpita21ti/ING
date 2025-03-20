package quiz

import (
	"net/http"
	"strconv"

	"server/internal/model"
	"server/internal/service"
	"server/pkg/logger"

	"github.com/gin-gonic/gin"
)

// QuizHandler handles HTTP requests related to quizzes
type QuizHandler struct {
	quizService service.QuizService
	logger      logger.Logger
}

// NewQuizHandler creates a new QuizHandler instance
func NewQuizHandler(quizService service.QuizService, logger logger.Logger) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
		logger:      logger,
	}
}

// GetAllQuizzes retrieves all quizzes
func (h *QuizHandler) GetAllQuizzes(c *gin.Context) {
	quizzes, err := h.quizService.GetAllQuizzes(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get quizzes", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quizzes"})
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

// GetQuizByID retrieves a quiz by ID
func (h *QuizHandler) GetQuizByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	quiz, err := h.quizService.GetQuizByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get quiz", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quiz"})
		return
	}

	if quiz == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

// CreateQuiz creates a new quiz
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	var quiz model.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.quizService.CreateQuiz(c.Request.Context(), &quiz)
	if err != nil {
		h.logger.Error("Failed to create quiz", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz"})
		return
	}

	quiz.ID = id
	c.JSON(http.StatusCreated, quiz)
}

// UpdateQuiz updates an existing quiz
func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	var quiz model.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz.ID = id
	if err := h.quizService.UpdateQuiz(c.Request.Context(), &quiz); err != nil {
		h.logger.Error("Failed to update quiz", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quiz"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

// DeleteQuiz deletes a quiz
func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	if err := h.quizService.DeleteQuiz(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete quiz", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quiz"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiz deleted successfully"})
}

// GetQuizQuestions retrieves all questions for a quiz
func (h *QuizHandler) GetQuizQuestions(c *gin.Context) {
	quizID, err := strconv.ParseInt(c.Param("quizId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	questions, err := h.quizService.GetQuizQuestions(c.Request.Context(), quizID)
	if err != nil {
		h.logger.Error("Failed to get quiz questions", "quizID", quizID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// AddQuizQuestion adds a question to a quiz
func (h *QuizHandler) AddQuizQuestion(c *gin.Context) {
	quizID, err := strconv.ParseInt(c.Param("quizId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	var question model.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question.QuizID = quizID
	id, err := h.quizService.AddQuizQuestion(c.Request.Context(), &question)
	if err != nil {
		h.logger.Error("Failed to add question", "quizID", quizID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add question"})
		return
	}

	question.ID = id
	c.JSON(http.StatusCreated, question)
}

// UpdateQuizQuestion updates a question
func (h *QuizHandler) UpdateQuizQuestion(c *gin.Context) {
	quizID, err := strconv.ParseInt(c.Param("quizId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	questionID, err := strconv.ParseInt(c.Param("questionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var question model.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question.ID = questionID
	question.QuizID = quizID
	if err := h.quizService.UpdateQuizQuestion(c.Request.Context(), &question); err != nil {
		h.logger.Error("Failed to update question", "questionID", questionID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
		return
	}

	c.JSON(http.StatusOK, question)
}

// DeleteQuizQuestion deletes a question
func (h *QuizHandler) DeleteQuizQuestion(c *gin.Context) {
	questionID, err := strconv.ParseInt(c.Param("questionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	if err := h.quizService.DeleteQuizQuestion(c.Request.Context(), questionID); err != nil {
		h.logger.Error("Failed to delete question", "questionID", questionID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}