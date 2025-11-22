package helper

import (
	"prestasi-mahasiswa/service"

	"github.com/gin-gonic/gin"
)

type AchievementHelper struct {
	AchievementService *service.AchievementService
	FileService        *service.FileService
}

func NewAchievementHelper(achievementSvc *service.AchievementService, fileSvc *service.FileService) *AchievementHelper {
	return &AchievementHelper{
		AchievementService: achievementSvc,
		FileService:        fileSvc,
	}
}

func (h *AchievementHelper) GetAchievements(c *gin.Context) {
	achievements, err := h.AchievementService.GetAllAchievements()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Get all achievements",
		"data":    achievements,
		"status":  "success",
	})
}

func (h *AchievementHelper) CreateAchievement(c *gin.Context) {
	var achievement service.Achievement
	if err := c.ShouldBindJSON(&achievement); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.AchievementService.CreateAchievement(achievement); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Achievement created successfully",
		"status":  "success",
	})
}

func (h *AchievementHelper) GetAchievement(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Get achievement by ID: " + id,
		"status":  "not_implemented",
	})
}

func (h *AchievementHelper) UpdateAchievement(c *gin.Context) {
	id := c.Param("id")
	var achievement service.Achievement
	if err := c.ShouldBindJSON(&achievement); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.AchievementService.UpdateAchievement(id, achievement); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement updated successfully",
		"status":  "success",
	})
}

func (h *AchievementHelper) DeleteAchievement(c *gin.Context) {
	id := c.Param("id")
	if err := h.AchievementService.DeleteAchievement(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement deleted successfully",
		"status":  "success",
	})
}

func (h *AchievementHelper) UploadFile(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Upload file for achievement ID: " + id,
		"status":  "not_implemented",
	})
}

func (h *AchievementHelper) GetFiles(c *gin.Context) {
	id := c.Param("id")
	files, err := h.FileService.GetFiles(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Get files for achievement ID: " + id,
		"data":    files,
		"status":  "success",
	})
}

func (h *AchievementHelper) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	fileId := c.Param("fileId")

	if err := h.FileService.DeleteFile(id, fileId); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "File deleted successfully",
		"status":  "success",
	})
}
