package helper

import (
	"prestasi-mahasiswa/service"

	"github.com/gin-gonic/gin"
)

type ReportHelper struct {
	ReportService *service.ReportService
}

func NewReportHelper(reportService *service.ReportService) *ReportHelper {
	return &ReportHelper{
		ReportService: reportService,
	}
}

// GetSystemStatistics 
func (h *ReportHelper) GetSystemStatistics(c *gin.Context) {
	stats, err := h.ReportService.GetSystemStatistics()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve statistics", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "System statistics retrieved successfully",
		"data":    stats,
	})
}

// GetStudentReport
func (h *ReportHelper) GetStudentReport(c *gin.Context) {
	studentID := c.Param("id")
	if studentID == "" {
		c.JSON(400, gin.H{"error": "Student ID is required"})
		return
	}

	report, err := h.ReportService.GetStudentReport(studentID)
	if err != nil {
		if err.Error() == "student not found" {
			c.JSON(404, gin.H{"error": "Student not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve report", "details": err.Error()})
		}
		return
	}

	c.JSON(200, gin.H{
		"message": "Student report retrieved successfully",
		"data":    report,
	})
}

// GetAchievementHistory retrieves status change history for an achievement
func (h *ReportHelper) GetAchievementHistory(c *gin.Context) {
	achievementID := c.Param("id")
	if achievementID == "" {
		c.JSON(400, gin.H{"error": "Achievement ID is required"})
		return
	}

	history, err := h.ReportService.GetAchievementHistory(achievementID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve history", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement history retrieved successfully",
		"data":    history,
		"total":   len(history),
	})
}
