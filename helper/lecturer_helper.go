package helper

import (
	"prestasi-mahasiswa/service"

	"github.com/gin-gonic/gin"
)

type LecturerHelper struct {
	UserService *service.UserService
}

func NewLecturerHelper(userService *service.UserService) *LecturerHelper {
	return &LecturerHelper{
		UserService: userService,
	}
}

// GetLecturers godoc
// @Tags Lecturers
// @Summary List all lecturers
// @Description Get list of all dosen_wali (lecturers/advisors)
// @Accept json
// @Produce json
// @Param search query string false "Search by name"
// @Param page query string false "Page number" default(1)
// @Param limit query string false "Items per page" default(10)
// @Success 200 {object} object
// @Failure 500 {object} object
// @Router /lecturers [get]
func (h *LecturerHelper) GetLecturers(c *gin.Context) {
	search := c.Query("search")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	lecturers, err := h.UserService.GetAllUsers("dosen_wali", search, false)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve lecturers", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Lecturers retrieved successfully",
		"data":    lecturers,
		"page":    page,
		"limit":   limit,
		"total":   len(lecturers),
	})
}

// GetLecturer godoc
// @Tags Lecturers
// @Summary Get lecturer profile
// @Description Get a specific lecturer's profile
// @Accept json
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /lecturers/{id} [get]
func (h *LecturerHelper) GetLecturer(c *gin.Context) {
	lecturerID := c.Param("id")
	if lecturerID == "" {
		c.JSON(400, gin.H{"error": "Lecturer ID is required"})
		return
	}

	lecturer, err := h.UserService.GetUserByID(lecturerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Lecturer not found", "details": err.Error()})
		return
	}

	// Verify it's a dosen_wali
	if lecturer.Role != "dosen_wali" {
		c.JSON(404, gin.H{"error": "Lecturer not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Lecturer profile retrieved successfully",
		"data":    lecturer,
	})
}

// GetLecturerAdvisees godoc
// @Tags Lecturers
// @Summary Get lecturer's advisees
// @Description Get list of students advised by a specific lecturer
// @Accept json
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} object
// @Failure 404 {object} object
// @Router /lecturers/{id}/advisees [get]
func (h *LecturerHelper) GetLecturerAdvisees(c *gin.Context) {
	lecturerID := c.Param("id")
	if lecturerID == "" {
		c.JSON(400, gin.H{"error": "Lecturer ID is required"})
		return
	}

	// Verify lecturer exists and is dosen_wali
	lecturer, err := h.UserService.GetUserByID(lecturerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Lecturer not found", "details": err.Error()})
		return
	}

	if lecturer.Role != "dosen_wali" {
		c.JSON(404, gin.H{"error": "Lecturer not found"})
		return
	}

	// Get students by advisor
	students, err := h.UserService.GetStudentsByAdvisor(lecturerID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve advisees", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Lecturer advisees retrieved successfully",
		"lecturer": gin.H{
			"id":   lecturerID,
			"name": lecturer.Name,
		},
		"advisees": students,
		"total":    len(students),
	})
}
