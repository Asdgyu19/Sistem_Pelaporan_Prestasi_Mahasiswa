package helper

import (
	"prestasi-mahasiswa/service"

	"github.com/gin-gonic/gin"
)

type StudentHelper struct {
	UserService        *service.UserService
	AchievementService *service.AchievementService
}

func NewStudentHelper(userService *service.UserService, achievementService *service.AchievementService) *StudentHelper {
	return &StudentHelper{
		UserService:        userService,
		AchievementService: achievementService,
	}
}

// GetStudents godoc
// @Tags Students
// @Summary List all students
// @Description Get list of all mahasiswa (students) with optional search filter
// @Accept json
// @Produce json
// @Param search query string false "Search by name or NIM"
// @Param page query string false "Page number" default(1)
// @Param limit query string false "Items per page" default(10)
// @Success 200 {object} object
// @Failure 500 {object} object
// @Router /students [get]
func (h *StudentHelper) GetStudents(c *gin.Context) {
	// Get filter parameters
	search := c.Query("search")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	students, err := h.UserService.GetAllUsers("mahasiswa", search, false)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve students", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Students retrieved successfully",
		"data":    students,
		"page":    page,
		"limit":   limit,
		"total":   len(students),
	})
}

// GetStudent godoc
// @Tags Students
// @Summary Get student profile
// @Description Get a specific student's public profile
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /students/{id} [get]
func (h *StudentHelper) GetStudent(c *gin.Context) {
	studentID := c.Param("id")
	if studentID == "" {
		c.JSON(400, gin.H{"error": "Student ID is required"})
		return
	}

	student, err := h.UserService.GetUserByID(studentID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Student not found", "details": err.Error()})
		return
	}

	// Verify it's a mahasiswa
	if student.Role != "mahasiswa" {
		c.JSON(404, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Student profile retrieved successfully",
		"data":    student,
	})
}

// GetStudentAchievements godoc
// @Tags Students
// @Summary Get student achievements
// @Description Get all achievements for a specific student
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /students/{id}/achievements [get]
func (h *StudentHelper) GetStudentAchievements(c *gin.Context) {
	studentID := c.Param("id")
	if studentID == "" {
		c.JSON(400, gin.H{"error": "Student ID is required"})
		return
	}

	// Verify student exists
	student, err := h.UserService.GetUserByID(studentID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Student not found", "details": err.Error()})
		return
	}

	if student.Role != "mahasiswa" {
		c.JSON(404, gin.H{"error": "Student not found"})
		return
	}

	// Get achievements
	achievements, err := h.AchievementService.GetAchievementsByMahasiswa(studentID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve achievements", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":      "Student achievements retrieved successfully",
		"student_id":   studentID,
		"achievements": achievements,
		"total":        len(achievements),
	})
}

// AssignAdvisor godoc
// @Tags Students
// @Summary Assign advisor to student
// @Description Assign a dosen_wali (advisor) to a student (student can only assign to themselves)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Param request body object true "Advisor ID"
// @Success 200 {object} object
// @Failure 403 {object} object
// @Router /students/{id}/advisor [put]
func (h *StudentHelper) AssignAdvisor(c *gin.Context) {
	studentID := c.Param("id")
	if studentID == "" {
		c.JSON(400, gin.H{"error": "Student ID is required"})
		return
	}

	// Get current user from context
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Students can only assign advisor to themselves
	if currentUserID.(string) != studentID {
		c.JSON(403, gin.H{"error": "Cannot assign advisor to another student"})
		return
	}

	var req struct {
		AdvisorID string `json:"advisor_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	err := h.UserService.AssignAdvisor(service.AssignAdvisorRequest{
		MahasiswaID: studentID,
		AdvisorID:   req.AdvisorID,
	})

	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to assign advisor", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":    "Advisor assigned successfully",
		"advisor_id": req.AdvisorID,
		"status":     "success",
	})
}
