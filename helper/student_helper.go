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

// GetStudents 
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

// GetStudent 
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

// GetStudent 
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

// AssignAdvisor
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
