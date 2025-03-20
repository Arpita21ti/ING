package student

import (
	"net/http"

	"server/internal/model"
	"server/internal/service"
	"server/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ProfileHandler handles HTTP requests related to student profiles
type ProfileHandler struct {
	profileService service.ProfileService
	logger         logger.Logger
}

// NewProfileHandler creates a new ProfileHandler instance
func NewProfileHandler(profileService service.ProfileService, logger logger.Logger) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		logger:         logger,
	}
}

// GetProfile retrieves a student's profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), studentID.(int64))
	if err != nil {
		h.logger.Error("Failed to get profile", "studentID", studentID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile updates a student's profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var profile model.StudentProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.StudentID = studentID.(int64)
	if err := h.profileService.UpdateProfile(c.Request.Context(), &profile); err != nil {
		h.logger.Error("Failed to update profile", "studentID", studentID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ChangePassword changes a student's password
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var request struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.profileService.ChangePassword(
		c.Request.Context(),
		studentID.(int64),
		request.CurrentPassword,
		request.NewPassword,
	); err != nil {
		h.logger.Error("Failed to change password", "studentID", studentID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// UploadProfilePicture uploads a profile picture
func (h *ProfileHandler) UploadProfilePicture(c *gin.Context) {
	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	file, err := c.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5MB limit"})
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		h.logger.Error("Failed to open uploaded file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process uploaded file"})
		return
	}
	defer src.Close()

	// Upload the file
	picturePath, err := h.profileService.UploadProfilePicture(c.Request.Context(), studentID.(int64), src, file.Filename)
	if err != nil {
		h.logger.Error("Failed to upload profile picture", "studentID", studentID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile picture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile picture uploaded successfully",
		"path":    picturePath,
	})
}

// GetAcademicRecords retrieves a student's academic records
func (h *ProfileHandler) GetAcademicRecords(c *gin.Context) {
	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	records, err := h.profileService.GetAcademicRecords(c.Request.Context(), studentID.(int64))
	if err != nil {
		h.logger.Error("Failed to get academic records", "studentID", studentID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve academic records"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// GetQuizHistory retrieves a student's quiz history
func (h *ProfileHandler) GetQuizHistory(c *gin.Context) {
	studentID, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	history, err := h.profileService.GetQuizHistory(c.Request.Context(), studentID.(int64))
	if err != nil {
		h.logger.Error("Failed to get quiz history", "studentID", studentID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quiz history"})
		return
	}

	c.JSON(http.StatusOK, history)
}