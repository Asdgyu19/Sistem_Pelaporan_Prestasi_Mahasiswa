package helper

import (
	"prestasi-mahasiswa/service"
	"time"

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
	// Check user role from middleware
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(401, gin.H{"error": "User role not found"})
		return
	}

	userID, userIDExists := c.Get("user_id")
	if !userIDExists {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return
	}

	var achievements []service.Achievement
	var err error

	// Mahasiswa only see their own achievements
	if userRole.(string) == "mahasiswa" {
		achievements, err = h.AchievementService.GetAchievementsByMahasiswa(userID.(string))
	} else {
		// Dosen/Admin see all achievements
		achievements, err = h.AchievementService.GetAllAchievements()
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Get achievements",
		"data":    achievements,
		"count":   len(achievements),
		"status":  "success",
	})
}

func (h *AchievementHelper) CreateAchievement(c *gin.Context) {
	// Get mahasiswa ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return
	}

	userRole, roleExists := c.Get("user_role")
	if !roleExists || userRole.(string) != "mahasiswa" {
		c.JSON(403, gin.H{"error": "Only mahasiswa can create achievements"})
		return
	}

	var req struct {
		Title           string    `json:"title" binding:"required"`
		Description     string    `json:"description" binding:"required"`
		Category        string    `json:"category" binding:"required"`
		AchievementDate time.Time `json:"achievement_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Create achievement object
	achievement := service.Achievement{
		MahasiswaID:     userID.(string),
		Title:           req.Title,
		Description:     req.Description,
		Category:        req.Category,
		AchievementDate: req.AchievementDate,
	}

	createdAchievement, err := h.AchievementService.CreateAchievement(achievement)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Achievement created successfully",
		"data":    createdAchievement,
		"status":  "success",
	})
}

func (h *AchievementHelper) GetAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Achievement ID is required"})
		return
	}

	achievement, err := h.AchievementService.GetAchievementByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check access rights
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	// Mahasiswa can only see their own achievements
	if userRole.(string) == "mahasiswa" && achievement.MahasiswaID != userID.(string) {
		c.JSON(403, gin.H{"error": "Access denied: You can only view your own achievements"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Get achievement by ID: " + id,
		"data":    achievement,
		"status":  "success",
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

// SubmitAchievement submits achievement for verification (mahasiswa only)
func (h *AchievementHelper) SubmitAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Achievement ID is required"})
		return
	}

	// Check if user is mahasiswa
	userRole, _ := c.Get("user_role")
	if userRole.(string) != "mahasiswa" {
		c.JSON(403, gin.H{"error": "Only mahasiswa can submit achievements"})
		return
	}

	err := h.AchievementService.SubmitAchievement(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement submitted for verification",
		"status":  "success",
	})
}

// VerifyAchievement verifies achievement (dosen/admin only)
func (h *AchievementHelper) VerifyAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Achievement ID is required"})
		return
	}

	// Check if user is dosen or admin
	userRole, _ := c.Get("user_role")
	if userRole.(string) == "mahasiswa" {
		c.JSON(403, gin.H{"error": "Only dosen and admin can verify achievements"})
		return
	}

	userID, _ := c.Get("user_id")
	err := h.AchievementService.VerifyAchievement(id, userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement verified successfully",
		"status":  "success",
	})
}

// RejectAchievement rejects achievement (dosen/admin only)
func (h *AchievementHelper) RejectAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Achievement ID is required"})
		return
	}

	// Check if user is dosen or admin
	userRole, _ := c.Get("user_role")
	if userRole.(string) == "mahasiswa" {
		c.JSON(403, gin.H{"error": "Only dosen and admin can reject achievements"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Rejection reason is required"})
		return
	}

	userID, _ := c.Get("user_id")
	err := h.AchievementService.RejectAchievement(id, userID.(string), req.Reason)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement rejected",
		"reason":  req.Reason,
		"status":  "success",
	})
}
